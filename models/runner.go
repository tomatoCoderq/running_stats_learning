package models

type Runner struct {
	ID string `json:"id"`
	FirstName string `json:"first_name"`
	LastName string `json:"last_name"`
	Age int `json:"age,omitempty"`
	IsActive bool `json:"is_active,omitempty"`
	Country string `json:"country"`
	PersonalBest string `json:"personal_best,omitempty"`
	SeasonBest string `json:"season_best,omitempty"`
	Results []*Result `json:"results,omitempty"`
	// `json:"age, omitempty"` are struct tags.
}