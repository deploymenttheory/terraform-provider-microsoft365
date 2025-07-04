---
page_title: "microsoft365_powershell_microsoft_teams_teams_meeting_policy Resource - terraform-provider-microsoft365"
subcategory: "Microsoft Teams"
description: |-
  Manages a Microsoft Teams Meeting Policy using PowerShell cmdlets. The CsTeamsMeetingPolicy cmdlets enable administrators to control the type of meetings that users can create or the features that they can access while in a meeting. It also helps determine how meetings deal with anonymous or external users.
---

# microsoft365_powershell_microsoft_teams_teams_meeting_policy (Resource)

Manages a Microsoft Teams Meeting Policy using PowerShell cmdlets. The CsTeamsMeetingPolicy cmdlets enable administrators to control the type of meetings that users can create or the features that they can access while in a meeting. It also helps determine how meetings deal with anonymous or external users.

## Microsoft Documentation

- [Get CsTeamsMeetingPolicy](https://learn.microsoft.com/en-us/powershell/module/teams/get-csteamsmeetingpolicy?view=teams-ps)
- [Set CsTeamsMeetingPolicy](https://learn.microsoft.com/en-us/powershell/module/teams/set-csteamsmeetingpolicy?view=teams-ps)
- [New CsTeamsMeetingPolicy](https://learn.microsoft.com/en-us/powershell/module/teams/new-csteamsmeetingpolicy?view=teams-ps)

## API Permissions

The following API permissions are required in order to use this resource.

### Microsoft Teams

- **Application**: `TeamsAdministration.ReadWrite.All`

## Version History

| Version | Status | Notes |
|---------|--------|-------|
| v0.21.0-alpha | Experimental | Initial release |

## Example Usage

```terraform
resource "microsoft365_powershell_microsoft_teams_teams_meeting_policy" "example" {
  id                                               = "ExampleMeetingPolicy"
  xds_identity                                     = "ExampleMeetingPolicy"
  ai_interpreter                                   = "Enabled"
  allow_annotations                                = true
  allow_anonymous_users_to_dial_out                = false
  allow_anonymous_users_to_join_meeting            = true
  allow_anonymous_users_to_start_meeting           = false
  allow_avatars_in_gallery                         = true
  allow_breakout_rooms                             = true
  allow_carbon_summary                             = false
  allow_cart_captions_scheduling                   = "Enabled"
  allow_channel_meeting_scheduling                 = true
  allow_cloud_recording                            = true
  allow_document_collaboration                     = "Enabled"
  allow_engagement_report                          = "Enabled"
  allow_external_non_trusted_meeting_chat          = false
  allow_external_participant_give_request_control  = true
  allow_immersive_view                             = true
  allow_ip_audio                                   = true
  allow_ip_video                                   = true
  allow_local_recording                            = false
  allow_meeting_coach                              = false
  allow_meet_now                                   = true
  allow_meeting_reactions                          = true
  allow_meeting_registration                       = false
  allow_ndi_streaming                              = false
  allow_network_configuration_settings_lookup      = true
  allow_organizers_to_override_lobby_settings      = false
  allow_outlook_add_in                             = true
  allow_pstn_users_to_bypass_lobby                 = false
  allow_participant_give_request_control           = true
  allow_power_point_sharing                        = true
  allow_private_meet_now                           = false
  allow_private_meeting_scheduling                 = true
  allow_recording_storage_outside_region           = false
  allow_screen_content_digitization                = false
  allow_shared_notes                               = true
  allow_tasks_from_transcript                      = "Enabled"
  allow_tracking_in_report                         = true
  allow_transcription                              = true
  allowed_users_for_meeting_context                = "Everyone"
  allow_user_to_join_external_meeting              = "Enabled"
  allow_watermark_customization_for_camera_video   = false
  allow_watermark_customization_for_screen_sharing = false
  allow_watermark_for_camera_video                 = false
  allow_watermark_for_screen_sharing               = false
  allow_whiteboard                                 = true
  allowed_streaming_media_input                    = "All"
  anonymous_user_authentication_method             = "Default"
  attendee_identity_masking                        = "None"
  audible_recording_notification                   = "Enabled"
  auto_admitted_users                              = "Everyone"
  auto_recording                                   = "Disabled"
  automatically_start_copilot                      = "Disabled"
  blocked_anonymous_join_client_types              = ["Skype", "Teams"]
  captcha_verification_for_meeting_join            = "Disabled"
  channel_recording_download                       = "Enabled"
  connect_to_meeting_controls                      = "Default"
  content_sharing_in_external_meetings             = "Enabled"
  copilot                                          = "Enabled"
  copy_restriction                                 = false
  description                                      = "Example Teams Meeting Policy for demonstration."
  designated_presenter_role_mode                   = "Default"
  detect_sensitive_content_during_screen_sharing   = false
  enroll_user_override                             = "None"
  explicit_recording_consent                       = "Required"
  external_meeting_join                            = "Enabled"
  info_shown_in_report_mode                        = "Default"
  ip_audio_mode                                    = "Default"
  ip_video_mode                                    = "Default"
  live_captions_enabled_type                       = "DisabledUserOverride"
  live_interpretation_enabled_type                 = "Disabled"
  live_streaming_mode                              = "Disabled"
  lobby_chat                                       = "Enabled"
  media_bit_rate_kb                                = 2048
  meeting_chat_enabled_type                        = "Enabled"
  meeting_invite_languages                         = "en-US"
  new_meeting_recording_expiration_days            = 30
  noise_suppression_for_dial_in_participants       = "Default"
  participant_name_change                          = "Allowed"
  preferred_meeting_provider_for_islands_mode      = "Teams"
  qna_engagement_mode                              = "Enabled"
  recording_storage_mode                           = "Default"
  room_attribute_user_override                     = "None"
  room_people_name_user_override                   = "None"
  screen_sharing_mode                              = "EntireScreen"
  sms_notifications                                = "Enabled"
  speaker_attribution_mode                         = "Default"
  streaming_attendee_mode                          = "Disabled"
  teams_camera_far_end_ptz_mode                    = "Disabled"
  tenant                                           = "exampletenant.onmicrosoft.com"
  users_can_admit_from_lobby                       = "Everyone"
  video_filters_mode                               = "Default"
  voice_isolation                                  = "Disabled"
  voice_simulation_in_interpreter                  = "Disabled"
  watermark_for_anonymous_users                    = "None"
  watermark_for_camera_video_opacity               = 0
  watermark_for_camera_video_pattern               = "None"
  watermark_for_screen_sharing_opacity             = 0
  watermark_for_screen_sharing_pattern             = "None"
  allowed_users_for_meeting_details                = "Everyone"
  real_time_text                                   = "Disabled"
  participant_slide_control                        = "Disabled"
  who_can_register                                 = "Everyone"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `id` (String) The name (Identity) of the Teams Meeting Policy.
- `xds_identity` (String) The XDS identity of the Teams Meeting Policy.

### Optional

- `ai_interpreter` (String) Enables the user to use the AI Interpreter related features.
- `allow_annotations` (Boolean) Allow users to use the Annotation feature.
- `allow_anonymous_users_to_dial_out` (Boolean) Allow anonymous users to dial out to a PSTN number.
- `allow_anonymous_users_to_join_meeting` (Boolean) Allow anonymous users to join meetings.
- `allow_anonymous_users_to_start_meeting` (Boolean) Allow anonymous users to start meetings.
- `allow_avatars_in_gallery` (Boolean) Allow avatars in gallery view.
- `allow_breakout_rooms` (Boolean) Allow breakout rooms in meetings.
- `allow_carbon_summary` (Boolean) Allow carbon summary in meetings.
- `allow_cart_captions_scheduling` (String) Allow scheduling of CART captions.
- `allow_channel_meeting_scheduling` (Boolean) Allow scheduling of channel meetings.
- `allow_cloud_recording` (Boolean) Allow cloud recording in meetings.
- `allow_document_collaboration` (String) Allow document collaboration in meetings.
- `allow_engagement_report` (String) Allow engagement report in meetings.
- `allow_external_non_trusted_meeting_chat` (Boolean) Allow chat for external non-trusted users in meetings.
- `allow_external_participant_give_request_control` (Boolean) Allow external participants to give/request control.
- `allow_immersive_view` (Boolean) Allow immersive view in meetings.
- `allow_ip_audio` (Boolean) Allow IP audio in meetings.
- `allow_ip_video` (Boolean) Allow IP video in meetings.
- `allow_local_recording` (Boolean) Allow local recording in meetings.
- `allow_meet_now` (Boolean) Allow Meet Now in meetings.
- `allow_meeting_coach` (Boolean) Allow meeting coach in meetings.
- `allow_meeting_reactions` (Boolean) Allow meeting reactions in meetings.
- `allow_meeting_registration` (Boolean) Allow meeting registration.
- `allow_ndi_streaming` (Boolean) Allow NDI streaming in meetings.
- `allow_network_configuration_settings_lookup` (Boolean) Allow network configuration settings lookup.
- `allow_organizers_to_override_lobby_settings` (Boolean) Allow organizers to override lobby settings.
- `allow_outlook_add_in` (Boolean) Allow Outlook add-in for meetings.
- `allow_participant_give_request_control` (Boolean) Allow participants to give/request control.
- `allow_power_point_sharing` (Boolean) Allow PowerPoint sharing in meetings.
- `allow_private_meet_now` (Boolean) Allow private Meet Now in meetings.
- `allow_private_meeting_scheduling` (Boolean) Allow private meeting scheduling.
- `allow_pstn_users_to_bypass_lobby` (Boolean) Allow PSTN users to bypass lobby.
- `allow_recording_storage_outside_region` (Boolean) Allow recording storage outside region.
- `allow_screen_content_digitization` (Boolean) Allow screen content digitization.
- `allow_shared_notes` (Boolean) Allow shared notes in meetings.
- `allow_tasks_from_transcript` (String) Allow tasks from transcript.
- `allow_tracking_in_report` (Boolean) Allow tracking in report.
- `allow_transcription` (Boolean) Allow transcription in meetings.
- `allow_user_to_join_external_meeting` (String) Allow user to join external meeting.
- `allow_watermark_customization_for_camera_video` (Boolean) Allow watermark customization for camera video.
- `allow_watermark_customization_for_screen_sharing` (Boolean) Allow watermark customization for screen sharing.
- `allow_watermark_for_camera_video` (Boolean) Allow watermark for camera video.
- `allow_watermark_for_screen_sharing` (Boolean) Allow watermark for screen sharing.
- `allow_whiteboard` (Boolean) Allow whiteboard in meetings.
- `allowed_streaming_media_input` (String) Allowed streaming media input.
- `allowed_users_for_meeting_context` (String) Allowed users for meeting context.
- `allowed_users_for_meeting_details` (String) Allowed users for meeting details.
- `anonymous_user_authentication_method` (String) Anonymous user authentication method.
- `attendee_identity_masking` (String) Attendee identity masking.
- `audible_recording_notification` (String) Audible recording notification.
- `auto_admitted_users` (String) Auto admitted users.
- `auto_recording` (String) Auto recording setting.
- `automatically_start_copilot` (String) Automatically start Copilot.
- `blocked_anonymous_join_client_types` (List of String) Blocked anonymous join client types.
- `captcha_verification_for_meeting_join` (String) Captcha verification for meeting join.
- `channel_recording_download` (String) Channel recording download setting.
- `connect_to_meeting_controls` (String) Connect to meeting controls.
- `content_sharing_in_external_meetings` (String) Content sharing in external meetings.
- `copilot` (String) Copilot setting.
- `copy_restriction` (Boolean) Copy restriction setting.
- `description` (String) Description of the Teams Meeting Policy.
- `designated_presenter_role_mode` (String) Designated presenter role mode.
- `detect_sensitive_content_during_screen_sharing` (Boolean) Detect sensitive content during screen sharing.
- `enroll_user_override` (String) Enroll user override setting.
- `explicit_recording_consent` (String) Explicit recording consent setting.
- `external_meeting_join` (String) External meeting join setting.
- `info_shown_in_report_mode` (String) Info shown in report mode.
- `ip_audio_mode` (String) IP audio mode.
- `ip_video_mode` (String) IP video mode.
- `live_captions_enabled_type` (String) Live captions enabled type.
- `live_interpretation_enabled_type` (String) Live interpretation enabled type.
- `live_streaming_mode` (String) Live streaming mode.
- `lobby_chat` (String) Lobby chat setting.
- `media_bit_rate_kb` (Number) Maximum media bit rate in Kb.
- `meeting_chat_enabled_type` (String) Meeting chat enabled type.
- `meeting_invite_languages` (String) Meeting invite languages.
- `new_meeting_recording_expiration_days` (Number) Number of days before new meeting recordings expire.
- `noise_suppression_for_dial_in_participants` (String) Noise suppression for dial-in participants.
- `participant_name_change` (String) Participant name change setting.
- `participant_slide_control` (String) Participant slide control setting.
- `preferred_meeting_provider_for_islands_mode` (String) Preferred meeting provider for Islands mode.
- `qna_engagement_mode` (String) QnA engagement mode.
- `real_time_text` (String) Real time text setting.
- `recording_storage_mode` (String) Recording storage mode.
- `room_attribute_user_override` (String) Room attribute user override.
- `room_people_name_user_override` (String) Room people name user override.
- `screen_sharing_mode` (String) Screen sharing mode.
- `sms_notifications` (String) SMS notifications setting.
- `speaker_attribution_mode` (String) Speaker attribution mode.
- `streaming_attendee_mode` (String) Streaming attendee mode.
- `teams_camera_far_end_ptz_mode` (String) Teams camera far end PTZ mode.
- `tenant` (String) Tenant GUID.
- `users_can_admit_from_lobby` (String) Users who can admit from lobby.
- `video_filters_mode` (String) Video filters mode.
- `voice_isolation` (String) Voice isolation setting.
- `voice_simulation_in_interpreter` (String) Voice simulation in interpreter setting.
- `watermark_for_anonymous_users` (String) Watermark for anonymous users.
- `watermark_for_camera_video_opacity` (Number) Watermark opacity for camera video.
- `watermark_for_camera_video_pattern` (String) Watermark pattern for camera video.
- `watermark_for_screen_sharing_opacity` (Number) Watermark opacity for screen sharing.
- `watermark_for_screen_sharing_pattern` (String) Watermark pattern for screen sharing.
- `who_can_register` (String) Who can register for meetings.

## Important Notes

- **Teams Calling Policy**: This resource manages the calling policy for Microsoft Teams.
- **AI Interpreter**: Enables the user to use the AI Interpreter related features.
- **Allow Call Forwarding to Phone**: Allows users to forward calls to a phone number.
- **Allow Call Forwarding to User**: Allows users to forward calls to another user.
- **Allow Call Groups**: Allows users to create call groups.
- **Allow Call Redirect**: Allows users to redirect calls to another user.
- **Allow Cloud Recording for Calls**: Allows users to record calls.
- **Allow Delegation**: Allows users to delegate calls to another user.

## Import

Import is supported using the following syntax:

```shell
#! /bin/sh

terraform import microsoft365_powershell_microsoft_teams_teams_meeting_policy.example ExampleMeetingPolicy
```

