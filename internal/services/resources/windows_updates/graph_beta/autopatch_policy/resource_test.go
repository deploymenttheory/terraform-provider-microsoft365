package graphBetaWindowsUpdatesAutopatchPolicy_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaWindowsAutopatchPolicy "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/windows_updates/graph_beta/autopatch_policy"
	autopatchPolicyMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/windows_updates/graph_beta/autopatch_policy/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

func setupMockEnvironment() (*mocks.Mocks, *autopatchPolicyMocks.WindowsAutopatchPolicyMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	autopatchMock := &autopatchPolicyMocks.WindowsAutopatchPolicyMock{}
	autopatchMock.RegisterMocks()
	return mockClient, autopatchMock
}

func loadUnitTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/unit/" + filename)
	if err != nil {
		panic("failed to load unit test config " + filename + ": " + err.Error())
	}
	return config
}

var resourceType = graphBetaWindowsAutopatchPolicy.ResourceName

// WAP001: Minimal policy - no approval rules
func TestUnitResourceWindowsAutopatchPolicy_01_WAP001_Minimal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, autopatchMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer autopatchMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_wap001-minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".wap001_minimal").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".wap001_minimal").Key("display_name").HasValue("WAP001: Minimal Windows Autopatch Policy-v1.0"),
					check.That(resourceType+".wap001_minimal").Key("description").HasValue("Minimal policy with no approval rules"),
					check.That(resourceType+".wap001_minimal").Key("created_date_time").IsNotEmpty(),
					check.That(resourceType+".wap001_minimal").Key("last_modified_date_time").IsNotEmpty(),
					check.That(resourceType+".wap001_minimal").Key("approval_rules.#").HasValue("0"),
				),
			},
			{
				ResourceName:      resourceType + ".wap001_minimal",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// WAP002: Policy with approval rules
func TestUnitResourceWindowsAutopatchPolicy_02_WAP002_ApprovalRules(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, autopatchMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer autopatchMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_wap002-approval-rules.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".wap002_approval_rules").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".wap002_approval_rules").Key("display_name").HasValue("WAP002: Windows Autopatch Policy With Approval Rules-v1.0"),
					check.That(resourceType+".wap002_approval_rules").Key("description").HasValue("Policy with all approval rule types"),
					check.That(resourceType+".wap002_approval_rules").Key("created_date_time").IsNotEmpty(),
					check.That(resourceType+".wap002_approval_rules").Key("last_modified_date_time").IsNotEmpty(),
					check.That(resourceType+".wap002_approval_rules").Key("approval_rules.#").HasValue("3"),
				),
			},
			{
				ResourceName:      resourceType + ".wap002_approval_rules",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// WAP003: Lifecycle - minimal to with-rules
func TestUnitResourceWindowsAutopatchPolicy_03_WAP003_Lifecycle(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, autopatchMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer autopatchMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_wap003-lifecycle-step1.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".wap003_lifecycle").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".wap003_lifecycle").Key("display_name").HasValue("WAP003: Windows Autopatch Policy Lifecycle-v1.0"),
					check.That(resourceType+".wap003_lifecycle").Key("description").HasValue("Lifecycle step 1: no approval rules"),
					check.That(resourceType+".wap003_lifecycle").Key("approval_rules.#").HasValue("0"),
				),
			},
			{
				Config: loadUnitTestTerraform("resource_wap003-lifecycle-step2.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".wap003_lifecycle").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".wap003_lifecycle").Key("display_name").HasValue("WAP003: Windows Autopatch Policy Lifecycle-v1.0"),
					check.That(resourceType+".wap003_lifecycle").Key("description").HasValue("Lifecycle step 2: with approval rules added"),
					check.That(resourceType+".wap003_lifecycle").Key("approval_rules.#").HasValue("2"),
				),
			},
		},
	})
}
