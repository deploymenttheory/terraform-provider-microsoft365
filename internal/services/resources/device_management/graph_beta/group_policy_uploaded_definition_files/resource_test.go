package graphBetaGroupPolicyUploadedDefinitionFiles_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaGroupPolicyUploadedDefinitionFiles "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/group_policy_uploaded_definition_files"
	gpudfMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/group_policy_uploaded_definition_files/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

var (
	// Resource type name from the resource package
	resourceType = graphBetaGroupPolicyUploadedDefinitionFiles.ResourceName

	// testResource is the test resource implementation for group policy uploaded definition files
	testResource = graphBetaGroupPolicyUploadedDefinitionFiles.GroupPolicyUploadedDefinitionFilesTestResource{}
)

// Helper function to load test configs from unit directory
func loadUnitTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/unit/" + filename)
	if err != nil {
		panic("failed to load unit config " + filename + ": " + err.Error())
	}
	return config
}

func setupMockEnvironment() (*mocks.Mocks, *gpudfMocks.GroupPolicyUploadedDefinitionFilesMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	gpudfMock := &gpudfMocks.GroupPolicyUploadedDefinitionFilesMock{}
	gpudfMock.RegisterMocks()
	return mockClient, gpudfMock
}

func setupErrorMockEnvironment() (*mocks.Mocks, *gpudfMocks.GroupPolicyUploadedDefinitionFilesMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	gpudfMock := &gpudfMocks.GroupPolicyUploadedDefinitionFilesMock{}
	gpudfMock.RegisterErrorMocks()
	return mockClient, gpudfMock
}

func TestUnitResourceGroupPolicyUploadedDefinitionFiles_01_Schema(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, gpudfMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer gpudfMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_group_policy_uploaded_definition_files_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".minimal").Key("file_name").HasValue("unit-test-firefox.admx"),
					check.That(resourceType+".minimal").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".minimal").Key("default_language_code").HasValue("en-US"),
					check.That(resourceType+".minimal").Key("group_policy_uploaded_language_files.#").HasValue("1"),
				),
			},
		},
	})
}

func TestUnitResourceGroupPolicyUploadedDefinitionFiles_02_LanguageFiles(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, gpudfMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer gpudfMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_group_policy_uploaded_definition_files_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".maximal").Key("id").Exists(),
					check.That(resourceType+".maximal").Key("file_name").HasValue("unit-test-msedge.admx"),
					check.That(resourceType+".maximal").Key("default_language_code").HasValue("en-US"),
					check.That(resourceType+".maximal").Key("group_policy_uploaded_language_files.#").HasValue("2"),
				),
			},
		},
	})
}

func TestUnitResourceGroupPolicyUploadedDefinitionFiles_03_ErrorHandling(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, gpudfMock := setupErrorMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer gpudfMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadUnitTestTerraform("resource_group_policy_uploaded_definition_files_minimal.tf"),
				ExpectError: regexp.MustCompile("Invalid ADMX file format"),
			},
		},
	})
}
