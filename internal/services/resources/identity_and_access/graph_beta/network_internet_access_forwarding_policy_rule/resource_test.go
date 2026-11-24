package graphBetaNetworkInternetAccessForwardingPolicyRule_test

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
	target "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_beta/network_internet_access_forwarding_policy_rule"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/jarcoal/httpmock"
)

const resourceType = target.ResourceName

// ruleMock is a stateful mock of /networkAccess/forwardingPolicies/{id}/policyRules.
type ruleMock struct {
	sync.Mutex
	rules map[string]map[string]any
	next  int
	// preconditionFailures makes the next write requests fail with 412 this many times.
	preconditionFailures int
}

func setupMockEnvironment(t *testing.T) *ruleMock {
	t.Helper()
	mocks.SetupUnitTestEnvironment(t)
	httpmock.Activate()
	t.Cleanup(httpmock.DeactivateAndReset)
	mocks.NewMocks().AuthMocks.RegisterMocks()

	m := &ruleMock{rules: map[string]map[string]any{}}
	collection := regexp.MustCompile(`^https://graph\.microsoft\.com/beta/networkAccess/forwardingPolicies/([^/]+)/policyRules$`)
	item := regexp.MustCompile(`^https://graph\.microsoft\.com/beta/networkAccess/forwardingPolicies/([^/]+)/policyRules/([^/?]+)$`)
	preconditionFailed := func(req *http.Request) (*http.Response, bool) {
		if m.preconditionFailures == 0 {
			return nil, false
		}
		m.preconditionFailures--
		response, _ := factories.ErrorResponse(http.StatusPreconditionFailed, "PreconditionFailed", "The policy changed")(req)
		return response, true
	}
	decode := func(req *http.Request) (map[string]any, *http.Response) {
		var body map[string]any
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			return nil, httpmock.NewStringResponse(http.StatusBadRequest, `{"error":{"code":"BadRequest","message":"invalid JSON"}}`)
		}
		return body, nil
	}

	httpmock.RegisterRegexpResponder(http.MethodPost, collection, func(req *http.Request) (*http.Response, error) {
		body, invalid := decode(req)
		if invalid != nil {
			return invalid, nil
		}
		m.Lock()
		defer m.Unlock()
		if response, failed := preconditionFailed(req); failed {
			return response, nil
		}
		m.next++
		body["id"] = fmt.Sprintf("00000000-0000-0000-0000-%012d", m.next)
		if _, ok := body["clientFallbackAction"]; !ok {
			body["clientFallbackAction"] = "default"
		}
		m.rules[body["id"].(string)] = body
		return factories.SuccessResponse(http.StatusCreated, body)(req)
	})
	httpmock.RegisterRegexpResponder(http.MethodGet, item, func(req *http.Request) (*http.Response, error) {
		m.Lock()
		defer m.Unlock()
		rule, ok := m.rules[item.FindStringSubmatch(req.URL.String())[2]]
		if !ok {
			return factories.ErrorResponse(http.StatusNotFound, "ResourceNotFound", "rule not found")(req)
		}
		return factories.SuccessResponse(http.StatusOK, rule)(req)
	})
	httpmock.RegisterRegexpResponder(http.MethodPatch, item, func(req *http.Request) (*http.Response, error) {
		body, invalid := decode(req)
		if invalid != nil {
			return invalid, nil
		}
		m.Lock()
		defer m.Unlock()
		if response, failed := preconditionFailed(req); failed {
			return response, nil
		}
		rule, ok := m.rules[item.FindStringSubmatch(req.URL.String())[2]]
		if !ok {
			return factories.ErrorResponse(http.StatusNotFound, "ResourceNotFound", "rule not found")(req)
		}
		for key, value := range body {
			rule[key] = value
		}
		return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
	})
	httpmock.RegisterRegexpResponder(http.MethodDelete, item, func(req *http.Request) (*http.Response, error) {
		m.Lock()
		defer m.Unlock()
		if response, failed := preconditionFailed(req); failed {
			return response, nil
		}
		delete(m.rules, item.FindStringSubmatch(req.URL.String())[2])
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

func importID(s *terraform.State) (string, error) {
	rule := s.RootModule().Resources[resourceType+".test"]
	return rule.Primary.Attributes["forwarding_policy_id"] + "/" + rule.Primary.ID, nil
}

func TestUnitResourceNetworkInternetAccessForwardingPolicyRule_01_Lifecycle(t *testing.T) {
	m := setupMockEnvironment(t)
	capture := mocks.CaptureGraphRequests(t)
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigHelper(t, "resource_fqdn.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("id").Exists(),
					check.That(resourceType+".test").Key("rule_type").HasValue("fqdn"),
					check.That(resourceType+".test").Key("client_fallback_action").HasValue("block"),
					check.That(resourceType+".test").Key("destinations.0.value").HasValue("www.example.com"),
				),
			},
			{ResourceName: resourceType + ".test", ImportState: true, ImportStateVerify: true, ImportStateIdFunc: importID, ImportStateVerifyIgnore: []string{"timeouts"}},
			{
				Config: testConfigHelper(t, "resource_fqdn_updated.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("protocol").HasValue("udp"),
					check.That(resourceType+".test").Key("destinations.#").HasValue("2"),
				),
			},
			{
				Config: testConfigHelper(t, "resource_ip_range.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("rule_type").HasValue("ip_range"),
					check.That(resourceType+".test").Key("destinations.0.begin_address").HasValue("192.0.2.10"),
					check.That(resourceType+".test").Key("destinations.0.end_address").HasValue("192.0.2.20"),
				),
			},
			{
				Config: testConfigHelper(t, "resource_ip_subnet.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("rule_type").HasValue("ip_subnet"),
					check.That(resourceType+".test").Key("client_fallback_action").HasValue("block"),
				),
			},
			{
				Config: testConfigHelper(t, "resource_ip_address.tf"),
				Check:  check.That(resourceType+".test").Key("rule_type").HasValue("ip_address"),
			},
		},
	})
	capture.AssertMatchesFixture(t, "tests/requests/lifecycle.json")
	m.Lock()
	defer m.Unlock()
	if len(m.rules) != 0 {
		t.Fatalf("destroy left %d rules", len(m.rules))
	}
}

func TestUnitResourceNetworkInternetAccessForwardingPolicyRule_02_PreconditionRetry(t *testing.T) {
	m := setupMockEnvironment(t)
	m.preconditionFailures = 1
	capture := mocks.CaptureGraphRequests(t)
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigHelper(t, "resource_ip_address.tf"),
				Check:  check.That(resourceType+".test").Key("client_fallback_action").HasValue("default"),
			},
			{
				PreConfig: func() {
					m.Lock()
					defer m.Unlock()
					m.preconditionFailures = 1
				},
				Config:  testConfigHelper(t, "resource_ip_address.tf"),
				Destroy: true,
			},
		},
	})
	capture.AssertMatchesFixture(t, "tests/requests/precondition_retry.json")
	m.Lock()
	defer m.Unlock()
	if len(m.rules) != 0 {
		t.Fatalf("destroy left %d rules", len(m.rules))
	}
}
