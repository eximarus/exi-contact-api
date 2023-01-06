package main

import (
	"context"
	"os"

	"github.com/eximarus/exi-contact-api/pkg/handlers"
	"github.com/eximarus/exi-contact-api/pkg/setup"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	setupLocalEnv()
	setup.InitDynamo(context.Background())
	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()
	r.POST("/submit", handlers.HandleSubmit)
	r.POST("/guestbook", handlers.HandleSubmit)
	r.Run(":8080")
}

func setupLocalEnv() {
	if os.Getenv("RUN_IN_DOCKER") == "true" {
		return
	}
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	os.Setenv("DYNAMO_ENDPOINT", "http://localhost:8000")
	os.Setenv("AWS_REGION", "eu-central-1")
	os.Setenv("AWS_ACCESS_KEY_ID", "local")
	os.Setenv("AWS_SECRET_ACCESS_KEY", "local")
}
