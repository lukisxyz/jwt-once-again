package domain

import (
	"errors"
	"net/mail"
	"time"

	"github.com/flukis/inboice/services/utils/hashing"
	"github.com/flukis/inboice/services/utils/random"
	"github.com/oklog/ulid/v2"
	"gopkg.in/guregu/null.v4"
)

var (
	ErrAccountNotFound               = errors.New("account: not found")
	ErrAccountEmailAlreadyRegistered = errors.New("account: email already registered")
)

type Account struct {
	ID               ulid.ULID `json:"id"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        null.Time `json:"updated_at"`
	DeletedAt        null.Time `json:"deleted_at"`
	Name             string    `json:"name"`
	Password         string    `json:"-"`
	Email            string    `json:"email"`
	CodeVerification string    `json:"-"`
	EmailVerifiedAt  null.Time `json:"email_verified_at"`
}

func NewAccount(email, password string) (Account, error) {
	id := ulid.Make()
	addr, err := mail.ParseAddress(email)
	if err != nil {
		return Account{}, err
	}
	username := addr.Name
	hashedPassword, err := hashing.HashPassword(password)
	if err != nil {
		return Account{}, err
	}
	code := random.RandString(6)
	acc := Account{
		ID:               id,
		CreatedAt:        time.Now(),
		Name:             username,
		Password:         hashedPassword,
		Email:            email,
		CodeVerification: code,
	}
	return acc, nil
}
