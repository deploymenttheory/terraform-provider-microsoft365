package mocks

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/google/uuid"
	"github.com/jarcoal/httpmock"
)

const BaseURL = "https://graph.microsoft.com/beta/networkaccess/tlsInspectionPolicies"

// TLSInspectionPolicyMock models independent policy and nested rule identities, including auto-created rules.
type TLSInspectionPolicyMock struct {
	sync.Mutex
	Policies map[string]map[string]any
	Rules    map[string]map[string]map[string]any
	Requests []RecordedRequest
}
type RecordedRequest struct {
	Method, Path string
	Body         map[string]any
}

var _ mocks.MockRegistrar = (*TLSInspectionPolicyMock)(nil)

func init() {
	mocks.GlobalRegistry.Register("network_tls_inspection_policy", &TLSInspectionPolicyMock{})
}
func (m *TLSInspectionPolicyMock) CleanupMockState() {
	m.Lock()
	defer m.Unlock()
	m.Policies = make(map[string]map[string]any)
	m.Rules = make(map[string]map[string]map[string]any)
	m.Requests = nil
}
func (m *TLSInspectionPolicyMock) RegisterMocks() {
	m.CleanupMockState()
	for _, method := range []string{"POST", "GET", "PATCH", "DELETE"} {
		httpmock.RegisterResponder(method, `=~^https://graph\.microsoft\.com/beta/networkaccess/tlsInspectionPolicies(?:/.*)?$`, m.respond)
	}
}
func (m *TLSInspectionPolicyMock) RegisterErrorMocks() {
	httpmock.RegisterResponder("POST", BaseURL, func(*http.Request) (*http.Response, error) { return apiError(400, "Invalid TLS inspection policy") })
}
func apiError(code int, message string) (*http.Response, error) {
	return httpmock.NewJsonResponse(code, map[string]any{"error": map[string]any{"code": fmt.Sprint(code), "message": message}})
}
func (m *TLSInspectionPolicyMock) respond(req *http.Request) (*http.Response, error) {
	m.Lock()
	defer m.Unlock()
	body := map[string]any{}
	if req.Method == "POST" || req.Method == "PATCH" {
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			return apiError(400, "Invalid JSON")
		}
	}
	m.Requests = append(m.Requests, RecordedRequest{req.Method, req.URL.Path, body})
	path := strings.TrimPrefix(req.URL.Path, "/beta/networkaccess/tlsInspectionPolicies")
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if path == "" {
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
			policy := map[string]any{"id": id, "description": nil, "version": "1.0.0", "lastModifiedDateTime": time.Now().UTC().Format(time.RFC3339Nano)}
			for k, v := range body {
				policy[k] = v
			}
			m.Policies[id] = policy
			systemID, recommendedID := uuid.NewString(), uuid.NewString()
			m.Rules[id] = map[string]map[string]any{
				systemID:      {"id": systemID, "name": "System Bypass TLS inspection rule", "priority": float64(50), "action": "bypass", "settings": map[string]any{"status": "enabled"}, "matchingConditions": nil},
				recommendedID: {"id": recommendedID, "name": "Recommended TLS inspection bypass categories rule", "priority": float64(65000), "action": "bypass", "settings": map[string]any{"status": "enabled"}, "matchingConditions": map[string]any{"destinations": []any{map[string]any{"@odata.type": "#microsoft.graph.networkaccess.tlsInspectionWebCategoryDestination", "values": []any{"Education", "Finance", "Government", "HealthAndMedicine"}}}}},
			}
			return httpmock.NewJsonResponse(201, policy)
		}
		return apiError(405, "unsupported collection request")
	}
	if _, err := uuid.Parse(parts[0]); err != nil {
		return apiError(400, "Invalid request parameters")
	}
	policy, exists := m.Policies[parts[0]]
	if !exists {
		return apiError(404, "Policy not found")
	}
	if len(parts) == 1 {
		switch req.Method {
		case "GET":
			return httpmock.NewJsonResponse(200, policy)
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
			return httpmock.NewStringResponse(204, ""), nil
		case "DELETE":
			delete(m.Policies, parts[0])
			delete(m.Rules, parts[0])
			return httpmock.NewStringResponse(204, ""), nil
		}
	}
	if len(parts) < 2 || parts[1] != "policyRules" {
		return apiError(404, "Unknown path")
	}
	rules := m.Rules[parts[0]]
	if req.Method == "POST" || req.Method == "PATCH" {
		if body["@odata.type"] != "#microsoft.graph.networkaccess.tlsInspectionRule" {
			return apiError(400, "Explicit TLS rule discriminator required")
		}
		if priority, ok := body["priority"].(float64); ok {
			if priority < 100 || priority > 2147483647 {
				return apiError(400, "Invalid priority")
			}
			for id, rule := range rules {
				if rule["priority"] == priority && (len(parts) < 3 || id != parts[2]) {
					return apiError(400, "A rule with this priority already exists")
				}
			}
		}
	}
	if len(parts) == 2 {
		if req.Method == "GET" {
			values := []any{}
			for _, v := range rules {
				values = append(values, v)
			}
			return httpmock.NewJsonResponse(200, map[string]any{"value": values})
		}
		if req.Method == "POST" {
			for _, key := range []string{"name", "action", "priority", "settings", "matchingConditions"} {
				if body[key] == nil {
					return apiError(400, key+" required")
				}
			}
			id := uuid.NewString()
			rule := map[string]any{"id": id, "description": nil}
			for k, v := range body {
				rule[k] = v
			}
			rules[id] = rule
			policy["lastModifiedDateTime"] = time.Now().UTC().Format(time.RFC3339Nano)
			response, err := httpmock.NewJsonResponse(201, rule)
			response.Header.Set("Location", BaseURL+"(%7BpolicyId%7D)/policyRules/"+id)
			return response, err
		}
	}
	if len(parts) == 3 {
		if _, err := uuid.Parse(parts[2]); err != nil {
			return apiError(400, "Invalid request parameters")
		}
		rule, exists := rules[parts[2]]
		if !exists {
			return apiError(404, "Rule not found")
		}
		switch req.Method {
		case "GET":
			return httpmock.NewJsonResponse(200, rule)
		case "PATCH":
			for k, v := range body {
				rule[k] = v
			}
			policy["lastModifiedDateTime"] = time.Now().UTC().Format(time.RFC3339Nano)
			return httpmock.NewStringResponse(204, ""), nil
		case "DELETE":
			delete(rules, parts[2])
			policy["lastModifiedDateTime"] = time.Now().UTC().Format(time.RFC3339Nano)
			return httpmock.NewStringResponse(204, ""), nil
		}
	}
	return apiError(405, "Unsupported request")
}
