package main

import (
	"GoCamera/internal/api"
	"GoCamera/internal/organiser"
	"GoCamera/internal/transfer"
	"GoCamera/pkg/utils"
	"fmt"
	"log"
	"os"
	"os/exec"
	"sync"
)

func sendNotification(title, message string) error {
	script := fmt.Sprintf(`display notification "%s" with title "%s"`, message, title)
	_, err := exec.Command("osascript", "-e", script).Output()
	return err
}

func main() {
	cfg, err := utils.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	if !cfg.Server && !cfg.Transfer && !cfg.Organise {
		log.Fatal("Please specify at least one mode to run: server, transfer, or organise")
	}

	if cfg.Transfer && (cfg.Source == "" || cfg.Dest == "") {
		log.Fatal("Source and destination directories are required for transfer mode")
	}
	if (cfg.Organise || cfg.Server) && cfg.Dest == "" {
		log.Fatal("Destination directory is required for organise and server modes")
	}

	var wg sync.WaitGroup

	// Run server if selected
	if cfg.Server {
		runServerComponent(cfg.Dest, cfg.Port)

	}

	// Run transfer if selected
	if cfg.Transfer {
		wg.Add(1)
		go func() {
			defer wg.Done()
			runTransferComponent(cfg.Source, cfg.Dest)
		}()
	}

	// Run organise if selected
	if cfg.Organise {
		wg.Add(1)
		go func() {
			defer wg.Done()
			runOrganiseComponent(cfg.Dest)
		}()
	}

	wg.Wait()

	// Notify completion
	err = sendNotification("Process Complete", "All selected operations have been completed successfully.")
	if err != nil {
		fmt.Printf("Error sending notification: %v\n", err)
	}

	fmt.Println("All selected operations completed.")
}

func runServerComponent(dir, port string) {
	server := api.NewServer(port)

	if err := server.Start(); err != nil {
		log.Fatalf("Error starting server: %v\n", err)
	}
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
