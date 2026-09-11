package mocks

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/google/uuid"
	"github.com/jarcoal/httpmock"
)

// CloudFirewallMock models independent policy/rule lifecycles and nested PATCH merging.
type CloudFirewallMock struct {
	sync.Mutex
	objects  map[string]map[string]any
	requests []string
}

var _ mocks.MockRegistrar = (*CloudFirewallMock)(nil)

func init() { mocks.GlobalRegistry.Register("network_cloud_firewall_policy", &CloudFirewallMock{}) }
func (m *CloudFirewallMock) CleanupMockState() {
	m.Lock()
	defer m.Unlock()
	m.objects = make(map[string]map[string]any)
	m.requests = nil
}
func (m *CloudFirewallMock) RegisterMocks() {
	m.CleanupMockState()
	for _, method := range []string{"GET", "POST", "PATCH", "DELETE"} {
		httpmock.RegisterResponder(method, `=~^https://graph\.microsoft\.com/beta/networkAccess/cloudFirewallPolicies(?:/.*)?$`, m.respond)
	}
}
func (m *CloudFirewallMock) RegisterErrorMocks() {
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/networkAccess/cloudFirewallPolicies", httpmock.NewStringResponder(400, `{"error":{"code":"BadRequest","message":"Invalid cloud firewall policy"}}`))
}

// Count returns request counts for assertions that local-only changes send no PATCH.
func (m *CloudFirewallMock) Count(method string) int {
	m.Lock()
	defer m.Unlock()
	n := 0
	for _, v := range m.requests {
		if strings.HasPrefix(v, method+" ") {
			n++
		}
	}
	return n
}
func (m *CloudFirewallMock) respond(req *http.Request) (*http.Response, error) {
	m.Lock()
	defer m.Unlock()
	path := strings.TrimPrefix(req.URL.Path, "/beta/networkAccess/cloudFirewallPolicies")
	m.requests = append(m.requests, req.Method+" "+path)
	fail := func(status int, msg string) (*http.Response, error) {
		return jsonResponse(status, map[string]any{"error": map[string]any{"code": "BadRequest", "message": msg}})
	}
	switch req.Method {
	case "POST":
		return m.create(req, path, fail)
	case "GET":
		body, ok := m.objects[path]
		if !ok {
			return fail(404, "Not found")
		}
		return jsonResponse(200, body)
	case "PATCH":
		body, ok := m.objects[path]
		if !ok {
			return fail(404, "Not found")
		}
		var patch map[string]any
		if err := json.NewDecoder(req.Body).Decode(&patch); err != nil {
			return fail(400, "Invalid JSON")
		}
		if strings.Contains(path, "/policyRules/") && patch["@odata.type"] != "#microsoft.graph.networkaccess.cloudFirewallRule" {
			return fail(400, "Missing rule discriminator")
		}
		merge(body, patch)
		return httpmock.NewStringResponse(204, ""), nil
	case "DELETE":
		if m.objects[path] == nil {
			return fail(404, "Not found")
		}
		for key := range m.objects {
			if key == path || strings.HasPrefix(key, path+"/") {
				delete(m.objects, key)
			}
		}
		return httpmock.NewStringResponse(204, ""), nil
	}
	return fail(405, "Method not allowed")
}
func merge(target, patch map[string]any) {
	for key, value := range patch {
		nested, ok := value.(map[string]any)
		old, exists := target[key].(map[string]any)
		if ok && exists {
			merge(old, nested)
		} else {
			target[key] = value
		}
	}
}

func (m *CloudFirewallMock) create(req *http.Request, path string, fail func(int, string) (*http.Response, error)) (*http.Response, error) {

	var body map[string]any
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		return fail(400, "Invalid JSON")
	}
	rule := strings.HasSuffix(path, "/policyRules")
	if rule {
		parent := strings.TrimSuffix(path, "/policyRules")
		if m.objects[parent] == nil {
			return fail(404, "Parent not found")
		}
	}
	if rule {
		if body["@odata.type"] != "#microsoft.graph.networkaccess.cloudFirewallRule" {
			return fail(400, "Missing rule discriminator")
		}
		settings, ok := body["settings"].(map[string]any)
		if !ok || (settings["status"] != "enabled" && settings["status"] != "disabled") {
			return fail(400, "Invalid status")
		}
		for key, existing := range m.objects {
			if strings.HasPrefix(key, path+"/") && (existing["name"] == body["name"] || existing["priority"] == body["priority"]) {
				return fail(400, "Duplicate name or priority")
			}
		}
	} else {
		settings, ok := body["settings"].(map[string]any)
		if !ok || settings["defaultAction"] != "allow" {
			return fail(400, "Invalid defaultAction")
		}
		if _, ok := body["policyRules"]; ok {
			return fail(400, "Inline rules must not be authored")
		}
	}
	id := uuid.NewString()
	body["id"] = id
	if _, ok := body["description"]; !ok {
		body["description"] = nil
	}
	if !rule {
		body["version"] = "1.0.0"
		body["lastModifiedDateTime"] = "2026-09-07T00:00:00Z"
	}
	m.objects[path+"/"+id] = body
	return jsonResponse(201, body)
}

func jsonResponse(status int, body any) (*http.Response, error) {
	response, err := httpmock.NewJsonResponse(status, body)
	if err != nil {
		return nil, fmt.Errorf("encode cloud firewall mock response: %w", err)
	}
	return response, nil
}
