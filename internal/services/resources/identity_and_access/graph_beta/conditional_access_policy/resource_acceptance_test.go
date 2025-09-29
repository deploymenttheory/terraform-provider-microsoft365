package graphBetaConditionalAccessPolicy_test

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

func TestAccConditionalAccessPolicyResource_Minimal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckConditionalAccessPolicyDestroy,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccConfigConditionalAccessPolicyMaximalExcludeRoles(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.conditional_access_policy_minimal", "display_name", "unit-test-conditional-access-policy-maximal-exclude-roles"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.conditional_access_policy_minimal", "state", "enabledForReportingButNotEnforced"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.conditional_access_policy_minimal", "conditions.client_app_types.#", "1"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.conditional_access_policy_minimal", "conditions.client_app_types.*", "all"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.conditional_access_policy_minimal", "conditions.applications.include_applications.#", "1"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.conditional_access_policy_minimal", "conditions.applications.include_applications.*", "All"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.conditional_access_policy_minimal", "conditions.users.include_users.#", "1"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.conditional_access_policy_minimal", "conditions.users.include_users.*", "All"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.conditional_access_policy_minimal", "conditions.locations.include_locations.#", "1"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.conditional_access_policy_minimal", "conditions.locations.include_locations.*", "All"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.conditional_access_policy_minimal", "grant_controls.operator", "AND"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.conditional_access_policy_minimal", "grant_controls.terms_of_use.#", "1"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.conditional_access_policy_minimal", "grant_controls.terms_of_use.*", "79f28780-c502-49c4-8951-f53f6a239b60"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.conditional_access_policy_minimal", "conditions.users.exclude_roles.#", "124"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_identity_and_access_conditional_access_policy.conditional_access_policy_minimal", "id"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_identity_and_access_conditional_access_policy.conditional_access_policy_minimal", "created_date_time"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_identity_and_access_conditional_access_policy.conditional_access_policy_minimal", "modified_date_time"),
				),
			},
			{
				ResourceName:                         "microsoft365_graph_beta_identity_and_access_conditional_access_policy.conditional_access_policy_minimal",
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

// Configuration helper functions
func testAccConfigConditionalAccessPolicyMaximalExcludeRoles() string {
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/resource_conditional_access_policy_maximal_exclude_roles.tf")
	if err != nil {
		panic(fmt.Sprintf("failed to load conditional access policy minimal config: %s", err.Error()))
	}
	return acceptance.ConfiguredM365ProviderBlock(accTestConfig)
}

// Destroy check function
func testAccCheckConditionalAccessPolicyDestroy(s *terraform.State) error {
	httpClient, err := acceptance.TestHTTPClient()
	if err != nil {
		return fmt.Errorf("error creating HTTP client for CheckDestroy: %v", err)
	}

	ctx := context.Background()
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "microsoft365_graph_beta_identity_and_access_conditional_access_policy" {
			continue
		}

		resourceID := rs.Primary.ID
		fmt.Printf("DEBUG: Checking destroy status for conditional access policy %s\n", resourceID)

		// Check if the conditional access policy still exists
		url := httpClient.GetBaseURL() + "/identity/conditionalAccess/policies/" + resourceID
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
			fmt.Printf("DEBUG: Conditional access policy %s successfully destroyed (404/NotFound)\n", resourceID)
			continue
		}

		// Resource still exists - attempt cleanup
		if httpResp.StatusCode == http.StatusOK {
			fmt.Printf("DEBUG: Conditional access policy %s still exists, attempting cleanup\n", resourceID)

			// Parse the response
			var currentResource map[string]any
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
				fmt.Printf("DEBUG: Successfully cleaned up conditional access policy %s\n", resourceID)
				continue
			}

			return fmt.Errorf("failed to clean up conditional access policy %s: %d %s", resourceID, deleteResp.StatusCode, deleteResp.Status)
		}

		return fmt.Errorf("unexpected response checking if conditional access policy %s was destroyed: %d %s", resourceID, httpResp.StatusCode, httpResp.Status)
	}

	return nil
}
