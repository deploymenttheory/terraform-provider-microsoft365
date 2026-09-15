package graphBetaNetworkCustomBlockPage

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
)

const (
	ResourceName  = "microsoft365_graph_beta_identity_and_access_network_custom_block_page"
	singletonID   = "customBlockPage"
	CreateTimeout = 180
	ReadTimeout   = 180
	UpdateTimeout = 180
)

var (
	_ resource.Resource                = &NetworkCustomBlockPageResource{}
	_ resource.ResourceWithConfigure   = &NetworkCustomBlockPageResource{}
	_ resource.ResourceWithImportState = &NetworkCustomBlockPageResource{}
)

type NetworkCustomBlockPageResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
}

func NewNetworkCustomBlockPageResource() resource.Resource {
	return &NetworkCustomBlockPageResource{
		ReadPermissions:  []string{"NetworkAccess.Read.All"},
		WritePermissions: []string{"NetworkAccess.ReadWrite.All"},
	}
}

func (r *NetworkCustomBlockPageResource) Metadata(
	_ context.Context,
	_ resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = ResourceName
}

func (r *NetworkCustomBlockPageResource) Configure(
	ctx context.Context,
	req resource.ConfigureRequest,
	resp *resource.ConfigureResponse,
) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

func (r *NetworkCustomBlockPageResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	if req.ID != singletonID {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			"This tenant-wide singleton must be imported using the ID customBlockPage.",
		)
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), singletonID)...)
}

func (r *NetworkCustomBlockPageResource) Schema(
	ctx context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages the Microsoft Entra Global Secure Access custom block page using `/networkAccess/settings/customBlockPage`. " +
			"This is a tenant-wide singleton: initial apply and updates use PATCH to configure existing settings. " +
			"Destroy removes the resource from Terraform state only, leaving the remote settings unchanged. Manage only one instance per tenant.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The fixed Terraform identifier `customBlockPage`. The observed API response does not include an ID.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"state": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Whether the custom message is enabled. Values: `enabled`, `disabled`, `unknownFutureValue`. Use enabled or disabled for normal configuration.",
				Validators: []validator.String{
					stringvalidator.OneOf("enabled", "disabled", "unknownFutureValue"),
				},
			},
			"configuration": schema.SingleNestedAttribute{
				Optional:            true,
				Computed:            true,
				PlanModifiers:       []planmodifier.Object{objectplanmodifier.UseStateForUnknown()},
				MarkdownDescription: "Custom block message configuration. Omission preserves the existing message. The current API uses a Markdown block message; its type discriminator is set internally.",
				Attributes: map[string]schema.Attribute{
					"body": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "Custom notification body. Must contain non-whitespace text and be 1–1024 UTF-16 code units long. Most characters count as one unit; supplementary characters such as emoji count as two. Supports plain text and Markdown links. Disabling the page preserves its configured body.",
						Validators:          []validator.String{blockMessageValidator{}},
					},
				},
			},
			"timeouts": commonschema.ResourceTimeouts(ctx),
		},
	}
}
