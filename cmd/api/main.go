package main

import (
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
		log.Fatalln("failed loading config")
	}

	// Create logs file.
	logsFileName := path.Join(config.Logs.Dir, time.Now().Format(time.RFC3339))

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
