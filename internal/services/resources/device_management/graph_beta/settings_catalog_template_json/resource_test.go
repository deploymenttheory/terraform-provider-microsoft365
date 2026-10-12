package graphBetaDeviceManagementTemplateJson_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	testresource "github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/shared_models/graph_beta/device_management"
	templatejson "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/settings_catalog_template_json"
)

const (
	testPolicyID     = "11111111-1111-1111-1111-111111111111"
	testTemplateType = "windows_anti_virus_defender_update_controls"
	testTemplateID   = "e3f74c5a-a6de-411d-aef6-eb15628f3a0a_1"
	testPolicyURL    = "https://graph.microsoft.com/beta/deviceManagement/configurationPolicies/" + testPolicyID
	testPutURL       = "https://graph.microsoft.com/beta/deviceManagement/configurationPolicies('" + testPolicyID + "')"
	testSettings     = `{"settings":[{"id":"0","settingInstance":{"@odata.type":"#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance","choiceSettingValue":{"children":[],"settingValueTemplateReference":{"settingValueTemplateId":"afc8df70-7b19-4335-b200-bf4b7e098f67","useTemplateDefault":false},"value":"device_vendor_msft_defender_configuration_engineupdateschannel_6"},"settingDefinitionId":"device_vendor_msft_defender_configuration_engineupdateschannel","settingInstanceTemplateReference":{"settingInstanceTemplateId":"f7e1409d-9c85-4a3f-85a6-ad05cc8ccf13"}}}]}`
)

func setupTemplateMocks(t *testing.T, includeTemplateMetadata bool) *mocks.Mocks {
	t.Helper()
	httpmock.Activate()
	t.Cleanup(httpmock.DeactivateAndReset)
	clients := mocks.NewMocks()
	clients.AuthMocks.RegisterMocks()

	policy := map[string]any{
		"id":                   testPolicyID,
		"name":                 "Template policy",
		"description":          "Updated description",
		"platforms":            "windows10",
		"technologies":         "mdm,microsoftSense",
		"roleScopeTagIds":      []string{"0"},
		"settingCount":         1,
		"isAssigned":           false,
		"createdDateTime":      "2024-01-01T00:00:00Z",
		"lastModifiedDateTime": "2024-01-02T00:00:00Z",
	}
	if includeTemplateMetadata {
		policy["templateReference"] = map[string]any{
			"templateId":             testTemplateID,
			"templateFamily":         "endpointSecurityAntivirus",
			"templateDisplayName":    "Defender Update controls",
			"templateDisplayVersion": "Version 1",
		}
	}
	responder, err := httpmock.NewJsonResponder(http.StatusOK, policy)
	require.NoError(t, err)
	httpmock.RegisterResponder(http.MethodGet, testPolicyURL, responder)
	createResponder, err := httpmock.NewJsonResponder(http.StatusCreated, policy)
	require.NoError(t, err)
	httpmock.RegisterResponder(
		http.MethodPost,
		"https://graph.microsoft.com/beta/deviceManagement/configurationPolicies",
		createResponder,
	)

	var settings map[string]any
	require.NoError(t, json.Unmarshal([]byte(testSettings), &settings))
	settingsResponder, err := httpmock.NewJsonResponder(
		http.StatusOK,
		map[string]any{"value": settings["settings"]},
	)
	require.NoError(t, err)
	httpmock.RegisterResponder(
		http.MethodGet,
		testPutURL+"/settings",
		settingsResponder,
	)
	emptyAssignments, err := httpmock.NewJsonResponder(
		http.StatusOK,
		map[string]any{"value": []any{}},
	)
	require.NoError(t, err)
	httpmock.RegisterResponder(
		http.MethodGet,
		testPolicyURL+"/assignments",
		emptyAssignments,
	)
	httpmock.RegisterResponder(
		http.MethodPost,
		testPolicyURL+"/assign",
		emptyAssignments,
	)
	httpmock.RegisterResponder(
		http.MethodPut,
		testPutURL,
		httpmock.NewStringResponder(http.StatusNoContent, ""),
	)
	httpmock.RegisterResponder(
		http.MethodDelete,
		testPolicyURL,
		httpmock.NewStringResponder(http.StatusNoContent, ""),
	)
	return clients
}

func configuredTemplateResource(
	t *testing.T,
	clients *mocks.Mocks,
) *templatejson.DeviceManagementTemplateJsonResource {
	t.Helper()
	r := templatejson.NewDeviceManagementTemplateJsonResource().(*templatejson.DeviceManagementTemplateJsonResource)
	var response resource.ConfigureResponse
	r.Configure(
		context.Background(),
		resource.ConfigureRequest{ProviderData: clients.Clients},
		&response,
	)
	require.False(t, response.Diagnostics.HasError(), "%v", response.Diagnostics)
	return r
}

func templateTestModel() sharedmodels.SettingsCatalogJsonResourceModel {
	return sharedmodels.SettingsCatalogJsonResourceModel{
		ID:                          types.StringValue(testPolicyID),
		Name:                        types.StringValue("Template policy"),
		Description:                 types.StringValue("Updated description"),
		SettingsCatalogTemplateType: types.StringValue(testTemplateType),
		Platforms:                   types.StringValue("windows10"),
		Technologies: types.ListValueMust(types.StringType, []attr.Value{
			types.StringValue("mdm"), types.StringValue("microsoftSense"),
		}),
		RoleScopeTagIds: types.SetValueMust(
			types.StringType,
			[]attr.Value{types.StringValue("0")},
		),
		Settings:             types.StringValue(testSettings),
		SettingsCount:        types.Int32Value(1),
		IsAssigned:           types.BoolValue(false),
		CreatedDateTime:      types.StringValue("2024-01-01T00:00:00Z"),
		LastModifiedDateTime: types.StringValue("2024-01-02T00:00:00Z"),
		Assignments: types.SetNull(
			templatejson.SettingsCatalogConfigurationPolicyAssignmentType(),
		),
		Timeouts: timeouts.Value{Object: types.ObjectNull(map[string]attr.Type{
			"create": types.StringType, "read": types.StringType,
			"update": types.StringType, "delete": types.StringType,
		})},
	}
}

func templateTestState(
	t *testing.T,
	r *templatejson.DeviceManagementTemplateJsonResource,
	data sharedmodels.SettingsCatalogJsonResourceModel,
) tfsdk.State {
	t.Helper()
	var response resource.SchemaResponse
	r.Schema(context.Background(), resource.SchemaRequest{}, &response)
	state := tfsdk.State{Schema: response.Schema}
	require.False(t, state.Set(context.Background(), &data).HasError())
	return state
}

func TestUnitSettingsCatalogTemplateJsonUpdateRecoversNullTemplateType(t *testing.T) {
	clients := setupTemplateMocks(t, false)
	r := configuredTemplateResource(t, clients)
	prior := templateTestModel()
	prior.SettingsCatalogTemplateType = types.StringNull()
	prior.Description = types.StringValue("Old description")
	state := templateTestState(t, r, prior)
	planned := templateTestModel()
	planState := templateTestState(t, r, planned)

	var putBody map[string]any
	httpmock.RegisterResponder(
		http.MethodPut,
		testPutURL,
		func(request *http.Request) (*http.Response, error) {
			if err := json.NewDecoder(request.Body).Decode(&putBody); err != nil {
				return nil, err
			}
			return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
		},
	)

	response := resource.UpdateResponse{State: state}
	r.Update(context.Background(), resource.UpdateRequest{
		State: state,
		Plan:  tfsdk.Plan{Schema: planState.Schema, Raw: planState.Raw},
	}, &response)
	require.False(t, response.Diagnostics.HasError(), "%v", response.Diagnostics)
	var result sharedmodels.SettingsCatalogJsonResourceModel
	require.False(t, response.State.Get(context.Background(), &result).HasError())
	require.Equal(t, planned.SettingsCatalogTemplateType, result.SettingsCatalogTemplateType)
	require.Equal(t, planned.ID, result.ID)
	require.Equal(t, planned.Description, result.Description)
	require.JSONEq(t, planned.Settings.ValueString(), result.Settings.ValueString())
	require.Equal(t, planned.Assignments, result.Assignments)
	require.Equal(t, testTemplateID, putBody["templateReference"].(map[string]any)["templateId"])
	require.Equal(t, 1, httpmock.GetCallCountInfo()[http.MethodPost+" "+testPolicyURL+"/assign"])
}

func TestUnitSettingsCatalogTemplateJsonUpdateWriteFailureRetainsState(t *testing.T) {
	for _, method := range []string{http.MethodPut, http.MethodPost} {
		t.Run(method, func(t *testing.T) {
			clients := setupTemplateMocks(t, true)
			r := configuredTemplateResource(t, clients)
			prior := templateTestModel()
			prior.SettingsCatalogTemplateType = types.StringNull()
			prior.Description = types.StringValue("Old description")
			state := templateTestState(t, r, prior)
			planState := templateTestState(t, r, templateTestModel())
			url := testPutURL
			if method == http.MethodPost {
				url = testPolicyURL + "/assign"
			}
			httpmock.RegisterResponder(
				method,
				url,
				httpmock.NewStringResponder(http.StatusBadRequest,
					`{"error":{"code":"BadRequest","message":"Synthetic write failure"}}`),
			)

			response := resource.UpdateResponse{State: state}
			r.Update(context.Background(), resource.UpdateRequest{
				State: state,
				Plan:  tfsdk.Plan{Schema: planState.Schema, Raw: planState.Raw},
			}, &response)
			require.True(t, response.Diagnostics.HasError())
			require.True(t, state.Raw.Equal(response.State.Raw))
			require.Zero(t, httpmock.GetCallCountInfo()[http.MethodGet+" "+testPolicyURL])
		})
	}
}

func TestUnitSettingsCatalogTemplateJsonImportConverges(t *testing.T) {
	t.Setenv("TF_ACC", "0")
	t.Setenv("MS365_TEST_MODE", "true")
	setupTemplateMocks(t, true)
	config := fmt.Sprintf(`resource "%s" "test" {
  name                           = "Template policy"
  description                    = "Updated description"
  settings_catalog_template_type = %q
  settings                       = %q
}`, templatejson.ResourceName, testTemplateType, testSettings)

	t.Run("create_and_verify_import", func(t *testing.T) {
		testresource.UnitTest(t, testresource.TestCase{
			ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
			Steps: []testresource.TestStep{
				{
					Config: config,
					Check: testresource.ComposeTestCheckFunc(
						testresource.TestCheckResourceAttr(
							templatejson.ResourceName+".test",
							"id",
							testPolicyID,
						),
						testresource.TestCheckResourceAttr(
							templatejson.ResourceName+".test",
							"settings_catalog_template_type",
							testTemplateType,
						),
					),
				},
				{
					ResourceName:      templatejson.ResourceName + ".test",
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	})

	// Persist imports into a fresh workspace so the next plan uses the imported state.
	t.Run("import_existing_policy_and_plan", func(t *testing.T) {
		testresource.UnitTest(t, testresource.TestCase{
			ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
			Steps: []testresource.TestStep{
				{
					Config:             config,
					ResourceName:       templatejson.ResourceName + ".test",
					ImportStateId:      testPolicyID,
					ImportState:        true,
					ImportStatePersist: true,
					ImportStateCheck: func(states []*terraform.InstanceState) error {
						require.Len(t, states, 1)
						require.Equal(t, testPolicyID, states[0].ID)
						require.Equal(
							t,
							testTemplateType,
							states[0].Attributes["settings_catalog_template_type"],
						)
						return nil
					},
				},
				{
					Config:   config,
					PlanOnly: true,
					ConfigPlanChecks: testresource.ConfigPlanChecks{
						PostApplyPreRefresh: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
					},
				},
			},
		})
	})
	require.Zero(t, httpmock.GetCallCountInfo()[http.MethodPut+" "+testPutURL])
}
