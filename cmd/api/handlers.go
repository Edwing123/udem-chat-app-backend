package main

import (
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/Edwing123/udem-chat-app/pkg/images/profile"
	"github.com/Edwing123/udem-chat-app/pkg/models"
	"github.com/Edwing123/udem-chat-app/pkg/validations/hashing"
	"github.com/gofiber/fiber/v2"
	"github.com/h2non/bimg"
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

// Handler for updating the user information.
func (g *Global) UserUpdate(c *fiber.Ctx) error {
	var updateImage bool

	imageFile, err := c.FormFile("profileImage")
	if err == nil {
		updateImage = true
	}

	var imageId string

	// Update the image only if required.
	if updateImage {
		const maxImageSize = 1024 * 1024 * 1.5 // 1.5MB

		if imageFile.Size > maxImageSize {
			return SendErrorMessage(
				c,
				fiber.StatusBadRequest,
				ErrProfileImageTooBig,
				"El tamaño de la imagen es muy grande",
			)
		}

		// Find the corresponding type of the image based on the `bimg` package.
		var imageType bimg.ImageType
		imageFileType := imageFile.Header.Get(fiber.HeaderContentType)

		switch imageFileType {
		case "image/jpeg":
			imageType = bimg.JPEG
		case "image/webp":
			imageType = bimg.WEBP
		case "image/png":
			imageType = bimg.PNG
		default:
			imageType = bimg.UNKNOWN
		}

		// Decode the crop details.
		cropJSONString := c.FormValue("crop", "")
		crop, err := ReadJSONBody[profile.Crop]([]byte(cropJSONString))
		if err != nil {
			return SendErrorMessage(c, fiber.StatusBadRequest, ErrCannotDecodeJSON, err.Error())
		}

		// Read the image data.
		imageFile, err := imageFile.Open()
		if err != nil {
			g.Logger.Error("Open profile image", err)
			return g.ServerError(
				c,
				nil,
			)
		}
		defer imageFile.Close()

		imageBuffer, err := ioutil.ReadAll(imageFile)
		if err != nil {
			g.Logger.Error("Read profile image", err)
			return g.ServerError(
				c,
				nil,
			)
		}

		imageId, err = g.ProfileManager.New(profile.Image{
			Type:   imageType,
			Buffer: imageBuffer,
		}, crop)

		if err != nil {
			if errors.Is(err, profile.ErrImageTypeNotSupported) {
				return SendErrorMessage(
					c,
					fiber.StatusBadRequest,
					err,
					fmt.Sprintf("El tipo de archivo %s no es soportado", imageFileType),
				)
			}

			if errors.Is(err, profile.ErrCannotGetImageSize) {
				return SendErrorMessage(
					c,
					fiber.StatusBadRequest,
					err,
					"La imagen no puede ser procesada",
				)
			}

			g.Logger.Error("New profile image", err)
			return g.ServerError(c, nil)
		}
	}

	// Get user id from session.
	sess := g.GetSession(c)

	// TODO: handle panic when the casting fails.
	id := sess.Get(UserIdKey).(int)

	// Get the other values from the request body.
	userName := c.FormValue("name", "")
	birthdate := c.FormValue("birthdate", "")

	updatedUser, oldImageId, err := g.Database.UserManager.Update(id, models.User{
		Name:             userName,
		Birthdate:        birthdate,
		ProfilePictureId: imageId,
	})
	if err != nil {
		if errors.Is(err, models.ErrNoUpdates) {
			return SendErrorMessage(
				c,
				fiber.StatusBadRequest,
				err,
				"No hay nada que actualizar :|",
			)
		}

		if errors.Is(err, models.ErrUserNameExists) {
			return SendErrorMessage(
				c,
				fiber.StatusConflict,
				err,
				fmt.Sprintf("El nombre de usuario %s ya existe", userName),
			)
		}

		if errors.Is(err, models.ErrDatabaseServerFail) {
			g.Logger.Error("Update user", err)
			return g.ServerError(c, nil)
		}

		return SendErrorMessage(c, fiber.StatusBadRequest, err, "")
	}

	if oldImageId != "" {
		err = g.ProfileManager.Archive(oldImageId)
		if err != nil {
			return g.ServerError(c, nil)
		}
	}

	return SendSucessMessage(c, fiber.StatusOK, updatedUser)
}

// Handler for getting the information
// of the logged-in user.
func (g *Global) UserGet(c *fiber.Ctx) error {
	sess := g.GetSession(c)
	id := sess.Get(UserIdKey).(int)

	user, err := g.Database.UserManager.Get(id)
	if err != nil {
		return g.ServerError(c, nil)
	}

	return SendSucessMessage(c, fiber.StatusOK, user)
}
