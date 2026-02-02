package graphBetaAgentIdentityBlueprintCertificateCredential_test

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
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func loadAccTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/acceptance/" + filename)
	if err != nil {
		panic("failed to load acceptance test config " + filename + ": " + err.Error())
	}
	return config
}

func TestAccResourceAgentIdentityBlueprintCertificateCredential_01_PEM(t *testing.T) {
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
			"tls": {
				Source:            "hashicorp/tls",
				VersionConstraint: constants.ExternalProviderTLSVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Step 1: Creating certificate credential with PEM encoding")
				},
				Config: loadAccTestTerraform("resource_pem.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(s *terraform.State) error {
						testlog.WaitForConsistency("certificate credential", 10*time.Second)
						time.Sleep(10 * time.Second)
						return nil
					},
					check.That(resourceType+".test_pem").Key("key_id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_pem").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-certificate-pem-`)),
					check.That(resourceType+".test_pem").Key("encoding").HasValue("pem"),
					check.That(resourceType+".test_pem").Key("type").HasValue("AsymmetricX509Cert"),
					check.That(resourceType+".test_pem").Key("usage").HasValue("Verify"),
				),
			},
			// Note: Import is not supported for certificate credentials
		},
	})
}

func TestAccResourceAgentIdentityBlueprintCertificateCredential_02_DER(t *testing.T) {
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
			"tls": {
				Source:            "hashicorp/tls",
				VersionConstraint: constants.ExternalProviderTLSVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Step 1: Creating certificate credential with DER encoding (base64)")
				},
				Config: loadAccTestTerraform("resource_der.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(s *terraform.State) error {
						testlog.WaitForConsistency("certificate credential", 10*time.Second)
						time.Sleep(10 * time.Second)
						return nil
					},
					check.That(resourceType+".test_der").Key("key_id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_der").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-certificate-der-`)),
					check.That(resourceType+".test_der").Key("encoding").HasValue("base64"),
					check.That(resourceType+".test_der").Key("type").HasValue("AsymmetricX509Cert"),
					check.That(resourceType+".test_der").Key("usage").HasValue("Verify"),
				),
			},
		},
	})
}

func TestAccResourceAgentIdentityBlueprintCertificateCredential_03_HEX(t *testing.T) {
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
			"tls": {
				Source:            "hashicorp/tls",
				VersionConstraint: constants.ExternalProviderTLSVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Step 1: Creating certificate credential with HEX encoding")
				},
				Config: loadAccTestTerraform("resource_hex.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(s *terraform.State) error {
						testlog.WaitForConsistency("certificate credential", 10*time.Second)
						time.Sleep(10 * time.Second)
						return nil
					},
					check.That(resourceType+".test_hex").Key("key_id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_hex").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-certificate-hex-`)),
					check.That(resourceType+".test_hex").Key("encoding").HasValue("hex"),
					check.That(resourceType+".test_hex").Key("type").HasValue("AsymmetricX509Cert"),
					check.That(resourceType+".test_hex").Key("usage").HasValue("Verify"),
				),
			},
		},
	})
}
