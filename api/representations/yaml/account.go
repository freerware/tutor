package yaml

import (
	"time"

	r "github.com/freerware/tutor/api/representations"
	"github.com/freerware/tutor/domain"
	u "github.com/gofrs/uuid"
)

type Account struct {
	r.Representation `yaml:"-"`

	UUID              u.UUID     `yaml:"uuid"`
	PrimaryCredential string     `yaml:"primaryCredential"`
	GivenName         string     `yaml:"givenName"`
	Surname           string     `yaml:"surname"`
	Posts             []Post     `yaml:"posts"`
	CreatedAt         time.Time  `yaml:"createdAt"`
	UpdatedAt         time.Time  `yaml:"updatedAt"`
	DeletedAt         *time.Time `yaml:"deletedAt"`
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
	account.SetContentType("application/yaml")
	account.SetSourceQuality(1.0)
	account.SetContentEncoding([]string{"identity"})
	return account
}
