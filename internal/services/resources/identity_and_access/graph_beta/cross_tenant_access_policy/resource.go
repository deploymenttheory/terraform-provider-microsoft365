package graphBetaCrossTenantAccessPolicy

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName  = "microsoft365_graph_beta_identity_and_access_cross_tenant_access_policy"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180

	// singletonID is the static identifier used for this singleton resource in Terraform state.
	// The crossTenantAccessPolicy API has no server-assigned ID; it is addressed by a fixed path.
	singletonID = "crossTenantAccessPolicy"
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &CrossTenantAccessPolicyResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &CrossTenantAccessPolicyResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &CrossTenantAccessPolicyResource{}
)

func NewCrossTenantAccessPolicyResource() resource.Resource {
	return &CrossTenantAccessPolicyResource{
		ReadPermissions: []string{
			"Policy.Read.All",
		},
		WritePermissions: []string{
			"Policy.ReadWrite.CrossTenantAccess",
		},
	}
}

type CrossTenantAccessPolicyResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
}

// Metadata returns the resource type name.
func (r *CrossTenantAccessPolicyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

// Configure sets the client for the resource.
func (r *CrossTenantAccessPolicyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

// ImportState imports the resource state. Since this is a singleton resource, any import ID is accepted
// and normalised to the static identifier "crossTenantAccessPolicy".
func (r *CrossTenantAccessPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Schema defines the schema for the resource.
func (r *CrossTenantAccessPolicyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages the tenant-wide cross-tenant access policy in Microsoft Entra ID using the `/policies/crossTenantAccessPolicy` endpoint.\n\n" +
			"This is a **singleton resource** — one policy exists per tenant and cannot be created or deleted via the Microsoft Graph API. " +
			"The `create` operation uses a PATCH request to configure the policy. " +
			"On `destroy`, the resource can optionally restore the policy to its service defaults (empty `allowed_cloud_endpoints` and the default `display_name`) " +
			"by setting `restore_defaults_on_destroy = true`, or simply remove it from Terraform state while leaving the configuration in place (the default behaviour).\n\n" +
			"See the [Microsoft Graph API documentation](https://learn.microsoft.com/en-us/graph/api/resources/crosstenantaccesspolicy?view=graph-rest-beta) for details.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier for the cross-tenant access policy. This is a singleton resource; the value is always `crossTenantAccessPolicy`.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"display_name": schema.StringAttribute{
				MarkdownDescription: "The display name of the cross-tenant access policy. This is a read-only value set by the service and cannot be modified via Terraform.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"allowed_cloud_endpoints": schema.SetAttribute{
				MarkdownDescription: "Specifies which Microsoft clouds your organization would like to collaborate with. Use an empty set `[]` to disable all cross-cloud collaboration. " +
					"Supported values: `microsoftonline.com` (Azure commercial), `microsoftonline.us` (Azure US Government), `partner.microsoftonline.cn` (Azure China operated by 21Vianet).",
				ElementType: types.StringType,
				Required:    true,
				Validators: []validator.Set{
					setvalidator.ValueStringsAre(
						stringvalidator.OneOf(
							"microsoftonline.us",
							"partner.microsoftonline.cn",
						),
					),
				},
			},
			"restore_defaults_on_destroy": schema.BoolAttribute{
				MarkdownDescription: "Controls behaviour when this resource is destroyed. " +
					"When `true`, Terraform will issue a PATCH request to reset `allowed_cloud_endpoints` to an empty collection and `display_name` to `CrossTenantAccessPolicy`, " +
					"restoring the policy to its service defaults. " +
					"When `false` (the default), Terraform removes the resource from state only — the existing policy configuration is left unchanged in Microsoft Entra ID.",
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(false),
			},
			"timeouts": commonschema.ResourceTimeouts(ctx),
		},
	}
}
