package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/insabelter/IWS_GO/middleware"
	"github.com/insabelter/IWS_GO/models"
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

// route to get a feedback based on its id
func MakeGetFeedbackHandler(ctx context.Context, repository repository.Repository) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// get the id from the url
		id := mux.Vars(r)["id"]

		// get the feedback from the repository
		if feedback, err := repository.GetFeedback(ctx, id); err == nil {
			// transform the feedback into json
			if json, err := json.Marshal(feedback); err == nil {
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
			w.WriteHeader(http.StatusNotFound)
		}
	}
}

// route to add a new feedback
// uses the validation middleware to validate the new feedback
func MakeAddFeedbackHandler(ctx context.Context, repository repository.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// read the request body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			// error response if the request body can't be read
			fmt.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var feedback models.Feedback
		// transform the request body into a feedback struct
		if err = json.Unmarshal(body, &feedback); err != nil {
			// error response if the tranformation fails
			fmt.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// validate the feedback
		if err = middleware.ValidateFeedback(feedback); err != nil {
			// error response if the validation fails
			w.WriteHeader(http.StatusBadRequest)
			// the custom validation error is added to the response
			fmt.Fprintf(w, err.Error())
			return
		}

		// generate a new uuid for the feedback
		feedback.ID = uuid.New().String()

		// add the feedback to the repository
		if createdFeedback, err := repository.CreateFeedback(ctx, feedback); err == nil {
			// transform the created feedback into json
			if json, err := json.Marshal(createdFeedback); err == nil {
				// successfully send the new feedback as json response
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusCreated)
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

// route to calculate the average of all overall satisfaction feedback scores
// uses goroutines and channels
func MakeAverageOverallSatisfactionHandler(ctx context.Context, repository repository.Repository) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		//get all feedback documents
		feedbacks, err := repository.GetAllFeedbacks(ctx)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		//initalize channel (unbuffered)
		cOverallSatisfaction := make(chan int)

		//start a goroutine for each feedback document
		for _, feedback := range feedbacks {
			go func(feedback models.Feedback) {
				//send the overall satisfaction score to the channel
				cOverallSatisfaction <- readOverallSatisfactionRating(feedback)
			}(feedback)
		}

		sum := 0

		//receive the overall satisfaction feedback scores for each feedback document and add them up
		for i := 0; i < len(feedbacks); i++ {
			overallSatisfactionScore := <-cOverallSatisfaction
			sum += overallSatisfactionScore
		}

		//calculate the average
		average := float64(sum) / float64(len(feedbacks))

		//construct the response
		response := struct {
			AverageOverallSatisfaction float64 `json:"averageOverallSatisfaction"`
		}{average}

		//transform the response into json
		json, err := json.Marshal(response)

		//error response if the json transformation fails
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		//send the response
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, string(json))

	}
}

func readOverallSatisfactionRating(feedback models.Feedback) int {
	time.Sleep(time.Second)
	fmt.Println("reading overall satisfaction score...")
	return feedback.Ratings.OverallSatisfaction.Rating
}

// route to calculate the average of all support feedback scores
// uses goroutines, shared memory and WaitGroups
func MakeAverageSupportHandler(ctx context.Context, repository repository.Repository) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		//get all feedback documents
		feedbacks, err := repository.GetAllFeedbacks(ctx)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		type Sum struct {
			Value int
			Mutex sync.Mutex
		}

		sum := Sum{
			Value: 0,
			Mutex: sync.Mutex{},
		}

		wg := sync.WaitGroup{}

		//register the number of goroutines to wait for (one goroutine per feedback document)
		wg.Add(len(feedbacks))

		//start a goroutine for each feedback document
		for _, feedback := range feedbacks {
			go func(sum *Sum, wg *sync.WaitGroup, feedback models.Feedback) {
				//register the goroutine as finished when it's done
				defer wg.Done()

				//fmt.Println("locking...")
				sum.Mutex.Lock()
				//fmt.Println("locked")

				sum.Value += readSupportRating(feedback)

				//uncomment this line and the prints to showcase that the mutex prevents multiple threads from accessing the shared memory at the same time
				//time.Sleep(time.Second)

				//fmt.Println("unlocking...")
				sum.Mutex.Unlock()
				//fmt.Println("unlocked")
			}(&sum, &wg, feedback)

		}
		//wait for all goroutines to finish
		wg.Wait()

		//calculate the average
		average := float64(sum.Value) / float64(len(feedbacks))

		//construct the response
		response := struct {
			AverageSupportRating float64 `json:"averageSupportRating"`
		}{average}

		//transform the response into json
		json, err := json.Marshal(response)

		//error response if the json transformation fails
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		//send the response
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, string(json))

	}
}

func readSupportRating(feedback models.Feedback) int {
	time.Sleep(time.Second)
	fmt.Println("reading support score...")
	return feedback.Ratings.Support.Rating
}

// route to test if the server is running -> health check
func MakePingHandler(ctx context.Context, repository repository.Repository) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "pong")
	}
}
