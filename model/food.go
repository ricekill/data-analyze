package model

type FoodData struct {
	Name       string `json:"name"`
	Score      string `json:"score"`
	Evaluation string `json:"evaluation"`
	FoodType   string `json:"food_type"`
	Area       string `json:"area"`
	Place      string `json:"place"`
	Banner     string `json:"banner"`
	CreatedAt  int64  `json:"created_at"`
}
