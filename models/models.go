package models

type Rating struct {
	Rating  int
	Comment string
}

type Ratings struct {
	Interesting         Rating
	Learning            Rating
	Pacing              Rating
	ExerciseDifficulty  Rating
	Support             Rating
	OverallSatisfaction Rating
}

type Author struct {
	Name  string
	Email string
}

type Feedback struct {
	ID      string
	Author  Author
	Ratings Ratings
}

func validateRating(rating Rating) string {
	// Check ratings number between 1 and 10
	if rating.Rating < 1 || rating.Rating > 10 {
		return "Rating must be between 1 and 10"
		// Check comment length max 2000 characters
	} else if len(rating.Comment) > 2000 {
		return "Comment can contain a maximum of 2000 characters"
	}
	return "Valid"
}

func NewRating(rating int, comment string) (result Rating, message string) {
	if validateRating(Rating{rating, comment}) == "Valid" {
		message = "Erstellung erfolgreich!"
		result = Rating{rating, comment}
	} else {
		message = "Bei der Erstellung ist ein Fehler Aufgetreten"
	}
	return
}

// Verbesserte Version der Funktion mit Pointer, hier wird bei einem Fehler <nil> zurückgegeben
func NewRatingPointer(rating int, comment string) (result *Rating, message string) {
	if validateRating(Rating{rating, comment}) == "Valid" {
		message = "Erstellung erfolgreich!"
		result = &Rating{rating, comment}
	} else {
		message = "Bei der Erstellung ist ein Fehler Aufgetreten"
	}
	return
}

func (r *Rating) ChangeRating(rating int, comment string) {
	r.Rating = rating
	r.Comment = comment
}
