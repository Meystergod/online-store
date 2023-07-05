package controller

import (
	"net/http"

	"github.com/Meystergod/online-store/internal/domain/dto"
	"github.com/Meystergod/online-store/internal/repository"
	"github.com/Meystergod/online-store/internal/utils"

	"github.com/labstack/echo/v4"
)

type SubcategoryController struct {
	subcategoryRepository repository.SubcategoryRepository
}

func NewSubcategoryController(subcategoryRepository repository.SubcategoryRepository) *SubcategoryController {
	return &SubcategoryController{subcategoryRepository: subcategoryRepository}
}

func (subcategoryController *SubcategoryController) CreateSubcategory(c echo.Context) error {
	var payload dto.CreateSubcategory

	if err := utils.BindAndValidate(c, &payload); err != nil {
		return utils.Negotiate(c, http.StatusBadRequest, utils.ErrorBindAndValidatePayload.Error())
	}

	_, err := subcategoryController.subcategoryRepository.GetSubcategoryByTitle(c.Request().Context(), payload.Title)
	if err == nil {
		return utils.Negotiate(c, http.StatusConflict, "category with this title is exist")
	}

	createdSubcategoryID, err := subcategoryController.subcategoryRepository.CreateSubcategory(c.Request().Context(), payload.ToModel())
	if err != nil {
		return utils.Negotiate(c, http.StatusInternalServerError, err.Error())
	}

	return utils.Negotiate(c, http.StatusCreated, createdSubcategoryID)
}

func (subcategoryController *SubcategoryController) GetAllSubcategories(c echo.Context) error {
	subcategories, err := subcategoryController.subcategoryRepository.GetAllSubcategories(c.Request().Context())
	if err != nil {
		return utils.Negotiate(c, http.StatusInternalServerError, err.Error())
	}

	return utils.Negotiate(c, http.StatusOK, subcategories)
}

func (subcategoryController *SubcategoryController) GetSubcategory(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.Negotiate(c, http.StatusBadRequest, utils.ErrorGetUrlParams.Error())
	}

	category, err := subcategoryController.subcategoryRepository.GetSubcategory(c.Request().Context(), id)
	if err != nil {
		return utils.Negotiate(c, http.StatusInternalServerError, err.Error())
	}

	return utils.Negotiate(c, http.StatusOK, category)
}

func (subcategoryController *SubcategoryController) GetSubcategoryByTitle(c echo.Context) error {
	title := c.Param("title")
	if title == "" {
		return utils.Negotiate(c, http.StatusBadRequest, utils.ErrorGetUrlParams.Error())
	}

	subcategory, err := subcategoryController.subcategoryRepository.GetSubcategoryByTitle(c.Request().Context(), title)
	if err != nil {
		return utils.Negotiate(c, http.StatusInternalServerError, err.Error())
	}

	return utils.Negotiate(c, http.StatusOK, subcategory)
}

func (subcategoryController *SubcategoryController) UpdateSubcategory(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.Negotiate(c, http.StatusBadRequest, utils.ErrorGetUrlParams.Error())
	}

	var payload dto.UpdateSubcategory

	if err := utils.BindAndValidate(c, &payload); err != nil {
		return utils.Negotiate(c, http.StatusBadRequest, utils.ErrorBindAndValidatePayload.Error())
	}

	subcategory := payload.ToModel()
	subcategory.ID = id

	err := subcategoryController.subcategoryRepository.UpdateSubcategory(c.Request().Context(), subcategory)
	if err != nil {
		return utils.Negotiate(c, http.StatusInternalServerError, err.Error())
	}

	return utils.Negotiate(c, http.StatusOK, subcategory)
}

func (subcategoryController *SubcategoryController) DeleteSubcategory(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.Negotiate(c, http.StatusBadRequest, utils.ErrorGetUrlParams.Error())
	}

	err := subcategoryController.subcategoryRepository.DeleteSubcategory(c.Request().Context(), id)
	if err != nil {
		return utils.Negotiate(c, http.StatusInternalServerError, err.Error())
	}

	return utils.Negotiate(c, http.StatusNoContent, nil)
}
