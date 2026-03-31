package graphBetaWindowsUpdatesAutopatchPolicy

import (
	"context"

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
	ResourceName  = "microsoft365_graph_beta_windows_updates_autopatch_policy"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &WindowsUpdatesAutopatchPolicyResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &WindowsUpdatesAutopatchPolicyResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &WindowsUpdatesAutopatchPolicyResource{}

	// Enables plan modification/diff suppression
	_ resource.ResourceWithModifyPlan = &WindowsUpdatesAutopatchPolicyResource{}

	// Enables identity schema for list resource support
	_ resource.ResourceWithIdentity = &WindowsUpdatesAutopatchPolicyResource{}
)

func NewWindowsUpdatesAutopatchPolicyResource() resource.Resource {
	return &WindowsUpdatesAutopatchPolicyResource{
		ReadPermissions: []string{
			"WindowsUpdates.Read.All",
		},
		WritePermissions: []string{
			"WindowsUpdates.ReadWrite.All",
		},
		ResourcePath: "/admin/windows/updates/policies",
	}
}

type WindowsUpdatesAutopatchPolicyResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

// Metadata returns the resource type name.
func (r *WindowsUpdatesAutopatchPolicyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

// Configure sets the client for the resource.
func (r *WindowsUpdatesAutopatchPolicyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

// ImportState imports the resource state.
func (r *WindowsUpdatesAutopatchPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// IdentitySchema defines the identity schema for this resource, used by list operations to uniquely identify instances.
func (r *WindowsUpdatesAutopatchPolicyResource) IdentitySchema(ctx context.Context, req resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	resp.IdentitySchema = identityschema.Schema{
		Attributes: map[string]identityschema.Attribute{
			"id": identityschema.StringAttribute{
				RequiredForImport: true,
			},
		},
	}
}

func (r *WindowsUpdatesAutopatchPolicyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a Windows Autopatch quality update policy using the `/admin/windows/updates/policies` endpoint. " +
			"A quality update policy defines approval rules that determine which published Windows quality updates are automatically approved " +
			"for deployment based on classification (security/nonSecurity) and cadence (monthly/outOfBand).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The unique identifier for the Windows Autopatch policy.",
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
			},
			"display_name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The policy display name. The maximum length is 200 characters.",
				Validators: []validator.String{
					stringvalidator.LengthAtMost(200),
				},
			},
			"description": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The policy description. The maximum length is 1,500 characters.",
				Validators: []validator.String{
					stringvalidator.LengthAtMost(1500),
				},
			},
			"created_date_time": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The date and time when the policy was created. Read-only.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"last_modified_date_time": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The date and time when the policy was last modified. Read-only.",
			},
			"approval_rules": schema.SetNestedAttribute{
				Optional:            true,
				MarkdownDescription: "The approval rules that determine which published content matches the rule on an ongoing basis.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"deferral_in_days": schema.Int32Attribute{
							Required:            true,
							MarkdownDescription: "The Windows update deferral period in days. The value must be between 0 and 30.",
							Validators: []validator.Int32{
								int32validator.Between(0, 30),
							},
						},
						"classification": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: "The classification of the quality update. Possible values are: `all`, `security`, `nonSecurity`.",
							Validators: []validator.String{
								stringvalidator.OneOf("all", "security", "nonSecurity"),
							},
						},
						"cadence": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: "The cadence of the quality update. Possible values are: `monthly`, `outOfBand`.",
							Validators: []validator.String{
								stringvalidator.OneOf("monthly", "outOfBand"),
							},
						},
					},
				},
			},
			"timeouts": commonschema.ResourceTimeouts(ctx),
		},
	}
}
