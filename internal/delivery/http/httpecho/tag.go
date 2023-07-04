package httpecho

import (
	"online-store/internal/controller"

	"github.com/labstack/echo/v4"
)

func SetTagApiRoutes(e *echo.Echo, tagController *controller.TagController) {
	v1 := e.Group("/api/v1")
	{
		v1.POST("/tag", tagController.CreateTag)
		v1.GET("/tags", tagController.GetAllTags)
		v1.GET("/tag/:id", tagController.GetTag)
		v1.PUT("/tag/:id", tagController.UpdateTag)
		v1.DELETE("/tag/:id", tagController.DeleteTag)
	}
}
