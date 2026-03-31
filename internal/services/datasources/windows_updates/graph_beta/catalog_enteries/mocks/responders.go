package mocks

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/jarcoal/httpmock"
)

var mockState struct {
	sync.Mutex
	catalogEntries map[string]map[string]any
}

func init() {
	mockState.catalogEntries = make(map[string]map[string]any)
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))
	mocks.GlobalRegistry.Register("windows_update_catalog", &WindowsUpdateCatalogMock{})
}

type WindowsUpdateCatalogMock struct{}

var _ mocks.MockRegistrar = (*WindowsUpdateCatalogMock)(nil)

func (m *WindowsUpdateCatalogMock) RegisterMocks() {
	mockState.Lock()
	mockState.catalogEntries = make(map[string]map[string]any)
	mockState.Unlock()

	m.registerListCatalogEntriesResponder()
	m.registerGetCatalogEntryByIdResponder()
}

func (m *WindowsUpdateCatalogMock) registerListCatalogEntriesResponder() {
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/admin/windows/updates/catalog/entries`, func(req *http.Request) (*http.Response, error) {
		queryParams, _ := url.ParseQuery(req.URL.RawQuery)
		filter := queryParams.Get("$filter")

		// Load JSON response from file
		_, filename, _, _ := runtime.Caller(0)
		sourceDir := filepath.Dir(filename)
		responsesPath := filepath.Join(sourceDir, "..", "tests", "responses", "validate_get", "get_catalog_entries_all.json")

		jsonData, err := os.ReadFile(responsesPath)
		if err != nil {
			return httpmock.NewStringResponse(500, fmt.Sprintf(`{"error":{"code":"InternalServerError","message":"Failed to load JSON response file: %s"}}`, err.Error())), nil
		}

		var responseObj map[string]any
		if err := json.Unmarshal(jsonData, &responseObj); err != nil {
			return httpmock.NewStringResponse(500, fmt.Sprintf(`{"error":{"code":"InternalServerError","message":"Failed to parse JSON response: %s"}}`, err.Error())), nil
		}

		// Apply filter if present
		if filter != "" {
			entries := responseObj["value"].([]any)
			var filteredEntries []any

			for _, entry := range entries {
				entryMap := entry.(map[string]any)

				// Handle different filter types
				if strings.Contains(filter, "microsoft.graph.windowsUpdates.featureUpdateCatalogEntry") {
					if odataType, ok := entryMap["@odata.type"].(string); ok && strings.Contains(odataType, "featureUpdateCatalogEntry") {
						filteredEntries = append(filteredEntries, entry)
					}
				} else if strings.Contains(filter, "microsoft.graph.windowsUpdates.qualityUpdateCatalogEntry") {
					if odataType, ok := entryMap["@odata.type"].(string); ok && strings.Contains(odataType, "qualityUpdateCatalogEntry") {
						filteredEntries = append(filteredEntries, entry)
					}
				} else if strings.Contains(filter, "displayName") {
					// Handle display name filter
					filteredEntries = append(filteredEntries, entry)
				}
			}

			responseObj["value"] = filteredEntries
		}

		resp, err := httpmock.NewJsonResponse(200, responseObj)
		if err != nil {
			return nil, fmt.Errorf("failed to create mock JSON response: %w", err)
		}
		return resp, nil
	})
}

func (m *WindowsUpdateCatalogMock) registerGetCatalogEntryByIdResponder() {
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/admin/windows/updates/catalog/entries/[0-9a-fA-F-]+$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		entryId := parts[len(parts)-1]

		// Load all entries and find the matching one
		_, filename, _, _ := runtime.Caller(0)
		sourceDir := filepath.Dir(filename)
		responsesPath := filepath.Join(sourceDir, "..", "tests", "responses", "validate_get", "get_catalog_entries_all.json")

		jsonData, err := os.ReadFile(responsesPath)
		if err != nil {
			return httpmock.NewStringResponse(500, fmt.Sprintf(`{"error":{"code":"InternalServerError","message":"Failed to load JSON response file: %s"}}`, err.Error())), nil
		}

		var responseObj map[string]any
		if err := json.Unmarshal(jsonData, &responseObj); err != nil {
			return httpmock.NewStringResponse(500, fmt.Sprintf(`{"error":{"code":"InternalServerError","message":"Failed to parse JSON response: %s"}}`, err.Error())), nil
		}

		entries := responseObj["value"].([]any)
		for _, entry := range entries {
			entryMap := entry.(map[string]any)
			if id, ok := entryMap["id"].(string); ok && id == entryId {
				resp, err := httpmock.NewJsonResponse(200, entryMap)
				if err != nil {
					return nil, fmt.Errorf("failed to create mock JSON response: %w", err)
				}
				return resp, nil
			}
		}

		return httpmock.NewStringResponse(404, `{"error":{"code":"NotFound","message":"Catalog entry not found"}}`), nil
	})
}

func (m *WindowsUpdateCatalogMock) RegisterErrorMocks() {
	mockState.Lock()
	mockState.catalogEntries = make(map[string]map[string]any)
	mockState.Unlock()

	m.registerListCatalogEntriesErrorResponder()
}

func (m *WindowsUpdateCatalogMock) registerListCatalogEntriesErrorResponder() {
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/admin/windows/updates/catalog/entries`, func(req *http.Request) (*http.Response, error) {
		errorObj := map[string]any{
			"error": map[string]any{
				"code":    "Forbidden",
				"message": "Insufficient privileges to complete the operation.",
			},
		}
		resp, err := httpmock.NewJsonResponse(403, errorObj)
		if err != nil {
			return nil, fmt.Errorf("failed to create mock JSON response: %w", err)
		}
		return resp, nil
	})
}

func (m *WindowsUpdateCatalogMock) CleanupMockState() {
	mockState.Lock()
	mockState.catalogEntries = make(map[string]map[string]any)
	mockState.Unlock()
}
