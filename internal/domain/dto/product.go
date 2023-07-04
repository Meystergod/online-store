package dto

import "github.com/Meystergod/online-store/internal/domain/model"

type CreateProduct struct {
	Title       string         `json:"title" bson:"title" validate:"required"`
	Description string         `json:"description" bson:"description" validate:"required"`
	Price       string         `json:"price" bson:"price" validate:"required"`
	Quantity    int            `json:"quantity" bson:"quantity" validate:"required"`
	Category    model.Category `json:"category" bson:"category,omitempty"`
	Discount    model.Discount `json:"discount" bson:"discount,omitempty"`
	Tags        []model.Tag    `json:"tags" bson:"description,omitempty"`
}

type UpdateProduct struct {
	Title       string         `json:"title" bson:"title" validate:"required"`
	Description string         `json:"description" bson:"description" validate:"required"`
	Price       string         `json:"price" bson:"price" validate:"required"`
	Quantity    int            `json:"quantity" bson:"quantity" validate:"required"`
	Category    model.Category `json:"category" bson:"category,omitempty"`
	Discount    model.Discount `json:"discount" bson:"discount,omitempty"`
	Tags        []model.Tag    `json:"tags" bson:"description,omitempty"`
}

func (createDiscount *CreateProduct) ToModel() *model.Product {
	return &model.Product{
		Title:       createDiscount.Title,
		Description: createDiscount.Description,
		Price:       createDiscount.Price,
		Quantity:    createDiscount.Quantity,
		Category:    createDiscount.Category,
		Discount:    createDiscount.Discount,
		Tags:        createDiscount.Tags,
	}
}

func (updateDiscount *UpdateProduct) ToModel() *model.Product {
	return &model.Product{
		Title:       updateDiscount.Title,
		Description: updateDiscount.Description,
		Price:       updateDiscount.Price,
		Quantity:    updateDiscount.Quantity,
		Category:    updateDiscount.Category,
		Discount:    updateDiscount.Discount,
		Tags:        updateDiscount.Tags,
	}
}
