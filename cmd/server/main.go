package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	sc "google-bussines-scraper/internal/scraper"
	wa "google-bussines-scraper/internal/whatsapp"
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
	ID         string  `json:"id"`
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

type addLeadRequest struct {
	Name     string  `json:"name"`
	Category string  `json:"category"`
	Rating   float64 `json:"rating"`
	Reviews  int     `json:"reviews"`
	Address  string  `json:"address"`
	Phone    string  `json:"phone"`
	Website  string  `json:"website"`
	MapsURL  string  `json:"maps_url"`
	Status   string  `json:"status"`
}

type updateLeadRequest struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Category string  `json:"category"`
	Rating   float64 `json:"rating"`
	Reviews  int     `json:"reviews"`
	Address  string  `json:"address"`
	Phone    string  `json:"phone"`
	Website  string  `json:"website"`
	MapsURL  string  `json:"maps_url"`
	Status   string  `json:"status"`
}

type deleteLeadRequest struct {
	ID string `json:"id"`
}

// ── Main ─────────────────────────────────────────────────────────────────────

func main() {
	// Load initial state from persistent JSON database
	if _, err := os.Stat("scraped_businesses.json"); err == nil {
		businesses, err := sc.LoadPersistentJSON()
		if err == nil && len(businesses) > 0 {
			state.mu.Lock()
			reloadMemoryState(businesses)
			state.mu.Unlock()
			fmt.Printf("Loaded %d leads from persistent database (scraped_businesses.json)\n", len(businesses))
		}
	}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/scrape", handleScrape)
	mux.HandleFunc("GET /api/status", handleStatus)
	mux.HandleFunc("GET /api/results", handleResults)
	mux.HandleFunc("POST /api/leads/add", handleAddLead)
	mux.HandleFunc("POST /api/leads/update", handleUpdateLead)
	mux.HandleFunc("POST /api/leads/delete", handleDeleteLead)
	mux.HandleFunc("POST /api/leads/complete", handleCompleteLead)

	// WhatsApp endpoints
	mux.HandleFunc("GET /api/wa/status", handleWaStatus)
	mux.HandleFunc("GET /api/wa/qr", handleWaQR)
	mux.HandleFunc("POST /api/wa/logout", handleWaLogout)
	mux.HandleFunc("POST /api/wa/send", handleWaSend)

	// Initialize whatsapp
	go func() {
		if err := wa.Init(); err != nil {
			log.Printf("Failed to init whatsapp: %v\n", err)
		}
	}()

	addr := "127.0.0.1:8080"
	fmt.Printf("\nG-Lead Scraper dashboard running at http://%s\n\n", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}

// ── Handlers ──────────────────────────────────────────────────────────────────

type waSendRequest struct {
	Phone   string `json:"phone"`
	Message string `json:"message"`
}

func handleWaStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"connected": wa.IsConnected()})
}

func handleWaQR(w http.ResponseWriter, r *http.Request) {
	if wa.IsConnected() {
		jsonError(w, "Already connected", http.StatusBadRequest)
		return
	}
	qr := wa.GetQR()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"qr": qr})
}

func handleWaLogout(w http.ResponseWriter, r *http.Request) {
	err := wa.Logout()
	if err != nil {
		jsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

func handleWaSend(w http.ResponseWriter, r *http.Request) {
	var req waSendRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if req.Phone == "" || req.Message == "" {
		jsonError(w, "phone and message are required", http.StatusBadRequest)
		return
	}
	err := wa.SendMessage(req.Phone, req.Message)
	if err != nil {
		jsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
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
	// Reset fields individually — never replace the whole struct while the
	// mutex is held, as that would zero out the mutex itself and panic on Unlock.
	state.Running = true
	state.Done = false
	state.Error = ""
	state.Results = nil
	state.Total = 0
	state.Keyword = req.Keyword
	state.StartedAt = time.Now()
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

		// Save files in background
		ts := time.Now().Format("20060102_150405")
		csvPath := fmt.Sprintf("output/google_maps_%s_%s.csv", sc.CleanFilename(req.Keyword), ts)
		xlsxPath := fmt.Sprintf("output/google_maps_%s_%s.xlsx", sc.CleanFilename(req.Keyword), ts)
		sc.SaveCSV(csvPath, businesses)   //nolint:errcheck
		sc.SaveXLSX(xlsxPath, businesses) //nolint:errcheck
		fmt.Printf("Saved: %s | %s\n", csvPath, xlsxPath)

		// Save to persistent JSON
		_ = sc.SavePersistentJSON(businesses)

		// Reload all accumulated businesses to state so dashboard reflects full history
		accumulated, err := sc.LoadPersistentJSON()
		if err == nil {
			businesses = accumulated
		}

		reloadMemoryState(businesses)
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

// reloadMemoryState updates the in-memory state.Results from a list of sc.Business records.
// Expects state.mu to be locked already by the caller.
func reloadMemoryState(businesses []sc.Business) {
	leads := sc.BusinessesToLeads(businesses)
	state.Total = len(leads)
	state.Results = make([]leadJSON, 0, len(leads))
	for _, l := range leads {
		state.Results = append(state.Results, leadJSON{
			ID:         l.ID,
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
	state.Done = true
}

func handleAddLead(w http.ResponseWriter, r *http.Request) {
	var req addLeadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if req.Name == "" {
		jsonError(w, "business name is required", http.StatusBadRequest)
		return
	}

	state.mu.Lock()
	defer state.mu.Unlock()

	// Convert requests to Business
	newBiz := sc.Business{
		Name:         req.Name,
		Category:     req.Category,
		Rating:       fmt.Sprintf("%.1f", req.Rating),
		ReviewsCount: fmt.Sprintf("%d", req.Reviews),
		Address:      req.Address,
		Phone:        req.Phone,
		Website:      req.Website,
		MapsURL:      req.MapsURL,
		Status:       req.Status,
	}
	newBiz.ID = sc.GenerateID(newBiz)

	// Load existing persistent JSON
	businesses, err := sc.LoadPersistentJSON()
	if err != nil {
		// If loading failed (e.g. no file), start empty
		businesses = []sc.Business{}
	}

	// Append and save back
	businesses = append(businesses, newBiz)
	if err := sc.SaveAllPersistentJSON(businesses); err != nil {
		jsonError(w, fmt.Sprintf("failed to save: %v", err), http.StatusInternalServerError)
		return
	}

	// Reload all to memory state
	reloadMemoryState(businesses)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"status": "success", "id": newBiz.ID})
}

func handleUpdateLead(w http.ResponseWriter, r *http.Request) {
	var req updateLeadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if req.ID == "" {
		jsonError(w, "id is required", http.StatusBadRequest)
		return
	}

	state.mu.Lock()
	defer state.mu.Unlock()

	businesses, err := sc.LoadPersistentJSON()
	if err != nil {
		jsonError(w, "database empty or unreadable", http.StatusNotFound)
		return
	}

	found := false
	for i, b := range businesses {
		if b.ID == req.ID {
			businesses[i].Name = req.Name
			businesses[i].Category = req.Category
			businesses[i].Rating = fmt.Sprintf("%.1f", req.Rating)
			businesses[i].ReviewsCount = fmt.Sprintf("%d", req.Reviews)
			businesses[i].Address = req.Address
			businesses[i].Phone = req.Phone
			businesses[i].Website = req.Website
			businesses[i].MapsURL = req.MapsURL
			businesses[i].Status = req.Status
			found = true
			break
		}
	}

	if !found {
		jsonError(w, "business not found", http.StatusNotFound)
		return
	}

	if err := sc.SaveAllPersistentJSON(businesses); err != nil {
		jsonError(w, fmt.Sprintf("failed to save: %v", err), http.StatusInternalServerError)
		return
	}

	reloadMemoryState(businesses)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func handleDeleteLead(w http.ResponseWriter, r *http.Request) {
	var req deleteLeadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if req.ID == "" {
		jsonError(w, "id is required", http.StatusBadRequest)
		return
	}

	state.mu.Lock()
	defer state.mu.Unlock()

	businesses, err := sc.LoadPersistentJSON()
	if err != nil {
		jsonError(w, "database empty or unreadable", http.StatusNotFound)
		return
	}

	found := false
	var updated []sc.Business
	for _, b := range businesses {
		if b.ID == req.ID {
			found = true
			continue
		}
		updated = append(updated, b)
	}

	if !found {
		jsonError(w, "business not found", http.StatusNotFound)
		return
	}

	if err := sc.SaveAllPersistentJSON(updated); err != nil {
		jsonError(w, fmt.Sprintf("failed to save: %v", err), http.StatusInternalServerError)
		return
	}

	reloadMemoryState(updated)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

type completeLeadRequest struct {
	ID string `json:"id"`
}

func handleCompleteLead(w http.ResponseWriter, r *http.Request) {
	var req completeLeadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if req.ID == "" {
		jsonError(w, "id is required", http.StatusBadRequest)
		return
	}

	state.mu.Lock()
	defer state.mu.Unlock()

	businesses, err := sc.LoadPersistentJSON()
	if err != nil {
		jsonError(w, "database empty or unreadable", http.StatusNotFound)
		return
	}

	found := false
	for i, b := range businesses {
		if b.ID == req.ID {
			businesses[i].Status = "COMPLETED"
			found = true
			break
		}
	}

	if !found {
		jsonError(w, "lead not found", http.StatusNotFound)
		return
	}

	if err := sc.SaveAllPersistentJSON(businesses); err != nil {
		jsonError(w, fmt.Sprintf("failed to save: %v", err), http.StatusInternalServerError)
		return
	}

	reloadMemoryState(businesses)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}
