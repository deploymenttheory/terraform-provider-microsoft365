package utilityMicrosoft365Endpoints

import (
	"context"

	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	DataSourceName = "microsoft365_utility_microsoft_365_endpoints"
	ReadTimeout    = 180
)

var (
	// Basic datasource interface (Read operations)
	_ datasource.DataSource = &microsoft365EndpointsDataSource{}

	// Allows the datasource to be configured with the provider client
	_ datasource.DataSourceWithConfigure = &microsoft365EndpointsDataSource{}
)

func NewMicrosoft365EndpointsDataSource() datasource.DataSource {
	return &microsoft365EndpointsDataSource{}
}

type microsoft365EndpointsDataSource struct{}

func (d *microsoft365EndpointsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = DataSourceName
}

// Configure implements the DataSourceWithConfigure interface.
// For utility datasources that fetch data from external APIs (not Microsoft Graph),
// this method doesn't need to extract Microsoft Graph clients from ProviderData.
// However, it's required for interface compliance and maintains consistency across datasources.
func (d *microsoft365EndpointsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
}

func (d *microsoft365EndpointsDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieves Microsoft 365 network endpoints from the official Microsoft 365 IP Address and URL Web Service. " +
			"This datasource queries `https://endpoints.office.com` to get current IP addresses, URLs, and ports for Microsoft 365 services. " +
			"Useful for configuring firewalls, proxy servers, SD-WAN devices, and PAC files. " +
			"Data is filtered by cloud instance (Worldwide, US Government, China) and can be narrowed by service area and category.\n\n" +
			"See [Managing Microsoft 365 endpoints](https://learn.microsoft.com/en-us/microsoft-365/enterprise/managing-office-365-endpoints) for configuration guidance.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The ID of this datasource (format: `{instance}_{hash}`)",
			},
			"instance": schema.StringAttribute{
				Required: true,
				MarkdownDescription: "The Microsoft 365 cloud instance to query. Valid values:\n" +
					"  - `worldwide` - Worldwide commercial cloud (including US GCC)\n" +
					"  - `usgov-dod` - US Government DoD cloud\n" +
					"  - `usgov-gcchigh` - US Government GCC High cloud\n" +
					"  - `china` - Microsoft 365 operated by 21Vianet (China)\n\n" +
					"See [Microsoft 365 endpoints](https://learn.microsoft.com/en-us/microsoft-365/enterprise/microsoft-365-endpoints) for cloud details.",
				Validators: []validator.String{
					stringvalidator.OneOf("worldwide", "usgov-dod", "usgov-gcchigh", "china"),
				},
			},
			"service_areas": schema.SetAttribute{
				ElementType: types.StringType,
				Optional:    true,
				MarkdownDescription: "Filter endpoints by service area. Valid values:\n" +
					"  - `Exchange` - Exchange Online and Exchange Online Protection\n" +
					"  - `SharePoint` - SharePoint Online and OneDrive for Business\n" +
					"  - `Skype` - Skype for Business Online and Microsoft Teams\n" +
					"  - `Common` - Microsoft 365 Common (Microsoft Entra ID, Office in browser, etc.)\n\n" +
					"If omitted, returns endpoints for all service areas.",
				Validators: []validator.Set{
					setvalidator.ValueStringsAre(
						stringvalidator.OneOf("Exchange", "SharePoint", "Skype", "Common"),
					),
				},
			},
			"categories": schema.SetAttribute{
				ElementType: types.StringType,
				Optional:    true,
				MarkdownDescription: "Filter endpoints by network optimization category. Valid values:\n" +
					"  - `Optimize` - Required endpoints with highest traffic volume, latency sensitive. Direct routing recommended.\n" +
					"  - `Allow` - Required endpoints with lower traffic volume. Direct routing recommended, proxy acceptable.\n" +
					"  - `Default` - Optional endpoints or low-priority traffic. Can be routed through proxy.\n\n" +
					"If omitted, returns endpoints for all categories. " +
					"See [Microsoft 365 Network Connectivity Principles](https://learn.microsoft.com/en-us/microsoft-365/enterprise/microsoft-365-network-connectivity-principles) for category guidance.",
				Validators: []validator.Set{
					setvalidator.ValueStringsAre(
						stringvalidator.OneOf("Optimize", "Allow", "Default"),
					),
				},
			},
			"required_only": schema.BoolAttribute{
				Optional: true,
				MarkdownDescription: "If `true`, only returns endpoints marked as required by Microsoft. " +
					"Optional endpoints provide enhanced functionality but are not necessary for core service operation. " +
					"Defaults to `false` (returns all endpoints).",
			},
			"express_route": schema.BoolAttribute{
				Optional: true,
				MarkdownDescription: "If `true`, only returns endpoints that support Azure ExpressRoute for Microsoft 365. " +
					"Useful for organizations using ExpressRoute for optimized connectivity. " +
					"Defaults to `false` (returns all endpoints regardless of ExpressRoute support).",
			},
			"endpoints": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "List of Microsoft 365 endpoint sets matching the specified filters.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "Unique identifier for this endpoint set from Microsoft's service.",
						},
						"service_area": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The service area: `Exchange`, `SharePoint`, `Skype`, or `Common`.",
						},
						"service_area_display_name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Human-readable display name for the service area.",
						},
						"urls": schema.ListAttribute{
							ElementType:         types.StringType,
							Computed:            true,
							MarkdownDescription: "List of URL patterns (FQDNs) for this endpoint set. May include wildcards (e.g., `*.office.com`).",
						},
						"ips": schema.ListAttribute{
							ElementType:         types.StringType,
							Computed:            true,
							MarkdownDescription: "List of IP address ranges in CIDR notation (e.g., `40.96.0.0/13`). May be empty for URL-only endpoints.",
						},
						"tcp_ports": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "TCP ports used by this endpoint (comma-separated, e.g., `80,443` or ranges like `1024-65535`).",
						},
						"udp_ports": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "UDP ports used by this endpoint (comma-separated, e.g., `3478-3481`).",
						},
						"express_route": schema.BoolAttribute{
							Computed:            true,
							MarkdownDescription: "Whether this endpoint supports Azure ExpressRoute for Microsoft 365.",
						},
						"category": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Network routing category: `Optimize`, `Allow`, or `Default`. See [Network Connectivity Principles](https://learn.microsoft.com/en-us/microsoft-365/enterprise/microsoft-365-network-connectivity-principles).",
						},
						"required": schema.BoolAttribute{
							Computed:            true,
							MarkdownDescription: "Whether this endpoint is required for core Microsoft 365 functionality.",
						},
						"notes": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Additional notes about this endpoint from Microsoft, such as third-party services or optional features.",
						},
					},
				},
			},
			"timeouts": commonschema.DatasourceTimeouts(ctx),
		},
	}
}
