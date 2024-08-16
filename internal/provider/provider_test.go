package provider

import (
	"context"
	"testing"

	graphBetaAssignmentFilter "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/deviceandappmanagement/beta/assignmentFilter"
	graphCloudPcProvisioningPolicy "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/devicemanagement/v1.0/cloudPcProvisioningPolicy"
	graphBetaConditionalAccessPolicy "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/identityandaccess/beta/conditionalaccesspolicy"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	frameworkResource "github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	testingResource "github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
)

const (
	TestsUnitProviderConfig = `
provider "microsoft365" {
  // Unit test provider config
}
`

	TestsAcceptanceProviderConfig = `
provider "microsoft365" {
  // Acceptance test provider config
}
`
)

var TestUnitTestProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"microsoft365": providerserver.NewProtocol6WithError(New("1.0.0")()),
}

var TestAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"microsoft365": providerserver.NewProtocol6WithError(New("1.0.0")()),
}

func TestUnitM365ProviderHasChildResources_Basic(t *testing.T) {
	expectedResources := []frameworkResource.Resource{
		graphCloudPcProvisioningPolicy.NewCloudPcProvisioningPolicyResource(),
		graphBetaAssignmentFilter.NewAssignmentFilterResource(),
		graphBetaConditionalAccessPolicy.NewConditionalAccessPolicyResource(),
	}

	providerFunc := New("1.0.0")
	providerInstance, ok := providerFunc().(*M365Provider)
	require.True(t, ok, "Failed to cast provider to *M365Provider")

	resources := providerInstance.Resources(context.Background())

	require.Equal(t, len(expectedResources), len(resources), "There are an unexpected number of registered resources")
	for _, r := range resources {
		require.Contains(t, expectedResources, r(), "An unexpected resource was registered")
	}
}

func TestUnitM365ProviderHasChildDataSources_Basic(t *testing.T) {
	expectedDataSources := []datasource.DataSource{
		// Add expected data sources here as they are implemented
	}

	providerFunc := New("1.0.0")
	providerInstance, ok := providerFunc().(*M365Provider)
	require.True(t, ok, "Failed to cast provider to *M365Provider")

	datasources := providerInstance.DataSources(context.Background())

	require.Equal(t, len(expectedDataSources), len(datasources), "There are an unexpected number of registered data sources")
	for _, d := range datasources {
		require.Contains(t, expectedDataSources, d(), "An unexpected data source was registered")
	}
}

func TestUnitM365Provider_Validate_Telemetry_Optout_Is_False(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Register HTTP mocks here

	testingResource.Test(t, testingResource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: TestUnitTestProtoV6ProviderFactories,
		Steps: []testingResource.TestStep{
			{
				Config: `provider "microsoft365" {
					telemetry_optout = false
				}	
				// Add your data source or resource configuration here
				`,
			},
		},
	})
}

func TestUnitM365Provider_Validate_Telemetry_Optout_Is_True(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Register HTTP mocks here

	testingResource.Test(t, testingResource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: TestUnitTestProtoV6ProviderFactories,
		Steps: []testingResource.TestStep{
			{
				Config: `provider "microsoft365" {
					telemetry_optout = true
				}	
				// Add your data source or resource configuration here
				`,
			},
		},
	})
}

func TestAccM365Provider_Basic(t *testing.T) {
	testingResource.Test(t, testingResource.TestCase{
		PreCheck: func() {
			// Check if environment variables like M365_TENANT_ID are set
		},
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []testingResource.TestStep{
			{
				Config: `provider "microsoft365" {
					// Provide the necessary configuration for an acceptance test
				}`,
				Check: testingResource.ComposeTestCheckFunc(
				// Add any checks you need to verify the provider's behavior
				),
			},
		},
	})
}
