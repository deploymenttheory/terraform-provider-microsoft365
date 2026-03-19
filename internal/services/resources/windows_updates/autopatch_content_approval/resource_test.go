package graphBetaWindowsUpdatesAutopatchContentApproval_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	contentApprovalMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/windows_updates/autopatch_content_approval/mocks"
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

func setupMockEnvironment() (*mocks.Mocks, *contentApprovalMocks.WindowsUpdateContentApprovalMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	contentApprovalMock := &contentApprovalMocks.WindowsUpdateContentApprovalMock{}
	contentApprovalMock.RegisterMocks()
	return mockClient, contentApprovalMock
}

func setupErrorMockEnvironment() (*mocks.Mocks, *contentApprovalMocks.WindowsUpdateContentApprovalMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	contentApprovalMock := &contentApprovalMocks.WindowsUpdateContentApprovalMock{}
	contentApprovalMock.RegisterErrorMocks()
	return mockClient, contentApprovalMock
}

func TestUnitResourceWindowsUpdateContentApproval_01_FeatureUpdateApproval(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, contentApprovalMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer contentApprovalMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("01_feature_update_approval.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("id").HasValue("bba2a340-1e32-b5ed-186e-678e16033319"),
					check.That(resourceType+".test").Key("update_policy_id").HasValue("983f03cd-03cd-983f-cd03-3f98cd033f98"),
					check.That(resourceType+".test").Key("catalog_entry_id").HasValue("c1dec151-c151-c1de-51c1-dec151c1dec1"),
					check.That(resourceType+".test").Key("catalog_entry_type").HasValue("featureUpdate"),
					check.That(resourceType+".test").Key("is_revoked").HasValue("false"),
					check.That(resourceType+".test").Key("created_date_time").Exists(),
					check.That(resourceType+".test").Key("deployment_settings.schedule.start_date_time").HasValue("2026-03-10T00:00:00Z"),
					check.That(resourceType+".test").Key("deployment_settings.schedule.gradual_rollout.end_date_time").HasValue("2026-03-20T00:00:00Z"),
				),
			},
		},
	})
}

func TestUnitResourceWindowsUpdateContentApproval_02_QualityUpdateApproval(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, contentApprovalMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer contentApprovalMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("02_quality_update_approval.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("id").Exists(),
					check.That(resourceType+".test").Key("catalog_entry_type").HasValue("qualityUpdate"),
					check.That(resourceType+".test").Key("deployment_settings.schedule.start_date_time").HasValue("2026-03-11T00:00:00Z"),
				),
			},
		},
	})
}

func TestUnitResourceWindowsUpdateContentApproval_03_MinimalApproval(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, contentApprovalMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer contentApprovalMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("03_minimal_approval.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("id").Exists(),
					check.That(resourceType+".test").Key("update_policy_id").HasValue("983f03cd-03cd-983f-cd03-3f98cd033f98"),
					check.That(resourceType+".test").Key("catalog_entry_id").HasValue("c1dec151-c151-c1de-51c1-dec151c1dec1"),
					check.That(resourceType+".test").Key("catalog_entry_type").HasValue("featureUpdate"),
				),
			},
		},
	})
}

func TestUnitResourceWindowsUpdateContentApproval_04_Update(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, contentApprovalMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer contentApprovalMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("01_feature_update_approval.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType + ".test").Key("is_revoked").HasValue("false"),
				),
			},
			{
				Config: loadUnitTestTerraform("01_feature_update_approval.tf") + `
resource "microsoft365_graph_beta_windows_updates_autopatch_content_approval" "test_revoked" {
  update_policy_id    = "983f03cd-03cd-983f-cd03-3f98cd033f98"
  catalog_entry_id    = "c1dec151-c151-c1de-51c1-dec151c1dec1"
  catalog_entry_type  = "featureUpdate"
  is_revoked          = true
}
`,
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_revoked").Key("is_revoked").HasValue("true"),
					check.That(resourceType+".test_revoked").Key("revoked_date_time").Exists(),
				),
			},
		},
	})
}

func TestUnitResourceWindowsUpdateContentApproval_05_Error(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, contentApprovalMock := setupErrorMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer contentApprovalMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadUnitTestTerraform("01_feature_update_approval.tf"),
				ExpectError: regexp.MustCompile("Forbidden|403|Insufficient privileges"),
			},
		},
	})
}
