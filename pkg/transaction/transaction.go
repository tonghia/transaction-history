package transaction

import (
	"fmt"
	"log"
	"os"
	"sort"
	"time"
)

// Transaction represents a single deposit or withdrawal.
type Transaction struct {
	Date    string `json:"date"`
	Amount  int    `json:"amount"`
	Content string `json:"content"`
}

// Summary represents the JSON output structure.
type Summary struct {
	Period           string        `json:"period"`
	TotalIncome      int           `json:"total_income"`
	TotalExpenditure int           `json:"total_expenditure"`
	Transactions     []Transaction `json:"transactions"`
}

// filterTransactions filters transactions based on the specified year and month.
func filterTransactions(transactions []Transaction, year int, month time.Month) []Transaction {
	var filtered []Transaction

	for _, tx := range transactions {
		txDate, err := time.Parse("2006/01/02", tx.Date)
		if err != nil {
			log.Fatalf("Error parsing transaction date '%s': %v", tx.Date, err)
		}

		if txDate.Year() == year && txDate.Month() == month {
			filtered = append(filtered, tx)
		}
	}

	return filtered
}

// calculateTotals calculates the total income and total expenditure.
func calculateTotals(transactions []Transaction) (int, int) {
	totalIncome := 0
	totalExpenditure := 0

	for _, tx := range transactions {
		if tx.Amount > 0 {
			totalIncome += tx.Amount
		} else {
			totalExpenditure += tx.Amount
		}
	}

	return totalIncome, totalExpenditure
}

// sortTransactions sorts transactions in descending order by date.
func sortTransactions(transactions []Transaction) {
	sort.Slice(transactions, func(i, j int) bool {
		dateI, _ := time.Parse("2006/01/02", transactions[i].Date)
		dateJ, _ := time.Parse("2006/01/02", transactions[j].Date)
		return dateI.After(dateJ)
	})
}

func ProcessData(file *os.File, yearMonth string) Summary {
	// Parse the test period
	year, month := parseYearMonth(yearMonth)

	// Read and parse the CSV file.
	transactions := readCSV(file)

	// Filter transactions based on the specified year and month.
	filteredTransactions := filterTransactions(transactions, year, month)

	// Calculate total income and expenditure.
	totalIncome, totalExpenditure := calculateTotals(filteredTransactions)

	// Sort transactions in descending order by date.
	sortTransactions(filteredTransactions)

	return Summary{
		Period:           fmt.Sprintf("%04d/%02d", year, month),
		TotalIncome:      totalIncome,
		TotalExpenditure: totalExpenditure,
		Transactions:     filteredTransactions,
	}
}
