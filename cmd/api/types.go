package main

import (
	"github.com/gofiber/fiber/v2/middleware/session"
	"golang.org/x/exp/slog"
)

type Global struct {
	Logger *slog.Logger
	Store  *session.Store
}

// ConnectionDetails represents the information
// needed to connect to database server.
type ConnectionDetails struct {
	User     string `json:"user"`
	Password string `json:"password"`
	Port     uint16 `json:"port"`
	Host     string `json:"host"`
}

// Config represents the configuration file.
type Config struct {
	// Redis connection details.
	Redis ConnectionDetails `"json: redis"`

	// Database connection details.
	Database ConnectionDetails `"database: redis"`

	// Server options.
	Server struct {
		// The address where the HTTP will listen on.
		Addr string
	} `"json: server"`

	// Logs options.
	Logs struct {
		// Path of the dir where logs will be stored.
		Dir string `"json: dir"`
	} `"json: logs"`
}

// Flags represents the command line flags passed
// to the executable.
type Flags struct {
	ConfigPath string // Path of the configuration file.
}
