package graphBetaNetworkProxyAutoConfiguration

import (
	"context"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
)

const (
	ResourceName  = "microsoft365_graph_beta_identity_and_access_network_proxy_auto_configuration"
	CreateTimeout = 180
	ReadTimeout   = 180
	UpdateTimeout = 180
	DeleteTimeout = 180
)

var (
	_ resource.Resource                = &NetworkProxyAutoConfigurationResource{}
	_ resource.ResourceWithConfigure   = &NetworkProxyAutoConfigurationResource{}
	_ resource.ResourceWithImportState = &NetworkProxyAutoConfigurationResource{}
)

// NetworkProxyAutoConfigurationResource manages the Graph configuration.
type NetworkProxyAutoConfigurationResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

func NewNetworkProxyAutoConfigurationResource() resource.Resource {
	return &NetworkProxyAutoConfigurationResource{
		ReadPermissions:  []string{"NetworkAccess.ReadWrite.All"},
		WritePermissions: []string{"NetworkAccess.ReadWrite.All"},
		ResourcePath:     "/networkAccess/explicitForwardProxyConfig/proxyAutoConfigurationFiles",
	}
}

func (r *NetworkProxyAutoConfigurationResource) Metadata(
	ctx context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = ResourceName
}

func (r *NetworkProxyAutoConfigurationResource) Configure(
	ctx context.Context,
	req resource.ConfigureRequest,
	resp *resource.ConfigureResponse,
) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

func (r *NetworkProxyAutoConfigurationResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	if _, err := uuid.Parse(req.ID); err != nil {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			"Use the custom PAC object UUID returned by Graph, not a name or public URL.",
		)
		return
	}
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *NetworkProxyAutoConfigurationResource) Schema(
	ctx context.Context,
	req resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages the name, JavaScript content and availability of a hosted custom PAC file using `/networkAccess/explicitForwardProxyConfig/proxyAutoConfigurationFiles`. Initial apply creates a file; destroy deletes it.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Graph custom PAC object UUID.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Unique hosted PAC name, without the automatically appended `.pac` extension.",
				Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
			},
			"content": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "PAC JavaScript body, supplied inline or with `file()`. Use `$${GSAEFP}` in Terraform strings to preserve the service placeholder. Required and nonempty; omission, null and empty strings do not delete or reset the file. Microsoft documents a 950 KB service limit and recommends no more than 250 KB. The service validates JavaScript syntax.",
				Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
			},
			"is_enabled": schema.BoolAttribute{
				Required:            true,
				MarkdownDescription: "Whether clients can retrieve this custom PAC. `false` uploads a disabled file.",
			},
			"created_date_time": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Creation timestamp returned by Graph.",
			},
			"last_modified_date_time": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Last modification timestamp returned by Graph.",
			},
			"timeouts": commonschema.ResourceTimeouts(ctx),
		},
	}
}
