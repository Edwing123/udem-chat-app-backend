package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"path"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

func (g *Global) GetSession(c *fiber.Ctx) *session.Session {
	return c.Locals(SessionKey).(*session.Session)
}

// Helper function to create server errors.
func (g *Global) ServerError(c *fiber.Ctx, err error) error {
	return SendErrorMessage(c, fiber.StatusInternalServerError, err)
}

// Helper function to create error response message.
func SendErrorMessage(c *fiber.Ctx, status int, err error) error {
	return c.Status(status).JSON(ErrorMessage{
		Ok:  false,
		Err: err,
	})
}

// Helper function to create response sucess message.
func SendSucessMessage[T any](c *fiber.Ctx, status int, data T) error {
	return c.Status(status).JSON(SuccessMessage[T]{
		Ok:   true,
		Data: data,
	})
}

// Parses the body of the request.
func ReadJSONBody[T any](c *fiber.Ctx) (T, error) {
	var v T

	err := c.BodyParser(&v)
	if err != nil {
		return v, err
	}

	return v, nil
}

// Reads the JSON configuration file from the provided
// file system path, then parses the configuration into
// the `Config` struct, and finally returns it with a
// nill error if everything goes fine. Otherwise, it returns
// a non-nil error, which should be checked.
func LoadConfig(path string) (Config, error) {
	fd, err := os.Open(path)
	if err != nil {
		return Config{}, err
	}
	defer fd.Close()

	decoder := json.NewDecoder(fd)
	decoder.DisallowUnknownFields()

	var config Config

	err = decoder.Decode(&config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}

// Defines, parses and returns the command line flags.
func GetFlags() Flags {
	config := flag.String("config", "", "The path of the configuration file")

	flag.Parse()

	if *config == "" {
		log.Fatalln("The flag [config] is required")
	}

	return Flags{
		ConfigPath: *config,
	}
}

// Validates the configuration, returns a slice of strings
// with the errors. If none were found, the return slice
// is nil.
func ValidateConfig(config Config) []string {
	redis := config.Redis
	database := config.Database

	var validationsErrors []string

	redisErrors := ValidateConnectionDetails(redis)
	if redisErrors != nil {
		for _, err := range redisErrors {
			validationsErrors = append(validationsErrors, fmt.Sprintf("redis: %s", err))
		}
	}

	databaseErrors := ValidateConnectionDetails(database)
	if databaseErrors != nil {
		for _, err := range databaseErrors {
			validationsErrors = append(validationsErrors, fmt.Sprintf("database: %s", err))
		}
	}

	addrError := ValidateAddr(config.Server.Addr)
	if addrError != nil {
		validationsErrors = append(validationsErrors, fmt.Sprintf("server: %s", addrError.Error()))
	}

	if config.AppData == "" {
		validationsErrors = append(validationsErrors, "field required: appdata")
	}

	return validationsErrors
}

// Validates fields User, Password and Host are not empty.
func ValidateConnectionDetails(details ConnectionDetails) []string {
	var validationsErrors []string

	if details.User == "" {
		validationsErrors = append(validationsErrors, "field required: user")
	}

	if details.Password == "" {
		validationsErrors = append(validationsErrors, "field required: password")
	}

	if details.Host == "" {
		validationsErrors = append(validationsErrors, "field required: host")
	}

	if len(validationsErrors) > 0 {
		return validationsErrors
	}

	return nil
}

// Validates whether the passed address is valid
// based on the rules of the function `net.SplitHostPort`.
func ValidateAddr(addr string) error {
	_, _, err := net.SplitHostPort(addr)
	return err
}

// Create the dirs where data generated by the
// application will be stored.
func CreateAppDataDirs(appDataPath string) error {
	stats, err := os.Stat(appDataPath)

	// Return an error if it does exist but it's not a directory.
	if err == nil && !stats.IsDir() {
		return fmt.Errorf("%s already exists, but it's not a directory", appDataPath)
	} else {
		// Create it if it doesn't exist.
		if errors.Is(err, os.ErrNotExist) {
			err := os.MkdirAll(appDataPath, 0o755)
			if err != nil {
				return err
			}
		} else {
			// TODO: investigate what other errors can happen here.
			return err
		}
	}

	// Create logs and images directories.
	err = os.MkdirAll(path.Join(appDataPath, "logs"), 0o755)
	if err != nil {
		return err
	}

	err = os.MkdirAll(path.Join(appDataPath, "images"), 0o755)
	if err != nil {
		return err
	}

	return nil
}
