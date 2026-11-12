package mocks_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/stretchr/testify/require"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	providercore "github.com/deploymenttheory/terraform-provider-microsoft365/internal/provider/core"
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
		response, err := server.GetProviderSchema(
			context.Background(),
			&tfprotov6.GetProviderSchemaRequest{},
		)
		require.NoError(t, err)
		require.Len(t, response.ResourceSchemas, testCase.count)
	}
}
