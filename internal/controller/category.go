package controller

import (
	"net/http"

	"github.com/Meystergod/online-store/internal/domain/dto"
	"github.com/Meystergod/online-store/internal/repository"
	"github.com/Meystergod/online-store/internal/utils"

	"github.com/labstack/echo/v4"
)

type CategoryController struct {
	categoryRepository repository.CategoryRepository
}

func NewCategoryController(categoryRepository repository.CategoryRepository) *CategoryController {
	return &CategoryController{categoryRepository: categoryRepository}
}

func (categoryController *CategoryController) CreateCategory(c echo.Context) error {
	var payload dto.CreateCategory

	if err := utils.BindAndValidate(c, &payload); err != nil {
		return utils.Negotiate(c, http.StatusBadRequest, utils.ErrorBindAndValidatePayload.Error())
	}

	_, err := categoryController.categoryRepository.GetCategoryByTitle(c.Request().Context(), payload.Title)
	if err == nil {
		return utils.Negotiate(c, http.StatusConflict, "category with this title is exist")
	}

	createdCategoryID, err := categoryController.categoryRepository.CreateCategory(c.Request().Context(), payload.ToModel())
	if err != nil {
		return utils.Negotiate(c, http.StatusInternalServerError, err.Error())
	}

	return utils.Negotiate(c, http.StatusCreated, createdCategoryID)
}

func (categoryController *CategoryController) GetAllCategories(c echo.Context) error {
	categories, err := categoryController.categoryRepository.GetAllCategories(c.Request().Context())
	if err != nil {
		return utils.Negotiate(c, http.StatusInternalServerError, err.Error())
	}

	return utils.Negotiate(c, http.StatusOK, categories)
}

func (categoryController *CategoryController) GetCategory(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.Negotiate(c, http.StatusBadRequest, utils.ErrorGetUrlParams.Error())
	}

	category, err := categoryController.categoryRepository.GetCategory(c.Request().Context(), id)
	if err != nil {
		return utils.Negotiate(c, http.StatusInternalServerError, err.Error())
	}

	return utils.Negotiate(c, http.StatusOK, category)
}

func (categoryController *CategoryController) GetCategoryByTitle(c echo.Context) error {
	title := c.Param("title")
	if title == "" {
		return utils.Negotiate(c, http.StatusBadRequest, utils.ErrorGetUrlParams.Error())
	}

	category, err := categoryController.categoryRepository.GetCategoryByTitle(c.Request().Context(), title)
	if err != nil {
		return utils.Negotiate(c, http.StatusInternalServerError, err.Error())
	}

	return utils.Negotiate(c, http.StatusOK, category)
}

func (categoryController *CategoryController) UpdateCategory(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.Negotiate(c, http.StatusBadRequest, utils.ErrorGetUrlParams.Error())
	}

	var payload dto.UpdateCategory

	if err := utils.BindAndValidate(c, &payload); err != nil {
		return utils.Negotiate(c, http.StatusBadRequest, utils.ErrorBindAndValidatePayload.Error())
	}

	category := payload.ToModel()
	category.ID = id

	err := categoryController.categoryRepository.UpdateCategory(c.Request().Context(), category)
	if err != nil {
		return utils.Negotiate(c, http.StatusInternalServerError, err.Error())
	}

	return utils.Negotiate(c, http.StatusOK, category)
}

func (categoryController *CategoryController) DeleteCategory(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.Negotiate(c, http.StatusBadRequest, utils.ErrorGetUrlParams.Error())
	}

	err := categoryController.categoryRepository.DeleteCategory(c.Request().Context(), id)
	if err != nil {
		return utils.Negotiate(c, http.StatusInternalServerError, err.Error())
	}

	return utils.Negotiate(c, http.StatusNoContent, nil)
}
