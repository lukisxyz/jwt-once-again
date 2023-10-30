package querier

import (
	"context"

	"github.com/flukis/inboice/services/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
)

type authDb struct {
	db *pgxpool.Pool
}

func (d *authDb) FindByUserId(ctx context.Context, id ulid.ULID) (*domain.RefreshToken, error) {
	query := `
		SELECT
			id,
			token_value,
			user_id,
			created_at,
			expires_at,
			user
		FROM
			refresh_tokens
		WHERE
			user_id = $1
			AND expires_at > NOW()  
			AND revoked = false;   
	`
	row := d.db.QueryRow(
		ctx,
		query,
		id,
	)
	var item domain.RefreshToken
	if err := row.Scan(
		item.ID,
		item.TokenValue,
		item.UserID,
		item.CreatedAt,
		item.ExpiresAt,
		item.Revoked,
	); err != nil {
		if err == pgx.ErrNoRows {
			log.Debug().Err(err).Msg("can't find any item")
			return &domain.RefreshToken{}, domain.ErrTokenNotFound
		}
		return &domain.RefreshToken{}, err
	}
	return &item, nil
}

func (d *authDb) FindByToken(ctx context.Context, token string) (*domain.RefreshToken, error) {
	query := `
		SELECT
			id,
			token_value,
			user_id,
			created_at,
			expires_at,
			user
		FROM
			refresh_tokens
		WHERE
			token = $1
			AND expires_at > NOW()  
			AND revoked = false;   
	`
	row := d.db.QueryRow(
		ctx,
		query,
		token,
	)
	var item domain.RefreshToken
	if err := row.Scan(
		item.ID,
		item.TokenValue,
		item.UserID,
		item.CreatedAt,
		item.ExpiresAt,
		item.Revoked,
	); err != nil {
		if err == pgx.ErrNoRows {
			log.Debug().Err(err).Msg("can't find any item")
			return &domain.RefreshToken{}, domain.ErrTokenNotFound
		}
		return &domain.RefreshToken{}, err
	}
	return &item, nil
}

func (d *authDb) FindById(ctx context.Context, id ulid.ULID) (*domain.RefreshToken, error) {
	query := `
		SELECT
			id,
			token_value,
			user_id,
			created_at,
			expires_at,
			user
		FROM
			refresh_tokens
		WHERE
			id = $1
			AND expires_at > NOW()  
			AND revoked = false;   
	`
	row := d.db.QueryRow(
		ctx,
		query,
		id,
	)
	var item domain.RefreshToken
	if err := row.Scan(
		item.ID,
		item.TokenValue,
		item.UserID,
		item.CreatedAt,
		item.ExpiresAt,
		item.Revoked,
	); err != nil {
		if err == pgx.ErrNoRows {
			log.Debug().Err(err).Msg("can't find any item")
			return &domain.RefreshToken{}, domain.ErrTokenNotFound
		}
		return &domain.RefreshToken{}, err
	}
	return &item, nil
}

func (d *authDb) Save(ctx context.Context, data *domain.RefreshToken) (*domain.RefreshToken, error) {
	query := `
		INSERT INTO refresh_tokens
			(id, token_value, user_id, created_at, expires_at, revoked)
		VALUES
			($1, $2, $3, $4, $5, $6);
	`

	if _, err := d.db.Exec(
		ctx,
		query,
		data.ID,
		data.TokenValue,
		data.UserID,
		data.CreatedAt,
		data.ExpiresAt,
		data.Revoked,
	); err != nil {
		return &domain.RefreshToken{}, err
	}

	return data, nil
}

func (d *authDb) Revoke(ctx context.Context, id ulid.ULID) error {
	query := `
		UPDATE refresh_tokens
		SET revoked = TRUE
		WHERE id = $1;
	`
	if _, err := d.db.Exec(
		ctx,
		query,
		id,
	); err != nil {
		return err
	}
	return nil
}

type RefreshTokenQuerier interface {
	FindById(ctx context.Context, id ulid.ULID) (*domain.RefreshToken, error)
	FindByUserId(ctx context.Context, id ulid.ULID) (*domain.RefreshToken, error)
	FindByToken(ctx context.Context, token string) (*domain.RefreshToken, error)
	Revoke(ctx context.Context, id ulid.ULID) error
	Save(ctx context.Context, data *domain.RefreshToken) (*domain.RefreshToken, error)
}

func NewRefreshToken(db *pgxpool.Pool) RefreshTokenQuerier {
	return &authDb{db: db}
}
