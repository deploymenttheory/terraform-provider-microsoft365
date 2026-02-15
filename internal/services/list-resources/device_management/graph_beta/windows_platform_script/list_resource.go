package graphBetaWindowsPlatformScript

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/list"
	listschema "github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ListResourceName = "microsoft365_graph_beta_device_management_windows_platform_script"
)

var (
	_ list.ListResource              = &WindowsPlatformScriptListResource{}
	_ list.ListResourceWithConfigure = &WindowsPlatformScriptListResource{}
)

func NewWindowsPlatformScriptListResource() list.ListResource {
	return &WindowsPlatformScriptListResource{
		ReadPermissions: []string{
			"DeviceManagementConfiguration.Read.All",
		},
		ResourcePath: "/deviceManagement/deviceManagementScripts",
	}
}

type WindowsPlatformScriptListResource struct {
	client          *msgraphbetasdk.GraphServiceClient
	ReadPermissions []string
	ResourcePath    string
}

// Metadata returns the list resource type name.
func (r *WindowsPlatformScriptListResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ListResourceName
}

// Configure sets the client for the list resource.
func (r *WindowsPlatformScriptListResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ListResourceName)
	if r.client != nil {
		tflog.Debug(ctx, "Successfully configured list resource client", map[string]any{
			"list_resource": ListResourceName,
		})
	} else {
		tflog.Error(ctx, "Failed to configure list resource client - client is nil", map[string]any{
			"list_resource": ListResourceName,
		})
	}
}

// ListResourceConfigSchema defines the schema for the list resource configuration.
func (r *WindowsPlatformScriptListResource) ListResourceConfigSchema(ctx context.Context, req list.ListResourceSchemaRequest, resp *list.ListResourceSchemaResponse) {
	resp.Schema = listschema.Schema{
		MarkdownDescription: "Lists Windows PowerShell scripts from Microsoft Intune using the `/deviceManagement/deviceManagementScripts` endpoint. " +
			"This list resource is used to automatically retrieve all scripts across multiple pages with advanced filtering capabilities for script discovery and import. " +
			"For full resource details, use Terraform's import functionality with `terraform plan -generate-config-out`.",
		Attributes: map[string]listschema.Attribute{
			"display_name_filter": listschema.StringAttribute{
				MarkdownDescription: "Filter scripts by display name using partial matching. Supports the OData `contains` operator. " +
					"Example: `display_name_filter = \"Baseline\"` will match \"Windows Baseline Script\".",
				Optional: true,
			},
			"file_name_filter": listschema.StringAttribute{
				MarkdownDescription: "Filter scripts by file name using partial matching. Supports the OData `contains` operator. " +
					"Example: `file_name_filter = \"setup.ps1\"` will match scripts with \"setup.ps1\" in the filename.",
				Optional: true,
			},
			"run_as_account_filter": listschema.StringAttribute{
				MarkdownDescription: "Filter scripts by execution context. Valid values: `system`, `user`. " +
					"Example: `run_as_account_filter = \"system\"` returns only scripts running as system.",
				Optional: true,
				Validators: []validator.String{
					stringvalidator.OneOf("system", "user"),
				},
			},
			"is_assigned_filter": listschema.BoolAttribute{
				MarkdownDescription: "Filter scripts by assignment status. Set to `true` to return only scripts with " +
					"assignments, `false` for scripts without assignments. This filter queries the assignments endpoint " +
					"for each script (the API's `isAssigned` field is unreliable) and may take 20-30 seconds for large tenants.",
				Optional: true,
			},
			"odata_filter": listschema.StringAttribute{
				MarkdownDescription: "Advanced: Custom OData $filter query for complex filtering scenarios. " +
					"Allows direct control over the API filter expression. " +
					"Example: `odata_filter = \"runAsAccount eq 'system' and contains(displayName, 'Baseline')\"`. " +
					"When specified, this overrides individual filter parameters. " +
					"See Microsoft Graph API documentation for supported operators and syntax.",
				Optional: true,
			},
		},
	}
}
