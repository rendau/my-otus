package hello

import (
	"github.com/gorilla/mux"
)

// Router - router for hello sub-route
func Router(r *mux.Router) {
	r.HandleFunc("", hGet).Methods("GET")
}
