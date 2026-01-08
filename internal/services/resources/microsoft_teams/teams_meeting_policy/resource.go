package teamsMeetingPolicy

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"os"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/powershell"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type TeamsMeetingPolicyResource struct{}

func NewTeamsMeetingPolicyResource() resource.Resource {
	return &TeamsMeetingPolicyResource{}
}

func (r *TeamsMeetingPolicyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_powershell_microsoft_teams_teams_meeting_policy"
}

func (r *TeamsMeetingPolicyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// 1. Ensure PowerShell and Teams module are available
	if err := powershell.EnsurePowerShellAndTeamsModule(); err != nil {
		resp.Diagnostics.AddError("PowerShell/Teams Module Check Failed", err.Error())
		return
	}

	// Get authentication details from environment variables
	tenantId := os.Getenv("M365_TENANT_ID")
	authMethod := os.Getenv("M365_AUTH_METHOD")
	applicationId := os.Getenv("M365_CLIENT_ID")
	clientSecret := os.Getenv("M365_CLIENT_SECRET")
	certificatePath := os.Getenv("M365_CLIENT_CERTIFICATE_FILE_PATH")
	certificatePassword := os.Getenv("M365_CLIENT_CERTIFICATE_PASSWORD")

	tflog.Debug(ctx, "Configuring Teams Meeting Policy resource", map[string]any{
		"tenantId":      tenantId,
		"authMethod":    authMethod,
		"applicationId": applicationId,
	})

	// For validation phase, don't try to connect if we don't have credentials
	if req.ProviderData == nil {
		tflog.Debug(ctx, "Provider data is nil, skipping Teams Meeting Policy authentication")
		return
	}

	// If auth method is empty, try to determine it based on available credentials
	if authMethod == "" {
		if clientSecret != "" {
			authMethod = "client_secret"
			tflog.Debug(ctx, "No auth method specified, using client_secret based on available credentials")
		} else if certificatePath != "" {
			authMethod = "client_certificate"
			tflog.Debug(ctx, "No auth method specified, using client_certificate based on available credentials")
		} else {
			resp.Diagnostics.AddError(
				"Teams PowerShell Authentication Failed",
				"No authentication method specified and couldn't determine one from available credentials. Please set M365_AUTH_METHOD to 'client_secret' or 'client_certificate' and provide the necessary credentials.",
			)
			return
		}
	}

	switch authMethod {
	case "client_secret":
		if applicationId == "" || clientSecret == "" || tenantId == "" {
			resp.Diagnostics.AddError("Teams PowerShell Authentication Failed", "App ID, Client Secret, and Tenant ID must be provided for client_secret authentication.")
			return
		}
		if err := powershell.ConnectMicrosoftTeams(tenantId, applicationId, clientSecret, "", ""); err != nil {
			resp.Diagnostics.AddError("Teams PowerShell Authentication Failed", err.Error())
			return
		}
	case "client_certificate":
		if applicationId == "" || certificatePath == "" || tenantId == "" {
			resp.Diagnostics.AddError("Teams PowerShell Authentication Failed", "App ID, Certificate Path, and Tenant ID must be provided for client_certificate authentication.")
			return
		}
		// Try to extract thumbprint from the certificate file
		certificateThumbprint := ""
		certData, err := os.ReadFile(certificatePath)
		if err == nil {
			certs, _, err := helpers.ParseCertificateData(ctx, certData, []byte(certificatePassword))
			if err == nil && len(certs) > 0 {
				h := sha1.New()
				h.Write(certs[0].Raw)
				certificateThumbprint = hex.EncodeToString(h.Sum(nil))
			}
		}
		// Prefer thumbprint-based authentication if available
		if certificateThumbprint != "" {
			if err := powershell.ConnectMicrosoftTeams(tenantId, applicationId, "", certificateThumbprint, ""); err != nil {
				resp.Diagnostics.AddError("Teams PowerShell Authentication Failed", err.Error())
				return
			}
		} else {
			if err := powershell.ConnectMicrosoftTeams(tenantId, applicationId, "", "", certificatePath); err != nil {
				resp.Diagnostics.AddError("Teams PowerShell Authentication Failed", err.Error())
				return
			}
		}
	default:
		resp.Diagnostics.AddError(
			"Teams PowerShell Authentication Failed",
			"Only auth_method values 'client_secret' and 'client_certificate' are supported for Teams PowerShell resources.",
		)
		return
	}
}

func (r *TeamsMeetingPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *TeamsMeetingPolicyResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a Microsoft Teams Meeting Policy using PowerShell cmdlets. The CsTeamsMeetingPolicy cmdlets " +
			"enable administrators to control the type of meetings that users can create or the features that they can access while in " +
			"a meeting. It also helps determine how meetings deal with anonymous or external users.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name (Identity) of the Teams Meeting Policy.",
			},
			"xds_identity": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The XDS identity of the Teams Meeting Policy.",
			},
			"ai_interpreter": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Enables the user to use the AI Interpreter related features.",
			},
			"allow_annotations": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Allow users to use the Annotation feature.",
			},
			"allow_anonymous_users_to_dial_out": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Allow anonymous users to dial out to a PSTN number.",
			},
			"allow_anonymous_users_to_join_meeting": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Allow anonymous users to join meetings.",
			},
			"allow_anonymous_users_to_start_meeting": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Allow anonymous users to start meetings.",
			},
			"allow_avatars_in_gallery": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Allow avatars in gallery view.",
			},
			"allow_breakout_rooms": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Allow breakout rooms in meetings.",
			},
			"allow_carbon_summary": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Allow carbon summary in meetings.",
			},
			"allow_cart_captions_scheduling": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Allow scheduling of CART captions.",
			},
			"allow_channel_meeting_scheduling": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Allow scheduling of channel meetings.",
			},
			"allow_cloud_recording": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Allow cloud recording in meetings.",
			},
			"allow_document_collaboration": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Allow document collaboration in meetings.",
			},
			"allow_engagement_report": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Allow engagement report in meetings.",
			},
			"allow_external_non_trusted_meeting_chat": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Allow chat for external non-trusted users in meetings.",
			},
			"allow_external_participant_give_request_control": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Allow external participants to give/request control.",
			},
			"allow_immersive_view": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Allow immersive view in meetings.",
			},
			"allow_ip_audio": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Allow IP audio in meetings.",
			},
			"allow_ip_video": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Allow IP video in meetings.",
			},
			"allow_local_recording": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Allow local recording in meetings.",
			},
			"allow_meeting_coach": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Allow meeting coach in meetings.",
			},
			"allow_meet_now": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Allow Meet Now in meetings.",
			},
			"allow_meeting_reactions": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Allow meeting reactions in meetings.",
			},
			"allow_meeting_registration": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Allow meeting registration.",
			},
			"allow_ndi_streaming": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Allow NDI streaming in meetings.",
			},
			"allow_network_configuration_settings_lookup": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Allow network configuration settings lookup.",
			},
			"allow_organizers_to_override_lobby_settings": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Allow organizers to override lobby settings.",
			},
			"allow_outlook_add_in": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Allow Outlook add-in for meetings.",
			},
			"allow_pstn_users_to_bypass_lobby": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Allow PSTN users to bypass lobby.",
			},
			"allow_participant_give_request_control": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Allow participants to give/request control.",
			},
			"allow_power_point_sharing": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Allow PowerPoint sharing in meetings.",
			},
			"allow_private_meet_now": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Allow private Meet Now in meetings.",
			},
			"allow_private_meeting_scheduling": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Allow private meeting scheduling.",
			},
			"allow_recording_storage_outside_region": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Allow recording storage outside region.",
			},
			"allow_screen_content_digitization": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Allow screen content digitization.",
			},
			"allow_shared_notes": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Allow shared notes in meetings.",
			},
			"allow_tasks_from_transcript": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Allow tasks from transcript.",
			},
			"allow_tracking_in_report": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Allow tracking in report.",
			},
			"allow_transcription": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Allow transcription in meetings.",
			},
			"allowed_users_for_meeting_context": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Allowed users for meeting context.",
			},
			"allow_user_to_join_external_meeting": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Allow user to join external meeting.",
			},
			"allow_watermark_customization_for_camera_video": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Allow watermark customization for camera video.",
			},
			"allow_watermark_customization_for_screen_sharing": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Allow watermark customization for screen sharing.",
			},
			"allow_watermark_for_camera_video": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Allow watermark for camera video.",
			},
			"allow_watermark_for_screen_sharing": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Allow watermark for screen sharing.",
			},
			"allow_whiteboard": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Allow whiteboard in meetings.",
			},
			"allowed_streaming_media_input": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Allowed streaming media input.",
			},
			"anonymous_user_authentication_method": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Anonymous user authentication method.",
			},
			"attendee_identity_masking": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Attendee identity masking.",
			},
			"audible_recording_notification": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Audible recording notification.",
			},
			"auto_admitted_users": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Auto admitted users.",
			},
			"auto_recording": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Auto recording setting.",
			},
			"automatically_start_copilot": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Automatically start Copilot.",
			},
			"blocked_anonymous_join_client_types": schema.ListAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				MarkdownDescription: "Blocked anonymous join client types.",
			},
			"captcha_verification_for_meeting_join": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Captcha verification for meeting join.",
			},
			"channel_recording_download": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Channel recording download setting.",
			},
			"connect_to_meeting_controls": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Connect to meeting controls.",
			},
			"content_sharing_in_external_meetings": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Content sharing in external meetings.",
			},
			"copilot": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Copilot setting.",
			},
			"copy_restriction": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Copy restriction setting.",
			},
			"description": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Optional description of the resource. Maximum length is 1500 characters.",
				Validators: []validator.String{
					stringvalidator.LengthAtMost(1500),
				},
			},
			"designated_presenter_role_mode": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Designated presenter role mode.",
			},
			"detect_sensitive_content_during_screen_sharing": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Detect sensitive content during screen sharing.",
			},
			"enroll_user_override": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Enroll user override setting.",
			},
			"explicit_recording_consent": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Explicit recording consent setting.",
			},
			"external_meeting_join": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "External meeting join setting.",
			},
			"info_shown_in_report_mode": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Info shown in report mode.",
			},
			"ip_audio_mode": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "IP audio mode.",
			},
			"ip_video_mode": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "IP video mode.",
			},
			"live_captions_enabled_type": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Live captions enabled type.",
			},
			"live_interpretation_enabled_type": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Live interpretation enabled type.",
			},
			"live_streaming_mode": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Live streaming mode.",
			},
			"lobby_chat": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Lobby chat setting.",
			},
			"media_bit_rate_kb": schema.Int64Attribute{
				Optional:            true,
				MarkdownDescription: "Maximum media bit rate in Kb.",
			},
			"meeting_chat_enabled_type": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Meeting chat enabled type.",
			},
			"meeting_invite_languages": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Meeting invite languages.",
			},
			"new_meeting_recording_expiration_days": schema.Int64Attribute{
				Optional:            true,
				MarkdownDescription: "Number of days before new meeting recordings expire.",
			},
			"noise_suppression_for_dial_in_participants": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Noise suppression for dial-in participants.",
			},
			"participant_name_change": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Participant name change setting.",
			},
			"preferred_meeting_provider_for_islands_mode": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Preferred meeting provider for Islands mode.",
			},
			"qna_engagement_mode": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "QnA engagement mode.",
			},
			"recording_storage_mode": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Recording storage mode.",
			},
			"room_attribute_user_override": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Room attribute user override.",
			},
			"room_people_name_user_override": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Room people name user override.",
			},
			"screen_sharing_mode": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Screen sharing mode.",
			},
			"sms_notifications": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "SMS notifications setting.",
			},
			"speaker_attribution_mode": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Speaker attribution mode.",
			},
			"streaming_attendee_mode": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Streaming attendee mode.",
			},
			"teams_camera_far_end_ptz_mode": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Teams camera far end PTZ mode.",
			},
			"tenant": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Tenant GUID.",
			},
			"users_can_admit_from_lobby": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Users who can admit from lobby.",
			},
			"video_filters_mode": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Video filters mode.",
			},
			"voice_isolation": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Voice isolation setting.",
			},
			"voice_simulation_in_interpreter": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Voice simulation in interpreter setting.",
			},
			"watermark_for_anonymous_users": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Watermark for anonymous users.",
			},
			"watermark_for_camera_video_opacity": schema.Int64Attribute{
				Optional:            true,
				MarkdownDescription: "Watermark opacity for camera video.",
			},
			"watermark_for_camera_video_pattern": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Watermark pattern for camera video.",
			},
			"watermark_for_screen_sharing_opacity": schema.Int64Attribute{
				Optional:            true,
				MarkdownDescription: "Watermark opacity for screen sharing.",
			},
			"watermark_for_screen_sharing_pattern": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Watermark pattern for screen sharing.",
			},
			"allowed_users_for_meeting_details": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Allowed users for meeting details.",
			},
			"real_time_text": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Real time text setting.",
			},
			"participant_slide_control": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Participant slide control setting.",
			},
			"who_can_register": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Who can register for meetings.",
			},
		},
	}
}
