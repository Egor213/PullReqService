package service_test

import (
	"app/internal/service"
	"context"
	"errors"
	"testing"

	e "app/internal/entity"
	repomocks "app/internal/mocks/repomock"
	repoerrs "app/internal/repo/errors"

	sd "app/internal/service/dto"
	serverrs "app/internal/service/errors"

	"github.com/avito-tech/go-transaction-manager/trm/v2/drivers/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestTeamsService_CreateOrUpdateTeam(t *testing.T) {
	var (
		ctx      = context.Background()
		teamName = "teamA"
	)

	memberDTOs := []sd.TeamMemberDTO{
		{UserID: "u1", Username: "user1"},
		{UserID: "u2", Username: "user2"},
	}
	members := []e.TeamMember{
		{UserID: "u1", Username: "user1"},
		{UserID: "u2", Username: "user2"},
	}

	type args struct {
		in sd.CrOrUpTeamInput
	}

	type mockBehavior func(t *repomocks.MockTeams, u *repomocks.MockUsers, args args)

	testCases := []struct {
		name         string
		args         args
		mockBehavior mockBehavior
		wantErr      error
		want         e.Team
	}{
		{
			name: "Create new team successfully",
			args: args{
				in: sd.CrOrUpTeamInput{
					TeamName: teamName,
					Members:  memberDTOs,
				},
			},
			mockBehavior: func(t *repomocks.MockTeams, u *repomocks.MockUsers, args args) {
				t.EXPECT().GetTeam(ctx, teamName).Return(e.Team{}, repoerrs.ErrNotFound)
				t.EXPECT().CreateTeam(ctx, teamName).Return(e.Team{
					TeamName: teamName,
					Members:  members,
				}, nil)
				u.EXPECT().UpsetBulk(ctx, gomock.AssignableToTypeOf([]e.User{})).Return(nil)
			},
			want: e.Team{
				TeamName: teamName,
				Members:  members,
			},
		},
		{
			name: "Team update",
			args: args{
				in: sd.CrOrUpTeamInput{
					TeamName: teamName,
					Members: []sd.TeamMemberDTO{
						{UserID: "u3", Username: "user1"},
						{UserID: "u2", Username: "user2"},
					},
				},
			},
			mockBehavior: func(t *repomocks.MockTeams, u *repomocks.MockUsers, args args) {
				t.EXPECT().GetTeam(ctx, teamName).Return(e.Team{
					TeamName: teamName,
					Members:  members,
				}, nil)
				t.EXPECT().DeleteUsersFromTeam(ctx, teamName).Return(nil)
				u.EXPECT().UpsetBulk(ctx, gomock.AssignableToTypeOf([]e.User{})).Return(nil)
			},
			want: e.Team{
				TeamName: teamName,
				Members: []e.TeamMember{
					{UserID: "u3", Username: "user1"},
					{UserID: "u2", Username: "user2"},
				},
			},
		},
		{
			name: "Cannot get team",
			args: args{
				in: sd.CrOrUpTeamInput{
					TeamName: teamName,
					Members:  memberDTOs,
				},
			},
			mockBehavior: func(t *repomocks.MockTeams, u *repomocks.MockUsers, args args) {
				t.EXPECT().GetTeam(ctx, teamName).Return(e.Team{}, errors.New("Other error"))
			},
			wantErr: serverrs.ErrCannotGetTeam,
		},
		{
			name: "Team with users already exists",
			args: args{
				in: sd.CrOrUpTeamInput{
					TeamName: teamName,
					Members:  memberDTOs,
				},
			},
			mockBehavior: func(t *repomocks.MockTeams, u *repomocks.MockUsers, args args) {
				t.EXPECT().GetTeam(ctx, teamName).Return(e.Team{
					TeamName: teamName,
					Members:  members,
				}, nil)
			},
			wantErr: serverrs.ErrTeamWithUsersExists,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockTeams := repomocks.NewMockTeams(ctrl)
			mockUsers := repomocks.NewMockUsers(ctrl)

			tc.mockBehavior(mockTeams, mockUsers, tc.args)

			m := mock.NewMockManager(ctrl)
			m.EXPECT().Do(gomock.Any(), gomock.Any()).DoAndReturn(
				func(ctx context.Context, fn func(ctx context.Context) error) error {
					return fn(ctx)
				},
			).AnyTimes()

			s := service.NewTeamsService(mockTeams, mockUsers, m)

			team, err := s.CreateOrUpdateTeam(ctx, tc.args.in)

			assert.ErrorIs(t, err, tc.wantErr)
			assert.Equal(t, tc.want, team)
		})
	}
}

func TestTeamsService_ReplaceTeamMembers(t *testing.T) {
	var (
		ctx      = context.Background()
		teamName = "teamA"
	)

	memberDTOs := []sd.TeamMemberDTO{
		{UserID: "u1", Username: "user1"},
		{UserID: "u2", Username: "user2"},
	}

	type args struct {
		in sd.ReplaceMembersInput
	}

	type mockBehavior func(t *repomocks.MockTeams, u *repomocks.MockUsers, args args)

	testCases := []struct {
		name         string
		args         args
		mockBehavior mockBehavior
		wantErr      error
		want         e.Team
	}{
		{
			name: "OK",
			args: args{
				in: sd.ReplaceMembersInput{
					TeamName: teamName,
					Members:  memberDTOs,
				},
			},
			mockBehavior: func(t *repomocks.MockTeams, u *repomocks.MockUsers, args args) {
				t.EXPECT().DeleteUsersFromTeam(ctx, teamName).Return(nil)
				u.EXPECT().UpsetBulk(ctx, gomock.AssignableToTypeOf([]e.User{})).Return(nil)
			},
		},
		{
			name: "Cannot del users from the team",
			args: args{
				in: sd.ReplaceMembersInput{
					TeamName: teamName,
					Members:  memberDTOs,
				},
			},
			mockBehavior: func(t *repomocks.MockTeams, u *repomocks.MockUsers, args args) {
				t.EXPECT().DeleteUsersFromTeam(ctx, teamName).Return(errors.New("other error"))
			},
			wantErr: serverrs.ErrCannotDelUsersFromTeam,
		},
		{
			name: "Cannot update or set users",
			args: args{
				in: sd.ReplaceMembersInput{
					TeamName: teamName,
					Members:  memberDTOs,
				},
			},
			mockBehavior: func(t *repomocks.MockTeams, u *repomocks.MockUsers, args args) {
				t.EXPECT().DeleteUsersFromTeam(ctx, teamName).Return(nil)
				u.EXPECT().UpsetBulk(ctx, gomock.AssignableToTypeOf([]e.User{})).Return(errors.New("other error"))
			},
			wantErr: serverrs.ErrCannotUpsetUsers,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockTeams := repomocks.NewMockTeams(ctrl)
			mockUsers := repomocks.NewMockUsers(ctrl)

			tc.mockBehavior(mockTeams, mockUsers, tc.args)

			m := mock.NewMockManager(ctrl)
			m.EXPECT().Do(gomock.Any(), gomock.Any()).DoAndReturn(
				func(ctx context.Context, fn func(ctx context.Context) error) error {
					return fn(ctx)
				},
			).AnyTimes()

			s := service.NewTeamsService(mockTeams, mockUsers, m)

			err := s.ReplaceTeamMembers(ctx, tc.args.in)

			assert.ErrorIs(t, err, tc.wantErr)
		})
	}
}

func TestTeamsService_GetTeam(t *testing.T) {
	var (
		ctx      = context.Background()
		teamName = "teamA"
	)

	members := []e.TeamMember{
		{UserID: "u1", Username: "user1"},
		{UserID: "u2", Username: "user2"},
	}

	type args struct {
		teamName string
	}

	type mockBehavior func(t *repomocks.MockTeams, args args)

	testCases := []struct {
		name         string
		args         args
		mockBehavior mockBehavior
		wantErr      error
		want         e.Team
	}{
		{
			name: "OK",
			args: args{teamName: teamName},
			mockBehavior: func(t *repomocks.MockTeams, args args) {
				t.EXPECT().GetTeam(ctx, args.teamName).Return(e.Team{
					TeamName: args.teamName,
					Members:  members,
				}, nil)
			},
			want: e.Team{
				TeamName: teamName,
				Members:  members,
			},
		},
		{
			name: "Team not found",
			args: args{teamName: teamName},
			mockBehavior: func(t *repomocks.MockTeams, args args) {
				t.EXPECT().GetTeam(ctx, args.teamName).Return(e.Team{}, repoerrs.ErrNotFound)
			},
			wantErr: serverrs.ErrNotFoundTeam,
		},
		{
			name: "Other repo error",
			args: args{teamName: teamName},
			mockBehavior: func(t *repomocks.MockTeams, args args) {
				t.EXPECT().GetTeam(ctx, args.teamName).Return(e.Team{}, errors.New("other error"))
			},
			wantErr: serverrs.ErrCannotGetTeam,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			tm := repomocks.NewMockTeams(ctrl)

			tc.mockBehavior(tm, tc.args)

			s := service.NewTeamsService(tm, nil, nil)
			got, err := s.GetTeam(ctx, tc.args.teamName)

			assert.ErrorIs(t, err, tc.wantErr)
			assert.Equal(t, tc.want, got)
		})
	}
}
