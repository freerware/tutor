package infrastructure

import (
	"errors"

	"github.com/freerware/tutor/domain"
	"github.com/freerware/work/v4/unit"
	u "github.com/gofrs/uuid"
)

// AccountRepository represents a collection of all
// accounts within the application.
type AccountRepository interface {
	Get(u.UUID) (*domain.Account, error)
	Put(domain.Account) error
	Remove(domain.Account) error
	Add(domain.Account) error
	Find(AccountQuery) ([]domain.Account, error)
	Size() (int, error)
}

type accountRepository struct {
	unit    unit.Unit
	queryer Queryer
}

func NewAccountRepository(unit unit.Unit, queryer Queryer) AccountRepository {
	return &accountRepository{unit: unit, queryer: queryer}
}

func (r *accountRepository) Find(query AccountQuery) ([]domain.Account, error) {
	return query.Execute()
}

func (r *accountRepository) Get(uuid u.UUID) (*domain.Account, error) {
	query := r.queryer.Query(uuid)
	matches, err := r.Find(query)
	if err != nil {
		return nil, err
	}
	if len(matches) == 0 {
		return nil, nil
	}
	a := matches[0]
	return &a, nil
}

func (r *accountRepository) Put(account domain.Account) error {

	// check if the account exists.
	c, e := r.Get(account.UUID())
	if e != nil {
		return e
	}

	// if the account is not within the repository, add it.
	if c == nil {
		return r.Add(account)
	}

	// otherwise, replace the existing state.
	return r.unit.Alter(account)
}

func (r *accountRepository) Remove(account domain.Account) error {

	// check if the account exists.
	c, e := r.Get(account.UUID())
	if e != nil {
		return e
	}

	// if the account is not within the repository, throw an error.
	if c == nil {
		return errors.New("could not find the account")
	}

	// otherwise, remove the account.
	r.unit.Remove(*c)
	return nil
}

func (r *accountRepository) Add(account domain.Account) error {

	// check if the account exists.
	c, e := r.Get(account.UUID())
	if e != nil {
		return e
	}

	// if the account is within the repository, throw an error.
	if c != nil {
		return errors.New("account already exists")
	}

	// otherwise, remove the account.
	r.unit.Add(account)
	return nil
}

func (r *accountRepository) Size() (int, error) {
	//TODO(FREER) use queryer to get query for select all.
	return 0, nil
}
