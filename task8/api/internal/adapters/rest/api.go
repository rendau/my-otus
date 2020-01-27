package rest

import (
	"context"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rendau/my-otus/task8/api/internal/domain/usecases"
	"github.com/rendau/my-otus/task8/api/internal/interfaces"
	"net/http"
	"time"
)

// API - is type for rest API adapter
type API struct {
	log     interfaces.Logger
	lAddr   string
	mlAddr  string
	server  *http.Server
	mServer *http.Server
	ucs     *usecases.Usecases
}

// CreateAPI - creates new instance
func CreateAPI(log interfaces.Logger, lAddr, mlAddr string, ucs *usecases.Usecases) *API {
	return &API{
		log:    log,
		lAddr:  lAddr,
		mlAddr: mlAddr,
		ucs:    ucs,
	}
}

// Start - starts api-server
func (a *API) Start() {
	a.server = &http.Server{
		Addr:         a.lAddr,
		Handler:      a.createRouter(a.mlAddr != ""),
		ReadTimeout:  20 * time.Second,
		WriteTimeout: 20 * time.Second,
	}
	go func() {
		err := a.server.ListenAndServe()
		if err != http.ErrServerClosed {
			a.log.Fatalw("Fail to start http server", "error", err)
		}
	}()

	if a.mlAddr != "" {
		a.mServer = &http.Server{
			Addr:    a.mlAddr,
			Handler: promhttp.Handler(),
		}
		go func() {
			err := a.mServer.ListenAndServe()
			if err != http.ErrServerClosed {
				a.log.Fatalw("Fail to start metrics http server", "error", err)
			}
		}()
	}
}

// Shutdown - stops api-server
func (a *API) Shutdown(ctx context.Context) error {
	err := a.server.Shutdown(ctx)
	if err != nil {
		return err
	}
	if a.mServer != nil {
		err = a.mServer.Shutdown(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}
