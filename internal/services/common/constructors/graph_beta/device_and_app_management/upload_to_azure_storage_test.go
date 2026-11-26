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
