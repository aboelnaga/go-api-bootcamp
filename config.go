package main

// 1. Import "os" and "github.com/joho/godotenv"
import (
	"os"

	"github.com/joho/godotenv"
)

// 2. Create a Config struct with fields: Port (string) and DBPath (string)
type Config struct {
	Port   string
	DBPath string
}

//  3. Write a func LoadConfig() Config that:
//     a. Calls godotenv.Load() — but DON'T panic if it fails!
//     (The .env file is optional; in production you set real env vars)
//     b. Reads os.Getenv("PORT") — if empty, default to "3000"
//     c. Reads os.Getenv("DB_PATH") — if empty, default to "tasks.db"
//     d. Returns the populated Config struct
func LoadConfig() Config {
	godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "tasks.db"
	}

	return Config{
		Port:   port,
		DBPath: dbPath,
	}
}
