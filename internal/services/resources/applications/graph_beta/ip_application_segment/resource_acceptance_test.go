package graphBetaApplicationsIpApplicationSegment_test

import (
	"regexp"
	"testing"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/destroy"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/testlog"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaApplicationsIpApplicationSegment "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/applications/graph_beta/ip_application_segment"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

var (
	// Resource type name from the resource package
	resourceType = graphBetaApplicationsIpApplicationSegment.ResourceName

	// testResource is the test resource implementation for IP application segments
	testResource = graphBetaApplicationsIpApplicationSegment.IpApplicationSegmentTestResource{}
)

// Helper function to load test configs from acceptance directory
func loadAcceptanceTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/acceptance/" + filename)
	if err != nil {
		panic("failed to load acceptance config " + filename + ": " + err.Error())
	}
	return config
}

// TestAccIpApplicationSegmentResource_Minimal tests the minimal IP application segment configuration
func TestAccIpApplicationSegmentResource_Minimal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			resourceType,
			30*time.Second,
		),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating minimal IP application segment")
				},
				Config: loadAcceptanceTestTerraform("resource_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("IP application segment", 10*time.Second)
						time.Sleep(10 * time.Second)
						return nil
					},
					check.That(resourceType+".ip_segment_minimal").ExistsInGraph(testResource),
					check.That(resourceType+".ip_segment_minimal").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".ip_segment_minimal").Key("application_object_id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".ip_segment_minimal").Key("destination_host").HasValue("192.168.1.100"),
					check.That(resourceType+".ip_segment_minimal").Key("destination_type").HasValue("ipAddress"),
					check.That(resourceType+".ip_segment_minimal").Key("protocol").HasValue("tcp"),

					// Ports
					check.That(resourceType+".ip_segment_minimal").Key("ports.#").HasValue("1"),
					check.That(resourceType+".ip_segment_minimal").Key("ports.*").ContainsTypeSetElement("80-80"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing minimal IP application segment")
				},
				ResourceName:            resourceType + ".ip_segment_minimal",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// TestAccIpApplicationSegmentResource_Maximal tests the maximal IP application segment configuration
func TestAccIpApplicationSegmentResource_Maximal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			resourceType,
			30*time.Second,
		),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating maximal IP application segment")
				},
				Config: loadAcceptanceTestTerraform("resource_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("IP application segment", 10*time.Second)
						time.Sleep(10 * time.Second)
						return nil
					},
					check.That(resourceType+".ip_segment_maximal").ExistsInGraph(testResource),
					check.That(resourceType+".ip_segment_maximal").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".ip_segment_maximal").Key("application_object_id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".ip_segment_maximal").Key("destination_host").HasValue("*.example.com"),
					check.That(resourceType+".ip_segment_maximal").Key("destination_type").HasValue("dnsSuffix"),
					check.That(resourceType+".ip_segment_maximal").Key("protocol").HasValue("tcp"),

					// Ports
					check.That(resourceType+".ip_segment_maximal").Key("ports.#").HasValue("4"),
					check.That(resourceType+".ip_segment_maximal").Key("ports.*").ContainsTypeSetElement("80-80"),
					check.That(resourceType+".ip_segment_maximal").Key("ports.*").ContainsTypeSetElement("443-443"),
					check.That(resourceType+".ip_segment_maximal").Key("ports.*").ContainsTypeSetElement("8080-8080"),
					check.That(resourceType+".ip_segment_maximal").Key("ports.*").ContainsTypeSetElement("8443-8443"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing maximal IP application segment")
				},
				ResourceName:            resourceType + ".ip_segment_maximal",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// TestAccIpApplicationSegmentResource_IpRange tests IP range configuration
func TestAccIpApplicationSegmentResource_IpRange(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			resourceType,
			30*time.Second,
		),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating IP range application segment")
				},
				Config: loadAcceptanceTestTerraform("resource_ip_range.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("IP application segment", 10*time.Second)
						time.Sleep(10 * time.Second)
						return nil
					},
					check.That(resourceType+".ip_segment_range").ExistsInGraph(testResource),
					check.That(resourceType+".ip_segment_range").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".ip_segment_range").Key("destination_host").HasValue("192.168.1.0/24"),
					check.That(resourceType+".ip_segment_range").Key("destination_type").HasValue("ipRangeCidr"),
					check.That(resourceType+".ip_segment_range").Key("protocol").HasValue("tcp"),
					check.That(resourceType+".ip_segment_range").Key("ports.#").HasValue("1"),
					check.That(resourceType+".ip_segment_range").Key("ports.*").ContainsTypeSetElement("443-443"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing IP range application segment")
				},
				ResourceName:            resourceType + ".ip_segment_range",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// TestAccIpApplicationSegmentResource_FQDN tests FQDN configuration
func TestAccIpApplicationSegmentResource_FQDN(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			resourceType,
			30*time.Second,
		),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating FQDN application segment")
				},
				Config: loadAcceptanceTestTerraform("resource_fqdn.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("IP application segment", 10*time.Second)
						time.Sleep(10 * time.Second)
						return nil
					},
					check.That(resourceType+".ip_segment_fqdn").ExistsInGraph(testResource),
					check.That(resourceType+".ip_segment_fqdn").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".ip_segment_fqdn").Key("destination_host").HasValue("app.example.com"),
					check.That(resourceType+".ip_segment_fqdn").Key("destination_type").HasValue("fqdn"),
					check.That(resourceType+".ip_segment_fqdn").Key("protocol").HasValue("tcp"),
					check.That(resourceType+".ip_segment_fqdn").Key("ports.#").HasValue("2"),
					check.That(resourceType+".ip_segment_fqdn").Key("ports.*").ContainsTypeSetElement("443-443"),
					check.That(resourceType+".ip_segment_fqdn").Key("ports.*").ContainsTypeSetElement("8443-8443"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing FQDN application segment")
				},
				ResourceName:            resourceType + ".ip_segment_fqdn",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}
