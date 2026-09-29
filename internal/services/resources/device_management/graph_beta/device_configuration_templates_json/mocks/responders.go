package mocks

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/jarcoal/httpmock"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
)

const endpoint = "https://graph.microsoft.com/beta/deviceManagement/deviceConfigurations"

var errInvalidFixture = errors.New("invalid device configuration fixture")

type DeviceConfigurationTemplatesJsonMock struct {
	sync.Mutex
	profiles    map[string]map[string]any
	assignments map[string][]any
	secrets     map[string]string
}

var _ mocks.MockRegistrar = (*DeviceConfigurationTemplatesJsonMock)(nil)

func init() {
	mocks.GlobalRegistry.Register(
		"device_configuration_templates_json",
		&DeviceConfigurationTemplatesJsonMock{},
	)
}

func (m *DeviceConfigurationTemplatesJsonMock) RegisterMocks() {
	m.CleanupMockState()
	httpmock.RegisterResponder("POST", endpoint, m.create)
	for _, method := range []string{"GET", "PATCH", "DELETE", "POST", "PUT"} {
		httpmock.RegisterResponder(
			method,
			`=~^https://graph\.microsoft\.com/beta/deviceManagement/deviceConfigurations/[^/]+.*$`,
			m.request,
		)
	}
}

func (m *DeviceConfigurationTemplatesJsonMock) RegisterErrorMocks() {
	m.RegisterMocks()
	httpmock.RegisterResponder("POST", endpoint, func(_ *http.Request) (*http.Response, error) {
		return errorResponse(http.StatusForbidden)
	})
}

func (m *DeviceConfigurationTemplatesJsonMock) CleanupMockState() {
	m.Lock()
	defer m.Unlock()
	m.profiles = make(map[string]map[string]any)
	m.assignments = make(map[string][]any)
	m.secrets = make(map[string]string)
}

func (m *DeviceConfigurationTemplatesJsonMock) create(req *http.Request) (*http.Response, error) {
	var body map[string]any
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		return nil, fmt.Errorf("decode profile request: %w", err)
	}
	typ, _ := body["@odata.type"].(string)
	// Graph rejects embedded certificate bindings for iOS wired profiles;
	// the provider must establish them through the typed $ref endpoints.
	if typ == "#microsoft.graph.iosWiredNetworkConfiguration" {
		for _, name := range []string{"rootCertificateForServerValidation", "identityCertificateForClientAuthentication"} {
			if _, exists := body[name+"@odata.bind"]; exists {
				return errorResponse(http.StatusBadRequest)
			}
		}
	}
	fixture, err := loadProfileFixture(typ)
	if err != nil {
		return nil, err
	}
	m.Lock()
	defer m.Unlock()
	for key, value := range body {
		fixture[key] = value
	}
	if typ == "#microsoft.graph.androidDeviceOwnerEnterpriseWiFiConfiguration" {
		delete(fixture, "rootCertificatesForServerValidation@odata.bind")
		delete(fixture, "rootCertificateForServerValidation@odata.bind")
	}
	id := uuid.NewString()
	fixture["id"] = id
	m.profiles[id] = fixture
	m.assignments[id] = []any{}
	response, err := m.profileResponse(fixture)
	if response != nil {
		response.StatusCode = http.StatusCreated
	}
	return response, err
}

func loadProfileFixture(typ string) (map[string]any, error) {
	files, err := filepath.Glob("tests/responses/validate_create/post_device_configuration_*.json")
	if err != nil {
		return nil, fmt.Errorf("find profile fixtures: %w", err)
	}
	for _, file := range files {
		content, err := helpers.ParseJSONFile(filepath.Join("..", file))
		if err != nil {
			return nil, fmt.Errorf("load profile fixture: %w", err)
		}
		var response map[string]any
		if err := json.Unmarshal([]byte(content), &response); err != nil {
			return nil, fmt.Errorf("decode profile fixture: %w", err)
		}
		if response["@odata.type"] == typ {
			return response, nil
		}
	}
	return nil, fmt.Errorf("%w: no fixture for %s", errInvalidFixture, typ)
}

func (m *DeviceConfigurationTemplatesJsonMock) request(req *http.Request) (*http.Response, error) {
	m.Lock()
	defer m.Unlock()
	segments := strings.Split(
		strings.TrimPrefix(req.URL.Path, "/beta/deviceManagement/deviceConfigurations/"),
		"/",
	)
	id := segments[0]
	profile, exists := m.profiles[id]
	if !exists {
		return errorResponse(http.StatusNotFound)
	}
	if len(segments) > 1 {
		return m.relatedRequest(req, id, profile, strings.Join(segments[1:], "/"))
	}
	switch req.Method {
	case "GET":
		return m.profileResponse(profile)
	case "PATCH":
		var body map[string]any
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			return nil, fmt.Errorf("decode profile patch: %w", err)
		}
		for key, value := range body {
			profile[key] = value
		}
		return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
	case "DELETE":
		delete(m.profiles, id)
		delete(m.assignments, id)
		return httpmock.NewStringResponse(http.StatusOK, ""), nil
	default:
		return errorResponse(http.StatusBadRequest)
	}
}

func (m *DeviceConfigurationTemplatesJsonMock) relatedRequest(
	req *http.Request,
	id string,
	profile map[string]any,
	suffix string,
) (*http.Response, error) {
	switch {
	case suffix == "assign" && req.Method == "POST":
		var body struct {
			Assignments []any `json:"assignments"`
		}
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			return nil, fmt.Errorf("decode assignment request: %w", err)
		}
		m.assignments[id] = body.Assignments
		return httpmock.NewStringResponse(http.StatusOK, ""), nil
	case suffix == "assignments" && req.Method == "GET":
		return jsonResponse(http.StatusOK, map[string]any{"value": m.assignments[id]})
	case strings.HasPrefix(suffix, "microsoft.graph."):
		name := strings.Split(suffix, "/")[1]
		if name == "rootCertificate" &&
			(profile["@odata.type"] == "#microsoft.graph.androidDeviceOwnerImportedPFXCertificateProfile" ||
				profile["@odata.type"] == "#microsoft.graph.androidForWorkImportedPFXCertificateProfile") {
			return errorResponse(http.StatusBadRequest)
		}
		binding := profile[name+"@odata.bind"]
		if strings.HasSuffix(suffix, "/$ref") {
			if req.Method == "DELETE" {
				parts := strings.Split(suffix, "/")
				if len(parts) == 4 {
					remaining := []any{}
					values, _ := binding.([]any)
					for _, item := range values {
						if !strings.Contains(item.(string), parts[2]) {
							remaining = append(remaining, item)
						}
					}
					profile[name+"@odata.bind"] = remaining
				} else {
					delete(profile, name+"@odata.bind")
				}
			} else {
				var body map[string]string
				if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
					return nil, fmt.Errorf("decode relationship reference request: %w", err)
				}
				if req.Method == "PUT" {
					profile[name+"@odata.bind"] = body["@odata.id"]
				} else {
					values, _ := binding.([]any)
					profile[name+"@odata.bind"] = append(values, body["@odata.id"])
				}
			}
			return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
		}
		switch value := binding.(type) {
		case string:
			match := regexp.MustCompile(`[0-9a-fA-F]{8}-[0-9a-fA-F-]{27}`).FindString(value)
			root, exists := m.profiles[match]
			if !exists {
				return errorResponse(http.StatusNotFound)
			}
			return jsonResponse(http.StatusOK, root)
		case []any:
			values := make([]any, 0, len(value))
			for _, entry := range value {
				text, _ := entry.(string)
				match := regexp.MustCompile(`[0-9a-fA-F]{8}-[0-9a-fA-F-]{27}`).FindString(text)
				root, exists := m.profiles[match]
				if !exists {
					return errorResponse(http.StatusNotFound)
				}
				values = append(values, root)
			}
			return jsonResponse(http.StatusOK, map[string]any{"value": values})
		default:
			if strings.HasPrefix(name, "rootCertificates") {
				return jsonResponse(http.StatusOK, map[string]any{"value": []any{}})
			}
			return errorResponse(http.StatusNotFound)
		}
	case strings.HasPrefix(suffix, "getOmaSettingPlainTextValue("):
		reference := strings.TrimSuffix(
			strings.TrimPrefix(suffix, "getOmaSettingPlainTextValue(secretReferenceValueId='"),
			"')",
		)
		value, exists := m.secrets[reference]
		if !exists {
			return errorResponse(http.StatusNotFound)
		}
		return jsonResponse(http.StatusOK, map[string]any{"value": value})
	default:
		return errorResponse(http.StatusBadRequest)
	}
}

func (m *DeviceConfigurationTemplatesJsonMock) profileResponse(
	profile map[string]any,
) (*http.Response, error) {
	content, err := json.Marshal(profile)
	if err != nil {
		return nil, fmt.Errorf("encode mock profile: %w", err)
	}
	var response map[string]any
	if err := json.Unmarshal(content, &response); err != nil {
		return nil, fmt.Errorf("copy mock profile: %w", err)
	}
	if response["@odata.type"] == "#microsoft.graph.androidDeviceOwnerImportedPFXCertificateProfile" ||
		response["@odata.type"] == "#microsoft.graph.androidForWorkImportedPFXCertificateProfile" {
		usages, _ := response["extendedKeyUsages"].([]any)
		response["extendedKeyUsages"] = append(usages, map[string]any{
			"name": "Any Purpose", "objectIdentifier": "2.5.29.37.0",
		})
	}
	for key := range response {
		if strings.HasSuffix(key, "@odata.bind") {
			delete(response, key)
		}
	}
	if entries, ok := response["omaSettings"].([]any); ok {
		for _, entry := range entries {
			setting, ok := entry.(map[string]any)
			if !ok {
				return nil, errInvalidFixture
			}
			setting["isEncrypted"] = false
			setting["secretReferenceValueId"] = nil
			setting["isReadOnly"] = false
			if typ := setting["@odata.type"]; typ == "#microsoft.graph.omaSettingString" ||
				typ == "#microsoft.graph.omaSettingBase64" ||
				typ == "#microsoft.graph.omaSettingStringXml" {
				reference := fmt.Sprint(profile["id"], ":", setting["omaUri"])
				value, _ := setting["value"].(string)
				m.secrets[reference] = value
				if setting["@odata.type"] == "#microsoft.graph.omaSettingStringXml" {
					decoded, err := base64.StdEncoding.DecodeString(value)
					if err != nil {
						return nil, fmt.Errorf("decode mock XML setting: %w", err)
					}
					m.secrets[reference] = string(decoded)
				}
				setting["value"] = "****"
				if setting["@odata.type"] != "#microsoft.graph.omaSettingString" {
					setting["value"] = "KioqKg=="
				}
				setting["isEncrypted"] = true
				setting["secretReferenceValueId"] = reference
				delete(setting, "isReadOnly")
			}
		}
	}
	if key, exists := response["preSharedKey"]; exists {
		response["preSharedKey"] = nil
		if _, hasFlag := response["preSharedKeyIsSet"]; hasFlag {
			response["preSharedKeyIsSet"] = key != nil
		}
	}
	return jsonResponse(http.StatusOK, response)
}

func (m *DeviceConfigurationTemplatesJsonMock) SetRemoteProperty(name string, value any) {
	m.Lock()
	defer m.Unlock()
	for _, profile := range m.profiles {
		profile[name] = value
	}
}

func errorResponse(status int) (*http.Response, error) {
	content, err := helpers.ParseJSONFile(
		fmt.Sprintf("../tests/responses/validate_error/get_device_configuration_%d.json", status),
	)
	if err != nil {
		return nil, fmt.Errorf("load error response: %w", err)
	}
	return httpmock.NewStringResponse(status, content), nil
}

func jsonResponse(status int, body any) (*http.Response, error) {
	response, err := httpmock.NewJsonResponse(status, body)
	if err != nil {
		return nil, fmt.Errorf("encode mock response: %w", err)
	}
	return response, nil
}
