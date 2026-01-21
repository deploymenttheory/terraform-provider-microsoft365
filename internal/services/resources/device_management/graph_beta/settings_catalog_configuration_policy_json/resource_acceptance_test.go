package graphBetaSettingsCatalogConfigurationPolicyJson_test

import (
	"regexp"
	"testing"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/destroy"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/testlog"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaSettingsCatalogConfigurationPolicyJson "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/settings_catalog_configuration_policy_json"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

var (
	resourceType = graphBetaSettingsCatalogConfigurationPolicyJson.ResourceName
	testResource = graphBetaSettingsCatalogConfigurationPolicyJson.SettingsCatalogJsonTestResource{}
)

// loadAcceptanceTestTerraform loads terraform test files from tests/terraform/acceptance
func loadAcceptanceTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/acceptance/" + filename)
	if err != nil {
		panic("failed to load acceptance config " + filename + ": " + err.Error())
	}
	return config
}

// TestAccSettingsCatalogPolicyResource_01_Camera tests a simple choice setting
func TestAccSettingsCatalogPolicyResource_01_Camera(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			resourceType,
			30*time.Second,
		),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
			"time": {
				Source:            "hashicorp/time",
				VersionConstraint: constants.ExternalProviderTimeVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating camera policy")
				},
				Config: loadAcceptanceTestTerraform("resource_01_camera.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("settings catalog policy", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".camera").ExistsInGraph(testResource),
					check.That(resourceType+".camera").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".camera").Key("name").MatchesRegex(regexp.MustCompile(`^acc-test-01-camera-`)),
					check.That(resourceType+".camera").Key("description").HasValue("Acceptance test policy for camera settings"),
					check.That(resourceType+".camera").Key("platforms").HasValue("windows10"),
					check.That(resourceType+".camera").Key("technologies.#").HasValue("1"),
					check.That(resourceType+".camera").Key("settings_count").HasValue("1"),
					check.That(resourceType+".camera").Key("settings").Exists(),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing camera policy")
				},
				ResourceName:            resourceType + ".camera",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// TestAccSettingsCatalogPolicyResource_02_TaskManager tests a simple choice setting
func TestAccSettingsCatalogPolicyResource_02_TaskManager(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			resourceType,
			30*time.Second,
		),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
			"time": {
				Source:            "hashicorp/time",
				VersionConstraint: constants.ExternalProviderTimeVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating task manager policy")
				},
				Config: loadAcceptanceTestTerraform("resource_02_task_manager.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("settings catalog policy", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".task_manager").ExistsInGraph(testResource),
					check.That(resourceType+".task_manager").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".task_manager").Key("name").MatchesRegex(regexp.MustCompile(`^acc-test-02-task-manager-`)),
					check.That(resourceType+".task_manager").Key("settings_count").HasValue("1"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing task manager policy")
				},
				ResourceName:            resourceType + ".task_manager",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// TestAccSettingsCatalogPolicyResource_03_AppPrivacy tests a simple choice setting
func TestAccSettingsCatalogPolicyResource_03_AppPrivacy(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			resourceType,
			30*time.Second,
		),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
			"time": {
				Source:            "hashicorp/time",
				VersionConstraint: constants.ExternalProviderTimeVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating app privacy policy")
				},
				Config: loadAcceptanceTestTerraform("resource_03_app_privacy.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("settings catalog policy", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".app_privacy").ExistsInGraph(testResource),
					check.That(resourceType+".app_privacy").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".app_privacy").Key("name").MatchesRegex(regexp.MustCompile(`^acc-test-03-app-privacy-`)),
					check.That(resourceType+".app_privacy").Key("settings_count").HasValue("1"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing app privacy policy")
				},
				ResourceName:            resourceType + ".app_privacy",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// TestAccSettingsCatalogPolicyResource_04_Cryptography tests a simple choice setting
func TestAccSettingsCatalogPolicyResource_04_Cryptography(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			resourceType,
			30*time.Second,
		),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
			"time": {
				Source:            "hashicorp/time",
				VersionConstraint: constants.ExternalProviderTimeVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating cryptography policy")
				},
				Config: loadAcceptanceTestTerraform("resource_04_cryptography.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("settings catalog policy", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".cryptography").ExistsInGraph(testResource),
					check.That(resourceType+".cryptography").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".cryptography").Key("name").MatchesRegex(regexp.MustCompile(`^acc-test-04-cryptography-`)),
					check.That(resourceType+".cryptography").Key("settings_count").HasValue("1"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing cryptography policy")
				},
				ResourceName:            resourceType + ".cryptography",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// TestAccSettingsCatalogPolicyResource_05_Notifications tests a simple choice setting
func TestAccSettingsCatalogPolicyResource_05_Notifications(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			resourceType,
			30*time.Second,
		),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
			"time": {
				Source:            "hashicorp/time",
				VersionConstraint: constants.ExternalProviderTimeVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating notifications policy")
				},
				Config: loadAcceptanceTestTerraform("resource_05_notifications.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("settings catalog policy", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".notifications").ExistsInGraph(testResource),
					check.That(resourceType+".notifications").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".notifications").Key("name").MatchesRegex(regexp.MustCompile(`^acc-test-05-notifications-`)),
					check.That(resourceType+".notifications").Key("settings_count").HasValue("1"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing notifications policy")
				},
				ResourceName:            resourceType + ".notifications",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// TestAccSettingsCatalogPolicyResource_06_AttachmentManager tests multiple choice settings
func TestAccSettingsCatalogPolicyResource_06_AttachmentManager(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			resourceType,
			30*time.Second,
		),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
			"time": {
				Source:            "hashicorp/time",
				VersionConstraint: constants.ExternalProviderTimeVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating attachment manager policy")
				},
				Config: loadAcceptanceTestTerraform("resource_06_attachment_manager.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("settings catalog policy", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".attachment_manager").ExistsInGraph(testResource),
					check.That(resourceType+".attachment_manager").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".attachment_manager").Key("name").MatchesRegex(regexp.MustCompile(`^acc-test-06-attachment-manager-`)),
					check.That(resourceType+".attachment_manager").Key("settings_count").HasValue("2"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing attachment manager policy")
				},
				ResourceName:            resourceType + ".attachment_manager",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// TestAccSettingsCatalogPolicyResource_07_CredentialUserInterface tests multiple choice settings
func TestAccSettingsCatalogPolicyResource_07_CredentialUserInterface(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			resourceType,
			30*time.Second,
		),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
			"time": {
				Source:            "hashicorp/time",
				VersionConstraint: constants.ExternalProviderTimeVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating credential user interface policy")
				},
				Config: loadAcceptanceTestTerraform("resource_07_credential_user_interface.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("settings catalog policy", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".credential_user_interface").ExistsInGraph(testResource),
					check.That(resourceType+".credential_user_interface").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".credential_user_interface").Key("name").MatchesRegex(regexp.MustCompile(`^acc-test-07-credential-ui-`)),
					check.That(resourceType+".credential_user_interface").Key("settings_count").HasValue("2"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing credential user interface policy")
				},
				ResourceName:            resourceType + ".credential_user_interface",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// TestAccSettingsCatalogPolicyResource_08_RemoteDesktopAVDURL tests simple collection
func TestAccSettingsCatalogPolicyResource_08_RemoteDesktopAVDURL(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			resourceType,
			30*time.Second,
		),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
			"time": {
				Source:            "hashicorp/time",
				VersionConstraint: constants.ExternalProviderTimeVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating remote desktop AVD URL policy")
				},
				Config: loadAcceptanceTestTerraform("resource_08_remote_desktop_avd_url.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("settings catalog policy", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".remote_desktop_avd_url").ExistsInGraph(testResource),
					check.That(resourceType+".remote_desktop_avd_url").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".remote_desktop_avd_url").Key("name").MatchesRegex(regexp.MustCompile(`^acc-test-08-avd-url-`)),
					check.That(resourceType+".remote_desktop_avd_url").Key("settings_count").HasValue("1"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing remote desktop AVD URL policy")
				},
				ResourceName:            resourceType + ".remote_desktop_avd_url",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// TestAccSettingsCatalogPolicyResource_09_StorageSense tests integer settings
func TestAccSettingsCatalogPolicyResource_09_StorageSense(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			resourceType,
			30*time.Second,
		),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
			"time": {
				Source:            "hashicorp/time",
				VersionConstraint: constants.ExternalProviderTimeVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating storage sense policy")
				},
				Config: loadAcceptanceTestTerraform("resource_09_storage_sense.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("settings catalog policy", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".storage_sense").ExistsInGraph(testResource),
					check.That(resourceType+".storage_sense").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".storage_sense").Key("name").MatchesRegex(regexp.MustCompile(`^acc-test-09-storage-sense-`)),
					check.That(resourceType+".storage_sense").Key("settings_count").Exists(),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing storage sense policy")
				},
				ResourceName:            resourceType + ".storage_sense",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// TestAccSettingsCatalogPolicyResource_11_AutoPlayPolicies tests nested choice settings
func TestAccSettingsCatalogPolicyResource_11_AutoPlayPolicies(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			resourceType,
			30*time.Second,
		),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
			"time": {
				Source:            "hashicorp/time",
				VersionConstraint: constants.ExternalProviderTimeVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating AutoPlay policies policy")
				},
				Config: loadAcceptanceTestTerraform("resource_11_autoplay_policies.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("settings catalog policy", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".autoplay").ExistsInGraph(testResource),
					check.That(resourceType+".autoplay").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".autoplay").Key("name").MatchesRegex(regexp.MustCompile(`^acc-test-11-autoplay-`)),
					check.That(resourceType+".autoplay").Key("settings_count").HasValue("3"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing AutoPlay policies policy")
				},
				ResourceName:            resourceType + ".autoplay",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// TestAccSettingsCatalogPolicyResource_12_DefenderSmartscreen tests choice with collection child
func TestAccSettingsCatalogPolicyResource_12_DefenderSmartscreen(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			resourceType,
			30*time.Second,
		),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
			"time": {
				Source:            "hashicorp/time",
				VersionConstraint: constants.ExternalProviderTimeVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating Defender Smartscreen policy")
				},
				Config: loadAcceptanceTestTerraform("resource_12_defender_smartscreen.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("settings catalog policy", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".defender_smartscreen").ExistsInGraph(testResource),
					check.That(resourceType+".defender_smartscreen").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".defender_smartscreen").Key("name").MatchesRegex(regexp.MustCompile(`^acc-test-12-smartscreen-`)),
					check.That(resourceType+".defender_smartscreen").Key("settings_count").HasValue("3"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing Defender Smartscreen policy")
				},
				ResourceName:            resourceType + ".defender_smartscreen",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// TestAccSettingsCatalogPolicyResource_13_EdgeExtensionsMacOS tests multiple collections on macOS
func TestAccSettingsCatalogPolicyResource_13_EdgeExtensionsMacOS(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			resourceType,
			30*time.Second,
		),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
			"time": {
				Source:            "hashicorp/time",
				VersionConstraint: constants.ExternalProviderTimeVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating Edge Extensions macOS policy")
				},
				Config: loadAcceptanceTestTerraform("resource_13_edge_extensions_macos.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("settings catalog policy", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".edge_extensions_macos").ExistsInGraph(testResource),
					check.That(resourceType+".edge_extensions_macos").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".edge_extensions_macos").Key("name").MatchesRegex(regexp.MustCompile(`^acc-test-13-edge-extensions-`)),
					check.That(resourceType+".edge_extensions_macos").Key("settings_count").HasValue("4"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing Edge Extensions macOS policy")
				},
				ResourceName:            resourceType + ".edge_extensions_macos",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// TestAccSettingsCatalogPolicyResource_14_OfficeConfigurationMacOS tests nested group collections on macOS
func TestAccSettingsCatalogPolicyResource_14_OfficeConfigurationMacOS(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			resourceType,
			30*time.Second,
		),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
			"time": {
				Source:            "hashicorp/time",
				VersionConstraint: constants.ExternalProviderTimeVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating Office Configuration macOS policy")
				},
				Config: loadAcceptanceTestTerraform("resource_14_office_configuration_macos.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("settings catalog policy", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".office_configuration_macos").ExistsInGraph(testResource),
					check.That(resourceType+".office_configuration_macos").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".office_configuration_macos").Key("name").MatchesRegex(regexp.MustCompile(`^acc-test-14-office-macos-`)),
					check.That(resourceType+".office_configuration_macos").Key("settings_count").HasValue("3"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing Office Configuration macOS policy")
				},
				ResourceName:            resourceType + ".office_configuration_macos",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// TestAccSettingsCatalogPolicyResource_10_WindowsConnectionManager tests choice with children
func TestAccSettingsCatalogPolicyResource_10_WindowsConnectionManager(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			resourceType,
			30*time.Second,
		),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
			"time": {
				Source:            "hashicorp/time",
				VersionConstraint: constants.ExternalProviderTimeVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating Windows Connection Manager policy")
				},
				Config: loadAcceptanceTestTerraform("resource_10_windows_connection_manager.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("settings catalog policy", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".windows_connection_manager").ExistsInGraph(testResource),
					check.That(resourceType+".windows_connection_manager").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".windows_connection_manager").Key("name").MatchesRegex(regexp.MustCompile(`^acc-test-10-wcm-`)),
					check.That(resourceType+".windows_connection_manager").Key("settings_count").HasValue("4"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing Windows Connection Manager policy")
				},
				ResourceName:            resourceType + ".windows_connection_manager",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// TestAccSettingsCatalogPolicyResource_15_DefenderAntivirusBaseline tests complex nested structures
func TestAccSettingsCatalogPolicyResource_15_DefenderAntivirusBaseline(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			resourceType,
			30*time.Second,
		),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
			"time": {
				Source:            "hashicorp/time",
				VersionConstraint: constants.ExternalProviderTimeVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating Defender Antivirus baseline policy")
				},
				Config: loadAcceptanceTestTerraform("resource_15_defender_antivirus_baseline.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("settings catalog policy", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".defender_antivirus_baseline").ExistsInGraph(testResource),
					check.That(resourceType+".defender_antivirus_baseline").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".defender_antivirus_baseline").Key("name").MatchesRegex(regexp.MustCompile(`^acc-test-15-defender-baseline-`)),
					check.That(resourceType+".defender_antivirus_baseline").Key("settings_count").HasValue("9"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing Defender Antivirus baseline policy")
				},
				ResourceName:            resourceType + ".defender_antivirus_baseline",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// TestAccSettingsCatalogPolicyResource_16_FileExplorerMinimalAssignments tests minimal group assignments
func TestAccSettingsCatalogPolicyResource_16_FileExplorerMinimalAssignments(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			resourceType,
			30*time.Second,
		),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
			"time": {
				Source:            "hashicorp/time",
				VersionConstraint: constants.ExternalProviderTimeVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating File Explorer policy with minimal assignments")
				},
				Config: loadAcceptanceTestTerraform("resource_16_file_explorer_minimal_assignments.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("settings catalog policy with assignments", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".file_explorer_minimal_assignments").ExistsInGraph(testResource),
					check.That(resourceType+".file_explorer_minimal_assignments").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".file_explorer_minimal_assignments").Key("name").MatchesRegex(regexp.MustCompile(`^acc-test-16-file-explorer-`)),
					check.That(resourceType+".file_explorer_minimal_assignments").Key("assignments.#").HasValue("1"),
					check.That(resourceType+".file_explorer_minimal_assignments").Key("is_assigned").HasValue("true"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing File Explorer policy")
				},
				ResourceName:            resourceType + ".file_explorer_minimal_assignments",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// TestAccSettingsCatalogPolicyResource_17_LocalPoliciesMaximalAssignments tests policy with maximal group assignments
func TestAccSettingsCatalogPolicyResource_17_LocalPoliciesMaximalAssignments(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			resourceType,
			30*time.Second,
		),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
			"time": {
				Source:            "hashicorp/time",
				VersionConstraint: constants.ExternalProviderTimeVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating Local Policies with maximal assignments")
				},
				Config: loadAcceptanceTestTerraform("resource_17_local_policies_maximal_assignments.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("settings catalog policy with assignments", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".local_policies_maximal").ExistsInGraph(testResource),
					check.That(resourceType+".local_policies_maximal").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".local_policies_maximal").Key("name").MatchesRegex(regexp.MustCompile(`^acc-test-17-local-policies-`)),
					check.That(resourceType+".local_policies_maximal").Key("assignments.#").HasValue("4"),
					check.That(resourceType+".local_policies_maximal").Key("is_assigned").HasValue("true"),
					check.That(resourceType+".local_policies_maximal").Key("settings_count").HasValue("6"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing Local Policies with maximal assignments")
				},
				ResourceName:            resourceType + ".local_policies_maximal",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}
