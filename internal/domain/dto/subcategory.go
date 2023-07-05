package dto

import (
	"github.com/Meystergod/online-store/internal/domain/model"
)

type CreateSubcategory struct {
	Title       string `json:"title" bson:"title" validate:"required"`
	Description string `json:"description" bson:"description" validate:"required"`
}

type UpdateSubcategory struct {
	Title       string `json:"title" bson:"title" validate:"required"`
	Description string `json:"description" bson:"description" validate:"required"`
}

func (createSubcategory *CreateSubcategory) ToModel() *model.Subcategory {
	return &model.Subcategory{
		Title:       createSubcategory.Title,
		Description: createSubcategory.Description,
	}
}

func (updateSubcategory *UpdateSubcategory) ToModel() *model.Subcategory {
	return &model.Subcategory{
		Title:       updateSubcategory.Title,
		Description: updateSubcategory.Description,
	}
}
