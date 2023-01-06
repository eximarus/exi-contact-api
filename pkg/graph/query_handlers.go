package graph

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/eximarus/exi-contact-api/pkg/dy"
	"github.com/eximarus/exi-contact-api/pkg/graph/model"
)

func (r *queryResolver) getGuestbookImpl(
	ctx context.Context,
	cursor *string,
	inputLimit *int,
) (*model.GetGuestbookOutput, error) {
	var limit *int32
	if inputLimit != nil {
		limit = aws.Int32(int32(*inputLimit))
	} else {
		limit = aws.Int32(100)
	}

	key, err := dy.CursorDeserialize(cursor)
	if err != nil {
		return nil, err
	}

	dbInput := &dynamodb.ScanInput{
		TableName:         aws.String(dy.TableName),
		FilterExpression:  aws.String("attribute_not_exists(DeletedAt)"),
		ExclusiveStartKey: key,
		Limit:             limit,
	}

	out, err := r.Db.Scan(context.Background(), dbInput)
	if err != nil {
		return nil, err
	}

	var dbEntries []dy.GuestbookEntry
	attributevalue.UnmarshalListOfMaps(out.Items, &dbEntries)

	cursor, err = dy.CursorSerialize(out.LastEvaluatedKey)
	if err != nil {
		return nil, err
	}

	entries := make([]*model.GuestbookEntry, len(dbEntries))
	for i, entry := range dbEntries {
		entries[i] = DbGuestbookEntryToGuestbookEntry(entry)
	}

	return &model.GetGuestbookOutput{
		Cursor:  cursor,
		Entries: entries,
	}, nil
}

func DbGuestbookEntryToGuestbookEntry(dbEntry dy.GuestbookEntry,
) *model.GuestbookEntry {
	return &model.GuestbookEntry{
		Email:     dbEntry.Email,
		Name:      dbEntry.Name,
		Message:   dbEntry.Message,
		Company:   dbEntry.Company,
		CreatedAt: dbEntry.CreatedAt.String(),
	}
}
