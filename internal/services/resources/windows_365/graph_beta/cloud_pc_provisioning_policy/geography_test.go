package graphBetaCloudPcProvisioningPolicy_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
)

const geographyResourceName = "microsoft365_graph_beta_windows_365_cloud_pc_provisioning_policy.geography"

// geographyGraph is a request-capturing Graph fake for the provisioning policy endpoints. It stores the last POST and
// PATCH bodies and answers GET with the stored policy, after readBack has rewritten its domain join configurations the
// way Graph returns them.
type geographyGraph struct {
	sync.Mutex
	policy   map[string]any
	posts    []map[string]any
	patches  []map[string]any
	readBack func(configs []any) []any
}

func (g *geographyGraph) register(t *testing.T) {
	t.Helper()
	httpmock.Activate()
	t.Cleanup(httpmock.DeactivateAndReset)
	base := "https://graph.microsoft.com/beta/deviceManagement/virtualEndpoint/provisioningPolicies"
	httpmock.RegisterResponder(http.MethodGet, base, func(*http.Request) (*http.Response, error) {
		return httpmock.NewJsonResponse(http.StatusOK, map[string]any{"value": []any{}})
	})
	httpmock.RegisterResponder(http.MethodPost, base, func(req *http.Request) (*http.Response, error) {
		var body map[string]any
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			return httpmock.NewStringResponse(http.StatusBadRequest, `{"error":{"code":"BadRequest","message":"invalid body"}}`), nil
		}
		g.Lock()
		defer g.Unlock()
		g.posts = append(g.posts, body)
		g.policy = map[string]any{"id": "7f8b4e1c-1f2a-4c3d-9e8f-0a1b2c3d4e5f", "gracePeriodInHours": 4, "microsoftManagedDesktop": map[string]any{"managedType": "notManaged", "profile": ""}}
		for key, value := range body {
			g.policy[key] = value
		}
		return httpmock.NewJsonResponse(http.StatusCreated, g.policy)
	})
	item := regexp.MustCompile(`^` + regexp.QuoteMeta(base) + `/[^/]+$`)
	httpmock.RegisterRegexpResponder(http.MethodGet, item, func(*http.Request) (*http.Response, error) {
		g.Lock()
		defer g.Unlock()
		if g.policy == nil {
			return httpmock.NewStringResponse(http.StatusNotFound, `{"error":{"code":"ResourceNotFound","message":"not found"}}`), nil
		}
		response := map[string]any{"assignments": []any{}}
		for key, value := range g.policy {
			response[key] = value
		}
		if configs, ok := g.policy["domainJoinConfigurations"].([]any); ok && g.readBack != nil {
			response["domainJoinConfigurations"] = g.readBack(configs)
		}
		return httpmock.NewJsonResponse(http.StatusOK, response)
	})
	httpmock.RegisterRegexpResponder(http.MethodPatch, item, func(req *http.Request) (*http.Response, error) {
		var body map[string]any
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			return httpmock.NewStringResponse(http.StatusBadRequest, `{"error":{"code":"BadRequest","message":"invalid body"}}`), nil
		}
		g.Lock()
		defer g.Unlock()
		g.patches = append(g.patches, body)
		for key, value := range body {
			g.policy[key] = value
		}
		return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
	})
	httpmock.RegisterRegexpResponder(http.MethodDelete, item, func(*http.Request) (*http.Response, error) {
		g.Lock()
		defer g.Unlock()
		g.policy = nil
		return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
	})
	httpmock.RegisterRegexpResponder(http.MethodPost, regexp.MustCompile(`^`+regexp.QuoteMeta(base)+`/[^/]+/assign$`), func(*http.Request) (*http.Response, error) {
		return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
	})
}

// lastDomainJoinConfigurations returns the domainJoinConfigurations of the last captured POST or PATCH.
func (g *geographyGraph) lastDomainJoinConfigurations(t *testing.T, patch bool) []any {
	t.Helper()
	g.Lock()
	defer g.Unlock()
	requests := g.posts
	if patch {
		requests = g.patches
	}
	require.NotEmpty(t, requests)
	configs, ok := requests[len(requests)-1]["domainJoinConfigurations"].([]any)
	require.True(t, ok, "request has no domainJoinConfigurations")
	return configs
}

func geographyConfig(displayName, domainJoin string) string {
	return fmt.Sprintf(`resource "microsoft365_graph_beta_windows_365_cloud_pc_provisioning_policy" "geography" {
  display_name          = %q
  image_id              = "microsoftwindowsdesktop_windows-ent-cpc_win11-26h2-ent-cpc-m365"
  image_type            = "gallery"
  enable_single_sign_on = true

  windows_setting = {
    locale = "en-US"
  }

  microsoft_managed_desktop = {
    managed_type = "notManaged"
  }

  domain_join_configurations = [
    {
%s
    }
  ]

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}`, displayName, domainJoin)
}

// TestUnitResourceCloudPcProvisioningPolicy_GeographyRequests asserts the exact domainJoinConfigurations sent on
// create and update for the Intune portal selections captured in #4904, and that each selection reads back, imports
// and plans with no diff. The portal also sends the deprecated "type" property; the provider sends its replacement
// domainJoinType with the same value, so the geography fields are compared field-for-field and domainJoinType is
// the only other key.
func TestUnitResourceCloudPcProvisioningPolicy_GeographyRequests(t *testing.T) {
	for _, tc := range []struct {
		name, domainJoin string
		portal           map[string]any
	}{
		{
			name:       "middle_east_whole_geography",
			domainJoin: "domain_join_type = \"azureADJoin\"\ngeographic_location_type = \"middleEast\"\nregion_group = \"automatic\"\nregion_name = \"automatic\"",
			portal:     map[string]any{"geographicLocationType": "middleEast", "regionName": "automatic", "regionGroup": "automatic"},
		},
		{
			name:       "middle_east_uae",
			domainJoin: "domain_join_type = \"azureADJoin\"\ngeographic_location_type = \"middleEast\"\nregion_group = \"uae\"\nregion_name = \"automatic\"",
			portal:     map[string]any{"geographicLocationType": "middleEast", "regionName": "automatic", "regionGroup": "uae"},
		},
		{
			name:       "africa_whole_geography",
			domainJoin: "domain_join_type = \"azureADJoin\"\ngeographic_location_type = \"africa\"\nregion_group = \"automatic\"\nregion_name = \"automatic\"",
			portal:     map[string]any{"geographicLocationType": "africa", "regionName": "automatic", "regionGroup": "automatic"},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			expected := map[string]any{"domainJoinType": "azureADJoin"}
			for key, value := range tc.portal {
				expected[key] = value
			}
			graph := &geographyGraph{readBack: func(configs []any) []any {
				// Graph returns the deprecated type alongside domainJoinType.
				out := make([]any, 0, len(configs))
				for _, config := range configs {
					returned := map[string]any{"type": "azureADJoin", "onPremisesConnectionId": nil}
					for key, value := range config.(map[string]any) {
						returned[key] = value
					}
					out = append(out, returned)
				}
				return out
			}}
			graph.register(t)
			created := geographyConfig("Geography "+tc.name, tc.domainJoin)
			updated := geographyConfig("Geography "+tc.name+" updated", tc.domainJoin)
			checks := resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr(geographyResourceName, "domain_join_configurations.0.geographic_location_type", tc.portal["geographicLocationType"].(string)),
				resource.TestCheckResourceAttr(geographyResourceName, "domain_join_configurations.0.region_group", tc.portal["regionGroup"].(string)),
				resource.TestCheckResourceAttr(geographyResourceName, "domain_join_configurations.0.region_name", "automatic"),
			)
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{Config: created, Check: resource.ComposeTestCheckFunc(checks, func(*terraform.State) error {
						require.Equal(t, []any{expected}, graph.lastDomainJoinConfigurations(t, false))
						return nil
					})},
					{Config: created, PlanOnly: true},
					{ResourceName: geographyResourceName, ImportState: true, ImportStateVerify: true},
					{Config: updated, Check: resource.ComposeTestCheckFunc(checks, func(*terraform.State) error {
						require.Equal(t, []any{expected}, graph.lastDomainJoinConfigurations(t, true))
						return nil
					})},
					{Config: updated, PlanOnly: true},
				},
			})
		})
	}
}

// TestUnitResourceCloudPcProvisioningPolicy_GeographyComputed checks that a policy configured without
// geographic_location_type keeps the geography Graph reports (middleEast for the israel region group) with no diff,
// including after import.
func TestUnitResourceCloudPcProvisioningPolicy_GeographyComputed(t *testing.T) {
	graph := &geographyGraph{readBack: func(configs []any) []any {
		out := make([]any, 0, len(configs))
		for _, config := range configs {
			returned := map[string]any{"geographicLocationType": "middleEast"}
			for key, value := range config.(map[string]any) {
				returned[key] = value
			}
			out = append(out, returned)
		}
		return out
	}}
	graph.register(t)
	config := geographyConfig("Geography israel", "domain_join_type = \"azureADJoin\"\nregion_group = \"israel\"\nregion_name = \"automatic\"")
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: config, Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr(geographyResourceName, "domain_join_configurations.0.geographic_location_type", "middleEast"),
				resource.TestCheckResourceAttr(geographyResourceName, "domain_join_configurations.0.region_group", "israel"),
				func(*terraform.State) error {
					require.Equal(t, []any{map[string]any{"domainJoinType": "azureADJoin", "regionGroup": "israel", "regionName": "automatic"}}, graph.lastDomainJoinConfigurations(t, false))
					return nil
				},
			)},
			{Config: config, PlanOnly: true},
			{ResourceName: geographyResourceName, ImportState: true, ImportStateVerify: true},
		},
	})
}

// TestUnitResourceCloudPcProvisioningPolicy_LegacyMiddleEastRegionGroup checks that the legacy middleEast region
// group, which Windows 365 reads back as israel, fails at plan time with guidance instead of an inconsistent result.
func TestUnitResourceCloudPcProvisioningPolicy_LegacyMiddleEastRegionGroup(t *testing.T) {
	config := geographyConfig("Geography legacy", "domain_join_type = \"azureADJoin\"\nregion_group = \"middleEast\"\nregion_name = \"automatic\"")
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps:                    []resource.TestStep{{Config: config, PlanOnly: true, ExpectError: regexp.MustCompile(`Unsupported region group`)}},
	})
}
