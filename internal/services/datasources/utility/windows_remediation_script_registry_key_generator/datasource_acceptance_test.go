package utilityWindowsRemediationScriptRegistryKeyGenerator_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccWindowsRemediationScriptRegistryKeyGeneratorDataSource_CurrentUserDword(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccConfigCurrentUserDword(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.microsoft365_windows_remediation_script_registry_key_generator.test", "context", "current_user"),
					resource.TestCheckResourceAttr("data.microsoft365_windows_remediation_script_registry_key_generator.test", "registry_key_path", "Software\\Policies\\Microsoft\\WindowsStore\\"),
					resource.TestCheckResourceAttr("data.microsoft365_windows_remediation_script_registry_key_generator.test", "value_name", "RequirePrivateStoreOnly"),
					resource.TestCheckResourceAttr("data.microsoft365_windows_remediation_script_registry_key_generator.test", "value_type", "REG_DWORD"),
					resource.TestCheckResourceAttr("data.microsoft365_windows_remediation_script_registry_key_generator.test", "value_data", "1"),
					resource.TestMatchResourceAttr("data.microsoft365_windows_remediation_script_registry_key_generator.test", "detection_script", regexp.MustCompile("Get-CimInstance win32_computersystem")),
					resource.TestMatchResourceAttr("data.microsoft365_windows_remediation_script_registry_key_generator.test", "detection_script", regexp.MustCompile("RequirePrivateStoreOnly")),
					resource.TestMatchResourceAttr("data.microsoft365_windows_remediation_script_registry_key_generator.test", "remediation_script", regexp.MustCompile("New-ItemProperty")),
					resource.TestMatchResourceAttr("data.microsoft365_windows_remediation_script_registry_key_generator.test", "remediation_script", regexp.MustCompile("PropertyType DWord")),
					resource.TestCheckResourceAttrSet("data.microsoft365_windows_remediation_script_registry_key_generator.test", "id"),
				),
			},
		},
	})
}

func TestAccWindowsRemediationScriptRegistryKeyGeneratorDataSource_AllUsersString(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccConfigAllUsersString(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.microsoft365_windows_remediation_script_registry_key_generator.test", "context", "all_users"),
					resource.TestMatchResourceAttr("data.microsoft365_windows_remediation_script_registry_key_generator.test", "detection_script", regexp.MustCompile(`foreach \(\$User in \$Users\)`)),
					resource.TestMatchResourceAttr("data.microsoft365_windows_remediation_script_registry_key_generator.test", "detection_script", regexp.MustCompile("S-1-5-18")),
					resource.TestMatchResourceAttr("data.microsoft365_windows_remediation_script_registry_key_generator.test", "remediation_script", regexp.MustCompile("PropertyType String")),
				),
			},
		},
	})
}

func TestAccWindowsRemediationScriptRegistryKeyGeneratorDataSource_QWord(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccConfigQWord(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.microsoft365_windows_remediation_script_registry_key_generator.test", "value_type", "REG_QWORD"),
					resource.TestCheckResourceAttr("data.microsoft365_windows_remediation_script_registry_key_generator.test", "value_data", "9223372036854775807"),
					resource.TestMatchResourceAttr("data.microsoft365_windows_remediation_script_registry_key_generator.test", "remediation_script", regexp.MustCompile("PropertyType QWord")),
					resource.TestMatchResourceAttr("data.microsoft365_windows_remediation_script_registry_key_generator.test", "remediation_script", regexp.MustCompile("9223372036854775807")),
				),
			},
		},
	})
}

func TestAccWindowsRemediationScriptRegistryKeyGeneratorDataSource_ExpandString(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccConfigExpandString(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.microsoft365_windows_remediation_script_registry_key_generator.test", "value_type", "REG_EXPAND_SZ"),
					resource.TestMatchResourceAttr("data.microsoft365_windows_remediation_script_registry_key_generator.test", "remediation_script", regexp.MustCompile("PropertyType ExpandString")),
					resource.TestMatchResourceAttr("data.microsoft365_windows_remediation_script_registry_key_generator.test", "remediation_script", regexp.MustCompile("%ProgramFiles%")),
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
