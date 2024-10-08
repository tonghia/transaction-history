package parser

import (
	"fmt"
	"strconv"
	"time"
)

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
