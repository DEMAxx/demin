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

	for _, v := range s {
		a, err := strconv.Atoi(string(v))

		if err == nil {
			if saveValue == "\\" {
				saveValue = strconv.Itoa(a)
				continue
			} else {
				if saveValue == "" {
					return "", ErrInvalidString
				}
				fmt.Fprintf(&b, "%s", strings.Repeat(saveValue, a))
			}

			saveValue = ""
		} else {
			if saveValue == "\\" && string(v) == "\\" {
				fmt.Fprintf(&b, "%s", saveValue)
				saveValue = ""
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
