package models

type Rating struct {
	Rating  int    `json:"rating"`
	Comment string `json:"comment"`
}

type Ratings struct {
	Interesting         Rating `json:"interesting"`
	Learning            Rating `json:"learning"`
	Pacing              Rating `json:"pacing"`
	ExerciseDifficulty  Rating `json:"exercise_difficulty"`
	Support             Rating `json:"support"`
	OverallSatisfaction Rating `json:"overall_satisfaction"`
}

type Author struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Feedback struct {
	ID      string  `json:"id"`
	Author  Author  `json:"author"`
	Ratings Ratings `json:"ratings"`
}
