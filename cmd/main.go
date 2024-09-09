package main

import (
	"GoCamera/internal/organiser"
	"GoCamera/internal/transfer"
	"fmt"
	"log"
	"os"
	"sync"
)

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

	processor := transfer.NewFileProcessor(sourceDir, destinationDir)
	organiser := organiser.NewFileOrganiser(destinationDir)

	var wg sync.WaitGroup
	var transferDone = make(chan bool)

	// Start file transfer in a goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := processor.ProcessDirectory(); err != nil {
			fmt.Println("Error in file transfer:", err)
			close(transferDone)
			return
		}
		transferDone <- true
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		<-transferDone
		if err := organiser.ProcessFiles(); err != nil {
			fmt.Println("Error in file organization:", err)
			return
		}
	}()

	wg.Wait()
	fmt.Println("Processing complete.")
}
