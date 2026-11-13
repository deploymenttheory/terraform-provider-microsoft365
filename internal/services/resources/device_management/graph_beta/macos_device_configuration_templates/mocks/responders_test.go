package mocks_test

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"

	configmocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/macos_device_configuration_templates/mocks"
)

func TestUnitMacosConfigurationMock_PreservesPayloadBytes(t *testing.T) {
	for _, configuration := range []struct {
		name      string
		odataType string
		field     string
	}{
		{"custom", "#microsoft.graph.macOSCustomConfiguration", "payload"},
		{"preference", "#microsoft.graph.macOSCustomAppConfiguration", "configurationXml"},
	} {
		for _, lineEnding := range []struct {
			name  string
			value string
		}{
			{"LF", "\n"},
			{"CRLF", "\r\n"},
		} {
			t.Run(configuration.name+"/"+lineEnding.name, func(t *testing.T) {
				httpmock.Activate()
				t.Cleanup(httpmock.DeactivateAndReset)
				mock := &configmocks.MacosDeviceConfigurationTemplatesMock{}
				mock.RegisterMocks()
				t.Cleanup(mock.CleanupMockState)

				payload := "<?xml version=\"1.0\"?>" + lineEnding.value +
					"<plist version=\"1.0\"><dict/></plist>" + lineEnding.value
				encoded := base64.StdEncoding.EncodeToString([]byte(payload))
				body, err := json.Marshal(map[string]any{
					"@odata.type":       configuration.odataType,
					"displayName":       "unit-test-payload-roundtrip",
					configuration.field: encoded,
				})
				if err != nil {
					t.Fatal(err)
				}

				const endpoint = "https://graph.microsoft.com/beta/deviceManagement/deviceConfigurations"
				request, err := http.NewRequestWithContext(
					t.Context(), http.MethodPost, endpoint, bytes.NewReader(body),
				)
				if err != nil {
					t.Fatal(err)
				}
				request.Header.Set("Content-Type", "application/json")
				response, err := http.DefaultClient.Do(request)
				if err != nil {
					t.Fatal(err)
				}
				defer response.Body.Close()
				if response.StatusCode != http.StatusCreated {
					t.Fatalf("create returned HTTP %d", response.StatusCode)
				}
				var created map[string]any
				if err := json.NewDecoder(response.Body).Decode(&created); err != nil {
					t.Fatal(err)
				}
				if created[configuration.field] != encoded {
					t.Errorf("create changed %s payload bytes", configuration.field)
				}

				id, ok := created["id"].(string)
				if !ok {
					t.Fatal("create did not return an ID")
				}
				request, err = http.NewRequestWithContext(
					t.Context(), http.MethodGet, endpoint+"/"+id+"?$expand=assignments", nil,
				)
				if err != nil {
					t.Fatal(err)
				}
				readResponse, err := http.DefaultClient.Do(request)
				if err != nil {
					t.Fatal(err)
				}
				defer readResponse.Body.Close()
				if readResponse.StatusCode != http.StatusOK {
					t.Fatalf("read returned HTTP %d", readResponse.StatusCode)
				}
				var read map[string]any
				if err := json.NewDecoder(readResponse.Body).Decode(&read); err != nil {
					t.Fatal(err)
				}
				if read[configuration.field] != encoded {
					t.Errorf("read changed %s payload bytes", configuration.field)
				}
			})
		}
	}
}
