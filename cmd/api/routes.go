package main

import (
	"budget-app/cmd/api/handlers"
)

func (app *Application) routues(handler handlers.Handler) {
	apiGroup := app.server.Group("/v1")
	publicRoute := apiGroup.Group("/users/public")
	{
		publicRoute.POST("/register", handler.RegisterHandler)
		publicRoute.POST("/login", handler.LoginHandler)
		publicRoute.POST("/forgot-password", handler.ForgotPasswordHandler)
		publicRoute.POST("/reset-password", handler.ResetPasswordHandler)
	}

	profileRoute := apiGroup.Group("/profile", app.appMiddleware.AuthenticateMiddleware)
	{
		profileRoute.GET("/authenticated/user", handler.GetAuthenticatedUserHandler)
		profileRoute.PATCH("/change-password", handler.ChangePasswordHandler)
	}

	categoryRoute := apiGroup.Group("/categories", app.appMiddleware.AuthenticateMiddleware)
	{
		categoryRoute.GET("", handler.GetCategories)
		categoryRoute.POST("", handler.CreateCategory)
		categoryRoute.DELETE("/:id", handler.DeleteCategory)
	}
	app.server.GET("/", handler.HealthCheck)
}
