package mocks

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/jarcoal/httpmock"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks/factories"
)

const BaseURL = "https://graph.microsoft.com/beta/networkAccess/threatIntelligencePolicies"

// ThreatIntelligencePolicyMock models independent policy and nested rule identities, including auto-created rules.
type ThreatIntelligencePolicyMock struct {
	sync.Mutex
	Policies map[string]map[string]any
	Rules    map[string]map[string]map[string]any
	Requests []RecordedRequest
}
type RecordedRequest struct {
	Method, Path string
	Body         map[string]any
}

var _ mocks.MockRegistrar = (*ThreatIntelligencePolicyMock)(nil)

func init() {
	mocks.GlobalRegistry.Register(
		"network_threat_intelligence_policy",
		&ThreatIntelligencePolicyMock{},
	)
}

func (m *ThreatIntelligencePolicyMock) CleanupMockState() {
	m.Lock()
	defer m.Unlock()
	m.Policies = make(map[string]map[string]any)
	m.Rules = make(map[string]map[string]map[string]any)
	m.Requests = nil
}

func (m *ThreatIntelligencePolicyMock) RegisterMocks() {
	m.CleanupMockState()
	for _, method := range []string{"POST", "GET", "PATCH", "DELETE"} {
		httpmock.RegisterResponder(
			method,
			`=~^https://graph\.microsoft\.com/beta/networkAccess/threatIntelligencePolicies(?:/.*)?$`,
			m.respond,
		)
	}
}

func (m *ThreatIntelligencePolicyMock) RegisterErrorMocks() {
	httpmock.RegisterResponder("POST", BaseURL, func(*http.Request) (*http.Response, error) {
		return apiError(400, "Invalid threat intelligence policy")
	})
}

func apiError(code int, message string) (*http.Response, error) {
	return factories.ErrorResponse(code, fmt.Sprint(code), message)(nil)
}

func (m *ThreatIntelligencePolicyMock) respond(req *http.Request) (*http.Response, error) {
	m.Lock()
	defer m.Unlock()
	body := map[string]any{}
	if req.Method == "POST" || req.Method == "PATCH" {
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			return apiError(400, "Invalid JSON")
		}
	}
	m.Requests = append(m.Requests, RecordedRequest{req.Method, req.URL.Path, body})
	path := strings.TrimPrefix(req.URL.Path, "/beta/networkAccess/threatIntelligencePolicies")
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if path == "" {
		return m.respondPolicyCollection(req, body)
	}
	if _, err := uuid.Parse(parts[0]); err != nil {
		return apiError(400, "Invalid request parameters")
	}
	policy, exists := m.Policies[parts[0]]
	if !exists {
		return apiError(404, "Policy not found")
	}
	if len(parts) == 1 {
		return m.respondPolicy(req, parts[0], policy, body)
	}
	if len(parts) < 2 || parts[1] != "policyRules" {
		return apiError(404, "Unknown path")
	}
	rules := m.Rules[parts[0]]
	if req.Method == "POST" || req.Method == "PATCH" {
		if message := validateRuleWrite(parts, rules, body); message != "" {
			return apiError(400, message)
		}
	}
	if len(parts) == 2 {
		return m.respondRuleCollection(req, policy, rules, body)
	}
	if len(parts) == 3 {
		return m.respondRule(req, parts[2], policy, rules, body)
	}
	return apiError(405, "Unsupported request")
}

func (m *ThreatIntelligencePolicyMock) respondPolicyCollection(
	req *http.Request,
	body map[string]any,
) (*http.Response, error) {
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
		m.Policies[id] = policy
		defaultID := uuid.NewString()
		m.Rules[id] = map[string]map[string]any{defaultID: {
			"@odata.type": "#microsoft.graph.networkaccess.threatIntelligenceRule",
			"id":          defaultID,
			"name":        "Default threat intel rule",
			"description": "Auto-created rule blocking access to sites with high severity threat detected",
			"priority":    float64(65000),
			"action":      "block",
			"settings":    map[string]any{"status": "enabled"},
			"matchingConditions": map[string]any{
				"severity": "high",
				"sources":  nil,
				"destinations": []any{
					map[string]any{
						"@odata.type":       "#microsoft.graph.networkaccess.threatIntelligenceFqdnDestination",
						"values":            []any{"*"},
						"httpRequestMethod": nil,
					},
				},
			},
		}}
		return factories.SuccessResponse(201, policy)(req)
	}
	return apiError(405, "unsupported collection request")
}

func (m *ThreatIntelligencePolicyMock) respondPolicy(
	req *http.Request,
	id string,
	policy, body map[string]any,
) (*http.Response, error) {
	switch req.Method {
	case "GET":
		return factories.SuccessResponse(200, policy)(req)
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
		return factories.EmptySuccessResponse(204)(req)
	case "DELETE":
		delete(m.Policies, id)
		delete(m.Rules, id)
		return factories.EmptySuccessResponse(204)(req)
	}

	return apiError(404, "Unknown path")
}

func validateRuleWrite(
	parts []string,
	rules map[string]map[string]any,
	body map[string]any,
) string {
	if body["@odata.type"] != "#microsoft.graph.networkaccess.threatIntelligenceRule" {
		return "Explicit threat intelligence rule discriminator required"
	}
	if conditions, ok := body["matchingConditions"].(map[string]any); ok {
		destinations, ok := conditions["destinations"].([]any)
		if !ok || len(destinations) > 1 {
			return "Destinations required; each type may occur once"
		}
		if severity, exists := conditions["severity"]; exists && severity != "high" {
			return "Only high severity supported"
		}
		for _, item := range destinations {
			destination, ok := item.(map[string]any)
			if !ok {
				return "Invalid destination"
			}
			values, ok := destination["values"].([]any)
			if !ok || len(values) == 0 {
				return "At least one value required"
			}
		}
	}
	for id, rule := range rules {
		if body["name"] != nil && body["name"] == rule["name"] &&
			(len(parts) < 3 || id != parts[2]) {
			return "A rule with this name already exists"
		}
	}
	if priority, ok := body["priority"].(float64); ok {
		if priority < 100 || priority > 2147483647 {
			return "Invalid priority"
		}
		for id, rule := range rules {
			if rule["priority"] == priority && (len(parts) < 3 || id != parts[2]) {
				return "A rule with this priority already exists"
			}
		}
	}

	return ""
}

func (m *ThreatIntelligencePolicyMock) respondRuleCollection(
	req *http.Request,
	policy map[string]any,
	rules map[string]map[string]any,
	body map[string]any,
) (*http.Response, error) {
	if req.Method == "GET" {
		values := []any{}
		for _, v := range rules {
			values = append(values, v)
		}
		return factories.SuccessResponse(200, map[string]any{"value": values})(req)
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
		response, err := factories.SuccessResponse(201, rule)(req)
		response.Header.Set("Location", BaseURL+"(%7BpolicyId%7D)/policyRules/"+id)
		return response, err
	}

	return apiError(405, "Unsupported request")
}

func (m *ThreatIntelligencePolicyMock) respondRule(
	req *http.Request,
	id string,
	policy map[string]any,
	rules map[string]map[string]any,
	body map[string]any,
) (*http.Response, error) {
	if _, err := uuid.Parse(id); err != nil {
		return apiError(400, "Invalid request parameters")
	}
	rule, exists := rules[id]
	if !exists {
		return apiError(404, "Rule not found")
	}
	switch req.Method {
	case "GET":
		return factories.SuccessResponse(200, rule)(req)
	case "PATCH":
		if rule["priority"] == float64(65000) {
			return apiError(400, "Default rule cannot be changed")
		}
		for k, v := range body {
			rule[k] = v
		}
		policy["lastModifiedDateTime"] = time.Now().UTC().Format(time.RFC3339Nano)
		return factories.EmptySuccessResponse(204)(req)
	case "DELETE":
		if rule["priority"] == float64(65000) {
			return apiError(400, "Deleting default rule is not permitted")
		}
		delete(rules, id)
		policy["lastModifiedDateTime"] = time.Now().UTC().Format(time.RFC3339Nano)
		return factories.EmptySuccessResponse(204)(req)
	}

	return apiError(405, "Unsupported request")
}
