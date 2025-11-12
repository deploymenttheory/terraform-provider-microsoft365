package utilityWindowsMSIAppMetadata_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestWindowsMSIAppMetadataDataSource_LocalFile(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigLocalFile(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.microsoft365_utility_windows_msi_app_metadata.test", "id"),
					resource.TestCheckResourceAttr("data.microsoft365_utility_windows_msi_app_metadata.test", "installer_file_path_source", "testdata/sample.msi"),
					resource.TestCheckResourceAttrSet("data.microsoft365_utility_windows_msi_app_metadata.test", "metadata.product_code"),
					resource.TestCheckResourceAttrSet("data.microsoft365_utility_windows_msi_app_metadata.test", "metadata.product_version"),
					resource.TestCheckResourceAttrSet("data.microsoft365_utility_windows_msi_app_metadata.test", "metadata.product_name"),
					resource.TestCheckResourceAttrSet("data.microsoft365_utility_windows_msi_app_metadata.test", "metadata.publisher"),
					resource.TestCheckResourceAttrSet("data.microsoft365_utility_windows_msi_app_metadata.test", "metadata.sha256_checksum"),
					resource.TestCheckResourceAttrSet("data.microsoft365_utility_windows_msi_app_metadata.test", "metadata.size_mb"),
				),
			},
		},
	})
}

func TestWindowsMSIAppMetadataDataSource_MissingInput(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testConfigMissingInput(),
				ExpectError: regexp.MustCompile("Missing Input Parameter"),
			},
		},
	})
}

func TestWindowsMSIAppMetadataDataSource_BothInputsProvided(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testConfigBothInputs(),
				ExpectError: regexp.MustCompile("Multiple Input Parameters|Conflicting Attribute Values"),
			},
		},
	})
}

func TestWindowsMSIAppMetadataDataSource_InvalidFilePath(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testConfigInvalidFilePath(),
				ExpectError: regexp.MustCompile("Error Reading MSI File|no such file or directory"),
			},
		},
	})
}

func TestWindowsMSIAppMetadataDataSource_InvalidMSIFormat(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testConfigInvalidMSIFormat(),
				ExpectError: regexp.MustCompile("Error Extracting MSI Metadata"),
			},
		},
	})
}

// Configuration functions
func testConfigLocalFile() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/01_local_file.tf")
	if err != nil {
		panic("failed to load local_file config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigMissingInput() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/08_missing_input.tf")
	if err != nil {
		panic("failed to load missing_input config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigBothInputs() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/09_both_inputs.tf")
	if err != nil {
		panic("failed to load both_inputs config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigInvalidFilePath() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/10_invalid_file_path.tf")
	if err != nil {
		panic("failed to load invalid_file_path config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigInvalidMSIFormat() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/11_invalid_msi_format.tf")
	if err != nil {
		panic("failed to load invalid_msi_format config: " + err.Error())
	}
	return unitTestConfig
}
