package graphBetaAuthorizationPolicy_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	policy "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_beta/authorization_policy"
	policymocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_beta/authorization_policy/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
)

func setupMockEnvironment(t *testing.T) *policymocks.AuthorizationPolicyMock {
	t.Helper()
	mocks.SetupUnitTestEnvironment(t)
	httpmock.Activate()
	mocks.NewMocks().AuthMocks.RegisterMocks()
	m := &policymocks.AuthorizationPolicyMock{}
	m.RegisterMocks()
	t.Cleanup(httpmock.DeactivateAndReset)
	t.Cleanup(m.CleanupMockState)
	return m
}

// loadUnitTestTerraform loads a configuration from the unit test fixtures.
func loadUnitTestTerraform(t *testing.T, name string) string {
	t.Helper()
	hcl, err := helpers.ParseHCLFile("tests/terraform/unit/resource_" + name + ".tf")
	require.NoError(t, err)
	return hcl
}

// TestUnitResourceAuthorizationPolicy_01_Lifecycle verifies lifecycle behavior.
func TestUnitResourceAuthorizationPolicy_01_Lifecycle(t *testing.T) {
	m := setupMockEnvironment(t)
	address := policy.ResourceName + ".test"
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: loadUnitTestTerraform(t, "minimal"), Check: resource.ComposeTestCheckFunc(
				check.That(address).Key("id").HasValue("authorizationPolicy"),
				check.That(address).Key("allow_invites_from").HasValue("adminsAndGuestInviters"),
				check.That(address).Key("allowed_to_use_sspr").HasValue("true"),
				check.That(address).Key("default_user_role_permissions.allowed_to_create_apps").HasValue("true"))},
			{Config: loadUnitTestTerraform(t, "complete"), Check: resource.ComposeTestCheckFunc(
				check.That(address).Key("allow_invites_from").HasValue("none"),
				check.That(address).Key("allow_user_consent_for_risky_apps").HasValue("false"),
				check.That(address).Key("guest_user_role_id").HasValue("2af84b1e-32c8-42b7-82bc-daa82404023b"),
				check.That(address).Key("enabled_preview_features.#").HasValue("1"),
				check.That(address).Key("default_user_role_permissions.allowed_to_create_apps").HasValue("false"))},
			{ResourceName: address, ImportState: true, ImportStateId: "authorizationPolicy", ImportStateVerify: true},
			{Config: loadUnitTestTerraform(t, "empty_sets"), Check: resource.ComposeTestCheckFunc(
				check.That(address).Key("enabled_preview_features.#").HasValue("0"),
				check.That(address).Key("permission_grant_policy_ids_assigned_to_default_user_role.#").HasValue("0"),
				check.That(address).Key("default_user_role_permissions.allowed_to_create_security_groups").HasValue("false"))},
			{Config: loadUnitTestTerraform(t, "empty_sets"), PlanOnly: true},
			{Config: loadUnitTestTerraform(t, "empty_sets"), Destroy: true},
		},
	})
	invites, patches := m.Snapshot()
	require.Equal(t, "none", invites)
	require.Equal(t, 3, patches)
	for _, url := range []string{policymocks.URL, policymocks.PatchURL} {
		require.Zero(t, httpmock.GetCallCountInfo()["DELETE "+url])
		require.Zero(t, httpmock.GetCallCountInfo()["POST "+url])
	}
}

// TestUnitResourceAuthorizationPolicy_02_Validation verifies validation behavior.
func TestUnitResourceAuthorizationPolicy_02_Validation(t *testing.T) {
	for _, name := range []string{"invalid_invites", "invalid_role"} {
		t.Run(name, func(t *testing.T) {
			setupMockEnvironment(t)
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
				Steps:                    []resource.TestStep{{Config: loadUnitTestTerraform(t, name), ExpectError: regexp.MustCompile(`value must be one of`)}},
			})
		})
	}
}

// TestUnitResourceAuthorizationPolicy_03_Drift verifies drift behavior.
func TestUnitResourceAuthorizationPolicy_03_Drift(t *testing.T) {
	m := setupMockEnvironment(t)
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: loadUnitTestTerraform(t, "minimal")},
			{PreConfig: func() { m.SetInvites("everyone") }, Config: loadUnitTestTerraform(t, "minimal"), Check: check.That(policy.ResourceName + ".test").Key("allow_invites_from").HasValue("adminsAndGuestInviters")},
		},
	})
}

// TestUnitResourceAuthorizationPolicy_04_Forbidden verifies forbidden behavior.
func TestUnitResourceAuthorizationPolicy_04_Forbidden(t *testing.T) {
	m := setupMockEnvironment(t)
	m.RegisterErrorMocks()
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps:                    []resource.TestStep{{Config: loadUnitTestTerraform(t, "minimal"), ExpectError: regexp.MustCompile(`(?i)(permission|privilege|forbidden|access)`)}},
	})
}

// TestUnitResourceAuthorizationPolicy_05_RequiredBooleans rejects omitted Boolean settings.
func TestUnitResourceAuthorizationPolicy_05_RequiredBooleans(t *testing.T) {
	for _, name := range []string{"missing_bool", "missing_nested_bool"} {
		t.Run(name, func(t *testing.T) {
			setupMockEnvironment(t)
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
				Steps:                    []resource.TestStep{{Config: loadUnitTestTerraform(t, name), ExpectError: regexp.MustCompile(`Missing Configuration for Required Attribute|Missing required argument`)}},
			})
		})
	}

}
