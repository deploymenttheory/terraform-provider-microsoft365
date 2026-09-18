package graphBetaNetworkCloudFirewallPolicyRule_test

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	cloudMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_beta/network_cloud_firewall_policy/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/jarcoal/httpmock"
)

const resourceType = "microsoft365_graph_beta_identity_and_access_network_cloud_firewall_policy_rule"

func TestUnitResourceNetworkCloudFirewallPolicyRule_01_Lifecycle(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	mocks.NewMocks().AuthMocks.RegisterMocks()
	mock := &cloudMocks.CloudFirewallMock{}
	mock.RegisterMocks()
	defer mock.CleanupMockState()
	config, err := helpers.ParseHCLFile("tests/terraform/unit/resource.tf")
	if err != nil {
		t.Fatal(err)
	}
	updated := strings.ReplaceAll(config, "Managed by Terraform", "Updated description")
	updated = strings.ReplaceAll(updated, "enabled     = false", "enabled     = true")
	noDescription := strings.ReplaceAll(updated, `  description    = "Updated description"
`, "")
	noDescription = strings.ReplaceAll(noDescription, `  description = "Match example destination addresses"
`, "")
	emptyDescriptions := strings.ReplaceAll(updated, "Updated description", "")
	emptyDescriptions = strings.ReplaceAll(emptyDescriptions, "Match example destination addresses", "")
	clearedMatching := regexp.MustCompile(`(?s)
  (sources|destinations) = \{.*?
  \}
`).ReplaceAllString(noDescription, "\n")
	var id string
	resource.UnitTest(t, resource.TestCase{ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories, Steps: []resource.TestStep{
		{Config: config, Check: resource.ComposeTestCheckFunc(check.That(resourceType+".test").Key("id").Exists(), func(s *terraform.State) error {
			id = s.RootModule().Resources[resourceType+".test"].Primary.ID
			return nil
		})},
		{Config: updated, Check: func(s *terraform.State) error {
			if s.RootModule().Resources[resourceType+".test"].Primary.ID != id {
				return fmt.Errorf("update replaced ID")
			}
			return nil
		}},
		{Config: emptyDescriptions, Check: check.That(resourceType + ".test").Key("description").HasValue("")},
		{Config: noDescription, Check: check.That(resourceType + ".test").Key("description").DoesNotExist()},
		{ResourceName: resourceType + ".test", ImportState: true, ImportStateVerify: true, ImportStateIdFunc: func(s *terraform.State) (string, error) {
			v := s.RootModule().Resources[resourceType+".test"].Primary
			return v.Attributes["policy_id"] + "/" + v.ID, nil
		}},
		{ResourceName: resourceType + ".test", ImportState: true, ImportStateKind: resource.ImportBlockWithResourceIdentity, ImportPlanChecks: resource.ImportPlanChecks{PreApply: []plancheck.PlanCheck{plancheck.ExpectResourceAction(resourceType+".test", plancheck.ResourceActionNoop)}}},
		{Config: noDescription, PlanOnly: true},
		{Config: clearedMatching, Check: resource.ComposeTestCheckFunc(
			check.That(resourceType+".test").Key("sources.%").DoesNotExist(),
			check.That(resourceType+".test").Key("destinations.%").DoesNotExist(),
		)},
		{Config: clearedMatching, PlanOnly: true},
	}})
}

func TestUnitResourceNetworkCloudFirewallPolicyRule_04_SetNoOpAndParentReplacement(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	mocks.NewMocks().AuthMocks.RegisterMocks()
	mock := &cloudMocks.CloudFirewallMock{}
	mock.RegisterMocks()
	defer mock.CleanupMockState()
	config, err := helpers.ParseHCLFile("tests/terraform/unit/resource.tf")
	if err != nil {
		t.Fatal(err)
	}
	reordered := strings.ReplaceAll(config, `["198.51.100.1", "198.51.100.2"]`, `["198.51.100.2", "198.51.100.1", "198.51.100.2"]`)
	reordered = strings.ReplaceAll(reordered, `["tcp", "udp"]`, `["udp", "tcp", "tcp"]`)
	withTimeout := strings.Replace(reordered, `"microsoft365_graph_beta_identity_and_access_network_cloud_firewall_policy_rule" "test" {`, `"microsoft365_graph_beta_identity_and_access_network_cloud_firewall_policy_rule" "test" {
  timeouts = { read = "3m" }`, 1)
	replacement := reordered + `
resource "microsoft365_graph_beta_identity_and_access_network_cloud_firewall_policy" "replacement" {
 name = "new-parent"
}
`
	replacement = strings.ReplaceAll(replacement, "policy.test.id", "policy.replacement.id")
	var oldID string
	resource.UnitTest(t, resource.TestCase{ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories, Steps: []resource.TestStep{
		{Config: config, Check: func(s *terraform.State) error {
			oldID = s.RootModule().Resources[resourceType+".test"].Primary.ID
			return nil
		}},
		{Config: reordered, PlanOnly: true, Check: func(_ *terraform.State) error {
			if mock.Count("PATCH") != 0 {
				return fmt.Errorf("set reorder caused PATCH")
			}
			return nil
		}},
		{Config: withTimeout, Check: func(_ *terraform.State) error {
			if mock.Count("PATCH") != 0 {
				return fmt.Errorf("timeout-only change sent PATCH")
			}
			return nil
		}},
		{Config: replacement, Check: func(s *terraform.State) error {
			if s.RootModule().Resources[resourceType+".test"].Primary.ID == oldID {
				return fmt.Errorf("unknown parent change did not replace rule")
			}
			oldID = s.RootModule().Resources[resourceType+".test"].Primary.ID
			return nil
		}},
		{Config: reordered, Check: func(s *terraform.State) error {
			if s.RootModule().Resources[resourceType+".test"].Primary.ID == oldID {
				return fmt.Errorf("known parent change did not replace rule")
			}
			return nil
		}},
	}})
}

func TestUnitResourceNetworkCloudFirewallPolicyRule_05_Validation(t *testing.T) {
	config, err := helpers.ParseHCLFile("tests/terraform/unit/resource.tf")
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct{ name, from, to, message string }{
		{"port_leading_zero", `["443"]`, `["0443"]`, "Invalid port"},
		{"port_range_reversed", `["443"]`, `["443-80"]`, "Invalid port"},
		{"port_out_of_range", `["443"]`, `["65536"]`, "Invalid port"},
		{"ipv6", `["192.0.2.0/24"]`, `["2001:db8::1"]`, "Invalid IPv4 condition"},
		{"cidr_invalid", `["192.0.2.0/24"]`, `["192.0.2.0/33"]`, "Invalid IPv4 condition"},
		{"protocol_invalid", `["tcp", "udp"]`, `["icmp"]`, "must be one of"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			mocks.SetupUnitTestEnvironment(t)
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()
			mocks.NewMocks().AuthMocks.RegisterMocks()
			resource.UnitTest(t, resource.TestCase{ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories, Steps: []resource.TestStep{{Config: strings.ReplaceAll(config, tc.from, tc.to), ExpectError: regexp.MustCompile(tc.message)}}})
		})
	}
}
