package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/insabelter/IWS_GO/middleware"
	"github.com/insabelter/IWS_GO/models"
	"github.com/insabelter/IWS_GO/repository"
)

func MakeGetFeedbacksHandler(ctx context.Context, repository repository.Repository) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
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
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, err.Error())
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
		id := mux.Vars(r)["id"]
		if repository.DeleteFeedback(ctx, id) == nil {
			w.WriteHeader(http.StatusNoContent)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}
}

// this route uses goroutines and channels
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
				cOverallSatisfaction <- feedback.Ratings.OverallSatisfaction.Rating
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

// this route uses goroutines, shared memory and WaitGroups
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

				sum.Value += feedback.Ratings.Support.Rating

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

func MakePingHandler(ctx context.Context, repository repository.Repository) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "pong")
	}
}
