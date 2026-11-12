package graphBetaAzureNetworkConnection_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/stretchr/testify/require"

	azureConnection "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/windows_365/graph_beta/azure_network_connection"
)

func TestUnitResourceAzureNetworkConnection_IsolatedProviderSchema(t *testing.T) {
	server, err := testUnitTestProtoV6ProviderFactories["microsoft365"]()
	require.NoError(t, err)
	response, err := server.GetProviderSchema(context.Background(), &tfprotov6.GetProviderSchemaRequest{})
	require.NoError(t, err)
	for _, diagnostic := range response.Diagnostics {
		require.NotEqual(t, tfprotov6.DiagnosticSeverityError, diagnostic.Severity, diagnostic.Detail)
	}
	require.Len(t, response.ResourceSchemas, 1)
	require.Contains(t, response.ResourceSchemas, azureConnection.ResourceName)
	require.Empty(t, response.DataSourceSchemas)
	require.Empty(t, response.ActionSchemas)
	require.Empty(t, response.ListResourceSchemas)
	require.Empty(t, response.EphemeralResourceSchemas)
}
