package graphBetaDeviceConfigurationTemplatesJson_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/jarcoal/httpmock"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
)

func TestUnitAssignmentGeneratedImport(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, profileMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer profileMock.CleanupMockState()
	config := loadUnitTestTerraform("resource_61_assignments_minimal_to_maximal.tf")
	resource.UnitTest(t, resource.TestCase{
		TerraformVersionChecks:   []tfversion.TerraformVersionCheck{tfversion.SkipBelow(tfversion.Version1_5_0)},
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: config},
			{
				ResourceName: resourceType + ".test",
				ImportState:  true, GenerateConfig: true,
				ImportStateKind: resource.ImportBlockWithID,
				ImportPlanChecks: resource.ImportPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceType+".generated", plancheck.ResourceActionNoop),
					},
				},
			},
			{Config: config, PlanOnly: true},
		},
	})
}
