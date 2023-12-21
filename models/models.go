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
