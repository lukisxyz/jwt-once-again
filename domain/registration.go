package domain

import "github.com/oklog/ulid/v2"

type RegistrationRequest struct {
	Email    string `json:"email" binding:"email"`
	Password string `json:"password" binding:"required"`
}

type RegistrationResponse struct {
	Id ulid.ULID `json:"id"`
}
