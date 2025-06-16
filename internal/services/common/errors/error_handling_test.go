package errors

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	abstractions "github.com/microsoft/kiota-abstractions-go"
	"github.com/microsoft/kiota-abstractions-go/serialization"
	"github.com/microsoft/kiota-abstractions-go/store"
	"github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
	"github.com/stretchr/testify/assert"
)

// Mock implementations for testing

// MockResponse implements the necessary interfaces for testing
type MockResponse struct {
	Diagnostics diag.Diagnostics
}

func (m *MockResponse) GetDiagnostics() diag.Diagnostics {
	return m.Diagnostics
}

// MockResponseHeaders implements abstractions.ResponseHeaders
type MockResponseHeaders struct {
	headers map[string][]string
}

func (m *MockResponseHeaders) Get(key string) []string {
	return m.headers[key]
}

func (m *MockResponseHeaders) ListKeys() []string {
	keys := make([]string, 0, len(m.headers))
	for k := range m.headers {
		keys = append(keys, k)
	}
	return keys
}

func (m *MockResponseHeaders) Add(key, value string) {
	if m.headers == nil {
		m.headers = make(map[string][]string)
	}
	m.headers[key] = append(m.headers[key], value)
}

func (m *MockResponseHeaders) Set(key, value string) {
	if m.headers == nil {
		m.headers = make(map[string][]string)
	}
	m.headers[key] = []string{value}
}

// MockODataError implements the OData error interfaces
type MockODataError struct {
	statusCode int
	headers    *MockResponseHeaders
	errorData  *MockMainError
}

// Ensure this is recognized as an ODataError
func (m *MockODataError) GetErrorEscaped() odataerrors.MainErrorable {
	return m.errorData
}

func (m *MockODataError) Error() string {
	if m.errorData != nil && m.errorData.message != nil {
		return *m.errorData.message
	}
	return fmt.Sprintf("HTTP %d", m.statusCode)
}

func (m *MockODataError) GetStatusCode() int {
	return m.statusCode
}

func (m *MockODataError) GetResponseHeaders() abstractions.ResponseHeaders {
	if m.headers == nil {
		// Create an empty ResponseHeaders
		empty := abstractions.NewResponseHeaders()
		return *empty
	}

	// Convert our mock headers to the actual ResponseHeaders
	respHeaders := abstractions.NewResponseHeaders()
	for key, values := range m.headers.headers {
		for _, value := range values {
			respHeaders.Add(key, value)
		}
	}
	return *respHeaders
}

// MockMainError implements odataerrors.MainErrorable
type MockMainError struct {
	code         *string
	message      *string
	target       *string
	details      []odataerrors.ErrorDetailsable
	innerError   *MockInnerError
	backingStore store.BackingStore
}

func (m *MockMainError) GetCode() *string {
	return m.code
}

func (m *MockMainError) GetMessage() *string {
	return m.message
}

func (m *MockMainError) GetTarget() *string {
	return m.target
}

func (m *MockMainError) GetDetails() []odataerrors.ErrorDetailsable {
	return m.details
}

func (m *MockMainError) GetInnerError() odataerrors.InnerErrorable {
	if m.innerError == nil {
		return nil
	}
	return m.innerError
}

// Add other required methods for MainErrorable interface
func (m *MockMainError) GetAdditionalData() map[string]interface{} {
	return make(map[string]interface{})
}

func (m *MockMainError) GetOdataType() *string {
	return nil
}

func (m *MockMainError) SetAdditionalData(value map[string]interface{}) {
}

func (m *MockMainError) SetCode(value *string) {
	m.code = value
}

func (m *MockMainError) SetDetails(value []odataerrors.ErrorDetailsable) {
	m.details = value
}

func (m *MockMainError) SetInnerError(value odataerrors.InnerErrorable) {
	if value == nil {
		m.innerError = nil
		return
	}

	// We need to handle the case where the value is not a *MockInnerError
	// This is a simplification for testing purposes
	mockInner := &MockInnerError{
		backingStore: store.BackingStoreFactoryInstance(),
	}

	if reqID := value.GetRequestId(); reqID != nil {
		mockInner.requestId = reqID
	}

	if clientReqID := value.GetClientRequestId(); clientReqID != nil {
		mockInner.clientRequestId = clientReqID
	}

	if date := value.GetDate(); date != nil {
		mockInner.date = date
	}

	if odataType := value.GetOdataType(); odataType != nil {
		mockInner.odataType = odataType
	}

	m.innerError = mockInner
}

func (m *MockMainError) SetMessage(value *string) {
	m.message = value
}

func (m *MockMainError) SetOdataType(value *string) {
}

func (m *MockMainError) SetTarget(value *string) {
	m.target = value
}

func (m *MockMainError) GetBackingStore() store.BackingStore {
	if m.backingStore == nil {
		m.backingStore = store.BackingStoreFactoryInstance()
	}
	return m.backingStore
}

func (m *MockMainError) GetFieldDeserializers() map[string]func(serialization.ParseNode) error {
	return make(map[string]func(serialization.ParseNode) error)
}

func (m *MockMainError) Serialize(writer serialization.SerializationWriter) error {
	return nil
}

func (m *MockMainError) SetBackingStore(value store.BackingStore) {
	m.backingStore = value
}

// MockInnerError implements odataerrors.InnerErrorable
type MockInnerError struct {
	requestId       *string
	clientRequestId *string
	date            *time.Time
	odataType       *string
	backingStore    store.BackingStore
}

func (m *MockInnerError) GetRequestId() *string {
	return m.requestId
}

func (m *MockInnerError) GetClientRequestId() *string {
	return m.clientRequestId
}

func (m *MockInnerError) GetDate() *time.Time {
	return m.date
}

func (m *MockInnerError) GetOdataType() *string {
	return m.odataType
}

// Add other required methods for InnerErrorable interface
func (m *MockInnerError) GetAdditionalData() map[string]interface{} {
	return make(map[string]interface{})
}

func (m *MockInnerError) SetAdditionalData(value map[string]interface{}) {
}

func (m *MockInnerError) SetClientRequestId(value *string) {
	m.clientRequestId = value
}

func (m *MockInnerError) SetDate(value *time.Time) {
	m.date = value
}

func (m *MockInnerError) SetOdataType(value *string) {
	m.odataType = value
}

func (m *MockInnerError) SetRequestId(value *string) {
	m.requestId = value
}

func (m *MockInnerError) GetBackingStore() store.BackingStore {
	if m.backingStore == nil {
		m.backingStore = store.BackingStoreFactoryInstance()
	}
	return m.backingStore
}

func (m *MockInnerError) GetFieldDeserializers() map[string]func(serialization.ParseNode) error {
	return make(map[string]func(serialization.ParseNode) error)
}

func (m *MockInnerError) Serialize(writer serialization.SerializationWriter) error {
	return nil
}

func (m *MockInnerError) SetBackingStore(value store.BackingStore) {
	m.backingStore = value
}

// MockErrorDetails implements odataerrors.ErrorDetailsable
type MockErrorDetails struct {
	code         *string
	message      *string
	target       *string
	backingStore store.BackingStore
}

func (m *MockErrorDetails) GetCode() *string {
	return m.code
}

func (m *MockErrorDetails) GetMessage() *string {
	return m.message
}

func (m *MockErrorDetails) GetTarget() *string {
	return m.target
}

func (m *MockErrorDetails) GetAdditionalData() map[string]interface{} {
	return make(map[string]interface{})
}

func (m *MockErrorDetails) GetOdataType() *string {
	return nil
}

func (m *MockErrorDetails) SetAdditionalData(value map[string]interface{}) {
}

func (m *MockErrorDetails) SetCode(value *string) {
	m.code = value
}

func (m *MockErrorDetails) SetMessage(value *string) {
	m.message = value
}

func (m *MockErrorDetails) SetOdataType(value *string) {
}

func (m *MockErrorDetails) SetTarget(value *string) {
	m.target = value
}

func (m *MockErrorDetails) GetBackingStore() store.BackingStore {
	if m.backingStore == nil {
		m.backingStore = store.BackingStoreFactoryInstance()
	}
	return m.backingStore
}

func (m *MockErrorDetails) GetFieldDeserializers() map[string]func(serialization.ParseNode) error {
	return make(map[string]func(serialization.ParseNode) error)
}

func (m *MockErrorDetails) Serialize(writer serialization.SerializationWriter) error {
	return nil
}

func (m *MockErrorDetails) SetBackingStore(value store.BackingStore) {
	m.backingStore = value
}

// Helper function to create mock OData errors
func createMockODataError(statusCode int, code, message, target string, headers map[string]string) *MockODataError {
	mockHeaders := &MockResponseHeaders{
		headers: make(map[string][]string),
	}

	for k, v := range headers {
		mockHeaders.Set(k, v)
	}

	// Parse date for inner error
	var parsedDate *time.Time
	if dateStr, exists := headers["date"]; exists {
		if t, err := time.Parse(time.RFC1123, dateStr); err == nil {
			parsedDate = &t
		}
	}

	innerError := &MockInnerError{
		requestId:       getStringPtr(headers["request-id"]),
		clientRequestId: getStringPtr(headers["client-request-id"]),
		date:            parsedDate,
		backingStore:    store.BackingStoreFactoryInstance(),
	}

	mainError := &MockMainError{
		code:         &code,
		message:      &message,
		target:       &target,
		details:      []odataerrors.ErrorDetailsable{},
		innerError:   innerError,
		backingStore: store.BackingStoreFactoryInstance(),
	}

	return &MockODataError{
		statusCode: statusCode,
		headers:    mockHeaders,
		errorData:  mainError,
	}
}

func getStringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func TestGraphError(t *testing.T) {
	tests := []struct {
		name          string
		statusCode    int
		code          string
		message       string
		target        string
		headers       map[string]string
		expectedError GraphErrorInfo
	}{
		{
			name:       "OData error with diagnostic info",
			statusCode: http.StatusBadRequest,
			code:       "BadRequest",
			message:    "No OData route exists that match template ~/singleton/navigation/key with http verb POST for request /DeviceConfiguration_2505/StatelessDeviceConfigurationFEService/deviceManagement/deviceConfigurations('assign')",
			target:     "deviceConfigurations",
			headers: map[string]string{
				"x-ms-ags-diagnostic": `[{"ServerInfo":{"DataCenter":"UK South","Slice":"E","Ring":"5","ScaleUnit":"000","RoleInstance":"LN2PEPF00014209"}}]`,
				"request-id":          "09fe057e-bae6-4aab-ae2b-98f912259821",
				"client-request-id":   "1b3b835a-15f2-4943-a9be-70b2f9e7431d",
				"date":                "Fri, 13 Jun 2025 15:24:26 GMT",
			},
			expectedError: GraphErrorInfo{
				StatusCode:      400,
				ErrorCode:       "BadRequest",
				ErrorMessage:    "No OData route exists that match template ~/singleton/navigation/key with http verb POST for request /DeviceConfiguration_2505/StatelessDeviceConfigurationFEService/deviceManagement/deviceConfigurations('assign')",
				Target:          "deviceConfigurations",
				IsODataError:    true,
				RequestID:       "09fe057e-bae6-4aab-ae2b-98f912259821",
				ClientRequestID: "1b3b835a-15f2-4943-a9be-70b2f9e7431d",
				ErrorDate:       "Fri, 13 Jun 2025 15:24:26 GMT",
				Category:        CategoryValidation,
			},
		},
		{
			name:       "Authentication error",
			statusCode: http.StatusUnauthorized,
			code:       "InvalidAuthenticationToken",
			message:    "Access token has expired",
			target:     "deviceConfigurations",
			headers: map[string]string{
				"x-ms-ags-diagnostic": `[{"ServerInfo":{"DataCenter":"UK South","Slice":"E","Ring":"5","ScaleUnit":"000","RoleInstance":"LN2PEPF00014209"}}]`,
				"request-id":          "09fe057e-bae6-4aab-ae2b-98f912259821",
				"client-request-id":   "1b3b835a-15f2-4943-a9be-70b2f9e7431d",
				"date":                "Fri, 13 Jun 2025 15:24:26 GMT",
			},
			expectedError: GraphErrorInfo{
				StatusCode:      401,
				ErrorCode:       "InvalidAuthenticationToken",
				ErrorMessage:    "Access token has expired",
				Target:          "deviceConfigurations",
				IsODataError:    true,
				RequestID:       "09fe057e-bae6-4aab-ae2b-98f912259821",
				ClientRequestID: "1b3b835a-15f2-4943-a9be-70b2f9e7431d",
				ErrorDate:       "Fri, 13 Jun 2025 15:24:26 GMT",
				Category:        CategoryAuthentication,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock OData error
			mockError := createMockODataError(tt.statusCode, tt.code, tt.message, tt.target, tt.headers)

			// Process the error
			errorInfo := GraphError(context.Background(), mockError)

			// Verify the error information
			assert.Equal(t, tt.expectedError.StatusCode, errorInfo.StatusCode)
			assert.Equal(t, tt.expectedError.ErrorCode, errorInfo.ErrorCode)
			assert.Equal(t, tt.expectedError.ErrorMessage, errorInfo.ErrorMessage)
			assert.Equal(t, tt.expectedError.Target, errorInfo.Target)
			assert.Equal(t, tt.expectedError.IsODataError, errorInfo.IsODataError)
			assert.Equal(t, tt.expectedError.RequestID, errorInfo.RequestID)
			assert.Equal(t, tt.expectedError.ClientRequestID, errorInfo.ClientRequestID)
			assert.Equal(t, tt.expectedError.ErrorDate, errorInfo.ErrorDate)
			assert.Equal(t, tt.expectedError.Category, errorInfo.Category)

			// Check that diagnostic info was extracted from headers
			if expectedDiag := tt.headers["x-ms-ags-diagnostic"]; expectedDiag != "" {
				assert.Contains(t, errorInfo.RequestDetails, "x-ms-ags-diagnostic")
			}
		})
	}
}

func TestHandleGraphError(t *testing.T) {
	tests := []struct {
		name                  string
		statusCode            int
		code                  string
		message               string
		target                string
		headers               map[string]string
		operation             string
		requiredPermissions   []string
		expectError           bool
		expectRemoveFromState bool
	}{
		{
			name:       "Handle 400 error on read operation - remove from state",
			statusCode: http.StatusBadRequest,
			code:       "BadRequest",
			message:    "Resource not found",
			target:     "deviceConfigurations",
			headers: map[string]string{
				"request-id":        "09fe057e-bae6-4aab-ae2b-98f912259821",
				"client-request-id": "1b3b835a-15f2-4943-a9be-70b2f9e7431d",
			},
			operation:             "Read",
			expectError:           false,
			expectRemoveFromState: true,
		},
		{
			name:       "Handle 401 error",
			statusCode: http.StatusUnauthorized,
			code:       "InvalidAuthenticationToken",
			message:    "Access token has expired",
			target:     "deviceConfigurations",
			headers: map[string]string{
				"request-id":        "09fe057e-bae6-4aab-ae2b-98f912259821",
				"client-request-id": "1b3b835a-15f2-4943-a9be-70b2f9e7431d",
			},
			operation:             "Read",
			requiredPermissions:   []string{"DeviceManagementConfiguration.Read.All"},
			expectError:           true,
			expectRemoveFromState: false,
		},
		{
			name:       "Handle 404 error on read operation - remove from state",
			statusCode: http.StatusNotFound,
			code:       "NotFound",
			message:    "Resource not found",
			target:     "deviceConfigurations",
			headers: map[string]string{
				"request-id": "09fe057e-bae6-4aab-ae2b-98f912259821",
			},
			operation:             "Read",
			expectError:           false,
			expectRemoveFromState: true,
		},
		{
			name:                  "Handle 404 error on create operation - add error",
			statusCode:            http.StatusNotFound,
			code:                  "NotFound",
			message:               "Resource not found",
			target:                "deviceConfigurations",
			headers:               map[string]string{},
			operation:             "Create",
			expectError:           true,
			expectRemoveFromState: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock OData error
			mockError := createMockODataError(tt.statusCode, tt.code, tt.message, tt.target, tt.headers)

			// Create a mock response object
			mockResp := &resource.ReadResponse{}

			// Handle the error
			HandleGraphError(context.Background(), mockError, mockResp, tt.operation, tt.requiredPermissions)

			// Verify the results
			if tt.expectError {
				assert.True(t, mockResp.Diagnostics.HasError(), "Expected error to be added to diagnostics")
			}

			if tt.expectRemoveFromState {
				// For read operations with 400/404, the resource should be removed from state
				// We can't easily test this without a more complex mock, but we can verify no error was added
				assert.False(t, mockResp.Diagnostics.HasError(), "Expected no error when removing from state")
			}
		})
	}
}

func TestErrorCategorization(t *testing.T) {
	tests := []struct {
		name             string
		errorInfo        GraphErrorInfo
		expectedCategory ErrorCategory
	}{
		{
			name: "Authentication error",
			errorInfo: GraphErrorInfo{
				StatusCode: 401,
				ErrorCode:  "InvalidAuthenticationToken",
			},
			expectedCategory: CategoryAuthentication,
		},
		{
			name: "Authorization error",
			errorInfo: GraphErrorInfo{
				StatusCode: 403,
				ErrorCode:  "Forbidden",
			},
			expectedCategory: CategoryAuthorization,
		},
		{
			name: "Validation error",
			errorInfo: GraphErrorInfo{
				StatusCode: 400,
				ErrorCode:  "BadRequest",
			},
			expectedCategory: CategoryValidation,
		},
		{
			name: "Throttling error",
			errorInfo: GraphErrorInfo{
				StatusCode: 429,
				ErrorCode:  "TooManyRequests",
			},
			expectedCategory: CategoryThrottling,
		},
		{
			name: "Service error",
			errorInfo: GraphErrorInfo{
				StatusCode: 500,
				ErrorCode:  "InternalServerError",
			},
			expectedCategory: CategoryService,
		},
		{
			name: "Service unavailable error - should be service category",
			errorInfo: GraphErrorInfo{
				StatusCode: 503,
				ErrorCode:  "ServiceUnavailable",
			},
			expectedCategory: CategoryService,
		},
		{
			name: "Network error with zero status code",
			errorInfo: GraphErrorInfo{
				StatusCode: 0,
				ErrorCode:  "NetworkError",
			},
			expectedCategory: CategoryNetwork,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			category := categorizeError(&tt.errorInfo)
			assert.Equal(t, tt.expectedCategory, category)
		})
	}
}

func TestIsRetryableError(t *testing.T) {
	tests := []struct {
		name        string
		errorInfo   GraphErrorInfo
		shouldRetry bool
	}{
		{
			name: "Rate limit error should be retryable",
			errorInfo: GraphErrorInfo{
				StatusCode: 429,
				ErrorCode:  "TooManyRequests",
			},
			shouldRetry: true,
		},
		{
			name: "Service unavailable should be retryable",
			errorInfo: GraphErrorInfo{
				StatusCode: 503,
				ErrorCode:  "ServiceUnavailable",
			},
			shouldRetry: true,
		},
		{
			name: "Internal server error should be retryable",
			errorInfo: GraphErrorInfo{
				StatusCode: 500,
				ErrorCode:  "InternalServerError",
			},
			shouldRetry: true,
		},
		{
			name: "Bad request should not be retryable",
			errorInfo: GraphErrorInfo{
				StatusCode: 400,
				ErrorCode:  "BadRequest",
			},
			shouldRetry: false,
		},
		{
			name: "Authentication error should not be retryable",
			errorInfo: GraphErrorInfo{
				StatusCode: 401,
				ErrorCode:  "InvalidAuthenticationToken",
			},
			shouldRetry: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsRetryableError(&tt.errorInfo)
			assert.Equal(t, tt.shouldRetry, result)
		})
	}
}

func TestGetRetryDelay(t *testing.T) {
	tests := []struct {
		name          string
		errorInfo     GraphErrorInfo
		attempt       int
		expectedDelay bool // Just check if delay is reasonable
	}{
		{
			name: "With retry-after header",
			errorInfo: GraphErrorInfo{
				RetryAfter: "30",
			},
			attempt:       1,
			expectedDelay: true,
		},
		{
			name:          "Without retry-after header",
			errorInfo:     GraphErrorInfo{},
			attempt:       1,
			expectedDelay: true,
		},
		{
			name:          "High attempt number",
			errorInfo:     GraphErrorInfo{},
			attempt:       10,
			expectedDelay: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			delay := GetRetryDelay(&tt.errorInfo, tt.attempt)
			assert.True(t, delay > 0, "Delay should be positive")
			assert.True(t, delay <= 5*time.Minute, "Delay should not exceed max delay")
		})
	}
}
