package infrastructure

import (
	"database/sql"

	u "github.com/gofrs/uuid"
	"go.uber.org/fx"
)

type Queryer interface {
	Query(u.UUID) AccountQuery
}

type queryer struct {
	db *sql.DB
}

type QueryerParameters struct {
	fx.In

	DB *sql.DB `name:"rwDB"`
}

func NewQueryer(parameters QueryerParameters) Queryer {
	return &queryer{
		db: parameters.DB,
	}
}

func (f *queryer) Query(uuid u.UUID) AccountQuery {
	return NewFindAccountByUUIDQuery(f.db, uuid)
}
