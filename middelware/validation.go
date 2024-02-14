package middleware

import (
	"fmt"
	//"github.com/insabelter/IWS_GO/models"
)

// custom error type for validation errors
type ValidationError struct {
	Message string
}

// implement the error interface for the custom validation error type
// returns the custom error message and marks it as a validation error
func (m *ValidationError) Error() string {
	return fmt.Sprintf("Validation error: %s", m.Message)
}

//Schreibe eine Funktion die einen ValidationError zurückgibt, wenn das Rating nicht zwischen 1-10 liegt oder der Kommentar mehr als 2000 Zeichen hat sonst gebe nil zurück
//Binde danach die Funktion im Handler, welcher das Rating empfängt ein.
