package validation

import (
	"fmt"
	"net/mail"

	"github.com/insabelter/IWS_GO/models"
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

func ValidateFeedback(feedback models.Feedback) error {
	// Check email address
	if _, err := mail.ParseAddress(feedback.Author.Email); err != nil {
		return &ValidationError{
			Message: "Invalid email address",
		}
	}

	// Create list of ratings to validate all at once
	ratings := []models.Rating{
		feedback.Ratings.Interesting,
		feedback.Ratings.Learning,
		feedback.Ratings.Pacing,
		feedback.Ratings.ExerciseDifficulty,
		feedback.Ratings.Support,
		feedback.Ratings.OverallSatisfaction,
	}

	// Iterate list of ratings
	// If one rating is invalid, returns the error
	for _, rating := range ratings {
		if err := validateRatings(rating); err != nil {
			return err
		}
	}
	return nil
}

func validateRatings(rating models.Rating) error {
	// Check ratings number between 1 and 10
	if rating.Rating < 1 || rating.Rating > 10 {
		return &ValidationError{
			Message: "Rating must be between 1 and 10",
		}
		// Check comment length max 2000 characters
	} else if len(rating.Comment) > 2000 {
		return &ValidationError{
			Message: "Comment can contain a maximum of 2000 characters",
		}
	}
	return nil
}
