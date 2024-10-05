package transaction

import (
	"fmt"
	"io"
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

// FilterTransactions filters transactions based on the specified year and month.
func FilterTransactions(transactions []Transaction, year int, month time.Month) []Transaction {
	var filtered []Transaction

	for _, tx := range transactions {
		txDate, err := time.Parse("2006/01/02", tx.Date)
		if err != nil {
			continue
		}

		if txDate.Year() == year && txDate.Month() == month {
			filtered = append(filtered, tx)
		}
	}

	return filtered
}

// CalculateTotals calculates the total income and total expenditure.
func CalculateTotals(transactions []Transaction) (int, int) {
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

// SortTransactions sorts transactions in descending order by date.
func SortTransactions(transactions []Transaction) {
	sort.Slice(transactions, func(i, j int) bool {
		dateI, _ := time.Parse("2006/01/02", transactions[i].Date)
		dateJ, _ := time.Parse("2006/01/02", transactions[j].Date)
		return dateI.After(dateJ)
	})
}

func ProcessData(file io.Reader, yearMonth string) (Summary, error) {
	// Parse the test period
	year, month, err := ParseYearMonth(yearMonth)
	if err != nil {
		return Summary{}, nil
	}

	// Read and parse the CSV file.
	transactions, err := CSVtoTransactions(file)
	if err != nil {
		return Summary{}, err
	}

	// Filter transactions based on the specified year and month.
	filteredTransactions := FilterTransactions(transactions, year, month)

	// Calculate total income and expenditure.
	totalIncome, totalExpenditure := CalculateTotals(filteredTransactions)

	// Sort transactions in descending order by date.
	SortTransactions(filteredTransactions)

	return Summary{
		Period:           fmt.Sprintf("%04d/%02d", year, month),
		TotalIncome:      totalIncome,
		TotalExpenditure: totalExpenditure,
		Transactions:     filteredTransactions,
	}, nil
}
