package models

import "github.com/golang-jwt/jwt/v5"

type UserRefreshToken struct {
	UserID       string `json:"user_id"`
	RefreshToken string `json:"refresh_token"`
}

type JwtUserRefreshToken struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

type JwtUserAccessToken struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}
