package graph

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/eximarus/exi-contact-api/pkg/dy"
	"github.com/eximarus/exi-contact-api/pkg/graph/model"

	gomail "gopkg.in/mail.v2"
)

func (r *mutationResolver) submitContactInfoImpl(ctx context.Context, input *model.ContactInfoInput) (*bool, error) {
	m := gomail.NewMessage()
	m.SetHeader("From", os.Getenv("SMTP_USER"))
	m.SetHeader("To", os.Getenv("TARGET_EMAIL"))
	m.SetHeader("Subject", input.Subject)
	m.SetBody("text/plain", fmt.Sprintf(`You got a message from
Name: %q
Email: %q
Message:
%q`, input.Name, input.Email, input.Message))

	port, err := strconv.ParseInt(os.Getenv("SMTP_PORT"), 10, 32)
	if err != nil {
		return nil, err
	}
	d := gomail.NewDialer(
		os.Getenv("SMTP_HOST"), int(port),
		os.Getenv("SMTP_USER"), os.Getenv("SMTP_PASSWORD"),
	)

	if err := d.DialAndSend(m); err != nil {
		return nil, err
	}

	return nil, nil
}

func (r *mutationResolver) createGuestbookEntryImpl(ctx context.Context, input *model.CreateGuestbookEntryInput) (*model.GuestbookEntry, error) {
	item := &dy.GuestbookEntry{
		Email:     input.Email,
		Name:      input.Name,
		Company:   input.Company,
		Message:   input.Message,
		CreatedAt: time.Now(),
	}

	av, err := attributevalue.MarshalMap(item)
	if err != nil {
		return nil, err
	}
	dbInput := &dynamodb.PutItemInput{
		Item:                av,
		TableName:           aws.String(dy.TableName),
		ConditionExpression: aws.String("attribute_not_exists(Email)"),
	}

	_, err = r.Db.PutItem(context.Background(), dbInput)
	if err != nil {
		return nil, err
	}
	return &model.GuestbookEntry{
		Email:     item.Email,
		Name:      item.Name,
		Company:   item.Company,
		Message:   item.Message,
		CreatedAt: item.CreatedAt.String(),
	}, nil
}
