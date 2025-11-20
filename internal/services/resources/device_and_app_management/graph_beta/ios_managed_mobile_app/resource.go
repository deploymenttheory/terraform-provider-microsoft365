package graphBetaDeviceAndAppManagementIOSManagedMobileApp

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
	ResourceName  = "microsoft365_graph_beta_device_and_app_management_ios_managed_mobile_app"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	_ resource.Resource                = &IOSManagedMobileAppResource{}
	_ resource.ResourceWithConfigure   = &IOSManagedMobileAppResource{}
	_ resource.ResourceWithImportState = &IOSManagedMobileAppResource{}
	_ resource.ResourceWithModifyPlan  = &IOSManagedMobileAppResource{}
)

func NewIOSManagedMobileAppResource() resource.Resource {
	return &IOSManagedMobileAppResource{
		ReadPermissions: []string{
			"DeviceManagementApps.Read.All",
		},
		WritePermissions: []string{
			"DeviceManagementApps.ReadWrite.All",
		},
		ResourcePath: "/deviceAppManagement/iosManagedAppProtections/{id}/apps",
	}
}

type IOSManagedMobileAppResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

func (r *IOSManagedMobileAppResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

func (r *IOSManagedMobileAppResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

func (r *IOSManagedMobileAppResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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

func (r *IOSManagedMobileAppResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages iOS managed mobile apps in Microsoft Intune iOS managed app protection policies. This resource associates iOS apps with iOS managed app protection policies for Mobile Application Management (MAM).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
				MarkdownDescription: "The unique identifier for this iOS managed mobile app",
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
				MarkdownDescription: "The ID of the iOS managed app protection policy to associate this app with",
			},
			"mobile_app_identifier": schema.SingleNestedAttribute{
				Required:            true,
				MarkdownDescription: "The iOS app identifier information",
				Attributes: map[string]schema.Attribute{
					"bundle_id": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "The iOS bundle identifier (e.g., com.company.myapp)",
					},
				},
			},
			"timeouts": commonschema.Timeouts(ctx),
		},
	}
}
