package graphBetaApplication_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/destroy"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
)

func TestAccResourceApplication_EmptyCollectionsGeneratedImport(t *testing.T) {
	config := loadAccTestTerraform("resource_empty_collections.tf")
	updated := strings.Replace(config, "Empty collection acceptance test", "Updated empty collection acceptance test", 1)
	address := resourceType + ".test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		TerraformVersionChecks:   []tfversion.TerraformVersionCheck{tfversion.SkipBelow(tfversion.Version1_5_0)},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 20*time.Second),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {Source: "hashicorp/random", VersionConstraint: constants.ExternalProviderRandomVersion},
		},
		Steps: []resource.TestStep{
			{Config: config, Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr(address, "api.known_client_applications.#", "0"),
				resource.TestCheckResourceAttr(address, "api.oauth2_permission_scopes.#", "0"),
				resource.TestCheckResourceAttr(address, "api.pre_authorized_applications.#", "0"),
				resource.TestCheckResourceAttr(address, "app_roles.#", "0"),
				resource.TestCheckResourceAttr(address, "required_resource_access.#", "0"),
				resource.TestCheckResourceAttr(address, "is_device_only_auth_supported", "false"),
				resource.TestCheckResourceAttr(address, "is_fallback_public_client", "false"),
				resource.TestCheckResourceAttr(address, "api.requested_access_token_version", "1"),
			)},
			{
				ResourceName: address, ImportState: true, GenerateConfig: true,
				ImportStateKind: resource.ImportBlockWithID,
				ImportStateIdFunc: func(state *terraform.State) (string, error) {
					application, ok := state.RootModule().Resources[address]
					if !ok {
						return "", fmt.Errorf("resource not found: %s", address)
					}
					return application.Primary.ID + ":hard_delete=true", nil
				},
				ImportPlanChecks: resource.ImportPlanChecks{PreApply: []plancheck.PlanCheck{
					plancheck.ExpectResourceAction(address, plancheck.ResourceActionNoop),
				}},
			},
			{Config: updated, Check: resource.TestCheckResourceAttr(address, "description", "Updated empty collection acceptance test")},
			{Config: updated, PlanOnly: true},
		},
	})
}
