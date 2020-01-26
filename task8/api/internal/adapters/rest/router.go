package rest

import (
	"github.com/gorilla/mux"
	"github.com/rendau/my-otus/task8/api/internal/adapters/rest/events"
	"net/http"
)

func createRouter(a *API) http.Handler {
	r := mux.NewRouter()

	events.Router(r.PathPrefix("/events").Subrouter())

	r.HandleFunc("/ready", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	})

	return middleware(r, a)
}
