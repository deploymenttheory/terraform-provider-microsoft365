package graphBetaNetworkMCPPolicy_test

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/jarcoal/httpmock"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	target "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_beta/network_mcp_policy"
	policyMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_beta/network_mcp_policy/mocks"
)

const resourceType = target.ResourceName

func setupMockEnvironment(t *testing.T) *policyMocks.MCPPolicyMock {
	mocks.SetupUnitTestEnvironment(t)
	httpmock.Activate()
	mocks.NewMocks().AuthMocks.RegisterMocks()
	m := &policyMocks.MCPPolicyMock{}
	m.RegisterMocks()
	t.Cleanup(func() { httpmock.DeactivateAndReset(); m.CleanupMockState() })
	return m
}

func testConfigHelper(filename string) string {
	v, e := helpers.ParseHCLFile("tests/terraform/unit/" + filename)
	if e != nil {
		panic(e)
	}
	return v
}

func TestUnitResourceNetworkMCPPolicy_01_Lifecycle(t *testing.T) {
	m := setupMockEnvironment(t)
	var savedID string
	sameID := func(s *terraform.State) error {
		id := s.RootModule().Resources[resourceType+".test"].Primary.ID
		if savedID == "" {
			savedID = id
		} else if id != savedID {
			return fmt.Errorf("in-place update replaced resource: %s -> %s", savedID, id)
		}
		return nil
	}
	resource.UnitTest(
		t,
		resource.TestCase{
			ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: testConfigHelper("resource.tf"),
					Check: resource.ComposeTestCheckFunc(
						sameID,
						check.That(resourceType+".test").Key("id").Exists(),
					),
				},
				{ResourceName: resourceType + ".test", ImportState: true, ImportStateVerify: true},
				{
					Config: testConfigHelper("resource_updated.tf"),
					Check: resource.ComposeTestCheckFunc(
						sameID,
						check.That(resourceType+".test").Key("version").HasValue("1.0.0"),
						check.That(resourceType+".test").Key("description").DoesNotExist(),
					),
				},
				{
					Config: testConfigHelper("resource_empty_description.tf"),
					Check: resource.ComposeTestCheckFunc(
						sameID,
						check.That(resourceType+".test").Key("description").HasValue(""),
					),
				},
			},
		},
	)
	m.Lock()
	defer m.Unlock()
	if len(m.Policies) != 0 || len(m.Rules) != 0 {
		t.Fatal("destroy left policies or rules")
	}
	for _, r := range m.Requests {
		if r.Method == "PATCH" || r.Method == "POST" {
			if _, ok := r.Body["policyRules"]; ok {
				t.Fatal("inline rules sent")
			}
		}
	}
}

func TestUnitResourceNetworkMCPPolicy_02_InvalidAction(t *testing.T) {
	setupMockEnvironment(t)
	resource.UnitTest(
		t,
		resource.TestCase{
			ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config:      invalidActionConfig(),
					ExpectError: regexp.MustCompile(`value must be one of`),
				},
			},
		},
	)
}

func invalidActionConfig() string {
	return strings.ReplaceAll(testConfigHelper("resource.tf"), "\"allow\"", "\"invalid\"")
}

func TestUnitResourceNetworkMCPPolicy_03_Drift(t *testing.T) {
	m := setupMockEnvironment(t)
	resource.UnitTest(
		t,
		resource.TestCase{
			ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{Config: testConfigHelper("resource.tf")},
				{PreConfig: func() {
					m.Lock()
					defer m.Unlock()
					for _, policy := range m.Policies {
						policy["name"] = "out-of-band-name"
					}
				}, Config: testConfigHelper("resource.tf"), PlanOnly: true, ExpectNonEmptyPlan: true},
				{
					Config: testConfigHelper("resource.tf"),
					Check: check.That(resourceType + ".test").
						Key("name").
						HasValue("tf-api-probe-policy"),
				},
			},
		},
	)
}
