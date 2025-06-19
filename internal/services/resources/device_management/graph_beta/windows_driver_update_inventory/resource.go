package graphBetaWindowsDriverUpdateInventory

import (
	"context"
	"regexp"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common"
	planmodifiers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/plan_modifiers"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName  = "graph_beta_device_management_windows_driver_update_inventory"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &WindowsDriverUpdateInventoryResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &WindowsDriverUpdateInventoryResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &WindowsDriverUpdateInventoryResource{}

	// Enables plan modification/diff suppression
	_ resource.ResourceWithModifyPlan = &WindowsDriverUpdateInventoryResource{}
)

func NewWindowsDriverUpdateInventoryResource() resource.Resource {
	return &WindowsDriverUpdateInventoryResource{
		ReadPermissions: []string{
			"DeviceManagementConfiguration.Read.All",
		},
		WritePermissions: []string{
			"DeviceManagementConfiguration.ReadWrite.All",
		},
		ResourcePath: "POST /deviceManagement/windowsDriverUpdateProfiles/{windowsDriverUpdateProfileId}/driverInventories",
	}
}

type WindowsDriverUpdateInventoryResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ProviderTypeName string
	TypeName         string
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

// Metadata returns the resource type name.
func (r *WindowsDriverUpdateInventoryResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	r.ProviderTypeName = req.ProviderTypeName
	r.TypeName = ResourceName
	resp.TypeName = r.FullTypeName()
}

// FullTypeName returns the full type name of the resource for logging purposes.
func (r *WindowsDriverUpdateInventoryResource) FullTypeName() string {
	return r.ProviderTypeName + "_" + r.TypeName
}

// Configure sets the client for the resource.
func (r *WindowsDriverUpdateInventoryResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = common.SetGraphBetaClientForResource(ctx, req, resp, constants.PROVIDER_NAME+"_"+ResourceName)
}

// ImportState imports the resource state.
func (r *WindowsDriverUpdateInventoryResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Schema defines the schema for the resource.
func (r *WindowsDriverUpdateInventoryResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages Windows Driver Update Inventory in Microsoft Intune.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
				MarkdownDescription: "The id of the driver.",
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of the driver.",
			},
			"version": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The version of the driver.",
			},
			"manufacturer": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The manufacturer of the driver.",
			},
			"release_date_time": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The release date time of the driver.",
			},
			"driver_class": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The class of the driver.",
			},
			"applicable_device_count": schema.Int32Attribute{
				Computed:            true,
				MarkdownDescription: "The number of devices for which this driver is applicable.",
			},
			"approval_status": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The approval status for this driver. Possible values are: `needsReview`, `declined`, `approved`, `suspended`.",
				Validators: []validator.String{
					stringvalidator.OneOf("needsReview", "declined", "approved", "suspended"),
				},
			},
			"category": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The category for this driver. Possible values are: `recommended`, `previouslyApproved`, `other`.",
				Validators: []validator.String{
					stringvalidator.OneOf("recommended", "previouslyApproved", "other"),
				},
			},
			"deploy_date_time": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The date time when a driver should be deployed if approvalStatus is approved.",
			},
			"windows_driver_update_profile_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The ID of the Windows Driver Update Profile this inventory belongs to.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(constants.GuidRegex),
						"must be a valid GUID in the format 00000000-0000-0000-0000-000000000000",
					),
				},
			},
			"timeouts": commonschema.Timeouts(ctx),
		},
	}
}
