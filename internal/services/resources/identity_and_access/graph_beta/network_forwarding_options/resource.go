package graphBetaNetworkForwardingOptions

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
	ResourceName  = "microsoft365_graph_beta_identity_and_access_network_forwarding_options"
	CreateTimeout = 180
	ReadTimeout   = 180
	UpdateTimeout = 180
	DeleteTimeout = 180
	singletonID   = "forwardingOptions"
)

var (
	_ resource.Resource                = &NetworkForwardingOptionsResource{}
	_ resource.ResourceWithConfigure   = &NetworkForwardingOptionsResource{}
	_ resource.ResourceWithImportState = &NetworkForwardingOptionsResource{}
)

// NetworkForwardingOptionsResource manages the existing tenant forwarding options.
type NetworkForwardingOptionsResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

func NewNetworkForwardingOptionsResource() resource.Resource {
	return &NetworkForwardingOptionsResource{
		ReadPermissions:  []string{"NetworkAccess.ReadWrite.All"},
		WritePermissions: []string{"NetworkAccess.ReadWrite.All"},
		ResourcePath:     "/networkAccess/settings/forwardingOptions",
	}
}

func (r *NetworkForwardingOptionsResource) Metadata(
	ctx context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = ResourceName
}

func (r *NetworkForwardingOptionsResource) Configure(
	ctx context.Context,
	req resource.ConfigureRequest,
	resp *resource.ConfigureResponse,
) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

func (r *NetworkForwardingOptionsResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	if req.ID != singletonID {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			"Use forwardingOptions to import the tenant-wide singleton.",
		)
		return
	}
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *NetworkForwardingOptionsResource) Schema(
	ctx context.Context,
	req resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages Microsoft Entra Global Secure Access DNS forwarding using Microsoft Graph beta `/networkAccess/settings/forwardingOptions`. Initial apply updates the existing tenant-wide settings; destroy removes Terraform state only and leaves the remote settings unchanged. Manage one instance per tenant. Updates preserve unmanaged properties.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The fixed singleton identifier `forwardingOptions`.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"skip_dns_lookup_state": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Whether to skip service-side DNS lookup and forward Microsoft 365 traffic using the client-resolved destination IP. Values: `enabled`, `disabled`. Explicit configuration is required; omission, null and empty strings are rejected rather than interpreted as a reset.",
				Validators: []validator.String{
					stringvalidator.OneOf("enabled", "disabled"),
				},
			},
			"timeouts": commonschema.ResourceTimeouts(ctx),
		},
	}
}
