package graphBetaNetworkMCPPolicyRule_test

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
	policyMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_beta/network_mcp_policy/mocks"
	target "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_beta/network_mcp_policy_rule"
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

func TestUnitResourceNetworkMCPPolicyRule_01_Lifecycle(t *testing.T) {
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
				{
					ResourceName:      resourceType + ".test",
					ImportState:       true,
					ImportStateVerify: true,
					ImportStateIdFunc: func(s *terraform.State) (string, error) {
						v := s.RootModule().Resources[resourceType+".test"]
						return v.Primary.Attributes["mcp_policy_id"] + "/" + v.Primary.ID, nil
					},
				},
				{
					Config: testConfigHelper("resource_updated.tf"),
					Check: resource.ComposeTestCheckFunc(
						sameID,
						check.That(resourceType+".test").Key("enabled").HasValue("false"),
						check.That(resourceType+".test").Key("status").HasValue("disabled"),
						check.That(resourceType+".test").Key("priority").HasValue("65001"),
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
				{
					Config: testConfigHelper("resource_no_conditions.tf"),
					Check: resource.ComposeTestCheckFunc(
						sameID,
						check.That(resourceType+".test").Key("enabled").HasValue("true"),
						check.That(resourceType+".test").Key("matching_conditions").DoesNotExist(),
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

func TestUnitResourceNetworkMCPPolicyRule_02_InvalidAction(t *testing.T) {
	setupMockEnvironment(t)
	resource.UnitTest(
		t,
		resource.TestCase{
			ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config:      invalidActionConfig(),
					ExpectError: regexp.MustCompile(`value must be\s+one of`),
				},
			},
		},
	)
}

func invalidActionConfig() string {
	return strings.ReplaceAll(testConfigHelper("resource.tf"), "\"allow\"", "\"invalid\"")
}

func TestUnitResourceNetworkMCPPolicyRule_03_Drift(t *testing.T) {
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
					for _, rules := range m.Rules {
						for _, rule := range rules {
							if rule["name"] == "tf-api-probe-rule" {
								rule["settings"] = map[string]any{"status": "disabled"}
							}
						}
					}
				}, Config: testConfigHelper("resource.tf"), PlanOnly: true, ExpectNonEmptyPlan: true},
				{
					Config: testConfigHelper("resource.tf"),
					Check:  check.That(resourceType + ".test").Key("enabled").HasValue("true"),
				},
			},
		},
	)
}

func TestUnitResourceNetworkMCPPolicyRule_04_PriorityAndConditionValidation(t *testing.T) {
	for _, tc := range []struct{ name, old, new, pattern string }{
		{"priority-too-low", "1000", "99", `value must be at least 100`},
		{"priority-overflow", "1000", "2147483648", `2147483647|32-bit|32 bit|Int32|int32`},
		{"invalid-match", "exactMatch", "regex", `value must be\s+one of`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			setupMockEnvironment(t)
			resource.UnitTest(
				t,
				resource.TestCase{
					ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
					Steps: []resource.TestStep{
						{
							Config: strings.ReplaceAll(
								testConfigHelper("resource.tf"),
								tc.old,
								tc.new,
							),
							ExpectError: regexp.MustCompile(tc.pattern),
						},
					},
				},
			)
		})
	}
}

func TestUnitResourceNetworkMCPPolicyRule_05_PriorityConflict(t *testing.T) {
	setupMockEnvironment(t)
	config := testConfigHelper("resource.tf") + `
resource "microsoft365_graph_beta_identity_and_access_network_mcp_policy_rule" "conflict" {
 mcp_policy_id = microsoft365_graph_beta_identity_and_access_network_mcp_policy.test.id
 name = "tf-api-probe-conflict"
 priority = 1000
 action = "block"
 enabled = false
 depends_on = [microsoft365_graph_beta_identity_and_access_network_mcp_policy_rule.test]
 }`
	resource.UnitTest(
		t,
		resource.TestCase{
			ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{Config: config, ExpectError: regexp.MustCompile("priority already exists")},
			},
		},
	)
}

func TestUnitResourceNetworkMCPPolicyRule_06_ParentReplacement(t *testing.T) {
	m := setupMockEnvironment(t)
	const policyType = "microsoft365_graph_beta_identity_and_access_network_mcp_policy"
	extraPolicy := `resource "microsoft365_graph_beta_identity_and_access_network_mcp_policy" "other" {
 name = "tf-api-probe-other-parent"
 default_action = "allow"
 }`
	original := testConfigHelper("resource.tf") + "\n" + extraPolicy
	moved := strings.ReplaceAll(original, policyType+".test.id", policyType+".other.id")
	var oldID, oldParent string
	resource.UnitTest(
		t,
		resource.TestCase{
			ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{Config: original, Check: func(s *terraform.State) error {
					v := s.RootModule().Resources[resourceType+".test"].Primary
					oldID = v.ID
					oldParent = v.Attributes["mcp_policy_id"]
					return nil
				}},
				{Config: moved, Check: func(s *terraform.State) error {
					v := s.RootModule().Resources[resourceType+".test"].Primary
					if oldID == v.ID || oldParent == v.Attributes["mcp_policy_id"] {
						return fmt.Errorf("changing the parent must replace the nested rule")
					}
					m.Lock()
					defer m.Unlock()
					if len(m.Policies) != 2 || len(m.Rules[oldParent]) != 0 {
						return fmt.Errorf("replacement modified an unrelated parent")
					}
					return nil
				}},
			},
		},
	)
}

func TestUnitResourceNetworkMCPPolicyRule_06_ParentReplacementWithNewParent(t *testing.T) {
	m := setupMockEnvironment(t)
	const policyType = "microsoft365_graph_beta_identity_and_access_network_mcp_policy"
	extraPolicy := `resource "microsoft365_graph_beta_identity_and_access_network_mcp_policy" "other" {
 name = "tf-api-probe-other-parent"
 default_action = "allow"
 }`
	original := testConfigHelper("resource.tf")
	moved := strings.ReplaceAll(
		original,
		policyType+".test.id",
		policyType+".other.id",
	) + "\n" + extraPolicy
	var oldID, oldParent string
	resource.UnitTest(
		t,
		resource.TestCase{
			ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{Config: original, Check: func(s *terraform.State) error {
					v := s.RootModule().Resources[resourceType+".test"].Primary
					oldID = v.ID
					oldParent = v.Attributes["mcp_policy_id"]
					return nil
				}},
				{Config: moved, Check: func(s *terraform.State) error {
					v := s.RootModule().Resources[resourceType+".test"].Primary
					if oldID == v.ID || oldParent == v.Attributes["mcp_policy_id"] {
						return fmt.Errorf("changing the parent must replace the nested rule")
					}
					m.Lock()
					defer m.Unlock()
					if len(m.Policies) != 2 || len(m.Rules[oldParent]) != 0 {
						return fmt.Errorf("replacement modified an unrelated parent")
					}
					return nil
				}},
			},
		},
	)
}
