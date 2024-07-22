package provider

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

// configureProxy sets up the HTTP client with proxy settings if useProxy is true and proxyURL is provided.
func configureProxy(useProxy bool, proxyURL string, diagnostics *diag.Diagnostics) *http.Client {
	var transport *http.Transport

	if useProxy {
		proxyUrlParsed, err := url.Parse(proxyURL)
		if err != nil {
			diagnostics.AddError(
				"Invalid Proxy URL",
				fmt.Sprintf("Failed to parse the provided proxy URL '%s': %s. "+
					"Ensure the URL is correctly formatted.", proxyURL, err.Error()),
			)
			return nil
		}
		transport = &http.Transport{
			Proxy: http.ProxyURL(proxyUrlParsed),
		}
	} else {
		transport = &http.Transport{}
	}

	return &http.Client{
		Transport: transport,
	}
}
