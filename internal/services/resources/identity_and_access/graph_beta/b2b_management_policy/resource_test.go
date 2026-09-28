package graphBetaIdentityAndAccessB2bManagementPolicy_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	b2bManagementPolicyMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_beta/b2b_management_policy/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

var resourceType = "microsoft365_graph_beta_identity_and_access_b2b_management_policy"

func setupMockEnvironment() (*mocks.Mocks, *b2bManagementPolicyMocks.B2bManagementPolicyMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	policyMock := &b2bManagementPolicyMocks.B2bManagementPolicyMock{}
	policyMock.RegisterMocks()
	return mockClient, policyMock
}

// TestUnitResourceB2bManagementPolicy_01_Minimal also guards against sending description: the
// mock rejects it with the same 404 Microsoft Graph returns.
func TestUnitResourceB2bManagementPolicy_01_Minimal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, policyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer policyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfig("tests/terraform/unit/resource_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceType+".minimal", "display_name", "unit-test-b2b-management-policy-min"),
					resource.TestCheckResourceAttr(resourceType+".minimal", "is_organization_default", "false"),
					resource.TestCheckResourceAttr(resourceType+".minimal", "definition.#", "1"),
					resource.TestCheckResourceAttr(resourceType+".minimal", "definition.0", `{"B2BManagementPolicy":{"InvitationsAllowedAndBlockedDomainsPolicy":{"BlockedDomains":["example.com"]}}}`),
					resource.TestMatchResourceAttr(resourceType+".minimal", "id", regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
				),
			},
			{
				ResourceName:      resourceType + ".minimal",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"timeouts",
				},
			},
		},
	})
}

func TestUnitResourceB2bManagementPolicy_02_Maximal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, policyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer policyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfig("tests/terraform/unit/resource_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceType+".maximal", "display_name", "unit-test-b2b-management-policy-max"),
					resource.TestCheckResourceAttr(resourceType+".maximal", "is_organization_default", "true"),
					resource.TestCheckResourceAttr(resourceType+".maximal", "definition.#", "1"),
					resource.TestMatchResourceAttr(resourceType+".maximal", "id", regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
				),
			},
		},
	})
}

func TestUnitResourceB2bManagementPolicy_03_Update(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, policyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer policyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfig("tests/terraform/unit/resource_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceType+".minimal", "display_name", "unit-test-b2b-management-policy-min"),
				),
			},
			{
				Config: testConfig("tests/terraform/unit/resource_minimal_updated.tf"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceType+".minimal", "display_name", "unit-test-b2b-management-policy-min-updated"),
					resource.TestCheckResourceAttr(resourceType+".minimal", "definition.0", `{"B2BManagementPolicy":{"InvitationsAllowedAndBlockedDomainsPolicy":{"AllowedDomains":["contoso.com"]}}}`),
				),
			},
		},
	})
}

// TestUnitResourceB2bManagementPolicy_04_EventualConsistency reproduces the replica flapping seen
// live after PATCH: without the consistency predicate the stale read would be stored and the
// post-apply plan would not be empty.
func TestUnitResourceB2bManagementPolicy_04_EventualConsistency(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, policyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer policyMock.CleanupMockState()

	policyMock.RegisterEventualConsistencyMocks(2)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfig("tests/terraform/unit/resource_minimal.tf"),
			},
			{
				Config: testConfig("tests/terraform/unit/resource_minimal_updated.tf"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceType+".minimal", "display_name", "unit-test-b2b-management-policy-min-updated"),
				),
			},
		},
	})
}

func TestUnitResourceB2bManagementPolicy_05_InvalidDefinition(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, policyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer policyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testConfig("tests/terraform/unit/resource_invalid_definition.tf"),
				ExpectError: regexp.MustCompile(`Invalid JSON String`),
			},
		},
	})
}

func TestUnitResourceB2bManagementPolicy_06_CreateError(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, policyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer policyMock.CleanupMockState()

	policyMock.RegisterErrorMocks()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testConfig("tests/terraform/unit/resource_minimal.tf"),
				ExpectError: regexp.MustCompile(`Property definition has an invalid value`),
			},
		},
	})
}

// TestUnitResourceB2bManagementPolicy_07_StaleNotFoundOnRefresh reproduces the stale-replica 404
// observed live on the first refresh after create. Removing the policy from state on that 404
// would make the plan create a duplicate and orphan the existing policy; the refresh must keep
// re-reading until the policy is visible and produce an empty plan.
func TestUnitResourceB2bManagementPolicy_07_StaleNotFoundOnRefresh(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, policyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer policyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfig("tests/terraform/unit/resource_minimal.tf"),
			},
			{
				PreConfig: func() {
					policyMock.RegisterStaleNotFoundMocks(2)
				},
				Config:   testConfig("tests/terraform/unit/resource_minimal.tf"),
				PlanOnly: true,
			},
		},
	})
}

func testConfig(path string) string {
	unitTestConfig, err := helpers.ParseHCLFile(path)
	if err != nil {
		panic("failed to load B2B management policy config " + path + ": " + err.Error())
	}
	return unitTestConfig
}
