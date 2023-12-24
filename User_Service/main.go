package main

import (
	"context"
	"errors"
	"log"
	"syscall"

	"net/http"
	"os"
	"os/signal"
	"time"

	"User_Service/config"
	"User_Service/config/database"
	v1 "User_Service/delivery/v1"
	"User_Service/entity"
	"User_Service/lib/logger"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title ERP core Service API Documentation.
// @version 1.0
// @description This is a sample api documentation.

// @host localhost:1327
// @BasePath /api/v1

func main() {
	conf := config.NewConfig("config.env")
	appLogger := logger.NewApiLogger(conf)

	appLogger.InitLogger()
	appLogger.Info("Starting the API Server")
	db := database.NewDB(conf)

	e := echo.New()
	//e.Logger.SetLevel(log.INFO)
	// Enable HTTP compression
	e.Use(middleware.Gzip())

	// Recover from panics
	e.Use(middleware.Recover())

	// Allow requests from *
	e.Use(middleware.CORS())
	// Print http request and response log to stdout if debug is enabled
	if conf.Debug {
		e.Use(middleware.Logger())
	}

	// JWT Middleware
	jwtConfig := middleware.JWTConfig{
		Claims:       &entity.JwtClaims{},
		SigningKey:   []byte(conf.JwtSecret),
		ErrorHandler: v1.InvalidJwt,
	}

	v1.SetupRouters(e, conf, db, jwtConfig, appLogger)

	e.GET("/swagger/*", echoSwagger.WrapHandler)
	go v1.SubscribeToRedis(db)
	go httpServer(e, conf.HTTP)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	log.Println("Shutting down HTTP server...")
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
	log.Println("HTTP server stopped!")
}

func httpServer(e *echo.Echo, httpConfig config.HTTP) {
	if err := e.Start(httpConfig.HTTPAddress); err != nil && !errors.Is(err, http.ErrServerClosed) {
		e.Logger.Fatal("shutting down the server")
	}
}
