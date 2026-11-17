package graphBetaM365AdminSharepointSettings

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	jsonserialization "github.com/microsoft/kiota-serialization-json-go"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// unknownModel returns a plan where nothing is configured.
func unknownModel() *SharepointSettingsResourceModel {
	return &SharepointSettingsResourceModel{
		ID:                                   types.StringUnknown(),
		AllowedDomainGuidsForSyncApp:         types.SetUnknown(types.StringType),
		AvailableManagedPathsForSiteCreation: types.SetUnknown(types.StringType),
		DeletedUserPersonalSiteRetentionPeriodInDays:    types.Int32Unknown(),
		ExcludedFileExtensionsForSyncApp:                types.SetUnknown(types.StringType),
		IdleSessionSignOut:                              types.ObjectUnknown(idleSessionSignOutAttrTypes),
		ImageTaggingOption:                              types.StringUnknown(),
		IsCommentingOnSitePagesEnabled:                  types.BoolUnknown(),
		IsFileActivityNotificationEnabled:               types.BoolUnknown(),
		IsLegacyAuthProtocolsEnabled:                    types.BoolUnknown(),
		IsLoopEnabled:                                   types.BoolUnknown(),
		IsMacSyncAppEnabled:                             types.BoolUnknown(),
		IsRequireAcceptingUserToMatchInvitedUserEnabled: types.BoolUnknown(),
		IsResharingByExternalUsersEnabled:               types.BoolUnknown(),
		IsSharePointMobileNotificationEnabled:           types.BoolUnknown(),
		IsSharePointNewsfeedEnabled:                     types.BoolUnknown(),
		IsSiteCreationEnabled:                           types.BoolUnknown(),
		IsSiteCreationUIEnabled:                         types.BoolUnknown(),
		IsSitePagesCreationEnabled:                      types.BoolUnknown(),
		IsSitesStorageLimitAutomatic:                    types.BoolUnknown(),
		IsSyncButtonHiddenOnPersonalSite:                types.BoolUnknown(),
		IsUnmanagedSyncAppForTenantRestricted:           types.BoolUnknown(),
		PersonalSiteDefaultStorageLimitInMB:             types.Int64Unknown(),
		SharingAllowedDomainList:                        types.SetUnknown(types.StringType),
		SharingBlockedDomainList:                        types.SetUnknown(types.StringType),
		SharingCapability:                               types.StringUnknown(),
		SharingDomainRestrictionMode:                    types.StringUnknown(),
		SiteCreationDefaultManagedPath:                  types.StringUnknown(),
		SiteCreationDefaultStorageLimitInMB:             types.Int32Unknown(),
		TenantDefaultTimezone:                           types.StringUnknown(),
	}
}

func idleSessionSignOutObject(isEnabled types.Bool, signOutAfter, warnAfter types.Int64) types.Object {
	return types.ObjectValueMust(idleSessionSignOutAttrTypes, map[string]attr.Value{
		"is_enabled":                isEnabled,
		"sign_out_after_in_seconds": signOutAfter,
		"warn_after_in_seconds":     warnAfter,
	})
}

func stringSet(values ...string) types.Set {
	elements := make([]attr.Value, 0, len(values))
	for _, value := range values {
		elements = append(elements, types.StringValue(value))
	}
	return types.SetValueMust(types.StringType, elements)
}

func serialize(t *testing.T, body graphmodels.SharepointSettingsable) map[string]any {
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

// Configured values are sent with their Graph names, and the read-only managed paths are never sent.
func TestUnitConstructRequestBodyAppliesConfiguredValues(t *testing.T) {
	data := unknownModel()
	data.AllowedDomainGuidsForSyncApp = stringSet("bdd1ab9b-3fd0-4def-a761-ec8d7471732c")
	data.AvailableManagedPathsForSiteCreation = stringSet("/sites/")
	data.DeletedUserPersonalSiteRetentionPeriodInDays = types.Int32Value(365)
	data.ExcludedFileExtensionsForSyncApp = stringSet()
	data.ImageTaggingOption = types.StringValue("enhanced")
	data.IsSharePointNewsfeedEnabled = types.BoolValue(true)
	data.PersonalSiteDefaultStorageLimitInMB = types.Int64Value(120000)
	data.SharingCapability = types.StringValue("existingExternalUserSharingOnly")
	data.SharingDomainRestrictionMode = types.StringValue("blockList")
	data.SharingBlockedDomainList = stringSet("fabrikam.com")
	data.IdleSessionSignOut = idleSessionSignOutObject(types.BoolValue(true), types.Int64Value(3600), types.Int64Value(1800))

	body, err := constructRequestBody(context.Background(), data)
	if err != nil {
		t.Fatalf("constructRequestBody returned error: %v", err)
	}
	result := serialize(t, body)

	if guids, _ := result["allowedDomainGuidsForSyncApp"].([]any); len(guids) != 1 || guids[0] != "bdd1ab9b-3fd0-4def-a761-ec8d7471732c" {
		t.Errorf("expected the configured sync app domain GUID, got %v", result["allowedDomainGuidsForSyncApp"])
	}
	if _, ok := result["availableManagedPathsForSiteCreation"]; ok {
		t.Errorf("expected read-only availableManagedPathsForSiteCreation to be omitted, got %v", result)
	}
	if result["deletedUserPersonalSiteRetentionPeriodInDays"] != float64(365) || result["personalSiteDefaultStorageLimitInMB"] != float64(120000) {
		t.Errorf("expected configured numbers, got %v", result)
	}
	if extensions, ok := result["excludedFileExtensionsForSyncApp"].([]any); !ok || len(extensions) != 0 {
		t.Errorf("expected an empty excludedFileExtensionsForSyncApp to be sent as [], got %v", result["excludedFileExtensionsForSyncApp"])
	}
	if result["imageTaggingOption"] != "enhanced" || result["sharingCapability"] != "existingExternalUserSharingOnly" || result["sharingDomainRestrictionMode"] != "blockList" {
		t.Errorf("expected configured enums, got %v", result)
	}
	if result["isSharePointNewsfeedEnabled"] != true {
		t.Errorf("expected isSharePointNewsfeedEnabled true, got %v", result)
	}
	if _, ok := result["sharingAllowedDomainList"]; ok {
		t.Errorf("expected unconfigured sharingAllowedDomainList to be omitted, got %v", result)
	}
	policy, _ := result["idleSessionSignOut"].(map[string]any)
	if policy["isEnabled"] != true || policy["signOutAfterInSeconds"] != float64(3600) || policy["warnAfterInSeconds"] != float64(1800) {
		t.Errorf("expected configured idleSessionSignOut, got %v", result["idleSessionSignOut"])
	}
	if _, ok := result["id"]; ok {
		t.Errorf("expected id to be omitted, got %v", result)
	}
}

// An idle_session_sign_out object with no known values is omitted, and partial objects only send known values.
func TestUnitConstructRequestBodyIdleSessionSignOut(t *testing.T) {
	data := unknownModel()
	data.IdleSessionSignOut = idleSessionSignOutObject(types.BoolUnknown(), types.Int64Null(), types.Int64Unknown())
	body, err := constructRequestBody(context.Background(), data)
	if err != nil {
		t.Fatalf("constructRequestBody returned error: %v", err)
	}
	if _, ok := serialize(t, body)["idleSessionSignOut"]; ok {
		t.Errorf("expected idleSessionSignOut to be omitted")
	}

	data.IdleSessionSignOut = idleSessionSignOutObject(types.BoolValue(false), types.Int64Unknown(), types.Int64Unknown())
	body, err = constructRequestBody(context.Background(), data)
	if err != nil {
		t.Fatalf("constructRequestBody returned error: %v", err)
	}
	policy, _ := serialize(t, body)["idleSessionSignOut"].(map[string]any)
	if len(policy) != 1 || policy["isEnabled"] != false {
		t.Errorf("expected only isEnabled false, got %v", policy)
	}
}

// An invalid sync app domain GUID is reported instead of being sent.
func TestUnitConstructRequestBodyRejectsInvalidDomainGuid(t *testing.T) {
	data := unknownModel()
	data.AllowedDomainGuidsForSyncApp = stringSet("not-a-guid")
	if _, err := constructRequestBody(context.Background(), data); err == nil {
		t.Errorf("expected an error for an invalid GUID")
	}
}

// State mapping uses the fixed singleton ID because Graph does not return one, and keeps empty collections as empty sets.
func TestUnitMapRemoteResourceStateToTerraform(t *testing.T) {
	retention := int32(30)
	enabled := true
	warn := int64(3300)
	signOut := int64(3600)
	capability := graphmodels.EXTERNALUSERANDGUESTSHARING_SHARINGCAPABILITIES
	policy := graphmodels.NewIdleSessionSignOut()
	policy.SetIsEnabled(&enabled)
	policy.SetWarnAfterInSeconds(&warn)
	policy.SetSignOutAfterInSeconds(&signOut)
	remote := graphmodels.NewSharepointSettings()
	remote.SetDeletedUserPersonalSiteRetentionPeriodInDays(&retention)
	remote.SetSharingCapability(&capability)
	remote.SetAllowedDomainGuidsForSyncApp([]uuid.UUID{uuid.MustParse("BDD1AB9B-3FD0-4DEF-A761-EC8D7471732C")})
	remote.SetSharingAllowedDomainList([]string{})
	remote.SetIdleSessionSignOut(policy)

	mapped := MapRemoteResourceStateToTerraform(context.Background(), *unknownModel(), remote)

	if mapped.ID.ValueString() != singletonID {
		t.Errorf("expected id %s, got %v", singletonID, mapped.ID)
	}
	if mapped.DeletedUserPersonalSiteRetentionPeriodInDays.ValueInt32() != 30 || mapped.SharingCapability.ValueString() != "externalUserAndGuestSharing" {
		t.Errorf("expected mapped scalar values, got %v and %v", mapped.DeletedUserPersonalSiteRetentionPeriodInDays, mapped.SharingCapability)
	}
	if !mapped.AllowedDomainGuidsForSyncApp.Equal(stringSet("bdd1ab9b-3fd0-4def-a761-ec8d7471732c")) {
		t.Errorf("expected lower-case sync app domain GUIDs, got %v", mapped.AllowedDomainGuidsForSyncApp)
	}
	if mapped.SharingAllowedDomainList.IsNull() || len(mapped.SharingAllowedDomainList.Elements()) != 0 {
		t.Errorf("expected an empty sharing allowed domain list, got %v", mapped.SharingAllowedDomainList)
	}
	if !mapped.IsRequireAcceptingUserToMatchInvitedUserEnabled.IsNull() {
		t.Errorf("expected a missing value to map to null, got %v", mapped.IsRequireAcceptingUserToMatchInvitedUserEnabled)
	}
	attributes := mapped.IdleSessionSignOut.Attributes()
	if !attributes["is_enabled"].(types.Bool).ValueBool() || attributes["warn_after_in_seconds"].(types.Int64).ValueInt64() != 3300 {
		t.Errorf("expected mapped idle session sign-out, got %v", mapped.IdleSessionSignOut)
	}

	remote.SetIdleSessionSignOut(nil)
	if mapped := MapRemoteResourceStateToTerraform(context.Background(), *unknownModel(), remote); !mapped.IdleSessionSignOut.IsNull() {
		t.Errorf("expected null idle session sign-out when Graph returns none, got %v", mapped.IdleSessionSignOut)
	}
}
