package main

import (
	"errors"
	"fmt"

	"github.com/Edwing123/udem-chat-app/pkg/models"
	"github.com/Edwing123/udem-chat-app/pkg/validations/hashing"
	"github.com/gofiber/fiber/v2"
)

// Handler for authenticating user.
func (g *Global) UserLogIn(c *fiber.Ctx) error {
	credentials, err := ReadBodyFromRequest[models.User](c)
	if err != nil {
		return SendErrorMessage(c, fiber.StatusBadRequest, ErrCannotDecodeJSON, err)
	}

	id, err := g.Database.UserManager.Login(credentials)
	if err != nil {
		if errors.Is(err, models.ErrLoginFail) {
			return SendErrorMessage(c, fiber.StatusUnauthorized, err, "")
		}

		if errors.Is(err, models.ErrDatabaseServerFail) {
			return g.ServerError(c, nil)
		}

		return SendErrorMessage(c, fiber.StatusBadRequest, err, "")
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
	user, err := ReadBodyFromRequest[models.User](c)
	if err != nil {
		return SendErrorMessage(c, fiber.StatusBadRequest, ErrCannotDecodeJSON, err.Error())
	}

	err = ValidatePassword(user.Password)
	if err != nil {
		return SendErrorMessage(
			c,
			fiber.StatusBadRequest,
			ErrPasswordNotValid,
			"Contraseña no cumple las reglas de validacion",
		)
	}

	err = g.Database.UserManager.New(user)
	if err != nil {
		if errors.Is(err, models.ErrUserNameExists) {
			return SendErrorMessage(
				c,
				fiber.StatusConflict,
				err,
				fmt.Sprintf("Usuario %s ya existe", user.Name),
			)
		}

		if errors.Is(err, models.ErrDatabaseServerFail) || errors.Is(err, hashing.ErrPasswordHashingFail) {
			return g.ServerError(c, nil)
		}

		return SendErrorMessage(c, fiber.StatusBadRequest, err, "")
	}

	return SendSucessMessage(c, fiber.StatusCreated, "Usuario registrado")
}

// Handler for getting the status of the user.
// whether or not it's logged-in.
func (g *Global) UserStatus(c *fiber.Ctx) error {
	sess := g.GetSession(c)
	isLoggedIn, _ := sess.Get(IsLoggedInKey).(bool)

	return SendSucessMessage(c, fiber.StatusOK, fiber.Map{
		"isActive": isLoggedIn,
	})
}
