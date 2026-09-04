package graphBetaNetworkManagedTLSCertificate_test

import (
	"regexp"
	"testing"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/destroy"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/testlog"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaNetworkManagedTLSCertificate "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_beta/network_managed_tls_certificate"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var (
	resourceType = graphBetaNetworkManagedTLSCertificate.ResourceName
	testResource = graphBetaNetworkManagedTLSCertificate.NetworkManagedTLSCertificateTestResource{}
)

func loadAcceptanceTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/acceptance/" + filename)
	if err != nil {
		panic("failed to load acceptance config " + filename + ": " + err.Error())
	}
	return acceptance.ConfiguredM365ProviderBlock(config)
}

func TestAccResourceNetworkManagedTLSCertificate_01_Lifecycle(t *testing.T) {
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
					testlog.StepAction(resourceType, "Creating disabled Microsoft-managed TLS certificate authority")
				},
				Config: loadAcceptanceTestTerraform("resource_01_disabled.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").ExistsInGraph(testResource),
					check.That(resourceType+".test").Key("id").Exists(),
					check.That(resourceType+".test").Key("name").MatchesRegex(regexp.MustCompile(`^M-TLSi-[0-9a-z]{5}$`)),
					check.That(resourceType+".test").Key("enabled").HasValue("false"),
					// Graph currently reports unknownFutureValue immediately after a
					// disabled create, while some tenants return disabled directly.
					check.That(resourceType+".test").Key("status").MatchesRegex(regexp.MustCompile(`^(disabled|unknownFutureValue)$`)),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing")
				},
				ResourceName:      resourceType + ".test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Enabling and waiting for activation")
				},
				Config: loadAcceptanceTestTerraform("resource_02_enabled.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").ExistsInGraph(testResource),
					check.That(resourceType+".test").Key("enabled").HasValue("true"),
					check.That(resourceType+".test").Key("status").HasValue("active"),
					check.That(resourceType+".test").Key("certificate").IsNotEmpty(),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Reading the root certificate with the data source")
				},
				Config: loadAcceptanceTestTerraform("resource_03_enabled_with_data_source.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That("data."+resourceType+".test").Key("id").Exists(),
					check.That("data."+resourceType+".test").Key("status").HasValue("active"),
					check.That("data."+resourceType+".test").Key("certificate").IsNotEmpty(),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Replacing after certificate subject changes")
				},
				Config: loadAcceptanceTestTerraform("resource_04_replaced.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").ExistsInGraph(testResource),
					check.That(resourceType+".test").Key("name").HasValue("M-TLSi-acc04"),
					check.That(resourceType+".test").Key("common_name").HasValue("Terraform Acceptance Updated TLS Inspection Root CA"),
					check.That(resourceType+".test").Key("organization_name").HasValue("Deployment Theory Acceptance"),
					check.That(resourceType+".test").Key("enabled").HasValue("true"),
					check.That(resourceType+".test").Key("status").HasValue("active"),
					check.That("data."+resourceType+".test").Key("certificate").IsNotEmpty(),
				),
			},
		},
	})
}
