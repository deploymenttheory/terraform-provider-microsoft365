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

func TestAccResourceDeviceConfigurationTemplatesJson_65_WindowsDeliveryOptimization(t *testing.T) {
	config := loadAcceptanceTestTerraform("resource_65_windows_delivery_optimization.tf")
	updated := loadAcceptanceTestTerraform("resource_65_windows_delivery_optimization_updated.tf")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 30*time.Second),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {Source: "hashicorp/random", VersionConstraint: constants.ExternalProviderRandomVersion},
		},
		Steps: []resource.TestStep{
			{PreConfig: func() { testlog.StepAction(resourceType, "Creating windows_delivery_optimization profile") },
				Config: config, Check: resource.ComposeTestCheckFunc(check.That(resourceType+".test").ExistsInGraph(testResource), check.That(resourceType+".test").Key("settings").MatchesRegex(regexp.MustCompile(`"@odata.type":"#microsoft.graph.windowsDeliveryOptimizationConfiguration"`)))},
			importStep(),
			{Config: config, PlanOnly: true},
			{Config: updated, Check: resource.ComposeTestCheckFunc(check.That(resourceType+".test").Key("description").HasValue("Updated template description"),
				check.That(resourceType+".test").Key("settings").MatchesRegex(regexp.MustCompile(`"backgroundDownloadFromHttpDelayInSeconds":30`)))},
			importStep(),
			{Config: updated, PlanOnly: true},
		},
	})
}

func TestAccResourceDeviceConfigurationTemplatesJson_66_WindowsDfci(t *testing.T) {
	config := loadAcceptanceTestTerraform("resource_66_windows_dfci.tf")
	updated := loadAcceptanceTestTerraform("resource_66_windows_dfci_updated.tf")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 30*time.Second),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {Source: "hashicorp/random", VersionConstraint: constants.ExternalProviderRandomVersion},
		},
		Steps: []resource.TestStep{
			{PreConfig: func() { testlog.StepAction(resourceType, "Creating windows_dfci profile") },
				Config: config, Check: resource.ComposeTestCheckFunc(check.That(resourceType+".test").ExistsInGraph(testResource), check.That(resourceType+".test").Key("settings").MatchesRegex(regexp.MustCompile(`"@odata.type":"#microsoft.graph.windows10DeviceFirmwareConfigurationInterface"`)))},
			importStep(),
			{Config: config, PlanOnly: true},
			{Config: updated, Check: resource.ComposeTestCheckFunc(check.That(resourceType+".test").Key("description").HasValue("Updated template description"),
				check.That(resourceType+".test").Key("settings").MatchesRegex(regexp.MustCompile(`"cameras":"disabled"`)))},
			importStep(),
			{Config: updated, PlanOnly: true},
		},
	})
}

func TestAccResourceDeviceConfigurationTemplatesJson_67_WindowsDeviceRestrictions(t *testing.T) {
	config := loadAcceptanceTestTerraform("resource_67_windows_device_restrictions.tf")
	updated := loadAcceptanceTestTerraform("resource_67_windows_device_restrictions_updated.tf")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 30*time.Second),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {Source: "hashicorp/random", VersionConstraint: constants.ExternalProviderRandomVersion},
		},
		Steps: []resource.TestStep{
			{PreConfig: func() { testlog.StepAction(resourceType, "Creating windows_device_restrictions profile") },
				Config: config, Check: resource.ComposeTestCheckFunc(check.That(resourceType+".test").ExistsInGraph(testResource), check.That(resourceType+".test").Key("settings").MatchesRegex(regexp.MustCompile(`"@odata.type":"#microsoft.graph.windows10GeneralConfiguration"`)))},
			importStep(),
			{Config: config, PlanOnly: true},
			{Config: updated, Check: resource.ComposeTestCheckFunc(check.That(resourceType+".test").Key("description").HasValue("Updated template description"),
				check.That(resourceType+".test").Key("settings").MatchesRegex(regexp.MustCompile(`"cameraBlocked":true`)))},
			importStep(),
			{Config: updated, PlanOnly: true},
		},
	})
}

func TestAccResourceDeviceConfigurationTemplatesJson_68_WindowsDomainJoin(t *testing.T) {
	config := loadAcceptanceTestTerraform("resource_68_windows_domain_join.tf")
	updated := loadAcceptanceTestTerraform("resource_68_windows_domain_join_updated.tf")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 30*time.Second),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {Source: "hashicorp/random", VersionConstraint: constants.ExternalProviderRandomVersion},
		},
		Steps: []resource.TestStep{
			{PreConfig: func() { testlog.StepAction(resourceType, "Creating windows_domain_join profile") },
				Config: config, Check: resource.ComposeTestCheckFunc(check.That(resourceType+".test").ExistsInGraph(testResource), check.That(resourceType+".test").Key("settings").MatchesRegex(regexp.MustCompile(`"@odata.type":"#microsoft.graph.windowsDomainJoinConfiguration"`)))},
			importStep(),
			{Config: config, PlanOnly: true},
			{Config: updated, Check: resource.ComposeTestCheckFunc(check.That(resourceType+".test").Key("description").HasValue("Updated template description"),
				check.That(resourceType+".test").Key("settings").MatchesRegex(regexp.MustCompile(`"computerNameStaticPrefix":"LAB"`)))},
			importStep(),
			{Config: updated, PlanOnly: true},
		},
	})
}

func TestAccResourceDeviceConfigurationTemplatesJson_69_WindowsEditionUpgrade(t *testing.T) {
	config := loadAcceptanceTestTerraform("resource_69_windows_edition_upgrade.tf")
	updated := loadAcceptanceTestTerraform("resource_69_windows_edition_upgrade_updated.tf")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 30*time.Second),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {Source: "hashicorp/random", VersionConstraint: constants.ExternalProviderRandomVersion},
		},
		Steps: []resource.TestStep{
			{PreConfig: func() { testlog.StepAction(resourceType, "Creating windows_edition_upgrade profile") },
				Config: config, Check: resource.ComposeTestCheckFunc(check.That(resourceType+".test").ExistsInGraph(testResource), check.That(resourceType+".test").Key("settings").MatchesRegex(regexp.MustCompile(`"@odata.type":"#microsoft.graph.editionUpgradeConfiguration"`)))},
			importStep(),
			{Config: config, PlanOnly: true},
			{Config: updated, Check: resource.ComposeTestCheckFunc(check.That(resourceType+".test").Key("description").HasValue("Updated template description"),
				check.That(resourceType+".test").Key("settings").MatchesRegex(regexp.MustCompile(`"windowsSMode":"unlock"`)))},
			importStep(),
			{Config: updated, PlanOnly: true},
		},
	})
}

func TestAccResourceDeviceConfigurationTemplatesJson_70_WindowsEmail(t *testing.T) {
	config := loadAcceptanceTestTerraform("resource_70_windows_email.tf")
	updated := loadAcceptanceTestTerraform("resource_70_windows_email_updated.tf")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 30*time.Second),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {Source: "hashicorp/random", VersionConstraint: constants.ExternalProviderRandomVersion},
		},
		Steps: []resource.TestStep{
			{PreConfig: func() { testlog.StepAction(resourceType, "Creating windows_email profile") },
				Config: config, Check: resource.ComposeTestCheckFunc(check.That(resourceType+".test").ExistsInGraph(testResource), check.That(resourceType+".test").Key("settings").MatchesRegex(regexp.MustCompile(`"@odata.type":"#microsoft.graph.windows10EasEmailProfileConfiguration"`)))},
			importStep(),
			{Config: config, PlanOnly: true},
			{Config: updated, Check: resource.ComposeTestCheckFunc(check.That(resourceType+".test").Key("description").HasValue("Updated template description"),
				check.That(resourceType+".test").Key("settings").MatchesRegex(regexp.MustCompile(`"syncCalendar":true`)))},
			importStep(),
			{Config: updated, PlanOnly: true},
		},
	})
}

func TestAccResourceDeviceConfigurationTemplatesJson_71_WindowsEndpointProtection(t *testing.T) {
	config := loadAcceptanceTestTerraform("resource_71_windows_endpoint_protection.tf")
	updated := loadAcceptanceTestTerraform("resource_71_windows_endpoint_protection_updated.tf")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 30*time.Second),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {Source: "hashicorp/random", VersionConstraint: constants.ExternalProviderRandomVersion},
		},
		Steps: []resource.TestStep{
			{PreConfig: func() { testlog.StepAction(resourceType, "Creating windows_endpoint_protection profile") },
				Config: config, Check: resource.ComposeTestCheckFunc(check.That(resourceType+".test").ExistsInGraph(testResource), check.That(resourceType+".test").Key("settings").MatchesRegex(regexp.MustCompile(`"@odata.type":"#microsoft.graph.windows10EndpointProtectionConfiguration"`)))},
			importStep(),
			{Config: config, PlanOnly: true},
			{Config: updated, Check: resource.ComposeTestCheckFunc(check.That(resourceType+".test").Key("description").HasValue("Updated template description"),
				check.That(resourceType+".test").Key("settings").MatchesRegex(regexp.MustCompile(`"firewallBlockStatefulFTP":true`)))},
			importStep(),
			{Config: updated, PlanOnly: true},
		},
	})
}

func TestAccResourceDeviceConfigurationTemplatesJson_72_WindowsKiosk(t *testing.T) {
	config := loadAcceptanceTestTerraform("resource_72_windows_kiosk.tf")
	updated := loadAcceptanceTestTerraform("resource_72_windows_kiosk_updated.tf")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 30*time.Second),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {Source: "hashicorp/random", VersionConstraint: constants.ExternalProviderRandomVersion},
		},
		Steps: []resource.TestStep{
			{PreConfig: func() { testlog.StepAction(resourceType, "Creating windows_kiosk profile") },
				Config: config, Check: resource.ComposeTestCheckFunc(check.That(resourceType+".test").ExistsInGraph(testResource), check.That(resourceType+".test").Key("settings").MatchesRegex(regexp.MustCompile(`"@odata.type":"#microsoft.graph.windowsKioskConfiguration"`)))},
			importStep(),
			{Config: config, PlanOnly: true},
			{Config: updated, Check: resource.ComposeTestCheckFunc(check.That(resourceType + ".test").Key("description").HasValue("Updated template description"))},
			importStep(),
			{Config: updated, PlanOnly: true},
		},
	})
}

func TestAccResourceDeviceConfigurationTemplatesJson_73_WindowsImportedPkcs(t *testing.T) {
	config := loadAcceptanceTestTerraform("resource_73_windows_imported_pkcs.tf")
	updated := loadAcceptanceTestTerraform("resource_73_windows_imported_pkcs_updated.tf")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 30*time.Second),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {Source: "hashicorp/random", VersionConstraint: constants.ExternalProviderRandomVersion},
		},
		Steps: []resource.TestStep{
			{PreConfig: func() { testlog.StepAction(resourceType, "Creating windows_imported_pkcs profile") },
				Config: config, Check: resource.ComposeTestCheckFunc(check.That(resourceType+".test").ExistsInGraph(testResource), check.That(resourceType+".test").Key("settings").MatchesRegex(regexp.MustCompile(`"@odata.type":"#microsoft.graph.windows10ImportedPFXCertificateProfile"`)))},
			importStep(),
			{Config: config, PlanOnly: true},
			{Config: updated, Check: resource.ComposeTestCheckFunc(check.That(resourceType+".test").Key("description").HasValue("Updated template description"),
				check.That(resourceType+".test").Key("settings").MatchesRegex(regexp.MustCompile(`"intendedPurpose":"smimeSigning"`)))},
			importStep(),
			{Config: updated, PlanOnly: true},
		},
	})
}

func TestAccResourceDeviceConfigurationTemplatesJson_74_WindowsScep(t *testing.T) {
	config := loadAcceptanceTestTerraform("resource_74_windows_scep.tf")
	updated := loadAcceptanceTestTerraform("resource_74_windows_scep_updated.tf")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 30*time.Second),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {Source: "hashicorp/random", VersionConstraint: constants.ExternalProviderRandomVersion},
		},
		Steps: []resource.TestStep{
			{PreConfig: func() { testlog.StepAction(resourceType, "Creating windows_scep profile") },
				Config: config, Check: resource.ComposeTestCheckFunc(check.That(resourceType+".test").ExistsInGraph(testResource), check.That(resourceType+".test").Key("settings").MatchesRegex(regexp.MustCompile(`"@odata.type":"#microsoft.graph.windows81SCEPCertificateProfile"`)))},
			importStep(),
			{Config: config, PlanOnly: true},
			{Config: updated, Check: resource.ComposeTestCheckFunc(check.That(resourceType+".test").Key("description").HasValue("Updated template description"),
				check.That(resourceType+".test").Key("settings").MatchesRegex(regexp.MustCompile(`"renewalThresholdPercentage":30`)))},
			importStep(),
			{Config: updated, PlanOnly: true},
		},
	})
}

func TestAccResourceDeviceConfigurationTemplatesJson_75_WindowsSecureAssessment(t *testing.T) {
	config := loadAcceptanceTestTerraform("resource_75_windows_secure_assessment.tf")
	updated := loadAcceptanceTestTerraform("resource_75_windows_secure_assessment_updated.tf")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 30*time.Second),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {Source: "hashicorp/random", VersionConstraint: constants.ExternalProviderRandomVersion},
		},
		Steps: []resource.TestStep{
			{PreConfig: func() { testlog.StepAction(resourceType, "Creating windows_secure_assessment profile") },
				Config: config, Check: resource.ComposeTestCheckFunc(check.That(resourceType+".test").ExistsInGraph(testResource), check.That(resourceType+".test").Key("settings").MatchesRegex(regexp.MustCompile(`"@odata.type":"#microsoft.graph.windows10SecureAssessmentConfiguration"`)))},
			importStep(),
			{Config: config, PlanOnly: true},
			{Config: updated, Check: resource.ComposeTestCheckFunc(check.That(resourceType+".test").Key("description").HasValue("Updated template description"),
				check.That(resourceType+".test").Key("settings").MatchesRegex(regexp.MustCompile(`"allowScreenCapture":true`)))},
			importStep(),
			{Config: updated, PlanOnly: true},
		},
	})
}

func TestAccResourceDeviceConfigurationTemplatesJson_76_WindowsSharedDevice(t *testing.T) {
	config := loadAcceptanceTestTerraform("resource_76_windows_shared_device.tf")
	updated := loadAcceptanceTestTerraform("resource_76_windows_shared_device_updated.tf")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 30*time.Second),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {Source: "hashicorp/random", VersionConstraint: constants.ExternalProviderRandomVersion},
		},
		Steps: []resource.TestStep{
			{PreConfig: func() { testlog.StepAction(resourceType, "Creating windows_shared_device profile") },
				Config: config, Check: resource.ComposeTestCheckFunc(check.That(resourceType+".test").ExistsInGraph(testResource), check.That(resourceType+".test").Key("settings").MatchesRegex(regexp.MustCompile(`"@odata.type":"#microsoft.graph.sharedPCConfiguration"`)))},
			importStep(),
			{Config: config, PlanOnly: true},
			{Config: updated, Check: resource.ComposeTestCheckFunc(check.That(resourceType+".test").Key("description").HasValue("Updated template description"),
				check.That(resourceType+".test").Key("settings").MatchesRegex(regexp.MustCompile(`"disableAccountManager":true`)))},
			importStep(),
			{Config: updated, PlanOnly: true},
		},
	})
}

func TestAccResourceDeviceConfigurationTemplatesJson_77_WindowsHealthMonitoring(t *testing.T) {
	config := loadAcceptanceTestTerraform("resource_77_windows_health_monitoring.tf")
	updated := loadAcceptanceTestTerraform("resource_77_windows_health_monitoring_updated.tf")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 30*time.Second),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {Source: "hashicorp/random", VersionConstraint: constants.ExternalProviderRandomVersion},
		},
		Steps: []resource.TestStep{
			{PreConfig: func() { testlog.StepAction(resourceType, "Creating windows_health_monitoring profile") },
				Config: config, Check: resource.ComposeTestCheckFunc(check.That(resourceType+".test").ExistsInGraph(testResource), check.That(resourceType+".test").Key("settings").MatchesRegex(regexp.MustCompile(`"@odata.type":"#microsoft.graph.windowsHealthMonitoringConfiguration"`)))},
			importStep(),
			{Config: config, PlanOnly: true},
			{Config: updated, Check: resource.ComposeTestCheckFunc(check.That(resourceType+".test").Key("description").HasValue("Updated template description"),
				check.That(resourceType+".test").Key("settings").MatchesRegex(regexp.MustCompile(`"allowDeviceHealthMonitoring":"enabled"`)))},
			importStep(),
			{Config: updated, PlanOnly: true},
		},
	})
}

func TestAccResourceDeviceConfigurationTemplatesJson_78_WindowsWiredNetwork(t *testing.T) {
	config := loadAcceptanceTestTerraform("resource_78_windows_wired_network.tf")
	updated := loadAcceptanceTestTerraform("resource_78_windows_wired_network_updated.tf")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 30*time.Second),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {Source: "hashicorp/random", VersionConstraint: constants.ExternalProviderRandomVersion},
		},
		Steps: []resource.TestStep{
			{PreConfig: func() { testlog.StepAction(resourceType, "Creating windows_wired_network profile") },
				Config: config, Check: resource.ComposeTestCheckFunc(check.That(resourceType+".test").ExistsInGraph(testResource), check.That(resourceType+".test").Key("settings").MatchesRegex(regexp.MustCompile(`"@odata.type":"#microsoft.graph.windowsWiredNetworkConfiguration"`)))},
			importStep(),
			{Config: config, PlanOnly: true},
			{Config: updated, Check: resource.ComposeTestCheckFunc(check.That(resourceType+".test").Key("description").HasValue("Updated template description"),
				check.That(resourceType+".test").Key("settings").MatchesRegex(regexp.MustCompile(`"maximumAuthenticationFailures":2`)))},
			importStep(),
			{Config: updated, PlanOnly: true},
		},
	})
}

func TestAccResourceDeviceConfigurationTemplatesJson_79_IosDerivedCredential(t *testing.T) {
	config := loadAcceptanceTestTerraform("resource_79_ios_derived_credential.tf")
	updated := loadAcceptanceTestTerraform("resource_79_ios_derived_credential_updated.tf")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 30*time.Second),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {Source: "hashicorp/random", VersionConstraint: constants.ExternalProviderRandomVersion},
		},
		Steps: []resource.TestStep{
			{PreConfig: func() { testlog.StepAction(resourceType, "Creating ios_derived_credential profile") },
				Config: config, Check: resource.ComposeTestCheckFunc(check.That(resourceType+".test").ExistsInGraph(testResource), check.That(resourceType+".test").Key("settings").MatchesRegex(regexp.MustCompile(`"@odata.type":"#microsoft.graph.iosDerivedCredentialAuthenticationConfiguration"`)))},
			importStep(),
			{Config: config, PlanOnly: true},
			{Config: updated, Check: resource.ComposeTestCheckFunc(check.That(resourceType + ".test").Key("description").HasValue("Updated template description"))},
			importStep(),
			{Config: updated, PlanOnly: true},
		},
	})
}

func TestAccResourceDeviceConfigurationTemplatesJson_80_IosImportedPkcs(t *testing.T) {
	config := loadAcceptanceTestTerraform("resource_80_ios_imported_pkcs.tf")
	updated := loadAcceptanceTestTerraform("resource_80_ios_imported_pkcs_updated.tf")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 30*time.Second),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {Source: "hashicorp/random", VersionConstraint: constants.ExternalProviderRandomVersion},
		},
		Steps: []resource.TestStep{
			{PreConfig: func() { testlog.StepAction(resourceType, "Creating ios_imported_pkcs profile") },
				Config: config, Check: resource.ComposeTestCheckFunc(check.That(resourceType+".test").ExistsInGraph(testResource), check.That(resourceType+".test").Key("settings").MatchesRegex(regexp.MustCompile(`"@odata.type":"#microsoft.graph.iosImportedPFXCertificateProfile"`)))},
			importStep(),
			{Config: config, PlanOnly: true},
			{Config: updated, Check: resource.ComposeTestCheckFunc(check.That(resourceType+".test").Key("description").HasValue("Updated template description"),
				check.That(resourceType+".test").Key("settings").MatchesRegex(regexp.MustCompile(`"intendedPurpose":"smimeSigning"`)))},
			importStep(),
			{Config: updated, PlanOnly: true},
		},
	})
}

func TestAccResourceDeviceConfigurationTemplatesJson_81_IosEducation(t *testing.T) {
	config := loadAcceptanceTestTerraform("resource_81_ios_education.tf")
	updated := loadAcceptanceTestTerraform("resource_81_ios_education_updated.tf")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 30*time.Second),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {Source: "hashicorp/random", VersionConstraint: constants.ExternalProviderRandomVersion},
		},
		Steps: []resource.TestStep{
			{PreConfig: func() { testlog.StepAction(resourceType, "Creating ios_education profile") },
				Config: config, Check: resource.ComposeTestCheckFunc(check.That(resourceType+".test").ExistsInGraph(testResource), check.That(resourceType+".test").Key("settings").MatchesRegex(regexp.MustCompile(`"@odata.type":"#microsoft.graph.iosEduDeviceConfiguration"`)))},
			importStep(),
			{Config: config, PlanOnly: true},
			{Config: updated, Check: resource.ComposeTestCheckFunc(check.That(resourceType + ".test").Key("description").HasValue("Updated template description"))},
			importStep(),
			{Config: updated, PlanOnly: true},
		},
	})
}

func TestAccResourceDeviceConfigurationTemplatesJson_82_IosWiredNetwork(t *testing.T) {
	config := loadAcceptanceTestTerraform("resource_82_ios_wired_network.tf")
	updated := loadAcceptanceTestTerraform("resource_82_ios_wired_network_updated.tf")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 30*time.Second),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {Source: "hashicorp/random", VersionConstraint: constants.ExternalProviderRandomVersion},
		},
		Steps: []resource.TestStep{
			{PreConfig: func() { testlog.StepAction(resourceType, "Creating ios_wired_network profile") },
				Config: config, Check: resource.ComposeTestCheckFunc(check.That(resourceType+".test").ExistsInGraph(testResource), check.That(resourceType+".test").Key("settings").MatchesRegex(regexp.MustCompile(`"@odata.type":"#microsoft.graph.iosWiredNetworkConfiguration"`)))},
			importStep(),
			{Config: config, PlanOnly: true},
			{Config: updated, Check: resource.ComposeTestCheckFunc(check.That(resourceType+".test").Key("description").HasValue("Updated template description"),
				check.That(resourceType+".test").Key("settings").MatchesRegex(regexp.MustCompile(`"networkName":"Updated wired network"`)))},
			importStep(),
			{Config: updated, PlanOnly: true},
		},
	})
}

func TestAccResourceDeviceConfigurationTemplatesJson_83_MacosDeviceFeatures(t *testing.T) {
	config := loadAcceptanceTestTerraform("resource_83_macos_device_features.tf")
	updated := loadAcceptanceTestTerraform("resource_83_macos_device_features_updated.tf")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 30*time.Second),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {Source: "hashicorp/random", VersionConstraint: constants.ExternalProviderRandomVersion},
		},
		Steps: []resource.TestStep{
			{PreConfig: func() { testlog.StepAction(resourceType, "Creating macos_device_features profile") },
				Config: config, Check: resource.ComposeTestCheckFunc(check.That(resourceType+".test").ExistsInGraph(testResource), check.That(resourceType+".test").Key("settings").MatchesRegex(regexp.MustCompile(`"@odata.type":"#microsoft.graph.macOSDeviceFeaturesConfiguration"`)))},
			importStep(),
			{Config: config, PlanOnly: true},
			{Config: updated, Check: resource.ComposeTestCheckFunc(check.That(resourceType+".test").Key("description").HasValue("Updated template description"),
				check.That(resourceType+".test").Key("settings").MatchesRegex(regexp.MustCompile(`"adminShowHostInfo":true`)))},
			importStep(),
			{Config: updated, PlanOnly: true},
		},
	})
}

func TestAccResourceDeviceConfigurationTemplatesJson_84_MacosDeviceRestrictions(t *testing.T) {
	config := loadAcceptanceTestTerraform("resource_84_macos_device_restrictions.tf")
	updated := loadAcceptanceTestTerraform("resource_84_macos_device_restrictions_updated.tf")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 30*time.Second),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {Source: "hashicorp/random", VersionConstraint: constants.ExternalProviderRandomVersion},
		},
		Steps: []resource.TestStep{
			{PreConfig: func() { testlog.StepAction(resourceType, "Creating macos_device_restrictions profile") },
				Config: config, Check: resource.ComposeTestCheckFunc(check.That(resourceType+".test").ExistsInGraph(testResource), check.That(resourceType+".test").Key("settings").MatchesRegex(regexp.MustCompile(`"@odata.type":"#microsoft.graph.macOSGeneralDeviceConfiguration"`)))},
			importStep(),
			{Config: config, PlanOnly: true},
			{Config: updated, Check: resource.ComposeTestCheckFunc(check.That(resourceType+".test").Key("description").HasValue("Updated template description"),
				check.That(resourceType+".test").Key("settings").MatchesRegex(regexp.MustCompile(`"cameraBlocked":true`)))},
			importStep(),
			{Config: updated, PlanOnly: true},
		},
	})
}

func TestAccResourceDeviceConfigurationTemplatesJson_85_MacosEndpointProtection(t *testing.T) {
	config := loadAcceptanceTestTerraform("resource_85_macos_endpoint_protection.tf")
	updated := loadAcceptanceTestTerraform("resource_85_macos_endpoint_protection_updated.tf")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 30*time.Second),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {Source: "hashicorp/random", VersionConstraint: constants.ExternalProviderRandomVersion},
		},
		Steps: []resource.TestStep{
			{PreConfig: func() { testlog.StepAction(resourceType, "Creating macos_endpoint_protection profile") },
				Config: config, Check: resource.ComposeTestCheckFunc(check.That(resourceType+".test").ExistsInGraph(testResource), check.That(resourceType+".test").Key("settings").MatchesRegex(regexp.MustCompile(`"@odata.type":"#microsoft.graph.macOSEndpointProtectionConfiguration"`)))},
			importStep(),
			{Config: config, PlanOnly: true},
			{Config: updated, Check: resource.ComposeTestCheckFunc(check.That(resourceType+".test").Key("description").HasValue("Updated template description"),
				check.That(resourceType+".test").Key("settings").MatchesRegex(regexp.MustCompile(`"firewallEnabled":true`)))},
			importStep(),
			{Config: updated, PlanOnly: true},
		},
	})
}

func TestAccResourceDeviceConfigurationTemplatesJson_86_MacosExtensions(t *testing.T) {
	config := loadAcceptanceTestTerraform("resource_86_macos_extensions.tf")
	updated := loadAcceptanceTestTerraform("resource_86_macos_extensions_updated.tf")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 30*time.Second),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {Source: "hashicorp/random", VersionConstraint: constants.ExternalProviderRandomVersion},
		},
		Steps: []resource.TestStep{
			{PreConfig: func() { testlog.StepAction(resourceType, "Creating macos_extensions profile") },
				Config: config, Check: resource.ComposeTestCheckFunc(check.That(resourceType+".test").ExistsInGraph(testResource), check.That(resourceType+".test").Key("settings").MatchesRegex(regexp.MustCompile(`"@odata.type":"#microsoft.graph.macOSExtensionsConfiguration"`)))},
			importStep(),
			{Config: config, PlanOnly: true},
			{Config: updated, Check: resource.ComposeTestCheckFunc(check.That(resourceType+".test").Key("description").HasValue("Updated template description"),
				check.That(resourceType+".test").Key("settings").MatchesRegex(regexp.MustCompile(`"kernelExtensionOverridesAllowed":true`)))},
			importStep(),
			{Config: updated, PlanOnly: true},
		},
	})
}

func TestAccResourceDeviceConfigurationTemplatesJson_87_MacosImportedPkcs(t *testing.T) {
	config := loadAcceptanceTestTerraform("resource_87_macos_imported_pkcs.tf")
	updated := loadAcceptanceTestTerraform("resource_87_macos_imported_pkcs_updated.tf")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 30*time.Second),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {Source: "hashicorp/random", VersionConstraint: constants.ExternalProviderRandomVersion},
		},
		Steps: []resource.TestStep{
			{PreConfig: func() { testlog.StepAction(resourceType, "Creating macos_imported_pkcs profile") },
				Config: config, Check: resource.ComposeTestCheckFunc(check.That(resourceType+".test").ExistsInGraph(testResource), check.That(resourceType+".test").Key("settings").MatchesRegex(regexp.MustCompile(`"@odata.type":"#microsoft.graph.macOSImportedPFXCertificateProfile"`)))},
			importStep(),
			{Config: config, PlanOnly: true},
			{Config: updated, Check: resource.ComposeTestCheckFunc(check.That(resourceType+".test").Key("description").HasValue("Updated template description"),
				check.That(resourceType+".test").Key("settings").MatchesRegex(regexp.MustCompile(`"intendedPurpose":"smimeSigning"`)))},
			importStep(),
			{Config: updated, PlanOnly: true},
		},
	})
}

func TestAccResourceDeviceConfigurationTemplatesJson_88_MacosWiredNetwork(t *testing.T) {
	config := loadAcceptanceTestTerraform("resource_88_macos_wired_network.tf")
	updated := loadAcceptanceTestTerraform("resource_88_macos_wired_network_updated.tf")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 30*time.Second),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {Source: "hashicorp/random", VersionConstraint: constants.ExternalProviderRandomVersion},
		},
		Steps: []resource.TestStep{
			{PreConfig: func() { testlog.StepAction(resourceType, "Creating macos_wired_network profile") },
				Config: config, Check: resource.ComposeTestCheckFunc(check.That(resourceType+".test").ExistsInGraph(testResource), check.That(resourceType+".test").Key("settings").MatchesRegex(regexp.MustCompile(`"@odata.type":"#microsoft.graph.macOSWiredNetworkConfiguration"`)))},
			importStep(),
			{Config: config, PlanOnly: true},
			{Config: updated, Check: resource.ComposeTestCheckFunc(check.That(resourceType+".test").Key("description").HasValue("Updated template description"),
				check.That(resourceType+".test").Key("settings").MatchesRegex(regexp.MustCompile(`"networkName":"Updated wired network"`)))},
			importStep(),
			{Config: updated, PlanOnly: true},
		},
	})
}

func TestAccResourceDeviceConfigurationTemplatesJson_89_AospPkcs(t *testing.T) {
	config := loadAcceptanceTestTerraform("resource_89_aosp_pkcs.tf")
	updated := loadAcceptanceTestTerraform("resource_89_aosp_pkcs_updated.tf")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 30*time.Second),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {Source: "hashicorp/random", VersionConstraint: constants.ExternalProviderRandomVersion},
		},
		Steps: []resource.TestStep{
			{PreConfig: func() { testlog.StepAction(resourceType, "Creating aosp_pkcs profile") },
				Config: config, Check: resource.ComposeTestCheckFunc(check.That(resourceType+".test").ExistsInGraph(testResource), check.That(resourceType+".test").Key("settings").MatchesRegex(regexp.MustCompile(`"@odata.type":"#microsoft.graph.aospDeviceOwnerPkcsCertificateProfile"`)))},
			importStep(),
			{Config: config, PlanOnly: true},
			{Config: updated, Check: resource.ComposeTestCheckFunc(check.That(resourceType+".test").Key("description").HasValue("Updated template description"),
				check.That(resourceType+".test").Key("settings").MatchesRegex(regexp.MustCompile(`"renewalThresholdPercentage":30`)))},
			importStep(),
			{Config: updated, PlanOnly: true},
		},
	})
}

func TestAccResourceDeviceConfigurationTemplatesJson_90_AndroidOwnerPkcs(t *testing.T) {
	config := loadAcceptanceTestTerraform("resource_90_android_owner_pkcs.tf")
	updated := loadAcceptanceTestTerraform("resource_90_android_owner_pkcs_updated.tf")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 30*time.Second),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {Source: "hashicorp/random", VersionConstraint: constants.ExternalProviderRandomVersion},
		},
		Steps: []resource.TestStep{
			{PreConfig: func() { testlog.StepAction(resourceType, "Creating android_owner_pkcs profile") },
				Config: config, Check: resource.ComposeTestCheckFunc(check.That(resourceType+".test").ExistsInGraph(testResource), check.That(resourceType+".test").Key("settings").MatchesRegex(regexp.MustCompile(`"@odata.type":"#microsoft.graph.androidDeviceOwnerPkcsCertificateProfile"`)))},
			importStep(),
			{Config: config, PlanOnly: true},
			{Config: updated, Check: resource.ComposeTestCheckFunc(check.That(resourceType+".test").Key("description").HasValue("Updated template description"),
				check.That(resourceType+".test").Key("settings").MatchesRegex(regexp.MustCompile(`"renewalThresholdPercentage":30`)))},
			importStep(),
			{Config: updated, PlanOnly: true},
		},
	})
}

func TestAccResourceDeviceConfigurationTemplatesJson_91_AndroidOwnerImportedPkcs(t *testing.T) {
	config := loadAcceptanceTestTerraform("resource_91_android_owner_imported_pkcs.tf")
	updated := loadAcceptanceTestTerraform("resource_91_android_owner_imported_pkcs_updated.tf")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 30*time.Second),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {Source: "hashicorp/random", VersionConstraint: constants.ExternalProviderRandomVersion},
		},
		Steps: []resource.TestStep{
			{PreConfig: func() { testlog.StepAction(resourceType, "Creating android_owner_imported_pkcs profile") },
				Config: config, Check: resource.ComposeTestCheckFunc(check.That(resourceType+".test").ExistsInGraph(testResource), check.That(resourceType+".test").Key("settings").MatchesRegex(regexp.MustCompile(`"@odata.type":"#microsoft.graph.androidDeviceOwnerImportedPFXCertificateProfile"`)))},
			importStep(),
			{Config: config, PlanOnly: true},
			{Config: updated, Check: resource.ComposeTestCheckFunc(check.That(resourceType+".test").Key("description").HasValue("Updated template description"),
				check.That(resourceType+".test").Key("settings").MatchesRegex(regexp.MustCompile(`"intendedPurpose":"smimeSigning"`)))},
			importStep(),
			{Config: updated, PlanOnly: true},
		},
	})
}

func TestAccResourceDeviceConfigurationTemplatesJson_92_AndroidOwnerDerivedCredential(t *testing.T) {
	config := loadAcceptanceTestTerraform("resource_92_android_owner_derived_credential.tf")
	updated := loadAcceptanceTestTerraform("resource_92_android_owner_derived_credential_updated.tf")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 30*time.Second),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {Source: "hashicorp/random", VersionConstraint: constants.ExternalProviderRandomVersion},
		},
		Steps: []resource.TestStep{
			{PreConfig: func() { testlog.StepAction(resourceType, "Creating android_owner_derived_credential profile") },
				Config: config, Check: resource.ComposeTestCheckFunc(check.That(resourceType+".test").ExistsInGraph(testResource), check.That(resourceType+".test").Key("settings").MatchesRegex(regexp.MustCompile(`"@odata.type":"#microsoft.graph.androidDeviceOwnerDerivedCredentialAuthenticationConfiguration"`)))},
			importStep(),
			{Config: config, PlanOnly: true},
			{Config: updated, Check: resource.ComposeTestCheckFunc(check.That(resourceType + ".test").Key("description").HasValue("Updated template description"))},
			importStep(),
			{Config: updated, PlanOnly: true},
		},
	})
}

func TestAccResourceDeviceConfigurationTemplatesJson_93_AndroidWorkDeviceRestrictions(t *testing.T) {
	config := loadAcceptanceTestTerraform("resource_93_android_work_device_restrictions.tf")
	updated := loadAcceptanceTestTerraform("resource_93_android_work_device_restrictions_updated.tf")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 30*time.Second),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {Source: "hashicorp/random", VersionConstraint: constants.ExternalProviderRandomVersion},
		},
		Steps: []resource.TestStep{
			{PreConfig: func() { testlog.StepAction(resourceType, "Creating android_work_device_restrictions profile") },
				Config: config, Check: resource.ComposeTestCheckFunc(check.That(resourceType+".test").ExistsInGraph(testResource), check.That(resourceType+".test").Key("settings").MatchesRegex(regexp.MustCompile(`"@odata.type":"#microsoft.graph.androidWorkProfileGeneralDeviceConfiguration"`)))},
			importStep(),
			{Config: config, PlanOnly: true},
			{Config: updated, Check: resource.ComposeTestCheckFunc(check.That(resourceType+".test").Key("description").HasValue("Updated template description"),
				check.That(resourceType+".test").Key("settings").MatchesRegex(regexp.MustCompile(`"workProfileBlockCamera":true`)))},
			importStep(),
			{Config: updated, PlanOnly: true},
		},
	})
}

func TestAccResourceDeviceConfigurationTemplatesJson_94_AndroidWorkGmail(t *testing.T) {
	config := loadAcceptanceTestTerraform("resource_94_android_work_gmail.tf")
	updated := loadAcceptanceTestTerraform("resource_94_android_work_gmail_updated.tf")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 30*time.Second),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {Source: "hashicorp/random", VersionConstraint: constants.ExternalProviderRandomVersion},
		},
		Steps: []resource.TestStep{
			{PreConfig: func() { testlog.StepAction(resourceType, "Creating android_work_gmail profile") },
				Config: config, Check: resource.ComposeTestCheckFunc(check.That(resourceType+".test").ExistsInGraph(testResource), check.That(resourceType+".test").Key("settings").MatchesRegex(regexp.MustCompile(`"@odata.type":"#microsoft.graph.androidWorkProfileGmailEasConfiguration"`)))},
			importStep(),
			{Config: config, PlanOnly: true},
			{Config: updated, Check: resource.ComposeTestCheckFunc(check.That(resourceType+".test").Key("description").HasValue("Updated template description"),
				check.That(resourceType+".test").Key("settings").MatchesRegex(regexp.MustCompile(`"hostName":"updated-mail.example.invalid"`)))},
			importStep(),
			{Config: updated, PlanOnly: true},
		},
	})
}

func TestAccResourceDeviceConfigurationTemplatesJson_95_AndroidWorkNineEmail(t *testing.T) {
	config := loadAcceptanceTestTerraform("resource_95_android_work_nine_email.tf")
	updated := loadAcceptanceTestTerraform("resource_95_android_work_nine_email_updated.tf")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 30*time.Second),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {Source: "hashicorp/random", VersionConstraint: constants.ExternalProviderRandomVersion},
		},
		Steps: []resource.TestStep{
			{PreConfig: func() { testlog.StepAction(resourceType, "Creating android_work_nine_email profile") },
				Config: config, Check: resource.ComposeTestCheckFunc(check.That(resourceType+".test").ExistsInGraph(testResource), check.That(resourceType+".test").Key("settings").MatchesRegex(regexp.MustCompile(`"@odata.type":"#microsoft.graph.androidWorkProfileNineWorkEasConfiguration"`)))},
			importStep(),
			{Config: config, PlanOnly: true},
			{Config: updated, Check: resource.ComposeTestCheckFunc(check.That(resourceType+".test").Key("description").HasValue("Updated template description"),
				check.That(resourceType+".test").Key("settings").MatchesRegex(regexp.MustCompile(`"syncCalendar":true`)))},
			importStep(),
			{Config: updated, PlanOnly: true},
		},
	})
}

func TestAccResourceDeviceConfigurationTemplatesJson_96_AndroidWorkVpn(t *testing.T) {
	config := loadAcceptanceTestTerraform("resource_96_android_work_vpn.tf")
	updated := loadAcceptanceTestTerraform("resource_96_android_work_vpn_updated.tf")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 30*time.Second),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {Source: "hashicorp/random", VersionConstraint: constants.ExternalProviderRandomVersion},
		},
		Steps: []resource.TestStep{
			{PreConfig: func() { testlog.StepAction(resourceType, "Creating android_work_vpn profile") },
				Config: config, Check: resource.ComposeTestCheckFunc(check.That(resourceType+".test").ExistsInGraph(testResource), check.That(resourceType+".test").Key("settings").MatchesRegex(regexp.MustCompile(`"@odata.type":"#microsoft.graph.androidWorkProfileVpnConfiguration"`)))},
			importStep(),
			{Config: config, PlanOnly: true},
			{Config: updated, Check: resource.ComposeTestCheckFunc(check.That(resourceType+".test").Key("description").HasValue("Updated template description"),
				check.That(resourceType+".test").Key("settings").MatchesRegex(regexp.MustCompile(`"connectionName":"Updated VPN"`)))},
			importStep(),
			{Config: updated, PlanOnly: true},
		},
	})
}

func TestAccResourceDeviceConfigurationTemplatesJson_97_AndroidWorkImportedPkcs(t *testing.T) {
	config := loadAcceptanceTestTerraform("resource_97_android_work_imported_pkcs.tf")
	updated := loadAcceptanceTestTerraform("resource_97_android_work_imported_pkcs_updated.tf")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 30*time.Second),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {Source: "hashicorp/random", VersionConstraint: constants.ExternalProviderRandomVersion},
		},
		Steps: []resource.TestStep{
			{PreConfig: func() { testlog.StepAction(resourceType, "Creating android_work_imported_pkcs profile") },
				Config: config, Check: resource.ComposeTestCheckFunc(check.That(resourceType+".test").ExistsInGraph(testResource), check.That(resourceType+".test").Key("settings").MatchesRegex(regexp.MustCompile(`"@odata.type":"#microsoft.graph.androidForWorkImportedPFXCertificateProfile"`)))},
			importStep(),
			{Config: config, PlanOnly: true},
			{Config: updated, Check: resource.ComposeTestCheckFunc(check.That(resourceType+".test").Key("description").HasValue("Updated template description"),
				check.That(resourceType+".test").Key("settings").MatchesRegex(regexp.MustCompile(`"intendedPurpose":"smimeSigning"`)))},
			importStep(),
			{Config: updated, PlanOnly: true},
		},
	})
}

func TestAccResourceDeviceConfigurationTemplatesJson_98_AospRoot(t *testing.T) {
	config := loadAcceptanceTestTerraform("resource_98_aosp_root.tf")
	updated := loadAcceptanceTestTerraform("resource_98_aosp_root_updated.tf")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 30*time.Second),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {Source: "hashicorp/random", VersionConstraint: constants.ExternalProviderRandomVersion},
		},
		Steps: []resource.TestStep{
			{PreConfig: func() { testlog.StepAction(resourceType, "Creating aosp_root profile") },
				Config: config, Check: resource.ComposeTestCheckFunc(check.That(resourceType+".test").ExistsInGraph(testResource), check.That(resourceType+".test").Key("settings").MatchesRegex(regexp.MustCompile(`"@odata.type":"#microsoft.graph.aospDeviceOwnerTrustedRootCertificate"`)))},
			importStep(),
			{Config: config, PlanOnly: true},
			{Config: updated, Check: resource.ComposeTestCheckFunc(check.That(resourceType+".test").Key("description").HasValue("Updated template description"),
				check.That(resourceType+".test").Key("settings").MatchesRegex(regexp.MustCompile(`"certFileName":"updated-root.cer"`)))},
			importStep(),
			{Config: updated, PlanOnly: true},
		},
	})
}

func TestAccResourceDeviceConfigurationTemplatesJson_99_AospScep(t *testing.T) {
	config := loadAcceptanceTestTerraform("resource_99_aosp_scep.tf")
	updated := loadAcceptanceTestTerraform("resource_99_aosp_scep_updated.tf")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 30*time.Second),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {Source: "hashicorp/random", VersionConstraint: constants.ExternalProviderRandomVersion},
		},
		Steps: []resource.TestStep{
			{PreConfig: func() { testlog.StepAction(resourceType, "Creating aosp_scep profile") },
				Config: config, Check: resource.ComposeTestCheckFunc(check.That(resourceType+".test").ExistsInGraph(testResource), check.That(resourceType+".test").Key("settings").MatchesRegex(regexp.MustCompile(`"@odata.type":"#microsoft.graph.aospDeviceOwnerScepCertificateProfile"`)))},
			importStep(),
			{Config: config, PlanOnly: true},
			{Config: updated, Check: resource.ComposeTestCheckFunc(check.That(resourceType+".test").Key("description").HasValue("Updated template description"),
				check.That(resourceType+".test").Key("settings").MatchesRegex(regexp.MustCompile(`"renewalThresholdPercentage":30`)))},
			importStep(),
			{Config: updated, PlanOnly: true},
		},
	})
}

func TestAccResourceDeviceConfigurationTemplatesJson_100_AndroidOwnerRoot(t *testing.T) {
	config := loadAcceptanceTestTerraform("resource_100_android_owner_root.tf")
	updated := loadAcceptanceTestTerraform("resource_100_android_owner_root_updated.tf")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 30*time.Second),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {Source: "hashicorp/random", VersionConstraint: constants.ExternalProviderRandomVersion},
		},
		Steps: []resource.TestStep{
			{PreConfig: func() { testlog.StepAction(resourceType, "Creating android_owner_root profile") },
				Config: config, Check: resource.ComposeTestCheckFunc(check.That(resourceType+".test").ExistsInGraph(testResource), check.That(resourceType+".test").Key("settings").MatchesRegex(regexp.MustCompile(`"@odata.type":"#microsoft.graph.androidDeviceOwnerTrustedRootCertificate"`)))},
			importStep(),
			{Config: config, PlanOnly: true},
			{Config: updated, Check: resource.ComposeTestCheckFunc(check.That(resourceType+".test").Key("description").HasValue("Updated template description"),
				check.That(resourceType+".test").Key("settings").MatchesRegex(regexp.MustCompile(`"certFileName":"updated-root.cer"`)))},
			importStep(),
			{Config: updated, PlanOnly: true},
		},
	})
}

func TestAccResourceDeviceConfigurationTemplatesJson_101_AndroidOwnerScep(t *testing.T) {
	config := loadAcceptanceTestTerraform("resource_101_android_owner_scep.tf")
	updated := loadAcceptanceTestTerraform("resource_101_android_owner_scep_updated.tf")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 30*time.Second),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {Source: "hashicorp/random", VersionConstraint: constants.ExternalProviderRandomVersion},
		},
		Steps: []resource.TestStep{
			{PreConfig: func() { testlog.StepAction(resourceType, "Creating android_owner_scep profile") },
				Config: config, Check: resource.ComposeTestCheckFunc(check.That(resourceType+".test").ExistsInGraph(testResource), check.That(resourceType+".test").Key("settings").MatchesRegex(regexp.MustCompile(`"@odata.type":"#microsoft.graph.androidDeviceOwnerScepCertificateProfile"`)))},
			importStep(),
			{Config: config, PlanOnly: true},
			{Config: updated, Check: resource.ComposeTestCheckFunc(check.That(resourceType+".test").Key("description").HasValue("Updated template description"),
				check.That(resourceType+".test").Key("settings").MatchesRegex(regexp.MustCompile(`"renewalThresholdPercentage":30`)))},
			importStep(),
			{Config: updated, PlanOnly: true},
		},
	})
}

func TestAccResourceDeviceConfigurationTemplatesJson_102_AndroidWorkRoot(t *testing.T) {
	config := loadAcceptanceTestTerraform("resource_102_android_work_root.tf")
	updated := loadAcceptanceTestTerraform("resource_102_android_work_root_updated.tf")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 30*time.Second),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {Source: "hashicorp/random", VersionConstraint: constants.ExternalProviderRandomVersion},
		},
		Steps: []resource.TestStep{
			{PreConfig: func() { testlog.StepAction(resourceType, "Creating android_work_root profile") },
				Config: config, Check: resource.ComposeTestCheckFunc(check.That(resourceType+".test").ExistsInGraph(testResource), check.That(resourceType+".test").Key("settings").MatchesRegex(regexp.MustCompile(`"@odata.type":"#microsoft.graph.androidWorkProfileTrustedRootCertificate"`)))},
			importStep(),
			{Config: config, PlanOnly: true},
			{Config: updated, Check: resource.ComposeTestCheckFunc(check.That(resourceType+".test").Key("description").HasValue("Updated template description"),
				check.That(resourceType+".test").Key("settings").MatchesRegex(regexp.MustCompile(`"certFileName":"updated-root.cer"`)))},
			importStep(),
			{Config: updated, PlanOnly: true},
		},
	})
}

func TestAccResourceDeviceConfigurationTemplatesJson_103_AndroidWorkScep(t *testing.T) {
	config := loadAcceptanceTestTerraform("resource_103_android_work_scep.tf")
	updated := loadAcceptanceTestTerraform("resource_103_android_work_scep_updated.tf")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 30*time.Second),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {Source: "hashicorp/random", VersionConstraint: constants.ExternalProviderRandomVersion},
		},
		Steps: []resource.TestStep{
			{PreConfig: func() { testlog.StepAction(resourceType, "Creating android_work_scep profile") },
				Config: config, Check: resource.ComposeTestCheckFunc(check.That(resourceType+".test").ExistsInGraph(testResource), check.That(resourceType+".test").Key("settings").MatchesRegex(regexp.MustCompile(`"@odata.type":"#microsoft.graph.androidWorkProfileScepCertificateProfile"`)))},
			importStep(),
			{Config: config, PlanOnly: true},
			{Config: updated, Check: resource.ComposeTestCheckFunc(check.That(resourceType+".test").Key("description").HasValue("Updated template description"),
				check.That(resourceType+".test").Key("settings").MatchesRegex(regexp.MustCompile(`"renewalThresholdPercentage":30`)))},
			importStep(),
			{Config: updated, PlanOnly: true},
		},
	})
}

func TestAccResourceDeviceConfigurationTemplatesJson_104_AndroidWorkMigration(t *testing.T) {
	config := loadAcceptanceTestTerraform("resource_104_android_work_migration.tf")
	updated := loadAcceptanceTestTerraform("resource_104_android_work_migration_updated.tf")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 30*time.Second),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {Source: "hashicorp/random", VersionConstraint: constants.ExternalProviderRandomVersion},
		},
		Steps: []resource.TestStep{
			{PreConfig: func() { testlog.StepAction(resourceType, "Creating android_work_migration profile") },
				Config: config, Check: resource.ComposeTestCheckFunc(check.That(resourceType+".test").ExistsInGraph(testResource), check.That(resourceType+".test").Key("settings").MatchesRegex(regexp.MustCompile(`"@odata.type":"#microsoft.graph.androidWorkProfileMigrationConfiguration"`)))},
			importStep(),
			{Config: config, PlanOnly: true},
			{Config: updated, Check: resource.ComposeTestCheckFunc(check.That(resourceType+".test").Key("description").HasValue("Updated template description"),
				check.That(resourceType+".test").Key("settings").MatchesRegex(regexp.MustCompile(`"disableMigration":true`)))},
			importStep(),
			{Config: updated, PlanOnly: true},
		},
	})
}

func TestAccResourceDeviceConfigurationTemplatesJson_105_IosUpdates(t *testing.T) {
	config := loadAcceptanceTestTerraform("resource_105_ios_updates.tf")
	updated := loadAcceptanceTestTerraform("resource_105_ios_updates_updated.tf")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 30*time.Second),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {Source: "hashicorp/random", VersionConstraint: constants.ExternalProviderRandomVersion},
		},
		Steps: []resource.TestStep{
			{PreConfig: func() { testlog.StepAction(resourceType, "Creating ios_updates profile") },
				Config: config, Check: resource.ComposeTestCheckFunc(check.That(resourceType+".test").ExistsInGraph(testResource), check.That(resourceType+".test").Key("settings").MatchesRegex(regexp.MustCompile(`"@odata.type":"#microsoft.graph.iosUpdateConfiguration"`)))},
			importStep(),
			{Config: config, PlanOnly: true},
			{Config: updated, Check: resource.ComposeTestCheckFunc(check.That(resourceType+".test").Key("description").HasValue("Updated template description"),
				check.That(resourceType+".test").Key("settings").MatchesRegex(regexp.MustCompile(`"isEnabled":true`)))},
			importStep(),
			{Config: updated, PlanOnly: true},
		},
	})
}
