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
