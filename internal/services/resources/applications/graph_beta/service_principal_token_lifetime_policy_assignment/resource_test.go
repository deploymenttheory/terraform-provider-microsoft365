package graphBetaApplicationsServicePrincipalTokenLifetimePolicyAssignment_test

import (
	"regexp"
	"strings"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	spAssignmentMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/applications/graph_beta/service_principal_token_lifetime_policy_assignment/mocks"
	tokenLifetimePolicyMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/applications/graph_beta/token_lifetime_policy/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

var resourceType = "microsoft365_graph_beta_applications_service_principal_token_lifetime_policy_assignment"

func setupMockEnvironment() (*mocks.Mocks, *spAssignmentMocks.ServicePrincipalTokenLifetimePolicyAssignmentMock, *tokenLifetimePolicyMocks.TokenLifetimePolicyMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	tlpMock := &tokenLifetimePolicyMocks.TokenLifetimePolicyMock{}
	tlpMock.RegisterMocks()
	// Registered after tlpMock so the assignment mock's referenced-policy GET responder
	// (which serves a canned policy for any id) wins over tlpMock's stateful one.
	assignmentMock := &spAssignmentMocks.ServicePrincipalTokenLifetimePolicyAssignmentMock{}
	assignmentMock.RegisterMocks()
	return mockClient, assignmentMock, tlpMock
}

func TestUnitResourceServicePrincipalTokenLifetimePolicyAssignment_01_Basic(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, assignmentMock, tlpMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer assignmentMock.CleanupMockState()
	defer tlpMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceType+".basic", "service_principal_id", "00000000-0000-0000-0000-000000000020"),
					resource.TestCheckResourceAttr(resourceType+".basic", "token_lifetime_policy_id", "00000000-0000-0000-0000-000000000010"),
					resource.TestMatchResourceAttr(resourceType+".basic", "id", regexp.MustCompile(`^[0-9a-fA-F-]+/[0-9a-fA-F-]+$`)),
				),
			},
			{
				ResourceName:      resourceType + ".basic",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     "00000000-0000-0000-0000-000000000020/00000000-0000-0000-0000-000000000010",
			},
		},
	})
}

// TestUnitResourceServicePrincipalTokenLifetimePolicyAssignment_02_EventualConsistency
// simulates Microsoft Entra replication lag during create:
//   - the first GET of the referenced policy returns the 404 seen while a just-created
//     policy has not yet propagated ("Unable to read the company information from the
//     directory") — Create's propagation wait must absorb it before the $ref POST
//   - after the POST succeeds, the first list GET returns an empty collection (stale
//     replica) — without the ConsistencyPredicate this produces "Missing Resource State
//     After Create"
//
// The non-idempotent $ref POST must be issued exactly once.
func TestUnitResourceServicePrincipalTokenLifetimePolicyAssignment_02_EventualConsistency(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, assignmentMock, tlpMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer assignmentMock.CleanupMockState()
	defer tlpMock.CleanupMockState()

	assignmentMock.RegisterEventualConsistencyMocks(1, 1)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceType+".basic", "service_principal_id", "00000000-0000-0000-0000-000000000020"),
					resource.TestCheckResourceAttr(resourceType+".basic", "token_lifetime_policy_id", "00000000-0000-0000-0000-000000000010"),
					resource.TestMatchResourceAttr(resourceType+".basic", "id", regexp.MustCompile(`^[0-9a-fA-F-]+/[0-9a-fA-F-]+$`)),
				),
			},
		},
	})

	postCalls := 0
	for key, count := range httpmock.GetCallCountInfo() {
		if strings.HasPrefix(key, "POST ") && strings.Contains(key, `/tokenLifetimePolicies/\$ref`) {
			postCalls += count
		}
	}
	if postCalls != 1 {
		t.Fatalf("expected the non-idempotent $ref POST to be issued exactly once, got %d", postCalls)
	}
}

func testConfigBasic() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_basic.tf")
	if err != nil {
		panic("failed to load service principal token lifetime policy assignment config: " + err.Error())
	}
	return unitTestConfig
}
