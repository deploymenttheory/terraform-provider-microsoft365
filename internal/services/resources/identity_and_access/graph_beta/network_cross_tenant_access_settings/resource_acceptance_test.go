package graphBetaNetworkCrossTenantAccessSettings_test

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	settings "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_beta/network_cross_tenant_access_settings"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/require"
)

// This test adopts the current value without toggling a tenant-wide access control.
// Enable/disable transitions require a separate, explicitly approved live validation.
func TestAccResourceNetworkCrossTenantAccessSettings_01_AdoptCurrentValue(t *testing.T) {
	if os.Getenv("TF_ACC") != "1" {
		t.Skip("acceptance test requires TF_ACC=1")
	}
	mocks.TestAccPreCheck(t)
	client, err := acceptance.TestGraphClient()
	require.NoError(t, err)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()
	original, err := client.NetworkAccess().Settings().CrossTenantAccess().Get(ctx, nil)
	require.NoError(t, err)
	require.NotNil(t, original)
	require.NotNil(t, original.GetNetworkPacketTaggingStatus())
	status := original.GetNetworkPacketTaggingStatus().String()
	require.Contains(t, []string{"enabled", "disabled"}, status)
	hcl, err := os.ReadFile("tests/terraform/acceptance/resource.tf")
	require.NoError(t, err)
	config := acceptance.ConfiguredM365ProviderBlock(strings.ReplaceAll(string(hcl), "{{STATUS}}", status))
	address := resourceType + ".test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: config, Check: resource.ComposeTestCheckFunc(check.That(address).Key("network_packet_tagging_status").HasValue(status), check.That(address).ExistsInGraph(settings.NetworkCrossTenantAccessSettingsTestResource{}))},
			{ResourceName: address, ImportState: true, ImportStateId: "crossTenantAccess", ImportStateVerify: true, ImportStateVerifyIgnore: []string{"timeouts"}},
			{Config: config, PlanOnly: true},
			{Config: config, Destroy: true},
		},
	})
	// Destroy must leave the backed-up value untouched; the settings must still exist.
	finalCtx, finalCancel := context.WithTimeout(context.Background(), time.Minute)
	defer finalCancel()
	remote, err := client.NetworkAccess().Settings().CrossTenantAccess().Get(finalCtx, nil)
	require.NoError(t, err)
	require.NotNil(t, remote)
	require.NotNil(t, remote.GetNetworkPacketTaggingStatus())
	require.Equal(t, status, remote.GetNetworkPacketTaggingStatus().String())
}
