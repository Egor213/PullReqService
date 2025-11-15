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
				u.EXPECT().GetActiveUsersTeam(ctx, authorTeam, gomock.AssignableToTypeOf([]string{})).Return([]string{"r1", "r2"}, nil)
				p.EXPECT().AssignReviewers(ctx, prID, gomock.AssignableToTypeOf([]string{})).Return([]string{"r1", "r2"}, nil)
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
