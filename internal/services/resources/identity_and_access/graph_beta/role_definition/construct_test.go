package graphBetaIdentityAndAccessRoleDefinition

import (
	"context"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	jsonserialization "github.com/microsoft/kiota-serialization-json-go"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

func serializeRoleDefinition(t *testing.T, roleDefinition graphmodels.UnifiedRoleDefinitionable) string {
	t.Helper()
	writer := jsonserialization.NewJsonSerializationWriter()
	if err := writer.WriteObjectValue("", roleDefinition); err != nil {
		t.Fatalf("serialize role definition: %v", err)
	}
	content, err := writer.GetSerializedContent()
	if err != nil {
		t.Fatalf("read serialized role definition: %v", err)
	}
	return string(content)
}

func testPermissions(actions ...string) []RoleDefinitionPermission {
	values := make([]attr.Value, 0, len(actions))
	for _, action := range actions {
		values = append(values, types.StringValue(action))
	}
	return []RoleDefinitionPermission{{AllowedResourceActions: types.SetValueMust(types.StringType, values)}}
}

func testModel() *RoleDefinitionResourceModel {
	return &RoleDefinitionResourceModel{
		DisplayName:     types.StringValue("role"),
		Description:     types.StringNull(),
		IsEnabled:       types.BoolValue(true),
		TemplateId:      types.StringValue("33333333-3333-3333-3333-333333333333"),
		Version:         types.StringNull(),
		RolePermissions: testPermissions("microsoft.directory/applications/standard/read"),
	}
}

func buildBody(t *testing.T, model *RoleDefinitionResourceModel, isCreate bool) string {
	t.Helper()
	requestBody, err := constructRequestBody(context.Background(), model, isCreate)
	if err != nil {
		t.Fatalf("constructRequestBody returned error: %v", err)
	}
	return serializeRoleDefinition(t, requestBody)
}

// templateId becomes the role ID, so it is only sent on create.
func TestConstructRequestBodyTemplateIdOnlyOnCreate(t *testing.T) {
	model := testModel()

	if body := buildBody(t, model, true); !strings.Contains(body, `"templateId":"33333333-3333-3333-3333-333333333333"`) {
		t.Errorf("expected templateId on create, got %s", body)
	}
	if body := buildBody(t, model, false); strings.Contains(body, `"templateId"`) {
		t.Errorf("expected no templateId on update, got %s", body)
	}

	model.TemplateId = types.StringUnknown()
	if body := buildBody(t, model, true); strings.Contains(body, `"templateId"`) {
		t.Errorf("expected no templateId on create when unknown, got %s", body)
	}
}

// On create, unset optional strings are omitted rather than sent as null.
func TestConstructRequestBodyOmitsUnsetOptionalStringsOnCreate(t *testing.T) {
	body := buildBody(t, testModel(), true)
	for _, field := range []string{"description", "version"} {
		if strings.Contains(body, `"`+field+`"`) {
			t.Errorf("expected %s to be omitted on create, got %s", field, body)
		}
	}
	if !strings.Contains(body, `"rolePermissions"`) || !strings.Contains(body, `"isEnabled":true`) {
		t.Errorf("expected rolePermissions and isEnabled, got %s", body)
	}
}

// On update, PATCH keeps omitted properties, so removed optional strings are cleared with an explicit null (Graph rejects "").
func TestConstructRequestBodyClearsRemovedOptionalStringsOnUpdate(t *testing.T) {
	body := buildBody(t, testModel(), false)
	for _, field := range []string{`"description":null`, `"version":null`} {
		if !strings.Contains(body, field) {
			t.Errorf("expected %s on update, got %s", field, body)
		}
	}

	model := testModel()
	model.Description = types.StringValue("kept")
	model.Version = types.StringValue("2")
	body = buildBody(t, model, false)
	if !strings.Contains(body, `"description":"kept"`) || !strings.Contains(body, `"version":"2"`) {
		t.Errorf("expected configured description and version, got %s", body)
	}
}

func TestConstructRolePermissions(t *testing.T) {
	permissions, err := constructRolePermissions(context.Background(), testPermissions("microsoft.directory/applications/standard/read", "microsoft.directory/groups/standard/read"))
	if err != nil {
		t.Fatalf("constructRolePermissions returned error: %v", err)
	}
	if len(permissions) != 1 {
		t.Fatalf("expected 1 permission, got %d", len(permissions))
	}
	if got := len(permissions[0].GetAllowedResourceActions()); got != 2 {
		t.Errorf("expected 2 allowed resource actions, got %d", got)
	}
}

func TestConstructRolePermissionsRejectsNonStringActions(t *testing.T) {
	numbers := types.SetValueMust(types.Int64Type, []attr.Value{types.Int64Value(1)})
	model := testModel()
	model.RolePermissions = []RoleDefinitionPermission{{AllowedResourceActions: numbers}}
	_, err := constructRequestBody(context.Background(), model, true)
	if err == nil || !strings.Contains(err.Error(), "failed to convert allowed resource actions for role permission 0") {
		t.Fatalf("expected conversion error, got %v", err)
	}
}
