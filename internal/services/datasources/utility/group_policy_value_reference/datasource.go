package utilityGroupPolicyValueReference

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	DataSourceName = "microsoft365_utility_group_policy_value_reference"
	ReadTimeout    = 180
)

var (
	_ datasource.DataSource              = &groupPolicyValueReferenceDataSource{}
	_ datasource.DataSourceWithConfigure = &groupPolicyValueReferenceDataSource{}
)

func NewGroupPolicyValueReferenceDataSource() datasource.DataSource {
	return &groupPolicyValueReferenceDataSource{
		ReadPermissions: []string{
			"DeviceManagementConfiguration.Read.All",
		},
	}
}

type groupPolicyValueReferenceDataSource struct {
	client          *msgraphbetasdk.GraphServiceClient
	ReadPermissions []string
}

func (d *groupPolicyValueReferenceDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = DataSourceName
}

func (d *groupPolicyValueReferenceDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = client.SetGraphBetaClientForDataSource(ctx, req, resp, DataSourceName)
}

func (d *groupPolicyValueReferenceDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Queries Microsoft Graph API for group policy definition metadata. " +
			"This utility data source retrieves detailed information about group policy definitions including " +
			"class type, category path, and presentation (checkbox/setting) details based on a policy display name.\n\n" +
			"Use this data source to discover the exact metadata needed for configuring group policy boolean values, " +
			"text values, list values, and other group policy settings in Microsoft Intune.\n\n" +
			"**Search Behavior:** Requires an exact match (case-insensitive, whitespace-normalized). " +
			"If no exact match is found, returns a helpful error message listing similar policy names ranked by similarity.\n\n" +
			"**Key Features:**\n\n" +
			"- Returns all definitions matching the policy name\n" +
			"- Provides class_type (`user` or `machine`) for each definition\n" +
			"- Shows the full category_path for policy organization\n" +
			"- Lists all presentations (individual settings) available for the policy\n" +
			"- Returns presentation types (checkbox, text, list, etc.) and their template IDs\n" +
			"- Provides helpful suggestions when policy name doesn't match exactly\n\n" +
			"**Common Use Cases:**\n\n" +
			"- Discovering the correct `class_type` for a policy\n" +
			"- Finding the exact `category_path` string\n" +
			"- Identifying which presentations are available for boolean configuration\n" +
			"- Distinguishing between multiple policies with similar names\n\n" +
			"**Reference:** [Group Policy Definitions API](https://learn.microsoft.com/en-us/graph/api/intune-grouppolicy-grouppolicydefinition-get?view=graph-rest-beta)",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The ID of this data source operation.",
			},
			"policy_name": schema.StringAttribute{
				Required: true,
				MarkdownDescription: "The display name of the group policy definition to search for. " +
					"Requires an exact match (case-insensitive, whitespace-normalized). " +
					"If no exact match is found, the error message will suggest similar policy names ranked by similarity using fuzzy matching. " +
					"Example: `\"Allow Cloud Policy Management\"`, `\"Enable Profile Containers\"`.",
			},
			"definitions": schema.ListNestedAttribute{
				Computed: true,
				MarkdownDescription: "List of group policy definitions matching the policy name. " +
					"Multiple definitions may be returned if the same policy name exists in different categories or for different class types.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The unique identifier (GUID) of the group policy definition template.",
						},
						"display_name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The display name of the group policy definition.",
						},
						"class_type": schema.StringAttribute{
							Computed: true,
							MarkdownDescription: "The class type of the policy. " +
								"Valid values: `user` (applies to user settings), `machine` (applies to computer settings).",
						},
						"category_path": schema.StringAttribute{
							Computed: true,
							MarkdownDescription: "The full category path of the policy in the Group Policy hierarchy. " +
								"Format: `\\Category\\Subcategory`. " +
								"Example: `\\FSLogix\\Profile Containers`, `\\Microsoft Edge\\Content settings`.",
						},
						"explain_text": schema.StringAttribute{
							Computed: true,
							MarkdownDescription: "The detailed explanation text describing what the policy does. " +
								"This is the same text shown in the Group Policy Management Console.",
						},
						"supported_on": schema.StringAttribute{
							Computed: true,
							MarkdownDescription: "The supported platforms/versions for this policy. " +
								"Example: `\"At least Windows 10 Server, Windows 10 or Windows 10 RT\"`.",
						},
						"policy_type": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The type of policy (e.g., `admxBacked`, `admxIngested`).",
						},
						"presentations": schema.ListNestedAttribute{
							Computed: true,
							MarkdownDescription: "List of presentations (individual settings/controls) available for this policy definition. " +
								"Each presentation represents a configurable element like a checkbox, text box, dropdown, etc.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Computed:            true,
										MarkdownDescription: "The unique identifier (GUID) of the presentation template.",
									},
									"label": schema.StringAttribute{
										Computed: true,
										MarkdownDescription: "The label text displayed for this presentation in the UI. " +
											"Example: `\"Enable Profile Containers\"`, `\"Set maximum size\"`.",
									},
									"type": schema.StringAttribute{
										Computed: true,
										MarkdownDescription: "The OData type of the presentation, indicating the control type. " +
											"Examples: " +
											"`#microsoft.graph.groupPolicyPresentationCheckBox` (boolean on/off), " +
											"`#microsoft.graph.groupPolicyPresentationText` (text input), " +
											"`#microsoft.graph.groupPolicyPresentationDecimalTextBox` (numeric input), " +
											"`#microsoft.graph.groupPolicyPresentationDropdownList` (dropdown selection), " +
											"`#microsoft.graph.groupPolicyPresentationListBox` (list of values).",
									},
									"required": schema.BoolAttribute{
										Computed: true,
										MarkdownDescription: "Indicates whether this presentation is required when the policy is enabled. " +
											"If true, the presentation must have a value set.",
									},
								},
							},
						},
					},
				},
			},
			"timeouts": commonschema.DatasourceTimeouts(ctx),
		},
	}
}
