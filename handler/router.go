package handler

import (
	"context"

	"github.com/gorilla/mux"
	"github.com/insabelter/IWS_GO/repository"
)

func NewRouter(ctx context.Context, repository repository.Repository) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/ping", MakePingHandler(ctx, repository)).Methods("GET")
	//CRUD
	r.HandleFunc("/feedback", MakeGetFeedbacksHandler(ctx, repository)).Methods("GET")
	r.HandleFunc("/feedback", MakeAddFeedbackHandler(ctx, repository)).Methods("POST")
	r.HandleFunc("/feedback/{id}", MakeGetFeedbackHandler(ctx, repository)).Methods("GET")
	r.HandleFunc("/feedback/{id}", MakeDeleteFeedbackHandler(ctx, repository)).Methods("DELETE")

	//Analytics
	r.HandleFunc("/analytics/average/overallsatisfaction", MakeAverageOverallSatisfactionHandler(ctx, repository)).Methods("GET")
	r.HandleFunc("/analytics/average/support", MakeAverageSupportHandler(ctx, repository)).Methods("GET")

	return r
}
