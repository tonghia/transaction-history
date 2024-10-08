package processor

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/tonghia/transaction-history/internal/parser"
	"github.com/tonghia/transaction-history/internal/transaction"
	"github.com/tonghia/transaction-history/pkg/tconv"
)

type part struct {
	offset, size int64
}

var expectedHeaders = []string{"date", "amount", "content"}

func Process(filePath string, yearMonth string, workerNum int) (json.RawMessage, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error opening CSV file: %v", err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	header, err := reader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("error reading header: %v", err)
	}
	for i, header := range strings.Split(header, ",") {
		if strings.TrimSpace(strings.ToLower(header)) != expectedHeaders[i] {
			return nil, fmt.Errorf("unexpected header: expected '%s', got '%s'", expectedHeaders[i], header)
		}
	}

	var summary transaction.Summary
	if workerNum <= 1 {
		rs, err := ProcessData(reader, yearMonth)
		if err != nil {
			return nil, fmt.Errorf("error processing CSV file: %v", err)
		}
		summary = rs
	} else {
		// Determine non-overlapping parts for file split (each part has offset and size).
		parts, err := splitFile(filePath, workerNum, len(header))
		if err != nil {
			return nil, fmt.Errorf("error spliting file")
		}
		// Start a goroutine to process each part, returning results on a channel.
		resultsCh := make(chan transaction.Summary)
		for _, part := range parts {
			go processPart(filePath, part.offset, part.size, yearMonth, resultsCh)
		}

		for i := 0; i < len(parts); i++ {
			result := <-resultsCh
			summary.TotalIncome = summary.TotalIncome + result.TotalIncome
			summary.TotalExpenditure = summary.TotalExpenditure + result.TotalExpenditure
			summary.Transactions = append(summary.Transactions, result.Transactions...)
		}

		summary.Period = yearMonth
		transaction.SortTransactions(summary.Transactions)
	}

	// Generate JSON output.
	jsonData, err := json.MarshalIndent(summary, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("error marshaling JSON: %v", err)
	}

	return jsonData, nil
}

func ProcessData(file io.Reader, yearMonth string) (transaction.Summary, error) {
	// Parse the test period
	year, month, err := tconv.ParseYearMonth(yearMonth)
	if err != nil {
		return transaction.Summary{}, nil
	}

	// Read and parse the CSV file.
	transactions, err := parser.CSVtoTransactions(file, expectedHeaders)
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

// splitFile splits a file into multiple parts based on the specified number of parts
// refer to https://github.com/benhoyt/go-1brc/blob/fafba3256ea28631f6b3739f6d3b711a91199861/r8.go#L124
func splitFile(inputPath string, numParts int, initOffset int) ([]part, error) {
	const maxLineLength = 100 // Assumption

	f, err := os.Open(inputPath)
	if err != nil {
		return nil, err
	}
	st, err := f.Stat()
	if err != nil {
		return nil, err
	}
	size := st.Size()
	splitSize := size / int64(numParts)

	buf := make([]byte, maxLineLength)

	parts := make([]part, 0, numParts)
	offset := int64(initOffset)
	for offset < size {
		seekOffset := max(offset+splitSize-maxLineLength, offset)
		if seekOffset > size {
			break
		}
		_, err := f.Seek(seekOffset, io.SeekStart)
		if err != nil {
			return nil, err
		}
		n, _ := io.ReadFull(f, buf)
		chunk := buf[:n]
		newline := bytes.LastIndexByte(chunk, '\n')
		if newline < 0 {
			// Case 1: there is no `\n` at the end of the file, we stop
			// Case 2: maxLineLength is too small for the line, we accept there will be a huge chunk and improve it later
			break
		}
		remaining := len(chunk) - newline - 1
		nextOffset := seekOffset + int64(len(chunk)) - int64(remaining)
		parts = append(parts, part{offset, nextOffset - offset})
		offset = nextOffset
	}

	return parts, nil
}

func processPart(inputPath string, fileOffset int64, fileSize int64, yearMonth string, resultsCh chan<- transaction.Summary) {
	file, err := os.Open(inputPath)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	if _, err := file.Seek(fileOffset, io.SeekStart); err != nil {
		panic(err)
	}

	f := io.LimitReader(file, fileSize)

	summary, err := ProcessData(f, yearMonth)
	if err != nil {
		log.Fatalf("Error processing CSV file: %v", err)
	}

	resultsCh <- summary
}
