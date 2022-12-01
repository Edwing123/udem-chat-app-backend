package main

import (
	"errors"
	"fmt"

	"github.com/Edwing123/udem-chat-app/pkg/codes"
	"github.com/Edwing123/udem-chat-app/pkg/models"
	"github.com/gofiber/fiber/v2"
)

// Handler for authenticating user.
func (g *Global) UserLogIn(c *fiber.Ctx) error {
	credentials, err := ReadJSONBody[models.User](c)
	if err != nil {
		return SendErrorMessage(c, fiber.StatusBadRequest, codes.ErrClient, err)
	}

	id, err := g.Database.UserManager.Login(credentials)
	if err != nil {
		if errors.Is(err, codes.ErrLoginFail) {
			return SendErrorMessage(c, fiber.StatusUnauthorized, codes.ErrLoginFail, err)
		}

		return g.ServerError(c, err)
	}

	// Save user id and status inside its session.
	sess := g.GetSession(c)
	sess.Set(UserIdKey, id)
	sess.Set(IsLoggedInKey, true)

	return SendSucessMessage(c, fiber.StatusOK, fiber.Map{
		"id": id,
	})
}

// Handler for logging out user.
func (g *Global) UserLogout(c *fiber.Ctx) error {
	sess := g.GetSession(c)

	err := sess.Destroy()
	if err != nil {
		return g.ServerError(c, err)
	}

	return SendSucessMessage(c, fiber.StatusOK, "Sesion cerrada")
}

// Handler for signing up a new user.
func (g *Global) UserSignUp(c *fiber.Ctx) error {
	user, err := ReadJSONBody[models.User](c)
	if err != nil {
		return SendErrorMessage(c, fiber.StatusBadRequest, codes.ErrClient, err)
	}

	err = ValidatePassword(user.Password)
	if err != nil {
		return SendErrorMessage(
			c,
			fiber.StatusBadRequest,
			codes.ErrPasswordNotValid,
			"Contraseña no cumple las reglas de validacion",
		)
	}

	err = g.Database.UserManager.New(user)
	if err != nil {
		if errors.Is(err, codes.ErrUserNameAlreadyExists) {
			return SendErrorMessage(
				c,
				fiber.StatusConflict,
				codes.ErrUserNameAlreadyExists,
				fmt.Sprintf("Usuario %s ya existe", user.Name),
			)
		}

		return g.ServerError(c, err)
	}

	return SendSucessMessage(c, fiber.StatusCreated, "Usuario registrado")
}

// Handler for getting the status of the user.
// whether or not its' logged-in.
func (g *Global) UserStatus(c *fiber.Ctx) error {
	sess := g.GetSession(c)
	userId, ok := sess.Get(UserIdKey).(int)

	if !ok {
		return SendErrorMessage(c, fiber.StatusBadRequest, codes.ErrNotLoggedIn, "")
	}

	return SendSucessMessage(c, fiber.StatusOK, fiber.Map{
		"id": userId,
	})
}
