package graphBetaNetworkDefaultForwardingProfile_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	target "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_beta/network_default_forwarding_profile"
	testmocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_beta/network_default_forwarding_profile/mocks"
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

func TestUnitResourceNetworkDefaultForwardingProfile_01_Lifecycle(t *testing.T) {
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
					check.That(addr).Key("id").HasValue(testmocks.InternetID),
					check.That(addr).Key("name").HasValue("Internet traffic forwarding profile"),
					check.That(addr).Key("priority").HasValue("300"),
					check.That(addr).Key("assignment_rules.membership").HasValue("enumerated"),
					check.That(addr).Key("assignment_rules.platforms.#").HasValue("4"),
				),
			},
			{
				ResourceName:            addr,
				ImportState:             true,
				ImportStateId:           "internet",
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
			{
				ResourceName:            addr,
				ImportState:             true,
				ImportStateId:           testmocks.InternetID,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
			{Config: initial, PlanOnly: true},
			{
				Config: updated,
				Check: resource.ComposeTestCheckFunc(
					check.That(addr).Key("state").HasValue("enabled"),
					check.That(addr).Key("assignment_rules.membership").HasValue("none"),
					check.That(addr).Key("assignment_rules.group_ids.#").DoesNotExist(),
				),
			},
			{Config: updated, PlanOnly: true},
		},
	})

	require.Len(t, m.Patches, 1, "adopting with matching values sends no request")
	require.Equal(t, map[string]any{
		"state":                    "enabled",
		"clientFallbackAction":     "block",
		"forwardingSourceMatching": map[string]any{"sessionType": "agent"},
		"assignmentRules": map[string]any{
			"assignments": map[string]any{
				"@odata.type": "#microsoft.graph.networkaccess.forwardingProfileAssignmentsNoMembership",
			},
			"devicePropertyFilters": map[string]any{"platforms": "windows"},
		},
	}, m.Patches[0])
	require.Zero(t, m.Writes, "destroy must not delete the default profile")
	require.Equal(t, "enabled", m.Internet["state"], "destroy leaves the remote profile unchanged")
}

func TestUnitResourceNetworkDefaultForwardingProfile_02_RejectsLocalNetworkOnInternet(
	t *testing.T,
) {
	setup(t)
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      hcl(t, "resource_invalid.tf"),
				ExpectError: regexp.MustCompile(`Unsupported local network settings`),
			},
		},
	})
}

func TestUnitResourceNetworkDefaultForwardingProfile_03_RejectsCustomProfileImport(t *testing.T) {
	setup(t)
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:        hcl(t, "resource.tf"),
				ResourceName:  resourceType + ".test",
				ImportState:   true,
				ImportStateId: testmocks.CustomID,
				ExpectError:   regexp.MustCompile(`Not a default forwarding profile`),
			},
		},
	})
}
