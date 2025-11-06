package graphBetaWindowsEnrollmentStatusPage

import (
	"context"
	"regexp"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	planmodifiers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/plan_modifiers"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	commonschemagraphbeta "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema/graph_beta/device_management"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/validators"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
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
				MarkdownDescription: "Unique identifier for the enrollment status page configuration.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"display_name": schema.StringAttribute{
				MarkdownDescription: "The display name of the enrollment status page configuration.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					stringvalidator.LengthAtMost(1000),
				},
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description of the enrollment status page configuration.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(1000),
				},
			},
			"show_installation_progress": schema.BoolAttribute{
				MarkdownDescription: "Show or hide installation progress to user.",
				Required:            true,
			},
			"custom_error_message": schema.StringAttribute{
				MarkdownDescription: "Set custom message when time limit or error occurs during initial device setup.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(1000),
				},
			},
			"install_quality_updates": schema.BoolAttribute{
				MarkdownDescription: "Whether to install quality updates during the enrollment status page experience.",
				Required:            true,
			},
			"install_progress_timeout_in_minutes": schema.Int32Attribute{
				MarkdownDescription: "Set installation progress timeout in minutes. Valid values are 1 to 1440 (24 hours).",
				Required:            true,
				Validators: []validator.Int32{
					int32validator.Between(1, 1440),
				},
			},
			"allow_log_collection_on_install_failure": schema.BoolAttribute{
				MarkdownDescription: "Allow or block log collection on installation failure.",
				Required:            true,
			},
			"only_show_page_to_devices_provisioned_by_out_of_box_experience_oobe": schema.BoolAttribute{
				MarkdownDescription: "Only show autopilot status page during initial device setup and during first user sign in to devices provisioned by oobe.",
				Required:            true,
			},
			"block_device_use_until_all_apps_and_profiles_are_installed": schema.BoolAttribute{
				MarkdownDescription: "Allow the user to retry the setup on installation failure.",
				Required:            true,
			},
			"allow_device_reset_on_install_failure": schema.BoolAttribute{
				MarkdownDescription: "Allow or block device reset on installation failure. When block_device_use_until_all_apps_and_profiles_are_installed is true, this must be false. When block_device_use_until_all_apps_and_profiles_are_installed is false, this can be true or false.",
				Required:            true,
				Validators: []validator.Bool{
					validators.BoolCanOnlyBeFalseWhen(
						"block_device_use_until_all_apps_and_profiles_are_installed",
						true,
						"allow_device_reset_on_install_failure must be false when block_device_use_until_all_apps_and_profiles_are_installed is true",
					),
				},
			},
			"allow_device_use_on_install_failure": schema.BoolAttribute{
				MarkdownDescription: "Allow the user to continue using the device on installation failure. When block_device_use_until_all_apps_and_profiles_are_installed is true, this must be false. When block_device_use_until_all_apps_and_profiles_are_installed is false, this can be true or false.",
				Required:            true,
				Validators: []validator.Bool{
					validators.BoolCanOnlyBeFalseWhen(
						"block_device_use_until_all_apps_and_profiles_are_installed",
						true,
						"allow_device_use_on_install_failure must be false when block_device_use_until_all_apps_and_profiles_are_installed is true",
					),
				},
			},
			"selected_mobile_app_ids": schema.SetAttribute{
				MarkdownDescription: "Selected applications to track the installation status. This collection can contain a maximum of 100 elements. Can only be used effectively when block_device_use_until_all_apps_and_profiles_are_installed is true.",
				ElementType:         types.StringType,
				Optional:            true,
				Validators: []validator.Set{
					setvalidator.SizeAtMost(100),
					setvalidator.ValueStringsAre(
						stringvalidator.RegexMatches(
							regexp.MustCompile(constants.GuidRegex),
							"must be a valid GUID in the format 00000000-0000-0000-0000-000000000000",
						),
					),
					validators.SetRequiresBoolValue(
						"block_device_use_until_all_apps_and_profiles_are_installed",
						false,
						"selected_mobile_app_ids can only be specified when block_device_use_until_all_apps_and_profiles_are_installed is false",
					),
				},
			},
			"only_fail_selected_blocking_apps_in_technician_phase": schema.BoolAttribute{
				MarkdownDescription: "When true, only the selected blocking apps will be failed in the technician phase. When block_device_use_until_all_apps_and_profiles_are_installed is true, this must be false. When block_device_use_until_all_apps_and_profiles_are_installed is false, this can be true or false.",
				Required:            true,
				Validators: []validator.Bool{
					validators.BoolCanOnlyBeFalseWhen(
						"block_device_use_until_all_apps_and_profiles_are_installed",
						true,
						"only_fail_selected_blocking_apps_in_technician_phase must be false when block_device_use_until_all_apps_and_profiles_are_installed is true",
					),
				},
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
			"timeouts":    commonschema.Timeouts(ctx),
		},
	}
}
