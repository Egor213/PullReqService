package service

import (
	"app/internal/entity"
	e "app/internal/entity"
	"app/internal/repo"
	"app/internal/repo/repoerrs"
	"app/internal/service/servdto"
	"app/internal/service/serverrs"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

type AuthService struct {
	userRepo repo.Users
	signKey  string
	tokenTTL time.Duration
}

func NewAuthService(userRepo repo.Users, signKey string, tokenTTL time.Duration) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		signKey:  signKey,
		tokenTTL: tokenTTL,
	}
}

func (s *AuthService) GenerateToken(ctx context.Context, in servdto.GenTokenInput) (string, error) {
	user, err := s.userRepo.GetUserByID(ctx, in.UserID)
	if err != nil {
		if errors.Is(err, repoerrs.ErrNotFound) {
			return "", serverrs.ErrUserNotFound
		}
		return "", serverrs.ErrCannotGetUser
	}

	claims := &entity.TokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(s.tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		UserID: user.UserID,
		Role:   in.Role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(s.signKey))
	if err != nil {
		return "", serverrs.ErrCannotSignToken
	}

	return signed, nil
}

func (s *AuthService) ParseToken(accessToken string) (e.ParsedToken, error) {
	token, err := jwt.ParseWithClaims(accessToken, &e.TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.signKey), nil
	})

	print("AAAAAAAAA")
	if err != nil {
		return e.ParsedToken{}, serverrs.ErrCannotParseToken
	}
	print("BBBBBBB")

	claims, ok := token.Claims.(*e.TokenClaims)
	if !ok || !token.Valid {
		return e.ParsedToken{}, serverrs.ErrCannotParseToken
	}
	print("CCCCCCCC")

	return e.ParsedToken{
		UserID: claims.UserID,
		Role:   claims.Role,
	}, nil
}
