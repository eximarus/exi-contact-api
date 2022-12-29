package main

import (
	"os"

	handlers "github.com/eximarus/exi-contact-api/pkg"
	"github.com/gin-gonic/gin"
)

func main() {
	setupLocalEnv()
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.POST("/submit", handlers.HandleSubmit)
	r.Run(":8080")
}

func setupLocalEnv() {
	if os.Getenv("RUN_IN_DOCKER") == "true" {
		return
	}

	os.Setenv("TARGET_EMAIL", "")
	os.Setenv("GMAIL_USER", "")
	os.Setenv("GMAIL_PASSWORD", "")
}
