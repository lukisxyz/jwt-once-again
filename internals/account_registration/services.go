package accountRegistration

import (
	"context"

	"github.com/flukis/inboice/services/domain"
	"github.com/flukis/inboice/services/infrastructure/querier"
	"github.com/oklog/ulid/v2"
)

type registerAccount struct {
	query querier.AccountQuerier
}

// Delete implements RegisterAccount.
func (*registerAccount) Delete(ctx context.Context, id ulid.ULID) (err error) {
	panic("unimplemented")
}

// Register implements RegisterAccount.
func (r *registerAccount) Register(ctx context.Context, data domain.RegistrationRequest) (res domain.RegistrationResponse, err error) {
	newAcc, err := domain.NewAccount(data.Email, data.Password)
	if err != nil {
		return
	}

	acc, err := r.query.Save(ctx, &newAcc)
	if err != nil {
		return
	}

	res.Id = acc.ID

	return
}

type RegisterAccount interface {
	Register(ctx context.Context, data domain.RegistrationRequest) (res domain.RegistrationResponse, err error)
	Delete(ctx context.Context, id ulid.ULID) (err error)
}

func New(query querier.AccountQuerier) RegisterAccount {
	return &registerAccount{query: query}
}
