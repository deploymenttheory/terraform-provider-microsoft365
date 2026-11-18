package graphBetaDirectoryRoleAssignment_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/stretchr/testify/require"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/destroy"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	assignment "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_beta/directory_role_assignment"
)

func TestAccResourceDirectoryRoleAssignment_Lifecycle(t *testing.T) {
	principalID := directoryRoleAcceptancePrincipal(t)
	template, err := helpers.ParseHCLFile("tests/terraform/acceptance/resource_minimal.tf")
	require.NoError(t, err)
	config := fmt.Sprintf(template, principalID)
	address := resourceType + ".test"
	resource.Test(
		t,
		resource.TestCase{
			PreCheck:                 func() { mocks.TestAccPreCheck(t) },
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			CheckDestroy:             destroy.CheckDestroyedAllFunc(assignment.DirectoryRoleAssignmentTestResource{}, resourceType, 20*time.Second),
			Steps: []resource.TestStep{
				{
					Config: config,
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(address, "principal_id", principalID),
						resource.TestCheckResourceAttr(address, "directory_scope_id", "/"),
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
				{Config: config, PlanOnly: true},
			},
		},
	)
}
