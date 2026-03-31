package graphBetaWindowsUpdatesAutopatchDeploymentAudience_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaWindowsAutopatchDeploymentAudience "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/windows_updates/graph_beta/autopatch_deployment_audience"
	audienceMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/windows_updates/graph_beta/autopatch_deployment_audience/mocks"
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

// loadUnitTestTerraform loads a unit test HCL config from the unit test directory
func loadUnitTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/unit/" + filename)
	if err != nil {
		panic("failed to load unit test config " + filename + ": " + err.Error())
	}
	return config
}

// Test 001: Basic audience creation (no members or exclusions)
//
// API calls exercised:
//   - POST /admin/windows/updates/deploymentAudiences
//   - GET  /admin/windows/updates/deploymentAudiences/{id}
//   - GET  /admin/windows/updates/deploymentAudiences/{id}/members
//   - GET  /admin/windows/updates/deploymentAudiences/{id}/exclusions
//   - DELETE /admin/windows/updates/deploymentAudiences/{id}
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
					check.That(graphBetaWindowsAutopatchDeploymentAudience.ResourceName+".test").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaWindowsAutopatchDeploymentAudience.ResourceName+".test").Key("member_type").HasValue("azureADDevice"),
				),
			},
			{
				ResourceName:      graphBetaWindowsAutopatchDeploymentAudience.ResourceName + ".test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// Test 002: Audience with azureADDevice members
//
// API calls exercised:
//   - POST /admin/windows/updates/deploymentAudiences
//   - POST /admin/windows/updates/deploymentAudiences/{id}/microsoft.graph.windowsUpdates.updateAudience (addMembers)
//   - GET  /admin/windows/updates/deploymentAudiences/{id}
//   - GET  /admin/windows/updates/deploymentAudiences/{id}/members
//   - GET  /admin/windows/updates/deploymentAudiences/{id}/exclusions
//   - DELETE /admin/windows/updates/deploymentAudiences/{id}
func TestUnitResourceWindowsUpdateDeploymentAudience_02_WithMembers(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, audienceMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer audienceMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("02_with_members.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaWindowsAutopatchDeploymentAudience.ResourceName+".test").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaWindowsAutopatchDeploymentAudience.ResourceName+".test").Key("member_type").HasValue("azureADDevice"),
					check.That(graphBetaWindowsAutopatchDeploymentAudience.ResourceName+".test").Key("members.#").HasValue("2"),
				),
			},
			{
				ResourceName:      graphBetaWindowsAutopatchDeploymentAudience.ResourceName + ".test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// Test 003: Audience with azureADDevice members and exclusions
//
// API calls exercised:
//   - POST /admin/windows/updates/deploymentAudiences
//   - POST /admin/windows/updates/deploymentAudiences/{id}/microsoft.graph.windowsUpdates.updateAudience (addMembers + addExclusions)
//   - GET  /admin/windows/updates/deploymentAudiences/{id}
//   - GET  /admin/windows/updates/deploymentAudiences/{id}/members
//   - GET  /admin/windows/updates/deploymentAudiences/{id}/exclusions
//   - DELETE /admin/windows/updates/deploymentAudiences/{id}
func TestUnitResourceWindowsUpdateDeploymentAudience_03_WithExclusions(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, audienceMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer audienceMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("03_with_exclusions.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaWindowsAutopatchDeploymentAudience.ResourceName+".test").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaWindowsAutopatchDeploymentAudience.ResourceName+".test").Key("member_type").HasValue("azureADDevice"),
					check.That(graphBetaWindowsAutopatchDeploymentAudience.ResourceName+".test").Key("members.#").HasValue("2"),
					check.That(graphBetaWindowsAutopatchDeploymentAudience.ResourceName+".test").Key("exclusions.#").HasValue("1"),
				),
			},
			{
				ResourceName:      graphBetaWindowsAutopatchDeploymentAudience.ResourceName + ".test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// Test 004: Lifecycle - add a member via updateAudience (diff-based update)
//
// Step 1: Create audience with 2 azureADDevice members.
// Step 2: Add a third member - exercises the diff calculation in constructUpdateMembersRequest
//         which sends only the delta (addMembers=[id3]) to updateAudience.
//
// API calls exercised (step 1):
//   - POST /admin/windows/updates/deploymentAudiences
//   - POST /admin/windows/updates/deploymentAudiences/{id}/microsoft.graph.windowsUpdates.updateAudience (addMembers x2)
//   - GET  /admin/windows/updates/deploymentAudiences/{id}
//   - GET  /admin/windows/updates/deploymentAudiences/{id}/members
//   - GET  /admin/windows/updates/deploymentAudiences/{id}/exclusions
//
// API calls exercised (step 2):
//   - POST /admin/windows/updates/deploymentAudiences/{id}/microsoft.graph.windowsUpdates.updateAudience (addMembers x1)
//   - GET  /admin/windows/updates/deploymentAudiences/{id}
//   - GET  /admin/windows/updates/deploymentAudiences/{id}/members
//   - GET  /admin/windows/updates/deploymentAudiences/{id}/exclusions
//
// API calls exercised (destroy):
//   - POST /admin/windows/updates/deploymentAudiences/{id}/microsoft.graph.windowsUpdates.updateAudience (removeMembers x3)
//   - DELETE /admin/windows/updates/deploymentAudiences/{id}
func TestUnitResourceWindowsUpdateDeploymentAudience_04_Lifecycle(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, audienceMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer audienceMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("04_lifecycle_step1.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaWindowsAutopatchDeploymentAudience.ResourceName+".test").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaWindowsAutopatchDeploymentAudience.ResourceName+".test").Key("member_type").HasValue("azureADDevice"),
					check.That(graphBetaWindowsAutopatchDeploymentAudience.ResourceName+".test").Key("members.#").HasValue("2"),
				),
			},
			{
				Config: loadUnitTestTerraform("04_lifecycle_step2.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaWindowsAutopatchDeploymentAudience.ResourceName+".test").Key("member_type").HasValue("azureADDevice"),
					check.That(graphBetaWindowsAutopatchDeploymentAudience.ResourceName+".test").Key("members.#").HasValue("3"),
				),
			},
			{
				ResourceName:      graphBetaWindowsAutopatchDeploymentAudience.ResourceName + ".test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// Test 005: Error handling - verifies that a server-side error during audience creation
// surfaces as a Terraform error with the expected HTTP status code.
//
// API calls exercised:
//   - POST /admin/windows/updates/deploymentAudiences → 400 BadRequest
func TestUnitResourceWindowsUpdateDeploymentAudience_05_Error(t *testing.T) {
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
