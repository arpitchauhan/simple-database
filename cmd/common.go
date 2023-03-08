package cmd

import (
	"errors"
	"strings"
)

func validateKey(key string) error {
	if len(strings.TrimSpace(key)) == 0 {
		return errors.New("key cannot be empty")
	}

	return nil
}
