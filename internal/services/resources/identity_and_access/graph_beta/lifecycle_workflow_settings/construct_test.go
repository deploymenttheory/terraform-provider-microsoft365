package graphBetaIdentityAndAccessLifecycleWorkflowSettings

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	jsonserialization "github.com/microsoft/kiota-serialization-json-go"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
	identitygovernance "github.com/microsoftgraph/msgraph-beta-sdk-go/models/identitygovernance"
)

// unknownModel returns a plan where nothing is configured.
func unknownModel() *LifecycleWorkflowSettingsResourceModel {
	return &LifecycleWorkflowSettingsResourceModel{
		ID:                              types.StringUnknown(),
		WorkflowScheduleIntervalInHours: types.Int32Unknown(),
		EmailSettings:                   types.ObjectUnknown(emailSettingsAttrTypes),
	}
}

func emailSettingsObject(senderDomain types.String, useCompanyBranding types.Bool) types.Object {
	return types.ObjectValueMust(emailSettingsAttrTypes, map[string]attr.Value{
		"sender_domain":        senderDomain,
		"use_company_branding": useCompanyBranding,
	})
}

func serialize(t *testing.T, body identitygovernance.LifecycleManagementSettingsable) map[string]any {
	t.Helper()
	writer := jsonserialization.NewJsonSerializationWriter()
	if err := writer.WriteObjectValue("", body); err != nil {
		t.Fatalf("serialize settings: %v", err)
	}
	content, err := writer.GetSerializedContent()
	if err != nil {
		t.Fatalf("read serialized settings: %v", err)
	}
	var result map[string]any
	if err := json.Unmarshal(content, &result); err != nil {
		t.Fatalf("parse serialized settings %s: %v", content, err)
	}
	return result
}

// An empty plan sends an empty PATCH body: Graph keeps every current value.
func TestUnitConstructRequestBodyOmitsUnconfiguredValues(t *testing.T) {
	body, err := constructRequestBody(context.Background(), unknownModel())
	if err != nil {
		t.Fatalf("constructRequestBody returned error: %v", err)
	}
	result := serialize(t, body)
	if len(result) != 0 {
		t.Errorf("expected an empty body, got %v", result)
	}
}

// Configured values are sent, and emailSettings has no @odata.type (Graph rejects microsoft.graph.emailSettings).
func TestUnitConstructRequestBodyAppliesConfiguredValues(t *testing.T) {
	data := unknownModel()
	data.WorkflowScheduleIntervalInHours = types.Int32Value(8)
	data.EmailSettings = emailSettingsObject(types.StringValue("contoso.com"), types.BoolValue(true))

	body, err := constructRequestBody(context.Background(), data)
	if err != nil {
		t.Fatalf("constructRequestBody returned error: %v", err)
	}
	result := serialize(t, body)

	if result["workflowScheduleIntervalInHours"] != float64(8) {
		t.Errorf("expected workflowScheduleIntervalInHours 8, got %v", result)
	}
	emailSettings, _ := result["emailSettings"].(map[string]any)
	if emailSettings["senderDomain"] != "contoso.com" || emailSettings["useCompanyBranding"] != true {
		t.Errorf("expected configured emailSettings, got %v", result["emailSettings"])
	}
	if _, ok := emailSettings["@odata.type"]; ok {
		t.Errorf("expected emailSettings without @odata.type, got %v", emailSettings)
	}
	if _, ok := result["@odata.type"]; ok {
		t.Errorf("expected body without @odata.type, got %v", result)
	}
	if _, ok := result["id"]; ok {
		t.Errorf("expected id to be omitted, got %v", result)
	}
}

// Graph rejects a null senderDomain, so an unset sender domain is omitted rather than sent as null.
func TestUnitConstructRequestBodyOmitsUnknownSenderDomain(t *testing.T) {
	data := unknownModel()
	data.EmailSettings = emailSettingsObject(types.StringUnknown(), types.BoolValue(false))

	body, err := constructRequestBody(context.Background(), data)
	if err != nil {
		t.Fatalf("constructRequestBody returned error: %v", err)
	}
	result := serialize(t, body)

	emailSettings, _ := result["emailSettings"].(map[string]any)
	if _, ok := emailSettings["senderDomain"]; ok {
		t.Errorf("expected unknown senderDomain to be omitted, got %v", emailSettings)
	}
	if emailSettings["useCompanyBranding"] != false {
		t.Errorf("expected useCompanyBranding false, got %v", emailSettings)
	}
}

// Graph rejects a null emailSettings, so an email_settings object with no known values is omitted.
func TestUnitConstructRequestBodyOmitsEmptyEmailSettings(t *testing.T) {
	data := unknownModel()
	data.EmailSettings = emailSettingsObject(types.StringNull(), types.BoolUnknown())

	body, err := constructRequestBody(context.Background(), data)
	if err != nil {
		t.Fatalf("constructRequestBody returned error: %v", err)
	}
	if _, ok := serialize(t, body)["emailSettings"]; ok {
		t.Errorf("expected emailSettings to be omitted")
	}
}

// State mapping uses the fixed singleton ID because Graph does not return one.
func TestUnitMapRemoteResourceStateToTerraform(t *testing.T) {
	interval := int32(4)
	domain := "microsoft.com"
	branding := true
	emailSettings := graphmodels.NewEmailSettings()
	emailSettings.SetSenderDomain(&domain)
	emailSettings.SetUseCompanyBranding(&branding)
	remote := identitygovernance.NewLifecycleManagementSettings()
	remote.SetWorkflowScheduleIntervalInHours(&interval)
	remote.SetEmailSettings(emailSettings)

	mapped := MapRemoteResourceStateToTerraform(context.Background(), *unknownModel(), remote)

	if mapped.ID.ValueString() != singletonID {
		t.Errorf("expected id %s, got %v", singletonID, mapped.ID)
	}
	if mapped.WorkflowScheduleIntervalInHours.ValueInt32() != 4 {
		t.Errorf("expected interval 4, got %v", mapped.WorkflowScheduleIntervalInHours)
	}
	attributes := mapped.EmailSettings.Attributes()
	if attributes["sender_domain"].(types.String).ValueString() != domain || !attributes["use_company_branding"].(types.Bool).ValueBool() {
		t.Errorf("expected mapped email settings, got %v", mapped.EmailSettings)
	}

	remote.SetEmailSettings(nil)
	if mapped := MapRemoteResourceStateToTerraform(context.Background(), *unknownModel(), remote); !mapped.EmailSettings.IsNull() {
		t.Errorf("expected null email settings when Graph returns none, got %v", mapped.EmailSettings)
	}
}
