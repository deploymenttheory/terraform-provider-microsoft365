package graphBetaServicePrincipalSynchronizationJob_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	job "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/applications/graph_beta/service_principal_synchronization_job"
	jobmocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/applications/graph_beta/service_principal_synchronization_job/mocks"
)

var resourceType = job.ResourceName

func setupMockEnvironment(t *testing.T) *jobmocks.ServicePrincipalSynchronizationJobMock {
	t.Helper()
	mocks.SetupUnitTestEnvironment(t)
	httpmock.Activate()
	mocks.NewMocks().AuthMocks.RegisterMocks()
	mock := &jobmocks.ServicePrincipalSynchronizationJobMock{}
	mock.CleanupMockState()
	mock.RegisterMocks()
	t.Cleanup(func() { mock.CleanupMockState(); httpmock.DeactivateAndReset() })
	return mock
}

func loadUnitTestTerraform(t *testing.T, file string) string {
	t.Helper()
	value, err := helpers.ParseHCLFile("tests/terraform/unit/" + file)
	require.NoError(t, err)
	return value
}

func TestUnitResourceServicePrincipalSynchronizationJob_Lifecycle(t *testing.T) {
	mock := setupMockEnvironment(t)
	minimal := loadUnitTestTerraform(t, "resource_minimal.tf")
	address := resourceType + ".test"
	resource.UnitTest(t, resource.TestCase{ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories, Steps: []resource.TestStep{
		{
			Config: minimal,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr(address, "id", jobmocks.ServicePrincipalID+"/"+jobmocks.JobID),
				resource.TestCheckResourceAttr(address, "job_id", jobmocks.JobID),
				resource.TestCheckResourceAttr(address, "enabled", "false"),
				resource.TestCheckResourceAttr(address, "schedule.state", "Disabled"),
			),
		},
		{ResourceName: address, ImportState: true, ImportStateVerify: true},
		{Config: loadUnitTestTerraform(t, "resource_enabled.tf"), Check: resource.TestCheckResourceAttr(address, "schedule.state", "Active")},
		{Config: minimal, Check: resource.TestCheckResourceAttr(address, "schedule.state", "Paused")},
		{Config: minimal, PlanOnly: true},
	}})
	creates, starts, pauses := mock.Actions()
	require.Equal(t, 1, creates)
	require.Equal(t, 1, starts)
	require.Equal(t, 1, pauses, "refresh must never start/pause or recreate a job")
}

func TestUnitResourceServicePrincipalSynchronizationJob_EnabledCreate(t *testing.T) {
	mock := setupMockEnvironment(t)
	resource.UnitTest(
		t,
		resource.TestCase{
			ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{Config: loadUnitTestTerraform(t, "resource_enabled.tf"), Check: resource.TestCheckResourceAttr(resourceType+".test", "enabled", "true")},
			},
		},
	)
	creates, starts, pauses := mock.Actions()
	require.Equal(t, 1, creates)
	require.Equal(t, 1, starts)
	require.Zero(t, pauses)
}

func TestUnitResourceServicePrincipalSynchronizationJob_GeneratedImport(t *testing.T) {
	mock := setupMockEnvironment(t)
	require.NoError(t, mock.SeedFixture())
	address := resourceType + ".test"
	resource.UnitTest(
		t,
		resource.TestCase{
			ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config:                 loadUnitTestTerraform(t, "import.tf"),
					ResourceName:           address,
					ImportState:            true,
					ImportStateId:          jobmocks.ServicePrincipalID + "/" + jobmocks.JobID,
					GenerateConfig:         true,
					ImportStateKind:        resource.ImportBlockWithID,
					ImportStateConfigExact: true,
					ImportPlanChecks:       resource.ImportPlanChecks{PreApply: []plancheck.PlanCheck{plancheck.ExpectResourceAction(address, plancheck.ResourceActionNoop)}},
				},
			},
		},
	)
}
