package graphBetaServicePrincipalSynchronizationJob_test

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
	job "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/applications/graph_beta/service_principal_synchronization_job"
)

func TestAccResourceServicePrincipalSynchronizationJob_Lifecycle(t *testing.T) {
	parent := synchronizationJobAcceptancePrincipal(t)
	template, err := helpers.ParseHCLFile("tests/terraform/acceptance/resource_minimal.tf")
	require.NoError(t, err)
	config := fmt.Sprintf(template, parent)
	address := resourceType + ".test"
	resource.Test(
		t,
		resource.TestCase{
			PreCheck:                 func() { mocks.TestAccPreCheck(t) },
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			CheckDestroy:             destroy.CheckDestroyedAllFunc(job.ServicePrincipalSynchronizationJobTestResource{}, resourceType, 20*time.Second),
			Steps: []resource.TestStep{
				{
					Config: config,
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(address, "service_principal_id", parent),
						resource.TestCheckResourceAttr(address, "enabled", "false"),
						resource.TestCheckResourceAttr(address, "schedule.state", "Disabled"),
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
