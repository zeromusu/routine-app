package router

import (
	"routine-app-server/internal/interfaces/handler"

	"github.com/gin-gonic/gin"

	_ "routine-app-server/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type AppHandlers struct {
	Routine handler.RoutineHandler
}

func NewRouter(h AppHandlers) *gin.Engine {
	r := gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := r.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			routines := v1.Group("/routines")
			{
				routines.GET("", h.Routine.GetAll)
				routines.GET("/:id", h.Routine.GetOne)
				routines.POST("/create", h.Routine.Create)
				routines.PUT("/:id", h.Routine.Update)
				routines.DELETE("/:id", h.Routine.Delete)
			}
		}
	}

	return r
}
