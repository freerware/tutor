package infrastructure

import (
	"database/sql"

	"github.com/freerware/tutor/domain"
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
	matches := []domain.Account{}
	statement, err := q.db.Prepare("SELECT CREATED_AT, DELETED_AT, GIVEN_NAME, PRIMARY_CREDENTIAL, SURNAME, UPDATED_AT, UUID FROM ACCOUNT WHERE UUID = ?")
	if err != nil {
		return matches, err
	}
	defer statement.Close()

	rows, err := statement.Query(q.uuid.String())
	if err != nil {
		return matches, err
	}
	defer rows.Close()

	for rows.Next() {
		var params domain.AccountParameters
		err = rows.Scan(
			&params.CreatedAt,
			&params.DeletedAt,
			&params.GivenName,
			&params.Username,
			&params.Surname,
			&params.UpdatedAt,
			&params.UUID,
		)
		if err != nil {
			return matches, err
		}
		a := domain.ReconstituteAccount(params)

		pStatement, err := q.db.Prepare("SELECT AUTHOR_UUID, CREATED_AT, DELETED_AT, DRAFT, LIKE_COUNT, UPDATED_AT, UUID, TITLE, CONTENT FROM POST WHERE AUTHOR_UUID = ?;")
		if err != nil {
			return matches, err
		}
		defer pStatement.Close()

		pRows, err := pStatement.Query(q.uuid.String())
		if err != nil {
			return matches, err
		}
		defer pRows.Close()

		for pRows.Next() {
			var params domain.PostParameters
			err = pRows.Scan(
				&params.AuthorUUID,
				&params.CreatedAt,
				&params.DeletedAt,
				&params.Draft,
				&params.Likes,
				&params.UpdatedAt,
				&params.UUID,
				&params.Title,
				&params.Content,
			)
			if err != nil {
				return matches, err
			}
			a.AddPost(domain.ReconstitutePost(params))
		}

		matches = append(matches, a)
	}
	return matches, nil
}
