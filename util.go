package main

import (
	"encoding/csv"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
)

// Business holds all extracted data for a single Google Maps listing.
type Business struct {
	Name         string
	Category     string
	Rating       string
	ReviewsCount string
	Address      string
	Phone        string
	Website      string
	MapsURL      string
}

// saveCSV writes a slice of Business records to a UTF-8 CSV file.
func saveCSV(path string, businesses []Business) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	// Write UTF-8 BOM so Excel opens it correctly.
	f.WriteString("\xef\xbb\xbf")

	w := csv.NewWriter(f)
	defer w.Flush()

	headers := []string{
		"Nama Bisnis", "Kategori", "Rating", "Jumlah Ulasan",
		"Alamat", "Telepon", "Website", "Google Maps URL",
	}
	if err := w.Write(headers); err != nil {
		return err
	}
	for _, b := range businesses {
		row := []string{
			b.Name, b.Category, b.Rating, b.ReviewsCount,
			b.Address, b.Phone, b.Website, b.MapsURL,
		}
		if err := w.Write(row); err != nil {
			return err
		}
	}
	return nil
}

// saveXLSX writes a slice of Business records to an Excel .xlsx file.
func saveXLSX(path string, businesses []Business) error {
	xf := excelize.NewFile()
	defer xf.Close()

	sheet := "Sheet1"
	headers := []string{
		"Nama Bisnis", "Kategori", "Rating", "Jumlah Ulasan",
		"Alamat", "Telepon", "Website", "Google Maps URL",
	}

	// Style: bold header
	boldStyle, _ := xf.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true},
	})

	for col, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(col+1, 1)
		xf.SetCellValue(sheet, cell, h)
		xf.SetCellStyle(sheet, cell, cell, boldStyle)
	}

	for row, b := range businesses {
		values := []string{
			b.Name, b.Category, b.Rating, b.ReviewsCount,
			b.Address, b.Phone, b.Website, b.MapsURL,
		}
		for col, v := range values {
			cell, _ := excelize.CoordinatesToCellName(col+1, row+2)
			xf.SetCellValue(sheet, cell, v)
		}
	}

	// Auto-fit columns (approximate widths)
	colWidths := []float64{30, 20, 8, 14, 40, 18, 35, 60}
	cols := []string{"A", "B", "C", "D", "E", "F", "G", "H"}
	for i, col := range cols {
		xf.SetColWidth(sheet, col, col, colWidths[i])
	}

	return xf.SaveAs(path)
}

// cleanFilename converts a string to a safe lowercase filename component.
var nonAlphanumRE = regexp.MustCompile(`[^a-zA-Z0-9_-]`)

func cleanFilename(name string) string {
	return strings.ToLower(nonAlphanumRE.ReplaceAllString(name, "_"))
}

// encodeKeyword URL-encodes a search keyword for the Google Maps search URL.
func encodeKeyword(keyword string) string {
	return url.QueryEscape(keyword)
}

// outputPaths returns the csv and xlsx output paths given a base name or generates
// one from the keyword and current timestamp.
func outputPaths(outputBase, keyword string) (csvPath, xlsxPath string) {
	if err := os.MkdirAll("output", 0755); err != nil {
		fmt.Fprintf(os.Stderr, "⚠️  Could not create output dir: %v\n", err)
	}
	if outputBase != "" {
		csvPath = filepath.Join("output", outputBase+".csv")
		xlsxPath = filepath.Join("output", outputBase+".xlsx")
	} else {
		ts := time.Now().Format("20060102_150405")
		base := "google_maps_" + cleanFilename(keyword) + "_" + ts
		csvPath = filepath.Join("output", base+".csv")
		xlsxPath = filepath.Join("output", base+".xlsx")
	}
	return
}
