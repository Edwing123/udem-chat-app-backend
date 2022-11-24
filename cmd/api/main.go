package main

import (
	"fmt"
	"log"
)

func main() {
	flags := GetFlags()

	config, err := LoadConfig(flags.ConfigPath)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(config)
}
