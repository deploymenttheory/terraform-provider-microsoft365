package mocks

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/jarcoal/httpmock"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
)

const BaseURL = "https://graph.microsoft.com/beta/networkaccess/mcpPolicies"

// MCPPolicyMock maintains independent parent/child identities and Graph's nested PATCH merge semantics.
type MCPPolicyMock struct {
	sync.Mutex
	Policies map[string]map[string]any
	Rules    map[string]map[string]map[string]any
	Requests []RecordedRequest
}

// RecordedRequest records a decoded request without authentication data.
type RecordedRequest struct {
	Method, Path string
	Body         map[string]any
}

var _ mocks.MockRegistrar = (*MCPPolicyMock)(nil)

func init() { mocks.GlobalRegistry.Register("network_mcp_policy", &MCPPolicyMock{}) }
func (m *MCPPolicyMock) CleanupMockState() {
	m.Lock()
	defer m.Unlock()
	m.Policies = map[string]map[string]any{}
	m.Rules = map[string]map[string]map[string]any{}
	m.Requests = nil
}

func (m *MCPPolicyMock) RegisterMocks() {
	m.CleanupMockState()
	for _, method := range []string{"POST", "GET", "PATCH", "DELETE"} {
		httpmock.RegisterResponder(
			method,
			`=~^https://graph\.microsoft\.com/beta/networkaccess/mcpPolicies(?:/.*)?$`,
			m.respond,
		)
	}
}

func (m *MCPPolicyMock) RegisterErrorMocks() {
	httpmock.RegisterResponder(
		"POST",
		BaseURL,
		func(*http.Request) (*http.Response, error) { return apiError(400, "Invalid MCP policy") },
	)
}

func apiError(code int, message string) (*http.Response, error) {
	return httpmock.NewJsonResponse(
		code,
		map[string]any{"error": map[string]any{"code": fmt.Sprint(code), "message": message}},
	)
}

func fixture(rule bool) (map[string]any, error) {
	_, file, _, _ := runtime.Caller(0)
	dir := filepath.Join(
		filepath.Dir(file),
		"..",
		"tests",
		"responses",
		"validate_create",
		"post_mcp_policy_success.json",
	)
	if rule {
		dir = filepath.Join(
			filepath.Dir(file),
			"..",
			"..",
			"network_mcp_policy_rule",
			"tests",
			"responses",
			"validate_create",
			"post_mcp_policy_rule_success.json",
		)
	}
	raw, err := os.ReadFile(dir)
	if err != nil {
		return nil, err
	}
	var out map[string]any
	err = json.Unmarshal(raw, &out)
	return out, err
}

func merge(dst, src map[string]any) {
	for k, v := range src {
		if obj, ok := v.(map[string]any); ok {
			if existing, ok := dst[k].(map[string]any); ok {
				merge(existing, obj)
			} else {
				dst[k] = obj
			}
		} else {
			dst[k] = v
		}
	}
}

func (m *MCPPolicyMock) respond(req *http.Request) (*http.Response, error) {
	m.Lock()
	defer m.Unlock()
	body := map[string]any{}
	if req.Method == "POST" || req.Method == "PATCH" {
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			return apiError(400, "Invalid JSON")
		}
	}
	m.Requests = append(m.Requests, RecordedRequest{req.Method, req.URL.Path, body})
	path := strings.TrimPrefix(req.URL.Path, "/beta/networkaccess/mcpPolicies")
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if path == "" {
		if req.Method != "POST" {
			return apiError(405, "Unsupported collection request")
		}
		if _, ok := body["policyRules"]; ok {
			return apiError(400, "Policy resource must not own inline rules")
		}
		if body["name"] == nil || body["settings"] == nil {
			return apiError(400, "Name and Settings required")
		}
		policy, err := fixture(false)
		if err != nil {
			return nil, err
		}
		id := uuid.NewString()
		policy["id"] = id
		policy["description"] = nil
		merge(policy, body)
		policy["lastModifiedDateTime"] = time.Now().UTC().Format(time.RFC3339Nano)
		m.Policies[id] = policy
		m.Rules[id] = map[string]map[string]any{}
		return httpmock.NewJsonResponse(200, policy)
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
			if _, ok := body["policyRules"]; ok {
				return apiError(400, "Policy resource must not update rules")
			}
			merge(policy, body)
			policy["lastModifiedDateTime"] = time.Now().UTC().Format(time.RFC3339Nano)
			return httpmock.NewJsonResponse(200, policy)
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
		if body["@odata.type"] != "#microsoft.graph.networkaccess.mcpPolicyRule" {
			if req.Method == "POST" {
				return apiError(400, "Required property Priority is missing")
			}
			if len(parts) != 3 || rules[parts[2]] == nil {
				return apiError(404, "Rule not found")
			}
			// Without the derived discriminator Graph silently changes only base properties.
			if name, ok := body["name"]; ok {
				rules[parts[2]]["name"] = name
			}
			return httpmock.NewJsonResponse(200, rules[parts[2]])
		}
		if p, ok := body["priority"].(float64); ok {
			if p < 100 || p > 2147483647 {
				return apiError(400, "Invalid priority")
			}
		}
		for id, rule := range rules {
			if len(parts) == 3 && id == parts[2] {
				continue
			}
			for _, key := range []string{"name", "priority"} {
				if v, ok := body[key]; ok && v == rule[key] {
					return apiError(400, key+" already exists")
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
			for _, key := range []string{"name", "priority", "settings"} {
				if body[key] == nil {
					return apiError(400, key+" required")
				}
			}
			rule, err := fixture(true)
			if err != nil {
				return nil, err
			}
			id := uuid.NewString()
			rule["id"] = id
			rule["description"] = nil
			merge(rule, body)
			rules[id] = rule
			return httpmock.NewJsonResponse(200, rule)
		}
	}
	if len(parts) == 3 {
		rule, ok := rules[parts[2]]
		if !ok {
			return apiError(404, "Rule not found")
		}
		switch req.Method {
		case "GET":
			return httpmock.NewJsonResponse(200, rule)
		case "PATCH":
			merge(rule, body)
			if rule["matchingConditions"] == nil {
				rule["matchingConditions"] = map[string]any{"sources": nil, "destinations": nil}
			}
			return httpmock.NewJsonResponse(200, rule)
		case "DELETE":
			delete(rules, parts[2])
			return httpmock.NewStringResponse(204, ""), nil
		}
	}
	return apiError(405, "Unsupported request")
}
