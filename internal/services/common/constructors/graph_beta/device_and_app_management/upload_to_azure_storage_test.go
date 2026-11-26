package sharedConstructors

import (
	"context"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"testing/synctest"
	"time"

	"github.com/stretchr/testify/require"
)

type azureUploadTimeoutTransport func(*http.Request) (*http.Response, error)

func (f azureUploadTimeoutTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

// A single timed-out request must not consume the entire retry budget. The
// PowerShell uploader retries the same block, with the same ID and bytes.
func TestUnitAzureStorageBlockRetryBoundaries(t *testing.T) {
	for _, test := range []struct {
		name      string
		timeouts  int
		status    int
		deadline  time.Duration
		wantCalls int
		wantError bool
	}{
		{"first request times out", 1, 0, 0, 2, false},
		{"all requests time out", 10, 0, 0, 3, true},
		{"parent deadline", 10, 0, time.Second, 1, true},
		{"HTTP retry limit", 0, 500, 0, 3, true},
	} {
		t.Run(test.name, func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				ctx := context.Background()
				if test.deadline > 0 {
					var cancel context.CancelFunc
					ctx, cancel = context.WithTimeout(ctx, test.deadline)
					defer cancel()
				}
				path := filepath.Join(t.TempDir(), "package.bin")
				require.NoError(t, os.WriteFile(path, []byte("test encrypted package bytes"), 0600))
				attempts, commits := 0, 0
				previous := http.DefaultClient
				defer func() { http.DefaultClient = previous }()
				http.DefaultClient = &http.Client{Transport: azureUploadTimeoutTransport(func(req *http.Request) (*http.Response, error) {
					require.Equal(t, http.MethodPut, req.Method)
					switch req.URL.Query().Get("comp") {
					case "block":
						attempts++
						require.Equal(t, "MDAwMDAw", req.URL.Query().Get("blockid"))
						data, err := io.ReadAll(req.Body)
						require.NoError(t, err)
						require.Equal(t, "test encrypted package bytes", string(data), "every attempt must replay the whole block")
						if attempts <= test.timeouts {
							<-req.Context().Done()
							return nil, req.Context().Err()
						}
						if test.status != 0 {
							return &http.Response{StatusCode: test.status, Header: make(http.Header), Body: io.NopCloser(strings.NewReader("transient storage failure"))}, nil
						}
					case "blocklist":
						commits++
						data, err := io.ReadAll(req.Body)
						require.NoError(t, err)
						require.Contains(t, string(data), "<Latest>MDAwMDAw</Latest>")
					default:
						t.Fatalf("unexpected upload request")
					}
					return &http.Response{StatusCode: http.StatusCreated, Header: make(http.Header), Body: io.NopCloser(strings.NewReader(""))}, nil
				})}
				err := UploadToAzureStorage(ctx, "https://storage.example.test/package?sig=synthetic-signed-token", path)
				if test.wantError {
					require.Error(t, err)
					require.Equal(t, test.status == 0, errors.Is(err, context.DeadlineExceeded), err)
					require.NotContains(t, err.Error(), "synthetic-signed-token")
				} else {
					require.NoError(t, err)
					require.Equal(t, 1, commits)
				}
				require.Equal(t, test.wantCalls, attempts)
			})
		})
	}
}

// Azure documents up to 30 seconds before a freshly created stored access policy
// becomes effective. Do not spend all attempts before a valid SAS can propagate.
func TestUnitAzureStorageStoredAccessPolicyPropagation(t *testing.T) {
	for _, test := range []struct {
		name                       string
		readyAfter, deadline       time.Duration
		permissionError, wantError bool
		wantCalls                  int
	}{
		{name: "policy becomes active", readyAfter: 30 * time.Second, wantCalls: 2},
		{name: "parent deadline while policy settles", readyAfter: 30 * time.Second, deadline: 5 * time.Second, wantError: true, wantCalls: 1},
		{name: "policy never becomes active", readyAfter: time.Hour, wantError: true, wantCalls: 3},
		{name: "unrelated forbidden response", permissionError: true, wantError: true, wantCalls: 3},
	} {
		t.Run(test.name, func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				started := time.Now()
				ctx := context.Background()
				if test.deadline > 0 {
					var cancel context.CancelFunc
					ctx, cancel = context.WithTimeout(ctx, test.deadline)
					defer cancel()
				}
				source := filepath.Join(t.TempDir(), "package.bin")
				require.NoError(t, os.WriteFile(source, []byte("same encrypted block"), 0600))
				calls, commits := 0, 0
				previous := http.DefaultClient
				defer func() { http.DefaultClient = previous }()
				http.DefaultClient = &http.Client{Transport: azureUploadTimeoutTransport(func(req *http.Request) (*http.Response, error) {
					if req.URL.Query().Get("comp") == "blocklist" {
						commits++
						return &http.Response{StatusCode: http.StatusCreated, Header: make(http.Header), Body: io.NopCloser(strings.NewReader(""))}, nil
					}
					calls++
					data, err := io.ReadAll(req.Body)
					require.NoError(t, err)
					require.Equal(t, "same encrypted block", string(data))
					require.Equal(t, "MDAwMDAw", req.URL.Query().Get("blockid"))
					if test.permissionError {
						return &http.Response{StatusCode: http.StatusForbidden, Header: make(http.Header), Body: io.NopCloser(strings.NewReader("<Error><Code>AuthorizationPermissionMismatch</Code></Error>"))}, nil
					}
					if time.Since(started) < test.readyAfter {
						return &http.Response{StatusCode: http.StatusForbidden, Header: http.Header{"X-Ms-Error-Code": []string{"AuthenticationFailed"}}, Body: io.NopCloser(strings.NewReader("<Error><Code>AuthenticationFailed</Code><AuthenticationErrorDetail>SAS identifier cannot be found for specified signed identifier</AuthenticationErrorDetail></Error>"))}, nil
					}
					return &http.Response{StatusCode: http.StatusCreated, Header: make(http.Header), Body: io.NopCloser(strings.NewReader(""))}, nil
				})}
				err := UploadToAzureStorage(ctx, "https://storage.example.test/package?sig=synthetic-signed-token", source)
				if test.wantError {
					require.Error(t, err)
					require.Zero(t, commits)
					require.NotContains(t, err.Error(), "synthetic-signed-token")
					if test.deadline > 0 {
						require.ErrorIs(t, err, context.DeadlineExceeded)
					}
				} else {
					require.NoError(t, err, "a fresh SAS must tolerate the documented policy propagation interval")
					require.Equal(t, 1, commits)
				}
				require.Equal(t, test.wantCalls, calls)
				if test.permissionError {
					require.Less(t, time.Since(started), 30*time.Second, "other403s must not wait for policy propagation")
				}
				if test.readyAfter == time.Hour {
					require.Less(t, time.Since(started), 2*time.Minute, "permanent policy failures must remain bounded")
				}
			})
		})
	}
}
