package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/rwcarlsen/goexif/exif"
)

type FileProcessor struct {
	sourceDir      string
	destinationDir string
	concurrency    int
}

func NewFileProcessor(sourceDir, destinationDir string, concurrency int) *FileProcessor {
	return &FileProcessor{sourceDir: sourceDir, destinationDir: destinationDir, concurrency: concurrency}
}

func (fp *FileProcessor) ProcessDirectory() error {
	startTime := time.Now()
	err := filepath.Walk(fp.sourceDir, fp.processFile)
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

	err = fp.processImageFile(path)
	if err != nil {
		return err
	}
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
	f, err := os.Open(path)
	if err != nil {
		log.Printf("Error opening file: %s\n", path)
		return nil
	}
	defer f.Close()

	x, err := exif.Decode(f)
	if err != nil {
		log.Printf("Error decoding EXIF data: %s\n", path)
		return nil
	}

	if _, err := f.Seek(0, io.SeekStart); err != nil {
		log.Printf("Error seeking file: %s\n", path)
		return nil
	}

	return fp.renameAndCopyFileWithDate(x, f, path)
}

func (fp *FileProcessor) renameAndCopyFileWithDate(x *exif.Exif, file *os.File, path string) error {
	date, err := x.Get(exif.DateTime)
	if err != nil {
		log.Printf("Error getting date: %s\n", path)
		return nil
	}

	dateStr := date.String()
	newName := sanitizeFilename(dateStr) + filepath.Ext(path)
	newPath := filepath.Join(fp.destinationDir, newName)

	err = copyFileWithBuffer(file, newPath)
	if err != nil {
		log.Printf("Error copying file: %s\n", err)
		return err
	}
	err = os.Remove(path)
	if err != nil {
		log.Printf("Error deleting file: %s\n", err)
		return err
	}

	return nil
}

func copyFileWithBuffer(src *os.File, dst string) error {
	startTime := time.Now()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	// Define a reusable buffer of 32KB
	buffer := make([]byte, 32*1024)

	// Use io.CopyBuffer to copy from src to dst
	_, err = io.CopyBuffer(out, src, buffer)
	if err != nil {
		return err
	}

	// Ensure everything is written to disk
	err = out.Sync()
	if err != nil {
		return err
	}

	log.Printf("File copied in %v\n", time.Since(startTime))
	return nil
}

func sanitizeFilename(dateStr string) string {
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
	if len(os.Args) < 3 {
		log.Println("Usage: go run main.go <source_directory> <destination_directory>")
		os.Exit(1)
	}

	sourceDir := os.Args[1]
	destinationDir := os.Args[2]

	if err := os.MkdirAll(destinationDir, os.ModePerm); err != nil {
		log.Fatalf("Error creating destination directory: %v\n", err)
	}

	processor := NewFileProcessor(sourceDir, destinationDir, 4)
	if err := processor.ProcessDirectory(); err != nil {
		log.Fatal(err)
	}
}
