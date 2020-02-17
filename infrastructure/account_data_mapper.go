package infrastructure

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/freerware/tutor/domain"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

// Errors that are potentially thrown during data mapper interactions.
var (

	// ErrInvalidType represents an error that indicates a unexpected type
	// was provided to the data mapper.
	ErrInvalidType = errors.New("infrastructure: invalid type provided to data mapper")
)

type AccountDataMapperParameters struct {
	fx.In

	DB     *sql.DB `name:"rwDB"`
	Logger *zap.Logger
}

type AccountDataMapper struct {
	db     *sql.DB
	logger *zap.Logger
}

func NewAccountDataMapper(
	parameters AccountDataMapperParameters) AccountDataMapper {
	return AccountDataMapper{db: parameters.DB, logger: parameters.Logger}
}

func (dm *AccountDataMapper) toAccount(accounts ...interface{}) ([]domain.Account, error) {
	accs := []domain.Account{}
	for _, account := range accounts {
		var acc domain.Account
		var ok bool
		if acc, ok = account.(domain.Account); !ok {
			return []domain.Account{}, ErrInvalidType
		}
		accs = append(accs, acc)
	}
	return accs, nil
}

func (dm *AccountDataMapper) Insert(tx *sql.Tx, accounts ...interface{}) error {
	if len(accounts) == 0 {
		return nil
	}
	accs, err := dm.toAccount(accounts...)
	if err != nil {
		return err
	}
	return dm.insert(tx, accs...)
}

func (dm *AccountDataMapper) insertSQL(accounts ...domain.Account) (sql string, args []interface{}) {
	sql =
		`INSERT INTO ACCOUNT
		(
			UUID,
			GIVEN_NAME,
			SURNAME,
			PRIMARY_CREDENTIAL,
			CREATED_AT,
			UPDATED_AT,
			DELETED_AT
		) VALUES `
	var vals []string
	for _, account := range accounts {
		vals = append(vals, "(?, ?, ?, ?, ?, ?, ?)")
		args = append(args,
			account.UUID().String(),
			account.GivenName(),
			account.Surname(),
			account.Username(),
			account.CreatedAt(),
			account.UpdatedAt(),
			account.DeletedAt(),
		)
	}
	sql = sql + strings.Join(vals, ", ") + ";"
	return
}

func (dm *AccountDataMapper) insert(tx *sql.Tx, accounts ...domain.Account) error {
	// insert accounts.
	sql, args := dm.insertSQL(accounts...)
	return dm.prepareAndExec(tx, sql, args)
}

func (dm *AccountDataMapper) Update(tx *sql.Tx, accounts ...interface{}) error {
	if len(accounts) == 0 {
		return nil
	}
	accs, err := dm.toAccount(accounts...)
	if err != nil {
		return err
	}
	return dm.update(tx, accs...)
}

func (dm *AccountDataMapper) updateSQL(
	accounts ...domain.Account) (sql []string, args [][]interface{}) {
	for _, account := range accounts {
		s :=
			`UPDATE
				ACCOUNT
			SET
				GIVEN_NAME = ?,
				SURNAME = ?,
				PRIMARY_CREDENTIAL = ?,
				CREATED_AT = ?,
				UPDATED_AT = ?,
				DELETED_AT = ?
			WHERE
				UUID = ?`
		sql = append(sql, s)
		sArgs := []interface{}{
			account.GivenName(),
			account.Surname(),
			account.Username(),
			account.CreatedAt(),
			account.UpdatedAt(),
			account.DeletedAt(),
			account.UUID().String(),
		}
		args = append(args, sArgs)
	}
	return
}

func (dm *AccountDataMapper) update(tx *sql.Tx, accounts ...domain.Account) error {
	sql, args := dm.updateSQL(accounts...)
	for idx, s := range sql {
		if err := dm.prepareAndExec(tx, s, args[idx]); err != nil {
			return err
		}
	}
	return nil
}

func (dm *AccountDataMapper) Delete(tx *sql.Tx, accounts ...interface{}) error {
	if len(accounts) == 0 {
		return nil
	}
	accs, err := dm.toAccount(accounts...)
	if err != nil {
		return err
	}
	return dm.delete(tx, accs...)
}

func (dm *AccountDataMapper) deleteSQL(
	accounts ...domain.Account) (sql string, args []interface{}) {
	sql = "DELETE FROM ACCOUNT WHERE UUID IN (?)"
	for _, account := range accounts {
		args = append(args, account.UUID().String())
	}
	return
}

func (dm *AccountDataMapper) delete(tx *sql.Tx, accounts ...domain.Account) error {
	sql, args := dm.deleteSQL(accounts...)
	return dm.prepareAndExec(tx, sql, args)
}

func (dm *AccountDataMapper) prepareAndExec(
	tx *sql.Tx, sql string, args []interface{}) error {
	s, err := tx.Prepare(sql)
	if err != nil {
		return err
	}
	defer s.Close()
	_, err = s.Exec(args...)
	return err
}
