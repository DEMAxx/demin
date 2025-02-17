package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	panic("implement me")
}

func Validate(v interface{}) error {
	var validationErrors ValidationErrors

	wg := sync.WaitGroup{}
	rv := reflect.ValueOf(v)
	rt := rv.Type()

	if rt.Kind() == reflect.Ptr {
		validationErrors = append(validationErrors, ValidationError{
			Field: rt.String(),
			Err:   errors.New("v must be a pointer to structure"),
		})
		return validationErrors
	}

	if rt.Kind() == reflect.Struct {
		fmt.Println("rt.NumField()", rt.NumField())
		for i := 0; i < rt.NumField(); i++ {
			wg.Add(1)

			go func() {
				defer wg.Done()
				reflectValue := rv.Field(i)
				reflectKind := reflectValue.Kind()

				switch reflectKind {
				case reflect.String:
				case reflect.Int:
					return
				default:
					validationErrors = append(validationErrors, ValidationError{
						Field: rt.String(),
						Err:   errors.New(rt.Name() + " not string or int"),
					})
				}
			}()
		}
	}

	wg.Wait()

	if len(validationErrors) > 0 {
		return validationErrors
	}

	return nil
}
