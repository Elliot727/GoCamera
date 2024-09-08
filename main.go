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

// FileProcessor handles the processing of image files.
type FileProcessor struct {
	dir string
}

// NewFileProcessor creates a new FileProcessor.
func NewFileProcessor(dir string) *FileProcessor {
	return &FileProcessor{dir: dir}
}

// ProcessDirectory processes the directory and its files.
func (fp *FileProcessor) ProcessDirectory() error {
	startTime := time.Now()

	err := filepath.Walk(fp.dir, fp.processFile)
	if err != nil {
		return err
	}

	totalTime := time.Since(startTime)
	fmt.Printf("Total processing time: %v\n", totalTime)

	return nil
}

// processFile is called for each file or directory during the walk.
func (fp *FileProcessor) processFile(path string, info os.FileInfo, err error) error {
	if err != nil {
		return handleFileError(path, err)
	}

	if shouldSkip(info) {
		return nil
	}

	if !isSupportedFile(path) {
		log.Printf("Skipping non-JPG/ARW file: %s\n", path)
		return nil
	}

	startTime := time.Now()

	err = fp.processImageFile(path)
	if err != nil {
		return err
	}

	fileProcessingTime := time.Since(startTime) // Calculate duration for the file
	fmt.Printf("Processing time for file %s: %v\n", path, fileProcessingTime)

	return nil
}

// handleFileError handles errors encountered while processing files.
func handleFileError(path string, err error) error {
	if os.IsPermission(err) {
		log.Printf("Skipping restricted file: %s (Permission denied)\n", path)
		return nil
	}
	return err
}

// shouldSkip determines if a file or directory should be skipped.
func shouldSkip(info os.FileInfo) bool {
	return info.IsDir() && strings.HasPrefix(info.Name(), ".") || strings.HasPrefix(info.Name(), ".")
}

// isSupportedFile checks if the file has a supported extension.
func isSupportedFile(path string) bool {
	ext := filepath.Ext(path)
	return ext == ".JPG" || ext == ".ARW"
}

// processImageFile processes an image file and prints its metadata.
func (fp *FileProcessor) processImageFile(path string) error {
	log.Printf("Processing file: %s\n", path)

	f, err := os.Open(path)
	if err != nil {
		log.Printf("Error opening file: %s\n", path)
		return nil
	}
	defer f.Close()

	if filepath.Ext(path) == ".JPG" || filepath.Ext(path) == ".ARW" {
		return processJPGFile(f)
	}

	log.Printf("Raw file detected, handling not implemented: %s\n", path)
	return nil
}

// processJPGFile processes a JPG file and prints its EXIF data.
func processJPGFile(f *os.File) error {
	x, err := exif.Decode(f)
	if err != nil {
		log.Printf("Error decoding EXIF data: %s\n", f.Name())
		return nil
	}

	fi, err := f.Stat()
	if err != nil {
		log.Printf("Error getting file size: %s\n", f.Name())
		return nil
	}

	fmt.Printf("The file is %.2f MB long", float64(fi.Size())/(1024*1024))

	if camModel, err := x.Get(exif.Model); err == nil {
		fmt.Printf("Camera model: %s\n", camModel.String())
	} else {
		log.Printf("Error getting camera model: %s\n", f.Name())
	}

	if date, err := x.Get(exif.DateTime); err == nil {
		fmt.Printf("Date: %s\n", date.String())
	} else {
		log.Printf("Error getting date: %s\n", f.Name())
	}

	fmt.Println()
	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <directory>")
		os.Exit(1)
	}

	dir := os.Args[1]
	processor := NewFileProcessor(dir)

	if err := processor.ProcessDirectory(); err != nil {
		log.Fatal(err)
	}
}
