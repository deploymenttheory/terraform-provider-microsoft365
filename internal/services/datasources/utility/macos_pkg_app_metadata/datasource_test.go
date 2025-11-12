package utilityMacOSPKGAppMetadata_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestMacOSPKGAppMetadataDataSource_InvalidFileExtension(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testConfigInvalidFileExtension(),
				ExpectError: regexp.MustCompile("File path must point to a valid \\.pkg"),
			},
		},
	})
}

func TestMacOSPKGAppMetadataDataSource_InvalidURL(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testConfigInvalidURL(),
				ExpectError: regexp.MustCompile("Must be a valid URL"),
			},
		},
	})
}

func TestMacOSPKGAppMetadataDataSource_BothParametersProvided(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testConfigBothProvided(),
				ExpectError: regexp.MustCompile("Multiple Input Sources"),
			},
		},
	})
}

func TestMacOSPKGAppMetadataDataSource_NeitherParameterProvided(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testConfigNeitherProvided(),
				ExpectError: regexp.MustCompile("Missing Input Source"),
			},
		},
	})
}

func TestMacOSPKGAppMetadataDataSource_ValidFilePathFormat(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigValidFilePath(),
				// This will fail when trying to read the actual file, but validates the path format
				ExpectError: regexp.MustCompile("Error Extracting Metadata from File|no such file"),
			},
		},
	})
}

func TestMacOSPKGAppMetadataDataSource_ValidURLFormat(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigValidURL(),
				// This will fail when trying to download, but validates the URL format
				ExpectError: regexp.MustCompile("Error Extracting Metadata from URL|connection refused|no such host"),
			},
		},
	})
}

// Configuration functions
func testConfigInvalidFileExtension() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/01_invalid_file_extension.tf")
	if err != nil {
		panic("failed to load invalid_file_extension config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigInvalidURL() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/02_invalid_url.tf")
	if err != nil {
		panic("failed to load invalid_url config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigBothProvided() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/03_both_provided.tf")
	if err != nil {
		panic("failed to load both_provided config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigNeitherProvided() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/04_neither_provided.tf")
	if err != nil {
		panic("failed to load neither_provided config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigValidFilePath() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/05_valid_file_path.tf")
	if err != nil {
		panic("failed to load valid_file_path config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigValidURL() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/06_valid_url.tf")
	if err != nil {
		panic("failed to load valid_url config: " + err.Error())
	}
	return unitTestConfig
}
