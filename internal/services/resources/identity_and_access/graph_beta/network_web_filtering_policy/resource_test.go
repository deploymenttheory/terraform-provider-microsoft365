package graphBetaNetworkWebFilteringPolicy_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"sync"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks/factories"
	target "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_beta/network_web_filtering_policy"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

const resourceType = target.ResourceName

const policiesURL = "https://graph.microsoft.com/beta/networkaccess/webFilteringPolicies"

// policyMock is a stateful mock of /networkaccess/webFilteringPolicies.
type policyMock struct {
	sync.Mutex
	policies map[string]map[string]any
	next     int
}

func setupMockEnvironment(t *testing.T) *policyMock {
	t.Helper()
	mocks.SetupUnitTestEnvironment(t)
	httpmock.Activate()
	t.Cleanup(httpmock.DeactivateAndReset)
	mocks.NewMocks().AuthMocks.RegisterMocks()

	m := &policyMock{policies: map[string]map[string]any{}}
	item := regexp.MustCompile(`^` + regexp.QuoteMeta(policiesURL) + `/([^/?]+)$`)

	httpmock.RegisterResponder(http.MethodPost, policiesURL, func(req *http.Request) (*http.Response, error) {
		var body map[string]any
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			return httpmock.NewStringResponse(http.StatusBadRequest, `{"error":{"code":"BadRequest","message":"invalid JSON"}}`), nil
		}
		m.Lock()
		defer m.Unlock()
		m.next++
		id := fmt.Sprintf("00000000-0000-0000-0000-%012d", m.next)
		policy := map[string]any{"id": id, "name": body["name"], "description": body["description"], "settings": body["settings"]}
		m.policies[id] = policy
		return factories.SuccessResponse(http.StatusCreated, policy)(req)
	})
	httpmock.RegisterRegexpResponder(http.MethodGet, item, func(req *http.Request) (*http.Response, error) {
		m.Lock()
		defer m.Unlock()
		policy, ok := m.policies[item.FindStringSubmatch(req.URL.String())[1]]
		if !ok {
			return factories.ErrorResponse(http.StatusNotFound, "ResourceNotFound", "web filtering policy not found")(req)
		}
		return factories.SuccessResponse(http.StatusOK, policy)(req)
	})
	httpmock.RegisterRegexpResponder(http.MethodPatch, item, func(req *http.Request) (*http.Response, error) {
		var body map[string]any
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			return httpmock.NewStringResponse(http.StatusBadRequest, `{"error":{"code":"BadRequest","message":"invalid JSON"}}`), nil
		}
		m.Lock()
		defer m.Unlock()
		policy, ok := m.policies[item.FindStringSubmatch(req.URL.String())[1]]
		if !ok {
			return factories.ErrorResponse(http.StatusNotFound, "ResourceNotFound", "web filtering policy not found")(req)
		}
		for key, value := range body {
			policy[key] = value
		}
		return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
	})
	httpmock.RegisterRegexpResponder(http.MethodDelete, item, func(req *http.Request) (*http.Response, error) {
		m.Lock()
		defer m.Unlock()
		delete(m.policies, item.FindStringSubmatch(req.URL.String())[1])
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

func TestUnitResourceNetworkWebFilteringPolicy_01_Lifecycle(t *testing.T) {
	m := setupMockEnvironment(t)
	capture := mocks.CaptureGraphRequests(t)
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigHelper(t, "resource.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("id").Exists(),
					check.That(resourceType+".test").Key("name").HasValue("unit-test-web-filtering-policy"),
					check.That(resourceType+".test").Key("description").HasValue("Unit test web filtering policy"),
					check.That(resourceType+".test").Key("default_action").HasValue("allow"),
				),
			},
			{ResourceName: resourceType + ".test", ImportState: true, ImportStateVerify: true, ImportStateVerifyIgnore: []string{"timeouts"}},
			{
				Config: testConfigHelper(t, "resource_updated.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("name").HasValue("unit-test-web-filtering-policy-updated"),
					check.That(resourceType+".test").Key("description").HasValue("Updated unit test web filtering policy"),
					check.That(resourceType+".test").Key("default_action").HasValue("block"),
				),
			},
		},
	})
	capture.AssertMatchesFixture(t, "tests/requests/lifecycle.json")
	m.Lock()
	defer m.Unlock()
	if len(m.policies) != 0 {
		t.Fatalf("destroy left %d policies", len(m.policies))
	}
}

func TestUnitResourceNetworkWebFilteringPolicy_02_RemovedOutOfBand(t *testing.T) {
	m := setupMockEnvironment(t)
	capture := mocks.CaptureGraphRequests(t)
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: testConfigHelper(t, "resource.tf")},
			{
				PreConfig: func() {
					m.Lock()
					defer m.Unlock()
					clear(m.policies)
				},
				Config:             testConfigHelper(t, "resource.tf"),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
	capture.AssertMatchesFixture(t, "tests/requests/removed_out_of_band.json")
}
