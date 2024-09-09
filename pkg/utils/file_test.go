package utils

import (
	"testing"
)

func TestIsSupported(t *testing.T) {
	tests := []struct {
		filename    string
		isSupported bool
	}{
		// Valid cases
		{"photo.jpg", true},
		{"image.arw", true},
		{"PHOTO.JPG", true},
		{"image.ARW", true},

		// Invalid cases
		{"image.png", false},
		{"photo.jpeg", false},
		{"document.txt", false},
		{"", false},
		{"no_extension", false},
		{"folder/photo.jpg/", false},
	}

	for _, test := range tests {
		result := IsSupportedFile(test.filename)
		if result != test.isSupported {
			t.Errorf("IsSupportedFile(%q) = %v, expected %v", test.filename, result, test.isSupported)
		}
	}
}
