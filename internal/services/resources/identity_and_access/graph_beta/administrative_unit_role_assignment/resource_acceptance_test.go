package graphBetaAdministrativeUnitRoleAssignment_test

import (
	"testing"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/destroy"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/testlog"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func loadAcceptanceTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/acceptance/" + filename)
	if err != nil {
		panic("failed to load acceptance config " + filename + ": " + err.Error())
	}
	return config
}

func acceptanceImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", nil
		}
		auID := rs.Primary.Attributes["administrative_unit_id"]
		id := rs.Primary.Attributes["id"]
		return auID + "/" + id, nil
	}
}

// AURA001: Basic scoped role assignment — creates an AU and assigns a User Administrator role
func TestAccResourceAdministrativeUnitRoleAssignment_01_AURA001(t *testing.T) {
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
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
			"time": {
				Source:            "hashicorp/time",
				VersionConstraint: constants.ExternalProviderTimeVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating AURA001 basic scoped role assignment")
				},
				Config: loadAcceptanceTestTerraform("resource_aura001_basic.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("administrative unit role assignment", 15*time.Second)
						time.Sleep(15 * time.Second)
						return nil
					},
					check.That(resourceType+".aura001_basic").ExistsInGraph(testResource),
					check.That(resourceType+".aura001_basic").Key("id").IsNotEmpty(),
					check.That(resourceType+".aura001_basic").Key("role_id").IsNotEmpty(),
					check.That(resourceType+".aura001_basic").Key("role_member_id").IsNotEmpty(),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing AURA001 scoped role assignment")
				},
				Config:            loadAcceptanceTestTerraform("resource_aura001_basic.tf"),
				ResourceName:      resourceType + ".aura001_basic",
				ImportState:       true,
				ImportStateIdFunc: acceptanceImportStateIDFunc(resourceType + ".aura001_basic"),
				ImportStateVerify: true,
			},
		},
	})
}

// AURA002: Different role — assigns a Helpdesk Administrator role to a user
func TestAccResourceAdministrativeUnitRoleAssignment_02_AURA002(t *testing.T) {
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
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
			"time": {
				Source:            "hashicorp/time",
				VersionConstraint: constants.ExternalProviderTimeVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating AURA002 helpdesk administrator role assignment")
				},
				Config: loadAcceptanceTestTerraform("resource_aura002_helpdesk_admin.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("administrative unit role assignment", 15*time.Second)
						time.Sleep(15 * time.Second)
						return nil
					},
					check.That(resourceType+".aura002_helpdesk_admin").ExistsInGraph(testResource),
					check.That(resourceType+".aura002_helpdesk_admin").Key("id").IsNotEmpty(),
					check.That(resourceType+".aura002_helpdesk_admin").Key("role_id").IsNotEmpty(),
					check.That(resourceType+".aura002_helpdesk_admin").Key("role_member_id").IsNotEmpty(),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing AURA002 scoped role assignment")
				},
				Config:            loadAcceptanceTestTerraform("resource_aura002_helpdesk_admin.tf"),
				ResourceName:      resourceType + ".aura002_helpdesk_admin",
				ImportState:       true,
				ImportStateIdFunc: acceptanceImportStateIDFunc(resourceType + ".aura002_helpdesk_admin"),
				ImportStateVerify: true,
			},
		},
	})
}

// AURA003: Replace test — changing role_id forces a new resource
func TestAccResourceAdministrativeUnitRoleAssignment_03_AURA003_Replace(t *testing.T) {
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
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
			"time": {
				Source:            "hashicorp/time",
				VersionConstraint: constants.ExternalProviderTimeVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating AURA003 - initial role assignment")
				},
				Config: loadAcceptanceTestTerraform("resource_aura003_replace_step1.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("administrative unit role assignment", 15*time.Second)
						time.Sleep(15 * time.Second)
						return nil
					},
					check.That(resourceType+".aura003_replace").ExistsInGraph(testResource),
					check.That(resourceType+".aura003_replace").Key("id").IsNotEmpty(),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Replacing AURA003 - changing role (forces new resource)")
					testlog.WaitForConsistency("administrative unit role assignment", 20*time.Second)
					time.Sleep(20 * time.Second)
				},
				Config: loadAcceptanceTestTerraform("resource_aura003_replace_step2.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("administrative unit role assignment", 15*time.Second)
						time.Sleep(15 * time.Second)
						return nil
					},
					check.That(resourceType+".aura003_replace").ExistsInGraph(testResource),
					check.That(resourceType+".aura003_replace").Key("id").IsNotEmpty(),
				),
			},
		},
	})
}
