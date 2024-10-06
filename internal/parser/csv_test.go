package parser

import (
	"reflect"
	"strings"
	"testing"

	"github.com/tonghia/transaction-history/internal/transaction"
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

	expectedTransactions := []transaction.Transaction{
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
