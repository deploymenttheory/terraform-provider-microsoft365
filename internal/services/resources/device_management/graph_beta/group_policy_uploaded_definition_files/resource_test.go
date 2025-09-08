package graphBetaGroupPolicyUploadedDefinitionFiles_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	gpudfMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/group_policy_uploaded_definition_files/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

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

func testCheckExists(resourceName string) resource.TestCheckFunc {
	return resource.TestCheckResourceAttrSet(resourceName, "id")
}

func TestGroupPolicyUploadedDefinitionFilesResource_Schema(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, gpudfMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer gpudfMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_group_policy_uploaded_definition_files.minimal", "file_name", "unit-test-firefox.admx"),
					resource.TestMatchResourceAttr("microsoft365_graph_beta_device_management_group_policy_uploaded_definition_files.minimal", "id", regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_group_policy_uploaded_definition_files.minimal", "default_language_code", "en-US"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_group_policy_uploaded_definition_files.minimal", "group_policy_uploaded_language_files.#", "1"),
				),
			},
		},
	})
}

func TestGroupPolicyUploadedDefinitionFilesResource_LanguageFiles(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, gpudfMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer gpudfMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMaximal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_group_policy_uploaded_definition_files.maximal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_group_policy_uploaded_definition_files.maximal", "file_name", "unit-test-msedge.admx"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_group_policy_uploaded_definition_files.maximal", "default_language_code", "en-US"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_group_policy_uploaded_definition_files.maximal", "group_policy_uploaded_language_files.#", "2"),
				),
			},
		},
	})
}

func TestGroupPolicyUploadedDefinitionFilesResource_ErrorHandling(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, gpudfMock := setupErrorMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer gpudfMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testConfigMinimal(),
				ExpectError: regexp.MustCompile("Invalid ADMX file format"),
			},
		},
	})
}

func testConfigMinimal() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_group_policy_uploaded_definition_files_minimal.tf")
	if err != nil {
		panic("failed to load minimal config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigMaximal() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_group_policy_uploaded_definition_files_maximal.tf")
	if err != nil {
		panic("failed to load maximal config: " + err.Error())
	}
	return unitTestConfig
}
