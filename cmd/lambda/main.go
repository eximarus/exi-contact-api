package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/eximarus/exi-contact-api/pkg/handlers"
	"github.com/eximarus/exi-contact-api/pkg/setup"
	"github.com/gin-gonic/gin"
)

func main() {
	setup.InitDynamo()
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.POST("/submit", handlers.HandleSubmit)

	ginLambda := ginadapter.New(r)
	lambda.Start(ginLambda.ProxyWithContext)
}
