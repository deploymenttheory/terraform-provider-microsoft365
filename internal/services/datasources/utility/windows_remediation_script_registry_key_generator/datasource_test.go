package utilityWindowsRemediationScriptRegistryKeyGenerator_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	utilityWindowsRemediationScriptRegistryKeyGenerator "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/datasources/utility/windows_remediation_script_registry_key_generator"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var (
	// DataSource type name from the datasource package
	dataSourceType = utilityWindowsRemediationScriptRegistryKeyGenerator.DataSourceName
)

func TestWindowsRemediationScriptRegistryKeyGeneratorDataSource_CurrentUserDword(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigCurrentUserDword(),
				Check: resource.ComposeTestCheckFunc(
					check.That("data."+dataSourceType+".test").Key("context").HasValue("current_user"),
					check.That("data."+dataSourceType+".test").Key("registry_key_path").HasValue("Software\\Policies\\Microsoft\\WindowsStore\\"),
					check.That("data."+dataSourceType+".test").Key("value_name").HasValue("RequirePrivateStoreOnly"),
					check.That("data."+dataSourceType+".test").Key("value_type").HasValue("REG_DWORD"),
					check.That("data."+dataSourceType+".test").Key("value_data").HasValue("1"),
					check.That("data."+dataSourceType+".test").Key("detection_script").MatchesRegex(regexp.MustCompile("Get-CimInstance win32_computersystem")),
					check.That("data."+dataSourceType+".test").Key("detection_script").MatchesRegex(regexp.MustCompile("RequirePrivateStoreOnly")),
					check.That("data."+dataSourceType+".test").Key("remediation_script").MatchesRegex(regexp.MustCompile("New-ItemProperty")),
					check.That("data."+dataSourceType+".test").Key("remediation_script").MatchesRegex(regexp.MustCompile("PropertyType DWord")),
					check.That("data."+dataSourceType+".test").Key("id").IsSet(),
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
					check.That("data."+dataSourceType+".test").Key("context").HasValue("all_users"),
					check.That("data."+dataSourceType+".test").Key("value_type").HasValue("REG_SZ"),
					check.That("data."+dataSourceType+".test").Key("detection_script").MatchesRegex(regexp.MustCompile(`foreach \(\$User in \$Users\)`)),
					check.That("data."+dataSourceType+".test").Key("detection_script").MatchesRegex(regexp.MustCompile("S-1-5-18")),
					check.That("data."+dataSourceType+".test").Key("remediation_script").MatchesRegex(regexp.MustCompile("PropertyType String")),
					check.That("data."+dataSourceType+".test").Key("remediation_script").MatchesRegex(regexp.MustCompile("'Enabled'")),
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
					check.That("data."+dataSourceType+".test").Key("value_type").HasValue("REG_MULTI_SZ"),
					check.That("data."+dataSourceType+".test").Key("remediation_script").MatchesRegex(regexp.MustCompile("PropertyType MultiString")),
					check.That("data."+dataSourceType+".test").Key("remediation_script").MatchesRegex(regexp.MustCompile(`@\(`)),
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
					check.That("data."+dataSourceType+".test").Key("value_type").HasValue("REG_BINARY"),
					check.That("data."+dataSourceType+".test").Key("remediation_script").MatchesRegex(regexp.MustCompile("PropertyType Binary")),
					check.That("data."+dataSourceType+".test").Key("remediation_script").MatchesRegex(regexp.MustCompile("0x01")),
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
