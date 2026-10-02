package graphBetaAndroidDeviceOwnerCompliancePolicy

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	providerschema "github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	frameworkresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/microsoft/kiota-abstractions-go/authentication"
	msgraph "github.com/microsoftgraph/msgraph-beta-sdk-go"
	"github.com/stretchr/testify/require"
)

// The isolated provider exercises the real resource through Terraform without
// credentials or registering the provider's unrelated resources.
type scheduleTestProvider struct {
	graph *msgraph.GraphServiceClient
}

func (p *scheduleTestProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "microsoft365"
}

func (p *scheduleTestProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = providerschema.Schema{}
}

func (p *scheduleTestProvider) Configure(_ context.Context, _ provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	resp.ResourceData = &client.GraphClients{KiotaGraphBetaClient: p.graph}
}

func (p *scheduleTestProvider) Resources(context.Context) []func() frameworkresource.Resource {
	return []func() frameworkresource.Resource{NewAndroidDeviceOwnerCompliancePolicyResource}
}

func (p *scheduleTestProvider) DataSources(context.Context) []func() datasource.DataSource {
	return nil
}

// Graph persists the actions but drops ruleName and supplies action defaults.
// Keep the mock stateful so a changed configuration has to reach the write API.
type scheduleTestGraph struct {
	mu     sync.Mutex
	policy map[string]any
	writes int
}

func (g *scheduleTestGraph) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	g.mu.Lock()
	defer g.mu.Unlock()
	w.Header().Set("Content-Type", "application/json")
	const collection = "/beta/deviceManagement/deviceCompliancePolicies"
	const item = collection + "/policy-id"
	var body map[string]any
	if req.Method == http.MethodPost || req.Method == http.MethodPatch {
		var reader io.Reader = req.Body
		if req.Header.Get("Content-Encoding") == "gzip" {
			decoded, err := gzip.NewReader(req.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			defer decoded.Close()
			reader = decoded
		}
		if err := json.NewDecoder(reader).Decode(&body); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		g.writes++
	}
	switch {
	case req.Method == http.MethodPost && req.URL.Path == collection:
		g.policy = body
		g.policy["id"] = "policy-id"
		g.policy["roleScopeTagIds"] = []string{"0"}
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(g.policy)
	case req.Method == http.MethodPatch && req.URL.Path == item:
		for key, value := range body {
			g.policy[key] = value
		}
		w.WriteHeader(http.StatusNoContent)
	case req.Method == http.MethodPost && req.URL.Path == item+"/scheduleActionsForRules":
		g.policy["scheduledActionsForRule"] = body["deviceComplianceScheduledActionForRules"]
		w.WriteHeader(http.StatusNoContent)
	case req.Method == http.MethodPost && req.URL.Path == item+"/assign":
		g.policy["assignments"] = body["assignments"]
		_ = json.NewEncoder(w).Encode(map[string]any{"value": body["assignments"]})
	case req.Method == http.MethodGet && req.URL.Path == item && g.policy != nil:
		expand := req.URL.Query().Get("$expand")
		if !strings.Contains(expand, "assignments") || !strings.Contains(expand, "scheduledActionsForRule($expand=scheduledActionConfigurations)") {
			http.Error(w, "missing navigation expansion", http.StatusBadRequest)
			return
		}
		for _, ruleValue := range g.policy["scheduledActionsForRule"].([]any) {
			rule := ruleValue.(map[string]any)
			rule["ruleName"] = nil
			for _, actionValue := range rule["scheduledActionConfigurations"].([]any) {
				action := actionValue.(map[string]any)
				if action["gracePeriodHours"] == nil {
					action["gracePeriodHours"] = 0
				}
				if action["notificationTemplateId"] == nil || action["notificationTemplateId"] == "" {
					action["notificationTemplateId"] = "00000000-0000-0000-0000-000000000000"
				}
				if action["notificationMessageCCList"] == nil {
					action["notificationMessageCCList"] = []any{}
				}
			}
		}
		_ = json.NewEncoder(w).Encode(g.policy)
	case req.Method == http.MethodDelete && req.URL.Path == item:
		g.policy = nil
		w.WriteHeader(http.StatusNoContent)
	default:
		w.WriteHeader(http.StatusNotFound)
		_, _ = fmt.Fprint(w, `{"error":{"code":"ResourceNotFound","message":"Policy not found"}}`)
	}
}

func TestUnitAndroidComplianceScheduledActionLifecycle(t *testing.T) {
	for _, tc := range []struct {
		name         string
		rule         string
		template     string
		initialGrace string
		wantRule     string
	}{
		{name: "explicit_rule", rule: `rule_name = "PasswordRequired"`, initialGrace: "grace_period_hours = 0", wantRule: "PasswordRequired"},
		{name: "omitted_defaults", wantRule: "PasswordRequired"},
		{name: "legacy_rule_migration", rule: `rule_name = "unavailable"`, wantRule: "unavailable"},
		{name: "explicit_empty_template", rule: `rule_name = "PasswordRequired"`, template: `notification_template_id = ""`, wantRule: "PasswordRequired"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			graph := &scheduleTestGraph{}
			server := httptest.NewServer(graph)
			t.Cleanup(server.Close)
			adapter, err := msgraph.NewGraphRequestAdapter(&authentication.AnonymousAuthenticationProvider{})
			require.NoError(t, err)
			adapter.SetBaseUrl(server.URL + "/beta")
			p := &scheduleTestProvider{graph: msgraph.NewGraphServiceClient(adapter)}
			config := func(rule, grace string) string {
				return fmt.Sprintf(`resource "%s" "test" {
  display_name = "Example Android policy"
  scheduled_actions_for_rule = [{
    %s
    scheduled_action_configurations = [{
      action_type = "block"
      %s
      %s
    }]
  }]
}`, ResourceName, rule, grace, tc.template)
			}
			initial := config(tc.rule, tc.initialGrace)
			updated := config(`rule_name = "PasswordRequired"`, "grace_period_hours = 24")
			address := ResourceName + ".test"
			noChanges := resource.ConfigPlanChecks{PostApplyPostRefresh: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()}}
			steps := []resource.TestStep{
				{
					Config: initial,
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(address, "scheduled_actions_for_rule.0.rule_name", tc.wantRule),
						resource.TestCheckTypeSetElemNestedAttrs(address, "scheduled_actions_for_rule.0.scheduled_action_configurations.*", map[string]string{"action_type": "block", "grace_period_hours": "0"}),
					),
				},
				{Config: initial, PlanOnly: true, ConfigPlanChecks: noChanges},
				{
					Config:           updated,
					ConfigPlanChecks: resource.ConfigPlanChecks{PreApply: []plancheck.PlanCheck{plancheck.ExpectResourceAction(address, plancheck.ResourceActionUpdate)}},
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(address, "scheduled_actions_for_rule.0.rule_name", "PasswordRequired"),
						resource.TestCheckTypeSetElemNestedAttrs(address, "scheduled_actions_for_rule.0.scheduled_action_configurations.*", map[string]string{"action_type": "block", "grace_period_hours": "24"}),
					),
				},
				{Config: updated, PlanOnly: true, ConfigPlanChecks: noChanges},
			}
			if tc.template == "" {
				steps = append(steps, resource.TestStep{
					ResourceName:            address,
					ImportState:             true,
					ImportStateVerify:       true,
					ImportStateVerifyIgnore: []string{"timeouts"},
				})
			}
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
					"microsoft365": providerserver.NewProtocol6WithError(p),
				},
				Steps: steps,
			})
			graph.mu.Lock()
			defer graph.mu.Unlock()
			require.Nil(t, graph.policy, "Terraform must delete the test policy")
			require.Equal(t, 5, graph.writes, "only create/assign and update/schedule/assign should write")
		})
	}
}
