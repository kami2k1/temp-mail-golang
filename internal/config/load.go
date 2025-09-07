package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadConfig() {
    err := godotenv.Load()
    if err != nil {
        log.Fatalf("Error loading .env file")
    }

    Config = &Dataconfig{
        HOST: os.Getenv("HOST" ),
        PORT: os.Getenv("PORT"),
        APP_ENV: os.Getenv("APP_ENV"),
        IMAP_HOST: os.Getenv("STMP_IMAP"),
        IMAP_PORT: os.Getenv("STMP_IMAP_PORT"),
        STMP_USER: os.Getenv("STMP_USER"),
        STMP_PASS: os.Getenv("STMP_PASS"),
        JWT_SECRET: func() string {
            if s := os.Getenv("JWT_SECRET"); s != "" {
                return s
            }
           
            return os.Getenv("STMP_PASS")
        }(),
    }
}
	
