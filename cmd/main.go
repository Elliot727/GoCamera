package main

import (
	"GoCamera/internal/organiser"
	"GoCamera/internal/transfer"
	"flag"
	"fmt"
	"log"
	"os"
	"sync"
)

func main() {
	// Define flags
	runServer := flag.Bool("server", false, "Run the server")
	runTransfer := flag.Bool("transfer", false, "Run the file transfer")
	runOrganise := flag.Bool("organise", false, "Run the file organiser")
	sourceDir := flag.String("source", "", "Source directory for transfer mode")
	destDir := flag.String("dest", "", "Destination directory for transfer, organise, and server modes")
	port := flag.String("port", "8080", "Port for server mode")

	flag.Parse()

	// Validate that at least one mode is selected
	if !*runServer && !*runTransfer && !*runOrganise {
		log.Fatal("Please specify at least one mode to run: --server, --transfer, or --organise")
	}

	// Validate required directories
	if *runTransfer && (*sourceDir == "" || *destDir == "") {
		log.Fatal("Source and destination directories are required for transfer mode")
	}
	if (*runOrganise || *runServer) && *destDir == "" {
		log.Fatal("Destination directory is required for organise and server modes")
	}

	var wg sync.WaitGroup

	// Run server if selected
	if *runServer {
		wg.Add(1)
		go func() {
			defer wg.Done()
			runServerComponent(*destDir, *port)
		}()
	}

	// Run transfer if selected
	if *runTransfer {
		wg.Add(1)
		go func() {
			defer wg.Done()
			runTransferComponent(*sourceDir, *destDir)
		}()
	}

	// Run organise if selected
	if *runOrganise {
		wg.Add(1)
		go func() {
			defer wg.Done()
			runOrganiseComponent(*destDir)
		}()
	}

	// Wait for all selected components to complete
	wg.Wait()
	fmt.Println("All selected operations completed.")
}

func runServerComponent(dir, port string) {

}

func runTransferComponent(src, dest string) {
	if err := os.MkdirAll(dest, os.ModePerm); err != nil {
		log.Printf("Error creating destination directory: %v\n", err)
		return
	}

	processor := transfer.NewFileProcessor(src, dest)
	if err := processor.ProcessDirectory(); err != nil {
		log.Printf("Error in file transfer: %v\n", err)
		return
	}
	fmt.Println("Transfer complete.")
}

func runOrganiseComponent(dir string) {
	organiser := organiser.NewFileOrganiser(dir)
	if err := organiser.ProcessFiles(); err != nil {
		log.Printf("Error in file organization: %v\n", err)
		return
	}
	fmt.Println("Organisation complete.")
}
