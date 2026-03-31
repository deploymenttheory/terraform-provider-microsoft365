package graphBetaAdministrativeUnit

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/identityschema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName  = "microsoft365_graph_beta_identity_and_access_administrative_unit"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &AdministrativeUnitResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &AdministrativeUnitResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &AdministrativeUnitResource{}

	// Enables plan modification/diff suppression
	_ resource.ResourceWithModifyPlan = &AdministrativeUnitResource{}

	// Enables identity schema for list resource support
	_ resource.ResourceWithIdentity = &AdministrativeUnitResource{}
)

func NewAdministrativeUnitResource() resource.Resource {
	return &AdministrativeUnitResource{
		ReadPermissions: []string{
			"Directory.Read.All",
		},
		WritePermissions: []string{
			"AdministrativeUnit.ReadWrite.All",
		},
	}
}

type AdministrativeUnitResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
}

// Metadata returns the resource type name.
func (r *AdministrativeUnitResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

// Configure sets the client for the resource.
func (r *AdministrativeUnitResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
//	terraform import microsoft365_graph_beta_identity_and_access_administrative_unit.example "11111111-1111-1111-1111-111111111111:hard_delete=true"
func (r *AdministrativeUnitResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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

// IdentitySchema defines the identity schema for this resource, used by list operations to uniquely identify instances
func (r *AdministrativeUnitResource) IdentitySchema(ctx context.Context, req resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	resp.IdentitySchema = identityschema.Schema{
		Attributes: map[string]identityschema.Attribute{
			"id": identityschema.StringAttribute{
				RequiredForImport: true,
			},
		},
	}
}

// Schema defines the schema for the resource.
func (r *AdministrativeUnitResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages an administrative unit in Microsoft Entra ID using the `/directory/administrativeUnits` endpoint. " +
			"Administrative units provide a conceptual container for user, group, and device directory objects, allowing delegation of administrative responsibilities.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Unique identifier for the administrative unit. Read-only.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(constants.GuidRegex),
						"must be a valid GUID in the format '00000000-0000-0000-0000-000000000000'",
					),
				},
			},
			"display_name": schema.StringAttribute{
				MarkdownDescription: "Display name for the administrative unit. Maximum length is 256 characters.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(256),
				},
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "An optional description for the administrative unit.",
				Optional:            true,
			},
			"is_member_management_restricted": schema.BoolAttribute{
				MarkdownDescription: "`true` if members of this administrative unit should be treated as sensitive, which requires specific permissions to manage. " +
					"If not set, the default value is `false`. Use this property to define administrative units with roles that don't inherit from tenant-level administrators, " +
					"and where the management of individual member objects is limited to administrators scoped to a restricted management administrative unit. " +
					"This property is immutable and can't be changed later.",
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(false),
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
			"membership_rule": schema.StringAttribute{
				MarkdownDescription: "The dynamic membership rule for the administrative unit. For more information about the rules you can use for dynamic administrative units and dynamic groups, " +
					"see [Manage rules for dynamic membership groups in Microsoft Entra ID](https://learn.microsoft.com/en-us/entra/identity/users/groups-dynamic-membership).",
				Optional: true,
			},
			"membership_rule_processing_state": schema.StringAttribute{
				MarkdownDescription: "Controls whether the dynamic membership rule is actively processed. Set to `On` to activate the dynamic membership rule, or `Paused` to stop updating membership dynamically.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("On", "Paused"),
				},
			},
			"membership_type": schema.StringAttribute{
				MarkdownDescription: "Indicates the membership type for the administrative unit. The possible values are: `Dynamic`, `Assigned`. If not set, the default behavior is assigned.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("Dynamic", "Assigned"),
				},
			},
			"visibility": schema.StringAttribute{
				MarkdownDescription: "Controls whether the administrative unit and its members are hidden or public. Can be set to `HiddenMembership` or `Public`. " +
					"If not set, the default behavior is public. When set to `HiddenMembership`, only members of the administrative unit can list other members of the administrative unit. " +
					"This property is immutable and can't be changed later.",
				Optional: true,
				Validators: []validator.String{
					stringvalidator.OneOf("HiddenMembership", "Public"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"hard_delete": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(false),
				MarkdownDescription: "When `true`, the administrative unit will be permanently deleted (hard delete) during destroy. " +
					"When `false` (default), the administrative unit will only be soft deleted and moved to the deleted items container where it can be restored within 30 days. " +
					"Note: This field defaults to `false` on import since the API does not return this value.",
			},
			"timeouts": commonschema.ResourceTimeouts(ctx),
		},
	}
}
