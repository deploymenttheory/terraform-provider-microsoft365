package graphBetaIdentityAndAccessNetworkRemoteNetwork

import (
	"context"
	"regexp"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
)

const (
	ResourceName  = "microsoft365_graph_beta_identity_and_access_network_remote_network"
	CreateTimeout = 180
	ReadTimeout   = 180
	UpdateTimeout = 180
	DeleteTimeout = 180
)

// regions is the microsoft.graph.networkaccess.region enum, excluding unknownFutureValue.
var regions = []string{
	"eastUS", "eastUS2", "westUS", "westUS2", "westUS3", "centralUS", "northCentralUS", "southCentralUS",
	"northEurope", "westEurope", "franceCentral", "germanyWestCentral", "switzerlandNorth", "ukSouth",
	"canadaEast", "canadaCentral", "southAfricaWest", "southAfricaNorth", "uaeNorth", "australiaEast",
	"westCentralUS", "centralIndia", "southEastAsia", "swedenCentral", "southIndia", "australiaSouthEast",
	"koreaCentral", "polandCentral", "brazilSouth", "japanEast", "japanWest", "koreaSouth", "italyNorth",
	"franceSouth", "israelCentral", "taiwanNorth", "mexicoCentral", "spainCentral", "jioIndiaCentral",
	"brazilSouthEast",
}

var (
	_ resource.Resource                = &NetworkRemoteNetworkResource{}
	_ resource.ResourceWithConfigure   = &NetworkRemoteNetworkResource{}
	_ resource.ResourceWithImportState = &NetworkRemoteNetworkResource{}
)

// NetworkRemoteNetworkResource manages a Global Secure Access remote network.
type NetworkRemoteNetworkResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
}

// https://learn.microsoft.com/en-us/graph/api/networkaccess-connectivity-post-remotenetworks?view=graph-rest-beta#permissions
func NewNetworkRemoteNetworkResource() resource.Resource {
	return &NetworkRemoteNetworkResource{
		ReadPermissions:  []string{"NetworkAccess.Read.All"},
		WritePermissions: []string{"NetworkAccess.ReadWrite.All"},
	}
}

func (r *NetworkRemoteNetworkResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

func (r *NetworkRemoteNetworkResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

func (r *NetworkRemoteNetworkResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if _, err := uuid.Parse(req.ID); err != nil {
		resp.Diagnostics.AddError("Invalid import ID", "Use the remote network ID (a GUID) returned by Microsoft Graph.")
		return
	}
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *NetworkRemoteNetworkResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a Microsoft Entra Global Secure Access remote network and its forwarding profile associations using the Microsoft Graph beta " +
			"`/networkaccess/connectivity/remoteNetworks` endpoint. Device links and connectivity configuration are not managed by this resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The unique identifier of the remote network.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of the remote network.",
				Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
			},
			"region": schema.StringAttribute{
				Required: true,
				MarkdownDescription: "The Azure region of the remote network. Changing this forces a new remote network to be created. " +
					"Possible values are: `" + strings.Join(regions, "`, `") + "`.",
				Validators:    []validator.String{stringvalidator.OneOf(regions...)},
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"forwarding_profile_ids": schema.SetAttribute{
				Required:    true,
				ElementType: types.StringType,
				MarkdownDescription: "IDs of the forwarding profiles associated with the remote network. At least one is required. " +
					"Only the Microsoft 365 and Internet Access default forwarding profiles can be assigned, and they must be enabled.",
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
					setvalidator.ValueStringsAre(
						stringvalidator.RegexMatches(
							regexp.MustCompile(constants.GuidRegex),
							"must be a valid GUID in the format '00000000-0000-0000-0000-000000000000'",
						),
					),
				},
			},
			"version": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The version of the remote network returned by Microsoft Graph.",
			},
			"timeouts": commonschema.ResourceTimeouts(ctx),
		},
	}
}
