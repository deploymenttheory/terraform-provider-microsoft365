package graphBetaApplicationsIpApplicationSegment

import (
	"context"
	"regexp"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
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
)

const (
	ResourceName  = "microsoft365_graph_beta_applications_ip_application_segment"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &IpApplicationSegmentResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &IpApplicationSegmentResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &IpApplicationSegmentResource{}
)

func NewIpApplicationSegmentResource() resource.Resource {
	return &IpApplicationSegmentResource{
		ReadPermissions: []string{
			"Application.Read.All",
		},
		WritePermissions: []string{
			"Application.ReadWrite.All",
		},
		ResourcePath: "/applications/{application-id}/onPremisesPublishing/segmentsConfiguration/microsoft.graph.ipSegmentConfiguration/applicationSegments",
	}
}

type IpApplicationSegmentResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

// Metadata returns the resource type name.
func (r *IpApplicationSegmentResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

// Configure sets the client for the resource.
func (r *IpApplicationSegmentResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

// ImportState imports the resource state.
func (r *IpApplicationSegmentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Schema defines the schema for the resource.
func (r *IpApplicationSegmentResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages an IP application segment for on-premises publishing. " +
			"IP application segments define the destination hosts, ports, and protocols for applications published through Azure AD Application Proxy.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier of the application segment.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"application_id": schema.StringAttribute{
				MarkdownDescription: "The unique object identifier of the application.",
				Required:            true,
			},
			"destination_host": schema.StringAttribute{
				MarkdownDescription: "Either the IP address, IP range, or FQDN of the application segment, with or without wildcards.",
				Required:            true,
			},
			"destination_type": schema.StringAttribute{
				MarkdownDescription: "The type of destination for the application segment." +
					"The possible values are: `ipAddress`, `ipRange`, `ipRangeCidr`, `fqdn`, `dnsSuffix`, `unknownFutureValue`.",
				Required: true,
				Validators: []validator.String{
					stringvalidator.OneOf("ipAddress", "ipRange", "ipRangeCidr", "fqdn", "dnsSuffix", "unknownFutureValue"),
				},
			},
			"ports": schema.SetAttribute{
				MarkdownDescription: "List of ports supported for the application segment.",
				Required:            true,
				ElementType:         types.StringType,
				Validators: []validator.Set{
					setvalidator.ValueStringsAre(
						stringvalidator.RegexMatches(regexp.MustCompile(constants.PortRangeRegex), "Each port defined in the set must be a valid format (xxxx-xxxx) e.g 80-80, 443-443, 8080-8080, 8443-8443"),
					),
				},
			},
			"protocol": schema.StringAttribute{
				MarkdownDescription: "Indicates the protocol of the network traffic acquired for the application segment." +
					"The possible values are: `tcp`, `udp`, `unknownFutureValue`.",
				Required: true,
				Validators: []validator.String{
					stringvalidator.OneOf("tcp", "udp"),
				},
			},
			"timeouts": commonschema.Timeouts(ctx),
		},
	}
}
