package guestbook

import (
	"context"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type GuestbookEntry struct {
	Email     string
	Message   string
	Name      string
	Company   *string
	CreatedAt time.Time
}

const tableName string = "guestbook"

func CreateGuestbookEntry(db *dynamodb.Client) {
	item := &GuestbookEntry{
		Email:     "eximarus@gmail.com",
		Name:      "Alexander Gaugusch",
		Company:   aws.String("Naughty Cult Ltd."),
		Message:   "Hello dynamo",
		CreatedAt: time.Now(),
	}

	av, err := attributevalue.MarshalMap(item)
	if err != nil {
		log.Fatalf("Got error marshalling map: %s", err)
	}
	input := &dynamodb.PutItemInput{
		Item:                av,
		TableName:           aws.String(tableName),
		ConditionExpression: aws.String("attribute_not_exists(pk)"),
	}

	_, err = db.PutItem(context.TODO(), input)
	if err != nil {
		log.Fatalf("Got error calling PutItem: %s", err)
	}
}
