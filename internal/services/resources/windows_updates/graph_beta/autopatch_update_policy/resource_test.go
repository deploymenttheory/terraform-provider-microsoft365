package graphBetaWindowsUpdatesAutopatchUpdatePolicy_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaWindowsUpdatesAutopatchUpdatePolicy "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/windows_updates/graph_beta/autopatch_update_policy"
	updatePolicyMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/windows_updates/graph_beta/autopatch_update_policy/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

var resourceType = graphBetaWindowsUpdatesAutopatchUpdatePolicy.ResourceName

func setupMockEnvironment() (*mocks.Mocks, *updatePolicyMocks.WindowsUpdatePolicyMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	updatePolicyMock := &updatePolicyMocks.WindowsUpdatePolicyMock{}
	updatePolicyMock.RegisterMocks()
	return mockClient, updatePolicyMock
}

func setupErrorMockEnvironment() (*mocks.Mocks, *updatePolicyMocks.WindowsUpdatePolicyMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	updatePolicyMock := &updatePolicyMocks.WindowsUpdatePolicyMock{}
	updatePolicyMock.RegisterErrorMocks()
	return mockClient, updatePolicyMock
}

func loadUnitTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/unit/" + filename)
	if err != nil {
		panic("failed to load unit test config " + filename + ": " + err.Error())
	}
	return config
}

func TestUnitResourceWindowsUpdatesUpdatePolicy_01_CreateUpdatePolicy(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, updatePolicyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer updatePolicyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("01_create_update_policy.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test").Key("created_date_time").IsNotEmpty(),
				),
			},
		{
			ResourceName:            resourceType + ".test",
			ImportState:             true,
			ImportStateVerify:       true,
			ImportStateVerifyIgnore: []string{"compliance_changes"},
		},
		},
	})
}

func TestUnitResourceWindowsUpdatesUpdatePolicy_02_UpdatePolicySettings(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, updatePolicyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer updatePolicyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("01_create_update_policy.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test").Key("audience_id").IsNotEmpty(),
					check.That(resourceType+".test").Key("compliance_change_rules.#").HasValue("1"),
				),
			},
			{
				Config: loadUnitTestTerraform("02_update_policy_settings.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test").Key("compliance_change_rules.#").HasValue("1"),
				),
			},
		},
	})
}

func TestUnitResourceWindowsUpdatesUpdatePolicy_03_MinimalPolicy(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, updatePolicyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer updatePolicyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("03_minimal_policy.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test").Key("audience_id").IsNotEmpty(),
				),
			},
		},
	})
}

func TestUnitResourceWindowsUpdatesUpdatePolicy_04_ErrorOnCreate(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, updatePolicyMock := setupErrorMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer updatePolicyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadUnitTestTerraform("01_create_update_policy.tf"),
				ExpectError: regexp.MustCompile("Forbidden|403|Insufficient privileges"),
			},
		},
	})
}
