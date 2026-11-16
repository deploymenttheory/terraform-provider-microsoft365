package graphBetaIdentityAndAccessLifecycleWorkflowCustomTaskExtension

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/microsoft/kiota-abstractions-go/serialization"
	jsonserialization "github.com/microsoft/kiota-serialization-json-go"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

func serializeBody(t *testing.T, body serialization.Parsable) map[string]any {
	t.Helper()
	writer := jsonserialization.NewJsonSerializationWriter()
	if err := writer.WriteObjectValue("", body); err != nil {
		t.Fatalf("serialize body: %v", err)
	}
	content, err := writer.GetSerializedContent()
	if err != nil {
		t.Fatalf("read serialized body: %v", err)
	}
	var result map[string]any
	if err := json.Unmarshal(content, &result); err != nil {
		t.Fatalf("decode serialized body %s: %v", string(content), err)
	}
	return result
}

func minimalModel() *LifecycleWorkflowCustomTaskExtensionResourceModel {
	return &LifecycleWorkflowCustomTaskExtensionResourceModel{
		DisplayName: types.StringValue("extension"),
		Description: types.StringValue(""),
		EndpointConfiguration: &EndpointConfiguration{
			SubscriptionId:       types.StringValue("00000000-0000-0000-0000-000000000000"),
			ResourceGroupName:    types.StringValue("rg"),
			LogicAppWorkflowName: types.StringValue("la"),
			Url:                  types.StringValue("https://example.logic.azure.com/invoke"),
		},
		AuthenticationConfiguration: &AuthenticationConfiguration{Type: types.StringValue(authTypeProofOfPossession), ResourceId: types.StringNull()},
	}
}

func TestConstructResource_CreateOmitsCallbackAndDefaultsClientConfiguration(t *testing.T) {
	body, err := constructResource(context.Background(), minimalModel(), false)
	if err != nil {
		t.Fatalf("construct: %v", err)
	}
	result := serializeBody(t, body)

	if _, ok := result["callbackConfiguration"]; ok {
		t.Fatalf("callbackConfiguration must be omitted on create, got %v", result["callbackConfiguration"])
	}

	endpoint, _ := result["endpointConfiguration"].(map[string]any)
	if endpoint["@odata.type"] != "#microsoft.graph.logicAppTriggerEndpointConfiguration" || endpoint["url"] != "https://example.logic.azure.com/invoke" {
		t.Fatalf("unexpected endpointConfiguration %v", endpoint)
	}

	authentication, _ := result["authenticationConfiguration"].(map[string]any)
	if authentication["@odata.type"] != "#microsoft.graph.azureAdPopTokenAuthentication" {
		t.Fatalf("unexpected authenticationConfiguration %v", authentication)
	}

	client, _ := result["clientConfiguration"].(map[string]any)
	if client["maximumRetries"] != float64(defaultMaximumRetries) || client["timeoutInMilliseconds"] != float64(defaultTimeoutInMilliseconds) {
		t.Fatalf("unexpected clientConfiguration %v", client)
	}
}

func TestConstructResource_UpdateSendsNullCallback(t *testing.T) {
	body, err := constructResource(context.Background(), minimalModel(), true)
	if err != nil {
		t.Fatalf("construct: %v", err)
	}
	result := serializeBody(t, body)

	value, ok := result["callbackConfiguration"]
	if !ok || value != nil {
		t.Fatalf("callbackConfiguration must be null on update, got present=%v value=%v", ok, value)
	}
}

func TestConstructResource_TokenAuthenticationAndCallback(t *testing.T) {
	model := minimalModel()
	model.AuthenticationConfiguration = &AuthenticationConfiguration{Type: types.StringValue(authTypeAzureAd), ResourceId: types.StringValue("00000003-0000-0000-c000-000000000000")}
	model.ClientConfiguration = &ClientConfiguration{MaximumRetries: types.Int64Value(0), TimeoutInMilliseconds: types.Int64Value(2000)}
	model.CallbackConfiguration = &CallbackConfiguration{
		TimeoutDuration:  types.StringValue("PT1H"),
		AuthorizedAppIds: types.SetValueMust(types.StringType, []attr.Value{types.StringValue("00000003-0000-0000-c000-000000000000")}),
	}

	body, err := constructResource(context.Background(), model, true)
	if err != nil {
		t.Fatalf("construct: %v", err)
	}
	result := serializeBody(t, body)

	authentication, _ := result["authenticationConfiguration"].(map[string]any)
	if authentication["@odata.type"] != "#microsoft.graph.azureAdTokenAuthentication" || authentication["resourceId"] != "00000003-0000-0000-c000-000000000000" {
		t.Fatalf("unexpected authenticationConfiguration %v", authentication)
	}

	client, _ := result["clientConfiguration"].(map[string]any)
	if client["maximumRetries"] != float64(0) || client["timeoutInMilliseconds"] != float64(2000) {
		t.Fatalf("unexpected clientConfiguration %v", client)
	}

	callback, _ := result["callbackConfiguration"].(map[string]any)
	if callback["@odata.type"] != "#microsoft.graph.identityGovernance.customTaskExtensionCallbackConfiguration" || callback["timeoutDuration"] != "PT1H" {
		t.Fatalf("unexpected callbackConfiguration %v", callback)
	}
	apps, _ := callback["authorizedApps"].([]any)
	if len(apps) != 1 {
		t.Fatalf("expected one authorized app, got %v", callback["authorizedApps"])
	}
	if app, _ := apps[0].(map[string]any); app["id"] != "00000003-0000-0000-c000-000000000000" {
		t.Fatalf("unexpected authorized app %v", apps[0])
	}
}

func TestConstructResource_CallbackWithoutAppsSendsEmptyArray(t *testing.T) {
	model := minimalModel()
	model.CallbackConfiguration = &CallbackConfiguration{TimeoutDuration: types.StringValue("PT30M"), AuthorizedAppIds: types.SetNull(types.StringType)}

	body, err := constructResource(context.Background(), model, false)
	if err != nil {
		t.Fatalf("construct: %v", err)
	}
	callback, _ := serializeBody(t, body)["callbackConfiguration"].(map[string]any)

	apps, ok := callback["authorizedApps"].([]any)
	if !ok || len(apps) != 0 {
		t.Fatalf("authorizedApps must be an empty array, got %v", callback["authorizedApps"])
	}
}

func TestConstructResource_ValidationErrors(t *testing.T) {
	model := minimalModel()
	model.AuthenticationConfiguration = &AuthenticationConfiguration{Type: types.StringValue(authTypeAzureAd), ResourceId: types.StringNull()}
	if _, err := constructResource(context.Background(), model, false); err == nil {
		t.Fatal("expected an error for azureAdTokenAuthentication without resource_id")
	}

	model = minimalModel()
	model.EndpointConfiguration = nil
	if _, err := constructResource(context.Background(), model, false); err == nil {
		t.Fatal("expected an error without endpoint_configuration")
	}
}

func TestStateCallbackConfiguration_KeepsEquivalentDurationAndEmptySet(t *testing.T) {
	remote := graphmodels.NewCustomExtensionCallbackConfiguration()
	duration, err := serialization.ParseISODuration("PT1H")
	if err != nil {
		t.Fatalf("parse duration: %v", err)
	}
	remote.SetTimeoutDuration(duration)
	remote.SetAdditionalData(map[string]any{"authorizedApps": []any{}})

	prior := &CallbackConfiguration{TimeoutDuration: types.StringValue("PT60M"), AuthorizedAppIds: types.SetValueMust(types.StringType, []attr.Value{})}
	result := stateCallbackConfiguration(remote, prior)
	if result.TimeoutDuration.ValueString() != "PT60M" {
		t.Fatalf("expected the prior equivalent duration PT60M, got %s", result.TimeoutDuration.ValueString())
	}
	if result.AuthorizedAppIds.IsNull() || len(result.AuthorizedAppIds.Elements()) != 0 {
		t.Fatalf("expected an empty set, got %v", result.AuthorizedAppIds)
	}

	result = stateCallbackConfiguration(remote, nil)
	if result.TimeoutDuration.ValueString() != "PT1H" || !result.AuthorizedAppIds.IsNull() {
		t.Fatalf("unexpected mapping without prior state: %v", result)
	}
}

func TestStateCallbackConfiguration_ReadsAuthorizedAppsFromAdditionalData(t *testing.T) {
	remote := graphmodels.NewCustomExtensionCallbackConfiguration()
	id := "00000003-0000-0000-c000-000000000000"
	remote.SetAdditionalData(map[string]any{"authorizedApps": []any{map[string]any{"id": &id}}})

	result := stateCallbackConfiguration(remote, nil)
	if len(result.AuthorizedAppIds.Elements()) != 1 || result.AuthorizedAppIds.Elements()[0].(types.String).ValueString() != id {
		t.Fatalf("unexpected authorized_app_ids %v", result.AuthorizedAppIds)
	}
}
