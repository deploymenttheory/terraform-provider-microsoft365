package graphBetaNamedLocation_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccNamedLocationResource_IPMinimal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckNamedLocationDestroy,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccConfigIPMinimal(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("microsoft365_graph_beta_identity_and_access_named_location.ip_minimal", "display_name", regexp.MustCompile(`^acc-test-named-location-ip-minimal-[0-9a-fA-F-]+$`)),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_named_location.ip_minimal", "is_trusted", "false"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_named_location.ip_minimal", "ipv4_ranges.#", "1"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_identity_and_access_named_location.ip_minimal", "ipv4_ranges.*", "192.168.1.0/24"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_identity_and_access_named_location.ip_minimal", "id"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_identity_and_access_named_location.ip_minimal", "created_date_time"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_identity_and_access_named_location.ip_minimal", "modified_date_time"),
				),
			},
			{
				ResourceName:                         "microsoft365_graph_beta_identity_and_access_named_location.ip_minimal",
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "id",
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
		CheckDestroy:             testAccCheckNamedLocationDestroy,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccConfigIPMaximal(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("microsoft365_graph_beta_identity_and_access_named_location.ip_maximal", "display_name", regexp.MustCompile(`^acc-test-named-location-ip-maximal-[0-9a-fA-F-]+$`)),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_named_location.ip_maximal", "is_trusted", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_named_location.ip_maximal", "ipv4_ranges.#", "2"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_identity_and_access_named_location.ip_maximal", "ipv4_ranges.*", "192.168.0.0/16"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_identity_and_access_named_location.ip_maximal", "ipv4_ranges.*", "172.16.0.0/12"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_named_location.ip_maximal", "ipv6_ranges.#", "3"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_identity_and_access_named_location.ip_maximal", "ipv6_ranges.*", "2001:db8::/32"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_identity_and_access_named_location.ip_maximal", "ipv6_ranges.*", "fe80::/10"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_identity_and_access_named_location.ip_maximal", "ipv6_ranges.*", "2001:4860:4860::/48"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_identity_and_access_named_location.ip_maximal", "id"),
				),
			},
		},
	})
}

func TestAccNamedLocationResource_IPv6Only(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckNamedLocationDestroy,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccConfigIPv6Only(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("microsoft365_graph_beta_identity_and_access_named_location.ipv6_only", "display_name", regexp.MustCompile(`^acc-test-named-location-ipv6-only-[0-9a-fA-F-]+$`)),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_named_location.ipv6_only", "is_trusted", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_named_location.ipv6_only", "ipv6_ranges.#", "2"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_identity_and_access_named_location.ipv6_only", "ipv6_ranges.*", "2001:db8::/32"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_identity_and_access_named_location.ipv6_only", "ipv6_ranges.*", "fe80::/10"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_identity_and_access_named_location.ipv6_only", "id"),
				),
			},
		},
	})
}

func TestAccNamedLocationResource_CountryClientIP(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckNamedLocationDestroy,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccConfigCountryClientIP(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("microsoft365_graph_beta_identity_and_access_named_location.country_client_ip", "display_name", regexp.MustCompile(`^acc-test-named-location-country-client-ip-[0-9a-fA-F-]+$`)),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_named_location.country_client_ip", "country_lookup_method", "clientIpAddress"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_named_location.country_client_ip", "include_unknown_countries_and_regions", "false"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_named_location.country_client_ip", "countries_and_regions.#", "3"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_identity_and_access_named_location.country_client_ip", "countries_and_regions.*", "US"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_identity_and_access_named_location.country_client_ip", "countries_and_regions.*", "CA"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_identity_and_access_named_location.country_client_ip", "countries_and_regions.*", "GB"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_identity_and_access_named_location.country_client_ip", "id"),
				),
			},
			{
				ResourceName:                         "microsoft365_graph_beta_identity_and_access_named_location.country_client_ip",
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "id",
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
		CheckDestroy:             testAccCheckNamedLocationDestroy,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccConfigCountryAuthenticatorGPS(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("microsoft365_graph_beta_identity_and_access_named_location.country_authenticator_gps", "display_name", regexp.MustCompile(`^acc-test-named-location-country-authenticator-gps-[0-9a-fA-F-]+$`)),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_named_location.country_authenticator_gps", "country_lookup_method", "authenticatorAppGps"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_named_location.country_authenticator_gps", "include_unknown_countries_and_regions", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_named_location.country_authenticator_gps", "countries_and_regions.#", "4"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_identity_and_access_named_location.country_authenticator_gps", "countries_and_regions.*", "AD"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_identity_and_access_named_location.country_authenticator_gps", "countries_and_regions.*", "AO"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_identity_and_access_named_location.country_authenticator_gps", "countries_and_regions.*", "AI"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_identity_and_access_named_location.country_authenticator_gps", "countries_and_regions.*", "AQ"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_identity_and_access_named_location.country_authenticator_gps", "id"),
				),
			},
		},
	})
}

// Configuration helper functions
func testAccConfigIPMinimal() string {
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/named_location_ip_minimal.tf")
	if err != nil {
		panic(fmt.Sprintf("failed to load IP minimal config: %s", err.Error()))
	}
	return acceptance.ConfiguredM365ProviderBlock(accTestConfig)
}

func testAccConfigIPMaximal() string {
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/named_location_ip_maximal.tf")
	if err != nil {
		panic(fmt.Sprintf("failed to load IP maximal config: %s", err.Error()))
	}
	return acceptance.ConfiguredM365ProviderBlock(accTestConfig)
}

func testAccConfigIPv6Only() string {
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/named_location_ipv6_only.tf")
	if err != nil {
		panic(fmt.Sprintf("failed to load IPv6 only config: %s", err.Error()))
	}
	return acceptance.ConfiguredM365ProviderBlock(accTestConfig)
}

func testAccConfigCountryClientIP() string {
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/named_location_country_client_ip.tf")
	if err != nil {
		panic(fmt.Sprintf("failed to load country client IP config: %s", err.Error()))
	}
	return acceptance.ConfiguredM365ProviderBlock(accTestConfig)
}

func testAccConfigCountryAuthenticatorGPS() string {
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/named_location_country_authenticator_gps.tf")
	if err != nil {
		panic(fmt.Sprintf("failed to load country authenticator GPS config: %s", err.Error()))
	}
	return acceptance.ConfiguredM365ProviderBlock(accTestConfig)
}

// Destroy check function
func testAccCheckNamedLocationDestroy(s *terraform.State) error {
	httpClient, err := acceptance.TestHTTPClient()
	if err != nil {
		return fmt.Errorf("error creating HTTP client for CheckDestroy: %v", err)
	}

	ctx := context.Background()
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "microsoft365_graph_beta_identity_and_access_named_location" {
			continue
		}

		resourceID := rs.Primary.ID
		fmt.Printf("DEBUG: Checking destroy status for named location %s\n", resourceID)

		// First, check if the named location still exists
		url := httpClient.GetBaseURL() + "/identity/conditionalAccess/namedLocations/" + resourceID
		httpReq, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			return fmt.Errorf("error creating request for CheckDestroy: %v", err)
		}

		httpResp, err := httpClient.Do(httpReq)
		if err != nil {
			return fmt.Errorf("error making request for CheckDestroy: %v", err)
		}
		defer httpResp.Body.Close()

		// Resource successfully destroyed
		if httpResp.StatusCode == http.StatusNotFound {
			fmt.Printf("DEBUG: Named location %s successfully destroyed (404/NotFound)\n", resourceID)
			continue
		}

		// Resource still exists - attempt cleanup
		if httpResp.StatusCode == http.StatusOK {
			fmt.Printf("DEBUG: Named location %s still exists, attempting cleanup\n", resourceID)

			// Parse the response to check if it's a trusted IP location
			var currentResource map[string]any
			if err := json.NewDecoder(httpResp.Body).Decode(&currentResource); err != nil {
				return fmt.Errorf("error parsing resource for cleanup: %v", err)
			}
			httpResp.Body.Close()

			// Check if this is an IP named location with isTrusted=true
			odataType, _ := currentResource["@odata.type"].(string)
			isTrusted, _ := currentResource["isTrusted"].(bool)

			if odataType == "#microsoft.graph.ipNamedLocation" && isTrusted {
				fmt.Printf("DEBUG: Named location %s is trusted, patching to untrusted before cleanup\n", resourceID)

				// Patch to set isTrusted to false
				patchBody := map[string]any{
					"@odata.type": "#microsoft.graph.ipNamedLocation",
					"isTrusted":   false,
				}

				jsonBytes, err := json.Marshal(patchBody)
				if err != nil {
					return fmt.Errorf("error marshaling patch request for cleanup: %v", err)
				}

				patchReq, err := http.NewRequestWithContext(ctx, "PATCH", url, bytes.NewReader(jsonBytes))
				if err != nil {
					return fmt.Errorf("error creating patch request for cleanup: %v", err)
				}

				patchResp, err := httpClient.Do(patchReq)
				if err != nil {
					return fmt.Errorf("error making patch request for cleanup: %v", err)
				}
				defer patchResp.Body.Close()

				if patchResp.StatusCode != http.StatusNoContent && patchResp.StatusCode != http.StatusOK {
					return fmt.Errorf("error patching trusted location for cleanup: %d %s", patchResp.StatusCode, patchResp.Status)
				}

				fmt.Printf("DEBUG: Successfully patched named location %s to untrusted\n", resourceID)
			}

			// Now attempt to delete the resource
			deleteReq, err := http.NewRequestWithContext(ctx, "DELETE", url, nil)
			if err != nil {
				return fmt.Errorf("error creating delete request for cleanup: %v", err)
			}

			deleteResp, err := httpClient.Do(deleteReq)
			if err != nil {
				return fmt.Errorf("error making delete request for cleanup: %v", err)
			}
			defer deleteResp.Body.Close()

			if deleteResp.StatusCode == http.StatusNoContent || deleteResp.StatusCode == http.StatusNotFound {
				fmt.Printf("DEBUG: Successfully cleaned up named location %s\n", resourceID)
				continue
			}

			return fmt.Errorf("failed to clean up named location %s: %d %s", resourceID, deleteResp.StatusCode, deleteResp.Status)
		}

		return fmt.Errorf("unexpected response checking if named location %s was destroyed: %d %s", resourceID, httpResp.StatusCode, httpResp.Status)
	}

	return nil
}
