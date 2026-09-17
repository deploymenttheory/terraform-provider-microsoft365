package graphBetaNetworkCrossTenantAccessSettings

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName  = "microsoft365_graph_beta_identity_and_access_network_cross_tenant_access_settings"
	CreateTimeout = 180
	ReadTimeout   = 180
	UpdateTimeout = 180
	DeleteTimeout = 180
	singletonID   = "crossTenantAccess"
)

var (
	_ resource.Resource                = &NetworkCrossTenantAccessSettingsResource{}
	_ resource.ResourceWithConfigure   = &NetworkCrossTenantAccessSettingsResource{}
	_ resource.ResourceWithImportState = &NetworkCrossTenantAccessSettingsResource{}
)

type NetworkCrossTenantAccessSettingsResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

func NewNetworkCrossTenantAccessSettingsResource() resource.Resource {
	return &NetworkCrossTenantAccessSettingsResource{
		ReadPermissions:  []string{"NetworkAccess.Read.All"},
		WritePermissions: []string{"NetworkAccess.ReadWrite.All"},
		ResourcePath:     "/networkAccess/settings/crossTenantAccess",
	}
}

func (r *NetworkCrossTenantAccessSettingsResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

func (r *NetworkCrossTenantAccessSettingsResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

func (r *NetworkCrossTenantAccessSettingsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if req.ID != singletonID {
		resp.Diagnostics.AddError("Invalid import ID", "Import this singleton using the literal ID crossTenantAccess.")
		return
	}
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *NetworkCrossTenantAccessSettingsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages Microsoft Entra Global Secure Access universal tenant restrictions using the Microsoft Graph beta `/networkAccess/settings/crossTenantAccess` endpoint. This is separate from `/policies/crossTenantAccessPolicy`.\n\n" +
			"This is a tenant-wide singleton with documented GET and PATCH operations only. Initial apply updates the existing settings using PATCH. Destroy removes the resource from Terraform state only; it does not disable, delete, or reset the remote settings. Manage only one instance per tenant across all Terraform states.\n\n" +
			"Enabling packet tagging applies configured tenant restrictions to affected traffic and can block access to external tenants. Configure and review the tenant restrictions policy before enabling this setting.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The fixed Terraform identifier `crossTenantAccess`, representing the singleton endpoint.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"network_packet_tagging_status": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Whether network packets are tagged to enforce tenant restrictions. Allowed values: `enabled`, `disabled`. To disable tagging, apply `disabled` before removing this resource; destroy leaves this value unchanged.",
				Validators:          []validator.String{stringvalidator.OneOf("enabled", "disabled")},
			},
			"timeouts": commonschema.ResourceTimeouts(ctx),
		},
	}
}
