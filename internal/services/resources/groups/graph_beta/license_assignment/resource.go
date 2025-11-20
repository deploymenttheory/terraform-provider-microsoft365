package graphBetaGroupLicenseAssignment

import (
	"context"
	"regexp"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	planmodifiers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/plan_modifiers"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName  = "microsoft365_graph_beta_groups_license_assignment"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &GroupLicenseAssignmentResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &GroupLicenseAssignmentResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &GroupLicenseAssignmentResource{}
)

func NewGroupLicenseAssignmentResource() resource.Resource {
	return &GroupLicenseAssignmentResource{
		ReadPermissions: []string{
			"Group.Read.All",
			"Directory.Read.All",
		},
		WritePermissions: []string{
			"LicenseAssignment.ReadWrite.All",
			"Group.ReadWrite.All",
			"Directory.ReadWrite.All",
		},
		ResourcePath: "/groups",
	}
}

type GroupLicenseAssignmentResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

// Metadata returns the resource type name.
func (r *GroupLicenseAssignmentResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

// Configure sets the client for the resource.
func (r *GroupLicenseAssignmentResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

// ImportState imports the resource state.
func (r *GroupLicenseAssignmentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("group_id"), req, resp)
}

// Schema returns the schema for the resource.
func (r *GroupLicenseAssignmentResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages group-based license assignments in Microsoft 365 using the `/groups/{groupId}/assignLicense` endpoint. This resource enables automatic license inheritance where all current and future group members receive the assigned licenses, providing centralized license management through Azure AD group membership.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
				MarkdownDescription: "The unique identifier for this license assignment resource. This is the same as the group_id.",
			},
			"group_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The unique identifier (UUID) for the group.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile(constants.GuidRegex), "Must be a valid UUID format (xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx)"),
				},
			},
			"display_name": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
				MarkdownDescription: "The display name of the group. This is computed and read-only.",
			},
			"add_licenses": schema.ListNestedAttribute{
				Optional:            true,
				MarkdownDescription: "A collection of licenses to assign to the group.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"sku_id": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: "The unique identifier (GUID) for the license SKU.",
							Validators: []validator.String{
								stringvalidator.RegexMatches(regexp.MustCompile(constants.GuidRegex), "Must be a valid UUID format (xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx)"),
							},
						},
						"disabled_plans": schema.SetAttribute{
							ElementType:         types.StringType,
							Optional:            true,
							MarkdownDescription: "A collection of the unique identifiers for service plans to disable.",
							Validators: []validator.Set{
								setvalidator.ValueStringsAre(
									stringvalidator.RegexMatches(regexp.MustCompile(constants.GuidRegex), "Each disabled plan must be a valid UUID format (xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx)"),
								),
							},
						},
					},
				},
			},
			"remove_licenses": schema.SetAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				MarkdownDescription: "A collection of SKU IDs that identify the licenses to remove from the group.",
				Validators: []validator.Set{
					setvalidator.ValueStringsAre(
						stringvalidator.RegexMatches(regexp.MustCompile(constants.GuidRegex), "Each license ID must be a valid UUID format (xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx)"),
					),
				},
			},
			"assigned_licenses": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "The current licenses assigned to the group. This is read-only.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"sku_id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The unique identifier (GUID) for the license SKU.",
						},
						"sku_part_number": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The string identifier of the license SKU, for example 'AAD_Premium'.",
						},
						"service_plans": schema.ListNestedAttribute{
							Computed:            true,
							MarkdownDescription: "The service plans available with this license.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"service_plan_id": schema.StringAttribute{
										Computed:            true,
										MarkdownDescription: "The unique identifier of the service plan.",
									},
									"service_plan_name": schema.StringAttribute{
										Computed:            true,
										MarkdownDescription: "The name of the service plan.",
									},
									"provisioning_status": schema.StringAttribute{
										Computed:            true,
										MarkdownDescription: "The provisioning status of the service plan.",
									},
									"applies_to": schema.StringAttribute{
										Computed:            true,
										MarkdownDescription: "The object the service plan can be assigned to.",
									},
								},
							},
						},
					},
				},
			},
			"timeouts": commonschema.Timeouts(ctx),
		},
	}
}
