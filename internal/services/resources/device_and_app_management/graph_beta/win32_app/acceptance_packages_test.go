package graphBetaWin32App_test

import (
	"archive/zip"
	"context"
	"crypto/sha256"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	construct "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors/graph_beta/device_and_app_management"
	"github.com/stretchr/testify/require"
)

// Mozilla publishes EXE installers, not ZIP or .intunewin packages. Download
// pinned public releases over HTTPS and prepare both formats on the CI runner.
// A loopback server exercises installer_url_source without publishing fixtures
// or requiring Windows, additional environment variables, or a setup script.
var firefoxInstallers = map[string]struct {
	url    string
	sha256 string
}{
	"140.0.2": {
		url:    "https://archive.mozilla.org/pub/firefox/releases/140.0.2/win64/en-US/Firefox%20Setup%20140.0.2.exe",
		sha256: "f7e9d9f30c0e96c502a53b25797f69b7d64fa9dbb5416d94d2a5ac6c51689142",
	},
	"140.0.4": {
		url:    "https://archive.mozilla.org/pub/firefox/releases/140.0.4/win64/en-US/Firefox%20Setup%20140.0.4.exe",
		sha256: "b4a7ec96b1ae9539c1cb7a49bc103f267bdda8051c9ec0574cd2acc2da0e45a7",
	},
}

type acceptancePackages struct {
	URL       string
	directory string
	prepared  map[string]bool
}

func newAcceptancePackages(t *testing.T) *acceptancePackages {
	t.Helper()
	directory := t.TempDir()
	server := httptest.NewServer(http.FileServer(http.Dir(directory)))
	t.Cleanup(server.Close)
	return &acceptancePackages{URL: server.URL, directory: directory, prepared: make(map[string]bool)}
}

func (p *acceptancePackages) prepare(t *testing.T, versions ...string) {
	t.Helper()
	client := &http.Client{Timeout: 5 * time.Minute}
	for _, version := range versions {
		if p.prepared[version] {
			continue
		}
		installer, ok := firefoxInstallers[version]
		require.True(t, ok, "no pinned Firefox installer for %s", version)
		t.Logf("Preparing Firefox %s from %s", version, installer.url)
		func() {
			response, err := client.Get(installer.url)
			require.NoError(t, err)
			defer response.Body.Close()
			require.Equal(t, http.StatusOK, response.StatusCode, "download Firefox %s", version)
			zipPath := filepath.Join(p.directory, "firefox-"+version+".zip")
			file, err := os.Create(zipPath)
			require.NoError(t, err)
			defer file.Close()
			archive := zip.NewWriter(file)
			defer archive.Close()
			entry, err := archive.CreateHeader(&zip.FileHeader{Name: "setup.exe", Method: zip.Store})
			require.NoError(t, err)
			digest := sha256.New()
			const maxInstallerSize = 250 << 20
			size, err := io.Copy(io.MultiWriter(entry, digest), io.LimitReader(response.Body, maxInstallerSize+1))
			require.NoError(t, err)
			require.LessOrEqual(t, size, int64(maxInstallerSize), "Firefox installer is larger than expected")
			require.Equal(t, installer.sha256, fmt.Sprintf("%x", digest.Sum(nil)), "Firefox SHA-256 does not match Mozilla's SHA256SUMS")
			require.NoError(t, archive.Close())
			require.NoError(t, file.Close())
			packageAcceptanceZip(t, zipPath, filepath.Join(p.directory, "firefox-"+version+".intunewin"))
			p.prepared[version] = true
		}()
	}
}

type acceptanceDetectionMetadata struct {
	XMLName                xml.Name `xml:"ApplicationInfo"`
	Name                   string
	FileName               string
	SetupFile              string
	UnencryptedContentSize int64
	EncryptionInfo         *construct.EncryptionInfo
}

func packageAcceptanceZip(t *testing.T, zipPath, packagePath string) {
	t.Helper()
	content, encryption, err := construct.EncryptMobileAppAndConstructFileContentMetadata(context.Background(), zipPath)
	require.NoError(t, err)
	encryptedPath := zipPath + ".bin"
	defer os.Remove(encryptedPath)
	payloadName := filepath.Base(packagePath)
	metadata, err := xml.Marshal(acceptanceDetectionMetadata{
		Name: "Mozilla Firefox", FileName: payloadName, SetupFile: "setup.exe",
		UnencryptedContentSize: *content.GetSize(), EncryptionInfo: encryption,
	})
	require.NoError(t, err)
	file, err := os.Create(packagePath)
	require.NoError(t, err)
	defer file.Close()
	archive := zip.NewWriter(file)
	defer archive.Close()
	entry, err := archive.Create("IntuneWinPackage/Metadata/Detection.xml")
	require.NoError(t, err)
	_, err = entry.Write(metadata)
	require.NoError(t, err)
	entry, err = archive.CreateHeader(&zip.FileHeader{Name: "IntuneWinPackage/Contents/" + payloadName, Method: zip.Store})
	require.NoError(t, err)
	encrypted, err := os.Open(encryptedPath)
	require.NoError(t, err)
	defer encrypted.Close()
	_, err = io.Copy(entry, encrypted)
	require.NoError(t, err)
	require.NoError(t, archive.Close())
	require.NoError(t, file.Close())
}
