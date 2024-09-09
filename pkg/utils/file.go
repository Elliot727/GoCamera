package utils

import (
	"fmt"
	"path/filepath"
	"strings"
)

func IsSupportedFile(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	return ext == ".jpg" || ext == ".arw"
}

func SanitizeFilename(dateStr string) string {
	sanitized := strings.NewReplacer(`"`, "", ":", "_", " ", "_").Replace(dateStr)

	// Extract components and format as yyyy_mm_dd_hh_mm_ss
	return fmt.Sprintf("%s_%s_%s_%s_%s_%s",
		sanitized[:4],    // year
		sanitized[5:7],   // month
		sanitized[8:10],  // day
		sanitized[11:13], // hour
		sanitized[14:16], // minute
		sanitized[17:19]) // second
}
