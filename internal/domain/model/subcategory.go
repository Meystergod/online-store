package model

type Subcategory struct {
	ID          string `json:"uuid" bson:"_id,omitempty"`
	Title       string `json:"title" bson:"title" validate:"required"`
	Description string `json:"description" bson:"description" validate:"required"`
}
