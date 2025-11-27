package utilityWindowsMSIAppMetadata_test

import (
	"os"
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	utilityWindowsMSIAppMetadata "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/datasources/utility/windows_msi_app_metadata"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var (
	// DataSource type name from the datasource package
	dataSourceType = utilityWindowsMSIAppMetadata.DataSourceName
)

func TestWindowsMSIAppMetadataDataSource_LocalFile(t *testing.T) {
	// Check if test MSI file exists
	if _, err := os.Stat("testdata/sample.msi"); os.IsNotExist(err) {
		t.Skip("Skipping test: testdata/sample.msi not found. This test requires a valid MSI file for testing.")
	}

	mocks.SetupUnitTestEnvironment(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigLocalFile(),
				Check: resource.ComposeTestCheckFunc(
					check.That("data."+dataSourceType+".test").Key("id").IsSet(),
					check.That("data."+dataSourceType+".test").Key("installer_file_path_source").HasValue("testdata/sample.msi"),
					check.That("data."+dataSourceType+".test").Key("metadata.product_code").IsSet(),
					check.That("data."+dataSourceType+".test").Key("metadata.product_version").IsSet(),
					check.That("data."+dataSourceType+".test").Key("metadata.product_name").IsSet(),
					check.That("data."+dataSourceType+".test").Key("metadata.publisher").IsSet(),
					check.That("data."+dataSourceType+".test").Key("metadata.sha256_checksum").IsSet(),
					check.That("data."+dataSourceType+".test").Key("metadata.size_mb").IsSet(),
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
				ExpectError: regexp.MustCompile("Invalid Attribute Combination"),
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
				ExpectError: regexp.MustCompile("Error Reading MSI File|Error Extracting MSI Metadata"),
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
