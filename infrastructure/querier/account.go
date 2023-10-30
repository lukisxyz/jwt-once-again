package querier

import (
	"context"

	"github.com/flukis/inboice/services/domain"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
)

type accountDb struct {
	db *pgxpool.Pool
}

// delete implements AccountQuerier.
func (*accountDb) Delete(ctx context.Context, id ulid.ULID) error {
	panic("unimplemented")
}

// findById implements AccountQuerier.
func (*accountDb) FindById(ctx context.Context, id ulid.ULID) (*domain.Account, error) {
	panic("unimplemented")
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
		return &domain.Account{}, err
	}
	log.Info().Msg("Sampai")

	return data, nil
}

type AccountQuerier interface {
	Save(ctx context.Context, data *domain.Account) (*domain.Account, error)
	FindById(ctx context.Context, id ulid.ULID) (*domain.Account, error)
	Delete(ctx context.Context, id ulid.ULID) error
}

func NewAccount(db *pgxpool.Pool) AccountQuerier {
	return &accountDb{db: db}
}
