package httpecho

import (
	"online-store/internal/controller"

	"github.com/labstack/echo/v4"
)

func SetDiscountApiRoutes(e *echo.Echo, discountController *controller.DiscountController) {
	v1 := e.Group("/api/v1")
	{
		v1.POST("/discount", discountController.CreateDiscount)
		v1.GET("/discounts", discountController.GetAllDiscounts)
		v1.GET("/discount/:id", discountController.GetDiscount)
		v1.PUT("/discount/:id", discountController.UpdateDiscount)
		v1.DELETE("/discount/:id", discountController.DeleteDiscount)
	}
}
