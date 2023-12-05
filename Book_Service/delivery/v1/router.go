package v1

import (
	"Book_Service/config"
	"Book_Service/lib/logger"
	"Book_Service/repository"
	"Book_Service/service"
	"github.com/labstack/echo/v4"

	"github.com/labstack/echo/v4/middleware"
	"github.com/uptrace/bun"
)

// Setup all routers
func SetupRouters(c *echo.Echo, conf *config.Config, db *bun.DB, jwtConfig middleware.JWTConfig, logger logger.Logger) {

	booksRepository := repository.NewBooksRepository(db)
	booksService := service.NewBooksService(booksRepository)
	booksHandler := NewBooksHandler(booksService, logger)

	authenticated := middleware.JWTWithConfig(jwtConfig)

	v1 := c.Group("/api/v1")

	booksGroup := v1.Group("/books")

	booksHandler.MapBooksRoutes(booksGroup, authenticated)

}
