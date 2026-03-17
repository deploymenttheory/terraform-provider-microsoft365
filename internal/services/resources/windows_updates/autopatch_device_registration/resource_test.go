package graphBetaWindowsUpdatesAutopatchDeviceRegistration_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaWindowsAutopatchDeviceRegistration "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/windows_updates/autopatch_device_registration"
	registrationMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/windows_updates/autopatch_device_registration/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

func setupMockEnvironment() (*mocks.Mocks, *registrationMocks.WindowsAutopatchDeviceRegistrationMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()

	registrationMock := &registrationMocks.WindowsAutopatchDeviceRegistrationMock{}
	registrationMock.RegisterMocks()
	return mockClient, registrationMock
}

func setupErrorMockEnvironment() (*mocks.Mocks, *registrationMocks.WindowsAutopatchDeviceRegistrationMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()

	registrationMock := &registrationMocks.WindowsAutopatchDeviceRegistrationMock{}
	registrationMock.RegisterErrorMocks()
	return mockClient, registrationMock
}

func loadUnitTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/unit/" + filename)
	if err != nil {
		panic("failed to load unit test config " + filename + ": " + err.Error())
	}
	return config
}

func TestUnitResourceWindowsAutopatchDeviceRegistration_01_Scenario_Minimal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, registrationMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer registrationMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("001_scenario_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaWindowsAutopatchDeviceRegistration.ResourceName+".test_001").Key("id").HasValue("feature"),
					check.That(graphBetaWindowsAutopatchDeviceRegistration.ResourceName+".test_001").Key("update_category").HasValue("feature"),
					check.That(graphBetaWindowsAutopatchDeviceRegistration.ResourceName+".test_001").Key("device_ids.#").HasValue("1"),
				),
			},
			{
				ResourceName:      graphBetaWindowsAutopatchDeviceRegistration.ResourceName + ".test_001",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestUnitResourceWindowsAutopatchDeviceRegistration_02_Scenario_Maximal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, registrationMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer registrationMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("002_scenario_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaWindowsAutopatchDeviceRegistration.ResourceName+".test_002").Key("id").HasValue("quality"),
					check.That(graphBetaWindowsAutopatchDeviceRegistration.ResourceName+".test_002").Key("update_category").HasValue("quality"),
					check.That(graphBetaWindowsAutopatchDeviceRegistration.ResourceName+".test_002").Key("device_ids.#").HasValue("3"),
				),
			},
			{
				ResourceName:      graphBetaWindowsAutopatchDeviceRegistration.ResourceName + ".test_002",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestUnitResourceWindowsAutopatchDeviceRegistration_03_Lifecycle_AddDevices(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, registrationMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer registrationMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("003_lifecycle_add_devices_step_1.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaWindowsAutopatchDeviceRegistration.ResourceName+".test_003").Key("id").HasValue("feature"),
					check.That(graphBetaWindowsAutopatchDeviceRegistration.ResourceName+".test_003").Key("device_ids.#").HasValue("1"),
				),
			},
			{
				Config: loadUnitTestTerraform("003_lifecycle_add_devices_step_2.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaWindowsAutopatchDeviceRegistration.ResourceName+".test_003").Key("id").HasValue("feature"),
					check.That(graphBetaWindowsAutopatchDeviceRegistration.ResourceName+".test_003").Key("device_ids.#").HasValue("3"),
				),
			},
		},
	})
}

func TestUnitResourceWindowsAutopatchDeviceRegistration_04_Lifecycle_RemoveDevices(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, registrationMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer registrationMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("004_lifecycle_remove_devices_step_1.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaWindowsAutopatchDeviceRegistration.ResourceName+".test_004").Key("id").HasValue("feature"),
					check.That(graphBetaWindowsAutopatchDeviceRegistration.ResourceName+".test_004").Key("device_ids.#").HasValue("3"),
				),
			},
			{
				Config: loadUnitTestTerraform("004_lifecycle_remove_devices_step_2.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaWindowsAutopatchDeviceRegistration.ResourceName+".test_004").Key("id").HasValue("feature"),
					check.That(graphBetaWindowsAutopatchDeviceRegistration.ResourceName+".test_004").Key("device_ids.#").HasValue("1"),
				),
			},
		},
	})
}

func TestUnitResourceWindowsAutopatchDeviceRegistration_05_ErrorHandling(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, registrationMock := setupErrorMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer registrationMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadUnitTestTerraform("005_error_scenario.tf"),
				ExpectError: regexp.MustCompile("Invalid request body"),
			},
		},
	})
}
