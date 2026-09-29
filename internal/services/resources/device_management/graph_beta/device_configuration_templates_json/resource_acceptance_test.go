package graphBetaDeviceConfigurationTemplatesJson_test

import (
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/destroy"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/testlog"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func loadAcceptanceTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/acceptance/" + filename)
	if err != nil {
		panic("failed to load acceptance config " + filename + ": " + err.Error())
	}
	return config
}

func TestAccResourceDeviceConfigurationTemplatesJson_01_IosGeneral(t *testing.T) {
	config := loadAcceptanceTestTerraform("resource_01_ios_general.tf")
	updated := strings.ReplaceAll(config, "Device configuration template test", "Updated template description")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 30*time.Second),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {Source: "hashicorp/random", VersionConstraint: constants.ExternalProviderRandomVersion},
		},
		Steps: []resource.TestStep{
			{PreConfig: func() { testlog.StepAction(resourceType, "Creating ios_general profile") },
				Config: config, Check: check.That(resourceType + ".test").ExistsInGraph(testResource)},
			importStep(),
			{Config: updated, Check: check.That(resourceType + ".test").Key("description").HasValue("Updated template description")},
			importStep(),
			{Config: updated, PlanOnly: true},
		},
	})
}

func TestAccResourceDeviceConfigurationTemplatesJson_02_IosFeatures(t *testing.T) {
	config := loadAcceptanceTestTerraform("resource_02_ios_features.tf")
	updated := strings.ReplaceAll(config, "Device configuration template test", "Updated template description")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 30*time.Second),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {Source: "hashicorp/random", VersionConstraint: constants.ExternalProviderRandomVersion},
		},
		Steps: []resource.TestStep{
			{PreConfig: func() { testlog.StepAction(resourceType, "Creating ios_features profile") },
				Config: config, Check: check.That(resourceType + ".test").ExistsInGraph(testResource)},
			importStep(),
			{Config: updated, Check: check.That(resourceType + ".test").Key("description").HasValue("Updated template description")},
			importStep(),
			{Config: updated, PlanOnly: true},
		},
	})
}

func TestAccResourceDeviceConfigurationTemplatesJson_03_IosCustom(t *testing.T) {
	config := loadAcceptanceTestTerraform("resource_03_ios_custom.tf")
	updated := strings.ReplaceAll(config, "Device configuration template test", "Updated template description")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 30*time.Second),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {Source: "hashicorp/random", VersionConstraint: constants.ExternalProviderRandomVersion},
		},
		Steps: []resource.TestStep{
			{PreConfig: func() { testlog.StepAction(resourceType, "Creating ios_custom profile") },
				Config: config, Check: check.That(resourceType + ".test").ExistsInGraph(testResource)},
			importStep(),
			{Config: updated, Check: check.That(resourceType + ".test").Key("description").HasValue("Updated template description")},
			importStep(),
			{Config: updated, PlanOnly: true},
		},
	})
}

func TestAccResourceDeviceConfigurationTemplatesJson_04_IosRoot(t *testing.T) {
	config := loadAcceptanceTestTerraform("resource_04_ios_root.tf")
	updated := strings.ReplaceAll(config, "Device configuration template test", "Updated template description")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 30*time.Second),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {Source: "hashicorp/random", VersionConstraint: constants.ExternalProviderRandomVersion},
		},
		Steps: []resource.TestStep{
			{PreConfig: func() { testlog.StepAction(resourceType, "Creating ios_root profile") },
				Config: config, Check: check.That(resourceType + ".test").ExistsInGraph(testResource)},
			importStep(),
			{Config: updated, Check: check.That(resourceType + ".test").Key("description").HasValue("Updated template description")},
			importStep(),
			{Config: updated, PlanOnly: true},
		},
	})
}

func TestAccResourceDeviceConfigurationTemplatesJson_05_IosScep(t *testing.T) {
	config := loadAcceptanceTestTerraform("resource_05_ios_scep.tf")
	updated := strings.ReplaceAll(config, "Device configuration template test", "Updated template description")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 30*time.Second),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {Source: "hashicorp/random", VersionConstraint: constants.ExternalProviderRandomVersion},
		},
		Steps: []resource.TestStep{
			{PreConfig: func() { testlog.StepAction(resourceType, "Creating ios_scep profile") },
				Config: config, Check: check.That(resourceType + ".test").ExistsInGraph(testResource)},
			importStep(),
			{Config: updated, Check: check.That(resourceType + ".test").Key("description").HasValue("Updated template description")},
			importStep(),
			{Config: updated, PlanOnly: true},
		},
	})
}

func TestAccResourceDeviceConfigurationTemplatesJson_06_MacosCustom(t *testing.T) {
	config := loadAcceptanceTestTerraform("resource_06_macos_custom.tf")
	updated := strings.ReplaceAll(config, "Device configuration template test", "Updated template description")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 30*time.Second),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {Source: "hashicorp/random", VersionConstraint: constants.ExternalProviderRandomVersion},
		},
		Steps: []resource.TestStep{
			{PreConfig: func() { testlog.StepAction(resourceType, "Creating macos_custom profile") },
				Config: config, Check: check.That(resourceType + ".test").ExistsInGraph(testResource)},
			importStep(),
			{Config: updated, Check: check.That(resourceType + ".test").Key("description").HasValue("Updated template description")},
			importStep(),
			{Config: updated, PlanOnly: true},
		},
	})
}

func TestAccResourceDeviceConfigurationTemplatesJson_07_WindowsCustom(t *testing.T) {
	config := loadAcceptanceTestTerraform("resource_07_windows_custom.tf")
	updated := strings.ReplaceAll(config, "Device configuration template test", "Updated template description")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 30*time.Second),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {Source: "hashicorp/random", VersionConstraint: constants.ExternalProviderRandomVersion},
		},
		Steps: []resource.TestStep{
			{PreConfig: func() { testlog.StepAction(resourceType, "Creating windows_custom profile") },
				Config: config, Check: check.That(resourceType + ".test").ExistsInGraph(testResource)},
			importStep(),
			{Config: updated, Check: check.That(resourceType + ".test").Key("description").HasValue("Updated template description")},
			importStep(),
			{Config: updated, PlanOnly: true},
		},
	})
}

func TestAccResourceDeviceConfigurationTemplatesJson_08_WindowsRoot(t *testing.T) {
	config := loadAcceptanceTestTerraform("resource_08_windows_root.tf")
	updated := strings.ReplaceAll(config, "Device configuration template test", "Updated template description")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 30*time.Second),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {Source: "hashicorp/random", VersionConstraint: constants.ExternalProviderRandomVersion},
		},
		Steps: []resource.TestStep{
			{PreConfig: func() { testlog.StepAction(resourceType, "Creating windows_root profile") },
				Config: config, Check: check.That(resourceType + ".test").ExistsInGraph(testResource)},
			importStep(),
			{Config: updated, Check: check.That(resourceType + ".test").Key("description").HasValue("Updated template description")},
			importStep(),
			{Config: updated, PlanOnly: true},
		},
	})
}

func TestAccResourceDeviceConfigurationTemplatesJson_09_AndroidGeneral(t *testing.T) {
	config := loadAcceptanceTestTerraform("resource_09_android_general.tf")
	updated := strings.ReplaceAll(config, "Device configuration template test", "Updated template description")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 30*time.Second),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {Source: "hashicorp/random", VersionConstraint: constants.ExternalProviderRandomVersion},
		},
		Steps: []resource.TestStep{
			{PreConfig: func() { testlog.StepAction(resourceType, "Creating android_general profile") },
				Config: config, Check: check.That(resourceType + ".test").ExistsInGraph(testResource)},
			importStep(),
			{Config: updated, Check: check.That(resourceType + ".test").Key("description").HasValue("Updated template description")},
			importStep(),
			{Config: updated, PlanOnly: true},
		},
	})
}

func TestAccResourceDeviceConfigurationTemplatesJson_10_WindowsEncrypted(t *testing.T) {
	config := loadAcceptanceTestTerraform("resource_10_windows_encrypted.tf")
	updated := strings.ReplaceAll(config, "Device configuration template test", "Updated template description")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 30*time.Second),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {Source: "hashicorp/random", VersionConstraint: constants.ExternalProviderRandomVersion},
		},
		Steps: []resource.TestStep{
			{PreConfig: func() { testlog.StepAction(resourceType, "Creating windows_encrypted profile") },
				Config: config, Check: check.That(resourceType + ".test").ExistsInGraph(testResource)},
			importStep(),
			{Config: updated, Check: check.That(resourceType + ".test").Key("description").HasValue("Updated template description")},
			importStep(),
			{Config: updated, PlanOnly: true},
		},
	})
}

func TestAccResourceDeviceConfigurationTemplatesJson_11_IosNestedFolder(t *testing.T) {
	config := loadAcceptanceTestTerraform("resource_11_ios_nested_folder.tf")
	updated := strings.ReplaceAll(config, "Device configuration template test", "Updated template description")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 30*time.Second),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {Source: "hashicorp/random", VersionConstraint: constants.ExternalProviderRandomVersion},
		},
		Steps: []resource.TestStep{
			{PreConfig: func() { testlog.StepAction(resourceType, "Creating ios_nested_folder profile") },
				Config: config, Check: check.That(resourceType + ".test").ExistsInGraph(testResource)},
			importStep(),
			{Config: updated, Check: check.That(resourceType + ".test").Key("description").HasValue("Updated template description")},
			importStep(),
			{Config: updated, PlanOnly: true},
		},
	})
}

func TestAccResourceDeviceConfigurationTemplatesJson_12_MacosRoot(t *testing.T) {
	config := loadAcceptanceTestTerraform("resource_12_macos_root.tf")
	updated := strings.ReplaceAll(config, "Device configuration template test", "Updated template description")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 30*time.Second),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {Source: "hashicorp/random", VersionConstraint: constants.ExternalProviderRandomVersion},
		},
		Steps: []resource.TestStep{
			{PreConfig: func() { testlog.StepAction(resourceType, "Creating macos_root profile") },
				Config: config, Check: check.That(resourceType + ".test").ExistsInGraph(testResource)},
			importStep(),
			{Config: updated, Check: check.That(resourceType + ".test").Key("description").HasValue("Updated template description")},
			importStep(),
			{Config: updated, PlanOnly: true},
		},
	})
}

func TestAccResourceDeviceConfigurationTemplatesJson_13_MacosScep(t *testing.T) {
	config := loadAcceptanceTestTerraform("resource_13_macos_scep.tf")
	updated := strings.ReplaceAll(config, "Device configuration template test", "Updated template description")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 30*time.Second),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {Source: "hashicorp/random", VersionConstraint: constants.ExternalProviderRandomVersion},
		},
		Steps: []resource.TestStep{
			{PreConfig: func() { testlog.StepAction(resourceType, "Creating macos_scep profile") },
				Config: config, Check: check.That(resourceType + ".test").ExistsInGraph(testResource)},
			importStep(),
			{Config: updated, Check: check.That(resourceType + ".test").Key("description").HasValue("Updated template description")},
			importStep(),
			{Config: updated, PlanOnly: true},
		},
	})
}

func TestAccResourceDeviceConfigurationTemplatesJson_22_WindowsOmaInteger(t *testing.T) {
	config := loadAcceptanceTestTerraform("resource_22_windows_oma_integer.tf")
	updated := strings.ReplaceAll(config, "Device configuration template test", "Updated template description")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 30*time.Second),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {Source: "hashicorp/random", VersionConstraint: constants.ExternalProviderRandomVersion},
		},
		Steps: []resource.TestStep{
			{PreConfig: func() { testlog.StepAction(resourceType, "Creating windows_oma_integer profile") },
				Config: config, Check: check.That(resourceType + ".test").ExistsInGraph(testResource)},
			importStep(),
			{Config: updated, Check: check.That(resourceType + ".test").Key("description").HasValue("Updated template description")},
			importStep(),
			{Config: updated, PlanOnly: true},
		},
	})
}

func TestAccResourceDeviceConfigurationTemplatesJson_23_WindowsOmaBoolean(t *testing.T) {
	config := loadAcceptanceTestTerraform("resource_23_windows_oma_boolean.tf")
	updated := strings.ReplaceAll(config, "Device configuration template test", "Updated template description")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 30*time.Second),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {Source: "hashicorp/random", VersionConstraint: constants.ExternalProviderRandomVersion},
		},
		Steps: []resource.TestStep{
			{PreConfig: func() { testlog.StepAction(resourceType, "Creating windows_oma_boolean profile") },
				Config: config, Check: check.That(resourceType + ".test").ExistsInGraph(testResource)},
			importStep(),
			{Config: updated, Check: check.That(resourceType + ".test").Key("description").HasValue("Updated template description")},
			importStep(),
			{Config: updated, PlanOnly: true},
		},
	})
}

func TestAccResourceDeviceConfigurationTemplatesJson_24_WindowsOmaFloatingpoint(t *testing.T) {
	config := loadAcceptanceTestTerraform("resource_24_windows_oma_floatingpoint.tf")
	updated := strings.ReplaceAll(config, "Device configuration template test", "Updated template description")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 30*time.Second),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {Source: "hashicorp/random", VersionConstraint: constants.ExternalProviderRandomVersion},
		},
		Steps: []resource.TestStep{
			{PreConfig: func() { testlog.StepAction(resourceType, "Creating windows_oma_floatingpoint profile") },
				Config: config, Check: check.That(resourceType + ".test").ExistsInGraph(testResource)},
			importStep(),
			{Config: updated, Check: check.That(resourceType + ".test").Key("description").HasValue("Updated template description")},
			importStep(),
			{Config: updated, PlanOnly: true},
		},
	})
}

func TestAccResourceDeviceConfigurationTemplatesJson_25_WindowsOmaDatetime(t *testing.T) {
	config := loadAcceptanceTestTerraform("resource_25_windows_oma_datetime.tf")
	updated := strings.ReplaceAll(config, "Device configuration template test", "Updated template description")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 30*time.Second),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {Source: "hashicorp/random", VersionConstraint: constants.ExternalProviderRandomVersion},
		},
		Steps: []resource.TestStep{
			{PreConfig: func() { testlog.StepAction(resourceType, "Creating windows_oma_datetime profile") },
				Config: config, Check: check.That(resourceType + ".test").ExistsInGraph(testResource)},
			importStep(),
			{Config: updated, Check: check.That(resourceType + ".test").Key("description").HasValue("Updated template description")},
			importStep(),
			{Config: updated, PlanOnly: true},
		},
	})
}

func TestAccResourceDeviceConfigurationTemplatesJson_26_PlaintextString(t *testing.T) {
	config := loadAcceptanceTestTerraform("resource_26_plaintext_string.tf")
	updated := strings.ReplaceAll(config, "Device configuration template test", "Updated template description")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 30*time.Second),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {Source: "hashicorp/random", VersionConstraint: constants.ExternalProviderRandomVersion},
		},
		Steps: []resource.TestStep{
			{PreConfig: func() { testlog.StepAction(resourceType, "Creating plaintext_string profile") },
				Config: config, Check: check.That(resourceType + ".test").ExistsInGraph(testResource)},
			importStep(),
			{Config: updated, Check: check.That(resourceType + ".test").Key("description").HasValue("Updated template description")},
			importStep(),
			{Config: updated, PlanOnly: true},
		},
	})
}

func TestAccResourceDeviceConfigurationTemplatesJson_27_PlaintextBase64(t *testing.T) {
	config := loadAcceptanceTestTerraform("resource_27_plaintext_base64.tf")
	updated := strings.ReplaceAll(config, "Device configuration template test", "Updated template description")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 30*time.Second),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {Source: "hashicorp/random", VersionConstraint: constants.ExternalProviderRandomVersion},
		},
		Steps: []resource.TestStep{
			{PreConfig: func() { testlog.StepAction(resourceType, "Creating plaintext_base64 profile") },
				Config: config, Check: check.That(resourceType + ".test").ExistsInGraph(testResource)},
			importStep(),
			{Config: updated, Check: check.That(resourceType + ".test").Key("description").HasValue("Updated template description")},
			importStep(),
			{Config: updated, PlanOnly: true},
		},
	})
}

func TestAccResourceDeviceConfigurationTemplatesJson_28_PlaintextStringxml(t *testing.T) {
	config := loadAcceptanceTestTerraform("resource_28_plaintext_stringxml.tf")
	updated := loadAcceptanceTestTerraform("resource_28_plaintext_stringxml_updated.tf")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 30*time.Second),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {Source: "hashicorp/random", VersionConstraint: constants.ExternalProviderRandomVersion},
		},
		Steps: []resource.TestStep{
			{PreConfig: func() { testlog.StepAction(resourceType, "Creating plaintext_stringxml profile") },
				Config: config, Check: resource.ComposeTestCheckFunc(check.That(resourceType+".test").ExistsInGraph(testResource), checkOMASettingValue("#microsoft.graph.omaSettingStringXml", `<test enabled="true"/>`))},
			importStep(),
			{Config: updated, Check: checkOMASettingValue("#microsoft.graph.omaSettingStringXml", "<test>\n  <name>é 日本語</name>\n</test>\n")},
			importStep(),
			{Config: updated, PlanOnly: true},
		},
	})
}

func TestAccResourceDeviceConfigurationTemplatesJson_29_WindowsUpdates(t *testing.T) {
	config := loadAcceptanceTestTerraform("resource_29_windows_updates.tf")
	updated := strings.ReplaceAll(config, "Device configuration template test", "Updated template description")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 30*time.Second),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {Source: "hashicorp/random", VersionConstraint: constants.ExternalProviderRandomVersion},
		},
		Steps: []resource.TestStep{
			{PreConfig: func() { testlog.StepAction(resourceType, "Creating windows_updates profile") },
				Config: config, Check: check.That(resourceType + ".test").ExistsInGraph(testResource)},
			importStep(),
			{Config: updated, Check: check.That(resourceType + ".test").Key("description").HasValue("Updated template description")},
			importStep(),
			{Config: updated, PlanOnly: true},
		},
	})
}

func TestAccResourceDeviceConfigurationTemplatesJson_30_IosWifi(t *testing.T) {
	config := loadAcceptanceTestTerraform("resource_30_ios_wifi.tf")
	updated := strings.ReplaceAll(config, "Device configuration template test", "Updated template description")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 30*time.Second),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {Source: "hashicorp/random", VersionConstraint: constants.ExternalProviderRandomVersion},
		},
		Steps: []resource.TestStep{
			{PreConfig: func() { testlog.StepAction(resourceType, "Creating ios_wifi profile") },
				Config: config, Check: check.That(resourceType + ".test").ExistsInGraph(testResource)},
			{ResourceName: resourceType + ".test", ImportState: true, ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"timeouts", "settings"},
				ImportStateCheck:        checkImportedWifiSettings("get_device_configuration_ios_wifi_import_settings.json")},
			{Config: updated, Check: check.That(resourceType + ".test").Key("description").HasValue("Updated template description")},
			{ResourceName: resourceType + ".test", ImportState: true, ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"timeouts", "settings"},
				ImportStateCheck:        checkImportedWifiSettings("get_device_configuration_ios_wifi_import_settings.json")},
			{Config: updated, PlanOnly: true},
		},
	})
}

func TestAccResourceDeviceConfigurationTemplatesJson_31_MacosWifi(t *testing.T) {
	config := loadAcceptanceTestTerraform("resource_31_macos_wifi.tf")
	updated := strings.ReplaceAll(config, "Device configuration template test", "Updated template description")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 30*time.Second),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {Source: "hashicorp/random", VersionConstraint: constants.ExternalProviderRandomVersion},
		},
		Steps: []resource.TestStep{
			{PreConfig: func() { testlog.StepAction(resourceType, "Creating macos_wifi profile") },
				Config: config, Check: check.That(resourceType + ".test").ExistsInGraph(testResource)},
			{ResourceName: resourceType + ".test", ImportState: true, ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"timeouts", "settings"},
				ImportStateCheck:        checkImportedWifiSettings("get_device_configuration_macos_wifi_import_settings.json")},
			{Config: updated, Check: check.That(resourceType + ".test").Key("description").HasValue("Updated template description")},
			{ResourceName: resourceType + ".test", ImportState: true, ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"timeouts", "settings"},
				ImportStateCheck:        checkImportedWifiSettings("get_device_configuration_macos_wifi_import_settings.json")},
			{Config: updated, PlanOnly: true},
		},
	})
}

func TestAccResourceDeviceConfigurationTemplatesJson_32_WindowsWifi(t *testing.T) {
	config := loadAcceptanceTestTerraform("resource_32_windows_wifi.tf")
	updated := strings.ReplaceAll(config, "Device configuration template test", "Updated template description")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 30*time.Second),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {Source: "hashicorp/random", VersionConstraint: constants.ExternalProviderRandomVersion},
		},
		Steps: []resource.TestStep{
			{PreConfig: func() { testlog.StepAction(resourceType, "Creating windows_wifi profile") },
				Config: config, Check: check.That(resourceType + ".test").ExistsInGraph(testResource)},
			{ResourceName: resourceType + ".test", ImportState: true, ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"timeouts", "settings"},
				ImportStateCheck:        checkImportedWifiSettings("get_device_configuration_windows_wifi_import_settings.json")},
			{Config: updated, Check: check.That(resourceType + ".test").Key("description").HasValue("Updated template description")},
			{ResourceName: resourceType + ".test", ImportState: true, ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"timeouts", "settings"},
				ImportStateCheck:        checkImportedWifiSettings("get_device_configuration_windows_wifi_import_settings.json")},
			{Config: updated, PlanOnly: true},
		},
	})
}

func TestAccResourceDeviceConfigurationTemplatesJson_33_AndroidOwnerWifi(t *testing.T) {
	config := loadAcceptanceTestTerraform("resource_33_android_owner_wifi.tf")
	updated := strings.ReplaceAll(config, "Device configuration template test", "Updated template description")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 30*time.Second),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {Source: "hashicorp/random", VersionConstraint: constants.ExternalProviderRandomVersion},
		},
		Steps: []resource.TestStep{
			{PreConfig: func() { testlog.StepAction(resourceType, "Creating android_owner_wifi profile") },
				Config: config, Check: check.That(resourceType + ".test").ExistsInGraph(testResource)},
			{ResourceName: resourceType + ".test", ImportState: true, ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"timeouts", "settings"},
				ImportStateCheck:        checkImportedWifiSettings("get_device_configuration_android_owner_wifi_import_settings.json")},
			{Config: updated, Check: check.That(resourceType + ".test").Key("description").HasValue("Updated template description")},
			{ResourceName: resourceType + ".test", ImportState: true, ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"timeouts", "settings"},
				ImportStateCheck:        checkImportedWifiSettings("get_device_configuration_android_owner_wifi_import_settings.json")},
			{Config: updated, PlanOnly: true},
		},
	})
}

func TestAccResourceDeviceConfigurationTemplatesJson_34_AndroidWorkWifi(t *testing.T) {
	config := loadAcceptanceTestTerraform("resource_34_android_work_wifi.tf")
	updated := strings.ReplaceAll(config, "Device configuration template test", "Updated template description")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 30*time.Second),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {Source: "hashicorp/random", VersionConstraint: constants.ExternalProviderRandomVersion},
		},
		Steps: []resource.TestStep{
			{PreConfig: func() { testlog.StepAction(resourceType, "Creating android_work_wifi profile") },
				Config: config, Check: check.That(resourceType + ".test").ExistsInGraph(testResource)},
			{ResourceName: resourceType + ".test", ImportState: true, ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"timeouts", "settings"},
				ImportStateCheck:        checkImportedWifiSettings("get_device_configuration_android_work_wifi_import_settings.json")},
			{Config: updated, Check: check.That(resourceType + ".test").Key("description").HasValue("Updated template description")},
			{ResourceName: resourceType + ".test", ImportState: true, ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"timeouts", "settings"},
				ImportStateCheck:        checkImportedWifiSettings("get_device_configuration_android_work_wifi_import_settings.json")},
			{Config: updated, PlanOnly: true},
		},
	})
}

func TestAccResourceDeviceConfigurationTemplatesJson_35_AospWifi(t *testing.T) {
	config := loadAcceptanceTestTerraform("resource_35_aosp_wifi.tf")
	updated := strings.ReplaceAll(config, "Device configuration template test", "Updated template description")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 30*time.Second),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {Source: "hashicorp/random", VersionConstraint: constants.ExternalProviderRandomVersion},
		},
		Steps: []resource.TestStep{
			{PreConfig: func() { testlog.StepAction(resourceType, "Creating aosp_wifi profile") },
				Config: config, Check: check.That(resourceType + ".test").ExistsInGraph(testResource)},
			{ResourceName: resourceType + ".test", ImportState: true, ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"timeouts", "settings"},
				ImportStateCheck:        checkImportedWifiSettings("get_device_configuration_aosp_wifi_import_settings.json")},
			{Config: updated, Check: check.That(resourceType + ".test").Key("description").HasValue("Updated template description")},
			{ResourceName: resourceType + ".test", ImportState: true, ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"timeouts", "settings"},
				ImportStateCheck:        checkImportedWifiSettings("get_device_configuration_aosp_wifi_import_settings.json")},
			{Config: updated, PlanOnly: true},
		},
	})
}

func TestAccResourceDeviceConfigurationTemplatesJson_36_IosEnterpriseWifi(t *testing.T) {
	config := loadAcceptanceTestTerraform("resource_36_ios_enterprise_wifi.tf")
	updated := strings.ReplaceAll(config, "Device configuration template test", "Updated template description")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 30*time.Second),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {Source: "hashicorp/random", VersionConstraint: constants.ExternalProviderRandomVersion},
		},
		Steps: []resource.TestStep{
			{PreConfig: func() { testlog.StepAction(resourceType, "Creating ios_enterprise_wifi profile") },
				Config: config, Check: check.That(resourceType + ".test").ExistsInGraph(testResource)},
			importStep(),
			{Config: updated, Check: check.That(resourceType + ".test").Key("description").HasValue("Updated template description")},
			importStep(),
			{Config: updated, PlanOnly: true},
		},
	})
}

func TestAccResourceDeviceConfigurationTemplatesJson_37_MacosEnterpriseWifi(t *testing.T) {
	config := loadAcceptanceTestTerraform("resource_37_macos_enterprise_wifi.tf")
	updated := strings.ReplaceAll(config, "Device configuration template test", "Updated template description")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 30*time.Second),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {Source: "hashicorp/random", VersionConstraint: constants.ExternalProviderRandomVersion},
		},
		Steps: []resource.TestStep{
			{PreConfig: func() { testlog.StepAction(resourceType, "Creating macos_enterprise_wifi profile") },
				Config: config, Check: check.That(resourceType + ".test").ExistsInGraph(testResource)},
			importStep(),
			{Config: updated, Check: check.That(resourceType + ".test").Key("description").HasValue("Updated template description")},
			importStep(),
			{Config: updated, PlanOnly: true},
		},
	})
}

func TestAccResourceDeviceConfigurationTemplatesJson_38_WindowsEnterpriseWifi(t *testing.T) {
	config := loadAcceptanceTestTerraform("resource_38_windows_enterprise_wifi.tf")
	updated := strings.ReplaceAll(config, "Device configuration template test", "Updated template description")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 30*time.Second),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {Source: "hashicorp/random", VersionConstraint: constants.ExternalProviderRandomVersion},
		},
		Steps: []resource.TestStep{
			{PreConfig: func() { testlog.StepAction(resourceType, "Creating windows_enterprise_wifi profile") },
				Config: config, Check: check.That(resourceType + ".test").ExistsInGraph(testResource)},
			importStep(),
			{Config: updated, Check: check.That(resourceType + ".test").Key("description").HasValue("Updated template description")},
			importStep(),
			{Config: updated, PlanOnly: true},
		},
	})
}

func TestAccResourceDeviceConfigurationTemplatesJson_39_AndroidOwnerEnterpriseWifi(t *testing.T) {
	config := loadAcceptanceTestTerraform("resource_39_android_owner_enterprise_wifi.tf")
	updated := strings.ReplaceAll(config, "Device configuration template test", "Updated template description")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 30*time.Second),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {Source: "hashicorp/random", VersionConstraint: constants.ExternalProviderRandomVersion},
		},
		Steps: []resource.TestStep{
			{PreConfig: func() { testlog.StepAction(resourceType, "Creating android_owner_enterprise_wifi profile") },
				Config: config, Check: check.That(resourceType + ".test").ExistsInGraph(testResource)},
			importStep(),
			{Config: updated, Check: check.That(resourceType + ".test").Key("description").HasValue("Updated template description")},
			importStep(),
			{Config: updated, PlanOnly: true},
		},
	})
}

func TestAccResourceDeviceConfigurationTemplatesJson_40_AndroidWorkEnterpriseWifi(t *testing.T) {
	config := loadAcceptanceTestTerraform("resource_40_android_work_enterprise_wifi.tf")
	updated := strings.ReplaceAll(config, "Device configuration template test", "Updated template description")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 30*time.Second),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {Source: "hashicorp/random", VersionConstraint: constants.ExternalProviderRandomVersion},
		},
		Steps: []resource.TestStep{
			{PreConfig: func() { testlog.StepAction(resourceType, "Creating android_work_enterprise_wifi profile") },
				Config: config, Check: check.That(resourceType + ".test").ExistsInGraph(testResource)},
			importStep(),
			{Config: updated, Check: check.That(resourceType + ".test").Key("description").HasValue("Updated template description")},
			importStep(),
			{Config: updated, PlanOnly: true},
		},
	})
}

func TestAccResourceDeviceConfigurationTemplatesJson_41_AospEnterpriseWifi(t *testing.T) {
	config := loadAcceptanceTestTerraform("resource_41_aosp_enterprise_wifi.tf")
	updated := strings.ReplaceAll(config, "Device configuration template test", "Updated template description")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 30*time.Second),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {Source: "hashicorp/random", VersionConstraint: constants.ExternalProviderRandomVersion},
		},
		Steps: []resource.TestStep{
			{PreConfig: func() { testlog.StepAction(resourceType, "Creating aosp_enterprise_wifi profile") },
				Config: config, Check: check.That(resourceType + ".test").ExistsInGraph(testResource)},
			importStep(),
			{Config: updated, Check: check.That(resourceType + ".test").Key("description").HasValue("Updated template description")},
			importStep(),
			{Config: updated, PlanOnly: true},
		},
	})
}

func TestAccResourceDeviceConfigurationTemplatesJson_42_IosVpn(t *testing.T) {
	config := loadAcceptanceTestTerraform("resource_42_ios_vpn.tf")
	updated := strings.ReplaceAll(config, "Device configuration template test", "Updated template description")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 30*time.Second),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {Source: "hashicorp/random", VersionConstraint: constants.ExternalProviderRandomVersion},
		},
		Steps: []resource.TestStep{
			{PreConfig: func() { testlog.StepAction(resourceType, "Creating ios_vpn profile") },
				Config: config, Check: check.That(resourceType + ".test").ExistsInGraph(testResource)},
			importStep(),
			{Config: updated, Check: check.That(resourceType + ".test").Key("description").HasValue("Updated template description")},
			importStep(),
			{Config: updated, PlanOnly: true},
		},
	})
}

func TestAccResourceDeviceConfigurationTemplatesJson_43_MacosVpn(t *testing.T) {
	config := loadAcceptanceTestTerraform("resource_43_macos_vpn.tf")
	updated := strings.ReplaceAll(config, "Device configuration template test", "Updated template description")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 30*time.Second),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {Source: "hashicorp/random", VersionConstraint: constants.ExternalProviderRandomVersion},
		},
		Steps: []resource.TestStep{
			{PreConfig: func() { testlog.StepAction(resourceType, "Creating macos_vpn profile") },
				Config: config, Check: check.That(resourceType + ".test").ExistsInGraph(testResource)},
			importStep(),
			{Config: updated, Check: check.That(resourceType + ".test").Key("description").HasValue("Updated template description")},
			importStep(),
			{Config: updated, PlanOnly: true},
		},
	})
}

func TestAccResourceDeviceConfigurationTemplatesJson_44_WindowsVpn(t *testing.T) {
	config := loadAcceptanceTestTerraform("resource_44_windows_vpn.tf")
	updated := strings.ReplaceAll(config, "Device configuration template test", "Updated template description")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 30*time.Second),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {Source: "hashicorp/random", VersionConstraint: constants.ExternalProviderRandomVersion},
		},
		Steps: []resource.TestStep{
			{PreConfig: func() { testlog.StepAction(resourceType, "Creating windows_vpn profile") },
				Config: config, Check: check.That(resourceType + ".test").ExistsInGraph(testResource)},
			importStep(),
			{Config: updated, Check: check.That(resourceType + ".test").Key("description").HasValue("Updated template description")},
			importStep(),
			{Config: updated, PlanOnly: true},
		},
	})
}

func TestAccResourceDeviceConfigurationTemplatesJson_45_AndroidOwnerVpn(t *testing.T) {
	config := loadAcceptanceTestTerraform("resource_45_android_owner_vpn.tf")
	updated := strings.ReplaceAll(config, "Device configuration template test", "Updated template description")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 30*time.Second),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {Source: "hashicorp/random", VersionConstraint: constants.ExternalProviderRandomVersion},
		},
		Steps: []resource.TestStep{
			{PreConfig: func() { testlog.StepAction(resourceType, "Creating android_owner_vpn profile") },
				Config: config, Check: check.That(resourceType + ".test").ExistsInGraph(testResource)},
			importStep(),
			{Config: updated, Check: check.That(resourceType + ".test").Key("description").HasValue("Updated template description")},
			importStep(),
			{Config: updated, PlanOnly: true},
		},
	})
}

func TestAccResourceDeviceConfigurationTemplatesJson_46_IosEmail(t *testing.T) {
	config := loadAcceptanceTestTerraform("resource_46_ios_email.tf")
	updated := strings.ReplaceAll(config, "Device configuration template test", "Updated template description")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 30*time.Second),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {Source: "hashicorp/random", VersionConstraint: constants.ExternalProviderRandomVersion},
		},
		Steps: []resource.TestStep{
			{PreConfig: func() { testlog.StepAction(resourceType, "Creating ios_email profile") },
				Config: config, Check: check.That(resourceType + ".test").ExistsInGraph(testResource)},
			importStep(),
			{Config: updated, Check: check.That(resourceType + ".test").Key("description").HasValue("Updated template description")},
			importStep(),
			{Config: updated, PlanOnly: true},
		},
	})
}

func TestAccResourceDeviceConfigurationTemplatesJson_47_IosPkcs(t *testing.T) {
	config := loadAcceptanceTestTerraform("resource_47_ios_pkcs.tf")
	updated := strings.ReplaceAll(config, "Device configuration template test", "Updated template description")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 30*time.Second),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {Source: "hashicorp/random", VersionConstraint: constants.ExternalProviderRandomVersion},
		},
		Steps: []resource.TestStep{
			{PreConfig: func() { testlog.StepAction(resourceType, "Creating ios_pkcs profile") },
				Config: config, Check: check.That(resourceType + ".test").ExistsInGraph(testResource)},
			importStep(),
			{Config: updated, Check: check.That(resourceType + ".test").Key("description").HasValue("Updated template description")},
			importStep(),
			{Config: updated, PlanOnly: true},
		},
	})
}

func TestAccResourceDeviceConfigurationTemplatesJson_48_MacosPkcs(t *testing.T) {
	config := loadAcceptanceTestTerraform("resource_48_macos_pkcs.tf")
	updated := strings.ReplaceAll(config, "Device configuration template test", "Updated template description")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 30*time.Second),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {Source: "hashicorp/random", VersionConstraint: constants.ExternalProviderRandomVersion},
		},
		Steps: []resource.TestStep{
			{PreConfig: func() { testlog.StepAction(resourceType, "Creating macos_pkcs profile") },
				Config: config, Check: check.That(resourceType + ".test").ExistsInGraph(testResource)},
			importStep(),
			{Config: updated, Check: check.That(resourceType + ".test").Key("description").HasValue("Updated template description")},
			importStep(),
			{Config: updated, PlanOnly: true},
		},
	})
}

func TestAccResourceDeviceConfigurationTemplatesJson_49_WindowsPkcs(t *testing.T) {
	config := loadAcceptanceTestTerraform("resource_49_windows_pkcs.tf")
	updated := strings.ReplaceAll(config, "Device configuration template test", "Updated template description")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 30*time.Second),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {Source: "hashicorp/random", VersionConstraint: constants.ExternalProviderRandomVersion},
		},
		Steps: []resource.TestStep{
			{PreConfig: func() { testlog.StepAction(resourceType, "Creating windows_pkcs profile") },
				Config: config, Check: check.That(resourceType + ".test").ExistsInGraph(testResource)},
			importStep(),
			{Config: updated, Check: check.That(resourceType + ".test").Key("description").HasValue("Updated template description")},
			importStep(),
			{Config: updated, PlanOnly: true},
		},
	})
}

func TestAccResourceDeviceConfigurationTemplatesJson_50_AndroidWorkPkcs(t *testing.T) {
	config := loadAcceptanceTestTerraform("resource_50_android_work_pkcs.tf")
	updated := strings.ReplaceAll(config, "Device configuration template test", "Updated template description")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 30*time.Second),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {Source: "hashicorp/random", VersionConstraint: constants.ExternalProviderRandomVersion},
		},
		Steps: []resource.TestStep{
			{PreConfig: func() { testlog.StepAction(resourceType, "Creating android_work_pkcs profile") },
				Config: config, Check: check.That(resourceType + ".test").ExistsInGraph(testResource)},
			importStep(),
			{Config: updated, Check: check.That(resourceType + ".test").Key("description").HasValue("Updated template description")},
			importStep(),
			{Config: updated, PlanOnly: true},
		},
	})
}

func TestAccResourceDeviceConfigurationTemplatesJson_51_MacosPreferences(t *testing.T) {
	config := loadAcceptanceTestTerraform("resource_51_macos_preferences.tf")
	updated := strings.ReplaceAll(config, "Device configuration template test", "Updated template description")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 30*time.Second),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {Source: "hashicorp/random", VersionConstraint: constants.ExternalProviderRandomVersion},
		},
		Steps: []resource.TestStep{
			{PreConfig: func() { testlog.StepAction(resourceType, "Creating macos_preferences profile") },
				Config: config, Check: check.That(resourceType + ".test").ExistsInGraph(testResource)},
			importStep(),
			{Config: updated, Check: check.That(resourceType + ".test").Key("description").HasValue("Updated template description")},
			importStep(),
			{Config: updated, PlanOnly: true},
		},
	})
}

func TestAccResourceDeviceConfigurationTemplatesJson_52_MacosUpdates(t *testing.T) {
	config := loadAcceptanceTestTerraform("resource_52_macos_updates.tf")
	updated := strings.ReplaceAll(config, "Device configuration template test", "Updated template description")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 30*time.Second),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {Source: "hashicorp/random", VersionConstraint: constants.ExternalProviderRandomVersion},
		},
		Steps: []resource.TestStep{
			{PreConfig: func() { testlog.StepAction(resourceType, "Creating macos_updates profile") },
				Config: config, Check: check.That(resourceType + ".test").ExistsInGraph(testResource)},
			importStep(),
			{Config: updated, Check: check.That(resourceType + ".test").Key("description").HasValue("Updated template description")},
			importStep(),
			{Config: updated, PlanOnly: true},
		},
	})
}

func TestAccResourceDeviceConfigurationTemplatesJson_53_AospRestrictions(t *testing.T) {
	config := loadAcceptanceTestTerraform("resource_53_aosp_restrictions.tf")
	updated := strings.ReplaceAll(config, "Device configuration template test", "Updated template description")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 30*time.Second),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {Source: "hashicorp/random", VersionConstraint: constants.ExternalProviderRandomVersion},
		},
		Steps: []resource.TestStep{
			{PreConfig: func() { testlog.StepAction(resourceType, "Creating aosp_restrictions profile") },
				Config: config, Check: check.That(resourceType + ".test").ExistsInGraph(testResource)},
			importStep(),
			{Config: updated, Check: check.That(resourceType + ".test").Key("description").HasValue("Updated template description")},
			importStep(),
			{Config: updated, PlanOnly: true},
		},
	})
}

func TestAccResourceDeviceConfigurationTemplatesJson_54_SingletonRelationships(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 30*time.Second),
		ExternalProviders:        map[string]resource.ExternalProvider{"random": {Source: "hashicorp/random", VersionConstraint: constants.ExternalProviderRandomVersion}},
		Steps: []resource.TestStep{
			{Config: loadAcceptanceTestTerraform("resource_54_singleton_initial.tf")},
			importStep(),
			{Config: loadAcceptanceTestTerraform("resource_54_singleton_replaced.tf")},
			importStep(),
			{Config: loadAcceptanceTestTerraform("resource_54_singleton_removed.tf")},
			importStep(),
			{Config: loadAcceptanceTestTerraform("resource_54_singleton_initial.tf")},
			importStep(),
		},
	})
}

func TestAccResourceDeviceConfigurationTemplatesJson_55_CollectionRelationships(t *testing.T) {
	var settingsBeforeImport string

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 30*time.Second),
		ExternalProviders:        map[string]resource.ExternalProvider{"random": {Source: "hashicorp/random", VersionConstraint: constants.ExternalProviderRandomVersion}},
		Steps: []resource.TestStep{
			{Config: loadAcceptanceTestTerraform("resource_55_collection_initial.tf"), Check: func(state *terraform.State) error {
				settingsBeforeImport = state.RootModule().Resources[resourceType+".test"].Primary.Attributes["settings"]
				return nil
			}},
			{ResourceName: resourceType + ".test", ImportState: true, ImportStateVerify: true, ImportStateVerifyIgnore: []string{"timeouts", "settings"}, ImportStateCheck: checkImportedJSONSettings(&settingsBeforeImport)},
			{Config: loadAcceptanceTestTerraform("resource_55_collection_multiple.tf"), Check: func(state *terraform.State) error {
				settingsBeforeImport = state.RootModule().Resources[resourceType+".test"].Primary.Attributes["settings"]
				return nil
			}},
			{ResourceName: resourceType + ".test", ImportState: true, ImportStateVerify: true, ImportStateVerifyIgnore: []string{"timeouts", "settings"}, ImportStateCheck: checkImportedJSONSettings(&settingsBeforeImport)},
			{Config: loadAcceptanceTestTerraform("resource_55_collection_reordered.tf"), Check: func(state *terraform.State) error {
				settingsBeforeImport = state.RootModule().Resources[resourceType+".test"].Primary.Attributes["settings"]
				return nil
			}},
			{ResourceName: resourceType + ".test", ImportState: true, ImportStateVerify: true, ImportStateVerifyIgnore: []string{"timeouts", "settings"}, ImportStateCheck: checkImportedJSONSettings(&settingsBeforeImport)},
			{Config: loadAcceptanceTestTerraform("resource_55_collection_initial.tf"), Check: func(state *terraform.State) error {
				settingsBeforeImport = state.RootModule().Resources[resourceType+".test"].Primary.Attributes["settings"]
				return nil
			}},
			{ResourceName: resourceType + ".test", ImportState: true, ImportStateVerify: true, ImportStateVerifyIgnore: []string{"timeouts", "settings"}, ImportStateCheck: checkImportedJSONSettings(&settingsBeforeImport)},
			{Config: loadAcceptanceTestTerraform("resource_55_collection_removed.tf"), Check: func(state *terraform.State) error {
				settingsBeforeImport = state.RootModule().Resources[resourceType+".test"].Primary.Attributes["settings"]
				return nil
			}},
			{ResourceName: resourceType + ".test", ImportState: true, ImportStateVerify: true, ImportStateVerifyIgnore: []string{"timeouts", "settings"}, ImportStateCheck: checkImportedJSONSettings(&settingsBeforeImport)},
			{Config: loadAcceptanceTestTerraform("resource_55_collection_multiple.tf"), Check: func(state *terraform.State) error {
				settingsBeforeImport = state.RootModule().Resources[resourceType+".test"].Primary.Attributes["settings"]
				return nil
			}},
			{ResourceName: resourceType + ".test", ImportState: true, ImportStateVerify: true, ImportStateVerifyIgnore: []string{"timeouts", "settings"}, ImportStateCheck: checkImportedJSONSettings(&settingsBeforeImport)},
		},
	})
}

func TestAccResourceDeviceConfigurationTemplatesJson_14_AssignmentTransitions(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 30*time.Second),
		ExternalProviders:        map[string]resource.ExternalProvider{"random": {Source: "hashicorp/random", VersionConstraint: constants.ExternalProviderRandomVersion}, "time": {Source: "hashicorp/time", VersionConstraint: constants.ExternalProviderTimeVersion}},
		Steps: []resource.TestStep{
			{Config: loadAcceptanceTestTerraform("resource_14_assignments_none.tf")},
			{Config: loadAcceptanceTestTerraform("resource_14_assignments_groups.tf"), Check: check.That(resourceType + ".test").Key("assignments.#").HasValue("2")},
			importStep(),
			{Config: loadAcceptanceTestTerraform("resource_14_assignments_empty.tf"), Check: check.That(resourceType + ".test").Key("assignments.#").HasValue("0")},
			{Config: loadAcceptanceTestTerraform("resource_14_assignments_groups.tf")},
			{Config: loadAcceptanceTestTerraform("resource_14_assignments_none.tf"), Check: check.That(resourceType + ".test").Key("assignments.#").DoesNotExist()},
			importStep(),
		},
	})
}

func TestAccResourceDeviceConfigurationTemplatesJson_15_ExplicitResets(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 30*time.Second),
		ExternalProviders:        map[string]resource.ExternalProvider{"random": {Source: "hashicorp/random", VersionConstraint: constants.ExternalProviderRandomVersion}},
		Steps: []resource.TestStep{
			{Config: loadAcceptanceTestTerraform("resource_01_ios_general.tf")},
			{Config: loadAcceptanceTestTerraform("resource_15_explicit_resets.tf"), Check: check.That(resourceType + ".test").Key("settings").MatchesRegex(regexp.MustCompile(`"cameraBlocked":false`))},
			importStep(),
			{Config: loadAcceptanceTestTerraform("resource_01_ios_general.tf"), Check: check.That(resourceType + ".test").Key("settings").MatchesRegex(regexp.MustCompile(`"cameraBlocked":true`))},
			importStep(),
		},
	})
}

func TestAccResourceDeviceConfigurationTemplatesJson_21_FormattedJSON(t *testing.T) {
	var settingsBeforeImport string
	config := loadAcceptanceTestTerraform("resource_21_formatted_json.tf")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 30*time.Second),
		ExternalProviders:        map[string]resource.ExternalProvider{"random": {Source: "hashicorp/random", VersionConstraint: constants.ExternalProviderRandomVersion}},
		Steps: []resource.TestStep{
			{Config: config, Check: func(state *terraform.State) error {
				settingsBeforeImport = state.RootModule().Resources[resourceType+".test"].Primary.Attributes["settings"]
				return nil
			}},
			{ResourceName: resourceType + ".test", ImportState: true, ImportStateVerify: true, ImportStateVerifyIgnore: []string{"timeouts", "settings"}, ImportStateCheck: checkImportedJSONSettings(&settingsBeforeImport)},
			{Config: config, PlanOnly: true},
		},
	})
}

func TestAccResourceDeviceConfigurationTemplatesJson_57_WifiKeyClear(t *testing.T) {
	config := loadAcceptanceTestTerraform("resource_33_android_owner_wifi.tf")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 30*time.Second),
		ExternalProviders:        map[string]resource.ExternalProvider{"random": {Source: "hashicorp/random", VersionConstraint: constants.ExternalProviderRandomVersion}},
		Steps: []resource.TestStep{
			{Config: config},
			{Config: loadAcceptanceTestTerraform("resource_57_wifi_key_clear.tf"), Check: check.That(resourceType + ".test").Key("settings").MatchesRegex(regexp.MustCompile(`"preSharedKey":null`))},
			importStep(),
			{Config: config},
			{Config: config, PlanOnly: true},
		},
	})
}

func TestAccResourceDeviceConfigurationTemplatesJson_61_AssignmentsMinimalToMaximal(t *testing.T) {
	var profileID string
	initial := loadAcceptanceTestTerraform("resource_61_assignments_minimal_to_maximal.tf")
	updated := loadAcceptanceTestTerraform("resource_61_assignments_minimal_to_maximal_updated.tf")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 30*time.Second),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {Source: "hashicorp/random", VersionConstraint: constants.ExternalProviderRandomVersion},
			"time":   {Source: "hashicorp/time", VersionConstraint: constants.ExternalProviderTimeVersion},
		},
		Steps: []resource.TestStep{
			{Config: loadAcceptanceTestTerraform("resource_61_assignment_dependencies.tf"), Check: resource.ComposeTestCheckFunc(check.That("microsoft365_graph_beta_groups_group.include").Key("id").Exists(), check.That("microsoft365_graph_beta_groups_group.exclude").Key("id").Exists())},
			{Config: initial, Check: checkTemplateAssignments(&profileID, false, true)},
			importStep(),
			{Config: initial, PlanOnly: true},
			{Config: updated, Check: checkTemplateAssignments(&profileID, true, true)},
			importStep(),
			{Config: updated, PlanOnly: true},
		},
	})
}

func TestAccResourceDeviceConfigurationTemplatesJson_62_AssignmentsMaximalToMinimal(t *testing.T) {
	var profileID string
	initial := loadAcceptanceTestTerraform("resource_62_assignments_maximal_to_minimal.tf")
	updated := loadAcceptanceTestTerraform("resource_62_assignments_maximal_to_minimal_updated.tf")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 30*time.Second),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {Source: "hashicorp/random", VersionConstraint: constants.ExternalProviderRandomVersion},
			"time":   {Source: "hashicorp/time", VersionConstraint: constants.ExternalProviderTimeVersion},
		},
		Steps: []resource.TestStep{
			{Config: loadAcceptanceTestTerraform("resource_61_assignment_dependencies.tf"), Check: resource.ComposeTestCheckFunc(check.That("microsoft365_graph_beta_groups_group.include").Key("id").Exists(), check.That("microsoft365_graph_beta_groups_group.exclude").Key("id").Exists())},
			{Config: initial, Check: checkTemplateAssignments(&profileID, true, true)},
			importStep(),
			{Config: initial, PlanOnly: true},
			{Config: updated, Check: checkTemplateAssignments(&profileID, false, true)},
			importStep(),
			{Config: updated, PlanOnly: true},
		},
	})
}

func TestAccResourceDeviceConfigurationTemplatesJson_63_MixedOmaValues(t *testing.T) {
	config := loadAcceptanceTestTerraform("resource_63_mixed_oma.tf")
	updated := loadAcceptanceTestTerraform("resource_63_mixed_oma_updated.tf")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 30*time.Second),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {Source: "hashicorp/random", VersionConstraint: constants.ExternalProviderRandomVersion},
		},
		Steps: []resource.TestStep{
			{PreConfig: func() { testlog.StepAction(resourceType, "Creating plaintext_stringxml profile") },
				Config: config, Check: resource.ComposeTestCheckFunc(check.That(resourceType+".test").ExistsInGraph(testResource), resource.ComposeTestCheckFunc(checkOMASettingValue("#microsoft.graph.omaSettingStringXml", `<test enabled="true"/>`), checkOMASettingValue("#microsoft.graph.omaSettingString", "dGVzdA=="), checkOMASettingValue("#microsoft.graph.omaSettingBase64", "AAECA//+/Q==")))},
			importStep(),
			{Config: updated, Check: resource.ComposeTestCheckFunc(checkOMASettingValue("#microsoft.graph.omaSettingStringXml", "<test>\n  <name>é 日本語</name>\n</test>\n"), checkOMASettingValue("#microsoft.graph.omaSettingString", "test"), checkOMASettingValue("#microsoft.graph.omaSettingBase64", "AAECA//+/Q=="))},
			importStep(),
			{Config: updated, PlanOnly: true},
		},
	})
}
