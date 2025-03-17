package resources

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/freerware/negotiator"
	"github.com/freerware/negotiator/proactive"
	"github.com/freerware/negotiator/representation"
	j "github.com/freerware/tutor/api/representations/json"
	p "github.com/freerware/tutor/api/representations/protobuf"
	x "github.com/freerware/tutor/api/representations/xml"
	y "github.com/freerware/tutor/api/representations/yaml"
	"github.com/freerware/tutor/api/server"
	app "github.com/freerware/tutor/application"
	"github.com/freerware/tutor/domain"
	u "github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type AccountResourceResult struct {
	fx.Out

	AccountResource  AccountResource
	MuxConfiguration server.MuxConfiguration
}

type AccountResourceParameters struct {
	fx.In

	AccountService app.AccountService
	Logger         *zap.Logger
}

type AccountResource struct {
	accountService app.AccountService
	logger         *zap.Logger
}

func NewAccountResource(
	parameters AccountResourceParameters,
) AccountResourceResult {
	a := AccountResource{
		accountService: parameters.AccountService,
		logger:         parameters.Logger,
	}
	return AccountResourceResult{
		AccountResource:  a,
		MuxConfiguration: a.MuxConfiguration(),
	}
}

func (ar *AccountResource) Get(w http.ResponseWriter, request *http.Request) {

	// retrieve the account uuid.
	vars := mux.Vars(request)
	uuid, err := u.FromString(vars["uuid"])
	if err != nil {
		http.Error(w, err.Error(), 400)
	}

	// retrieve the account.
	account, err := ar.accountService.Get(uuid)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	jacc := j.NewAccount(account)
	jacc.SetContentLocation(*request.URL)
	gjacc := j.NewAccount(account)
	gjacc.SetContentLocation(*request.URL)
	gjacc.SetContentEncoding([]string{"gzip"})
	yacc := y.NewAccount(account)
	yacc.SetContentLocation(*request.URL)
	xacc := x.NewAccount(account)
	xacc.SetContentLocation(*request.URL)
	pacc := p.NewAccount(account)
	pacc.SetContentLocation(*request.URL)
	representations := []representation.Representation{jacc, yacc, xacc, gjacc, pacc}

	// negotiate.
	ctx := negotiator.NegotiationContext{Request: request, ResponseWriter: w}
	if err = proactive.Default.Negotiate(ctx, representations...); err != nil {
		http.Error(w, err.Error(), 500)
	}
}

func (ar *AccountResource) CreateAndAppend(
	w http.ResponseWriter, request *http.Request) {

	// clean this part up.
	representation := j.Account{}
	if err := json.NewDecoder(request.Body).Decode(&representation); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	account, err := domain.NewAccount(domain.AccountParameters{
		UUID:      u.Must(u.NewV4()),
		GivenName: representation.GivenName,
		Surname:   representation.Surname,
		Username:  representation.PrimaryCredential,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		DeletedAt: nil,
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// create the account.
	err = ar.accountService.Create(request.Context(), account)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	uri, _ := request.URL.Parse("/" + account.UUID().String())
	w.Header().Add("Content-Location", uri.String())
	w.WriteHeader(201)
}

func (ar *AccountResource) Replace(w http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)

	// clean this part up.
	representation := j.Account{}
	if err := json.NewDecoder(request.Body).Decode(&representation); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	account, err := domain.NewAccount(domain.AccountParameters{
		UUID:      representation.UUID,
		GivenName: representation.GivenName,
		Surname:   representation.Surname,
		Username:  representation.PrimaryCredential,
		CreatedAt: representation.CreatedAt,
		UpdatedAt: time.Now(),
		DeletedAt: nil,
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	uuid := u.Must(u.FromString(vars["uuid"]))
	if account.UUID() != uuid {
		http.Error(w, fmt.Errorf("mismatching UUIDs").Error(), 400)
		return
	}

	// retrieve the account.
	existing, err := ar.accountService.Get(uuid)
	if err != nil {
		http.Error(w, err.Error(), 404)
		return
	}

	// upsert the account.
	account.SetCreatedAt(existing.CreatedAt())
	err = ar.accountService.Put(request.Context(), account)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.WriteHeader(204)
}

func (ar *AccountResource) Delete(w http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)

	// retrieve the account uuid.
	uuid := u.Must(u.FromString(vars["uuid"]))
	account, err := ar.accountService.Get(uuid)
	if err != nil {
		http.Error(w, err.Error(), 404)
		return
	}

	// delete the account.
	err = ar.accountService.Delete(request.Context(), account)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.WriteHeader(204)
}
