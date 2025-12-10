package graphBetaGroup

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	planmodifiers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/plan_modifiers"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	validate "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/validate/attribute"
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
	"github.com/hashicorp/terraform-plugin-log/tflog"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName  = "microsoft365_graph_beta_groups_group"
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
			"RoleManagement.ReadWrite.Directory",
		},
		ResourcePath: "/groups",
	}
}

type GroupResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

// Metadata returns the resource type name.
func (r *GroupResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

// Configure sets the client for the resource.
func (r *GroupResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

// ImportState handles importing the resource with an extended ID format.
//
// Supported formats:
//   - Simple:   "resource_id" (hard_delete defaults to false)
//   - Extended: "resource_id:hard_delete=true" or "resource_id:hard_delete=false"
//
// Example:
//
//	terraform import microsoft365_graph_beta_groups_group.example "12345678-1234-1234-1234-123456789012:hard_delete=true"
func (r *GroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.Split(req.ID, ":")
	resourceID := idParts[0]
	hardDelete := false // Default to soft delete for safety

	if len(idParts) > 1 {
		for _, part := range idParts[1:] {
			if strings.HasPrefix(part, "hard_delete=") {
				value := strings.TrimPrefix(part, "hard_delete=")
				switch strings.ToLower(value) {
				case "true":
					hardDelete = true
				case "false":
					hardDelete = false
				default:
					resp.Diagnostics.AddError(
						"Invalid Import ID",
						fmt.Sprintf("Invalid hard_delete value '%s'. Must be 'true' or 'false'.", value),
					)
					return
				}
			}
		}
	}

	tflog.Info(ctx, fmt.Sprintf("Importing %s with ID: %s, hard_delete: %t", ResourceName, resourceID, hardDelete))

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), resourceID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("hard_delete"), hardDelete)...)
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
					validate.StringLengthAtMost(256),
				},
			},
			"description": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "An optional description for the group. May be auto-populated by the API for certain group types.",
			},
			"mail_nickname": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The mail alias for the group, unique for Microsoft 365 groups in the organization. Maximum length is 64 characters. This property can contain only characters in the ASCII character set 0 - 127 except the following: @ () \\ [] \" ; : <> , SPACE.",
				Validators: []validator.String{
					validate.StringLengthAtMost(64),
					validate.ASCIIOnly(),
					validate.IllegalCharactersInString([]rune{'@', '(', ')', '\\', '[', ']', '"', ';', ':', '<', '>', ',', ' '}, "mail nickname cannot contain: @ () \\ [] \" ; : <> , SPACE"),
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
						validate.RegexMatches(regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`), "value must be a valid UUID/GUID"),
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
						validate.RegexMatches(regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`), "value must be a valid UUID/GUID"),
					),
				},
			},
			"hard_delete": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(false),
				MarkdownDescription: "When `true`, the group will be permanently deleted (hard delete) during destroy. " +
					"When `false` (default), the group will only be soft deleted and moved to the deleted items container where it can be restored within 30 days. " +
					"Note: This field defaults to `false` on import since the API does not return this value.",
			},
			"timeouts": commonschema.Timeouts(ctx),
		},
	}
}
