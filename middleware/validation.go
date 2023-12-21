package middleware

import (
	"net/mail"

	"github.com/insabelter/IWS_GO/models"
)

type ValidationError struct{}

func (m *ValidationError) Error() string {
	return "The feedback contains invalid data"
}

func ValidateFeedback(feedback models.Feedback) error {
	// Check email address
	if _, err := mail.ParseAddress(feedback.Author.Email); err != nil {
		return &ValidationError{}
	}

	// Check all ratings
	if err := checkRating(feedback.Ratings.Interesting); err != nil {
		return err
	} else if err := checkRating(feedback.Ratings.Learning); err != nil {
		return err
	} else if err := checkRating(feedback.Ratings.Pacing); err != nil {
		return err
	} else if err := checkRating(feedback.Ratings.ExerciseDifficulty); err != nil {
		return err
	} else if err := checkRating(feedback.Ratings.Support); err != nil {
		return err
	} else if err := checkRating(feedback.Ratings.OverallSatisfaction); err != nil {
		return err
	}
	return nil
}

func checkRating(rating models.Rating) error {
	// Check ratings number between 1 and 10
	// And Check comment length max 2000 characters
	if rating.Rating < 1 || rating.Rating > 10 || len(rating.Comment) > 2000 {
		return &ValidationError{}
	}
	return nil
}
