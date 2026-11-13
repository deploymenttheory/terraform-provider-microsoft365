package graphBetaNetworkForwardingProfile_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	target "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_beta/network_forwarding_profile"
	testmocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_beta/network_forwarding_profile/mocks"
)

const resourceType = target.ResourceName

func setup(t *testing.T) *testmocks.Mock {
	t.Helper()
	mocks.SetupUnitTestEnvironment(t)
	httpmock.Activate()
	mocks.NewMocks().AuthMocks.RegisterMocks()
	m := &testmocks.Mock{}
	m.RegisterMocks()
	t.Cleanup(func() { httpmock.DeactivateAndReset(); m.CleanupMockState() })
	return m
}

func hcl(t *testing.T, name string) string {
	t.Helper()
	config, err := helpers.ParseHCLFile("tests/terraform/unit/" + name)
	require.NoError(t, err)
	return config
}

func TestUnitResourceNetworkForwardingProfile_01_Lifecycle(t *testing.T) {
	m := setup(t)
	initial := hcl(t, "resource.tf")
	updated := hcl(t, "resource_updated.tf")
	addr := resourceType + ".test"
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: initial,
				Check: resource.ComposeTestCheckFunc(
					check.That(addr).Key("id").HasValue(testmocks.ProfileID),
					check.That(addr).Key("client_fallback_action").HasValue("bypass"),
					check.That(addr).
						Key("local_network_settings.allow_local_subnet_bypass").
						HasValue("true"),
					check.That(addr).Key("assignment_rules.membership").HasValue("enumerated"),
					check.That(addr).Key("assignment_rules.group_ids.#").HasValue("1"),
					check.That(addr).Key("assignment_rules.platforms.#").HasValue("2"),
				),
			},
			{
				ResourceName:            addr,
				ImportState:             true,
				ImportStateId:           testmocks.ProfileID,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
			{Config: initial, PlanOnly: true},
			{
				Config: updated,
				Check: resource.ComposeTestCheckFunc(
					check.That(addr).Key("name").HasValue("unit-home-network-updated"),
					check.That(addr).Key("description").DoesNotExist(),
					check.That(addr).Key("assignment_rules.membership").HasValue("all"),
					check.That(addr).Key("assignment_rules.group_ids.#").DoesNotExist(),
					check.That(addr).
						Key("forwarding_source_matching.session_type").
						HasValue("agent"),
				),
			},
			{Config: updated, PlanOnly: true},
		},
	})

	require.Len(t, m.Posts, 1)
	require.Equal(t, true, m.Posts[0]["isCustomProfile"])
	require.Equal(t, map[string]any{
		"assignments": map[string]any{
			"@odata.type": "#microsoft.graph.networkaccess.forwardingProfileAssignmentsEnumeratedMembership",
			"users":       []any{},
			"groups":      []any{map[string]any{"id": "a6ed41d3-45ac-4e12-8601-bdb56216e84e"}},
			"devices":     []any{},
		},
		"devicePropertyFilters": map[string]any{"platforms": "windows,macOS"},
	}, m.Posts[0]["assignmentRules"])

	require.Len(t, m.Patches, 1)
	patch := m.Patches[0]
	require.NotContains(t, patch, "trafficForwardingType")
	require.Contains(t, patch, "description")
	require.Nil(t, patch["description"], "removing description sends an explicit null")
	require.Equal(t, map[string]any{"platforms": "windows,android,iOS"},
		patch["assignmentRules"].(map[string]any)["devicePropertyFilters"])
	require.Equal(t, 1, m.Deletes)
	require.Nil(t, m.State)
}

func TestUnitResourceNetworkForwardingProfile_02_MinimalKeepsServerDefaults(t *testing.T) {
	m := setup(t)
	minimal := hcl(t, "resource_minimal.tf")
	addr := resourceType + ".test"
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: minimal,
				Check: resource.ComposeTestCheckFunc(
					check.That(addr).Key("client_fallback_action").HasValue("bypass"),
					check.That(addr).Key("assignment_rules.membership").HasValue("none"),
					check.That(addr).Key("assignment_rules.platforms.#").HasValue("0"),
					check.That(addr).Key("local_network_settings").DoesNotExist(),
				),
			},
			{Config: minimal, PlanOnly: true},
		},
	})
	require.Len(t, m.Posts, 1)
	for _, key := range []string{"assignmentRules", "localNetworkSettings", "forwardingSourceMatching", "clientFallbackAction", "description"} {
		require.NotContains(t, m.Posts[0], key)
	}
	require.Empty(t, m.Patches)
}

func TestUnitResourceNetworkForwardingProfile_03_AdoptsProfileCreatedBeforeGatewayTimeout(
	t *testing.T,
) {
	m := setup(t)
	t.Cleanup(target.ShortenCreatedLookup())
	m.PostFailures = 1
	minimal := hcl(t, "resource_minimal.tf")
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: minimal,
				Check:  check.That(resourceType + ".test").Key("id").HasValue(testmocks.ProfileID),
			},
		},
	})
	require.Len(t, m.Posts, 1, "the timed-out POST is not retried")
	require.Equal(t, 1, m.Deletes)
}

func TestUnitResourceNetworkForwardingProfile_04_ValidationErrors(t *testing.T) {
	setup(t)
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: hcl(t, "resource_invalid.tf"),
				ExpectError: regexp.MustCompile(
					`(?s)Unsupported local network settings.*Missing assignment IDs|(?s)Missing assignment IDs.*Unsupported local network settings`,
				),
			},
		},
	})
}

func TestUnitResourceNetworkForwardingProfile_05_RejectsDefaultProfileImport(t *testing.T) {
	setup(t)
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:        hcl(t, "resource_minimal.tf"),
				ResourceName:  resourceType + ".test",
				ImportState:   true,
				ImportStateId: testmocks.DefaultID,
				ExpectError:   regexp.MustCompile(`Not a custom forwarding profile`),
			},
		},
	})
}
