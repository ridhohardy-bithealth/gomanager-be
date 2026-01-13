package main

import (
	"log"
	"os"
	"ps-gogo-manajer/internal/config"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("failed to load env", err.Error())
		return
	}

	log := config.NewLogger()
	validator := config.NewValidator()
	app := echo.New()
	s3Client := config.NewS3Client()
	pg := config.NewDatabase(log)
	defer pg.Pool.Close()

	config.Bootstrap(&config.BootstrapConfig{
		App:       app,
		DB:        pg,
		Log:       log,
		Validator: validator,
		S3Client:  s3Client,
	})

	PORT := os.Getenv("PORT")
	log.Fatal(app.Start(PORT))
}
