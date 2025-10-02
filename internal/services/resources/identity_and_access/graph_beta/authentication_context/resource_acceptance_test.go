package graphBetaAuthenticationContext_test

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

func TestAccAuthenticationContextResource_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckAuthenticationContextDestroy,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccAuthenticationContextConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_authentication_context.test", "id", "c90"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_authentication_context.test", "display_name", "Acceptance Test Context"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_authentication_context.test", "description", "Context for acceptance testing"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_authentication_context.test", "is_available", "true"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_identity_and_access_authentication_context.test", "id"),
				),
			},
		},
	})
}

func TestAccAuthenticationContextResource_Update(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckAuthenticationContextDestroy,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccAuthenticationContextConfigUpdate1(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_authentication_context.test", "display_name", "Initial Context"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_authentication_context.test", "is_available", "true"),
				),
			},
			{
				Config: testAccAuthenticationContextConfigUpdate2(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_authentication_context.test", "display_name", "Updated Context"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_authentication_context.test", "description", "Updated description"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_authentication_context.test", "is_available", "false"),
				),
			},
		},
	})
}

func TestAccAuthenticationContextResource_Import(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckAuthenticationContextDestroy,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccAuthenticationContextConfigImport(),
			},
			{
				ResourceName:      "microsoft365_graph_beta_identity_and_access_authentication_context.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAuthenticationContextResource_Minimal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckAuthenticationContextDestroy,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccAuthenticationContextConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_authentication_context.test", "id", "c93"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_authentication_context.test", "display_name", "Minimal Context"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_authentication_context.test", "is_available", "false"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_identity_and_access_authentication_context.test", "id"),
				),
			},
		},
	})
}

func testAccCheckAuthenticationContextDestroy(s *terraform.State) error {
	graphClient, err := acceptance.TestGraphClient()
	if err != nil {
		return fmt.Errorf("error creating Graph client for CheckDestroy: %v", err)
	}
	ctx := context.Background()
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "microsoft365_graph_beta_identity_and_access_authentication_context" {
			continue
		}
		_, err := graphClient.
			Identity().
			ConditionalAccess().
			AuthenticationContextClassReferences().
			ByAuthenticationContextClassReferenceId(rs.Primary.ID).
			Get(ctx, nil)

		if err != nil {
			errorInfo := errors.GraphError(ctx, err)

			if errorInfo.StatusCode == 404 ||
				errorInfo.ErrorCode == "ResourceNotFound" ||
				errorInfo.ErrorCode == "ItemNotFound" {
				fmt.Printf("DEBUG: Resource %s successfully destroyed (404/NotFound)\n", rs.Primary.ID)
				continue
			}
			return fmt.Errorf("error checking if authentication context %s was destroyed: %v", rs.Primary.ID, err)
		}

		return fmt.Errorf("authentication context %s still exists", rs.Primary.ID)
	}

	return nil
}

func testAccAuthenticationContextConfigBasic() string {
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/01_basic.tf")
	if err != nil {
		panic(fmt.Sprintf("failed to load acceptance test config: %s", err.Error()))
	}
	return acceptance.ConfiguredM365ProviderBlock(accTestConfig)
}

func testAccAuthenticationContextConfigUpdate1() string {
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/02_update_1.tf")
	if err != nil {
		panic(fmt.Sprintf("failed to load acceptance test config: %s", err.Error()))
	}
	return acceptance.ConfiguredM365ProviderBlock(accTestConfig)
}

func testAccAuthenticationContextConfigUpdate2() string {
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/02_update_2.tf")
	if err != nil {
		panic(fmt.Sprintf("failed to load acceptance test config: %s", err.Error()))
	}
	return acceptance.ConfiguredM365ProviderBlock(accTestConfig)
}

func testAccAuthenticationContextConfigImport() string {
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/03_import.tf")
	if err != nil {
		panic(fmt.Sprintf("failed to load acceptance test config: %s", err.Error()))
	}
	return acceptance.ConfiguredM365ProviderBlock(accTestConfig)
}

func testAccAuthenticationContextConfigMinimal() string {
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/04_minimal.tf")
	if err != nil {
		panic(fmt.Sprintf("failed to load acceptance test config: %s", err.Error()))
	}
	return acceptance.ConfiguredM365ProviderBlock(accTestConfig)
}
