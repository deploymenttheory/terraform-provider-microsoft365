package graphBetaAppControlForBusinessManagedInstaller

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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName  = "graph_beta_device_management_app_control_for_business_managed_installer"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &AppControlForBusinessManagedInstallerResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &AppControlForBusinessManagedInstallerResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &AppControlForBusinessManagedInstallerResource{}

	// Enables plan modification/diff suppression
	_ resource.ResourceWithModifyPlan = &AppControlForBusinessManagedInstallerResource{}
)

func NewAppControlForBusinessManagedInstallerResource() resource.Resource {
	return &AppControlForBusinessManagedInstallerResource{
		ReadPermissions: []string{
			"DeviceManagementConfiguration.Read.All",
			"DeviceManagementManagedDevices.Read.All",
		},
		WritePermissions: []string{
			"DeviceManagementConfiguration.ReadWrite.All",
			"DeviceManagementManagedDevices.ReadWrite.All",
		},
		ResourcePath: "/deviceAppManagement/windowsManagementApp",
	}
}

type AppControlForBusinessManagedInstallerResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ProviderTypeName string
	TypeName         string
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

// Metadata returns the resource type name.
func (r *AppControlForBusinessManagedInstallerResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	r.ProviderTypeName = req.ProviderTypeName
	r.TypeName = ResourceName
	resp.TypeName = r.FullTypeName()
}

// FullTypeName returns the full resource type name in the format "providername_resourcename".
func (r *AppControlForBusinessManagedInstallerResource) FullTypeName() string {
	return r.ProviderTypeName + "_" + r.TypeName
}

// Configure sets the client for the resource.
func (r *AppControlForBusinessManagedInstallerResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, constants.PROVIDER_NAME+"_"+ResourceName)
}

// ImportState imports the resource state.
func (r *AppControlForBusinessManagedInstallerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Schema returns the schema for the resource.
func (r *AppControlForBusinessManagedInstallerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages the Intune Management Extension as a Managed Installer for App Control for Business. This resource configures whether the Intune Management Extension should be trusted as a managed installer, allowing it to install applications that would otherwise be blocked by App Control for Business policies.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
				MarkdownDescription: "The unique identifier for the Windows Management App",
			},
			"intune_management_extension_as_managed_installer": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.OneOf("Enabled", "Disabled"),
				},
				MarkdownDescription: "Valid values are 'Enabled' or 'Disabled'. 'Enabled' will grant Microsoft permission to configure " +
					"Intune Management Extension as a managed installer (an authourized source for application deployment) on applicable devices. " +
					"'Disabled' will pause any scheduled policy to set the Intune Management Extension as a managed installer on applicable devices. " +
					"Existing policies already deployed to devices will not be changed when configuring to 'Disabled'. If removal of existing policies " +
					"is required, a cleanup script may be used. Learn more about how Intune sets the managed installer at 'https://learn.microsoft.com/en-us/intune/intune-service/protect/endpoint-security-app-control-policy'.",
			},
			"available_version": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The available version of the Windows Management App.",
			},
			"managed_installer_configured_date_time": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The date and time when the managed installer was last configured.",
			},
			"timeouts": commonschema.Timeouts(ctx),
		},
	}
}
