package mocks

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/google/uuid"
	"github.com/jarcoal/httpmock"
)

const policyURL = "https://graph.microsoft.com/beta/deviceManagement/compliancePolicies"

// LinuxDeviceCompliancePolicyMock models the separate policy, settings, and
// assignment collections. The individual policy expansion deliberately stays
// empty, matching the live API behavior that motivated PR #3914.
type LinuxDeviceCompliancePolicyMock struct {
	mu             sync.Mutex
	policies       map[string]map[string]any
	assignments    map[string][]any
	writes         int
	ExpectedWrites []string
}

var _ mocks.MockRegistrar = (*LinuxDeviceCompliancePolicyMock)(nil)

func init() {
	mocks.GlobalRegistry.Register("linux_device_compliance_policy", &LinuxDeviceCompliancePolicyMock{})
}

func fixture(name string) (map[string]any, error) {
	raw, err := helpers.ParseJSONFile("../tests/responses/" + name)
	if err != nil {
		return nil, fmt.Errorf("load Linux compliance fixture: %w", err)
	}
	var body map[string]any
	if err := json.Unmarshal([]byte(raw), &body); err != nil {
		return nil, fmt.Errorf("decode Linux compliance fixture: %w", err)
	}
	return body, nil
}

func fixtureResponse(status int, name string) (*http.Response, error) {
	body, err := fixture(name)
	if err != nil {
		return nil, err
	}
	return httpmock.NewJsonResponse(status, body)
}

func (m *LinuxDeviceCompliancePolicyMock) RegisterMocks() {
	m.CleanupMockState()
	httpmock.RegisterResponder("POST", policyURL, m.writePolicy)
	// The shared custom PUT helper uses OData key syntax, unlike SDK GET/DELETE.
	httpmock.RegisterResponder("PUT", `=~^https://graph\.microsoft\.com/beta/deviceManagement/compliancePolicies\('[^']+'\)$`, m.writePolicy)
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/compliancePolicies/[^/?]+(?:\?.*)?$`, m.readPolicy)
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/compliancePolicies/[^/?]+/(settings|assignments)(?:\?.*)?$`, m.readCollection)
	httpmock.RegisterResponder("POST", `=~^https://graph\.microsoft\.com/beta/deviceManagement/compliancePolicies/[^/?]+/assign$`, m.assign)
	httpmock.RegisterResponder("DELETE", `=~^https://graph\.microsoft\.com/beta/deviceManagement/compliancePolicies/[^/?]+$`, func(req *http.Request) (*http.Response, error) {
		m.mu.Lock()
		defer m.mu.Unlock()
		id := policyID(req)
		if _, ok := m.policies[id]; !ok {
			return fixtureResponse(400, "validate_delete/get_policy_not_found.json")
		}
		delete(m.policies, id)
		delete(m.assignments, id)
		return httpmock.NewStringResponse(200, ""), nil
	})
}

func (m *LinuxDeviceCompliancePolicyMock) RegisterErrorMocks() {
	m.RegisterMocks()
	httpmock.RegisterResponder("POST", policyURL, func(_ *http.Request) (*http.Response, error) {
		return fixtureResponse(400, "validate_create/post_009_error_scenario.json")
	})
}

func (m *LinuxDeviceCompliancePolicyMock) CleanupMockState() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.policies = make(map[string]map[string]any)
	m.assignments = make(map[string][]any)
	m.writes = 0
}

// CheckDestroyed also verifies that each expected POST/PUT was actually made.
func (m *LinuxDeviceCompliancePolicyMock) CheckDestroyed() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if len(m.policies) != 0 || len(m.assignments) != 0 {
		return fmt.Errorf("Linux compliance mock still contains %d policies and %d assignment collections", len(m.policies), len(m.assignments))
	}
	if m.writes != len(m.ExpectedWrites) {
		return fmt.Errorf("expected %d policy writes, got %d", len(m.ExpectedWrites), m.writes)
	}
	return nil
}

func policyID(req *http.Request) string {
	if _, key, ok := strings.Cut(req.URL.Path, "('"); ok {
		return strings.TrimSuffix(key, "')")
	}
	parts := strings.Split(strings.Trim(req.URL.Path, "/"), "/")
	return parts[3]
}

// canonicalSettings compares the wire settings independently of Terraform state.
// Setting and child ordering is not significant to Graph; distribution list
// ordering is exercised separately by the Terraform state assertions.
func canonicalSettings(value any) string {
	switch v := value.(type) {
	case []any:
		items := make([]string, 0, len(v))
		for _, item := range v {
			items = append(items, canonicalSettings(item))
		}
		sort.Strings(items)
		return "[" + strings.Join(items, ",") + "]"
	case map[string]any:
		items := make([]string, 0, len(v))
		for key, item := range v {
			items = append(items, strconv.Quote(key)+":"+canonicalSettings(item))
		}
		sort.Strings(items)
		return "{" + strings.Join(items, ",") + "}"
	default:
		encoded, _ := json.Marshal(value)
		return string(encoded)
	}
}

func (m *LinuxDeviceCompliancePolicyMock) writePolicy(req *http.Request) (*http.Response, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	var body map[string]any
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		return nil, fmt.Errorf("decode policy request: %w", err)
	}
	if body["platforms"] != "linux" || body["technologies"] != "linuxMdm" {
		return nil, fmt.Errorf("unexpected Linux policy platform or technology: %v", body)
	}
	if m.writes >= len(m.ExpectedWrites) {
		return nil, fmt.Errorf("unexpected policy write %s %s", req.Method, req.URL)
	}
	expected, err := fixture(m.ExpectedWrites[m.writes])
	if err != nil {
		return nil, err
	}
	if canonicalSettings(body["settings"]) != canonicalSettings(expected["settings"]) {
		return nil, fmt.Errorf("policy settings do not match %s: got %s", m.ExpectedWrites[m.writes], canonicalSettings(body["settings"]))
	}
	id := uuid.NewString()
	if req.Method == http.MethodPut {
		id = policyID(req)
		if _, ok := m.policies[id]; !ok {
			return fixtureResponse(400, "validate_delete/get_policy_not_found.json")
		}
	}
	response, err := fixture("validate_read/get_policy.json")
	if err != nil {
		return nil, err
	}
	for k, v := range body {
		response[k] = v
	}
	response["id"] = id
	response["settingCount"] = len(body["settings"].([]any))
	if req.Method == http.MethodPut {
		response["lastModifiedDateTime"] = "2026-09-17T01:00:00Z"
	} else {
		m.assignments[id] = []any{}
	}
	m.policies[id] = response
	m.writes++
	if req.Method == http.MethodPut {
		return httpmock.NewStringResponse(204, ""), nil
	}
	return httpmock.NewJsonResponse(201, response)
}

func (m *LinuxDeviceCompliancePolicyMock) readPolicy(req *http.Request) (*http.Response, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	stored, ok := m.policies[policyID(req)]
	if !ok {
		return fixtureResponse(400, "validate_delete/get_policy_not_found.json")
	}
	response := make(map[string]any, len(stored))
	for key, value := range stored {
		if key != "settings" {
			response[key] = value
		}
	}
	// Keep the expansion empty even after assignment, just like the live API.
	response["assignments"] = []any{}
	return httpmock.NewJsonResponse(200, response)
}

func (m *LinuxDeviceCompliancePolicyMock) readCollection(req *http.Request) (*http.Response, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	id := policyID(req)
	policy, ok := m.policies[id]
	if !ok {
		return fixtureResponse(400, "validate_delete/get_policy_not_found.json")
	}
	items := m.assignments[id]
	if strings.HasSuffix(req.URL.Path, "/settings") {
		items = policy["settings"].([]any)
	}
	start := 0
	if cursor := req.URL.Query().Get("$skiptoken"); cursor != "" {
		var err error
		start, err = strconv.Atoi(cursor)
		if err != nil || start < 0 || start > len(items) {
			return nil, fmt.Errorf("invalid collection cursor %q", cursor)
		}
	}
	// Small pages make maximal Terraform scenarios exercise both SDK iterators.
	end := min(start+2, len(items))
	response := map[string]any{"value": items[start:end]}
	if end < len(items) {
		response["@odata.nextLink"] = policyURL + "/" + id + "/" + req.URL.Path[strings.LastIndex(req.URL.Path, "/")+1:] + "?%24skiptoken=" + strconv.Itoa(end)
	}
	return httpmock.NewJsonResponse(200, response)
}

func (m *LinuxDeviceCompliancePolicyMock) assign(req *http.Request) (*http.Response, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	id := policyID(req)
	if _, ok := m.policies[id]; !ok {
		return fixtureResponse(400, "validate_delete/get_policy_not_found.json")
	}
	var body struct {
		Assignments []any `json:"assignments"`
	}
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		return nil, fmt.Errorf("decode assignment request: %w", err)
	}
	if body.Assignments == nil {
		return nil, fmt.Errorf("assignment request must contain an assignments array")
	}
	for _, value := range body.Assignments {
		assignment := value.(map[string]any)
		assignment["id"] = uuid.NewString()
		assignment["source"] = "direct"
		target := assignment["target"].(map[string]any)
		if _, ok := target["deviceAndAppManagementAssignmentFilterType"]; !ok {
			target["deviceAndAppManagementAssignmentFilterType"] = "none"
		}
	}
	m.assignments[id] = body.Assignments
	return httpmock.NewStringResponse(200, ""), nil
}
