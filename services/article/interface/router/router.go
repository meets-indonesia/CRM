package router

import (
	"github.com/gin-gonic/gin"
	"github.com/kevinnaserwan/crm-be/services/article/config"
	"github.com/kevinnaserwan/crm-be/services/article/interface/handler"
	"github.com/kevinnaserwan/crm-be/services/article/middleware"
)

// Setup initializes the router
func Setup(cfg *config.Config, articleHandler *handler.ArticleHandler) *gin.Engine {
	// Set Gin mode
	gin.SetMode(cfg.Server.Mode)

	// Create a new Gin router
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Tambahkan config ke context
	r.Use(func(c *gin.Context) {
		c.Set("config", cfg)
		c.Next()
	})

	// Health check endpoint
	r.GET("/health", articleHandler.HealthCheck)

	// Article endpoints
	articles := r.Group("/articles")
	{
		// Public endpoints
		articles.GET("", articleHandler.ListArticles)
		articles.GET("/:id", articleHandler.ViewArticle)
		articles.GET("/search", articleHandler.SearchArticles)

		// Admin only endpoints
		admin := articles.Group("")
		admin.Use(middleware.AuthMiddleware(), middleware.AdminOnlyMiddleware())
		admin.POST("", articleHandler.CreateArticle)
		admin.PUT("/:id", articleHandler.UpdateArticle)
		admin.DELETE("/:id", articleHandler.DeleteArticle)
	}

	return r
}
