package graphBetaWindowsUpdatesAutopatchContentApproval

import (
	"context"
	"fmt"
	"strings"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	planmodifiers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/plan_modifiers"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/identityschema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName  = "microsoft365_graph_beta_windows_updates_autopatch_content_approval"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &WindowsUpdatesAutopatchContentApprovalResource{}
	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &WindowsUpdatesAutopatchContentApprovalResource{}
	// Enables import functionality
	_ resource.ResourceWithImportState = &WindowsUpdatesAutopatchContentApprovalResource{}
	// Enables identity schema for list resource support
	_ resource.ResourceWithIdentity = &WindowsUpdatesAutopatchContentApprovalResource{}
)

func NewWindowsUpdatesAutopatchContentApprovalResource() resource.Resource {
	return &WindowsUpdatesAutopatchContentApprovalResource{
		ReadPermissions: []string{
			"WindowsUpdates.Read.All",
		},
		WritePermissions: []string{
			"WindowsUpdates.ReadWrite.All",
		},
		ResourcePath: "/admin/windows/updates/updatePolicies",
	}
}

type WindowsUpdatesAutopatchContentApprovalResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

func (r *WindowsUpdatesAutopatchContentApprovalResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

func (r *WindowsUpdatesAutopatchContentApprovalResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

// ImportState handles import by splitting "{update_policy_id}/{compliance_change_id}" on the first "/".
func (r *WindowsUpdatesAutopatchContentApprovalResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.SplitN(req.ID, "/", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		resp.Diagnostics.AddError(
			"Invalid import ID format",
			fmt.Sprintf("Expected format: {update_policy_id}/{compliance_change_id}, got: %q", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("update_policy_id"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), parts[1])...)
}

func (r *WindowsUpdatesAutopatchContentApprovalResource) IdentitySchema(ctx context.Context, req resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	resp.IdentitySchema = identityschema.Schema{
		Attributes: map[string]identityschema.Attribute{
			"id": identityschema.StringAttribute{
				RequiredForImport: true,
			},
		},
	}
}

func (r *WindowsUpdatesAutopatchContentApprovalResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages Windows Update content approvals using the `/admin/windows/updates/updatePolicies/{updatePolicyId}/complianceChanges` endpoint. " +
			"Content approvals define which updates (feature or quality) should be deployed to devices according to a specific update policy. " +
			"This resource requires an existing update policy and references catalog entries from the Windows Update catalog.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
				MarkdownDescription: "Unique identifier for the content approval (compliance change ID).",
			},
			"update_policy_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The ID of the Windows Update policy to which this content approval belongs.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
				PlanModifiers: []planmodifier.String{
					planmodifiers.RequiresReplaceString(),
				},
			},
			"catalog_entry_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The ID of the catalog entry to approve for deployment. This should reference a feature or quality update from the Windows Update catalog.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"catalog_entry_type": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The type of catalog entry being approved. Valid values are: `featureUpdate`, `qualityUpdate`.",
				Validators: []validator.String{
					stringvalidator.OneOf("featureUpdate", "qualityUpdate"),
				},
			},
			"is_revoked": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Set to `true` to revoke the content approval, preventing further deployment. Revoking is a final action and cannot be undone.",
			},
			"created_date_time": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The date and time when the content approval was created.",
			},
			"revoked_date_time": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The date and time when the content approval was revoked.",
			},
			"deployment_settings": schema.SingleNestedAttribute{
				Optional:            true,
				MarkdownDescription: "Settings for governing how to deploy the approved content.",
				Attributes: map[string]schema.Attribute{
					"schedule": schema.SingleNestedAttribute{
						Optional:            true,
						MarkdownDescription: "Schedule settings for the deployment.",
						Attributes: map[string]schema.Attribute{
							"start_date_time": schema.StringAttribute{
								Optional:            true,
								MarkdownDescription: "The date and time when the deployment should start.",
							},
							"gradual_rollout": schema.SingleNestedAttribute{
								Optional:            true,
								MarkdownDescription: "Settings for gradual rollout of the deployment.",
								Attributes: map[string]schema.Attribute{
									"end_date_time": schema.StringAttribute{
										Optional:            true,
										MarkdownDescription: "The date and time when the gradual rollout should complete.",
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
