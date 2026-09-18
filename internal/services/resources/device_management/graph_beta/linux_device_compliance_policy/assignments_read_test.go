package graphBetaLinuxDeviceCompliancePolicy

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/microsoft/kiota-abstractions-go/authentication"
	msgraph "github.com/microsoftgraph/msgraph-beta-sdk-go"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
	"github.com/stretchr/testify/require"

	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/shared_models/graph_beta/device_management"
)

func TestUnitAssignmentReadUsesDedicatedEndpointAndFollowsPages(t *testing.T) {
	for _, secondPageFails := range []bool{false, true} {
		t.Run(fmt.Sprintf("second_page_fails_%t", secondPageFails), func(t *testing.T) {
			calls := 0
			var server *httptest.Server
			server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				calls++
				w.Header().Set("Content-Type", "application/json")
				if req.Method != http.MethodGet || req.URL.Path != "/beta/deviceManagement/compliancePolicies/policy-id/assignments" {
					t.Errorf("unexpected request: %s %s", req.Method, req.URL)
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				if req.URL.Query().Get("page") == "2" {
					if secondPageFails {
						w.WriteHeader(http.StatusNotFound)
						fmt.Fprint(w, `{"error":{"code":"NotFound","message":"second page unavailable"}}`)
						return
					}
					fmt.Fprint(w, `{"value":[{"id":"exclude","target":{"@odata.type":"#microsoft.graph.exclusionGroupAssignmentTarget","groupId":"excluded-group","deviceAndAppManagementAssignmentFilterType":"none"}}]}`)
					return
				}
				fmt.Fprintf(w, `{"@odata.nextLink":%q,"value":[{"id":"include","target":{"@odata.type":"#microsoft.graph.groupAssignmentTarget","groupId":"included-group","deviceAndAppManagementAssignmentFilterType":"none"}}]}`, server.URL+"/beta/deviceManagement/compliancePolicies/policy-id/assignments?page=2")
			}))
			t.Cleanup(server.Close)

			adapter, err := msgraph.NewGraphRequestAdapter(&authentication.AnonymousAuthenticationProvider{})
			if err != nil {
				t.Fatal(err)
			}
			adapter.SetBaseUrl(server.URL + "/beta")
			r := &LinuxDeviceCompliancePolicyResource{client: msgraph.NewGraphServiceClient(adapter)}
			assignments, err := r.getAllPolicyAssignmentsWithPageIterator(context.Background(), "policy-id")
			if secondPageFails {
				if err == nil || assignments != nil {
					t.Fatalf("failed pagination must not return partial assignments: %v, %v", assignments, err)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if calls != 2 || len(assignments) != 2 {
				t.Fatalf("expected both pages, got %d calls and %d assignments", calls, len(assignments))
			}
			policy := models.NewDeviceManagementCompliancePolicy()
			policy.SetAssignments(assignments)
			var state LinuxDeviceCompliancePolicyResourceModel
			MapRemoteStateToTerraform(context.Background(), &state, policy)
			if state.Assignments.IsNull() || len(state.Assignments.Elements()) != 2 {
				t.Fatalf("assignment scope was lost: %v", state.Assignments)
			}
		})
	}
}

func TestUnitLinuxComplianceReadAssignments(t *testing.T) {
	for _, tc := range []struct {
		name        string
		policyCode  int
		firstCode   int
		secondCode  int
		unassigned  bool
		expectError bool
	}{
		{name: "assigned", policyCode: 200, firstCode: 200, secondCode: 200},
		{name: "unassigned", policyCode: 200, firstCode: 200, unassigned: true},
		{name: "first_page_bad_request", policyCode: 200, firstCode: 400, expectError: true},
		{name: "first_page_forbidden", policyCode: 200, firstCode: 403, expectError: true},
		{name: "first_page_not_found", policyCode: 200, firstCode: 404, expectError: true},
		{name: "second_page_not_found", policyCode: 200, firstCode: 200, secondCode: 404, expectError: true},
		{name: "empty_assignment_response", policyCode: 200, firstCode: 204, expectError: true},
		{name: "empty_policy_response", policyCode: 204, expectError: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			assignmentCalls := 0
			var server *httptest.Server
			server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				if req.Method != http.MethodGet {
					t.Errorf("read issued an unexpected method: %s", req.Method)
					w.WriteHeader(http.StatusMethodNotAllowed)
					return
				}
				switch req.URL.Path {
				case "/beta/deviceManagement/compliancePolicies/policy-id":
					w.WriteHeader(tc.policyCode)
					if tc.policyCode == http.StatusOK {
						fmt.Fprint(w, `{"id":"policy-id","name":"Example Linux policy","assignments":[],"roleScopeTagIds":["0"]}`)
					}
				case "/beta/deviceManagement/compliancePolicies/policy-id/settings":
					fmt.Fprint(w, `{"value":[]}`)
				case "/beta/deviceManagement/compliancePolicies/policy-id/assignments":
					assignmentCalls++
					secondPage := req.URL.Query().Get("page") == "2"
					code := tc.firstCode
					if secondPage {
						code = tc.secondCode
					}
					w.WriteHeader(code)
					switch {
					case code == http.StatusNoContent:
						return
					case code != http.StatusOK:
						fmt.Fprint(w, `{"error":{"code":"AssignmentReadFailure","message":"Cannot read assignments"}}`)
					case tc.unassigned:
						fmt.Fprint(w, `{"value":[]}`)
					case secondPage:
						fmt.Fprint(w, `{"value":[{"id":"exclude","target":{"@odata.type":"#microsoft.graph.exclusionGroupAssignmentTarget","groupId":"excluded-group","deviceAndAppManagementAssignmentFilterType":"none"}}]}`)
					default:
						fmt.Fprintf(w, `{"@odata.nextLink":%q,"value":[{"id":"include","target":{"@odata.type":"#microsoft.graph.groupAssignmentTarget","groupId":"included-group","deviceAndAppManagementAssignmentFilterType":"none"}}]}`, server.URL+"/beta/deviceManagement/compliancePolicies/policy-id/assignments?page=2")
					}
				default:
					t.Errorf("unexpected path: %s", req.URL.Path)
					w.WriteHeader(http.StatusNotFound)
				}
			}))
			t.Cleanup(server.Close)

			adapter, err := msgraph.NewGraphRequestAdapter(&authentication.AnonymousAuthenticationProvider{})
			require.NoError(t, err)
			adapter.SetBaseUrl(server.URL + "/beta")
			r := &LinuxDeviceCompliancePolicyResource{client: msgraph.NewGraphServiceClient(adapter)}

			var schema resource.SchemaResponse
			r.Schema(ctx, resource.SchemaRequest{}, &schema)
			stateType := schema.Schema.Type().TerraformType(ctx).(tftypes.Object)
			attributes := make(map[string]tftypes.Value, len(stateType.AttributeTypes))
			for name, attributeType := range stateType.AttributeTypes {
				attributes[name] = tftypes.NewValue(attributeType, nil)
			}
			attributes["id"] = tftypes.NewValue(tftypes.String, "policy-id")
			attributes["name"] = tftypes.NewValue(tftypes.String, "Previous policy name")
			state := tfsdk.State{Schema: schema.Schema, Raw: tftypes.NewValue(stateType, attributes)}
			response := resource.ReadResponse{State: state}
			r.Read(ctx, resource.ReadRequest{State: state}, &response)

			require.Equal(t, tc.expectError, response.Diagnostics.HasError(), "%v", response.Diagnostics)
			if tc.expectError {
				require.True(t, state.Raw.Equal(response.State.Raw), "an incomplete read must preserve prior state")
				return
			}
			var actual LinuxDeviceCompliancePolicyResourceModel
			require.False(t, response.State.Get(ctx, &actual).HasError())
			require.Equal(t, "policy-id", actual.ID.ValueString())
			if tc.unassigned {
				require.Equal(t, 1, assignmentCalls)
				require.True(t, actual.Assignments.IsNull())
				return
			}
			require.Equal(t, 2, assignmentCalls)
			var assignments []sharedmodels.DeviceCompliancePolicyAssignmentResourceModel
			require.False(t, actual.Assignments.ElementsAs(ctx, &assignments, false).HasError())
			require.Len(t, assignments, 2)
			actualTargets := make(map[string]string, len(assignments))
			for _, assignment := range assignments {
				actualTargets[assignment.GroupId.ValueString()] = assignment.Type.ValueString()
			}
			require.Equal(t, map[string]string{
				"included-group": "groupAssignmentTarget",
				"excluded-group": "exclusionGroupAssignmentTarget",
			}, actualTargets)
		})
	}
}
