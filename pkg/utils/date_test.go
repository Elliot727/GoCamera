package utils

import (
	"testing"
)

func TestExtractDateFromFilename(t *testing.T) {
	tests := []struct {
		filename     string
		expectedDate string
		expectError  bool
	}{
		{"2024_09_09_13_05_16.jpg", "2024/09/09", false},
		{"2024_09_09_13_05_16.png", "2024/09/09", false},
		{"invalid_date.jpg", "", true},
		{"2024_09_09_13_05.jpg", "", true},
	}

	for _, test := range tests {
		date, err := ExtractDateFromFilename(test.filename)
		if (err != nil) != test.expectError {
			t.Errorf("ExtractDateFromFilename(%q) error = %v, expectError %v", test.filename, err, test.expectError)
			continue
		}
		if date != test.expectedDate {
			t.Errorf("ExtractDateFromFilename(%q) = %q, expected %q", test.filename, date, test.expectedDate)
		}
	}
}

func TestSanitizeFilename(t *testing.T) {
	tests := []struct {
		input          string
		expectedOutput string
	}{
		{
			input:          `2024:09:09 13:05:16`,
			expectedOutput: "2024_09_09_13_05_16",
		},
		{
			input:          `2022-01-01 10:20:30`,
			expectedOutput: "2022_01_01_10_20_30",
		},
		{
			input:          `1999:12:31 23:59:59"`,
			expectedOutput: "1999_12_31_23_59_59",
		},
	}

	for _, test := range tests {
		output := SanitizeFilename(test.input)
		if output != test.expectedOutput {
			t.Errorf("sanitizeFilename(%q) = %q, expected %q", test.input, output, test.expectedOutput)
		}
	}
}
