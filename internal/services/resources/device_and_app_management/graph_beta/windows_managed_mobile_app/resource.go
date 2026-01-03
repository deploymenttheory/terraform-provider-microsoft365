package graphBetaDeviceAndAppManagementWindowsManagedMobileApp

import (
	"context"
	"fmt"
	"strings"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	planmodifiers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/plan_modifiers"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName  = "microsoft365_graph_beta_device_and_app_management_windows_managed_mobile_app"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	_ resource.Resource                = &WindowsManagedMobileAppResource{}
	_ resource.ResourceWithConfigure   = &WindowsManagedMobileAppResource{}
	_ resource.ResourceWithImportState = &WindowsManagedMobileAppResource{}
	_ resource.ResourceWithModifyPlan  = &WindowsManagedMobileAppResource{}
)

func NewWindowsManagedMobileAppResource() resource.Resource {
	return &WindowsManagedMobileAppResource{
		ReadPermissions: []string{
			"DeviceManagementApps.Read.All",
		},
		WritePermissions: []string{
			"DeviceManagementApps.ReadWrite.All",
		},
		ResourcePath: "/deviceAppManagement/windowsManagedAppProtections/{id}/apps",
	}
}

type WindowsManagedMobileAppResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

func (r *WindowsManagedMobileAppResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

func (r *WindowsManagedMobileAppResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

func (r *WindowsManagedMobileAppResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// The import ID is expected to be in the format "managedAppProtectionId/appId"
	idParts := strings.Split(req.ID, "/")
	if len(idParts) != 2 {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			fmt.Sprintf("Expected import ID in format 'managedAppProtectionId/appId', got: %s", req.ID),
		)
		return
	}

	managedAppProtectionId := idParts[0]
	appId := idParts[1]

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("managed_app_protection_id"), managedAppProtectionId)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), appId)...)
}

func (r *WindowsManagedMobileAppResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages Windows managed mobile apps in Microsoft Intune Windows managed app protection policies. This resource associates Windows apps with Windows managed app protection policies for Mobile Application Management (MAM).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
				MarkdownDescription: "The unique identifier for this Windows managed mobile app",
			},
			"version": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				MarkdownDescription: "The version of the managed mobile app",
			},
			"managed_app_protection_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				MarkdownDescription: "The ID of the Windows managed app protection policy to associate this app with",
			},
			"mobile_app_identifier": schema.SingleNestedAttribute{
				Required:            true,
				MarkdownDescription: "The Windows app identifier information",
				Attributes: map[string]schema.Attribute{
					"windows_app_id": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "The Windows app identifier (e.g., Microsoft.WindowsApp_8wekyb3d8bbwe)",
					},
				},
			},
			"timeouts": commonschema.ResourceTimeouts(ctx),
		},
	}
}
