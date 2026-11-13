package mocks

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"net/http"
	"regexp"
	"sort"
	"strings"
	"sync"

	"github.com/jarcoal/httpmock"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
)

const (
	// PolicyID is the app protection policy every unit test targets.
	PolicyID = "T_00000000-0000-0000-0000-000000000001"
	// MissingPolicyID returns 404 like a deleted policy.
	MissingPolicyID = "T_00000000-0000-0000-0000-000000000404"
)

var (
	policyAppsURL = regexp.MustCompile(
		`^https://graph\.microsoft\.com/beta/deviceAppManagement/iosManagedAppProtections/([^/?]+)/apps(\?.*)?$`,
	)
	targetAppsURL = regexp.MustCompile(
		`^https://graph\.microsoft\.com/beta/deviceAppManagement/managedAppPolicies/([^/?]+)/targetApps$`,
	)
	errNotIOS = fmt.Errorf("targetApps entry is not an iOS app")
	mockState = struct {
		sync.Mutex
		policies map[string][]string // policy ID -> targeted bundle IDs
	}{policies: map[string][]string{}}
)

func init() {
	mocks.GlobalRegistry.Register("ios_managed_mobile_app", &IOSManagedMobileAppMock{})
}

// IOSManagedMobileAppMock mirrors the live API: per-app routes do not exist,
// targetApps replaces the whole list, and Intune assigns each entry's id and version.
type IOSManagedMobileAppMock struct{}

var _ mocks.MockRegistrar = (*IOSManagedMobileAppMock)(nil)

// AppID is the identifier Intune assigns to an iOS bundle entry.
func AppID(pkg string) string { return pkg + ".ios" }

// version imitates the server-assigned version: stable per bundle, never client-set.
func version(pkg string) string {
	h := fnv.New32a()
	_, _ = h.Write([]byte(pkg))
	return fmt.Sprintf("%d", h.Sum32())
}

// SeedApps sets the bundles a policy targets before a test runs.
func SeedApps(policyID string, bundles ...string) {
	mockState.Lock()
	defer mockState.Unlock()
	mockState.policies[policyID] = append([]string(nil), bundles...)
}

// PolicyAppIDs returns the app IDs the policy currently targets, sorted.
func PolicyAppIDs(policyID string) []string {
	mockState.Lock()
	defer mockState.Unlock()
	ids := make([]string, 0, len(mockState.policies[policyID]))
	for _, pkg := range mockState.policies[policyID] {
		ids = append(ids, AppID(pkg))
	}
	sort.Strings(ids)
	return ids
}

func jsonResponse(status int, body any) (*http.Response, error) {
	resp, err := httpmock.NewJsonResponse(status, body)
	if err != nil {
		return nil, fmt.Errorf("create mock JSON response: %w", err)
	}
	return resp, nil
}

func graphError(status int, code, message string) (*http.Response, error) {
	return jsonResponse(
		status,
		map[string]any{"error": map[string]any{"code": code, "message": message}},
	)
}

func appEntry(pkg string) map[string]any {
	return map[string]any{
		"id":      AppID(pkg),
		"version": version(pkg),
		"mobileAppIdentifier": map[string]any{
			"@odata.type": "#microsoft.graph.iosMobileAppIdentifier",
			"bundleId":    pkg,
		},
	}
}

func decodeBundles(req *http.Request) ([]string, error) {
	var payload struct {
		Apps []map[string]any `json:"apps"`
	}
	if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
		return nil, fmt.Errorf("decode targetApps request: %w", err)
	}
	bundles := make([]string, 0, len(payload.Apps))
	for _, app := range payload.Apps {
		identifier, _ := app["mobileAppIdentifier"].(map[string]any)
		odataType, _ := identifier["@odata.type"].(string)
		pkg, _ := identifier["bundleId"].(string)
		if odataType != "#microsoft.graph.iosMobileAppIdentifier" || pkg == "" {
			return nil, errNotIOS
		}
		bundles = append(bundles, pkg)
	}
	return bundles, nil
}

func (m *IOSManagedMobileAppMock) RegisterMocks() {
	m.CleanupMockState()

	// GET /deviceAppManagement/iosManagedAppProtections/{id}/apps
	httpmock.RegisterRegexpResponder(
		"GET",
		policyAppsURL,
		func(req *http.Request) (*http.Response, error) {
			policyID := policyAppsURL.FindStringSubmatch(req.URL.String())[1]
			if policyID == MissingPolicyID {
				return graphError(http.StatusNotFound, "ResourceNotFound", "Resource not found")
			}
			mockState.Lock()
			bundles := append([]string(nil), mockState.policies[policyID]...)
			mockState.Unlock()
			apps := make([]map[string]any, 0, len(bundles))
			for _, pkg := range bundles {
				apps = append(apps, appEntry(pkg))
			}
			return jsonResponse(http.StatusOK, map[string]any{"value": apps})
		},
	)

	// POST /deviceAppManagement/managedAppPolicies/{id}/targetApps replaces the list.
	httpmock.RegisterRegexpResponder(
		"POST",
		targetAppsURL,
		func(req *http.Request) (*http.Response, error) {
			policyID := targetAppsURL.FindStringSubmatch(req.URL.String())[1]
			if policyID == MissingPolicyID {
				return graphError(http.StatusNotFound, "ResourceNotFound", "Resource not found")
			}
			bundles, err := decodeBundles(req)
			if err != nil {
				return graphError(http.StatusBadRequest, "BadRequest", err.Error())
			}
			seen := map[string]bool{}
			unique := make([]string, 0, len(bundles))
			for _, pkg := range bundles {
				if !seen[strings.ToLower(pkg)] {
					seen[strings.ToLower(pkg)] = true
					unique = append(unique, pkg)
				}
			}
			mockState.Lock()
			mockState.policies[policyID] = unique
			mockState.Unlock()
			return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
		},
	)
}

func (m *IOSManagedMobileAppMock) RegisterErrorMocks() {
	m.CleanupMockState()
	forbidden := func(*http.Request) (*http.Response, error) {
		return graphError(
			http.StatusForbidden,
			"Forbidden",
			"Insufficient privileges to complete the operation.",
		)
	}
	httpmock.RegisterRegexpResponder("GET", policyAppsURL, forbidden)
	httpmock.RegisterRegexpResponder("POST", targetAppsURL, forbidden)
}

func (m *IOSManagedMobileAppMock) CleanupMockState() {
	mockState.Lock()
	defer mockState.Unlock()
	mockState.policies = map[string][]string{}
}
