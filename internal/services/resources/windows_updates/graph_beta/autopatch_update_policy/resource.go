package graphBetaWindowsUpdatesAutopatchUpdatePolicy

import (
	"context"
	"regexp"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	planmodifiers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/plan_modifiers"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/identityschema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName  = "microsoft365_graph_beta_windows_updates_update_policy"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &WindowsUpdatesAutopatchUpdatePolicyResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &WindowsUpdatesAutopatchUpdatePolicyResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &WindowsUpdatesAutopatchUpdatePolicyResource{}

	// Enables identity schema for list resource support
	_ resource.ResourceWithIdentity = &WindowsUpdatesAutopatchUpdatePolicyResource{}
)

func NewWindowsUpdatesAutopatchUpdatePolicyResource() resource.Resource {
	return &WindowsUpdatesAutopatchUpdatePolicyResource{
		ReadPermissions: []string{
			"WindowsUpdates.Read.All",
		},
		WritePermissions: []string{
			"WindowsUpdates.ReadWrite.All",
		},
		ResourcePath: "/admin/windows/updates/updatePolicies",
	}
}

type WindowsUpdatesAutopatchUpdatePolicyResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

func (r *WindowsUpdatesAutopatchUpdatePolicyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

func (r *WindowsUpdatesAutopatchUpdatePolicyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

func (r *WindowsUpdatesAutopatchUpdatePolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *WindowsUpdatesAutopatchUpdatePolicyResource) IdentitySchema(ctx context.Context, req resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	resp.IdentitySchema = identityschema.Schema{
		Attributes: map[string]identityschema.Attribute{
			"id": identityschema.StringAttribute{
				RequiredForImport: true,
			},
		},
	}
}

func (r *WindowsUpdatesAutopatchUpdatePolicyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a Windows Update policy using the `/admin/windows/updates/updatePolicies` endpoint. " +
			"An update policy serves as a container for compliance changes (content approvals) that define which updates " +
			"should be deployed to devices. This resource is a prerequisite for creating content approvals.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The unique identifier for the update policy.",
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
			},
			"created_date_time": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The date and time when the update policy was created. Read-only.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"audience_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The ID of the deployment audience to target with this policy.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"compliance_changes": schema.BoolAttribute{
				Required:            true,
				MarkdownDescription: "Enable compliance changes (content approvals) for this policy. Must be set to `true` to create content approvals.",
			},
		"compliance_change_rules": schema.SetNestedAttribute{
			Optional:            true,
			Computed:            true,
			MarkdownDescription: "Rules for governing the automatic creation of compliance changes. Cannot be updated after creation - changes require resource replacement.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"content_filter": schema.SingleNestedAttribute{
							Optional:            true,
							MarkdownDescription: "The content filter for the compliance change rule.",
							Attributes: map[string]schema.Attribute{
								"filter_type": schema.StringAttribute{
									Required:            true,
									MarkdownDescription: "The type of content filter. Only `driverUpdateFilter` is supported. Note: `windowsUpdateFilter` is not valid despite appearing in API documentation.",
									Validators: []validator.String{
										stringvalidator.OneOf("driverUpdateFilter"),
									},
									PlanModifiers: []planmodifier.String{
										stringplanmodifier.RequiresReplace(),
									},
								},
							},
						},
					"duration_before_deployment_start": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "The duration before deployment starts (ISO 8601 duration format). Valid range: P1D to P30D (1 to 30 days). Cannot be updated after creation. Examples: 'P7D' for 7 days, 'P14D' for 14 days.",
						Validators: []validator.String{
							stringvalidator.RegexMatches(
								regexp.MustCompile(`^P([1-9]|[12][0-9]|30)D$`),
								"must be in ISO 8601 duration format P#D where # is between 1 and 30 (e.g., P7D, P14D, P30D)",
							),
						},
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
						"created_date_time": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The date and time when the rule was created. Read-only.",
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
						},
						"last_evaluated_date_time": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The date and time when the rule was last evaluated. Read-only.",
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
						},
						"last_modified_date_time": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The date and time when the rule was last modified. Read-only.",
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
						},
					},
				},
			},
			"deployment_settings": schema.SingleNestedAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Settings for governing how to deploy content.",
				Attributes: map[string]schema.Attribute{
					"schedule": schema.SingleNestedAttribute{
						Optional:            true,
						MarkdownDescription: "Settings for the schedule of the deployment.",
						Attributes: map[string]schema.Attribute{
							"start_date_time": schema.StringAttribute{
								Optional:            true,
								MarkdownDescription: "The start date and time for the deployment (ISO 8601 format).",
							},
							"gradual_rollout": schema.SingleNestedAttribute{
								Optional:            true,
								MarkdownDescription: "Settings for gradual rollout.",
								Attributes: map[string]schema.Attribute{
								"duration_between_offers": schema.StringAttribute{
									Required:            true,
									MarkdownDescription: "The duration between offers (ISO 8601 duration format). Valid range: P1D to P30D (1 to 30 days). Examples: 'P1D' for 1 day, 'P7D' for 7 days.",
									Validators: []validator.String{
										stringvalidator.RegexMatches(
											regexp.MustCompile(`^P([1-9]|[12][0-9]|30)D$`),
											"must be in ISO 8601 duration format P#D where # is between 1 and 30 (e.g., P1D, P7D, P30D)",
										),
									},
								},
									"devices_per_offer": schema.Int32Attribute{
										Required:            true,
										MarkdownDescription: "The number of devices to offer the update to in each batch. Must be a positive integer (minimum 1).",
										Validators: []validator.Int32{
											int32validator.AtLeast(1),
										},
									},
								},
							},
						},
					},
				},
			},
			"timeouts": commonschema.ResourceTimeouts(ctx),
		},
	}
}
