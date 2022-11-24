package main

// Config represents the configuration file.
type Config struct {
	// Redis connection details.
	RUser     string `json:"rUser"`
	RPassword string `json:"rPassword"`
	RPort     uint   `json:"rPort"`
	RHost     string `json:"rHost"`

	// Database connection details.
	DbUser     string `json:"dbUser"`
	DbPassword string `json:"dbPassword"`
	DbPort     uint   `json:"dbPort"`
	DbHost     string `json:"dbHost"`
}

// Flags represents the command line flags passed
// to the executable.
type Flags struct {
	ConfigPath  string // Path of the configuration file.
	LogsDirPath string // Path of the dir where the logs file will be stored.
}
