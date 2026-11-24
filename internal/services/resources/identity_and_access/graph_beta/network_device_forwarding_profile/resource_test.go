package graphBetaNetworkDeviceForwardingProfile_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	target "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_beta/network_device_forwarding_profile"
	testmocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_beta/network_device_forwarding_profile/mocks"
)

const resourceType = target.ResourceName

func setup(t *testing.T) *testmocks.Mock {
	t.Helper()
	mocks.SetupUnitTestEnvironment(t)
	httpmock.Activate()
	mocks.NewMocks().AuthMocks.RegisterMocks()
	m := &testmocks.Mock{}
	m.RegisterMocks()
	crud.ResetClaims(target.ResourceName)
	t.Cleanup(func() { httpmock.DeactivateAndReset(); m.CleanupMockState() })
	return m
}

func hcl(t *testing.T, name string) string {
	t.Helper()
	config, err := helpers.ParseHCLFile("tests/terraform/unit/" + name)
	require.NoError(t, err)
	return config
}

func TestUnitResourceNetworkDeviceForwardingProfile_01_Lifecycle(t *testing.T) {
	m := setup(t)
	capture := mocks.CaptureGraphRequests(t)
	initial := hcl(t, "resource.tf")
	updated := hcl(t, "resource_updated.tf")
	addr := resourceType + ".test"
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
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
				),
			},
			{Config: updated, PlanOnly: true},
		},
	})
	capture.AssertMatchesFixture(t, "tests/requests/lifecycle.json")

	require.Len(t, m.Posts, 1)
	require.Equal(t, true, m.Posts[0]["isCustomProfile"])
	require.Equal(t, map[string]any{
		"assignments": map[string]any{
			"@odata.type": "#microsoft.graph.networkaccess.forwardingProfileAssignmentsEnumeratedMembership",
			"users":       []any{},
			"groups":      []any{map[string]any{"id": "aaaaaaaa-1111-4111-8111-111111111111"}},
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

func TestUnitResourceNetworkDeviceForwardingProfile_02_MinimalKeepsServerDefaults(t *testing.T) {
	m := setup(t)
	capture := mocks.CaptureGraphRequests(t)
	minimal := hcl(t, "resource_minimal.tf")
	addr := resourceType + ".test"
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
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
	capture.AssertMatchesFixture(t, "tests/requests/minimal.json")
	require.Len(t, m.Posts, 1)
	for _, key := range []string{"assignmentRules", "localNetworkSettings", "forwardingSourceMatching", "clientFallbackAction", "description"} {
		require.NotContains(t, m.Posts[0], key)
	}
	require.Empty(t, m.Patches)
}

func TestUnitResourceNetworkDeviceForwardingProfile_03_AdoptsProfileCreatedBeforeGatewayTimeout(
	t *testing.T,
) {
	m := setup(t)
	capture := mocks.CaptureGraphRequests(t)
	t.Cleanup(target.ShortenCreatedLookup())
	m.PostFailures = 1
	minimal := hcl(t, "resource_minimal.tf")
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: minimal,
				Check:  check.That(resourceType + ".test").Key("id").HasValue(testmocks.ProfileID),
			},
		},
	})
	capture.AssertMatchesFixture(t, "tests/requests/adopt_after_gateway_timeout.json")
	require.Len(t, m.Posts, 1, "the timed-out POST is not retried")
	require.Equal(t, 1, m.Deletes)
}

// A retried POST fails as a duplicate; the profile the first one created is adopted.
func TestUnitResourceNetworkDeviceForwardingProfile_03b_AdoptsProfileAfterRetriedPostFails(
	t *testing.T,
) {
	m := setup(t)
	capture := mocks.CaptureGraphRequests(t)
	t.Cleanup(target.ShortenCreatedLookup())
	m.PostDuplicates = 1
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: hcl(t, "resource_minimal.tf"),
				Check:  check.That(resourceType + ".test").Key("id").HasValue(testmocks.ProfileID),
			},
		},
	})
	capture.AssertMatchesFixture(t, "tests/requests/adopt_after_duplicate.json")
	require.Equal(t, 1, m.Deletes)
}

// A profile that existed before the POST is never adopted.
func TestUnitResourceNetworkDeviceForwardingProfile_03c_DoesNotAdoptExistingProfile(t *testing.T) {
	m := setup(t)
	capture := mocks.CaptureGraphRequests(t)
	t.Cleanup(target.ShortenCreatedLookup())
	m.Existing = map[string]any{
		"id":                    "99999999-8888-7777-6666-555555555555",
		"name":                  "unit-minimal",
		"trafficForwardingType": "private",
		"priority":              150,
		"isCustomProfile":       true,
		"state":                 "disabled",
	}
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      hcl(t, "resource_minimal.tf"),
				ExpectError: regexp.MustCompile(`Bad Request - 400`),
			},
		},
	})
	capture.AssertMatchesFixture(t, "tests/requests/no_adopt_existing.json")
	require.Zero(t, m.Deletes, "the pre-existing profile must not be adopted and destroyed")
}

// A profile another resource created in the same apply is not adopted.
func TestUnitResourceNetworkDeviceForwardingProfile_03f_DoesNotAdoptProfileClaimedInSameApply(t *testing.T) {
	m := setup(t)
	capture := mocks.CaptureGraphRequests(t)
	t.Cleanup(target.ShortenCreatedLookup())
	crud.RecordClaim(target.ResourceName, testmocks.ProfileID)
	m.PostDuplicates = 1
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      hcl(t, "resource_minimal.tf"),
				ExpectError: regexp.MustCompile(`Bad Request - 400`),
			},
		},
	})
	capture.AssertMatchesFixture(t, "tests/requests/no_adopt_claimed.json")
}

// Without the list before the POST, a failed create is reported, not adopted.
func TestUnitResourceNetworkDeviceForwardingProfile_03d_NoAdoptionWithoutPriorList(t *testing.T) {
	m := setup(t)
	capture := mocks.CaptureGraphRequests(t)
	t.Cleanup(target.ShortenCreatedLookup())
	m.ListFailures = 1
	m.PostFailures = 1
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      hcl(t, "resource_minimal.tf"),
				ExpectError: regexp.MustCompile(`504`),
			},
		},
	})
	capture.AssertMatchesFixture(t, "tests/requests/no_adopt_without_list.json")
	require.Len(t, m.Posts, 1)
}

// A transport error while listing does not stop the search.
func TestUnitResourceNetworkDeviceForwardingProfile_03e_KeepsPollingAfterListTransportError(
	t *testing.T,
) {
	m := setup(t)
	capture := mocks.CaptureGraphRequests(t)
	t.Cleanup(target.ShortenCreatedLookup())
	m.PostFailures = 1
	m.ListResetsAfterCreate = 1
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: hcl(t, "resource_minimal.tf"),
				Check:  check.That(resourceType + ".test").Key("id").HasValue(testmocks.ProfileID),
			},
		},
	})
	capture.AssertMatchesFixture(t, "tests/requests/adopt_after_list_transport_error.json")
	require.Equal(t, 1, m.Deletes)
}

func TestUnitResourceNetworkDeviceForwardingProfile_04_ValidationErrors(t *testing.T) {
	setup(t)
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: hcl(t, "resource_invalid.tf"),
				ExpectError: regexp.MustCompile(
					`(?s)(Invalid Attribute Value|User assignments are not supported).*(Invalid Attribute Value|User assignments are not supported)`,
				),
			},
		},
	})
}

func TestUnitResourceNetworkDeviceForwardingProfile_05_RejectsBuiltInProfileImport(t *testing.T) {
	setup(t)
	capture := mocks.CaptureGraphRequests(t)
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
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
	capture.AssertMatchesFixture(t, "tests/requests/reject_built_in_profile_import.json")
}
