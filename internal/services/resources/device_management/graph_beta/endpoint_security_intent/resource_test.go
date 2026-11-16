package graphBetaEndpointSecurityIntent_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/jarcoal/httpmock"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	endpointSecurityIntentMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/endpoint_security_intent/mocks"
)

const (
	unitResourceLabel    = ".test"
	unitBitLockerTemplId = "d1174162-1dd2-4976-affc-6667049ab0ae"
	encryptDeviceId      = "deviceConfiguration--windows10EndpointProtectionConfiguration_bitLockerEncryptDevice"
	recoveryRotationId   = "deviceConfiguration--windows10EndpointProtectionConfiguration_bitLockerRecoveryPasswordRotation"
	standardUserId       = "deviceConfiguration--windows10EndpointProtectionConfiguration_bitLockerAllowStandardUserEncryption"
)

func setupMockEnvironment() (*mocks.Mocks, *endpointSecurityIntentMocks.EndpointSecurityIntentMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	intentMock := &endpointSecurityIntentMocks.EndpointSecurityIntentMock{}
	intentMock.RegisterMocks()
	return mockClient, intentMock
}

func setupErrorMockEnvironment() (*mocks.Mocks, *endpointSecurityIntentMocks.EndpointSecurityIntentMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	intentMock := &endpointSecurityIntentMocks.EndpointSecurityIntentMock{}
	intentMock.RegisterErrorMocks()
	return mockClient, intentMock
}

func loadUnitTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/unit/" + filename)
	if err != nil {
		panic("failed to load unit test config " + filename + ": " + err.Error())
	}
	return config
}

// importIdWithSettings builds an import ID that names the settings to manage, so import reads them back too.
func importIdWithSettings(address string, definitionIds ...string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[address]
		if !ok {
			return "", fmt.Errorf("resource %s not found in state", address)
		}
		importId := rs.Primary.ID + "/"
		for i, definitionId := range definitionIds {
			if i > 0 {
				importId += ","
			}
			importId += definitionId
		}
		return importId, nil
	}
}

// checkMinimalIntent verifies the defaults applied when only the template and display name are configured.
func checkMinimalIntent(address string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		check.That(address).Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-f-]{36}$`)),
		check.That(address).Key("template_id").HasValue(unitBitLockerTemplId),
		check.That(address).Key("display_name").HasValue("unit-test-endpoint-security-intent-minimal"),
		check.That(address).Key("description").HasValue(""),
		check.That(address).Key("role_scope_tag_ids.#").HasValue("1"),
		check.That(address).Key("role_scope_tag_ids.*").ContainsTypeSetElement("0"),
		check.That(address).Key("settings.#").DoesNotExist(),
		check.That(address).Key("assignments.#").DoesNotExist(),
		check.That(address).Key("is_assigned").HasValue("false"),
		check.That(address).Key("last_modified_date_time").IsSet(),
	)
}

// checkMaximalIntent verifies the configured settings and assignments.
func checkMaximalIntent(address string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		check.That(address).Key("template_id").HasValue(unitBitLockerTemplId),
		check.That(address).Key("display_name").HasValue("unit-test-endpoint-security-intent-maximal"),
		check.That(address).Key("description").HasValue("Unit test BitLocker intent"),
		check.That(address).Key("role_scope_tag_ids.#").HasValue("2"),
		check.That(address).Key("settings.#").HasValue("3"),
		resource.TestCheckTypeSetElemNestedAttrs(address, "settings.*", map[string]string{"definition_id": encryptDeviceId, "setting_type": "boolean", "value_json": "true"}),
		resource.TestCheckTypeSetElemNestedAttrs(address, "settings.*", map[string]string{"definition_id": recoveryRotationId, "setting_type": "string", "value_json": `"enabledForAzureAd"`}),
		resource.TestCheckTypeSetElemNestedAttrs(address, "settings.*", map[string]string{"definition_id": standardUserId, "setting_type": "boolean", "value_json": "true"}),
		check.That(address).Key("assignments.#").HasValue("3"),
		resource.TestCheckTypeSetElemNestedAttrs(address, "assignments.*", map[string]string{"type": "groupAssignmentTarget", "group_id": "11111111-1111-1111-1111-111111111111", "filter_type": "none"}),
		resource.TestCheckTypeSetElemNestedAttrs(address, "assignments.*", map[string]string{"type": "exclusionGroupAssignmentTarget", "group_id": "22222222-2222-2222-2222-222222222222"}),
		resource.TestCheckTypeSetElemNestedAttrs(address, "assignments.*", map[string]string{"type": "allDevicesAssignmentTarget", "filter_id": "33333333-3333-3333-3333-333333333333", "filter_type": "include"}),
		check.That(address).Key("is_assigned").HasValue("true"),
	)
}

func TestUnitResourceEndpointSecurityIntent_01_Minimal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, intentMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer intentMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_minimal.tf"),
				Check:  checkMinimalIntent(resourceType + unitResourceLabel),
			},
			{
				ResourceName:      resourceType + unitResourceLabel,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"timeouts",
				},
			},
		},
	})
}

func TestUnitResourceEndpointSecurityIntent_02_Maximal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, intentMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer intentMock.CleanupMockState()

	address := resourceType + unitResourceLabel

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_maximal.tf"),
				Check:  checkMaximalIntent(address),
			},
			{
				ResourceName:      address,
				ImportState:       true,
				ImportStateIdFunc: importIdWithSettings(address, encryptDeviceId, recoveryRotationId, standardUserId),
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"timeouts",
				},
			},
		},
	})
}

func TestUnitResourceEndpointSecurityIntent_03_MinimalToMaximalToMinimal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, intentMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer intentMock.CleanupMockState()

	address := resourceType + unitResourceLabel

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_minimal.tf"),
				Check:  checkMinimalIntent(address),
			},
			{
				Config: loadUnitTestTerraform("resource_maximal.tf"),
				Check:  checkMaximalIntent(address),
			},
			{
				Config: loadUnitTestTerraform("resource_minimal.tf"),
				Check:  checkMinimalIntent(address),
			},
		},
	})
}

// Changing a setting goes through updateSettings and keeps the same intent.
func TestUnitResourceEndpointSecurityIntent_04_SettingsUpdateInPlace(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, intentMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer intentMock.CleanupMockState()

	address := resourceType + unitResourceLabel
	var intentId string

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					checkMaximalIntent(address),
					resource.TestCheckResourceAttrWith(address, "id", func(value string) error {
						intentId = value
						return nil
					}),
				),
			},
			{
				Config: loadUnitTestTerraform("resource_settings_updated.tf"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPtr(address, "id", &intentId),
					resource.TestCheckTypeSetElemNestedAttrs(address, "settings.*", map[string]string{"definition_id": encryptDeviceId, "value_json": "false"}),
					resource.TestCheckTypeSetElemNestedAttrs(address, "settings.*", map[string]string{"definition_id": recoveryRotationId, "value_json": `"enabledForAzureAdAndHybrid"`}),
					func(_ *terraform.State) error {
						if value := intentMock.SettingValue(intentId, recoveryRotationId); value != `"enabledForAzureAdAndHybrid"` {
							return fmt.Errorf("expected updateSettings to write the new rotation value, mock holds %s", value)
						}
						return nil
					},
				),
			},
		},
	})
}

// A value that differs from Graph only in JSON formatting keeps its configured form, so there is no drift.
func TestUnitResourceEndpointSecurityIntent_05_ValueJsonFormatting(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, intentMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer intentMock.CleanupMockState()

	address := resourceType + unitResourceLabel

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_value_json_formatting.tf"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckTypeSetElemNestedAttrs(address, "settings.*", map[string]string{"definition_id": encryptDeviceId, "value_json": " true "}),
					resource.TestCheckTypeSetElemNestedAttrs(address, "settings.*", map[string]string{"definition_id": recoveryRotationId, "value_json": "\"enabledForAzureAd\"\n"}),
				),
			},
			{
				Config:   loadUnitTestTerraform("resource_value_json_formatting.tf"),
				PlanOnly: true,
			},
		},
	})
}

// An intent deleted outside Terraform reads as 404, so it is removed from state and planned for creation again.
func TestUnitResourceEndpointSecurityIntent_06_DeletedOutsideTerraform(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, intentMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer intentMock.CleanupMockState()

	address := resourceType + unitResourceLabel
	var intentId string

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_minimal.tf"),
				Check: resource.TestCheckResourceAttrWith(address, "id", func(value string) error {
					intentId = value
					return nil
				}),
			},
			{
				PreConfig: func() {
					intentMock.RemoveIntent(intentId)
				},
				Config:             loadUnitTestTerraform("resource_minimal.tf"),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestUnitResourceEndpointSecurityIntent_07_Forbidden(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, intentMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer intentMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadUnitTestTerraform("resource_error_forbidden.tf"),
				ExpectError: regexp.MustCompile(`Forbidden - 403`),
			},
		},
	})
}

func TestUnitResourceEndpointSecurityIntent_08_UnknownTemplate(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, intentMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer intentMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadUnitTestTerraform("resource_error_unknown_template.tf"),
				ExpectError: regexp.MustCompile(`(?s)Bad Request - 400.*Template not\s+found`),
			},
		},
	})
}

func TestUnitResourceEndpointSecurityIntent_09_UnknownSetting(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, intentMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer intentMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadUnitTestTerraform("resource_error_unknown_setting.tf"),
				ExpectError: regexp.MustCompile(`(?s)Bad Request - 400.*not\s+part\s+of\s+the\s+template`),
			},
		},
	})
}

func TestUnitResourceEndpointSecurityIntent_10_CreateError(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, intentMock := setupErrorMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer intentMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadUnitTestTerraform("resource_minimal.tf"),
				ExpectError: regexp.MustCompile(`Bad Request - 400`),
			},
		},
	})
}

func TestUnitResourceEndpointSecurityIntent_11_DuplicateSetting(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, intentMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer intentMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadUnitTestTerraform("resource_error_duplicate_setting.tf"),
				ExpectError: regexp.MustCompile(`Duplicate setting`),
			},
		},
	})
}

func TestUnitResourceEndpointSecurityIntent_12_InvalidValueJson(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, intentMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer intentMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadUnitTestTerraform("resource_error_invalid_json.tf"),
				ExpectError: regexp.MustCompile(`Invalid setting value`),
			},
		},
	})
}

func TestUnitResourceEndpointSecurityIntent_13_MissingGroupId(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, intentMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer intentMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadUnitTestTerraform("resource_error_missing_group_id.tf"),
				ExpectError: regexp.MustCompile(`Missing group ID`),
			},
		},
	})
}

func TestUnitResourceEndpointSecurityIntent_14_InvalidTemplateId(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, intentMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer intentMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `resource "microsoft365_graph_beta_device_management_endpoint_security_intent" "test" {
  template_id  = "not-a-guid"
  display_name = "unit-test-endpoint-security-intent-invalid-template-id"
}`,
				ExpectError: regexp.MustCompile(`must be a valid GUID`),
			},
		},
	})
}

func TestUnitResourceEndpointSecurityIntent_15_InvalidImportId(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, intentMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer intentMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:        loadUnitTestTerraform("resource_minimal.tf"),
				ResourceName:  resourceType + unitResourceLabel,
				ImportState:   true,
				ImportStateId: "not-a-guid",
				ExpectError:   regexp.MustCompile(`Invalid import ID`),
			},
		},
	})
}
