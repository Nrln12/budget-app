package main

import (
	"budget-app/cmd/api/handlers"
)

func (app *Application) routues(handler handlers.Handler) {
	apiGroup := app.server.Group("/v1/users")
	publicRoute := apiGroup.Group("/public")
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
	app.server.GET("/", handler.HealthCheck)
}
