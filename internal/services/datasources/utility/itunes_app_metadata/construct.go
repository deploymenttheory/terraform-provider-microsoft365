package itunes_app_metadata

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// constructRequest builds the HTTP request for the iTunes Search API
func constructRequest(ctx context.Context, countryCode, searchTerm string) (*http.Request, error) {

	searchURL, err := buildItunesSearchURL(countryCode, searchTerm)
	if err != nil {
		return nil, err
	}

	tflog.Debug(ctx, fmt.Sprintf("Constructing iTunes app metadata request for URL: %s", searchURL))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, searchURL, nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// buildItunesSearchURL constructs the URL for the iTunes Search API
func buildItunesSearchURL(countryCode, searchTerm string) (string, error) {
	params := url.Values{}
	params.Add("country", countryCode)
	params.Add("media", "software")
	params.Add("entity", "software,iPadSoftware")
	params.Add("term", searchTerm)

	return fmt.Sprintf("%s?%s", itunesSearchBaseURL, params.Encode()), nil
}
