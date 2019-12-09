package rest

import (
	"context"
	"github.com/rendau/my-otus/task8/api/internal/domain/usecases"
	"github.com/rendau/my-otus/task8/api/internal/interfaces"
	"net/http"
	"time"
)

// API - is type for rest API adapter
type API struct {
	log    interfaces.Logger
	lAddr  string
	server *http.Server
	ucs    *usecases.Usecases
}

// CreateAPI - creates new instance
func CreateAPI(log interfaces.Logger, lAddr string, ucs *usecases.Usecases) *API {
	return &API{
		log:   log,
		lAddr: lAddr,
		ucs:   ucs,
	}
}

// Start - starts api-server
func (a *API) Start() {
	a.server = &http.Server{
		Addr:         a.lAddr,
		Handler:      createRouter(a),
		ReadTimeout:  20 * time.Second,
		WriteTimeout: 20 * time.Second,
	}
	go func() {
		err := a.server.ListenAndServe()
		if err != http.ErrServerClosed {
			a.log.Fatalw("Fail to start http server", "error", err)
		}
	}()
}

// Shutdown - stops api-server
func (a *API) Shutdown(ctx context.Context) error {
	return a.server.Shutdown(ctx)
}
