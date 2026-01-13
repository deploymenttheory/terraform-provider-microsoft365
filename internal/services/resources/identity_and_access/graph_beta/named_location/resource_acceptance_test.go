package graphBetaNamedLocation_test

import (
	"regexp"
	"testing"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/destroy"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/testlog"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaNamedLocation "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_beta/named_location"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

var (
	// Resource type name from the resource package
	resourceType = graphBetaNamedLocation.ResourceName

	// testResource is the test resource implementation for named locations
	testResource = graphBetaNamedLocation.NamedLocationTestResource{}
)

// Helper function to load test configs from acceptance directory
func loadAcceptanceTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/acceptance/" + filename)
	if err != nil {
		panic("failed to load acceptance config " + filename + ": " + err.Error())
	}
	return acceptance.ConfiguredM365ProviderBlock(config)
}

func TestAccNamedLocationResource_IPMinimal(t *testing.T) {
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
					testlog.StepAction(resourceType, "Creating IP minimal named location")
				},
				Config: loadAcceptanceTestTerraform("named_location_ip_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("named location", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".ip_minimal").ExistsInGraph(testResource),
					check.That(resourceType+".ip_minimal").Key("id").Exists(),
					check.That(resourceType+".ip_minimal").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-named-location-ip-minimal-[0-9a-fA-F-]+$`)),
					check.That(resourceType+".ip_minimal").Key("is_trusted").HasValue("false"),
					check.That(resourceType+".ip_minimal").Key("ipv4_ranges.#").HasValue("1"),
					check.That(resourceType+".ip_minimal").Key("ipv4_ranges.*").ContainsTypeSetElement("192.168.1.0/24"),
					check.That(resourceType+".ip_minimal").Key("created_date_time").Exists(),
					check.That(resourceType+".ip_minimal").Key("modified_date_time").Exists(),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing IP minimal named location")
				},
				ResourceName:      resourceType + ".ip_minimal",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"timeouts",
				},
			},
		},
	})
}

func TestAccNamedLocationResource_IPMaximal(t *testing.T) {
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
					testlog.StepAction(resourceType, "Creating IP maximal named location")
				},
				Config: loadAcceptanceTestTerraform("named_location_ip_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("named location", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".ip_maximal").ExistsInGraph(testResource),
					check.That(resourceType+".ip_maximal").Key("id").Exists(),
					check.That(resourceType+".ip_maximal").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-named-location-ip-maximal-[0-9a-fA-F-]+$`)),
					check.That(resourceType+".ip_maximal").Key("is_trusted").HasValue("true"),
					check.That(resourceType+".ip_maximal").Key("ipv4_ranges.#").HasValue("2"),
					check.That(resourceType+".ip_maximal").Key("ipv4_ranges.*").ContainsTypeSetElement("192.168.0.0/16"),
					check.That(resourceType+".ip_maximal").Key("ipv4_ranges.*").ContainsTypeSetElement("172.16.0.0/12"),
					check.That(resourceType+".ip_maximal").Key("ipv6_ranges.#").HasValue("3"),
					check.That(resourceType+".ip_maximal").Key("ipv6_ranges.*").ContainsTypeSetElement("2001:db8::/32"),
					check.That(resourceType+".ip_maximal").Key("ipv6_ranges.*").ContainsTypeSetElement("fe80::/10"),
					check.That(resourceType+".ip_maximal").Key("ipv6_ranges.*").ContainsTypeSetElement("2001:4860:4860::/48"),
				),
			},
		},
	})
}

func TestAccNamedLocationResource_IPv6Only(t *testing.T) {
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
					testlog.StepAction(resourceType, "Creating IPv6 only named location")
				},
				Config: loadAcceptanceTestTerraform("named_location_ipv6_only.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("named location", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".ipv6_only").ExistsInGraph(testResource),
					check.That(resourceType+".ipv6_only").Key("id").Exists(),
					check.That(resourceType+".ipv6_only").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-named-location-ipv6-only-[0-9a-fA-F-]+$`)),
					check.That(resourceType+".ipv6_only").Key("is_trusted").HasValue("true"),
					check.That(resourceType+".ipv6_only").Key("ipv6_ranges.#").HasValue("2"),
					check.That(resourceType+".ipv6_only").Key("ipv6_ranges.*").ContainsTypeSetElement("2001:db8::/32"),
					check.That(resourceType+".ipv6_only").Key("ipv6_ranges.*").ContainsTypeSetElement("fe80::/10"),
				),
			},
		},
	})
}

func TestAccNamedLocationResource_CountryClientIP(t *testing.T) {
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
					testlog.StepAction(resourceType, "Creating country client IP named location")
				},
				Config: loadAcceptanceTestTerraform("named_location_country_client_ip.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("named location", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".country_client_ip").ExistsInGraph(testResource),
					check.That(resourceType+".country_client_ip").Key("id").Exists(),
					check.That(resourceType+".country_client_ip").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-named-location-country-client-ip-[0-9a-fA-F-]+$`)),
					check.That(resourceType+".country_client_ip").Key("country_lookup_method").HasValue("clientIpAddress"),
					check.That(resourceType+".country_client_ip").Key("include_unknown_countries_and_regions").HasValue("false"),
					check.That(resourceType+".country_client_ip").Key("countries_and_regions.#").HasValue("3"),
					check.That(resourceType+".country_client_ip").Key("countries_and_regions.*").ContainsTypeSetElement("US"),
					check.That(resourceType+".country_client_ip").Key("countries_and_regions.*").ContainsTypeSetElement("CA"),
					check.That(resourceType+".country_client_ip").Key("countries_and_regions.*").ContainsTypeSetElement("GB"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing country client IP named location")
				},
				ResourceName:      resourceType + ".country_client_ip",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"timeouts",
				},
			},
		},
	})
}

func TestAccNamedLocationResource_CountryAuthenticatorGPS(t *testing.T) {
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
					testlog.StepAction(resourceType, "Creating country authenticator GPS named location")
				},
				Config: loadAcceptanceTestTerraform("named_location_country_authenticator_gps.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("named location", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".country_authenticator_gps").ExistsInGraph(testResource),
					check.That(resourceType+".country_authenticator_gps").Key("id").Exists(),
					check.That(resourceType+".country_authenticator_gps").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-named-location-country-authenticator-gps-[0-9a-fA-F-]+$`)),
					check.That(resourceType+".country_authenticator_gps").Key("country_lookup_method").HasValue("authenticatorAppGps"),
					check.That(resourceType+".country_authenticator_gps").Key("include_unknown_countries_and_regions").HasValue("true"),
					check.That(resourceType+".country_authenticator_gps").Key("countries_and_regions.#").HasValue("4"),
					check.That(resourceType+".country_authenticator_gps").Key("countries_and_regions.*").ContainsTypeSetElement("AD"),
					check.That(resourceType+".country_authenticator_gps").Key("countries_and_regions.*").ContainsTypeSetElement("AO"),
					check.That(resourceType+".country_authenticator_gps").Key("countries_and_regions.*").ContainsTypeSetElement("AI"),
					check.That(resourceType+".country_authenticator_gps").Key("countries_and_regions.*").ContainsTypeSetElement("AQ"),
				),
			},
		},
	})
}
