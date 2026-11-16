package graphBetaNetworkForwardingProfilePrivateAccessPolicyLink_test

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	target "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_beta/network_forwarding_profile_private_access_policy_link"
	testmocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_beta/network_forwarding_profile_private_access_policy_link/mocks"
)

const resourceAddress = target.ResourceName + ".test"

func setup(t *testing.T) *testmocks.Mock {
	t.Helper()
	mocks.SetupUnitTestEnvironment(t)
	httpmock.Activate()
	mocks.NewMocks().AuthMocks.RegisterMocks()
	m := &testmocks.Mock{}
	m.RegisterMocks()
	t.Cleanup(func() {
		httpmock.DeactivateAndReset()
		m.CleanupMockState()
	})
	return m
}

func config(t *testing.T, name string) string {
	t.Helper()
	cfg, err := helpers.ParseHCLFile("tests/terraform/unit/" + name + ".tf")
	require.NoError(t, err)
	return cfg
}

func importID(prefix string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceAddress]
		if !ok {
			return "", fmt.Errorf("resource %s not found", resourceAddress)
		}
		return prefix + rs.Primary.ID, nil
	}
}

func removals(m *testmocks.Mock) int {
	count := 0
	for _, p := range m.Patches {
		for _, raw := range p["policies@delta"].([]any) {
			if _, ok := raw.(map[string]any)["@removed"]; ok {
				count++
			}
		}
	}
	return count
}

func TestUnitResourceNetworkForwardingProfilePrivateAccessPolicyLink_01_UserLifecycle(t *testing.T) {
	m := setup(t)
	m.StaleReads = 2
	var firstLinkID string
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config(t, "user_app_policy"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceAddress).Key("profile_type").HasValue("user"),
					check.That(resourceAddress).Key("state").HasValue("enabled"),
					check.That(resourceAddress).Key("priority").HasValue("120"),
					check.That(resourceAddress).Key("config_type").HasValue("appPolicy"),
					check.That(resourceAddress).Key("policy_link_id").MatchesRegex(regexp.MustCompile(`^bbbbbbbb-`)),
					func(s *terraform.State) error {
						rs := s.RootModule().Resources[resourceAddress]
						firstLinkID = rs.Primary.Attributes["policy_link_id"]
						if rs.Primary.ID != testmocks.UserProfileID+"/"+firstLinkID {
							return fmt.Errorf("unexpected id %q", rs.Primary.ID)
						}
						if len(m.Patches) != 1 {
							return fmt.Errorf("expected one add PATCH, got %d", len(m.Patches))
						}
						return nil
					},
				),
			},
			{
				ResourceName:      resourceAddress,
				ImportState:       true,
				ImportStateIdFunc: importID(""),
				ImportStateVerify: true,
			},
			{
				ResourceName:      resourceAddress,
				ImportState:       true,
				ImportStateIdFunc: importID("user/"),
				ImportStateVerify: true,
			},
			{Config: config(t, "user_app_policy"), PlanOnly: true},
			{
				// Priority and state changes are a remove then a separate add.
				Config: config(t, "user_app_policy_updated"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceAddress).Key("priority").HasValue("99"),
					check.That(resourceAddress).Key("state").HasValue("disabled"),
					func(s *terraform.State) error {
						linkID := s.RootModule().Resources[resourceAddress].Primary.Attributes["policy_link_id"]
						if linkID == firstLinkID {
							return fmt.Errorf("expected a new link ID after replacement")
						}
						if len(m.Patches) != 3 || removals(m) != 1 {
							return fmt.Errorf("expected add, remove, add; got %v", m.Patches)
						}
						return nil
					},
				),
			},
		},
	})
	links := m.Links(testmocks.UserCollection, testmocks.UserProfileID)
	require.Len(t, links, 1)
	require.Equal(t, testmocks.AutoLinkID, links[0].ID)
}

func TestUnitResourceNetworkForwardingProfilePrivateAccessPolicyLink_02_AdoptAutoLink(t *testing.T) {
	m := setup(t)
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config(t, "adopt_auto_link"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceAddress).Key("policy_link_id").HasValue(testmocks.AutoLinkID),
					check.That(resourceAddress).Key("priority").HasValue("65000"),
					func(*terraform.State) error {
						if len(m.Patches) != 0 {
							return fmt.Errorf("adoption must not write, got %v", m.Patches)
						}
						return nil
					},
				),
			},
			{Config: config(t, "adopt_auto_link"), PlanOnly: true},
		},
	})
	require.Empty(t, m.Links(testmocks.UserCollection, testmocks.UserProfileID))
	require.Equal(t, 1, removals(m))
}

func TestUnitResourceNetworkForwardingProfilePrivateAccessPolicyLink_03_ReplaceMismatchedAutoLink(t *testing.T) {
	m := setup(t)
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config(t, "replace_auto_link"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceAddress).Key("priority").HasValue("140"),
					func(s *terraform.State) error {
						linkID := s.RootModule().Resources[resourceAddress].Primary.Attributes["policy_link_id"]
						if linkID == testmocks.AutoLinkID {
							return fmt.Errorf("expected the auto link to be replaced")
						}
						if len(m.Patches) != 2 || removals(m) != 1 {
							return fmt.Errorf("expected separate remove and add PATCHes, got %v", m.Patches)
						}
						return nil
					},
				),
			},
		},
	})
}

func TestUnitResourceNetworkForwardingProfilePrivateAccessPolicyLink_04_LinkIDChangesOutOfBand(t *testing.T) {
	m := setup(t)
	const reAddedID = "cccccccc-0000-0000-0000-000000000001"
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: config(t, "user_app_policy")},
			{
				PreConfig: func() {
					m.SetLinks(testmocks.UserCollection, testmocks.UserProfileID, []testmocks.Link{{
						ID: reAddedID, PolicyID: testmocks.AppPolicyID, State: "enabled", ConfigType: "appPolicy", Priority: 120,
					}})
				},
				Config: config(t, "user_app_policy"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceAddress).Key("policy_link_id").HasValue(reAddedID),
					check.That(resourceAddress).Key("id").HasValue(testmocks.UserProfileID+"/"+reAddedID),
				),
			},
		},
	})
}

func TestUnitResourceNetworkForwardingProfilePrivateAccessPolicyLink_05_DeviceProfile(t *testing.T) {
	m := setup(t)
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config(t, "device_dns_policy"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceAddress).Key("profile_type").HasValue("device"),
					check.That(resourceAddress).Key("config_type").HasValue("dnsPolicy"),
					check.That(resourceAddress).Key("state").HasValue("disabled"),
					check.That(resourceAddress).Key("priority").HasValue("65000"),
				),
			},
			{
				ResourceName:      resourceAddress,
				ImportState:       true,
				ImportStateIdFunc: importID("device/"),
				ImportStateVerify: true,
			},
		},
	})
	require.Empty(t, m.Links(testmocks.DeviceCollection, testmocks.DeviceProfileID))
	require.Len(t, m.Links(testmocks.UserCollection, testmocks.UserProfileID), 1)
}

func TestUnitResourceNetworkForwardingProfilePrivateAccessPolicyLink_06_RemovedOutOfBand(t *testing.T) {
	m := setup(t)
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: config(t, "user_app_policy")},
			{
				PreConfig: func() {
					m.SetLinks(testmocks.UserCollection, testmocks.UserProfileID, []testmocks.Link{})
				},
				Config:             config(t, "user_app_policy"),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
			{
				PreConfig: func() {
					m.DeleteProfile(testmocks.UserCollection, testmocks.UserProfileID)
				},
				Config:             config(t, "user_app_policy"),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestUnitResourceNetworkForwardingProfilePrivateAccessPolicyLink_07_Validation(t *testing.T) {
	setup(t)
	cases := map[string]string{
		"invalid_priority":     `must be between 1 and 65000`,
		"invalid_config_type":  `value must be one of`,
		"invalid_profile_type": `value must be one of`,
		"invalid_ids":          `must be a valid UUID`,
	}
	steps := make([]resource.TestStep, 0, len(cases))
	for name, pattern := range cases {
		steps = append(steps, resource.TestStep{
			Config:      config(t, name),
			PlanOnly:    true,
			ExpectError: regexp.MustCompile(pattern),
		})
	}
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps:                    steps,
	})
}

func TestUnitResourceNetworkForwardingProfilePrivateAccessPolicyLink_08_GraphErrors(t *testing.T) {
	setup(t)
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      config(t, "unknown_policy"),
				ExpectError: regexp.MustCompile(`of traffic type Private was not found`),
			},
		},
	})
}

func TestUnitResourceNetworkForwardingProfilePrivateAccessPolicyLink_09_Forbidden(t *testing.T) {
	m := setup(t)
	m.RegisterErrorMocks()
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      config(t, "user_app_policy"),
				ExpectError: regexp.MustCompile(`(?i)forbidden|permission`),
			},
		},
	})
}

func TestUnitResourceNetworkForwardingProfilePrivateAccessPolicyLink_10_InvalidImportID(t *testing.T) {
	setup(t)
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:        config(t, "user_app_policy"),
				ResourceName:  resourceAddress,
				ImportState:   true,
				ImportStateId: "tenant/" + testmocks.UserProfileID + "/bbbbbbbb-0000-0000-0000-000000000001",
				ExpectError:   regexp.MustCompile(`Invalid import ID`),
			},
		},
	})
}

// A failed add after removing a mismatched link restores the old link.
func TestUnitResourceNetworkForwardingProfilePrivateAccessPolicyLink_11_RestoresLinkWhenAddFails(t *testing.T) {
	m := setup(t)
	m.FailAdds = 1
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      config(t, "replace_auto_link"),
				ExpectError: regexp.MustCompile(`Simulated add failure`),
			},
		},
	})
	links := m.Profiles[testmocks.Key(testmocks.UserCollection, testmocks.UserProfileID)]
	require.Len(t, links, 1)
	require.Equal(t, testmocks.AutoLinkedPolicyID, links[0].PolicyID)
	require.Equal(t, int64(65000), links[0].Priority)
	require.Equal(t, "enabled", links[0].State)
}

// A link with a different config type is reported, not removed.
func TestUnitResourceNetworkForwardingProfilePrivateAccessPolicyLink_12_ConfigTypeMismatch(t *testing.T) {
	m := setup(t)
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      config(t, "auto_link_config_type_mismatch"),
				ExpectError: regexp.MustCompile(`already linked with a different config type`),
			},
		},
	})
	require.Empty(t, m.Patches)
}

// Import IDs are lowercased to match Graph.
func TestUnitResourceNetworkForwardingProfilePrivateAccessPolicyLink_13_ImportLowercasesIDs(t *testing.T) {
	setup(t)
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: config(t, "adopt_auto_link")},
			{
				Config:                  config(t, "adopt_auto_link"),
				ResourceName:            resourceAddress,
				ImportState:             true,
				ImportStateId:           strings.ToUpper(testmocks.UserProfileID + "/" + testmocks.AutoLinkID),
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}
