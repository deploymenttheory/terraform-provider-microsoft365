package graphBetaIdentityAndAccessRoleDefinition

import (
	"context"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/identityschema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
)

const (
	ResourceName  = "microsoft365_graph_beta_identity_and_access_role_definition"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180

	// builtInRoleDetail explains why built-in roles cannot be managed (role ID, resource name)
	builtInRoleDetail = "Role definition %s is a built-in Microsoft Entra role. %s only manages custom roles; use the role_definitions data source to read built-in roles."
)

var (
	// resourceActionRegex matches a resource action in the microsoft.directory namespace
	resourceActionRegex   = regexp.MustCompile(`^microsoft\.directory/[A-Za-z0-9.]+(/[A-Za-z0-9.]+)+$`)
	resourceActionMessage = "must be a microsoft.directory resource action, for example 'microsoft.directory/applications/standard/read'"
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &RoleDefinitionResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &RoleDefinitionResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &RoleDefinitionResource{}

	// Enables plan modification/diff suppression
	_ resource.ResourceWithModifyPlan = &RoleDefinitionResource{}

	// Enables identity schema for list resource support
	_ resource.ResourceWithIdentity = &RoleDefinitionResource{}
)

// https://learn.microsoft.com/en-us/graph/api/rbacapplication-post-roledefinitions?view=graph-rest-beta
func NewRoleDefinitionResource() resource.Resource {
	return &RoleDefinitionResource{
		ReadPermissions: []string{
			"RoleManagement.Read.Directory",
		},
		WritePermissions: []string{
			"RoleManagement.ReadWrite.Directory",
		},
	}
}

type RoleDefinitionResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
}

// Metadata returns the resource type name.
func (r *RoleDefinitionResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

// Configure sets the client for the resource.
func (r *RoleDefinitionResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

// ImportState imports the resource state.
func (r *RoleDefinitionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// IdentitySchema defines the identity schema for this resource, used by list operations to uniquely identify instances
func (r *RoleDefinitionResource) IdentitySchema(ctx context.Context, req resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	resp.IdentitySchema = identityschema.Schema{
		Attributes: map[string]identityschema.Attribute{
			"id": identityschema.StringAttribute{
				RequiredForImport: true,
			},
		},
	}
}

// Schema defines the schema for the resource.
func (r *RoleDefinitionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages custom Microsoft Entra directory role definitions using the `/roleManagement/directory/roleDefinitions` endpoint. " +
			"Custom roles grant a chosen set of `microsoft.directory/*` resource actions and require Microsoft Entra ID P1 or P2. Built-in roles cannot be managed or imported by this resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier of the role definition. When `template_id` is set, Microsoft Entra uses it as the role definition ID.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"display_name": schema.StringAttribute{
				MarkdownDescription: "The display name of the custom role. Must be unique among the tenant's role definitions.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 256),
				},
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description of the custom role.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(1024),
				},
			},
			"is_enabled": schema.BoolAttribute{
				MarkdownDescription: "Whether the role can be assigned. A disabled role keeps its definition but cannot be used in new assignments.",
				Required:            true,
			},
			"template_id": schema.StringAttribute{
				MarkdownDescription: "A custom template identifier for the role. Microsoft Entra uses it as the role definition ID, so it must be a GUID that is unique in the tenant. " +
					"When omitted, Microsoft Entra sets it to the generated role ID. Changing it replaces the role.",
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplaceIfConfigured(),
				},
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile(constants.GuidRegex), "must be a valid GUID in the format '00000000-0000-0000-0000-000000000000'"),
				},
			},
			"version": schema.StringAttribute{
				MarkdownDescription: "An optional version label for the role definition, for example `1`.",
				Optional:            true,
			},
			"role_permissions": schema.SetNestedAttribute{
				MarkdownDescription: "The permissions granted by the role. Each entry lists `microsoft.directory/*` resource actions that the role allows.",
				Required:            true,
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"allowed_resource_actions": schema.SetAttribute{
							MarkdownDescription: "The resource actions the role allows, for example `microsoft.directory/applications/standard/read`. " +
								"Each action must exist in the `microsoft.directory` resource namespace; the provider checks this before writing because Microsoft Entra silently rewrites unknown actions.",
							ElementType: types.StringType,
							Required:    true,
							Validators: []validator.Set{
								setvalidator.SizeAtLeast(1),
								setvalidator.ValueStringsAre(
									stringvalidator.RegexMatches(resourceActionRegex, resourceActionMessage),
								),
							},
						},
					},
				},
			},
			"resource_scopes": schema.SetAttribute{
				MarkdownDescription: "The scopes the role definition applies to. Microsoft Entra sets this to `/` (tenant-wide) for custom directory roles.",
				ElementType:         types.StringType,
				Computed:            true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"is_built_in": schema.BoolAttribute{
				MarkdownDescription: "Whether the role definition is built in. Always `false` for roles managed by this resource.",
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"is_privileged": schema.BoolAttribute{
				MarkdownDescription: "Whether Microsoft Entra classifies the role as privileged, based on the resource actions it grants.",
				Computed:            true,
			},
			"timeouts": commonschema.ResourceTimeouts(ctx),
		},
	}
}
