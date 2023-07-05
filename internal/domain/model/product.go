package model

type Product struct {
	ID          string      `json:"uuid" bson:"_id,omitempty"`
	Title       string      `json:"title" bson:"title" validate:"required"`
	Description string      `json:"description" bson:"description" validate:"required"`
	Price       string      `json:"price" bson:"price" validate:"required"`
	Quantity    int         `json:"quantity" bson:"quantity" validate:"required"`
	Category    Category    `json:"category" bson:"category,omitempty"`
	Subcategory Subcategory `json:"subcategory" bson:"subcategory,omitempty"`
	Discount    Discount    `json:"discount" bson:"discount,omitempty"`
	Tags        []Tag       `json:"tags" bson:"description,omitempty"`
}
