package mocks

import (
	"bytes"
	"compress/gzip"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
)

func TestUnitRequestCapture_RecordsNormalisesAndComparesFixtures(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(204, ""))
	capture := CaptureGraphRequests(t)

	send := func(method, url, body string, headers map[string]string) {
		t.Helper()
		var reader *bytes.Reader
		if body != "" {
			reader = bytes.NewReader([]byte(body))
		}
		var req *http.Request
		var err error
		if reader != nil {
			req, err = http.NewRequest(method, url, reader)
		} else {
			req, err = http.NewRequest(method, url, nil)
		}
		require.NoError(t, err)
		for name, value := range headers {
			req.Header.Set(name, value)
		}
		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		require.NoError(t, resp.Body.Close())
	}

	id := "AAAAAAAA-0000-0000-0000-000000000001"
	other := "bbbbbbbb-0000-0000-0000-000000000002"
	send("PATCH", "https://graph.microsoft.com/beta/things/"+other, `{"z":1,"a":{"y":true,"b":null}}`, map[string]string{"Content-Type": "application/json", "If-Match": "*", "Accept": "application/json"})
	send("GET", "https://graph.microsoft.com/beta/things/"+id+"?$select=id", "", nil)
	send("POST", "https://login.microsoftonline.com/token", "grant_type=x", nil)
	send("POST", "https://graph.microsoft.com/beta/things", "not json", map[string]string{"Prefer": "return=minimal"})
	var gzipped bytes.Buffer
	writer := gzip.NewWriter(&gzipped)
	_, err := writer.Write([]byte(`{"id":"` + strings.ToLower(id) + `"}`))
	require.NoError(t, err)
	require.NoError(t, writer.Close())
	send("POST", "https://graph.microsoft.com/beta/things", gzipped.String(), map[string]string{"Content-Encoding": "gzip"})
	send("DELETE", "https://graph.microsoft.com/beta/things/"+other, " ", nil)

	requests, err := capture.Requests()
	require.NoError(t, err)
	require.Len(t, requests, 5)
	require.Equal(t, "DELETE", requests[0].Method)
	require.Equal(t, "/beta/things/{uuid-1}", requests[0].URL)
	require.Nil(t, requests[0].Body)
	require.Equal(t, "/beta/things/{uuid-2}?$select=id", requests[1].URL)
	require.Equal(t, "PATCH", requests[2].Method)
	require.Equal(t, map[string]string{"Content-Type": "application/json", "If-Match": "*"}, requests[2].Headers)
	require.Equal(t, map[string]any{"id": "{uuid-2}"}, requests[3].Body)
	require.Equal(t, "not json", requests[4].Body)
	require.Equal(t, map[string]string{"Prefer": "return=minimal"}, requests[4].Headers)

	fixture := filepath.Join(t.TempDir(), "requests", "lifecycle.json")
	t.Setenv(UpdateRequestFixturesEnv, "1")
	capture.AssertMatchesFixture(t, fixture)
	written, err := os.ReadFile(fixture)
	require.NoError(t, err)
	require.Contains(t, string(written), `"a": {`)
	t.Setenv(UpdateRequestFixturesEnv, "")
	capture.AssertMatchesFixture(t, fixture)
	require.NoError(t, os.WriteFile(fixture, bytes.ReplaceAll(written, []byte("\n"), []byte("\r\n")), 0o600))
	capture.AssertMatchesFixture(t, fixture)
}

func TestUnitRequestCapture_MismatchAndMissingFixtureFail(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(204, ""))
	capture := CaptureGraphRequests(t)
	req, err := http.NewRequest("GET", "https://graph.microsoft.com/beta/things", nil)
	require.NoError(t, err)
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	require.NoError(t, resp.Body.Close())

	dir := t.TempDir()
	mismatch := filepath.Join(dir, "mismatch.json")
	require.NoError(t, os.WriteFile(mismatch, []byte("[]\n"), 0o600))
	for _, path := range []string{mismatch, filepath.Join(dir, "missing.json")} {
		recorder := &failureRecorder{TB: t}
		func() {
			defer func() { _ = recover() }()
			capture.AssertMatchesFixture(recorder, path)
		}()
		require.True(t, recorder.failed, path)
	}
}

func TestUnitRequestCapture_CleanupRestoresTransport(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	before := http.DefaultTransport
	t.Run("capture", func(t *testing.T) {
		capture := CaptureGraphRequests(t)
		require.Same(t, capture, http.DefaultTransport)
	})
	require.Equal(t, before, http.DefaultTransport)
}

// failureRecorder turns Fatalf into a recorded failure so a test can assert that a comparison fails.
type failureRecorder struct {
	testing.TB
	failed bool
}

func (f *failureRecorder) Helper() {}

func (f *failureRecorder) Fatal(...any) {
	f.failed = true
	panic("fatal")
}

func (f *failureRecorder) Fatalf(string, ...any) {
	f.failed = true
	panic("fatal")
}
