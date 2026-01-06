package teams_calling_policy

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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// TeamsCallingPolicyResource implements the Terraform resource for Teams Calling Policy

type TeamsCallingPolicyResource struct{}

func NewTeamsCallingPolicyResource() resource.Resource {
	return &TeamsCallingPolicyResource{}
}

func (r *TeamsCallingPolicyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_powershell_microsoft_teams_teams_calling_policy"
}

func (r *TeamsCallingPolicyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	tflog.Debug(ctx, "Configuring Teams Calling Policy resource", map[string]any{
		"tenantId":      tenantId,
		"authMethod":    authMethod,
		"applicationId": applicationId,
	})

	// For validation phase, don't try to connect if we don't have credentials
	if req.ProviderData == nil {
		tflog.Debug(ctx, "Provider data is nil, skipping Teams Calling Policy authentication")
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

func (r *TeamsCallingPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *TeamsCallingPolicyResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a Microsoft Teams Calling Policy using PowerShell cmdlets. See [Set-CsTeamsCallingPolicy](https://learn.microsoft.com/en-us/powershell/module/teams/set-csteamscallingpolicy?view=teams-ps) and [New-CsTeamsCallingPolicy](https://learn.microsoft.com/en-us/powershell/module/teams/new-csteamscallingpolicy?view=teams-ps).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name (Identity) of the Teams Calling Policy.",
			},
			"ai_interpreter": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Enables the user to use the AI Interpreter related features.",
			},
			"allow_call_forwarding_to_phone": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Allow call forwarding to phone.",
			},
			"allow_call_forwarding_to_user": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Allow call forwarding to user.",
			},
			"allow_call_groups": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Allow call groups.",
			},
			"allow_call_redirect": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Allow call redirect. Valid values: Enabled, Disabled.",
			},
			"allow_cloud_recording_for_calls": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Allow cloud recording for calls.",
			},
			"allow_delegation": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Allow delegation.",
			},
			"allow_private_calling": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Allow private calling.",
			},
			"allow_sip_devices_calling": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Allow SIP devices calling.",
			},
			"allow_transcription_for_calling": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Allow transcription for calling.",
			},
			"allow_voicemail": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Allow voicemail. Valid values: AlwaysEnabled, AlwaysDisabled, UserOverride.",
			},
			"allow_web_pstn_calling": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Allow web PSTN calling.",
			},
			"auto_answer_enabled_type": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Auto answer enabled type.",
			},
			"busy_on_busy_enabled_type": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Busy on busy enabled type.",
			},
			"calling_spend_user_limit": schema.Int64Attribute{
				Optional:            true,
				MarkdownDescription: "Calling spend user limit.",
			},
			"call_recording_expiration_days": schema.Int64Attribute{
				Optional:            true,
				MarkdownDescription: "Call recording expiration days.",
			},
			"copilot": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Copilot setting.",
			},
			"description": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
				MarkdownDescription: "Optional description of the resource. Maximum length is 1500 characters.",
				Validators: []validator.String{
					stringvalidator.LengthAtMost(1500),
				},
			},
			"enable_spend_limits": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Enable spend limits.",
			},
			"enable_web_pstn_media_bypass": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Enable web PSTN media bypass.",
			},
			"inbound_federated_call_routing_treatment": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Inbound federated call routing treatment.",
			},
			"inbound_pstn_call_routing_treatment": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Inbound PSTN call routing treatment.",
			},
			"live_captions_enabled_type_for_calling": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Live captions enabled type for calling.",
			},
			"music_on_hold_enabled_type": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Music on hold enabled type.",
			},
			"popout_app_path_for_incoming_pstn_calls": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Popout app path for incoming PSTN calls.",
			},
			"popout_for_incoming_pstn_calls": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Popout for incoming PSTN calls.",
			},
			"prevent_toll_bypass": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Prevent toll bypass.",
			},
			"spam_filtering_enabled_type": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Spam filtering enabled type.",
			},
			"voice_simulation_in_interpreter": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Voice simulation in interpreter.",
			},
			"real_time_text": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Real time text.",
			},
		},
	}
}
