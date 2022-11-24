package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
)

// Reads the JSON configuration file from the provided
// file system path, then parses the configuration into
// the `Config` struct, and finally returns it with a
// nill error if everything goes fine. Otherwise, it returns
// a non-nil error.
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
	logs := flag.String("logsDir", "", "The path of the logs dir")
	addr := flag.String("addr", "", "The addr for the server to listen on")

	flag.Parse()

	if *config == "" {
		log.Fatalln("The flag [config] is required")
	}

	if *logs == "" {
		log.Fatalln("The flag [logsDir] is required")
	}

	if *addr == "" {
		log.Fatalln("The flag [addr] is required")
	}

	return Flags{
		ConfigPath:  *config,
		LogsDirPath: *logs,
		Addr:        *addr,
	}
}
