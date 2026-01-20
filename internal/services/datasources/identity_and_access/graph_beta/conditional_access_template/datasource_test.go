package graphBetaConditionalAccessTemplate_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaConditionalAccessTemplate "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/datasources/identity_and_access/graph_beta/conditional_access_template"
	conditionalAccessTemplateMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/datasources/identity_and_access/graph_beta/conditional_access_template/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

var (
	dataSourceType = graphBetaConditionalAccessTemplate.DataSourceName
)

func setupMockEnvironment() (*mocks.Mocks, *conditionalAccessTemplateMocks.ConditionalAccessTemplateMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	catMock := &conditionalAccessTemplateMocks.ConditionalAccessTemplateMock{}
	catMock.RegisterMocks()
	return mockClient, catMock
}

func setupErrorMockEnvironment() (*mocks.Mocks, *conditionalAccessTemplateMocks.ConditionalAccessTemplateMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	catMock := &conditionalAccessTemplateMocks.ConditionalAccessTemplateMock{}
	catMock.RegisterErrorMocks()
	return mockClient, catMock
}

func TestConditionalAccessTemplateDataSource_ByTemplateId(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, catMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer catMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("01_by_template_id.tf"),
				Check: resource.ComposeTestCheckFunc(
					// Top level attributes
					check.That("data."+dataSourceType+".by_template_id").Key("template_id").HasValue("c7503427-338e-4c5e-902d-abe252abfb43"),
					check.That("data."+dataSourceType+".by_template_id").Key("name").HasValue("Require multifactor authentication for admins"),
					check.That("data."+dataSourceType+".by_template_id").Key("description").HasValue("Require multifactor authentication for privileged administrative accounts to reduce risk of compromise. This policy will target the same roles as security defaults."),
					check.That("data."+dataSourceType+".by_template_id").Key("scenarios.#").HasValue("3"),
					check.That("data."+dataSourceType+".by_template_id").Key("id").IsSet(),

					// Details - Conditions
					check.That("data."+dataSourceType+".by_template_id").Key("details.conditions.client_app_types.#").HasValue("1"),
					check.That("data."+dataSourceType+".by_template_id").Key("details.conditions.user_risk_levels.#").HasValue("0"),
					check.That("data."+dataSourceType+".by_template_id").Key("details.conditions.sign_in_risk_levels.#").HasValue("0"),
					check.That("data."+dataSourceType+".by_template_id").Key("details.conditions.service_principal_risk_levels.#").HasValue("0"),

					// Details - Conditions - Applications
					check.That("data."+dataSourceType+".by_template_id").Key("details.conditions.applications.include_applications.#").HasValue("1"),
					check.That("data."+dataSourceType+".by_template_id").Key("details.conditions.applications.exclude_applications.#").HasValue("0"),
					check.That("data."+dataSourceType+".by_template_id").Key("details.conditions.applications.include_user_actions.#").HasValue("0"),
					check.That("data."+dataSourceType+".by_template_id").Key("details.conditions.applications.include_authentication_context_class_references.#").HasValue("0"),

					// Details - Conditions - Users
					check.That("data."+dataSourceType+".by_template_id").Key("details.conditions.users.include_users.#").HasValue("0"),
					check.That("data."+dataSourceType+".by_template_id").Key("details.conditions.users.exclude_users.#").HasValue("1"),
					check.That("data."+dataSourceType+".by_template_id").Key("details.conditions.users.exclude_users.0").HasValue("Current administrator will be excluded"),
					check.That("data."+dataSourceType+".by_template_id").Key("details.conditions.users.include_groups.#").HasValue("0"),
					check.That("data."+dataSourceType+".by_template_id").Key("details.conditions.users.exclude_groups.#").HasValue("0"),
					check.That("data."+dataSourceType+".by_template_id").Key("details.conditions.users.include_roles.#").HasValue("14"),
					check.That("data."+dataSourceType+".by_template_id").Key("details.conditions.users.exclude_roles.#").HasValue("0"),

					// Details - Grant Controls
					check.That("data."+dataSourceType+".by_template_id").Key("details.grant_controls.operator").HasValue("OR"),
					check.That("data."+dataSourceType+".by_template_id").Key("details.grant_controls.built_in_controls.#").HasValue("1"),
					check.That("data."+dataSourceType+".by_template_id").Key("details.grant_controls.custom_authentication_factors.#").HasValue("0"),
					check.That("data."+dataSourceType+".by_template_id").Key("details.grant_controls.terms_of_use.#").HasValue("0"),
				),
			},
		},
	})
}

func TestConditionalAccessTemplateDataSource_ByName(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, catMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer catMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("02_by_name.tf"),
				Check: resource.ComposeTestCheckFunc(
					// Top level attributes
					check.That("data."+dataSourceType+".by_name").Key("template_id").HasValue("c7503427-338e-4c5e-902d-abe252abfb43"),
					check.That("data."+dataSourceType+".by_name").Key("name").HasValue("Require multifactor authentication for admins"),
					check.That("data."+dataSourceType+".by_name").Key("description").HasValue("Require multifactor authentication for privileged administrative accounts to reduce risk of compromise. This policy will target the same roles as security defaults."),
					check.That("data."+dataSourceType+".by_name").Key("scenarios.#").HasValue("3"),
					check.That("data."+dataSourceType+".by_name").Key("id").IsSet(),

					// Details - Conditions
					check.That("data."+dataSourceType+".by_name").Key("details.conditions.client_app_types.#").HasValue("1"),
					check.That("data."+dataSourceType+".by_name").Key("details.conditions.user_risk_levels.#").HasValue("0"),
					check.That("data."+dataSourceType+".by_name").Key("details.conditions.sign_in_risk_levels.#").HasValue("0"),
					check.That("data."+dataSourceType+".by_name").Key("details.conditions.service_principal_risk_levels.#").HasValue("0"),

					// Details - Conditions - Applications
					check.That("data."+dataSourceType+".by_name").Key("details.conditions.applications.include_applications.#").HasValue("1"),
					check.That("data."+dataSourceType+".by_name").Key("details.conditions.applications.exclude_applications.#").HasValue("0"),
					check.That("data."+dataSourceType+".by_name").Key("details.conditions.applications.include_user_actions.#").HasValue("0"),
					check.That("data."+dataSourceType+".by_name").Key("details.conditions.applications.include_authentication_context_class_references.#").HasValue("0"),

					// Details - Conditions - Users
					check.That("data."+dataSourceType+".by_name").Key("details.conditions.users.include_users.#").HasValue("0"),
					check.That("data."+dataSourceType+".by_name").Key("details.conditions.users.exclude_users.#").HasValue("1"),
					check.That("data."+dataSourceType+".by_name").Key("details.conditions.users.exclude_users.0").HasValue("Current administrator will be excluded"),
					check.That("data."+dataSourceType+".by_name").Key("details.conditions.users.include_groups.#").HasValue("0"),
					check.That("data."+dataSourceType+".by_name").Key("details.conditions.users.exclude_groups.#").HasValue("0"),
					check.That("data."+dataSourceType+".by_name").Key("details.conditions.users.include_roles.#").HasValue("14"),
					check.That("data."+dataSourceType+".by_name").Key("details.conditions.users.exclude_roles.#").HasValue("0"),

					// Details - Grant Controls
					check.That("data."+dataSourceType+".by_name").Key("details.grant_controls.operator").HasValue("OR"),
					check.That("data."+dataSourceType+".by_name").Key("details.grant_controls.built_in_controls.#").HasValue("1"),
					check.That("data."+dataSourceType+".by_name").Key("details.grant_controls.custom_authentication_factors.#").HasValue("0"),
					check.That("data."+dataSourceType+".by_name").Key("details.grant_controls.terms_of_use.#").HasValue("0"),
				),
			},
		},
	})
}

func TestConditionalAccessTemplateDataSource_ValidationError(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, catMock := setupErrorMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer catMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadUnitTestTerraform("01_by_template_id.tf"),
				ExpectError: regexp.MustCompile("Failed to Retrieve Templates"),
			},
		},
	})
}

func TestConditionalAccessTemplateDataSource_FuzzyMatchSuggestion(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, catMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer catMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadUnitTestTerraform("03_invalid_name_fuzzy.tf"),
				ExpectError: regexp.MustCompile("(?s)Invalid Template Name.*No conditional access template found with name: Require MFA for admin.*Did you mean one of these.*Require multifactor authentication for admins"),
			},
		},
	})
}

// Helper function to load test configs from unit directory
func loadUnitTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/unit/" + filename)
	if err != nil {
		panic("failed to load unit test config " + filename + ": " + err.Error())
	}
	return config
}
