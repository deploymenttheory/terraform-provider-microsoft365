package graphBetaWindowsUpdatesAutopatchOperationalInsightsConnection

import (
	"context"
	"regexp"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
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
	ResourceName  = "microsoft365_graph_beta_windows_updates_autopatch_operational_insights_connection"
	CreateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	_ resource.Resource                = &WindowsUpdatesAutopatchOperationalInsightsConnectionResource{}
	_ resource.ResourceWithConfigure   = &WindowsUpdatesAutopatchOperationalInsightsConnectionResource{}
	_ resource.ResourceWithImportState = &WindowsUpdatesAutopatchOperationalInsightsConnectionResource{}
)

func NewWindowsUpdatesAutopatchOperationalInsightsConnectionResource() resource.Resource {
	return &WindowsUpdatesAutopatchOperationalInsightsConnectionResource{
		ReadPermissions: []string{
			"WindowsUpdates.Read.All",
		},
		WritePermissions: []string{
			"WindowsUpdates.ReadWrite.All",
		},
		ResourcePath: "/admin/windows/updates/resourceConnections",
	}
}

type WindowsUpdatesAutopatchOperationalInsightsConnectionResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

func (r *WindowsUpdatesAutopatchOperationalInsightsConnectionResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

func (r *WindowsUpdatesAutopatchOperationalInsightsConnectionResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

func (r *WindowsUpdatesAutopatchOperationalInsightsConnectionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *WindowsUpdatesAutopatchOperationalInsightsConnectionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a Windows Update Operational Insights connection using the `/admin/windows/updates/resourceConnections` endpoint. " +
			"An Operational Insights connection links Windows Update for Business deployment service to an Azure Log Analytics workspace, " +
			"enabling deployment reporting and insights. " +
			"This resource does not support in-place updates — any change to a field will destroy and recreate the connection.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
				MarkdownDescription: "The unique identifier for the Operational Insights connection.",
			},
			"azure_resource_group_name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of the Azure resource group that contains the Log Analytics workspace.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
				PlanModifiers: []planmodifier.String{
					planmodifiers.RequiresReplaceString(),
				},
			},
			"azure_subscription_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The Azure subscription ID that contains the Log Analytics workspace.",
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
			"workspace_name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of the Log Analytics workspace.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
				PlanModifiers: []planmodifier.String{
					planmodifiers.RequiresReplaceString(),
				},
			},
			"state": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The state of the connection. Possible values: `connected`, `notAuthorized`, `notFound`, `unknownFutureValue` (read-only).",
			},
			"timeouts": commonschema.ResourceTimeouts(ctx),
		},
	}
}
