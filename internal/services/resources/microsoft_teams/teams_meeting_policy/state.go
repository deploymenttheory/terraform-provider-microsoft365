package teamsMeetingPolicy

import (
	"context"

	attr "github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// mapAllPolicyFieldsFromJson maps all fields from the JSON output to the resource model
func mapAllPolicyFieldsFromJson(ctx context.Context, p map[string]any, data *TeamsMeetingPolicyResourceModel) {
	// Helper for booleans
	boolVal := func(name string) types.Bool {
		if v, ok := p[name]; ok && v != nil {
			if b, ok := v.(bool); ok {
				return types.BoolValue(b)
			}
		}
		return types.BoolNull()
	}
	// Helper for strings
	strVal := func(name string) types.String {
		if v, ok := p[name]; ok && v != nil {
			if s, ok := v.(string); ok {
				return types.StringValue(s)
			}
		}
		return types.StringNull()
	}
	// Helper for int64
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
	// Helper for lists
	listVal := func(name string) types.List {
		if v, ok := p[name]; ok && v != nil {
			if arr, ok := v.([]interface{}); ok {
				var result []attr.Value
				for _, item := range arr {
					if s, ok := item.(string); ok {
						result = append(result, types.StringValue(s))
					}
				}
				return types.ListValueMust(types.StringType, result)
			}
		}
		return types.ListNull(types.StringType)
	}
	// Map all fields
	data.AIInterpreter = strVal("AIInterpreter")
	data.AllowAnnotations = boolVal("AllowAnnotations")
	data.AllowAnonymousUsersToDialOut = boolVal("AllowAnonymousUsersToDialOut")
	data.AllowAnonymousUsersToJoinMeeting = boolVal("AllowAnonymousUsersToJoinMeeting")
	data.AllowAnonymousUsersToStartMeeting = boolVal("AllowAnonymousUsersToStartMeeting")
	data.AllowAvatarsInGallery = boolVal("AllowAvatarsInGallery")
	data.AllowBreakoutRooms = boolVal("AllowBreakoutRooms")
	data.AllowCarbonSummary = boolVal("AllowCarbonSummary")
	data.AllowCartCaptionsScheduling = strVal("AllowCartCaptionsScheduling")
	data.AllowChannelMeetingScheduling = boolVal("AllowChannelMeetingScheduling")
	data.AllowCloudRecording = boolVal("AllowCloudRecording")
	data.AllowDocumentCollaboration = strVal("AllowDocumentCollaboration")
	data.AllowEngagementReport = strVal("AllowEngagementReport")
	data.AllowExternalNonTrustedMeetingChat = boolVal("AllowExternalNonTrustedMeetingChat")
	data.AllowExternalParticipantGiveRequestControl = boolVal("AllowExternalParticipantGiveRequestControl")
	data.AllowImmersiveView = boolVal("AllowImmersiveView")
	data.AllowIPAudio = boolVal("AllowIPAudio")
	data.AllowIPVideo = boolVal("AllowIPVideo")
	data.AllowLocalRecording = boolVal("AllowLocalRecording")
	data.AllowMeetingCoach = boolVal("AllowMeetingCoach")
	data.AllowMeetNow = boolVal("AllowMeetNow")
	data.AllowMeetingReactions = boolVal("AllowMeetingReactions")
	data.AllowMeetingRegistration = boolVal("AllowMeetingRegistration")
	data.AllowNDIStreaming = boolVal("AllowNDIStreaming")
	data.AllowNetworkConfigurationSettingsLookup = boolVal("AllowNetworkConfigurationSettingsLookup")
	data.AllowOrganizersToOverrideLobbySettings = boolVal("AllowOrganizersToOverrideLobbySettings")
	data.AllowOutlookAddIn = boolVal("AllowOutlookAddIn")
	data.AllowPSTNUsersToBypassLobby = boolVal("AllowPSTNUsersToBypassLobby")
	data.AllowParticipantGiveRequestControl = boolVal("AllowParticipantGiveRequestControl")
	data.AllowPowerPointSharing = boolVal("AllowPowerPointSharing")
	data.AllowPrivateMeetNow = boolVal("AllowPrivateMeetNow")
	data.AllowPrivateMeetingScheduling = boolVal("AllowPrivateMeetingScheduling")
	data.AllowRecordingStorageOutsideRegion = boolVal("AllowRecordingStorageOutsideRegion")
	data.AllowScreenContentDigitization = boolVal("AllowScreenContentDigitization")
	data.AllowSharedNotes = boolVal("AllowSharedNotes")
	data.AllowTasksFromTranscript = strVal("AllowTasksFromTranscript")
	data.AllowTrackingInReport = boolVal("AllowTrackingInReport")
	data.AllowTranscription = boolVal("AllowTranscription")
	data.AllowedUsersForMeetingContext = strVal("AllowedUsersForMeetingContext")
	data.AllowUserToJoinExternalMeeting = strVal("AllowUserToJoinExternalMeeting")
	data.AllowWatermarkCustomizationForCameraVideo = boolVal("AllowWatermarkCustomizationForCameraVideo")
	data.AllowWatermarkCustomizationForScreenSharing = boolVal("AllowWatermarkCustomizationForScreenSharing")
	data.AllowWatermarkForCameraVideo = boolVal("AllowWatermarkForCameraVideo")
	data.AllowWatermarkForScreenSharing = boolVal("AllowWatermarkForScreenSharing")
	data.AllowWhiteboard = boolVal("AllowWhiteboard")
	data.AllowedStreamingMediaInput = strVal("AllowedStreamingMediaInput")
	data.AnonymousUserAuthenticationMethod = strVal("AnonymousUserAuthenticationMethod")
	data.AttendeeIdentityMasking = strVal("AttendeeIdentityMasking")
	data.AudibleRecordingNotification = strVal("AudibleRecordingNotification")
	data.AutoAdmittedUsers = strVal("AutoAdmittedUsers")
	data.AutoRecording = strVal("AutoRecording")
	data.AutomaticallyStartCopilot = strVal("AutomaticallyStartCopilot")
	data.BlockedAnonymousJoinClientTypes = listVal("BlockedAnonymousJoinClientTypes")
	data.CaptchaVerificationForMeetingJoin = strVal("CaptchaVerificationForMeetingJoin")
	data.ChannelRecordingDownload = strVal("ChannelRecordingDownload")
	data.ConnectToMeetingControls = strVal("ConnectToMeetingControls")
	data.ContentSharingInExternalMeetings = strVal("ContentSharingInExternalMeetings")
	data.Copilot = strVal("Copilot")
	data.CopyRestriction = boolVal("CopyRestriction")
	data.Description = strVal("Description")
	data.DesignatedPresenterRoleMode = strVal("DesignatedPresenterRoleMode")
	data.DetectSensitiveContentDuringScreenSharing = boolVal("DetectSensitiveContentDuringScreenSharing")
	data.EnrollUserOverride = strVal("EnrollUserOverride")
	data.ExplicitRecordingConsent = strVal("ExplicitRecordingConsent")
	data.ExternalMeetingJoin = strVal("ExternalMeetingJoin")
	data.InfoShownInReportMode = strVal("InfoShownInReportMode")
	data.IPAudioMode = strVal("IPAudioMode")
	data.IPVideoMode = strVal("IPVideoMode")
	data.LiveCaptionsEnabledType = strVal("LiveCaptionsEnabledType")
	data.LiveInterpretationEnabledType = strVal("LiveInterpretationEnabledType")
	data.LiveStreamingMode = strVal("LiveStreamingMode")
	data.LobbyChat = strVal("LobbyChat")
	data.MediaBitRateKb = intVal("MediaBitRateKb")
	data.MeetingChatEnabledType = strVal("MeetingChatEnabledType")
	data.MeetingInviteLanguages = strVal("MeetingInviteLanguages")
	data.NewMeetingRecordingExpirationDays = intVal("NewMeetingRecordingExpirationDays")
	data.NoiseSuppressionForDialInParticipants = strVal("NoiseSuppressionForDialInParticipants")
	data.ParticipantNameChange = strVal("ParticipantNameChange")
	data.PreferredMeetingProviderForIslandsMode = strVal("PreferredMeetingProviderForIslandsMode")
	data.QnAEngagementMode = strVal("QnAEngagementMode")
	data.RecordingStorageMode = strVal("RecordingStorageMode")
	data.RoomAttributeUserOverride = strVal("RoomAttributeUserOverride")
	data.RoomPeopleNameUserOverride = strVal("RoomPeopleNameUserOverride")
	data.ScreenSharingMode = strVal("ScreenSharingMode")
	data.SmsNotifications = strVal("SmsNotifications")
	data.SpeakerAttributionMode = strVal("SpeakerAttributionMode")
	data.StreamingAttendeeMode = strVal("StreamingAttendeeMode")
	data.TeamsCameraFarEndPTZMode = strVal("TeamsCameraFarEndPTZMode")
	data.Tenant = strVal("Tenant")
	data.UsersCanAdmitFromLobby = strVal("UsersCanAdmitFromLobby")
	data.VideoFiltersMode = strVal("VideoFiltersMode")
	data.VoiceIsolation = strVal("VoiceIsolation")
	data.VoiceSimulationInInterpreter = strVal("VoiceSimulationInInterpreter")
	data.WatermarkForAnonymousUsers = strVal("WatermarkForAnonymousUsers")
	data.WatermarkForCameraVideoOpacity = intVal("WatermarkForCameraVideoOpacity")
	data.WatermarkForCameraVideoPattern = strVal("WatermarkForCameraVideoPattern")
	data.WatermarkForScreenSharingOpacity = intVal("WatermarkForScreenSharingOpacity")
	data.WatermarkForScreenSharingPattern = strVal("WatermarkForScreenSharingPattern")
	data.AllowedUsersForMeetingDetails = strVal("AllowedUsersForMeetingDetails")
	data.RealTimeText = strVal("RealTimeText")
	data.ParticipantSlideControl = strVal("ParticipantSlideControl")
	data.WhoCanRegister = strVal("WhoCanRegister")
}
