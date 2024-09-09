package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/rwcarlsen/goexif/exif"
)

type FileProcessor struct {
	sourceDir      string
	destinationDir string
}

func NewFileProcessor(sourceDir, destinationDir string) *FileProcessor {
	return &FileProcessor{sourceDir: sourceDir, destinationDir: destinationDir}
}

func (fp *FileProcessor) ProcessDirectory() error {
	startTime := time.Now()
	var wg sync.WaitGroup
	sem := make(chan struct{}, runtime.NumCPU()) // Semaphore for limiting concurrency

	err := filepath.Walk(fp.sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if shouldSkip(info) {
			return nil
		}
		if !isSupportedFile(path) {
			log.Printf("Skipping unsupported file: %s\n", path)
			return nil
		}

		wg.Add(1)
		sem <- struct{}{} // Acquire semaphore
		go func() {
			defer wg.Done()
			defer func() { <-sem }() // Release semaphore
			if err := fp.processImageFile(path); err != nil {
				log.Printf("Error processing file %s: %v\n", path, err)
			}
		}()

		return nil
	})
	if err != nil {
		return err
	}

	wg.Wait()

	log.Printf("Total processing time: %v\n", time.Since(startTime))
	return nil
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

	processor := NewFileProcessor(sourceDir, destinationDir)
	if err := processor.ProcessDirectory(); err != nil {
		log.Fatal(err)
	}
}
