package graphBetaIdentityAndAccessNetworkRemoteNetwork_test

import (
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	target "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_beta/network_remote_network"
	remoteNetworkMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_beta/network_remote_network/mocks"
)

const resourceAddress = target.ResourceName + ".test"

func setupMock(t *testing.T) *remoteNetworkMocks.RemoteNetworkMock {
	t.Helper()
	mocks.SetupUnitTestEnvironment(t)
	httpmock.Activate()
	mocks.NewMocks().AuthMocks.RegisterMocks()
	m := &remoteNetworkMocks.RemoteNetworkMock{}
	m.RegisterMocks()
	restore := target.SetDeletionPollInterval(10 * time.Millisecond)
	t.Cleanup(func() {
		restore()
		httpmock.DeactivateAndReset()
		m.CleanupMockState()
	})
	return m
}

func loadConfig(t *testing.T, filename string) string {
	t.Helper()
	config, err := helpers.ParseHCLFile("tests/terraform/unit/" + filename)
	require.NoError(t, err)
	return config
}

// Create, import, rename in place, then destroy.
func TestUnitResourceNetworkRemoteNetwork_01_CreateRenameImport(t *testing.T) {
	m := setupMock(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadConfig(t, "resource_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceAddress).Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-f-]{36}$`)),
					check.That(resourceAddress).Key("name").HasValue("unit-test-remote-network"),
					check.That(resourceAddress).Key("region").HasValue("eastUS"),
					check.That(resourceAddress).Key("version").HasValue("1.0.0"),
					check.That(resourceAddress).Key("forwarding_profile_ids.#").HasValue("1"),
					resource.TestCheckTypeSetElemAttr(resourceAddress, "forwarding_profile_ids.*", remoteNetworkMocks.M365ProfileID),
				),
			},
			{
				ResourceName:            resourceAddress,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
			{Config: loadConfig(t, "resource_minimal.tf"), PlanOnly: true},
			{
				Config: loadConfig(t, "resource_renamed.tf"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceAddress, plancheck.ResourceActionUpdate),
					},
				},
				Check: check.That(resourceAddress).Key("name").HasValue("unit-test-remote-network-renamed"),
			},
			{Config: loadConfig(t, "resource_renamed.tf"), PlanOnly: true},
		},
	})

	require.Len(t, m.WritesTo("POST", "/remoteNetworks"), 1)
	renames := m.WritesTo("PATCH", "")
	require.Len(t, renames, 1, "a rename must not touch forwarding profiles")
	require.Equal(t, map[string]any{"name": "unit-test-remote-network-renamed"}, renames[0].Body)
	require.Len(t, m.WritesTo("DELETE", ""), 1)
	require.Positive(t, m.StaleGets, "read-back should have observed replica lag")
}

// Forwarding profile changes use a $delta PATCH carrying the full desired set.
func TestUnitResourceNetworkRemoteNetwork_02_UpdateForwardingProfiles(t *testing.T) {
	m := setupMock(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: loadConfig(t, "resource_minimal.tf")},
			{
				Config: loadConfig(t, "resource_both_profiles.tf"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceAddress, plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceAddress).Key("forwarding_profile_ids.#").HasValue("2"),
					resource.TestCheckTypeSetElemAttr(resourceAddress, "forwarding_profile_ids.*", remoteNetworkMocks.M365ProfileID),
					resource.TestCheckTypeSetElemAttr(resourceAddress, "forwarding_profile_ids.*", remoteNetworkMocks.InternetProfileID),
				),
			},
			{
				Config: loadConfig(t, "resource_internet_only.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceAddress).Key("forwarding_profile_ids.#").HasValue("1"),
					resource.TestCheckTypeSetElemAttr(resourceAddress, "forwarding_profile_ids.*", remoteNetworkMocks.InternetProfileID),
				),
			},
			{Config: loadConfig(t, "resource_internet_only.tf"), PlanOnly: true},
		},
	})

	deltas := m.WritesTo("PATCH", "/forwardingProfiles")
	require.Len(t, deltas, 2)
	require.Equal(t, "#$delta", deltas[0].Body["@context"])
	require.Len(t, deltas[0].Body["value"], 2)
	require.Equal(t, []any{map[string]any{"id": remoteNetworkMocks.InternetProfileID}}, deltas[1].Body["value"])
	for _, patch := range m.WritesTo("PATCH", "") {
		require.NotContains(t, patch.Body, "forwardingProfiles")
	}
}

// Region is immutable, so changing it replaces the remote network.
func TestUnitResourceNetworkRemoteNetwork_03_RegionForcesReplace(t *testing.T) {
	m := setupMock(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: loadConfig(t, "resource_minimal.tf")},
			{
				Config: loadConfig(t, "resource_region_changed.tf"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceAddress, plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Check: check.That(resourceAddress).Key("region").HasValue("westEurope"),
			},
		},
	})

	// The replacement reuses the name, so Delete must wait for the old network.
	require.Len(t, m.WritesTo("POST", "/remoteNetworks"), 2)
}

// A remote network deleted outside Terraform is removed from state on refresh.
func TestUnitResourceNetworkRemoteNetwork_04_NotFoundRemovesFromState(t *testing.T) {
	m := setupMock(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: loadConfig(t, "resource_minimal.tf")},
			{
				PreConfig:          m.DeleteAll,
				Config:             loadConfig(t, "resource_minimal.tf"),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

// Graph rejects disabled and private forwarding profiles; the 400 message is surfaced.
func TestUnitResourceNetworkRemoteNetwork_05_GraphRejectsProfile(t *testing.T) {
	setupMock(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadConfig(t, "resource_disabled_profile.tf"),
				ExpectError: regexp.MustCompile(`as\s+it\s+is\s+disabled`),
			},
			{
				Config:      loadConfig(t, "resource_private_profile.tf"),
				ExpectError: regexp.MustCompile(`not allowed to be\s+assigned to a\s+remote network`),
			},
		},
	})
}

// A failed $delta PATCH keeps the prior associations in state, so the next plan retries.
func TestUnitResourceNetworkRemoteNetwork_06_FailedProfileUpdateKeepsState(t *testing.T) {
	setupMock(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: loadConfig(t, "resource_minimal.tf")},
			{
				Config:      loadConfig(t, "resource_disabled_profile.tf"),
				ExpectError: regexp.MustCompile(`as\s+it\s+is\s+disabled`),
			},
			{
				Config: loadConfig(t, "resource_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceAddress).Key("forwarding_profile_ids.#").HasValue("1"),
					resource.TestCheckTypeSetElemAttr(resourceAddress, "forwarding_profile_ids.*", remoteNetworkMocks.M365ProfileID),
				),
			},
		},
	})
}

// Plan-time validation rejects invalid configuration before any Graph call.
func TestUnitResourceNetworkRemoteNetwork_07_Validation(t *testing.T) {
	m := setupMock(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadConfig(t, "resource_invalid_empty_profiles.tf"),
				ExpectError: regexp.MustCompile(`set must contain at least 1 elements`),
			},
			{
				Config:      loadConfig(t, "resource_invalid_region.tf"),
				ExpectError: regexp.MustCompile(`value must be one of`),
			},
			{
				Config:      loadConfig(t, "resource_invalid_guid.tf"),
				ExpectError: regexp.MustCompile(`must\s+be\s+a\s+lowercase\s+GUID`),
			},
			{
				Config:      loadConfig(t, "resource_invalid_uppercase_guid.tf"),
				ExpectError: regexp.MustCompile(`must\s+be\s+a\s+lowercase\s+GUID`),
			},
		},
	})

	require.Empty(t, m.Requests)
}

// Import requires a GUID.
func TestUnitResourceNetworkRemoteNetwork_08_InvalidImportID(t *testing.T) {
	setupMock(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:        loadConfig(t, "resource_minimal.tf"),
				ResourceName:  resourceAddress,
				ImportState:   true,
				ImportStateId: "unit-test-remote-network",
				ExpectError:   regexp.MustCompile(`Invalid import ID`),
			},
		},
	})
}

// Permission errors on create are surfaced.
func TestUnitResourceNetworkRemoteNetwork_09_ErrorHandling(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	httpmock.Activate()
	mocks.NewMocks().AuthMocks.RegisterMocks()
	m := &remoteNetworkMocks.RemoteNetworkMock{}
	m.RegisterErrorMocks()
	t.Cleanup(func() { httpmock.DeactivateAndReset() })

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadConfig(t, "resource_minimal.tf"),
				ExpectError: regexp.MustCompile(`(?i)forbidden|insufficient privileges|403`),
			},
		},
	})
}

// Changing only timeouts updates state without calling Graph.
func TestUnitResourceNetworkRemoteNetwork_10_TimeoutsOnlyUpdate(t *testing.T) {
	m := setupMock(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: loadConfig(t, "resource_minimal.tf")},
			{
				Config: loadConfig(t, "resource_timeouts_changed.tf"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceAddress, plancheck.ResourceActionUpdate),
					},
				},
				Check: check.That(resourceAddress).Key("timeouts.delete").HasValue("240s"),
			},
		},
	})

	require.Empty(t, m.WritesTo("PATCH", ""))
}

// The deletion wait survives a transient error and a 404 from DELETE.
func TestUnitResourceNetworkRemoteNetwork_11_DeletionWaitEdgeCases(t *testing.T) {
	for name, configure := range map[string]func(*remoteNetworkMocks.RemoteNetworkMock){
		"transient error":            func(m *remoteNetworkMocks.RemoteNetworkMock) { m.DeletionErrors = 1 },
		"delete reports in progress": func(m *remoteNetworkMocks.RemoteNetworkMock) { m.DeleteReportsDeleting = true },
	} {
		t.Run(name, func(t *testing.T) {
			m := setupMock(t)
			configure(m)
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{Config: loadConfig(t, "resource_minimal.tf")},
					{Config: loadConfig(t, "resource_region_changed.tf")},
				},
			})
			require.Len(t, m.WritesTo("POST", "/remoteNetworks"), 2)
		})
	}
}
