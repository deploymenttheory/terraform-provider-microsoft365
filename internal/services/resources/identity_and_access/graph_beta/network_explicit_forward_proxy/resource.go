package graphBetaNetworkExplicitForwardProxy

import (
	"context"

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
	ResourceName  = "microsoft365_graph_beta_identity_and_access_network_explicit_forward_proxy"
	CreateTimeout = 180
	ReadTimeout   = 180
	UpdateTimeout = 180
	DeleteTimeout = 180
)

var (
	_ resource.Resource                   = &NetworkExplicitForwardProxyResource{}
	_ resource.ResourceWithValidateConfig = &NetworkExplicitForwardProxyResource{}
	_ resource.ResourceWithConfigure      = &NetworkExplicitForwardProxyResource{}
	_ resource.ResourceWithImportState    = &NetworkExplicitForwardProxyResource{}
)

// NetworkExplicitForwardProxyResource manages the Graph configuration.
type NetworkExplicitForwardProxyResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

func NewNetworkExplicitForwardProxyResource() resource.Resource {
	return &NetworkExplicitForwardProxyResource{
		ReadPermissions:  []string{"NetworkAccess.ReadWrite.All"},
		WritePermissions: []string{"NetworkAccess.ReadWrite.All"},
		ResourcePath:     "/networkAccess/explicitForwardProxyConfig",
	}
}

func (r *NetworkExplicitForwardProxyResource) Metadata(
	ctx context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = ResourceName
}

func (r *NetworkExplicitForwardProxyResource) Configure(
	ctx context.Context,
	req resource.ConfigureRequest,
	resp *resource.ConfigureResponse,
) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

func (r *NetworkExplicitForwardProxyResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	if req.ID != "explicitForwardProxyConfig" {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			"Use explicitForwardProxyConfig to import the tenant singleton.",
		)
		return
	}
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *NetworkExplicitForwardProxyResource) Schema(
	ctx context.Context,
	req resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages Microsoft Entra Global Secure Access Explicit Forward Proxy using `/networkAccess/explicitForwardProxyConfig`. Initial apply updates the existing tenant-wide settings; destroy removes Terraform state only and leaves the remote settings unchanged. Manage one instance per tenant. Updates preserve Private Access, mTLS and other unmanaged properties.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Fixed singleton identifier `explicitForwardProxyConfig`.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"internet_access": schema.SingleNestedAttribute{
				Required:            true,
				MarkdownDescription: "Internet Access configuration, corresponding to the API `internetAccess` object. All managed attributes are required; omission and null do not reset settings.",
				Attributes: map[string]schema.Attribute{
					"is_enabled": schema.BoolAttribute{
						Required:            true,
						MarkdownDescription: "API `isEnabled`. Enables Internet Access through Explicit Forward Proxy; changing this affects proxy clients.",
					},
					"is_source_ip_session_affinity_enabled": schema.BoolAttribute{
						Required:            true,
						MarkdownDescription: "API `isSourceIpSessionAffinityEnabled`. Set `true` when using session affinity; `false` requires `source_ip_session_affinity_options = \"none\"`.",
					},
					"source_ip_session_affinity_options": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "API `sourceIpSessionAffinityOptions`. Use `none` when affinity is disabled; otherwise use `useSessionId` (session ID affinity), `useHttpHeader` (HTTP header affinity), or `useSessionId,useHttpHeader` (both, in this canonical order). Session ID affinity requires service-hosted PAC files; HTTP header affinity requires the egress proxy to supply the documented header.",
						Validators: []validator.String{
							stringvalidator.OneOf(
								"none",
								"useSessionId",
								"useHttpHeader",
								"useSessionId,useHttpHeader",
							),
						},
					},
				},
			},
			"proxy_auto_configuration_file_url": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Read-only service-generated default PAC URL. This is not a configurable PAC reference URL or custom PAC body.",
			},
			"timeouts": commonschema.ResourceTimeouts(ctx),
		},
	}
}
