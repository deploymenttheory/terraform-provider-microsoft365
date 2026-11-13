package graphBetaNetworkDefaultForwardingProfile_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	abstractions "github.com/microsoft/kiota-abstractions-go"
	"github.com/stretchr/testify/require"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	customrequests "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/custom_requests"
)

// TestAccResourceNetworkDefaultForwardingProfile_01_Lifecycle changes the
// tenant's default private profile and restores it afterwards; it requires
// explicit approval for the tenant.
func TestAccResourceNetworkDefaultForwardingProfile_01_Lifecycle(t *testing.T) {
	if os.Getenv("TF_ACC") != "1" {
		t.Skip("TF_ACC=1 and approved tenant required")
	}
	mocks.TestAccPreCheck(t)
	client, err := acceptance.TestGraphClient()
	require.NoError(t, err)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	var profiles struct {
		Value []map[string]any `json:"value"`
	}
	require.NoError(t, customrequests.JSONRequest(ctx, client.GetAdapter(), abstractions.GET,
		"/networkaccess/forwardingProfiles", nil, nil, &profiles))
	var original map[string]any
	for _, p := range profiles.Value {
		if p["trafficForwardingType"] == "private" && p["isCustomProfile"] == false {
			original = p
		}
	}
	require.NotNil(t, original, "default private profile not found")
	t.Cleanup(func() {
		restoreCtx, restoreCancel := context.WithTimeout(context.Background(), 3*time.Minute)
		defer restoreCancel()
		restore := map[string]any{
			"state":                    original["state"],
			"localNetworkSettings":     original["localNetworkSettings"],
			"forwardingSourceMatching": original["forwardingSourceMatching"],
		}
		require.NoError(
			t,
			customrequests.JSONRequest(
				restoreCtx,
				client.GetAdapter(),
				abstractions.PATCH,
				"/networkaccess/forwardingProfiles/{id}",
				map[string]string{"id": fmt.Sprint(original["id"])},
				restore,
				nil,
			),
		)
	})

	source, err := helpers.ParseHCLFile("tests/terraform/acceptance/resource.tf")
	require.NoError(t, err)
	config := acceptance.ConfiguredM365ProviderBlock(
		fmt.Sprintf("variable \"state\" { default = %q }\n", original["state"]) + source,
	)
	addr := resourceType + ".test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					check.That(addr).Key("id").HasValue(fmt.Sprint(original["id"])),
					check.That(addr).
						Key("local_network_settings.allow_local_subnet_bypass").
						HasValue("true"),
				),
			},
			{
				ResourceName:            addr,
				ImportState:             true,
				ImportStateId:           "private",
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
			{Config: config, PlanOnly: true},
		},
	})
}
