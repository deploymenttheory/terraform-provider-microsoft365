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
	"strings"
	"sync"
	"time"

	"github.com/jarcoal/httpmock"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
)

const (
	// PolicyID is the ID of the role management policy in the mock tenant.
	PolicyID = "DirectoryRole_00000000-0000-0000-0000-000000000001_22222222-2222-2222-2222-222222222222"

	// RoleDefinitionID is the role definition the mock policy is assigned to. It differs from the policy ID suffix
	// as it does for some policies in real tenants, so the policy can only be found through its assignment.
	RoleDefinitionID = "11111111-1111-1111-1111-111111111111"
)

var (
	policyURLPattern      = regexp.MustCompile(`^https://graph\.microsoft\.com/beta/policies/roleManagementPolicies/([^/?]+)(\?.*)?$`)
	assignmentsURLPattern = regexp.MustCompile(`^https://graph\.microsoft\.com/beta/policies/roleManagementPolicyAssignments(\?.*)?$`)
	roleDefinitionFilter  = regexp.MustCompile(`roleDefinitionId eq '([^']*)'`)
	scopeIdFilter         = regexp.MustCompile(`scopeId eq '([^']*)'`)
	scopeTypeFilter       = regexp.MustCompile(`scopeType eq '([^']*)'`)
)

var mockState struct {
	sync.Mutex
	policies    map[string]map[string]any
	assignments []any
	patches     []map[string]any
	loadErr     error
}

func init() {
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))
	mocks.GlobalRegistry.Register("role_management_policy", &RoleManagementPolicyMock{})
}

// RoleManagementPolicyMock provides stateful httpmock responders for role management policies.
type RoleManagementPolicyMock struct{}

var _ mocks.MockRegistrar = (*RoleManagementPolicyMock)(nil)

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

// normaliseDuration mirrors Graph, which returns whole days as days and shorter durations in hours and minutes.
func normaliseDuration(value string) string {
	pattern := regexp.MustCompile(`^P(?:(\d+)D)?(?:T(?:(\d+)H)?(?:(\d+)M)?(?:(\d+)S)?)?$`)
	matches := pattern.FindStringSubmatch(value)
	if matches == nil {
		return value
	}
	units := []time.Duration{24 * time.Hour, time.Hour, time.Minute, time.Second}
	var total time.Duration
	for index, unit := range units {
		var count int
		if matches[index+1] != "" {
			_, _ = fmt.Sscanf(matches[index+1], "%d", &count)
		}
		total += time.Duration(count) * unit
	}
	if total == 0 {
		return "PT0S"
	}
	days := total / (24 * time.Hour)
	rest := total % (24 * time.Hour)
	result := "P"
	if days > 0 {
		result += fmt.Sprintf("%dD", days)
	}
	if rest > 0 {
		result += "T"
		if hours := rest / time.Hour; hours > 0 {
			result += fmt.Sprintf("%dH", hours)
		}
		if minutes := (rest % time.Hour) / time.Minute; minutes > 0 {
			result += fmt.Sprintf("%dM", minutes)
		}
		if seconds := (rest % time.Minute) / time.Second; seconds > 0 {
			result += fmt.Sprintf("%dS", seconds)
		}
	}
	return result
}

// applyRule merges one rule from a PATCH body into the stored policy, validating it the way Graph does.
func applyRule(policy map[string]any, update map[string]any) (int, string, string) {
	ruleId, _ := update["id"].(string)
	rules, _ := policy["rules"].([]any)
	for _, item := range rules {
		stored, _ := item.(map[string]any)
		if stored["id"] != ruleId {
			continue
		}
		if update["@odata.type"] != stored["@odata.type"] {
			return 400, "InvalidPolicyRuleProperty", fmt.Sprintf("The property value of '@odata.type' of rule id '%s' is invalid.", ruleId)
		}
		if stored["@odata.type"] == "#microsoft.graph.unifiedRoleManagementPolicyNotificationRule" {
			for _, field := range []string{"recipientType", "notificationType"} {
				if update[field] != stored[field] {
					return 400, "InvalidPolicyRuleProperty", fmt.Sprintf("The property value of '%s' of rule id '%s' is invalid.", field, ruleId)
				}
			}
		}
		if stored["@odata.type"] == "#microsoft.graph.unifiedRoleManagementPolicyAuthenticationContextRule" {
			if update["isEnabled"] == true && update["claimValue"] == nil {
				return 400, "InvalidPolicyRule", "The policy rule is invalid."
			}
			if update["isEnabled"] == false {
				update["claimValue"] = nil
			}
		}
		for key, value := range update {
			if key == "id" || key == "@odata.type" {
				continue
			}
			if key == "maximumDuration" {
				if text, ok := value.(string); ok {
					value = normaliseDuration(text)
				}
			}
			stored[key] = value
		}
		return 0, "", ""
	}
	return 400, "InvalidPolicyRuleId", fmt.Sprintf("The policy rule id %s is invalid.", ruleId)
}

func (m *RoleManagementPolicyMock) RegisterMocks() {
	policy, policyErr := loadJSONResponse("get_role_management_policy_success.json")
	assignments, assignmentsErr := loadJSONResponse("get_role_management_policy_assignments_success.json")

	mockState.Lock()
	mockState.policies = map[string]map[string]any{}
	if policy != nil {
		mockState.policies[PolicyID] = policy
	}
	mockState.assignments = nil
	if assignments != nil {
		mockState.assignments, _ = assignments["value"].([]any)
	}
	mockState.patches = nil
	mockState.loadErr = policyErr
	if mockState.loadErr == nil {
		mockState.loadErr = assignmentsErr
	}
	mockState.Unlock()

	// List policy assignments - GET /policies/roleManagementPolicyAssignments?$filter=...
	httpmock.RegisterRegexpResponder("GET", assignmentsURLPattern, func(req *http.Request) (*http.Response, error) {
		mockState.Lock()
		defer mockState.Unlock()

		if mockState.loadErr != nil {
			return errorResponse(500, "InternalServerError", mockState.loadErr.Error())
		}

		filter := req.URL.Query().Get("$filter")
		scopeId := scopeIdFilter.FindStringSubmatch(filter)
		scopeType := scopeTypeFilter.FindStringSubmatch(filter)
		if scopeId == nil || scopeType == nil {
			return errorResponse(400, "UnknownError", `{"errorCode":"MissingProvider","message":"The provider is missing."}`)
		}
		roleDefinition := roleDefinitionFilter.FindStringSubmatch(filter)

		values := make([]any, 0)
		for _, item := range mockState.assignments {
			assignment, _ := item.(map[string]any)
			if assignment["scopeId"] != scopeId[1] || assignment["scopeType"] != scopeType[1] {
				continue
			}
			if roleDefinition != nil && assignment["roleDefinitionId"] != roleDefinition[1] {
				continue
			}
			values = append(values, assignment)
		}
		return httpmock.NewJsonResponse(200, map[string]any{"value": values})
	})

	// Get role management policy - GET /policies/roleManagementPolicies/{id}?$expand=rules
	httpmock.RegisterRegexpResponder("GET", policyURLPattern, func(req *http.Request) (*http.Response, error) {
		mockState.Lock()
		defer mockState.Unlock()

		if mockState.loadErr != nil {
			return errorResponse(500, "InternalServerError", mockState.loadErr.Error())
		}

		policyId, _ := url.PathUnescape(policyURLPattern.FindStringSubmatch(req.URL.String())[1])
		stored, ok := mockState.policies[policyId]
		if !ok {
			return errorResponse(400, "ArgumentNullException", "Value cannot be null.\r\nParameter name: policyId")
		}

		response := map[string]any{}
		for key, value := range stored {
			response[key] = value
		}
		if !strings.Contains(req.URL.Query().Get("$expand"), "rules") {
			delete(response, "rules")
		}
		return httpmock.NewJsonResponse(200, response)
	})

	// Update role management policy rules - PATCH /policies/roleManagementPolicies/{id}
	httpmock.RegisterRegexpResponder("PATCH", policyURLPattern, func(req *http.Request) (*http.Response, error) {
		content, err := io.ReadAll(req.Body)
		if err != nil {
			return errorResponse(400, "BadRequest", "Invalid request body")
		}
		var requestBody, recorded map[string]any
		if err := json.Unmarshal(content, &requestBody); err != nil {
			return errorResponse(400, "BadRequest", "Invalid request body")
		}
		_ = json.Unmarshal(content, &recorded)

		mockState.Lock()
		defer mockState.Unlock()

		policyId, _ := url.PathUnescape(policyURLPattern.FindStringSubmatch(req.URL.String())[1])
		stored, ok := mockState.policies[policyId]
		if !ok {
			return errorResponse(400, "ArgumentNullException", "Value cannot be null.\r\nParameter name: policyId")
		}

		mockState.patches = append(mockState.patches, recorded)

		// Validate every rule before applying any of them, as Graph rejects the whole request.
		var staged map[string]any
		stagedContent, _ := json.Marshal(stored)
		_ = json.Unmarshal(stagedContent, &staged)
		updates, _ := requestBody["rules"].([]any)
		for _, item := range updates {
			update, _ := item.(map[string]any)
			if status, code, message := applyRule(staged, update); status != 0 {
				return errorResponse(status, code, message)
			}
		}

		staged["lastModifiedDateTime"] = time.Now().UTC().Format(time.RFC3339)
		mockState.policies[policyId] = staged

		response := map[string]any{}
		for key, value := range staged {
			if key != "rules" {
				response[key] = value
			}
		}
		return httpmock.NewJsonResponse(200, response)
	})
}

func (m *RoleManagementPolicyMock) RegisterErrorMocks() {
	// Error scenarios for testing
	forbidden := `{"error":{"code":"Authorization_RequestDenied","message":"Insufficient privileges to complete the operation."}}`
	httpmock.RegisterRegexpResponder("GET", assignmentsURLPattern, httpmock.NewStringResponder(403, forbidden))
	httpmock.RegisterRegexpResponder("GET", policyURLPattern, httpmock.NewStringResponder(403, forbidden))
	httpmock.RegisterRegexpResponder("PATCH", policyURLPattern, httpmock.NewStringResponder(400, `{"error":{"code":"InvalidPolicy","message":"The policy is invalid."}}`))
}

func (m *RoleManagementPolicyMock) CleanupMockState() {
	mockState.Lock()
	mockState.policies = nil
	mockState.assignments = nil
	mockState.patches = nil
	mockState.loadErr = nil
	mockState.Unlock()
}

// PatchRequests returns the request bodies received by the PATCH responder, oldest first.
func (m *RoleManagementPolicyMock) PatchRequests() []map[string]any {
	mockState.Lock()
	defer mockState.Unlock()
	return append([]map[string]any{}, mockState.patches...)
}

// SetRule replaces fields of a rule in the stored policy, simulating a change made outside Terraform.
func (m *RoleManagementPolicyMock) SetRule(ruleId string, fields map[string]any) {
	mockState.Lock()
	defer mockState.Unlock()
	policy := mockState.policies[PolicyID]
	rules, _ := policy["rules"].([]any)
	for _, item := range rules {
		rule, _ := item.(map[string]any)
		if rule["id"] == ruleId {
			for key, value := range fields {
				rule[key] = value
			}
		}
	}
}
