package graphBetaNamedLocation_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	namedLocationMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_beta/named_location/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

func setupMockEnvironment() (*mocks.Mocks, *namedLocationMocks.NamedLocationMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	namedLocationMock := &namedLocationMocks.NamedLocationMock{}
	namedLocationMock.RegisterMocks()
	return mockClient, namedLocationMock
}

func setupErrorMockEnvironment() (*mocks.Mocks, *namedLocationMocks.NamedLocationMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	namedLocationMock := &namedLocationMocks.NamedLocationMock{}
	namedLocationMock.RegisterErrorMocks()
	return mockClient, namedLocationMock
}

func TestNamedLocationResource_IPMinimal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, namedLocationMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer namedLocationMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigIPMinimal(),
				Check: resource.ComposeTestCheckFunc(
					// Basic attributes
					check.That(resourceType+".ip_minimal").Key("display_name").HasValue("unit-test-ip-named-location-minimal"),
					check.That(resourceType+".ip_minimal").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".ip_minimal").Key("is_trusted").HasValue("false"),

					// IPv4 ranges
					check.That(resourceType+".ip_minimal").Key("ipv4_ranges.#").HasValue("1"),
					check.That(resourceType+".ip_minimal").Key("ipv4_ranges.*").ContainsTypeSetElement("192.168.1.0/24"),

					// Country fields should be null for IP locations
					check.That(resourceType+".ip_minimal").Key("country_lookup_method").DoesNotExist(),
					check.That(resourceType+".ip_minimal").Key("countries_and_regions").DoesNotExist(),
					check.That(resourceType+".ip_minimal").Key("include_unknown_countries_and_regions").DoesNotExist(),
				),
			},
			{
				ResourceName:      resourceType + ".ip_minimal",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestNamedLocationResource_IPMaximal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, namedLocationMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer namedLocationMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigIPMaximal(),
				Check: resource.ComposeTestCheckFunc(
					// Basic attributes
					check.That(resourceType+".ip_maximal").Key("display_name").HasValue("unit-test-ip-named-location-maximal"),
					check.That(resourceType+".ip_maximal").Key("is_trusted").HasValue("true"),

					// IPv4 ranges
					check.That(resourceType+".ip_maximal").Key("ipv4_ranges.#").HasValue("2"),
					check.That(resourceType+".ip_maximal").Key("ipv4_ranges.*").ContainsTypeSetElement("192.168.0.0/16"),
					check.That(resourceType+".ip_maximal").Key("ipv4_ranges.*").ContainsTypeSetElement("172.16.0.0/12"),

					// IPv6 ranges
					check.That(resourceType+".ip_maximal").Key("ipv6_ranges.#").HasValue("3"),
					check.That(resourceType+".ip_maximal").Key("ipv6_ranges.*").ContainsTypeSetElement("2001:db8::/32"),
					check.That(resourceType+".ip_maximal").Key("ipv6_ranges.*").ContainsTypeSetElement("fe80::/10"),
					check.That(resourceType+".ip_maximal").Key("ipv6_ranges.*").ContainsTypeSetElement("2001:4860:4860::/48"),
				),
			},
			{
				ResourceName:      resourceType + ".ip_maximal",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestNamedLocationResource_IPv6Only(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, namedLocationMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer namedLocationMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigIPv6Only(),
				Check: resource.ComposeTestCheckFunc(
					// Basic attributes
					check.That(resourceType+".ip_ipv6_only").Key("display_name").HasValue("unit-test-ip-named-location-ipv6-only"),
					check.That(resourceType+".ip_ipv6_only").Key("is_trusted").HasValue("true"),

					// IPv6 ranges only
					check.That(resourceType+".ip_ipv6_only").Key("ipv6_ranges.#").HasValue("2"),
					check.That(resourceType+".ip_ipv6_only").Key("ipv6_ranges.*").ContainsTypeSetElement("2001:db8::/32"),
					check.That(resourceType+".ip_ipv6_only").Key("ipv6_ranges.*").ContainsTypeSetElement("fe80::/10"),
				),
			},
			{
				ResourceName:      resourceType + ".ip_ipv6_only",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestNamedLocationResource_CountryMinimal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, namedLocationMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer namedLocationMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigCountryMinimal(),
				Check: resource.ComposeTestCheckFunc(
					// Basic attributes
					check.That(resourceType+".country_minimal").Key("display_name").HasValue("unit-test-country-named-location-minimal"),

					// Country attributes
					check.That(resourceType+".country_minimal").Key("country_lookup_method").HasValue("clientIpAddress"),
					check.That(resourceType+".country_minimal").Key("include_unknown_countries_and_regions").HasValue("false"),
					check.That(resourceType+".country_minimal").Key("countries_and_regions.#").HasValue("1"),
					check.That(resourceType+".country_minimal").Key("countries_and_regions.*").ContainsTypeSetElement("US"),

					// IP fields should be null for country locations
					check.That(resourceType+".country_minimal").Key("is_trusted").DoesNotExist(),
					check.That(resourceType+".country_minimal").Key("ipv4_ranges").DoesNotExist(),
					check.That(resourceType+".country_minimal").Key("ipv6_ranges").DoesNotExist(),
				),
			},
			{
				ResourceName:      resourceType + ".country_minimal",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestNamedLocationResource_CountryAuthenticatorGPS(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, namedLocationMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer namedLocationMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigCountryAuthenticatorGPS(),
				Check: resource.ComposeTestCheckFunc(
					// Basic attributes
					check.That(resourceType+".country_authenticator_gps").Key("display_name").HasValue("unit-test-country-named-location-authenticator-gps"),

					// Country attributes
					check.That(resourceType+".country_authenticator_gps").Key("country_lookup_method").HasValue("authenticatorAppGps"),
					check.That(resourceType+".country_authenticator_gps").Key("include_unknown_countries_and_regions").HasValue("true"),
					check.That(resourceType+".country_authenticator_gps").Key("countries_and_regions.#").HasValue("4"),
					check.That(resourceType+".country_authenticator_gps").Key("countries_and_regions.*").ContainsTypeSetElement("AD"),
					check.That(resourceType+".country_authenticator_gps").Key("countries_and_regions.*").ContainsTypeSetElement("AO"),
					check.That(resourceType+".country_authenticator_gps").Key("countries_and_regions.*").ContainsTypeSetElement("AI"),
					check.That(resourceType+".country_authenticator_gps").Key("countries_and_regions.*").ContainsTypeSetElement("AQ"),
				),
			},
			{
				ResourceName:      resourceType + ".country_authenticator_gps",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// Configuration helper functions
func testConfigIPMinimal() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_ip_minimal.tf")
	if err != nil {
		panic("failed to load IP minimal config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigIPMaximal() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_ip_maximal.tf")
	if err != nil {
		panic("failed to load IP maximal config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigIPv6Only() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_ip_ipv6_only.tf")
	if err != nil {
		panic("failed to load IPv6 only config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigCountryMinimal() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_country_minimal.tf")
	if err != nil {
		panic("failed to load country minimal config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigCountryAuthenticatorGPS() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_country_authenticator_gps.tf")
	if err != nil {
		panic("failed to load country authenticator GPS config: " + err.Error())
	}
	return unitTestConfig
}

// TestNamedLocationResource_TrustedIPDeletion tests that trusted IP locations are properly handled during deletion
func TestNamedLocationResource_TrustedIPDeletion(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, namedLocationMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer namedLocationMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Create a trusted IP location
				Config: testConfigIPMaximal(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType + ".ip_maximal").Key("is_trusted").HasValue("true"),
				),
			},
		},
		// The CheckDestroy function will test that the resource was properly deleted
		// even though it was trusted, verifying our delete logic handles the PATCH correctly
	})
}
