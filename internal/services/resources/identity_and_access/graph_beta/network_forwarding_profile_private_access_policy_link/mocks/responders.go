package mocks

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"sync"

	"github.com/jarcoal/httpmock"

	common "github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
)

const (
	UserCollection   = "forwardingProfiles"
	DeviceCollection = "deviceForwardingProfiles"

	UserProfileID      = "11111111-1111-1111-1111-111111111111"
	DeviceProfileID    = "22222222-2222-2222-2222-222222222222"
	AppPolicyID        = "33333333-3333-3333-3333-333333333333"
	AutoLinkedPolicyID = "44444444-4444-4444-4444-444444444444"
	DNSPolicyID        = "55555555-5555-5555-5555-555555555555"
	AutoLinkID         = "aaaaaaaa-0000-0000-0000-000000000001"
)

var profilePattern = regexp.MustCompile(
	`^https://graph\.microsoft\.com/beta/networkAccess/(forwardingProfiles|deviceForwardingProfiles)/([0-9a-fA-F-]{36})(\?.*)?$`,
)

// Link is the mock's view of a forwardingPolicyLink.
type Link struct {
	ID         string
	PolicyID   string
	State      string
	ConfigType string
	Priority   int64
}

// Mock models a profile's policies with the live add/remove-only delta rules.
type Mock struct {
	sync.Mutex
	Profiles map[string][]Link
	Patches  []map[string]any
	// StaleReads makes the next N profile GETs after each write return the pre-write links.
	StaleReads int
	// FailAdds makes the next N PATCHes that add a link fail with 400.
	FailAdds   int
	stale      map[string][]Link
	staleLeft  int
	nextLinkID int
}

var _ common.MockRegistrar = (*Mock)(nil)

func init() {
	common.GlobalRegistry.Register("network_forwarding_profile_private_access_policy_link", &Mock{})
}

func Key(collection, id string) string { return collection + "/" + id }

func policyName(id string) (string, bool) {
	switch id {
	case AppPolicyID:
		return "App One", true
	case AutoLinkedPolicyID:
		return "Auto Linked App", true
	case DNSPolicyID:
		return "Private DNS", true
	}
	return "", false
}

func decode(q *http.Request) (map[string]any, error) {
	var reader io.Reader = q.Body
	if q.Header.Get("Content-Encoding") == "gzip" {
		gz, err := gzip.NewReader(q.Body)
		if err != nil {
			return nil, fmt.Errorf("decompress request: %w", err)
		}
		defer gz.Close()
		reader = gz
	}
	raw, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("read request: %w", err)
	}
	var body map[string]any
	if err := json.Unmarshal(raw, &body); err != nil {
		return nil, fmt.Errorf("decode request: %w", err)
	}
	body[rawBodyKey] = string(raw)
	return body, nil
}

// rawBodyKey keeps the request bytes so patch can check key order.
const rawBodyKey = "__raw"

// removedFirst reports whether every delta entry carrying @removed has it as
// its first key; Graph rejects the removal otherwise.
func removedFirst(raw string) bool {
	return strings.Count(raw, `"@removed"`) == strings.Count(raw, `{"@removed"`)
}

func jsonResponse(status int, body any) (*http.Response, error) {
	resp, err := httpmock.NewJsonResponse(status, body)
	if err != nil {
		return nil, fmt.Errorf("encode mock response: %w", err)
	}
	return resp, nil
}

func graphError(status int, code, message string) (*http.Response, error) {
	return jsonResponse(status, map[string]any{
		"error": map[string]any{"code": code, "message": message},
	})
}

func profileJSON(id string, links []Link) map[string]any {
	policies := make([]map[string]any, 0, len(links))
	for _, l := range links {
		name, _ := policyName(l.PolicyID)
		policies = append(policies, map[string]any{
			"@odata.type": "#microsoft.graph.networkaccess.forwardingPolicyLink",
			"id":          l.ID,
			"state":       l.State,
			"priority":    l.Priority,
			"configType":  l.ConfigType,
			"version":     "1.0.0",
			"policy": map[string]any{
				"@odata.type":           "#microsoft.graph.networkaccess.forwardingPolicy",
				"id":                    l.PolicyID,
				"name":                  name,
				"trafficForwardingType": "private",
			},
		})
	}
	return map[string]any{"id": id, "name": "Private profile", "policies": policies}
}

// RegisterMocks seeds a user profile auto-linked to AutoLinkedPolicyID and an empty device profile.
func (m *Mock) RegisterMocks() {
	m.CleanupMockState()
	m.Lock()
	m.Profiles[Key(UserCollection, UserProfileID)] = []Link{{
		ID: AutoLinkID, PolicyID: AutoLinkedPolicyID, State: "enabled", ConfigType: "appPolicy", Priority: 65000,
	}}
	m.Profiles[Key(DeviceCollection, DeviceProfileID)] = []Link{}
	m.Unlock()

	httpmock.RegisterRegexpResponder("GET", profilePattern, m.get)
	httpmock.RegisterRegexpResponder("PATCH", profilePattern, m.patch)
	httpmock.RegisterRegexpResponder("POST", regexp.MustCompile(`/policies$`),
		func(*http.Request) (*http.Response, error) {
			return graphError(400, "BadRequest", "POST to policies is not supported")
		})
}

func (m *Mock) profileKey(q *http.Request) string {
	match := profilePattern.FindStringSubmatch(q.URL.String())
	if match == nil {
		return ""
	}
	return Key(match[1], strings.ToLower(match[2]))
}

func (m *Mock) get(q *http.Request) (*http.Response, error) {
	key := m.profileKey(q)
	m.Lock()
	defer m.Unlock()
	links, ok := m.Profiles[key]
	if m.staleLeft > 0 {
		m.staleLeft--
		links, ok = m.stale[key]
	}
	if !ok {
		return graphError(404, "NotFound", "Forwarding profile not found")
	}
	return jsonResponse(200, profileJSON(key[strings.Index(key, "/")+1:], links))
}

func (m *Mock) patch(q *http.Request) (*http.Response, error) {
	key := m.profileKey(q)
	body, err := decode(q)
	if err != nil {
		return nil, err
	}
	m.Lock()
	defer m.Unlock()
	raw, _ := body[rawBodyKey].(string)
	delete(body, rawBodyKey)
	m.Patches = append(m.Patches, body)
	links, ok := m.Profiles[key]
	if !ok {
		return graphError(404, "NotFound", "Forwarding profile not found")
	}
	delta, _ := body["policies@delta"].([]any)
	if len(delta) == 0 {
		return graphError(400, "BadRequest", "policies@delta is required")
	}
	if !removedFirst(raw) {
		return graphError(400, "UnknownError", `{"error":"Invalid request parameters"}`)
	}

	next := append([]Link(nil), links...)
	for _, raw := range delta {
		entry, _ := raw.(map[string]any)
		id, _ := entry["id"].(string)
		if _, removed := entry["@removed"]; removed {
			idx := indexByID(next, id)
			if idx < 0 {
				return graphError(400, "BadRequest", fmt.Sprintf("PolicyLink with Id %s was not found", id))
			}
			next = append(next[:idx], next[idx+1:]...)
			continue
		}
		if id != "" {
			return graphError(400, "BadRequest",
				"Updating a PolicyLink is not supported. Only Deleting or Adding policy links are currently supported")
		}
		if m.FailAdds > 0 {
			m.FailAdds--
			return graphError(400, "BadRequest", "Simulated add failure")
		}
		policy, _ := entry["policy"].(map[string]any)
		policyID, _ := policy["id"].(string)
		if _, known := policyName(policyID); !known {
			return graphError(400, "BadRequest",
				fmt.Sprintf("Policy with Id %s of traffic type Private was not found", policyID))
		}
		// Checked against the pre-request links: remove+add of one policy in a single delta fails live.
		if indexByPolicy(links, policyID) >= 0 || indexByPolicy(next, policyID) >= 0 {
			return graphError(400, "BadRequest",
				fmt.Sprintf("Policy with Id %s is already linked to profile with id %s", policyID, key))
		}
		link := Link{PolicyID: policyID, Priority: 65000}
		link.State, _ = entry["state"].(string)
		link.ConfigType, _ = entry["configType"].(string)
		if p, ok := entry["priority"].(float64); ok {
			link.Priority = int64(p)
		}
		m.nextLinkID++
		link.ID = fmt.Sprintf("bbbbbbbb-0000-0000-0000-%012d", m.nextLinkID)
		next = append(next, link)
	}

	if m.StaleReads > 0 {
		m.stale = map[string][]Link{}
		for k, v := range m.Profiles {
			m.stale[k] = append([]Link(nil), v...)
		}
		m.staleLeft = m.StaleReads
	}
	m.Profiles[key] = next
	return jsonResponse(200, profileJSON(key[strings.Index(key, "/")+1:], next))
}

func indexByID(links []Link, id string) int {
	for i, l := range links {
		if strings.EqualFold(l.ID, id) {
			return i
		}
	}
	return -1
}

func indexByPolicy(links []Link, policyID string) int {
	for i, l := range links {
		if strings.EqualFold(l.PolicyID, policyID) {
			return i
		}
	}
	return -1
}

// SetLinks replaces a profile's links; use it to simulate out-of-band changes.
func (m *Mock) SetLinks(collection, profileID string, links []Link) {
	m.Lock()
	defer m.Unlock()
	m.Profiles[Key(collection, profileID)] = links
}

// DeleteProfile simulates the profile being removed out of band.
func (m *Mock) DeleteProfile(collection, profileID string) {
	m.Lock()
	defer m.Unlock()
	delete(m.Profiles, Key(collection, profileID))
}

// Links returns a copy of a profile's links.
func (m *Mock) Links(collection, profileID string) []Link {
	m.Lock()
	defer m.Unlock()
	return append([]Link(nil), m.Profiles[Key(collection, profileID)]...)
}

func (m *Mock) RegisterErrorMocks() {
	for _, method := range []string{"GET", "PATCH"} {
		httpmock.RegisterRegexpResponder(method, profilePattern, httpmock.NewStringResponder(
			403, `{"error":{"code":"Forbidden","message":"Synthetic denied request"}}`,
		))
	}
}

func (m *Mock) CleanupMockState() {
	m.Lock()
	defer m.Unlock()
	m.Profiles = map[string][]Link{}
	m.Patches = nil
	m.StaleReads = 0
	m.stale = nil
	m.staleLeft = 0
	m.nextLinkID = 0
}
