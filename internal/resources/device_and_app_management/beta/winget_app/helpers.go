package graphBetaWinGetApp

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	utils "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/utilities"

	"github.com/PuerkitoBio/goquery"
)

// FetchStoreAppDetails fetches and parses details from the Microsoft Store webpage based on the packageIdentifier
// It also extracts the icon URL, app description, and publisher
func FetchStoreAppDetails(packageIdentifier string) (string, string, string, string, error) {
	// Construct the URL using the packageIdentifier
	storeURL := fmt.Sprintf("https://apps.microsoft.com/detail/%s?hl=en-gb&gl=GB", strings.ToLower(packageIdentifier))

	// Fetch the webpage
	resp, err := http.Get(storeURL)
	if err != nil {
		return "", "", "", "", fmt.Errorf("failed to fetch store page: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", "", "", "", fmt.Errorf("received non-OK response code: %d", resp.StatusCode)
	}

	// Read the raw HTML content
	rawHTML, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", "", "", fmt.Errorf("failed to read raw HTML content: %v", err)
	}

	// Print raw HTML for debugging purposes
	if utils.IsDebugMode() {
		fmt.Println("---- Start of Raw HTML ----")
		fmt.Println(string(rawHTML))
		fmt.Println("---- End of Raw HTML ----")
	}

	// Parse the webpage using goquery
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(rawHTML)))
	if err != nil {
		return "", "", "", "", fmt.Errorf("failed to parse store page: %v", err)
	}

	// Extract the title and publisher
	title, publisher, err := extractTitleAndPublisher(doc)
	if err != nil {
		return "", "", "", "", err
	}

	// Extract the image link
	imageURL, err := extractIconURL(doc)
	if err != nil {
		return "", "", "", "", err
	}

	// Extract the app description
	description, err := extractMicrosoftStoreAppDescription(doc)
	if err != nil {
		return "", "", "", "", err
	}

	return title, imageURL, description, publisher, nil
}

// extractTitleAndPublisher extracts the title and publisher from the parsed HTML document
func extractTitleAndPublisher(doc *goquery.Document) (string, string, error) {
	var title string
	var publisher string

	// Try to find title and publisher normally
	doc.Find("h1, h5").EachWithBreak(func(i int, s *goquery.Selection) bool {
		title = strings.TrimSpace(s.Text())
		if title != "" {
			// Try to find the next sibling <span>
			siblingSpan := s.Next()
			if siblingSpan.Is("span") {
				publisher = strings.TrimSpace(siblingSpan.Text())
			}
			return false // Stop iterating
		}
		return true // Continue iterating
	})

	// If title or publisher is empty, check inside <noscript> tags
	if title == "" || publisher == "" {
		doc.Find("noscript").EachWithBreak(func(i int, s *goquery.Selection) bool {
			// Parse the content inside <noscript>
			innerDoc, err := goquery.NewDocumentFromReader(strings.NewReader(s.Text()))
			if err != nil {
				return true // Continue iterating
			}
			innerDoc.Find("h1, h5").EachWithBreak(func(i int, s *goquery.Selection) bool {
				title = strings.TrimSpace(s.Text())
				if title != "" {
					// Try to find the next sibling <span>
					siblingSpan := s.Next()
					if siblingSpan.Is("span") {
						publisher = strings.TrimSpace(siblingSpan.Text())
					}
					return false // Stop iterating
				}
				return true // Continue iterating
			})
			return title == "" || publisher == ""
		})
	}

	if title == "" {
		return "", "", fmt.Errorf("title not found")
	}

	if publisher == "" {
		return "", "", fmt.Errorf("publisher not found")
	}

	return title, publisher, nil
}

// extractIconURL extracts the icon URL from the parsed HTML document
func extractIconURL(doc *goquery.Document) (string, error) {
	var imageURL string
	imageRegex := regexp.MustCompile(`^https://store-images\.s-microsoft\.com/image/apps\.[a-zA-Z0-9\.\-]+`)

	// Check for <img> tags normally
	doc.Find("img").EachWithBreak(func(i int, s *goquery.Selection) bool {
		src, exists := s.Attr("src")
		if exists && imageRegex.MatchString(src) {
			imageURL = src
			return false // Stop iterating
		}
		return true // Continue iterating
	})

	// If imageURL is still empty, check inside <noscript> tags
	if imageURL == "" {
		doc.Find("noscript").EachWithBreak(func(i int, s *goquery.Selection) bool {
			// Parse the content inside <noscript>
			innerDoc, err := goquery.NewDocumentFromReader(strings.NewReader(s.Text()))
			if err != nil {
				return true // Continue iterating
			}
			innerDoc.Find("img").EachWithBreak(func(j int, imgTag *goquery.Selection) bool {
				src, exists := imgTag.Attr("src")
				if exists && imageRegex.MatchString(src) {
					imageURL = src
					return false // Stop iterating
				}
				return true // Continue iterating
			})
			return imageURL == "" // Continue iterating if imageURL is still empty
		})
	}

	if imageURL == "" {
		return "", fmt.Errorf("image link matching pattern not found")
	}

	return imageURL, nil
}

// extractMicrosoftStoreAppDescription extracts the app description from the parsed HTML document
func extractMicrosoftStoreAppDescription(doc *goquery.Document) (string, error) {
	var description string

	// Try to find the description directly
	doc.Find("pre").EachWithBreak(func(i int, s *goquery.Selection) bool {
		descText := strings.TrimSpace(s.Text())
		if descText != "" {
			description = descText
			return false // Stop iterating
		}
		return true // Continue iterating
	})

	// If description is still empty, check inside <noscript> tags
	if description == "" {
		doc.Find("noscript").EachWithBreak(func(i int, s *goquery.Selection) bool {
			// Parse the content inside <noscript>
			innerDoc, err := goquery.NewDocumentFromReader(strings.NewReader(s.Text()))
			if err != nil {
				return true // Continue iterating
			}
			innerDoc.Find("pre").EachWithBreak(func(j int, preTag *goquery.Selection) bool {
				descText := strings.TrimSpace(preTag.Text())
				if descText != "" {
					description = descText
					return false // Stop iterating
				}
				return true // Continue iterating
			})
			return description == "" // Continue iterating if description is still empty
		})
	}

	if description == "" {
		return "", fmt.Errorf("description not found")
	}

	return description, nil
}

// DownloadImage downloads an image from a given URL and returns it as a byte slice
func DownloadImage(url string) ([]byte, error) {
	// Perform an HTTP GET request to download the image
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// Allow up to 10 redirects
			if len(via) >= 10 {
				return fmt.Errorf("too many redirects")
			}
			return nil
		},
	}

	// Retry the download if redirected
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch image: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-OK response code: %d", resp.StatusCode)
	}

	// Read the image data
	imageData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read image data: %v", err)
	}

	return imageData, nil
}