package graphBetaIdentityAndAccessLifecycleWorkflow

import (
	"context"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/identityschema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
)

const (
	ResourceName  = "microsoft365_graph_beta_identity_and_access_lifecycle_workflow"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &LifecycleWorkflowResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &LifecycleWorkflowResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &LifecycleWorkflowResource{}

	// Enables plan modification/diff suppression
	_ resource.ResourceWithModifyPlan = &LifecycleWorkflowResource{}

	// Enables cross-attribute validation of execution conditions
	_ resource.ResourceWithValidateConfig = &LifecycleWorkflowResource{}

	// Enables identity schema for list resource support
	_ resource.ResourceWithIdentity = &LifecycleWorkflowResource{}
)

// NewLifecycleWorkflowResource returns a new lifecycle workflow resource.
//
// LifecycleWorkflows-Workflow.ReadBasic.All is the least privileged read permission in Graph metadata,
// but it only allows listing workflows ("Allows the app to list all workflows"); reading tasks requires
// LifecycleWorkflows-Workflow.Read.All.
func NewLifecycleWorkflowResource() resource.Resource {
	return &LifecycleWorkflowResource{
		ReadPermissions: []string{
			"LifecycleWorkflows-Workflow.Read.All",
		},
		WritePermissions: []string{
			"LifecycleWorkflows-Workflow.ReadWrite.All",
		},
	}
}

type LifecycleWorkflowResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
}

// Metadata returns the resource type name.
func (r *LifecycleWorkflowResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

// Configure sets the client for the resource.
func (r *LifecycleWorkflowResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

// ImportState imports the resource state.
func (r *LifecycleWorkflowResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// IdentitySchema defines the identity schema for this resource, used by list operations to uniquely identify instances
func (r *LifecycleWorkflowResource) IdentitySchema(ctx context.Context, req resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	resp.IdentitySchema = identityschema.Schema{
		Attributes: map[string]identityschema.Attribute{
			"id": identityschema.StringAttribute{
				RequiredForImport: true,
			},
		},
	}
}

// Schema defines the schema for the resource.
func (r *LifecycleWorkflowResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages Microsoft Entra ID Governance lifecycle workflows using the `/identityGovernance/lifecycleWorkflows/workflows` endpoint. " +
			"Lifecycle workflows automate joiner, mover and leaver tasks for users. Changes to `display_name`, `description`, `is_enabled` and " +
			"`is_scheduling_enabled` are applied with PATCH; changes to `category`, `execution_conditions` or `tasks` create a new workflow version " +
			"in place (`createNewVersion`). On destroy the workflow is deleted and then permanently removed from the lifecycle workflows deleted items container.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier of the workflow.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"category": schema.StringAttribute{
				MarkdownDescription: "The category of the workflow. Possible values are: `joiner`, `mover`, `leaver`. Changing the category creates a new workflow version.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("joiner", "mover", "leaver"),
				},
			},
			"display_name": schema.StringAttribute{
				MarkdownDescription: "The display name of the workflow.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "A description of the workflow. Defaults to an empty string, which is what Graph returns when no description is set.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"is_enabled": schema.BoolAttribute{
				MarkdownDescription: "Whether the workflow is enabled. Defaults to `false`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"is_scheduling_enabled": schema.BoolAttribute{
				MarkdownDescription: "Whether the workflow runs on its schedule. When `true`, the workflow processes in-scope users every schedule interval. Defaults to `false`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"execution_conditions": schema.SingleNestedAttribute{
				MarkdownDescription: "The conditions that determine when and for whom the workflow runs. Set `on_demand_execution_only = true` for a workflow that only runs on demand, " +
					"or set both `scope` and `trigger` for a trigger and scope based workflow.",
				Required: true,
				Attributes: map[string]schema.Attribute{
					"on_demand_execution_only": schema.BoolAttribute{
						MarkdownDescription: "When `true`, the workflow can only be run on demand (`onDemandExecutionOnly`) and `scope` and `trigger` must not be set. Defaults to `false`.",
						Optional:            true,
						Computed:            true,
						Default:             booldefault.StaticBool(false),
					},
					"scope": schema.SingleNestedAttribute{
						MarkdownDescription: "The users the workflow applies to. Set exactly one of `rule` or `group_ids`.",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"rule": schema.StringAttribute{
								MarkdownDescription: "A filter rule that selects users (`ruleBasedSubjectSet`), for example `(department eq 'Marketing')`.",
								Optional:            true,
								Validators: []validator.String{
									stringvalidator.LengthAtLeast(1),
									stringvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("group_ids")),
								},
							},
							"group_ids": schema.SetAttribute{
								MarkdownDescription: "The IDs of the groups whose members are in scope (`groupBasedSubjectSet`). Required by a `membership_change` trigger.",
								ElementType:         types.StringType,
								Optional:            true,
								Validators: []validator.Set{
									setvalidator.SizeAtLeast(1),
									setvalidator.ValueStringsAre(
										stringvalidator.RegexMatches(regexp.MustCompile(constants.GuidRegex), "must be a valid GUID"),
									),
								},
							},
						},
					},
					"trigger": schema.SingleNestedAttribute{
						MarkdownDescription: "The event that starts the workflow. Set exactly one of `time_based_attribute`, `attribute_change` or `membership_change`.",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"time_based_attribute": schema.SingleNestedAttribute{
								MarkdownDescription: "Runs the workflow relative to a date attribute of the user (`timeBasedAttributeTrigger`).",
								Optional:            true,
								Validators: []validator.Object{
									objectvalidator.ExactlyOneOf(
										path.MatchRelative().AtParent().AtName("attribute_change"),
										path.MatchRelative().AtParent().AtName("membership_change"),
									),
								},
								Attributes: map[string]schema.Attribute{
									"attribute": schema.StringAttribute{
										MarkdownDescription: "The user date attribute. Possible values are: `employeeHireDate`, `employeeLeaveDateTime`, `createdDateTime`.",
										Required:            true,
										Validators: []validator.String{
											stringvalidator.OneOf("employeeHireDate", "employeeLeaveDateTime", "createdDateTime"),
										},
									},
									"offset_in_days": schema.Int64Attribute{
										MarkdownDescription: "The number of days before (negative) or after (positive) the attribute date to run the workflow. Must be between -180 and 180.",
										Required:            true,
										Validators: []validator.Int64{
											int64validator.Between(-180, 180),
										},
									},
								},
							},
							"attribute_change": schema.SingleNestedAttribute{
								MarkdownDescription: "Runs the workflow when a user attribute changes (`attributeChangeTrigger`). Typically used by `mover` workflows.",
								Optional:            true,
								Attributes: map[string]schema.Attribute{
									"trigger_attribute": schema.StringAttribute{
										MarkdownDescription: "The name of the user attribute whose change triggers the workflow, for example `department` or `jobTitle`. Graph accepts a single trigger attribute.",
										Required:            true,
										Validators: []validator.String{
											stringvalidator.LengthAtLeast(1),
										},
									},
								},
							},
							"membership_change": schema.SingleNestedAttribute{
								MarkdownDescription: "Runs the workflow when a user is added to or removed from a group in `scope.group_ids` (`membershipChangeTrigger`).",
								Optional:            true,
								Attributes: map[string]schema.Attribute{
									"change_type": schema.StringAttribute{
										MarkdownDescription: "The membership change that triggers the workflow. Possible values are: `add`, `remove`.",
										Required:            true,
										Validators: []validator.String{
											stringvalidator.OneOf("add", "remove"),
										},
									},
								},
							},
						},
					},
				},
			},
			"tasks": schema.ListNestedAttribute{
				MarkdownDescription: "The built-in tasks the workflow runs, in execution order. Custom task extensions are not supported by this resource. " +
					"Changing tasks creates a new workflow version.",
				Required: true,
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"task_definition_id": schema.StringAttribute{
							MarkdownDescription: "The ID of the built-in task definition, from `GET /identityGovernance/lifecycleWorkflows/taskDefinitions`. " +
								"For example `6fc52c9d-398b-4305-9763-15f42c1676fc` (Enable user account).",
							Required: true,
							Validators: []validator.String{
								stringvalidator.RegexMatches(regexp.MustCompile(constants.GuidRegex), "must be a valid GUID"),
							},
						},
						"display_name": schema.StringAttribute{
							MarkdownDescription: "The display name of the task.",
							Required:            true,
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(1),
							},
						},
						"description": schema.StringAttribute{
							MarkdownDescription: "A description of the task. Defaults to an empty string.",
							Optional:            true,
							Computed:            true,
							Default:             stringdefault.StaticString(""),
						},
						"is_enabled": schema.BoolAttribute{
							MarkdownDescription: "Whether the task runs when the workflow runs. Defaults to `true`.",
							Optional:            true,
							Computed:            true,
							Default:             booldefault.StaticBool(true),
						},
						"continue_on_error": schema.BoolAttribute{
							MarkdownDescription: "Whether the workflow continues with the next task when this task fails. Defaults to `false`.",
							Optional:            true,
							Computed:            true,
							Default:             booldefault.StaticBool(false),
						},
						"execution_sequence": schema.Int64Attribute{
							MarkdownDescription: "The position of the task in the workflow (1-based). Set from the order of `tasks`.",
							Computed:            true,
						},
						"arguments": schema.SetNestedAttribute{
							MarkdownDescription: "The arguments of the task, as defined by the task definition's parameters (for example `enableOnPremisesAccount` or `groupID`).",
							Optional:            true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"name": schema.StringAttribute{
										MarkdownDescription: "The argument name.",
										Required:            true,
									},
									"value": schema.StringAttribute{
										MarkdownDescription: "The argument value.",
										Required:            true,
									},
								},
							},
						},
					},
				},
			},
			"version": schema.Int64Attribute{
				MarkdownDescription: "The current version number of the workflow. Increases each time `category`, `execution_conditions` or `tasks` change.",
				Computed:            true,
			},
			"created_date_time": schema.StringAttribute{
				MarkdownDescription: "When the workflow was created.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"last_modified_date_time": schema.StringAttribute{
				MarkdownDescription: "When the workflow was last modified.",
				Computed:            true,
			},
			"timeouts": commonschema.ResourceTimeouts(ctx),
		},
	}
}
