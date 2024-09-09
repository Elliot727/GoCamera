package utils

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"
)

func ExtractDateFromFilename(filename string) (string, error) {
	base := filepath.Base(filename)

	dateStr := strings.TrimSuffix(base, filepath.Ext(base))
	if len(dateStr) != 19 {
		return "", fmt.Errorf("filename format is incorrect: %s", dateStr)
	}

	layout := "2006_01_02_15_04_05"

	date, err := time.Parse(layout, dateStr)
	if err != nil {
		return "", fmt.Errorf("error parsing date: %w", err)
	}

	formattedDate := date.Format("2006/01/02")
	return formattedDate, nil
}
