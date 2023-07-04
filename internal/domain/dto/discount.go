package dto

import "online-store/internal/domain/model"

type CreateDiscount struct {
	Title    string `json:"title" bson:"title" validate:"required"`
	Percent  int    `json:"percent" bson:"percent" validate:"required"`
	IsActive bool   `json:"is-active" bson:"is-active" validate:"required"`
}

type UpdateDiscount struct {
	Title    string `json:"title" bson:"title" validate:"required"`
	Percent  int    `json:"percent" bson:"percent" validate:"required"`
	IsActive bool   `json:"is-active" bson:"is-active" validate:"required"`
}

func (createDiscount *CreateDiscount) ToModel() *model.Discount {
	return &model.Discount{
		Title:    createDiscount.Title,
		Percent:  createDiscount.Percent,
		IsActive: createDiscount.IsActive,
	}
}

func (updateDiscount *UpdateDiscount) ToModel() *model.Discount {
	return &model.Discount{
		Title:    updateDiscount.Title,
		Percent:  updateDiscount.Percent,
		IsActive: updateDiscount.IsActive,
	}
}
