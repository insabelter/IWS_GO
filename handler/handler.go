package handler

import (
	//"encoding/json"
	"fmt"
	//"io"
	"net/http"
	//validation "github.com/insabelter/IWS_GO/validation"
	//"github.com/insabelter/IWS_GO/models"
)

// route to test if the server is running -> health check
func MakePingHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "pong")

	}
}

//Schreibe eine  Funktion, welche aus dem http.Request das Rating liest

// Dekodiere JSON-Daten aus dem Anfragekörper
//Versuche den Body des http.Request mithilfe io.ReadAll auszulesen

//Versuchen den JSON Body in ein Rating umzuwandeln.

// Drucke das empfangene Rating

// Antworte mit einer Bestätigung, welche das Rating zurückgibt.
