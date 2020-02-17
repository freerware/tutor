package application

import (
	"errors"

	"github.com/freerware/tutor/domain"
	"github.com/freerware/tutor/infrastructure"
	"github.com/freerware/work"
	u "github.com/gofrs/uuid"
	"go.uber.org/fx"
)

// AccountService encapsulates the various operations
// our application offers for user accounts.
type AccountService struct {
	uniter  work.Uniter
	queryer infrastructure.Queryer
}

type AccountServiceParameters struct {
	fx.In

	Uniter  work.Uniter `name:"sqlWorkUniter"`
	Queryer infrastructure.Queryer
}

func NewAccountService(
	parameters AccountServiceParameters) AccountService {
	return AccountService{
		uniter:  parameters.Uniter,
		queryer: parameters.Queryer,
	}
}

// Get retrieves an existing account.
func (a *AccountService) Get(uuid u.UUID) (domain.Account, error) {
	unit, err := a.uniter.Unit()
	if err != nil {
		return domain.Account{}, err
	}
	repository := infrastructure.NewAccountRepository(unit, a.queryer)
	account, err := repository.Get(uuid)
	if err != nil {
		return domain.Account{}, err
	}
	if account == nil {
		return domain.Account{}, errors.New("application: account not found")
	}
	return *account, nil
}

// Create creates a new account.
func (a *AccountService) Create(account domain.Account) error {
	unit, err := a.uniter.Unit()
	if err != nil {
		return err
	}
	repository := infrastructure.NewAccountRepository(unit, a.queryer)
	if err = repository.Add(account); err != nil {
		return err
	}
	return unit.Save()
}

// Put upserts an account.
func (a *AccountService) Put(account domain.Account) error {
	unit, err := a.uniter.Unit()
	if err != nil {
		return err
	}
	repository := infrastructure.NewAccountRepository(unit, a.queryer)
	if err = repository.Put(account); err != nil {
		return err
	}
	return unit.Save()
}

// Delete deletes an existing account.
func (a *AccountService) Delete(account domain.Account) error {
	unit, err := a.uniter.Unit()
	if err != nil {
		return err
	}
	repository := infrastructure.NewAccountRepository(unit, a.queryer)
	if err = repository.Remove(account); err != nil {
		return err
	}
	return unit.Save()
}
