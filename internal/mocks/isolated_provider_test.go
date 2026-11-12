package mocks_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/stretchr/testify/require"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	providercore "github.com/deploymenttheory/terraform-provider-microsoft365/internal/provider/core"
	syncdevice "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/actions/device_management/graph_beta/managed_device/sync_device"
	group "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/datasources/groups/graph_beta/group"
	auditevents "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/ephemerals/multitenant_management/graph_beta/audit_events"
	windowscripts "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/list-resources/device_management/graph_beta/windows_platform_script"
	windowscriptresource "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/windows_platform_script"
	guidsharder "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/utility/guid_list_sharder"
)

func TestUnit_IsolatedProviderFactorySchema(t *testing.T) {
	factories := mocks.NewProviderFactories(true, providercore.Registry{
		Resources: []func() resource.Resource{guidsharder.NewGuidListSharderResource},
	})
	server, err := factories["microsoft365"]()
	require.NoError(t, err)
	other, err := factories["microsoft365"]()
	require.NoError(t, err)
	require.NotSame(t, server, other)
	response, err := server.GetProviderSchema(context.Background(), &tfprotov6.GetProviderSchemaRequest{})
	require.NoError(t, err)
	for _, diagnostic := range response.Diagnostics {
		require.NotEqual(t, tfprotov6.DiagnosticSeverityError, diagnostic.Severity, diagnostic.Detail)
	}
	require.Contains(t, response.ResourceSchemas, guidsharder.ResourceName)
	require.Len(t, response.ResourceSchemas, 1)
	require.Empty(t, response.DataSourceSchemas)
	require.NotNil(t, response.Provider)
}

func TestUnit_ProviderRegistriesDoNotLeakBetweenPartitions(t *testing.T) {
	withResource := mocks.NewProviderFactories(true, providercore.Registry{
		Resources: []func() resource.Resource{guidsharder.NewGuidListSharderResource},
	})
	empty := mocks.NewProviderFactories(true, providercore.Registry{})
	for _, testCase := range []struct {
		factories map[string]func() (tfprotov6.ProviderServer, error)
		count     int
	}{{withResource, 1}, {empty, 0}, {withResource, 1}} {
		server, err := testCase.factories["microsoft365"]()
		require.NoError(t, err)
		response, err := server.GetProviderSchema(context.Background(), &tfprotov6.GetProviderSchemaRequest{})
		require.NoError(t, err)
		require.Len(t, response.ResourceSchemas, testCase.count)
	}
}

func TestUnit_IsolatedProviderObjectKinds(t *testing.T) {
	testCases := []struct {
		name     string
		registry providercore.Registry
	}{
		{"resource", providercore.Registry{Resources: []func() resource.Resource{guidsharder.NewGuidListSharderResource}}},
		{"data_source", providercore.Registry{DataSources: []func() datasource.DataSource{group.NewGroupDataSource}}},
		{"action", providercore.Registry{Actions: []func() action.Action{syncdevice.NewSyncManagedDeviceAction}}},
		{"list_resource", providercore.Registry{
			Resources:     []func() resource.Resource{windowscriptresource.NewWindowsPlatformScriptResource},
			ListResources: []func() list.ListResource{windowscripts.NewWindowsPlatformScriptListResource},
		}},
		{"ephemeral_resource", providercore.Registry{EphemeralResources: []func() ephemeral.EphemeralResource{auditevents.NewAuditEventsEphemeralResource}}},
	}
	for index, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			factories := mocks.NewProviderFactories(true, testCase.registry)
			server, err := factories["microsoft365"]()
			require.NoError(t, err)
			response, err := server.GetProviderSchema(context.Background(), &tfprotov6.GetProviderSchemaRequest{})
			require.NoError(t, err)
			for _, diagnostic := range response.Diagnostics {
				require.NotEqual(t, tfprotov6.DiagnosticSeverityError, diagnostic.Severity, diagnostic.Detail)
			}
			var expected [5]int
			expected[index] = 1
			if testCase.name == "list_resource" {
				expected[0] = 1
			}
			actual := [5]int{len(response.ResourceSchemas), len(response.DataSourceSchemas), len(response.ActionSchemas),
				len(response.ListResourceSchemas), len(response.EphemeralResourceSchemas)}
			require.Equal(t, expected, actual)
		})
	}
}
