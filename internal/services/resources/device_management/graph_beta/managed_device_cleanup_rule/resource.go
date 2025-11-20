package graphBetaManagedDeviceCleanupRule

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName  = "microsoft365_graph_beta_device_management_managed_device_cleanup_rule"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &ManagedDeviceCleanupRuleResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &ManagedDeviceCleanupRuleResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &ManagedDeviceCleanupRuleResource{}

	// Enables plan modification/diff suppression
	_ resource.ResourceWithModifyPlan = &ManagedDeviceCleanupRuleResource{}
)

func NewManagedDeviceCleanupRuleResource() resource.Resource {
	return &ManagedDeviceCleanupRuleResource{
		ReadPermissions: []string{
			"DeviceManagementConfiguration.Read.All",
			"DeviceManagementManagedDevices.Read.All,",
		},
		WritePermissions: []string{
			"DeviceManagementConfiguration.ReadWrite.All",
			"DeviceManagementManagedDevices.ReadWrite.All",
		},
		ResourcePath: "deviceManagement/managedDeviceCleanupRules",
	}
}

type ManagedDeviceCleanupRuleResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

// Metadata returns the resource type name.
func (r *ManagedDeviceCleanupRuleResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

// Configure sets the client for the resource.
func (r *ManagedDeviceCleanupRuleResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

// ImportState imports the resource state.
func (r *ManagedDeviceCleanupRuleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Schema returns the schema for the resource.
func (r *ManagedDeviceCleanupRuleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages Managed Device Cleanup Rules in Microsoft Intune. Device cleanup rules define when the admin wants devices to be automatically removed from Intune management based on inactivity periods and platform types.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier of the managed device cleanup rule.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"display_name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The display name of the device cleanup rule.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 256),
				},
			},
			"description": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The description for the device cleanup rule.",
				Validators: []validator.String{
					stringvalidator.LengthAtMost(1024),
				},
			},
			"device_cleanup_rule_platform_type": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The managed device platform for which the admin wants to create the device cleanup rule. Possible values are: `all`, `androidAOSP`, `androidDeviceAdministrator`, `androidDedicatedAndFullyManagedCorporateOwnedWorkProfile`, `chromeOS`, `androidPersonallyOwnedWorkProfile`, `ios`, `macOS`, `windows`, `windowsHolographic`, `unknownFutureValue`.",
				Validators: []validator.String{
					stringvalidator.OneOf(
						"all",
						"androidAOSP",
						"androidDeviceAdministrator",
						"androidDedicatedAndFullyManagedCorporateOwnedWorkProfile",
						"chromeOS",
						"androidPersonallyOwnedWorkProfile",
						"ios",
						"macOS",
						"windows",
						"windowsHolographic",
					),
				},
			},
			"device_inactivity_before_retirement_in_days": schema.Int32Attribute{
				Required:            true,
				MarkdownDescription: "Number of days when the device has not contacted Intune before it gets automatically removed. Valid values are 30 to 270.",
				Validators: []validator.Int32{
					int32validator.Between(30, 270),
				},
			},
			"last_modified_date_time": schema.StringAttribute{
				MarkdownDescription: "The date and time when the device cleanup rule was last modified. This property is read-only.",
				Computed:            true,
			},
			"timeouts": commonschema.Timeouts(ctx),
		},
	}
}
