package handler

import (
	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
	r := mux.NewRouter()

	// Health check
	r.HandleFunc("/ping", MakePingHandler()).Methods("GET")
	r.HandleFunc("/postRating", MakeRatingHandler()).Methods("POST")
	return r
}
