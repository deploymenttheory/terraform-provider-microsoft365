package graphBetaIdentityAndAccessDeviceRegistrationPolicy

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	jsonserialization "github.com/microsoft/kiota-serialization-json-go"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// remotePolicy returns a tenant policy like the one Graph returns when Intune manages registration.
func remotePolicy() graphmodels.DeviceRegistrationPolicyable {
	quota := int32(50)
	notRequired := graphmodels.NOTREQUIRED_MULTIFACTORAUTHCONFIGURATION
	enabled := true
	notConfigurable := false
	configurable := true

	policy := graphmodels.NewDeviceRegistrationPolicy()
	policy.SetUserDeviceQuota(&quota)
	policy.SetMultiFactorAuthConfiguration(&notRequired)

	laps := graphmodels.NewLocalAdminPasswordSettings()
	laps.SetIsEnabled(&enabled)
	policy.SetLocalAdminPassword(laps)

	registration := graphmodels.NewAzureADRegistrationPolicy()
	registration.SetIsAdminConfigurable(&notConfigurable)
	registration.SetAllowedToRegister(graphmodels.NewAllDeviceRegistrationMembership())
	policy.SetAzureADRegistration(registration)

	allowedToJoin := graphmodels.NewEnumeratedDeviceRegistrationMembership()
	allowedToJoin.SetUsers([]string{})
	allowedToJoin.SetGroups([]string{"11111111-1111-1111-1111-111111111111"})
	localAdmins := graphmodels.NewLocalAdminSettings()
	localAdmins.SetEnableGlobalAdmins(&enabled)
	localAdmins.SetRegisteringUsers(graphmodels.NewAllDeviceRegistrationMembership())
	join := graphmodels.NewAzureADJoinPolicy()
	join.SetIsAdminConfigurable(&configurable)
	join.SetAllowedToJoin(allowedToJoin)
	join.SetLocalAdmins(localAdmins)
	policy.SetAzureADJoin(join)

	return policy
}

// unknownModel returns a plan where nothing but the given overrides is configured.
func unknownModel() *DeviceRegistrationPolicyResourceModel {
	return &DeviceRegistrationPolicyResourceModel{
		ID:                           types.StringUnknown(),
		DisplayName:                  types.StringUnknown(),
		Description:                  types.StringUnknown(),
		UserDeviceQuota:              types.Int32Unknown(),
		MultiFactorAuthConfiguration: types.StringUnknown(),
		LocalAdminPassword:           types.ObjectUnknown(localAdminPasswordAttrTypes),
		AzureADRegistration:          types.ObjectUnknown(azureADRegistrationAttrTypes),
		AzureADJoin:                  types.ObjectUnknown(azureADJoinAttrTypes),
	}
}

func membershipObject(membershipType string, users, groups types.Set) types.Object {
	return types.ObjectValueMust(membershipAttrTypes, map[string]attr.Value{
		"type":   types.StringValue(membershipType),
		"users":  users,
		"groups": groups,
	})
}

func serialize(t *testing.T, body graphmodels.DeviceRegistrationPolicyable) map[string]any {
	t.Helper()
	writer := jsonserialization.NewJsonSerializationWriter()
	if err := writer.WriteObjectValue("", body); err != nil {
		t.Fatalf("serialize policy: %v", err)
	}
	content, err := writer.GetSerializedContent()
	if err != nil {
		t.Fatalf("read serialized policy: %v", err)
	}
	var result map[string]any
	if err := json.Unmarshal(content, &result); err != nil {
		t.Fatalf("parse serialized policy %s: %v", content, err)
	}
	return result
}

// An empty plan must send every updatable property back with its current value (full-replace PUT).
func TestUnitConstructResourceCopiesUnconfiguredValuesFromRemote(t *testing.T) {
	body, err := constructResource(context.Background(), unknownModel(), remotePolicy())
	if err != nil {
		t.Fatalf("constructResource returned error: %v", err)
	}
	result := serialize(t, body)

	if result["userDeviceQuota"] != float64(50) || result["multiFactorAuthConfiguration"] != "notRequired" {
		t.Errorf("expected quota and MFA to be copied, got %v", result)
	}
	if laps, _ := result["localAdminPassword"].(map[string]any); laps["isEnabled"] != true {
		t.Errorf("expected localAdminPassword to be copied, got %v", result["localAdminPassword"])
	}
	registration, _ := result["azureADRegistration"].(map[string]any)
	if registration["isAdminConfigurable"] != false {
		t.Errorf("expected isAdminConfigurable to be sent back, got %v", registration)
	}
	if allowed, _ := registration["allowedToRegister"].(map[string]any); allowed["@odata.type"] != odataTypeAllMembership {
		t.Errorf("expected allowedToRegister to be copied, got %v", registration)
	}
	join, _ := result["azureADJoin"].(map[string]any)
	allowedToJoin, _ := join["allowedToJoin"].(map[string]any)
	if allowedToJoin["@odata.type"] != odataTypeEnumeratedMembership {
		t.Errorf("expected enumerated allowedToJoin to be copied, got %v", allowedToJoin)
	}
	if groups, _ := allowedToJoin["groups"].([]any); len(groups) != 1 {
		t.Errorf("expected allowedToJoin groups to be copied, got %v", allowedToJoin)
	}
	localAdmins, _ := join["localAdmins"].(map[string]any)
	if localAdmins["enableGlobalAdmins"] != true {
		t.Errorf("expected localAdmins to be copied, got %v", localAdmins)
	}
	if _, ok := result["id"]; ok {
		t.Errorf("expected id to be omitted, got %v", result)
	}
}

// Configured values replace the remote values, and unset users and groups are omitted rather than sent as null.
func TestUnitConstructResourceAppliesConfiguredValues(t *testing.T) {
	data := unknownModel()
	data.UserDeviceQuota = types.Int32Value(10)
	data.MultiFactorAuthConfiguration = types.StringValue("required")
	data.LocalAdminPassword = types.ObjectValueMust(localAdminPasswordAttrTypes, map[string]attr.Value{"is_enabled": types.BoolValue(false)})
	data.AzureADJoin = types.ObjectValueMust(azureADJoinAttrTypes, map[string]attr.Value{
		"is_admin_configurable": types.BoolUnknown(),
		"allowed_to_join":       membershipObject(membershipTypeEnumerated, types.SetNull(types.StringType), types.SetValueMust(types.StringType, []attr.Value{types.StringValue("22222222-2222-2222-2222-222222222222")})),
		"local_admins":          types.ObjectUnknown(localAdminsAttrTypes),
	})

	body, err := constructResource(context.Background(), data, remotePolicy())
	if err != nil {
		t.Fatalf("constructResource returned error: %v", err)
	}
	result := serialize(t, body)

	if result["userDeviceQuota"] != float64(10) || result["multiFactorAuthConfiguration"] != "required" {
		t.Errorf("expected configured quota and MFA, got %v", result)
	}
	if laps, _ := result["localAdminPassword"].(map[string]any); laps["isEnabled"] != false {
		t.Errorf("expected configured localAdminPassword, got %v", result["localAdminPassword"])
	}
	join, _ := result["azureADJoin"].(map[string]any)
	if join["isAdminConfigurable"] != true {
		t.Errorf("expected azureADJoin.isAdminConfigurable to be copied, got %v", join)
	}
	allowedToJoin, _ := join["allowedToJoin"].(map[string]any)
	if _, ok := allowedToJoin["users"]; ok {
		t.Errorf("expected unset users to be omitted, got %v", allowedToJoin)
	}
	if groups, _ := allowedToJoin["groups"].([]any); len(groups) != 1 || groups[0] != "22222222-2222-2222-2222-222222222222" {
		t.Errorf("expected configured groups, got %v", allowedToJoin)
	}
	if localAdmins, _ := join["localAdmins"].(map[string]any); localAdmins["enableGlobalAdmins"] != true {
		t.Errorf("expected unconfigured localAdmins to be copied, got %v", join["localAdmins"])
	}
}

func TestUnitConstructMembershipTypes(t *testing.T) {
	cases := map[string]string{
		membershipTypeAll:  odataTypeAllMembership,
		membershipTypeNone: odataTypeNoMembership,
	}
	for membershipType, odataType := range cases {
		t.Run(membershipType, func(t *testing.T) {
			membership, err := constructMembership(context.Background(), membershipObject(membershipType, types.SetNull(types.StringType), types.SetNull(types.StringType)), nil)
			if err != nil {
				t.Fatalf("constructMembership returned error: %v", err)
			}
			if membership.GetOdataType() == nil || *membership.GetOdataType() != odataType {
				t.Errorf("expected %s, got %v", odataType, membership.GetOdataType())
			}
		})
	}

	_, err := constructMembership(context.Background(), membershipObject("unknownKind", types.SetNull(types.StringType), types.SetNull(types.StringType)), nil)
	if err == nil || !strings.Contains(err.Error(), "unsupported membership type") {
		t.Errorf("expected unsupported membership type error, got %v", err)
	}
}

// A configured registration membership that differs from the tenant fails when Intune manages registration.
func TestUnitConstructResourceRejectsRegistrationChangeWhenNotAdminConfigurable(t *testing.T) {
	data := unknownModel()
	data.AzureADRegistration = types.ObjectValueMust(azureADRegistrationAttrTypes, map[string]attr.Value{
		"is_admin_configurable": types.BoolUnknown(),
		"allowed_to_register":   membershipObject(membershipTypeNone, types.SetNull(types.StringType), types.SetNull(types.StringType)),
	})

	_, err := constructResource(context.Background(), data, remotePolicy())
	if err == nil || !strings.Contains(err.Error(), "isAdminConfigurable = false") {
		t.Fatalf("expected isAdminConfigurable error, got %v", err)
	}

	// The current value is accepted.
	data.AzureADRegistration = types.ObjectValueMust(azureADRegistrationAttrTypes, map[string]attr.Value{
		"is_admin_configurable": types.BoolUnknown(),
		"allowed_to_register":   membershipObject(membershipTypeAll, types.SetNull(types.StringType), types.SetNull(types.StringType)),
	})
	if _, err := constructResource(context.Background(), data, remotePolicy()); err != nil {
		t.Fatalf("expected unchanged registration to be accepted, got %v", err)
	}
}

func TestUnitConstructResourceRejectsUsersOnNonEnumeratedMembership(t *testing.T) {
	data := unknownModel()
	data.AzureADJoin = types.ObjectValueMust(azureADJoinAttrTypes, map[string]attr.Value{
		"is_admin_configurable": types.BoolUnknown(),
		"allowed_to_join":       membershipObject(membershipTypeNone, types.SetValueMust(types.StringType, []attr.Value{types.StringValue("33333333-3333-3333-3333-333333333333")}), types.SetNull(types.StringType)),
		"local_admins":          types.ObjectUnknown(localAdminsAttrTypes),
	})

	_, err := constructResource(context.Background(), data, remotePolicy())
	if err == nil || !strings.Contains(err.Error(), `users can only be set when type is "enumerated"`) {
		t.Fatalf("expected users validation error, got %v", err)
	}
}

func TestUnitConstructResourceRequiresRemotePolicy(t *testing.T) {
	if _, err := constructResource(context.Background(), unknownModel(), nil); err == nil {
		t.Fatal("expected an error without the current policy")
	}
}

// State mapping keeps an unset users list null when Graph returns it empty, and keeps an explicit empty set.
func TestUnitMapMembershipPreservesNullVersusEmpty(t *testing.T) {
	remote := graphmodels.NewEnumeratedDeviceRegistrationMembership()
	remote.SetUsers([]string{})
	remote.SetGroups([]string{"11111111-1111-1111-1111-111111111111"})

	mapped := mapMembership(context.Background(), nil, remote)
	if users := mapped.Attributes()["users"].(types.Set); !users.IsNull() {
		t.Errorf("expected null users without a prior value, got %v", users)
	}

	prior := membershipObject(membershipTypeEnumerated, types.SetValueMust(types.StringType, []attr.Value{}), types.SetNull(types.StringType))
	mapped = mapMembership(context.Background(), prior, remote)
	if users := mapped.Attributes()["users"].(types.Set); users.IsNull() || len(users.Elements()) != 0 {
		t.Errorf("expected an empty users set for an explicit empty prior value, got %v", users)
	}
	if mapped.Attributes()["type"].(types.String).ValueString() != membershipTypeEnumerated {
		t.Errorf("expected enumerated type, got %v", mapped)
	}
}
