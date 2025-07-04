package teams_calling_policy

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// mapAllCallingPolicyFieldsFromJson maps all fields from the JSON output to the resource model
func mapAllCallingPolicyFieldsFromJson(ctx context.Context, p map[string]interface{}, data *TeamsCallingPolicyResourceModel) {
	boolVal := func(name string) types.Bool {
		if v, ok := p[name]; ok && v != nil {
			if b, ok := v.(bool); ok {
				return types.BoolValue(b)
			}
		}
		return types.BoolNull()
	}
	strVal := func(name string) types.String {
		if v, ok := p[name]; ok && v != nil {
			if s, ok := v.(string); ok {
				return types.StringValue(s)
			}
		}
		return types.StringNull()
	}
	intVal := func(name string) types.Int64 {
		if v, ok := p[name]; ok && v != nil {
			switch t := v.(type) {
			case float64:
				return types.Int64Value(int64(t))
			case int64:
				return types.Int64Value(t)
			}
		}
		return types.Int64Null()
	}
	// Map all fields
	data.AIInterpreter = strVal("AIInterpreter")
	data.AllowCallForwardingToPhone = boolVal("AllowCallForwardingToPhone")
	data.AllowCallForwardingToUser = boolVal("AllowCallForwardingToUser")
	data.AllowCallGroups = boolVal("AllowCallGroups")
	data.AllowCallRedirect = strVal("AllowCallRedirect")
	data.AllowCloudRecordingForCalls = boolVal("AllowCloudRecordingForCalls")
	data.AllowDelegation = boolVal("AllowDelegation")
	data.AllowPrivateCalling = boolVal("AllowPrivateCalling")
	data.AllowSIPDevicesCalling = boolVal("AllowSIPDevicesCalling")
	data.AllowTranscriptionForCalling = boolVal("AllowTranscriptionForCalling")
	data.AllowVoicemail = strVal("AllowVoicemail")
	data.AllowWebPSTNCalling = boolVal("AllowWebPSTNCalling")
	data.AutoAnswerEnabledType = strVal("AutoAnswerEnabledType")
	data.BusyOnBusyEnabledType = strVal("BusyOnBusyEnabledType")
	data.CallingSpendUserLimit = intVal("CallingSpendUserLimit")
	data.CallRecordingExpirationDays = intVal("CallRecordingExpirationDays")
	data.Copilot = strVal("Copilot")
	data.Description = strVal("Description")
	data.EnableSpendLimits = boolVal("EnableSpendLimits")
	data.EnableWebPstnMediaBypass = boolVal("EnableWebPstnMediaBypass")
	data.InboundFederatedCallRoutingTreatment = strVal("InboundFederatedCallRoutingTreatment")
	data.InboundPstnCallRoutingTreatment = strVal("InboundPstnCallRoutingTreatment")
	data.LiveCaptionsEnabledTypeForCalling = strVal("LiveCaptionsEnabledTypeForCalling")
	data.MusicOnHoldEnabledType = strVal("MusicOnHoldEnabledType")
	data.PopoutAppPathForIncomingPstnCalls = strVal("PopoutAppPathForIncomingPstnCalls")
	data.PopoutForIncomingPstnCalls = strVal("PopoutForIncomingPstnCalls")
	data.PreventTollBypass = boolVal("PreventTollBypass")
	data.SpamFilteringEnabledType = strVal("SpamFilteringEnabledType")
	data.VoiceSimulationInInterpreter = strVal("VoiceSimulationInInterpreter")
	data.RealTimeText = strVal("RealTimeText")
}
