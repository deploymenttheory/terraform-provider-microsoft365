package graphBetaM365AdminSharepointSettings

import (
	"context"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/identityschema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
)

const (
	ResourceName  = "microsoft365_graph_beta_m365_admin_sharepoint_settings"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180

	// singletonID is the fixed identifier for the tenant-wide SharePoint settings.
	// Graph does not return an id for this singleton.
	singletonID = "sharepointSettings"
)

// lowerCaseGUIDRegex matches a GUID in the lower-case form Graph returns, so configured values match state.
var lowerCaseGUIDRegex = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &SharepointSettingsResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &SharepointSettingsResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &SharepointSettingsResource{}

	// Enables plan modification/diff suppression
	_ resource.ResourceWithModifyPlan = &SharepointSettingsResource{}

	// Enables identity schema for list resource support
	_ resource.ResourceWithIdentity = &SharepointSettingsResource{}
)

// https://learn.microsoft.com/en-us/graph/api/sharepointsettings-get?view=graph-rest-beta#permissions
// https://learn.microsoft.com/en-us/graph/api/sharepointsettings-update?view=graph-rest-beta#permissions
func NewSharepointSettingsResource() resource.Resource {
	return &SharepointSettingsResource{
		ReadPermissions: []string{
			"SharePointTenantSettings.Read.All",
		},
		WritePermissions: []string{
			"SharePointTenantSettings.ReadWrite.All",
		},
	}
}

type SharepointSettingsResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
}

// Metadata returns the resource type name.
func (r *SharepointSettingsResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

// Configure sets the client for the resource.
func (r *SharepointSettingsResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

// ImportState imports the resource state using the fixed singleton identifier.
func (r *SharepointSettingsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if req.ID != singletonID {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			"Import this singleton using the literal ID "+singletonID+".",
		)
		return
	}
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// IdentitySchema defines the identity schema for this resource, used by list operations to uniquely identify instances
func (r *SharepointSettingsResource) IdentitySchema(ctx context.Context, req resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	resp.IdentitySchema = identityschema.Schema{
		Attributes: map[string]identityschema.Attribute{
			"id": identityschema.StringAttribute{
				RequiredForImport: true,
			},
		},
	}
}

// Schema defines the schema for the resource.
func (r *SharepointSettingsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages the tenant-wide SharePoint and OneDrive settings using the Microsoft Graph beta `/admin/sharepoint/settings` endpoint. " +
			"The settings control external sharing, site creation, OneDrive storage and retention, the OneDrive sync app and idle session sign-out.\n\n" +
			"This is a **singleton resource**: one settings object exists per tenant. `create` adopts the existing settings and `create`/`update` " +
			"apply the configuration with PATCH. Settings that are not configured keep their current tenant values and are not sent to Graph. " +
			"On `destroy`, Terraform only removes the resource from state and leaves the tenant settings as they are. Manage only one instance per tenant across all Terraform states.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The fixed singleton identifier `" + singletonID + "`.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"allowed_domain_guids_for_sync_app": schema.SetAttribute{
				MarkdownDescription: "The trusted domain GUIDs for the OneDrive sync app, in lower case. Used with `is_unmanaged_sync_app_for_tenant_restricted`. " +
					"Set to an empty set to clear the list. When not configured, the current tenant value is kept.",
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.Set{
					setvalidator.ValueStringsAre(stringvalidator.RegexMatches(lowerCaseGUIDRegex, "must be a lower-case GUID, as returned by Graph")),
				},
			},
			"available_managed_paths_for_site_creation": schema.SetAttribute{
				MarkdownDescription: "The managed paths available for site creation, for example `/sites/` and `/teams/`. Read-only.",
				ElementType:         types.StringType,
				Computed:            true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"deleted_user_personal_site_retention_period_in_days": schema.Int32Attribute{
				MarkdownDescription: "The number of days a deleted user's OneDrive is kept. Must be between `30` and `3650`. " +
					"When not configured, the current tenant value is kept.",
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Int32{
					int32planmodifier.UseStateForUnknown(),
				},
				Validators: []validator.Int32{
					int32validator.Between(30, 3650),
				},
			},
			"excluded_file_extensions_for_sync_app": schema.SetAttribute{
				MarkdownDescription: "The file extensions that the OneDrive sync app does not upload, for example `*.pst`. " +
					"Set to an empty set to clear the list. When not configured, the current tenant value is kept.",
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"idle_session_sign_out": schema.SingleNestedAttribute{
				MarkdownDescription: "The idle session sign-out policy for SharePoint and OneDrive. When not configured, the current tenant value is kept.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{
					"is_enabled": schema.BoolAttribute{
						MarkdownDescription: "Whether the idle session sign-out policy is enabled. When not configured, the current tenant value is kept.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.UseStateForUnknown(),
						},
					},
					"sign_out_after_in_seconds": schema.Int64Attribute{
						MarkdownDescription: "The number of seconds of inactivity after which a user is signed out. When not configured, the current tenant value is kept.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.Int64{
							int64planmodifier.UseStateForUnknown(),
						},
						Validators: []validator.Int64{
							int64validator.AtLeast(0),
						},
					},
					"warn_after_in_seconds": schema.Int64Attribute{
						MarkdownDescription: "The number of seconds of inactivity after which a user is warned that they will be signed out. " +
							"Must be less than `sign_out_after_in_seconds`. When not configured, the current tenant value is kept.",
						Optional: true,
						Computed: true,
						PlanModifiers: []planmodifier.Int64{
							int64planmodifier.UseStateForUnknown(),
						},
						Validators: []validator.Int64{
							int64validator.AtLeast(0),
						},
					},
				},
			},
			"image_tagging_option": schema.StringAttribute{
				MarkdownDescription: "The image tagging option for the tenant. Possible values are `disabled`, `basic` and `enhanced`. " +
					"When not configured, the current tenant value is kept.",
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf("disabled", "basic", "enhanced"),
				},
			},
			"is_commenting_on_site_pages_enabled": schema.BoolAttribute{
				MarkdownDescription: "Whether comments are allowed on modern site pages in SharePoint. When not configured, the current tenant value is kept.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"is_file_activity_notification_enabled": schema.BoolAttribute{
				MarkdownDescription: "Whether push notifications are enabled for OneDrive events. When not configured, the current tenant value is kept.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"is_legacy_auth_protocols_enabled": schema.BoolAttribute{
				MarkdownDescription: "Whether legacy authentication protocols are enabled for the tenant. Enabling legacy authentication weakens conditional access enforcement for SharePoint and OneDrive. When not configured, the current tenant value is kept.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"is_loop_enabled": schema.BoolAttribute{
				MarkdownDescription: "Whether the Fluid Framework (Loop components) is allowed on SharePoint sites. When not configured, the current tenant value is kept.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"is_mac_sync_app_enabled": schema.BoolAttribute{
				MarkdownDescription: "Whether files can be synced using the OneDrive sync app for Mac. When not configured, the current tenant value is kept.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"is_require_accepting_user_to_match_invited_user_enabled": schema.BoolAttribute{
				MarkdownDescription: "Whether guests must sign in using the same account to which sharing invitations are sent. Graph may return no value for this setting, in which case it is null in state. When not configured, the current tenant value is kept.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"is_resharing_by_external_users_enabled": schema.BoolAttribute{
				MarkdownDescription: "Whether guests are allowed to reshare files, folders and sites they don't own. When not configured, the current tenant value is kept.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"is_sharepoint_mobile_notification_enabled": schema.BoolAttribute{
				MarkdownDescription: "Whether mobile push notifications are enabled for SharePoint. When not configured, the current tenant value is kept.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"is_sharepoint_newsfeed_enabled": schema.BoolAttribute{
				MarkdownDescription: "Whether the newsfeed is allowed on modern site pages in SharePoint. When not configured, the current tenant value is kept.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"is_site_creation_enabled": schema.BoolAttribute{
				MarkdownDescription: "Whether users are allowed to create sites. When not configured, the current tenant value is kept.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"is_site_creation_ui_enabled": schema.BoolAttribute{
				MarkdownDescription: "Whether the UI commands for creating sites are shown. When not configured, the current tenant value is kept.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"is_site_pages_creation_enabled": schema.BoolAttribute{
				MarkdownDescription: "Whether creating new modern pages is allowed on SharePoint sites. When not configured, the current tenant value is kept.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"is_sites_storage_limit_automatic": schema.BoolAttribute{
				MarkdownDescription: "Whether site storage space is managed automatically (`true`) or with specific storage limits per site (`false`). When not configured, the current tenant value is kept.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"is_sync_button_hidden_on_personal_site": schema.BoolAttribute{
				MarkdownDescription: "Whether the sync button in OneDrive is hidden. When not configured, the current tenant value is kept.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"is_unmanaged_sync_app_for_tenant_restricted": schema.BoolAttribute{
				MarkdownDescription: "Whether users are allowed to sync files only on PCs joined to specific domains (see `allowed_domain_guids_for_sync_app`). When not configured, the current tenant value is kept.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"personal_site_default_storage_limit_in_mb": schema.Int64Attribute{
				MarkdownDescription: "The default OneDrive storage limit in megabytes for all new and existing users who are assigned a qualifying licence. " +
					"When not configured, the current tenant value is kept.",
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
				Validators: []validator.Int64{
					int64validator.AtLeast(1),
				},
			},
			"sharing_allowed_domain_list": schema.SetAttribute{
				MarkdownDescription: "The email domains that are allowed for sharing outside the organization. Applies when `sharing_domain_restriction_mode` is `allowList`. " +
					"Set to an empty set to clear the list. When not configured, the current tenant value is kept.",
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"sharing_blocked_domain_list": schema.SetAttribute{
				MarkdownDescription: "The email domains that are blocked for sharing outside the organization. Applies when `sharing_domain_restriction_mode` is `blockList`. " +
					"Set to an empty set to clear the list. When not configured, the current tenant value is kept.",
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"sharing_capability": schema.StringAttribute{
				MarkdownDescription: "The external sharing capability for the tenant. Possible values are `disabled` (only people in the organization), " +
					"`existingExternalUserSharingOnly` (existing guests), `externalUserSharingOnly` (new and existing guests) and " +
					"`externalUserAndGuestSharing` (anyone, including links that don't require sign-in). When not configured, the current tenant value is kept.",
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf("disabled", "externalUserSharingOnly", "externalUserAndGuestSharing", "existingExternalUserSharingOnly"),
				},
			},
			"sharing_domain_restriction_mode": schema.StringAttribute{
				MarkdownDescription: "The external sharing mode for domains. Possible values are `none`, `allowList` and `blockList`. " +
					"When not configured, the current tenant value is kept.",
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf("none", "allowList", "blockList"),
				},
			},
			"site_creation_default_managed_path": schema.StringAttribute{
				MarkdownDescription: "The managed path under which new team sites are created, for example `/sites/`. It must be one of " +
					"`available_managed_paths_for_site_creation`. When not configured, the current tenant value is kept.",
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"site_creation_default_storage_limit_in_mb": schema.Int32Attribute{
				MarkdownDescription: "The default storage quota in megabytes for a new site. When not configured, the current tenant value is kept.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Int32{
					int32planmodifier.UseStateForUnknown(),
				},
				Validators: []validator.Int32{
					int32validator.AtLeast(1),
				},
			},
			"tenant_default_timezone": schema.StringAttribute{
				MarkdownDescription: "The default time zone for newly created sites, for example `(UTC-08:00) Pacific Time (US and Canada)`. " +
					"See [SPRegionalSettings.TimeZones](https://learn.microsoft.com/en-us/sharepoint/dev/schema/regional-settings-schema) for the possible values. " +
					"When not configured, the current tenant value is kept.",
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"timeouts": commonschema.ResourceTimeouts(ctx),
		},
	}
}
