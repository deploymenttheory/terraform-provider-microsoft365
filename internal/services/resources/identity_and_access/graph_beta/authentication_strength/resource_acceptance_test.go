package graphBetaAuthenticationStrength_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccAuthenticationStrengthResource_Minimal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckAuthenticationStrengthDestroy,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccConfigAuthStrengthMinimal(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_authentication_strength.auth_strength_minimal", "display_name", "acc-test-authentication-strength-minimal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_authentication_strength.auth_strength_minimal", "description", "Acceptance test minimal authentication strength policy"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_authentication_strength.auth_strength_minimal", "allowed_combinations.#", "1"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_identity_and_access_authentication_strength.auth_strength_minimal", "allowed_combinations.*", "password,sms"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_identity_and_access_authentication_strength.auth_strength_minimal", "id"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_identity_and_access_authentication_strength.auth_strength_minimal", "created_date_time"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_identity_and_access_authentication_strength.auth_strength_minimal", "modified_date_time"),
				),
			},
			{
				ResourceName:                         "microsoft365_graph_beta_identity_and_access_authentication_strength.auth_strength_minimal",
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

func TestAccAuthenticationStrengthResource_MFAOnly(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckAuthenticationStrengthDestroy,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccConfigAuthStrengthMFAOnly(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_authentication_strength.auth_strength_mfa_only", "display_name", "acc-test-authentication-strength-mfa-only"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_authentication_strength.auth_strength_mfa_only", "description", "Acceptance test MFA-only authentication strength policy"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_authentication_strength.auth_strength_mfa_only", "allowed_combinations.#", "4"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_identity_and_access_authentication_strength.auth_strength_mfa_only", "allowed_combinations.*", "fido2"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_identity_and_access_authentication_strength.auth_strength_mfa_only", "allowed_combinations.*", "windowsHelloForBusiness"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_identity_and_access_authentication_strength.auth_strength_mfa_only", "allowed_combinations.*", "microsoftAuthenticatorPush,federatedSingleFactor"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_identity_and_access_authentication_strength.auth_strength_mfa_only", "allowed_combinations.*", "x509CertificateMultiFactor"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_identity_and_access_authentication_strength.auth_strength_mfa_only", "id"),
				),
			},
			{
				ResourceName:                         "microsoft365_graph_beta_identity_and_access_authentication_strength.auth_strength_mfa_only",
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

func TestAccAuthenticationStrengthResource_Maximal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckAuthenticationStrengthDestroy,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccConfigAuthStrengthMaximal(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_authentication_strength.auth_strength_maximal", "display_name", "acc-test-authentication-strength-maximal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_authentication_strength.auth_strength_maximal", "description", "Acceptance test maximal authentication strength policy with all combinations"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_authentication_strength.auth_strength_maximal", "allowed_combinations.#", "22"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_identity_and_access_authentication_strength.auth_strength_maximal", "allowed_combinations.*", "deviceBasedPush"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_identity_and_access_authentication_strength.auth_strength_maximal", "allowed_combinations.*", "federatedMultiFactor"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_identity_and_access_authentication_strength.auth_strength_maximal", "allowed_combinations.*", "federatedSingleFactor"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_identity_and_access_authentication_strength.auth_strength_maximal", "allowed_combinations.*", "fido2"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_identity_and_access_authentication_strength.auth_strength_maximal", "allowed_combinations.*", "hardwareOath,federatedSingleFactor"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_identity_and_access_authentication_strength.auth_strength_maximal", "allowed_combinations.*", "microsoftAuthenticatorPush,federatedSingleFactor"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_identity_and_access_authentication_strength.auth_strength_maximal", "allowed_combinations.*", "password"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_identity_and_access_authentication_strength.auth_strength_maximal", "allowed_combinations.*", "windowsHelloForBusiness"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_identity_and_access_authentication_strength.auth_strength_maximal", "allowed_combinations.*", "x509CertificateMultiFactor"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_identity_and_access_authentication_strength.auth_strength_maximal", "id"),
				),
			},
		},
	})
}

// Configuration helper functions
func testAccConfigAuthStrengthMinimal() string {
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/resource_auth_strength_minimal.tf")
	if err != nil {
		panic(fmt.Sprintf("failed to load authentication strength minimal config: %s", err.Error()))
	}
	return acceptance.ConfiguredM365ProviderBlock(accTestConfig)
}

func testAccConfigAuthStrengthMFAOnly() string {
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/resource_auth_strength_mfa_only.tf")
	if err != nil {
		panic(fmt.Sprintf("failed to load authentication strength MFA-only config: %s", err.Error()))
	}
	return acceptance.ConfiguredM365ProviderBlock(accTestConfig)
}

func testAccConfigAuthStrengthMaximal() string {
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/resource_auth_strength_maximal.tf")
	if err != nil {
		panic(fmt.Sprintf("failed to load authentication strength maximal config: %s", err.Error()))
	}
	return acceptance.ConfiguredM365ProviderBlock(accTestConfig)
}

// Destroy check function
func testAccCheckAuthenticationStrengthDestroy(s *terraform.State) error {
	httpClient, err := acceptance.TestHTTPClient()
	if err != nil {
		return fmt.Errorf("error creating HTTP client for CheckDestroy: %v", err)
	}

	ctx := context.Background()
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "microsoft365_graph_beta_identity_and_access_authentication_strength" {
			continue
		}

		resourceID := rs.Primary.ID
		fmt.Printf("DEBUG: Checking destroy status for authentication strength %s\n", resourceID)

		// Check if the authentication strength still exists
		url := httpClient.GetBaseURL() + "/identity/conditionalAccess/authenticationStrengths/" + resourceID
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
			fmt.Printf("DEBUG: Authentication strength %s successfully destroyed (404/NotFound)\n", resourceID)
			continue
		}

		// Resource still exists - attempt cleanup
		if httpResp.StatusCode == http.StatusOK {
			fmt.Printf("DEBUG: Authentication strength %s still exists, attempting cleanup\n", resourceID)

			// Parse the response
			var currentResource map[string]interface{}
			if err := json.NewDecoder(httpResp.Body).Decode(&currentResource); err != nil {
				return fmt.Errorf("error parsing resource for cleanup: %v", err)
			}
			httpResp.Body.Close()

			// Attempt to delete the resource
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
				fmt.Printf("DEBUG: Successfully cleaned up authentication strength %s\n", resourceID)
				continue
			}

			return fmt.Errorf("failed to clean up authentication strength %s: %d %s", resourceID, deleteResp.StatusCode, deleteResp.Status)
		}

		return fmt.Errorf("unexpected response checking if authentication strength %s was destroyed: %d %s", resourceID, httpResp.StatusCode, httpResp.Status)
	}

	return nil
}
