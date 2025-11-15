package service_test

import (
	"context"
	"errors"
	"testing"

	e "app/internal/entity"
	repomocks "app/internal/mocks/repomock"
	repoerrs "app/internal/repo/errors"
	"app/internal/service"
	sd "app/internal/service/dto"
	serverrs "app/internal/service/errors"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestUsersService_SetIsActive(t *testing.T) {
	var (
		ctx    = context.Background()
		uID    = "user1"
		active = true
	)

	user := e.User{
		UserID:   uID,
		Username: "us1",
		IsActive: &active,
	}

	type args struct {
		in sd.SetIsActiveInput
	}

	type mockBehavior func(u *repomocks.MockUsers, a args)

	testCases := []struct {
		name         string
		args         args
		mockBehavior mockBehavior
		wantErr      error
		want         e.User
	}{
		{
			name: "OK",
			args: args{
				in: sd.SetIsActiveInput{
					UserID:   uID,
					IsActive: &active,
				},
			},
			mockBehavior: func(u *repomocks.MockUsers, a args) {
				u.EXPECT().SetIsActive(ctx, uID, &active).Return(user, nil)
			},
			want: user,
		},
		{
			name: "User not found",
			args: args{
				in: sd.SetIsActiveInput{
					UserID:   uID,
					IsActive: &active,
				},
			},
			mockBehavior: func(u *repomocks.MockUsers, a args) {
				u.EXPECT().SetIsActive(ctx, uID, &active).Return(e.User{}, repoerrs.ErrNotFound)
			},
			wantErr: serverrs.ErrNotFoundUser,
		},
		{
			name: "Cannot set param",
			args: args{
				in: sd.SetIsActiveInput{
					UserID:   uID,
					IsActive: &active,
				},
			},
			mockBehavior: func(u *repomocks.MockUsers, a args) {
				u.EXPECT().SetIsActive(ctx, uID, &active).Return(e.User{}, errors.New("other error"))
			},
			wantErr: serverrs.ErrCannotSetParam,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			u := repomocks.NewMockUsers(ctrl)

			tc.mockBehavior(u, tc.args)

			s := service.NewUsersService(u)

			got, err := s.SetIsActive(ctx, tc.args.in)

			assert.ErrorIs(t, err, tc.wantErr)
			assert.Equal(t, tc.want, got)
		})
	}
}
