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

	// Read Connection String from connection.json
	var settings, jsonErr = readSettings()

	if jsonErr != nil {
		fmt.Println("Error while reading configuration:", jsonErr.Error)
	} else {
		// Connecting to ATLAS cluster
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()
		client, err := mongo.Connect(ctx, options.Client().ApplyURI(
			settings.DbConnectionString,
		))

		// Disconnect from the client when the program ends
		defer func() {
			if err = client.Disconnect(ctx); err != nil {
				panic(err)
			}
		}()

		// Check the connection
		err = client.Ping(ctx, nil)
		if err != nil {
			fmt.Println("There was a problem connecting to your Atlas cluster.")
			panic(err)
		}
		fmt.Println("Connected to MongoDB!")

		// Create a MongoDB repository
		repository := repository.NewMongoDBRepository(client.Database("FeedbackDB"))

		// Start the server
		fmt.Println("Server started on port 3000...")
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
