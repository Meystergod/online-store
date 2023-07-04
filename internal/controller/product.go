package controller

import (
	"net/http"

	"github.com/Meystergod/online-store/internal/domain/dto"
	"github.com/Meystergod/online-store/internal/repository"
	"github.com/Meystergod/online-store/internal/utils"

	"github.com/labstack/echo/v4"
)

type ProductController struct {
	productRepository repository.ProductRepository
}

func NewProductController(productRepository repository.ProductRepository) *ProductController {
	return &ProductController{productRepository: productRepository}
}

func (productController *ProductController) CreateProduct(c echo.Context) error {
	var payload dto.CreateProduct

	if err := utils.BindAndValidate(c, &payload); err != nil {
		return utils.Negotiate(c, http.StatusBadRequest, utils.ErrorBindAndValidatePayload.Error())
	}

	_, err := productController.productRepository.GetProductByTitle(c.Request().Context(), payload.Title)
	if err == nil {
		return utils.Negotiate(c, http.StatusConflict, "product with this title is exist")
	}

	createdProductID, err := productController.productRepository.CreateProduct(c.Request().Context(), payload.ToModel())
	if err != nil {
		return utils.Negotiate(c, http.StatusInternalServerError, err.Error())
	}

	return utils.Negotiate(c, http.StatusCreated, createdProductID)
}

func (productController *ProductController) GetAllProducts(c echo.Context) error {
	products, err := productController.productRepository.GetAllProducts(c.Request().Context())
	if err != nil {
		return utils.Negotiate(c, http.StatusInternalServerError, err.Error())
	}

	return utils.Negotiate(c, http.StatusOK, products)
}

func (productController *ProductController) GetProduct(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.Negotiate(c, http.StatusBadRequest, utils.ErrorGetUrlParams.Error())
	}

	product, err := productController.productRepository.GetProduct(c.Request().Context(), id)
	if err != nil {
		return utils.Negotiate(c, http.StatusInternalServerError, err.Error())
	}

	return utils.Negotiate(c, http.StatusOK, product)
}

func (productController *ProductController) UpdateProduct(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.Negotiate(c, http.StatusBadRequest, utils.ErrorGetUrlParams.Error())
	}

	var payload dto.UpdateProduct

	if err := utils.BindAndValidate(c, &payload); err != nil {
		return utils.Negotiate(c, http.StatusBadRequest, utils.ErrorBindAndValidatePayload.Error())
	}

	product := payload.ToModel()
	product.ID = id

	err := productController.productRepository.UpdateProduct(c.Request().Context(), product)
	if err != nil {
		return utils.Negotiate(c, http.StatusInternalServerError, err.Error())
	}

	return utils.Negotiate(c, http.StatusOK, product)
}

func (productController *ProductController) DeleteProduct(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.Negotiate(c, http.StatusBadRequest, utils.ErrorGetUrlParams.Error())
	}

	err := productController.productRepository.DeleteProduct(c.Request().Context(), id)
	if err != nil {
		return utils.Negotiate(c, http.StatusInternalServerError, err.Error())
	}

	return utils.Negotiate(c, http.StatusNoContent, nil)
}
