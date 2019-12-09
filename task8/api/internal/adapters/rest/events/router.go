package events

import (
	"github.com/gorilla/mux"
)

// Router - router for hello sub-route
func Router(r *mux.Router) {
	r.HandleFunc("", hCreate).Methods("POST")
	r.HandleFunc("/{id:[0-9]+}", hUpdate).Methods("PUT")
	r.HandleFunc("/{id:[0-9]+}", hDelete).Methods("DELETE")
	r.HandleFunc("/for_day", hListForDay).Methods("GET")
	r.HandleFunc("/for_week", hListForWeek).Methods("GET")
	r.HandleFunc("/for_month", hListForMonth).Methods("GET")
}
