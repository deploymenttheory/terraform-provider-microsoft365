package mocks

import (
	"encoding/json"
	"fmt"
	"maps"
	"net/http"
	"strings"
	"sync"

	"github.com/jarcoal/httpmock"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	commonMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
)

// MobileAppSupersedenceMock stores the relationship collection between Terraform steps.
type MobileAppSupersedenceMock struct {
	mu            sync.Mutex
	relationships map[string]map[string]any
}

var _ commonMocks.MockRegistrar = &MobileAppSupersedenceMock{}

func init() {
	commonMocks.GlobalRegistry.Register("mobile_app_supersedence", &MobileAppSupersedenceMock{})
}

func (m *MobileAppSupersedenceMock) RegisterMocks() {
	m.CleanupMockState()
	template := loadResponse("validate_get/get_supersedence.json")
	badRoute := loadResponse("validate_get/get_supersedence_bad_route.json")
	httpmock.RegisterResponder(http.MethodGet,
		`=~^https://graph\.microsoft\.com/beta/deviceAppManagement/mobileApps/[^/]+/relationships$`,
		func(req *http.Request) (*http.Response, error) {
			m.mu.Lock()
			defer m.mu.Unlock()
			sourceID := appIDFromRequest(req)
			records := make([]map[string]any, 0)
			for _, record := range m.relationships {
				if record["sourceId"] == sourceID {
					records = append(records, record)
				}
			}
			return jsonResponse(http.StatusOK, map[string]any{"value": records})
		})
	httpmock.RegisterResponder(
		http.MethodPost,
		`=~^https://graph\.microsoft\.com/beta/deviceAppManagement/mobileApps/[^/]+/updateRelationships$`,
		func(req *http.Request) (*http.Response, error) {
			var body struct {
				Relationships []map[string]any `json:"relationships"`
			}
			if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
				return nil, fmt.Errorf("decode relationship update: %w", err)
			}
			m.mu.Lock()
			defer m.mu.Unlock()
			sourceID := appIDFromRequest(req)
			for id, record := range m.relationships {
				if record["sourceId"] == sourceID {
					delete(m.relationships, id)
				}
			}
			for _, relationship := range body.Relationships {
				record := maps.Clone(template)
				maps.Copy(record, relationship)
				targetID, _ := record["targetId"].(string)
				id, _ := relationship["id"].(string)
				if id == "" {
					id = sourceID + "_" + targetID
				}
				record["id"], record["sourceId"] = id, sourceID
				m.relationships[id] = record
			}
			return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
		},
	)
	httpmock.RegisterResponder(http.MethodGet,
		`=~^https://graph\.microsoft\.com/beta/deviceAppManagement/mobileAppRelationships/[^/]+$`,
		func(_ *http.Request) (*http.Response, error) {
			return jsonResponse(http.StatusBadRequest, badRoute)
		})
}

func (m *MobileAppSupersedenceMock) RegisterErrorMocks() {
	m.RegisterMocks()
	response := loadResponse("validate_create/post_update_relationships_error.json")
	httpmock.RegisterResponder(
		http.MethodPost,
		`=~^https://graph\.microsoft\.com/beta/deviceAppManagement/mobileApps/[^/]+/updateRelationships$`,
		func(_ *http.Request) (*http.Response, error) { return jsonResponse(http.StatusBadRequest, response) },
	)
}

func (m *MobileAppSupersedenceMock) CleanupMockState() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.relationships = make(map[string]map[string]any)
}

func appIDFromRequest(req *http.Request) string {
	segments := strings.Split(req.URL.Path, "/")
	return segments[len(segments)-2]
}

func loadResponse(filename string) map[string]any {
	content, err := helpers.ParseJSONFile("../tests/responses/" + filename)
	if err != nil {
		panic("failed to load supersedence response fixture: " + err.Error())
	}
	var response map[string]any
	if err := json.Unmarshal([]byte(content), &response); err != nil {
		panic("invalid supersedence response fixture: " + err.Error())
	}
	return response
}

func jsonResponse(status int, value any) (*http.Response, error) {
	response, err := httpmock.NewJsonResponse(status, value)
	if err != nil {
		return nil, fmt.Errorf("create supersedence mock response: %w", err)
	}
	return response, nil
}
