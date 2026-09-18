package mocks

import (
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/jarcoal/httpmock"

	common "github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
)

const endpoint = "https://graph.microsoft.com/beta/networkAccess/settings/forwardingOptions"

type ForwardingOptionsMock struct {
	sync.Mutex
	State   map[string]any
	Patches []map[string]any
}

var _ common.MockRegistrar = (*ForwardingOptionsMock)(nil)

var errUnexpectedPatchFields = errors.New("unexpected patch fields")

func init() { common.GlobalRegistry.Register("network_forwarding_options", &ForwardingOptionsMock{}) }
func (m *ForwardingOptionsMock) RegisterMocks() {
	m.CleanupMockState()
	httpmock.RegisterResponder("GET", endpoint, func(q *http.Request) (*http.Response, error) {
		m.Lock()
		defer m.Unlock()
		return httpmock.NewJsonResponse(200, m.State)
	})
	httpmock.RegisterResponder("PATCH", endpoint, func(q *http.Request) (*http.Response, error) {
		var reader io.Reader = q.Body
		if q.Header.Get("Content-Encoding") == "gzip" {
			gz, err := gzip.NewReader(q.Body)
			if err != nil {
				return nil, fmt.Errorf("decompress forwarding options request: %w", err)
			}
			defer gz.Close()
			reader = gz
		}
		var body map[string]any
		if err := json.NewDecoder(reader).Decode(&body); err != nil {
			return nil, fmt.Errorf("decode forwarding options request: %w", err)
		}
		if len(body) != 1 ||
			(body["skipDnsLookupState"] != "enabled" && body["skipDnsLookupState"] != "disabled") {
			return nil, fmt.Errorf("%w: %v", errUnexpectedPatchFields, body)
		}
		m.Lock()
		defer m.Unlock()
		m.Patches = append(m.Patches, body)
		m.State["skipDnsLookupState"] = body["skipDnsLookupState"]
		return httpmock.NewStringResponse(204, ""), nil
	})
}

func (m *ForwardingOptionsMock) RegisterErrorMocks() {
	for _, method := range []string{"GET", "PATCH"} {
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

func (m *ForwardingOptionsMock) CleanupMockState() {
	m.Lock()
	defer m.Unlock()
	m.State = map[string]any{
		"skipDnsLookupState":     "enabled",
		"unmanagedFutureSetting": "preserve-me",
	}
	m.Patches = nil
}
