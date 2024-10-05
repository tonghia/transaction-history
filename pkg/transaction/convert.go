package transaction

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
)

// CSVtoTransactions reads and parses the CSV file into a slice of Transactions.
func CSVtoTransactions(file io.Reader) ([]Transaction, error) {

	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true

	// Read the header line.
	headers, err := reader.Read()
	if err != nil {
		return []Transaction{}, fmt.Errorf("error reading CSV header: %v", err)
	}

	// Validate headers.
	expectedHeaders := []string{"date", "amount", "content"}
	for i, header := range headers {
		if strings.TrimSpace(strings.ToLower(header)) != expectedHeaders[i] {
			return []Transaction{}, fmt.Errorf("unexpected header: expected '%s', got '%s'", expectedHeaders[i], header)
		}
	}

	var transactions []Transaction

	// Read each record.
	recordNumber := 1 // Including header
	for {
		record, err := reader.Read()
		recordNumber++
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return []Transaction{}, fmt.Errorf("error reading CSV record at line %d: %v", recordNumber, err)
		}

		// Validate that no columns are empty.
		for i, field := range record {
			if strings.TrimSpace(field) == "" {
				return []Transaction{}, fmt.Errorf("empty field in column '%s' at line %d", expectedHeaders[i], recordNumber)
			}
		}

		// Parse the date to ensure correct format.
		dateStr := record[0]
		_, err = time.Parse("2006/01/02", dateStr)
		if err != nil {
			return []Transaction{}, fmt.Errorf("invalid date format at line %d: %v", recordNumber, err)
		}

		// Parse the amount.
		amount, err := strconv.Atoi(record[1])
		if err != nil {
			return []Transaction{}, fmt.Errorf("invalid amount at line %d: %v", recordNumber, err)
		}

		content := record[2]

		transaction := Transaction{
			Date:    dateStr,
			Amount:  amount,
			Content: content,
		}

		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

// ParseYearMonth validates and parses the YYYYMM input.
func ParseYearMonth(yearMonth string) (int, time.Month, error) {
	if len(yearMonth) != 6 {
		return 0, 0, fmt.Errorf("invalid period format. Expected YYYYMM, got: %s", yearMonth)
	}

	year, err := strconv.Atoi(yearMonth[:4])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid year in period: %v", err)
	}

	monthInt, err := strconv.Atoi(yearMonth[4:6])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid month in period: %v", err)
	}

	if monthInt < 1 || monthInt > 12 {
		return 0, 0, fmt.Errorf("month must be between 01 and 12. Got: %02d", monthInt)
	}

	return year, time.Month(monthInt), nil
}
