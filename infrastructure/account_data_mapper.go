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

	return AccountDataMapper{
		db:     parameters.DB,
		logger: parameters.Logger,
		table:  t,
	}
}

func (dm *AccountDataMapper) Insert(ctx context.Context, mCtx unit.MapperContext, accounts ...any) error {
	for _, account := range accounts {
		sql, args, err := dm.table.InsertQueryWithArgs(account)
		if err != nil {
			return err
		}

		stmt, err := mCtx.Tx.Prepare(sql)
		if err != nil {
			return err
		}
		defer stmt.Close()

		_, err = stmt.ExecContext(ctx, args...)
		if err != nil {
			return err
		}
	}

	return nil
}

func (dm *AccountDataMapper) Update(ctx context.Context, mCtx unit.MapperContext, accounts ...any) error {
	for _, account := range accounts {
		sql, args, err := dm.table.UpdateQueryWithArgs(account)
		if err != nil {
			return err
		}

		stmt, err := mCtx.Tx.Prepare(sql)
		if err != nil {
			return err
		}
		defer stmt.Close()

		_, err = stmt.ExecContext(ctx, args...)
		if err != nil {
			return err
		}
	}

	return nil
}

func (dm *AccountDataMapper) Delete(ctx context.Context, mCtx unit.MapperContext, accounts ...any) error {
	for _, account := range accounts {
		sql, args, err := dm.table.DeleteQueryWithArgs(account)
		if err != nil {
			return err
		}

		stmt, err := mCtx.Tx.Prepare(sql)
		if err != nil {
			return err
		}
		defer stmt.Close()

		_, err = stmt.ExecContext(ctx, args...)
		if err != nil {
			return err
		}
	}

	return nil
}
