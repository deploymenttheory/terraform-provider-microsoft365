package mocks

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strings"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/jarcoal/httpmock"
)

// RegisterApplicationCategoryMockResponders registers all the mock HTTP responders for application category tests
func RegisterApplicationCategoryMockResponders() {
	// Mock for getting all application categories
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/deviceAppManagement/mobileAppCategories",
		func(req *http.Request) (*http.Response, error) {
			queryParams := req.URL.Query()
			filter := queryParams.Get("$filter")
			top := queryParams.Get("$top")

			// Handle ID-based request (single item)
			if strings.Contains(req.URL.Path, "/mobileAppCategories/") {
				jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_get/get_application_category_by_id.json")
				var responseObj map[string]any
				json.Unmarshal([]byte(jsonStr), &responseObj)
				return httpmock.NewJsonResponse(200, responseObj)
			}

			// Handle OData filter requests
			if filter != "" {
				if strings.Contains(filter, "displayName eq") || strings.Contains(filter, "startswith(displayName") {
					jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_get/get_application_categories_by_display_name.json")
					var responseObj map[string]any
					json.Unmarshal([]byte(jsonStr), &responseObj)
					return httpmock.NewJsonResponse(200, responseObj)
				}
			}

			// Handle $top parameter for OData
			if top != "" {
				jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_get/get_application_categories_all.json")
				var responseObj map[string]any
				json.Unmarshal([]byte(jsonStr), &responseObj)
				return httpmock.NewJsonResponse(200, responseObj)
			}

			// Default: Return all application categories
			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_get/get_application_categories_all.json")
			var responseObj map[string]any
			json.Unmarshal([]byte(jsonStr), &responseObj)
			return httpmock.NewJsonResponse(200, responseObj)
		},
	)

	// Mock for getting a specific application category by ID
	httpmock.RegisterRegexpResponder("GET", regexp.MustCompile(`https://graph\.microsoft\.com/beta/deviceAppManagement/mobileAppCategories/[a-fA-F0-9-]+$`),
		func(req *http.Request) (*http.Response, error) {
			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_get/get_application_category_by_id.json")
			var responseObj map[string]any
			json.Unmarshal([]byte(jsonStr), &responseObj)
			return httpmock.NewJsonResponse(200, responseObj)
		},
	)
}
