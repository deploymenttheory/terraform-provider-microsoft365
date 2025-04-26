package graphBetaWindowsUpdateCatalogItems

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/crud"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

const (
	DataSourceName = "graph_beta_windows_update_catalog_items"
	ReadTimeout    = 180
)

var (
	_ datasource.DataSource              = &WindowsUpdateCatalogItemsDataSource{}
	_ datasource.DataSourceWithConfigure = &WindowsUpdateCatalogItemsDataSource{}
)

// NewWindowsUpdateCatalogItemsDataSource creates a new data source for Windows Update Catalog Items
func NewWindowsUpdateCatalogItemsDataSource() datasource.DataSource {
	return &WindowsUpdateCatalogItemsDataSource{
		ReadPermissions: []string{
			"DeviceManagementConfiguration.Read.All",
		},
	}
}

// WindowsUpdateCatalogItemsDataSource defines the data source implementation
type WindowsUpdateCatalogItemsDataSource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ProviderTypeName string
	TypeName         string
	ReadPermissions  []string
}

// WindowsUpdateCatalogItemModel represents a single catalog item
type WindowsUpdateCatalogItemModel struct {
	ID               types.String `tfsdk:"id"`
	DisplayName      types.String `tfsdk:"display_name"`
	ReleaseDateTime  types.String `tfsdk:"release_date_time"`
	EndOfSupportDate types.String `tfsdk:"end_of_support_date"`
}

// WindowsUpdateCatalogItemsDataSourceModel defines the data source model
type WindowsUpdateCatalogItemsDataSourceModel struct {
	FilterType       types.String                    `tfsdk:"filter_type"`  // Required field to specify how to filter
	FilterValue      types.String                    `tfsdk:"filter_value"` // Value to filter by (not used for "all")
	Items            []WindowsUpdateCatalogItemModel `tfsdk:"items"`        // List of catalog items that match the filters
	Timeouts         types.Object                    `tfsdk:"timeouts"`
}

// Metadata returns the data source type name
func (d *WindowsUpdateCatalogItemsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + DataSourceName
}

// Configure configures the data source with the provider client
func (d *WindowsUpdateCatalogItemsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = common.SetGraphBetaClientForDataSource(ctx, req, resp, d.ProviderTypeName)
}

// Schema defines the schema for the data source
func (d *WindowsUpdateCatalogItemsDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieves Windows Update Catalog Items from Microsoft Intune with explicit filtering options.",
		Attributes: map[string]schema.Attribute{
			"filter_type": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Type of filter to apply. Valid values are: `all`, `id`, `display_name`, `release_date_time`, `end_of_support_date`.",
				Validators: []validator.String{
					stringvalidator.OneOf("all", "id", "display_name", "release_date_time", "end_of_support_date"),
				},
			},
			"filter_value": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Value to filter by. Not required when filter_type is 'all'. For date filters, use RFC3339 format (e.g., '2023-01-01T00:00:00Z').",
			},
			"items": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "The list of Windows Update Catalog Items that match the filter criteria.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The ID of the catalog item.",
						},
						"display_name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The display name of the catalog item.",
						},
						"release_date_time": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The release date time of the catalog item.",
						},
						"end_of_support_date": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The end of support date of the catalog item.",
						},
					},
				},
			},
			"timeouts": common.TimeoutSchema(),
		},
	}
}

// MapCatalogItemToModel maps a Windows Update Catalog Item to a model
func MapCatalogItemToModel(item graphmodels.WindowsUpdateCatalogItemable) WindowsUpdateCatalogItemModel {
	model := WindowsUpdateCatalogItemModel{
		ID:               types.StringPointerValue(item.GetId()),
		DisplayName:      types.StringPointerValue(item.GetDisplayName()),
	}

	if releaseDateTime := item.GetReleaseDateTime(); releaseDateTime != nil {
		model.ReleaseDateTime = types.StringValue(releaseDateTime.Format(time.RFC3339))
	}

	if endOfSupportDate := item.GetEndOfSupportDate(); endOfSupportDate != nil {
		model.EndOfSupportDate = types.StringValue(endOfSupportDate.Format(time.RFC3339))
	}

	return model
}

// Read handles the Read operation for Windows Update Catalog Items data source.
func (d *WindowsUpdateCatalogItemsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var object WindowsUpdateCatalogItemsDataSourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting Read method for: %s_%s", d.ProviderTypeName, d.TypeName))

	resp.Diagnostics.Append(req.Config.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	filterType := object.FilterType.ValueString()
	tflog.Debug(ctx, fmt.Sprintf("Reading %s_%s with filter_type: %s", d.ProviderTypeName, d.TypeName, filterType))

	// Validate filter_value is provided when filter_type is not "all"
	if filterType != "all" && (object.FilterValue.IsNull() || object.FilterValue.ValueString() == "") {
		resp.Diagnostics.AddError(
			"Missing Required Parameter",
			fmt.Sprintf("filter_value must be provided when filter_type is '%s'", filterType),
		)
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Read, ReadTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	// Fetch all catalog items
	catalogItemsResult, err := d.client.
		DeviceManagement().
		WindowsUpdateCatalogItems().
		Get(ctx, nil)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Windows Update Catalog Items",
			fmt.Sprintf("Could not read Windows Update Catalog Items: %s", err),
		)
		return
	}

	// Parse date filter if necessary
	var releaseDateTime, endOfSupportDate *time.Time
	
	if filterType == "release_date_time" || filterType == "end_of_support_date" {
		parsedTime, err := time.Parse(time.RFC3339, object.FilterValue.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Invalid Date Format",
				fmt.Sprintf("Could not parse date value as RFC3339 format: %s", err),
			)
			return
		}
		
		if filterType == "release_date_time" {
			releaseDateTime = &parsedTime
		} else {
			endOfSupportDate = &parsedTime
		}
	}

	// Filter the results based on the specified filter_type and filter_value
	var filteredItems []WindowsUpdateCatalogItemModel
	filterValue := object.FilterValue.ValueString()
	
	for _, item := range catalogItemsResult.GetValue() {
		switch filterType {
		case "all":
			// No filtering, include all items
			filteredItems = append(filteredItems, MapCatalogItemToModel(item))
			
		case "id":
			// Filter by ID (exact match)
			if item.GetId() != nil && *item.GetId() == filterValue {
				filteredItems = append(filteredItems, MapCatalogItemToModel(item))
			}
			
		case "display_name":
			// Filter by display name (case-insensitive substring match)
			if item.GetDisplayName() != nil && strings.Contains(
				strings.ToLower(*item.GetDisplayName()),
				strings.ToLower(filterValue)) {
				filteredItems = append(filteredItems, MapCatalogItemToModel(item))
			}
			
		case "release_date_time":
			// Filter by release date time (exact match)
			itemReleaseDate := item.GetReleaseDateTime()
			if itemReleaseDate != nil && releaseDateTime != nil && itemReleaseDate.Equal(*releaseDateTime) {
				filteredItems = append(filteredItems, MapCatalogItemToModel(item))
			}
			
		case "end_of_support_date":
			// Filter by end of support date (exact match)
			itemEndOfSupportDate := item.GetEndOfSupportDate()
			if itemEndOfSupportDate != nil && endOfSupportDate != nil && itemEndOfSupportDate.Equal(*endOfSupportDate) {
				filteredItems = append(filteredItems, MapCatalogItemToModel(item))
			}
		}
	}

	// Update the model with the filtered items
	object.Items = filteredItems

	// Set the data in the response
	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Datasource Read Method: %s_%s, found %d items", d.ProviderTypeName, d.TypeName, len(filteredItems)))
}
