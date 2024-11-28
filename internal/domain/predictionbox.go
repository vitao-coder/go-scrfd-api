package domain

type PredictionBox struct {
	Score     float32   `json:"score"`
	BoxLeft   float32   `json:"box_left"`
	BoxTop    float32   `json:"box_top"`
	BoxRight  float32   `json:"box_right"`
	BoxBottom float32   `json:"box_bottom"`
	KeyPoints []float32 `json:"key_points"`
}
