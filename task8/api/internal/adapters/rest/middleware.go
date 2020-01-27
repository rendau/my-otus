package rest

import (
	"context"
	"github.com/rendau/my-otus/task8/api/internal/adapters/rest/constants"
	"github.com/rendau/my-otus/task8/api/internal/adapters/rest/entities"
	"github.com/rendau/my-otus/task8/api/internal/adapters/rest/util"
	"github.com/slok/go-http-metrics/metrics/prometheus"
	"github.com/slok/go-http-metrics/middleware"
	"net/http"
)

func (a *API) middleware(h http.Handler, withMetrics bool) http.Handler {
	h = a.mwLog(h)
	h = a.mwAPICtx(h, &entities.APICtx{
		Ucs: a.ucs,
		Log: a.log,
	})
	if withMetrics {
		h = middleware.New(middleware.Config{
			Recorder: prometheus.NewRecorder(prometheus.Config{}),
		}).Handler("", h)
	}
	h = a.mwRecovery(h)
	return h
}

func (a *API) mwLog(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiCtx := util.GetAPICtx(r)
		apiCtx.Log.Infow(
			"HTTP request",
			"method", r.Method,
			"url", r.URL,
		)
		h.ServeHTTP(w, r)
	})
}

func (a *API) mwRecovery(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				apiCtx := util.GetAPICtx(r)
				apiCtx.Log.Errorw(
					"Fail to handle http request",
					"method", r.Method,
					"url", r.URL,
					"error", err,
				)
			}
		}()
		h.ServeHTTP(w, r)
	})
}

func (a *API) mwAPICtx(h http.Handler, ctx *entities.APICtx) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r = r.WithContext(context.WithValue(r.Context(), constants.APICtxKey, ctx))
		h.ServeHTTP(w, r)
	})
}
