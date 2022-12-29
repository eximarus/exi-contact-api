package main

import (
	"github.com/eximarus/exi-contact-api/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
	// targetEmail := os.Getenv("TARGET_EMAIL")
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.POST("/submit", handlers.HandleSubmit)
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
