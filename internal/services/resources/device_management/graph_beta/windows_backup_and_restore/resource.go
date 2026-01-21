package graphBetaWindowsBackupAndRestore

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
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
	ResourceName  = "microsoft365_graph_beta_device_management_windows_backup_and_restore"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &WindowsBackupAndRestoreResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &WindowsBackupAndRestoreResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &WindowsBackupAndRestoreResource{}
)

func NewWindowsBackupAndRestoreResource() resource.Resource {
	return &WindowsBackupAndRestoreResource{
		ReadPermissions: []string{
			"DeviceManagementConfiguration.Read.All",
		},
		WritePermissions: []string{
			"DeviceManagementConfiguration.ReadWrite.All",
		},
		ResourcePath: "/deviceManagement/deviceEnrollmentConfigurations",
	}
}

type WindowsBackupAndRestoreResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

// Metadata returns the resource type name.
func (r *WindowsBackupAndRestoreResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

// Configure sets the client for the resource.
func (r *WindowsBackupAndRestoreResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

// ImportState imports the resource state.
func (r *WindowsBackupAndRestoreResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Schema returns the schema for the resource.
func (r *WindowsBackupAndRestoreResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages Windows Backup and Restore device enrollment configuration using the `/deviceManagement/deviceEnrollmentConfigurations` endpoint. This resource is used to configuration controls the Windows Restore feature for device enrollment. Learn more here: 'https://learn.microsoft.com/en-us/windows/configuration/windows-backup/?tabs=intune'..",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
				MarkdownDescription: "The unique identifier of the Windows Backup and Restore configuration.",
			},
			"state": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The enablement state of Windows Backup and Restore. Possible values are: 'enabled', 'disabled', 'notConfigured'.",
				Validators: []validator.String{
					stringvalidator.OneOf("enabled", "disabled", "notConfigured"),
				},
			},
			"timeouts": commonschema.ResourceTimeouts(ctx),
		},
	}
}
