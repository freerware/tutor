package infrastructure

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/avast/retry-go"
	"github.com/freerware/tutor/config"
	"github.com/freerware/tutor/domain"
	"github.com/freerware/work/v4/unit"
	"github.com/freerware/workfx/v4"
	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type UnitResult struct {
	fx.Out

	Options []unit.Option `group:"unitOptions"`
}

type DBResult struct {
	fx.Out

	DB *sql.DB `name:"rwDB"`
}

type UnitParameters struct {
	fx.In

	Logger *zap.Logger
	DB     *sql.DB `name:"rwDB"`
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
			retry.OnRetry(func(attempt uint, retryErr error) {
				fmt.Println("Attempting to connect:", c.Database.DSN())
			}),
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
		return DBResult{DB: db}, nil
	}),
	fx.Provide(func(parameters UnitParameters) UnitResult {
		dataMappers := make(map[unit.TypeName]unit.DataMapper)
		accountTN := unit.TypeNameOf(domain.Account{})
		dm := NewAccountDataMapper(AccountDataMapperParameters{
			Logger: parameters.Logger,
		})
		dataMappers[accountTN] = &dm
		result := UnitResult{Options: []unit.Option{
			unit.DataMappers(dataMappers),
			unit.Logger(parameters.Logger),
			unit.DB(parameters.DB),
		}}
		return result
	}),
	workfx.Module,
)
