package transaction

import (
	"reflect"
	"strings"
	"testing"
	"time"
)

// Successfully converts a well-formed CSV file to a list of Transactions
func TestCSVtoTransactionsSuccess(t *testing.T) {
	csvContent := `date,amount,content
		2023/10/01,100,Groceries
		2023/10/02,200,Rent`

	file := strings.NewReader(csvContent)

	transactions, err := CSVtoTransactions(file)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expectedTransactions := []Transaction{
		{Date: "2023/10/01", Amount: 100, Content: "Groceries"},
		{Date: "2023/10/02", Amount: 200, Content: "Rent"},
	}

	if !reflect.DeepEqual(transactions, expectedTransactions) {
		t.Errorf("expected %v, got %v", expectedTransactions, transactions)
	}
}

// Handles CSV files with unexpected headers gracefully
func TestCSVtoTransactionsUnexpectedHeaders(t *testing.T) {
	csvContent := `date,amount,description
    2023/10/01,100,Groceries`

	file := strings.NewReader(csvContent)

	_, err := CSVtoTransactions(file)

	if err == nil {
		t.Fatal("expected an error due to unexpected headers, got none")
	}

	expectedError := "unexpected header: expected 'content', got 'description'"

	if err.Error() != expectedError {
		t.Errorf("expected error %v, got %v", expectedError, err.Error())
	}
}

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
