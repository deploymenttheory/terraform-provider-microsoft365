package graphBetaDeviceEnrollmentLimitConfiguration

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common"
	planmodifiers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/plan_modifiers"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName  = "graph_beta_device_management_device_enrollment_limit_configuration"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &DeviceEnrollmentLimitConfigurationResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &DeviceEnrollmentLimitConfigurationResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &DeviceEnrollmentLimitConfigurationResource{}

	// Enables plan modification/diff suppression
	_ resource.ResourceWithModifyPlan = &DeviceEnrollmentLimitConfigurationResource{}
)

func NewDeviceEnrollmentLimitConfigurationResource() resource.Resource {
	return &DeviceEnrollmentLimitConfigurationResource{
		ReadPermissions: []string{
			"DeviceManagementServiceConfig.Read.All",
			"DeviceManagementConfiguration.Read.All",
		},
		WritePermissions: []string{
			"DeviceManagementServiceConfig.ReadWrite.All",
			"DeviceManagementConfiguration.ReadWrite.All",
		},
		ResourcePath: "/deviceManagement/deviceEnrollmentConfigurations",
	}
}

type DeviceEnrollmentLimitConfigurationResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ProviderTypeName string
	TypeName         string
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

// Metadata returns the resource type name.
func (r *DeviceEnrollmentLimitConfigurationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	r.ProviderTypeName = req.ProviderTypeName
	r.TypeName = ResourceName
	resp.TypeName = r.FullTypeName()
}

// FullTypeName returns the full type name of the resource for logging purposes.
func (r *DeviceEnrollmentLimitConfigurationResource) FullTypeName() string {
	return r.ProviderTypeName + "_" + r.TypeName
}

// Configure sets the client for the resource.
func (r *DeviceEnrollmentLimitConfigurationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = common.SetGraphBetaClientForResource(ctx, req, resp, constants.PROVIDER_NAME+"_"+ResourceName)
}

// ImportState imports the resource state.
func (r *DeviceEnrollmentLimitConfigurationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *DeviceEnrollmentLimitConfigurationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Resource is currently broken as per this issue: `https://github.com/microsoft/Microsoft365DSC/issues/5127`. Manages device enrollment limit configurations using the `/deviceManagement/deviceEnrollmentConfigurations` endpoint. Device enrollment limit configurations restrict the number of devices a user can enroll in the organization.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
				MarkdownDescription: "The unique identifier of the device enrollment limit configuration.",
			},
			"display_name": schema.StringAttribute{
				MarkdownDescription: "The display name of the device enrollment configuration.",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description of the device enrollment configuration.",
				Optional:            true,
			},
			"priority": schema.Int32Attribute{
				MarkdownDescription: "Priority is used when a user exists in multiple groups that are assigned enrollment configuration. Users are subject only to the configuration with the lowest priority value.",
				Optional:            true,
				Validators: []validator.Int32{
					int32validator.AtLeast(0),
				},
			},
			"created_date_time": schema.StringAttribute{
				MarkdownDescription: "Created date time in UTC of the device enrollment configuration. This property is read-only.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"last_modified_date_time": schema.StringAttribute{
				MarkdownDescription: "Last modified date time in UTC of the device enrollment configuration. This property is read-only.",
				Computed:            true,
			},
			"version": schema.Int32Attribute{
				MarkdownDescription: "The version of the device enrollment configuration. This property is read-only.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int32{
					int32planmodifier.UseStateForUnknown(),
				},
			},
			"role_scope_tag_ids": schema.SetAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "Optional role scope tags for the enrollment restrictions.",
				Optional:            true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"device_enrollment_configuration_type": schema.StringAttribute{
				MarkdownDescription: "Support for Enrollment Configuration Type. Possible values are: `unknown`, `limit`, `platformRestrictions`, `windowsHelloForBusiness`, `defaultLimit`, `defaultPlatformRestrictions`, `defaultWindowsHelloForBusiness`, `defaultWindows10EnrollmentCompletionPageConfiguration`, `windows10EnrollmentCompletionPageConfiguration`, `deviceComanagementAuthorityConfiguration`, `singlePlatformRestriction`, `unknownFutureValue`, `enrollmentNotificationsConfiguration`, `windowsRestore`.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf("unknown", "limit", "platformRestrictions", "windowsHelloForBusiness", "defaultLimit", "defaultPlatformRestrictions", "defaultWindowsHelloForBusiness", "defaultWindows10EnrollmentCompletionPageConfiguration", "windows10EnrollmentCompletionPageConfiguration", "deviceComanagementAuthorityConfiguration", "singlePlatformRestriction", "unknownFutureValue", "enrollmentNotificationsConfiguration", "windowsRestore"),
				},
			},
			"limit": schema.Int32Attribute{
				MarkdownDescription: "The maximum number of devices that a user can enroll.",
				Required:            true,
				Validators: []validator.Int32{
					int32validator.AtLeast(1),
				},
			},
			"timeouts": commonschema.Timeouts(ctx),
		},
	}
}
