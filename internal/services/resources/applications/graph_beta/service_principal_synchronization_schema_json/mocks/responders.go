package mocks

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"sync"

	"github.com/jarcoal/httpmock"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
)

const (
	// ServicePrincipalID is the ID of the service principal in the mock tenant.
	ServicePrincipalID = "11111111-1111-1111-1111-111111111111"

	// JobID is the ID of the synchronization job in the mock tenant.
	JobID = "scim.00000000000000000000000000000001.22222222-2222-2222-2222-222222222222"

	// MissingJobID is a synchronization job ID that does not exist in the mock tenant.
	MissingJobID = "scim.00000000000000000000000000000001.33333333-3333-3333-3333-333333333333"
)

var schemaURLPattern = regexp.MustCompile(`^https://graph\.microsoft\.com/beta/servicePrincipals/([^/?]+)/synchronization/jobs/([^/?]+)/schema(\?.*)?$`)

var mockState struct {
	sync.Mutex
	defaultSchema map[string]any
	schemas       map[string]map[string]any
	puts          []map[string]any
	deletes       int
	revision      int
	loadErr       error
}

func init() {
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))
	mocks.GlobalRegistry.Register("service_principal_synchronization_schema_json", &ServicePrincipalSynchronizationSchemaJsonMock{})
}

// ServicePrincipalSynchronizationSchemaJsonMock provides stateful httpmock responders for synchronization schemas.
type ServicePrincipalSynchronizationSchemaJsonMock struct{}

var _ mocks.MockRegistrar = (*ServicePrincipalSynchronizationSchemaJsonMock)(nil)

// loadJSONResponse loads a JSON fixture from tests/responses/validate_get.
func loadJSONResponse(filename string) (map[string]any, error) {
	jsonData, err := os.ReadFile(filepath.Join("tests", "responses", "validate_get", filename))
	if err != nil {
		return nil, fmt.Errorf("failed to load JSON response file: %w", err)
	}

	var response map[string]any
	if err := json.Unmarshal(jsonData, &response); err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}
	return response, nil
}

func errorResponse(status int, code, message string) (*http.Response, error) {
	resp, err := httpmock.NewJsonResponse(status, map[string]any{"error": map[string]any{"code": code, "message": message}})
	if err != nil {
		return nil, fmt.Errorf("build error response: %w", err)
	}
	return resp, nil
}

// deepCopy returns an independent copy of a decoded JSON value.
func deepCopy(value map[string]any) map[string]any {
	content, _ := json.Marshal(value)
	var copied map[string]any
	_ = json.Unmarshal(content, &copied)
	return copied
}

func schemaKey(req *http.Request) (string, string, string) {
	matches := schemaURLPattern.FindStringSubmatch(req.URL.String())
	servicePrincipalId, _ := url.PathUnescape(matches[1])
	jobId, _ := url.PathUnescape(matches[2])
	return servicePrincipalId + "/" + jobId, servicePrincipalId, jobId
}

// applyServerManagedValues mirrors Graph, which sets a new version on every directory and the Unsynchronized object
// mapping metadata from the mapping's enabled flag when it stores a schema.
func applyServerManagedValues(schema map[string]any, revision int) {
	version := fmt.Sprintf("Date:2026-01-01T00:00:%02d.0000000Z, ActivityId:00000000-0000-0000-0000-%012d", revision%60, revision)
	directories, _ := schema["directories"].([]any)
	for _, item := range directories {
		if directory, ok := item.(map[string]any); ok {
			directory["version"] = version
		}
	}
	rules, _ := schema["synchronizationRules"].([]any)
	for _, item := range rules {
		rule, _ := item.(map[string]any)
		objectMappings, _ := rule["objectMappings"].([]any)
		for _, mappingItem := range objectMappings {
			mapping, _ := mappingItem.(map[string]any)
			enabled, _ := mapping["enabled"].(bool)
			metadata, _ := mapping["metadata"].([]any)
			for _, entryItem := range metadata {
				entry, _ := entryItem.(map[string]any)
				if entry["key"] == "Unsynchronized" {
					entry["value"] = fmt.Sprintf("%t", !enabled)
				}
			}
		}
	}
	schema["version"] = version
}

func (m *ServicePrincipalSynchronizationSchemaJsonMock) RegisterMocks() {
	defaultSchema, loadErr := loadJSONResponse("get_synchronization_schema_default.json")

	mockState.Lock()
	mockState.defaultSchema = defaultSchema
	mockState.schemas = map[string]map[string]any{}
	if defaultSchema != nil {
		mockState.schemas[ServicePrincipalID+"/"+JobID] = deepCopy(defaultSchema)
	}
	mockState.puts = nil
	mockState.deletes = 0
	mockState.revision = 0
	mockState.loadErr = loadErr
	mockState.Unlock()

	// Get synchronization schema - GET /servicePrincipals/{id}/synchronization/jobs/{jobId}/schema
	httpmock.RegisterRegexpResponder("GET", schemaURLPattern, func(req *http.Request) (*http.Response, error) {
		mockState.Lock()
		defer mockState.Unlock()

		if mockState.loadErr != nil {
			return errorResponse(500, "InternalServerError", mockState.loadErr.Error())
		}

		key, servicePrincipalId, jobId := schemaKey(req)
		stored, ok := mockState.schemas[key]
		if !ok {
			return errorResponse(404, "NotFound", fmt.Sprintf("Synchronization job '%s' does not exist under service principal '%s'", jobId, servicePrincipalId))
		}
		return httpmock.NewJsonResponse(200, deepCopy(stored))
	})

	// Replace synchronization schema - PUT /servicePrincipals/{id}/synchronization/jobs/{jobId}/schema
	httpmock.RegisterRegexpResponder("PUT", schemaURLPattern, func(req *http.Request) (*http.Response, error) {
		content, err := io.ReadAll(req.Body)
		if err != nil {
			return errorResponse(400, "BadRequest", "Invalid request body")
		}
		var requestBody map[string]any
		if err := json.Unmarshal(content, &requestBody); err != nil {
			return errorResponse(400, "BadRequest", "Invalid request body")
		}

		mockState.Lock()
		defer mockState.Unlock()

		key, servicePrincipalId, jobId := schemaKey(req)
		stored, ok := mockState.schemas[key]
		if !ok {
			return errorResponse(404, "NotFound", fmt.Sprintf("Synchronization job '%s' does not exist under service principal '%s'", jobId, servicePrincipalId))
		}

		for property := range requestBody {
			if property != "synchronizationRules" && property != "directories" {
				message := fmt.Sprintf("schema : The property '%s' does not exist on type 'microsoft.synchronization.synchronizationSchema'.", property)
				return errorResponse(406, "RequestParameterInvalid", message)
			}
		}

		mockState.puts = append(mockState.puts, deepCopy(requestBody))
		mockState.revision++

		updated := deepCopy(stored)
		updated["synchronizationRules"] = deepCopy(requestBody)["synchronizationRules"]
		updated["directories"] = deepCopy(requestBody)["directories"]
		applyServerManagedValues(updated, mockState.revision)
		mockState.schemas[key] = updated

		return httpmock.NewStringResponse(204, ""), nil
	})

	// Reset synchronization schema - DELETE /servicePrincipals/{id}/synchronization/jobs/{jobId}/schema
	httpmock.RegisterRegexpResponder("DELETE", schemaURLPattern, func(req *http.Request) (*http.Response, error) {
		mockState.Lock()
		defer mockState.Unlock()

		key, servicePrincipalId, jobId := schemaKey(req)
		if _, ok := mockState.schemas[key]; !ok {
			return errorResponse(404, "NotFound", fmt.Sprintf("Synchronization job '%s' does not exist under service principal '%s'", jobId, servicePrincipalId))
		}

		mockState.deletes++
		mockState.schemas[key] = deepCopy(mockState.defaultSchema)
		return httpmock.NewStringResponse(204, ""), nil
	})
}

func (m *ServicePrincipalSynchronizationSchemaJsonMock) RegisterErrorMocks() {
	// Error scenarios for testing
	forbidden := `{"error":{"code":"Authorization_RequestDenied","message":"Insufficient privileges to complete the operation."}}`
	httpmock.RegisterRegexpResponder("GET", schemaURLPattern, httpmock.NewStringResponder(403, forbidden))
	badRequest := `{"error":{"code":"BadRequest","message":"The synchronization schema is invalid."}}`
	httpmock.RegisterRegexpResponder("PUT", schemaURLPattern, httpmock.NewStringResponder(400, badRequest))
	httpmock.RegisterRegexpResponder("DELETE", schemaURLPattern, httpmock.NewStringResponder(403, forbidden))
}

func (m *ServicePrincipalSynchronizationSchemaJsonMock) CleanupMockState() {
	mockState.Lock()
	mockState.defaultSchema = nil
	mockState.schemas = nil
	mockState.puts = nil
	mockState.deletes = 0
	mockState.revision = 0
	mockState.loadErr = nil
	mockState.Unlock()
}

// PutRequests returns the request bodies received by the PUT responder, oldest first.
func (m *ServicePrincipalSynchronizationSchemaJsonMock) PutRequests() []map[string]any {
	mockState.Lock()
	defer mockState.Unlock()
	return append([]map[string]any{}, mockState.puts...)
}

// DeleteCount returns the number of DELETE requests that reset a schema.
func (m *ServicePrincipalSynchronizationSchemaJsonMock) DeleteCount() int {
	mockState.Lock()
	defer mockState.Unlock()
	return mockState.deletes
}

// IsDefaultSchema reports whether the mock job schema equals the template default.
func (m *ServicePrincipalSynchronizationSchemaJsonMock) IsDefaultSchema() bool {
	mockState.Lock()
	defer mockState.Unlock()
	stored, ok := mockState.schemas[ServicePrincipalID+"/"+JobID]
	if !ok {
		return false
	}
	left, _ := json.Marshal(stored)
	right, _ := json.Marshal(mockState.defaultSchema)
	return string(left) == string(right)
}

// SetObjectMappingEnabled changes the enabled flag of an object mapping in the stored schema, simulating a change made
// outside Terraform.
func (m *ServicePrincipalSynchronizationSchemaJsonMock) SetObjectMappingEnabled(name string, enabled bool) {
	mockState.Lock()
	defer mockState.Unlock()
	stored := mockState.schemas[ServicePrincipalID+"/"+JobID]
	rules, _ := stored["synchronizationRules"].([]any)
	for _, item := range rules {
		rule, _ := item.(map[string]any)
		objectMappings, _ := rule["objectMappings"].([]any)
		for _, mappingItem := range objectMappings {
			mapping, _ := mappingItem.(map[string]any)
			if mapping["name"] == name {
				mapping["enabled"] = enabled
			}
		}
	}
	mockState.revision++
	applyServerManagedValues(stored, mockState.revision)
}

// RemoveJob deletes the mock synchronization job, simulating a job removed outside Terraform.
func (m *ServicePrincipalSynchronizationSchemaJsonMock) RemoveJob() {
	mockState.Lock()
	defer mockState.Unlock()
	delete(mockState.schemas, ServicePrincipalID+"/"+JobID)
}
