package graphCloudPcUserSetting

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/schema"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
)

var _ resource.Resource = &CloudPcUserSettingResource{}
var _ resource.ResourceWithConfigure = &CloudPcUserSettingResource{}
var _ resource.ResourceWithImportState = &CloudPcUserSettingResource{}

func NewCloudPcUserSettingResource() resource.Resource {
	return &CloudPcUserSettingResource{}
}

type CloudPcUserSettingResource struct {
	client           *msgraphsdk.GraphServiceClient
	ProviderTypeName string
	TypeName         string
}

// GetID returns the ID of a resource from the state model.
func (s *CloudPcUserSettingResourceModel) GetID() string {
	return s.ID.ValueString()
}

// GetTypeName returns the type name of the resource from the state model.
func (r *CloudPcUserSettingResource) GetTypeName() string {
	return r.TypeName
}

// Metadata returns the resource type name.
func (r *CloudPcUserSettingResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_graph_cloud_pc_user_setting"
}

// Configure sets the client for the resource.
func (r *CloudPcUserSettingResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = common.SetGraphStableClient(ctx, req, resp, r.TypeName)
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
				Computed:    true,
				Description: "Unique identifier for the Cloud PC user setting. Read-only.",
			},
			"created_date_time": schema.StringAttribute{
				Computed:    true,
				Description: "The date and time when the setting was created. The timestamp type represents the date and time information using ISO 8601 format and is always in UTC. For example, midnight UTC on Jan 1, 2014 is `2014-01-01T00:00:00Z`. Read-only.",
			},
			"display_name": schema.StringAttribute{
				Required:    true,
				Description: "The setting name displayed in the user interface.",
			},
			"last_modified_date_time": schema.StringAttribute{
				Computed:    true,
				Description: "The date and time when the setting was last modified. The timestamp type represents the date and time information using ISO 8601 format and is always in UTC. For example, midnight UTC on Jan 1, 2014 is `2014-01-01T00:00:00Z`. Read-only.",
			},
			"local_admin_enabled": schema.BoolAttribute{
				Optional:    true,
				Description: "Indicates whether the local admin option is enabled. The default value is `false`. To enable the local admin option, change the setting to `true`. If the local admin option is enabled, the end user can be an admin of the Cloud PC device.",
			},
			"reset_enabled": schema.BoolAttribute{
				Optional:    true,
				Description: "Indicates whether an end user is allowed to reset their Cloud PC. When `true`, the user is allowed to reset their Cloud PC. When `false`, end-user initiated reset is not allowed. The default value is `false`.",
			},
			"restore_point_setting": schema.SingleNestedAttribute{
				Optional:    true,
				Description: "Defines how frequently a restore point is created (that is, a snapshot is taken) for users' provisioned Cloud PCs (default is 12 hours), and whether the user is allowed to restore their own Cloud PCs to a backup made at a specific point in time.",
				Attributes:  r.cloudPcRestorePointSettingSchema(),
			},
			"timeouts": commonschema.Timeouts(ctx),
		},
	}
}

func (r *CloudPcUserSettingResource) cloudPcRestorePointSettingSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"frequency_type": schema.StringAttribute{
			Optional:    true,
			Description: "The time interval in hours to take snapshots (restore points) of a Cloud PC automatically. Possible values are: `default`, `fourHours`, `sixHours`, `twelveHours`, `sixteenHours`, `twentyFourHours`, `unknownFutureValue`. The default value is `default` which indicates that the time interval for automatic capturing of restore point snapshots is set to 12 hours.",
			Validators: []validator.String{
				stringvalidator.OneOf("default", "fourHours", "sixHours", "twelveHours", "sixteenHours", "twentyFourHours", "unknownFutureValue"),
			},
		},
		"user_restore_enabled": schema.BoolAttribute{
			Optional:    true,
			Description: "If `true`, the user has the ability to use snapshots to restore Cloud PCs. If `false`, non-admin users can't use snapshots to restore the Cloud PC.",
		},
	}
}
