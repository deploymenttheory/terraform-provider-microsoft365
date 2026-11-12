package mocks

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"maps"
	"net/http"
	"sync"

	"github.com/jarcoal/httpmock"

	common "github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
)

const (
	endpoint   = "https://graph.microsoft.com/beta/networkaccess/forwardingProfiles"
	InternetID = "72661c0d-027e-4dff-8c76-af103f200903"
	CustomID   = "11111111-2222-3333-4444-555555555555"
)

// Mock serves one default internet profile and one custom profile. Default
// profiles cannot be created or deleted, so POST and DELETE fail the test.
type Mock struct {
	sync.Mutex
	Internet map[string]any
	Patches  []map[string]any
	Writes   int
}

var _ common.MockRegistrar = (*Mock)(nil)

func init() { common.GlobalRegistry.Register("network_default_forwarding_profile", &Mock{}) }

func decode(q *http.Request) (map[string]any, error) {
	var reader io.Reader = q.Body
	if q.Header.Get("Content-Encoding") == "gzip" {
		gz, err := gzip.NewReader(q.Body)
		if err != nil {
			return nil, fmt.Errorf("decompress request: %w", err)
		}
		defer gz.Close()
		reader = gz
	}
	var body map[string]any
	if err := json.NewDecoder(reader).Decode(&body); err != nil {
		return nil, fmt.Errorf("decode request: %w", err)
	}
	return body, nil
}

func initialInternet() map[string]any {
	return map[string]any{
		"id":                       InternetID,
		"name":                     "Internet traffic forwarding profile",
		"description":              "Default traffic forwarding profile for Internet traffic acquisition.",
		"state":                    "disabled",
		"version":                  "1.0.0",
		"trafficForwardingType":    "internet",
		"priority":                 300,
		"isCustomProfile":          false,
		"clientFallbackAction":     "bypass",
		"localNetworkSettings":     nil,
		"forwardingSourceMatching": map[string]any{"sessionType": "any"},
		"assignmentRules": map[string]any{
			"assignments": map[string]any{
				"@odata.type": "#microsoft.graph.networkaccess.forwardingProfileAssignmentsEnumeratedMembership",
				"users":       []any{},
				"groups":      []any{map[string]any{"id": "a6ed41d3-45ac-4e12-8601-bdb56216e84e"}},
				"devices":     []any{},
			},
			"devicePropertyFilters": map[string]any{"platforms": "windows,macOS,android,iOS"},
		},
	}
}

func custom() map[string]any {
	return map[string]any{
		"id":                    CustomID,
		"name":                  "custom",
		"state":                 "disabled",
		"trafficForwardingType": "private",
		"priority":              150,
		"isCustomProfile":       true,
	}
}

func (m *Mock) RegisterMocks() {
	m.CleanupMockState()
	item := endpoint + "/" + InternetID
	httpmock.RegisterResponder("GET", endpoint, func(q *http.Request) (*http.Response, error) {
		m.Lock()
		defer m.Unlock()
		return httpmock.NewJsonResponse(200, map[string]any{"value": []any{m.Internet, custom()}})
	})
	httpmock.RegisterResponder("GET", item, func(q *http.Request) (*http.Response, error) {
		m.Lock()
		defer m.Unlock()
		return httpmock.NewJsonResponse(200, m.Internet)
	})
	httpmock.RegisterResponder(
		"GET",
		endpoint+"/"+CustomID,
		func(q *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(200, custom())
		},
	)
	httpmock.RegisterResponder("PATCH", item, func(q *http.Request) (*http.Response, error) {
		body, err := decode(q)
		if err != nil {
			return nil, err
		}
		m.Lock()
		defer m.Unlock()
		m.Patches = append(m.Patches, maps.Clone(body))
		maps.Copy(m.Internet, body)
		return httpmock.NewJsonResponse(200, m.Internet)
	})
	unexpected := func(q *http.Request) (*http.Response, error) {
		m.Lock()
		defer m.Unlock()
		m.Writes++
		return httpmock.NewStringResponse(
			400,
			`{"error":{"code":"BadRequest","message":"unexpected request"}}`,
		), nil
	}
	httpmock.RegisterResponder("POST", endpoint, unexpected)
	httpmock.RegisterResponder("DELETE", item, unexpected)
}

func (m *Mock) RegisterErrorMocks() {
	httpmock.RegisterResponder("GET", endpoint, httpmock.NewStringResponder(
		403,
		`{"error":{"code":"Forbidden","message":"Synthetic denied request"}}`,
	))
}

func (m *Mock) CleanupMockState() {
	m.Lock()
	defer m.Unlock()
	m.Internet = initialInternet()
	m.Patches = nil
	m.Writes = 0
}
