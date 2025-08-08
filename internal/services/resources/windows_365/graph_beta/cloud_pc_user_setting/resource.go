package graphBetaCloudPcUserSetting

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	planmodifiers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/plan_modifiers"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName  = "graph_beta_windows_365_cloud_pc_user_setting"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &CloudPcUserSettingResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &CloudPcUserSettingResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &CloudPcUserSettingResource{}

	// Enables plan modification/diff suppression
	_ resource.ResourceWithModifyPlan = &CloudPcUserSettingResource{}
)

func NewCloudPcUserSettingResource() resource.Resource {
	return &CloudPcUserSettingResource{
		ReadPermissions: []string{
			"CloudPC.Read.All",
		},
		WritePermissions: []string{
			"CloudPC.ReadWrite.All",
		},
	}
}

type CloudPcUserSettingResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ProviderTypeName string
	TypeName         string
	ReadPermissions  []string
	WritePermissions []string
}

// Metadata returns the resource type name.
func (r *CloudPcUserSettingResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	r.ProviderTypeName = req.ProviderTypeName
	r.TypeName = ResourceName
	resp.TypeName = r.FullTypeName()
}

// FullTypeName returns the full resource type name in the format "providername_resourcename".
func (r *CloudPcUserSettingResource) FullTypeName() string {
	return r.ProviderTypeName + "_" + r.TypeName
}

// Configure sets the client for the resource.
func (r *CloudPcUserSettingResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, constants.PROVIDER_NAME+"_"+ResourceName)
}

// ImportState imports the resource state.
func (r *CloudPcUserSettingResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Schema returns the schema for the resource.
func (r *CloudPcUserSettingResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Unique identifier for the Cloud PC user setting. Read-only.",
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
			},
			"display_name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The setting name displayed in the user interface.",
			},
			"created_date_time": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The date and time the setting was created. The timestamp type represents the date and time information using ISO 8601 format and is always in UTC. For example, midnight UTC on Jan 1, 2014 looks like this: '2014-01-01T00:00:00Z'.",
			},
			"last_modified_date_time": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The last date and time the setting was modified. The timestamp type represents the date and time information using ISO 8601 format and is always in UTC. For example, midnight UTC on Jan 1, 2014 looks like this: '2014-01-01T00:00:00Z'.",
			},
			"local_admin_enabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Indicates whether the local admin option is enabled. Default value is false. To enable the local admin option, change the setting to true. If the local admin option is enabled, the end user can be an admin of the Cloud PC device.",
			},
			"reset_enabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Indicates whether an end user is allowed to reset their Cloud PC. When true, the user is allowed to reset their Cloud PC. When false, end-user initiated reset isn't allowed. The default value is false.",
			},
			"self_service_enabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Indicates whether the self-service option is enabled. Default value is false. To enable the self-service option, change the setting to true. If the self-service option is enabled, the end user is allowed to perform some self-service operations, such as upgrading the Cloud PC through the end user portal. The **selfServiceEnabled** property is deprecated and will stop returning data on December 1, 2023.",
			},
			"restore_point_setting": schema.SingleNestedAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Defines how frequently a restore point is created that is, a snapshot is taken) for users' provisioned Cloud PCs (default is 12 hours), and whether the user is allowed to restore their own Cloud PCs to a backup made at a specific point in time.",
				Attributes: map[string]schema.Attribute{
					"frequency_in_hours": schema.Int32Attribute{
						Optional:            true,
						Computed:            true,
						Default:             int32default.StaticInt32(12),
						MarkdownDescription: "The frequency in hours at which restore points are created.",
					},
					"frequency_type": schema.StringAttribute{
						Optional:            true,
						Computed:            true,
						Default:             stringdefault.StaticString("default"),
						MarkdownDescription: "The time interval in hours to take snapshots (restore points) of a Cloud PC automatically. Possible values are: default, fourHours, sixHours, twelveHours, sixteenHours, twentyFourHours, unknownFutureValue. The default value is default that indicates that the time interval for automatic capturing of restore point snapshots is set to 12 hours.",
						Validators: []validator.String{
							stringvalidator.OneOf("default", "fourHours", "sixHours", "twelveHours", "sixteenHours", "twentyFourHours", "unknownFutureValue"),
						},
					},
					"user_restore_enabled": schema.BoolAttribute{
						Optional:            true,
						Computed:            true,
						Default:             booldefault.StaticBool(false),
						MarkdownDescription: "Indicates whether the user is allowed to restore their own Cloud PCs.",
					},
				},
			},
			"cross_region_disaster_recovery_setting": schema.SingleNestedAttribute{
				Optional:            true,
				MarkdownDescription: "Defines whether the user's Cloud PC enables cross-region disaster recovery and specifies the network for the disaster recovery.",
				Attributes: map[string]schema.Attribute{
					"maintain_cross_region_restore_point_enabled": schema.BoolAttribute{
						Optional: true,
						Computed: true,
						Default:  booldefault.StaticBool(false),
						MarkdownDescription: "Indicates whether Windows 365 maintain the cross-region disaster recovery function generated restore points." +
							"If true, the Windows 365 stored restore points; false indicates that Windows 365 doesn't generate or keep the restore point from " +
							"the original Cloud PC. If a disaster occurs, the new Cloud PC can only be provisioned using the initial image." +
							"This limitation can result in the loss of some user data on the original Cloud PC. The default value is false.",
					},
					"user_initiated_disaster_recovery_allowed": schema.BoolAttribute{
						Optional: true,
						Computed: true,
						Default:  booldefault.StaticBool(false),
						MarkdownDescription: "Indicates whether the client allows the end user to initiate a disaster recovery activation." +
							"True indicates that the client includes the option for the end user to activate Backup Cloud PC. When false, the end" +
							"user doesn't have the option to activate disaster recovery. The default value is false. Currently, only premium disaster recovery is supported.",
					},
					"disaster_recovery_type": schema.StringAttribute{
						Optional: true,
						Computed: true,
						Default:  stringdefault.StaticString("notConfigured"),
						MarkdownDescription: "Indicates the type of disaster recovery to perform when a disaster occurs on the user's Cloud PC." +
							"The possible values are: notConfigured, crossRegion, premium, unknownFutureValue. The default value is notConfigured.",
						Validators: []validator.String{
							stringvalidator.OneOf("notConfigured", "crossRegion", "premium", "unknownFutureValue"),
						},
					},
					"disaster_recovery_network_setting": schema.SingleNestedAttribute{
						Optional:            true,
						MarkdownDescription: "Indicates the network settings of the Cloud PC during a cross-region disaster recovery operation.",
						Attributes: map[string]schema.Attribute{
							"network_type": schema.StringAttribute{
								Optional:            true,
								Computed:            true,
								Default:             stringdefault.StaticString("microsoftHosted"),
								MarkdownDescription: "The type of network for disaster recovery.",
								Validators: []validator.String{
									stringvalidator.OneOf("microsoftHosted", "azureNetworkConnection"),
								},
							},
							"region_name": schema.StringAttribute{
								Optional:            true,
								Computed:            true,
								Default:             stringdefault.StaticString("automatic"),
								MarkdownDescription: "The region name for the disaster recovery network. Default is 'automatic'.",
							},
							"region_group": schema.StringAttribute{
								Optional:            true,
								Computed:            true,
								Default:             stringdefault.StaticString("usWest"),
								MarkdownDescription: "The region group for the disaster recovery network. Default is 'usWest'.",
								Validators: []validator.String{
									stringvalidator.OneOf("asia", "australia", "canada", "europeUnion", "france", "germany", "india", "japan", "mexico", "usEast", "middleEast", "norway", "southAfrica", "southKorea", "switzerland", "uae", "unitedKingdom", "usCentral", "usWest", "usEast"),
								},
							},
						},
					},
				},
			},
			"notification_setting": schema.SingleNestedAttribute{
				Optional:            true,
				MarkdownDescription: "Defines the setting of the Cloud PC notification prompts for the Cloud PC user.",
				Attributes: map[string]schema.Attribute{
					"restart_prompts_disabled": schema.BoolAttribute{
						Optional:            true,
						Computed:            true,
						Default:             booldefault.StaticBool(false),
						MarkdownDescription: "If true, doesn't prompt the user to restart the Cloud PC. If false, prompts the user to restart Cloud PC. The default value is false.",
					},
				},
			},
			"assignments": Windows365UserSettingsAssignmentSchema(),
			"timeouts":    commonschema.Timeouts(ctx),
		},
	}
}

// Windows365UserSettingsAssignmentSchema returns the schema for the assignments attribute
func Windows365UserSettingsAssignmentSchema() schema.SetNestedAttribute {
	return schema.SetNestedAttribute{
		Optional:            true,
		MarkdownDescription: "Assignments of the Cloud PC user setting to groups. Only Microsoft 365 groups and security groups in Microsoft Entra ID are currently supported.",
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"type": schema.StringAttribute{
					Required:            true,
					MarkdownDescription: "The type of assignment target. Valid values are 'groupAssignmentTarget'.",
					Validators: []validator.String{
						stringvalidator.OneOf(
							"groupAssignmentTarget",
						),
					},
				},
				"group_id": schema.StringAttribute{
					Required:            true,
					MarkdownDescription: "The ID of the Microsoft 365 group or security group in Microsoft Entra ID. Required when type is 'groupAssignmentTarget' or 'exclusionGroupAssignmentTarget'.",
				},
			},
		},
	}
}
