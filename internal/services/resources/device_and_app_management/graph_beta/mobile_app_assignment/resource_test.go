package graphBetaDeviceAndAppManagementAppAssignment_test

import (
	"fmt"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaMobileAppAssignment "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_and_app_management/graph_beta/mobile_app_assignment"
	mobileAppAssignmentMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_and_app_management/graph_beta/mobile_app_assignment/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/jarcoal/httpmock"
)

const resourceType = graphBetaMobileAppAssignment.ResourceName

// setupMockEnvironment sets up the mock environment using centralized mocks
func setupMockEnvironment() (*mocks.Mocks, *mobileAppAssignmentMocks.MobileAppAssignmentMock) {
	// Activate httpmock
	httpmock.Activate()

	// Create a new Mocks instance and register authentication mocks
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()

	// Register local mocks directly
	assignmentMock := &mobileAppAssignmentMocks.MobileAppAssignmentMock{}
	assignmentMock.RegisterMocks()

	return mockClient, assignmentMock
}

// testConfig returns a unit test configuration read from tests/terraform/unit
func testConfig(t *testing.T, name string) string {
	t.Helper()

	content, err := helpers.ParseHCLFile(filepath.Join("tests", "terraform", "unit", name))
	if err != nil {
		t.Fatalf("failed to read test configuration %s: %v", name, err)
	}
	return content
}

// checkIsRemovableSent asserts whether isRemovable was present in the settings payload the
// provider actually sent for the given settings block, rather than only in Terraform state.
func checkIsRemovableSent(t *testing.T, assignmentMock *mobileAppAssignmentMocks.MobileAppAssignmentMock, block string, wantPresent bool, wantValue bool) resource.TestCheckFunc {
	t.Helper()

	return func(*terraform.State) error {
		settings, ok := assignmentMock.SettingsSent(block)
		if !ok {
			return fmt.Errorf("no assignment request was captured for settings block %s", block)
		}

		value, present := settings["isRemovable"]
		if present != wantPresent {
			return fmt.Errorf("isRemovable present in request = %v, want %v (payload: %v)", present, wantPresent, settings)
		}
		if wantPresent && value != wantValue {
			return fmt.Errorf("isRemovable sent as %v, want %v", value, wantValue)
		}
		return nil
	}
}

// TestUnitResourceMobileAppAssignment_01_Create_Available_IosVpp verifies that an iOS VPP
// assignment with intent "available" can be created when is_removable is omitted from the
// configuration. The Intune service rejects isRemovable for any non-required intent, so the
// provider must not send the field at all.
func TestUnitResourceMobileAppAssignment_01_Create_Available_IosVpp(t *testing.T) {
	_, assignmentMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer assignmentMock.CleanupMockState()

	mocks.SetupUnitTestEnvironment(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfig(t, "resource_available_ios_vpp.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".available_ios_vpp").Key("id").Exists(),
					check.That(resourceType+".available_ios_vpp").Key("intent").HasValue("available"),
					check.That(resourceType+".available_ios_vpp").Key("settings.ios_vpp.use_device_licensing").HasValue("false"),
					check.That(resourceType+".available_ios_vpp").Key("settings.ios_vpp.is_removable").DoesNotExist(),
					checkIsRemovableSent(t, assignmentMock, "iosVpp", false, false),
				),
			},
		},
	})
}

// TestUnitResourceMobileAppAssignment_02_Create_Available_IosStore verifies the same behaviour
// for the ios_store settings block.
func TestUnitResourceMobileAppAssignment_02_Create_Available_IosStore(t *testing.T) {
	_, assignmentMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer assignmentMock.CleanupMockState()

	mocks.SetupUnitTestEnvironment(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfig(t, "resource_available_ios_store.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".available_ios_store").Key("id").Exists(),
					check.That(resourceType+".available_ios_store").Key("intent").HasValue("available"),
					check.That(resourceType+".available_ios_store").Key("settings.ios_store.is_removable").DoesNotExist(),
					checkIsRemovableSent(t, assignmentMock, "iosStore", false, false),
				),
			},
		},
	})
}

// TestUnitResourceMobileAppAssignment_03_Create_Available_IosLob verifies the same behaviour
// for the ios_lob settings block.
func TestUnitResourceMobileAppAssignment_03_Create_Available_IosLob(t *testing.T) {
	_, assignmentMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer assignmentMock.CleanupMockState()

	mocks.SetupUnitTestEnvironment(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfig(t, "resource_available_ios_lob.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".available_ios_lob").Key("id").Exists(),
					check.That(resourceType+".available_ios_lob").Key("intent").HasValue("available"),
					check.That(resourceType+".available_ios_lob").Key("settings.ios_lob.is_removable").DoesNotExist(),
					checkIsRemovableSent(t, assignmentMock, "iosLob", false, false),
				),
			},
		},
	})
}

// TestUnitResourceMobileAppAssignment_04_Create_Required_IosVpp verifies that explicitly
// setting is_removable alongside intent "required" continues to work unchanged.
func TestUnitResourceMobileAppAssignment_04_Create_Required_IosVpp(t *testing.T) {
	_, assignmentMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer assignmentMock.CleanupMockState()

	mocks.SetupUnitTestEnvironment(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfig(t, "resource_required_ios_vpp.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".required_ios_vpp").Key("id").Exists(),
					check.That(resourceType+".required_ios_vpp").Key("intent").HasValue("required"),
					check.That(resourceType+".required_ios_vpp").Key("settings.ios_vpp.is_removable").HasValue("true"),
					checkIsRemovableSent(t, assignmentMock, "iosVpp", true, true),
				),
			},
		},
	})
}

// TestUnitResourceMobileAppAssignment_05_Create_Required_IsRemovable_Omitted verifies that the
// static schema default still applies for required intent when the attribute is omitted, so
// that required-intent behaviour is unchanged.
func TestUnitResourceMobileAppAssignment_05_Create_Required_IsRemovable_Omitted(t *testing.T) {
	_, assignmentMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer assignmentMock.CleanupMockState()

	mocks.SetupUnitTestEnvironment(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfig(t, "resource_required_ios_vpp_omitted.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".required_omitted").Key("intent").HasValue("required"),
					check.That(resourceType+".required_omitted").Key("settings.ios_vpp.is_removable").HasValue("false"),
					checkIsRemovableSent(t, assignmentMock, "iosVpp", true, false),
				),
			},
		},
	})
}

// TestUnitResourceMobileAppAssignment_06_Create_Uninstall_IosVpp verifies that intents other
// than "available" which the service also rejects are handled, not just the available case.
func TestUnitResourceMobileAppAssignment_06_Create_Uninstall_IosVpp(t *testing.T) {
	_, assignmentMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer assignmentMock.CleanupMockState()

	mocks.SetupUnitTestEnvironment(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfig(t, "resource_uninstall_ios_vpp.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".uninstall_ios_vpp").Key("intent").HasValue("uninstall"),
					check.That(resourceType+".uninstall_ios_vpp").Key("settings.ios_vpp.is_removable").DoesNotExist(),
					checkIsRemovableSent(t, assignmentMock, "iosVpp", false, false),
				),
			},
		},
	})
}

// TestUnitResourceMobileAppAssignment_07_Validate_IsRemovable_UnsupportedIntent verifies that
// explicitly setting is_removable with a non-required intent fails at plan time with a clear
// message, rather than with an opaque HTTP 400 at apply time. Every Apple settings block is
// covered, for both boolean values, across each unsupported intent.
func TestUnitResourceMobileAppAssignment_07_Validate_IsRemovable_UnsupportedIntent(t *testing.T) {
	testCases := []struct {
		name   string
		block  string
		value  string
		intent string
	}{
		{name: "ios_vpp_true_available", block: "ios_vpp", value: "true", intent: "available"},
		{name: "ios_vpp_false_available", block: "ios_vpp", value: "false", intent: "available"},
		{name: "ios_store_true_available", block: "ios_store", value: "true", intent: "available"},
		{name: "ios_lob_false_uninstall", block: "ios_lob", value: "false", intent: "uninstall"},
		{name: "ios_vpp_true_available_without_enrollment", block: "ios_vpp", value: "true", intent: "availableWithoutEnrollment"},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			_, assignmentMock := setupMockEnvironment()
			defer httpmock.DeactivateAndReset()
			defer assignmentMock.CleanupMockState()

			mocks.SetupUnitTestEnvironment(t)

			config := fmt.Sprintf(`
resource "microsoft365_graph_beta_device_and_app_management_mobile_app_assignment" "invalid" {
  mobile_app_id = "00000000-0000-0000-0000-000000000001"
  intent        = %q
  source        = "direct"

  target = {
    target_type = "groupAssignment"
    group_id    = "11111111-1111-1111-1111-111111111111"
  }

  settings = {
    %s = {
      is_removable = %s
    }
  }
}
`, testCase.intent, testCase.block, testCase.value)

			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config:      config,
						ExpectError: regexp.MustCompile(`(?s)is_removable cannot be set when intent is .` + testCase.intent + `.*only supports this setting when intent is .required`),
					},
				},
			})
		})
	}
}

// TestUnitResourceMobileAppAssignment_08_Create_Unknown_Intent verifies the case where intent
// is not known at plan time, for example when interpolated from another resource. ModifyPlan
// cannot evaluate the intent then, so the defaulted is_removable must be planned as unknown:
// Terraform requires a value known in the initial plan to be identical in the final plan, so
// planning the static default and then nulling it once intent resolves would be rejected with
// "Provider produced inconsistent final plan".
func TestUnitResourceMobileAppAssignment_08_Create_Unknown_Intent(t *testing.T) {
	_, assignmentMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer assignmentMock.CleanupMockState()

	mocks.SetupUnitTestEnvironment(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfig(t, "resource_unknown_intent.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".unknown_intent").Key("intent").HasValue("available"),
					check.That(resourceType+".unknown_intent").Key("settings.ios_vpp.is_removable").DoesNotExist(),
					checkIsRemovableSent(t, assignmentMock, "iosVpp", false, false),
				),
			},
		},
	})
}
