package graphBetaGroup

import (
	"context"
	"regexp"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
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
			"Group.Create",
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
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, constants.PROVIDER_NAME+"_"+ResourceName)
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
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("Private"),
				MarkdownDescription: "Specifies the group join policy and group content visibility for groups. " +
					"Possible values are: `Private`, `Public`, or `HiddenMembership`. `HiddenMembership` can be set " +
					"only for Microsoft 365 groups when the groups are created and cannot be updated later. Other values " +
					"of visibility can be updated after group creation. If visibility value is not specified during group " +
					"creation, a security group is created as `Private` by default, and a Microsoft 365 group is `Public`. " +
					"Groups assignable to roles are always `Private`. Returned by default. Nullable.",
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
				MarkdownDescription: "Indicates whether the dynamic membership processing is on or paused. Possible values are 'On' or 'Paused'. Only applicable for dynamic groups (when groupTypes contains DynamicMembership).",
				Validators: []validator.String{
					stringvalidator.OneOf("On", "Paused"),
				},
			},
			"created_date_time": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Timestamp of when the group was created. The value can't be modified and is automatically populated when the group is created. Read-only.",
			},
			"group_owners": schema.SetAttribute{
				ElementType: types.StringType,
				Optional:    true,
				MarkdownDescription: "The owners of the group at creation time. A maximum of 20 relationships, such as owners and members, can be added as part of group creation. " +
					"Specify the user IDs (GUIDs) of the users who should be owners of the group. Note: A non-admin user cannot add themselves to the group owners collection. " +
					"Owners can be added after creation using the `/groups/{id}/owners/$ref` endpoint.",
				Validators: []validator.Set{
					setvalidator.SizeAtMost(20),
					setvalidator.ValueStringsAre(
						validators.RegexMatches(regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`), "value must be a valid UUID/GUID"),
					),
				},
			},
			"group_members": schema.SetAttribute{
				ElementType: types.StringType,
				Optional:    true,
				MarkdownDescription: "The members of the group at creation time. A maximum of 20 relationships, such as owners and members, can be added as part of group creation. " +
					"Specify the user IDs (GUIDs) of the users who should be members of the group. " +
					"Additional members can be added after creation using the `/groups/{id}/members/$ref` endpoint or JSON batching.",
				Validators: []validator.Set{
					setvalidator.SizeAtMost(20),
					setvalidator.ValueStringsAre(
						validators.RegexMatches(regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`), "value must be a valid UUID/GUID"),
					),
				},
			},
			"timeouts": commonschema.Timeouts(ctx),
		},
	}
}
