package httpecho

import (
	"online-store/internal/controller"

	"github.com/labstack/echo/v4"
)

func SetProductApiRoutes(e *echo.Echo, productController *controller.ProductController) {
	v1 := e.Group("/api/v1")
	{
		v1.POST("/product", productController.CreateProduct)
		v1.GET("/products", productController.GetAllProducts)
		v1.GET("/product/:id", productController.GetProduct)
		v1.PUT("/product/:id", productController.UpdateProduct)
		v1.DELETE("/product/:id", productController.DeleteProduct)
	}
}
