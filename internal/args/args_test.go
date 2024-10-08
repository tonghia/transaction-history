package args

import (
	"path/filepath"
	"testing"
)

func TestParsePeriodWithValidInput(t *testing.T) {
	input := "2023-10"
	expected := "2023-10"

	result, err := ParsePeriod(input)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result != expected {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestParsePeriodWithEmptyInput(t *testing.T) {
	input := ""

	_, err := ParsePeriod(input)
	if err == nil {
		t.Fatal("expected an error, got none")
	}

	expectedError := "-period argument are required"
	if err.Error() != expectedError {
		t.Errorf("expected error %v, got %v", expectedError, err.Error())
	}
}

// Parse valid file path and return absolute path
func TestParseFilePathValid(t *testing.T) {
	validFilePath := "args_test.go"
	absPath, err := ParseFilePath(validFilePath)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	expectedAbsPath, _ := filepath.Abs(validFilePath)
	if absPath != expectedAbsPath {
		t.Errorf("expected %s, got %s", expectedAbsPath, absPath)
	}
}

// Handle empty file path string
func TestParseFilePathEmpty(t *testing.T) {
	_, err := ParseFilePath("")
	if err == nil {
		t.Fatal("expected an error for empty file path, got nil")
	}
	expectedError := "-file arguments are required"
	if err.Error() != expectedError {
		t.Errorf("expected error message %q, got %q", expectedError, err.Error())
	}
}
