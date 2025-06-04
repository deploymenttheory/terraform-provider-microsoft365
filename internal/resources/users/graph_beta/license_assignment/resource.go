package graphBetaUserLicenseAssignment

import (
	"context"
	"regexp"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common"
	planmodifiers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/plan_modifiers"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/schema"
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
	ResourceName  = "graph_beta_users_user_license_assignment"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &UserLicenseAssignmentResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &UserLicenseAssignmentResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &UserLicenseAssignmentResource{}

	// Compiled regex for UUID validation
	uuidRegex = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)

	// Compiled regex for UPN validation (user principal name)
	upnRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
)

func NewUserLicenseAssignmentResource() resource.Resource {
	return &UserLicenseAssignmentResource{
		ReadPermissions: []string{
			"User.Read.All",
			"Directory.Read.All",
		},
		WritePermissions: []string{
			"User.ReadWrite.All",
			"Directory.ReadWrite.All",
		},
		ResourcePath: "/users",
	}
}

type UserLicenseAssignmentResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ProviderTypeName string
	TypeName         string
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

// Metadata returns the resource type name.
func (r *UserLicenseAssignmentResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	r.ProviderTypeName = req.ProviderTypeName
	r.TypeName = ResourceName
	resp.TypeName = r.FullTypeName()
}

// FullTypeName returns the full resource type name in the format "providername_resourcename".
func (r *UserLicenseAssignmentResource) FullTypeName() string {
	return r.ProviderTypeName + "_" + r.TypeName
}

// Configure sets the client for the resource.
func (r *UserLicenseAssignmentResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = common.SetGraphBetaClientForResource(ctx, req, resp, constants.PROVIDER_NAME+"_"+ResourceName)
}

// ImportState imports the resource state.
func (r *UserLicenseAssignmentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("user_id"), req, resp)
}

// Schema returns the schema for the resource.
func (r *UserLicenseAssignmentResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages Microsoft 365 license assignments for individual users using the `/users/{userId}/assignLicense` endpoint. This resource enables direct license assignment to users, allowing administrators to grant or revoke access to Microsoft cloud services like Office 365, Enterprise Mobility + Security, and Windows licenses.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
				MarkdownDescription: "The unique identifier for this license assignment resource. This is the same as the user_id.",
			},
			"user_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The unique identifier for the user. Can be either the object ID (UUID) or user principal name (UPN).",
				Validators: []validator.String{
					stringvalidator.Any(
						stringvalidator.RegexMatches(uuidRegex, "Must be a valid UUID format (xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx)"),
						stringvalidator.RegexMatches(upnRegex, "Must be a valid User Principal Name format (user@domain.com)"),
					),
				},
			},
			"user_principal_name": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
				MarkdownDescription: "The user principal name (UPN) of the user. This is computed and read-only.",
			},
			"add_licenses": schema.ListNestedAttribute{
				Optional:            true,
				MarkdownDescription: "A collection of licenses to assign to the user.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"sku_id": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: "The unique identifier (GUID) for the license SKU.",
							Validators: []validator.String{
								stringvalidator.RegexMatches(uuidRegex, "Must be a valid UUID format (xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx)"),
							},
						},
						"disabled_plans": schema.SetAttribute{
							ElementType:         types.StringType,
							Optional:            true,
							MarkdownDescription: "A collection of the unique identifiers for service plans to disable.",
							Validators: []validator.Set{
								setvalidator.ValueStringsAre(
									stringvalidator.RegexMatches(uuidRegex, "Each disabled plan must be a valid UUID format (xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx)"),
								),
							},
						},
					},
				},
			},
			"remove_licenses": schema.SetAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				MarkdownDescription: "A collection of SKU IDs that identify the licenses to remove from the user.",
				Validators: []validator.Set{
					setvalidator.ValueStringsAre(
						stringvalidator.RegexMatches(uuidRegex, "Each license ID must be a valid UUID format (xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx)"),
					),
				},
			},
			"assigned_licenses": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "The current licenses assigned to the user. This is read-only.",
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
