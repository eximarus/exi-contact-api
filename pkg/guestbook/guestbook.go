package guestbook

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type CreateGuestbookEntryRequest struct {
	Email   string  `json:"email"`
	Message string  `json:"message"`
	Name    string  `json:"name"`
	Company *string `json:"company"`
}

type GuestbookEntry struct {
	Email     string
	Message   string
	Name      string
	Company   *string
	CreatedAt time.Time
}

const tableName string = "Guestbook"

func CreateGuestbookEntry(db *dynamodb.Client, req *CreateGuestbookEntryRequest) error {
	item := &GuestbookEntry{
		Email:     req.Email,
		Name:      req.Name,
		Company:   req.Company,
		Message:   req.Message,
		CreatedAt: time.Now(),
	}

	av, err := attributevalue.MarshalMap(item)
	if err != nil {
		return err
	}
	input := &dynamodb.PutItemInput{
		Item:                av,
		TableName:           aws.String(tableName),
		ConditionExpression: aws.String("attribute_not_exists(pk)"),
	}

	_, err = db.PutItem(context.TODO(), input)
	if err != nil {
		return err
	}
	return nil
}
