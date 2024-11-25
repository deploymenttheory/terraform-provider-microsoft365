package graphbetadevicemanagementscript

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

var (
	// Basic resource interface (CRUD operations)
	_ datasource.DataSource = &DeviceManagementScriptDataSource{}

	// Allows the resource to be configured with the provider client
	_ datasource.DataSourceWithConfigure = &DeviceManagementScriptDataSource{}
)

func NewDeviceManagementScriptDataSource() datasource.DataSource {
	return &DeviceManagementScriptDataSource{
		ReadPermissions: []string{
			"DeviceManagementConfiguration.Read.All",
		},
	}
}

type DeviceManagementScriptDataSource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ProviderTypeName string
	TypeName         string
	ReadPermissions  []string
}

func (d *DeviceManagementScriptDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_graph_beta_device_and_app_management_device_management_script"
}

func (d *DeviceManagementScriptDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves information about a device management script.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Unique identifier for the device management script.",
				Required:    true,
			},
			"display_name": schema.StringAttribute{
				Description: "Name of the device management script.",
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "Description of the device management script.",
				Computed:    true,
			},
			"created_date_time": schema.StringAttribute{
				Description: "The date and time the device management script was created.",
				Computed:    true,
			},
			"last_modified_date_time": schema.StringAttribute{
				Description: "The date and time the device management script was last modified.",
				Computed:    true,
			},
			"run_as_account": schema.StringAttribute{
				Description: "Indicates the type of execution context.",
				Computed:    true,
			},
			"enforce_signature_check": schema.BoolAttribute{
				Description: "Indicate whether the script signature needs be checked.",
				Computed:    true,
			},
			"file_name": schema.StringAttribute{
				Description: "Script file name.",
				Computed:    true,
			},
			"run_as_32_bit": schema.BoolAttribute{
				Description: "A value indicating whether the PowerShell script should run as 32-bit.",
				Computed:    true,
			},
			"role_scope_tag_ids": schema.ListAttribute{
				Description: "List of Scope Tag IDs for this PowerShellScript instance.",
				Computed:    true,
				ElementType: types.StringType,
			},
			"script_content": schema.StringAttribute{
				Description: "The script content.",
				Computed:    true,
				Sensitive:   true,
			},
			"assignments": schema.ListNestedAttribute{
				Description: "The assignments of the device management script.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "Key of the device management script assignment entity.",
							Computed:    true,
						},
						"target": schema.SingleNestedAttribute{
							Description: "The target of the assignment.",
							Computed:    true,
							Attributes: map[string]schema.Attribute{
								"device_and_app_management_assignment_filter_id": schema.StringAttribute{
									Description: "The Id of the filter for the target assignment.",
									Computed:    true,
								},
								"device_and_app_management_assignment_filter_type": schema.StringAttribute{
									Description: "The type of filter of the target assignment.",
									Computed:    true,
								},
								"target_type": schema.StringAttribute{
									Description: "The target type of the assignment.",
									Computed:    true,
								},
								"entra_object_id": schema.StringAttribute{
									Description: "The ID of the Azure Active Directory object.",
									Computed:    true,
								},
							},
						},
					},
				},
			},
			"group_assignments": schema.ListNestedAttribute{
				Description: "The group assignments of the device management script.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "Key of the device management script group assignment entity.",
							Computed:    true,
						},
						"target_group_id": schema.StringAttribute{
							Description: "The Id of the Azure Active Directory group we are targeting the script to.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *DeviceManagementScriptDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = common.SetGraphBetaClientForDataSource(ctx, req, resp, d.TypeName)
}
