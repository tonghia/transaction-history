package transaction

import (
	"reflect"
	"testing"
	"time"
)

// Filters transactions correctly for a given year and month
func TestFilterTransactionsCorrectYearMonth(t *testing.T) {
	transactions := []Transaction{
		{Date: "2023/05/15"},
		{Date: "2023/06/20"},
		{Date: "2023/06/25"},
		{Date: "2022/06/10"},
	}

	year := 2023
	month := time.June

	filtered := FilterTransactions(transactions, year, month)

	expected := []Transaction{
		{Date: "2023/06/20"},
		{Date: "2023/06/25"},
	}

	if !reflect.DeepEqual(filtered, expected) {
		t.Errorf("Expected %v, but got %v", expected, filtered)
	}
}

// Handles transactions with invalid date formats gracefully
func TestFilterTransactionsInvalidDateFormat(t *testing.T) {
	transactions := []Transaction{
		{Date: "2023/05/15"},
		{Date: "invalid-date"},
		{Date: "2023/06/25"},
	}

	year := 2023
	month := time.June

	filtered := FilterTransactions(transactions, year, month)

	expected := []Transaction{
		{Date: "2023/06/25"},
	}

	if !reflect.DeepEqual(filtered, expected) {
		t.Errorf("Expected %v, but got %v", expected, filtered)
	}
}

// CalculateTotals correctly sums positive amounts as income
func TestCalculateTotalsSumsPositiveAmountsAsIncome(t *testing.T) {
	transactions := []Transaction{
		{Amount: 100},
		{Amount: 200},
		{Amount: -50},
	}
	totalIncome, totalExpenditure := CalculateTotals(transactions)
	if totalIncome != 300 {
		t.Errorf("Expected total income to be 300, got %d", totalIncome)
	}
	if totalExpenditure != -50 {
		t.Errorf("Expected total expenditure to be -50, got %d", totalExpenditure)
	}
}

// CalculateTotals handles transactions with zero amount correctly
func TestCalculateTotalsHandlesZeroAmount(t *testing.T) {
	transactions := []Transaction{
		{Amount: 0},
		{Amount: 100},
		{Amount: -100},
	}
	totalIncome, totalExpenditure := CalculateTotals(transactions)
	if totalIncome != 100 {
		t.Errorf("Expected total income to be 100, got %d", totalIncome)
	}
	if totalExpenditure != -100 {
		t.Errorf("Expected total expenditure to be -100, got %d", totalExpenditure)
	}
}

// Sorts transactions in descending order by date
func TestSortTransactionsDescendingOrder(t *testing.T) {
	transactions := []Transaction{
		{Date: "2023/10/01"},
		{Date: "2023/09/15"},
		{Date: "2023/10/05"},
	}

	SortTransactions(transactions)

	expectedOrder := []string{"2023/10/05", "2023/10/01", "2023/09/15"}
	for i, transaction := range transactions {
		if transaction.Date != expectedOrder[i] {
			t.Errorf("Expected %s but got %s", expectedOrder[i], transaction.Date)
		}
	}
}
