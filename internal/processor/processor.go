package processor

import (
	"fmt"
	"io"

	"github.com/tonghia/transaction-history/internal/parser"
	"github.com/tonghia/transaction-history/internal/transaction"
	"github.com/tonghia/transaction-history/pkg/tconv"
)

func ProcessData(file io.Reader, yearMonth string) (transaction.Summary, error) {
	// Parse the test period
	year, month, err := tconv.ParseYearMonth(yearMonth)
	if err != nil {
		return transaction.Summary{}, nil
	}

	// Read and parse the CSV file.
	transactions, err := parser.CSVtoTransactions(file)
	if err != nil {
		return transaction.Summary{}, err
	}

	// Filter transactions based on the specified year and month.
	filteredTransactions := transaction.FilterTransactions(transactions, year, month)

	// Calculate total income and expenditure.
	totalIncome, totalExpenditure := transaction.CalculateTotals(filteredTransactions)

	// Sort transactions in descending order by date.
	transaction.SortTransactions(filteredTransactions)

	return transaction.Summary{
		Period:           fmt.Sprintf("%04d/%02d", year, month),
		TotalIncome:      totalIncome,
		TotalExpenditure: totalExpenditure,
		Transactions:     filteredTransactions,
	}, nil
}

func ProcessLargeData(file io.Reader, yearMonth string) (transaction.Summary, error) {
	// Parse the test period
	year, month, err := tconv.ParseYearMonth(yearMonth)
	if err != nil {
		return transaction.Summary{}, nil
	}

	// Read and parse the CSV file.
	transactions, err := parser.ChunkCSVtoTransactions(file)
	if err != nil {
		return transaction.Summary{}, err
	}

	// Filter transactions based on the specified year and month.
	filteredTransactions := transaction.FilterTransactions(transactions, year, month)

	// Calculate total income and expenditure.
	totalIncome, totalExpenditure := transaction.CalculateTotals(filteredTransactions)

	// Sort transactions in descending order by date.
	transaction.SortTransactions(filteredTransactions)

	return transaction.Summary{
		Period:           fmt.Sprintf("%04d/%02d", year, month),
		TotalIncome:      totalIncome,
		TotalExpenditure: totalExpenditure,
		Transactions:     filteredTransactions,
	}, nil
}
