package graphBetaWindowsUpdatesAutopatchPolicyApproval

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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName  = "microsoft365_graph_beta_windows_updates_autopatch_policy_approval"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &WindowsUpdatesAutopatchPolicyApprovalResource{}
	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &WindowsUpdatesAutopatchPolicyApprovalResource{}
	// Enables import functionality
	_ resource.ResourceWithImportState = &WindowsUpdatesAutopatchPolicyApprovalResource{}
)

func NewWindowsUpdatesAutopatchPolicyApprovalResource() resource.Resource {
	return &WindowsUpdatesAutopatchPolicyApprovalResource{
		ReadPermissions: []string{
			"WindowsUpdates.ReadWrite.All",
		},
		WritePermissions: []string{
			"WindowsUpdates.ReadWrite.All",
		},
		ResourcePath: "/admin/windows/updates/policies",
	}
}

type WindowsUpdatesAutopatchPolicyApprovalResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

func (r *WindowsUpdatesAutopatchPolicyApprovalResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

func (r *WindowsUpdatesAutopatchPolicyApprovalResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

// ImportState handles import by splitting "{policy_id}/{approval_id}" on the first "/".
func (r *WindowsUpdatesAutopatchPolicyApprovalResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.SplitN(req.ID, "/", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		resp.Diagnostics.AddError(
			"Invalid import ID format",
			fmt.Sprintf("Expected format: {policy_id}/{approval_id}, got: %q", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("policy_id"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), parts[1])...)
}

func (r *WindowsUpdatesAutopatchPolicyApprovalResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a Windows Update policy approval using the `/admin/windows/updates/policies/{policyId}/approvals` endpoint. " +
			"A policy approval grants or suspends approval for a specific catalog entry (quality update) within a Windows Update policy. " +
			"The policy must already exist, managed by the `microsoft365_graph_beta_windows_updates_autopatch_policy` resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
				MarkdownDescription: "The unique identifier for the policy approval.",
			},
			"policy_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The ID of the Windows Update policy to which this approval belongs.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
				PlanModifiers: []planmodifier.String{
					planmodifiers.RequiresReplaceString(),
				},
			},
			"catalog_entry_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The catalog entry ID to approve. References a quality update from the Windows Update catalog. This field cannot be changed after creation — changes require the resource to be replaced.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
				PlanModifiers: []planmodifier.String{
					planmodifiers.RequiresReplaceString(),
				},
			},
			"status": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The approval status. Valid values are: `approved`, `suspended`.",
				Validators: []validator.String{
					stringvalidator.OneOf("approved", "suspended"),
				},
			},
			"created_date_time": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The date and time the policy approval was created (read-only).",
			},
			"last_modified_date_time": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The date and time the policy approval was last modified (read-only).",
			},
			"timeouts": commonschema.ResourceTimeouts(ctx),
		},
	}
}
