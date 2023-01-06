package setup

import (
	"context"
	"log"
	"os"
	"strings"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/eximarus/exi-contact-api/pkg/graph"
	"github.com/gin-gonic/gin"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/aws/aws-sdk-go-v2/aws/ratelimit"
	"github.com/aws/aws-sdk-go-v2/aws/retry"
)

const (
	// DefaultMaxAttempts is the maximum of attempts for an API request
	DefaultMaxAttempts int = 6 // sdk=3

	// DefaultMaxBackoff is the maximum back off delay between attempts
	DefaultMaxBackoff = 20 * time.Second // sdk default

	// retry token quota values.
	DefaultRetryRateTokens  uint = 10_000 // sdk=500
	DefaultRetryCost        uint = 5      // sdk default
	DefaultRetryTimeoutCost uint = 10     // sdk default
	DefaultNoRetryIncrement uint = 1      // sdk default
)

func InitGraphqlHandler(ctx context.Context, db *dynamodb.Client) gin.HandlerFunc {
	h := handler.NewDefaultServer(graph.NewExecutableSchema(
		graph.Config{Resolvers: &graph.Resolver{
			Db: db,
		}},
	))

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func InitDynamo(ctx context.Context) *dynamodb.Client {
	retryer := config.WithRetryer(func() aws.Retryer {
		return newRetryer()
	})

	cfg, err := config.LoadDefaultConfig(ctx, retryer)

	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	endpoint := os.Getenv("DYNAMO_ENDPOINT")
	db := dynamodb.NewFromConfig(cfg, func(options *dynamodb.Options) {
		if endpoint != "" {
			options.EndpointResolver = dynamodb.EndpointResolverFromURL(endpoint)
		}
	})

	err = ensureTable(db, createGuestbookTableInput())
	if err != nil {
		log.Fatalf("Got error calling CreateTable: %s", err)
	}

	return db
}

func newRetryer() aws.Retryer {
	return retry.NewStandard(func(options *retry.StandardOptions) {
		options.MaxAttempts = DefaultMaxAttempts
		options.MaxBackoff = DefaultMaxBackoff
		options.Retryables = retry.DefaultRetryables
		options.RateLimiter = ratelimit.NewTokenRateLimit(DefaultRetryRateTokens)
		options.RetryCost = DefaultRetryCost
		options.RetryTimeoutCost = DefaultRetryTimeoutCost
		options.NoRetryIncrement = DefaultNoRetryIncrement
	})
}

func ensureTable(db *dynamodb.Client, createTableInput *dynamodb.CreateTableInput) error {
	_, err := db.CreateTable(context.Background(), createTableInput)

	if err == nil {
		log.Println("CreateTable OK")
		return nil
	}
	str := err.Error()
	if strings.Contains(str, "ResourceInUseException") || strings.Contains(str, "Cannot create preexisting table") {
		log.Println("CreateTable OK. error ignored: " + str)
		return nil // ignore silently
	}

	return err
}

func createGuestbookTableInput() *dynamodb.CreateTableInput {
	return &dynamodb.CreateTableInput{
		TableName: aws.String("Guestbook"),
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("Email"),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("Email"),
				KeyType:       types.KeyTypeHash,
			},
		},
		ProvisionedThroughput: &types.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(1), // up to 25 on free tier
			WriteCapacityUnits: aws.Int64(1), // up to 25 on free tier
		},
	}
}
