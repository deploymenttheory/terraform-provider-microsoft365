package teamsMeetingPolicy

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// appendAllPolicyFieldsToCmd appends all non-null fields to the PowerShell command
func appendAllPolicyFieldsToCmd(ctx context.Context, data *TeamsMeetingPolicyResourceModel, cmd *[]string) {
	// Helper for booleans
	boolFlag := func(name string, v types.Bool) {
		if !v.IsNull() {
			*cmd = append(*cmd, fmt.Sprintf("-%s $%v", name, v.ValueBool()))
		}
	}
	// Helper for strings
	strFlag := func(name string, v types.String) {
		if !v.IsNull() {
			*cmd = append(*cmd, fmt.Sprintf("-%s '%s'", name, v.ValueString()))
		}
	}
	// Helper for int64
	intFlag := func(name string, v types.Int64) {
		if !v.IsNull() {
			*cmd = append(*cmd, fmt.Sprintf("-%s %d", name, v.ValueInt64()))
		}
	}
	// Helper for lists
	listFlag := func(name string, v types.List) {
		if !v.IsNull() && v.ElementType(ctx) == types.StringType {
			var arr []string
			v.ElementsAs(ctx, &arr, false)
			if len(arr) > 0 {
				*cmd = append(*cmd, fmt.Sprintf("-%s %s", name, strings.Join(arr, ",")))
			}
		}
	}
	// Add all fields (use PowerShell parameter names)
	strFlag("AIInterpreter", data.AIInterpreter)
	boolFlag("AllowAnnotations", data.AllowAnnotations)
	boolFlag("AllowAnonymousUsersToDialOut", data.AllowAnonymousUsersToDialOut)
	boolFlag("AllowAnonymousUsersToJoinMeeting", data.AllowAnonymousUsersToJoinMeeting)
	boolFlag("AllowAnonymousUsersToStartMeeting", data.AllowAnonymousUsersToStartMeeting)
	boolFlag("AllowAvatarsInGallery", data.AllowAvatarsInGallery)
	boolFlag("AllowBreakoutRooms", data.AllowBreakoutRooms)
	boolFlag("AllowCarbonSummary", data.AllowCarbonSummary)
	strFlag("AllowCartCaptionsScheduling", data.AllowCartCaptionsScheduling)
	boolFlag("AllowChannelMeetingScheduling", data.AllowChannelMeetingScheduling)
	boolFlag("AllowCloudRecording", data.AllowCloudRecording)
	strFlag("AllowDocumentCollaboration", data.AllowDocumentCollaboration)
	strFlag("AllowEngagementReport", data.AllowEngagementReport)
	boolFlag("AllowExternalNonTrustedMeetingChat", data.AllowExternalNonTrustedMeetingChat)
	boolFlag("AllowExternalParticipantGiveRequestControl", data.AllowExternalParticipantGiveRequestControl)
	boolFlag("AllowImmersiveView", data.AllowImmersiveView)
	boolFlag("AllowIPAudio", data.AllowIPAudio)
	boolFlag("AllowIPVideo", data.AllowIPVideo)
	boolFlag("AllowLocalRecording", data.AllowLocalRecording)
	boolFlag("AllowMeetingCoach", data.AllowMeetingCoach)
	boolFlag("AllowMeetNow", data.AllowMeetNow)
	boolFlag("AllowMeetingReactions", data.AllowMeetingReactions)
	boolFlag("AllowMeetingRegistration", data.AllowMeetingRegistration)
	boolFlag("AllowNDIStreaming", data.AllowNDIStreaming)
	boolFlag("AllowNetworkConfigurationSettingsLookup", data.AllowNetworkConfigurationSettingsLookup)
	boolFlag("AllowOrganizersToOverrideLobbySettings", data.AllowOrganizersToOverrideLobbySettings)
	boolFlag("AllowOutlookAddIn", data.AllowOutlookAddIn)
	boolFlag("AllowPSTNUsersToBypassLobby", data.AllowPSTNUsersToBypassLobby)
	boolFlag("AllowParticipantGiveRequestControl", data.AllowParticipantGiveRequestControl)
	boolFlag("AllowPowerPointSharing", data.AllowPowerPointSharing)
	boolFlag("AllowPrivateMeetNow", data.AllowPrivateMeetNow)
	boolFlag("AllowPrivateMeetingScheduling", data.AllowPrivateMeetingScheduling)
	boolFlag("AllowRecordingStorageOutsideRegion", data.AllowRecordingStorageOutsideRegion)
	boolFlag("AllowScreenContentDigitization", data.AllowScreenContentDigitization)
	boolFlag("AllowSharedNotes", data.AllowSharedNotes)
	strFlag("AllowTasksFromTranscript", data.AllowTasksFromTranscript)
	boolFlag("AllowTrackingInReport", data.AllowTrackingInReport)
	boolFlag("AllowTranscription", data.AllowTranscription)
	strFlag("AllowedUsersForMeetingContext", data.AllowedUsersForMeetingContext)
	strFlag("AllowUserToJoinExternalMeeting", data.AllowUserToJoinExternalMeeting)
	boolFlag("AllowWatermarkCustomizationForCameraVideo", data.AllowWatermarkCustomizationForCameraVideo)
	boolFlag("AllowWatermarkCustomizationForScreenSharing", data.AllowWatermarkCustomizationForScreenSharing)
	boolFlag("AllowWatermarkForCameraVideo", data.AllowWatermarkForCameraVideo)
	boolFlag("AllowWatermarkForScreenSharing", data.AllowWatermarkForScreenSharing)
	boolFlag("AllowWhiteboard", data.AllowWhiteboard)
	strFlag("AllowedStreamingMediaInput", data.AllowedStreamingMediaInput)
	strFlag("AnonymousUserAuthenticationMethod", data.AnonymousUserAuthenticationMethod)
	strFlag("AttendeeIdentityMasking", data.AttendeeIdentityMasking)
	strFlag("AudibleRecordingNotification", data.AudibleRecordingNotification)
	strFlag("AutoAdmittedUsers", data.AutoAdmittedUsers)
	strFlag("AutoRecording", data.AutoRecording)
	strFlag("AutomaticallyStartCopilot", data.AutomaticallyStartCopilot)
	listFlag("BlockedAnonymousJoinClientTypes", data.BlockedAnonymousJoinClientTypes)
	strFlag("CaptchaVerificationForMeetingJoin", data.CaptchaVerificationForMeetingJoin)
	strFlag("ChannelRecordingDownload", data.ChannelRecordingDownload)
	strFlag("ConnectToMeetingControls", data.ConnectToMeetingControls)
	strFlag("ContentSharingInExternalMeetings", data.ContentSharingInExternalMeetings)
	strFlag("Copilot", data.Copilot)
	boolFlag("CopyRestriction", data.CopyRestriction)
	strFlag("Description", data.Description)
	strFlag("DesignatedPresenterRoleMode", data.DesignatedPresenterRoleMode)
	boolFlag("DetectSensitiveContentDuringScreenSharing", data.DetectSensitiveContentDuringScreenSharing)
	strFlag("EnrollUserOverride", data.EnrollUserOverride)
	strFlag("ExplicitRecordingConsent", data.ExplicitRecordingConsent)
	strFlag("ExternalMeetingJoin", data.ExternalMeetingJoin)
	strFlag("InfoShownInReportMode", data.InfoShownInReportMode)
	strFlag("IPAudioMode", data.IPAudioMode)
	strFlag("IPVideoMode", data.IPVideoMode)
	strFlag("LiveCaptionsEnabledType", data.LiveCaptionsEnabledType)
	strFlag("LiveInterpretationEnabledType", data.LiveInterpretationEnabledType)
	strFlag("LiveStreamingMode", data.LiveStreamingMode)
	strFlag("LobbyChat", data.LobbyChat)
	intFlag("MediaBitRateKb", data.MediaBitRateKb)
	strFlag("MeetingChatEnabledType", data.MeetingChatEnabledType)
	strFlag("MeetingInviteLanguages", data.MeetingInviteLanguages)
	intFlag("NewMeetingRecordingExpirationDays", data.NewMeetingRecordingExpirationDays)
	strFlag("NoiseSuppressionForDialInParticipants", data.NoiseSuppressionForDialInParticipants)
	strFlag("ParticipantNameChange", data.ParticipantNameChange)
	strFlag("PreferredMeetingProviderForIslandsMode", data.PreferredMeetingProviderForIslandsMode)
	strFlag("QnAEngagementMode", data.QnAEngagementMode)
	strFlag("RecordingStorageMode", data.RecordingStorageMode)
	strFlag("RoomAttributeUserOverride", data.RoomAttributeUserOverride)
	strFlag("RoomPeopleNameUserOverride", data.RoomPeopleNameUserOverride)
	strFlag("ScreenSharingMode", data.ScreenSharingMode)
	strFlag("SmsNotifications", data.SmsNotifications)
	strFlag("SpeakerAttributionMode", data.SpeakerAttributionMode)
	strFlag("StreamingAttendeeMode", data.StreamingAttendeeMode)
	strFlag("TeamsCameraFarEndPTZMode", data.TeamsCameraFarEndPTZMode)
	strFlag("Tenant", data.Tenant)
	strFlag("UsersCanAdmitFromLobby", data.UsersCanAdmitFromLobby)
	strFlag("VideoFiltersMode", data.VideoFiltersMode)
	strFlag("VoiceIsolation", data.VoiceIsolation)
	strFlag("VoiceSimulationInInterpreter", data.VoiceSimulationInInterpreter)
	strFlag("WatermarkForAnonymousUsers", data.WatermarkForAnonymousUsers)
	intFlag("WatermarkForCameraVideoOpacity", data.WatermarkForCameraVideoOpacity)
	strFlag("WatermarkForCameraVideoPattern", data.WatermarkForCameraVideoPattern)
	intFlag("WatermarkForScreenSharingOpacity", data.WatermarkForScreenSharingOpacity)
	strFlag("WatermarkForScreenSharingPattern", data.WatermarkForScreenSharingPattern)
	strFlag("AllowedUsersForMeetingDetails", data.AllowedUsersForMeetingDetails)
	strFlag("RealTimeText", data.RealTimeText)
	strFlag("ParticipantSlideControl", data.ParticipantSlideControl)
	strFlag("WhoCanRegister", data.WhoCanRegister)
}
