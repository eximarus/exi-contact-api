package setup

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func InitDynamo() *dynamodb.Client {
	if os.Getenv("RUN_IN_DOCKER") != "true" {
		return nil
	}

	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		if service == dynamodb.ServiceID && region == "eu-central-1" {
			return aws.Endpoint{
				PartitionID:   "aws",
				URL:           os.Getenv("DYNAMO_ENDPOINT"),
				SigningRegion: "eu-central-1",
			}, nil
		}
		return aws.Endpoint{}, &aws.EndpointNotFoundError{}
	})

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("eu-central-1"), config.WithEndpointResolverWithOptions(customResolver))

	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	db := dynamodb.NewFromConfig(cfg)

	tableName := "Guestbook"
	input := &dynamodb.CreateTableInput{
		TableName: aws.String(tableName),
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
			ReadCapacityUnits:  aws.Int64(10), // up to 25 on free tier
			WriteCapacityUnits: aws.Int64(10), // up to 25 on free tier
		},
	}

	_, err = db.CreateTable(context.TODO(), input)
	if err != nil {
		log.Fatalf("Got error calling CreateTable: %s", err)
	}

	return db
}
