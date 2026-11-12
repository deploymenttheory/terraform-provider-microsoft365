package provider_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/stretchr/testify/require"

	providercore "github.com/deploymenttheory/terraform-provider-microsoft365/internal/provider/core"
	syncdevice "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/actions/device_management/graph_beta/managed_device/sync_device"
	group "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/datasources/groups/graph_beta/group"
	auditevents "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/ephemerals/multitenant_management/graph_beta/audit_events"
	windowscripts "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/list-resources/device_management/graph_beta/windows_platform_script"
	windowscriptresource "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/windows_platform_script"
	guidsharder "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/utility/guid_list_sharder"
)

func TestUnit_ProviderRegistryObjectKinds(t *testing.T) {
	testCases := []struct {
		name     string
		registry providercore.Registry
	}{
		{
			"resource",
			providercore.Registry{
				Resources: []func() resource.Resource{guidsharder.NewGuidListSharderResource},
			},
		},
		{
			"data_source",
			providercore.Registry{
				DataSources: []func() datasource.DataSource{group.NewGroupDataSource},
			},
		},
		{
			"action",
			providercore.Registry{
				Actions: []func() action.Action{syncdevice.NewSyncManagedDeviceAction},
			},
		},
		{
			"list_resource",
			providercore.Registry{
				Resources: []func() resource.Resource{
					windowscriptresource.NewWindowsPlatformScriptResource,
				},
				ListResources: []func() list.ListResource{
					windowscripts.NewWindowsPlatformScriptListResource,
				},
			},
		},
		{
			"ephemeral_resource",
			providercore.Registry{
				EphemeralResources: []func() ephemeral.EphemeralResource{
					auditevents.NewAuditEventsEphemeralResource,
				},
			},
		},
	}
	for index, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			server, err := providerserver.NewProtocol6WithError(
				providercore.NewWithRegistry("test", true, testCase.registry),
			)()
			require.NoError(t, err)
			response, err := server.GetProviderSchema(
				context.Background(),
				&tfprotov6.GetProviderSchemaRequest{},
			)
			require.NoError(t, err)
			for _, diagnostic := range response.Diagnostics {
				require.NotEqual(
					t,
					tfprotov6.DiagnosticSeverityError,
					diagnostic.Severity,
					diagnostic.Detail,
				)
			}
			var expected [5]int
			expected[index] = 1
			if testCase.name == "list_resource" {
				expected[0] = 1
			}
			actual := [5]int{
				len(
					response.ResourceSchemas,
				),
				len(response.DataSourceSchemas),
				len(response.ActionSchemas),
				len(response.ListResourceSchemas),
				len(response.EphemeralResourceSchemas),
			}
			require.Equal(t, expected, actual)
		})
	}
}
