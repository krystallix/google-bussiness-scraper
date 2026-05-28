package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	sc "google-bussines-scraper/internal/scraper"
)

// ── State ────────────────────────────────────────────────────────────────────

type jobState struct {
	mu        sync.RWMutex
	Running   bool        `json:"running"`
	Done      bool        `json:"done"`
	Error     string      `json:"error,omitempty"`
	Results   []leadJSON  `json:"results"`
	Total     int         `json:"total"`
	Keyword   string      `json:"keyword"`
	StartedAt time.Time   `json:"started_at"`
}

type leadJSON struct {
	Name       string  `json:"name"`
	Category   string  `json:"category"`
	Rating     float64 `json:"rating"`
	Reviews    int     `json:"reviews"`
	Address    string  `json:"address"`
	Phone      string  `json:"phone"`
	Website    string  `json:"website"`
	MapsURL    string  `json:"maps_url"`
	HasWebsite bool    `json:"has_website"`
	HasPhone   bool    `json:"has_phone"`
	Tier       string  `json:"tier"`
	Reason     string  `json:"reason"`
}

var state = &jobState{}

// ── Request types ─────────────────────────────────────────────────────────────

type scrapeRequest struct {
	Keyword string  `json:"keyword"`
	Limit   int     `json:"limit"`
	Headful bool    `json:"headful"`
	Delay   float64 `json:"delay"`
}

// ── Main ─────────────────────────────────────────────────────────────────────

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", serveIndex)
	mux.HandleFunc("POST /api/scrape", handleScrape)
	mux.HandleFunc("GET /api/status", handleStatus)
	mux.HandleFunc("GET /api/results", handleResults)

	addr := "127.0.0.1:8080"
	fmt.Printf("\nG-Lead Scraper dashboard running at http://%s\n\n", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}

// ── Handlers ──────────────────────────────────────────────────────────────────

func serveIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(indexHTML))
}

func handleScrape(w http.ResponseWriter, r *http.Request) {
	var req scrapeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if req.Keyword == "" {
		jsonError(w, "keyword is required", http.StatusBadRequest)
		return
	}
	if req.Limit <= 0 {
		req.Limit = 50
	}
	if req.Delay <= 0 {
		req.Delay = 1.5
	}

	state.mu.Lock()
	if state.Running {
		state.mu.Unlock()
		jsonError(w, "scraper already running", http.StatusConflict)
		return
	}
	*state = jobState{Running: true, Keyword: req.Keyword, StartedAt: time.Now()}
	state.mu.Unlock()

	go func() {
		businesses, err := sc.ScrapeGoogleMaps(req.Keyword, req.Limit, req.Headful, req.Delay)

		state.mu.Lock()
		defer state.mu.Unlock()
		state.Running = false
		state.Done = true

		if err != nil {
			state.Error = err.Error()
			return
		}

		leads := sc.BusinessesToLeads(businesses)
		state.Total = len(leads)
		state.Results = make([]leadJSON, 0, len(leads))
		for _, l := range leads {
			state.Results = append(state.Results, leadJSON{
				Name:       l.Name,
				Category:   l.Category,
				Rating:     l.RatingFloat,
				Reviews:    l.Reviews,
				Address:    l.Address,
				Phone:      l.Phone,
				Website:    l.Website,
				MapsURL:    l.MapsURL,
				HasWebsite: l.HasWebsite,
				HasPhone:   l.HasPhone,
				Tier:       l.Tier,
				Reason:     l.Reason,
			})
		}

		// Save files in background
		ts := time.Now().Format("20060102_150405")
		csvPath := fmt.Sprintf("output/google_maps_%s_%s.csv", sc.CleanFilename(req.Keyword), ts)
		xlsxPath := fmt.Sprintf("output/google_maps_%s_%s.xlsx", sc.CleanFilename(req.Keyword), ts)
		sc.SaveCSV(csvPath, businesses)   //nolint:errcheck
		sc.SaveXLSX(xlsxPath, businesses) //nolint:errcheck
		fmt.Printf("Saved: %s | %s\n", csvPath, xlsxPath)
	}()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "started"})
}

func handleStatus(w http.ResponseWriter, r *http.Request) {
	state.mu.RLock()
	defer state.mu.RUnlock()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"running": state.Running,
		"done":    state.Done,
		"error":   state.Error,
		"total":   state.Total,
		"keyword": state.Keyword,
	})
}

func handleResults(w http.ResponseWriter, r *http.Request) {
	state.mu.RLock()
	defer state.mu.RUnlock()
	w.Header().Set("Content-Type", "application/json")
	results := state.Results
	if results == nil {
		results = []leadJSON{}
	}
	json.NewEncoder(w).Encode(results)
}

func jsonError(w http.ResponseWriter, msg string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}
