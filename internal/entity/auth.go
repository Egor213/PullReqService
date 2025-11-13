package entity

import "github.com/golang-jwt/jwt"

type Role string

const (
	RoleAdmin Role = "admin"
	RoleUser  Role = "user"
)

type TokenClaims struct {
	jwt.StandardClaims
	UserID string
	Role   Role
}

type ParsedToken struct {
	UserID string
	Role   Role
}
