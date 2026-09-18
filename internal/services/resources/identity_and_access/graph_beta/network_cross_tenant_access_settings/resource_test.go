package graphBetaNetworkCrossTenantAccessSettings_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	settings "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_beta/network_cross_tenant_access_settings"
	settingsmocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_beta/network_cross_tenant_access_settings/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
)

var resourceType = settings.ResourceName

func setupMockEnvironment(t *testing.T) *settingsmocks.NetworkCrossTenantAccessSettingsMock {
	t.Helper()
	mocks.SetupUnitTestEnvironment(t)
	httpmock.Activate()
	mocks.NewMocks().AuthMocks.RegisterMocks()
	m := &settingsmocks.NetworkCrossTenantAccessSettingsMock{}
	m.RegisterMocks()
	t.Cleanup(httpmock.DeactivateAndReset)
	t.Cleanup(m.CleanupMockState)
	return m
}

func config(t *testing.T, status string) string {
	t.Helper()
	hcl, err := helpers.ParseHCLFile("tests/terraform/unit/resource_" + status + ".tf")
	require.NoError(t, err)
	return hcl
}

func TestUnitResourceNetworkCrossTenantAccessSettings_01_Lifecycle(t *testing.T) {
	m := setupMockEnvironment(t)
	address := resourceType + ".test"
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: config(t, "disabled"), Check: resource.ComposeTestCheckFunc(
				check.That(address).Key("id").HasValue("crossTenantAccess"),
				check.That(address).Key("network_packet_tagging_status").HasValue("disabled"))},
			{Config: config(t, "enabled"), Check: check.That(address).Key("network_packet_tagging_status").HasValue("enabled")},
			{ResourceName: address, ImportState: true, ImportStateId: "crossTenantAccess", ImportStateVerify: true, ImportStateVerifyIgnore: []string{"timeouts"}},
			{Config: config(t, "enabled"), PlanOnly: true},
			{Config: config(t, "enabled"), Destroy: true},
		},
	})
	status, patches := m.Snapshot()
	require.Equal(t, "enabled", status, "destroy must not reset tenant settings")
	require.Equal(t, 2, patches, "only initial apply and status update should PATCH")
	calls := httpmock.GetCallCountInfo()
	require.Zero(t, calls["DELETE "+settingsmocks.URL])
	require.Zero(t, calls["POST "+settingsmocks.URL])
}

func TestUnitResourceNetworkCrossTenantAccessSettings_02_InvalidStatus(t *testing.T) {
	setupMockEnvironment(t)
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps:                    []resource.TestStep{{Config: config(t, "invalid"), ExpectError: regexp.MustCompile(`Attribute network_packet_tagging_status value must be one of`)}},
	})
}

func TestUnitResourceNetworkCrossTenantAccessSettings_03_Drift(t *testing.T) {
	m := setupMockEnvironment(t)
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: config(t, "disabled")},
			{PreConfig: func() { m.SetStatus("enabled") }, Config: config(t, "disabled"), Check: check.That(resourceType + ".test").Key("network_packet_tagging_status").HasValue("disabled")},
		},
	})
}

func TestUnitResourceNetworkCrossTenantAccessSettings_04_Forbidden(t *testing.T) {
	m := setupMockEnvironment(t)
	m.RegisterErrorMocks()
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps:                    []resource.TestStep{{Config: config(t, "disabled"), ExpectError: regexp.MustCompile(`(?i)(permission|privilege|forbidden|access)`)}},
	})
}
