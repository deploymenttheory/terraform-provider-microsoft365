package graphBetaNamedLocation

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
)

const (
	ResourceName  = "graph_beta_identity_and_access_named_location"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &NamedLocationResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &NamedLocationResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &NamedLocationResource{}
)

func NewNamedLocationResource() resource.Resource {
	return &NamedLocationResource{
		ReadPermissions: []string{
			"Policy.Read.All",
			"Policy.Read.ConditionalAccess",
		},
		WritePermissions: []string{
			"Policy.ReadWrite.ConditionalAccess",
		},
		ResourcePath: "/identity/conditionalAccess/namedLocations",
	}
}

type NamedLocationResource struct {
	httpClient       *client.AuthenticatedHTTPClient
	ProviderTypeName string
	TypeName         string
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

// Metadata returns the resource type name.
func (r *NamedLocationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	r.ProviderTypeName = req.ProviderTypeName
	r.TypeName = ResourceName
	resp.TypeName = r.FullTypeName()
}

// FullTypeName returns the full resource type name in the format "providername_resourcename".
func (r *NamedLocationResource) FullTypeName() string {
	return r.ProviderTypeName + "_" + r.TypeName
}

// Configure sets the client for the resource.
func (r *NamedLocationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.httpClient = client.SetGraphBetaHTTPClientForResource(ctx, req, resp, constants.PROVIDER_NAME+"_"+ResourceName)
}

// ImportState imports the resource state.
func (r *NamedLocationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Schema defines the schema for the resource.
func (r *NamedLocationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages Microsoft 365 Named Locations using the `/identity/conditionalAccess/namedLocations` endpoint. Named Locations define network locations that can be used in Conditional Access policies. Supports both IP-based and country-based named locations.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "String (identifier)",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(constants.GuidRegex),
						"must be a valid GUID",
					),
				},
			},
			"display_name": schema.StringAttribute{
				MarkdownDescription: "The display name for the Named Location.",
				Required:            true,
			},
			"created_date_time": schema.StringAttribute{
				MarkdownDescription: "The creation date and time of the named location.",
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(constants.TimeFormatRFC3339Regex),
						"must be a valid RFC3339 date-time string",
					),
				},
			},
			"modified_date_time": schema.StringAttribute{
				MarkdownDescription: "The last modified date and time of the named location.",
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(constants.TimeFormatRFC3339Regex),
						"must be a valid RFC3339 date-time string",
					),
				},
			},

			// IP Named Location attributes
			"is_trusted": schema.BoolAttribute{
				MarkdownDescription: "Indicates whether the IP named location is trusted. Only applies to IP named locations.",
				Optional:            true,
				Computed:            true,
			},
			"ipv4_ranges": schema.SetAttribute{
				MarkdownDescription: "Set of IPv4 CIDR ranges that define this IP named location. Each range should be specified in CIDR notation (e.g., '192.168.1.0/24'). Used for IP named locations only.",
				ElementType:         types.StringType,
				Optional:            true,
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
					setvalidator.ValueStringsAre(
						stringvalidator.RegexMatches(
							regexp.MustCompile(constants.IPv4CIDRRegex),
							"must be a valid IPv4 CIDR range (e.g., '192.168.1.0/24')",
						),
					),
				},
			},
			"ipv6_ranges": schema.SetAttribute{
				MarkdownDescription: "Set of IPv6 CIDR ranges that define this IP named location. Each range should be specified in CIDR notation (e.g., '2001:db8::/32'). Used for IP named locations only.",
				ElementType:         types.StringType,
				Optional:            true,
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
					setvalidator.ValueStringsAre(
						stringvalidator.RegexMatches(
							regexp.MustCompile(constants.IPv6CIDRRegex),
							"must be a valid IPv6 CIDR range (e.g., '2001:db8::/32')",
						),
					),
				},
			},

			// Country Named Location attributes
			"country_lookup_method": schema.StringAttribute{
				MarkdownDescription: "Provides the method used to decide which country the user is located in. Possible values are `clientIpAddress` and `authenticatorAppGps`. Used for country named locations only.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("clientIpAddress", "authenticatorAppGps"),
				},
			},
			"include_unknown_countries_and_regions": schema.BoolAttribute{
				MarkdownDescription: "True if IP addresses that don't map to a country or region should be included in the named location. Used for country named locations only.",
				Optional:            true,
			},
			"countries_and_regions": schema.SetAttribute{
				MarkdownDescription: "Set of countries and/or regions in two-letter format specified by ISO 3166-2 (e.g., 'US', 'GB', 'CA'). Used for country named locations only.",
				ElementType:         types.StringType,
				Optional:            true,
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
					setvalidator.ValueStringsAre(
						stringvalidator.OneOf(
							"AF", "AX", "AL", "DZ", "AS", "AD", "AO", "AI", "AQ", "AG", "AR", "AM", "AW", "AU",
							"AT", "AZ", "BS", "BH", "BD", "BB", "BY", "BE", "BZ", "BJ", "BM", "BT", "BO", "BQ",
							"BA", "BW", "BV", "BR", "IO", "BN", "BG", "BF", "BI", "CV", "KH", "CM", "CA", "KY",
							"CF", "TD", "CL", "CN", "CX", "CC", "CO", "KM", "CK", "CR", "CI", "HR", "CU", "CW",
							"CY", "CZ", "CD", "DK", "DJ", "DM", "DO", "EC", "EG", "SV", "GQ", "ER", "EE", "SZ",
							"ET", "FK", "FO", "FJ", "FI", "FR", "GF", "PF", "TF", "GA", "GM", "GE", "DE", "GH",
							"GI", "GR", "GL", "GD", "GP", "GU", "GT", "GG", "GN", "GW", "GY", "HT", "HM", "VA",
							"HN", "HK", "HU", "IS", "IN", "ID", "IR", "IQ", "IE", "IM", "IL", "IT", "JM", "JP",
							"JE", "JO", "KZ", "KE", "KI", "KR", "XK", "KW", "KG", "LA", "LV", "LB", "LS", "LR",
							"LY", "LI", "LT", "LU", "MO", "MG", "MW", "MY", "MV", "ML", "MT", "MH", "MQ", "MR",
							"MU", "YT", "MX", "FM", "MD", "MC", "MN", "ME", "MS", "MA", "MZ", "MM", "NA", "NR",
							"NP", "NL", "NC", "NZ", "NI", "NE", "NG", "NU", "NF", "KP", "MK", "MP", "NO", "OM",
							"PK", "PW", "PS", "PA", "PG", "PY", "PE", "PH", "PN", "PL", "PT", "PR", "QA", "CG",
							"RE", "RO", "RU", "RW", "BL", "SH", "KN", "LC", "MF", "PM", "VC", "WS", "SM", "ST",
							"SA", "SN", "RS", "SC", "SL", "SG", "SX", "SK", "SI", "SB", "SO", "ZA", "GS", "SS",
							"ES", "LK", "SD", "SR", "SJ", "SE", "CH", "SY", "TW", "TJ", "TZ", "TH", "TL", "TG",
							"TK", "TO", "TT", "TN", "TR", "TM", "TC", "TV", "UG", "UA", "AE", "GB", "US", "UY",
							"UM", "UZ", "VU", "VE", "VN", "VG", "VI", "WF", "EH", "YE", "ZM", "ZW",
						),
					),
				},
			},

			"timeouts": commonschema.Timeouts(ctx),
		},
	}
}
