# google-bussiness-scraper

A command-line tool written in Go that scrapes business listings from Google Maps and saves the results to CSV and Excel files.

It uses [go-rod](https://github.com/go-rod/rod) to control a Chromium browser, which lets it handle the dynamic content that Google Maps renders through JavaScript. No API key required.

---

## What it collects

For each business found, the tool extracts:

- Business name
- Category
- Rating
- Review count
- Address
- Phone number
- Website
- Google Maps URL

---

## Requirements

- Go 1.22 or later
- Internet connection (to download Chromium on first run and to access Google Maps)

---

## Setup

Clone the repository and download dependencies:

```bash
git clone https://github.com/krystallix/google-bussiness-scraper.git
cd google-bussiness-scraper
go mod tidy
```

go-rod will automatically download a compatible Chromium browser to your local cache the first time you run the tool. There is no separate browser installation step.

---

## Usage

Build and run:

```bash
make run
```

Or manually:

```bash
go build -o scraper .
./scraper
```

The tool will prompt you for input:

```
Search keyword (e.g. coffee shop Bandung): hotel Jakarta

Max results to collect (default 50): 30

Show browser window? (y/N): n

Delay between pages in seconds (default 1.5): 1.5
```

After that it runs automatically and saves two files to the `output/` directory.

---

## Output

Results are saved as both CSV and Excel, named after the keyword and timestamp:

```
output/google_maps_hotel_jakarta_20240528_103000.csv
output/google_maps_hotel_jakarta_20240528_103000.xlsx
```

The CSV file uses UTF-8 with BOM encoding so it opens correctly in Excel without encoding issues.

---

## Build options

| Command | Description |
|---|---|
| `make run` | Build (if needed) then run |
| `make build` | Compile the binary only |
| `make clean` | Remove the compiled binary |

---

## Notes

- Setting a very high limit (above 500) increases RAM usage and the chance of a temporary IP block from Google. Use reasonable limits and add delay if scraping many keywords.
- If pages time out frequently, try increasing the delay value.
- This tool is intended for research and personal use. Do not use it for mass spam or anything that violates Google's terms of service.

---

## Dependencies

| Package | Purpose |
|---|---|
| [go-rod/rod](https://github.com/go-rod/rod) | Browser automation via Chrome DevTools Protocol |
| [xuri/excelize](https://github.com/xuri/excelize) | Read and write Excel (.xlsx) files |
| [schollz/progressbar](https://github.com/schollz/progressbar) | Terminal progress bar |
