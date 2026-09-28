package graphBetaIdentityAndAccessB2bManagementPolicyAssignment_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	assignmentMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_beta/b2b_management_policy_assignment/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

var resourceType = "microsoft365_graph_beta_identity_and_access_b2b_management_policy_assignment"

const policyID = "00000000-0000-0000-0000-000000000010"

func setupMockEnvironment() (*mocks.Mocks, *assignmentMocks.B2bManagementPolicyAssignmentMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	assignmentMock := &assignmentMocks.B2bManagementPolicyAssignmentMock{}
	assignmentMock.RegisterMocks()
	return mockClient, assignmentMock
}

func TestUnitResourceB2bManagementPolicyAssignment_01_ServicePrincipal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, assignmentMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer assignmentMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfig("tests/terraform/unit/resource_service_principal.tf"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceType+".service_principal", "b2b_management_policy_id", policyID),
					resource.TestCheckResourceAttr(resourceType+".service_principal", "directory_object_id", assignmentMocks.ServicePrincipalID),
					resource.TestCheckResourceAttr(resourceType+".service_principal", "directory_object_type", "servicePrincipal"),
					resource.TestCheckResourceAttr(resourceType+".service_principal", "id", policyID+"/"+assignmentMocks.ServicePrincipalID),
				),
			},
			{
				ResourceName:      resourceType + ".service_principal",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     policyID + "/" + assignmentMocks.ServicePrincipalID,
				ImportStateVerifyIgnore: []string{
					"timeouts",
				},
			},
		},
	})

	assertCallCount(t, assignmentMock, "POST", "servicePrincipals", 1)
	assertCallCount(t, assignmentMock, "DELETE", "servicePrincipals", 1)
}

func TestUnitResourceB2bManagementPolicyAssignment_02_Application(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, assignmentMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer assignmentMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfig("tests/terraform/unit/resource_application.tf"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceType+".application", "directory_object_id", assignmentMocks.ApplicationID),
					resource.TestCheckResourceAttr(resourceType+".application", "directory_object_type", "application"),
				),
			},
		},
	})

	assertCallCount(t, assignmentMock, "POST", "applications", 1)
	assertCallCount(t, assignmentMock, "DELETE", "applications", 1)
}

func TestUnitResourceB2bManagementPolicyAssignment_03_UnsupportedDirectoryObject(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, assignmentMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer assignmentMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testConfig("tests/terraform/unit/resource_user.tf"),
				ExpectError: regexp.MustCompile(`can\s+only\s+be\s+applied\s+to\s+applications\s+and\s+service\s+principals`),
			},
		},
	})
}

// TestUnitResourceB2bManagementPolicyAssignment_04_EventualConsistency simulates Microsoft Entra
// replication lag during create: the referenced policy is briefly not visible, the first $ref
// POST is rejected with 404 by a replica the policy has not reached yet, and the first appliesTo
// read after the POST comes from a stale replica. The rejected POST is retried once.
func TestUnitResourceB2bManagementPolicyAssignment_04_EventualConsistency(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, assignmentMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer assignmentMock.CleanupMockState()

	assignmentMock.RegisterEventualConsistencyMocks(1, 1)
	assignmentMock.RegisterPostNotFound(1)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfig("tests/terraform/unit/resource_service_principal.tf"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceType+".service_principal", "id", policyID+"/"+assignmentMocks.ServicePrincipalID),
				),
			},
		},
	})

	assertCallCount(t, assignmentMock, "POST", "servicePrincipals", 2)
}

func TestUnitResourceB2bManagementPolicyAssignment_05_CreateError(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, assignmentMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer assignmentMock.CleanupMockState()

	assignmentMock.RegisterErrorMocks()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testConfig("tests/terraform/unit/resource_service_principal.tf"),
				ExpectError: regexp.MustCompile(`Insufficient privileges|Forbidden|403`),
			},
		},
	})
}

// TestUnitResourceB2bManagementPolicyAssignment_06_TruncatedAppliesTo replays the response
// Microsoft Graph returns to callers without Application.Read.All once an application is in
// appliesTo: HTTP 200 with a truncated body. Read must fail with an actionable error instead of
// a bare JSON parse error, and must not drop the assignment from state.
func TestUnitResourceB2bManagementPolicyAssignment_06_TruncatedAppliesTo(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, assignmentMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer assignmentMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfig("tests/terraform/unit/resource_application.tf"),
			},
			{
				PreConfig: func() {
					assignmentMock.RegisterTruncatedAppliesToMocks()
				},
				Config:      testConfig("tests/terraform/unit/resource_application.tf"),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`Application\.Read\.All`),
			},
			{
				PreConfig: func() {
					assignmentMock.RegisterMocks()
				},
				Config: testConfig("tests/terraform/unit/resource_application.tf"),
			},
		},
	})
}

// TestUnitResourceB2bManagementPolicyAssignment_07_StaleAppliesToOnRefresh reproduces a stale
// replica returning an empty appliesTo on refresh. Removing the assignment from state on that read
// would make the plan re-create it; the refresh must keep re-reading and produce an empty plan.
func TestUnitResourceB2bManagementPolicyAssignment_07_StaleAppliesToOnRefresh(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, assignmentMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer assignmentMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfig("tests/terraform/unit/resource_service_principal.tf"),
			},
			{
				PreConfig: func() {
					assignmentMock.RegisterStaleAppliesToMocks(2)
				},
				Config:   testConfig("tests/terraform/unit/resource_service_principal.tf"),
				PlanOnly: true,
			},
		},
	})

	assertCallCount(t, assignmentMock, "POST", "servicePrincipals", 1)
}

// TestUnitResourceB2bManagementPolicyAssignment_08_AlreadyAssigned ensures an assignment that
// already exists before Create is reported (so it can be imported) instead of being adopted.
func TestUnitResourceB2bManagementPolicyAssignment_08_AlreadyAssigned(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, assignmentMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer assignmentMock.CleanupMockState()

	assignmentMock.SeedAssignment(policyID, assignmentMocks.ServicePrincipalID)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testConfig("tests/terraform/unit/resource_service_principal.tf"),
				ExpectError: regexp.MustCompile(`already\s+exist`),
			},
		},
	})

	assertCallCount(t, assignmentMock, "POST", "servicePrincipals", 1)
}

// assertCallCount checks how many policy reference requests the mock received.
func assertCallCount(t *testing.T, assignmentMock *assignmentMocks.B2bManagementPolicyAssignmentMock, method, collection string, expected int) {
	t.Helper()

	if calls := assignmentMock.CallCount(method, collection); calls != expected {
		t.Fatalf("expected %d %s call(s) to %s policies/$ref, got %d", expected, method, collection, calls)
	}
}

func testConfig(path string) string {
	unitTestConfig, err := helpers.ParseHCLFile(path)
	if err != nil {
		panic("failed to load B2B management policy assignment config " + path + ": " + err.Error())
	}
	return unitTestConfig
}
