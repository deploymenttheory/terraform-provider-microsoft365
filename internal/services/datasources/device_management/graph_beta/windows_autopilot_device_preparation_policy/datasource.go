package graphBetaWindowsAutopilotDevicePreparationPolicy

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	DataSourceName = "microsoft365_graph_beta_device_management_windows_autopilot_device_preparation_policy"
	ReadTimeout    = 180

	// Template family and base template ids of the Autopilot device preparation policy templates.
	// Used to ensure only device preparation policies are returned from the settings catalog collection.
	templateFamily           = "enrollmentConfiguration"
	baseTemplateIDAutomatic  = "a6157a7f-aa00-42d9-ac82-7d2479f545db"
	baseTemplateIDUserDriven = "80d33118-b7b4-40d8-b15f-81be745e053f"
)

var (
	// Basic resource interface (CRUD operations)
	_ datasource.DataSource = &WindowsAutopilotDevicePreparationPolicyDataSource{}

	// Allows the resource to be configured with the provider client
	_ datasource.DataSourceWithConfigure = &WindowsAutopilotDevicePreparationPolicyDataSource{}
)

func NewWindowsAutopilotDevicePreparationPolicyDataSource() datasource.DataSource {
	return &WindowsAutopilotDevicePreparationPolicyDataSource{
		ReadPermissions: []string{
			"DeviceManagementConfiguration.Read.All",
		},
	}
}

type WindowsAutopilotDevicePreparationPolicyDataSource struct {
	client *msgraphbetasdk.GraphServiceClient

	ReadPermissions []string
}

// Metadata returns the resource type name.
func (d *WindowsAutopilotDevicePreparationPolicyDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = DataSourceName
}

// Configure sets the client for the data source
func (d *WindowsAutopilotDevicePreparationPolicyDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = client.SetGraphBetaClientForDataSource(ctx, req, resp, DataSourceName)
}

// Schema defines the schema for the data source
func (d *WindowsAutopilotDevicePreparationPolicyDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieves a Windows Autopilot device preparation policy from Microsoft Intune using the `/deviceManagement/configurationPolicies` endpoint. " +
			"Device preparation policies are settings catalog policies created from the Autopilot device preparation templates. " +
			"This data source is typically used to resolve the policy id of a device preparation policy managed in a different Terraform state, " +
			"for example to populate `autopilot_configuration.device_preparation_profile_id` on a Windows 365 Cloud PC provisioning policy.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Unique identifier for the windows autopilot device preparation policy.",
				Optional:            true,
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the windows autopilot device preparation policy. The lookup is an exact match on the policy name.",
				Optional:            true,
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description of the windows autopilot device preparation policy.",
				Computed:            true,
			},
			"role_scope_tag_ids": schema.SetAttribute{
				MarkdownDescription: "List of Scope Tag IDs for this windows autopilot device preparation policy.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"timeouts": commonschema.ResourceTimeouts(ctx),
		},
	}
}
