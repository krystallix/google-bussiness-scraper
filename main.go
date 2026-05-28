package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const banner = `
+--------------------------------------------------+
|       Google Business Scraper -- Go Edition      |
|   Extract business data from Google Maps         |
+--------------------------------------------------+
`

func main() {
	fmt.Print(banner)

	reader := bufio.NewReader(os.Stdin)

	// 1. Keyword
	var keyword string
	for keyword == "" {
		fmt.Print("Search keyword (e.g. coffee shop Bandung): ")
		line, _ := reader.ReadString('\n')
		keyword = strings.TrimSpace(line)
		if keyword == "" {
			fmt.Println("Keyword cannot be empty, please try again.")
		}
	}

	// 2. Limit
	limit := 50
	fmt.Printf("\nMax results to collect (default %d): ", limit)
	line, _ := reader.ReadString('\n')
	line = strings.TrimSpace(line)
	if line != "" {
		if parsed, err := strconv.Atoi(line); err == nil && parsed > 0 {
			limit = parsed
		} else {
			fmt.Printf("Invalid input, using default: %d\n", limit)
		}
	}

	// 3. Headful
	headful := false
	fmt.Print("\nShow browser window? (y/N): ")
	line, _ = reader.ReadString('\n')
	line = strings.ToLower(strings.TrimSpace(line))
	if line == "y" || line == "yes" {
		headful = true
	}

	// 4. Delay
	delay := 1.5
	fmt.Print("\nDelay between pages in seconds (default 1.5): ")
	line, _ = reader.ReadString('\n')
	line = strings.TrimSpace(line)
	if line != "" {
		if parsed, err := strconv.ParseFloat(line, 64); err == nil && parsed >= 0 {
			delay = parsed
		} else {
			fmt.Printf("Invalid input, using default: %.1fs\n", delay)
		}
	}

	// Summary
	fmt.Printf("\n%s\n", strings.Repeat("-", 52))
	fmt.Printf("  Keyword : %s\n", keyword)
	fmt.Printf("  Limit   : %d results\n", limit)
	fmt.Printf("  Mode    : %s\n", func() string {
		if headful {
			return "headful (visible browser)"
		}
		return "headless (background)"
	}())
	fmt.Printf("  Delay   : %.1fs\n", delay)
	fmt.Printf("%s\n", strings.Repeat("-", 52))

	csvPath, xlsxPath := outputPaths("", keyword)

	startTime := time.Now()
	data, err := ScrapeGoogleMaps(keyword, limit, headful, delay)
	elapsed := time.Since(startTime)

	if err != nil {
		fmt.Fprintf(os.Stderr, "\nERROR: %v\n", err)
		os.Exit(1)
	}
	if len(data) == 0 {
		fmt.Fprintln(os.Stderr, "\nNo results found.")
		os.Exit(1)
	}

	if err := saveCSV(csvPath, data); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to save CSV: %v\n", err)
		os.Exit(1)
	}
	if err := saveXLSX(xlsxPath, data); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to save Excel: %v\n", err)
		os.Exit(1)
	}

	absCsv, _ := filepath.Abs(csvPath)
	absXlsx, _ := filepath.Abs(xlsxPath)

	h := int(elapsed.Hours())
	m := int(elapsed.Minutes()) % 60
	s := int(elapsed.Seconds()) % 60

	fmt.Printf("\n%s\n", strings.Repeat("=", 52))
	fmt.Println("Done.")
	fmt.Printf("%s\n", strings.Repeat("=", 52))
	fmt.Printf("  Time    : %02d:%02d:%02d\n", h, m, s)
	fmt.Printf("  Results : %d businesses\n", len(data))
	fmt.Printf("  CSV     : %s\n", absCsv)
	fmt.Printf("  Excel   : %s\n", absXlsx)
	fmt.Printf("%s\n\n", strings.Repeat("=", 52))
}
