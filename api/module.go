package api

import (
	"context"

	"github.com/freerware/tutor/api/resources"
	"github.com/freerware/tutor/api/server"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Module = fx.Options(
	fx.Provide(resources.NewAccountResource),
	fx.Provide(server.New),
	fx.Provide(zap.NewDevelopment),
	fx.Invoke(Start),
)

func Start(lc fx.Lifecycle, s server.Server) {

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go s.Start()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return s.Stop(ctx)
		},
	})
}
