package handlers

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	gomail "gopkg.in/mail.v2"
)

type SubmitRequest struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

func HandleSubmit(c *gin.Context) {
	m := gomail.NewMessage()

	var req SubmitRequest
	c.BindJSON(&req)

	m.SetHeader("From", os.Getenv("SMTP_USER"))
	m.SetHeader("To", os.Getenv("TARGET_EMAIL"))
	m.SetHeader("Subject", req.Subject)
	m.SetBody("text/plain", fmt.Sprintf(`You got a message from
Name: %q
Email: %q
Message:
%q`, req.Name, req.Email, req.Message))

	port, err := strconv.ParseInt(os.Getenv("SMTP_PORT"), 10, 32)
	if err != nil {
		panic(err)
	}

	d := gomail.NewDialer(
		os.Getenv("SMTP_HOST"), int(port),
		os.Getenv("SMTP_USER"), os.Getenv("SMTP_PASSWORD"),
	)

	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}

	c.Status(http.StatusOK)
}
