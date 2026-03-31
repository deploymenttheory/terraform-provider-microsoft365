package graphBetaWindowsUpdatesComplianceChanges_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaWindowsUpdatesComplianceChanges "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/datasources/windows_updates/graph_beta/compliance_changes"
	complianceChangesMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/datasources/windows_updates/graph_beta/compliance_changes/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

var (
	dataSourceType = "data." + graphBetaWindowsUpdatesComplianceChanges.DataSourceName
)

func loadUnitTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/unit/" + filename)
	if err != nil {
		panic("failed to load unit test config " + filename + ": " + err.Error())
	}
	return config
}

func setupMockEnvironment() (*mocks.Mocks, *complianceChangesMocks.ComplianceChangesMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	complianceChangesMock := &complianceChangesMocks.ComplianceChangesMock{}
	complianceChangesMock.RegisterMocks()
	return mockClient, complianceChangesMock
}

func setupErrorMockEnvironment() (*mocks.Mocks, *complianceChangesMocks.ComplianceChangesMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	complianceChangesMock := &complianceChangesMocks.ComplianceChangesMock{}
	complianceChangesMock.RegisterErrorMocks()
	return mockClient, complianceChangesMock
}

func TestUnitDatasourceComplianceChanges_01_Basic(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, complianceChangesMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer complianceChangesMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("01_basic.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".test").Key("update_policy_id").HasValue("d7a89208-17c5-4daf-a164-ce176b00e4ef"),
					check.That(dataSourceType+".test").Key("compliance_changes.#").HasValue("2"),
					check.That(dataSourceType+".test").Key("compliance_changes.0.id").Exists(),
					check.That(dataSourceType+".test").Key("compliance_changes.0.created_date_time").Exists(),
					check.That(dataSourceType+".test").Key("compliance_changes.0.is_revoked").HasValue("false"),
					check.That(dataSourceType+".test").Key("compliance_changes.0.content.catalog_entry_id").Exists(),
					check.That(dataSourceType+".test").Key("compliance_changes.0.content.catalog_entry_type").HasValue("driverUpdate"),
					check.That(dataSourceType+".test").Key("compliance_changes.1.is_revoked").HasValue("true"),
					check.That(dataSourceType+".test").Key("compliance_changes.1.revoked_date_time").Exists(),
				),
			},
		},
	})
}

func TestUnitDatasourceComplianceChanges_02_Error(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, complianceChangesMock := setupErrorMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer complianceChangesMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadUnitTestTerraform("01_basic.tf"),
				ExpectError: regexp.MustCompile("Forbidden|403|Insufficient privileges"),
			},
		},
	})
}
