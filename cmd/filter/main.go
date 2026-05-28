package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/xuri/excelize/v2"
)

// Tier categories
const (
	TierHigh     = "HIGH POTENTIAL"
	TierStandard = "STANDARD"
	TierSkip     = "SKIP"
)

type Lead struct {
	Name         string
	Category     string
	Rating       float64
	Reviews      int
	Address      string
	Phone        string
	Website      string
	MapsURL      string
	Tier         string
	Reason       string
	HasWebsite   bool
	HasPhone     bool
}

func main() {
	// Find all CSV files in output/
	files, err := filepath.Glob("output/*.csv")
	if err != nil || len(files) == 0 {
		fmt.Println("No CSV files found in output/")
		os.Exit(1)
	}

	fmt.Printf("Found %d CSV file(s):\n", len(files))
	for _, f := range files {
		fmt.Printf("  - %s\n", f)
	}
	fmt.Println()

	var allLeads []Lead

	for _, f := range files {
		leads, err := parseCSV(f)
		if err != nil {
			fmt.Printf("WARN: could not read %s: %v\n", f, err)
			continue
		}
		allLeads = append(allLeads, leads...)
	}

	// Remove duplicates by phone number
	allLeads = deduplicateByPhone(allLeads)

	// Classify each lead
	for i := range allLeads {
		allLeads[i].Tier, allLeads[i].Reason = classify(allLeads[i])
	}

	// Sort: HIGH first, then STANDARD, then SKIP; within tier sort by reviews desc
	sort.SliceStable(allLeads, func(i, j int) bool {
		ti := tierOrder(allLeads[i].Tier)
		tj := tierOrder(allLeads[j].Tier)
		if ti != tj {
			return ti < tj
		}
		return allLeads[i].Reviews > allLeads[j].Reviews
	})

	// Print summary to terminal
	printSummary(allLeads)

	// Save filtered output
	ts := time.Now().Format("20060102_150405")
	outCSV := fmt.Sprintf("output/filtered_leads_%s.csv", ts)
	outXLSX := fmt.Sprintf("output/filtered_leads_%s.xlsx", ts)

	if err := saveFilteredCSV(outCSV, allLeads); err != nil {
		fmt.Printf("ERROR saving CSV: %v\n", err)
	}
	if err := saveFilteredXLSX(outXLSX, allLeads); err != nil {
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

// classify assigns a tier and reason to a lead based on business signals.
//
// Scoring logic (tailored for selling website + business management systems):
//
//   HIGH POTENTIAL:
//     - Has phone, no website -> strong candidate (needs digital presence)
//     - High reviews (>=100) + no website -> proven business, missing online
//     - Rating >= 4.5 + reviews >= 50 -> active, quality business
//
//   SKIP:
//     - No phone and no rating -> incomplete listing, likely inactive
//     - Rating < 3.5 -> struggling business, low priority
//     - Very few reviews (<5) AND no website -> too small / just opened
//
//   STANDARD:
//     - Everything else
func classify(l Lead) (string, string) {
	// Skip conditions
	if !l.HasPhone && l.Rating == 0 {
		return TierSkip, "No phone, no rating — incomplete listing"
	}
	if l.Rating > 0 && l.Rating < 3.5 {
		return TierSkip, fmt.Sprintf("Low rating (%.1f) — struggling business", l.Rating)
	}
	if l.Reviews < 5 && !l.HasWebsite {
		return TierSkip, "Very few reviews and no website — too small or just opened"
	}

	// High potential conditions
	if l.HasPhone && !l.HasWebsite && l.Reviews >= 100 {
		return TierHigh, fmt.Sprintf("Active business (%d reviews), no website — prime target", l.Reviews)
	}
	if l.HasPhone && !l.HasWebsite && l.Rating >= 4.5 && l.Reviews >= 30 {
		return TierHigh, fmt.Sprintf("Rated %.1f with %d reviews, no website — easy upsell", l.Rating, l.Reviews)
	}
	if l.HasPhone && !l.HasWebsite && l.Reviews >= 50 {
		return TierHigh, fmt.Sprintf("Good volume (%d reviews), no digital presence", l.Reviews)
	}

	// Standard
	if l.HasWebsite {
		return TierStandard, "Has website — pitch business management system only"
	}
	if l.HasPhone && l.Reviews > 0 {
		return TierStandard, fmt.Sprintf("Active (%d reviews) but small — potential with nurturing", l.Reviews)
	}

	return TierStandard, "Moderate signals"
}

func parseCSV(path string) ([]Lead, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// Strip UTF-8 BOM if present
	buf := make([]byte, 3)
	n, _ := f.Read(buf)
	if n == 3 && buf[0] == 0xef && buf[1] == 0xbb && buf[2] == 0xbf {
		// BOM present, already consumed — continue reading from here
	} else {
		f.Seek(0, io.SeekStart)
	}

	r := csv.NewReader(f)
	r.LazyQuotes = true
	r.FieldsPerRecord = -1

	records, err := r.ReadAll()
	if err != nil {
		return nil, err
	}
	if len(records) < 2 {
		return nil, nil
	}

	// Map header names to column indices
	header := records[0]
	col := map[string]int{}
	for i, h := range header {
		col[strings.TrimSpace(h)] = i
	}

	get := func(row []string, key string) string {
		idx, ok := col[key]
		if !ok || idx >= len(row) {
			return ""
		}
		return strings.TrimSpace(row[idx])
	}

	var leads []Lead
	for _, row := range records[1:] {
		rating, _ := strconv.ParseFloat(strings.Replace(get(row, "Rating"), ",", ".", 1), 64)
		reviews, _ := strconv.Atoi(strings.Replace(get(row, "Jumlah Ulasan"), ".", "", -1))
		website := get(row, "Website")
		phone := get(row, "Telepon")

		leads = append(leads, Lead{
			Name:       get(row, "Nama Bisnis"),
			Category:   get(row, "Kategori"),
			Rating:     math.Round(rating*10) / 10,
			Reviews:    reviews,
			Address:    get(row, "Alamat"),
			Phone:      phone,
			Website:    website,
			MapsURL:    get(row, "Google Maps URL"),
			HasWebsite: website != "",
			HasPhone:   phone != "",
		})
	}
	return leads, nil
}

func deduplicateByPhone(leads []Lead) []Lead {
	seen := map[string]bool{}
	var out []Lead
	for _, l := range leads {
		key := l.Phone
		if key == "" {
			key = l.Name // fallback dedup by name
		}
		if !seen[key] {
			seen[key] = true
			out = append(out, l)
		}
	}
	return out
}

func tierOrder(t string) int {
	switch t {
	case TierHigh:
		return 0
	case TierStandard:
		return 1
	default:
		return 2
	}
}

func printSummary(leads []Lead) {
	counts := map[string]int{}
	for _, l := range leads {
		counts[l.Tier]++
	}

	fmt.Printf("%s\n", strings.Repeat("=", 60))
	fmt.Printf("Lead Analysis — %d total businesses\n", len(leads))
	fmt.Printf("%s\n", strings.Repeat("=", 60))
	fmt.Printf("  HIGH POTENTIAL : %d\n", counts[TierHigh])
	fmt.Printf("  STANDARD       : %d\n", counts[TierStandard])
	fmt.Printf("  SKIP           : %d\n", counts[TierSkip])
	fmt.Printf("%s\n\n", strings.Repeat("-", 60))

	for _, tier := range []string{TierHigh, TierStandard, TierSkip} {
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
				truncate(l.Name, 40), l.Rating, l.Reviews, website)
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

func saveFilteredCSV(path string, leads []Lead) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	f.WriteString("\xef\xbb\xbf")

	w := csv.NewWriter(f)
	defer w.Flush()

	headers := []string{"Tier", "Nama Bisnis", "Rating", "Jumlah Ulasan", "Has Website", "Website", "Telepon", "Alamat", "Alasan", "Google Maps URL"}
	w.Write(headers) //nolint:errcheck

	for _, l := range leads {
		hasWeb := "No"
		if l.HasWebsite {
			hasWeb = "Yes"
		}
		w.Write([]string{ //nolint:errcheck
			l.Tier,
			l.Name,
			fmt.Sprintf("%.1f", l.Rating),
			strconv.Itoa(l.Reviews),
			hasWeb,
			l.Website,
			l.Phone,
			l.Address,
			l.Reason,
			l.MapsURL,
		})
	}
	return nil
}

var tierColors = map[string]string{
	TierHigh:     "63BE7B", // green
	TierStandard: "FFEB84", // yellow
	TierSkip:     "F8696B", // red
}

func saveFilteredXLSX(path string, leads []Lead) error {
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
			l.Tier,
			l.Name,
			l.Rating,
			l.Reviews,
			hasWeb,
			l.Website,
			l.Phone,
			l.Address,
			l.Reason,
			l.MapsURL,
		}

		// Row background color based on tier
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

	// Freeze header row
	xf.SetPanes(sheet, &excelize.Panes{
		Freeze:      true,
		YSplit:      1,
		TopLeftCell: "A2",
		ActivePane:  "bottomLeft",
	})

	return xf.SaveAs(path)
}
