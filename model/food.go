package model

type Deliveroo struct {
	Name       string `json:"name"`
	Score      string `json:"score"`
	Evaluation string `json:"evaluation"`
	FoodType   string `json:"food_type"`
	Area       string `json:"area"`
	Place      string `json:"place"`
	Banner     string `json:"banner"`
	CreatedAt  int64  `json:"created_at"`
}

type Foodpanda struct {
	Name       string  `json:"name"`
	Score      float64 `json:"score"`
	Evaluation int     `json:"evaluation"`
	FoodType   string  `json:"food_type"`
	Address    string  `json:"address"`
	Latitude   string  `json:"latitude"`
	Longitude  string  `json:"longitude"`
	Banner     string  `json:"banner"`
	Img        string  `json:"img"`
	CreatedAt  int64   `json:"created_at"`
}
