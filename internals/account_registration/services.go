package accountRegistration

import (
	"context"

	"github.com/flukis/inboice/services/domain"
	"github.com/flukis/inboice/services/infra/querier"
	"github.com/oklog/ulid/v2"
)

type registerAccount struct {
	query querier.AccountQuerier
}

// Delete implements RegisterAccount.
func (r *registerAccount) Delete(ctx context.Context, id ulid.ULID) (err error) {
	acc, err := r.query.FindById(ctx, id)
	if err != nil {
		return
	}

	if !acc.DeletedAt.IsZero() {
		return domain.ErrAccountAlreadyDeleted
	}

	return r.query.Delete(ctx, acc.ID)
}

func (r *registerAccount) GetByID(ctx context.Context, id ulid.ULID) (account *domain.GetAccountResponse, err error) {
	acc, err := r.query.FindById(ctx, id)
	if err != nil {
		return
	}

	account.Email = acc.Email
	account.ID = acc.ID
	account.EmailIsVerified = !acc.EmailVerifiedAt.IsZero()
	account.Name = acc.Name
	return
}

func (r *registerAccount) GetByEmail(ctx context.Context, email string) (account *domain.GetAccountResponse, err error) {
	acc, err := r.query.FindByEmail(ctx, email)
	if err != nil {
		return
	}

	res := domain.GetAccountResponse{
		Name:            acc.Name,
		Email:           acc.Email,
		EmailIsVerified: !acc.EmailVerifiedAt.IsZero(),
		ID:              acc.ID,
	}

	account = &res

	return
}

// Register implements RegisterAccount.
func (r *registerAccount) Register(ctx context.Context, data domain.RegistrationRequest) (res *domain.RegistrationResponse, err error) {
	newAcc, err := domain.NewAccount(data.Email, data.Password)
	if err != nil {
		return
	}

	acc, err := r.query.Save(ctx, &newAcc)
	if err != nil {
		return
	}

	resuls := domain.RegistrationResponse{
		Id: acc.ID,
	}

	res = &resuls
	return
}

type RegisterAccount interface {
	Register(ctx context.Context, data domain.RegistrationRequest) (res *domain.RegistrationResponse, err error)
	GetByID(ctx context.Context, id ulid.ULID) (res *domain.GetAccountResponse, err error)
	GetByEmail(ctx context.Context, email string) (res *domain.GetAccountResponse, err error)
	Delete(ctx context.Context, id ulid.ULID) (err error)
}

func New(query querier.AccountQuerier) RegisterAccount {
	return &registerAccount{query: query}
}
