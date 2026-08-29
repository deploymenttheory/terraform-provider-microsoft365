package graphBetaWin32App

import (
	"archive/zip"
	"bytes"
	"compress/gzip"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema/graph_beta/device_and_app_management"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
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
			assert.Empty(t, attribute.PlanModifiers, "unknown configured sources must not reuse the previous package")

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
	for _, plainZip := range []bool{false, true} {
		t.Run(fmt.Sprintf("plain_zip=%t", plainZip), func(t *testing.T) {
			metadata := validIntuneWinMetadata()
			var installer *preparedWin32Content
			var err error
			if plainZip {
				installer, err = prepareZipContent(context.Background(), writePlainZipFixture(t, "payload.zip", map[string][]byte{"setup.cmd": []byte("exit /b 0")}))
			} else {
				installer, err = extractIntuneWinPackage(writeIntuneWinFixture(t, metadata, map[string][]byte{
					"IntuneWinPackage/Contents/" + metadata.FileName: []byte("preserve this exact encrypted payload"),
				}))
			}
			require.NoError(t, err)
			t.Cleanup(func() { require.NoError(t, installer.cleanup()) })
			encryptedContent, err := os.ReadFile(installer.encryptedFilePath)
			require.NoError(t, err)

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
					assert.Equal(t, installer.fileName, fileMetadata["name"])
					assert.Equal(t, float64(installer.unencryptedSize), fileMetadata["size"])
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
			assert.Equal(t, installer.encryptionInfo.EncryptionKey, committedMetadata["encryptionKey"])
			assert.Equal(t, installer.encryptionInfo.MacKey, committedMetadata["macKey"])
			assert.Equal(t, installer.encryptionInfo.InitializationVector, committedMetadata["initializationVector"])
			assert.Equal(t, installer.encryptionInfo.Mac, committedMetadata["mac"])
			assert.Equal(t, installer.encryptionInfo.FileDigest, committedMetadata["fileDigest"])
			assert.Equal(t, installer.encryptionInfo.FileDigestAlgorithm, committedMetadata["fileDigestAlgorithm"])
			assert.Equal(t, installer.encryptionInfo.ProfileIdentifier, committedMetadata["profileIdentifier"])
		})
	}
}

func TestUnitResourceWin32App_05_ConstructResourceUsesEncryptedInnerFilename(t *testing.T) {
	metadata := validIntuneWinMetadata()
	requestBody, err := constructResource(context.Background(), &Win32LobAppResourceModel{
		FileName:        types.StringValue("outer-package.intunewin"),
		RoleScopeTagIds: types.SetNull(types.StringType),
	}, metadata.FileName)
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

func writePlainZipFixture(t *testing.T, name string, entries map[string][]byte) string {
	t.Helper()
	filePath := filepath.Join(t.TempDir(), name)
	file, err := os.Create(filePath)
	require.NoError(t, err)
	archive := zip.NewWriter(file)
	for name, data := range entries {
		entry, err := archive.Create(name)
		require.NoError(t, err)
		_, err = entry.Write(data)
		require.NoError(t, err)
	}
	require.NoError(t, archive.Close())
	require.NoError(t, file.Close())
	return filePath
}

func assertDecryptsTo(t *testing.T, installer *preparedWin32Content, plain []byte) {
	t.Helper()
	decode := func(value string) []byte {
		decoded, err := base64.StdEncoding.DecodeString(value)
		require.NoError(t, err)
		return decoded
	}
	encrypted, err := os.ReadFile(installer.encryptedFilePath)
	require.NoError(t, err)
	require.Greater(t, len(encrypted), 48)
	info := installer.encryptionInfo
	mac := hmac.New(sha256.New, decode(info.MacKey))
	_, err = mac.Write(encrypted[32:])
	require.NoError(t, err)
	assert.Equal(t, mac.Sum(nil), encrypted[:32])
	assert.Equal(t, decode(info.Mac), encrypted[:32])
	assert.Equal(t, decode(info.InitializationVector), encrypted[32:48])
	block, err := aes.NewCipher(decode(info.EncryptionKey))
	require.NoError(t, err)
	decrypted := make([]byte, len(encrypted)-48)
	cipher.NewCBCDecrypter(block, encrypted[32:48]).CryptBlocks(decrypted, encrypted[48:])
	padding := int(decrypted[len(decrypted)-1])
	require.True(t, padding > 0 && padding <= aes.BlockSize)
	assert.Equal(t, bytes.Repeat([]byte{byte(padding)}, padding), decrypted[len(decrypted)-padding:])
	assert.Equal(t, plain, decrypted[:len(decrypted)-padding])
	digest := sha256.Sum256(plain)
	assert.Equal(t, digest[:], decode(info.FileDigest))
	assert.Equal(t, int64(len(plain)), installer.unencryptedSize)
	assert.Equal(t, int64(len(encrypted)), installer.encryptedSize)
}

func TestUnitResourceWin32App_07_PreparePlainZip(t *testing.T) {
	for _, name := range []string{"source.zip", "legacy.intunewin"} {
		t.Run(name, func(t *testing.T) {
			source := writePlainZipFixture(t, name, map[string][]byte{"setup.cmd": []byte("exit /b 0")})
			plain, err := os.ReadFile(source)
			require.NoError(t, err)
			// An existing sibling must never be overwritten or removed.
			require.NoError(t, os.WriteFile(source+".bin", []byte("unrelated"), 0600))
			installer, err := prepareZipContent(context.Background(), source)
			require.NoError(t, err)
			assertDecryptsTo(t, installer, plain)
			assert.Equal(t, strings.TrimSuffix(name, filepath.Ext(name))+".intunewin", installer.fileName)
			temporaryDirectory := installer.temporaryDirectory
			require.NoError(t, installer.cleanup())
			_, err = os.Stat(temporaryDirectory)
			assert.True(t, os.IsNotExist(err))
			after, err := os.ReadFile(source)
			require.NoError(t, err)
			assert.Equal(t, plain, after)
			sibling, err := os.ReadFile(source + ".bin")
			require.NoError(t, err)
			assert.Equal(t, "unrelated", string(sibling))
		})
	}
}

func TestUnitResourceWin32App_08_RejectWrongZipMode(t *testing.T) {
	for _, entry := range []string{"IntuneWinPackage/Contents/content.intunewin", "IntuneWinPackage/Metadata/Detection.xml", "../setup.cmd", "C:/setup.cmd"} {
		t.Run(entry, func(t *testing.T) {
			installer, err := prepareZipContent(context.Background(), writePlainZipFixture(t, "bad.zip", map[string][]byte{entry: []byte("bad")}))
			require.Error(t, err)
			assert.Nil(t, installer)
		})
	}
	_, err := prepareZipContent(context.Background(), writePlainZipFixture(t, "empty.zip", nil))
	require.ErrorContains(t, err, "no files")
	plain := writePlainZipFixture(t, "legacy.intunewin", map[string][]byte{"setup.cmd": []byte("exit /b 0")})
	_, err = prepareInstaller(context.Background(), &Win32LobAppResourceModel{AppInstaller: installerObject(types.StringValue(plain), types.StringNull())})
	require.ErrorContains(t, err, "use app_installer_zip")
}

func installerObject(local, remote types.String) types.Object {
	return types.ObjectValueMust(map[string]attr.Type{
		"installer_file_path_source": types.StringType, "installer_url_source": types.StringType,
	}, map[string]attr.Value{"installer_file_path_source": local, "installer_url_source": remote})
}

func TestUnitResourceWin32App_09_SelectSource(t *testing.T) {
	local := installerObject(types.StringValue("app.intunewin"), types.StringNull())
	remote := installerObject(types.StringNull(), types.StringValue("https://example.com/app.zip"))
	for _, test := range []struct {
		name               string
		model              Win32LobAppResourceModel
		wantZip, wantError bool
	}{
		{name: "prepackaged", model: Win32LobAppResourceModel{AppInstaller: local}},
		{name: "zip URL", model: Win32LobAppResourceModel{AppInstallerZip: remote}, wantZip: true},
		{name: "both formats", model: Win32LobAppResourceModel{AppInstaller: local, AppInstallerZip: remote}, wantError: true},
		{name: "no source", wantError: true},
		{name: "unknown", model: Win32LobAppResourceModel{AppInstaller: types.ObjectUnknown(local.AttributeTypes(context.Background()))}, wantError: true},
		{name: "unknown path", model: Win32LobAppResourceModel{AppInstaller: installerObject(types.StringUnknown(), types.StringNull())}, wantError: true},
		{name: "both transports", model: Win32LobAppResourceModel{AppInstaller: installerObject(types.StringValue("app.intunewin"), types.StringValue("https://example.com/app"))}, wantError: true},
		{name: "empty block", model: Win32LobAppResourceModel{AppInstaller: installerObject(types.StringNull(), types.StringNull())}, wantError: true},
	} {
		t.Run(test.name, func(t *testing.T) {
			_, zip, err := selectInstallerSource(&test.model)
			assert.Equal(t, test.wantError, err != nil)
			if !test.wantError {
				assert.Equal(t, test.wantZip, zip)
			}
		})
	}
}

func TestUnitResourceWin32App_10_PrepareURLSources(t *testing.T) {
	metadata := validIntuneWinMetadata()
	for _, plainZip := range []bool{false, true} {
		t.Run(fmt.Sprintf("plain_zip=%t", plainZip), func(t *testing.T) {
			var source string
			if plainZip {
				source = writePlainZipFixture(t, "payload.zip", map[string][]byte{"setup.cmd": []byte("exit /b 0")})
			} else {
				source = writeIntuneWinFixture(t, metadata, map[string][]byte{"IntuneWinPackage/Contents/" + metadata.FileName: []byte("ciphertext")})
			}
			data, err := os.ReadFile(source)
			require.NoError(t, err)
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { _, err := w.Write(data); assert.NoError(t, err) }))
			defer server.Close()
			object := installerObject(types.StringNull(), types.StringValue(server.URL+"/"+filepath.Base(source)))
			model := Win32LobAppResourceModel{AppInstaller: object}
			if plainZip {
				model.AppInstaller = types.Object{}
				model.AppInstallerZip = object
			}
			installer, err := prepareInstaller(context.Background(), &model)
			require.NoError(t, err)
			defer func() { require.NoError(t, installer.cleanup()) }()
			if plainZip {
				assertDecryptsTo(t, installer, data)
			} else {
				payload, err := os.ReadFile(installer.encryptedFilePath)
				require.NoError(t, err)
				assert.Equal(t, "ciphertext", string(payload))
			}
		})
	}
}

func TestUnitResourceWin32App_11_MetadataUpdateDoesNotChangeContent(t *testing.T) {
	body, err := constructResource(context.Background(), &Win32LobAppResourceModel{
		FileName: types.StringValue("outer.intunewin"), CommittedContentVersion: types.StringValue("old"), RoleScopeTagIds: types.SetNull(types.StringType),
	}, "")
	require.NoError(t, err)
	assert.Nil(t, body.GetFileName())
	assert.Nil(t, body.GetCommittedContentVersion())
}

func TestUnitResourceWin32App_12_SourceSchemaValidation(t *testing.T) {
	ctx := context.Background()
	local := installerObject(types.StringValue("app.intunewin"), types.StringNull())
	remote := installerObject(types.StringNull(), types.StringValue("https://example.com/app.zip"))
	null := types.ObjectNull(local.AttributeTypes(ctx))
	sourceSchema := schema.Schema{Attributes: map[string]schema.Attribute{
		"app_installer":     commonschema.MobileAppWin32LobInstallerMetadataSchema(),
		"app_installer_zip": commonschema.MobileAppWin32ZipInstallerMetadataSchema(),
	}}
	for _, test := range []struct {
		name          string
		packaged, zip types.Object
		wantError     bool
	}{
		{"package path", local, null, false}, {"ZIP URL", null, remote, false}, {"both modes", local, remote, true},
		{"both transports", installerObject(types.StringValue("app.intunewin"), types.StringValue("https://example.com/app")), null, true},
		{"empty block", installerObject(types.StringNull(), types.StringNull()), null, true},
		{"import without source", null, null, false},
		{"unknown source", types.ObjectUnknown(local.AttributeTypes(ctx)), null, false},
		{"unknown path", installerObject(types.StringUnknown(), types.StringNull()), null, false},
	} {
		t.Run(test.name, func(t *testing.T) {
			value := types.ObjectValueMust(map[string]attr.Type{"app_installer": local.Type(ctx), "app_installer_zip": local.Type(ctx)}, map[string]attr.Value{"app_installer": test.packaged, "app_installer_zip": test.zip})
			raw, err := value.ToTerraformValue(ctx)
			require.NoError(t, err)
			config := tfsdk.Config{Schema: sourceSchema, Raw: raw}
			hasError := false
			for name, value := range map[string]types.Object{"app_installer": test.packaged, "app_installer_zip": test.zip} {
				attribute := sourceSchema.Attributes[name].(schema.SingleNestedAttribute)
				for _, v := range attribute.Validators {
					response := validator.ObjectResponse{}
					v.ValidateObject(ctx, validator.ObjectRequest{Config: config, ConfigValue: value, Path: path.Root(name), PathExpression: path.MatchRoot(name)}, &response)
					hasError = hasError || response.Diagnostics.HasError()
				}
				if value.IsNull() || value.IsUnknown() {
					continue
				}
				for child, childAttribute := range attribute.Attributes {
					for _, v := range childAttribute.(schema.StringAttribute).Validators {
						response := validator.StringResponse{}
						v.ValidateString(ctx, validator.StringRequest{Config: config, ConfigValue: value.Attributes()[child].(types.String), Path: path.Root(name).AtName(child), PathExpression: path.MatchRoot(name).AtName(child)}, &response)
						hasError = hasError || response.Diagnostics.HasError()
					}
				}
			}
			assert.Equal(t, test.wantError, hasError)
		})
	}
}

func TestUnitResourceWin32App_13_PlanContentVersionChanges(t *testing.T) {
	ctx := context.Background()
	source := installerObject(types.StringValue("v1.intunewin"), types.StringNull())
	sourceType := source.Type(ctx)
	null := types.ObjectNull(source.AttributeTypes(ctx))
	oldVersion := types.ListValueMust(types.StringType, []attr.Value{types.StringValue("1")})
	planSchema := schema.Schema{Attributes: map[string]schema.Attribute{
		"app_installer":     commonschema.MobileAppWin32LobInstallerMetadataSchema(),
		"app_installer_zip": commonschema.MobileAppWin32ZipInstallerMetadataSchema(),
		"content_version":   schema.ListAttribute{Optional: true, Computed: true, ElementType: types.StringType},
	}}
	raw := func(packaged, zip types.Object, version types.List) tftypes.Value {
		value, err := types.ObjectValueMust(map[string]attr.Type{
			"app_installer": sourceType, "app_installer_zip": sourceType, "content_version": oldVersion.Type(ctx),
		}, map[string]attr.Value{"app_installer": packaged, "app_installer_zip": zip, "content_version": version}).ToTerraformValue(ctx)
		require.NoError(t, err)
		return value
	}
	for _, test := range []struct {
		name          string
		packaged, zip types.Object
		changed       bool
	}{
		{"unchanged", source, null, false},
		{"path change", installerObject(types.StringValue("v2.intunewin"), types.StringNull()), null, true},
		{"URL change", installerObject(types.StringNull(), types.StringValue("https://example.com/v2")), null, true},
		{"switch to ZIP", null, source, true},
		{"unknown new path", installerObject(types.StringUnknown(), types.StringNull()), null, true},
		{"forget local source", null, null, false},
	} {
		t.Run(test.name, func(t *testing.T) {
			planned := tfsdk.Plan{Schema: planSchema, Raw: raw(test.packaged, test.zip, oldVersion)}
			response := resource.ModifyPlanResponse{Plan: planned}
			(&Win32LobAppResource{}).ModifyPlan(ctx, resource.ModifyPlanRequest{
				Plan: planned, State: tfsdk.State{Schema: planSchema, Raw: raw(source, null, oldVersion)},
				Config: tfsdk.Config{Schema: planSchema, Raw: raw(test.packaged, test.zip, types.ListNull(types.StringType))},
			}, &response)
			require.False(t, response.Diagnostics.HasError(), "%v", response.Diagnostics)
			var version types.List
			require.False(t, response.Plan.GetAttribute(ctx, path.Root("content_version"), &version).HasError())
			assert.Equal(t, test.changed, version.IsUnknown())
		})
	}
	t.Run("creation requires a source", func(t *testing.T) {
		planned := tfsdk.Plan{Schema: planSchema, Raw: raw(null, null, types.ListNull(types.StringType))}
		response := resource.ModifyPlanResponse{Plan: planned}
		(&Win32LobAppResource{}).ModifyPlan(ctx, resource.ModifyPlanRequest{
			Plan:  planned,
			State: tfsdk.State{Schema: planSchema, Raw: tftypes.NewValue(planned.Raw.Type(), nil)},
		}, &response)
		require.True(t, response.Diagnostics.HasError())
		assert.Contains(t, response.Diagnostics.Errors()[0].Summary(), "Missing installer source")
	})

}

func TestUnitResourceWin32App_14_RejectPasswordProtectedZip(t *testing.T) {
	source := filepath.Join(t.TempDir(), "encrypted.zip")
	file, err := os.Create(source)
	require.NoError(t, err)
	archive := zip.NewWriter(file)
	entry, err := archive.CreateHeader(&zip.FileHeader{Name: "setup.cmd", Flags: 1, Method: zip.Store})
	require.NoError(t, err)
	_, err = entry.Write([]byte("password-protected ZIP entry"))
	require.NoError(t, err)
	require.NoError(t, archive.Close())
	require.NoError(t, file.Close())
	installer, err := prepareZipContent(context.Background(), source)
	require.ErrorContains(t, err, "encrypted")
	assert.Nil(t, installer)
}
