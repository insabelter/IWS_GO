package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/insabelter/IWS_GO/handler"
	"github.com/insabelter/IWS_GO/repository"
)

type Settings struct {
	DbConnectionString string `json:"connectionString"`
}

func main() {

	// var feedback = models.Feedback{
	// 	ID: "1",
	// 	IP: "test",
	// 	Author: models.Author{
	// 		ID:    "1",
	// 		Email: "test@test.de",
	// 	},
	// 	Ratings: models.Ratings{
	// 		Interesting: models.Rating{
	// 			ID:      "1",
	// 			Rating:  5,
	// 			Comment: "test",
	// 		},
	// 		Learning: models.Rating{
	// 			ID:      "2",
	// 			Rating:  5,
	// 			Comment: "test",
	// 		},
	// 		Pacing: models.Rating{
	// 			ID:      "3",
	// 			Rating:  5,
	// 			Comment: "test",
	// 		},
	// 		ExerciseDifficulty: models.Rating{
	// 			ID:      "4",
	// 			Rating:  5,
	// 			Comment: "test",
	// 		},
	// 		Support: models.Rating{
	// 			ID:      "5",
	// 			Rating:  5,
	// 			Comment: "test",
	// 		},
	// 		OverallSatisfaction: models.Rating{
	// 			ID:      "6",
	// 			Rating:  5,
	// 			Comment: "test",
	// 		},
	// 	},
	// }

	// Read Connection String from connection.json
	var settings, jsonErr = readSettings()

	if jsonErr != nil {
		fmt.Println("Error while reading configuration:", jsonErr.Error)
	} else {
		fmt.Println("connection-string:" + settings.DbConnectionString)

		// Connecting to ATLAS cluster
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()
		client, err := mongo.Connect(ctx, options.Client().ApplyURI(
			settings.DbConnectionString,
		))

		defer func() {
			if err = client.Disconnect(ctx); err != nil {
				panic(err)
			}
		}()

		err = client.Ping(ctx, nil)

		if err != nil {
			fmt.Println("There was a problem connecting to your Atlas cluster. Check that the URI includes a valid username and password, and that your IP address has been added to the access list. Error: ")
			panic(err)
		}

		fmt.Println("Connected to MongoDB!\n")

		// create a repository
		repository := repository.NewRepository(client.Database("FeedbackDB"))

		fmt.Printf("Server started on port 3000...\n")
		http.ListenAndServe(":3000", handler.NewRouter(ctx, repository))
	}
}

func readSettings() (settings *Settings, err error) {
	var fileName = "connection.json"
	file, _ := os.Open(fileName)
	settings = new(Settings)
	err = json.NewDecoder(file).Decode(settings)
	return
}
