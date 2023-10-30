package authentication

import (
	"context"
	"time"

	"github.com/flukis/inboice/services/domain"
	"github.com/flukis/inboice/services/infra/querier"
	"github.com/flukis/inboice/services/utils/hashing"
	"github.com/flukis/inboice/services/utils/random"
	"github.com/golang-jwt/jwt/v5"
)

type registerAccount struct {
	query          querier.AccountQuerier
	authQuery      querier.RefreshTokenQuerier
	secret         string
	refreshExpTime uint
	accessExpTime  uint
}

func (r *registerAccount) Login(ctx context.Context, email, password string) (account *domain.LoginResponse, err error) {
	acc, err := r.query.FindByEmail(ctx, email)
	if err != nil {
		return
	}

	if errHash := hashing.CheckPassword(
		password,
		acc.Password,
	); errHash != nil {
		err = domain.ErrPasswordNotMatch
		return
	}

	refreshExpTime := time.Now().Add(time.Duration(r.refreshExpTime) * 24 * time.Hour)

	jwtKey := []byte(r.secret)
	tokenRefreshString := random.RandString(24)

	accessExpTime := time.Now().Add(time.Duration(r.accessExpTime) * time.Minute)
	claims := &domain.ClaimResponse{
		ID:    acc.ID,
		Name:  acc.Name,
		Email: acc.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(accessExpTime),
		},
	}
	tokenAccess := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenAccessString, errSign := tokenAccess.SignedString(jwtKey)
	if errSign != nil {
		err = errSign
		return
	}

	oauthToken := domain.LoginResponse{
		AccessToken:  tokenAccessString,
		RefreshToken: tokenRefreshString,
		Type:         "Bearer",
		ExpiredAt:    refreshExpTime.Format(time.RFC3339),
		Scope:        "*",
	}

	refreshToken := domain.NewRefreshToken(
		acc.ID,
		tokenRefreshString,
		refreshExpTime,
	)

	_, err = r.authQuery.Save(ctx, &refreshToken)
	if err != nil {
		return
	}

	return &oauthToken, nil
}

type RegisterAccount interface {
	Login(ctx context.Context, email, password string) (res *domain.LoginResponse, err error)
}

func New(query querier.AccountQuerier, authQuery querier.RefreshTokenQuerier, secret string, refreshExpTime, accessExpTime uint) RegisterAccount {
	return &registerAccount{query: query, authQuery: authQuery, secret: secret, refreshExpTime: refreshExpTime, accessExpTime: accessExpTime}
}
