package mocks

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/jarcoal/httpmock"

	common "github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
)

const (
	collectionURL            = "https://graph.microsoft.com/beta/networkaccess/connectivity/remoteNetworks"
	guidPattern              = `[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}`
	itemURLRegex             = `=~^https://graph\.microsoft\.com/beta/networkaccess/connectivity/remoteNetworks/(` + guidPattern + `)$`
	forwardingProfilesURLRgx = `=~^https://graph\.microsoft\.com/beta/networkaccess/connectivity/remoteNetworks/(` + guidPattern + `)/forwardingProfiles$`

	// Forwarding profile IDs known to the mock tenant.
	M365ProfileID     = "11111111-1111-1111-1111-111111111111"
	InternetProfileID = "22222222-2222-2222-2222-222222222222"
	DisabledProfileID = "33333333-3333-3333-3333-333333333333"
	PrivateProfileID  = "44444444-4444-4444-4444-444444444444"

	notFoundBody = `{"error":{"code":"NotFound","message":"Remote network not found."}}`
)

type profile struct {
	name        string
	trafficType string
	enabled     bool
}

var profiles = map[string]profile{
	M365ProfileID:     {name: "Microsoft 365 traffic forwarding profile", trafficType: "m365", enabled: true},
	InternetProfileID: {name: "Internet traffic forwarding profile", trafficType: "internet", enabled: true},
	DisabledProfileID: {name: "M365", trafficType: "m365", enabled: false},
	PrivateProfileID:  {name: "Private access traffic forwarding profile", trafficType: "private", enabled: true},
}

type remoteNetwork struct {
	name       string
	region     string
	profileIDs []string
}

// Request records a write sent to the mock.
type Request struct {
	Method string
	Path   string
	Body   map[string]any
}

// RemoteNetworkMock is a stateful mock enforcing the live-verified Graph rules for remote networks.
// After each write, the next GET returns the pre-write snapshot to simulate replica lag.
//
// DELETE is asynchronous, as in Graph: the next DeletionGets GETs return 404 "on deletion
// process", and until then POST rejects the name with 409.
type RemoteNetworkMock struct {
	sync.Mutex
	networks     map[string]*remoteNetwork
	stale        map[string]*remoteNetwork
	deleting     map[string]*deletingNetwork
	Requests     []Request
	StaleGets    int
	DeletionGets int
	// DeletionErrors makes the next GETs of a deleting network return 503.
	DeletionErrors int
	// DeleteReportsDeleting makes DELETE start the deletion but return 404
	// "on deletion process", as after an earlier interrupted destroy.
	DeleteReportsDeleting bool
}

type deletingNetwork struct {
	name     string
	getsLeft int
}

// defaultDeletionGets is how many GETs see a deleted network as still being deleted.
const defaultDeletionGets = 2

var _ common.MockRegistrar = (*RemoteNetworkMock)(nil)

func init() {
	common.GlobalRegistry.Register("identity_and_access_network_remote_network", &RemoteNetworkMock{})
}

func decode(req *http.Request) (map[string]any, error) {
	var reader io.Reader = req.Body
	if req.Header.Get("Content-Encoding") == "gzip" {
		gz, err := gzip.NewReader(req.Body)
		if err != nil {
			return nil, fmt.Errorf("decompress request: %w", err)
		}
		defer gz.Close()
		reader = gz
	}
	var body map[string]any
	if err := json.NewDecoder(reader).Decode(&body); err != nil {
		return nil, fmt.Errorf("decode request: %w", err)
	}
	return body, nil
}

func badRequest(message string) *http.Response {
	resp, _ := httpmock.NewJsonResponse(400, map[string]any{"error": map[string]any{"code": "BadRequest", "message": message}})
	return resp
}

func lastSegment(path string, fromEnd int) string {
	parts := strings.Split(strings.TrimSuffix(path, "/"), "/")
	return parts[len(parts)-1-fromEnd]
}

// profileIDs extracts and validates forwarding profile references like Graph does.
func profileIDs(raw any) ([]string, *http.Response) {
	items, _ := raw.([]any)
	ids := make([]string, 0, len(items))
	for _, item := range items {
		ref, _ := item.(map[string]any)
		id, _ := ref["id"].(string)
		if id == "" {
			continue
		}
		p, ok := profiles[id]
		switch {
		case !ok:
			return nil, badRequest(fmt.Sprintf("Forwarding profile %s was not found.", id))
		case p.trafficType == "private":
			return nil, badRequest(fmt.Sprintf("Forwarding profile %s is not allowed to be assigned to a remote network.", p.name))
		case !p.enabled:
			return nil, badRequest(fmt.Sprintf("Can't assign forwarding profile %s as it is disabled", p.name))
		}
		ids = append(ids, id)
	}
	if len(ids) == 0 {
		return nil, badRequest("At least one forwarding traffic profile must be assigned to the remote network")
	}
	return ids, nil
}

func (n *remoteNetwork) clone() *remoteNetwork {
	c := *n
	c.profileIDs = append([]string(nil), n.profileIDs...)
	return &c
}

func (n *remoteNetwork) json(id string) map[string]any {
	refs := make([]map[string]any, 0, len(n.profileIDs))
	for _, pid := range n.profileIDs {
		p := profiles[pid]
		refs = append(refs, map[string]any{
			"id":                    pid,
			"name":                  p.name,
			"trafficForwardingType": p.trafficType,
			"state":                 "enabled",
		})
	}
	return map[string]any{
		"id":                   id,
		"name":                 n.name,
		"region":               n.region,
		"version":              "1.0.0",
		"lastModifiedDateTime": "2026-01-01T00:00:00Z",
		"forwardingProfiles":   refs,
	}
}

func (m *RemoteNetworkMock) record(req *http.Request, body map[string]any) {
	m.Requests = append(m.Requests, Request{Method: req.Method, Path: req.URL.Path, Body: body})
}

func (m *RemoteNetworkMock) RegisterMocks() {
	m.CleanupMockState()

	httpmock.RegisterResponder("POST", collectionURL, func(req *http.Request) (*http.Response, error) {
		body, err := decode(req)
		if err != nil {
			return badRequest(err.Error()), nil
		}
		m.Lock()
		defer m.Unlock()
		m.record(req, body)

		name, _ := body["name"].(string)
		region, _ := body["region"].(string)
		if name == "" || region == "" {
			return badRequest("name and region are required"), nil
		}
		ids, errResp := profileIDs(body["forwardingProfiles"])
		if errResp != nil {
			return errResp, nil
		}
		if m.nameInUse(name) {
			return jsonError(409, fmt.Sprintf(
				`{"error":{"code":"Conflict","message":"A branch with the same name of %s already exists."}}`, name)), nil
		}

		id := uuid.New().String()
		network := &remoteNetwork{name: name, region: region, profileIDs: ids}
		m.networks[id] = network
		m.stale[id] = nil

		response := network.json(id)
		delete(response, "forwardingProfiles")
		return httpmock.NewJsonResponse(201, response)
	})

	httpmock.RegisterResponder("GET", itemURLRegex, func(req *http.Request) (*http.Response, error) {
		id := lastSegment(req.URL.Path, 0)
		m.Lock()
		defer m.Unlock()
		if d, ok := m.deleting[id]; ok {
			if m.DeletionErrors > 0 {
				m.DeletionErrors--
				return jsonError(503, `{"error":{"code":"ServiceUnavailable","message":"Service unavailable."}}`), nil
			}
			d.getsLeft--
			if d.getsLeft <= 0 {
				delete(m.deleting, id)
			}
			return jsonError(404, fmt.Sprintf(
				`{"error":{"code":"NotFound","message":"The requested remote network %s is on deletion process."}}`, id)), nil
		}
		if snapshot, lagging := m.stale[id]; lagging {
			delete(m.stale, id)
			m.StaleGets++
			if snapshot == nil {
				return httpmock.NewStringResponse(404, notFoundBody), nil
			}
			return httpmock.NewJsonResponse(200, snapshot.json(id))
		}
		network, ok := m.networks[id]
		if !ok {
			return httpmock.NewStringResponse(404, notFoundBody), nil
		}
		return httpmock.NewJsonResponse(200, network.json(id))
	})

	// The navigation property GET returns an empty collection even when profiles are associated.
	httpmock.RegisterResponder("GET", forwardingProfilesURLRgx, func(req *http.Request) (*http.Response, error) {
		return httpmock.NewJsonResponse(200, map[string]any{"value": []any{}})
	})

	httpmock.RegisterResponder("PATCH", itemURLRegex, func(req *http.Request) (*http.Response, error) {
		id := lastSegment(req.URL.Path, 0)
		body, err := decode(req)
		if err != nil {
			return badRequest(err.Error()), nil
		}
		m.Lock()
		defer m.Unlock()
		m.record(req, body)

		network, ok := m.networks[id]
		if !ok {
			return httpmock.NewStringResponse(404, notFoundBody), nil
		}
		if _, ok := body["forwardingProfiles"]; ok {
			return badRequest("Updating the property ForwardingProfiles for entity RemoteNetwork is not allowed."), nil
		}
		m.stale[id] = network.clone()
		if name, ok := body["name"].(string); ok {
			network.name = name
		}
		return httpmock.NewStringResponse(204, ""), nil
	})

	httpmock.RegisterResponder("PATCH", forwardingProfilesURLRgx, func(req *http.Request) (*http.Response, error) {
		id := lastSegment(req.URL.Path, 1)
		body, err := decode(req)
		if err != nil {
			return badRequest(err.Error()), nil
		}
		m.Lock()
		defer m.Unlock()
		m.record(req, body)

		network, ok := m.networks[id]
		if !ok {
			return httpmock.NewStringResponse(404, notFoundBody), nil
		}
		if body["@context"] != "#$delta" {
			return badRequest("Expected a #$delta payload"), nil
		}
		ids, errResp := profileIDs(body["value"])
		if errResp != nil {
			return errResp, nil
		}
		m.stale[id] = network.clone()
		network.profileIDs = ids
		return httpmock.NewStringResponse(204, ""), nil
	})

	httpmock.RegisterResponder("DELETE", itemURLRegex, func(req *http.Request) (*http.Response, error) {
		id := lastSegment(req.URL.Path, 0)
		m.Lock()
		defer m.Unlock()
		m.record(req, nil)
		network, ok := m.networks[id]
		if !ok {
			return httpmock.NewStringResponse(404, notFoundBody), nil
		}
		delete(m.networks, id)
		delete(m.stale, id)
		if m.DeletionGets > 0 {
			m.deleting[id] = &deletingNetwork{name: network.name, getsLeft: m.DeletionGets}
		}
		if m.DeleteReportsDeleting {
			return jsonError(404, fmt.Sprintf(
				`{"error":{"code":"NotFound","message":"The requested remote network %s is on deletion process."}}`, id)), nil
		}
		return httpmock.NewStringResponse(200, ""), nil
	})
}

func (m *RemoteNetworkMock) RegisterErrorMocks() {
	forbidden := httpmock.NewStringResponder(403, `{"error":{"code":"Forbidden","message":"Insufficient privileges to complete the operation."}}`)
	httpmock.RegisterResponder("POST", collectionURL, forbidden)
	httpmock.RegisterResponder("GET", itemURLRegex, httpmock.NewStringResponder(404, notFoundBody))
	httpmock.RegisterResponder("PATCH", itemURLRegex, forbidden)
	httpmock.RegisterResponder("PATCH", forwardingProfilesURLRgx, forbidden)
	httpmock.RegisterResponder("DELETE", itemURLRegex, forbidden)
}

func (m *RemoteNetworkMock) CleanupMockState() {
	m.Lock()
	defer m.Unlock()
	m.networks = make(map[string]*remoteNetwork)
	m.stale = make(map[string]*remoteNetwork)
	m.deleting = make(map[string]*deletingNetwork)
	m.Requests = nil
	m.StaleGets = 0
	m.DeletionGets = defaultDeletionGets
	m.DeletionErrors = 0
	m.DeleteReportsDeleting = false
}

// jsonError returns a Graph error with a JSON content type so Kiota parses the message.
func jsonError(status int, body string) *http.Response {
	resp := httpmock.NewStringResponse(status, body)
	resp.Header.Set("Content-Type", "application/json")
	return resp
}

// nameInUse reports whether a live or still-deleting network has the name.
func (m *RemoteNetworkMock) nameInUse(name string) bool {
	for _, n := range m.networks {
		if n.name == name {
			return true
		}
	}
	for _, d := range m.deleting {
		if d.name == name {
			return true
		}
	}
	return false
}

// DeleteAll removes every remote network, simulating deletion outside Terraform.
func (m *RemoteNetworkMock) DeleteAll() {
	m.Lock()
	defer m.Unlock()
	m.networks = make(map[string]*remoteNetwork)
	m.stale = make(map[string]*remoteNetwork)
}

// WritesTo returns the recorded writes with the given method whose path ends with suffix.
func (m *RemoteNetworkMock) WritesTo(method, suffix string) []Request {
	m.Lock()
	defer m.Unlock()
	var out []Request
	for _, r := range m.Requests {
		if r.Method == method && strings.HasSuffix(r.Path, suffix) {
			out = append(out, r)
		}
	}
	return out
}
