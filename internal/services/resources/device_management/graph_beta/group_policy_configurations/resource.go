package graphBetaGroupPolicyConfigurations

import (
	"context"
	"regexp"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	commonschemagraphbeta "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema/graph_beta/device_management"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName  = "graph_beta_device_management_group_policy_configuration"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &GroupPolicyConfigurationResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &GroupPolicyConfigurationResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &GroupPolicyConfigurationResource{}
)

func NewGroupPolicyConfigurationResource() resource.Resource {
	return &GroupPolicyConfigurationResource{
		ReadPermissions: []string{
			"DeviceManagementConfiguration.Read.All",
		},
		WritePermissions: []string{
			"DeviceManagementConfiguration.ReadWrite.All",
		},
		ResourcePath: "/deviceManagement/groupPolicyConfigurations",
	}
}

type GroupPolicyConfigurationResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ProviderTypeName string
	TypeName         string
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

// Metadata returns the resource type name.
func (r *GroupPolicyConfigurationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	r.ProviderTypeName = req.ProviderTypeName
	r.TypeName = ResourceName
	resp.TypeName = r.FullTypeName()
}

// FullTypeName returns the full resource type name in the format "providername_resourcename".
func (r *GroupPolicyConfigurationResource) FullTypeName() string {
	return r.ProviderTypeName + "_" + r.TypeName
}

// Configure sets the client for the resource.
func (r *GroupPolicyConfigurationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, constants.PROVIDER_NAME+"_"+ResourceName)
}

// ImportState imports the resource state.
func (r *GroupPolicyConfigurationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Schema defines the schema for the resource.
func (r *GroupPolicyConfigurationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages Microsoft 365 Group Policy Configurations using the `/deviceManagement/groupPolicyConfigurations` endpoint. Group Policy Configurations define policy settings that can be applied to managed devices.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "String (identifier)",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(constants.GuidRegex),
						"must be a valid GUID",
					),
				},
			},
			"display_name": schema.StringAttribute{
				MarkdownDescription: "The display name for the Group Policy Configuration.",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description for the Group Policy Configuration.",
				Optional:            true,
			},
			"created_date_time": schema.StringAttribute{
				MarkdownDescription: "The creation date and time of the group policy configuration.",
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(constants.TimeFormatRFC3339Regex),
						"must be a valid RFC3339 date-time string",
					),
				},
			},
			"last_modified_date_time": schema.StringAttribute{
				MarkdownDescription: "The last modified date and time of the group policy configuration.",
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(constants.TimeFormatRFC3339Regex),
						"must be a valid RFC3339 date-time string",
					),
				},
			},
			"role_scope_tag_ids": schema.SetAttribute{
				MarkdownDescription: "List of role scope tag IDs for this configuration.",
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
			},
			"assignments": commonschemagraphbeta.DeviceConfigurationWithAllGroupAssignmentsAndFilterSchema(),
			"timeouts":    commonschema.Timeouts(ctx),
		},
		Blocks: map[string]schema.Block{
			"definition_values": schema.SetNestedBlock{
				MarkdownDescription: "Set of policy definition values that define the configuration settings.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "The ID of the definition value.",
							Computed:            true,
							Validators: []validator.String{
								stringvalidator.RegexMatches(
									regexp.MustCompile(constants.GuidRegex),
									"must be a valid GUID",
								),
							},
						},
						"enabled": schema.BoolAttribute{
							MarkdownDescription: "Whether this policy definition is enabled.",
							Required:            true,
						},
						"configuration_type": schema.StringAttribute{
							MarkdownDescription: "The configuration type (e.g., 'policy').",
							Computed:            true,
						},
						"created_date_time": schema.StringAttribute{
							MarkdownDescription: "The creation date and time of the definition value.",
							Computed:            true,
							Validators: []validator.String{
								stringvalidator.RegexMatches(
									regexp.MustCompile(constants.TimeFormatRFC3339Regex),
									"must be a valid RFC3339 date-time string",
								),
							},
						},
						"last_modified_date_time": schema.StringAttribute{
							MarkdownDescription: "The last modified date and time of the definition value.",
							Computed:            true,
							Validators: []validator.String{
								stringvalidator.RegexMatches(
									regexp.MustCompile(constants.TimeFormatRFC3339Regex),
									"must be a valid RFC3339 date-time string",
								),
							},
						},
						"display_name": schema.StringAttribute{
							MarkdownDescription: "The display name of the group policy definition (e.g., 'Allow users to contact Microsoft for feedback and support').",
							Required:            true,
						},
						"class_type": schema.StringAttribute{
							MarkdownDescription: "The class type of the group policy definition. Valid values are 'user' or 'machine'.",
							Required:            true,
							Validators: []validator.String{
								stringvalidator.OneOf("user", "machine"),
							},
						},
						"category_path": schema.StringAttribute{
							MarkdownDescription: "Optional category path for disambiguation when multiple definitions have the same display name (e.g., '\\\\OneDrive').",
							Optional:            true,
						},
						"definition_id": schema.StringAttribute{
							MarkdownDescription: "The ID of the group policy definition this value is for (computed from display_name and class_type).",
							Computed:            true,
							Validators: []validator.String{
								stringvalidator.RegexMatches(
									regexp.MustCompile(constants.GuidRegex),
									"must be a valid GUID",
								),
							},
						},
					},
					Blocks: map[string]schema.Block{
						"presentation_values": schema.SetNestedBlock{
							MarkdownDescription: "Set of presentation values that contain the actual configuration values.",
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										MarkdownDescription: "The ID of the presentation value.",
										Computed:            true,
										Validators: []validator.String{
											stringvalidator.RegexMatches(
												regexp.MustCompile(constants.GuidRegex),
												"must be a valid GUID",
											),
										},
									},
									"created_date_time": schema.StringAttribute{
										MarkdownDescription: "The creation date and time of the presentation value.",
										Computed:            true,
										Validators: []validator.String{
											stringvalidator.RegexMatches(
												regexp.MustCompile(constants.TimeFormatRFC3339Regex),
												"must be a valid RFC3339 date-time string",
											),
										},
									},
									"last_modified_date_time": schema.StringAttribute{
										MarkdownDescription: "The last modified date and time of the presentation value.",
										Computed:            true,
										Validators: []validator.String{
											stringvalidator.RegexMatches(
												regexp.MustCompile(constants.TimeFormatRFC3339Regex),
												"must be a valid RFC3339 date-time string",
											),
										},
									},
									"presentation_id": schema.StringAttribute{
										MarkdownDescription: "The ID of the group policy presentation this value is for.",
										Required:            true,
										Validators: []validator.String{
											stringvalidator.RegexMatches(
												regexp.MustCompile(constants.GuidRegex),
												"must be a valid GUID",
											),
										},
									},
									"odata_type": schema.StringAttribute{
										MarkdownDescription: "The OData type that determines the type of presentation value. Supported types: '#microsoft.graph.groupPolicyPresentationValueText', '#microsoft.graph.groupPolicyPresentationValueDecimal', '#microsoft.graph.groupPolicyPresentationValueBoolean', '#microsoft.graph.groupPolicyPresentationValueList', '#microsoft.graph.groupPolicyPresentationValueMultiText'",
										Required:            true,
										Validators: []validator.String{
											stringvalidator.OneOf(
												"#microsoft.graph.groupPolicyPresentationValueText",
												"#microsoft.graph.groupPolicyPresentationValueDecimal",
												"#microsoft.graph.groupPolicyPresentationValueBoolean",
												"#microsoft.graph.groupPolicyPresentationValueList",
												"#microsoft.graph.groupPolicyPresentationValueMultiText",
												"#microsoft.graph.groupPolicyPresentationValueDropdownList",
												"#microsoft.graph.groupPolicyPresentationValueLongDecimal",
											),
										},
									},
									"value": schema.StringAttribute{
										MarkdownDescription: "Generic value field. The interpretation depends on the odata_type.",
										Optional:            true,
									},
									"text_value": schema.StringAttribute{
										MarkdownDescription: "Text value (used when odata_type is '#microsoft.graph.groupPolicyPresentationValueText').",
										Optional:            true,
									},
									"decimal_value": schema.Int64Attribute{
										MarkdownDescription: "Decimal value (used when odata_type is '#microsoft.graph.groupPolicyPresentationValueDecimal' or '#microsoft.graph.groupPolicyPresentationValueLongDecimal').",
										Optional:            true,
									},
									"boolean_value": schema.BoolAttribute{
										MarkdownDescription: "Boolean value (used when odata_type is '#microsoft.graph.groupPolicyPresentationValueBoolean').",
										Optional:            true,
									},
									"list_values": schema.SetAttribute{
										MarkdownDescription: "Set of list values (used when odata_type is '#microsoft.graph.groupPolicyPresentationValueList').",
										ElementType:         types.StringType,
										Optional:            true,
									},
									"multi_text_values": schema.SetAttribute{
										MarkdownDescription: "Set of multi-text values (used when odata_type is '#microsoft.graph.groupPolicyPresentationValueMultiText').",
										ElementType:         types.StringType,
										Optional:            true,
									},
								},
							},
							Validators: []validator.Set{
								setvalidator.SizeAtLeast(0),
							},
						},
					},
				},
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
				},
			},
		},
	}
}
