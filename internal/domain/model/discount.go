package model

type Discount struct {
	ID       string `json:"uuid" bson:"_id,omitempty"`
	Title    string `json:"title" bson:"title" validate:"required"`
	Percent  int    `json:"percent" bson:"percent" validate:"required"`
	IsActive bool   `json:"is-active" bson:"is-active" validate:"required"`
}
