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
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaNamedLocation "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_beta/named_location"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const (
	testResourceName = "microsoft365_graph_beta_identity_and_access_named_location"
)

var (
	// Resource type name constructed from exported constants
	resourceType = constants.PROVIDER_NAME + "_" + graphBetaNamedLocation.ResourceName

	// testResource is the test resource implementation for named locations
	testResource = graphBetaNamedLocation.NamedLocationTestResource{}
)

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
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating IP minimal named location")
				},
				Config: testAccConfigIPMinimal(),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("named location", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(testResourceName+".ip_minimal").ExistsInGraph(testResource),
					check.That(testResourceName+".ip_minimal").Key("id").Exists(),
					check.That(testResourceName+".ip_minimal").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-named-location-ip-minimal-[0-9a-fA-F-]+$`)),
					check.That(testResourceName+".ip_minimal").Key("is_trusted").HasValue("false"),
					check.That(testResourceName+".ip_minimal").Key("ipv4_ranges.#").HasValue("1"),
					check.That(testResourceName+".ip_minimal").Key("ipv4_ranges.*").ContainsTypeSetElement("192.168.1.0/24"),
					check.That(testResourceName+".ip_minimal").Key("created_date_time").Exists(),
					check.That(testResourceName+".ip_minimal").Key("modified_date_time").Exists(),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing IP minimal named location")
				},
				ResourceName:      testResourceName + ".ip_minimal",
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
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating IP maximal named location")
				},
				Config: testAccConfigIPMaximal(),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("named location", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(testResourceName+".ip_maximal").ExistsInGraph(testResource),
					check.That(testResourceName+".ip_maximal").Key("id").Exists(),
					check.That(testResourceName+".ip_maximal").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-named-location-ip-maximal-[0-9a-fA-F-]+$`)),
					check.That(testResourceName+".ip_maximal").Key("is_trusted").HasValue("true"),
					check.That(testResourceName+".ip_maximal").Key("ipv4_ranges.#").HasValue("2"),
					check.That(testResourceName+".ip_maximal").Key("ipv4_ranges.*").ContainsTypeSetElement("192.168.0.0/16"),
					check.That(testResourceName+".ip_maximal").Key("ipv4_ranges.*").ContainsTypeSetElement("172.16.0.0/12"),
					check.That(testResourceName+".ip_maximal").Key("ipv6_ranges.#").HasValue("3"),
					check.That(testResourceName+".ip_maximal").Key("ipv6_ranges.*").ContainsTypeSetElement("2001:db8::/32"),
					check.That(testResourceName+".ip_maximal").Key("ipv6_ranges.*").ContainsTypeSetElement("fe80::/10"),
					check.That(testResourceName+".ip_maximal").Key("ipv6_ranges.*").ContainsTypeSetElement("2001:4860:4860::/48"),
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
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating IPv6 only named location")
				},
				Config: testAccConfigIPv6Only(),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("named location", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(testResourceName+".ipv6_only").ExistsInGraph(testResource),
					check.That(testResourceName+".ipv6_only").Key("id").Exists(),
					check.That(testResourceName+".ipv6_only").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-named-location-ipv6-only-[0-9a-fA-F-]+$`)),
					check.That(testResourceName+".ipv6_only").Key("is_trusted").HasValue("true"),
					check.That(testResourceName+".ipv6_only").Key("ipv6_ranges.#").HasValue("2"),
					check.That(testResourceName+".ipv6_only").Key("ipv6_ranges.*").ContainsTypeSetElement("2001:db8::/32"),
					check.That(testResourceName+".ipv6_only").Key("ipv6_ranges.*").ContainsTypeSetElement("fe80::/10"),
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
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating country client IP named location")
				},
				Config: testAccConfigCountryClientIP(),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("named location", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(testResourceName+".country_client_ip").ExistsInGraph(testResource),
					check.That(testResourceName+".country_client_ip").Key("id").Exists(),
					check.That(testResourceName+".country_client_ip").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-named-location-country-client-ip-[0-9a-fA-F-]+$`)),
					check.That(testResourceName+".country_client_ip").Key("country_lookup_method").HasValue("clientIpAddress"),
					check.That(testResourceName+".country_client_ip").Key("include_unknown_countries_and_regions").HasValue("false"),
					check.That(testResourceName+".country_client_ip").Key("countries_and_regions.#").HasValue("3"),
					check.That(testResourceName+".country_client_ip").Key("countries_and_regions.*").ContainsTypeSetElement("US"),
					check.That(testResourceName+".country_client_ip").Key("countries_and_regions.*").ContainsTypeSetElement("CA"),
					check.That(testResourceName+".country_client_ip").Key("countries_and_regions.*").ContainsTypeSetElement("GB"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing country client IP named location")
				},
				ResourceName:      testResourceName + ".country_client_ip",
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
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating country authenticator GPS named location")
				},
				Config: testAccConfigCountryAuthenticatorGPS(),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("named location", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(testResourceName+".country_authenticator_gps").ExistsInGraph(testResource),
					check.That(testResourceName+".country_authenticator_gps").Key("id").Exists(),
					check.That(testResourceName+".country_authenticator_gps").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-named-location-country-authenticator-gps-[0-9a-fA-F-]+$`)),
					check.That(testResourceName+".country_authenticator_gps").Key("country_lookup_method").HasValue("authenticatorAppGps"),
					check.That(testResourceName+".country_authenticator_gps").Key("include_unknown_countries_and_regions").HasValue("true"),
					check.That(testResourceName+".country_authenticator_gps").Key("countries_and_regions.#").HasValue("4"),
					check.That(testResourceName+".country_authenticator_gps").Key("countries_and_regions.*").ContainsTypeSetElement("AD"),
					check.That(testResourceName+".country_authenticator_gps").Key("countries_and_regions.*").ContainsTypeSetElement("AO"),
					check.That(testResourceName+".country_authenticator_gps").Key("countries_and_regions.*").ContainsTypeSetElement("AI"),
					check.That(testResourceName+".country_authenticator_gps").Key("countries_and_regions.*").ContainsTypeSetElement("AQ"),
				),
			},
		},
	})
}

// Test configuration functions
func testAccConfigIPMinimal() string {
	config := mocks.LoadTerraformConfigFile("named_location_ip_minimal.tf")
	return acceptance.ConfiguredM365ProviderBlock(config)
}

func testAccConfigIPMaximal() string {
	config := mocks.LoadTerraformConfigFile("named_location_ip_maximal.tf")
	return acceptance.ConfiguredM365ProviderBlock(config)
}

func testAccConfigIPv6Only() string {
	config := mocks.LoadTerraformConfigFile("named_location_ipv6_only.tf")
	return acceptance.ConfiguredM365ProviderBlock(config)
}

func testAccConfigCountryClientIP() string {
	config := mocks.LoadTerraformConfigFile("named_location_country_client_ip.tf")
	return acceptance.ConfiguredM365ProviderBlock(config)
}

func testAccConfigCountryAuthenticatorGPS() string {
	config := mocks.LoadTerraformConfigFile("named_location_country_authenticator_gps.tf")
	return acceptance.ConfiguredM365ProviderBlock(config)
}
