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
	endpoint = "https://graph.microsoft.com/beta/networkAccess/explicitForwardProxyConfig"
	objectID = "11111111-2222-3333-4444-555555555555"
)

type Mock struct {
	sync.Mutex
	State   map[string]any
	Patches []map[string]any
	Deletes int
}

var _ common.MockRegistrar = (*Mock)(nil)

func init() { common.GlobalRegistry.Register("network_explicit_forward_proxy", &Mock{}) }
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
	item := endpoint
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
		access := m.State["internetAccess"].(map[string]any)
		for k, v := range body["internetAccess"].(map[string]any) {
			access[k] = v
		}
		if access["isSourceIpSessionAffinityEnabled"] == false {
			access["sourceIpSessionAffinityOptions"] = "none"
		}
		return httpmock.NewStringResponse(204, ""), nil
	})
}

func (m *Mock) RegisterErrorMocks() {
	for _, method := range []string{"GET", "POST", "PATCH", "DELETE"} {
		httpmock.RegisterResponder(
			method,
			endpoint,
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
	m.State = map[string]any{}
	_ = json.Unmarshal(
		[]byte(
			`{"proxyAutoConfigurationFileUrl": "https://example.invalid/generated", "privateAccess": {"isEnabled": false}, "internetAccess": {"isEnabled": false, "isSourceIpSessionAffinityEnabled": true, "sourceIpSessionAffinityOptions": "useSessionId", "isMtlsRequired": true, "futureSetting": "preserve"}}`,
		),
		&m.State,
	)
}
