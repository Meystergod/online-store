package model

type Tag struct {
	ID    string `json:"uuid" bson:"_id,omitempty"`
	Title string `json:"title" bson:"title" validate:"required"`
}
