package main

import (
	"fmt"
	"net/http"

	"github.com/insabelter/IWS_GO/handler"
)

func main() {

	// Read Connection String from connection.json

	fmt.Println("Server started on port 3000...\n")
	http.ListenAndServe(":3000", handler.NewRouter())
}
