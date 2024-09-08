package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/rwcarlsen/goexif/exif"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <directory>")
		os.Exit(1)
	}

	dir := os.Args[1]

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			if os.IsPermission(err) {
				log.Printf("Skipping restricted file: %s (Permission denied)\n", path)
				return nil
			}
			return err
		}

		if info.IsDir() && strings.HasPrefix(info.Name(), ".") {
			return filepath.SkipDir
		}
		if strings.HasPrefix(info.Name(), ".") {
			return nil
		}

		if filepath.Ext(path) != ".JPG" && filepath.Ext(path) != ".ARW" {
			log.Printf("Skipping non-JPG/ARW file: %s\n", path)
			return nil
		}

		log.Printf("Processing file: %s\n", path)

		f, err := os.Open(path)
		if err != nil {
			log.Printf("Error opening file: %s\n", path)
			return nil
		}
		defer f.Close()

		if filepath.Ext(path) == ".JPG" {
			x, err := exif.Decode(f)
			if err != nil {
				log.Printf("Error decoding EXIF data: %s\n", path)
				return nil
			}

			camModel, err := x.Get(exif.Model)
			if err != nil {
				log.Printf("Error getting camera model: %s\n", path)
			} else {
				fmt.Printf("Camera model: %s\n", camModel.String())
			}

			date, err := x.Get(exif.DateTime)
			if err != nil {
				log.Printf("Error getting date: %s\n", path)
			} else {
				fmt.Printf("Date: %s\n", date.String())
			}
		} else {
			log.Printf("Raw file detected, handling not implemented: %s\n", path)
		}

		fmt.Println()

		return nil
	})

	if err != nil {
		log.Fatal(err)
	}
}
