package customrequests

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	abstractions "github.com/microsoft/kiota-abstractions-go"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// rawTestPolicy is a plain Go struct standing in for a Graph type the SDK lacks.
type rawTestPolicy struct {
	ID          string `json:"id,omitempty"`
	DisplayName string `json:"displayName"`
	State       string `json:"state,omitempty"`
}

// capturedRequest is one request that reached the mock transport.
type capturedRequest struct {
	method, url, contentType, prefer, accept, body string
}

// captureRequests records every request that reaches transport and answers each with respond.
func captureRequests(transport *httpmock.MockTransport, respond func(req *http.Request) (*http.Response, error)) *[]capturedRequest {
	captured := &[]capturedRequest{}
	transport.RegisterNoResponder(func(req *http.Request) (*http.Response, error) {
		var body []byte
		if req.Body != nil {
			body, _ = io.ReadAll(req.Body)
		}
		*captured = append(*captured, capturedRequest{method: req.Method, url: req.URL.String(), contentType: req.Header.Get("Content-Type"), prefer: req.Header.Get("Prefer"), accept: req.Header.Get("Accept"), body: string(body)})
		return respond(req)
	})
	return captured
}

func jsonResponse(status int, body string) func(req *http.Request) (*http.Response, error) {
	return func(req *http.Request) (*http.Response, error) {
		resp := httpmock.NewStringResponse(status, body)
		resp.Header.Set("Content-Type", "application/json")
		return resp, nil
	}
}

const (
	preferUnknownEnums = "include-unknown-enum-members"
	rawPolicyResponse  = `{"id":"p1","displayName":"Raw policy","state":"futureValue"}`
)

var rawPolicyHeaders = map[string]string{"Prefer": preferUnknownEnums}

// rawVerbCall sends one verb in raw JSON mode with the Prefer header, decoding any response into target.
type rawVerbCall struct {
	name, wantMethod, wantURL, wantBody string
	send                                func(ctx context.Context, adapter abstractions.RequestAdapter, target *rawTestPolicy) error
}

// rawVerbCalls returns every verb in raw JSON mode against the US Gov beta root.
func rawVerbCalls() []rawVerbCall {
	body := rawTestPolicy{DisplayName: "Raw policy"}
	wantBody := `{"displayName":"Raw policy"}`

	return []rawVerbCall{
		{name: "GetRequest", wantMethod: http.MethodGet, wantURL: usGovBetaRoot + "/deviceManagement/rawPolicies/" + testResourceID, send: func(ctx context.Context, adapter abstractions.RequestAdapter, target *rawTestPolicy) error {
			_, err := GetRequest(ctx, adapter, GetRequestConfig{APIVersion: GraphAPIBeta, Endpoint: "deviceManagement/rawPolicies", ResourceID: testResourceID, Headers: rawPolicyHeaders, RawResponse: target}, nil, nil)
			return err
		}},
		{name: "GetRequestByResourceId", wantMethod: http.MethodGet, wantURL: usGovBetaRoot + "/deviceManagement/rawPolicies('" + testResourceID + "')", send: func(ctx context.Context, adapter abstractions.RequestAdapter, target *rawTestPolicy) error {
			_, err := GetRequestByResourceId(ctx, adapter, GetRequestConfig{APIVersion: GraphAPIBeta, Endpoint: "deviceManagement/rawPolicies", ResourceIDPattern: ResourceIDPatternParentheses, ResourceID: testResourceID, Headers: rawPolicyHeaders, RawResponse: target})
			return err
		}},
		{name: "PostRequest", wantMethod: http.MethodPost, wantURL: usGovBetaRoot + "/deviceManagement/rawPolicies", wantBody: wantBody, send: func(ctx context.Context, adapter abstractions.RequestAdapter, target *rawTestPolicy) error {
			_, err := PostRequest(ctx, adapter, PostRequestConfig{APIVersion: GraphAPIBeta, Endpoint: "deviceManagement/rawPolicies", RawBody: body, Headers: rawPolicyHeaders, RawResponse: target}, nil, nil)
			return err
		}},
		{name: "PostRequestNoContent", wantMethod: http.MethodPost, wantURL: usGovBetaRoot + "/deviceManagement/rawPolicies/" + testResourceID + "/assign", wantBody: wantBody, send: func(ctx context.Context, adapter abstractions.RequestAdapter, target *rawTestPolicy) error {
			return PostRequestNoContent(ctx, adapter, PostRequestConfig{APIVersion: GraphAPIBeta, Endpoint: "deviceManagement/rawPolicies/" + testResourceID + "/assign", RawBody: body, Headers: rawPolicyHeaders, RawResponse: target})
		}},
		{name: "PatchRequestByResourceId", wantMethod: http.MethodPatch, wantURL: usGovBetaRoot + "/deviceManagement/rawPolicies/" + testResourceID, wantBody: wantBody, send: func(ctx context.Context, adapter abstractions.RequestAdapter, target *rawTestPolicy) error {
			return PatchRequestByResourceId(ctx, adapter, PatchRequestConfig{APIVersion: GraphAPIBeta, Endpoint: "deviceManagement/rawPolicies", ResourceIDPattern: ResourceIDPatternSegment, ResourceID: testResourceID, RawBody: body, Headers: rawPolicyHeaders, RawResponse: target})
		}},
		{name: "PutRequestByResourceId", wantMethod: http.MethodPut, wantURL: usGovBetaRoot + "/deviceManagement/rawPolicies('" + testResourceID + "')", wantBody: wantBody, send: func(ctx context.Context, adapter abstractions.RequestAdapter, target *rawTestPolicy) error {
			return PutRequestByResourceId(ctx, adapter, PutRequestConfig{APIVersion: GraphAPIBeta, Endpoint: "deviceManagement/rawPolicies", ResourceID: testResourceID, RawBody: body, Headers: rawPolicyHeaders, RawResponse: target})
		}},
		{name: "DeleteRequestByResourceId", wantMethod: http.MethodDelete, wantURL: usGovBetaRoot + "/deviceManagement/rawPolicies('" + testResourceID + "')", send: func(ctx context.Context, adapter abstractions.RequestAdapter, target *rawTestPolicy) error {
			return DeleteRequestByResourceId(ctx, adapter, DeleteRequestConfig{APIVersion: GraphAPIBeta, Endpoint: "deviceManagement/rawPolicies", ResourceIDPattern: ResourceIDPatternParentheses, ResourceID: testResourceID, Headers: rawPolicyHeaders, RawResponse: target})
		}},
	}
}

// TestRawModeEveryVerb proves every verb sends a raw JSON body and custom headers to a non-default (US Gov) base URL
// through the adapter, and decodes the raw JSON response into a plain Go struct.
func TestRawModeEveryVerb(t *testing.T) {
	for _, call := range rawVerbCalls() {
		t.Run(call.name, func(t *testing.T) {
			adapter, transport := newTestAdapter(t, usGovV1Root)
			captured := captureRequests(transport, jsonResponse(http.StatusOK, rawPolicyResponse))

			var target rawTestPolicy
			require.NoError(t, call.send(context.Background(), adapter, &target))

			require.Len(t, *captured, 1)
			sent := (*captured)[0]
			assert.Equal(t, call.wantMethod, sent.method)
			assert.Equal(t, call.wantURL, sent.url)
			assert.Equal(t, preferUnknownEnums, sent.prefer)
			assert.Equal(t, "application/json", sent.accept)
			if call.wantBody != "" {
				assert.Equal(t, "application/json", sent.contentType)
				assert.JSONEq(t, call.wantBody, sent.body)
			} else {
				assert.Empty(t, sent.body)
			}
			assert.Equal(t, rawTestPolicy{ID: "p1", DisplayName: "Raw policy", State: "futureValue"}, target)
		})
	}
}

// TestRawModeNoContent proves a 204 No Content with a raw response target returns no error and leaves the target empty.
func TestRawModeNoContent(t *testing.T) {
	for _, call := range rawVerbCalls() {
		t.Run(call.name, func(t *testing.T) {
			adapter, transport := newTestAdapter(t, usGovBetaRoot)
			captured := captureRequests(transport, func(req *http.Request) (*http.Response, error) {
				return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
			})

			var target rawTestPolicy
			require.NoError(t, call.send(context.Background(), adapter, &target))
			assert.Equal(t, rawTestPolicy{}, target)
			assert.Len(t, *captured, 1)
		})
	}
}

// TestRawModeReturnsODataErrors proves a 4xx response in raw mode returns an error wrapping an ODataError for every verb.
func TestRawModeReturnsODataErrors(t *testing.T) {
	for _, call := range rawVerbCalls() {
		for _, st := range []struct {
			status int
			code   string
		}{{status: http.StatusBadRequest, code: "BadRequest"}, {status: http.StatusNotFound, code: "ResourceNotFound"}, {status: http.StatusConflict, code: "Conflict"}} {
			t.Run(call.name+"/"+st.code, func(t *testing.T) {
				adapter, transport := newTestAdapter(t, usGovBetaRoot)
				transport.RegisterNoResponder(odataErrorResponder(st.status, st.code))

				var target rawTestPolicy
				requireODataError(t, call.send(context.Background(), adapter, &target), st.status, st.code)
				assert.Equal(t, rawTestPolicy{}, target)
				assert.Equal(t, 1, transport.GetTotalCallCount())
			})
		}
	}
}

// TestRawGetCollectsEveryPage proves a raw GET follows @odata.nextLink across 3 pages through the adapter, sends the
// custom headers with every page and decodes every page's items into the target.
func TestRawGetCollectsEveryPage(t *testing.T) {
	ctx := context.Background()
	firstURL := usGovBetaRoot + "/deviceManagement/rawPolicies"

	tests := []struct {
		name string
		send func(adapter abstractions.RequestAdapter, target any) error
	}{
		{name: "GetRequest", send: func(adapter abstractions.RequestAdapter, target any) error {
			_, err := GetRequest(ctx, adapter, GetRequestConfig{APIVersion: GraphAPIBeta, Endpoint: "deviceManagement/rawPolicies", QueryParameters: map[string]string{"$top": "2"}, Headers: rawPolicyHeaders, RawResponse: target}, nil, nil)
			return err
		}},
		{name: "GetRequestByResourceId", send: func(adapter abstractions.RequestAdapter, target any) error {
			_, err := GetRequestByResourceId(ctx, adapter, GetRequestConfig{APIVersion: GraphAPIBeta, Endpoint: "deviceManagement/rawPolicies", QueryParameters: map[string]string{"$top": "2"}, Headers: rawPolicyHeaders, RawResponse: target})
			return err
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adapter, transport := newTestAdapter(t, usGovV1Root)
			captured := captureRequests(transport, func(req *http.Request) (*http.Response, error) {
				switch req.URL.Query().Get("$skiptoken") {
				case "":
					return jsonResponse(http.StatusOK, `{"value":[{"displayName":"a"},{"displayName":"b"}],"@odata.nextLink":"`+firstURL+`?$skiptoken=2"}`)(req)
				case "2":
					return jsonResponse(http.StatusOK, `{"value":[{"displayName":"c"},{"displayName":"d"}],"@odata.nextLink":"`+firstURL+`?$skiptoken=3"}`)(req)
				default:
					return jsonResponse(http.StatusOK, `{"value":[{"displayName":"e"}]}`)(req)
				}
			})

			var page struct {
				Value []rawTestPolicy `json:"value"`
			}
			require.NoError(t, tt.send(adapter, &page))

			require.Len(t, *captured, 3)
			assert.Equal(t, firstURL+"?%24top=2", (*captured)[0].url)
			assert.Equal(t, firstURL+"?$skiptoken=2", (*captured)[1].url)
			assert.Equal(t, firstURL+"?$skiptoken=3", (*captured)[2].url)
			for _, sent := range *captured {
				assert.Equal(t, preferUnknownEnums, sent.prefer)
			}

			names := make([]string, 0, len(page.Value))
			for _, item := range page.Value {
				names = append(names, item.DisplayName)
			}
			assert.Equal(t, []string{"a", "b", "c", "d", "e"}, names)
		})
	}
}

// TestRawGetPagingError proves an error on a later page is returned as an ODataError.
func TestRawGetPagingError(t *testing.T) {
	adapter, transport := newTestAdapter(t, usGovBetaRoot)
	firstURL := usGovBetaRoot + "/deviceManagement/rawPolicies"
	transport.RegisterResponder(http.MethodGet, firstURL, jsonResponder(http.StatusOK, `{"value":[{"displayName":"a"}],"@odata.nextLink":"`+firstURL+`?$skiptoken=2"}`))
	transport.RegisterResponderWithQuery(http.MethodGet, firstURL, "$skiptoken=2", odataErrorResponder(http.StatusForbidden, "Forbidden"))

	var page struct {
		Value []rawTestPolicy `json:"value"`
	}
	_, err := GetRequest(context.Background(), adapter, GetRequestConfig{APIVersion: GraphAPIBeta, Endpoint: "deviceManagement/rawPolicies", RawResponse: &page}, nil, nil)
	requireODataError(t, err, http.StatusForbidden, "Forbidden")
}

// TestMixedModesAreRejected proves a config with both a Parsable and a raw body, or both a ParsableFactory and a raw
// response target, is rejected before anything is sent.
func TestMixedModesAreRejected(t *testing.T) {
	ctx := context.Background()
	raw := rawTestPolicy{DisplayName: "Raw policy"}
	var target rawTestPolicy

	tests := []struct {
		name string
		want error
		send func(adapter abstractions.RequestAdapter) error
	}{
		{name: "PostRequest body", want: ErrMixedBodyModes, send: func(adapter abstractions.RequestAdapter) error {
			_, err := PostRequest(ctx, adapter, PostRequestConfig{APIVersion: GraphAPIBeta, Endpoint: "deviceManagement/deviceCategories", RequestBody: newTestBody(), RawBody: raw}, models.CreateDeviceCategoryFromDiscriminatorValue, nil)
			return err
		}},
		{name: "PostRequestNoContent body", want: ErrMixedBodyModes, send: func(adapter abstractions.RequestAdapter) error {
			return PostRequestNoContent(ctx, adapter, PostRequestConfig{APIVersion: GraphAPIBeta, Endpoint: "deviceManagement/deviceCategories", RequestBody: newTestBody(), RawBody: raw})
		}},
		{name: "PatchRequestByResourceId body", want: ErrMixedBodyModes, send: func(adapter abstractions.RequestAdapter) error {
			return PatchRequestByResourceId(ctx, adapter, PatchRequestConfig{APIVersion: GraphAPIBeta, Endpoint: "deviceManagement/deviceCategories", ResourceIDPattern: ResourceIDPatternSegment, ResourceID: testResourceID, RequestBody: newTestBody(), RawBody: raw})
		}},
		{name: "PutRequestByResourceId body", want: ErrMixedBodyModes, send: func(adapter abstractions.RequestAdapter) error {
			return PutRequestByResourceId(ctx, adapter, PutRequestConfig{APIVersion: GraphAPIBeta, Endpoint: "deviceManagement/deviceCategories", ResourceID: testResourceID, RequestBody: newTestBody(), RawBody: raw})
		}},
		{name: "PostRequest response", want: ErrMixedResponseModes, send: func(adapter abstractions.RequestAdapter) error {
			_, err := PostRequest(ctx, adapter, PostRequestConfig{APIVersion: GraphAPIBeta, Endpoint: "deviceManagement/deviceCategories", RawBody: raw, RawResponse: &target}, models.CreateDeviceCategoryFromDiscriminatorValue, nil)
			return err
		}},
		{name: "GetRequest response", want: ErrMixedResponseModes, send: func(adapter abstractions.RequestAdapter) error {
			_, err := GetRequest(ctx, adapter, GetRequestConfig{APIVersion: GraphAPIBeta, Endpoint: "deviceManagement/deviceCategories", RawResponse: &target}, models.CreateDeviceCategoryFromDiscriminatorValue, nil)
			return err
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adapter, transport := newTestAdapter(t, usGovBetaRoot)
			transport.RegisterNoResponder(jsonResponder(http.StatusOK, `{}`))

			require.ErrorIs(t, tt.send(adapter), tt.want)
			assert.Equal(t, 0, transport.GetTotalCallCount())
		})
	}
}

// TestRawBodyEncodings proves pre-encoded JSON is sent as it is, a typed nil Parsable does not count as a body, and
// *json.RawMessage and *[]byte targets receive the response body as it is.
func TestRawBodyEncodings(t *testing.T) {
	ctx := context.Background()
	var nilParsable *models.DeviceCategory

	tests := []struct {
		name     string
		config   PostRequestConfig
		wantBody string
	}{
		{name: "json.RawMessage", config: PostRequestConfig{RawBody: json.RawMessage(`{"a":1}`)}, wantBody: `{"a":1}`},
		{name: "bytes", config: PostRequestConfig{RawBody: []byte(`{"b":[1,2]}`)}, wantBody: `{"b":[1,2]}`},
		{name: "map with typed nil Parsable", config: PostRequestConfig{RequestBody: nilParsable, RawBody: map[string]any{"c": "x"}}, wantBody: `{"c":"x"}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adapter, transport := newTestAdapter(t, usGovBetaRoot)
			captured := captureRequests(transport, jsonResponse(http.StatusOK, `{"id":"1"}`))

			var rawMessage json.RawMessage
			tt.config.APIVersion = GraphAPIBeta
			tt.config.Endpoint = "deviceManagement/rawPolicies"
			tt.config.RawResponse = &rawMessage
			_, err := PostRequest(ctx, adapter, tt.config, nil, nil)
			require.NoError(t, err)

			require.Len(t, *captured, 1)
			assert.Equal(t, tt.wantBody, (*captured)[0].body)
			assert.Equal(t, `{"id":"1"}`, string(rawMessage))
		})
	}

	t.Run("bytes target", func(t *testing.T) {
		adapter, transport := newTestAdapter(t, usGovBetaRoot)
		captureRequests(transport, jsonResponse(http.StatusOK, `{"id":"2"}`))

		var raw []byte
		require.NoError(t, PatchRequestByResourceId(ctx, adapter, PatchRequestConfig{APIVersion: GraphAPIBeta, Endpoint: "deviceManagement/rawPolicies", ResourceIDPattern: ResourceIDPatternSegment, ResourceID: testResourceID, RawBody: rawTestPolicy{DisplayName: "x"}, RawResponse: &raw}))
		assert.Equal(t, `{"id":"2"}`, string(raw))
	})

	t.Run("undecodable response", func(t *testing.T) {
		adapter, transport := newTestAdapter(t, usGovBetaRoot)
		captureRequests(transport, jsonResponse(http.StatusOK, `not json`))

		var target rawTestPolicy
		_, err := PostRequest(ctx, adapter, PostRequestConfig{APIVersion: GraphAPIBeta, Endpoint: "deviceManagement/rawPolicies", RawBody: rawTestPolicy{DisplayName: "x"}, RawResponse: &target}, nil, nil)
		require.ErrorContains(t, err, "decode Graph JSON response")
	})

	t.Run("unmarshalable body", func(t *testing.T) {
		adapter, transport := newTestAdapter(t, usGovBetaRoot)
		_, err := PostRequest(ctx, adapter, PostRequestConfig{APIVersion: GraphAPIBeta, Endpoint: "deviceManagement/rawPolicies", RawBody: make(chan int)}, nil, nil)
		require.ErrorContains(t, err, "marshal Graph JSON request")
		assert.Equal(t, 0, transport.GetTotalCallCount())
	})
}

// TestParsableModeSendsHeaders proves the Parsable mode keeps its body and response handling and also sends custom headers.
func TestParsableModeSendsHeaders(t *testing.T) {
	adapter, transport := newTestAdapter(t, usGovBetaRoot)
	captured := captureRequests(transport, jsonResponse(http.StatusOK, `{"displayName":"created"}`))

	result, err := PostRequest(context.Background(), adapter, PostRequestConfig{APIVersion: GraphAPIBeta, Endpoint: "deviceManagement/deviceCategories", RequestBody: newTestBody(), Headers: rawPolicyHeaders}, models.CreateDeviceCategoryFromDiscriminatorValue, nil)
	require.NoError(t, err)

	category, ok := result.(models.DeviceCategoryable)
	require.True(t, ok)
	assert.Equal(t, "created", *category.GetDisplayName())
	require.Len(t, *captured, 1)
	assert.Equal(t, preferUnknownEnums, (*captured)[0].prefer)
	assert.JSONEq(t, `{"displayName":"test"}`, (*captured)[0].body)
}
