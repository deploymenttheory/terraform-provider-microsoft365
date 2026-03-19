package graphBetaWindowsUpdatesAutopatchPolicyApproval_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	policyApprovalMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/windows_updates/autopatch_policy_approval/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

func loadUnitTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/unit/" + filename)
	if err != nil {
		panic("failed to load unit test config " + filename + ": " + err.Error())
	}
	return config
}

func setupMockEnvironment() (*mocks.Mocks, *policyApprovalMocks.WindowsUpdatePolicyApprovalMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	approvalMock := &policyApprovalMocks.WindowsUpdatePolicyApprovalMock{}
	approvalMock.RegisterMocks()
	return mockClient, approvalMock
}

func setupErrorMockEnvironment() (*mocks.Mocks, *policyApprovalMocks.WindowsUpdatePolicyApprovalMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	approvalMock := &policyApprovalMocks.WindowsUpdatePolicyApprovalMock{}
	approvalMock.RegisterErrorMocks()
	return mockClient, approvalMock
}

func TestUnitResourceWindowsUpdatePolicyApproval_01_Approved(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, approvalMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer approvalMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("01_approved.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("id").HasValue("a1b2c3d4-1234-5678-abcd-a1b2c3d4e5f6"),
					check.That(resourceType+".test").Key("policy_id").HasValue("983f03cd-03cd-983f-cd03-3f98cd033f98"),
					check.That(resourceType+".test").Key("catalog_entry_id").HasValue("c1dec151-c151-c1de-51c1-dec151c1dec1"),
					check.That(resourceType+".test").Key("status").HasValue("approved"),
					check.That(resourceType+".test").Key("created_date_time").Exists(),
					check.That(resourceType+".test").Key("last_modified_date_time").Exists(),
				),
			},
		},
	})
}

func TestUnitResourceWindowsUpdatePolicyApproval_02_Suspended(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, approvalMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer approvalMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("02_suspended.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("id").Exists(),
					check.That(resourceType+".test").Key("status").HasValue("suspended"),
				),
			},
		},
	})
}

func TestUnitResourceWindowsUpdatePolicyApproval_03_UpdateStatus(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, approvalMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer approvalMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("01_approved.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("status").HasValue("approved"),
				),
			},
			{
				Config: loadUnitTestTerraform("02_suspended.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("status").HasValue("suspended"),
				),
			},
		},
	})
}

func TestUnitResourceWindowsUpdatePolicyApproval_04_Import(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, approvalMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer approvalMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("01_approved.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("id").Exists(),
				),
			},
			{
				ResourceName:            resourceType + ".test",
				ImportState:             true,
				ImportStateId:           "983f03cd-03cd-983f-cd03-3f98cd033f98/a1b2c3d4-1234-5678-abcd-a1b2c3d4e5f6",
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

func TestUnitResourceWindowsUpdatePolicyApproval_05_Error(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, approvalMock := setupErrorMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer approvalMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadUnitTestTerraform("01_approved.tf"),
				ExpectError: regexp.MustCompile("Forbidden|403|Insufficient privileges"),
			},
		},
	})
}
