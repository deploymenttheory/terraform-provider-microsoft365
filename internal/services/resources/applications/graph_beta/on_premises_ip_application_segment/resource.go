package graphBetaApplicationsOnPremisesIpApplicationSegment

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"regexp"
	"strings"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/identityschema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName  = "microsoft365_graph_beta_applications_on_premises_ip_application_segment"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &OnPremisesIpApplicationSegmentResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &OnPremisesIpApplicationSegmentResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &OnPremisesIpApplicationSegmentResource{}

	// Enables identity schema for list resource support
	_ resource.ResourceWithIdentity = &OnPremisesIpApplicationSegmentResource{}

	// Enables migration of protocol from string to set(string)
	_ resource.ResourceWithUpgradeState = &OnPremisesIpApplicationSegmentResource{}

	// Validates cross-attribute constraints between destination_type and destination_host
	_ resource.ResourceWithValidateConfig = &OnPremisesIpApplicationSegmentResource{}
)

func NewOnPremisesIpApplicationSegmentResource() resource.Resource {
	return &OnPremisesIpApplicationSegmentResource{
		ReadPermissions: []string{
			"Application.Read.All",
		},
		WritePermissions: []string{
			"Application.ReadWrite.All",
		},
		ResourcePath: "/applications/{application-id}/onPremisesPublishing/segmentsConfiguration/microsoft.graph.ipSegmentConfiguration/applicationSegments",
	}
}

type OnPremisesIpApplicationSegmentResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

// Metadata returns the resource type name.
func (r *OnPremisesIpApplicationSegmentResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

// Configure sets the client for the resource.
func (r *OnPremisesIpApplicationSegmentResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

// ImportState imports the resource state.
func (r *OnPremisesIpApplicationSegmentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.Split(req.ID, "/")
	if len(parts) != 2 {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			fmt.Sprintf("Expected import ID in format 'application_object_id/ip_application_segment_id', got: %s", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("application_object_id"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), parts[1])...)
}

// IdentitySchema defines the identity schema for this resource, used by list operations to uniquely identify instances
func (r *OnPremisesIpApplicationSegmentResource) IdentitySchema(ctx context.Context, req resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	resp.IdentitySchema = identityschema.Schema{
		Attributes: map[string]identityschema.Attribute{
			"id": identityschema.StringAttribute{
				RequiredForImport: true,
			},
		},
	}
}

// ValidateConfig validates that destination_host matches the format required by
// the Graph API for the configured destination_type. The application-scoped
// beta endpoint infers the destination type from the host format regardless of
// the requested destinationType: an "ipRange" request whose host is a CIDR, a
// single IP address, or an FQDN succeeds but is stored as ipRangeCidr, ip, or
// fqdn respectively, which would cause a permanent diff between configuration
// and state. A reversed range returns 500 InternalServerError and other forms
// return 400 DestinationHost_InvalidIP, so only "start..end" IPv4 ranges with
// start <= end are accepted at plan time.
func (r *OnPremisesIpApplicationSegmentResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var data OnPremisesIpApplicationSegmentResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if data.DestinationType.IsNull() || data.DestinationType.IsUnknown() ||
		data.DestinationHost.IsNull() || data.DestinationHost.IsUnknown() {
		return
	}

	if data.DestinationType.ValueString() != "ipRange" {
		return
	}

	if err := validateIpRangeHost(data.DestinationHost.ValueString()); err != nil {
		resp.Diagnostics.AddAttributeError(
			path.Root("destination_host"),
			"Invalid IP Range Destination Host",
			fmt.Sprintf("%s. When destination_type is \"ipRange\", destination_host must be an IPv4 start and end address separated by \"..\" with start <= end, for example \"192.168.1.1..192.168.1.10\". "+
				"Use destination_type \"ipRangeCidr\" for CIDR notation, \"ipAddress\" for a single IP address, or \"fqdn\" for hostnames.", err),
		)
	}
}

// validateIpRangeHost returns an error unless host is an IPv4 "start..end"
// range with start <= end.
func validateIpRangeHost(host string) error {
	parts := strings.Split(host, "..")
	if len(parts) != 2 {
		return fmt.Errorf("destination_host %q is not in \"start..end\" format", host)
	}

	start := net.ParseIP(parts[0])
	end := net.ParseIP(parts[1])
	if start == nil || start.To4() == nil {
		return fmt.Errorf("start address %q is not a valid IPv4 address", parts[0])
	}
	if end == nil || end.To4() == nil {
		return fmt.Errorf("end address %q is not a valid IPv4 address", parts[1])
	}

	if bytes.Compare(start.To4(), end.To4()) > 0 {
		return fmt.Errorf("start address %q is greater than end address %q", parts[0], parts[1])
	}

	return nil
}

// Schema defines the schema for the resource.
func (r *OnPremisesIpApplicationSegmentResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Version: 1,
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
			"application_object_id": schema.StringAttribute{
				MarkdownDescription: "The unique object identifier of the application.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"destination_host": schema.StringAttribute{
				MarkdownDescription: "Either the IP address, IP range, or FQDN of the application segment, with or without wildcards. " +
					"For `destination_type = \"ipRange\"`, use the IPv4 start and end addresses separated by `..` with start <= end, for example `192.168.1.1..192.168.1.10`; " +
					"this format is validated at plan time because the Graph API infers the stored destination type from the host format (for example, a CIDR host such as `192.168.1.0/24` is stored as `ipRangeCidr` even when `ipRange` is requested), which would otherwise cause a permanent diff.",
				Required: true,
			},
			"destination_type": schema.StringAttribute{
				MarkdownDescription: "The type of destination for the application segment. " +
					"The supported values are: `ipAddress`, `ipRange`, `ipRangeCidr`, and `fqdn`. " +
					"Microsoft Learn also lists `dnsSuffix` for `ipApplicationSegment`, but this resource does not support it because the application-scoped Graph endpoint discards `ports` and `protocol` for `dnsSuffix` segments, which this resource requires.",
				Required: true,
				Validators: []validator.String{
					stringvalidator.OneOf("ipAddress", "ipRange", "ipRangeCidr", "fqdn"),
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
			"protocol": schema.SetAttribute{
				MarkdownDescription: "The protocols of the network traffic acquired for the application segment. " +
					"Supported values are `tcp` and `udp`; specify both values to enable both protocols.",
				Required:    true,
				ElementType: types.StringType,
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
					setvalidator.SizeAtMost(2),
					setvalidator.ValueStringsAre(
						stringvalidator.OneOf("tcp", "udp"),
					),
				},
			},
			"timeouts": commonschema.ResourceTimeouts(ctx),
		},
	}
}
