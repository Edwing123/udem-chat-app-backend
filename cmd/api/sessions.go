package main

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/redis"
)

// Creates a Fiber store for keeping track of
// sessions data.
func NewSessionStore(storage fiber.Storage) *session.Store {
	store := session.New(session.Config{
		Storage:        storage,
		Expiration:     time.Hour * 1,
		CookieHTTPOnly: true,
	})

	return store
}

// Creates a storage that stores sessions data to a Redis server.
func NewRedisStorage(details ConnectionDetails) fiber.Storage {
	redis := redis.New(redis.Config{
		Host:     details.Host,
		Port:     int(details.Port),
		Username: details.User,
		Password: details.Password,
	})

	return redis
}
