package graphBetaWindowsEnrollmentStatusPage

import (
	"context"
	"regexp"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	planmodifiers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/plan_modifiers"
	commonschemagraphbeta "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema/graph_beta/device_management"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName  = "graph_beta_device_management_windows_enrollment_status_page"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &WindowsEnrollmentStatusPageResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &WindowsEnrollmentStatusPageResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &WindowsEnrollmentStatusPageResource{}

	// Enables plan modification/diff suppression
	_ resource.ResourceWithModifyPlan = &WindowsEnrollmentStatusPageResource{}
)

func NewWindowsEnrollmentStatusPageResource() resource.Resource {
	return &WindowsEnrollmentStatusPageResource{
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

type WindowsEnrollmentStatusPageResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ProviderTypeName string
	TypeName         string
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

func (r *WindowsEnrollmentStatusPageResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	r.ProviderTypeName = req.ProviderTypeName
	r.TypeName = ResourceName
	resp.TypeName = r.FullTypeName()
}

func (r *WindowsEnrollmentStatusPageResource) FullTypeName() string {
	return r.ProviderTypeName + "_" + r.TypeName
}

func (r *WindowsEnrollmentStatusPageResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, constants.PROVIDER_NAME+"_"+ResourceName)
}

func (r *WindowsEnrollmentStatusPageResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *WindowsEnrollmentStatusPageResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a Windows 10 Enrollment Status Page configuration in Microsoft Intune." +
			" Using the `/deviceManagement/deviceEnrollmentConfigurations/{deviceEnrollmentConfigurationId}` endpoint." +
			"The Enrollment Status Page (ESP) displays the progress of the device setup process during " +
			"Windows Autopilot provisioning or when a device first enrolls in Microsoft Intune.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Unique identifier for the enrollment status page configuration.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},

			"display_name": schema.StringAttribute{
				Description: "The display name of the enrollment status page configuration.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					stringvalidator.LengthAtMost(1000),
				},
			},

			"description": schema.StringAttribute{
				Description: "The description of the enrollment status page configuration.",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(1000),
				},
			},

			"priority": schema.Int32Attribute{
				Description: "Priority is used when a user exists in multiple groups that are assigned enrollment " +
					"configuration. Users are subject only to the configuration with the lowest priority value.",
				Computed: true,
			},

			"show_installation_progress": schema.BoolAttribute{
				Description: "Show or hide installation progress to user.",
				Required:    true,
			},

			"block_device_setup_retry_by_user": schema.BoolAttribute{
				Description: "Allow the user to retry the setup on installation failure.",
				Required:    true,
			},

			"allow_device_reset_on_install_failure": schema.BoolAttribute{
				Description: "Allow or block device reset on installation failure.",
				Required:    true,
			},

			"allow_log_collection_on_install_failure": schema.BoolAttribute{
				Description: "Allow or block log collection on installation failure.",
				Required:    true,
			},

			"custom_error_message": schema.StringAttribute{
				Description: "Set custom error message to show upon installation failure.",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(1000),
				},
			},

			"install_progress_timeout_in_minutes": schema.Int32Attribute{
				Description: "Set installation progress timeout in minutes. Valid values are 1 to 1440 (24 hours).",
				Optional:    true,
				Computed:    true,
				Default:     int32default.StaticInt32(60),
				Validators: []validator.Int32{
					int32validator.Between(1, 1440),
				},
			},

			"allow_device_use_on_install_failure": schema.BoolAttribute{
				Description: "Allow the user to continue using the device on installation failure.",
				Required:    true,
			},

			"selected_mobile_app_ids": schema.SetAttribute{
				Description: "Selected applications to track the installation status. This collection can contain a maximum of 100 elements.",
				ElementType: types.StringType,
				Optional:    true,
				Validators: []validator.Set{
					setvalidator.SizeAtMost(100),
					setvalidator.ValueStringsAre(
						stringvalidator.RegexMatches(
							regexp.MustCompile(constants.GuidRegex),
							"must be a valid GUID in the format 00000000-0000-0000-0000-000000000000",
						),
					),
				},
			},

			"track_install_progress_for_autopilot_only": schema.BoolAttribute{
				Description: "Only show installation progress for Autopilot enrollment scenarios.",
				Required:    true,
			},

			"disable_user_status_tracking_after_first_user": schema.BoolAttribute{
				Description: "Only show installation progress for first user post enrollment.",
				Required:    true,
			},

			"role_scope_tag_ids": schema.SetAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Set of scope tag IDs for this Settings Catalog template profile.",
				PlanModifiers: []planmodifier.Set{
					planmodifiers.DefaultSetValue(
						[]attr.Value{types.StringValue("0")},
					),
				},
			},

			"assignments": commonschemagraphbeta.DeviceConfigurationWithAllLicensedUsersAllDevicesInclusionGroupAssignmentsSchema(),

			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
				Create: true,
				Read:   true,
				Update: true,
				Delete: true,
			}),
		},
	}
}
