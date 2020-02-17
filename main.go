package main

import (
	"time"
	"github.com/freerware/tutor/api"
	"github.com/freerware/tutor/application"
	"github.com/freerware/tutor/config"
	"github.com/freerware/tutor/infrastructure"

	"go.uber.org/fx"
)

func main() {
	fx.New(
		fx.StartTimeout(time.Second*45),
		api.Module,
		config.Module,
		application.Module,
		infrastructure.Module,
	).Run()
}
