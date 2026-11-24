package graphBetaNetworkForwardingProfilePolicyLink_test

import (
	"encoding/json"
	"net/http"
	"regexp"
	"sync"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks/factories"
	target "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_beta/network_forwarding_profile_policy_link"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

const (
	resourceType = target.ResourceName
	profileID    = "11111111-2222-3333-4444-555555555555"
	policyLinkID = "66666666-7777-8888-9999-000000000000"
)

// linkMock is a stateful mock of the Microsoft-managed /networkAccess/forwardingProfiles/{id}/policies links.
type linkMock struct {
	sync.Mutex
	links map[string]map[string]any
}

func setupMockEnvironment(t *testing.T) *linkMock {
	t.Helper()
	mocks.SetupUnitTestEnvironment(t)
	httpmock.Activate()
	t.Cleanup(httpmock.DeactivateAndReset)
	mocks.NewMocks().AuthMocks.RegisterMocks()

	m := &linkMock{links: map[string]map[string]any{
		policyLinkID: {
			"@odata.type": "#microsoft.graph.networkaccess.forwardingPolicyLink",
			"id":          policyLinkID,
			"state":       "enabled",
			"version":     "1.0.0",
			"priority":    4,
			"policy": map[string]any{
				"@odata.type":           "#microsoft.graph.networkaccess.forwardingPolicy",
				"id":                    "f0474b3e-307a-4230-bc1c-cd8ac2f1a2cf",
				"name":                  "Default Acquire",
				"description":           "Internet Access default acquire forwarding policy",
				"trafficForwardingType": "internet",
			},
		},
	}}
	item := regexp.MustCompile(`^https://graph\.microsoft\.com/beta/networkAccess/forwardingProfiles/([^/]+)/policies/([^/?]+)$`)

	httpmock.RegisterRegexpResponder(http.MethodGet, item, func(req *http.Request) (*http.Response, error) {
		m.Lock()
		defer m.Unlock()
		link, ok := m.links[item.FindStringSubmatch(req.URL.String())[2]]
		if !ok {
			return factories.ErrorResponse(http.StatusNotFound, "ResourceNotFound", "policy link not found")(req)
		}
		return factories.SuccessResponse(http.StatusOK, link)(req)
	})
	httpmock.RegisterRegexpResponder(http.MethodPatch, item, func(req *http.Request) (*http.Response, error) {
		var body map[string]any
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			return httpmock.NewStringResponse(http.StatusBadRequest, `{"error":{"code":"BadRequest","message":"invalid JSON"}}`), nil
		}
		m.Lock()
		defer m.Unlock()
		link, ok := m.links[item.FindStringSubmatch(req.URL.String())[2]]
		if !ok {
			return factories.ErrorResponse(http.StatusNotFound, "ResourceNotFound", "policy link not found")(req)
		}
		link["state"] = body["state"]
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

func TestUnitResourceNetworkForwardingProfilePolicyLink_01_Lifecycle(t *testing.T) {
	m := setupMockEnvironment(t)
	capture := mocks.CaptureGraphRequests(t)
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigHelper(t, "resource_disabled.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("id").HasValue(profileID+"/"+policyLinkID),
					check.That(resourceType+".test").Key("state").HasValue("disabled"),
					check.That(resourceType+".test").Key("priority").HasValue("4"),
					check.That(resourceType+".test").Key("policy_name").HasValue("Default Acquire"),
					check.That(resourceType+".test").Key("traffic_forwarding_type").HasValue("internet"),
				),
			},
			{ResourceName: resourceType + ".test", ImportState: true, ImportStateVerify: true, ImportStateVerifyIgnore: []string{"timeouts"}},
			{
				Config: testConfigHelper(t, "resource_enabled.tf"),
				Check:  check.That(resourceType+".test").Key("state").HasValue("enabled"),
			},
		},
	})
	capture.AssertMatchesFixture(t, "tests/requests/lifecycle.json")
	m.Lock()
	defer m.Unlock()
	if len(m.links) != 1 {
		t.Fatal("destroy deleted the Microsoft-managed policy link")
	}
}

func TestUnitResourceNetworkForwardingProfilePolicyLink_02_RemovedOutOfBand(t *testing.T) {
	m := setupMockEnvironment(t)
	capture := mocks.CaptureGraphRequests(t)
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: testConfigHelper(t, "resource_enabled.tf")},
			{
				PreConfig: func() {
					m.Lock()
					defer m.Unlock()
					clear(m.links)
				},
				Config:             testConfigHelper(t, "resource_enabled.tf"),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
	capture.AssertMatchesFixture(t, "tests/requests/removed_out_of_band.json")
}
