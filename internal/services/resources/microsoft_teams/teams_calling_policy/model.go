// REF: https://learn.microsoft.com/en-us/powershell/module/teams/set-csteamscallingpolicy?view=teams-ps
// REF: https://learn.microsoft.com/en-us/powershell/module/teams/new-csteamscallingpolicy?view=teams-ps
package teams_calling_policy

import "github.com/hashicorp/terraform-plugin-framework/types"

type TeamsCallingPolicyResourceModel struct {
	ID                                   types.String `tfsdk:"id"`
	AIInterpreter                        types.String `tfsdk:"ai_interpreter"`
	AllowCallForwardingToPhone           types.Bool   `tfsdk:"allow_call_forwarding_to_phone"`
	AllowCallForwardingToUser            types.Bool   `tfsdk:"allow_call_forwarding_to_user"`
	AllowCallGroups                      types.Bool   `tfsdk:"allow_call_groups"`
	AllowCallRedirect                    types.String `tfsdk:"allow_call_redirect"`
	AllowCloudRecordingForCalls          types.Bool   `tfsdk:"allow_cloud_recording_for_calls"`
	AllowDelegation                      types.Bool   `tfsdk:"allow_delegation"`
	AllowPrivateCalling                  types.Bool   `tfsdk:"allow_private_calling"`
	AllowSIPDevicesCalling               types.Bool   `tfsdk:"allow_sip_devices_calling"`
	AllowTranscriptionForCalling         types.Bool   `tfsdk:"allow_transcription_for_calling"`
	AllowVoicemail                       types.String `tfsdk:"allow_voicemail"`
	AllowWebPSTNCalling                  types.Bool   `tfsdk:"allow_web_pstn_calling"`
	AutoAnswerEnabledType                types.String `tfsdk:"auto_answer_enabled_type"`
	BusyOnBusyEnabledType                types.String `tfsdk:"busy_on_busy_enabled_type"`
	CallingSpendUserLimit                types.Int64  `tfsdk:"calling_spend_user_limit"`
	CallRecordingExpirationDays          types.Int64  `tfsdk:"call_recording_expiration_days"`
	Copilot                              types.String `tfsdk:"copilot"`
	Description                          types.String `tfsdk:"description"`
	EnableSpendLimits                    types.Bool   `tfsdk:"enable_spend_limits"`
	EnableWebPstnMediaBypass             types.Bool   `tfsdk:"enable_web_pstn_media_bypass"`
	InboundFederatedCallRoutingTreatment types.String `tfsdk:"inbound_federated_call_routing_treatment"`
	InboundPstnCallRoutingTreatment      types.String `tfsdk:"inbound_pstn_call_routing_treatment"`
	LiveCaptionsEnabledTypeForCalling    types.String `tfsdk:"live_captions_enabled_type_for_calling"`
	MusicOnHoldEnabledType               types.String `tfsdk:"music_on_hold_enabled_type"`
	PopoutAppPathForIncomingPstnCalls    types.String `tfsdk:"popout_app_path_for_incoming_pstn_calls"`
	PopoutForIncomingPstnCalls           types.String `tfsdk:"popout_for_incoming_pstn_calls"`
	PreventTollBypass                    types.Bool   `tfsdk:"prevent_toll_bypass"`
	SpamFilteringEnabledType             types.String `tfsdk:"spam_filtering_enabled_type"`
	VoiceSimulationInInterpreter         types.String `tfsdk:"voice_simulation_in_interpreter"`
	RealTimeText                         types.String `tfsdk:"real_time_text"`
}
