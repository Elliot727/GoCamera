package transfer

import (
	"GoCamera/pkg/utils"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
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
	sem := make(chan struct{}, 12) // Semaphore for limiting concurrency

	err := filepath.WalkDir(fp.sourceDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if shouldSkip(d) {
			return nil
		}
		if d.IsDir() || !utils.IsSupportedFile(path) {
			log.Printf("Skipping unsupported or directory: %s\n", path)
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

func shouldSkip(d fs.DirEntry) bool {
	name := d.Name()
	return d.IsDir() && (strings.HasPrefix(name, ".") || name == ".Trashes")
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

	if err := os.MkdirAll(fp.destinationDir, os.ModePerm); err != nil {
		log.Printf("Error creating directory: %s\n", fp.destinationDir)
		return err
	}

	newName := utils.SanitizeFilename(dateStr) + filepath.Ext(path)
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

	buffer := make([]byte, 32*1024)

	_, err = io.CopyBuffer(out, src, buffer)
	if err != nil {
		return err
	}

	err = out.Sync()
	if err != nil {
		return err
	}

	log.Printf("File copied in %v\n", time.Since(startTime))
	return nil
}
