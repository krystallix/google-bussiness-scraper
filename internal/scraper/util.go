package scraper

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"math"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
)

// Tier categories
const (
	TierHigh      = "HIGH POTENTIAL"
	TierStandard  = "STANDARD"
	TierSkip      = "SKIP"
	TierCompleted = "COMPLETED"
)

// Business holds all extracted data for a single Google Maps listing.
type Business struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Category     string `json:"category"`
	Rating       string `json:"rating"`
	ReviewsCount string `json:"reviews_count"`
	Address      string `json:"address"`
	Phone        string `json:"phone"`
	Website      string `json:"website"`
	MapsURL      string `json:"maps_url"`
	Status       string `json:"status"`
}

// Lead holds the classified lead data.
type Lead struct {
	ID           string  `json:"id"`
	Name         string  `json:"name"`
	Category     string  `json:"category"`
	RatingFloat  float64 `json:"rating"`
	Reviews      int     `json:"reviews"`
	Address      string  `json:"address"`
	Phone        string  `json:"phone"`
	Website      string  `json:"website"`
	MapsURL      string  `json:"maps_url"`
	Tier         string  `json:"tier"`
	Reason       string  `json:"reason"`
	HasWebsite   bool    `json:"has_website"`
	HasPhone     bool    `json:"has_phone"`
}

// SaveCSV writes a slice of Business records to a UTF-8 CSV file.
func SaveCSV(path string, businesses []Business) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

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

// SaveXLSX writes a slice of Business records to an Excel .xlsx file.
func SaveXLSX(path string, businesses []Business) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

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

// ParseCSV reads a CSV file containing saved Business records and returns a slice of Business.
func ParseCSV(path string) ([]Business, error) {
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
		f.Seek(0, 0)
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

	var businesses []Business
	for _, row := range records[1:] {
		businesses = append(businesses, Business{
			Name:         get(row, "Nama Bisnis"),
			Category:     get(row, "Kategori"),
			Rating:       get(row, "Rating"),
			ReviewsCount: get(row, "Jumlah Ulasan"),
			Address:      get(row, "Alamat"),
			Phone:        get(row, "Telepon"),
			Website:      get(row, "Website"),
			MapsURL:      get(row, "Google Maps URL"),
		})
	}
	return businesses, nil
}

// CleanFilename converts a string to a safe lowercase filename component.
var nonAlphanumRE = regexp.MustCompile(`[^a-zA-Z0-9_-]`)

func CleanFilename(name string) string {
	return strings.ToLower(nonAlphanumRE.ReplaceAllString(name, "_"))
}

// encodeKeyword URL-encodes a search keyword for the Google Maps search URL.
func encodeKeyword(keyword string) string {
	return url.QueryEscape(keyword)
}

// OutputPaths returns the csv and xlsx output paths given a base name or generates
// one from the keyword and current timestamp.
func OutputPaths(outputBase, keyword string) (csvPath, xlsxPath string) {
	if err := os.MkdirAll("output", 0755); err != nil {
		fmt.Fprintf(os.Stderr, "⚠️  Could not create output dir: %v\n", err)
	}
	if outputBase != "" {
		csvPath = filepath.Join("output", outputBase+".csv")
		xlsxPath = filepath.Join("output", outputBase+".xlsx")
	} else {
		ts := time.Now().Format("20060102_150405")
		base := "google_maps_" + CleanFilename(keyword) + "_" + ts
		csvPath = filepath.Join("output", base+".csv")
		xlsxPath = filepath.Join("output", base+".xlsx")
	}
	return
}

// BusinessesToLeads converts a slice of Business records into a slice of classified Leads.
func BusinessesToLeads(businesses []Business) []Lead {
	var leads []Lead
	for _, b := range businesses {
		rating, _ := strconv.ParseFloat(strings.Replace(b.Rating, ",", ".", 1), 64)
		reviews, _ := strconv.Atoi(strings.NewReplacer(".", "", ",", "").Replace(b.ReviewsCount))
		ratingVal := math.Round(rating*10) / 10

		l := Lead{
			ID:          b.ID,
			Name:        b.Name,
			Category:    b.Category,
			RatingFloat: ratingVal,
			Reviews:     reviews,
			Address:     b.Address,
			Phone:       b.Phone,
			Website:     b.Website,
			MapsURL:     b.MapsURL,
			HasWebsite:  b.Website != "",
			HasPhone:    b.Phone != "",
		}

		if b.Status == "COMPLETED" {
			l.Tier = TierCompleted
			l.Reason = "Message sent successfully via WhatsApp"
		} else {
			l.Tier, l.Reason = classify(l)
		}
		leads = append(leads, l)
	}
	return leads
}

// classify assigns a tier and reason to a lead based on business signals.
func classify(l Lead) (string, string) {
	// Skip conditions
	if !l.HasPhone && l.RatingFloat == 0 {
		return TierSkip, "No phone, no rating — incomplete listing"
	}
	if l.RatingFloat > 0 && l.RatingFloat < 3.5 {
		return TierSkip, fmt.Sprintf("Low rating (%.1f) — struggling business", l.RatingFloat)
	}
	if l.Reviews < 5 && !l.HasWebsite {
		return TierSkip, "Very few reviews and no website — too small or just opened"
	}

	// High potential conditions
	if l.HasPhone && !l.HasWebsite && l.Reviews >= 100 {
		return TierHigh, fmt.Sprintf("Active business (%d reviews), no website — prime target", l.Reviews)
	}
	if l.HasPhone && !l.HasWebsite && l.RatingFloat >= 4.5 && l.Reviews >= 30 {
		return TierHigh, fmt.Sprintf("Rated %.1f with %d reviews, no website — easy upsell", l.RatingFloat, l.Reviews)
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

// GenerateID creates a deterministic unique MD5 hash for a business.
func GenerateID(b Business) string {
	input := b.MapsURL
	if input == "" {
		input = b.Phone
	}
	if input == "" {
		input = b.Name + "|" + b.Address
	}
	hash := md5.Sum([]byte(input))
	return hex.EncodeToString(hash[:])
}

// SavePersistentJSON appends new businesses to the persistent scraped_businesses.json file.
// It loads existing businesses, merges the new ones, deduplicates, and saves them back.
func SavePersistentJSON(newBusinesses []Business) error {
	path := "scraped_businesses.json"

	// 1. Read existing businesses if file exists
	var existing []Business
	if _, err := os.Stat(path); err == nil {
		data, err := os.ReadFile(path)
		if err == nil {
			_ = json.Unmarshal(data, &existing)
		}
	}

	// 2. Generate IDs if missing
	for i, b := range existing {
		if b.ID == "" {
			existing[i].ID = GenerateID(b)
		}
	}
	for i, b := range newBusinesses {
		if b.ID == "" {
			newBusinesses[i].ID = GenerateID(b)
		}
	}

	// 3. Merge existing and new
	all := append(existing, newBusinesses...)

	// 4. Deduplicate by ID
	seen := map[string]bool{}
	var unique []Business
	for _, b := range all {
		if b.ID != "" && !seen[b.ID] {
			seen[b.ID] = true
			unique = append(unique, b)
		}
	}

	return SaveAllPersistentJSON(unique)
}

// SaveAllPersistentJSON overwrites the scraped_businesses.json file with a complete slice of Business.
func SaveAllPersistentJSON(businesses []Business) error {
	path := "scraped_businesses.json"

	// Ensure IDs are generated
	for i, b := range businesses {
		if b.ID == "" {
			businesses[i].ID = GenerateID(b)
		}
	}

	jsonData, err := json.MarshalIndent(businesses, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	err = os.WriteFile(path, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("failed to write persistent JSON: %w", err)
	}

	return nil
}

// LoadPersistentJSON loads all businesses from the persistent scraped_businesses.json file.
func LoadPersistentJSON() ([]Business, error) {
	path := "scraped_businesses.json"
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var businesses []Business
	if err := json.Unmarshal(data, &businesses); err != nil {
		return nil, err
	}

	// Ensure all have valid IDs
	changed := false
	for i, b := range businesses {
		if b.ID == "" {
			businesses[i].ID = GenerateID(b)
			changed = true
		}
	}
	if changed {
		_ = SaveAllPersistentJSON(businesses)
	}

	return businesses, nil
}

