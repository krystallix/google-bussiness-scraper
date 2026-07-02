# Google Maps Business Scraper

Scrape business listings from Google Maps. No API key required.

Uses [go-rod](https://github.com/go-rod/rod) (Chromium automation) to handle dynamic JS-rendered content.

Comes with a **Go CLI** (scraper + lead classifier + exporter) and a **Next.js web dashboard** (scrape, filter, manage leads, send WhatsApp messages).

---

## What it collects

| Field | Description |
|---|---|
| Business Name | Listing title |
| Category | Business category (restaurant, plumber, etc.) |
| Rating | Star rating (1.0–5.0) |
| Reviews | Number of reviews |
| Address | Full address |
| Phone | Contact number |
| Website | Business website URL |
| Google Maps URL | Direct link to the listing |

---

## Requirements

- Go 1.22+
- Node.js 18+ (for web dashboard)
- Internet connection (downloads Chromium automatically on first run)

---

## Setup

```bash
git clone <repo-url>
cd google-bussiness-scraper
go mod tidy
```

Chromium downloads automatically on first run.

---

## Usage

### CLI

```bash
make run
```

Or manually:

```bash
go build -o scraper ./cmd/server
./scraper
```

Prompts:

```
Search keyword (e.g. coffee shop Bandung): laundry Jakarta
Max results to collect (default 50): 30
Show browser window? (y/N): n
Delay between pages in seconds (default 1.5): 1.5
```

Output saved to `output/` as CSV + Excel:

```
output/google_maps_laundry_jakarta_20240701_120000.csv
output/google_maps_laundry_jakarta_20240701_120000.xlsx
```

### Lead Filter

```bash
go run ./cmd/filter
```

Loads data from `scraped_businesses.json` and CSV files in `output/`, classifies leads by tier (HIGH POTENTIAL / STANDARD / SKIP / COMPLETED), and exports filtered results.

### Web Dashboard

```bash
cd web
npm install
npm run dev
```

Open `http://localhost:3000` in browser.

Features:
- Scrape Google Maps via browser UI
- Browse, search, sort, filter leads
- Edit / delete / bulk delete leads
- Send WhatsApp messages (individual or bulk) with customizable templates
- WhatsApp connection status
- Real-time scraping progress

---

## Build options

| Command | Description |
|---|---|
| `make run` | Build and run CLI |
| `make build` | Compile binary |
| `make clean` | Remove binary |

---

## Notes

- Limit above 500 increases RAM usage and IP-block risk. Use reasonable limits and add delay.
- Intended for research and personal use. Respect Google's Terms of Service.

---

## Dependencies

| Package | Purpose |
|---|---|
| [go-rod/rod](https://github.com/go-rod/rod) | Browser automation |
| [xuri/excelize](https://github.com/xuri/excelize) | Excel file I/O |
| [schollz/progressbar](https://github.com/schollz/progressbar) | Terminal progress bar |
