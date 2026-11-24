package graphBetaNetworkFilteringProfilePolicyLink_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"sort"
	"sync"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks/factories"
	target "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_beta/network_filtering_profile_policy_link"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

const resourceType = target.ResourceName

// linkMock is a stateful mock of the policy links of /networkAccess/filteringProfiles/{id}.
type linkMock struct {
	sync.Mutex
	links map[string]map[string]any
	next  int
}

func setupMockEnvironment(t *testing.T) *linkMock {
	t.Helper()
	mocks.SetupUnitTestEnvironment(t)
	httpmock.Activate()
	t.Cleanup(httpmock.DeactivateAndReset)
	mocks.NewMocks().AuthMocks.RegisterMocks()

	m := &linkMock{links: map[string]map[string]any{}}
	profile := regexp.MustCompile(`^https://graph\.microsoft\.com/beta/networkAccess/filteringProfiles/([^/?]+)(\?.*)?$`)
	collection := regexp.MustCompile(`^https://graph\.microsoft\.com/beta/networkAccess/filteringProfiles/([^/]+)/policies$`)
	item := regexp.MustCompile(`^https://graph\.microsoft\.com/beta/networkAccess/filteringProfiles/([^/]+)/policies/([^/?]+)$`)
	decode := func(req *http.Request) (map[string]any, *http.Response) {
		var body map[string]any
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			return nil, httpmock.NewStringResponse(http.StatusBadRequest, `{"error":{"code":"BadRequest","message":"invalid JSON"}}`)
		}
		return body, nil
	}

	httpmock.RegisterRegexpResponder(http.MethodGet, profile, func(req *http.Request) (*http.Response, error) {
		m.Lock()
		defer m.Unlock()
		ids := make([]string, 0, len(m.links))
		for id := range m.links {
			ids = append(ids, id)
		}
		sort.Strings(ids)
		policies := make([]any, 0, len(ids))
		for _, id := range ids {
			policies = append(policies, m.links[id])
		}
		return factories.SuccessResponse(http.StatusOK, map[string]any{"id": profile.FindStringSubmatch(req.URL.String())[1], "name": "unit-test-profile", "policies": policies})(req)
	})
	httpmock.RegisterRegexpResponder(http.MethodPost, collection, func(req *http.Request) (*http.Response, error) {
		body, invalid := decode(req)
		if invalid != nil {
			return invalid, nil
		}
		m.Lock()
		defer m.Unlock()
		m.next++
		body["id"] = fmt.Sprintf("00000000-0000-0000-0000-%012d", m.next)
		body["version"] = "1.0.0"
		m.links[body["id"].(string)] = body
		return factories.SuccessResponse(http.StatusCreated, body)(req)
	})
	httpmock.RegisterRegexpResponder(http.MethodPatch, item, func(req *http.Request) (*http.Response, error) {
		body, invalid := decode(req)
		if invalid != nil {
			return invalid, nil
		}
		m.Lock()
		defer m.Unlock()
		link, ok := m.links[item.FindStringSubmatch(req.URL.String())[2]]
		if !ok {
			return factories.ErrorResponse(http.StatusNotFound, "ResourceNotFound", "policy link not found")(req)
		}
		for key, value := range body {
			link[key] = value
		}
		return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
	})
	httpmock.RegisterRegexpResponder(http.MethodDelete, item, func(req *http.Request) (*http.Response, error) {
		m.Lock()
		defer m.Unlock()
		delete(m.links, item.FindStringSubmatch(req.URL.String())[2])
		return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
	})
	return m
}

func testConfigHelper(t *testing.T, filename string) string {
	t.Helper()
	config, err := helpers.ParseHCLFile("tests/terraform/unit/" + filename)
	if err != nil {
		t.Fatalf("load unit test config %s: %v", filename, err)
	}
	return config
}

func assertNoLinks(t *testing.T, m *linkMock) {
	t.Helper()
	m.Lock()
	defer m.Unlock()
	if len(m.links) != 0 {
		t.Fatalf("destroy left %d policy links", len(m.links))
	}
}

func TestUnitResourceNetworkFilteringProfilePolicyLink_01_FilteringLifecycle(t *testing.T) {
	m := setupMockEnvironment(t)
	capture := mocks.CaptureGraphRequests(t)
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigHelper(t, "resource_filtering.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("policy_link_id").Exists(),
					check.That(resourceType+".test").Key("state").HasValue("enabled"),
					check.That(resourceType+".test").Key("priority").HasValue("100"),
					check.That(resourceType+".test").Key("logging_state").HasValue("enabled"),
				),
			},
			{ResourceName: resourceType + ".test", ImportState: true, ImportStateVerify: true, ImportStateVerifyIgnore: []string{"timeouts"}},
			{
				Config: testConfigHelper(t, "resource_filtering_disabled.tf"),
				Check:  check.That(resourceType+".test").Key("state").HasValue("disabled"),
			},
			{
				Config: testConfigHelper(t, "resource_filtering_custom.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("priority").HasValue("200"),
					check.That(resourceType+".test").Key("logging_state").HasValue("disabled"),
				),
			},
		},
	})
	capture.AssertMatchesFixture(t, "tests/requests/filtering_lifecycle.json")
	assertNoLinks(t, m)
}

func TestUnitResourceNetworkFilteringProfilePolicyLink_02_WebFilteringLifecycle(t *testing.T) {
	m := setupMockEnvironment(t)
	capture := mocks.CaptureGraphRequests(t)
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigHelper(t, "resource_web_filtering.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("policy_type").HasValue("web_filtering"),
					check.That(resourceType+".test").Key("priority").DoesNotExist(),
				),
			},
			{ResourceName: resourceType + ".test", ImportState: true, ImportStateVerify: true, ImportStateVerifyIgnore: []string{"timeouts"}},
			{
				Config: testConfigHelper(t, "resource_web_filtering_disabled.tf"),
				Check:  check.That(resourceType+".test").Key("state").HasValue("disabled"),
			},
		},
	})
	capture.AssertMatchesFixture(t, "tests/requests/web_filtering_lifecycle.json")
	assertNoLinks(t, m)
}

func TestUnitResourceNetworkFilteringProfilePolicyLink_03_TLSInspection(t *testing.T) {
	m := setupMockEnvironment(t)
	capture := mocks.CaptureGraphRequests(t)
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigHelper(t, "resource_tls_inspection.tf"),
				Check:  check.That(resourceType+".test").Key("policy_type").HasValue("tls_inspection"),
			},
		},
	})
	capture.AssertMatchesFixture(t, "tests/requests/tls_inspection.json")
	assertNoLinks(t, m)
}

func TestUnitResourceNetworkFilteringProfilePolicyLink_04_RemovedOutOfBand(t *testing.T) {
	m := setupMockEnvironment(t)
	capture := mocks.CaptureGraphRequests(t)
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: testConfigHelper(t, "resource_filtering.tf")},
			{
				PreConfig: func() {
					m.Lock()
					defer m.Unlock()
					clear(m.links)
				},
				Config:             testConfigHelper(t, "resource_filtering.tf"),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
	capture.AssertMatchesFixture(t, "tests/requests/removed_out_of_band.json")
}
