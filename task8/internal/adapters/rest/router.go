package rest

import (
	"github.com/gorilla/mux"
	"github.com/rendau/my-otus/task8/internal/adapters/rest/hello"
	"net/http"
)

func createRouter(a *API) http.Handler {
	r := mux.NewRouter()

	hello.Router(r.PathPrefix("/hello").Subrouter())

	return middleware(r, a)
}
