package graphBetaAdminConsentRequestPolicy_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	policy "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_beta/admin_consent_request_policy"
	policymocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_beta/admin_consent_request_policy/mocks"
)

var resourceType = policy.ResourceName

func setupMockEnvironment(t *testing.T) *policymocks.AdminConsentRequestPolicyMock {
	t.Helper()
	mocks.SetupUnitTestEnvironment(t)
	httpmock.Activate()
	mocks.NewMocks().AuthMocks.RegisterMocks()
	mock := &policymocks.AdminConsentRequestPolicyMock{}
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

func TestUnitResourceAdminConsentRequestPolicy_Lifecycle(t *testing.T) {
	mock := setupMockEnvironment(t)
	minimal := loadUnitTestTerraform(t, "resource_minimal.tf")
	address := resourceType + ".test"
	resource.UnitTest(t, resource.TestCase{ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories, Steps: []resource.TestStep{
		{
			Config: minimal,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr(address, "id", "adminConsentRequestPolicy"),
				resource.TestCheckResourceAttr(address, "request_duration_in_days", "0"),
				resource.TestCheckResourceAttr(address, "reviewers.#", "0"),
				resource.TestCheckResourceAttr(address, "version", "1"),
			),
		},
		{ResourceName: address, ImportState: true, ImportStateVerify: true},
		{
			Config: loadUnitTestTerraform(t, "resource_maximal.tf"),
			Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr(address, "reviewers.#", "1"), resource.TestCheckResourceAttr(address, "version", "2")),
		},
		{Config: minimal, Check: resource.TestCheckResourceAttr(address, "version", "3")},
		{Config: minimal, PlanOnly: true},
	}})
	require.Len(t, mock.Writes(), 3, "create/update must use PUT; destroy must retain the policy")
}

func TestUnitResourceAdminConsentRequestPolicy_GeneratedImport(t *testing.T) {
	setupMockEnvironment(t)
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
					ImportStateId:          "adminConsentRequestPolicy",
					GenerateConfig:         true,
					ImportStateKind:        resource.ImportBlockWithID,
					ImportStateConfigExact: true,
					ImportPlanChecks:       resource.ImportPlanChecks{PreApply: []plancheck.PlanCheck{plancheck.ExpectResourceAction(address, plancheck.ResourceActionNoop)}},
				},
			},
		},
	)
}
