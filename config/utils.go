package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
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
