package scraper

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
	"github.com/schollz/progressbar/v3"
)

const userAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) " +
	"AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"

// ScrapeGoogleMaps navigates Google Maps for the given keyword and returns
// up to limit Business records. headful shows the browser window; delay adds
// a pause (seconds) between visiting each place detail page.
func ScrapeGoogleMaps(keyword string, limit int, headful bool, delay float64) ([]Business, error) {
	fmt.Printf("\nStarting scraper for keyword: '%s'\n", keyword)
	mode := "headless"
	if headful {
		mode = "headful"
	}
	fmt.Printf("Target: %d results, browser mode: %s\n\n", limit, mode)

	// Build launcher
	l := launcher.New().
		Headless(!headful).
		Set("disable-gpu", "").
		Set("disable-dev-shm-usage", "").
		Set("no-sandbox", "").
		Set("disable-blink-features", "AutomationControlled")

	u, err := l.Launch()
	if err != nil {
		return nil, fmt.Errorf("failed to launch browser: %w", err)
	}

	browser := rod.New().ControlURL(u).MustConnect()
	defer browser.MustClose()

	page, err := browser.Page(proto.TargetCreateTarget{URL: "about:blank"})
	if err != nil {
		return nil, fmt.Errorf("failed to open page: %w", err)
	}

	// Set user-agent and viewport
	if err := page.SetUserAgent(&proto.NetworkSetUserAgentOverride{UserAgent: userAgent}); err != nil {
		return nil, fmt.Errorf("failed to set user agent: %w", err)
	}
	if err := page.SetViewport(&proto.EmulationSetDeviceMetricsOverride{
		Width: 1280, Height: 800, DeviceScaleFactor: 1,
	}); err != nil {
		return nil, fmt.Errorf("failed to set viewport: %w", err)
	}

	searchURL := "https://www.google.com/maps/search/" + encodeKeyword(keyword)
	fmt.Println("Opening Google Maps and searching...")

	if err := page.Navigate(searchURL); err != nil {
		return nil, fmt.Errorf("navigation error: %w", err)
	}

	// Wait for feed or h1 (single place redirect)
	err = rod.Try(func() {
		page.Timeout(20 * time.Second).MustElement(`div[role="feed"], h1`)
	})
	if err != nil {
		return nil, fmt.Errorf("timed out waiting for results: %w", err)
	}

	collectedURLs := map[string]struct{}{}
	feedSelector := `div[role="feed"]`

	hasFeed := false
	if els, e := page.Elements(feedSelector); e == nil && len(els) > 0 {
		hasFeed = true
	}

	if hasFeed {
		fmt.Println("Found business listing, scrolling to collect URLs...")

		bar := progressbar.NewOptions(limit,
			progressbar.OptionSetDescription("Collecting URLs"),
			progressbar.OptionSetWidth(40),
			progressbar.OptionShowCount(),
			progressbar.OptionSetTheme(progressbar.Theme{
				Saucer:        "█",
				SaucerPadding: "░",
				BarStart:      "[",
				BarEnd:        "]",
			}),
		)

		noChangeCount := 0
		lastHeight := 0

		for len(collectedURLs) < limit {
			links, _ := page.Elements(`a[href*="/maps/place/"]`)
			for _, link := range links {
				href, _ := link.Attribute("href")
				if href == nil {
					continue
				}
				clean := cleanURL(*href)
				if _, seen := collectedURLs[clean]; !seen {
					collectedURLs[clean] = struct{}{}
					bar.Add(1) //nolint:errcheck
					if len(collectedURLs) >= limit {
						break
					}
				}
			}
			if len(collectedURLs) >= limit {
				break
			}

			// Scroll the feed
			currentHeight := 0
			scrollJS := fmt.Sprintf(`() => {
				const el = document.querySelector('%s');
				if (!el) return 0;
				el.scrollTop = el.scrollHeight;
				return el.scrollHeight;
			}`, feedSelector)
			result, scrollErr := page.Eval(scrollJS)
			if scrollErr == nil && result != nil {
				currentHeight = int(result.Value.Int())
			}

			page.WaitIdle(5 * time.Second) //nolint:errcheck
			time.Sleep(2 * time.Second)

			if currentHeight == lastHeight {
				noChangeCount++
				if noChangeCount >= 4 {
					fmt.Println("\nReached end of search results.")
					break
				}
			} else {
				noChangeCount = 0
				lastHeight = currentHeight
			}
		}
		bar.Finish() //nolint:errcheck
		fmt.Println()

	} else {
		// Single place redirect
		currentURL := page.MustInfo().URL
		if strings.Contains(currentURL, "/maps/place/") {
			fmt.Println("Redirected to a single business page.")
			collectedURLs[cleanURL(currentURL)] = struct{}{}
		} else {
			return nil, fmt.Errorf("could not find business listing or detail page")
		}
	}

	totalURLs := make([]string, 0, len(collectedURLs))
	for u := range collectedURLs {
		totalURLs = append(totalURLs, u)
	}

	fmt.Printf("\nCollected %d business URLs. Extracting details...\n", len(totalURLs))

	bar := progressbar.NewOptions(len(totalURLs),
		progressbar.OptionSetDescription("Extracting details"),
		progressbar.OptionSetWidth(40),
		progressbar.OptionShowCount(),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "█",
			SaucerPadding: "░",
			BarStart:      "[",
			BarEnd:        "]",
		}),
	)

	var results []Business

	for i, u := range totalURLs {
		if i > 0 {
			time.Sleep(time.Duration(delay * float64(time.Second)))
		}

		b, extractErr := extractBusiness(page, u)
		if extractErr != nil {
			fmt.Printf("\nWARN: failed to extract %s: %v\n", u, extractErr)
			bar.Add(1) //nolint:errcheck
			continue
		}
		results = append(results, b)
		bar.Add(1) //nolint:errcheck
	}

	bar.Finish() //nolint:errcheck
	fmt.Println()

	return results, nil
}

// extractBusiness navigates to a single place URL and extracts all fields.
func extractBusiness(page *rod.Page, placeURL string) (Business, error) {
	if err := page.Navigate(placeURL); err != nil {
		return Business{}, fmt.Errorf("navigation failed: %w", err)
	}

	err := rod.Try(func() {
		page.Timeout(10 * time.Second).MustElement("h1")
	})
	if err != nil {
		return Business{}, fmt.Errorf("h1 not found: %w", err)
	}

	b := Business{MapsURL: placeURL}

	// 1. Business Name
	b.Name = innerText(page, "h1")

	// 2. Rating
	b.Rating = innerText(page, `div.F7nice span span[aria-hidden="true"]`)
	if b.Rating == "" {
		// Fallback: aria-label containing "star"
		if el := firstElement(page, `span[aria-label*="star"]`); el != nil {
			ariaLabel, _ := el.Attribute("aria-label")
			if ariaLabel != nil {
				re := regexp.MustCompile(`(\d[\.,]\d)`)
				if m := re.FindString(*ariaLabel); m != "" {
					b.Rating = strings.ReplaceAll(m, ",", ".")
				}
			}
		}
	}

	// 3. Reviews Count
	if el := firstElement(page, `div.F7nice span[aria-label*="review"]`); el != nil {
		ariaLabel, _ := el.Attribute("aria-label")
		if ariaLabel != nil {
			re := regexp.MustCompile(`([\d\.,]+)`)
			if m := re.FindString(*ariaLabel); m != "" {
				b.ReviewsCount = strings.NewReplacer(".", "", ",", "").Replace(m)
			}
		}
	}
	if b.ReviewsCount == "" {
		raw := innerText(page, `div.F7nice span:nth-child(2) span span`)
		b.ReviewsCount = strings.NewReplacer("(", "", ")", "", ".", "", ",", "").Replace(raw)
	}

	// 4. Category
	b.Category = innerText(page, `button[jsaction="pane.rating.category"]`)
	if b.Category == "" {
		b.Category = innerText(page, `button[class*="DkEaCb"]`)
	}

	// 5. Address
	b.Address = innerText(page, `button[data-item-id="address"]`)

	// 6. Website
	if el := firstElement(page, `*[data-item-id="authority"]`); el != nil {
		href, _ := el.Attribute("href")
		if href != nil && *href != "" {
			b.Website = *href
		} else {
			b.Website = strings.TrimSpace(mustInnerText(el))
		}
	}

	// 7. Phone
	phone := innerText(page, `button[data-item-id^="phone:tel:"]`)
	for _, prefix := range []string{"Telepon:", "Phone:", "Tel:"} {
		if strings.HasPrefix(strings.ToLower(phone), strings.ToLower(prefix)) {
			phone = strings.TrimSpace(phone[len(prefix):])
			break
		}
	}
	// Strip non-digit chars except leading +
	re := regexp.MustCompile(`[^\d+]`)
	phone = re.ReplaceAllString(phone, "")
	phone = strings.TrimPrefix(phone, "+")
	b.Phone = phone

	return b, nil
}

// --- helpers ---

func firstElement(page *rod.Page, selector string) *rod.Element {
	els, err := page.Elements(selector)
	if err != nil || len(els) == 0 {
		return nil
	}
	return els[0]
}

func innerText(page *rod.Page, selector string) string {
	el := firstElement(page, selector)
	if el == nil {
		return ""
	}
	return mustInnerText(el)
}

func mustInnerText(el *rod.Element) string {
	text, err := el.Text()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(text)
}

// cleanURL strips query parameters from a Google Maps place URL.
func cleanURL(raw string) string {
	if idx := strings.Index(raw, "?"); idx != -1 {
		return raw[:idx]
	}
	return raw
}
