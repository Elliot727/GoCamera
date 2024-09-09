# File Organiser

A Go-based tool for organizing files based on the date extracted from their filenames. This project processes files in a specified directory, categorizes them by date, and moves them into a structured folder hierarchy.

## Features

- **File Processing**: Scans a source directory for files.
- **Date Extraction**: Extracts date information from filenames.
- **File Organization**: Moves files into a directory structure based on the extracted date.
- **Concurrent Processing**: Utilizes goroutines to speed up file processing.

## Prerequisites

- Go 1.18 or higher
- Access to the source directory with files to process
- Destination directory (optional, used for copying files if applicable)

## Installation

1. **Clone the repository:**

    ```sh
    git clone https://github.com/yourusername/file-organiser.git
    cd file-organiser
    ```

2. **Ensure Go is installed and set up on your system.**

## Usage

To run the application:

```sh
go run main/*.go <source_directory> <destination_directory>
```

- `<source_directory>`: The path to the directory containing the files to be processed.
- `<destination_directory>`: The path to the directory where files will be organized (if applicable).

### Example

```sh
go run main/*.go "/Volumes/Untitled/DCIM/100MSDCF" "/Volumes/My Passport/Photos"
```

## Code Overview

- **`main.go`**: Entry point of the application. Handles argument parsing and coordinates file processing and organization.
- **`file_processor.go`**: Contains logic for copying files from the source directory to the destination.
- **`file_organiser.go`**: Contains logic for organizing files into directories based on extracted dates from filenames.

## Contributing

Contributions are welcome! Please fork the repository and submit a pull request with your proposed changes.
