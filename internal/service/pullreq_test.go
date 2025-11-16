package service_test

import (
	e "app/internal/entity"
	repomocks "app/internal/mocks/repomock"
	rd "app/internal/repo/dto"
	repoerrs "app/internal/repo/errors"
	"app/internal/service"
	sd "app/internal/service/dto"
	serverrs "app/internal/service/errors"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/avito-tech/go-transaction-manager/trm/v2/drivers/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestPullReqService_CreatePR(t *testing.T) {
	var (
		ctx        = context.Background()
		prID       = "pr1"
		authorID   = "user1"
		authorTeam = "teamA"
	)
	isActiveTrue := true
	isActiveFalse := false
	user := e.User{
		IsActive: &isActiveTrue,
		UserID:   authorID,
		Username: "us1",
		TeamName: authorTeam,
	}

	type args struct {
		in rd.CreatePRInput
	}

	type MockBehavior func(u *repomocks.MockUsers, p *repomocks.MockPullReq, args args)

	defPR := rd.CreatePRInput{
		PullReqID: prID,
		NamePR:    prID,
		AuthorID:  authorID,
		Status:    e.StatusOpen,
	}

	testCases := []struct {
		name         string
		args         args
		mockBehavior MockBehavior
		wantErr      error
		want         e.PullRequest
	}{
		{
			name: "OK",
			args: args{
				in: defPR,
			},
			mockBehavior: func(u *repomocks.MockUsers, p *repomocks.MockPullReq, args args) {
				u.EXPECT().GetUserByID(ctx, authorID).Return(user, nil)
				p.EXPECT().CreatePR(ctx, args.in).Return(e.PullRequest{
					PullReqID: prID,
					NamePR:    prID,
					AuthorID:  authorID,
					Status:    e.StatusOpen,
				}, nil)
				u.EXPECT().GetActiveUsersTeam(ctx, authorTeam, gomock.AssignableToTypeOf([]string{})).
					Return([]string{"r1", "r2"}, nil)
				p.EXPECT().AssignReviewers(ctx, prID, gomock.AssignableToTypeOf([]string{})).
					Return([]string{"r1", "r2"}, nil)
			},
			want: e.PullRequest{
				PullReqID: prID,
				NamePR:    prID,
				AuthorID:  authorID,
				Status:    e.StatusOpen,
				Reviewers: []string{"r1", "r2"},
			},
		},
		{
			name: "PR already exists",
			args: args{
				in: defPR,
			},
			mockBehavior: func(u *repomocks.MockUsers, p *repomocks.MockPullReq, args args) {
				u.EXPECT().GetUserByID(ctx, authorID).Return(user, nil)
				p.EXPECT().CreatePR(ctx, args.in).Return(e.PullRequest{}, repoerrs.ErrAlreadyExists)
			},
			wantErr: serverrs.ErrPRExists,
		},
		{
			name: "PR cannot create",
			args: args{
				in: defPR,
			},
			mockBehavior: func(u *repomocks.MockUsers, p *repomocks.MockPullReq, args args) {
				u.EXPECT().GetUserByID(ctx, authorID).Return(user, nil)
				p.EXPECT().CreatePR(ctx, args.in).Return(e.PullRequest{}, errors.New("other error"))
			},
			wantErr: serverrs.ErrCannotCreatePR,
		},
		{
			name: "Author pr not found",
			args: args{
				in: defPR,
			},
			mockBehavior: func(u *repomocks.MockUsers, p *repomocks.MockPullReq, args args) {
				u.EXPECT().GetUserByID(ctx, authorID).Return(e.User{}, repoerrs.ErrNotFound)
			},
			wantErr: serverrs.ErrNotFoundUser,
		},
		{
			name: "Author pr not found",
			args: args{
				in: defPR,
			},
			mockBehavior: func(u *repomocks.MockUsers, p *repomocks.MockPullReq, args args) {
				u.EXPECT().GetUserByID(ctx, authorID).Return(e.User{}, repoerrs.ErrNotFound)
			},
			wantErr: serverrs.ErrNotFoundUser,
		},
		{
			name: "Author pr not active",
			args: args{
				in: defPR,
			},
			mockBehavior: func(u *repomocks.MockUsers, p *repomocks.MockPullReq, args args) {
				u.EXPECT().GetUserByID(ctx, authorID).Return(e.User{
					IsActive: &isActiveFalse,
					UserID:   user.UserID,
					Username: user.Username,
					TeamName: user.TeamName,
				}, nil)
			},
			wantErr: serverrs.ErrInactiveCreator,
		},
		{
			name: "Author pr not have the team",
			args: args{
				in: defPR,
			},
			mockBehavior: func(u *repomocks.MockUsers, p *repomocks.MockPullReq, args args) {
				u.EXPECT().GetUserByID(ctx, authorID).Return(e.User{
					IsActive: user.IsActive,
					UserID:   user.UserID,
					Username: user.Username,
					TeamName: "",
				}, nil)
			},
			wantErr: serverrs.ErrNotFoundTeam,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUsers := repomocks.NewMockUsers(ctrl)
			mockPR := repomocks.NewMockPullReq(ctrl)

			tc.mockBehavior(mockUsers, mockPR, tc.args)

			m := mock.NewMockManager(ctrl)
			m.EXPECT().Do(gomock.Any(), gomock.Any()).DoAndReturn(
				func(ctx context.Context, fn func(ctx context.Context) error) error {
					return fn(ctx)
				},
			).AnyTimes()

			s := service.NewPullReqService(mockPR, mockUsers, m)

			inServ := sd.CreatePRInput{
				PullReqID: tc.args.in.PullReqID,
				NamePR:    tc.args.in.NamePR,
				AuthorID:  tc.args.in.AuthorID,
			}
			pr, err := s.CreatePR(ctx, inServ)

			assert.ErrorIs(t, err, tc.wantErr)
			assert.Equal(t, pr, tc.want)
		})
	}
}

func TestPullReqService_GetPR(t *testing.T) {
	var (
		ctx  = context.Background()
		prID = "pr1"
		pr   = e.PullRequest{
			PullReqID: prID,
			NamePR:    "TestPR",
			AuthorID:  "user1",
			Status:    e.StatusOpen,
		}
	)

	type args struct {
		prID string
	}

	type mockBehavior func(p *repomocks.MockPullReq, args args)

	testCases := []struct {
		name         string
		args         args
		mockBehavior mockBehavior
		wantErr      error
		want         e.PullRequest
	}{
		{
			name: "OK",
			args: args{
				prID: prID,
			},
			mockBehavior: func(p *repomocks.MockPullReq, args args) {
				p.EXPECT().GetPR(ctx, args.prID).Return(pr, nil)
			},
			want: pr,
		},
		{
			name: "PR not found",
			args: args{
				prID: prID,
			},
			mockBehavior: func(p *repomocks.MockPullReq, args args) {
				p.EXPECT().GetPR(ctx, args.prID).Return(e.PullRequest{}, repoerrs.ErrNotFound)
			},
			wantErr: serverrs.ErrNotFoundPR,
		},
		{
			name: "Other repo error",
			args: args{
				prID: prID,
			},
			mockBehavior: func(p *repomocks.MockPullReq, args args) {
				p.EXPECT().GetPR(ctx, args.prID).Return(e.PullRequest{}, errors.New("other error"))
			},
			wantErr: serverrs.ErrCannotGetPR,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockPR := repomocks.NewMockPullReq(ctrl)
			tc.mockBehavior(mockPR, tc.args)

			s := service.NewPullReqService(mockPR, nil, nil)

			got, err := s.GetPR(ctx, tc.args.prID)
			assert.ErrorIs(t, err, tc.wantErr)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestPullReqService_GetPRsByReviewer(t *testing.T) {
	var (
		ctx    = context.Background()
		userID = "user1"
		user   = e.User{
			UserID:   userID,
			Username: "us1",
			IsActive: nil,
			TeamName: "teamA",
		}
		prs = []e.PullRequestShort{
			{PullReqID: "pr1", NamePR: "TestPR1", Status: e.StatusOpen},
			{PullReqID: "pr2", NamePR: "TestPR2", Status: e.StatusMerged},
		}
	)

	type args struct {
		uID string
	}

	type mockBehavior func(u *repomocks.MockUsers, p *repomocks.MockPullReq, args args)

	testCases := []struct {
		name         string
		args         args
		mockBehavior mockBehavior
		wantErr      error
		want         []e.PullRequestShort
	}{
		{
			name: "OK",
			args: args{
				uID: userID,
			},
			mockBehavior: func(u *repomocks.MockUsers, p *repomocks.MockPullReq, args args) {
				u.EXPECT().GetUserByID(ctx, args.uID).Return(user, nil)
				p.EXPECT().GetPRsByReviewer(ctx, args.uID).Return(prs, nil)
			},
			want: prs,
		},
		{
			name: "User not found",
			args: args{
				uID: userID,
			},
			mockBehavior: func(u *repomocks.MockUsers, p *repomocks.MockPullReq, args args) {
				u.EXPECT().GetUserByID(ctx, args.uID).Return(e.User{}, repoerrs.ErrNotFound)
			},
			wantErr: serverrs.ErrNotFoundUser,
		},
		{
			name: "Cannot get PRs by reviewer",
			args: args{
				uID: userID,
			},
			mockBehavior: func(u *repomocks.MockUsers, p *repomocks.MockPullReq, args args) {
				u.EXPECT().GetUserByID(ctx, args.uID).Return(user, nil)
				p.EXPECT().GetPRsByReviewer(ctx, args.uID).Return(nil, errors.New("other error"))
			},
			wantErr: serverrs.ErrCannotGetPR,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUsers := repomocks.NewMockUsers(ctrl)
			mockPR := repomocks.NewMockPullReq(ctrl)

			tc.mockBehavior(mockUsers, mockPR, tc.args)

			s := service.NewPullReqService(mockPR, mockUsers, nil)

			got, err := s.GetPRsByReviewer(ctx, tc.args.uID)
			assert.ErrorIs(t, err, tc.wantErr)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestPullReqService_MergePR(t *testing.T) {
	ctx := context.Background()
	prID := "pr1"
	mergedAt := time.Date(2025, 11, 15, 23, 0, 0, 0, time.UTC)
	prOpen := e.PullRequest{
		PullReqID: prID,
		NamePR:    "TestPR",
		AuthorID:  "user1",
		Status:    e.StatusOpen,
	}
	prMerged := e.PullRequest{
		PullReqID: prID,
		NamePR:    "TestPR",
		AuthorID:  "user1",
		Status:    e.StatusMerged,
	}

	type args struct {
		prID string
	}

	type mockBehavior func(p *repomocks.MockPullReq, args args)

	testCases := []struct {
		name         string
		args         args
		mockBehavior mockBehavior
		wantErr      error
		want         e.PullRequest
	}{
		{
			name: "Merge OK",
			args: args{
				prID: prID,
			},
			mockBehavior: func(p *repomocks.MockPullReq, args args) {
				p.EXPECT().GetPR(ctx, args.prID).Return(prOpen, nil)
				p.EXPECT().MergePR(ctx, args.prID).Return(&mergedAt, nil)
			},
			want: e.PullRequest{
				PullReqID: prID,
				NamePR:    "TestPR",
				AuthorID:  "user1",
				Status:    e.StatusMerged,
				MergedAt:  &mergedAt,
			},
		},
		{
			name: "PR already merged",
			args: args{
				prID: prID,
			},
			mockBehavior: func(p *repomocks.MockPullReq, args args) {
				p.EXPECT().GetPR(ctx, args.prID).Return(prMerged, nil)
			},
			want: prMerged,
		},
		{
			name: "PR not found",
			args: args{
				prID: prID,
			},
			mockBehavior: func(p *repomocks.MockPullReq, args args) {
				p.EXPECT().GetPR(ctx, args.prID).Return(e.PullRequest{}, repoerrs.ErrNotFound)
			},
			wantErr: serverrs.ErrNotFoundPR,
		},
		{
			name: "MergePR not found after attempt",
			args: args{
				prID: prID,
			},
			mockBehavior: func(p *repomocks.MockPullReq, args args) {
				p.EXPECT().GetPR(ctx, args.prID).Return(prOpen, nil).AnyTimes()
				p.EXPECT().MergePR(ctx, args.prID).Return(&time.Time{}, repoerrs.ErrNotFound)
			},
			want: prOpen,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockPR := repomocks.NewMockPullReq(ctrl)
			tc.mockBehavior(mockPR, tc.args)

			s := service.NewPullReqService(mockPR, nil, nil)
			got, err := s.MergePR(ctx, tc.args.prID)

			assert.ErrorIs(t, err, tc.wantErr)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestPullReqService_ReassignReviewer(t *testing.T) {
	ctx := context.Background()
	prID := "pr1"
	oldRev := "r1"
	authorID := "user1"
	teamName := "teamA"

	type args struct {
		in sd.ReassignReviewerInput
	}

	type mockBehavior func(p *repomocks.MockPullReq, u *repomocks.MockUsers, args args)

	prOpen := e.PullRequest{
		PullReqID: prID,
		AuthorID:  authorID,
		Reviewers: []string{oldRev, "r2"},
		Status:    e.StatusOpen,
	}

	user := e.User{
		UserID:   authorID,
		TeamName: teamName,
	}

	testCases := []struct {
		name         string
		args         args
		mockBehavior mockBehavior
		wantErr      error
		want         sd.ReassignReviewerOutput
	}{
		{
			name: "PR not found",
			args: args{
				in: sd.ReassignReviewerInput{
					PullReqID: prID,
					RevID:     oldRev,
				},
			},
			mockBehavior: func(p *repomocks.MockPullReq, u *repomocks.MockUsers, args args) {
				p.EXPECT().GetPR(ctx, prID).Return(e.PullRequest{}, repoerrs.ErrNotFound)
			},
			wantErr: serverrs.ErrNotFoundPR,
		},
		{
			name: "No reviewers assigned",
			args: args{
				in: sd.ReassignReviewerInput{
					PullReqID: prID,
					RevID:     oldRev,
				},
			},
			mockBehavior: func(p *repomocks.MockPullReq, u *repomocks.MockUsers, args args) {
				p.EXPECT().GetPR(ctx, prID).Return(e.PullRequest{
					PullReqID: prID,
					AuthorID:  authorID,
					Reviewers: []string{},
					Status:    e.StatusOpen,
				}, nil)
			},
			wantErr: serverrs.ErrNotFoundReviewers,
		},
		{
			name: "Reviewer not assigned",
			args: args{
				in: sd.ReassignReviewerInput{
					PullReqID: prID,
					RevID:     "rX",
				},
			},
			mockBehavior: func(p *repomocks.MockPullReq, u *repomocks.MockUsers, args args) {
				p.EXPECT().GetPR(ctx, prID).Return(prOpen, nil)
			},
			wantErr: serverrs.ErrReviewerNotAssigned,
		},
		{
			name: "PR already merged",
			args: args{
				in: sd.ReassignReviewerInput{
					PullReqID: prID,
					RevID:     oldRev,
				},
			},
			mockBehavior: func(p *repomocks.MockPullReq, u *repomocks.MockUsers, args args) {
				p.EXPECT().GetPR(ctx, prID).Return(e.PullRequest{
					PullReqID: prID,
					AuthorID:  authorID,
					Reviewers: []string{oldRev},
					Status:    e.StatusMerged,
				}, nil)
			},
			wantErr: serverrs.ErrMergedPR,
		},
		{
			name: "No available reviewers",
			args: args{
				in: sd.ReassignReviewerInput{
					PullReqID: prID,
					RevID:     oldRev,
				},
			},
			mockBehavior: func(p *repomocks.MockPullReq, u *repomocks.MockUsers, args args) {
				p.EXPECT().GetPR(ctx, prID).Return(prOpen, nil)
				u.EXPECT().GetUserByID(ctx, authorID).Return(user, nil)
				u.EXPECT().GetActiveUsersTeam(ctx, teamName, gomock.AssignableToTypeOf([]string{})).
					Return([]string{}, nil)
			},
			wantErr: serverrs.ErrNoAvailableReviewers,
		},
		{
			name: "Cannot change reviewer",
			args: args{
				in: sd.ReassignReviewerInput{
					PullReqID: prID,
					RevID:     oldRev,
				},
			},
			mockBehavior: func(p *repomocks.MockPullReq, u *repomocks.MockUsers, args args) {
				p.EXPECT().GetPR(ctx, prID).Return(prOpen, nil)
				u.EXPECT().GetUserByID(ctx, authorID).Return(user, nil)
				u.EXPECT().GetActiveUsersTeam(ctx, teamName, gomock.AssignableToTypeOf([]string{})).
					Return([]string{"r3"}, nil)
				p.EXPECT().ChangeReviewer(ctx, gomock.Any()).Return(errors.New("other error"))
			},
			wantErr: serverrs.ErrCannotChangeReviewer,
		},
		{
			name: "OK",
			args: args{
				in: sd.ReassignReviewerInput{
					PullReqID: prID,
					RevID:     oldRev,
				},
			},
			mockBehavior: func(p *repomocks.MockPullReq, u *repomocks.MockUsers, args args) {
				p.EXPECT().GetPR(ctx, prID).Return(prOpen, nil)
				u.EXPECT().GetUserByID(ctx, authorID).Return(user, nil)
				u.EXPECT().GetActiveUsersTeam(ctx, teamName, gomock.AssignableToTypeOf([]string{})).
					Return([]string{"r3"}, nil)
				p.EXPECT().ChangeReviewer(ctx, gomock.Any()).Return(nil)
			},
			want: sd.ReassignReviewerOutput{
				NewRevID: "r3",
				PullReq: e.PullRequest{
					PullReqID: prID,
					AuthorID:  authorID,
					Reviewers: []string{"r3", "r2"},
					Status:    e.StatusOpen,
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			p := repomocks.NewMockPullReq(ctrl)
			u := repomocks.NewMockUsers(ctrl)

			tc.mockBehavior(p, u, tc.args)

			m := mock.NewMockManager(ctrl)
			m.EXPECT().DoWithSettings(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
				func(ctx context.Context, _ any, fn func(ctx context.Context) error) error {
					return fn(ctx)
				},
			).AnyTimes()

			svc := service.NewPullReqService(p, u, m)

			got, err := svc.ReassignReviewer(ctx, tc.args.in)

			assert.ErrorIs(t, err, tc.wantErr)
			assert.Equal(t, tc.want.PullReq.PullReqID, got.PullReq.PullReqID)
			assert.Equal(t, tc.want.PullReq.Reviewers, got.PullReq.Reviewers)
			assert.Equal(t, tc.want.NewRevID, got.NewRevID)
		})
	}
}

func TestPullReqService_AssignReviewers(t *testing.T) {
	ctx := context.Background()
	prID := "pr1"
	authorTeam := "teamA"

	type args struct {
		prID       string
		authorTeam string
		exclude    []string
	}

	type mockBehavior func(p *repomocks.MockPullReq, u *repomocks.MockUsers, args args)

	testCases := []struct {
		name         string
		args         args
		mockBehavior mockBehavior
		wantErr      error
		want         []string
	}{
		{
			name: "OK - enough reviewers",
			args: args{
				prID:       prID,
				authorTeam: authorTeam,
				exclude:    []string{"u1"},
			},
			mockBehavior: func(p *repomocks.MockPullReq, u *repomocks.MockUsers, args args) {
				u.EXPECT().GetActiveUsersTeam(ctx, args.authorTeam, args.exclude).
					Return([]string{"r2", "r3", "r4"}, nil)
				p.EXPECT().AssignReviewers(ctx, args.prID, gomock.AssignableToTypeOf([]string{})).
					Return([]string{"r2", "r3"}, nil)
			},
			want: []string{"r2", "r3"},
		},
		{
			name: "Not enough reviewers - set need more",
			args: args{
				prID:       prID,
				authorTeam: authorTeam,
				exclude:    []string{"u1"},
			},
			mockBehavior: func(p *repomocks.MockPullReq, u *repomocks.MockUsers, args args) {
				u.EXPECT().GetActiveUsersTeam(ctx, args.authorTeam, args.exclude).
					Return([]string{"r2"}, nil)
				p.EXPECT().AssignReviewers(ctx, args.prID, gomock.AssignableToTypeOf([]string{})).
					Return([]string{"r2"}, nil)
				p.EXPECT().SetNeedMoreReviewers(ctx, args.prID, true).Return(nil)
			},
			want: []string{"r2"},
		},
		{
			name: "No reviewers - set need more",
			args: args{
				prID:       prID,
				authorTeam: authorTeam,
				exclude:    []string{"u1"},
			},
			mockBehavior: func(p *repomocks.MockPullReq, u *repomocks.MockUsers, args args) {
				u.EXPECT().GetActiveUsersTeam(ctx, args.authorTeam, args.exclude).
					Return([]string{}, nil)
				p.EXPECT().SetNeedMoreReviewers(ctx, args.prID, true).Return(nil)
			},
			want: nil,
		},
		{
			name: "Cannot get users",
			args: args{
				prID:       prID,
				authorTeam: authorTeam,
				exclude:    []string{"u1"},
			},
			mockBehavior: func(p *repomocks.MockPullReq, u *repomocks.MockUsers, args args) {
				u.EXPECT().GetActiveUsersTeam(ctx, args.authorTeam, args.exclude).
					Return(nil, errors.New("other error"))
			},
			wantErr: serverrs.ErrCannotGetUser,
		},
		{
			name: "Cannot assign reviewers",
			args: args{
				prID:       prID,
				authorTeam: authorTeam,
				exclude:    []string{"u1"},
			},
			mockBehavior: func(p *repomocks.MockPullReq, u *repomocks.MockUsers, args args) {
				u.EXPECT().GetActiveUsersTeam(ctx, args.authorTeam, args.exclude).
					Return([]string{"r2", "r3"}, nil)
				p.EXPECT().AssignReviewers(ctx, args.prID, gomock.AssignableToTypeOf([]string{})).
					Return(nil, errors.New("other error"))
			},
			wantErr: serverrs.ErrCannotAssignReviewers,
		},
		{
			name: "Cannot change set need more reviewers",
			args: args{
				prID:       prID,
				authorTeam: authorTeam,
				exclude:    []string{"u1"},
			},
			mockBehavior: func(p *repomocks.MockPullReq, u *repomocks.MockUsers, args args) {
				u.EXPECT().GetActiveUsersTeam(ctx, args.authorTeam, args.exclude).
					Return([]string{"r2"}, nil)
				p.EXPECT().AssignReviewers(ctx, args.prID, gomock.AssignableToTypeOf([]string{})).
					Return([]string{"r2"}, nil)
				p.EXPECT().SetNeedMoreReviewers(ctx, args.prID, true).Return(errors.New("other error"))
			},
			wantErr: serverrs.ErrCannotChangeSetNeedMoRev,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			p := repomocks.NewMockPullReq(ctrl)
			u := repomocks.NewMockUsers(ctrl)

			tc.mockBehavior(p, u, tc.args)

			m := mock.NewMockManager(ctrl)
			m.EXPECT().Do(gomock.Any(), gomock.Any()).DoAndReturn(
				func(ctx context.Context, fn func(ctx context.Context) error) error {
					return fn(ctx)
				},
			).AnyTimes()

			svc := service.NewPullReqService(p, u, m)

			got, err := svc.AssignReviewers(ctx, sd.AssignReviewersInput{
				PullReqID:    tc.args.prID,
				AuthorTeam:   tc.args.authorTeam,
				ExcludeUsers: tc.args.exclude,
			})

			assert.ErrorIs(t, err, tc.wantErr)
			assert.Equal(t, tc.want, got)
		})
	}
}
