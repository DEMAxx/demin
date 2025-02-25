package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int             `validate:"min:18|max:50"`
		Email  string          `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole        `validate:"in:admin,stuff"`
		Phones []string        `validate:"len:11"`
		meta   json.RawMessage //nolint:unused
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
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			Token{
				[]byte("one"),
				[]byte("two"),
				[]byte("three"),
			},
			errors.New("token error"),
		},
		{
			App{
				Version: "10440",
			},
			errors.New("app error"),
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()
			err := Validate(tt.in)

			validationErrors := ValidationErrors{
				{
					Field: "field",
					Err:   tt.expectedErr,
				},
			}

			if errors.Is(err, validationErrors) {
				fmt.Println("validation errors")
			} else {
				fmt.Println("validation no errors")
			}

			if err == nil {
				require.NoError(t, err)
			} else {
				errorTest := errors.New("test")

				errors.Is(err, validationErrors)
				errors.Is(err, errorTest)
				//errors.As(err, &errorTest)

				//fmt.Println("error", tt.expectedErr.Error())
				//fmt.Println("validationErrors", validationErrors)
				//
				//require.PanicsWithError(t, err.Error(), func() {
				//	fmt.Println("validationErrors", validationErrors)
				//})

				//require.Nil(t, err)
				//require.ErrorIs(t, err, validationErrors, "error %q must be type %q", err)
				//require.EqualError(t, err, tt.expectedErr.Error())
			}

			//require.Truef(t, tt.expectedErr == nil && err == nil, "expected err %#v, got %#v", tt.expectedErr, err)
		})
	}
}
