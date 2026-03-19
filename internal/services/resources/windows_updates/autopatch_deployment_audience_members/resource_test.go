package graphBetaWindowsUpdatesAutopatchDeploymentAudienceMembers_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	WindowsUpdatesAutopatchDeploymentResourceAudienceMembers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/windows_updates/autopatch_deployment_audience_members"
	membersMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/windows_updates/autopatch_deployment_audience_members/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

// setupMockEnvironment sets up the mock environment using centralized mocks
func setupMockEnvironment() (*mocks.Mocks, *membersMocks.WindowsUpdateDeploymentAudienceMembersMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()

	membersMock := &membersMocks.WindowsUpdateDeploymentAudienceMembersMock{}
	membersMock.RegisterMocks()
	return mockClient, membersMock
}

// setupErrorMockEnvironment sets up the mock environment for error testing
func setupErrorMockEnvironment() (*mocks.Mocks, *membersMocks.WindowsUpdateDeploymentAudienceMembersMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()

	membersMock := &membersMocks.WindowsUpdateDeploymentAudienceMembersMock{}
	membersMock.RegisterErrorMocks()
	return mockClient, membersMock
}

// Helper function to load test configs from unit directory
func loadUnitTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/unit/" + filename)
	if err != nil {
		panic("failed to load unit test config " + filename + ": " + err.Error())
	}
	return config
}

// Test 001: Basic members with devices
func TestUnitResourceWindowsUpdateDeploymentAudienceMembers_01_BasicMembers(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, membersMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer membersMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("01_basic_members.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(WindowsUpdatesAutopatchDeploymentResourceAudienceMembers.ResourceName+".test").Key("id").MatchesRegex(regexp.MustCompile(`^.+_azureADDevice$`)),
					check.That(WindowsUpdatesAutopatchDeploymentResourceAudienceMembers.ResourceName+".test").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+_azureADDevice$`)),
					check.That(WindowsUpdatesAutopatchDeploymentResourceAudienceMembers.ResourceName+".test").Key("audience_id").HasValue("test-audience-id-001"),
					check.That(WindowsUpdatesAutopatchDeploymentResourceAudienceMembers.ResourceName+".test").Key("member_type").HasValue("azureADDevice"),
					check.That(WindowsUpdatesAutopatchDeploymentResourceAudienceMembers.ResourceName+".test").Key("members.#").HasValue("2"),
				),
			},
			{
				ResourceName:      WindowsUpdatesAutopatchDeploymentResourceAudienceMembers.ResourceName + ".test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// Test 002: Members with exclusions
func TestUnitResourceWindowsUpdateDeploymentAudienceMembers_02_MembersWithExclusions(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, membersMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer membersMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("02_members_with_exclusions.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(WindowsUpdatesAutopatchDeploymentResourceAudienceMembers.ResourceName+".test").Key("id").MatchesRegex(regexp.MustCompile(`^.+_azureADDevice$`)),
					check.That(WindowsUpdatesAutopatchDeploymentResourceAudienceMembers.ResourceName+".test").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+_azureADDevice$`)),
					check.That(WindowsUpdatesAutopatchDeploymentResourceAudienceMembers.ResourceName+".test").Key("member_type").HasValue("azureADDevice"),
					check.That(WindowsUpdatesAutopatchDeploymentResourceAudienceMembers.ResourceName+".test").Key("members.#").HasValue("2"),
					check.That(WindowsUpdatesAutopatchDeploymentResourceAudienceMembers.ResourceName+".test").Key("exclusions.#").HasValue("1"),
				),
			},
			{
				ResourceName:      WindowsUpdatesAutopatchDeploymentResourceAudienceMembers.ResourceName + ".test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// Test 003: Error handling
func TestUnitResourceWindowsUpdateDeploymentAudienceMembers_03_Error(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, membersMock := setupErrorMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer membersMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadUnitTestTerraform("01_basic_members.tf"),
				ExpectError: regexp.MustCompile("BadRequest|400|Invalid"),
			},
		},
	})
}
