package utilityWindowsRemediationScriptRegistryKeyGenerator_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDatasourceWindowsRemediationScriptRegistryKeyGenerator_01_CurrentUserDword(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccConfigCurrentUserDword(),
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

func TestAccDatasourceWindowsRemediationScriptRegistryKeyGenerator_02_AllUsersString(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccConfigAllUsersString(),
				Check: resource.ComposeTestCheckFunc(
					check.That("data."+dataSourceType+".test").Key("context").HasValue("all_users"),
					check.That("data."+dataSourceType+".test").Key("detection_script").MatchesRegex(regexp.MustCompile(`foreach \(\$User in \$Users\)`)),
					check.That("data."+dataSourceType+".test").Key("detection_script").MatchesRegex(regexp.MustCompile("S-1-5-18")),
					check.That("data."+dataSourceType+".test").Key("remediation_script").MatchesRegex(regexp.MustCompile("PropertyType String")),
				),
			},
		},
	})
}

func TestAccDatasourceWindowsRemediationScriptRegistryKeyGenerator_03_QWord(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccConfigQWord(),
				Check: resource.ComposeTestCheckFunc(
					check.That("data."+dataSourceType+".test").Key("value_type").HasValue("REG_QWORD"),
					check.That("data."+dataSourceType+".test").Key("value_data").HasValue("9223372036854775807"),
					check.That("data."+dataSourceType+".test").Key("remediation_script").MatchesRegex(regexp.MustCompile("PropertyType QWord")),
					check.That("data."+dataSourceType+".test").Key("remediation_script").MatchesRegex(regexp.MustCompile("9223372036854775807")),
				),
			},
		},
	})
}

func TestAccDatasourceWindowsRemediationScriptRegistryKeyGenerator_04_ExpandString(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccConfigExpandString(),
				Check: resource.ComposeTestCheckFunc(
					check.That("data."+dataSourceType+".test").Key("value_type").HasValue("REG_EXPAND_SZ"),
					check.That("data."+dataSourceType+".test").Key("remediation_script").MatchesRegex(regexp.MustCompile("PropertyType ExpandString")),
					check.That("data."+dataSourceType+".test").Key("remediation_script").MatchesRegex(regexp.MustCompile("%ProgramFiles%")),
				),
			},
		},
	})
}

// Acceptance test configuration functions
func testAccConfigCurrentUserDword() string {
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/01_current_user_dword.tf")
	if err != nil {
		panic(fmt.Sprintf("failed to load acceptance test config: %s", err.Error()))
	}
	return acceptance.ConfiguredM365ProviderBlock(accTestConfig)
}

func testAccConfigAllUsersString() string {
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/02_all_users_string.tf")
	if err != nil {
		panic(fmt.Sprintf("failed to load acceptance test config: %s", err.Error()))
	}
	return acceptance.ConfiguredM365ProviderBlock(accTestConfig)
}

func testAccConfigQWord() string {
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/03_qword.tf")
	if err != nil {
		panic(fmt.Sprintf("failed to load acceptance test config: %s", err.Error()))
	}
	return acceptance.ConfiguredM365ProviderBlock(accTestConfig)
}

func testAccConfigExpandString() string {
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/04_expand_string.tf")
	if err != nil {
		panic(fmt.Sprintf("failed to load acceptance test config: %s", err.Error()))
	}
	return acceptance.ConfiguredM365ProviderBlock(accTestConfig)
}
