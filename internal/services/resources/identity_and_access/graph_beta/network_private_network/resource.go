package graphBetaNetworkPrivateNetwork

import (
	"context"
	"regexp"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	planmodifiers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/plan_modifiers"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/identityschema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName  = "microsoft365_graph_beta_identity_and_access_network_private_network"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	_ resource.Resource                = &NetworkPrivateNetworkResource{}
	_ resource.ResourceWithConfigure   = &NetworkPrivateNetworkResource{}
	_ resource.ResourceWithImportState = &NetworkPrivateNetworkResource{}
	_ resource.ResourceWithIdentity    = &NetworkPrivateNetworkResource{}
)

func NewNetworkPrivateNetworkResource() resource.Resource {
	return &NetworkPrivateNetworkResource{
		ReadPermissions: []string{
			"NetworkAccess.Read.All",
		},
		WritePermissions: []string{
			"NetworkAccess.ReadWrite.All",
		},
		ResourcePath: "/networkaccess/privateNetworks",
	}
}

type NetworkPrivateNetworkResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

func (r *NetworkPrivateNetworkResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

func (r *NetworkPrivateNetworkResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

func (r *NetworkPrivateNetworkResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *NetworkPrivateNetworkResource) IdentitySchema(ctx context.Context, req resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	resp.IdentitySchema = identityschema.Schema{
		Attributes: map[string]identityschema.Attribute{
			"id": identityschema.StringAttribute{
				RequiredForImport: true,
			},
		},
	}
}

func (r *NetworkPrivateNetworkResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages Microsoft Entra Global Secure Access private networks using the Microsoft Graph beta `/networkaccess/privateNetworks` endpoint observed from the Entra portal.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier for the private network.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the private network.",
				Required:            true,
			},
			"app_ids": schema.SetAttribute{
				MarkdownDescription: "A set of Microsoft Entra application client IDs attached as target resources for the private network.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
				Default:             setdefault.StaticValue(types.SetValueMust(types.StringType, []attr.Value{})),
				Validators: []validator.Set{
					setvalidator.ValueStringsAre(
						stringvalidator.RegexMatches(
							regexp.MustCompile(constants.GuidRegex),
							"must be a valid UUID",
						),
					),
				},
			},
			"dns_resolution_identification": schema.SingleNestedAttribute{
				MarkdownDescription: "DNS resolution identification settings used to identify this private network.",
				Required:            true,
				Attributes: map[string]schema.Attribute{
					"dns_servers": schema.SetAttribute{
						MarkdownDescription: "DNS server IP addresses used to resolve the fully qualified domain name.",
						Required:            true,
						ElementType:         types.StringType,
						Validators: []validator.Set{
							setvalidator.SizeAtLeast(1),
						},
					},
					"fqdn_to_resolve": schema.StringAttribute{
						MarkdownDescription: "The fully qualified domain name that should resolve to the expected private IP address, CIDR, or IP range values.",
						Required:            true,
					},
					"expected_ip_resolutions": schema.SetNestedAttribute{
						MarkdownDescription: "Expected IP address, CIDR subnet, or IP range values returned when resolving the FQDN. The Entra portal presents these as comma-separated values, while Graph stores them as a collection.",
						Required:            true,
						Validators: []validator.Set{
							setvalidator.SizeAtLeast(1),
						},
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"type": schema.StringAttribute{
									MarkdownDescription: "The expected IP resolution type. Possible values are `ip_address`, `ip_subnet`, and `ip_range`.",
									Required:            true,
									Validators: []validator.String{
										stringvalidator.OneOf(expectedIPResolutionTypeIPAddress, expectedIPResolutionTypeIPSubnet, expectedIPResolutionTypeIPRange),
									},
								},
								"value": schema.StringAttribute{
									MarkdownDescription: "The IP address value when `type` is `ip_address`, or CIDR value when `type` is `ip_subnet`.",
									Optional:            true,
								},
								"begin_address": schema.StringAttribute{
									MarkdownDescription: "The beginning IP address when `type` is `ip_range`.",
									Optional:            true,
								},
								"end_address": schema.StringAttribute{
									MarkdownDescription: "The ending IP address when `type` is `ip_range`.",
									Optional:            true,
								},
							},
						},
					},
				},
			},
			"timeouts": commonschema.ResourceTimeouts(ctx),
		},
	}
}
