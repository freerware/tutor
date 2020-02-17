package infrastructure

import (
	"database/sql"
	"time"

	"github.com/avast/retry-go"
	"github.com/freerware/tutor/config"
	"github.com/freerware/tutor/domain"
	"github.com/freerware/work"
	"github.com/freerware/workfx"
	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type SQLDataMapperResult struct {
	fx.Out

	SQLDataMappers map[work.TypeName]work.SQLDataMapper
}

type DBResult struct {
	fx.Out

	DB *sql.DB `name:"rwDB"`
}

type SQLDataMapperParameters struct {
	fx.In

	Logger *zap.Logger
}

var Module = fx.Options(
	fx.Provide(NewQueryer),
	fx.Provide(func(c config.Configuration) (DBResult, error) {
		var db *sql.DB
		connect := func() (err error) {
			db, err = sql.Open("mysql", c.Database.DSN())
			return
		}
		options := []retry.Option{
			retry.Attempts(3),
			retry.Delay(time.Second * 5),
		}
		if err := retry.Do(connect, options...); err != nil {
			return DBResult{}, err
		}
		ping := func() (err error) {
			err = db.Ping()
			return
		}
		if err := retry.Do(ping, options...); err != nil {
			return DBResult{}, err
		}
		return DBResult{
			DB: db,
		}, nil
	}),
	fx.Provide(func(parameters SQLDataMapperParameters) SQLDataMapperResult {
		dataMappers := make(map[work.TypeName]work.SQLDataMapper)
		accountTN := work.TypeNameOf(domain.Account{})
		dm := NewAccountDataMapper(AccountDataMapperParameters{
			Logger: parameters.Logger,
		})
		dataMappers[accountTN] = &dm
		result := SQLDataMapperResult{
			SQLDataMappers: dataMappers,
		}
		return result
	}),
	workfx.Modules.SQLUnit,
)
