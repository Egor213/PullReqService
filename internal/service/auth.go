package service

import (
	e "app/internal/entity"
	"app/internal/repo"
	repoerrs "app/internal/repo/errors"
	sd "app/internal/service/dto"
	se "app/internal/service/errors"
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

func (s *AuthService) GenerateToken(ctx context.Context, in sd.GenTokenInput) (string, error) {
	_, err := s.userRepo.GetUserByID(ctx, in.UserID)
	if err != nil {
		if errors.Is(err, repoerrs.ErrNotFound) {
			return "", se.ErrUserNotFound
		}
		return "", se.ErrCannotGetUser
	}
	exTime := time.Now().Add(s.tokenTTL).Unix()
	claims := &e.TokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: exTime,
			IssuedAt:  time.Now().Unix(),
		},
		UserID: in.UserID,
		Role:   in.Role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(s.signKey))
	if err != nil {
		return "", se.ErrCannotSignToken
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

	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return e.ParsedToken{}, se.ErrTokenExpired
			}
		}
		return e.ParsedToken{}, se.ErrCannotParseToken
	}

	claims, ok := token.Claims.(*e.TokenClaims)
	if !ok || !token.Valid {
		return e.ParsedToken{}, se.ErrCannotParseToken
	}

	return e.ParsedToken{
		UserID: claims.UserID,
		Role:   claims.Role,
	}, nil
}
