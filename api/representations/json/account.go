package json

import (
	"time"

	r "github.com/freerware/tutor/api/representations"
	"github.com/freerware/tutor/domain"
	u "github.com/gofrs/uuid"
)

type Account struct {
	r.Representation `json:"-"`

	UUID              u.UUID     `json:"uuid"`
	PrimaryCredential string     `json:"primaryCredential"`
	GivenName         string     `json:"givenName"`
	Surname           string     `json:"surname"`
	Posts             []Post     `json:"posts"`
	CreatedAt         time.Time  `json:"createdAt"`
	UpdatedAt         time.Time  `json:"updatedAt"`
	DeletedAt         *time.Time `json:"deletedAt"`
}

// Bytes provides the representation as bytes.
func (a Account) Bytes() ([]byte, error) {
	return a.Base.Bytes(&a)
}

// FromBytes constructs the representation from bytes.
func (a Account) FromBytes(b []byte) error {
	return a.Base.FromBytes(b, &a)
}

// NewAccount constructs a new account representation.
func NewAccount(a domain.Account) Account {
	account := Account{
		UUID:              a.UUID(),
		GivenName:         a.GivenName(),
		Surname:           a.Surname(),
		PrimaryCredential: a.Username(),
		Posts:             NewPosts(a.Posts()...),
		CreatedAt:         a.CreatedAt(),
		UpdatedAt:         a.UpdatedAt(),
		DeletedAt:         a.DeletedAt(),
	}
	account.SetContentCharset("ascii")
	account.SetContentLanguage("en-US")
	account.SetContentType("application/json")
	account.SetSourceQuality(1.0)
	account.SetContentEncoding([]string{"identity"})
	return account
}
