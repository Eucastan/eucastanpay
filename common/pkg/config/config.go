package config

import "errors"

func RequireString(name, value string) error {
	if value == "" {
		return errors.New(name + " is required")
	}

	return nil
}

func RequireMinLength(name, value string, min int) error {
	if len(value) < min {
		return errors.New(name + " is invalid")
	}

	return nil
}
