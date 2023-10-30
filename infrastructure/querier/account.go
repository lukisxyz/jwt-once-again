package querier

import (
	"context"

	"github.com/flukis/inboice/services/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
)

type accountDb struct {
	db *pgxpool.Pool
}

// delete implements AccountQuerier.
func (d *accountDb) Delete(ctx context.Context, id ulid.ULID) error {
	query := `
		UPDATE accounts
		SET deleted_at = NOW()
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

// findById implements AccountQuerier.
func (d *accountDb) FindById(ctx context.Context, id ulid.ULID) (*domain.Account, error) {
	query := `
		SELECT
			id,
			email,
			name,
			password,
			email_verified_at,
			code_verification,
			created_at,
			updated_at,
			deleted_at
		FROM
			accounts
		WHERE
			id = $1
	`
	row := d.db.QueryRow(
		ctx,
		query,
		id,
	)
	var item domain.Account
	if err := row.Scan(
		&item.ID,
		&item.Email,
		&item.Name,
		&item.Password,
		&item.EmailVerifiedAt,
		&item.CodeVerification,
		&item.CreatedAt,
		&item.UpdatedAt,
		&item.DeletedAt,
	); err != nil {
		if err == pgx.ErrNoRows {
			log.Debug().Err(err).Msg("can't find any item")
			return &domain.Account{}, domain.ErrAccountNotFound
		}
		return &domain.Account{}, err
	}
	return &item, nil
}

// findById implements AccountQuerier.
func (d *accountDb) FindByEmail(ctx context.Context, email string) (*domain.Account, error) {
	query := `
		SELECT
			id,
			email,
			name,
			password,
			email_verified_at,
			code_verification,
			created_at,
			updated_at,
			deleted_at
		FROM
			accounts
		WHERE
			email = $1
	`
	row := d.db.QueryRow(
		ctx,
		query,
		email,
	)
	var item domain.Account
	if err := row.Scan(
		&item.ID,
		&item.Email,
		&item.Name,
		&item.Password,
		&item.EmailVerifiedAt,
		&item.CodeVerification,
		&item.CreatedAt,
		&item.UpdatedAt,
		&item.DeletedAt,
	); err != nil {
		if err == pgx.ErrNoRows {
			log.Debug().Err(err).Msg("can't find any item")
			return &domain.Account{}, domain.ErrAccountNotFound
		}
		return &domain.Account{}, err
	}
	return &item, nil
}

// save implements AccountQuerier.
func (d *accountDb) Save(ctx context.Context, data *domain.Account) (*domain.Account, error) {
	query := `
		INSERT INTO accounts
			(id, created_at, name, password, email, code_verification)
		VALUES
			($1, $2, $4, $5, $6, $7)
		ON CONFLICT (id) DO UPDATE
		SET
			updated_at = $3,
			name = $4,
			password = $5,
			email = $6,
			code_verification = $7,
			email_verified_at = $8;
	`

	if _, err := d.db.Exec(
		ctx,
		query,
		data.ID,
		data.CreatedAt,
		data.UpdatedAt,
		data.Name,
		data.Password,
		data.Email,
		data.CodeVerification,
		data.EmailVerifiedAt,
	); err != nil {
		pqErr := err.(*pgconn.PgError)
		if pqErr.Code == "23505" {
			return &domain.Account{}, domain.ErrAccountEmailAlreadyRegistered
		}
		return &domain.Account{}, err
	}

	return data, nil
}

type AccountQuerier interface {
	Save(ctx context.Context, data *domain.Account) (*domain.Account, error)
	FindById(ctx context.Context, id ulid.ULID) (*domain.Account, error)
	FindByEmail(ctx context.Context, email string) (*domain.Account, error)
	Delete(ctx context.Context, id ulid.ULID) error
}

func NewAccount(db *pgxpool.Pool) AccountQuerier {
	return &accountDb{db: db}
}
