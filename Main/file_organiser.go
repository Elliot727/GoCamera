package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type FileOrganiser struct {
	sourceDir string
}

func NewFileOrganiser(sourceDir string) *FileOrganiser {
	return &FileOrganiser{sourceDir: sourceDir}
}

func (fo *FileOrganiser) ProcessFiles() error {
	startTime := time.Now()
	var wg sync.WaitGroup
	sem := make(chan struct{}, 12)

	fileCount := 0
	dirCount := 0

	err := filepath.Walk(fo.sourceDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		fmt.Printf("Processing path: %s\n", path)

		if info.IsDir() {
			dirCount++
			return nil
		}

		if !isSupportedFile(path) {
			log.Printf("Skipping unsupported file: %s\n", path)
			return nil
		}

		wg.Add(1)
		sem <- struct{}{}

		go func(p string, file fs.FileInfo) {
			defer wg.Done()
			defer func() { <-sem }()

			startExtractTime := time.Now()
			date, err := extractDateFromFilename(p)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			fmt.Println("Extracted Date:", date)
			fmt.Println(p)
			extractDuration := time.Since(startExtractTime)

			// Move file
			startMoveTime := time.Now()
			if err := moveFileToDateDir(p, date, fo.sourceDir); err != nil {
				fmt.Println("Error moving file:", err)
				return
			}
			moveDuration := time.Since(startMoveTime)

			fileCount++
			fmt.Printf("Date extraction took: %v\n", extractDuration)
			fmt.Printf("File move took: %v\n", moveDuration)
		}(path, info)

		return nil
	})

	if err != nil {
		return err
	}

	wg.Wait()

	duration := time.Since(startTime)

	fmt.Printf("Total directories: %d\n", dirCount)
	fmt.Printf("Total files: %d\n", fileCount)
	fmt.Printf("Total processing time: %v\n", duration)

	return nil
}

func extractDateFromFilename(filename string) (string, error) {
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

func createDirectories(baseDir, year, month, day string) error {
	yearDir := filepath.Join(baseDir, year)
	monthDir := filepath.Join(yearDir, month)
	dayDir := filepath.Join(monthDir, day)

	if err := os.MkdirAll(dayDir, os.ModePerm); err != nil {
		return fmt.Errorf("error creating directories: %w", err)
	}

	return nil
}

func moveFileToDateDir(filePath, formattedDate, baseDir string) error {
	parts := strings.Split(formattedDate, "/")
	if len(parts) != 3 {
		return fmt.Errorf("date format is incorrect: %s", formattedDate)
	}
	year, month, day := parts[0], parts[1], parts[2]

	if err := createDirectories(baseDir, year, month, day); err != nil {
		return fmt.Errorf("error creating dir: %w", err)
	}

	baseName := filepath.Base(filePath)
	newPath := filepath.Join(baseDir, year, month, day, baseName)

	if err := os.Rename(filePath, newPath); err != nil {
		return fmt.Errorf("error moving file: %w", err)
	}

	return nil
}
