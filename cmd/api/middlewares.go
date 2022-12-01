package main

import "github.com/gofiber/fiber/v2"

const (
	SessionKey    string = "session_key"
	UserIdKey     string = "user_id_key"
	IsLoggedInKey string = "is_logged_in"
)

// This middleware has the following responsabilities:
// - Create/Get a session for the request.
// - Save the session to the context's locals.
// - Call the next middleware.
// - Check the returned error of the the previously called middleware.
// - Save the session state (which could've been modified by other middlewares).
func (g *Global) ManageSession(c *fiber.Ctx) error {
	sess, err := g.Store.Get(c)
	if err != nil {
		g.Logger.Error("session create/get", err)
		return err
	}

	c.Locals(SessionKey, sess)

	err = c.Next()
	if err != nil {
		return err
	}

	err = sess.Save()
	if err != nil {
		g.Logger.Error("session save", err)
		panic(err)
	}

	return nil
}
