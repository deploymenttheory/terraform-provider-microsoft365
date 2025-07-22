package graphIOSMobileAppConfiguration_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/jarcoal/httpmock"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	localMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/datasources/device_and_app_management/graph_v1.0/ios_mobile_app_configuration/mocks/datasource"
)

func TestUnitIOSMobileAppConfigurationDataSource_Read(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()

	// Register data source specific mocks
	dsMock := &localMocks.IOSMobileAppConfigurationDataSourceMock{}
	dsMock.RegisterMocks()

	dataSourceType := "microsoft365_graph_v1_device_and_app_management_ios_mobile_app_configuration"
	dataSourceName := "test"
	dataSourceID := fmt.Sprintf("data.%s.%s", dataSourceType, dataSourceName)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testIOSMobileAppConfigurationDataSourceConfigByID("00000000-0000-0000-0000-000000000001"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceID, "id", "00000000-0000-0000-0000-000000000001"),
					resource.TestCheckResourceAttr(dataSourceID, "display_name", "iOS App Config Test"),
					resource.TestCheckResourceAttr(dataSourceID, "description", "Test iOS app configuration"),
					resource.TestCheckResourceAttr(dataSourceID, "targeted_mobile_apps.#", "2"),
					resource.TestCheckResourceAttr(dataSourceID, "targeted_mobile_apps.0", "com.example.app1"),
					resource.TestCheckResourceAttr(dataSourceID, "targeted_mobile_apps.1", "com.example.app2"),
					resource.TestCheckResourceAttr(dataSourceID, "settings.#", "2"),
					resource.TestCheckResourceAttr(dataSourceID, "settings.0.app_config_key", "serverUrl"),
					resource.TestCheckResourceAttr(dataSourceID, "settings.0.app_config_key_type", "stringType"),
					resource.TestCheckResourceAttr(dataSourceID, "settings.0.app_config_key_value", "https://api.example.com"),
					resource.TestCheckResourceAttr(dataSourceID, "settings.1.app_config_key", "syncInterval"),
					resource.TestCheckResourceAttr(dataSourceID, "settings.1.app_config_key_type", "integerType"),
					resource.TestCheckResourceAttr(dataSourceID, "settings.1.app_config_key_value", "300"),
					resource.TestCheckResourceAttr(dataSourceID, "assignments.#", "1"),
					resource.TestCheckResourceAttr(dataSourceID, "assignments.0.id", "00000000-0000-0000-0000-000000000002"),
					resource.TestCheckResourceAttr(dataSourceID, "assignments.0.target.odata_type", "#microsoft.graph.groupAssignmentTarget"),
					resource.TestCheckResourceAttr(dataSourceID, "assignments.0.target.group_id", "00000000-0000-0000-0000-000000000003"),
					resource.TestCheckResourceAttr(dataSourceID, "version", "1"),
					testCheckResourceDateTimeAttr(dataSourceID, "created_date_time"),
					testCheckResourceDateTimeAttr(dataSourceID, "last_modified_date_time"),
				),
			},
		},
	})
}

func TestUnitIOSMobileAppConfigurationDataSource_ReadByDisplayName(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()

	// Register data source specific mocks
	dsMock := &localMocks.IOSMobileAppConfigurationDataSourceMock{}
	dsMock.RegisterMocks()

	dataSourceType := "microsoft365_graph_v1_device_and_app_management_ios_mobile_app_configuration"
	dataSourceName := "test"
	dataSourceID := fmt.Sprintf("data.%s.%s", dataSourceType, dataSourceName)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testIOSMobileAppConfigurationDataSourceConfigByDisplayName("iOS App Config Test"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceID, "id", "00000000-0000-0000-0000-000000000001"),
					resource.TestCheckResourceAttr(dataSourceID, "display_name", "iOS App Config Test"),
					resource.TestCheckResourceAttr(dataSourceID, "description", "Test iOS app configuration"),
					resource.TestCheckResourceAttr(dataSourceID, "targeted_mobile_apps.#", "2"),
					resource.TestCheckResourceAttr(dataSourceID, "settings.#", "2"),
					resource.TestCheckResourceAttr(dataSourceID, "assignments.#", "1"),
				),
			},
		},
	})
}

func TestUnitIOSMobileAppConfigurationDataSource_ReadWithEncodedXml(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()

	// Register data source specific mocks
	dsMock := &localMocks.IOSMobileAppConfigurationDataSourceMock{}
	dsMock.RegisterMocks()

	dataSourceType := "microsoft365_graph_v1_device_and_app_management_ios_mobile_app_configuration"
	dataSourceName := "test"
	dataSourceID := fmt.Sprintf("data.%s.%s", dataSourceType, dataSourceName)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testIOSMobileAppConfigurationDataSourceConfigByID("00000000-0000-0000-0000-000000000004"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceID, "id", "00000000-0000-0000-0000-000000000004"),
					resource.TestCheckResourceAttr(dataSourceID, "display_name", "iOS Config with XML"),
					resource.TestCheckResourceAttrSet(dataSourceID, "encoded_setting_xml"),
					resource.TestCheckResourceAttr(dataSourceID, "settings.#", "0"),
				),
			},
		},
	})
}

func TestUnitIOSMobileAppConfigurationDataSource_ReadNotFound(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()

	// Register data source specific mocks
	dsMock := &localMocks.IOSMobileAppConfigurationDataSourceMock{}
	dsMock.RegisterMocks()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testIOSMobileAppConfigurationDataSourceConfigByDisplayName("Non-existent Config"),
				ExpectError: regexp.MustCompile("(?i)no iOS mobile app configuration found with display name"),
			},
		},
	})
}

func TestUnitIOSMobileAppConfigurationDataSource_ReadMultipleFound(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()

	// Register data source specific mocks
	dsMock := &localMocks.IOSMobileAppConfigurationDataSourceMock{}
	dsMock.RegisterMocks()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testIOSMobileAppConfigurationDataSourceConfigByDisplayName("Duplicate Config"),
				ExpectError: regexp.MustCompile("(?i)multiple iOS mobile app configurations found with display name"),
			},
		},
	})
}

func TestUnitIOSMobileAppConfigurationDataSource_ReadMissingRequiredFields(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()

	// Register data source specific mocks
	dsMock := &localMocks.IOSMobileAppConfigurationDataSourceMock{}
	dsMock.RegisterMocks()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testIOSMobileAppConfigurationDataSourceConfigEmpty(),
				ExpectError: regexp.MustCompile("Either 'id' or 'display_name' must be provided"),
			},
		},
	})
}

// Helper functions

func testIOSMobileAppConfigurationDataSourceConfigByID(id string) string {
	return fmt.Sprintf(`
provider "microsoft365" {
}

data "microsoft365_graph_v1_device_and_app_management_ios_mobile_app_configuration" "test" {
  id = "%s"
}
`, id)
}

func testIOSMobileAppConfigurationDataSourceConfigByDisplayName(displayName string) string {
	return fmt.Sprintf(`
provider "microsoft365" {
}

data "microsoft365_graph_v1_device_and_app_management_ios_mobile_app_configuration" "test" {
  display_name = "%s"
}
`, displayName)
}

func testIOSMobileAppConfigurationDataSourceConfigEmpty() string {
	return `
provider "microsoft365" {
}

data "microsoft365_graph_v1_device_and_app_management_ios_mobile_app_configuration" "test" {
}
`
}

func testCheckResourceDateTimeAttr(resourceName, attributeName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		value := s.RootModule().Resources[resourceName].Primary.Attributes[attributeName]
		if value == "" {
			return fmt.Errorf("expected %s to be set", attributeName)
		}
		return nil
	}
}