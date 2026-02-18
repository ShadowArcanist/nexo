//go:build !tinygo

package main

import (
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/gen2brain/webp"
)

// Supported image extensions
var supportedExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".gif":  true,
}

type conversionResult struct {
	path    string
	success bool
	err     error
}

func main() {
	startTime := time.Now()

	// Default directory name
	inputDir := "converter"

	// Allow directory override via command line
	if len(os.Args) > 1 {
		inputDir = os.Args[1]
	}

	// Ensure directory exists
	if _, err := os.Stat(inputDir); os.IsNotExist(err) {
		fmt.Printf("Directory '%s' not found\n", inputDir)
		fmt.Println("Create it and put your images inside")
		os.Exit(1)
	}

	// Collect all image files
	var imageFiles []string

	err := filepath.WalkDir(inputDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		if supportedExtensions[ext] {
			imageFiles = append(imageFiles, path)
		}

		return nil
	})

	if err != nil {
		fmt.Printf("Error scanning directory: %v\n", err)
		os.Exit(1)
	}

	if len(imageFiles) == 0 {
		fmt.Printf("No supported images found in '%s'\n", inputDir)
		fmt.Println("Supported formats: JPG, JPEG, PNG, GIF")
		os.Exit(0)
	}

	fmt.Printf("Found %d image(s) to convert\n", len(imageFiles))

	// Process images concurrently using worker pool
	results := processImages(imageFiles)

	// Print results
	successCount := 0
	failCount := 0
	var failedImages []string

	for _, result := range results {
		if result.success {
			successCount++
		} else {
			failCount++
			failedImages = append(failedImages, filepath.Base(result.path))
		}
	}

	elapsed := time.Since(startTime)

	fmt.Println("Conversion completed!")
	fmt.Printf("Successful: %d\n", successCount)
	if failCount > 0 {
		fmt.Printf("Failure: %d\n", failCount)
		fmt.Printf("Failed Images: %s\n", strings.Join(failedImages, ", "))
	}
	fmt.Printf("Time taken: %s\n", formatDuration(elapsed))
}

func formatDuration(d time.Duration) string {
	if d < time.Second {
		return fmt.Sprintf("%dms", d.Milliseconds())
	}
	return fmt.Sprintf("%.2fs", d.Seconds())
}

func processImages(files []string) []conversionResult {
	// Use number of CPUs for worker count
	numWorkers := max(1, min(len(files), 8))

	results := make([]conversionResult, 0, len(files))
	resultsChan := make(chan conversionResult, len(files))

	// Create work channel
	workChan := make(chan string, len(files))

	// Start workers
	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for path := range workChan {
				result := convertImage(path)
				resultsChan <- result
			}
		}()
	}

	// Send work
	go func() {
		for _, file := range files {
			workChan <- file
		}
		close(workChan)
	}()

	// Wait for workers and close results channel
	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	// Collect results
	for result := range resultsChan {
		results = append(results, result)
	}

	return results
}

func convertImage(inputPath string) conversionResult {
	// Open input file
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return conversionResult{path: inputPath, success: false, err: err}
	}
	defer inputFile.Close()

	// Decode image
	var img image.Image
	ext := strings.ToLower(filepath.Ext(inputPath))

	switch ext {
	case ".jpg", ".jpeg":
		img, err = jpeg.Decode(inputFile)
	case ".png":
		img, err = png.Decode(inputFile)
	case ".gif":
		img, err = gif.Decode(inputFile)
	default:
		err = fmt.Errorf("unsupported format: %s", ext)
	}

	if err != nil {
		return conversionResult{path: inputPath, success: false, err: fmt.Errorf("decode failed: %w", err)}
	}

	inputFile.Close()

	// Determine output path
	outputPath := inputPath[:len(inputPath)-len(ext)] + ".webp"

	// Create output file
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return conversionResult{path: inputPath, success: false, err: fmt.Errorf("create output failed: %w", err)}
	}
	defer outputFile.Close()

	// Encode to WebP with quality 85 (good balance)
	op := webp.Options{Quality: 85}
	if err := webp.Encode(outputFile, img, op); err != nil {
		os.Remove(outputPath)
		return conversionResult{path: inputPath, success: false, err: fmt.Errorf("encode failed: %w", err)}
	}

	outputFile.Close()

	// Remove original file
	if err := os.Remove(inputPath); err != nil {
		fmt.Printf("Could not remove original: %s\n", inputPath)
	}

	return conversionResult{path: inputPath, success: true}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
