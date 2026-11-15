package graphBetaRoleManagementPolicy

import (
	"context"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/identityschema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
)

const (
	ResourceName  = "microsoft365_graph_beta_identity_and_access_role_management_policy"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180

	// importAllRulesKey is the private state key ImportState sets so the first Read after an import
	// populates every supported rule instead of only the rules already tracked in state.
	importAllRulesKey = "import_all_rules"
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &RoleManagementPolicyResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &RoleManagementPolicyResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &RoleManagementPolicyResource{}

	// Enables plan modification/diff suppression
	_ resource.ResourceWithModifyPlan = &RoleManagementPolicyResource{}

	// Enables common validations
	_ resource.ResourceWithConfigValidators = &RoleManagementPolicyResource{}

	// Enables rule-level validation that cannot be expressed with attribute validators
	_ resource.ResourceWithValidateConfig = &RoleManagementPolicyResource{}

	// Enables identity schema for list resource support
	_ resource.ResourceWithIdentity = &RoleManagementPolicyResource{}
)

// https://learn.microsoft.com/en-us/graph/api/unifiedrolemanagementpolicy-update?view=graph-rest-beta
func NewRoleManagementPolicyResource() resource.Resource {
	return &RoleManagementPolicyResource{
		ReadPermissions: []string{
			"RoleManagementPolicy.Read.Directory",
			"RoleManagementPolicy.Read.AzureADGroup",
		},
		WritePermissions: []string{
			"RoleManagementPolicy.ReadWrite.Directory",
			"RoleManagementPolicy.ReadWrite.AzureADGroup",
		},
	}
}

type RoleManagementPolicyResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
}

// Metadata returns the resource type name.
func (r *RoleManagementPolicyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

// Configure sets the client for the resource.
func (r *RoleManagementPolicyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

// ImportState imports the resource state by role management policy ID.
//
// The import marks private state so the following Read populates every supported rule of the policy.
func (r *RoleManagementPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	resp.Diagnostics.Append(resp.Private.SetKey(ctx, importAllRulesKey, []byte(`true`))...)
}

// IdentitySchema defines the identity schema for this resource, used by list operations to uniquely identify instances
func (r *RoleManagementPolicyResource) IdentitySchema(ctx context.Context, req resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	resp.IdentitySchema = identityschema.Schema{
		Attributes: map[string]identityschema.Attribute{
			"id": identityschema.StringAttribute{
				RequiredForImport: true,
			},
		},
	}
}

// ConfigValidators returns the resource level validators.
func (r *RoleManagementPolicyResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		resourcevalidator.ExactlyOneOf(path.MatchRoot("policy_id"), path.MatchRoot("role_definition_id")),
		resourcevalidator.Conflicting(path.MatchRoot("policy_id"), path.MatchRoot("scope_id")),
		resourcevalidator.Conflicting(path.MatchRoot("policy_id"), path.MatchRoot("scope_type")),
	}
}

// Schema defines the schema for the resource.
func (r *RoleManagementPolicyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages the rules of an existing Privileged Identity Management (PIM) role management policy using the `/policies/roleManagementPolicies` endpoint. " +
			"Microsoft Entra creates one policy for every role and scope, so this resource adopts the existing policy, identified either by `role_definition_id` " +
			"(with `scope_id` and `scope_type`) or by `policy_id`, and updates only the rules listed in the configuration. " +
			"Destroying the resource removes it from Terraform state and leaves the policy rules as they are.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The ID of the role management policy, for example `DirectoryRole_{tenantId}_{policyGuid}`.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"policy_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the role management policy to adopt. Conflicts with `role_definition_id`, `scope_id` and `scope_type`. " +
					"When the policy is identified by role definition, this is computed from the policy assignment. Changing it adopts a different policy.",
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplaceIfConfigured(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"role_definition_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the role definition whose policy is adopted. Use the role template ID for built-in Microsoft Entra roles, " +
					"or `member` / `owner` for PIM for Groups. The policy is found through `/policies/roleManagementPolicyAssignments`. " +
					"Computed when the policy is identified by `policy_id`. Changing it adopts a different policy.",
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplaceIfConfigured(),
				},
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile(roleDefinitionIdRegex), "must be a GUID, member or owner"),
				},
			},
			"scope_id": schema.StringAttribute{
				MarkdownDescription: "The scope of the policy. Use `/` for Microsoft Entra roles (the default when `role_definition_id` is set) or the group ID for PIM for Groups. " +
					"Changing it adopts a different policy.",
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplaceIfConfigured(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"scope_type": schema.StringAttribute{
				MarkdownDescription: "The type of the policy scope: `DirectoryRole` for Microsoft Entra roles (the default when `role_definition_id` is set) or `Group` for PIM for Groups. " +
					"Changing it adopts a different policy.",
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplaceIfConfigured(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf(scopeTypeDirectoryRole, scopeTypeGroup),
				},
			},
			"display_name": schema.StringAttribute{
				MarkdownDescription: "The display name of the policy, set by Microsoft Entra.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description of the policy, set by Microsoft Entra.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"is_organization_default": schema.BoolAttribute{
				MarkdownDescription: "Whether the policy is the organization default policy.",
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"last_modified_date_time": schema.StringAttribute{
				MarkdownDescription: "When the policy was last modified. Null until the policy is first changed.",
				Computed:            true,
			},
			"expiration_rules": schema.MapNestedAttribute{
				MarkdownDescription: "Expiration rules to manage, keyed by rule ID, for example `Expiration_Admin_Eligibility`, `Expiration_Admin_Assignment` or `Expiration_EndUser_Assignment`. " +
					"Rules that are not listed are left unchanged.",
				Optional: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"is_expiration_required": schema.BoolAttribute{
							MarkdownDescription: "Whether assignments or activations must have an expiration date. `false` allows permanent assignments.",
							Required:            true,
						},
						"maximum_duration": schema.StringAttribute{
							MarkdownDescription: "The maximum duration as an ISO 8601 duration made of days, hours, minutes and seconds, for example `P180D`, `PT8H` or `P1DT12H`. " +
								"Equivalent values, such as `PT24H` for `P1D`, are kept as written.",
							Required: true,
							Validators: []validator.String{
								isoDurationValidator{},
							},
						},
					},
				},
			},
			"enablement_rules": schema.MapNestedAttribute{
				MarkdownDescription: "Enablement rules to manage, keyed by rule ID, for example `Enablement_Admin_Eligibility`, `Enablement_Admin_Assignment` or `Enablement_EndUser_Assignment`. " +
					"Rules that are not listed are left unchanged.",
				Optional: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"enabled_rules": schema.SetAttribute{
							MarkdownDescription: "The requirements that are enabled: `MultiFactorAuthentication`, `Justification` and `Ticketing`. Use an empty set to require none.",
							Required:            true,
							ElementType:         types.StringType,
							Validators: []validator.Set{
								setvalidator.ValueStringsAre(stringvalidator.OneOf(enablementRuleValues...)),
							},
						},
					},
				},
			},
			"notification_rules": schema.MapNestedAttribute{
				MarkdownDescription: "Notification rules to manage, keyed by rule ID, for example `Notification_Admin_Admin_Eligibility` or `Notification_Approver_EndUser_Assignment`. " +
					"The notification and recipient types are fixed by the rule ID and are taken from the policy. Rules that are not listed are left unchanged.",
				Optional: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"notification_level": schema.StringAttribute{
							MarkdownDescription: "Which notifications are sent: `All`, `Critical` or `None`.",
							Required:            true,
							Validators: []validator.String{
								stringvalidator.OneOf(notificationLevelValues...),
							},
						},
						"is_default_recipients_enabled": schema.BoolAttribute{
							MarkdownDescription: "Whether the default recipients for this recipient type receive the notification.",
							Required:            true,
						},
						"notification_recipients": schema.SetAttribute{
							MarkdownDescription: "Additional email addresses that receive the notification. Defaults to an empty set.",
							Optional:            true,
							Computed:            true,
							ElementType:         types.StringType,
							Default:             setdefault.StaticValue(types.SetValueMust(types.StringType, []attr.Value{})),
							Validators: []validator.Set{
								setvalidator.ValueStringsAre(stringvalidator.LengthAtLeast(1)),
							},
						},
					},
				},
			},
			"authentication_context_rules": schema.MapNestedAttribute{
				MarkdownDescription: "Authentication context rules to manage, keyed by rule ID, for example `AuthenticationContext_EndUser_Assignment`. " +
					"Rules that are not listed are left unchanged.",
				Optional: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"is_enabled": schema.BoolAttribute{
							MarkdownDescription: "Whether activation requires the Conditional Access authentication context in `claim_value`.",
							Required:            true,
						},
						"claim_value": schema.StringAttribute{
							MarkdownDescription: "The authentication context class reference ID, `c1` to `c99`. Required when `is_enabled` is `true` and must be omitted when it is `false`.",
							Optional:            true,
							Validators: []validator.String{
								stringvalidator.RegexMatches(regexp.MustCompile(claimValueRegex), "must be an authentication context ID from c1 to c99"),
							},
						},
					},
				},
			},
			"timeouts": commonschema.ResourceTimeouts(ctx),
		},
	}
}
