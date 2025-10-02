package graphConditionalAccessTermsOfUse_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccConditionalAccessTermsOfUseResource_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckConditionalAccessTermsOfUseDestroy,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccConditionalAccessTermsOfUseConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("microsoft365_graph_identity_and_access_conditional_access_terms_of_use.test", "display_name", "Acceptance Test Terms of Use"),
					resource.TestCheckResourceAttr("microsoft365_graph_identity_and_access_conditional_access_terms_of_use.test", "is_viewing_before_acceptance_required", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_identity_and_access_conditional_access_terms_of_use.test", "is_per_device_acceptance_required", "false"),
					resource.TestCheckResourceAttr("microsoft365_graph_identity_and_access_conditional_access_terms_of_use.test", "user_reaccept_required_frequency", "P90D"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_identity_and_access_conditional_access_terms_of_use.test", "id"),
				),
			},
		},
	})
}

func TestAccConditionalAccessTermsOfUseResource_Update(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckConditionalAccessTermsOfUseDestroy,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccConditionalAccessTermsOfUseConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("microsoft365_graph_identity_and_access_conditional_access_terms_of_use.test", "display_name", "Acceptance Test Terms of Use"),
					resource.TestCheckResourceAttr("microsoft365_graph_identity_and_access_conditional_access_terms_of_use.test", "is_viewing_before_acceptance_required", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_identity_and_access_conditional_access_terms_of_use.test", "is_per_device_acceptance_required", "false"),
				),
			},
			{
				Config: testAccConditionalAccessTermsOfUseConfigUpdate(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("microsoft365_graph_identity_and_access_conditional_access_terms_of_use.test", "display_name", "Updated Acceptance Test Terms of Use"),
					resource.TestCheckResourceAttr("microsoft365_graph_identity_and_access_conditional_access_terms_of_use.test", "is_viewing_before_acceptance_required", "false"),
					resource.TestCheckResourceAttr("microsoft365_graph_identity_and_access_conditional_access_terms_of_use.test", "is_per_device_acceptance_required", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_identity_and_access_conditional_access_terms_of_use.test", "user_reaccept_required_frequency", "P180D"),
				),
			},
		},
	})
}

func TestAccConditionalAccessTermsOfUseResource_Import(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckConditionalAccessTermsOfUseDestroy,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccConditionalAccessTermsOfUseConfigBasic(),
			},
			{
				ResourceName:      "microsoft365_graph_identity_and_access_conditional_access_terms_of_use.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckConditionalAccessTermsOfUseDestroy(s *terraform.State) error {
	graphClient, err := acceptance.TestGraphClient()
	if err != nil {
		return fmt.Errorf("error creating Graph client for CheckDestroy: %v", err)
	}
	ctx := context.Background()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "microsoft365_graph_identity_and_access_conditional_access_terms_of_use" {
			continue
		}

		agreementId := rs.Primary.ID

		_, err := graphClient.
			Agreements().
			ByAgreementId(agreementId).
			Get(ctx, nil)

		if err != nil {
			errorInfo := errors.GraphError(ctx, err)

			if errorInfo.StatusCode == 404 ||
				errorInfo.ErrorCode == "ResourceNotFound" ||
				errorInfo.ErrorCode == "ItemNotFound" {
				fmt.Printf("DEBUG: Resource %s successfully destroyed (404/NotFound)\n", rs.Primary.ID)
				continue
			}
			return fmt.Errorf("error checking if conditional access terms of use %s was destroyed: %v", rs.Primary.ID, err)
		}

		return fmt.Errorf("conditional access terms of use %s still exists", rs.Primary.ID)
	}

	return nil
}

func testAccConditionalAccessTermsOfUseConfigBasic() string {
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/01_basic.tf")
	if err != nil {
		panic(fmt.Sprintf("failed to load acceptance test config: %s", err.Error()))
	}
	return acceptance.ConfiguredM365ProviderBlock(accTestConfig)
}

func testAccConditionalAccessTermsOfUseConfigUpdate() string {
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/02_update.tf")
	if err != nil {
		panic(fmt.Sprintf("failed to load acceptance test config: %s", err.Error()))
	}
	return acceptance.ConfiguredM365ProviderBlock(accTestConfig)
}
