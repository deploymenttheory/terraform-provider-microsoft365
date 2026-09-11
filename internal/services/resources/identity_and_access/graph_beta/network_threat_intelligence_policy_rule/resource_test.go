package graphBetaNetworkThreatIntelligencePolicyRule_test

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/jarcoal/httpmock"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	policyMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_beta/network_threat_intelligence_policy/mocks"
	target "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_beta/network_threat_intelligence_policy_rule"
)

const resourceType = target.ResourceName

func setupMockEnvironment(t *testing.T) *policyMocks.ThreatIntelligencePolicyMock {
	mocks.SetupUnitTestEnvironment(t)
	httpmock.Activate()
	mocks.NewMocks().AuthMocks.RegisterMocks()
	m := &policyMocks.ThreatIntelligencePolicyMock{}
	m.RegisterMocks()
	t.Cleanup(func() { httpmock.DeactivateAndReset(); m.CleanupMockState() })
	return m
}

func TestUnitResourceNetworkThreatIntelligencePolicyRule_01_Lifecycle(t *testing.T) {
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
					Config: mocks.LoadUnitTerraformConfig("resource.tf"),
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
						return v.Primary.Attributes["threat_intelligence_policy_id"] + "/" + v.Primary.ID, nil
					},
				},
				{
					Config: mocks.LoadUnitTerraformConfig("resource_updated.tf"),
					Check: resource.ComposeTestCheckFunc(
						sameID,
						check.That(resourceType+".test").Key("enabled").HasValue("false"),
						check.That(resourceType+".test").Key("status").HasValue("disabled"),
						check.That(resourceType+".test").Key("priority").HasValue("65001"),
						check.That(resourceType+".test").Key("description").DoesNotExist(),
					),
				},
				{
					Config: mocks.LoadUnitTerraformConfig("resource_empty_description.tf"),
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

func TestUnitResourceNetworkThreatIntelligencePolicyRule_02_InvalidAction(t *testing.T) {
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
	return regexp.MustCompile(`(?m)^(\s*action\s*=\s*)"allow"`).
		ReplaceAllString(mocks.LoadUnitTerraformConfig("resource.tf"), `${1}"invalid"`)
}

func TestUnitResourceNetworkThreatIntelligencePolicyRule_03_Drift(t *testing.T) {
	m := setupMockEnvironment(t)
	resource.UnitTest(
		t,
		resource.TestCase{
			ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{Config: mocks.LoadUnitTerraformConfig("resource.tf")},
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
				}, Config: mocks.LoadUnitTerraformConfig("resource.tf"), PlanOnly: true, ExpectNonEmptyPlan: true},
				{
					Config: mocks.LoadUnitTerraformConfig("resource.tf"),
					Check:  check.That(resourceType + ".test").Key("enabled").HasValue("true"),
				},
			},
		},
	)
}

func TestUnitResourceNetworkThreatIntelligencePolicyRule_04_PriorityAndDestinationValidation(
	t *testing.T,
) {
	for _, tc := range []struct{ name, old, new, pattern string }{
		{"priority-too-low", "1000", "99", `value must be at least 100`},
		{"priority-overflow", "1000", "2147483648", `2147483647|32-bit|32 bit|Int32|int32`},
		{"invalid-destination", "fqdn", "url", `value must be one of`},
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
								mocks.LoadUnitTerraformConfig("resource.tf"),
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

func TestUnitResourceNetworkThreatIntelligencePolicyRule_05_PriorityConflict(t *testing.T) {
	setupMockEnvironment(t)
	resource.UnitTest(
		t,
		resource.TestCase{
			ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: strings.ReplaceAll(
						mocks.LoadUnitTerraformConfig("resource.tf"),
						"1000",
						"65000",
					),
					ExpectError: regexp.MustCompile(`65000`),
				},
			},
		},
	)
}

func TestUnitResourceNetworkThreatIntelligencePolicyRule_06_ParentReplacement(t *testing.T) {
	m := setupMockEnvironment(t)
	const policyType = "microsoft365_graph_beta_identity_and_access_network_threat_intelligence_policy"
	extraPolicy := `resource "microsoft365_graph_beta_identity_and_access_network_threat_intelligence_policy" "other" {
 name = "tf-api-probe-other-parent"
 default_action = "allow"
 }`
	original := mocks.LoadUnitTerraformConfig("resource.tf") + "\n" + extraPolicy
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
					oldParent = v.Attributes["threat_intelligence_policy_id"]
					return nil
				}},
				{Config: moved, Check: func(s *terraform.State) error {
					v := s.RootModule().Resources[resourceType+".test"].Primary
					if oldID == v.ID || oldParent == v.Attributes["threat_intelligence_policy_id"] {
						return fmt.Errorf("changing the parent must replace the nested rule")
					}
					m.Lock()
					defer m.Unlock()
					if len(m.Policies) != 2 || len(m.Rules[oldParent]) != 1 {
						return fmt.Errorf(
							"replacement modified an unrelated parent or its automatic default rule",
						)
					}
					return nil
				}},
			},
		},
	)
}

func TestUnitResourceNetworkThreatIntelligencePolicyRule_06_ParentReplacementWithNewParent(
	t *testing.T,
) {
	m := setupMockEnvironment(t)
	const policyType = "microsoft365_graph_beta_identity_and_access_network_threat_intelligence_policy"
	extraPolicy := `resource "microsoft365_graph_beta_identity_and_access_network_threat_intelligence_policy" "other" {
 name = "tf-api-probe-other-parent"
 default_action = "allow"
 }`
	original := mocks.LoadUnitTerraformConfig("resource.tf")
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
					oldParent = v.Attributes["threat_intelligence_policy_id"]
					return nil
				}},
				{Config: moved, Check: func(s *terraform.State) error {
					v := s.RootModule().Resources[resourceType+".test"].Primary
					if oldID == v.ID || oldParent == v.Attributes["threat_intelligence_policy_id"] {
						return fmt.Errorf("changing the parent must replace the nested rule")
					}
					m.Lock()
					defer m.Unlock()
					if len(m.Policies) != 2 || len(m.Rules[oldParent]) != 1 {
						return fmt.Errorf(
							"replacement modified an unrelated parent or its automatic default rule",
						)
					}
					return nil
				}},
			},
		},
	)
}

func TestUnitResourceNetworkThreatIntelligencePolicyRule_07_EmptyDestinations(t *testing.T) {
	setupMockEnvironment(t)
	original := mocks.LoadUnitTerraformConfig("resource.tf")
	start := regexp.MustCompile(`(?m)^\s*destinations\s*=`).FindStringIndex(original)[0]
	empty := original[:start] + "destinations = []\n}\n"
	resource.UnitTest(
		t,
		resource.TestCase{
			ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{Config: original},
				{
					Config: empty,
					Check:  check.That(resourceType + ".test").Key("destinations.#").HasValue("0"),
				},
				{Config: empty, PlanOnly: true},
				{
					Config: original,
					Check: check.That(resourceType + ".test").
						Key("destinations.0.values.#").
						HasValue("3"),
				},
			},
		},
	)
}
