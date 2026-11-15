package graphBetaMobileThreatDefenseConnector_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	mobileThreatDefenseConnectorMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/mobile_threat_defense_connector/mocks"
)

const (
	unitPartnerId     = "8d0ed095-8191-4bd3-8a41-953b22d51ff7"
	unitDefenderId    = "fc780465-2017-40d4-a0c5-307022471b92"
	unitResourceLabel = ".test"
)

func setupMockEnvironment() (*mocks.Mocks, *mobileThreatDefenseConnectorMocks.MobileThreatDefenseConnectorMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	connectorMock := &mobileThreatDefenseConnectorMocks.MobileThreatDefenseConnectorMock{}
	connectorMock.RegisterMocks()
	return mockClient, connectorMock
}

func setupErrorMockEnvironment() (*mocks.Mocks, *mobileThreatDefenseConnectorMocks.MobileThreatDefenseConnectorMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	connectorMock := &mobileThreatDefenseConnectorMocks.MobileThreatDefenseConnectorMock{}
	connectorMock.RegisterErrorMocks()
	return mockClient, connectorMock
}

func loadUnitTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/unit/" + filename)
	if err != nil {
		panic("failed to load unit test config " + filename + ": " + err.Error())
	}
	return config
}

// checkMinimalConnector verifies the defaults applied when only the partner ID is configured.
func checkMinimalConnector(address string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		check.That(address).Key("id").HasValue(unitPartnerId),
		check.That(address).Key("android_enabled").HasValue("false"),
		check.That(address).Key("ios_enabled").HasValue("false"),
		check.That(address).Key("mac_enabled").HasValue("false"),
		check.That(address).Key("windows_enabled").HasValue("false"),
		check.That(address).Key("android_mobile_application_management_enabled").HasValue("false"),
		check.That(address).Key("ios_mobile_application_management_enabled").HasValue("false"),
		check.That(address).Key("windows_mobile_application_management_enabled").HasValue("false"),
		check.That(address).Key("partner_unsupported_os_version_blocked").HasValue("false"),
		check.That(address).Key("partner_unresponsiveness_threshold_in_days").HasValue("7"),
		check.That(address).Key("allow_partner_to_collect_ios_application_metadata").HasValue("false"),
		check.That(address).Key("allow_partner_to_collect_ios_personal_application_metadata").HasValue("false"),
		check.That(address).Key("allow_partner_to_collect_ios_certificate_metadata").HasValue("false"),
		check.That(address).Key("allow_partner_to_collect_ios_personal_certificate_metadata").HasValue("false"),
		check.That(address).Key("microsoft_defender_for_endpoint_attach_enabled").HasValue("false"),
		check.That(address).Key("grant_mobile_threat_defense_partner_role").HasValue("false"),
		check.That(address).Key("launch_mobile_threat_defense_partner_on_setup_enabled").HasValue("false"),
		check.That(address).Key("android_device_blocked_on_missing_partner_data").HasValue("false"),
		check.That(address).Key("ios_device_blocked_on_missing_partner_data").HasValue("false"),
		check.That(address).Key("mac_device_blocked_on_missing_partner_data").HasValue("false"),
		check.That(address).Key("windows_device_blocked_on_missing_partner_data").HasValue("false"),
		check.That(address).Key("partner_state").HasValue("notSetUp"),
		check.That(address).Key("last_heartbeat_date_time").DoesNotExist(),
		check.That(address).Key("partner_supports_android_mam").HasValue("false"),
		check.That(address).Key("partner_supports_ios_mam").HasValue("false"),
		check.That(address).Key("partner_supports_windows_mam").HasValue("false"),
		check.That(address).Key("partner_supports_ios_certificate_sync").HasValue("false"),
	)
}

// checkMaximalConnector verifies every writable toggle and the read-only fields Graph derives from them.
func checkMaximalConnector(address string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		check.That(address).Key("id").HasValue(unitPartnerId),
		check.That(address).Key("android_enabled").HasValue("true"),
		check.That(address).Key("ios_enabled").HasValue("true"),
		check.That(address).Key("mac_enabled").HasValue("true"),
		check.That(address).Key("windows_enabled").HasValue("true"),
		check.That(address).Key("android_mobile_application_management_enabled").HasValue("true"),
		check.That(address).Key("ios_mobile_application_management_enabled").HasValue("true"),
		check.That(address).Key("windows_mobile_application_management_enabled").HasValue("true"),
		check.That(address).Key("partner_unsupported_os_version_blocked").HasValue("true"),
		check.That(address).Key("partner_unresponsiveness_threshold_in_days").HasValue("14"),
		check.That(address).Key("allow_partner_to_collect_ios_application_metadata").HasValue("true"),
		check.That(address).Key("allow_partner_to_collect_ios_personal_application_metadata").HasValue("true"),
		check.That(address).Key("allow_partner_to_collect_ios_certificate_metadata").HasValue("true"),
		check.That(address).Key("allow_partner_to_collect_ios_personal_certificate_metadata").HasValue("true"),
		check.That(address).Key("microsoft_defender_for_endpoint_attach_enabled").HasValue("false"),
		check.That(address).Key("grant_mobile_threat_defense_partner_role").HasValue("true"),
		check.That(address).Key("launch_mobile_threat_defense_partner_on_setup_enabled").HasValue("false"),
		check.That(address).Key("android_device_blocked_on_missing_partner_data").HasValue("true"),
		check.That(address).Key("ios_device_blocked_on_missing_partner_data").HasValue("true"),
		check.That(address).Key("mac_device_blocked_on_missing_partner_data").HasValue("true"),
		check.That(address).Key("windows_device_blocked_on_missing_partner_data").HasValue("true"),
	)
}

func TestUnitResourceMobileThreatDefenseConnector_01_Minimal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, connectorMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer connectorMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_minimal.tf"),
				Check:  checkMinimalConnector(resourceType + unitResourceLabel),
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

func TestUnitResourceMobileThreatDefenseConnector_02_Maximal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, connectorMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer connectorMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_maximal.tf"),
				Check:  checkMaximalConnector(resourceType + unitResourceLabel),
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

func TestUnitResourceMobileThreatDefenseConnector_03_MinimalToMaximalToMinimal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, connectorMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer connectorMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_minimal.tf"),
				Check:  checkMinimalConnector(resourceType + unitResourceLabel),
			},
			{
				Config: loadUnitTestTerraform("resource_maximal.tf"),
				Check:  checkMaximalConnector(resourceType + unitResourceLabel),
			},
			{
				Config: loadUnitTestTerraform("resource_minimal.tf"),
				Check:  checkMinimalConnector(resourceType + unitResourceLabel),
			},
		},
	})
}

func TestUnitResourceMobileThreatDefenseConnector_04_MicrosoftDefender(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, connectorMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer connectorMock.CleanupMockState()

	address := resourceType + ".microsoft_defender"

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_microsoft_defender.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(address).Key("id").HasValue(unitDefenderId),
					check.That(address).Key("android_enabled").HasValue("true"),
					check.That(address).Key("ios_enabled").HasValue("true"),
					check.That(address).Key("mac_enabled").HasValue("false"),
					check.That(address).Key("windows_enabled").HasValue("true"),
					check.That(address).Key("partner_unsupported_os_version_blocked").HasValue("true"),
					check.That(address).Key("partner_unresponsiveness_threshold_in_days").HasValue("5"),
					check.That(address).Key("allow_partner_to_collect_ios_application_metadata").HasValue("true"),
					check.That(address).Key("microsoft_defender_for_endpoint_attach_enabled").HasValue("true"),
					check.That(address).Key("android_device_blocked_on_missing_partner_data").HasValue("true"),
					check.That(address).Key("mac_device_blocked_on_missing_partner_data").HasValue("false"),
				),
			},
			{
				ResourceName:      address,
				ImportState:       true,
				ImportStateId:     unitDefenderId,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"timeouts",
				},
			},
		},
	})
}

// A connector deleted outside Terraform reads as 404, so it is removed from state and planned for creation again.
func TestUnitResourceMobileThreatDefenseConnector_05_DeletedOutsideTerraform(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, connectorMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer connectorMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_minimal.tf"),
				Check:  checkMinimalConnector(resourceType + unitResourceLabel),
			},
			{
				PreConfig: func() {
					connectorMock.RemoveConnector(unitPartnerId)
				},
				Config:             loadUnitTestTerraform("resource_minimal.tf"),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestUnitResourceMobileThreatDefenseConnector_06_InvalidPartner(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, connectorMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer connectorMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadUnitTestTerraform("resource_error_invalid_partner.tf"),
				ExpectError: regexp.MustCompile(`(?s)Bad Request - 400.*not a valid\s+partner`),
			},
		},
	})
}

func TestUnitResourceMobileThreatDefenseConnector_07_Forbidden(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, connectorMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer connectorMock.CleanupMockState()

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

func TestUnitResourceMobileThreatDefenseConnector_08_CreateError(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, connectorMock := setupErrorMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer connectorMock.CleanupMockState()

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

func TestUnitResourceMobileThreatDefenseConnector_09_InvalidId(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, connectorMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer connectorMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `resource "microsoft365_graph_beta_device_management_mobile_threat_defense_connector" "test" {
  id = "not-a-guid"
}`,
				ExpectError: regexp.MustCompile(`must be a valid GUID`),
			},
		},
	})
}

func TestUnitResourceMobileThreatDefenseConnector_10_DefenderAttachOnOtherPartner(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, connectorMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer connectorMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadUnitTestTerraform("resource_error_defender_attach_other_partner.tf"),
				ExpectError: regexp.MustCompile(`Unsupported setting for partner`),
			},
		},
	})
}

// Graph's POST is an upsert, so create must refuse a partner that already has a connector instead of overwriting it.
func TestUnitResourceMobileThreatDefenseConnector_11_CreateWhenConnectorExists(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, connectorMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer connectorMock.CleanupMockState()

	if err := connectorMock.SeedConnector(unitPartnerId); err != nil {
		t.Fatalf("failed to seed connector: %v", err)
	}

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadUnitTestTerraform("resource_minimal.tf"),
				ExpectError: regexp.MustCompile(`(?s)already\s+exists.*terraform\s+import`),
			},
		},
	})
}

// A lookup error other than 404 is returned instead of being treated as "connector absent".
func TestUnitResourceMobileThreatDefenseConnector_12_CreateLookupError(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, connectorMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer connectorMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadUnitTestTerraform("resource_error_lookup_failed.tf"),
				ExpectError: regexp.MustCompile(`(?s)could not check whether mobile threat defense\s+connector`),
			},
		},
	})
}
