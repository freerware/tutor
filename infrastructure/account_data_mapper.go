package infrastructure

import (
	"context"
	"database/sql"
	"errors"

	"github.com/freerware/morph"
	"github.com/freerware/tutor/domain"
	"github.com/freerware/work/v4/unit"
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
	table  morph.Table
}

func NewAccountDataMapper(parameters AccountDataMapperParameters) AccountDataMapper {
	opts := []morph.ReflectOption{
		morph.WithPrimaryKeyColumn("UUID"),
		morph.WithInferredTableName(morph.ScreamingSnakeCaseStrategy, false),
		morph.WithInferredColumnNames(morph.ScreamingSnakeCaseStrategy),
		morph.WithInferredTableAlias(morph.UpperCaseStrategy, 1),
		morph.WithColumnNameMapping("Username", "PRIMARY_CREDENTIAL"),
	}
	t, err := morph.Reflect(domain.Account{}, opts...)
	if err != nil {
		panic(err)
	}

	return AccountDataMapper{db: parameters.DB, logger: parameters.Logger, table: t}
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

func (dm *AccountDataMapper) Insert(ctx context.Context, mCtx unit.MapperContext, accounts ...interface{}) error {
	if len(accounts) == 0 {
		return nil
	}
	accs, err := dm.toAccount(accounts...)
	if err != nil {
		return err
	}
	return dm.insert(mCtx.Tx, accs...)
}

func (dm *AccountDataMapper) insertSQL(accounts ...domain.Account) (sql string, args [][]interface{}) {
	var err error
	sql, err = dm.table.InsertQuery()
	if err != nil {
		panic(err)
	}

	for _, account := range accounts {
		sArgs := []interface{}{
			account.CreatedAt(),
			account.DeletedAt(),
			account.GivenName(),
			account.Username(),
			account.Surname(),
			account.UpdatedAt(),
			account.UUID().String(),
		}
		args = append(args, sArgs)
	}
	return
}

func (dm *AccountDataMapper) insert(tx *sql.Tx, accounts ...domain.Account) error {
	// insert accounts.
	sql, args := dm.insertSQL(accounts...)
	for _, a := range args {
		if err := dm.prepareAndExec(tx, sql, a); err != nil {
			return err
		}
	}
	return nil
}

func (dm *AccountDataMapper) Update(ctx context.Context, mCtx unit.MapperContext, accounts ...interface{}) error {
	if len(accounts) == 0 {
		return nil
	}
	accs, err := dm.toAccount(accounts...)
	if err != nil {
		return err
	}
	return dm.update(mCtx.Tx, accs...)
}

func (dm *AccountDataMapper) updateSQL(accounts ...domain.Account) (sql string, args [][]interface{}) {
	var err error
	sql, err = dm.table.UpdateQuery()
	if err != nil {
		panic(err)
	}

	for _, account := range accounts {
		sArgs := []interface{}{
			account.CreatedAt(),
			account.DeletedAt(),
			account.GivenName(),
			account.Username(),
			account.Surname(),
			account.UpdatedAt(),
			account.UUID().String(),
		}
		args = append(args, sArgs)
	}
	return
}

func (dm *AccountDataMapper) update(tx *sql.Tx, accounts ...domain.Account) error {
	sql, args := dm.updateSQL(accounts...)
	for _, a := range args {
		if err := dm.prepareAndExec(tx, sql, a); err != nil {
			return err
		}
	}
	return nil
}

func (dm *AccountDataMapper) Delete(ctx context.Context, mCtx unit.MapperContext, accounts ...interface{}) error {
	if len(accounts) == 0 {
		return nil
	}
	accs, err := dm.toAccount(accounts...)
	if err != nil {
		return err
	}
	return dm.delete(mCtx.Tx, accs...)
}

func (dm *AccountDataMapper) deleteSQL(accounts ...domain.Account) (sql string, args [][]interface{}) {
	var err error
	sql, err = dm.table.DeleteQuery()
	if err != nil {
		panic(err)
	}

	for _, account := range accounts {
		args = append(args, []interface{}{account.UUID().String()})
	}
	return
}

func (dm *AccountDataMapper) delete(tx *sql.Tx, accounts ...domain.Account) error {
	sql, args := dm.deleteSQL(accounts...)
	for _, a := range args {
		if err := dm.prepareAndExec(tx, sql, a); err != nil {
			return err
		}
	}
	return nil
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
