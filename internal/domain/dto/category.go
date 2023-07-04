package dto

import (
	"online-store/internal/domain/model"
)

type CreateCategory struct {
	Title       string `json:"title" bson:"title" validate:"required"`
	Description string `json:"description" bson:"description" validate:"required"`
}

type UpdateCategory struct {
	Title       string `json:"title" bson:"title" validate:"required"`
	Description string `json:"description" bson:"description" validate:"required"`
}

func (createCategory *CreateCategory) ToModel() *model.Category {
	return &model.Category{
		Title:       createCategory.Title,
		Description: createCategory.Description,
	}
}

func (updateCategory *UpdateCategory) ToModel() *model.Category {
	return &model.Category{
		Title:       updateCategory.Title,
		Description: updateCategory.Description,
	}
}
