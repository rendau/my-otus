package rest

import (
	"context"
	"github.com/rendau/my-otus/task8/api/internal/adapters/rest/constants"
	"github.com/rendau/my-otus/task8/api/internal/adapters/rest/entities"
	"github.com/rendau/my-otus/task8/api/internal/adapters/rest/util"
	"net/http"
)

func middleware(h http.Handler, a *API) http.Handler {
	h = mwLog(h)
	h = mwAPICtx(h, &entities.APICtx{
		Ucs: a.ucs,
		Log: a.log,
	})
	h = mwRecovery(h)
	return h
}

func mwLog(h http.Handler) http.Handler {
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

func mwRecovery(h http.Handler) http.Handler {
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

func mwAPICtx(h http.Handler, ctx *entities.APICtx) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r = r.WithContext(context.WithValue(r.Context(), constants.APICtxKey, ctx))
		h.ServeHTTP(w, r)
	})
}
