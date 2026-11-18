package graphBetaOAuth2PermissionGrant_test

import (
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	grant "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/applications/graph_beta/oauth2_permission_grant"
	grantmocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/applications/graph_beta/oauth2_permission_grant/mocks"
)

var (
	resourceType = grant.ResourceName
	testResource = grant.OAuth2PermissionGrantTestResource{}
)

func setupMockEnvironment(t *testing.T) *grantmocks.OAuth2PermissionGrantMock {
	t.Helper()
	mocks.SetupUnitTestEnvironment(t)
	httpmock.Activate()
	mocks.NewMocks().AuthMocks.RegisterMocks()
	mock := &grantmocks.OAuth2PermissionGrantMock{}
	mock.CleanupMockState()
	mock.RegisterMocks()
	t.Cleanup(func() { mock.CleanupMockState(); httpmock.DeactivateAndReset() })
	return mock
}

func loadUnitTestTerraform(t *testing.T, file string) string {
	t.Helper()
	config, err := helpers.ParseHCLFile("tests/terraform/unit/" + file)
	require.NoError(t, err)
	return config
}

func TestUnitResourceOAuth2PermissionGrant_Lifecycle(t *testing.T) {
	for _, fixture := range []string{"resource_minimal.tf", "resource_principal.tf"} {
		t.Run(fixture, func(t *testing.T) {
			mock := setupMockEnvironment(t)
			config := loadUnitTestTerraform(t, fixture)
			updated := strings.Replace(config, `["Probe.Read"]`, `["Probe.Write", "Probe.Read"]`, 1)
			address := resourceType + ".test"
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{Config: config, Check: resource.TestCheckResourceAttr(address, "id", grantmocks.GrantID)},
					{ResourceName: address, ImportState: true, ImportStateVerify: true},
					{Config: updated, Check: resource.TestCheckResourceAttr(address, "scopes.#", "2")},
					{Config: config, Check: resource.TestCheckResourceAttr(address, "scopes.#", "1")},
				},
			})
			require.Len(t, mock.Patches(), 2, "scope changes must PATCH without replacing the grant")
		})
	}
}

func TestUnitResourceOAuth2PermissionGrant_GeneratedImport(t *testing.T) {
	mock := setupMockEnvironment(t)
	require.NoError(t, mock.SeedFixture())
	resource.UnitTest(t, resource.TestCase{
		TerraformVersionChecks:   []tfversion.TerraformVersionCheck{tfversion.SkipBelow(tfversion.Version1_5_0)},
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:                 loadUnitTestTerraform(t, "import.tf"),
				ResourceName:           resourceType + ".test",
				ImportState:            true,
				ImportStateId:          grantmocks.GrantID,
				GenerateConfig:         true,
				ImportStateKind:        resource.ImportBlockWithID,
				ImportStateConfigExact: true,
				ImportPlanChecks: resource.ImportPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectResourceAction(resourceType+".test", plancheck.ResourceActionNoop)},
				},
			},
		},
	})
}

func TestUnitResourceOAuth2PermissionGrant_Validation(t *testing.T) {
	base := `resource "microsoft365_graph_beta_applications_oauth2_permission_grant" "test" {
client_id="11111111-1111-4111-8111-111111111111"
resource_id="22222222-2222-4222-8222-222222222222"
consent_type="AllPrincipals"
scopes=["Probe.Read"]
}`
	for _, tc := range []struct{ name, config, pattern string }{
		{"principal required", strings.Replace(base, `"AllPrincipals"`, `"Principal"`, 1), "Attribute Required"},
		{"principal forbidden", strings.Replace(base, "scopes=", `principal_id="33333333-3333-4333-8333-333333333333"`+"\nscopes=", 1), "Field Must Be Empty"},
		{"single scope token", strings.Replace(base, `"Probe.Read"`, `"Probe.Read Probe.Write"`, 1), "scope token"},
		{"combined length", strings.Replace(base, `"Probe.Read"`, `"`+strings.Repeat("a", 3851)+`"`, 1), "Scope string exceeds Graph limit"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			setupMockEnvironment(t)
			resource.UnitTest(
				t,
				resource.TestCase{
					ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
					Steps:                    []resource.TestStep{{Config: tc.config, ExpectError: regexp.MustCompile(tc.pattern)}},
				},
			)
		})
	}
}
