package test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"transaction-history/pkg/transaction"
)

// TestSummaryGeneration tests the summary generation for a specified period.
func TestSummaryGeneration(t *testing.T) {
	// Define the test period and file paths
	testPeriod := "202201"
	testDataDir := "testdata"
	transactionsFilePath := filepath.Join(testDataDir, "transactions.csv")
	expectedSummaryFilePath := filepath.Join(testDataDir, "summary.json")

	// Load expected summary from summary.json
	transactionsFile, err := os.Open(transactionsFilePath)
	if err != nil {
		t.Fatalf("Failed to read transactions file: %v", err)
	}
	defer transactionsFile.Close()

	// Load expected summary from summary.json
	expectedSummaryData, err := os.ReadFile(expectedSummaryFilePath)
	if err != nil {
		t.Fatalf("Failed to read expected summary file: %v", err)
	}

	var expectedSummary transaction.Summary
	err = json.Unmarshal(expectedSummaryData, &expectedSummary)
	if err != nil {
		t.Fatalf("Failed to unmarshal expected summary JSON: %v", err)
	}

	generatedSummary := transaction.ProcessData(transactionsFile, testPeriod)

	// Compare the generated summary with the expected summary
	if generatedSummary.Period != expectedSummary.Period {
		t.Errorf("Period mismatch: expected %s, got %s", expectedSummary.Period, generatedSummary.Period)
	}

	if generatedSummary.TotalIncome != expectedSummary.TotalIncome {
		t.Errorf("TotalIncome mismatch: expected %d, got %d", expectedSummary.TotalIncome, generatedSummary.TotalIncome)
	}

	if generatedSummary.TotalExpenditure != expectedSummary.TotalExpenditure {
		t.Errorf("TotalExpenditure mismatch: expected %d, got %d", expectedSummary.TotalExpenditure, generatedSummary.TotalExpenditure)
	}

	if len(generatedSummary.Transactions) != len(expectedSummary.Transactions) {
		t.Errorf("Number of transactions mismatch: expected %d, got %d", len(expectedSummary.Transactions), len(generatedSummary.Transactions))
	}

	// Compare each transaction
	for i, expectedTx := range expectedSummary.Transactions {
		generatedTx := generatedSummary.Transactions[i]

		if generatedTx.Date != expectedTx.Date {
			t.Errorf("Transaction %d date mismatch: expected %s, got %s", i, expectedTx.Date, generatedTx.Date)
		}

		if generatedTx.Amount != expectedTx.Amount {
			t.Errorf("Transaction %d amount mismatch: expected %d, got %d", i, expectedTx.Amount, generatedTx.Amount)
		}

		if generatedTx.Content != expectedTx.Content {
			t.Errorf("Transaction %d content mismatch: expected %s, got %s", i, expectedTx.Content, generatedTx.Content)
		}
	}
}
