package graphBetaDeviceAndAppManagementAppAssignment_test

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
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
	httpmock.Activate()

	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()

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

// diagnosticRegexp builds a whitespace tolerant pattern for a diagnostic phrase, since
// Terraform wraps diagnostic text across lines when rendering it.
func diagnosticRegexp(phrase string) *regexp.Regexp {
	fields := strings.Fields(phrase)
	for i, field := range fields {
		fields[i] = regexp.QuoteMeta(field)
	}
	return regexp.MustCompile(`(?s)` + strings.Join(fields, `\s+`))
}

// checkSettingSent asserts on the settings payload the provider actually sent for the given
// settings block, rather than only on what ended up in Terraform state. The whole point of
// this fix is which fields leave the provider, which state alone cannot show.
func checkSettingSent(t *testing.T, assignmentMock *mobileAppAssignmentMocks.MobileAppAssignmentMock, block, field string, wantPresent bool, wantValue any) resource.TestCheckFunc {
	t.Helper()

	return func(*terraform.State) error {
		settings, ok := assignmentMock.SettingsSent(block)
		if !ok {
			return fmt.Errorf("no assignment request was captured for settings block %s", block)
		}

		value, present := settings[field]
		if present != wantPresent {
			return fmt.Errorf("%s present in request = %v, want %v (payload: %v)", field, present, wantPresent, settings)
		}
		if wantPresent && value != wantValue {
			return fmt.Errorf("%s sent as %v, want %v", field, value, wantValue)
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
					checkSettingSent(t, assignmentMock, "iosVpp", "isRemovable", false, nil),
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
					checkSettingSent(t, assignmentMock, "iosStore", "isRemovable", false, nil),
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
					checkSettingSent(t, assignmentMock, "iosLob", "isRemovable", false, nil),
				),
			},
		},
	})
}

// TestUnitResourceMobileAppAssignment_04_Create_Required_IosVpp verifies that explicitly
// setting is_removable alongside intent "required" continues to work, and is still sent.
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
					checkSettingSent(t, assignmentMock, "iosVpp", "isRemovable", true, true),
				),
			},
		},
	})
}

// TestUnitResourceMobileAppAssignment_05_Create_Required_IsRemovable_Omitted verifies that an
// omitted is_removable is not sent even for the one intent that supports it, so that the
// Intune service default applies rather than a value the practitioner never asked for.
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
					check.That(resourceType+".required_omitted").Key("settings.ios_vpp.is_removable").DoesNotExist(),
					checkSettingSent(t, assignmentMock, "iosVpp", "isRemovable", false, nil),
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
					checkSettingSent(t, assignmentMock, "iosVpp", "isRemovable", false, nil),
					checkSettingSent(t, assignmentMock, "iosVpp", "uninstallOnDeviceRemoval", false, nil),
				),
			},
		},
	})
}

// TestUnitResourceMobileAppAssignment_07_Create_Uninstall_IosStore verifies that an uninstall
// intent with an empty Apple settings block sends no settings at all. The service rejects
// uninstallOnDeviceRemoval and preventManagedAppBackup for that intent as well as
// isRemovable, so any schema default on those attributes makes the intent unusable.
func TestUnitResourceMobileAppAssignment_07_Create_Uninstall_IosStore(t *testing.T) {
	_, assignmentMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer assignmentMock.CleanupMockState()

	mocks.SetupUnitTestEnvironment(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfig(t, "resource_uninstall_ios_store.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".uninstall_ios_store").Key("intent").HasValue("uninstall"),
					check.That(resourceType+".uninstall_ios_store").Key("settings.ios_store.is_removable").DoesNotExist(),
					check.That(resourceType+".uninstall_ios_store").Key("settings.ios_store.uninstall_on_device_removal").DoesNotExist(),
					check.That(resourceType+".uninstall_ios_store").Key("settings.ios_store.prevent_managed_app_backup").DoesNotExist(),
					checkSettingSent(t, assignmentMock, "iosStore", "isRemovable", false, nil),
					checkSettingSent(t, assignmentMock, "iosStore", "uninstallOnDeviceRemoval", false, nil),
					checkSettingSent(t, assignmentMock, "iosStore", "preventManagedAppBackup", false, nil),
				),
			},
		},
	})
}

// unsupportedSettingConfig renders an assignment that sets a single Apple settings field
// alongside the given intent.
func unsupportedSettingConfig(intent, block, field, value string) string {
	return fmt.Sprintf(`
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
      %s = %s
    }
  }
}
`, intent, block, field, value)
}

// TestUnitResourceMobileAppAssignment_08_Validate_IsRemovable_UnsupportedIntent verifies that
// explicitly setting is_removable with a non-required intent fails at plan time with a clear
// message, rather than with an opaque HTTP 400 part-way through an apply. Every Apple settings
// block is covered, for both boolean values, across each unsupported intent.
func TestUnitResourceMobileAppAssignment_08_Validate_IsRemovable_UnsupportedIntent(t *testing.T) {
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

			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config:      unsupportedSettingConfig(testCase.intent, testCase.block, "is_removable", testCase.value),
						ExpectError: diagnosticRegexp("`is_removable` can only be set when `intent` is `required`"),
					},
				},
			})
		})
	}
}

// TestUnitResourceMobileAppAssignment_09_Validate_UninstallIntent_UnsupportedSettings verifies
// the same plan time rejection for the two Apple settings the service refuses specifically for
// an uninstall intent.
func TestUnitResourceMobileAppAssignment_09_Validate_UninstallIntent_UnsupportedSettings(t *testing.T) {
	testCases := []struct {
		name  string
		block string
		field string
		value string
	}{
		{name: "ios_lob_uninstall_on_device_removal", block: "ios_lob", field: "uninstall_on_device_removal", value: "true"},
		{name: "ios_store_prevent_managed_app_backup", block: "ios_store", field: "prevent_managed_app_backup", value: "false"},
		{name: "ios_vpp_uninstall_on_device_removal", block: "ios_vpp", field: "uninstall_on_device_removal", value: "false"},
		{name: "ios_vpp_prevent_managed_app_backup", block: "ios_vpp", field: "prevent_managed_app_backup", value: "true"},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			_, assignmentMock := setupMockEnvironment()
			defer httpmock.DeactivateAndReset()
			defer assignmentMock.CleanupMockState()

			mocks.SetupUnitTestEnvironment(t)

			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config: unsupportedSettingConfig("uninstall", testCase.block, testCase.field, testCase.value),
						ExpectError: diagnosticRegexp(
							fmt.Sprintf("`%s` cannot be set when `intent` is `uninstall`", testCase.field)),
					},
				},
			})
		})
	}
}

// TestUnitResourceMobileAppAssignment_10_Create_Unknown_Intent verifies the case where intent
// is not known at plan time, for example when interpolated from another resource. Because the
// affected attributes carry no schema default, an omitted is_removable is simply null in the
// plan whether or not the intent is known, so no unknown value is introduced and no plan
// modification is needed to reconcile the two.
func TestUnitResourceMobileAppAssignment_10_Create_Unknown_Intent(t *testing.T) {
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
					checkSettingSent(t, assignmentMock, "iosVpp", "isRemovable", false, nil),
				),
			},
		},
	})
}
