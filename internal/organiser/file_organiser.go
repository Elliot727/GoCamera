package organiser

import (
	"GoCamera/pkg/utils"
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

	err := filepath.WalkDir(fo.sourceDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			if os.IsPermission(err) {
				log.Printf("Skipping file due to permissions: %s\n", path)
				return nil
			}
			return err
		}

		fmt.Println("Processing file:", d.Name())
		if d.Name() == ".Trashes" || d.Name()[0] == '.' {
			log.Printf("Skipping hidden directory or file: %s\n", path)
			return nil
		}

		if d.IsDir() {
			dirCount++
			return nil
		}

		if !utils.IsSupportedFile(path) {
			log.Printf("Skipping unsupported file: %s\n", path)
			return nil
		}

		wg.Add(1)
		sem <- struct{}{}

		go func(p string, entry fs.DirEntry) {
			defer wg.Done()
			defer func() { <-sem }()

			startExtractTime := time.Now()
			date, err := utils.ExtractDateFromFilename(p)
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
		}(path, d)

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
	photosDir := filepath.Join(baseDir, "photos")

	parts := strings.Split(formattedDate, "/")
	if len(parts) != 3 {
		return fmt.Errorf("date format is incorrect: %s", formattedDate)
	}
	year, month, day := parts[0], parts[1], parts[2]

	if err := createDirectories(photosDir, year, month, day); err != nil {
		return fmt.Errorf("error creating dir: %w", err)
	}

	baseName := filepath.Base(filePath)
	newPath := filepath.Join(photosDir, year, month, day, baseName)

	if err := os.Rename(filePath, newPath); err != nil {
		return fmt.Errorf("error moving file: %w", err)
	}

	return nil
}
