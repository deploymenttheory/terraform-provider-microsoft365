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
		`^https://graph\.microsoft\.com/beta/deviceAppManagement/androidManagedAppProtections/([^/?]+)/apps(\?.*)?$`,
	)
	targetAppsURL = regexp.MustCompile(
		`^https://graph\.microsoft\.com/beta/deviceAppManagement/managedAppPolicies/([^/?]+)/targetApps$`,
	)
	errNotAndroid = fmt.Errorf("targetApps entry is not an Android app")
	mockState     = struct {
		sync.Mutex
		policies map[string][]string // policy ID -> targeted package IDs
	}{policies: map[string][]string{}}
)

func init() {
	mocks.GlobalRegistry.Register("android_managed_mobile_app", &AndroidManagedMobileAppMock{})
}

// AndroidManagedMobileAppMock mirrors the live API: per-app routes do not exist,
// targetApps replaces the whole list, and Intune assigns each entry's id and version.
type AndroidManagedMobileAppMock struct{}

var _ mocks.MockRegistrar = (*AndroidManagedMobileAppMock)(nil)

// AppID is the identifier Intune assigns to an Android package entry.
func AppID(pkg string) string { return pkg + ".android" }

// version imitates the server-assigned version: stable per package, never client-set.
func version(pkg string) string {
	h := fnv.New32a()
	_, _ = h.Write([]byte(pkg))
	return fmt.Sprintf("%d", h.Sum32())
}

// SeedApps sets the packages a policy targets before a test runs.
func SeedApps(policyID string, packages ...string) {
	mockState.Lock()
	defer mockState.Unlock()
	mockState.policies[policyID] = append([]string(nil), packages...)
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
			"@odata.type": "#microsoft.graph.androidMobileAppIdentifier",
			"packageId":   pkg,
		},
	}
}

func decodePackages(req *http.Request) ([]string, error) {
	var payload struct {
		Apps []map[string]any `json:"apps"`
	}
	if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
		return nil, fmt.Errorf("decode targetApps request: %w", err)
	}
	packages := make([]string, 0, len(payload.Apps))
	for _, app := range payload.Apps {
		identifier, _ := app["mobileAppIdentifier"].(map[string]any)
		odataType, _ := identifier["@odata.type"].(string)
		pkg, _ := identifier["packageId"].(string)
		if odataType != "#microsoft.graph.androidMobileAppIdentifier" || pkg == "" {
			return nil, errNotAndroid
		}
		packages = append(packages, pkg)
	}
	return packages, nil
}

func (m *AndroidManagedMobileAppMock) RegisterMocks() {
	m.CleanupMockState()

	// GET /deviceAppManagement/androidManagedAppProtections/{id}/apps
	httpmock.RegisterRegexpResponder(
		"GET",
		policyAppsURL,
		func(req *http.Request) (*http.Response, error) {
			policyID := policyAppsURL.FindStringSubmatch(req.URL.String())[1]
			if policyID == MissingPolicyID {
				return graphError(http.StatusNotFound, "ResourceNotFound", "Resource not found")
			}
			mockState.Lock()
			packages := append([]string(nil), mockState.policies[policyID]...)
			mockState.Unlock()
			apps := make([]map[string]any, 0, len(packages))
			for _, pkg := range packages {
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
			packages, err := decodePackages(req)
			if err != nil {
				return graphError(http.StatusBadRequest, "BadRequest", err.Error())
			}
			seen := map[string]bool{}
			unique := make([]string, 0, len(packages))
			for _, pkg := range packages {
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

func (m *AndroidManagedMobileAppMock) RegisterErrorMocks() {
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

func (m *AndroidManagedMobileAppMock) CleanupMockState() {
	mockState.Lock()
	defer mockState.Unlock()
	mockState.policies = map[string][]string{}
}
