package utils

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type Validator struct {
	validator *validator.Validate
}

func NewValidator() echo.Validator {
	return &Validator{validator: validator.New()}
}

func (v *Validator) Validate(i interface{}) error {
	return v.validator.Struct(i)
}

func BindAndValidate(c echo.Context, i interface{}) error {
	if err := c.Bind(i); err != nil {
		return BadRequestException(err.Error())
	}

	if err := c.Validate(i); err != nil {
		return BadRequestException(err.Error())
	}

	return nil
}
