package hw02unpackstring

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(s string) (string, error) {
	var b strings.Builder
	saveValue := ""
	duplicateSlash := false

	for _, v := range s {
		a, err := strconv.Atoi(string(v))

		if err == nil {
			if saveValue == "\\" {
				saveValue = strconv.Itoa(a)
				continue
			} else {
				if saveValue == "" {
					if duplicateSlash {
						saveValue = "\\"
						a = a - 1
					} else {
						return "", ErrInvalidString
					}
				}
				fmt.Fprintf(&b, "%s", strings.Repeat(saveValue, a))
			}

			duplicateSlash = false
			saveValue = ""
		} else {
			if saveValue == "\\" && string(v) == "\\" {
				fmt.Fprintf(&b, "%s", saveValue)
				duplicateSlash = true
				saveValue = ""
				continue
			}

			if saveValue != "" {
				fmt.Fprintf(&b, "%s", saveValue)
			}

			saveValue = string(v)
		}
	}

	if saveValue != "" {
		fmt.Fprintf(&b, "%s", saveValue)
	}

	return b.String(), nil
}
