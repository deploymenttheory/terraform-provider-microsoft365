package graphBetaCrossTenantAccessPartnerGroupSyncSettings_test

import (
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	partnerSyncMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_beta/cross_tenant_access_partner_group_sync_settings/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

const resourceType = "microsoft365_graph_beta_identity_and_access_cross_tenant_access_partner_group_sync_settings"

func setupMockEnvironment() (*mocks.Mocks, *partnerSyncMocks.CrossTenantAccessPartnerGroupSyncSettingsMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	partnerSyncMock := &partnerSyncMocks.CrossTenantAccessPartnerGroupSyncSettingsMock{}
	partnerSyncMock.RegisterMocks()
	return mockClient, partnerSyncMock
}

func loadUnitConfig(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/unit/" + filename)
	if err != nil {
		panic("failed to load unit config " + filename + ": " + err.Error())
	}
	return config
}

// TestUnitResourceCrossTenantAccessPartnerGroupSyncSettings_01_Minimal verifies creation of partner group sync settings
// with minimal configuration (tenant_id and group_sync_inbound only).
func TestUnitResourceCrossTenantAccessPartnerGroupSyncSettings_01_Minimal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, partnerSyncMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer partnerSyncMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitConfig("resource_01_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("id").HasValue("12345678-1234-1234-1234-123456789012"),
					check.That(resourceType+".test").Key("tenant_id").HasValue("12345678-1234-1234-1234-123456789012"),
					check.That(resourceType+".test").Key("group_sync_inbound.is_sync_allowed").HasValue("true"),
				),
			},
			{
				ResourceName:      resourceType + ".test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"timeouts",
				},
			},
		},
	})
}

// TestUnitResourceCrossTenantAccessPartnerGroupSyncSettings_02_WithDisplayName verifies creation with display_name.
func TestUnitResourceCrossTenantAccessPartnerGroupSyncSettings_02_WithDisplayName(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, partnerSyncMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer partnerSyncMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitConfig("resource_02_with_display_name.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("id").HasValue("12345678-1234-1234-1234-123456789012"),
					check.That(resourceType+".test").Key("tenant_id").HasValue("12345678-1234-1234-1234-123456789012"),
					check.That(resourceType+".test").Key("display_name").HasValue("Partner Sync Configuration"),
					check.That(resourceType+".test").Key("group_sync_inbound.is_sync_allowed").HasValue("true"),
				),
			},
			{
				ResourceName:      resourceType + ".test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"timeouts",
				},
			},
		},
	})
}

// TestUnitResourceCrossTenantAccessPartnerGroupSyncSettings_03_SyncDisabled verifies sync can be disabled.
func TestUnitResourceCrossTenantAccessPartnerGroupSyncSettings_03_SyncDisabled(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, partnerSyncMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer partnerSyncMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitConfig("resource_03_sync_disabled.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("id").HasValue("12345678-1234-1234-1234-123456789012"),
					check.That(resourceType+".test").Key("tenant_id").HasValue("12345678-1234-1234-1234-123456789012"),
					check.That(resourceType+".test").Key("display_name").HasValue("Partner Sync Disabled"),
					check.That(resourceType+".test").Key("group_sync_inbound.is_sync_allowed").HasValue("false"),
				),
			},
		},
	})
}

// TestUnitResourceCrossTenantAccessPartnerGroupSyncSettings_04_Update verifies updating settings.
func TestUnitResourceCrossTenantAccessPartnerGroupSyncSettings_04_Update(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, partnerSyncMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer partnerSyncMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitConfig("resource_01_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("group_sync_inbound.is_sync_allowed").HasValue("true"),
				),
			},
			{
				Config: loadUnitConfig("resource_03_sync_disabled.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("display_name").HasValue("Partner Sync Disabled"),
					check.That(resourceType+".test").Key("group_sync_inbound.is_sync_allowed").HasValue("false"),
				),
			},
		},
	})
}
