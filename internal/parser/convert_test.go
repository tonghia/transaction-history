package parser

import (
	"testing"
	"time"
)

// Correctly parses a valid YYYYMM string into year and month
func TestParseYearMonthValidInput(t *testing.T) {
	yearMonth := "202309"
	expectedYear := 2023
	expectedMonth := time.September

	year, month, err := ParseYearMonth(yearMonth)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if year != expectedYear {
		t.Errorf("expected year %d, got %d", expectedYear, year)
	}

	if month != expectedMonth {
		t.Errorf("expected month %v, got %v", expectedMonth, month)
	}
}

// Handles input with incorrect length, such as "2023" or "20230912"
func TestParseYearMonthInvalidLength(t *testing.T) {
	invalidInputs := []string{"2023", "20230912"}

	for _, input := range invalidInputs {
		_, _, err := ParseYearMonth(input)

		if err == nil {
			t.Errorf("expected error for input %s, got nil", input)
		}
	}
}
