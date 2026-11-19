package graphBetaServicePrincipalSynchronizationJob

import (
	"context"
	"net/http"
	"testing"
	"testing/synctest"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"

	providerclient "github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
)

const (
	testParent  = "11111111-1111-4111-8111-111111111111"
	testJob     = "scim.opaque-job"
	testJobsURL = "https://graph.microsoft.com/beta/servicePrincipals/" + testParent + "/synchronization/jobs"
)

func testModel() ServicePrincipalSynchronizationJobResourceModel {
	return ServicePrincipalSynchronizationJobResourceModel{
		ID: types.StringValue(
			testParent + "/" + testJob,
		),
		ServicePrincipalID: types.StringValue(testParent),
		JobID:              types.StringValue(testJob),
		TemplateID:         types.StringValue("scim"),
		Enabled:            types.BoolValue(false),
		Schedule:           types.ObjectNull(scheduleAttributeTypes),
		StatusCode:         types.StringNull(),
		Timeouts: timeouts.Value{
			Object: types.ObjectNull(map[string]attr.Type{"create": types.StringType, "read": types.StringType, "update": types.StringType, "delete": types.StringType}),
		},
	}
}

func testState(t *testing.T, r *ServicePrincipalSynchronizationJobResource) tfsdk.State {
	t.Helper()
	var schema resource.SchemaResponse
	r.Schema(context.Background(), resource.SchemaRequest{}, &schema)
	state := tfsdk.State{Schema: schema.Schema}
	data := testModel()
	require.False(t, state.Set(context.Background(), &data).HasError())
	return state
}

func errorReply(status int) (*http.Response, error) {
	return httpmock.NewJsonResponse(status, map[string]any{"error": map[string]any{"code": "Request_BadRequest", "message": "probe error"}})
}

func TestUnitSynchronizationJobImportIDs(t *testing.T) {
	parent, job, err := parseResourceID(testParent + "/" + testJob)
	require.NoError(t, err)
	require.Equal(t, testParent, parent)
	require.Equal(t, testJob, job)
	for _, id := range []string{"", testJob, "bad/" + testJob, testParent + "/", testParent + "/a/b"} {
		_, _, err := parseResourceID(id)
		require.ErrorIs(t, err, errImportID)
	}
}

func TestUnitSynchronizationJobReadStateSafety(t *testing.T) {
	for _, tc := range []struct {
		name                     string
		status                   int
		postWrite, missing, fail bool
	}{
		{name: "bad request", status: 400, fail: true}, {name: "forbidden", status: 403, fail: true}, {name: "missing", status: 404, missing: true}, {name: "post-write missing", status: 404, postWrite: true, fail: true}, {name: "empty response", status: 204, fail: true}, {name: "incomplete schedule", status: 200, fail: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			transport := httpmock.NewMockTransport()
			transport.RegisterResponder("GET", testJobsURL+"/"+testJob, func(_ *http.Request) (*http.Response, error) {
				if tc.status == 204 {
					return httpmock.NewStringResponse(204, ""), nil
				}
				if tc.status == 200 {
					return httpmock.NewJsonResponse(200, map[string]any{"id": testJob, "templateId": "scim"})
				}
				return errorReply(tc.status)
			})
			r := NewServicePrincipalSynchronizationJobResource().(*ServicePrincipalSynchronizationJobResource)
			r.client = providerclient.NewMockGraphClients(&http.Client{Transport: transport}).MockBetaClient
			state := testState(t, r)
			resp := resource.ReadResponse{State: state}
			ctx := context.Background()
			if tc.postWrite {
				ctx = context.WithValue(ctx, "retry_operation", "Create")
			}
			r.Read(ctx, resource.ReadRequest{State: state}, &resp)
			require.Equal(t, tc.fail, resp.Diagnostics.HasError(), resp.Diagnostics)
			require.Equal(t, tc.missing, resp.State.Raw.IsNull())
		})
	}
}

func TestUnitSynchronizationJobCreateRecordsOwnershipBeforeStart(t *testing.T) {
	transport := httpmock.NewMockTransport()
	posts := 0
	transport.RegisterResponder("POST", testJobsURL, func(_ *http.Request) (*http.Response, error) {
		posts++
		return httpmock.NewJsonResponse(201, map[string]any{"id": testJob})
	})
	transport.RegisterResponder("POST", testJobsURL+"/"+testJob+"/start", func(_ *http.Request) (*http.Response, error) { return errorReply(403) })
	r := NewServicePrincipalSynchronizationJobResource().(*ServicePrincipalSynchronizationJobResource)
	r.client = providerclient.NewMockGraphClients(&http.Client{Transport: transport}).MockBetaClient
	state := testState(t, r)
	data := testModel()
	data.ID = types.StringNull()
	data.JobID = types.StringNull()
	data.Enabled = types.BoolValue(true)
	require.False(t, state.Set(context.Background(), &data).HasError())
	resp := resource.CreateResponse{State: tfsdk.State{Schema: state.Schema}}
	r.Create(context.Background(), resource.CreateRequest{Plan: tfsdk.Plan{Schema: state.Schema, Raw: state.Raw}}, &resp)
	require.True(t, resp.Diagnostics.HasError())
	require.Equal(t, 1, posts)
	var id types.String
	require.False(t, resp.State.GetAttribute(context.Background(), path.Root("id"), &id).HasError())
	require.Equal(t, testParent+"/"+testJob, id.ValueString())
}

func TestUnitSynchronizationJobDeletionSafety(t *testing.T) {
	for _, status := range []int{204, 404, 403} {
		t.Run(http.StatusText(status), func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				transport := httpmock.NewMockTransport()
				gets := 0
				transport.RegisterResponder("DELETE", testJobsURL+"/"+testJob, func(_ *http.Request) (*http.Response, error) {
					if status == 204 {
						return httpmock.NewStringResponse(204, ""), nil
					}
					return errorReply(status)
				})
				transport.RegisterResponder("GET", testJobsURL+"/"+testJob, func(_ *http.Request) (*http.Response, error) {
					gets++
					if gets == 1 {
						return httpmock.NewJsonResponse(200, map[string]any{"id": testJob})
					}
					return errorReply(404)
				})
				r := NewServicePrincipalSynchronizationJobResource().(*ServicePrincipalSynchronizationJobResource)
				r.client = providerclient.NewMockGraphClients(&http.Client{Transport: transport}).MockBetaClient
				state := testState(t, r)
				resp := resource.DeleteResponse{State: state}
				r.Delete(context.Background(), resource.DeleteRequest{State: state}, &resp)
				require.Equal(t, status == 403, resp.Diagnostics.HasError(), resp.Diagnostics)
				require.Equal(t, status != 403, resp.State.Raw.IsNull())
				if status == 403 {
					require.Zero(t, gets)
				} else {
					require.Equal(t, 4, gets)
				}
			})
		})
	}
}
