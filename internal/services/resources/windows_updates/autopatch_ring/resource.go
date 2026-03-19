package graphBetaWindowsUpdatesAutopatchRing

import (
	"context"
	"fmt"
	"strings"

	"regexp"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	planmodifiers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/plan_modifiers"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName  = "microsoft365_graph_beta_windows_updates_autopatch_ring"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	_ resource.Resource                = &WindowsUpdatesAutopatchRingResource{}
	_ resource.ResourceWithConfigure   = &WindowsUpdatesAutopatchRingResource{}
	_ resource.ResourceWithImportState = &WindowsUpdatesAutopatchRingResource{}
)

func NewWindowsUpdatesAutopatchRingResource() resource.Resource {
	return &WindowsUpdatesAutopatchRingResource{
		ReadPermissions: []string{
			"WindowsUpdates.ReadWrite.All",
		},
		WritePermissions: []string{
			"WindowsUpdates.ReadWrite.All",
		},
		ResourcePath: "/admin/windows/updates/policies",
	}
}

type WindowsUpdatesAutopatchRingResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

func (r *WindowsUpdatesAutopatchRingResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

func (r *WindowsUpdatesAutopatchRingResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

// ImportState handles import by splitting "{policy_id}/{ring_id}" on the first "/".
func (r *WindowsUpdatesAutopatchRingResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.SplitN(req.ID, "/", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		resp.Diagnostics.AddError(
			"Invalid import ID format",
			fmt.Sprintf("Expected format: {policy_id}/{ring_id}, got: %q", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("policy_id"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), parts[1])...)
}

func (r *WindowsUpdatesAutopatchRingResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a Windows Update policy ring using the `/admin/windows/updates/policies/{policyId}/rings` endpoint. " +
			"A ring defines the deployment audience, deferral, and pause settings for quality updates within a Windows Update policy. " +
			"The policy must already exist, managed by the `microsoft365_graph_beta_windows_updates_autopatch_policy` resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
				MarkdownDescription: "The unique identifier for the ring.",
			},
			"policy_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The ID of the Windows Update policy to which this ring belongs.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(constants.GuidRegex),
						"Must be a valid UUID format (xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx)",
					),
				},
				PlanModifiers: []planmodifier.String{
					planmodifiers.RequiresReplaceString(),
				},
			},
			"display_name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The ring display name. The maximum length is 200 characters.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 200),
				},
			},
			"description": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The ring description. The maximum length is 1,500 characters.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 1500),
				},
			},
			"is_paused": schema.BoolAttribute{
				Required:            true,
				MarkdownDescription: "Whether the ring is paused. When `true`, quality update deployment to devices in this ring is halted.",
			},
			"deferral_in_days": schema.Int32Attribute{
				Optional:            true,
				MarkdownDescription: "The quality update deferral period in days. The value must be between 0 and 30.",
				Validators: []validator.Int32{
					int32validator.Between(0, 30),
				},
			},
			"included_group_assignment": schema.SingleNestedAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Defines the Microsoft Entra groups whose devices are included in this ring's deployment audience. If not set, an empty assignment is sent.",
				Attributes: map[string]schema.Attribute{
					"assignments": schema.SetNestedAttribute{
						Optional:            true,
						Computed:            true,
						MarkdownDescription: "A set of group assignments governing the included deployment audience.",
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"group_id": schema.StringAttribute{
									Required:            true,
									MarkdownDescription: "The Microsoft Entra group ID to include in this ring.",
									Validators: []validator.String{
										stringvalidator.RegexMatches(
											regexp.MustCompile(constants.GuidRegex),
											"Must be a valid UUID format (xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx)",
										),
									},
								},
							},
						},
					},
				},
			},
			"excluded_group_assignment": schema.SingleNestedAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Defines the Microsoft Entra groups whose devices are excluded from this ring's deployment audience. If not set, an empty assignment is sent.",
				Attributes: map[string]schema.Attribute{
					"assignments": schema.SetNestedAttribute{
						Optional:            true,
						Computed:            true,
						MarkdownDescription: "A set of group assignments governing the excluded deployment audience.",
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"group_id": schema.StringAttribute{
									Required:            true,
									MarkdownDescription: "The Microsoft Entra group ID to exclude from this ring.",
									Validators: []validator.String{
										stringvalidator.RegexMatches(
											regexp.MustCompile(constants.GuidRegex),
											"Must be a valid UUID format (xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx)",
										),
									},
								},
							},
						},
					},
				},
			},
			"is_hotpatch_enabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Whether hotpatch updates are enabled for this ring (quality update rings only). Hotpatch updates apply without requiring a device restart.",
			},
			"created_date_time": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The date and time the ring was created (read-only).",
			},
			"last_modified_date_time": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The date and time the ring was last modified (read-only).",
			},
			"timeouts": commonschema.ResourceTimeouts(ctx),
		},
	}
}
