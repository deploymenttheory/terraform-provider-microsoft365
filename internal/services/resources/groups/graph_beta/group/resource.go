package graphBetaGroup

import (
	"context"
	"regexp"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common"
	planmodifiers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/plan_modifiers"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/validators"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName  = "graph_beta_groups_group"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &GroupResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &GroupResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &GroupResource{}

	// Enables plan modification/diff suppression
	_ resource.ResourceWithModifyPlan = &GroupResource{}
)

func NewGroupResource() resource.Resource {
	return &GroupResource{
		ReadPermissions: []string{
			"Group.Read.All",
			"Directory.Read.All",
		},
		WritePermissions: []string{
			"Group.ReadWrite.All",
			"Directory.ReadWrite.All",
		},
		ResourcePath: "/groups",
	}
}

type GroupResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ProviderTypeName string
	TypeName         string
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

// Metadata returns the resource type name.
func (r *GroupResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	r.ProviderTypeName = req.ProviderTypeName
	r.TypeName = ResourceName
	resp.TypeName = req.ProviderTypeName + "_" + ResourceName
}

// FullTypeName returns the full type name of the resource for logging purposes.
func (r *GroupResource) FullTypeName() string {
	return r.ProviderTypeName + "_" + r.TypeName
}

// Configure sets the client for the resource.
func (r *GroupResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = common.SetGraphBetaClientForResource(ctx, req, resp, constants.PROVIDER_NAME+"_"+ResourceName)
}

// ImportState imports the resource state.
func (r *GroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Schema defines the schema for the resource.
func (r *GroupResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages Azure AD/Entra groups using the `/groups` endpoint. This resource enables creation and management of security groups, Microsoft 365 groups, and distribution groups with support for dynamic membership, role assignment capabilities, and comprehensive group configuration options for organizational identity and access management.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
				MarkdownDescription: "The unique identifier for the group. Read-only.",
			},
			"display_name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The display name for the group. This property is required when a group is created and can't be cleared during updates. Maximum length is 256 characters.",
				Validators: []validator.String{
					validators.StringLengthAtMost(256),
				},
			},
			"description": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "An optional description for the group.",
			},
			"mail_nickname": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The mail alias for the group, unique for Microsoft 365 groups in the organization. Maximum length is 64 characters. This property can contain only characters in the ASCII character set 0 - 127 except the following: @ () \\ [] \" ; : <> , SPACE.",
				Validators: []validator.String{
					validators.StringLengthAtMost(64),
					validators.ASCIIOnly(),
					validators.IllegalCharactersInString([]rune{'@', '(', ')', '\\', '[', ']', '"', ';', ':', '<', '>', ',', ' '}, "mail nickname cannot contain: @ () \\ [] \" ; : <> , SPACE"),
				},
			},
			"mail_enabled": schema.BoolAttribute{
				Required:            true,
				MarkdownDescription: "Specifies whether the group is mail-enabled. Required.",
			},
			"security_enabled": schema.BoolAttribute{
				Required:            true,
				MarkdownDescription: "Specifies whether the group is a security group. Required.",
			},
			"group_types": schema.SetAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				MarkdownDescription: "Specifies the group type and its membership. If the collection contains 'Unified', the group is a Microsoft 365 group; otherwise, it's either a security group or a distribution group. If the collection includes 'DynamicMembership', the group has dynamic membership; otherwise, membership is static.",
				Validators: []validator.Set{
					setvalidator.ValueStringsAre(
						stringvalidator.OneOf("Unified", "DynamicMembership"),
					),
				},
			},
			"visibility": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("Private"),
				MarkdownDescription: "Specifies the group join policy and group content visibility for groups. Possible values are: Private, Public, or HiddenMembership. Default is 'Private'.",
				Validators: []validator.String{
					stringvalidator.OneOf("Private", "Public", "HiddenMembership"),
				},
			},
			"is_assignable_to_role": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Indicates whether this group can be assigned to a Microsoft Entra role. This property can only be set while creating the group and is immutable. If set to true, the securityEnabled property must also be set to true, visibility must be Hidden, and the group can't be a dynamic group. Default is false.",
			},
			"membership_rule": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The rule that determines members for this group if the group is a dynamic group (groupTypes contains DynamicMembership). For more information about the syntax of the membership rule, see Membership Rules syntax.",
			},
			"membership_rule_processing_state": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("Paused"),
				MarkdownDescription: "Indicates whether the dynamic membership processing is on or paused. Possible values are 'On' or 'Paused'. Default is 'Paused'.",
				Validators: []validator.String{
					stringvalidator.OneOf("On", "Paused"),
				},
			},
			"created_date_time": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Timestamp of when the group was created. The value can't be modified and is automatically populated when the group is created. Read-only.",
			},
			"mail": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The SMTP address for the group. Read-only.",
			},
			"proxy_addresses": schema.SetAttribute{
				ElementType:         types.StringType,
				Computed:            true,
				MarkdownDescription: "Email addresses for the group that direct to the same group mailbox. Read-only.",
			},
			"on_premises_sync_enabled": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "true if this group is synced from an on-premises directory; false if this group was originally synced from an on-premises directory but is no longer synced; null if this object has never synced from an on-premises directory. Read-only.",
			},
			"preferred_data_location": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The preferred data location for the Microsoft 365 group. By default, the group inherits the group creator's preferred data location.",
			},
			"preferred_language": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The preferred language for a Microsoft 365 group. Should follow ISO 639-1 Code; for example, en-US.",
				Validators: []validator.String{
					validators.RegexMatches(regexp.MustCompile(`^[a-z]{2}(-[A-Z]{2})?$`), "language code must follow ISO 639-1 format (e.g., en, en-US)"),
				},
			},
			"theme": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Specifies a Microsoft 365 group's color theme. Possible values are Teal, Purple, Green, Blue, Pink, Orange, or Red.",
				Validators: []validator.String{
					stringvalidator.OneOf("Teal", "Purple", "Green", "Blue", "Pink", "Orange", "Red"),
				},
			},
			"classification": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Describes a classification for the group (such as low, medium, or high business impact). Valid values for this property are defined by creating a ClassificationList setting value, based on the template definition.",
			},
			"expiration_date_time": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Timestamp of when the group is set to expire. It's null for security groups, but for Microsoft 365 groups, it represents when the group is set to expire as defined in the groupLifecyclePolicy. Read-only.",
			},
			"renewed_date_time": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Timestamp of when the group was last renewed. This value can't be modified directly and is only updated via the renew service action. Read-only.",
			},
			"security_identifier": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Security identifier of the group, used in Windows scenarios. Read-only.",
			},
			"timeouts": commonschema.Timeouts(ctx),
		},
	}
}
