package infrastructure

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/avast/retry-go"
	"github.com/cactus/go-statsd-client/statsd"
	"github.com/freerware/tutor/config"
	"github.com/freerware/tutor/infrastructure/models"
	"github.com/freerware/work/v4/unit"
	"github.com/freerware/workfx/v4"
	_ "github.com/go-sql-driver/mysql"
	"github.com/uber-go/tally"
	tstatsd "github.com/uber-go/tally/statsd"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type UnitResult struct {
	fx.Out

	Option unit.Option `group:"unitOptions"`
}

type DBResult struct {
	fx.Out

	DB *sql.DB `name:"rwDB"`
}

type DBParameters struct {
	fx.In

	DB *sql.DB `name:"rwDB"`
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
	fx.Provide(func(c config.Configuration) (tally.Scope, error) {
		addr := fmt.Sprintf("%s:%d", c.Metrics.Host, c.Metrics.Port)
		flushInterval := time.Duration(c.Metrics.MaxFlushInterval) * time.Millisecond
		statter, err :=
			statsd.NewBufferedClient(
				addr, c.Metrics.Prefix, flushInterval, c.Metrics.MaxFlushBytes)
		if err != nil {
			return nil, err
		}
		reporter := tstatsd.NewReporter(statter, tstatsd.Options{
			SampleRate: 1.0,
		})
		scope, _ := tally.NewRootScope(tally.ScopeOptions{
			Tags:     map[string]string{},
			Reporter: reporter,
		}, time.Second)
		return scope, nil
	}),
	fx.Provide(func(l *zap.Logger, c config.Configuration) UnitResult {
		dataMappers := make(map[unit.TypeName]unit.DataMapper)
		accountTN := unit.TypeNameOf(models.Account{})
		dm := NewAccountDataMapper(AccountDataMapperParameters{DSN: c.Database.DSN(), Logger: l})
		dataMappers[accountTN] = &dm
		return UnitResult{Option: unit.DataMappers(dataMappers)}
	}),
	fx.Provide(func(l *zap.Logger) UnitResult {
		return UnitResult{Option: unit.Logger(l)}
	}),
	fx.Provide(func(s tally.Scope) UnitResult {
		return UnitResult{Option: unit.Scope(s)}
	}),
	fx.Provide(func(parameters DBParameters) UnitResult {
		return UnitResult{Option: unit.DB(parameters.DB)}
	}),
	workfx.Module,
)
