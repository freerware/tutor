package infrastructure

import (
	"database/sql"
	"time"

	"github.com/freerware/tutor/domain"
	"github.com/go-sql-driver/mysql"
	u "github.com/gofrs/uuid"
)

type findAccountByUUID struct {
	accountQuery

	uuid u.UUID
}

func NewFindAccountByUUIDQuery(db *sql.DB, uuid u.UUID) AccountQuery {
	return &findAccountByUUID{
		accountQuery: accountQuery{
			db: db,
		},
		uuid: uuid,
	}
}

func (q *findAccountByUUID) Execute() ([]domain.Account, error) {

	// retrieve accounts.
	return q.accounts()
}

func (q findAccountByUUID) accounts() ([]domain.Account, error) {

	var matches []domain.Account
	statement, err := q.db.Prepare("SELECT * FROM ACCOUNT WHERE UUID = ?")
	if err != nil {
		return []domain.Account{}, err
	}
	defer statement.Close()

	rows, err := statement.Query(q.uuid.String())
	if err != nil {
		return matches, err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			uuid                                  u.UUID
			givenName, surname, primaryCredential string
			updatedAt, createdAt                  time.Time
			deletedAt                             mysql.NullTime
		)

		fields := []any{
			&uuid,
			&givenName,
			&surname,
			&primaryCredential,
			&createdAt,
			&updatedAt,
			&deletedAt,
		}
		if err := rows.Scan(fields...); err != nil {
			return matches, err
		}

		a := domain.Account{}
		a.SetUUID(uuid)
		a.SetUsername(primaryCredential)
		a.SetGivenName(givenName)
		a.SetSurname(surname)
		a.SetCreatedAt(createdAt)
		a.SetUpdatedAt(updatedAt)
		if deletedAt.Valid {
			a.SetDeletedAt(deletedAt.Time)
		}
		matches = append(matches, a)
	}
	return matches, nil
}
