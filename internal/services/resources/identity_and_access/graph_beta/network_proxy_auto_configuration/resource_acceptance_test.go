package graphBetaNetworkProxyAutoConfiguration_test

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
	target "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_beta/network_proxy_auto_configuration"
)

// TestAccResourceNetworkProxyAutoConfiguration_01_Lifecycle requires explicit approval to update tenant-wide proxy settings.
func TestAccResourceNetworkProxyAutoConfiguration_01_Lifecycle(t *testing.T) {
	if os.Getenv("TF_ACC") != "1" {
		t.Skip("TF_ACC=1 and approved tenant required")
	}
	mocks.TestAccPreCheck(t)
	client, err := acceptance.TestGraphClient()
	require.NoError(t, err)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()
	var original struct {
		InternetAccess struct {
			IsEnabled bool   `json:"isEnabled"`
			Affinity  bool   `json:"isSourceIpSessionAffinityEnabled"`
			Options   string `json:"sourceIpSessionAffinityOptions"`
		} `json:"internetAccess"`
	}
	require.NoError(
		t,
		customrequests.JSONRequest(
			ctx,
			client.GetAdapter(),
			abstractions.GET,
			"/networkAccess/explicitForwardProxyConfig",
			nil,
			nil,
			&original,
		),
	)
	access := original.InternetAccess
	require.Contains(
		t,
		[]string{
			"none",
			"useSessionId",
			"useHttpHeader",
			"useSessionId,useHttpHeader",
			"useHttpHeader,useSessionId",
		},
		access.Options,
	)
	t.Cleanup(func() {
		restoreCtx, restoreCancel := context.WithTimeout(context.Background(), 3*time.Minute)
		defer restoreCancel()
		restore := map[string]any{
			"internetAccess": map[string]any{
				"isEnabled":                        access.IsEnabled,
				"isSourceIpSessionAffinityEnabled": access.Affinity,
				"sourceIpSessionAffinityOptions":   access.Options,
			},
		}
		require.NoError(
			t,
			customrequests.JSONRequest(
				restoreCtx,
				client.GetAdapter(),
				abstractions.PATCH,
				"/networkAccess/explicitForwardProxyConfig",
				nil,
				restore,
				nil,
			),
		)
	})
	source, err := helpers.ParseHCLFile("tests/terraform/acceptance/resource.tf")
	require.NoError(t, err)
	source = fmt.Sprintf(
		"variable \"proxy_enabled\" { default = %t }\nvariable \"affinity_enabled\" { default = %t }\nvariable \"affinity_options\" { default = %q }\n",
		true,
		access.Affinity,
		canonicalAffinityOptions(access.Options),
	) + source

	config := acceptance.ConfiguredM365ProviderBlock(source)
	resource.Test(
		t,
		resource.TestCase{
			ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: config,
					Check: check.That(resourceType + ".test").
						ExistsInGraph(target.NetworkProxyAutoConfigurationTestResource{}),
				},
				{
					ResourceName:            resourceType + ".test",
					ImportState:             true,
					ImportStateVerify:       true,
					ImportStateVerifyIgnore: []string{"timeouts"},
				},
				{Config: config, PlanOnly: true},
			},
		},
	)
}

func canonicalAffinityOptions(options string) string {
	if options == "useHttpHeader,useSessionId" {
		return "useSessionId,useHttpHeader"
	}
	return options
}
