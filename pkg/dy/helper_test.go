package dy

import (
	"os"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func TestCursorSerialization(t *testing.T) {
	os.Setenv("AES_ENCRYPTION_SECRET", "abc&1*~#^2^#s0^=)^^7%b34")

	cursor := map[string]types.AttributeValue{
		"Email": &types.AttributeValueMemberS{
			Value: "example@example.com",
		},
	}
	cursorStr, err := CursorSerialize(cursor)
	t.Log(*cursorStr)

	if err != nil {
		t.Error(err)
	}
	if cursorStr == nil {
		t.Errorf("Cursor was nil after serialize")
	}

	deser, err := CursorDeserialize(cursorStr)

	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(cursor, deser) {
		t.Errorf("Cursor changed in deserialize")
	}
}
