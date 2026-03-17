package graphBetaWindowsUpdatesAutopatchDeploymentState_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	WindowsUpdatesAutopatchDeploymentStateResource "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/windows_updates/autopatch_deployment_state"
	deploymentStateMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/windows_updates/autopatch_deployment_state/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

var (
	resourceType = WindowsUpdatesAutopatchDeploymentStateResource.ResourceName
)

func loadUnitTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/unit/" + filename)
	if err != nil {
		panic("failed to load unit test config " + filename + ": " + err.Error())
	}
	return config
}

func setupMockEnvironment() (*mocks.Mocks, *deploymentStateMocks.WindowsUpdatesAutopatchDeploymentStateMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	stateMock := &deploymentStateMocks.WindowsUpdatesAutopatchDeploymentStateMock{}
	stateMock.RegisterMocks()
	return mockClient, stateMock
}

func setupErrorMockEnvironment() (*mocks.Mocks, *deploymentStateMocks.WindowsUpdatesAutopatchDeploymentStateMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	stateMock := &deploymentStateMocks.WindowsUpdatesAutopatchDeploymentStateMock{}
	stateMock.RegisterErrorMocks()
	return mockClient, stateMock
}

func TestUnitResourceWindowsUpdateDeploymentState_01_PauseDeployment(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, stateMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer stateMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("01_deployment_state_paused.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("id").Exists(),
					check.That(resourceType+".test").Key("deployment_id").Exists(),
					check.That(resourceType+".test").Key("requested_value").HasValue("paused"),
					check.That(resourceType+".test").Key("effective_value").Exists(),
				),
			},
		},
	})
}

func TestUnitResourceWindowsUpdateDeploymentState_02_UpdateState(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, stateMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer stateMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("01_deployment_state_paused.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("requested_value").HasValue("paused"),
				),
			},
			{
				Config: loadUnitTestTerraform("02_deployment_state_none.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("requested_value").HasValue("none"),
				),
			},
		},
	})
}

func TestUnitResourceWindowsUpdateDeploymentState_03_Import(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, stateMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer stateMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("01_deployment_state_paused.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("id").Exists(),
				),
			},
			{
				ResourceName:            resourceType + ".test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

func TestUnitResourceWindowsUpdateDeploymentState_04_Error(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, stateMock := setupErrorMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer stateMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadUnitTestTerraform("01_deployment_state_paused.tf"),
				ExpectError: regexp.MustCompile("BadRequest|400|Invalid"),
			},
		},
	})
}
