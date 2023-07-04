package dto

import "online-store/internal/domain/model"

type CreateTag struct {
	Title string `json:"title" bson:"title" validate:"required"`
}

type UpdateTag struct {
	Title string `json:"title" bson:"title" validate:"required"`
}

func (createTag *CreateTag) ToModel() *model.Tag {
	return &model.Tag{
		Title: createTag.Title,
	}
}

func (updateTag *UpdateTag) ToModel() *model.Tag {
	return &model.Tag{
		Title: updateTag.Title,
	}
}
