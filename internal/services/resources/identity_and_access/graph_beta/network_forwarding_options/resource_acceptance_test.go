package graphBetaNetworkForwardingOptions_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/require"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	target "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_beta/network_forwarding_options"
)

// TestAccResourceNetworkForwardingOptions_01_AdoptExisting preserves the existing DNS option.
// Run only against an approved tenant: even adoption sends PATCH to the shared singleton.
func TestAccResourceNetworkForwardingOptions_01_AdoptExisting(t *testing.T) {
	if os.Getenv("TF_ACC") != "1" {
		t.Skip("TF_ACC=1 and explicit tenant approval are required")
	}
	mocks.TestAccPreCheck(t)
	client, err := acceptance.TestGraphClient()
	require.NoError(t, err)
	original, err := client.NetworkAccess().
		Settings().
		ForwardingOptions().
		Get(context.Background(), nil)
	require.NoError(t, err)
	require.NotNil(t, original)
	require.NotNil(t, original.GetSkipDnsLookupState())
	value := original.GetSkipDnsLookupState().String()
	require.Contains(t, []string{"enabled", "disabled"}, value)
	source, err := helpers.ParseHCLFile("tests/terraform/acceptance/resource.tf")
	require.NoError(t, err)
	// Supply the baseline as the variable default, keeping the fixture reusable without fixed tenant values.
	source = fmt.Sprintf(
		"variable \"existing_dns_state\" {\n type = string\n default = %q\n }\n",
		value,
	) + source
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: func(_ *terraform.State) error {
			remote, err := client.NetworkAccess().
				Settings().
				ForwardingOptions().
				Get(context.Background(), nil)
			if err != nil {
				return err
			}
			if remote == nil || remote.GetSkipDnsLookupState() == nil ||
				remote.GetSkipDnsLookupState().String() != value {
				return fmt.Errorf("destroy changed or lost existing forwarding settings")
			}
			return nil
		},
		Steps: []resource.TestStep{
			{
				Config: acceptance.ConfiguredM365ProviderBlock(source),
				Check: resource.ComposeAggregateTestCheckFunc(
					check.That(resourceType+".test").
						ExistsInGraph(target.NetworkForwardingOptionsTestResource{}),
					check.That(resourceType+".test").Key("skip_dns_lookup_state").HasValue(value),
				),
			},
			{
				ResourceName:            resourceType + ".test",
				ImportState:             true,
				ImportStateId:           "forwardingOptions",
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
			{Config: acceptance.ConfiguredM365ProviderBlock(source), PlanOnly: true},
		},
	})
}
