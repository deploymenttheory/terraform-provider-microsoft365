package graphBetaWindowsAutopatchDeviceRegistration

import (
	"context"
	"regexp"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	planmodifiers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/plan_modifiers"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/identityschema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName  = "microsoft365_graph_beta_device_management_windows_autopatch_device_registration"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	_ resource.Resource                = &WindowsAutopatchDeviceRegistrationResource{}
	_ resource.ResourceWithConfigure   = &WindowsAutopatchDeviceRegistrationResource{}
	_ resource.ResourceWithImportState = &WindowsAutopatchDeviceRegistrationResource{}
	_ resource.ResourceWithIdentity    = &WindowsAutopatchDeviceRegistrationResource{}
)

func NewWindowsAutopatchDeviceRegistrationResource() resource.Resource {
	return &WindowsAutopatchDeviceRegistrationResource{
		ReadPermissions: []string{
			"WindowsUpdates.ReadWrite.All",
		},
		WritePermissions: []string{
			"WindowsUpdates.ReadWrite.All",
		},
		ResourcePath: "/admin/windows/updates/updatableAssets",
	}
}

type WindowsAutopatchDeviceRegistrationResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

func (r *WindowsAutopatchDeviceRegistrationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

func (r *WindowsAutopatchDeviceRegistrationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

func (r *WindowsAutopatchDeviceRegistrationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *WindowsAutopatchDeviceRegistrationResource) IdentitySchema(ctx context.Context, req resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	resp.IdentitySchema = identityschema.Schema{
		Attributes: map[string]identityschema.Attribute{
			"id": identityschema.StringAttribute{
				RequiredForImport: true,
			},
		},
	}
}

func (r *WindowsAutopatchDeviceRegistrationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages batch enrollment of Entra ID devices in Windows Autopatch update management using the `/admin/windows/updates/updatableAssets/enrollAssetsById` endpoint. " +
			"This resource enrolls multiple devices for a specific update category (driver, feature, or quality updates). " +
			"Devices are automatically created as azureADDevice resources in Windows Autopatch if they don't already exist. " +
			"When the resource is deleted, devices are unenrolled from the specified update category. " +
			"See the [Microsoft Graph API documentation](https://learn.microsoft.com/en-us/graph/api/windowsupdates-updatableasset-enrollassetsbyid?view=graph-rest-beta) for more information.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The composite identifier for this resource. Format: `{update_category}`.",
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
			},
			"update_category": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The category of updates for Windows Autopatch to manage. Valid values are: `driver`, `feature`, `quality`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf("driver", "feature", "quality"),
				},
			},
			"device_ids": schema.SetAttribute{
				ElementType:         types.StringType,
				Required:            true,
				MarkdownDescription: "Set of Entra ID device IDs to enroll in Windows Autopatch for the specified update category. These should be valid Entra ID device object IDs.",
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
					setvalidator.ValueStringsAre(
						stringvalidator.RegexMatches(regexp.MustCompile(constants.GuidRegex), "Must be a valid UUID format (xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx)"),
					),
				},
			},
			"timeouts": commonschema.ResourceTimeouts(ctx),
		},
	}
}
