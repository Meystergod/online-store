package httpecho

import (
	"github.com/Meystergod/online-store/internal/controller"

	"github.com/labstack/echo/v4"
)

func SetCategoryApiRoutes(e *echo.Echo, categoryController *controller.CategoryController) {
	v1 := e.Group("/api/v1")
	{
		v1.POST("/category", categoryController.CreateCategory)
		v1.GET("/categories", categoryController.GetAllCategories)
		v1.GET("/category/:id", categoryController.GetCategory)
		v1.PUT("/category/:id", categoryController.UpdateCategory)
		v1.DELETE("/category/:id", categoryController.DeleteCategory)
	}
}
