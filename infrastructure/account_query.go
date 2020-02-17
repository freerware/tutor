package infrastructure

import (
	"database/sql"
	"github.com/freerware/tutor/domain"
)

type AccountQuery interface {
	Execute() ([]domain.Account, error)
}

type accountQuery struct {
	db *sql.DB
}
