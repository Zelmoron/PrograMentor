package services

import (
	"fmt"
	"os"
)

func SaveUserCode(userID int, code string) (string, error) {
	dir := "./codes"

	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return "", fmt.Errorf("could not create code directory: %w", err)
	}

	filePath := fmt.Sprintf("%s/%d.go", dir, userID)

	if err := os.WriteFile(filePath, []byte(code), os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to save code to file: %w", err)
	}

	return filePath, nil
}
