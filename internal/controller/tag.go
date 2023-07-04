package controller

import (
	"net/http"

	"github.com/Meystergod/online-store/internal/domain/dto"
	"github.com/Meystergod/online-store/internal/repository"
	"github.com/Meystergod/online-store/internal/utils"

	"github.com/labstack/echo/v4"
)

type TagController struct {
	tagRepository repository.TagRepository
}

func NewTagController(tagRepository repository.TagRepository) *TagController {
	return &TagController{tagRepository: tagRepository}
}

func (tagController *TagController) CreateTag(c echo.Context) error {
	var payload dto.CreateTag

	if err := utils.BindAndValidate(c, &payload); err != nil {
		return utils.Negotiate(c, http.StatusBadRequest, utils.ErrorBindAndValidatePayload.Error())
	}

	createdTagID, err := tagController.tagRepository.CreateTag(c.Request().Context(), payload.ToModel())
	if err != nil {
		return utils.Negotiate(c, http.StatusInternalServerError, err.Error())
	}

	return utils.Negotiate(c, http.StatusCreated, createdTagID)
}

func (tagController *TagController) GetAllTags(c echo.Context) error {
	tags, err := tagController.tagRepository.GetAllTags(c.Request().Context())
	if err != nil {
		return utils.Negotiate(c, http.StatusInternalServerError, err.Error())
	}

	return utils.Negotiate(c, http.StatusOK, tags)
}

func (tagController *TagController) GetTag(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.Negotiate(c, http.StatusBadRequest, utils.ErrorGetUrlParams.Error())
	}

	tag, err := tagController.tagRepository.GetTag(c.Request().Context(), id)
	if err != nil {
		return utils.Negotiate(c, http.StatusInternalServerError, err.Error())
	}

	return utils.Negotiate(c, http.StatusOK, tag)
}

func (tagController *TagController) UpdateTag(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.Negotiate(c, http.StatusBadRequest, utils.ErrorGetUrlParams.Error())
	}

	var payload dto.UpdateTag

	if err := utils.BindAndValidate(c, &payload); err != nil {
		return utils.Negotiate(c, http.StatusBadRequest, utils.ErrorBindAndValidatePayload.Error())
	}

	tag := payload.ToModel()
	tag.ID = id

	err := tagController.tagRepository.UpdateTag(c.Request().Context(), tag)
	if err != nil {
		return utils.Negotiate(c, http.StatusInternalServerError, err.Error())
	}

	return utils.Negotiate(c, http.StatusOK, tag)
}

func (tagController *TagController) DeleteTag(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.Negotiate(c, http.StatusBadRequest, utils.ErrorGetUrlParams.Error())
	}

	err := tagController.tagRepository.DeleteTag(c.Request().Context(), id)
	if err != nil {
		return utils.Negotiate(c, http.StatusInternalServerError, err.Error())
	}

	return utils.Negotiate(c, http.StatusNoContent, nil)
}
