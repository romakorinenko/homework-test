package hw09structvalidator

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// Test the function on different structures and other types.
type (
	UserRole string

	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"len:11"`
		meta   json.RawMessage
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}
)

func TestValidate(t *testing.T) {
	testCases := []struct {
		name        string
		givenStruct interface{}
		expectedErr error
	}{
		{
			name: "struct is valid",
			givenStruct: User{
				ID:     "1q2w3e4r5t6y7u8i9o0pqwertyuiop123456",
				Name:   "Name",
				Age:    20,
				Email:  "mail@mail.com",
				Role:   UserRole("admin"),
				Phones: []string{"12345678901"},
				meta:   nil,
			},
			expectedErr: nil,
		},
		{
			name: "struct is valid, phones are empty",
			givenStruct: User{
				ID:     "1q2w3e4r5t6y7u8i9o0pqwertyuiop123456",
				Name:   "Name",
				Age:    20,
				Email:  "mail@mail.com",
				Role:   UserRole("admin"),
				Phones: []string{},
				meta:   nil,
			},
			expectedErr: nil,
		},
		{
			name:        "value is not a struct",
			givenStruct: 42,
			expectedErr: ErrTypeIsNotStruct,
		},

		{
			name: "struct does not have tags in fields",
			givenStruct: Token{
				[]byte{3, 6, 2, 5},
				[]byte{1, 8, 3, 6},
				[]byte{10, 11, 15, 16},
			},
			expectedErr: nil,
		},
		{
			name:        "empty struct",
			givenStruct: struct{}{},
			expectedErr: nil,
		},
		{
			name: "unknown field type",
			givenStruct: struct {
				Bool bool `validate:"1:0"`
			}{
				Bool: true,
			},
			expectedErr: ErrUnknownType,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			tt := testCase
			t.Parallel()
			err := Validate(tt.givenStruct)
			require.ErrorIs(t, err, tt.expectedErr)
		})
	}
}

func TestValidate_UserIsNotValid(t *testing.T) {
	user := User{
		ID:     "111111",
		Name:   "Name",
		Age:    17,
		Email:  "mail",
		Role:   UserRole("user"),
		Phones: []string{"212121"},
		meta:   nil,
	}
	err := Validate(user)
	require.Error(t, err)
	builder := strings.Builder{}
	builder.WriteString("ID: 111111 occurs err: 'ID' does not have length = '36'\n")
	builder.WriteString("Age: 17 occurs err: 'Age' is less than: '18'\n")
	builder.WriteString("Email: mail occurs err: 'Email' doesn't match regex: '^\\w+@\\w+\\.\\w+$'\n")
	builder.WriteString("Role: user occurs err: 'Role' not in array of values: 'admin,stuff'\n")
	builder.WriteString("Phones: 212121 occurs err: 'Phones' does not have length = '11'\n")
	require.Equal(t, builder.String(), err.Error())
}

func TestValidate_ParsingErr(t *testing.T) {
	givenStruct := struct {
		Value string `validate:"givenStruct"`
	}{
		Value: "for",
	}
	err := Validate(givenStruct)
	require.Error(t, err)
	require.Equal(t, "failed to validate field 'Value': invalid check in validate tag", err.Error())
}
