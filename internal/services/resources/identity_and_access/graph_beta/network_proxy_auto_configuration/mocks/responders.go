package mocks

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/jarcoal/httpmock"

	common "github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
)

const (
	endpoint = "https://graph.microsoft.com/beta/networkAccess/explicitForwardProxyConfig/proxyAutoConfigurationFiles"
	objectID = "11111111-2222-3333-4444-555555555555"
)

type Mock struct {
	sync.Mutex
	State   map[string]any
	Patches []map[string]any
	Deletes int
}

var _ common.MockRegistrar = (*Mock)(nil)

func init() { common.GlobalRegistry.Register("network_proxy_auto_configuration", &Mock{}) }
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

func (m *Mock) RegisterMocks() {
	m.CleanupMockState()
	item := endpoint + "/" + objectID
	httpmock.RegisterResponder("GET", item, func(q *http.Request) (*http.Response, error) {
		m.Lock()
		defer m.Unlock()
		if m.State == nil {
			return httpmock.NewStringResponse(
				404,
				`{"error":{"code":"NotFound","message":"Not found"}}`,
			), nil
		}
		return httpmock.NewJsonResponse(200, m.State)
	})
	httpmock.RegisterResponder("PATCH", item, func(q *http.Request) (*http.Response, error) {
		body, err := decode(q)
		if err != nil {
			return nil, err
		}
		m.Lock()
		defer m.Unlock()
		m.Patches = append(m.Patches, body)
		for k, v := range body {
			m.State[k] = v
		}
		m.State["lastModifiedDateTime"] = "2026-01-02T00:00:00Z"
		return httpmock.NewStringResponse(204, ""), nil
	})
	httpmock.RegisterResponder("POST", endpoint, func(q *http.Request) (*http.Response, error) {
		body, err := decode(q)
		if err != nil {
			return nil, err
		}
		m.Lock()
		defer m.Unlock()
		body["id"] = objectID
		body["createdDateTime"] = "2026-01-01T00:00:00Z"
		body["lastModifiedDateTime"] = "2026-01-01T00:00:00Z"
		m.State = body
		return httpmock.NewJsonResponse(201, body)
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
	for _, method := range []string{"GET", "POST", "PATCH", "DELETE"} {
		httpmock.RegisterResponder(
			method,
			endpoint+"/"+objectID,
			httpmock.NewStringResponder(
				403,
				`{"error":{"code":"Forbidden","message":"Synthetic denied request"}}`,
			),
		)
	}
}

func (m *Mock) CleanupMockState() {
	m.Lock()
	defer m.Unlock()
	m.Patches = nil
	m.Deletes = 0
	m.State = nil
}
