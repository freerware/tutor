package infrastructure

import (
	"context"
	"database/sql"

	"github.com/freerware/morph"
)

type DataMapper[T any] struct {
	table morph.Table
	tx    *sql.Tx
}

func NewDataMapper[T any](tx *sql.Tx, options ...morph.ReflectOption) (DataMapper[T], error) {
	t, err := morph.Reflect(new(T), options...)
	if err != nil {
		return DataMapper[T]{}, err
	}
	return DataMapper[T]{table: t, tx: tx}, nil
}

func (dm DataMapper[T]) Table() morph.Table {
	return dm.table
}

func (dm DataMapper[T]) Insert(ctx context.Context, entities ...T) error {
	for _, entity := range entities {
		sql, args, err := dm.table.InsertQueryWithArgs(entity)
		if err != nil {
			return err
		}

		_, err = dm.tx.ExecContext(ctx, sql, args...)
		if err != nil {
			return err
		}
	}

	return nil
}

func (dm DataMapper[T]) Update(ctx context.Context, entities ...T) error {
	for _, entity := range entities {
		sql, args, err := dm.table.UpdateQueryWithArgs(entity)
		if err != nil {
			return err
		}

		_, err = dm.tx.ExecContext(ctx, sql, args...)
		if err != nil {
			return err
		}
	}

	return nil
}

func (dm DataMapper[T]) Delete(ctx context.Context, entities ...T) error {
	for _, entity := range entities {
		sql, args, err := dm.table.DeleteQueryWithArgs(entity)
		if err != nil {
			return err
		}

		_, err = dm.tx.ExecContext(ctx, sql, args...)
		if err != nil {
			return err
		}
	}

	return nil
}
