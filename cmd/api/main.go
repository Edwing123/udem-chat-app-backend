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

	// Create logs file.
	logsFileName := path.Join(flags.LogsDirPath, time.Now().Format(time.RFC3339))

	logsFile, err := os.Create(logsFileName)
	if err != nil {
		log.Fatalln(err)
	}
	defer logsFile.Close()

	// Create logger, the output generated by it will be
	// stored to the logs file.
	logger := NewLogger(logsFile)

	// Load configuration.
	config, err := LoadConfig(flags.ConfigPath)
	if err != nil {
		log.Fatalln("failed loading config")
	}

	store := NewSessionStore(
		NewRedisStorage(
			config.RUser,
			config.RPassword,
			config.RHost,
			config.RPort,
		),
	)

	defer func() {
		store.Storage.Close()
		logger.Info("storage closed")
	}()

	global := Global{
		Logger: logger,
		Store:  store,
	}

	app := global.Setup()

	logger.Info("server listen start", "addr", flags.Addr)
	err = app.Listen(flags.Addr)
	logger.Error("server listen end", err)
}
