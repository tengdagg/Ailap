package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"ailap-backend/internal/database"
	"ailap-backend/internal/handler"
	"ailap-backend/internal/middleware"
)

func New() *gin.Engine {
	_ = database.Init()

	r := gin.New()
	r.Use(middleware.RequestLogger())
	r.Use(gin.Recovery())
	r.Use(cors.Default())

	authHandler := handler.NewAuthHandler()
	logsHandler := handler.NewLogsHandler()
	modelsHandler := handler.NewModelsHandler()
	dsHandler := handler.NewDataSourcesHandler()

	api := r.Group("/api")
	{
		auth := api.Group("/auth")
		auth.POST("/login", authHandler.Login)
		auth.POST("/logout", authHandler.Logout)
		auth.GET("/profile", authHandler.Profile)

		// Test datasource connection (no auth required)
		api.POST("/datasources/test", dsHandler.Test)

		api.Use(middleware.AuthRequired())

		logs := api.Group("/logs")
		logs.GET("/query", logsHandler.Query)
		logs.GET("/suggestions", logsHandler.Suggestions)
		logs.GET("/label-values", logsHandler.LabelValues)
		logs.GET("/history", logsHandler.History)
		logs.POST("/history/:id/favorite", logsHandler.ToggleFavorite)
		logs.PUT("/history/:id/note", logsHandler.UpdateNote)
		logs.DELETE("/history/:id", logsHandler.DeleteHistory)
		logs.GET("/inspect", logsHandler.Inspect)

		models := api.Group("/models")
		models.GET("", modelsHandler.List)
		models.POST("", modelsHandler.Create)
		models.POST("/test", modelsHandler.Test)
		models.PUT(":id", modelsHandler.Update)
		models.POST(":id/enabled", modelsHandler.ToggleEnabled)
		models.POST(":id/default", modelsHandler.SetDefault)
		models.DELETE(":id", modelsHandler.Delete)

		ds := api.Group("/datasources")
		ds.GET("", dsHandler.List)
		ds.POST("", dsHandler.Create)
		ds.PUT(":id", dsHandler.Update)
		ds.DELETE(":id", dsHandler.Delete)
		ds.POST(":id/test", dsHandler.Test) // optional test by id (requires auth)
	}

	return r
}
