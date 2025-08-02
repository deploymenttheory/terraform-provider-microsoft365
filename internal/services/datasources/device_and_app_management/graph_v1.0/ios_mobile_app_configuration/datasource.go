package graphIOSMobileAppConfiguration

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
)

const (
	DataSourceName = "graph_v1_device_and_app_management_ios_mobile_app_configuration"
	ReadTimeout    = 180
)

var (
	// Basic data source interface (Read operations)
	_ datasource.DataSource = &IOSMobileAppConfigurationDataSource{}

	// Allows the data source to be configured with the provider client
	_ datasource.DataSourceWithConfigure = &IOSMobileAppConfigurationDataSource{}
)

func NewIOSMobileAppConfigurationDataSource() datasource.DataSource {
	return &IOSMobileAppConfigurationDataSource{
		ReadPermissions: []string{
			"DeviceManagementApps.Read.All",
		},
	}
}

type IOSMobileAppConfigurationDataSource struct {
	client           *msgraphsdk.GraphServiceClient
	ProviderTypeName string
	TypeName         string
	ReadPermissions  []string
}

// Metadata returns the data source type name.
func (d *IOSMobileAppConfigurationDataSource) Metadata(
	ctx context.Context,
	req datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_" + DataSourceName
}

// Configure sets the client for the data source
func (d *IOSMobileAppConfigurationDataSource) Configure(
	ctx context.Context,
	req datasource.ConfigureRequest,
	resp *datasource.ConfigureResponse,
) {
	d.client = client.SetGraphStableClientForDataSource(ctx, req, resp, d.TypeName)
}

// Schema defines the schema for the data source
func (d *IOSMobileAppConfigurationDataSource) Schema(
	ctx context.Context,
	req datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieves information about iOS mobile app configurations in Microsoft Intune.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Unique identifier for the iOS mobile app configuration.",
				Optional:            true,
				Computed:            true,
			},
			"display_name": schema.StringAttribute{
				MarkdownDescription: "The display name of the iOS mobile app configuration.",
				Optional:            true,
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description of the iOS mobile app configuration.",
				Computed:            true,
			},
			"targeted_mobile_apps": schema.ListAttribute{
				MarkdownDescription: "The list of targeted mobile app IDs.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"encoded_setting_xml": schema.StringAttribute{
				MarkdownDescription: "Base64 encoded configuration XML.",
				Computed:            true,
				Sensitive:           true,
			},
			"created_date_time": schema.StringAttribute{
				MarkdownDescription: "DateTime the object was created.",
				Computed:            true,
			},
			"last_modified_date_time": schema.StringAttribute{
				MarkdownDescription: "DateTime the object was last modified.",
				Computed:            true,
			},
			"version": schema.Int32Attribute{
				MarkdownDescription: "Version of the device configuration.",
				Computed:            true,
			},
			"timeouts": commonschema.Timeouts(ctx),
		},
		Blocks: map[string]schema.Block{
			"settings": schema.ListNestedBlock{
				MarkdownDescription: "iOS app configuration settings.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"app_config_key": schema.StringAttribute{
							MarkdownDescription: "The application configuration key.",
							Computed:            true,
						},
						"app_config_key_type": schema.StringAttribute{
							MarkdownDescription: "The application configuration key type.",
							Computed:            true,
						},
						"app_config_key_value": schema.StringAttribute{
							MarkdownDescription: "The application configuration key value.",
							Computed:            true,
						},
					},
				},
			},
			"assignments": schema.ListNestedBlock{
				MarkdownDescription: "The list of assignments for this iOS mobile app configuration.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "Key of the entity.",
							Computed:            true,
						},
					},
					Blocks: map[string]schema.Block{
						"target": schema.SingleNestedBlock{
							MarkdownDescription: "Target for the assignment.",
							Attributes: map[string]schema.Attribute{
								"odata_type": schema.StringAttribute{
									MarkdownDescription: "The type of assignment target.",
									Computed:            true,
								},
								"group_id": schema.StringAttribute{
									MarkdownDescription: "The group Id that is the target of the assignment.",
									Computed:            true,
								},
							},
						},
					},
				},
			},
		},
	}
}
