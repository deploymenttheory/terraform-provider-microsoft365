package graphBetaWindowsAutopilotDevicePreparationPolicy

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	"github.com/hashicorp/terraform-plugin-framework-validators/boolvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
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
	_ datasource.DataSource              = &WindowsAutopilotDevicePreparationPolicyDataSource{}
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

func (d *WindowsAutopilotDevicePreparationPolicyDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = DataSourceName
}

func (d *WindowsAutopilotDevicePreparationPolicyDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = client.SetGraphBetaClientForDataSource(ctx, req, resp, DataSourceName)
}

func (d *WindowsAutopilotDevicePreparationPolicyDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieves Windows Autopilot device preparation policies from Microsoft Intune using the `/deviceManagement/configurationPolicies` endpoint. " +
			"Device preparation policies are settings catalog policies created from the Autopilot device preparation templates, so results are always restricted to policies " +
			"whose `templateReference` matches one of those templates. Supports flexible lookup by policy ID, name, custom OData queries, or listing all policies.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The unique identifier for the data source. This is a placeholder attribute required by Terraform.",
			},
			"policy_id": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The unique identifier of the device preparation policy. Conflicts with other lookup attributes.",
				Validators: []validator.String{
					stringvalidator.ConflictsWith(
						path.MatchRoot("name"),
						path.MatchRoot("list_all"),
						path.MatchRoot("odata_query"),
					),
					stringvalidator.AtLeastOneOf(
						path.MatchRoot("policy_id"),
						path.MatchRoot("name"),
						path.MatchRoot("list_all"),
						path.MatchRoot("odata_query"),
					),
				},
			},
			"name": schema.StringAttribute{
				Optional: true,
				MarkdownDescription: "The name of the device preparation policy. The lookup is an exact match on the policy name. " +
					"Conflicts with other lookup attributes.",
				Validators: []validator.String{
					stringvalidator.ConflictsWith(
						path.MatchRoot("policy_id"),
						path.MatchRoot("list_all"),
						path.MatchRoot("odata_query"),
					),
				},
			},
			"list_all": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Retrieve all Windows Autopilot device preparation policies in the tenant. Conflicts with specific lookup attributes.",
				Validators: []validator.Bool{
					boolvalidator.ConflictsWith(
						path.MatchRoot("policy_id"),
						path.MatchRoot("name"),
						path.MatchRoot("odata_query"),
					),
				},
			},
			"list_assignments": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "When true, retrieves the assignments of the matched policy. When several policies match, the assignments of the first policy are returned.",
			},
			"odata_query": schema.StringAttribute{
				Optional: true,
				MarkdownDescription: "Custom OData filter expression for advanced queries (e.g., `name eq 'Autopilot policy' and isAssigned eq true`). " +
					"Results are still restricted to device preparation policies. Conflicts with specific lookup attributes.",
				Validators: []validator.String{
					stringvalidator.ConflictsWith(
						path.MatchRoot("policy_id"),
						path.MatchRoot("name"),
						path.MatchRoot("list_all"),
					),
				},
			},
			"items": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "List of Windows Autopilot device preparation policies matching the query criteria.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The unique identifier of the device preparation policy. This is the value required by `autopilot_configuration.device_preparation_profile_id` on a Windows 365 Cloud PC provisioning policy.",
						},
						"name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The name of the device preparation policy.",
						},
						"description": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The description of the device preparation policy.",
						},
						"created_date_time": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The date and time the policy was created.",
						},
						"last_modified_date_time": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The date and time the policy was last modified.",
						},
						"creation_source": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The source that created the policy.",
						},
						"platforms": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The platforms the policy applies to.",
						},
						"technologies": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The technologies used to deliver the policy.",
						},
						"setting_count": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "The number of settings configured on the policy.",
						},
						"is_assigned": schema.BoolAttribute{
							Computed:            true,
							MarkdownDescription: "Indicates whether the policy is assigned to any group.",
						},
						"disable_entra_group_policy_assignment": schema.BoolAttribute{
							Computed:            true,
							MarkdownDescription: "Indicates whether Microsoft Entra group policy assignment is disabled for the policy.",
						},
						"priority": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "The priority of the policy, taken from its priority metadata.",
						},
						"role_scope_tag_ids": schema.ListAttribute{
							ElementType:         types.StringType,
							Computed:            true,
							MarkdownDescription: "The list of scope tag IDs applied to the policy.",
						},
						"template_reference": schema.SingleNestedAttribute{
							Computed:            true,
							MarkdownDescription: "The settings catalog template the policy was created from.",
							Attributes: map[string]schema.Attribute{
								"template_id": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "The identifier of the template, including its version suffix.",
								},
								"template_family": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "The template family. Always `enrollmentConfiguration` for device preparation policies.",
								},
								"template_display_name": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "The display name of the template.",
								},
								"template_display_version": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "The display version of the template.",
								},
								"deployment_mode": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "The Autopilot device preparation deployment mode derived from the template ID. Either `automatic` or `user_driven`.",
								},
							},
						},
					},
				},
			},
			"assignments": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "Assignments of the matched device preparation policy. Only populated when `list_assignments` is true.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The unique identifier of the assignment.",
						},
						"type": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The assignment target type, for example `groupAssignmentTarget`.",
						},
						"group_id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The group the policy is assigned to, when the target is a group assignment target.",
						},
						"filter_id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The assignment filter applied to the assignment, if any.",
						},
						"filter_type": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The assignment filter mode. One of `include`, `exclude` or `none`.",
						},
					},
				},
			},
			"timeouts": commonschema.DatasourceTimeouts(ctx),
		},
	}
}
