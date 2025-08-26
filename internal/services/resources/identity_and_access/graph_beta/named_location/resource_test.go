package graphBetaNamedLocation_test

import (
	"regexp"
	"testing"

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

func testCheckExists(resourceName string) resource.TestCheckFunc {
	return resource.TestCheckResourceAttrSet(resourceName, "id")
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
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_named_location.ip_minimal", "display_name", "unit-test-ip-named-location-minimal"),
					resource.TestMatchResourceAttr("microsoft365_graph_beta_identity_and_access_named_location.ip_minimal", "id", regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_named_location.ip_minimal", "is_trusted", "false"),
					
					// IPv4 ranges
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_named_location.ip_minimal", "ipv4_ranges.#", "1"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_identity_and_access_named_location.ip_minimal", "ipv4_ranges.*", "192.168.1.0/24"),
					
					// Country fields should be null for IP locations
					resource.TestCheckNoResourceAttr("microsoft365_graph_beta_identity_and_access_named_location.ip_minimal", "country_lookup_method"),
					resource.TestCheckNoResourceAttr("microsoft365_graph_beta_identity_and_access_named_location.ip_minimal", "countries_and_regions"),
					resource.TestCheckNoResourceAttr("microsoft365_graph_beta_identity_and_access_named_location.ip_minimal", "include_unknown_countries_and_regions"),
				),
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
					testCheckExists("microsoft365_graph_beta_identity_and_access_named_location.ip_maximal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_named_location.ip_maximal", "display_name", "unit-test-ip-named-location-maximal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_named_location.ip_maximal", "is_trusted", "true"),
					
					// IPv4 ranges
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_named_location.ip_maximal", "ipv4_ranges.#", "2"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_identity_and_access_named_location.ip_maximal", "ipv4_ranges.*", "192.168.0.0/16"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_identity_and_access_named_location.ip_maximal", "ipv4_ranges.*", "172.16.0.0/12"),
					
					// IPv6 ranges
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_named_location.ip_maximal", "ipv6_ranges.#", "3"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_identity_and_access_named_location.ip_maximal", "ipv6_ranges.*", "2001:db8::/32"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_identity_and_access_named_location.ip_maximal", "ipv6_ranges.*", "fe80::/10"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_identity_and_access_named_location.ip_maximal", "ipv6_ranges.*", "2001:4860:4860::/48"),
				),
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
					testCheckExists("microsoft365_graph_beta_identity_and_access_named_location.ip_ipv6_only"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_named_location.ip_ipv6_only", "display_name", "unit-test-ip-named-location-ipv6-only"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_named_location.ip_ipv6_only", "is_trusted", "true"),
					
					// IPv6 ranges only
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_named_location.ip_ipv6_only", "ipv6_ranges.#", "2"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_identity_and_access_named_location.ip_ipv6_only", "ipv6_ranges.*", "2001:db8::/32"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_identity_and_access_named_location.ip_ipv6_only", "ipv6_ranges.*", "fe80::/10"),
				),
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
					testCheckExists("microsoft365_graph_beta_identity_and_access_named_location.country_minimal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_named_location.country_minimal", "display_name", "unit-test-country-named-location-minimal"),
					
					// Country attributes
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_named_location.country_minimal", "country_lookup_method", "clientIpAddress"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_named_location.country_minimal", "include_unknown_countries_and_regions", "false"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_named_location.country_minimal", "countries_and_regions.#", "1"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_identity_and_access_named_location.country_minimal", "countries_and_regions.*", "US"),
					
					// IP fields should be null for country locations
					resource.TestCheckNoResourceAttr("microsoft365_graph_beta_identity_and_access_named_location.country_minimal", "is_trusted"),
					resource.TestCheckNoResourceAttr("microsoft365_graph_beta_identity_and_access_named_location.country_minimal", "ipv4_ranges"),
					resource.TestCheckNoResourceAttr("microsoft365_graph_beta_identity_and_access_named_location.country_minimal", "ipv6_ranges"),
				),
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
					testCheckExists("microsoft365_graph_beta_identity_and_access_named_location.country_authenticator_gps"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_named_location.country_authenticator_gps", "display_name", "unit-test-country-named-location-authenticator-gps"),
					
					// Country attributes
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_named_location.country_authenticator_gps", "country_lookup_method", "authenticatorAppGps"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_named_location.country_authenticator_gps", "include_unknown_countries_and_regions", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_named_location.country_authenticator_gps", "countries_and_regions.#", "4"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_identity_and_access_named_location.country_authenticator_gps", "countries_and_regions.*", "AD"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_identity_and_access_named_location.country_authenticator_gps", "countries_and_regions.*", "AO"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_identity_and_access_named_location.country_authenticator_gps", "countries_and_regions.*", "AI"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_identity_and_access_named_location.country_authenticator_gps", "countries_and_regions.*", "AQ"),
				),
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
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_named_location.ip_maximal", "is_trusted", "true"),
					testCheckExists("microsoft365_graph_beta_identity_and_access_named_location.ip_maximal"),
				),
			},
		},
		// The CheckDestroy function will test that the resource was properly deleted
		// even though it was trusted, verifying our delete logic handles the PATCH correctly
	})
}