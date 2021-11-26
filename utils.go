package main

import (
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
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

// Helper function to execute a shell command
func execCommand(name string, args ...string) {
	cmd := exec.Command(name, args...)
	stdout, err := cmd.Output()

	if err != nil {
		log.Printf("failed to exec command [%s]: %s\n",
			name + strings.Join(args, " "), err.Error())
		return
	}

	log.Printf("command [%s %s] stdout:\n%s",
		name, strings.Join(args, " "), string(stdout))
}