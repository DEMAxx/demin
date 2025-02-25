package hw09structvalidator

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"
)

type ValidationError struct {
	Field string
	Err   error
}

var timeType = reflect.TypeOf(time.Time{})

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	panic("implement me")
}

func Validate(v interface{}) error {
	var validationErrors []ValidationErrors

	wg := &sync.WaitGroup{}
	rv := reflect.ValueOf(v)

	if rv.Kind() == reflect.Ptr {
		validationErrors = append(
			validationErrors,
			[]ValidationError{
				{
					Field: rv.String(),
					Err:   errors.New("v must not be a Ptr"),
				},
			},
		)
		return errors.Join()
	}

	if rv.Kind() != reflect.Struct || rv.Type().ConvertibleTo(timeType) {
		validationErrors = append(validationErrors, []ValidationError{
			{
				Field: rv.String(),
				Err:   errors.New("v must be a pointer to structure"),
			},
		})
		return errorJoiner(validationErrors)
	}

	valid := validator.New(validator.WithRequiredStructEnabled())

	_ = valid

	//err := valid.Struct(reflect.New(rv.Type()).Interface())
	//
	//if err != nil {
	//	fmt.Println("error", err.Error())
	//	fmt.Println("type", rv.Type())
	//	//validationErrors = append(validationErrors, ValidationError{})
	//}
	fmt.Println("rt", rv)

	//if err != nil {
	//	validationErrors = append(validationErrors, ValidationError{})
	//}

	st := reflect.TypeOf(v)

	for i := 0; i < rv.Type().NumField(); i++ {
		wg.Add(1)

		field := st.Field(i)

		err := traverseField(field, rv.Field(i), &validationErrors, wg)

		if err != nil {
			return errorJoiner(validationErrors)
		}
	}

	wg.Wait()

	if len(validationErrors) > 0 {
		return errorJoiner(validationErrors)
	}

	return nil
}

func traverseField(field reflect.StructField, current reflect.Value, validationErrors *[]ValidationErrors, wg *sync.WaitGroup) []ValidationError {
	defer wg.Done()

	var errorSlice []ValidationError
	var typ reflect.Type
	var kind reflect.Kind

	typ = current.Type()
	kind = current.Kind()

	if alias, ok := field.Tag.Lookup("validate"); ok {
		if alias == "" {
			fmt.Println("(len)")
		} else {
			fmt.Println("len alias", alias)
		}
	} else {
		fmt.Println("(not specified)")
	}

	if kind == reflect.Ptr {
		return ValidationErrors{
			{
				Field: current.String(),
				Err:   errors.New("v must not be a Ptr"),
			},
		}
	}

	if typ == timeType {
		return nil
	}

	st := reflect.TypeOf(current)

	if kind == reflect.Struct {
		rv := reflect.ValueOf(current)

		for i := 0; i < rv.NumField(); i++ {
			wg.Add(1)

			field := st.Field(i)

			err := traverseField(field, rv, validationErrors, wg)

			if err != nil {
				errorSlice = append(errorSlice, ValidationError{
					Field: rv.String(),
					Err:   errors.New(fmt.Sprintf("%s not a struct", rv.Type().Field(i).Name)),
				})
			}
		}

		if len(errorSlice) > 0 {
			return errorSlice
		}
	}

	if kind == reflect.String {
		return validateSplitter(current.String(), current.String())
	}

	return ValidationErrors{
		{
			Field: current.String(),
			Err:   errors.New(fmt.Sprintf("%q: is incampatible value", field.Name)),
		},
	}
}

func validateSplitter(value, field string) ValidationErrors {
	var err ValidationErrors

	splitValues := strings.Split(value, "|")

	for _, v := range splitValues {
		e := validateTag(v, field)

		err = append(err, ValidationError{
			Field: field,
			Err:   e,
		})
	}

	return err
}

func validateTag(value, field string) error {

	splitValue := strings.Split(field, ":")

	if len(splitValue) != 2 {
		return errors.New("value must contain at least one ':' separator")
	}

	switch splitValue[0] {
	case "required":
		if len(value) == 0 {
			return errors.New("value must contain at least one ':' separator")
		}
	case "len":
		length, err := strconv.Atoi(splitValue[1])
		if err != nil {
			return errors.New("value must contain at least one ':' separator")
		}

		if len(value) != length {
			return errors.New("value must contain at least one ':' separator")
		}
	case "in":
		hasValue := false

		for _, v := range strings.Split(splitValue[1], ",") {
			if v == value {
				hasValue = true
			}
		}

		if !hasValue {
			return errors.New("value must contain at least one ':' separator")
		}
	case "min":

	}

	return nil
}

func errorJoiner(validationErrors []ValidationErrors) error {

	for _, err := range validationErrors {
		return errors.New(err.Error())
	}

	return nil
}
