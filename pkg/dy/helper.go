package dy

import (
	"encoding/base64"
	"encoding/json"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// TODO double check security
func CursorSerialize(cursor map[string]types.AttributeValue) (*string, error) {
	if len(cursor) < 1 {
		return nil, nil
	}

	var d map[string]interface{}
	attributevalue.UnmarshalMap(cursor, &d)

	marshalled, err := json.Marshal(d)
	if err != nil {
		return nil, err
	}

	return aws.String(base64.StdEncoding.EncodeToString(marshalled)), nil
}

func CursorDeserialize(cursor *string) (map[string]types.AttributeValue, error) {
	if cursor == nil {
		return nil, nil
	}

	jsonBytes, err := base64.StdEncoding.DecodeString(*cursor)
	if err != nil {
		return nil, err
	}

	var d map[string]interface{}
	err = json.Unmarshal(jsonBytes, &d)
	if err != nil {
		return nil, err
	}

	marshalled, err := attributevalue.MarshalMap(d)
	if err != nil {
		return nil, err
	}
	return marshalled, nil
}
