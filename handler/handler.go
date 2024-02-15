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

// Schreibe hier den Make FeedbackHandler


		//Mithilfe io.ReadAll() kannst du den Body speichern
		//Was soll passieren wenn dies nicht möglich ist? 
		//Schaut euch dazu die Rückgabewerte der Methode io.ReadAll()
		
		// Antworte mit einer Bestätigung. Diese soll das Feedback zurück geben
		// Der body ist von Typ byte[]. string() kann helfen
		
		// Für die Schnellen mithilfe von r.FromValue("key") kannst du HTTP Parameter ausgeben.
		//Übergebe die Gesamtbewertung und gebe diese aus.
