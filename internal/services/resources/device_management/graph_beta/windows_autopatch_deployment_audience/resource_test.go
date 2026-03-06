package graphBetaWindowsAutopatchDeploymentAudience_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaWindowsUpdateDeploymentAudience "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/windows_update_deployment_audience"
	audienceMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/windows_update_deployment_audience/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

// setupMockEnvironment sets up the mock environment using centralized mocks
func setupMockEnvironment() (*mocks.Mocks, *audienceMocks.WindowsUpdateDeploymentAudienceMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()

	audienceMock := &audienceMocks.WindowsUpdateDeploymentAudienceMock{}
	audienceMock.RegisterMocks()
	return mockClient, audienceMock
}

// setupErrorMockEnvironment sets up the mock environment for error testing
func setupErrorMockEnvironment() (*mocks.Mocks, *audienceMocks.WindowsUpdateDeploymentAudienceMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()

	audienceMock := &audienceMocks.WindowsUpdateDeploymentAudienceMock{}
	audienceMock.RegisterErrorMocks()
	return mockClient, audienceMock
}

// Helper function to load test configs from unit directory
func loadUnitTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/unit/" + filename)
	if err != nil {
		panic("failed to load unit test config " + filename + ": " + err.Error())
	}
	return config
}

// Test 001: Basic audience creation
func TestUnitResourceWindowsUpdateDeploymentAudience_01_BasicAudience(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, audienceMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer audienceMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("01_basic_audience.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaWindowsUpdateDeploymentAudience.ResourceName + ".test").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
				),
			},
			{
				ResourceName:      graphBetaWindowsUpdateDeploymentAudience.ResourceName + ".test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// Test 002: Error handling
func TestUnitResourceWindowsUpdateDeploymentAudience_02_Error(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, audienceMock := setupErrorMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer audienceMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadUnitTestTerraform("01_basic_audience.tf"),
				ExpectError: regexp.MustCompile("BadRequest|400|Invalid"),
			},
		},
	})
}
