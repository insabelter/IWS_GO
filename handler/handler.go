package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/insabelter/IWS_GO/repository"
)

// route to get all feedbacks as a list
func MakeGetFeedbacksHandler(ctx context.Context, repository repository.Repository) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// get all feedbacks from the repository
		if feedbacks, err := repository.GetAllFeedbacks(ctx); err == nil {
			// transform the feedbacks into json
			if json, err := json.Marshal(feedbacks); err == nil {
				// successfully send the json response
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprintf(w, string(json))
			} else {
				// error response if the json transformation fails
				fmt.Println(err)
				w.WriteHeader(http.StatusBadRequest)
			}
		} else {
			// error response if the repository returns an error
			fmt.Println(err)
			w.WriteHeader(http.StatusBadRequest)
		}
	}
}

// 1. TODO: Implement the get feedback by id handler

// route to get a feedback based on its id
func MakeGetFeedbackHandler(ctx context.Context, repository repository.Repository) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// get the id from the url

		// get the feedback from the repository (+ check for error)

		// transform the feedback into json (+ check for error)

		// successful: send the feedback as json response

		// repository or json error: send an error response

	}
}

// 2. TODO: Implement the add feedback handler

// route to add a new feedback
// uses the validation middleware to validate the new feedback
func MakeAddFeedbackHandler(ctx context.Context, repository repository.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// read the request body (+ check for error)

		// error response if the request body can't be read

		// create a variable for the new feedback

		// transform the request body into a feedback struct and save it to the variable (+ check for error)

		// error response if the json tranformation fails

		// validate the feedback using the ValidateFeedback function from the middleware (+ check for error)

		// error response if the validation fails (custom error message should be added to the response)

		// generate a new uuid for the feedback (save it to feedback.ID)

		// add the new feedback to the repository (+ check for error)

		// transform the created feedback repository response into json (+ check for error)

		// successful: send the feedback as json response

		// repository or json error: send an error response

	}
}

// route to delete a feedback based on its id
func MakeDeleteFeedbackHandler(ctx context.Context, repository repository.Repository) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// get the id from the url
		id := mux.Vars(r)["id"]

		// delete the feedback from the repository
		if repository.DeleteFeedback(ctx, id) == nil {
			// successful response for deleting the feedback
			w.WriteHeader(http.StatusNoContent)
		} else {
			// error response if the repository can not find a feedback with the given id
			w.WriteHeader(http.StatusNotFound)
		}
	}
}

// route to test if the server is running -> health check
func MakePingHandler(ctx context.Context, repository repository.Repository) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "pong")
	}
}
