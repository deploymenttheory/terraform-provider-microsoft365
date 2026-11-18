package graphBetaOAuth2PermissionGrant_test

import (
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/destroy"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
)

func loadAcceptanceTerraform(file string) (string, error) {
	return helpers.ParseHCLFile("tests/terraform/acceptance/" + file)
}

func TestAccResourceOAuth2PermissionGrant_Lifecycle(t *testing.T) {
	for _, principal := range []bool{false, true} {
		name := "AllPrincipals"
		if principal {
			name = "Principal"
		}
		t.Run(name, func(t *testing.T) {
			fixture := grantAcceptanceFixture(t, principal)
			config := fixture.config(t, principal)
			updated := strings.Replace(config, `["Probe.Read"]`, `["Probe.Write", "Probe.Read"]`, 1)
			address := resourceType + ".test"
			resource.Test(t, resource.TestCase{
				PreCheck:                 func() { mocks.TestAccPreCheck(t) },
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 20*time.Second),
				Steps: []resource.TestStep{
					{
						Config: config,
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr(address, "consent_type", name),
							resource.TestCheckResourceAttr(address, "scopes.#", "1"),
						),
					},
					{ResourceName: address, ImportState: true, ImportStateVerify: true},
					{
						ResourceName:     address,
						ImportState:      true,
						GenerateConfig:   true,
						ImportStateKind:  resource.ImportBlockWithID,
						ImportPlanChecks: resource.ImportPlanChecks{PreApply: []plancheck.PlanCheck{plancheck.ExpectResourceAction(address, plancheck.ResourceActionNoop)}},
					},
					{Config: updated, Check: resource.TestCheckResourceAttr(address, "scopes.#", "2")},
					{Config: config, Check: resource.TestCheckResourceAttr(address, "scopes.#", "1")},
					{Config: config, PlanOnly: true},
				},
			})
		})
	}
}
