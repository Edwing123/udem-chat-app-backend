package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"time"
)

func main() {
	// Get command line flags.
	flags := GetFlags()

	// Load configuration.
	config, err := LoadConfig(flags.ConfigPath)
	if err != nil {
		log.Fatalln("failed loading config: ", err)
	}

	// Validate the configuration.
	configValidationErrors := ValidateConfig(config)
	if configValidationErrors != nil {
		fmt.Println("Configuration validation failed with the following errors:")
		fmt.Println()

		for _, err := range configValidationErrors {
			fmt.Printf("\t- %s\n", err)
		}

		fmt.Println()
		os.Exit(1)
	}

	// Create appdata directories.
	err = CreateAppDataDirs(config.AppData)
	if err != nil {
		fmt.Println("An error occured while creating appdata dirs:")
		fmt.Println()

		fmt.Println(err)

		fmt.Println()
		os.Exit(1)
	}

	// Create logs file.
	logsFileName := path.Join(
		config.AppData,
		"logs",
		time.Now().Format(time.RFC3339),
	)

	logsFile, err := os.Create(logsFileName)
	if err != nil {
		log.Fatalln(err)
	}
	defer logsFile.Close()

	// Create logger, the output generated by it will be
	// stored to the logs file.
	logger := NewLogger(logsFile)

	// Create sessions store.
	store := NewSessionStore(
		NewRedisStorage(config.Redis),
	)

	// Remember to close the store's underlying storage.
	defer func() {
		store.Storage.Close()
		logger.Info("storage closed")
	}()

	global := Global{
		Logger: logger,
		Store:  store,
	}

	app := global.Setup()
	addr := config.Server.Addr

	logger.Info("server listen start", "addr", addr)

	err = app.Listen(addr)

	logger.Error("server listen end", err)
}
