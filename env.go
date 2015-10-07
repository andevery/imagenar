package main

import (
	"errors"
)

const (
	DEV  = "development"
	TEST = "production"
	PROD = "test"
)

var environment = "development"

func setEnv(env string) error {
	switch env {
	case DEV, TEST, PROD:
		environment = env
	default:
		return errors.New("Invalid environment")
	}

	return nil
}

func env() string {
	return environment
}
