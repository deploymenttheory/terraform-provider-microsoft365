package graphBetaDeviceAndAppManagementAndroidManagedMobileApp

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
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
	ResourceName  = "graph_beta_device_and_app_management_android_managed_mobile_app"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	_ resource.Resource                = &AndroidManagedMobileAppResource{}
	_ resource.ResourceWithConfigure   = &AndroidManagedMobileAppResource{}
	_ resource.ResourceWithImportState = &AndroidManagedMobileAppResource{}
	_ resource.ResourceWithModifyPlan  = &AndroidManagedMobileAppResource{}
)

func NewAndroidManagedMobileAppResource() resource.Resource {
	return &AndroidManagedMobileAppResource{
		ReadPermissions: []string{
			"DeviceManagementApps.Read.All",
		},
		WritePermissions: []string{
			"DeviceManagementApps.ReadWrite.All",
		},
		ResourcePath: "/deviceAppManagement/androidManagedAppProtections/{id}/apps",
	}
}

type AndroidManagedMobileAppResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ProviderTypeName string
	TypeName         string
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

func (r *AndroidManagedMobileAppResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	r.ProviderTypeName = req.ProviderTypeName
	r.TypeName = ResourceName
	resp.TypeName = r.FullTypeName()
}

func (r *AndroidManagedMobileAppResource) FullTypeName() string {
	return r.ProviderTypeName + "_" + r.TypeName
}

func (r *AndroidManagedMobileAppResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, constants.PROVIDER_NAME+"_"+ResourceName)
}

func (r *AndroidManagedMobileAppResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *AndroidManagedMobileAppResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages Android managed mobile apps in Microsoft Intune Android managed app protection policies. This resource associates Android apps with Android managed app protection policies for Mobile Application Management (MAM).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
				MarkdownDescription: "The unique identifier for this Android managed mobile app",
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
				MarkdownDescription: "The ID of the Android managed app protection policy to associate this app with",
			},
			"mobile_app_identifier": schema.SingleNestedAttribute{
				Required: true,
				MarkdownDescription: "The Android app identifier information",
				Attributes: map[string]schema.Attribute{
					"package_id": schema.StringAttribute{
						Required: true,
						MarkdownDescription: "The Android package identifier (e.g., com.company.myapp)",
					},
				},
			},
			"timeouts": commonschema.Timeouts(ctx),
		},
	}
}