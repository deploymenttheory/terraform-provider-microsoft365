package mocks

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"
	"testing"
)

// UpdateRequestFixturesEnv names the environment variable that makes AssertMatchesFixture write the fixture instead of comparing.
const UpdateRequestFixturesEnv = "UPDATE_REQUEST_FIXTURES"

// capturedHeaders are the request headers a request-capture fixture records when present.
var capturedHeaders = []string{"Content-Type", "Prefer", "If-Match"}

// uuidPattern matches the GUIDs that mocks generate, so fixtures do not depend on random IDs.
var uuidPattern = regexp.MustCompile(`(?i)[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}`)

// CapturedRequest is one Graph request recorded by a RequestCapture.
type CapturedRequest struct {
	Method  string            `json:"method"`
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers,omitempty"`
	Body    any               `json:"body,omitempty"`
}

// RequestCapture records the Microsoft Graph requests a unit test sends through httpmock, so request-capture regression
// tests can prove that a refactor leaves the requests on the wire unchanged.
type RequestCapture struct {
	mu       sync.Mutex
	next     http.RoundTripper
	host     string
	requests []CapturedRequest
}

// CaptureGraphRequests wraps http.DefaultTransport, which must already be httpmock's transport, and records every
// request sent to graph.microsoft.com. Call it after httpmock.Activate. The wrapper is removed when the test ends,
// unless httpmock has already restored the original transport.
func CaptureGraphRequests(t testing.TB) *RequestCapture {
	t.Helper()
	capture := &RequestCapture{next: http.DefaultTransport, host: "graph.microsoft.com"}
	http.DefaultTransport = capture
	t.Cleanup(func() {
		if http.DefaultTransport == capture {
			http.DefaultTransport = capture.next
		}
	})
	return capture
}

// RoundTrip records the request when it targets Microsoft Graph, then sends it to the wrapped transport.
func (c *RequestCapture) RoundTrip(req *http.Request) (*http.Response, error) {
	if req.URL.Host == c.host {
		captured, err := c.record(req)
		if err != nil {
			return nil, err
		}
		c.mu.Lock()
		c.requests = append(c.requests, captured)
		c.mu.Unlock()
	}
	return c.next.RoundTrip(req) //nolint:wrapcheck // The capture must not change the transport's errors.
}

// record reads the method, path and query, captured headers and JSON body of req, restoring the body for the next transport.
func (c *RequestCapture) record(req *http.Request) (CapturedRequest, error) {
	captured := CapturedRequest{Method: req.Method, URL: req.URL.RequestURI()}
	for _, name := range capturedHeaders {
		if value := req.Header.Get(name); value != "" {
			if captured.Headers == nil {
				captured.Headers = map[string]string{}
			}
			captured.Headers[name] = value
		}
	}
	if req.Body == nil || req.Body == http.NoBody {
		return captured, nil
	}
	raw, err := io.ReadAll(req.Body)
	if err != nil {
		return captured, fmt.Errorf("read captured request body: %w", err)
	}
	_ = req.Body.Close()
	req.Body = io.NopCloser(bytes.NewReader(raw))
	payload := raw
	if req.Header.Get("Content-Encoding") == "gzip" {
		reader, err := gzip.NewReader(bytes.NewReader(raw))
		if err != nil {
			return captured, fmt.Errorf("decompress captured request body: %w", err)
		}
		if payload, err = io.ReadAll(reader); err != nil {
			return captured, fmt.Errorf("decompress captured request body: %w", err)
		}
	}
	if len(bytes.TrimSpace(payload)) == 0 {
		return captured, nil
	}
	decoder := json.NewDecoder(bytes.NewReader(payload))
	decoder.UseNumber()
	var body any
	if err := decoder.Decode(&body); err != nil {
		captured.Body = string(payload)
		return captured, nil //nolint:nilerr // A non-JSON body is recorded verbatim.
	}
	captured.Body = body
	return captured, nil
}

// Requests returns the recorded requests in a stable form. GUIDs are replaced with {uuid-N} in the order they were
// first sent, and the requests are sorted, because Terraform refreshes independent resources in parallel and so the
// arrival order of their reads is not fixed. Repeated identical requests are kept.
func (c *RequestCapture) Requests() ([]CapturedRequest, error) {
	c.mu.Lock()
	recorded := append([]CapturedRequest(nil), c.requests...)
	c.mu.Unlock()
	ids := map[string]string{}
	normalise := func(value string) string {
		return uuidPattern.ReplaceAllStringFunc(value, func(id string) string {
			key := strings.ToLower(id)
			if _, ok := ids[key]; !ok {
				ids[key] = fmt.Sprintf("{uuid-%d}", len(ids)+1)
			}
			return ids[key]
		})
	}
	type keyed struct {
		key     string
		request CapturedRequest
	}
	normalised := make([]keyed, 0, len(recorded))
	for _, request := range recorded {
		raw, err := json.Marshal(request)
		if err != nil {
			return nil, fmt.Errorf("encode captured request: %w", err)
		}
		decoder := json.NewDecoder(strings.NewReader(normalise(string(raw))))
		decoder.UseNumber()
		var out CapturedRequest
		if err := decoder.Decode(&out); err != nil {
			return nil, fmt.Errorf("decode captured request: %w", err)
		}
		key, err := json.Marshal(out)
		if err != nil {
			return nil, fmt.Errorf("encode captured request: %w", err)
		}
		normalised = append(normalised, keyed{key: string(key), request: out})
	}
	sort.SliceStable(normalised, func(i, j int) bool { return normalised[i].key < normalised[j].key })
	result := make([]CapturedRequest, 0, len(normalised))
	for _, entry := range normalised {
		result = append(result, entry.request)
	}
	return result, nil
}

// AssertMatchesFixture compares the recorded requests with the JSON fixture at path. When UPDATE_REQUEST_FIXTURES=1
// it writes the fixture instead, creating its directory.
func (c *RequestCapture) AssertMatchesFixture(t testing.TB, path string) {
	t.Helper()
	requests, err := c.Requests()
	if err != nil {
		t.Fatal(err)
	}
	actual, err := json.MarshalIndent(requests, "", "  ")
	if err != nil {
		t.Fatalf("encode captured requests: %v", err)
	}
	actual = append(actual, '\n')
	if os.Getenv(UpdateRequestFixturesEnv) == "1" {
		if err := os.MkdirAll(filepath.Dir(path), 0o750); err != nil {
			t.Fatalf("create fixture directory: %v", err)
		}
		if err := os.WriteFile(path, actual, 0o600); err != nil {
			t.Fatalf("write request fixture: %v", err)
		}
		return
	}
	expected, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read request fixture %s (set %s=1 to create it): %v", path, UpdateRequestFixturesEnv, err)
	}
	expected = bytes.ReplaceAll(expected, []byte("\r\n"), []byte("\n"))
	if !bytes.Equal(expected, actual) {
		t.Fatalf("captured requests differ from %s\n--- expected\n%s\n--- actual\n%s", path, expected, actual)
	}
}
