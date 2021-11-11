package main

import (
	"os"
	"strconv"
)

// Get an environment string variable by name
// which is replaced by default value if it's empty
func getEnvString(key string, value string) (string, bool) {
	v := os.Getenv(key)
	if v == "" {
		return value, false
	}

	return v, true
}

// Get an environment integer variable by name
// which is replaced by default value if it's empty
func getEnvInt(key string, value int) (int, bool) {
	v := os.Getenv(key)
	if v == "" {
		return value, false
	}
	i, err := strconv.Atoi(v)
	if err != nil {
		return value, false
	}

	return i, true
}