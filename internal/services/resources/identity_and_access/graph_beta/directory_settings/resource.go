package graphBetaDirectorySettings

import (
	"context"
	"regexp"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	planmodifiers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/plan_modifiers"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName  = "microsoft365_graph_beta_identity_and_access_directory_settings"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &DirectorySettingsResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &DirectorySettingsResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &DirectorySettingsResource{}
)

func NewDirectorySettingsResource() resource.Resource {
	return &DirectorySettingsResource{
		ReadPermissions: []string{
			"Directory.Read.All",
		},
		WritePermissions: []string{
			"GroupSettings.ReadWrite.All",
			"Directory.ReadWrite.All",
			"Policy.ReadWrite.Authorization",
		},
		ResourcePath: "/settings",
	}
}

type DirectorySettingsResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

// Metadata returns the resource type name.
func (r *DirectorySettingsResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

// Configure sets the client for the resource.
func (r *DirectorySettingsResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

// ImportState imports the resource state.
// For this singleton resource, the ID is always the settings object ID that has the Group.Unified template.
func (r *DirectorySettingsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Schema defines the schema for the resource.
func (r *DirectorySettingsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages directory settings using the `/beta/settings` endpoint. " +
			"Directory settings allow you to instantiate and configure tenant-wide or object-specific controls based on directorySettingTemplate objects.\n\n" +
			"## Endpoint Behavior\n\n" +
			"The `/beta/settings` endpoint manages **directorySetting** objects, which are instantiated from **directorySettingTemplate** blueprints. " +
			"Microsoft provides default configurations via templates (e.g., Group.Unified, Password Rule Settings, Consent Policy Settings). " +
			"To customize these defaults, you create a directorySetting object from a template and update its values.\n\n" +
			"Each template typically allows only **one** settings object per tenant. However, some templates work differently:\n" +
			"- **Tenant-level settings**: Most templates (Group.Unified, Password Rule Settings, etc.) apply configurations tenant-wide\n" +
			"- **Object-specific settings**: The Group.Unified.Guest template (`08d542b9-071f-4e16-94b0-74abb372e3d9`) creates settings that must be " +
			"assigned to individual groups. You create the directory setting object from the template, then assign it to specific groups on a per-object basis\n\n" +
			"**Important**: Tenant-level settings always take precedence over object-specific settings. If you configure a setting at the directory level " +
			"(e.g., AllowToAddGuests set to \"false\" in Group.Unified), any group-level configuration (Group.Unified.Guest) for the same setting will be ignored. " +
			"To selectively apply settings to specific groups, you must set the tenant-level setting to \"true\" (open by default), then create group-level " +
			"settings objects for individual groups where you want to restrict access.\n\n" +
			"The `overwrite_existing_settings` flag controls behavior during resource creation. When `true`, the resource checks for existing settings: " +
			"if found, it updates them (PATCH); if not found, it creates them (POST). When `false` (default), it always attempts creation (POST), " +
			"which will fail if settings already exist.\n\n" +
			"**Licensing**: Using Microsoft 365 group settings requires Microsoft Entra ID P1 or Microsoft Entra Basic EDU license " +
			"for each unique user who is a member of one or more Microsoft 365 groups.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
				MarkdownDescription: "The unique identifier of the directory setting object. " +
					"This is a system-generated UUID that identifies the settings object.",
			},
			"template_type": schema.StringAttribute{
				Required: true,
				MarkdownDescription: "The type of directory setting template to instantiate. Valid values:\n" +
					"- `Group.Unified` - Tenant-wide Microsoft 365 group settings\n" +
					"- `Group.Unified.Guest` - Group-specific guest settings\n" +
					"- `Application` - Tenant-wide application behavior settings\n" +
					"- `Password Rule Settings` - Tenant-wide password rule settings\n" +
					"- `Prohibited Names Settings` - Prohibited names for applications\n" +
					"- `Custom Policy Settings` - Custom conditional access policy settings\n" +
					"- `Prohibited Names Restricted Settings` - Allowed names for applications\n" +
					"- `Consent Policy Settings` - Tenant-wide consent policy settings\n\n" +
					"This field determines which template is used and cannot be changed after creation.",
				Validators: []validator.String{
					stringvalidator.OneOf(
						TemplateTypeGroupUnifiedGuest,
						TemplateTypeApplication,
						TemplateTypePasswordRuleSettings,
						TemplateTypeGroupUnified,
						TemplateTypeProhibitedNamesSettings,
						TemplateTypeCustomPolicySettings,
						TemplateTypeProhibitedNamesRestricted,
						TemplateTypeConsentPolicySettings,
					),
				},
				PlanModifiers: []planmodifier.String{
					planmodifiers.RequiresReplaceString(),
				},
			},
			"overwrite_existing_settings": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(false),
				MarkdownDescription: "When set to `true`, Terraform will overwrite existing directory settings " +
					"with the configuration specified in this resource. This is useful for first-time adoption when " +
					"directory settings already exist with default values. Setting this to `true` forces a PATCH operation " +
					"to replace the existing settings. Defaults to `false`, which attempts to create (POST) new settings first. " +
					"If settings already exist and this is `false`, creation will fail.",
			},
			"timeouts": commonschema.Timeouts(ctx),
		},
		Blocks: map[string]schema.Block{
			"group_unified_guest": schema.SingleNestedBlock{
				MarkdownDescription: "Settings for a specific Unified Group. Template ID: `08d542b9-071f-4e16-94b0-74abb372e3d9`",
				Attributes: map[string]schema.Attribute{
					"allow_to_add_guests": schema.BoolAttribute{
						Optional:            true,
						Computed:            true,
						Default:             booldefault.StaticBool(true),
						MarkdownDescription: "Flag indicating if guests are allowed in a specific Unified Group.",
					},
				},
			},
			"application": schema.SingleNestedBlock{
				MarkdownDescription: "Settings for managing tenant-wide application behavior. Template ID: `4bc7f740-180e-4586-adb6-38b2e9024e6b`",
				Attributes: map[string]schema.Attribute{
					"enable_access_check_for_privileged_application_updates": schema.BoolAttribute{
						Optional:            true,
						Computed:            true,
						Default:             booldefault.StaticBool(false),
						MarkdownDescription: "Flag indicating if access check for application privileged updates is turned on.",
					},
				},
			},
			"password_rule_settings": schema.SingleNestedBlock{
				MarkdownDescription: "Settings for managing tenant-wide password rule settings. Template ID: `5cf42378-d67d-4f36-ba46-e8b86229381d`",
				Attributes: map[string]schema.Attribute{
					"banned_password_check_on_premises_mode": schema.StringAttribute{
						Optional:            true,
						Computed:            true,
						Default:             stringdefault.StaticString("Audit"),
						MarkdownDescription: "How should we enforce password policy check in on-premises system.",
					},
					"enable_banned_password_check_on_premises": schema.BoolAttribute{
						Optional:            true,
						Computed:            true,
						Default:             booldefault.StaticBool(true),
						MarkdownDescription: "Flag indicating if the banned password check is turned on or not for on-premises system.",
					},
					"enable_banned_password_check": schema.BoolAttribute{
						Optional:            true,
						Computed:            true,
						Default:             booldefault.StaticBool(true),
						MarkdownDescription: "Flag indicating if the banned password check for tenant specific banned password list is turned on or not.",
					},
					"lockout_duration_in_seconds": schema.Int32Attribute{
						Optional:            true,
						Computed:            true,
						MarkdownDescription: "The duration in seconds of the initial lockout period.",
					},
					"lockout_threshold": schema.Int32Attribute{
						Optional:            true,
						Computed:            true,
						MarkdownDescription: "The number of failed login attempts before the first lockout period begins.",
					},
					"banned_password_list": schema.StringAttribute{
						Optional:            true,
						Computed:            true,
						Default:             stringdefault.StaticString(""),
						MarkdownDescription: "A tab-delimited banned password list.",
					},
				},
			},
			"group_unified": schema.SingleNestedBlock{
				MarkdownDescription: "Settings for Unified Groups. Template ID: `62375ab9-6b52-47ed-826b-58e47e0e304b`",
				Attributes: map[string]schema.Attribute{
					"new_unified_group_writeback_default": schema.BoolAttribute{
						Optional:            true,
						Computed:            true,
						Default:             booldefault.StaticBool(true),
						MarkdownDescription: "Default value of IsWritebackEnabled property for newly created Unified Groups.",
					},
					"enable_mip_labels": schema.BoolAttribute{
						Optional:            true,
						Computed:            true,
						Default:             booldefault.StaticBool(false),
						MarkdownDescription: "Flag indicating whether Microsoft Information Protection labels can be assigned to Unified Groups.",
					},
					"custom_blocked_words_list": schema.StringAttribute{
						Optional:            true,
						Computed:            true,
						Default:             stringdefault.StaticString(""),
						MarkdownDescription: "A comma-delimited list of blocked words for Unified Group displayName and mailNickName.",
					},
					"enable_ms_standard_blocked_words": schema.BoolAttribute{
						Optional:            true,
						Computed:            true,
						Default:             booldefault.StaticBool(false),
						MarkdownDescription: "A flag indicating whether or not to enable the Microsoft Standard list of blocked words for Unified Group displayName and mailNickName.",
					},
					"classification_descriptions": schema.StringAttribute{
						Optional:            true,
						Computed:            true,
						Default:             stringdefault.StaticString(""),
						MarkdownDescription: "A comma-delimited list of structured strings describing the classification values in the ClassificationList. The structure of the string is: Value: Description",
					},
					"default_classification": schema.StringAttribute{
						Optional:            true,
						Computed:            true,
						Default:             stringdefault.StaticString(""),
						MarkdownDescription: "The classification value to be used by default for Unified Group creation.",
					},
					"prefix_suffix_naming_requirement": schema.StringAttribute{
						Optional:            true,
						Computed:            true,
						Default:             stringdefault.StaticString(""),
						MarkdownDescription: "A structured string describing how a Unified Group displayName and mailNickname should be structured. Please refer to docs to discover how to structure a valid requirement.",
					},
					"allow_guests_to_be_group_owner": schema.BoolAttribute{
						Optional:            true,
						Computed:            true,
						Default:             booldefault.StaticBool(false),
						MarkdownDescription: "Flag indicating if guests are allowed to be owner in any Unified Group.",
					},
					"allow_guests_to_access_groups": schema.BoolAttribute{
						Optional:            true,
						Computed:            true,
						Default:             booldefault.StaticBool(true),
						MarkdownDescription: "Flag indicating if guests are allowed to access any Unified Group resources.",
					},
					"guest_usage_guidelines_url": schema.StringAttribute{
						Optional:            true,
						Computed:            true,
						Default:             stringdefault.StaticString(""),
						MarkdownDescription: "A link to the Group Usage Guidelines for guests.",
					},
					"group_creation_allowed_group_id": schema.StringAttribute{
						Optional:            true,
						Computed:            true,
						Default:             stringdefault.StaticString(""),
						MarkdownDescription: "Guid of the security group that is always allowed to create Unified Groups.",
						Validators: []validator.String{
							stringvalidator.RegexMatches(
								regexp.MustCompile(constants.GuidRegex),
								"group_creation_allowed_group_id must be a valid GUID in the format '00000000-0000-0000-0000-000000000000'",
							),
						},
					},
					"allow_to_add_guests": schema.BoolAttribute{
						Optional:            true,
						Computed:            true,
						Default:             booldefault.StaticBool(true),
						MarkdownDescription: "Flag indicating if guests are allowed in any Unified Group.",
					},
					"usage_guidelines_url": schema.StringAttribute{
						Optional:            true,
						Computed:            true,
						Default:             stringdefault.StaticString(""),
						MarkdownDescription: "A link to the Group Usage Guidelines.",
					},
					"classification_list": schema.StringAttribute{
						Optional:            true,
						Computed:            true,
						Default:             stringdefault.StaticString(""),
						MarkdownDescription: "A comma-delimited list of valid classification values that can be applied to Unified Groups.",
					},
					"enable_group_creation": schema.BoolAttribute{
						Optional:            true,
						Computed:            true,
						Default:             booldefault.StaticBool(true),
						MarkdownDescription: "Flag indicating if group creation feature is on.",
					},
				},
			},
			"prohibited_names_settings": schema.SingleNestedBlock{
				MarkdownDescription: "Settings for managing tenant-wide prohibited names settings. Template ID: `80661d51-be2f-4d46-9713-98a2fcaec5bc`",
				Attributes: map[string]schema.Attribute{
					"custom_blocked_sub_strings_list": schema.StringAttribute{
						Optional:            true,
						Computed:            true,
						Default:             stringdefault.StaticString(""),
						MarkdownDescription: "A comma delimited list of substring reserved words to block for application display names.",
					},
					"custom_blocked_whole_words_list": schema.StringAttribute{
						Optional:            true,
						Computed:            true,
						Default:             stringdefault.StaticString(""),
						MarkdownDescription: "A comma delimited list of reserved words to block for application display names.",
					},
				},
			},
			"custom_policy_settings": schema.SingleNestedBlock{
				MarkdownDescription: "Settings for managing tenant-wide custom policy settings. Template ID: `898f1161-d651-43d1-805c-3b0b388a9fc2`",
				Attributes: map[string]schema.Attribute{
					"custom_conditional_access_policy_url": schema.StringAttribute{
						Optional:            true,
						Computed:            true,
						Default:             stringdefault.StaticString(""),
						MarkdownDescription: "Custom conditional access policy url.",
					},
				},
			},
			"prohibited_names_restricted_settings": schema.SingleNestedBlock{
				MarkdownDescription: "Settings for managing tenant-wide prohibited names restricted settings. Template ID: `aad3907d-1d1a-448b-b3ef-7bf7f63db63b`",
				Attributes: map[string]schema.Attribute{
					"custom_allowed_sub_strings_list": schema.StringAttribute{
						Optional:            true,
						Computed:            true,
						Default:             stringdefault.StaticString(""),
						MarkdownDescription: "A comma delimited list of substring reserved words to allow for application display names.",
					},
					"custom_allowed_whole_words_list": schema.StringAttribute{
						Optional:            true,
						Computed:            true,
						Default:             stringdefault.StaticString(""),
						MarkdownDescription: "A comma delimited list of whole reserved words to allow for application display names.",
					},
					"do_not_validate_against_trademark": schema.BoolAttribute{
						Optional:            true,
						Computed:            true,
						Default:             booldefault.StaticBool(false),
						MarkdownDescription: "Flag indicating if prohibited names validation against trademark global list is disabled.",
					},
				},
			},
			"consent_policy_settings": schema.SingleNestedBlock{
				MarkdownDescription: "Settings for managing tenant-wide consent policy. Template ID: `dffd5d46-495d-40a9-8e21-954ff55e198a`",
				Attributes: map[string]schema.Attribute{
					"enable_group_specific_consent": schema.BoolAttribute{
						Computed: true,
						MarkdownDescription: "Flag indicating if groups owners are allowed to grant group specific permissions. " +
							"**Note**: This field is read-only and computed by Microsoft Graph API. It cannot be set directly " +
							"and may be automatically enabled when `constrain_group_specific_consent_to_members_of_group_id` is configured.",
					},
					"block_user_consent_for_risky_apps": schema.BoolAttribute{
						Optional:            true,
						Computed:            true,
						Default:             booldefault.StaticBool(true),
						MarkdownDescription: "Flag indicating if user consent will be blocked when a risky request is detected. Administrators will still be able to consent to apps considered risky.",
					},
					"enable_admin_consent_requests": schema.BoolAttribute{
						Optional:            true,
						Computed:            true,
						Default:             booldefault.StaticBool(false),
						MarkdownDescription: "Flag indicating if users will be able to request admin consent when they are unable to grant consent to an app themselves.",
					},
					"constrain_group_specific_consent_to_members_of_group_id": schema.StringAttribute{
						Computed: true,
						MarkdownDescription: "If EnableGroupSpecificConsent is set to \"True\" and this is set to a security group object ID, members (both direct and transitive) of the group identified will be authorized to grant group-specific permissions to the groups they own. " +
							"**Note**: This field is read-only and cannot be set directly via this API endpoint. It may require special licensing or tenant configuration.",
					},
				},
			},
		},
	}
}
