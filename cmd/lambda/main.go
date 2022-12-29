package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	handlers "github.com/eximarus/exi-contact-api/pkg"
	"github.com/gin-gonic/gin"
)

var ginLambda *ginadapter.GinLambda

func init() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.POST("/submit", handlers.HandleSubmit)

	ginLambda = ginadapter.New(r)
}

func main() {
	lambda.Start(handleAPIGateway)
}

func handleAPIGateway(ctx context.Context, req events.APIGatewayProxyRequest) (resp events.APIGatewayProxyResponse, err error) {
	return ginLambda.ProxyWithContext(ctx, req)
}
