package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
)

func GetInstance(path string) {

	if path == "" {
		path = getConfig()
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Fatalf("Config file not found: %s", path)
	}
	if err := godotenv.Load(path); err != nil {
		panic(err)
	}
}

func getConfig() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	return dir
}
