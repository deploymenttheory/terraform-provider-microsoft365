package mocks

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/jarcoal/httpmock"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
)

const BaseURL = "https://graph.microsoft.com/beta/networkaccess/promptPolicies"

// PromptPolicyMock models independent policy and nested rule identities, without automatically creating rules.
type PromptPolicyMock struct {
	sync.Mutex
	Policies map[string]map[string]any
	Rules    map[string]map[string]map[string]any
	Requests []RecordedRequest
}
type RecordedRequest struct {
	Method, Path string
	Body         map[string]any
}

var _ mocks.MockRegistrar = (*PromptPolicyMock)(nil)

func init() {
	mocks.GlobalRegistry.Register("network_prompt_policy", &PromptPolicyMock{})
}

func (m *PromptPolicyMock) CleanupMockState() {
	m.Lock()
	defer m.Unlock()
	m.Policies = make(map[string]map[string]any)
	m.Rules = make(map[string]map[string]map[string]any)
	m.Requests = nil
}

func (m *PromptPolicyMock) RegisterMocks() {
	m.CleanupMockState()
	for _, method := range []string{"POST", "GET", "PATCH", "DELETE"} {
		httpmock.RegisterResponder(
			method,
			`=~^https://graph\.microsoft\.com/beta/networkaccess/promptPolicies(?:/.*)?$`,
			m.respond,
		)
	}
}

func (m *PromptPolicyMock) RegisterErrorMocks() {
	httpmock.RegisterResponder(
		"POST",
		BaseURL,
		func(*http.Request) (*http.Response, error) { return apiError(400, "Invalid prompt policy") },
	)
}

func apiError(code int, message string) (*http.Response, error) {
	return jsonResponse(code, map[string]any{"error": map[string]any{"code": fmt.Sprint(code), "message": message}})
}
func jsonResponse(code int, body any) (*http.Response, error) {
	response, err := httpmock.NewJsonResponse(code, body)
	if err != nil {
		return nil, fmt.Errorf("encode prompt mock response: %w", err)
	}
	return response, nil
}
func (m *PromptPolicyMock) respond(req *http.Request) (*http.Response, error) {
	m.Lock()
	defer m.Unlock()
	body := map[string]any{}
	if req.Method == "POST" || req.Method == "PATCH" {
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			return apiError(400, "Invalid JSON")
		}
	}
	m.Requests = append(m.Requests, RecordedRequest{req.Method, req.URL.Path, body})
	path := strings.TrimPrefix(req.URL.Path, "/beta/networkaccess/promptPolicies")
	if path == "" {
		return m.policyCollection(req, body)
	}
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if _, err := uuid.Parse(parts[0]); err != nil {
		return apiError(400, "Invalid request parameters")
	}
	policy, exists := m.Policies[parts[0]]
	if !exists {
		return apiError(404, "Policy not found")
	}
	if len(parts) == 1 {
		return m.policyItem(req, parts[0], policy, body)
	}
	if parts[1] != "policyRules" {
		return apiError(404, "Unknown path")
	}
	rules := m.Rules[parts[0]]
	ruleID := ""
	if len(parts) == 3 {
		ruleID = parts[2]
	}
	if req.Method == "POST" || req.Method == "PATCH" {
		if message := validateRuleFields(body); message != "" {
			return apiError(400, message)
		}
		if conditions, ok := body["matchingConditions"].(map[string]any); ok {
			if message := validateConditions(conditions); message != "" {
				return apiError(400, message)
			}
		}
		if message := validatePriority(body, rules, ruleID); message != "" {
			return apiError(400, message)
		}
	}
	if len(parts) == 2 {
		return ruleCollection(req, policy, rules, body)
	}
	if len(parts) == 3 {
		return ruleItem(req, ruleID, policy, rules, body)
	}
	return apiError(405, "Unsupported request")
}
func (m *PromptPolicyMock) policyCollection(req *http.Request, body map[string]any) (*http.Response, error) {
	if req.Method == "POST" {
		for _, key := range []string{"policyRules", "version"} {
			if _, ok := body[key]; ok {
				return apiError(400, key+" must not be sent by policy resource")
			}
		}
		if body["name"] == nil || body["settings"] == nil {
			return apiError(400, "Name and Settings required")
		}
		id := uuid.NewString()
		policy := map[string]any{
			"id":                   id,
			"description":          nil,
			"version":              "1.0.0",
			"lastModifiedDateTime": time.Now().UTC().Format(time.RFC3339Nano),
		}
		for k, v := range body {
			policy[k] = v
		}
		settings := policy["settings"].(map[string]any)
		if settings["defaultAction"] == nil {
			settings["defaultAction"] = "allow"
		}
		if settings["defaultAction"] != "allow" {
			return apiError(400, "Only allow is supported")
		}
		m.Policies[id] = policy
		m.Rules[id] = map[string]map[string]any{}
		return jsonResponse(200, policy)
	}
	return apiError(405, "unsupported collection request")
}
func (m *PromptPolicyMock) policyItem(req *http.Request, policyID string, policy, body map[string]any) (*http.Response, error) {
	switch req.Method {
	case "GET":
		return jsonResponse(200, policy)
	case "PATCH":
		for _, key := range []string{"policyRules", "version"} {
			if _, ok := body[key]; ok {
				return apiError(400, key+" must not be updated")
			}
		}
		for k, v := range body {
			policy[k] = v
		}
		policy["lastModifiedDateTime"] = time.Now().UTC().Format(time.RFC3339Nano)
		return jsonResponse(200, policy)
	case "DELETE":
		delete(m.Policies, policyID)
		delete(m.Rules, policyID)
		return httpmock.NewStringResponse(204, ""), nil
	}
	return apiError(405, "Unsupported request")
}
func validateRuleFields(body map[string]any) string {
	if body["@odata.type"] != "#microsoft.graph.networkaccess.promptRule" {
		return "Explicit prompt rule discriminator required"
	}
	if v, ok := body["action"]; ok && v != "allow" && v != "block" {
		return "Invalid action"
	}
	if v, ok := body["promptLogging"]; ok && v != "never" && v != "always" && v != "onBlock" {
		return "Invalid prompt logging"
	}
	if settings, ok := body["settings"].(map[string]any); ok {
		if v, ok := settings["status"]; ok && v != "enabled" && v != "disabled" {
			return "Invalid status"
		}
	}

	return ""
}
func validateConditions(conditions map[string]any) string {
	schemes, ok := conditions["conversationSchemes"].([]any)
	if !ok {
		return "ConversationSchemes required"
	}
	if v, ok := conditions["scanResult"]; ok && v != "maliciousPromptDetected" {
		return "Invalid scan result"
	}
	for _, value := range schemes {
		scheme, ok := value.(map[string]any)
		if !ok {
			return "Invalid conversation scheme"
		}
		switch scheme["@odata.type"] {
		case "#microsoft.graph.networkaccess.customConversationScheme":
			if scheme["url"] == nil || scheme["url"] == "" {
				return "Url required"
			}
			if scheme["jsonPath"] == nil {
				scheme["jsonPath"] = ""
			}
		case "#microsoft.graph.networkaccess.predefinedConversationScheme":
			if scheme["schemeName"] == nil {
				scheme["schemeName"] = "chatGpt"
			}
		default:
			return "Explicit scheme discriminator required"
		}
	}
	if conditions["scanResult"] == nil {
		conditions["scanResult"] = "maliciousPromptDetected"
	}

	return ""
}
func validatePriority(body map[string]any, rules map[string]map[string]any, ruleID string) string {
	if priority, ok := body["priority"].(float64); ok {
		if priority < 100 || priority > 2147483647 || math.Trunc(priority) != priority {
			return "Invalid priority"
		}
		for id, rule := range rules {
			if rule["priority"] == priority && id != ruleID {
				return "A rule with this priority already exists"
			}
		}
	}

	return ""
}
func ruleCollection(req *http.Request, policy map[string]any, rules map[string]map[string]any, body map[string]any) (*http.Response, error) {
	if req.Method == "GET" {
		values := []any{}
		for _, v := range rules {
			values = append(values, v)
		}
		return jsonResponse(200, map[string]any{"value": values})
	}
	if req.Method == "POST" {
		for _, key := range []string{"name", "priority", "settings", "matchingConditions"} {
			if body[key] == nil {
				return apiError(400, key+" required")
			}
		}
		id := uuid.NewString()
		rule := map[string]any{
			"id":            id,
			"description":   nil,
			"action":        "allow",
			"promptLogging": "never",
		}
		for k, v := range body {
			rule[k] = v
		}
		settings := rule["settings"].(map[string]any)
		if settings["status"] == nil {
			settings["status"] = "enabled"
		}
		rules[id] = rule
		policy["lastModifiedDateTime"] = time.Now().UTC().Format(time.RFC3339Nano)
		response, err := jsonResponse(200, rule)
		if err != nil {
			return nil, err
		}
		response.Header.Set("Location", "https://graph.microsoft.com")
		return response, err
	}
	return apiError(405, "Unsupported request")
}
func ruleItem(req *http.Request, ruleID string, policy map[string]any, rules map[string]map[string]any, body map[string]any) (*http.Response, error) {
	if _, err := uuid.Parse(ruleID); err != nil {
		return apiError(400, "Invalid request parameters")
	}
	rule, exists := rules[ruleID]
	if !exists {
		return apiError(404, "Rule not found")
	}
	switch req.Method {
	case "GET":
		return jsonResponse(200, rule)
	case "PATCH":
		for k, v := range body {
			rule[k] = v
		}
		policy["lastModifiedDateTime"] = time.Now().UTC().Format(time.RFC3339Nano)
		return jsonResponse(200, rule)
	case "DELETE":
		delete(rules, ruleID)
		policy["lastModifiedDateTime"] = time.Now().UTC().Format(time.RFC3339Nano)
		return httpmock.NewStringResponse(204, ""), nil
	}
	return apiError(405, "Unsupported request")
}
