package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/rwcarlsen/goexif/exif"
)

type FileProcessor struct {
	dir string
}

func NewFileProcessor(dir string) *FileProcessor {
	return &FileProcessor{dir: dir}
}

func (fp *FileProcessor) ProcessDirectory() error {
	startTime := time.Now()
	err := filepath.Walk(fp.dir, fp.processFile)
	if err != nil {
		return err
	}
	log.Printf("Total processing time: %v\n", time.Since(startTime))
	return nil
}

func (fp *FileProcessor) processFile(path string, info os.FileInfo, err error) error {
	if err != nil {
		return handleFileError(path, err)
	}
	if shouldSkip(info) {
		return nil
	}
	if !isSupportedFile(path) {
		log.Printf("Skipping unsupported file: %s\n", path)
		return nil
	}

	startTime := time.Now()
	err = fp.processImageFile(path)
	if err != nil {
		return err
	}
	log.Printf("Processing time for file %s: %v\n", path, time.Since(startTime))
	return nil
}

func handleFileError(path string, err error) error {
	if os.IsPermission(err) {
		log.Printf("Skipping restricted file: %s (Permission denied)\n", path)
		return nil
	}
	return err
}

func shouldSkip(info os.FileInfo) bool {
	return info.IsDir() && strings.HasPrefix(info.Name(), ".")
}

func isSupportedFile(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	return ext == ".jpg" || ext == ".arw"
}

func (fp *FileProcessor) processImageFile(path string) error {
	log.Printf("Processing file: %s\n", path)

	f, err := os.Open(path)
	if err != nil {
		log.Printf("Error opening file: %s\n", path)
		return nil
	}
	defer f.Close()

	x, err := exif.Decode(f)
	if err != nil {
		log.Printf("Error decoding EXIF data: %s\n", path)
		return err
	}

	printFileInfo(f)
	printCameraModel(x)
	return renameFileWithDate(x, path)
}

func printFileInfo(f *os.File) {
	fi, err := f.Stat()
	if err == nil {
		log.Printf("The file is %.2f MB long", float64(fi.Size())/(1024*1024))
	}
}

func printCameraModel(x *exif.Exif) {
	if model, err := x.Get(exif.Model); err == nil {
		log.Printf("Camera model: %s", model)
	}
}

func renameFileWithDate(x *exif.Exif, path string) error {
	date, err := x.Get(exif.DateTime)
	if err != nil {
		log.Printf("Error getting date: %s\n", path)
		return nil
	}

	dateStr := date.String()
	newName := sanitizeFilename(dateStr) + filepath.Ext(path)
	newPath := filepath.Join(filepath.Dir(path), newName)

	if err := os.Rename(path, newPath); err != nil {
		log.Printf("Error renaming file: %s\n", err)
		return err
	}
	log.Printf("File renamed to: %s\n", newName)
	return nil
}

func sanitizeFilename(dateStr string) string {
	// Remove any quotes and replace colons and spaces with underscores
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

func main() {
	if len(os.Args) < 2 {
		log.Println("Usage: go run main.go <directory>")
		os.Exit(1)
	}

	processor := NewFileProcessor(os.Args[1])
	if err := processor.ProcessDirectory(); err != nil {
		log.Fatal(err)
	}
}
