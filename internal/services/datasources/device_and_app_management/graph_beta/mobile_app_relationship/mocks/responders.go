package mocks

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strings"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/jarcoal/httpmock"
)

// RegisterMobileAppRelationshipMockResponders registers all the mock HTTP responders for mobile app relationship tests
func RegisterMobileAppRelationshipMockResponders() {
	// Mock for getting all mobile app relationships
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/deviceAppManagement/mobileAppRelationships",
		func(req *http.Request) (*http.Response, error) {
			queryParams := req.URL.Query()
			filter := queryParams.Get("$filter")
			top := queryParams.Get("$top")
			skip := queryParams.Get("$skip")

			// Handle OData filter requests
			if filter != "" {
				if strings.Contains(filter, "sourceId eq") {
					jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_get/get_mobile_app_relationships_by_source_id.json")
					var responseObj map[string]any
					json.Unmarshal([]byte(jsonStr), &responseObj)
					return httpmock.NewJsonResponse(200, responseObj)
				}
				if strings.Contains(filter, "targetId eq") {
					jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_get/get_mobile_app_relationships_by_target_id.json")
					var responseObj map[string]any
					json.Unmarshal([]byte(jsonStr), &responseObj)
					return httpmock.NewJsonResponse(200, responseObj)
				}
			}

			// Handle $top and $skip for OData pagination
			if top != "" || skip != "" {
				jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_get/get_mobile_app_relationships_all.json")
				var responseObj map[string]any
				json.Unmarshal([]byte(jsonStr), &responseObj)
				return httpmock.NewJsonResponse(200, responseObj)
			}

			// Default: Return all mobile app relationships
			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_get/get_mobile_app_relationships_all.json")
			var responseObj map[string]any
			json.Unmarshal([]byte(jsonStr), &responseObj)
			return httpmock.NewJsonResponse(200, responseObj)
		},
	)

	// Mock for getting a specific mobile app relationship by ID
	httpmock.RegisterRegexpResponder("GET", regexp.MustCompile(`https://graph\.microsoft\.com/beta/deviceAppManagement/mobileAppRelationships/[a-fA-F0-9-]+$`),
		func(req *http.Request) (*http.Response, error) {
			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_get/get_mobile_app_relationship_by_id.json")
			var responseObj map[string]any
			json.Unmarshal([]byte(jsonStr), &responseObj)
			return httpmock.NewJsonResponse(200, responseObj)
		},
	)
}
