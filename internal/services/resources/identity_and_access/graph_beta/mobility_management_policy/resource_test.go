package graphBetaIdentityAndAccessMobilityManagementPolicy_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/jarcoal/httpmock"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	mobilityMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_beta/mobility_management_policy/mocks"
)

func setupMockEnvironment() (*mocks.Mocks, *mobilityMocks.MobilityManagementPolicyMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	mobilityMock := &mobilityMocks.MobilityManagementPolicyMock{}
	mobilityMock.RegisterMocks()
	return mockClient, mobilityMock
}

func setupErrorMockEnvironment() (*mocks.Mocks, *mobilityMocks.MobilityManagementPolicyMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	mobilityMock := &mobilityMocks.MobilityManagementPolicyMock{}
	mobilityMock.RegisterErrorMocks()
	return mockClient, mobilityMock
}

func loadUnitTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/unit/" + filename)
	if err != nil {
		panic("failed to load unit test config " + filename + ": " + err.Error())
	}
	return config
}

// checkRefCalls verifies the total number of includedGroups $ref POST and DELETE requests so far
func checkRefCalls(posts, deletes int) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		if got := mobilityMocks.RefPostCount(); got != posts {
			return fmt.Errorf("expected %d includedGroups $ref POST requests, got %d", posts, got)
		}
		if got := mobilityMocks.RefDeleteCount(); got != deletes {
			return fmt.Errorf("expected %d includedGroups $ref DELETE requests, got %d", deletes, got)
		}
		return nil
	}
}

// checkPatchCount verifies the total number of PATCH requests so far
func checkPatchCount(expected int) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		if got := len(mobilityMocks.PatchBodies()); got != expected {
			return fmt.Errorf("expected %d PATCH requests, got %d: %v", expected, got, mobilityMocks.PatchBodies())
		}
		return nil
	}
}

// Minimal: adopting the policy without configured settings sends no write and keeps the tenant values
func TestUnitResourceMobilityManagementPolicy_01_Minimal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, mobilityMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer mobilityMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					checkPatchCount(0),
					checkRefCalls(0, 0),
					check.That(resourceType+".minimal").Key("id").HasValue("device/0000000a-0000-0000-c000-000000000000"),
					check.That(resourceType+".minimal").Key("policy_type").HasValue("device"),
					check.That(resourceType+".minimal").Key("policy_id").HasValue("0000000a-0000-0000-c000-000000000000"),
					check.That(resourceType+".minimal").Key("display_name").HasValue("Microsoft Intune"),
					check.That(resourceType+".minimal").Key("is_valid").HasValue("true"),
					check.That(resourceType+".minimal").Key("applies_to").HasValue("none"),
					check.That(resourceType+".minimal").Key("included_groups.#").DoesNotExist(),
					check.That(resourceType+".minimal").Key("compliance_url").HasValue("https://portal.manage.microsoft.com/?portalAction=Compliance"),
					check.That(resourceType+".minimal").Key("discovery_url").HasValue("https://enrollment.manage.microsoft.com/enrollmentserver/discovery.svc"),
					check.That(resourceType+".minimal").Key("terms_of_use_url").HasValue("https://portal.manage.microsoft.com/TermsofUse.aspx"),
					check.That(resourceType+".minimal").Key("is_mdm_enrollment_during_registration_disabled").HasValue("false"),
				),
			},
			{
				ResourceName:            resourceType + ".minimal",
				ImportState:             true,
				ImportStateId:           "device/0000000a-0000-0000-c000-000000000000",
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// Maximal: every MDM setting and a selected scope; appliesTo selected is applied through group references only
func TestUnitResourceMobilityManagementPolicy_02_Maximal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, mobilityMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer mobilityMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					checkPatchCount(1),
					checkRefCalls(2, 0),
					check.That(resourceType+".maximal").Key("applies_to").HasValue("selected"),
					check.That(resourceType+".maximal").Key("included_groups.#").HasValue("2"),
					check.That(resourceType+".maximal").Key("included_groups.*").ContainsTypeSetElement("11111111-1111-1111-1111-111111111111"),
					check.That(resourceType+".maximal").Key("included_groups.*").ContainsTypeSetElement("22222222-2222-2222-2222-222222222222"),
					check.That(resourceType+".maximal").Key("compliance_url").HasValue("https://portal.manage.microsoft.com/?portalAction=Compliance&unit-test=maximal"),
					check.That(resourceType+".maximal").Key("terms_of_use_url").HasValue("https://portal.manage.microsoft.com/TermsofUse.aspx?unit-test=maximal"),
					check.That(resourceType+".maximal").Key("is_mdm_enrollment_during_registration_disabled").HasValue("true"),
				),
			},
			{
				ResourceName:            resourceType + ".maximal",
				ImportState:             true,
				ImportStateId:           "device/0000000a-0000-0000-c000-000000000000",
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// Scope lifecycle: group changes add and remove only the differing references, then none and all are PATCHed on their own
func TestUnitResourceMobilityManagementPolicy_03_ScopeLifecycle(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, mobilityMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer mobilityMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					checkRefCalls(2, 0),
					check.That(resourceType+".maximal").Key("included_groups.#").HasValue("2"),
				),
			},
			{
				Config: loadUnitTestTerraform("resource_groups_changed.tf"),
				Check: resource.ComposeTestCheckFunc(
					checkPatchCount(1),
					checkRefCalls(3, 1),
					check.That(resourceType+".maximal").Key("applies_to").HasValue("selected"),
					check.That(resourceType+".maximal").Key("included_groups.#").HasValue("2"),
					check.That(resourceType+".maximal").Key("included_groups.*").ContainsTypeSetElement("22222222-2222-2222-2222-222222222222"),
					check.That(resourceType+".maximal").Key("included_groups.*").ContainsTypeSetElement("33333333-3333-3333-3333-333333333333"),
				),
			},
			{
				Config: loadUnitTestTerraform("resource_scope_none.tf"),
				Check: resource.ComposeTestCheckFunc(
					checkPatchCount(3),
					checkRefCalls(3, 1),
					check.That(resourceType+".maximal").Key("applies_to").HasValue("none"),
					check.That(resourceType+".maximal").Key("included_groups.#").DoesNotExist(),
					check.That(resourceType+".maximal").Key("compliance_url").HasValue("https://portal.manage.microsoft.com/?portalAction=Compliance&unit-test=none"),
					check.That(resourceType+".maximal").Key("is_mdm_enrollment_during_registration_disabled").HasValue("false"),
				),
			},
			{
				Config: loadUnitTestTerraform("resource_scope_all.tf"),
				Check: resource.ComposeTestCheckFunc(
					checkPatchCount(4),
					check.That(resourceType+".maximal").Key("applies_to").HasValue("all"),
					check.That(resourceType+".maximal").Key("included_groups.#").DoesNotExist(),
					check.That(resourceType+".maximal").Key("compliance_url").HasValue("https://portal.manage.microsoft.com/?portalAction=Compliance&unit-test=none"),
				),
			},
			{
				Config: loadUnitTestTerraform("resource_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					checkRefCalls(5, 1),
					check.That(resourceType+".maximal").Key("applies_to").HasValue("selected"),
					check.That(resourceType+".maximal").Key("included_groups.#").HasValue("2"),
				),
			},
		},
	})
}

// App policy: the MAM user scope uses the mobileAppManagementPolicies collection and has no MDM-only setting
func TestUnitResourceMobilityManagementPolicy_04_AppPolicy(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, mobilityMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer mobilityMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_app.tf"),
				Check: resource.ComposeTestCheckFunc(
					checkPatchCount(1),
					checkRefCalls(1, 0),
					check.That(resourceType+".app").Key("id").HasValue("app/0000000a-0000-0000-c000-000000000000"),
					check.That(resourceType+".app").Key("applies_to").HasValue("selected"),
					check.That(resourceType+".app").Key("included_groups.#").HasValue("1"),
					check.That(resourceType+".app").Key("discovery_url").HasValue("https://wip.mam.manage.microsoft.com/Enroll"),
					check.That(resourceType+".app").Key("is_mdm_enrollment_during_registration_disabled").DoesNotExist(),
				),
			},
			{
				ResourceName:            resourceType + ".app",
				ImportState:             true,
				ImportStateId:           "app/0000000a-0000-0000-c000-000000000000",
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// Validation: inconsistent scope settings and MDM-only settings on app policies are rejected at plan time
func TestUnitResourceMobilityManagementPolicy_05_ValidationErrors(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, mobilityMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer mobilityMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadUnitTestTerraform("resource_invalid_groups_with_all.tf"),
				ExpectError: regexp.MustCompile(`included_groups\s+can\s+only\s+be\s+set\s+when\s+applies_to\s+is\s+"selected"`),
			},
			{
				Config:      loadUnitTestTerraform("resource_invalid_selected_without_groups.tf"),
				ExpectError: regexp.MustCompile(`applies_to\s+"selected"\s+requires\s+at\s+least\s+one\s+group`),
			},
			{
				Config:      loadUnitTestTerraform("resource_invalid_mdm_setting_on_app.tf"),
				ExpectError: regexp.MustCompile(`only\s+valid\s+when\s+policy_type\s+is\s+"device"`),
			},
			{
				Config:      loadUnitTestTerraform("resource_invalid_policy_type.tf"),
				ExpectError: regexp.MustCompile(`value\s+must\s+be\s+one\s+of`),
			},
		},
	})
}

// Not found: adopting a policy that does not exist fails with a 404
func TestUnitResourceMobilityManagementPolicy_06_NotFound(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, mobilityMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer mobilityMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadUnitTestTerraform("resource_not_found.tf"),
				ExpectError: regexp.MustCompile(`(?i)not found|does not exist|404`),
			},
		},
	})
}

// Forbidden: a 403 on read surfaces as a permission error that names the required permissions
func TestUnitResourceMobilityManagementPolicy_07_Forbidden(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, mobilityMock := setupErrorMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer mobilityMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadUnitTestTerraform("resource_minimal.tf"),
				ExpectError: regexp.MustCompile(`(?i)Policy\.Read\.All|forbidden|403|permission`),
			},
		},
	})
}

// Bad request: a 400 on PATCH surfaces as a Terraform error
func TestUnitResourceMobilityManagementPolicy_08_BadRequest(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	mobilityMock := &mobilityMocks.MobilityManagementPolicyMock{}
	mobilityMock.RegisterBadRequestOnWriteMocks()
	defer httpmock.DeactivateAndReset()
	defer mobilityMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadUnitTestTerraform("resource_maximal.tf"),
				ExpectError: regexp.MustCompile(`Invalid request|Bad Request|400`),
			},
		},
	})
}

// Import: an ID that is not {policy_type}/{policy_id} is rejected
func TestUnitResourceMobilityManagementPolicy_09_ImportInvalidID(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, mobilityMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer mobilityMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:        loadUnitTestTerraform("resource_minimal.tf"),
				ResourceName:  resourceType + ".minimal",
				ImportState:   true,
				ImportStateId: "0000000a-0000-0000-c000-000000000000",
				ExpectError:   regexp.MustCompile(`expected\s+an\s+ID\s+in\s+the\s+format\s+\{policy_type\}/\{policy_id\}`),
			},
		},
	})
}
