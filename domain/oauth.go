package domain

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/oklog/ulid/v2"
)

type ClaimResponse struct {
	ID    ulid.ULID `json:"id"`
	Name  string    `json:"name"`
	Email string    `json:"email"`
	jwt.RegisteredClaims
}

type MapClaimResponse struct {
	ID    ulid.ULID `json:"id"`
	Name  string    `json:"name"`
	Email string    `json:"email"`
	jwt.MapClaims
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Type         string `json:"Bearer"`
	ExpiredAt    string `json:"expired_at"`
	Scope        string `json:"scope"`
}
