package handlers

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	gomail "gopkg.in/mail.v2"
)

func HandleSubmit(c *gin.Context) {
	m := gomail.NewMessage()
	m.SetHeader("From", c.Request.FormValue("email"))
	m.SetHeader("To", os.Getenv("TARGET_EMAIL"))
	m.SetHeader("Subject", c.Request.FormValue("subject"))
	m.SetBody("text/plain", fmt.Sprintf(`You got a message from 
    Name: %s
    Message: %s`, c.Request.FormValue("name"), c.Request.FormValue("message")))

	d := gomail.NewDialer(
		"smtp.gmail.com", 587,
		os.Getenv("GMAIL_USER"), os.Getenv("GMAIL_PASSWORD"),
	)

	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}

	c.Status(http.StatusOK)
}
