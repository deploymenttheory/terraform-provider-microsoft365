package graphBetaDeviceConfigurationTemplatesJson_test

import (
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"sort"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/jarcoal/httpmock"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaDeviceConfigurationTemplatesJson "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/device_configuration_templates_json"
	deviceConfigurationMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/device_configuration_templates_json/mocks"
)

var (
	resourceType = graphBetaDeviceConfigurationTemplatesJson.ResourceName
	testResource = graphBetaDeviceConfigurationTemplatesJson.DeviceConfigurationTemplatesJsonTestResource{}
)

func setupMockEnvironment() (*mocks.Mocks, *deviceConfigurationMocks.DeviceConfigurationTemplatesJsonMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	profileMock := &deviceConfigurationMocks.DeviceConfigurationTemplatesJsonMock{}
	profileMock.RegisterMocks()
	return mockClient, profileMock
}

func setupErrorMockEnvironment() (*mocks.Mocks, *deviceConfigurationMocks.DeviceConfigurationTemplatesJsonMock) {
	mockClient, profileMock := setupMockEnvironment()
	profileMock.RegisterErrorMocks()
	return mockClient, profileMock
}

func loadUnitTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/unit/" + filename)
	if err != nil {
		panic("failed to load unit config " + filename + ": " + err.Error())
	}
	return config
}

func importStep() resource.TestStep {
	return resource.TestStep{
		ResourceName:            resourceType + ".test",
		ImportState:             true,
		ImportStateVerify:       true,
		ImportStateVerifyIgnore: []string{"timeouts"},
	}
}

func TestUnitResourceDeviceConfigurationTemplatesJson_01_IosGeneral(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, profileMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer profileMock.CleanupMockState()
	config := loadUnitTestTerraform("resource_01_ios_general.tf")
	updated := strings.ReplaceAll(config, "Device configuration template test", "Updated template description")
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: config, Check: resource.ComposeTestCheckFunc(
				check.That(resourceType+".test").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
				check.That(resourceType+".test").Key("display_name").HasValue("unit-test-ios-general"),
				check.That(resourceType+".test").Key("settings").Exists(),
				check.That(resourceType+".test").Key("role_scope_tag_ids.#").HasValue("1"),
			)},
			importStep(),
			{Config: updated, Check: check.That(resourceType + ".test").Key("description").HasValue("Updated template description")},
			importStep(),
			{Config: updated, PlanOnly: true},
		},
	})
}

func TestUnitResourceDeviceConfigurationTemplatesJson_02_IosFeatures(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, profileMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer profileMock.CleanupMockState()
	config := loadUnitTestTerraform("resource_02_ios_features.tf")
	updated := strings.ReplaceAll(config, "Device configuration template test", "Updated template description")
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: config, Check: resource.ComposeTestCheckFunc(
				check.That(resourceType+".test").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
				check.That(resourceType+".test").Key("display_name").HasValue("unit-test-ios-features"),
				check.That(resourceType+".test").Key("settings").Exists(),
				check.That(resourceType+".test").Key("role_scope_tag_ids.#").HasValue("1"),
			)},
			importStep(),
			{Config: updated, Check: check.That(resourceType + ".test").Key("description").HasValue("Updated template description")},
			importStep(),
			{Config: updated, PlanOnly: true},
		},
	})
}

func TestUnitResourceDeviceConfigurationTemplatesJson_03_IosCustom(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, profileMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer profileMock.CleanupMockState()
	config := loadUnitTestTerraform("resource_03_ios_custom.tf")
	updated := strings.ReplaceAll(config, "Device configuration template test", "Updated template description")
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: config, Check: resource.ComposeTestCheckFunc(
				check.That(resourceType+".test").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
				check.That(resourceType+".test").Key("display_name").HasValue("unit-test-ios-custom"),
				check.That(resourceType+".test").Key("settings").Exists(),
				check.That(resourceType+".test").Key("role_scope_tag_ids.#").HasValue("1"),
			)},
			importStep(),
			{Config: updated, Check: check.That(resourceType + ".test").Key("description").HasValue("Updated template description")},
			importStep(),
			{Config: updated, PlanOnly: true},
		},
	})
}

func TestUnitResourceDeviceConfigurationTemplatesJson_04_IosRoot(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, profileMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer profileMock.CleanupMockState()
	config := loadUnitTestTerraform("resource_04_ios_root.tf")
	updated := strings.ReplaceAll(config, "Device configuration template test", "Updated template description")
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: config, Check: resource.ComposeTestCheckFunc(
				check.That(resourceType+".test").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
				check.That(resourceType+".test").Key("display_name").HasValue("unit-test-ios-root"),
				check.That(resourceType+".test").Key("settings").Exists(),
				check.That(resourceType+".test").Key("role_scope_tag_ids.#").HasValue("1"),
			)},
			importStep(),
			{Config: updated, Check: check.That(resourceType + ".test").Key("description").HasValue("Updated template description")},
			importStep(),
			{Config: updated, PlanOnly: true},
		},
	})
}

func TestUnitResourceDeviceConfigurationTemplatesJson_05_IosScep(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, profileMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer profileMock.CleanupMockState()
	config := loadUnitTestTerraform("resource_05_ios_scep.tf")
	updated := strings.ReplaceAll(config, "Device configuration template test", "Updated template description")
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: config, Check: resource.ComposeTestCheckFunc(
				check.That(resourceType+".test").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
				check.That(resourceType+".test").Key("display_name").HasValue("unit-test-ios-scep"),
				check.That(resourceType+".test").Key("settings").Exists(),
				check.That(resourceType+".test").Key("role_scope_tag_ids.#").HasValue("1"),
			)},
			importStep(),
			{Config: updated, Check: check.That(resourceType + ".test").Key("description").HasValue("Updated template description")},
			importStep(),
			{Config: updated, PlanOnly: true},
		},
	})
}

func TestUnitResourceDeviceConfigurationTemplatesJson_06_MacosCustom(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, profileMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer profileMock.CleanupMockState()
	config := loadUnitTestTerraform("resource_06_macos_custom.tf")
	updated := strings.ReplaceAll(config, "Device configuration template test", "Updated template description")
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: config, Check: resource.ComposeTestCheckFunc(
				check.That(resourceType+".test").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
				check.That(resourceType+".test").Key("display_name").HasValue("unit-test-macos-custom"),
				check.That(resourceType+".test").Key("settings").Exists(),
				check.That(resourceType+".test").Key("role_scope_tag_ids.#").HasValue("1"),
			)},
			importStep(),
			{Config: updated, Check: check.That(resourceType + ".test").Key("description").HasValue("Updated template description")},
			importStep(),
			{Config: updated, PlanOnly: true},
		},
	})
}

func TestUnitResourceDeviceConfigurationTemplatesJson_07_WindowsCustom(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, profileMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer profileMock.CleanupMockState()
	config := loadUnitTestTerraform("resource_07_windows_custom.tf")
	updated := strings.ReplaceAll(config, "Device configuration template test", "Updated template description")
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: config, Check: resource.ComposeTestCheckFunc(
				check.That(resourceType+".test").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
				check.That(resourceType+".test").Key("display_name").HasValue("unit-test-windows-custom"),
				check.That(resourceType+".test").Key("settings").Exists(),
				check.That(resourceType+".test").Key("role_scope_tag_ids.#").HasValue("1"),
			)},
			importStep(),
			{Config: updated, Check: check.That(resourceType + ".test").Key("description").HasValue("Updated template description")},
			importStep(),
			{Config: updated, PlanOnly: true},
		},
	})
}

func TestUnitResourceDeviceConfigurationTemplatesJson_08_WindowsRoot(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, profileMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer profileMock.CleanupMockState()
	config := loadUnitTestTerraform("resource_08_windows_root.tf")
	updated := strings.ReplaceAll(config, "Device configuration template test", "Updated template description")
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: config, Check: resource.ComposeTestCheckFunc(
				check.That(resourceType+".test").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
				check.That(resourceType+".test").Key("display_name").HasValue("unit-test-windows-root"),
				check.That(resourceType+".test").Key("settings").Exists(),
				check.That(resourceType+".test").Key("role_scope_tag_ids.#").HasValue("1"),
			)},
			importStep(),
			{Config: updated, Check: check.That(resourceType + ".test").Key("description").HasValue("Updated template description")},
			importStep(),
			{Config: updated, PlanOnly: true},
		},
	})
}

func TestUnitResourceDeviceConfigurationTemplatesJson_09_AndroidGeneral(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, profileMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer profileMock.CleanupMockState()
	config := loadUnitTestTerraform("resource_09_android_general.tf")
	updated := strings.ReplaceAll(config, "Device configuration template test", "Updated template description")
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: config, Check: resource.ComposeTestCheckFunc(
				check.That(resourceType+".test").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
				check.That(resourceType+".test").Key("display_name").HasValue("unit-test-android-general"),
				check.That(resourceType+".test").Key("settings").Exists(),
				check.That(resourceType+".test").Key("role_scope_tag_ids.#").HasValue("1"),
			)},
			importStep(),
			{Config: updated, Check: check.That(resourceType + ".test").Key("description").HasValue("Updated template description")},
			importStep(),
			{Config: updated, PlanOnly: true},
		},
	})
}

func TestUnitResourceDeviceConfigurationTemplatesJson_10_WindowsEncrypted(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, profileMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer profileMock.CleanupMockState()
	config := loadUnitTestTerraform("resource_10_windows_encrypted.tf")
	updated := strings.ReplaceAll(config, "Device configuration template test", "Updated template description")
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: config, Check: resource.ComposeTestCheckFunc(
				check.That(resourceType+".test").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
				check.That(resourceType+".test").Key("display_name").HasValue("unit-test-windows-encrypted"),
				check.That(resourceType+".test").Key("settings").Exists(),
				check.That(resourceType+".test").Key("role_scope_tag_ids.#").HasValue("1"),
			)},
			importStep(),
			{Config: updated, Check: check.That(resourceType + ".test").Key("description").HasValue("Updated template description")},
			importStep(),
			{Config: updated, PlanOnly: true},
		},
	})
}

func TestUnitResourceDeviceConfigurationTemplatesJson_11_IosNestedFolder(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, profileMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer profileMock.CleanupMockState()
	config := loadUnitTestTerraform("resource_11_ios_nested_folder.tf")
	updated := strings.ReplaceAll(config, "Device configuration template test", "Updated template description")
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: config, Check: resource.ComposeTestCheckFunc(
				check.That(resourceType+".test").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
				check.That(resourceType+".test").Key("display_name").HasValue("unit-test-ios-nested-folder"),
				check.That(resourceType+".test").Key("settings").Exists(),
				check.That(resourceType+".test").Key("role_scope_tag_ids.#").HasValue("1"),
			)},
			importStep(),
			{Config: updated, Check: check.That(resourceType + ".test").Key("description").HasValue("Updated template description")},
			importStep(),
			{Config: updated, PlanOnly: true},
		},
	})
}

func TestUnitResourceDeviceConfigurationTemplatesJson_12_MacosRoot(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, profileMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer profileMock.CleanupMockState()
	config := loadUnitTestTerraform("resource_12_macos_root.tf")
	updated := strings.ReplaceAll(config, "Device configuration template test", "Updated template description")
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: config, Check: resource.ComposeTestCheckFunc(
				check.That(resourceType+".test").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
				check.That(resourceType+".test").Key("display_name").HasValue("unit-test-macos-root"),
				check.That(resourceType+".test").Key("settings").Exists(),
				check.That(resourceType+".test").Key("role_scope_tag_ids.#").HasValue("1"),
			)},
			importStep(),
			{Config: updated, Check: check.That(resourceType + ".test").Key("description").HasValue("Updated template description")},
			importStep(),
			{Config: updated, PlanOnly: true},
		},
	})
}

func TestUnitResourceDeviceConfigurationTemplatesJson_13_MacosScep(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, profileMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer profileMock.CleanupMockState()
	config := loadUnitTestTerraform("resource_13_macos_scep.tf")
	updated := strings.ReplaceAll(config, "Device configuration template test", "Updated template description")
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: config, Check: resource.ComposeTestCheckFunc(
				check.That(resourceType+".test").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
				check.That(resourceType+".test").Key("display_name").HasValue("unit-test-macos-scep"),
				check.That(resourceType+".test").Key("settings").Exists(),
				check.That(resourceType+".test").Key("role_scope_tag_ids.#").HasValue("1"),
			)},
			importStep(),
			{Config: updated, Check: check.That(resourceType + ".test").Key("description").HasValue("Updated template description")},
			importStep(),
			{Config: updated, PlanOnly: true},
		},
	})
}

func TestUnitResourceDeviceConfigurationTemplatesJson_14_AssignmentTransitions(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, profileMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer profileMock.CleanupMockState()
	minimal := loadUnitTestTerraform("resource_14_assignments_minimal.tf")
	maximal := loadUnitTestTerraform("resource_14_assignments_maximal.tf")
	empty := loadUnitTestTerraform("resource_14_assignments_empty.tf")
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: minimal},
			{Config: maximal, Check: check.That(resourceType + ".test").Key("assignments.#").HasValue("4")},
			importStep(),
			{Config: empty, Check: check.That(resourceType + ".test").Key("assignments.#").HasValue("0")},
			{Config: maximal},
			{Config: minimal, Check: check.That(resourceType + ".test").Key("assignments.#").DoesNotExist()},
			importStep(),
		},
	})
}

func TestUnitResourceDeviceConfigurationTemplatesJson_15_ExplicitResets(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, profileMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer profileMock.CleanupMockState()
	config := loadUnitTestTerraform("resource_01_ios_general.tf")
	updated := loadUnitTestTerraform("resource_15_explicit_resets.tf")
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: config},
			{Config: updated, Check: check.That(resourceType + ".test").Key("settings").MatchesRegex(regexp.MustCompile(`"cameraBlocked":false`))},
			importStep(),
			{Config: config, Check: check.That(resourceType + ".test").Key("settings").MatchesRegex(regexp.MustCompile(`"cameraBlocked":true`))},
			importStep(),
		},
	})
}

func TestUnitResourceDeviceConfigurationTemplatesJson_16_InvalidJSON(t *testing.T) {
	for _, tc := range []struct{ name, message string }{
		{"null", "settings must be a JSON object"},
		{"array", "decode device configuration settings"},
		{"missing_type", "root @odata.type"},
		{"invalid_type", "root @odata.type"},
		{"metadata", "displayName belongs"},
		{"assignments", "assignments belongs"},
		{"read_only", "version belongs"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			mocks.SetupUnitTestEnvironment(t)
			_, profileMock := setupMockEnvironment()
			defer httpmock.DeactivateAndReset()
			defer profileMock.CleanupMockState()
			config := loadUnitTestTerraform("resource_16_invalid_" + tc.name + ".tf")
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
				Steps:                    []resource.TestStep{{Config: config, ExpectError: regexp.MustCompile(tc.message)}},
			})
		})
	}
}

func TestUnitResourceDeviceConfigurationTemplatesJson_17_PermissionDenied(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, profileMock := setupErrorMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer profileMock.CleanupMockState()
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps:                    []resource.TestStep{{Config: loadUnitTestTerraform("resource_03_ios_custom.tf"), ExpectError: regexp.MustCompile("403|Forbidden|permissions")}},
	})
}

func TestUnitResourceDeviceConfigurationTemplatesJson_18_RemoteDrift(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, profileMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer profileMock.CleanupMockState()
	config := loadUnitTestTerraform("resource_01_ios_general.tf")
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: config},
			{PreConfig: func() { profileMock.SetRemoteProperty("cameraBlocked", false) }, Config: config, PlanOnly: true, ExpectNonEmptyPlan: true},
			{Config: config, Check: check.That(resourceType + ".test").Key("settings").MatchesRegex(regexp.MustCompile(`"cameraBlocked":true`))},
		},
	})
}

func TestUnitResourceDeviceConfigurationTemplatesJson_19_RemoteDeletion(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, profileMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer profileMock.CleanupMockState()
	config := loadUnitTestTerraform("resource_03_ios_custom.tf")
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: config},
			{PreConfig: profileMock.CleanupMockState, Config: config, PlanOnly: true, ExpectNonEmptyPlan: true},
		},
	})
}

func TestUnitResourceDeviceConfigurationTemplatesJson_20_AssignmentFailureRetainsID(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, profileMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer profileMock.CleanupMockState()
	httpmock.RegisterResponder("POST", `=~^https://graph\.microsoft\.com/beta/deviceManagement/deviceConfigurations/[^/]+.*$`,
		httpmock.NewStringResponder(400, `{"error":{"code":"BadRequest","message":"Assignment rejected"}}`))
	config := loadUnitTestTerraform("resource_20_assignment_failure.tf")
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps:                    []resource.TestStep{{Config: config, ExpectError: regexp.MustCompile("400|Bad Request")}},
	})
	deleted := false
	for call, count := range httpmock.GetCallCountInfo() {
		if strings.HasPrefix(call, "DELETE ") && count > 0 {
			deleted = true
		}
	}
	if !deleted {
		t.Fatal("assignment failure lost the profile ID; cleanup could not delete it")
	}
}

func TestUnitResourceDeviceConfigurationTemplatesJson_21_FormattedJSON(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, profileMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer profileMock.CleanupMockState()
	config := loadUnitTestTerraform("resource_21_formatted_json.tf")
	updated := strings.ReplaceAll(config, "Device configuration template test", "Updated template description")
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: config, Check: check.That(resourceType + ".test").Key("settings").MatchesRegex(regexp.MustCompile(`(?s)^\{\n  "`))},
			{Config: config, PlanOnly: true},
			{Config: updated},
			{Config: updated, PlanOnly: true},
		},
	})
}

// Graph never returns Wi-Fi pre-shared keys on import. Verify every other JSON
// property against the captured response rather than ignoring the entire payload.
func checkImportedWifiSettings(filename string) resource.ImportStateCheckFunc {
	return func(states []*terraform.InstanceState) error {
		content, err := helpers.ParseJSONFile("tests/responses/validate_get/" + filename)
		if err != nil {
			return err
		}
		var expected map[string]any
		if err := json.Unmarshal([]byte(content), &expected); err != nil {
			return err
		}
		for _, state := range states {
			var actual map[string]any
			if err := json.Unmarshal([]byte(state.Attributes["settings"]), &actual); err != nil {
				return err
			}
			if !reflect.DeepEqual(expected, actual) {
				return fmt.Errorf("imported Wi-Fi settings differ from captured writable settings")
			}
		}
		if len(states) != 1 {
			return fmt.Errorf("expected one imported profile, got %d", len(states))
		}
		return nil
	}
}

func TestUnitResourceDeviceConfigurationTemplatesJson_22_WindowsOmaInteger(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, profileMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer profileMock.CleanupMockState()
	config := loadUnitTestTerraform("resource_22_windows_oma_integer.tf")
	updated := strings.ReplaceAll(config, "Device configuration template test", "Updated template description")
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: config, Check: check.That(resourceType + ".test").Key("settings").Exists()},
			importStep(),
			{Config: updated, Check: check.That(resourceType + ".test").Key("description").HasValue("Updated template description")},
			importStep(),
			{Config: updated, PlanOnly: true},
		},
	})
}

func TestUnitResourceDeviceConfigurationTemplatesJson_23_WindowsOmaBoolean(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, profileMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer profileMock.CleanupMockState()
	config := loadUnitTestTerraform("resource_23_windows_oma_boolean.tf")
	updated := strings.ReplaceAll(config, "Device configuration template test", "Updated template description")
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: config, Check: check.That(resourceType + ".test").Key("settings").Exists()},
			importStep(),
			{Config: updated, Check: check.That(resourceType + ".test").Key("description").HasValue("Updated template description")},
			importStep(),
			{Config: updated, PlanOnly: true},
		},
	})
}

func TestUnitResourceDeviceConfigurationTemplatesJson_24_WindowsOmaFloatingpoint(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, profileMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer profileMock.CleanupMockState()
	config := loadUnitTestTerraform("resource_24_windows_oma_floatingpoint.tf")
	updated := strings.ReplaceAll(config, "Device configuration template test", "Updated template description")
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: config, Check: check.That(resourceType + ".test").Key("settings").Exists()},
			importStep(),
			{Config: updated, Check: check.That(resourceType + ".test").Key("description").HasValue("Updated template description")},
			importStep(),
			{Config: updated, PlanOnly: true},
		},
	})
}

func TestUnitResourceDeviceConfigurationTemplatesJson_25_WindowsOmaDatetime(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, profileMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer profileMock.CleanupMockState()
	config := loadUnitTestTerraform("resource_25_windows_oma_datetime.tf")
	updated := strings.ReplaceAll(config, "Device configuration template test", "Updated template description")
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: config, Check: check.That(resourceType + ".test").Key("settings").Exists()},
			importStep(),
			{Config: updated, Check: check.That(resourceType + ".test").Key("description").HasValue("Updated template description")},
			importStep(),
			{Config: updated, PlanOnly: true},
		},
	})
}

func TestUnitResourceDeviceConfigurationTemplatesJson_26_PlaintextString(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, profileMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer profileMock.CleanupMockState()
	config := loadUnitTestTerraform("resource_26_plaintext_string.tf")
	updated := strings.ReplaceAll(config, "Device configuration template test", "Updated template description")
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: config, Check: check.That(resourceType + ".test").Key("settings").Exists()},
			importStep(),
			{Config: updated, Check: check.That(resourceType + ".test").Key("description").HasValue("Updated template description")},
			importStep(),
			{Config: updated, PlanOnly: true},
		},
	})
}

func TestUnitResourceDeviceConfigurationTemplatesJson_27_PlaintextBase64(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, profileMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer profileMock.CleanupMockState()
	config := loadUnitTestTerraform("resource_27_plaintext_base64.tf")
	updated := strings.ReplaceAll(config, "Device configuration template test", "Updated template description")
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: config, Check: check.That(resourceType + ".test").Key("settings").Exists()},
			importStep(),
			{Config: updated, Check: check.That(resourceType + ".test").Key("description").HasValue("Updated template description")},
			importStep(),
			{Config: updated, PlanOnly: true},
		},
	})
}

func TestUnitResourceDeviceConfigurationTemplatesJson_28_PlaintextStringxml(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, profileMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer profileMock.CleanupMockState()
	config := loadUnitTestTerraform("resource_28_plaintext_stringxml.tf")
	updated := strings.ReplaceAll(config, "Device configuration template test", "Updated template description")
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: config, Check: check.That(resourceType + ".test").Key("settings").Exists()},
			importStep(),
			{Config: updated, Check: check.That(resourceType + ".test").Key("description").HasValue("Updated template description")},
			importStep(),
			{Config: updated, PlanOnly: true},
		},
	})
}

func TestUnitResourceDeviceConfigurationTemplatesJson_29_WindowsUpdates(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, profileMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer profileMock.CleanupMockState()
	config := loadUnitTestTerraform("resource_29_windows_updates.tf")
	updated := strings.ReplaceAll(config, "Device configuration template test", "Updated template description")
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: config, Check: check.That(resourceType + ".test").Key("settings").Exists()},
			importStep(),
			{Config: updated, Check: check.That(resourceType + ".test").Key("description").HasValue("Updated template description")},
			importStep(),
			{Config: updated, PlanOnly: true},
		},
	})
}

func TestUnitResourceDeviceConfigurationTemplatesJson_30_IosWifi(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, profileMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer profileMock.CleanupMockState()
	config := loadUnitTestTerraform("resource_30_ios_wifi.tf")
	updated := strings.ReplaceAll(config, "Device configuration template test", "Updated template description")
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: config, Check: check.That(resourceType + ".test").Key("settings").Exists()},
			{ResourceName: resourceType + ".test", ImportState: true, ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"timeouts", "settings"},
				ImportStateCheck:        checkImportedWifiSettings("get_device_configuration_ios_wifi_import_settings.json")},
			{Config: updated, Check: check.That(resourceType + ".test").Key("description").HasValue("Updated template description")},
			{ResourceName: resourceType + ".test", ImportState: true, ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"timeouts", "settings"},
				ImportStateCheck:        checkImportedWifiSettings("get_device_configuration_ios_wifi_import_settings.json")},
			{Config: updated, PlanOnly: true},
		},
	})
}

func TestUnitResourceDeviceConfigurationTemplatesJson_31_MacosWifi(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, profileMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer profileMock.CleanupMockState()
	config := loadUnitTestTerraform("resource_31_macos_wifi.tf")
	updated := strings.ReplaceAll(config, "Device configuration template test", "Updated template description")
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: config, Check: check.That(resourceType + ".test").Key("settings").Exists()},
			{ResourceName: resourceType + ".test", ImportState: true, ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"timeouts", "settings"},
				ImportStateCheck:        checkImportedWifiSettings("get_device_configuration_macos_wifi_import_settings.json")},
			{Config: updated, Check: check.That(resourceType + ".test").Key("description").HasValue("Updated template description")},
			{ResourceName: resourceType + ".test", ImportState: true, ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"timeouts", "settings"},
				ImportStateCheck:        checkImportedWifiSettings("get_device_configuration_macos_wifi_import_settings.json")},
			{Config: updated, PlanOnly: true},
		},
	})
}

func TestUnitResourceDeviceConfigurationTemplatesJson_32_WindowsWifi(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, profileMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer profileMock.CleanupMockState()
	config := loadUnitTestTerraform("resource_32_windows_wifi.tf")
	updated := strings.ReplaceAll(config, "Device configuration template test", "Updated template description")
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: config, Check: check.That(resourceType + ".test").Key("settings").Exists()},
			{ResourceName: resourceType + ".test", ImportState: true, ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"timeouts", "settings"},
				ImportStateCheck:        checkImportedWifiSettings("get_device_configuration_windows_wifi_import_settings.json")},
			{Config: updated, Check: check.That(resourceType + ".test").Key("description").HasValue("Updated template description")},
			{ResourceName: resourceType + ".test", ImportState: true, ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"timeouts", "settings"},
				ImportStateCheck:        checkImportedWifiSettings("get_device_configuration_windows_wifi_import_settings.json")},
			{Config: updated, PlanOnly: true},
		},
	})
}

func TestUnitResourceDeviceConfigurationTemplatesJson_33_AndroidOwnerWifi(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, profileMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer profileMock.CleanupMockState()
	config := loadUnitTestTerraform("resource_33_android_owner_wifi.tf")
	updated := strings.ReplaceAll(config, "Device configuration template test", "Updated template description")
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: config, Check: check.That(resourceType + ".test").Key("settings").Exists()},
			{ResourceName: resourceType + ".test", ImportState: true, ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"timeouts", "settings"},
				ImportStateCheck:        checkImportedWifiSettings("get_device_configuration_android_owner_wifi_import_settings.json")},
			{Config: updated, Check: check.That(resourceType + ".test").Key("description").HasValue("Updated template description")},
			{ResourceName: resourceType + ".test", ImportState: true, ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"timeouts", "settings"},
				ImportStateCheck:        checkImportedWifiSettings("get_device_configuration_android_owner_wifi_import_settings.json")},
			{Config: updated, PlanOnly: true},
		},
	})
}

func TestUnitResourceDeviceConfigurationTemplatesJson_34_AndroidWorkWifi(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, profileMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer profileMock.CleanupMockState()
	config := loadUnitTestTerraform("resource_34_android_work_wifi.tf")
	updated := strings.ReplaceAll(config, "Device configuration template test", "Updated template description")
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: config, Check: check.That(resourceType + ".test").Key("settings").Exists()},
			{ResourceName: resourceType + ".test", ImportState: true, ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"timeouts", "settings"},
				ImportStateCheck:        checkImportedWifiSettings("get_device_configuration_android_work_wifi_import_settings.json")},
			{Config: updated, Check: check.That(resourceType + ".test").Key("description").HasValue("Updated template description")},
			{ResourceName: resourceType + ".test", ImportState: true, ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"timeouts", "settings"},
				ImportStateCheck:        checkImportedWifiSettings("get_device_configuration_android_work_wifi_import_settings.json")},
			{Config: updated, PlanOnly: true},
		},
	})
}

func TestUnitResourceDeviceConfigurationTemplatesJson_35_AospWifi(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, profileMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer profileMock.CleanupMockState()
	config := loadUnitTestTerraform("resource_35_aosp_wifi.tf")
	updated := strings.ReplaceAll(config, "Device configuration template test", "Updated template description")
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: config, Check: check.That(resourceType + ".test").Key("settings").Exists()},
			{ResourceName: resourceType + ".test", ImportState: true, ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"timeouts", "settings"},
				ImportStateCheck:        checkImportedWifiSettings("get_device_configuration_aosp_wifi_import_settings.json")},
			{Config: updated, Check: check.That(resourceType + ".test").Key("description").HasValue("Updated template description")},
			{ResourceName: resourceType + ".test", ImportState: true, ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"timeouts", "settings"},
				ImportStateCheck:        checkImportedWifiSettings("get_device_configuration_aosp_wifi_import_settings.json")},
			{Config: updated, PlanOnly: true},
		},
	})
}

func TestUnitResourceDeviceConfigurationTemplatesJson_36_IosEnterpriseWifi(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, profileMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer profileMock.CleanupMockState()
	config := loadUnitTestTerraform("resource_36_ios_enterprise_wifi.tf")
	updated := strings.ReplaceAll(config, "Device configuration template test", "Updated template description")
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: config, Check: check.That(resourceType + ".test").Key("settings").Exists()},
			importStep(),
			{Config: updated, Check: check.That(resourceType + ".test").Key("description").HasValue("Updated template description")},
			importStep(),
			{Config: updated, PlanOnly: true},
		},
	})
}

func TestUnitResourceDeviceConfigurationTemplatesJson_37_MacosEnterpriseWifi(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, profileMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer profileMock.CleanupMockState()
	config := loadUnitTestTerraform("resource_37_macos_enterprise_wifi.tf")
	updated := strings.ReplaceAll(config, "Device configuration template test", "Updated template description")
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: config, Check: check.That(resourceType + ".test").Key("settings").Exists()},
			importStep(),
			{Config: updated, Check: check.That(resourceType + ".test").Key("description").HasValue("Updated template description")},
			importStep(),
			{Config: updated, PlanOnly: true},
		},
	})
}

func TestUnitResourceDeviceConfigurationTemplatesJson_38_WindowsEnterpriseWifi(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, profileMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer profileMock.CleanupMockState()
	config := loadUnitTestTerraform("resource_38_windows_enterprise_wifi.tf")
	updated := strings.ReplaceAll(config, "Device configuration template test", "Updated template description")
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: config, Check: check.That(resourceType + ".test").Key("settings").Exists()},
			importStep(),
			{Config: updated, Check: check.That(resourceType + ".test").Key("description").HasValue("Updated template description")},
			importStep(),
			{Config: updated, PlanOnly: true},
		},
	})
}

func TestUnitResourceDeviceConfigurationTemplatesJson_39_AndroidOwnerEnterpriseWifi(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, profileMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer profileMock.CleanupMockState()
	config := loadUnitTestTerraform("resource_39_android_owner_enterprise_wifi.tf")
	updated := strings.ReplaceAll(config, "Device configuration template test", "Updated template description")
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: config, Check: check.That(resourceType + ".test").Key("settings").Exists()},
			importStep(),
			{Config: updated, Check: check.That(resourceType + ".test").Key("description").HasValue("Updated template description")},
			importStep(),
			{Config: updated, PlanOnly: true},
		},
	})
}

func TestUnitResourceDeviceConfigurationTemplatesJson_40_AndroidWorkEnterpriseWifi(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, profileMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer profileMock.CleanupMockState()
	config := loadUnitTestTerraform("resource_40_android_work_enterprise_wifi.tf")
	updated := strings.ReplaceAll(config, "Device configuration template test", "Updated template description")
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: config, Check: check.That(resourceType + ".test").Key("settings").Exists()},
			importStep(),
			{Config: updated, Check: check.That(resourceType + ".test").Key("description").HasValue("Updated template description")},
			importStep(),
			{Config: updated, PlanOnly: true},
		},
	})
}

func TestUnitResourceDeviceConfigurationTemplatesJson_41_AospEnterpriseWifi(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, profileMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer profileMock.CleanupMockState()
	config := loadUnitTestTerraform("resource_41_aosp_enterprise_wifi.tf")
	updated := strings.ReplaceAll(config, "Device configuration template test", "Updated template description")
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: config, Check: check.That(resourceType + ".test").Key("settings").Exists()},
			importStep(),
			{Config: updated, Check: check.That(resourceType + ".test").Key("description").HasValue("Updated template description")},
			importStep(),
			{Config: updated, PlanOnly: true},
		},
	})
}

func TestUnitResourceDeviceConfigurationTemplatesJson_42_IosVpn(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, profileMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer profileMock.CleanupMockState()
	config := loadUnitTestTerraform("resource_42_ios_vpn.tf")
	updated := strings.ReplaceAll(config, "Device configuration template test", "Updated template description")
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: config, Check: check.That(resourceType + ".test").Key("settings").Exists()},
			importStep(),
			{Config: updated, Check: check.That(resourceType + ".test").Key("description").HasValue("Updated template description")},
			importStep(),
			{Config: updated, PlanOnly: true},
		},
	})
}

func TestUnitResourceDeviceConfigurationTemplatesJson_43_MacosVpn(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, profileMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer profileMock.CleanupMockState()
	config := loadUnitTestTerraform("resource_43_macos_vpn.tf")
	updated := strings.ReplaceAll(config, "Device configuration template test", "Updated template description")
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: config, Check: check.That(resourceType + ".test").Key("settings").Exists()},
			importStep(),
			{Config: updated, Check: check.That(resourceType + ".test").Key("description").HasValue("Updated template description")},
			importStep(),
			{Config: updated, PlanOnly: true},
		},
	})
}

func TestUnitResourceDeviceConfigurationTemplatesJson_44_WindowsVpn(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, profileMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer profileMock.CleanupMockState()
	config := loadUnitTestTerraform("resource_44_windows_vpn.tf")
	updated := strings.ReplaceAll(config, "Device configuration template test", "Updated template description")
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: config, Check: check.That(resourceType + ".test").Key("settings").Exists()},
			importStep(),
			{Config: updated, Check: check.That(resourceType + ".test").Key("description").HasValue("Updated template description")},
			importStep(),
			{Config: updated, PlanOnly: true},
		},
	})
}

func TestUnitResourceDeviceConfigurationTemplatesJson_45_AndroidOwnerVpn(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, profileMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer profileMock.CleanupMockState()
	config := loadUnitTestTerraform("resource_45_android_owner_vpn.tf")
	updated := strings.ReplaceAll(config, "Device configuration template test", "Updated template description")
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: config, Check: check.That(resourceType + ".test").Key("settings").Exists()},
			importStep(),
			{Config: updated, Check: check.That(resourceType + ".test").Key("description").HasValue("Updated template description")},
			importStep(),
			{Config: updated, PlanOnly: true},
		},
	})
}

func TestUnitResourceDeviceConfigurationTemplatesJson_46_IosEmail(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, profileMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer profileMock.CleanupMockState()
	config := loadUnitTestTerraform("resource_46_ios_email.tf")
	updated := strings.ReplaceAll(config, "Device configuration template test", "Updated template description")
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: config, Check: check.That(resourceType + ".test").Key("settings").Exists()},
			importStep(),
			{Config: updated, Check: check.That(resourceType + ".test").Key("description").HasValue("Updated template description")},
			importStep(),
			{Config: updated, PlanOnly: true},
		},
	})
}

func TestUnitResourceDeviceConfigurationTemplatesJson_47_IosPkcs(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, profileMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer profileMock.CleanupMockState()
	config := loadUnitTestTerraform("resource_47_ios_pkcs.tf")
	updated := strings.ReplaceAll(config, "Device configuration template test", "Updated template description")
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: config, Check: check.That(resourceType + ".test").Key("settings").Exists()},
			importStep(),
			{Config: updated, Check: check.That(resourceType + ".test").Key("description").HasValue("Updated template description")},
			importStep(),
			{Config: updated, PlanOnly: true},
		},
	})
}

func TestUnitResourceDeviceConfigurationTemplatesJson_48_MacosPkcs(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, profileMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer profileMock.CleanupMockState()
	config := loadUnitTestTerraform("resource_48_macos_pkcs.tf")
	updated := strings.ReplaceAll(config, "Device configuration template test", "Updated template description")
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: config, Check: check.That(resourceType + ".test").Key("settings").Exists()},
			importStep(),
			{Config: updated, Check: check.That(resourceType + ".test").Key("description").HasValue("Updated template description")},
			importStep(),
			{Config: updated, PlanOnly: true},
		},
	})
}

func TestUnitResourceDeviceConfigurationTemplatesJson_49_WindowsPkcs(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, profileMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer profileMock.CleanupMockState()
	config := loadUnitTestTerraform("resource_49_windows_pkcs.tf")
	updated := strings.ReplaceAll(config, "Device configuration template test", "Updated template description")
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: config, Check: check.That(resourceType + ".test").Key("settings").Exists()},
			importStep(),
			{Config: updated, Check: check.That(resourceType + ".test").Key("description").HasValue("Updated template description")},
			importStep(),
			{Config: updated, PlanOnly: true},
		},
	})
}

func TestUnitResourceDeviceConfigurationTemplatesJson_50_AndroidWorkPkcs(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, profileMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer profileMock.CleanupMockState()
	config := loadUnitTestTerraform("resource_50_android_work_pkcs.tf")
	updated := strings.ReplaceAll(config, "Device configuration template test", "Updated template description")
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: config, Check: check.That(resourceType + ".test").Key("settings").Exists()},
			importStep(),
			{Config: updated, Check: check.That(resourceType + ".test").Key("description").HasValue("Updated template description")},
			importStep(),
			{Config: updated, PlanOnly: true},
		},
	})
}

func TestUnitResourceDeviceConfigurationTemplatesJson_51_MacosPreferences(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, profileMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer profileMock.CleanupMockState()
	config := loadUnitTestTerraform("resource_51_macos_preferences.tf")
	updated := strings.ReplaceAll(config, "Device configuration template test", "Updated template description")
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: config, Check: check.That(resourceType + ".test").Key("settings").Exists()},
			importStep(),
			{Config: updated, Check: check.That(resourceType + ".test").Key("description").HasValue("Updated template description")},
			importStep(),
			{Config: updated, PlanOnly: true},
		},
	})
}

func TestUnitResourceDeviceConfigurationTemplatesJson_52_MacosUpdates(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, profileMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer profileMock.CleanupMockState()
	config := loadUnitTestTerraform("resource_52_macos_updates.tf")
	updated := strings.ReplaceAll(config, "Device configuration template test", "Updated template description")
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: config, Check: check.That(resourceType + ".test").Key("settings").Exists()},
			importStep(),
			{Config: updated, Check: check.That(resourceType + ".test").Key("description").HasValue("Updated template description")},
			importStep(),
			{Config: updated, PlanOnly: true},
		},
	})
}

func TestUnitResourceDeviceConfigurationTemplatesJson_53_AospRestrictions(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, profileMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer profileMock.CleanupMockState()
	config := loadUnitTestTerraform("resource_53_aosp_restrictions.tf")
	updated := strings.ReplaceAll(config, "Device configuration template test", "Updated template description")
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: config, Check: check.That(resourceType + ".test").Key("settings").Exists()},
			importStep(),
			{Config: updated, Check: check.That(resourceType + ".test").Key("description").HasValue("Updated template description")},
			importStep(),
			{Config: updated, PlanOnly: true},
		},
	})
}

func TestUnitResourceDeviceConfigurationTemplatesJson_54_SingletonRelationships(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, profileMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer profileMock.CleanupMockState()
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: loadUnitTestTerraform("resource_54_singleton_initial.tf")},
			importStep(),
			{Config: loadUnitTestTerraform("resource_54_singleton_replaced.tf")},
			importStep(),
			{Config: loadUnitTestTerraform("resource_54_singleton_removed.tf")},
			importStep(),
			{Config: loadUnitTestTerraform("resource_54_singleton_initial.tf")},
			importStep(),
		},
	})
}

func TestUnitResourceDeviceConfigurationTemplatesJson_55_CollectionRelationships(t *testing.T) {
	var settingsBeforeImport string
	mocks.SetupUnitTestEnvironment(t)
	_, profileMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer profileMock.CleanupMockState()
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: loadUnitTestTerraform("resource_55_collection_initial.tf"), Check: func(state *terraform.State) error {
				settingsBeforeImport = state.RootModule().Resources[resourceType+".test"].Primary.Attributes["settings"]
				return nil
			}},
			{ResourceName: resourceType + ".test", ImportState: true, ImportStateVerify: true, ImportStateVerifyIgnore: []string{"timeouts", "settings"}, ImportStateCheck: checkImportedJSONSettings(&settingsBeforeImport)},
			{Config: loadUnitTestTerraform("resource_55_collection_multiple.tf"), Check: func(state *terraform.State) error {
				settingsBeforeImport = state.RootModule().Resources[resourceType+".test"].Primary.Attributes["settings"]
				return nil
			}},
			{ResourceName: resourceType + ".test", ImportState: true, ImportStateVerify: true, ImportStateVerifyIgnore: []string{"timeouts", "settings"}, ImportStateCheck: checkImportedJSONSettings(&settingsBeforeImport)},
			{Config: loadUnitTestTerraform("resource_55_collection_reordered.tf"), Check: func(state *terraform.State) error {
				settingsBeforeImport = state.RootModule().Resources[resourceType+".test"].Primary.Attributes["settings"]
				return nil
			}},
			{ResourceName: resourceType + ".test", ImportState: true, ImportStateVerify: true, ImportStateVerifyIgnore: []string{"timeouts", "settings"}, ImportStateCheck: checkImportedJSONSettings(&settingsBeforeImport)},
			{Config: loadUnitTestTerraform("resource_55_collection_initial.tf"), Check: func(state *terraform.State) error {
				settingsBeforeImport = state.RootModule().Resources[resourceType+".test"].Primary.Attributes["settings"]
				return nil
			}},
			{ResourceName: resourceType + ".test", ImportState: true, ImportStateVerify: true, ImportStateVerifyIgnore: []string{"timeouts", "settings"}, ImportStateCheck: checkImportedJSONSettings(&settingsBeforeImport)},
			{Config: loadUnitTestTerraform("resource_55_collection_removed.tf"), Check: func(state *terraform.State) error {
				settingsBeforeImport = state.RootModule().Resources[resourceType+".test"].Primary.Attributes["settings"]
				return nil
			}},
			{ResourceName: resourceType + ".test", ImportState: true, ImportStateVerify: true, ImportStateVerifyIgnore: []string{"timeouts", "settings"}, ImportStateCheck: checkImportedJSONSettings(&settingsBeforeImport)},
			{Config: loadUnitTestTerraform("resource_55_collection_multiple.tf"), Check: func(state *terraform.State) error {
				settingsBeforeImport = state.RootModule().Resources[resourceType+".test"].Primary.Attributes["settings"]
				return nil
			}},
			{ResourceName: resourceType + ".test", ImportState: true, ImportStateVerify: true, ImportStateVerifyIgnore: []string{"timeouts", "settings"}, ImportStateCheck: checkImportedJSONSettings(&settingsBeforeImport)},
		},
	})
}

func checkImportedJSONSettings(expected *string) resource.ImportStateCheckFunc {
	return func(states []*terraform.InstanceState) error {
		if len(states) != 1 {
			return fmt.Errorf("expected one imported profile, got %d", len(states))
		}
		var before, after map[string]any
		if err := json.Unmarshal([]byte(*expected), &before); err != nil {
			return err
		}
		if err := json.Unmarshal([]byte(states[0].Attributes["settings"]), &after); err != nil {
			return err
		}
		for _, settings := range []map[string]any{before, after} {
			for key, value := range settings {
				if !strings.HasSuffix(key, "@odata.bind") {
					continue
				}
				if references, ok := value.([]any); ok {
					sort.Slice(references, func(i, j int) bool { return references[i].(string) < references[j].(string) })
				}
			}
		}
		if !reflect.DeepEqual(before, after) {
			return fmt.Errorf("imported settings differ from applied settings after sorting unordered reference collections")
		}
		return nil
	}
}

func TestUnitResourceDeviceConfigurationTemplatesJson_56_RefreshErrorsPreserveState(t *testing.T) {
	for _, tc := range []struct {
		name, suffix string
		status       int
		body         string
	}{
		{"bad_request", "", 400, `{"error":{"code":"TestFailure","message":"Refresh failed"}}`},
		{"forbidden", "", 403, `{"error":{"code":"TestFailure","message":"Refresh failed"}}`},
		{"assignments_not_found", "/assignments", 404, `{"error":{"code":"TestFailure","message":"Refresh failed"}}`},
		{"empty_profile", "", 200, ""},
		{"missing_metadata", "", 200, `{}`},
		{"empty_assignments", "/assignments", 200, ""},
		{"missing_target", "/assignments", 200, `{"value":[{"id":"invalid"}]}`},
		{"unknown_target", "/assignments", 200, `{"value":[{"target":{"@odata.type":"#microsoft.graph.unknownTarget"}}]}`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			mocks.SetupUnitTestEnvironment(t)
			_, profileMock := setupMockEnvironment()
			defer httpmock.DeactivateAndReset()
			defer profileMock.CleanupMockState()
			var id string
			config := loadUnitTestTerraform("resource_03_ios_custom.tf")
			endpoint := "https://graph.microsoft.com/beta/deviceManagement/deviceConfigurations/"
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{Config: config, Check: func(state *terraform.State) error {
						id = state.RootModule().Resources[resourceType+".test"].Primary.ID
						return nil
					}},
					{PreConfig: func() {
						httpmock.RegisterResponder("GET", endpoint+id+tc.suffix, httpmock.NewStringResponder(tc.status, tc.body))
					}, Config: config, ExpectError: regexp.MustCompile("Refresh failed|TestFailure|Error")},
					{PreConfig: func() { httpmock.RegisterResponder("GET", endpoint+id+tc.suffix, nil) }, Config: config, Check: func(state *terraform.State) error {
						if state.RootModule().Resources[resourceType+".test"].Primary.ID != id {
							return fmt.Errorf("refresh error caused replacement of existing profile")
						}
						return nil
					}},
				},
			})
		})
	}
}

func TestUnitResourceDeviceConfigurationTemplatesJson_57_RemoteWifiKeyCleared(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, profileMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer profileMock.CleanupMockState()
	config := loadUnitTestTerraform("resource_33_android_owner_wifi.tf")
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: config},
			{PreConfig: func() { profileMock.SetRemoteProperty("preSharedKey", nil) }, Config: config, PlanOnly: true, ExpectNonEmptyPlan: true},
			{Config: config},
			{Config: config, PlanOnly: true},
		},
	})
}

func TestUnitResourceDeviceConfigurationTemplatesJson_58_ProfileTypeChange(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, profileMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer profileMock.CleanupMockState()
	config := loadUnitTestTerraform("resource_03_ios_custom.tf")
	updated := strings.ReplaceAll(config, "#microsoft.graph.iosCustomConfiguration", "#microsoft.graph.macOSCustomConfiguration")
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: config},
			{Config: updated, ExpectError: regexp.MustCompile("Profile type cannot be changed")},
			{Config: config, PlanOnly: true},
		},
	})
}

func TestUnitResourceDeviceConfigurationTemplatesJson_59_InvalidAssignments(t *testing.T) {
	for _, tc := range []struct{ name, message string }{
		{"missing_group", "group targets require group_id"},
		{"unexpected_group", "group_id only applies to group targets"},
		{"missing_filter", "Attribute Required"},
		{"exclusion_filter", `exclusion groups do not support assignment\s+filters`},
		{"unexpected_filter", "filter_id requires filter_type"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			mocks.SetupUnitTestEnvironment(t)
			_, profileMock := setupMockEnvironment()
			defer httpmock.DeactivateAndReset()
			defer profileMock.CleanupMockState()
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
				Steps:                    []resource.TestStep{{Config: loadUnitTestTerraform("resource_59_" + tc.name + ".tf"), ExpectError: regexp.MustCompile(tc.message)}},
			})
		})
	}
}

func TestUnitResourceDeviceConfigurationTemplatesJson_60_PaginatedRelationships(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, profileMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer profileMock.CleanupMockState()
	var endpoint, firstID, secondID string
	config := loadUnitTestTerraform("resource_55_collection_multiple.tf")
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: config, Check: func(state *terraform.State) error {
				objects := state.RootModule().Resources
				firstID = objects[resourceType+".ios_root"].Primary.ID
				secondID = objects[resourceType+".second_root"].Primary.ID
				endpoint = "https://graph.microsoft.com/beta/deviceManagement/deviceConfigurations/" + objects[resourceType+".test"].Primary.ID + "/microsoft.graph.iosEnterpriseWiFiConfiguration/rootCertificatesForServerValidation"
				return nil
			}},
			{PreConfig: func() {
				first, err := helpers.ParseJSONFile("tests/responses/validate_get/get_device_configuration_relationships_page_1.json")
				if err != nil {
					t.Fatal(err)
				}
				second, err := helpers.ParseJSONFile("tests/responses/validate_get/get_device_configuration_relationships_page_2.json")
				if err != nil {
					t.Fatal(err)
				}
				first = strings.NewReplacer("FIRST_ID", firstID, "NEXT_LINK", endpoint+"?$skiptoken=next").Replace(first)
				second = strings.ReplaceAll(second, "SECOND_ID", secondID)
				httpmock.RegisterResponder("GET", endpoint, httpmock.NewStringResponder(200, first))
				httpmock.RegisterResponder("GET", endpoint+"?$skiptoken=next", httpmock.NewStringResponder(200, second))
			}, Config: config, PlanOnly: true},
		},
	})
}
