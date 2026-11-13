package mocks

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/jarcoal/httpmock"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
)

const (
	// SourceAppID owns every relationship in the mock fixture.
	SourceAppID = "11111111-1111-1111-1111-111111111111"
)

var (
	appRelationshipsURL = regexp.MustCompile(
		`^https://graph\.microsoft\.com/beta/deviceAppManagement/mobileApps/([0-9a-fA-F-]{36})/relationships`,
	)
	rootRelationshipURL = regexp.MustCompile(
		`^https://graph\.microsoft\.com/beta/deviceAppManagement/mobileAppRelationships/([^/?]+)$`,
	)
	targetTypeFilter = regexp.MustCompile(`^targetType eq '(child|parent)'$`)
	isofFilter       = regexp.MustCompile(
		`^isof\('microsoft\.graph\.(mobileAppDependency|mobileAppSupersedence)'\)$`,
	)
)

func init() {
	mocks.GlobalRegistry.Register("mobile_app_relationship", &MobileAppRelationshipMock{})
}

// MobileAppRelationshipMock serves the per-app relationships navigation and the root by-ID route.
type MobileAppRelationshipMock struct{}

var _ mocks.MockRegistrar = (*MobileAppRelationshipMock)(nil)

func loadJSON(name string) (map[string]any, error) {
	content, err := helpers.ParseJSONFile("../tests/responses/validate_get/" + name)
	if err != nil {
		return nil, fmt.Errorf("load mock %s: %w", name, err)
	}
	var obj map[string]any
	if err := json.Unmarshal([]byte(content), &obj); err != nil {
		return nil, fmt.Errorf("decode mock %s: %w", name, err)
	}
	return obj, nil
}

// jsonResponse wraps httpmock's constructor so mock failures carry context.
func jsonResponse(status int, body any) (*http.Response, error) {
	resp, err := httpmock.NewJsonResponse(status, body)
	if err != nil {
		return nil, fmt.Errorf("create mock JSON response: %w", err)
	}
	return resp, nil
}

func graphError(status int, code, message string) (*http.Response, error) {
	return jsonResponse(status, map[string]any{
		"error": map[string]any{"code": code, "message": message},
	})
}

// applyFilter mimics the subset of $filter the Graph navigation supports in tests.
func applyFilter(items []any, filter string) ([]any, bool) {
	if filter == "" {
		return items, true
	}
	var key, want string
	if m := targetTypeFilter.FindStringSubmatch(filter); m != nil {
		key, want = "targetType", m[1]
	} else if m := isofFilter.FindStringSubmatch(filter); m != nil {
		key, want = "@odata.type", "#microsoft.graph."+m[1]
	} else {
		return nil, false
	}
	filtered := make([]any, 0, len(items))
	for _, item := range items {
		if entry, ok := item.(map[string]any); ok && entry[key] == want {
			filtered = append(filtered, entry)
		}
	}
	return filtered, true
}

func (m *MobileAppRelationshipMock) RegisterMocks() {
	// GET /deviceAppManagement/mobileApps/{id}/relationships[?$filter=...]
	httpmock.RegisterRegexpResponder(
		"GET",
		appRelationshipsURL,
		func(req *http.Request) (*http.Response, error) {
			appID := appRelationshipsURL.FindStringSubmatch(req.URL.String())[1]
			if !strings.EqualFold(appID, SourceAppID) {
				return jsonResponse(200, map[string]any{"value": []any{}})
			}
			obj, err := loadJSON("get_app_relationships.json")
			if err != nil {
				return nil, err
			}
			items, ok := applyFilter(obj["value"].([]any), req.URL.Query().Get("$filter"))
			if !ok {
				return graphError(400, "BadRequest", "Invalid filter clause")
			}
			obj["value"] = items
			return jsonResponse(200, obj)
		},
	)

	// GET /deviceAppManagement/mobileAppRelationships/{id} for opaque legacy IDs.
	httpmock.RegisterRegexpResponder(
		"GET",
		rootRelationshipURL,
		func(req *http.Request) (*http.Response, error) {
			if rootRelationshipURL.FindStringSubmatch(req.URL.String())[1] != "legacy-relationship-001" {
				return graphError(404, "ResourceNotFound", "Resource not found")
			}
			obj, err := loadJSON("get_relationship_by_opaque_id.json")
			if err != nil {
				return nil, err
			}
			return jsonResponse(200, obj)
		},
	)
}

func (m *MobileAppRelationshipMock) RegisterErrorMocks() {
	forbidden := func(*http.Request) (*http.Response, error) {
		return graphError(403, "Forbidden", "Insufficient privileges to complete the operation.")
	}
	httpmock.RegisterRegexpResponder("GET", appRelationshipsURL, forbidden)
	httpmock.RegisterRegexpResponder("GET", rootRelationshipURL, forbidden)
}

func (m *MobileAppRelationshipMock) CleanupMockState() {}
