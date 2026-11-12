package graphBetaNetworkForwardingProfile_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/require"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	target "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_beta/network_forwarding_profile"
)

// TestAccResourceNetworkForwardingProfile_01_Lifecycle needs a tenant on traffic
// forwarding profiles v2 with a free custom private priority (default 190).
func TestAccResourceNetworkForwardingProfile_01_Lifecycle(t *testing.T) {
	if os.Getenv("TF_ACC") != "1" {
		t.Skip("TF_ACC=1 and approved tenant required")
	}
	mocks.TestAccPreCheck(t)
	priority := os.Getenv("M365_ACC_FORWARDING_PROFILE_PRIORITY")
	if priority == "" {
		priority = "190"
	}
	source, err := helpers.ParseHCLFile("tests/terraform/acceptance/resource.tf")
	require.NoError(t, err)
	config := acceptance.ConfiguredM365ProviderBlock(
		fmt.Sprintf("variable \"priority\" { default = %s }\n", priority) + source,
	)
	addr := resourceType + ".test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					check.That(addr).ExistsInGraph(target.NetworkForwardingProfileTestResource{}),
					check.That(addr).
						Key("local_network_settings.allow_local_subnet_bypass").
						HasValue("true"),
					check.That(addr).Key("assignment_rules.membership").HasValue("all"),
				),
			},
			{
				ResourceName:            addr,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
			{Config: config, PlanOnly: true},
		},
	})
}
