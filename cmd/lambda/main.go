package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/eximarus/exi-contact-api/pkg/handlers"
	"github.com/eximarus/exi-contact-api/pkg/setup"
	"github.com/gin-gonic/gin"
)

var g *ginadapter.GinLambda

func main() {
	gin.SetMode(gin.ReleaseMode)
	lambda.Start(lambdaHandler)
}

func lambdaHandler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if g == nil {
		db := setup.InitDynamo(ctx)
		r := gin.Default()
		r.POST("/submit", handlers.HandleSubmit)
		r.POST("/guestbook", func(ctx *gin.Context) {
			handlers.HandleGuestbook(ctx, db)
		})

		g = ginadapter.New(r)
	}
	return g.ProxyWithContext(ctx, req)
}
