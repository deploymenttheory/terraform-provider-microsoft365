package graphBetaNetworkCloudFirewallPolicy

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	planmodifiers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/plan_modifiers"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/identityschema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	kiotahttp "github.com/microsoft/kiota-http-go"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName  = "microsoft365_graph_beta_identity_and_access_network_cloud_firewall_policy"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	_ resource.Resource                = &NetworkCloudFirewallPolicyResource{}
	_ resource.ResourceWithConfigure   = &NetworkCloudFirewallPolicyResource{}
	_ resource.ResourceWithImportState = &NetworkCloudFirewallPolicyResource{}
	_ resource.ResourceWithIdentity    = &NetworkCloudFirewallPolicyResource{}
)

func NewNetworkCloudFirewallPolicyResource() resource.Resource {
	return &NetworkCloudFirewallPolicyResource{
		ReadPermissions:  []string{"NetworkAccess.Read.All"},
		WritePermissions: []string{"NetworkAccess.ReadWrite.All"},
		ResourcePath:     "/networkaccess/cloudFirewallPolicies",
	}
}

type NetworkCloudFirewallPolicyResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	retryOptions     kiotahttp.RetryHandlerOptions
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

func (r *NetworkCloudFirewallPolicyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

func (r *NetworkCloudFirewallPolicyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
	r.retryOptions = client.GraphRetryOptionsForResource(req.ProviderData)
}

func (r *NetworkCloudFirewallPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughWithIdentity(ctx, path.Root("id"), path.Root("id"), req, resp)
}

func (r *NetworkCloudFirewallPolicyResource) IdentitySchema(ctx context.Context, req resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	resp.IdentitySchema = identityschema.Schema{Attributes: map[string]identityschema.Attribute{
		"id": identityschema.StringAttribute{RequiredForImport: true},
	}}
}

func (r *NetworkCloudFirewallPolicyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages Microsoft Entra Global Secure Access cloud firewall policies using the portal-backed Microsoft Graph beta `/networkaccess/cloudFirewallPolicies` endpoint.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier for the cloud firewall policy.",
				Computed:            true,
				PlanModifiers:       []planmodifier.String{planmodifiers.UseStateForUnknownString()},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the cloud firewall policy.",
				Required:            true,
			},
			"description": schema.StringAttribute{MarkdownDescription: "The policy description. Omitting it clears the description; an empty string is preserved.", Optional: true},
			"default_action": schema.StringAttribute{
				MarkdownDescription: "The action when no rule matches. Only `allow` is supported. Terraform defaults to `allow` and explicitly sends it; the API requires this setting.",
				Optional:            true, Computed: true, Default: stringdefault.StaticString("allow"),
				Validators: []validator.String{stringvalidator.OneOf("allow")},
			},
			"version":                 schema.StringAttribute{MarkdownDescription: "The service-generated policy version.", Computed: true},
			"last_modified_date_time": schema.StringAttribute{MarkdownDescription: "The last modification time, including changes to child rules.", Computed: true},
			"timeouts":                commonschema.ResourceTimeouts(ctx),
		},
	}
}
