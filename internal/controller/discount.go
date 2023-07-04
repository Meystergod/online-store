package controller

import (
	"net/http"

	"online-store/internal/domain/dto"
	"online-store/internal/repository"
	"online-store/internal/utils"

	"github.com/labstack/echo/v4"
)

type DiscountController struct {
	discountRepository repository.DiscountRepository
}

func NewDiscountController(discountRepository repository.DiscountRepository) *DiscountController {
	return &DiscountController{discountRepository: discountRepository}
}

func (discountController *DiscountController) CreateDiscount(c echo.Context) error {
	var payload dto.CreateDiscount

	if err := utils.BindAndValidate(c, &payload); err != nil {
		return utils.Negotiate(c, http.StatusBadRequest, utils.ErrorBindAndValidatePayload.Error())
	}

	createdDiscountID, err := discountController.discountRepository.CreateDiscount(c.Request().Context(), payload.ToModel())
	if err != nil {
		return utils.Negotiate(c, http.StatusInternalServerError, err.Error())
	}

	return utils.Negotiate(c, http.StatusCreated, createdDiscountID)
}

func (discountController *DiscountController) GetAllDiscounts(c echo.Context) error {
	discounts, err := discountController.discountRepository.GetAllDiscounts(c.Request().Context())
	if err != nil {
		return utils.Negotiate(c, http.StatusInternalServerError, err.Error())
	}

	return utils.Negotiate(c, http.StatusOK, discounts)
}

func (discountController *DiscountController) GetDiscount(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.Negotiate(c, http.StatusBadRequest, utils.ErrorGetUrlParams.Error())
	}

	discount, err := discountController.discountRepository.GetDiscount(c.Request().Context(), id)
	if err != nil {
		return utils.Negotiate(c, http.StatusInternalServerError, err.Error())
	}

	return utils.Negotiate(c, http.StatusOK, discount)
}

func (discountController *DiscountController) UpdateDiscount(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.Negotiate(c, http.StatusBadRequest, utils.ErrorGetUrlParams.Error())
	}

	var payload dto.UpdateDiscount

	if err := utils.BindAndValidate(c, &payload); err != nil {
		return utils.Negotiate(c, http.StatusBadRequest, utils.ErrorBindAndValidatePayload.Error())
	}

	discount := payload.ToModel()
	discount.ID = id

	err := discountController.discountRepository.UpdateDiscount(c.Request().Context(), discount)
	if err != nil {
		return utils.Negotiate(c, http.StatusInternalServerError, err.Error())
	}

	return utils.Negotiate(c, http.StatusOK, discount)
}

func (discountController *DiscountController) DeleteDiscount(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.Negotiate(c, http.StatusBadRequest, utils.ErrorGetUrlParams.Error())
	}

	err := discountController.discountRepository.DeleteDiscount(c.Request().Context(), id)
	if err != nil {
		return utils.Negotiate(c, http.StatusInternalServerError, err.Error())
	}

	return utils.Negotiate(c, http.StatusNoContent, nil)
}
