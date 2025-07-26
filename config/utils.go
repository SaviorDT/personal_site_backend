package config

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
)

var variableCache = make(map[string]string)

func Init() error {
	err := godotenv.Load()
	if err != nil {
		return fmt.Errorf("Error loading .env file: %v", err)
	}
	return nil
}

func GetVariableAsString(varName string) (string, error) {
	if value, found := variableCache[varName]; found {
		return value, nil
	}

	value := os.Getenv(varName)
	if value == "" {
		return "", fmt.Errorf("%s not set in .env", varName)
	}
	variableCache[varName] = value

	return value, nil
}

func GetVariableAsByteArr(varName string) ([]byte, error) {
	value, err := GetVariableAsString(varName)
	if err != nil {
		return nil, err
	}
	return []byte(value), nil
}

func GetVariableAsTimeDuration(varName string) (time.Duration, error) {
	value, err := GetVariableAsString(varName)
	if err != nil {
		return 0, err
	}
	duration, err := time.ParseDuration(value)
	if err != nil {
		return 0, fmt.Errorf("Error parsing %s as duration: %v", varName, err)
	}
	return duration, nil
}
