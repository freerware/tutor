package resources

import (
	"github.com/freerware/tutor/api/server"
)

func (ar *AccountResource) MuxConfiguration() (config server.MuxConfiguration) {
	config = server.MuxConfiguration{
		PathPrefix: "/accounts",
		Handlers: []server.HandlerConfiguration{
			{
				Path:        "/{uuid}/",
				HandlerFunc: ar.Get,
				Methods:     []string{"GET"},
			},
			{
				Path:        "/{uuid}",
				HandlerFunc: ar.Get,
				Methods:     []string{"GET"},
			},
			{
				Path:        "/{uuid}",
				HandlerFunc: ar.Replace,
				Methods:     []string{"PUT"},
			},
			{
				Path:        "/{uuid}/",
				HandlerFunc: ar.Replace,
				Methods:     []string{"PUT"},
			},
			{
				Path:        "/{uuid}",
				HandlerFunc: ar.Delete,
				Methods:     []string{"DELETE"},
			},
			{
				Path:        "/{uuid}/",
				HandlerFunc: ar.Delete,
				Methods:     []string{"DELETE"},
			},
			{
				Path:        "",
				HandlerFunc: ar.CreateAndAppend,
				Methods:     []string{"POST"},
			},
			{
				Path:        "/",
				HandlerFunc: ar.CreateAndAppend,
				Methods:     []string{"POST"},
			},
		},
	}
	return
}
