package graphBetaWindowsUpdatesAutopatchOperationalInsightsConnection_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	connectionMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/windows_updates/graph_beta/autopatch_operational_insights_connection/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

const resourceType = "microsoft365_graph_beta_windows_updates_autopatch_operational_insights_connection"

func loadUnitTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/unit/" + filename)
	if err != nil {
		panic("failed to load unit test config " + filename + ": " + err.Error())
	}
	return config
}

func setupMockEnvironment() (*mocks.Mocks, *connectionMocks.WindowsUpdateOperationalInsightsConnectionMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	connectionMock := &connectionMocks.WindowsUpdateOperationalInsightsConnectionMock{}
	connectionMock.RegisterMocks()
	return mockClient, connectionMock
}

func setupErrorMockEnvironment() (*mocks.Mocks, *connectionMocks.WindowsUpdateOperationalInsightsConnectionMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	connectionMock := &connectionMocks.WindowsUpdateOperationalInsightsConnectionMock{}
	connectionMock.RegisterErrorMocks()
	return mockClient, connectionMock
}

func TestUnitResourceWindowsUpdatesAutopatchOperationalInsightsConnection_01_Create(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, connectionMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer connectionMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("01_create.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("id").HasValue("c3d4e5f6-3456-7890-cdef-c3d4e5f6a7b8"),
					check.That(resourceType+".test").Key("azure_resource_group_name").HasValue("my-resource-group"),
					check.That(resourceType+".test").Key("azure_subscription_id").HasValue("12345678-1234-1234-1234-123456789012"),
					check.That(resourceType+".test").Key("workspace_name").HasValue("my-log-analytics-workspace"),
					check.That(resourceType+".test").Key("state").HasValue("connected"),
				),
			},
		},
	})
}

func TestUnitResourceWindowsUpdatesAutopatchOperationalInsightsConnection_02_Import(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, connectionMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer connectionMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("01_create.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("id").Exists(),
				),
			},
			{
				ResourceName:            resourceType + ".test",
				ImportState:             true,
				ImportStateId:           "c3d4e5f6-3456-7890-cdef-c3d4e5f6a7b8",
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

func TestUnitResourceWindowsUpdatesAutopatchOperationalInsightsConnection_03_Error(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, connectionMock := setupErrorMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer connectionMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadUnitTestTerraform("01_create.tf"),
				ExpectError: regexp.MustCompile("Forbidden|403|Insufficient privileges"),
			},
		},
	})
}
