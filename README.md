# File Organiser

A Go-based tool for organizing and transferring files. This project includes functionality for transferring files from a source directory to a destination directory and organizing them based on extracted date information from filenames.

## Features

- **File Transfer**: Moves files from a source directory to a destination directory.
- **File Organization**: Scans a directory, extracts date information from filenames, and organizes files into a structured folder hierarchy.
- **Concurrent Processing**: Utilizes goroutines and semaphore patterns to handle multiple files efficiently.
- **Server Mode**: (To be implemented) Placeholder for running a server if required.

## Prerequisites

- Go 1.18 or higher
- Access to the source and destination directories

## Installation

1. **Clone the repository:**

   ```sh
   git clone https://github.com/elliot727/GoCamera.git
   cd GoCamera
   ```

2. **Ensure Go is installed and set up on your system.**

## Usage

To run the application, use the following command:

```sh
cd cmd
go run main/*.go [options]
```

### Options

- `--server`: Run the server mode (currently not implemented).
- `--transfer`: Run the file transfer mode.
- `--organise`: Run the file organisation mode.
- `--source`: Source directory for file transfer mode.
- `--dest`: Destination directory for file transfer, organisation, and server modes.
- `--port`: Port for server mode (default is `8080`).

### Examples

1. **Run file transfer:**

   ```sh
   go run main/*.go --transfer --source "/path/to/source" --dest "/path/to/destination"
   ```

2. **Run file organization:**

   ```sh
   go run main/*.go --organise --dest "/path/to/destination"
   ```

3. **Run server mode (if implemented):**

   ```sh
   go run main/*.go --server --dest "/path/to/destination" --port "8080"
   ```

## Code Overview

- **`main.go`**: Entry point of the application. Handles argument parsing and coordinates file processing and organization.
- **`internal/transfer/transfer.go`**: Contains logic for copying files from the source directory to the destination.
- **`internal/organiser/organiser.go`**: Contains logic for organizing files into directories based on extracted dates from filenames.

### Directory Structure

- **`main.go`**: Entry point for the application.
- **`internal/transfer/transfer.go`**: File transfer logic.
- **`internal/organiser/organiser.go`**: File organization logic.

## Contributing

Contributions are welcome! Please fork the repository and submit a pull request with your proposed changes.
