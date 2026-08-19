package graphBetaWin32App

import (
	"archive/zip"
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema/graph_beta/device_and_app_management"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/microsoft/kiota-abstractions-go/authentication"
	kiotahttp "github.com/microsoft/kiota-http-go"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnitResourceWin32App_01_ExtractIntuneWinPackage(t *testing.T) {
	encryptedContent := []byte("the original encrypted inner payload")
	metadata := validIntuneWinMetadata()
	packagePath := writeIntuneWinFixture(t, metadata, map[string][]byte{
		"IntuneWinPackage/Contents/" + metadata.FileName: encryptedContent,
	})

	installer, err := extractIntuneWinPackage(packagePath)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, installer.cleanup()) })

	assert.Equal(t, metadata.FileName, installer.fileName)
	assert.Equal(t, metadata.UnencryptedContentSize, installer.unencryptedSize)
	assert.Equal(t, int64(len(encryptedContent)), installer.encryptedSize)
	assert.Equal(t, metadata.EncryptionInfo.EncryptionKey, installer.encryptionInfo.EncryptionKey)
	assert.Equal(t, metadata.EncryptionInfo.MacKey, installer.encryptionInfo.MacKey)
	assert.Equal(t, metadata.EncryptionInfo.InitializationVector, installer.encryptionInfo.InitializationVector)
	assert.Equal(t, metadata.EncryptionInfo.Mac, installer.encryptionInfo.Mac)
	assert.Equal(t, metadata.EncryptionInfo.FileDigest, installer.encryptionInfo.FileDigest)
	assert.Equal(t, metadata.EncryptionInfo.FileDigestAlgorithm, installer.encryptionInfo.FileDigestAlgorithm)
	assert.Equal(t, metadata.EncryptionInfo.ProfileIdentifier, installer.encryptionInfo.ProfileIdentifier)

	extractedContent, err := os.ReadFile(installer.encryptedFilePath)
	require.NoError(t, err)
	assert.Equal(t, encryptedContent, extractedContent, "the packaged ciphertext must not be encrypted again")

	contentFile := installer.contentFile()
	assert.Equal(t, metadata.FileName, *contentFile.GetName())
	assert.Equal(t, metadata.UnencryptedContentSize, *contentFile.GetSize())
	assert.Equal(t, int64(len(encryptedContent)), *contentFile.GetSizeEncrypted())
}

func TestUnitResourceWin32App_02_RejectInvalidIntuneWinPackages(t *testing.T) {
	tests := []struct {
		name        string
		mutate      func(*intuneWinDetectionMetadata)
		contents    func(intuneWinDetectionMetadata) map[string][]byte
		wantMessage string
	}{
		{
			name: "missing encrypted content",
			contents: func(metadata intuneWinDetectionMetadata) map[string][]byte {
				return nil
			},
			wantMessage: "encrypted content",
		},
		{
			name: "mismatched encrypted filename",
			contents: func(metadata intuneWinDetectionMetadata) map[string][]byte {
				return map[string][]byte{"IntuneWinPackage/Contents/different.intunewin": []byte("ciphertext")}
			},
			wantMessage: "encrypted content",
		},
		{
			name: "path traversal filename",
			mutate: func(metadata *intuneWinDetectionMetadata) {
				metadata.FileName = "../escape.intunewin"
			},
			wantMessage: "file name",
		},
		{
			name: "missing unencrypted size",
			mutate: func(metadata *intuneWinDetectionMetadata) {
				metadata.UnencryptedContentSize = 0
			},
			wantMessage: "unencrypted content size",
		},
		{
			name: "invalid base64 encryption key",
			mutate: func(metadata *intuneWinDetectionMetadata) {
				metadata.EncryptionInfo.EncryptionKey = "not-base64!"
			},
			wantMessage: "encryption key",
		},
		{
			name: "invalid initialization vector length",
			mutate: func(metadata *intuneWinDetectionMetadata) {
				metadata.EncryptionInfo.InitializationVector = base64.StdEncoding.EncodeToString([]byte("short"))
			},
			wantMessage: "initialization vector",
		},
		{
			name: "missing encryption profile",
			mutate: func(metadata *intuneWinDetectionMetadata) {
				metadata.EncryptionInfo.ProfileIdentifier = ""
			},
			wantMessage: "profile identifier",
		},
		{
			name: "unsupported digest algorithm",
			mutate: func(metadata *intuneWinDetectionMetadata) {
				metadata.EncryptionInfo.FileDigestAlgorithm = "SHA1"
			},
			wantMessage: "digest algorithm",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			metadata := validIntuneWinMetadata()
			if test.mutate != nil {
				test.mutate(&metadata)
			}

			contents := map[string][]byte{"IntuneWinPackage/Contents/" + metadata.FileName: []byte("ciphertext")}
			if test.contents != nil {
				contents = test.contents(metadata)
			}

			installer, err := extractIntuneWinPackage(writeIntuneWinFixture(t, metadata, contents))
			require.ErrorContains(t, err, test.wantMessage)
			assert.Nil(t, installer)
		})
	}

	t.Run("missing detection metadata", func(t *testing.T) {
		packagePath := filepath.Join(t.TempDir(), "missing-metadata.intunewin")
		file, err := os.Create(packagePath)
		require.NoError(t, err)
		archive := zip.NewWriter(file)
		entry, err := archive.Create("IntuneWinPackage/Contents/app.intunewin")
		require.NoError(t, err)
		_, err = entry.Write([]byte("ciphertext"))
		require.NoError(t, err)
		require.NoError(t, archive.Close())
		require.NoError(t, file.Close())

		_, err = extractIntuneWinPackage(packagePath)
		require.ErrorContains(t, err, "Detection.xml")
	})

	t.Run("not a zip archive", func(t *testing.T) {
		packagePath := filepath.Join(t.TempDir(), "invalid.intunewin")
		require.NoError(t, os.WriteFile(packagePath, []byte("not a zip"), 0o600))

		_, err := extractIntuneWinPackage(packagePath)
		require.Error(t, err)
	})
}

func TestUnitResourceWin32App_03_InstallerChangesDoNotRequireReplacement(t *testing.T) {
	installerSchema := commonschema.MobileAppWin32LobInstallerMetadataSchema()

	for _, attributeName := range []string{"installer_file_path_source", "installer_url_source"} {
		t.Run(attributeName, func(t *testing.T) {
			attribute, ok := installerSchema.Attributes[attributeName].(schema.StringAttribute)
			require.True(t, ok)
			assert.Len(t, attribute.PlanModifiers, 1, "only the existing unknown-value state modifier should remain")

			for _, modifier := range attribute.PlanModifiers {
				assert.NotContains(t, modifier.Description(context.Background()), "destroy and recreate", "installer changes must preserve the application ID")
			}
		})
	}

	macOSSchema := commonschema.MobileAppMacOSPkgInstallerMetadataSchema()
	macOSAttribute := macOSSchema.Attributes["installer_file_path_source"].(schema.StringAttribute)
	assert.Len(t, macOSAttribute.PlanModifiers, 2, "other application types must keep their existing replacement behavior")
}

func TestUnitResourceWin32App_04_PublishNewContentVersionWithoutReplacingApplication(t *testing.T) {
	metadata := validIntuneWinMetadata()
	encryptedContent := []byte("preserve this exact encrypted payload")
	installer, err := extractIntuneWinPackage(writeIntuneWinFixture(t, metadata, map[string][]byte{
		"IntuneWinPackage/Contents/" + metadata.FileName: encryptedContent,
	}))
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, installer.cleanup()) })

	const existingAppID = "existing-application-id"
	const contentVersionID = "2"
	const contentFileID = "content-file-id"

	var uploadedContent []byte
	var committedMetadata map[string]any
	var committedVersion string
	fileCommitted := false
	var server *httptest.Server
	server = httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		basePath := "/deviceAppManagement/mobileApps/" + existingAppID
		versionPath := basePath + "/graph.win32LobApp/contentVersions"
		filePath := versionPath + "/" + contentVersionID + "/files/" + contentFileID

		switch {
		case request.Method == http.MethodPost && request.URL.Path == versionPath:
			require.NoError(t, json.NewEncoder(writer).Encode(map[string]any{"id": contentVersionID}))
		case request.Method == http.MethodPost && request.URL.Path == versionPath+"/"+contentVersionID+"/files":
			var fileMetadata map[string]any
			decodeGraphRequestBody(t, request, &fileMetadata)
			assert.Equal(t, metadata.FileName, fileMetadata["name"])
			assert.Equal(t, float64(metadata.UnencryptedContentSize), fileMetadata["size"])
			assert.Equal(t, float64(len(encryptedContent)), fileMetadata["sizeEncrypted"])
			require.NoError(t, json.NewEncoder(writer).Encode(map[string]any{"id": contentFileID}))
		case request.Method == http.MethodGet && request.URL.Path == filePath:
			uploadState := "azureStorageUriRequestSuccess"
			if fileCommitted {
				uploadState = "commitFileSuccess"
			}
			require.NoError(t, json.NewEncoder(writer).Encode(map[string]any{
				"id":              contentFileID,
				"uploadState":     uploadState,
				"azureStorageUri": server.URL + "/blob?sig=test",
			}))
		case request.Method == http.MethodPut && request.URL.Path == "/blob":
			if request.URL.Query().Get("comp") == "block" {
				uploadedContent, err = io.ReadAll(request.Body)
				require.NoError(t, err)
			}
			writer.WriteHeader(http.StatusCreated)
		case request.Method == http.MethodPost && request.URL.Path == filePath+"/commit":
			var commitBody map[string]any
			decodeGraphRequestBody(t, request, &commitBody)
			var ok bool
			committedMetadata, ok = commitBody["fileEncryptionInfo"].(map[string]any)
			require.True(t, ok)
			fileCommitted = true
			writer.WriteHeader(http.StatusNoContent)
		case request.Method == http.MethodPatch && request.URL.Path == basePath:
			var updateBody map[string]any
			decodeGraphRequestBody(t, request, &updateBody)
			committedVersion, _ = updateBody["committedContentVersion"].(string)
			require.NoError(t, json.NewEncoder(writer).Encode(map[string]any{
				"@odata.type":             "#microsoft.graph.win32LobApp",
				"id":                      existingAppID,
				"committedContentVersion": committedVersion,
			}))
		default:
			t.Errorf("unexpected Graph operation: %s %s", request.Method, request.URL.Path)
			writer.WriteHeader(http.StatusNotFound)
		}
	}))
	t.Cleanup(server.Close)

	adapter, err := kiotahttp.NewNetHttpRequestAdapter(&authentication.AnonymousAuthenticationProvider{})
	require.NoError(t, err)
	adapter.SetBaseUrl(server.URL)
	resourceUnderTest := &Win32LobAppResource{client: msgraphbetasdk.NewGraphServiceClient(adapter)}
	response := &resource.UpdateResponse{}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	assert.True(t, resourceUnderTest.publishContentVersion(ctx, existingAppID, installer, response, constants.TfOperationUpdate))
	assert.False(t, response.Diagnostics.HasError())
	assert.Equal(t, encryptedContent, uploadedContent)
	assert.Equal(t, contentVersionID, committedVersion)
	require.NotNil(t, committedMetadata)
	assert.Equal(t, metadata.EncryptionInfo.EncryptionKey, committedMetadata["encryptionKey"])
	assert.Equal(t, metadata.EncryptionInfo.MacKey, committedMetadata["macKey"])
	assert.Equal(t, metadata.EncryptionInfo.InitializationVector, committedMetadata["initializationVector"])
	assert.Equal(t, metadata.EncryptionInfo.Mac, committedMetadata["mac"])
	assert.Equal(t, metadata.EncryptionInfo.FileDigest, committedMetadata["fileDigest"])
	assert.Equal(t, metadata.EncryptionInfo.FileDigestAlgorithm, committedMetadata["fileDigestAlgorithm"])
	assert.Equal(t, metadata.EncryptionInfo.ProfileIdentifier, committedMetadata["profileIdentifier"])
}

func TestUnitResourceWin32App_05_ConstructResourceUsesEncryptedInnerFilename(t *testing.T) {
	metadata := validIntuneWinMetadata()
	packagePath := writeIntuneWinFixture(t, metadata, map[string][]byte{
		"IntuneWinPackage/Contents/" + metadata.FileName: []byte("ciphertext"),
	})

	requestBody, err := constructResource(context.Background(), &Win32LobAppResourceModel{
		FileName:        types.StringValue("outer-package.intunewin"),
		RoleScopeTagIds: types.SetNull(types.StringType),
	}, packagePath, metadata.FileName)
	require.NoError(t, err)
	require.NotNil(t, requestBody.GetFileName())
	assert.Equal(t, metadata.FileName, *requestBody.GetFileName())
}

func TestUnitResourceWin32App_06_PreserveConfiguredOuterFilename(t *testing.T) {
	remoteApp := graphmodels.NewWin32LobApp()
	applicationID := "existing-application-id"
	innerFileName := "inner-installer.intunewin"
	remoteApp.SetId(&applicationID)
	remoteApp.SetFileName(&innerFileName)

	t.Run("existing configuration retains the outer archive filename", func(t *testing.T) {
		state := &Win32LobAppResourceModel{FileName: types.StringValue("outer-package.intunewin")}
		MapRemoteResourceStateToTerraform(context.Background(), state, remoteApp)
		assert.Equal(t, "outer-package.intunewin", state.FileName.ValueString())
	})

	t.Run("imports use the Graph filename", func(t *testing.T) {
		state := &Win32LobAppResourceModel{FileName: types.StringNull()}
		MapRemoteResourceStateToTerraform(context.Background(), state, remoteApp)
		assert.Equal(t, innerFileName, state.FileName.ValueString())
	})
}

func validIntuneWinMetadata() intuneWinDetectionMetadata {
	encoded := func(value byte, count int) string {
		return base64.StdEncoding.EncodeToString(bytes.Repeat([]byte{value}, count))
	}

	return intuneWinDetectionMetadata{
		XMLName:                xml.Name{Local: "ApplicationInfo"},
		FileName:               "inner-installer.intunewin",
		UnencryptedContentSize: 4096,
		EncryptionInfo: intuneWinEncryptionMetadata{
			EncryptionKey:        encoded(1, 32),
			MacKey:               encoded(2, 32),
			InitializationVector: encoded(3, 16),
			Mac:                  encoded(4, 32),
			FileDigest:           encoded(5, 32),
			FileDigestAlgorithm:  "SHA256",
			ProfileIdentifier:    "ProfileVersion1",
		},
	}
}

func decodeGraphRequestBody(t *testing.T, request *http.Request, target any) {
	t.Helper()

	var body io.Reader = request.Body
	if request.Header.Get("Content-Encoding") == "gzip" {
		gzipReader, err := gzip.NewReader(request.Body)
		require.NoError(t, err)
		defer gzipReader.Close()
		body = gzipReader
	}

	require.NoError(t, json.NewDecoder(body).Decode(target))
}

func writeIntuneWinFixture(t *testing.T, metadata intuneWinDetectionMetadata, contents map[string][]byte) string {
	t.Helper()

	packagePath := filepath.Join(t.TempDir(), "outer-package.intunewin")
	file, err := os.Create(packagePath)
	require.NoError(t, err)
	archive := zip.NewWriter(file)

	metadataEntry, err := archive.Create("IntuneWinPackage/Metadata/Detection.xml")
	require.NoError(t, err)
	require.NoError(t, xml.NewEncoder(metadataEntry).Encode(metadata))

	for name, content := range contents {
		entry, err := archive.Create(name)
		require.NoError(t, err)
		_, err = entry.Write(content)
		require.NoError(t, err)
	}

	require.NoError(t, archive.Close())
	require.NoError(t, file.Close())
	return packagePath
}
