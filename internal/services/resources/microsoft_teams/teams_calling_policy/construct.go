package teams_calling_policy

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// appendAllCallingPolicyFieldsToCmd appends all non-null fields to the PowerShell command
func appendAllCallingPolicyFieldsToCmd(ctx context.Context, data *TeamsCallingPolicyResourceModel, cmd *[]string) {
	boolFlag := func(name string, v types.Bool) {
		if !v.IsNull() {
			*cmd = append(*cmd, fmt.Sprintf("-%s $%v", name, v.ValueBool()))
		}
	}
	strFlag := func(name string, v types.String) {
		if !v.IsNull() {
			*cmd = append(*cmd, fmt.Sprintf("-%s '%s'", name, v.ValueString()))
		}
	}
	intFlag := func(name string, v types.Int64) {
		if !v.IsNull() {
			*cmd = append(*cmd, fmt.Sprintf("-%s %d", name, v.ValueInt64()))
		}
	}
	// Add all fields (use PowerShell parameter names)
	strFlag("AIInterpreter", data.AIInterpreter)
	boolFlag("AllowCallForwardingToPhone", data.AllowCallForwardingToPhone)
	boolFlag("AllowCallForwardingToUser", data.AllowCallForwardingToUser)
	boolFlag("AllowCallGroups", data.AllowCallGroups)
	strFlag("AllowCallRedirect", data.AllowCallRedirect)
	boolFlag("AllowCloudRecordingForCalls", data.AllowCloudRecordingForCalls)
	boolFlag("AllowDelegation", data.AllowDelegation)
	boolFlag("AllowPrivateCalling", data.AllowPrivateCalling)
	boolFlag("AllowSIPDevicesCalling", data.AllowSIPDevicesCalling)
	boolFlag("AllowTranscriptionForCalling", data.AllowTranscriptionForCalling)
	strFlag("AllowVoicemail", data.AllowVoicemail)
	boolFlag("AllowWebPSTNCalling", data.AllowWebPSTNCalling)
	strFlag("AutoAnswerEnabledType", data.AutoAnswerEnabledType)
	strFlag("BusyOnBusyEnabledType", data.BusyOnBusyEnabledType)
	intFlag("CallingSpendUserLimit", data.CallingSpendUserLimit)
	intFlag("CallRecordingExpirationDays", data.CallRecordingExpirationDays)
	strFlag("Copilot", data.Copilot)
	strFlag("Description", data.Description)
	boolFlag("EnableSpendLimits", data.EnableSpendLimits)
	boolFlag("EnableWebPstnMediaBypass", data.EnableWebPstnMediaBypass)
	strFlag("InboundFederatedCallRoutingTreatment", data.InboundFederatedCallRoutingTreatment)
	strFlag("InboundPstnCallRoutingTreatment", data.InboundPstnCallRoutingTreatment)
	strFlag("LiveCaptionsEnabledTypeForCalling", data.LiveCaptionsEnabledTypeForCalling)
	strFlag("MusicOnHoldEnabledType", data.MusicOnHoldEnabledType)
	strFlag("PopoutAppPathForIncomingPstnCalls", data.PopoutAppPathForIncomingPstnCalls)
	strFlag("PopoutForIncomingPstnCalls", data.PopoutForIncomingPstnCalls)
	boolFlag("PreventTollBypass", data.PreventTollBypass)
	strFlag("SpamFilteringEnabledType", data.SpamFilteringEnabledType)
	strFlag("VoiceSimulationInInterpreter", data.VoiceSimulationInInterpreter)
	strFlag("RealTimeText", data.RealTimeText)
}
