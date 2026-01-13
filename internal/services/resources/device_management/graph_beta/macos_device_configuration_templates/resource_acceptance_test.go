package graphBetaMacosDeviceConfigurationTemplates_test

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
	graphBetaMacosDeviceConfigurationTemplates "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/macos_device_configuration_templates"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

var (
	resourceType = graphBetaMacosDeviceConfigurationTemplates.ResourceName
	testResource = graphBetaMacosDeviceConfigurationTemplates.MacosDeviceConfigurationTemplatesTestResource{}
)

func loadAcceptanceTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/acceptance/" + filename)
	if err != nil {
		panic("failed to load acceptance config " + filename + ": " + err.Error())
	}
	return config
}

// Custom Configuration Tests
func TestAccMacosDeviceConfigurationTemplatesResource_CustomConfiguration(t *testing.T) {
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
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating custom configuration")
				},
				Config: loadAcceptanceTestTerraform("resource_custom_configuration_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("macOS device configuration", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".custom_configuration_example").ExistsInGraph(testResource),
					check.That(resourceType+".custom_configuration_example").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".custom_configuration_example").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-macOS-custom-config-[a-z0-9]{8}$`)),
					check.That(resourceType+".custom_configuration_example").Key("description").HasValue("Example custom configuration template for macOS devices"),
					check.That(resourceType+".custom_configuration_example").Key("custom_configuration.deployment_channel").HasValue("deviceChannel"),
					check.That(resourceType+".custom_configuration_example").Key("custom_configuration.payload_file_name").HasValue("com.example.custom.mobileconfig"),
					check.That(resourceType+".custom_configuration_example").Key("custom_configuration.payload_name").HasValue("Custom Configuration Example"),
					check.That(resourceType+".custom_configuration_example").Key("custom_configuration.payload").Exists(),
					check.That(resourceType+".custom_configuration_example").Key("assignments.#").HasValue("4"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing custom configuration")
				},
				ResourceName:      resourceType + ".custom_configuration_example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// Preference File Tests
func TestAccMacosDeviceConfigurationTemplatesResource_PreferenceFile(t *testing.T) {
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
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating preference file configuration")
				},
				Config: loadAcceptanceTestTerraform("resource_preference_file_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("macOS device configuration", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".preference_file_example").ExistsInGraph(testResource),
					check.That(resourceType+".preference_file_example").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".preference_file_example").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-macOS-preference-file-[a-z0-9]{8}$`)),
					check.That(resourceType+".preference_file_example").Key("description").HasValue("Configure Safari browser settings via preference file"),
					check.That(resourceType+".preference_file_example").Key("preference_file.file_name").HasValue("com.apple.Safari.plist"),
					check.That(resourceType+".preference_file_example").Key("preference_file.bundle_id").HasValue("com.apple.Safari"),
					check.That(resourceType+".preference_file_example").Key("preference_file.configuration_xml").Exists(),
					check.That(resourceType+".preference_file_example").Key("assignments.#").HasValue("4"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing preference file configuration")
				},
				ResourceName:      resourceType + ".preference_file_example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// Trusted Root Certificate Tests
func TestAccMacosDeviceConfigurationTemplatesResource_TrustedRootCertificate(t *testing.T) {
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
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating trusted root certificate configuration")
				},
				Config: loadAcceptanceTestTerraform("resource_trusted_root_certificate_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("macOS device configuration", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".trusted_cert_example").ExistsInGraph(testResource),
					check.That(resourceType+".trusted_cert_example").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".trusted_cert_example").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-macOS-trusted-root-cert-[a-z0-9]{8}$`)),
					check.That(resourceType+".trusted_cert_example").Key("description").HasValue("Install company root certificate for secure connections"),
					check.That(resourceType+".trusted_cert_example").Key("trusted_certificate.deployment_channel").HasValue("deviceChannel"),
					check.That(resourceType+".trusted_cert_example").Key("trusted_certificate.cert_file_name").HasValue("MicrosoftRootCertificateAuthority2011.cer"),
					check.That(resourceType+".trusted_cert_example").Key("trusted_certificate.trusted_root_certificate").Exists(),
					check.That(resourceType+".trusted_cert_example").Key("assignments.#").HasValue("4"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing trusted root certificate configuration")
				},
				ResourceName:      resourceType + ".trusted_cert_example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// PKCS Certificate Tests
func TestAccMacosDeviceConfigurationTemplatesResource_PkcsCertificate(t *testing.T) {
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
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating PKCS certificate configuration")
				},
				Config: loadAcceptanceTestTerraform("resource_pkcs_certificate_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("macOS device configuration", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".pkcs_cert_example").ExistsInGraph(testResource),
					check.That(resourceType+".pkcs_cert_example").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".pkcs_cert_example").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-macOS-pkcs-cert-[a-z0-9]{8}$`)),
					check.That(resourceType+".pkcs_cert_example").Key("description").HasValue("PKCS certificate profile for user authentication"),
					check.That(resourceType+".pkcs_cert_example").Key("pkcs_certificate.deployment_channel").HasValue("deviceChannel"),
					check.That(resourceType+".pkcs_cert_example").Key("pkcs_certificate.renewal_threshold_percentage").HasValue("20"),
					check.That(resourceType+".pkcs_cert_example").Key("pkcs_certificate.certificate_store").HasValue("machine"),
					check.That(resourceType+".pkcs_cert_example").Key("pkcs_certificate.certificate_validity_period_scale").HasValue("months"),
					check.That(resourceType+".pkcs_cert_example").Key("pkcs_certificate.certificate_validity_period_value").HasValue("12"),
					check.That(resourceType+".pkcs_cert_example").Key("pkcs_certificate.subject_name_format").HasValue("custom"),
					check.That(resourceType+".pkcs_cert_example").Key("pkcs_certificate.subject_name_format_string").HasValue("CN={{AAD_Device_ID}}"),
					check.That(resourceType+".pkcs_cert_example").Key("pkcs_certificate.custom_subject_alternative_names.#").HasValue("6"),
					check.That(resourceType+".pkcs_cert_example").Key("pkcs_certificate.allow_all_apps_access").HasValue("true"),
					check.That(resourceType+".pkcs_cert_example").Key("assignments.#").HasValue("4"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing PKCS certificate configuration")
				},
				ResourceName:      resourceType + ".pkcs_cert_example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// SCEP Certificate Tests
func TestAccMacosDeviceConfigurationTemplatesResource_ScepCertificate(t *testing.T) {
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
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating trusted root certificate and SCEP certificate configuration")
				},
				Config: loadAcceptanceTestTerraform("resource_trusted_root_certificate_maximal.tf") + "\n" + loadAcceptanceTestTerraform("resource_scep_certificate_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("macOS device configuration", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".scep_cert_example").ExistsInGraph(testResource),
					check.That(resourceType+".scep_cert_example").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".scep_cert_example").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-macOS-scep-cert-[a-z0-9]{8}$`)),
					check.That(resourceType+".scep_cert_example").Key("description").HasValue("SCEP certificate profile for device authentication"),
					check.That(resourceType+".scep_cert_example").Key("scep_certificate.deployment_channel").HasValue("deviceChannel"),
					check.That(resourceType+".scep_cert_example").Key("scep_certificate.renewal_threshold_percentage").HasValue("20"),
					check.That(resourceType+".scep_cert_example").Key("scep_certificate.certificate_store").HasValue("machine"),
					check.That(resourceType+".scep_cert_example").Key("scep_certificate.certificate_validity_period_scale").HasValue("years"),
					check.That(resourceType+".scep_cert_example").Key("scep_certificate.certificate_validity_period_value").HasValue("1"),
					check.That(resourceType+".scep_cert_example").Key("scep_certificate.subject_name_format").HasValue("custom"),
					check.That(resourceType+".scep_cert_example").Key("scep_certificate.subject_name_format_string").HasValue("CN={{AAD_Device_ID}}"),
					check.That(resourceType+".scep_cert_example").Key("scep_certificate.key_size").HasValue("size4096"),
					check.That(resourceType+".scep_cert_example").Key("scep_certificate.key_usage.#").HasValue("2"),
					check.That(resourceType+".scep_cert_example").Key("scep_certificate.custom_subject_alternative_names.#").HasValue("4"),
					check.That(resourceType+".scep_cert_example").Key("scep_certificate.extended_key_usages.#").HasValue("4"),
					check.That(resourceType+".scep_cert_example").Key("scep_certificate.scep_server_urls.#").HasValue("2"),
					check.That(resourceType+".scep_cert_example").Key("scep_certificate.allow_all_apps_access").HasValue("true"),
					check.That(resourceType+".scep_cert_example").Key("assignments.#").HasValue("4"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing SCEP certificate configuration")
				},
				ResourceName:            resourceType + ".scep_cert_example",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"scep_certificate.root_certificate_odata_bind"},
			},
		},
	})
}
