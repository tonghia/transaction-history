package parser

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/tonghia/transaction-history/internal/transaction"
)

// CSVtoTransactions reads and parses the CSV file into a slice of Transactions.
func CSVtoTransactions(file io.Reader, expectedHeaders []string) ([]transaction.Transaction, error) {

	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true

	var transactions []transaction.Transaction

	// Read each record.
	recordNumber := 1 // Including header
	for {
		record, err := reader.Read()
		recordNumber++
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return []transaction.Transaction{}, fmt.Errorf("error reading CSV record at line %d: %v", recordNumber, err)
		}

		// Validate that no columns are empty.
		for i, field := range record {
			if strings.TrimSpace(field) == "" {
				return []transaction.Transaction{}, fmt.Errorf("empty field in column '%s' at line %d", expectedHeaders[i], recordNumber)
			}
		}

		// Parse the date to ensure correct format.
		dateStr := record[0]
		if _, err := time.Parse("2006/01/02", dateStr); err != nil {
			return []transaction.Transaction{}, fmt.Errorf("invalid date format at line %d: %v", recordNumber, err)
		}

		// Parse the amount.
		amount, err := strconv.Atoi(record[1])
		if err != nil {
			return []transaction.Transaction{}, fmt.Errorf("invalid amount at line %d: %v", recordNumber, err)
		}

		content := record[2]

		transaction := transaction.Transaction{
			Date:    dateStr,
			Amount:  amount,
			Content: content,
		}

		transactions = append(transactions, transaction)
	}

	return transactions, nil
}
