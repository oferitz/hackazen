package main

import "github.com/gofiber/fiber/v2"

func publicRoutes(srv *fiber.App) {
	// Routes for GET method:
	srv.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hackazen 🤖")
	})
	srv.Get("/api", func(c *fiber.Ctx) error {
		return c.SendString("Hackazen API")
	})

}

func notFoundRoute(srv *fiber.App) {
	// Register new special route.
	srv.Use(
		// Anonymous function.
		func(c *fiber.Ctx) error {
			// Return HTTP 404 status and JSON response.
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": true,
				"msg":   "sorry, endpoint was not found",
			})
		},
	)
}

func (app *application) initRoutes() {
	srv := app.server
	publicRoutes(srv)
	route := srv.Group("/api")
	route.Get("/health", app.healthHandler)
	// Auth routes
	route.Post("/auth/signup", app.signupHandler)
	//route.Post("/auth/login", auth.Login)
	//route.Post("/auth/logout", middleware.AuthProtected, auth.Logout)
	//route.Get("/auth/user", middleware.AuthProtected, auth.Me)
	// 404
	notFoundRoute(srv)
}
