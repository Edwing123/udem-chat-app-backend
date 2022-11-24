package main

import (
	"io"

	"golang.org/x/exp/slog"
)

func NewLogger(output io.Writer) *slog.Logger {
	logger := slog.New(slog.NewTextHandler(output))
	return logger
}
