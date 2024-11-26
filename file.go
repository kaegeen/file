package main

import (
	"fmt"
	"os"
	"sync"
	"time"
)

// File structure to represent a download task
type File struct {
	Name     string
	Size     int // in MB
	Progress int // in MB downloaded
}

// Function to simulate file download
func downloadFile(file *File, wg *sync.WaitGroup, ch chan string) {
	defer wg.Done()

	for file.Progress < file.Size {
		// Simulate downloading by incrementing progress
		time.Sleep(time.Second * 1) // Simulating download delay
		file.Progress++

		// Report download progress
		ch <- fmt.Sprintf("Downloading %s: %d/%d MB", file.Name, file.Progress, file.Size)
	}

	// Notify when download is complete
	ch <- fmt.Sprintf("Download of %s completed!", file.Name)
}

func main() {
	// List of files to download
	files := []File{
		{Name: "File1", Size: 10},
		{Name: "File2", Size: 5},
		{Name: "File3", Size: 7},
	}

	// WaitGroup to synchronize goroutines
	var wg sync.WaitGroup

	// Channel to handle progress updates
	progressChannel := make(chan string)

	// Start downloading files concurrently
	for i := range files {
		wg.Add(1)
		go downloadFile(&files[i], &wg, progressChannel)
	}

	// Writing progress updates to file and printing to console
	go func() {
		for msg := range progressChannel {
			fmt.Println(msg)         // Print progress to console
			writeProgressToFile(msg) // Write to a log file
		}
	}()

	// Wait for all downloads to complete
	wg.Wait()
	close(progressChannel)
	fmt.Println("All downloads completed.")
}

// Function to write progress to a file
func writeProgressToFile(progress string) {
	// Open the file in append mode, or create if it doesn't exist
	file, err := os.OpenFile("download_progress.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Write progress message to the file
	if _, err := file.WriteString(progress + "\n"); err != nil {
		fmt.Println("Error writing to file:", err)
	}
}
