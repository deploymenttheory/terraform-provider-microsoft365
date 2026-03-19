package graphBetaWindowsUpdatesAutopatchDeviceRegistration

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
	ResourceName  = "microsoft365_graph_beta_windows_updates_autopatch_device_registration"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &WindowsUpdatesAutopatchDeviceRegistrationResource{}
	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &WindowsUpdatesAutopatchDeviceRegistrationResource{}
	// Enables import functionality
	_ resource.ResourceWithImportState = &WindowsUpdatesAutopatchDeviceRegistrationResource{}
	// Enables identity schema for list resource support
	_ resource.ResourceWithIdentity = &WindowsUpdatesAutopatchDeviceRegistrationResource{}
)

func NewWindowsUpdatesAutopatchDeviceRegistrationResource() resource.Resource {
	return &WindowsUpdatesAutopatchDeviceRegistrationResource{
		ReadPermissions: []string{
			"WindowsUpdates.ReadWrite.All",
		},
		WritePermissions: []string{
			"WindowsUpdates.ReadWrite.All",
		},
		ResourcePath: "/admin/windows/updates/updatableAssets",
	}
}

type WindowsUpdatesAutopatchDeviceRegistrationResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

func (r *WindowsUpdatesAutopatchDeviceRegistrationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

func (r *WindowsUpdatesAutopatchDeviceRegistrationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

func (r *WindowsUpdatesAutopatchDeviceRegistrationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *WindowsUpdatesAutopatchDeviceRegistrationResource) IdentitySchema(ctx context.Context, req resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	resp.IdentitySchema = identityschema.Schema{
		Attributes: map[string]identityschema.Attribute{
			"id": identityschema.StringAttribute{
				RequiredForImport: true,
			},
		},
	}
}

func (r *WindowsUpdatesAutopatchDeviceRegistrationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages batch enrollment of Entra ID devices in Windows Autopatch update management using the `/admin/windows/updates/updatableAssets/enrollAssetsById` endpoint. " +
			"This resource enrolls multiple devices for a specific update category (driver, feature, or quality updates). " +
			"Devices are automatically created as azureADDevice resources in Windows Autopatch if they don't already exist. " +
			"When the resource is deleted, devices are unenrolled from the specified update category. ",
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
			"entra_device_object_ids": schema.SetAttribute{
				ElementType:         types.StringType,
				Required:            true,
				MarkdownDescription: "Set of Entra ID device object IDs to enroll in Windows Autopatch for the specified update category. These are the 'id' property from the Entra ID device object, NOT the 'deviceId' property.",
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
