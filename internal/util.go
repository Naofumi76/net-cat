package internal

import (
	"errors"
)

func StringToInt(s string) (int, error) {
	n := 0
	for _, char := range s {
		if char < '0' || char > '9' {
			return 0, errors.New("invalid number")
		}
		n = n*10 + int(char-'0')
		if n < 0 {
			return 0, errors.New("integer overflow")
		}
	}
	return n, nil
}
