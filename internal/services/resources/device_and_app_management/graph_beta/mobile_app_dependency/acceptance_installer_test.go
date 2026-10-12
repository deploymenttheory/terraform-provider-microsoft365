package graphBetaMobileAppDependency_test

import (
	"archive/zip"
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// Prepare the same real Firefox installer used by the Win32 app acceptance tests.
// The sibling resource encrypts, uploads, and commits this ZIP through Graph.
func prepareAcceptanceInstaller(t *testing.T) string {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		"https://archive.mozilla.org/pub/firefox/releases/140.0.2/win64/en-US/Firefox%20Setup%20140.0.2.exe",
		nil,
	)
	require.NoError(t, err)
	response, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer response.Body.Close()
	require.Equal(t, http.StatusOK, response.StatusCode)

	path := filepath.Join(t.TempDir(), "firefox-140.0.2.zip")
	file, err := os.Create(path)
	require.NoError(t, err)
	defer file.Close()
	archive := zip.NewWriter(file)
	defer archive.Close()
	entry, err := archive.CreateHeader(&zip.FileHeader{Name: "setup.exe", Method: zip.Store})
	require.NoError(t, err)
	digest := sha256.New()
	const maxInstallerSize = 250 << 20
	size, err := io.Copy(
		io.MultiWriter(entry, digest),
		io.LimitReader(response.Body, maxInstallerSize+1),
	)
	require.NoError(t, err)
	require.LessOrEqual(t, size, int64(maxInstallerSize))
	require.Equal(
		t,
		"f7e9d9f30c0e96c502a53b25797f69b7d64fa9dbb5416d94d2a5ac6c51689142",
		fmt.Sprintf("%x", digest.Sum(nil)),
		"Firefox SHA-256 must match Mozilla's SHA256SUMS",
	)
	require.NoError(t, archive.Close())
	require.NoError(t, file.Close())
	return path
}
