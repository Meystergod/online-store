package httpecho

import (
	"github.com/Meystergod/online-store/internal/controller"

	"github.com/labstack/echo/v4"
)

func SetSubcategoryApiRoutes(e *echo.Echo, subcategoryController *controller.SubcategoryController) {
	v1 := e.Group("/api/v1")
	{
		v1.POST("/subcategory", subcategoryController.CreateSubcategory)
		v1.GET("/subcategories", subcategoryController.GetAllSubcategories)
		v1.GET("/subcategory/:id", subcategoryController.GetSubcategory)
		v1.PUT("/subcategory/:id", subcategoryController.UpdateSubcategory)
		v1.DELETE("/subcategory/:id", subcategoryController.DeleteSubcategory)
	}
}
