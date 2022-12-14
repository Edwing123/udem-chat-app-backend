package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func (g *Global) Setup() *fiber.App {
	app := fiber.New(fiber.Config{
		AppName:       "Nameless",
		ServerHeader:  "Go+FiberV2",
		CaseSensitive: true,
		StrictRouting: true,
	})

	// Define global middlewares.
	app.Use(
		recover.New(),
		logger.New(),
		g.ManageSession,
	)

	// Define profile images route.
	app.Get("/images/profile/:id<guid>", g.ProfileManager.ServeImage)

	// Group API endpoints under the same group.
	api := app.Group("/api")

	user := api.Group("/user")
	user.Post("/login", g.UserLogIn)
	user.Post("/signup", g.UserSignUp)
	user.Post("/logout", g.RequireAuth, g.UserLogout)
	user.Get("/status", g.UserStatus)
	user.Patch("/update", g.RequireAuth, g.UserUpdate)
	user.Get("/data", g.RequireAuth, g.UserGet)

	// TODO: remove later.
	api.Get("/hello", func(c *fiber.Ctx) error {
		sess := g.GetSession(c)

		visits, _ := sess.Get("visits").(int)
		visits++
		sess.Set("visits", visits)

		return c.SendString(fmt.Sprintf("Hello (%d)\n", visits))
	})

	return app
}
