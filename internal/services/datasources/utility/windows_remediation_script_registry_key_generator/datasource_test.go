package utilityWindowsRemediationScriptRegistryKeyGenerator_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestWindowsRemediationScriptRegistryKeyGeneratorDataSource_CurrentUserDword(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigCurrentUserDword(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.microsoft365_utility_windows_remediation_script_registry_key_generator.test", "context", "current_user"),
					resource.TestCheckResourceAttr("data.microsoft365_utility_windows_remediation_script_registry_key_generator.test", "registry_key_path", "Software\\Policies\\Microsoft\\WindowsStore\\"),
					resource.TestCheckResourceAttr("data.microsoft365_utility_windows_remediation_script_registry_key_generator.test", "value_name", "RequirePrivateStoreOnly"),
					resource.TestCheckResourceAttr("data.microsoft365_utility_windows_remediation_script_registry_key_generator.test", "value_type", "REG_DWORD"),
					resource.TestCheckResourceAttr("data.microsoft365_utility_windows_remediation_script_registry_key_generator.test", "value_data", "1"),
					resource.TestMatchResourceAttr("data.microsoft365_utility_windows_remediation_script_registry_key_generator.test", "detection_script", regexp.MustCompile("Get-CimInstance win32_computersystem")),
					resource.TestMatchResourceAttr("data.microsoft365_utility_windows_remediation_script_registry_key_generator.test", "detection_script", regexp.MustCompile("RequirePrivateStoreOnly")),
					resource.TestMatchResourceAttr("data.microsoft365_utility_windows_remediation_script_registry_key_generator.test", "remediation_script", regexp.MustCompile("New-ItemProperty")),
					resource.TestMatchResourceAttr("data.microsoft365_utility_windows_remediation_script_registry_key_generator.test", "remediation_script", regexp.MustCompile("PropertyType DWord")),
					resource.TestCheckResourceAttrSet("data.microsoft365_utility_windows_remediation_script_registry_key_generator.test", "id"),
				),
			},
		},
	})
}

func TestWindowsRemediationScriptRegistryKeyGeneratorDataSource_AllUsersString(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigAllUsersString(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.microsoft365_utility_windows_remediation_script_registry_key_generator.test", "context", "all_users"),
					resource.TestCheckResourceAttr("data.microsoft365_utility_windows_remediation_script_registry_key_generator.test", "value_type", "REG_SZ"),
					resource.TestMatchResourceAttr("data.microsoft365_utility_windows_remediation_script_registry_key_generator.test", "detection_script", regexp.MustCompile(`foreach \(\$User in \$Users\)`)),
					resource.TestMatchResourceAttr("data.microsoft365_utility_windows_remediation_script_registry_key_generator.test", "detection_script", regexp.MustCompile("S-1-5-18")),
					resource.TestMatchResourceAttr("data.microsoft365_utility_windows_remediation_script_registry_key_generator.test", "remediation_script", regexp.MustCompile("PropertyType String")),
					resource.TestMatchResourceAttr("data.microsoft365_utility_windows_remediation_script_registry_key_generator.test", "remediation_script", regexp.MustCompile("'Enabled'")),
				),
			},
		},
	})
}

func TestWindowsRemediationScriptRegistryKeyGeneratorDataSource_MultiString(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMultiString(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.microsoft365_utility_windows_remediation_script_registry_key_generator.test", "value_type", "REG_MULTI_SZ"),
					resource.TestMatchResourceAttr("data.microsoft365_utility_windows_remediation_script_registry_key_generator.test", "remediation_script", regexp.MustCompile("PropertyType MultiString")),
					resource.TestMatchResourceAttr("data.microsoft365_utility_windows_remediation_script_registry_key_generator.test", "remediation_script", regexp.MustCompile(`@\(`)),
				),
			},
		},
	})
}

func TestWindowsRemediationScriptRegistryKeyGeneratorDataSource_Binary(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigBinary(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.microsoft365_utility_windows_remediation_script_registry_key_generator.test", "value_type", "REG_BINARY"),
					resource.TestMatchResourceAttr("data.microsoft365_utility_windows_remediation_script_registry_key_generator.test", "remediation_script", regexp.MustCompile("PropertyType Binary")),
					resource.TestMatchResourceAttr("data.microsoft365_utility_windows_remediation_script_registry_key_generator.test", "remediation_script", regexp.MustCompile("0x01")),
				),
			},
		},
	})
}

func TestWindowsRemediationScriptRegistryKeyGeneratorDataSource_InvalidContext(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testConfigInvalidContext(),
				ExpectError: regexp.MustCompile("Invalid Attribute Value Match"),
			},
		},
	})
}

func TestWindowsRemediationScriptRegistryKeyGeneratorDataSource_InvalidValueType(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testConfigInvalidValueType(),
				ExpectError: regexp.MustCompile("Invalid Attribute Value Match"),
			},
		},
	})
}

func TestWindowsRemediationScriptRegistryKeyGeneratorDataSource_InvalidDwordValue(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testConfigInvalidDwordValue(),
				ExpectError: regexp.MustCompile("invalid REG_DWORD value"),
			},
		},
	})
}

// Configuration functions
func testConfigCurrentUserDword() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/01_current_user_dword.tf")
	if err != nil {
		panic("failed to load current_user_dword config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigAllUsersString() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/02_all_users_string.tf")
	if err != nil {
		panic("failed to load all_users_string config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigMultiString() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/03_multistring.tf")
	if err != nil {
		panic("failed to load multistring config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigBinary() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/04_binary.tf")
	if err != nil {
		panic("failed to load binary config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigInvalidContext() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/05_invalid_context.tf")
	if err != nil {
		panic("failed to load invalid_context config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigInvalidValueType() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/06_invalid_value_type.tf")
	if err != nil {
		panic("failed to load invalid_value_type config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigInvalidDwordValue() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/07_invalid_dword_value.tf")
	if err != nil {
		panic("failed to load invalid_dword_value config: " + err.Error())
	}
	return unitTestConfig
}
