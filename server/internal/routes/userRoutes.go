package routes

import (
	"net/http"

	"tanoserver/internal/auth"
	"tanoserver/internal/middleware"
)

func NewAuthHandler(c *auth.AuthController) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc(
		"POST /user/register",
		c.RegisterController,
	)

	mux.HandleFunc(
		"POST /user/login",
		c.LoginController,
	)

	mux.HandleFunc(
		"DELETE /user/delete",
		middleware.AuthMiddleware(
			c.DeleteController,
		),
	)

	mux.HandleFunc(
		"POST /user/logout",
		middleware.AuthMiddleware(
			c.LogoutController,
		),
	)

	return mux
}
