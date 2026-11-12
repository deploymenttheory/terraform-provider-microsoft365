package mocks

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"maps"
	"net/http"
	"strings"
	"sync"

	"github.com/jarcoal/httpmock"

	common "github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
)

const (
	endpoint  = "https://graph.microsoft.com/beta/networkaccess/forwardingProfiles"
	ProfileID = "11111111-2222-3333-4444-555555555555"
	DefaultID = "6e6a0de4-9100-4499-be13-9e9d72da9df0"
)

// Mock emulates the Graph behaviors observed on a live tenant: empty member
// lists are returned for enumerated assignments and platform names are
// normalized to Graph casing.
type Mock struct {
	sync.Mutex
	State   map[string]any
	Posts   []map[string]any
	Patches []map[string]any
	Deletes int
	// PostFailures makes the next POSTs return 504 after storing the profile.
	PostFailures int
}

var _ common.MockRegistrar = (*Mock)(nil)

func init() { common.GlobalRegistry.Register("network_forwarding_profile", &Mock{}) }

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

func normalize(profile map[string]any) {
	rules, ok := profile["assignmentRules"].(map[string]any)
	if !ok {
		return
	}
	if filters, ok := rules["devicePropertyFilters"].(map[string]any); ok {
		if p, ok := filters["platforms"].(string); ok {
			repl := strings.NewReplacer("macos", "macOS", "ios", "iOS")
			filters["platforms"] = repl.Replace(strings.ToLower(p))
		}
	}
}

func defaultProfile() map[string]any {
	return map[string]any{
		"id":                    DefaultID,
		"name":                  "Private access traffic forwarding profile",
		"state":                 "disabled",
		"version":               "1.0.0",
		"trafficForwardingType": "private",
		"priority":              200,
		"isCustomProfile":       false,
		"clientFallbackAction":  "bypass",
	}
}

func (m *Mock) RegisterMocks() {
	m.CleanupMockState()
	item := endpoint + "/" + ProfileID
	httpmock.RegisterResponder("GET", endpoint, func(q *http.Request) (*http.Response, error) {
		m.Lock()
		defer m.Unlock()
		value := []any{defaultProfile()}
		if m.State != nil {
			value = append(value, m.State)
		}
		return httpmock.NewJsonResponse(200, map[string]any{"value": value})
	})
	httpmock.RegisterResponder(
		"GET",
		endpoint+"/"+DefaultID,
		func(q *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(200, defaultProfile())
		},
	)
	httpmock.RegisterResponder("GET", item, func(q *http.Request) (*http.Response, error) {
		m.Lock()
		defer m.Unlock()
		if m.State == nil {
			return httpmock.NewStringResponse(
				404,
				`{"error":{"code":"NotFound","message":"Forwarding profile not found."}}`,
			), nil
		}
		return httpmock.NewJsonResponse(200, m.State)
	})
	httpmock.RegisterResponder("POST", endpoint, func(q *http.Request) (*http.Response, error) {
		body, err := decode(q)
		if err != nil {
			return nil, err
		}
		m.Lock()
		defer m.Unlock()
		m.Posts = append(m.Posts, maps.Clone(body))
		body["id"] = ProfileID
		body["version"] = "1.0.0"
		if _, ok := body["clientFallbackAction"]; !ok {
			body["clientFallbackAction"] = "bypass"
		}
		if _, ok := body["assignmentRules"]; !ok {
			body["assignmentRules"] = map[string]any{
				"assignments": map[string]any{
					"@odata.type": "#microsoft.graph.networkaccess.forwardingProfileAssignmentsNoMembership",
				},
				"devicePropertyFilters": map[string]any{"platforms": "none"},
			}
		}
		normalize(body)
		m.State = body
		if m.PostFailures > 0 {
			m.PostFailures--
			return httpmock.NewStringResponse(
				504,
				`{"error":{"code":"UnknownError","message":"Gateway Timeout"}}`,
			), nil
		}
		return httpmock.NewJsonResponse(201, body)
	})
	httpmock.RegisterResponder("PATCH", item, func(q *http.Request) (*http.Response, error) {
		body, err := decode(q)
		if err != nil {
			return nil, err
		}
		m.Lock()
		defer m.Unlock()
		m.Patches = append(m.Patches, maps.Clone(body))
		maps.Copy(m.State, body)
		normalize(m.State)
		return httpmock.NewJsonResponse(200, m.State)
	})
	httpmock.RegisterResponder("DELETE", item, func(q *http.Request) (*http.Response, error) {
		m.Lock()
		defer m.Unlock()
		m.Deletes++
		m.State = nil
		return httpmock.NewStringResponse(204, ""), nil
	})
}

func (m *Mock) RegisterErrorMocks() {
	httpmock.RegisterResponder("POST", endpoint, httpmock.NewStringResponder(
		403,
		`{"error":{"code":"Forbidden","message":"Synthetic denied request"}}`,
	))
}

func (m *Mock) CleanupMockState() {
	m.Lock()
	defer m.Unlock()
	m.State = nil
	m.Posts = nil
	m.Patches = nil
	m.Deletes = 0
	m.PostFailures = 0
}
