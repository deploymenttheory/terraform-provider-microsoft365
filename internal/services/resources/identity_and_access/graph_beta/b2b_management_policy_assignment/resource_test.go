package graphBetaIdentityAndAccessB2bManagementPolicyAssignment_test

import (
	"regexp"
	"testing"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaB2bManagementPolicyAssignment "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_beta/b2b_management_policy_assignment"
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

// Replication lag during create: policy GET 404, $ref POST 404 (retried once), stale appliesTo.
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

// Replays Graph's truncated appliesTo body (no Application.Read.All): Read must fail with a
// permission hint and keep the state.
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

// A stale empty appliesTo on refresh must not drop the assignment from state.
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

// A pre-existing assignment must be reported for import, not adopted.
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
				ExpectError: regexp.MustCompile(`import\s+it\s+with\s+ID`),
			},
		},
	})

	assertCallCount(t, assignmentMock, "POST", "servicePrincipals", 1)
}

// Same after a retried 404: the 404 wrote nothing, so "already exist" still means pre-existing.
func TestUnitResourceB2bManagementPolicyAssignment_09_AlreadyAssignedAfterNotFound(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, assignmentMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer assignmentMock.CleanupMockState()

	assignmentMock.SeedAssignment(policyID, assignmentMocks.ServicePrincipalID)
	assignmentMock.RegisterPostNotFound(1)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testConfig("tests/terraform/unit/resource_service_principal.tf"),
				ExpectError: regexp.MustCompile(`import\s+it\s+with\s+ID`),
			},
		},
	})

	assertCallCount(t, assignmentMock, "POST", "servicePrincipals", 2)
}

// An assignment that is really gone is dropped once the confirmation window elapses.
func TestUnitResourceB2bManagementPolicyAssignment_10_RemovedOutOfBand(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, assignmentMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer assignmentMock.CleanupMockState()
	defer graphBetaB2bManagementPolicyAssignment.SetNotFoundConfirmationForTesting(2*time.Second, 500*time.Millisecond)()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfig("tests/terraform/unit/resource_service_principal.tf"),
			},
			{
				PreConfig: func() {
					assignmentMock.CleanupMockState()
				},
				Config:             testConfig("tests/terraform/unit/resource_service_principal.tf"),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
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
