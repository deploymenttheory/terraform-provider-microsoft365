package graphBetaTargetedManagedAppConfigurations

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/microsoft/kiota-abstractions-go/serialization"
	"github.com/microsoft/kiota-abstractions-go/store"
	jsonserialization "github.com/microsoft/kiota-serialization-json-go"
)

// emptyChildrenSettingsCatalog holds a choice, a choice collection, a group collection and a group with nested choice
// children, none of which have children of their own.
func emptyChildrenSettingsCatalog() *DeviceConfigV2GraphServiceResourceModel {
	choice := func(value string) *ChoiceSettingStruct {
		return &ChoiceSettingStruct{Value: types.StringValue(value), Children: []ChoiceSettingChild{}}
	}
	choiceCollection := func(value string) []ChoiceSettingCollectionStruct {
		return []ChoiceSettingCollectionStruct{{Value: types.StringValue(value), Children: []ChoiceSettingCollectionChild{}}}
	}
	return &DeviceConfigV2GraphServiceResourceModel{Settings: []Setting{
		{ID: types.StringValue("0"), SettingInstance: SettingInstance{
			ODataType:           types.StringValue("#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"),
			SettingDefinitionId: types.StringValue("com.microsoft.intune.mam.managedbrowser.tunnelavailable"),
			ChoiceSettingValue:  choice("com.microsoft.intune.mam.managedbrowser.tunnelavailable_true"),
		}},
		{ID: types.StringValue("1"), SettingInstance: SettingInstance{
			ODataType:                    types.StringValue("#microsoft.graph.deviceManagementConfigurationChoiceSettingCollectionInstance"),
			SettingDefinitionId:          types.StringValue("com.microsoft.intune.mam.choicecollection"),
			ChoiceSettingCollectionValue: choiceCollection("com.microsoft.intune.mam.choicecollection_1"),
		}},
		{ID: types.StringValue("2"), SettingInstance: SettingInstance{
			ODataType:                   types.StringValue("#microsoft.graph.deviceManagementConfigurationGroupSettingCollectionInstance"),
			SettingDefinitionId:         types.StringValue("com.microsoft.intune.mam.groupcollection"),
			GroupSettingCollectionValue: []GroupSettingCollectionStruct{{Children: []GroupSettingCollectionChild{}}},
		}},
		{ID: types.StringValue("3"), SettingInstance: SettingInstance{
			ODataType:           types.StringValue("#microsoft.graph.deviceManagementConfigurationGroupSettingCollectionInstance"),
			SettingDefinitionId: types.StringValue("com.microsoft.intune.mam.nestedgroup"),
			GroupSettingCollectionValue: []GroupSettingCollectionStruct{{Children: []GroupSettingCollectionChild{
				{
					ODataType:           types.StringValue("#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"),
					SettingDefinitionId: types.StringValue("com.microsoft.intune.mam.nestedgroup_choice"),
					ChoiceSettingValue:  choice("com.microsoft.intune.mam.nestedgroup_choice_1"),
				},
				{
					ODataType:                    types.StringValue("#microsoft.graph.deviceManagementConfigurationChoiceSettingCollectionInstance"),
					SettingDefinitionId:          types.StringValue("com.microsoft.intune.mam.nestedgroup_choicecollection"),
					ChoiceSettingCollectionValue: choiceCollection("com.microsoft.intune.mam.nestedgroup_choicecollection_1"),
				},
				{
					ODataType:           types.StringValue("#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance"),
					SettingDefinitionId: types.StringValue("com.microsoft.intune.mam.nestedgroup_name"),
					SimpleSettingValue: &SimpleSettingStruct{
						ODataType: types.StringValue("#microsoft.graph.deviceManagementConfigurationStringSettingValue"),
						Value:     types.StringValue("contoso"),
					},
				},
			}}},
		}},
	}}
}

// assertRequestFixture serialises body with the backing-store proxy factory that the Graph request adapter uses, so a
// children slice set to nil is written as "children": null as it is on the wire, and compares it with the JSON fixture
// at path. With UPDATE_REQUEST_FIXTURES=1 it writes the fixture instead.
func assertRequestFixture(t *testing.T, body serialization.Parsable, path string) {
	t.Helper()
	factory := store.NewBackingStoreSerializationWriterProxyFactory(jsonserialization.NewJsonSerializationWriterFactory())
	writer, err := factory.GetSerializationWriter("application/json")
	if err != nil {
		t.Fatal(err)
	}
	defer writer.Close()
	if err := writer.WriteObjectValue("", body); err != nil {
		t.Fatal(err)
	}
	content, err := writer.GetSerializedContent()
	if err != nil {
		t.Fatal(err)
	}
	var actual bytes.Buffer
	if err := json.Indent(&actual, content, "", "  "); err != nil {
		t.Fatalf("indent %s: %v", content, err)
	}
	actual.WriteByte('\n')
	if os.Getenv(mocks.UpdateRequestFixturesEnv) == "1" {
		if err := os.MkdirAll(filepath.Dir(path), 0o750); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, actual.Bytes(), 0o600); err != nil {
			t.Fatal(err)
		}
		return
	}
	expected, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read request fixture %s (set %s=1 to create it): %v", path, mocks.UpdateRequestFixturesEnv, err)
	}
	expected = bytes.ReplaceAll(expected, []byte("\r\n"), []byte("\n"))
	if !bytes.Equal(expected, actual.Bytes()) {
		t.Fatalf("request body differs from %s\n--- expected\n%s\n--- actual\n%s", path, expected, actual.Bytes())
	}
}

// TestUnitResourceTargetedManagedAppConfiguration_01_EmptyChildrenRequests checks the create (POST) and changeSettings
// bodies for settings without children. Each such value must be sent as "children": [], because Graph rejects
// "children": null (#4956).
func TestUnitResourceTargetedManagedAppConfiguration_01_EmptyChildrenRequests(t *testing.T) {
	data := &TargetedManagedAppConfigurationResourceModel{
		DisplayName:     types.StringValue("Test Empty Children Configuration"),
		Description:     types.StringValue("Test configuration for settings without children"),
		RoleScopeTagIds: types.SetNull(types.StringType),
		AppGroupType:    types.StringValue("selectedPublicApps"),
		Apps:            types.SetNull(types.ObjectType{}),
		Assignments:     types.SetNull(types.ObjectType{}),
		SettingsCatalog: emptyChildrenSettingsCatalog(),
	}

	createBody, err := constructResource(context.Background(), data)
	if err != nil {
		t.Fatal(err)
	}
	assertRequestFixture(t, createBody, "tests/requests/empty_children_create.json")

	settingsBody, err := constructSettingsUpdate(context.Background(), data)
	if err != nil {
		t.Fatal(err)
	}
	assertRequestFixture(t, settingsBody, "tests/requests/empty_children_change_settings.json")
}
