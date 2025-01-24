package graphBetaApplications

import (
	"context"
	"encoding/json"
	"fmt"
	"html"
	"io"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// htmlSelector represents a CSS selector and its attribute to extract
type htmlSelector struct {
	selector string // CSS selector (e.g., "meta[name='description']")
	attr     string // HTML attribute to extract (e.g., "content"), empty for text content
}

// FetchStoreAppDetails fetches and parses details from the Microsoft Store webpage
func FetchStoreAppDetails(ctx context.Context, packageIdentifier string) (string, string, string, string, error) {
	storeURL := fmt.Sprintf("https://apps.microsoft.com/detail/%s?hl=en-gb&gl=GB", strings.ToLower(packageIdentifier))

	doc, err := getHTML(ctx, storeURL)
	if err != nil {
		return "", "", "", "", err
	}

	// Extract title
	titleSelectors := []htmlSelector{
		{"meta[property='og:title']", "content"},
		{"meta[name='twitter:title']", "content"},
		{"title", ""},
		{"h1", ""},
	}
	fullTitle, err := trySelectors(ctx, doc, titleSelectors, "title")
	if err != nil {
		return "", "", "", "", err
	}
	title := cleanTitle(fullTitle)
	tflog.Debug(ctx, "Cleaned title", map[string]interface{}{
		"fullTitle": fullTitle,
		"title":     title,
	})

	// Extract app icon image URL
	imageSelectors := []htmlSelector{
		{"meta[property='og:image']", "content"},
		{"meta[name='twitter:image']", "content"},
		{"img.product-image", "src"},
		{"img[src*='store-images.s-microsoft.com']", "src"},
	}
	imageURL, err := trySelectors(ctx, doc, imageSelectors, "image URL")
	if err != nil {
		return "", "", "", "", err
	}

	// Extract app description
	descriptionSelectors := []htmlSelector{
		{"meta[name='description']", "content"},
		{"meta[property='og:description']", "content"},
		{"meta[name='twitter:description']", "content"},
	}
	description, err := trySelectors(ctx, doc, descriptionSelectors, "description")
	if err != nil {
		// Try JSON-LD as fallback
		description = parseJSONLD(ctx, doc, "description")
		if description == "" {
			return "", "", "", "", err
		}
	}
	description = cleanDescription(description)

	// Extract app publisher
	publisher := parseJSONLD(ctx, doc, "author")
	if publisher == "" {
		publisherSelectors := []htmlSelector{
			{"meta[name='application-name']", "content"},
			{"meta[property='og:site_name']", "content"},
		}
		publisher, err = trySelectors(ctx, doc, publisherSelectors, "publisher")
		if err != nil {
			return "", "", "", "", err
		}
	}

	return title, imageURL, description, publisher, nil
}

// sharedClient returns a configured HTTP client with standard settings
func sharedClient(ctx context.Context) *http.Client {
	return &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 10 {
				tflog.Error(ctx, "Too many redirects", map[string]interface{}{
					"redirectCount": len(via),
				})
				return fmt.Errorf("too many redirects")
			}
			tflog.Debug(ctx, "Following redirect", map[string]interface{}{
				"redirectCount": len(via),
				"location":      req.URL.String(),
			})
			return nil
		},
	}
}

// getHTML makes an HTTP request and returns the parsed HTML document
func getHTML(ctx context.Context, url string) (*goquery.Document, error) {
	client := sharedClient(ctx)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		tflog.Error(ctx, "Failed to create request", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Add("Accept-Language", "en-US,en;q=0.5")

	resp, err := client.Do(req)
	if err != nil {
		tflog.Error(ctx, "Failed to fetch URL", map[string]interface{}{
			"error": err.Error(),
			"url":   url,
		})
		return nil, fmt.Errorf("failed to fetch URL: %v", err)
	}
	defer resp.Body.Close()

	tflog.Debug(ctx, "Received response", map[string]interface{}{
		"status":     resp.Status,
		"statusCode": resp.StatusCode,
		"headers":    fmt.Sprintf("%v", resp.Header),
	})

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-OK response code: %d", resp.StatusCode)
	}

	rawHTML, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	return goquery.NewDocumentFromReader(strings.NewReader(string(rawHTML)))
}

// trySelectors  is used to extract content from HTML using CSS selectors.
func trySelectors(ctx context.Context, doc *goquery.Document, selectors []htmlSelector, logDesc string) (string, error) {
	tflog.Debug(ctx, fmt.Sprintf("Trying selectors for %s", logDesc))

	var foundContent string
	for _, sel := range selectors {
		doc.Find(sel.selector).Each(func(i int, s *goquery.Selection) {
			if foundContent != "" {
				return
			}

			var content string
			if sel.attr != "" {
				var exists bool
				content, exists = s.Attr(sel.attr)
				if !exists {
					return
				}
			} else {
				content = s.Text()
			}

			content = strings.TrimSpace(content)
			if content != "" {
				tflog.Debug(ctx, fmt.Sprintf("Found %s", logDesc), map[string]interface{}{
					"selector": sel.selector,
					"content":  content,
				})
				foundContent = content
			}
		})

		if foundContent != "" {
			break
		}
	}

	if foundContent == "" {
		return "", fmt.Errorf("%s not found", logDesc)
	}

	return foundContent, nil
}

// parseJSONLD attempts to parse JSON-LD data and extract a specific field.
// The Microsoft Store uses JSON-LD to include standardized metadata about apps following schema.org vocabulary.
// The data typically includes: App name, Description, Publisher/Author, Images, etc.
func parseJSONLD(ctx context.Context, doc *goquery.Document, field string) string {
	var content string
	doc.Find("script[type='application/ld+json']").EachWithBreak(func(i int, s *goquery.Selection) bool {
		if content != "" {
			return false
		}

		var data map[string]interface{}
		if err := json.Unmarshal([]byte(s.Text()), &data); err == nil {
			if value, ok := data[field].(string); ok {
				content = value
				tflog.Debug(ctx, fmt.Sprintf("Found %s in JSON-LD", field), map[string]interface{}{
					"content": content,
				})
				return false
			}
		}
		return true
	})
	return content
}

// cleanDescription Decode all HTML entities normalizes and cleans up description text
func cleanDescription(description string) string {
	if description == "" {
		return ""
	}

	description = html.UnescapeString(description)

	description = strings.ReplaceAll(description, "\r\n", "\n")

	for strings.Contains(description, "\n\n\n") {
		description = strings.ReplaceAll(description, "\n\n\n", "\n\n")
	}

	return strings.TrimSpace(description)
}

// cleanTitle extracts just the app name from the full title
func cleanTitle(title string) string {
	title = strings.TrimSpace(title)
	if strings.Contains(title, " - ") {
		parts := strings.Split(title, " - ")
		return strings.TrimSpace(parts[0])
	}
	return title
}
