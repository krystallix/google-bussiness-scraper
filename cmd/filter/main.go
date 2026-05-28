package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	sc "google-bussines-scraper/internal/scraper"

	"github.com/xuri/excelize/v2"
)

func main() {
	files, err := filepath.Glob("output/*.csv")
	if err != nil || len(files) == 0 {
		fmt.Println("No CSV files found in output/")
		os.Exit(1)
	}

	fmt.Printf("Found %d CSV file(s):\n", len(files))
	for _, f := range files {
		// skip filtered output files
		if strings.Contains(filepath.Base(f), "filtered_leads") {
			continue
		}
		fmt.Printf("  - %s\n", f)
	}
	fmt.Println()

	var allBusinesses []sc.Business
	for _, f := range files {
		if strings.Contains(filepath.Base(f), "filtered_leads") {
			continue
		}
		bs, err := sc.ParseCSV(f)
		if err != nil {
			fmt.Printf("WARN: could not read %s: %v\n", f, err)
			continue
		}
		allBusinesses = append(allBusinesses, bs...)
	}

	// Deduplicate by phone
	seen := map[string]bool{}
	var unique []sc.Business
	for _, b := range allBusinesses {
		key := b.Phone
		if key == "" {
			key = b.Name
		}
		if !seen[key] {
			seen[key] = true
			unique = append(unique, b)
		}
	}

	leads := sc.BusinessesToLeads(unique)

	// Sort HIGH first, then STANDARD, then SKIP; within tier by reviews desc
	tierOrder := map[string]int{sc.TierHigh: 0, sc.TierStandard: 1, sc.TierSkip: 2}
	sort.SliceStable(leads, func(i, j int) bool {
		ti := tierOrder[leads[i].Tier]
		tj := tierOrder[leads[j].Tier]
		if ti != tj {
			return ti < tj
		}
		return leads[i].Reviews > leads[j].Reviews
	})

	printSummary(leads)

	ts := time.Now().Format("20060102_150405")
	outCSV := fmt.Sprintf("output/filtered_leads_%s.csv", ts)
	outXLSX := fmt.Sprintf("output/filtered_leads_%s.xlsx", ts)

	if err := saveFilteredCSV(outCSV, leads); err != nil {
		fmt.Printf("ERROR saving CSV: %v\n", err)
	}
	if err := saveFilteredXLSX(outXLSX, leads); err != nil {
		fmt.Printf("ERROR saving Excel: %v\n", err)
	}

	absCSV, _ := filepath.Abs(outCSV)
	absXLSX, _ := filepath.Abs(outXLSX)

	fmt.Printf("\n%s\n", strings.Repeat("=", 60))
	fmt.Println("Output saved:")
	fmt.Printf("  CSV   : %s\n", absCSV)
	fmt.Printf("  Excel : %s\n", absXLSX)
	fmt.Printf("%s\n\n", strings.Repeat("=", 60))
}

func printSummary(leads []sc.Lead) {
	counts := map[string]int{}
	for _, l := range leads {
		counts[l.Tier]++
	}

	fmt.Printf("%s\n", strings.Repeat("=", 60))
	fmt.Printf("Lead Analysis -- %d total businesses\n", len(leads))
	fmt.Printf("%s\n", strings.Repeat("=", 60))
	fmt.Printf("  HIGH POTENTIAL : %d\n", counts[sc.TierHigh])
	fmt.Printf("  STANDARD       : %d\n", counts[sc.TierStandard])
	fmt.Printf("  SKIP           : %d\n", counts[sc.TierSkip])
	fmt.Printf("%s\n\n", strings.Repeat("-", 60))

	for _, tier := range []string{sc.TierHigh, sc.TierStandard, sc.TierSkip} {
		fmt.Printf("[ %s ]\n", tier)
		for _, l := range leads {
			if l.Tier != tier {
				continue
			}
			website := "-"
			if l.HasWebsite {
				website = truncate(l.Website, 35)
			}
			fmt.Printf("  %-40s  Rating: %.1f  Reviews: %4d  Web: %s\n",
				truncate(l.Name, 40), l.RatingFloat, l.Reviews, website)
			fmt.Printf("  %s  Phone: %s\n", strings.Repeat(" ", 40), l.Phone)
			fmt.Printf("  Reason: %s\n\n", l.Reason)
		}
	}
}

func truncate(s string, max int) string {
	if utf8.RuneCountInString(s) <= max {
		return s
	}
	runes := []rune(s)
	return string(runes[:max-3]) + "..."
}

func saveFilteredCSV(path string, leads []sc.Lead) error {
	// Build businesses with tier info as extra columns
	var rows [][]string
	rows = append(rows, []string{"Tier", "Nama Bisnis", "Rating", "Jumlah Ulasan", "Has Website", "Website", "Telepon", "Alamat", "Alasan", "Google Maps URL"})
	for _, l := range leads {
		hasWeb := "No"
		if l.HasWebsite {
			hasWeb = "Yes"
		}
		rows = append(rows, []string{
			l.Tier, l.Name,
			fmt.Sprintf("%.1f", l.RatingFloat),
			strconv.Itoa(l.Reviews),
			hasWeb, l.Website, l.Phone, l.Address, l.Reason, l.MapsURL,
		})
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	f.WriteString("\xef\xbb\xbf")

	for _, row := range rows {
		f.WriteString(strings.Join(quoteCSV(row), ",") + "\n")
	}
	return nil
}

func quoteCSV(row []string) []string {
	out := make([]string, len(row))
	for i, v := range row {
		if strings.ContainsAny(v, ",\"\n\r") {
			v = "\"" + strings.ReplaceAll(v, "\"", "\"\"") + "\""
		}
		out[i] = v
	}
	return out
}

var tierColors = map[string]string{
	sc.TierHigh:     "63BE7B",
	sc.TierStandard: "FFEB84",
	sc.TierSkip:     "F8696B",
}

func saveFilteredXLSX(path string, leads []sc.Lead) error {
	xf := excelize.NewFile()
	defer xf.Close()

	sheet := "Leads"
	xf.NewSheet(sheet)
	xf.DeleteSheet("Sheet1")

	headers := []string{"Tier", "Nama Bisnis", "Rating", "Reviews", "Has Website", "Website", "Telepon", "Alamat", "Alasan", "Maps URL"}
	colWidths := []float64{18, 38, 8, 10, 12, 35, 18, 45, 50, 20}

	boldStyle, _ := xf.NewStyle(&excelize.Style{Font: &excelize.Font{Bold: true}})
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		xf.SetCellValue(sheet, cell, h)
		xf.SetCellStyle(sheet, cell, cell, boldStyle)
	}

	cols := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J"}
	for i, col := range cols {
		xf.SetColWidth(sheet, col, col, colWidths[i])
	}

	for row, l := range leads {
		excelRow := row + 2
		hasWeb := "No"
		if l.HasWebsite {
			hasWeb = "Yes"
		}
		values := []interface{}{
			l.Tier, l.Name, l.RatingFloat, l.Reviews,
			hasWeb, l.Website, l.Phone, l.Address, l.Reason, l.MapsURL,
		}
		hex := tierColors[l.Tier]
		rowStyle, _ := xf.NewStyle(&excelize.Style{
			Fill: excelize.Fill{Type: "pattern", Color: []string{hex}, Pattern: 1},
		})
		for col, v := range values {
			cell, _ := excelize.CoordinatesToCellName(col+1, excelRow)
			xf.SetCellValue(sheet, cell, v)
			xf.SetCellStyle(sheet, cell, cell, rowStyle)
		}
	}

	xf.SetPanes(sheet, &excelize.Panes{
		Freeze: true, YSplit: 1, TopLeftCell: "A2", ActivePane: "bottomLeft",
	})
	return xf.SaveAs(path)
}
