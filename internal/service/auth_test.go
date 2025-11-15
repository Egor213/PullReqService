package service_test

import (
	"app/internal/entity"
	repomocks "app/internal/mocks/repomock"
	repoerrs "app/internal/repo/errors"
	"app/internal/service"
	sd "app/internal/service/dto"
	se "app/internal/service/errors"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestAuthService_GenerateToken(t *testing.T) {
	const (
		secretKey = "secret"
		tokenTTL  = time.Minute
	)

	type args struct {
		ctx   context.Context
		input sd.GenTokenInput
	}

	type MockBehavior func(o *repomocks.MockUsers, args args)

	isActive := true
	user := entity.User{
		IsActive: &isActive,
		UserID:   "u1",
		Username: "test",
		TeamName: "test",
	}

	testCases := []struct {
		name         string
		args         args
		mockBehavior MockBehavior
		wantErr      error
	}{
		{
			name: "OK",
			args: args{
				ctx: context.Background(),
				input: sd.GenTokenInput{
					UserID: user.UserID,
					Role:   entity.RoleUser,
				},
			},
			mockBehavior: func(o *repomocks.MockUsers, args args) {
				o.EXPECT().GetUserByID(args.ctx, args.input.UserID).Return(user, nil)
			},
			wantErr: nil,
		},
		{
			name: "User not found",
			args: args{
				ctx: context.Background(),
				input: sd.GenTokenInput{
					UserID: "unknown",
					Role:   entity.RoleUser,
				},
			},
			mockBehavior: func(o *repomocks.MockUsers, args args) {
				o.EXPECT().GetUserByID(args.ctx, args.input.UserID).Return(entity.User{}, repoerrs.ErrNotFound)
			},
			wantErr: se.ErrNotFoundUser,
		},
		{
			name: "Repo error",
			args: args{
				ctx: context.Background(),
				input: sd.GenTokenInput{
					UserID: user.UserID,
					Role:   entity.RoleUser,
				},
			},
			mockBehavior: func(o *repomocks.MockUsers, args args) {
				o.EXPECT().GetUserByID(args.ctx, args.input.UserID).Return(entity.User{}, errors.New("other error"))
			},
			wantErr: se.ErrCannotGetUser,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUserRepo := repomocks.NewMockUsers(ctrl)
			tc.mockBehavior(mockUserRepo, tc.args)

			s := service.NewAuthService(mockUserRepo, secretKey, tokenTTL)
			got, err := s.GenerateToken(tc.args.ctx, tc.args.input)

			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
				assert.Empty(t, got)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, got)
			}
		})
	}
}

func TestAuthService_ParseToken(t *testing.T) {
	const (
		secretKey = "secret"
		tokenTTL  = time.Minute
	)

	isActive := true
	user := entity.User{
		IsActive: &isActive,
		UserID:   "u1",
		Username: "test",
		TeamName: "team",
	}

	claims := &entity.TokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		UserID: user.UserID,
		Role:   entity.RoleUser,
	}

	tokenStr, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secretKey))
	if err != nil {
		t.Fatal(err)
	}

	s := service.NewAuthService(nil, secretKey, tokenTTL)

	t.Run("Parse valid token", func(t *testing.T) {
		parsed, err := s.ParseToken(tokenStr)
		assert.NoError(t, err)
		assert.Equal(t, user.UserID, parsed.UserID)
		assert.Equal(t, entity.RoleUser, parsed.Role)
	})

	t.Run("Parse expired token", func(t *testing.T) {
		expiredClaims := &entity.TokenClaims{
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(-time.Minute).Unix(),
				IssuedAt:  time.Now().Add(-2 * time.Minute).Unix(),
			},
			UserID: user.UserID,
			Role:   entity.RoleUser,
		}
		expiredToken, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, expiredClaims).SignedString([]byte(secretKey))

		_, err := s.ParseToken(expiredToken)
		assert.ErrorIs(t, err, se.ErrTokenExpired)
	})

	t.Run("Parse invalid token", func(t *testing.T) {
		_, err := s.ParseToken("invalid_token")
		assert.ErrorIs(t, err, se.ErrCannotParseToken)
	})
}
