package infrastructure

import (
	"context"
	"errors"
	"log"
	"os"
	"time"

	"github.com/freerware/work/v4/unit"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Errors that are potentially thrown during data mapper interactions.
var (

	// ErrInvalidType represents an error that indicates a unexpected type
	// was provided to the data mapper.
	ErrInvalidType = errors.New("infrastructure: invalid type provided to data mapper")
)

type AccountDataMapperParameters struct {
	fx.In

	DSN    string
	Logger *zap.Logger
}

type AccountDataMapper struct {
	dsn    string
	logger *zap.Logger
}

func NewAccountDataMapper(
	parameters AccountDataMapperParameters) AccountDataMapper {
	return AccountDataMapper{logger: parameters.Logger, dsn: parameters.DSN}
}

func (dm *AccountDataMapper) db(mCtx unit.MapperContext) (*gorm.DB, error) {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // Slow SQL threshold
			LogLevel:      logger.Info, // Log level
			Colorful:      true,        // Disable color
		},
	)
	return gorm.Open(mysql.Open(dm.dsn), &gorm.Config{ConnPool: mCtx.Tx, Logger: newLogger})
}

func (dm *AccountDataMapper) Insert(ctx context.Context, mCtx unit.MapperContext, accounts ...interface{}) error {
	if len(accounts) == 0 {
		return nil
	}
	db, err := dm.db(mCtx)
	if err != nil {
		return err
	}
	for _, account := range accounts {
		if err = db.Create(account).Error; err != nil {
			return err
		}
	}
	return nil
}

func (dm *AccountDataMapper) Update(ctx context.Context, mCtx unit.MapperContext, accounts ...interface{}) error {
	if len(accounts) == 0 {
		return nil
	}
	db, err := dm.db(mCtx)
	if err != nil {
		return err
	}
	for _, account := range accounts {
		if err = db.Save(account).Error; err != nil {
			return err
		}
	}
	return nil
}

func (dm *AccountDataMapper) Delete(ctx context.Context, mCtx unit.MapperContext, accounts ...interface{}) error {
	if len(accounts) == 0 {
		return nil
	}
	db, err := dm.db(mCtx)
	if err != nil {
		return err
	}
	for _, account := range accounts {
		if err = db.Delete(account).Error; err != nil {
			return err
		}
	}
	return nil
}
