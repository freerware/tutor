package domain

import (
	"errors"
	"time"

	u "github.com/gofrs/uuid"
)

// Errors that are potentially thrown during account interactions.
var (
	ErrFutureCreatedAt  = errors.New("domain: account creation time cannot be in the future")
	ErrFutureUpdatedAt  = errors.New("domain: account modification time cannot be in the future")
	ErrInvalidUpdatedAt = errors.New("domain: account modification time cannot be prior to account creation time")
	ErrFutureDeletedAt  = errors.New("domain: account deleton time cannot be in the future")
	ErrInvalidDeletedAt = errors.New("domain: account deletion time cannot be prior to account creation or modification time")
)

type Account struct {
	uuid      u.UUID
	givenName string
	surname   string
	username  string
	createdAt time.Time
	updatedAt time.Time
	deletedAt *time.Time
}

type AccountParameters struct {
	UUID      u.UUID
	GivenName string
	Surname   string
	Username  string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

func NewAccount(parameters AccountParameters) (Account, error) {
	account := Account{}
	account.SetUUID(parameters.UUID)
	account.SetGivenName(parameters.GivenName)
	account.SetSurname(parameters.Surname)
	account.SetUsername(parameters.Username)
	if err := account.SetCreatedAt(parameters.CreatedAt); err != nil {
		return Account{}, err
	}
	if err := account.SetUpdatedAt(parameters.UpdatedAt); err != nil {
		return Account{}, err
	}
	if parameters.DeletedAt != nil {
		if err := account.SetDeletedAt(*parameters.DeletedAt); err != nil {
			return Account{}, err
		}
	}
	return account, nil
}

func ReconstituteAccount(parameters AccountParameters) Account {
	return Account{
		uuid:      parameters.UUID,
		givenName: parameters.GivenName,
		surname:   parameters.Surname,
		username:  parameters.Username,
		createdAt: parameters.CreatedAt,
		updatedAt: parameters.UpdatedAt,
		deletedAt: parameters.DeletedAt,
	}
}

func (a Account) UUID() u.UUID {
	return a.uuid
}

func (a *Account) SetUUID(uuid u.UUID) {
	a.uuid = uuid
}

func (a Account) GivenName() string {
	return a.givenName
}

func (a *Account) SetGivenName(name string) {
	a.givenName = name
}

func (a Account) Surname() string {
	return a.surname
}

func (a *Account) SetSurname(name string) {
	a.surname = name
}

func (a Account) Username() string {
	return a.username
}

func (a *Account) SetUsername(username string) {
	a.username = username
}

func (a Account) CreatedAt() time.Time {
	return a.createdAt
}

func (a *Account) SetCreatedAt(t time.Time) error {
	if t.After(time.Now()) {
		return ErrFutureCreatedAt
	}
	a.createdAt = t
	return nil
}

func (a Account) UpdatedAt() time.Time {
	return a.updatedAt
}

func (a *Account) SetUpdatedAt(t time.Time) error {
	if t.After(time.Now()) {
		return ErrFutureUpdatedAt
	}
	if t.Before(a.CreatedAt()) {
		return ErrInvalidUpdatedAt
	}
	a.updatedAt = t
	return nil
}

func (a Account) DeletedAt() *time.Time {
	return a.deletedAt
}

func (a *Account) SetDeletedAt(t time.Time) error {
	if t.After(time.Now()) {
		return ErrFutureDeletedAt
	}
	if t.Before(a.CreatedAt()) || t.Before(a.UpdatedAt()) {
		return ErrInvalidDeletedAt
	}
	a.deletedAt = &t
	return nil
}
