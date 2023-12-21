package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/insabelter/IWS_GO/middleware"
	"github.com/insabelter/IWS_GO/models"
	"github.com/insabelter/IWS_GO/repository"
)

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func MakeGetFeedbacksHandler(ctx context.Context, repository repository.Repository) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)

		if feedbacks, err := repository.GetAllFeedbacks(ctx); err == nil {
			if json, err := json.Marshal(feedbacks); err == nil {
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprintf(w, string(json))
			} else {
				fmt.Println(err)
				w.WriteHeader(http.StatusBadRequest)
			}
		} else {
			fmt.Println(err)
			w.WriteHeader(http.StatusBadRequest)
		}
	}
}

func MakeGetFeedbackHandler(ctx context.Context, repository repository.Repository) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)

		id := mux.Vars(r)["id"]
		if feedback, err := repository.GetFeedback(ctx, id); err == nil {
			if json, err := json.Marshal(feedback); err == nil {
				fmt.Fprintf(w, string(json))
			} else {
				fmt.Println(err)
				w.WriteHeader(http.StatusBadRequest)
			}
		} else {
			fmt.Println(err)
			w.WriteHeader(http.StatusNotFound)
		}
	}
}

func MakeAddFeedbackHandler(ctx context.Context, repository repository.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)

		body, err := io.ReadAll(r.Body)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var feedback models.Feedback
		if err = json.Unmarshal(body, &feedback); err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		feedback.ID = uuid.New().String()

		if err = middleware.ValidateFeedback(feedback); err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if createdFeedback, err := repository.CreateFeedback(ctx, feedback); err == nil {
			if json, err := json.Marshal(createdFeedback); err == nil {
				w.WriteHeader(http.StatusCreated)
				fmt.Fprintf(w, string(json))
			} else {
				fmt.Println(err)
				w.WriteHeader(http.StatusBadRequest)
			}
		} else {
			fmt.Println(err)
			w.WriteHeader(http.StatusBadRequest)
		}

	}
}

func MakeDeleteFeedbackHandler(ctx context.Context, repository repository.Repository) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)

		id := mux.Vars(r)["id"]
		if repository.DeleteFeedback(ctx, id) == nil {
			w.WriteHeader(http.StatusNoContent)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}
}
