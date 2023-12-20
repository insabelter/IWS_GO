package models

type Rating struct {
	ID      string
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
	ID    string
	Email string
}

type Feedback struct {
	ID      string
	IP      string
	Author  Author
	Ratings Ratings
}
