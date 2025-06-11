package main

import (
	"budget-app/cmd/api/handlers"
	"budget-app/cmd/api/middlewares"
	"budget-app/common"
	"budget-app/internal/mailer"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"os"
)

type Application struct {
	logger        echo.Logger
	server        *echo.Echo
	handler       handlers.Handler
	appMiddleware middlewares.AppMiddleware
}

func main() {
	e := echo.New()
	err := godotenv.Load()
	if err != nil {
		e.Logger.Fatal(err.Error())
	}

	db, err := common.NewMySql()

	if err != nil {
		e.Logger.Fatal("Error loading .env file")
	}

	appMailer := mailer.NewMailer(e.Logger)
	handler := handlers.Handler{
		Database: db,
		Logger:   e.Logger,
		Mailer:   appMailer,
	}

	appMiddleware := middlewares.AppMiddleware{
		Logger: e.Logger,
		Db:     db,
	}

	app := Application{
		logger:        e.Logger,
		server:        e,
		handler:       handler,
		appMiddleware: appMiddleware,
	}

	e.Use(middleware.Logger())
	e.Use(middlewares.SecondMiddleware)
	e.Use(middlewares.CustomMiddleware)
	app.routues(handler)
	e.GET("/", handler.HealthCheck)
	fmt.Println(app)
	port := os.Getenv("APP_PORT")
	appAddress := fmt.Sprintf("localhost:%s", port)
	e.Logger.Fatal(e.Start(appAddress))
}
