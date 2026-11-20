package graphBetaWindowsAutopilotDevicePreparationPolicy

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"

	builders "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors/graph_beta/device_management"
)

func TestUnitResourceWindowsAutopilotDevicePreparationPolicy_ImportOOBESettings(t *testing.T) {
	for _, test := range []struct {
		name   string
		suffix string
		want   bool
	}{
		{"false choices", "_0", false},
		{"true choices", "_1", true},
	} {
		t.Run(test.name, func(t *testing.T) {
			state := WindowsAutopilotDevicePreparationPolicyResourceModel{
				TemplateId: types.StringValue(TemplateIDUserDriven),
			}
			settings := models.NewDeviceManagementConfigurationSettingCollectionResponse()
			settings.SetValue([]models.DeviceManagementConfigurationSettingable{
				builders.ConstructChoiceSettingInstance(
					SettingDefAllowSkip, SettingDefAllowSkip+test.suffix, "", "",
				),
				builders.ConstructChoiceSettingInstance(
					SettingDefAllowDiagnostics, SettingDefAllowDiagnostics+test.suffix, "", "",
				),
				builders.ConstructStringSimpleSettingInstance(
					SettingDefCustomErrorMessage, "Contact support for help.", "", "",
				),
			})
			if err := mapSettingsToState(t.Context(), &state, settings); err != nil {
				t.Fatal(err)
			}
			for name, actual := range map[string]types.Bool{
				"allow_skip":        state.OOBESettings.AllowSkip,
				"allow_diagnostics": state.OOBESettings.AllowDiagnostics,
			} {
				if !actual.Equal(types.BoolValue(test.want)) {
					t.Errorf("%s = %s, want %t", name, actual.String(), test.want)
				}
			}
			if got := state.OOBESettings.CustomErrorMessage.ValueString(); got != "Contact support for help." {
				t.Errorf("custom_error_message = %q, want API setting value", got)
			}
		})
	}
}

func TestUnitResourceWindowsAutopilotDevicePreparationPolicy_MembershipTargetState(t *testing.T) {
	for _, test := range []struct {
		name    string
		body    string
		want    types.String
		wantErr bool
	}{
		{
			name: "configured target",
			body: `{"enrollmentTimeDeviceMembershipTargetValidationStatuses":[{"targetId":"00000000-0000-0000-0000-000000000001","validationSucceeded":true}]}`,
			want: types.StringValue("00000000-0000-0000-0000-000000000001"),
		},
		{
			name: "removed target",
			body: `{"enrollmentTimeDeviceMembershipTargetValidationStatuses":[]}`,
			want: types.StringNull(),
		},
		{
			name:    "Graph error response preserves state",
			body:    `{"error":{"code":"ResourceNotFound","message":"Invalid version"}}`,
			want:    types.StringValue("previous-target"),
			wantErr: true,
		},
		{
			name:    "invalid response preserves state",
			body:    `{`,
			want:    types.StringValue("previous-target"),
			wantErr: true,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			state := WindowsAutopilotDevicePreparationPolicyResourceModel{
				DeviceSecurityGroup: types.StringValue("previous-target"),
			}
			err := mapEnrollmentTimeDeviceMembershipTargetToState(
				t.Context(),
				&state,
				[]byte(test.body),
			)
			if (err != nil) != test.wantErr {
				t.Fatalf("mapping error = %v, want error %t", err, test.wantErr)
			}
			if !state.DeviceSecurityGroup.Equal(test.want) {
				t.Errorf(
					"device_security_group = %s, want %s",
					state.DeviceSecurityGroup,
					test.want,
				)
			}
		})
	}
}
