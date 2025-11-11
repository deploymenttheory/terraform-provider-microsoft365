package entra_id_sid_converter_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestEntraIdSidConverterDataSource_SidToObjectId(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigSidToObjectId(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.microsoft365_utility_entra_id_sid_converter.test", "sid", "S-1-12-1-1943430372-1249052806-2496021943-3034400218"),
					resource.TestCheckResourceAttr("data.microsoft365_utility_entra_id_sid_converter.test", "object_id", "73d664e4-0886-4a73-b745-c694da45ddb4"),
					resource.TestCheckResourceAttrSet("data.microsoft365_utility_entra_id_sid_converter.test", "id"),
				),
			},
		},
	})
}

func TestEntraIdSidConverterDataSource_ObjectIdToSid(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigObjectIdToSid(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.microsoft365_utility_entra_id_sid_converter.test", "object_id", "73d664e4-0886-4a73-b745-c694da45ddb4"),
					resource.TestCheckResourceAttr("data.microsoft365_utility_entra_id_sid_converter.test", "sid", "S-1-12-1-1943430372-1249052806-2496021943-3034400218"),
					resource.TestCheckResourceAttrSet("data.microsoft365_utility_entra_id_sid_converter.test", "id"),
				),
			},
		},
	})
}

func TestEntraIdSidConverterDataSource_InvalidRidExceedsMax(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testConfigInvalidRidExceedsMax(),
				ExpectError: regexp.MustCompile("RID Component Out of Range"),
			},
		},
	})
}

func TestEntraIdSidConverterDataSource_InvalidSidFormat(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testConfigInvalidSidFormat(),
				ExpectError: regexp.MustCompile("SID must be in the format"),
			},
		},
	})
}

func TestEntraIdSidConverterDataSource_InvalidGuidFormat(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testConfigInvalidGuidFormat(),
				ExpectError: regexp.MustCompile("Object ID must be a valid GUID"),
			},
		},
	})
}

func TestEntraIdSidConverterDataSource_BothParametersProvided(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testConfigBothProvided(),
				ExpectError: regexp.MustCompile("Invalid Attribute Combination"),
			},
		},
	})
}

func TestEntraIdSidConverterDataSource_NeitherParameterProvided(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testConfigNeitherProvided(),
				ExpectError: regexp.MustCompile("Invalid Attribute Combination"),
			},
		},
	})
}

// Configuration functions
func testConfigSidToObjectId() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/01_sid_to_objectid.tf")
	if err != nil {
		panic("failed to load sid_to_objectid config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigObjectIdToSid() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/02_objectid_to_sid.tf")
	if err != nil {
		panic("failed to load objectid_to_sid config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigInvalidRidExceedsMax() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/03_invalid_rid_exceeds_max.tf")
	if err != nil {
		panic("failed to load invalid_rid_exceeds_max config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigInvalidSidFormat() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/04_invalid_sid_format.tf")
	if err != nil {
		panic("failed to load invalid_sid_format config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigInvalidGuidFormat() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/05_invalid_guid_format.tf")
	if err != nil {
		panic("failed to load invalid_guid_format config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigBothProvided() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/06_both_provided.tf")
	if err != nil {
		panic("failed to load both_provided config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigNeitherProvided() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/07_neither_provided.tf")
	if err != nil {
		panic("failed to load neither_provided config: " + err.Error())
	}
	return unitTestConfig
}
