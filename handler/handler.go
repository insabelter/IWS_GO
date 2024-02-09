package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/insabelter/IWS_GO/models"
	"github.com/insabelter/IWS_GO/middelware"
)

// route to test if the server is running -> health check
func MakePingHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "pong")

	}
}
func MakeRatingHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		// Dekodiere JSON-Daten aus dem Anfragekörper
		var receivedRating models.Rating

		body, err := io.ReadAll(r.Body)
		if err != nil {
			// error response if the request body can't be read
			fmt.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		
		err = json.Unmarshal(body, &receivedRating)
		if err != nil {
			http.Error(w, "Failed to decode JSON", http.StatusBadRequest)
			return
		}

		if err = middleware.validateRating(receivedRating); err != nil {
			// error response if the validation fails
			w.WriteHeader(http.StatusBadRequest)
			// the custom validation error is added to the response
			fmt.Fprintf(w, err.Error())
			return
		}



		// Drucke das empfangene Rating
		fmt.Printf("Received Rating: %+v\n", receivedRating)

		// Antworte mit einer Bestätigung
		response := fmt.Sprintf("Rating received successfully. Rating: %+v", receivedRating)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response))
	}
}
