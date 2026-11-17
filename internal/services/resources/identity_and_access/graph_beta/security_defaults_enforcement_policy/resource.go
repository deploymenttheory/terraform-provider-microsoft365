package graphBetaIdentityAndAccessSecurityDefaultsEnforcementPolicy

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/identityschema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
)

const (
	ResourceName  = "microsoft365_graph_beta_identity_and_access_security_defaults_enforcement_policy"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180

	// singletonID is the fixed identifier for the tenant-wide security defaults enforcement policy.
	// It matches the Graph path segment so that import does not depend on the object id Graph returns.
	singletonID = "identitySecurityDefaultsEnforcementPolicy"
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &SecurityDefaultsEnforcementPolicyResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &SecurityDefaultsEnforcementPolicyResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &SecurityDefaultsEnforcementPolicyResource{}

	// Enables plan modification/diff suppression
	_ resource.ResourceWithModifyPlan = &SecurityDefaultsEnforcementPolicyResource{}

	// Enables identity schema for list resource support
	_ resource.ResourceWithIdentity = &SecurityDefaultsEnforcementPolicyResource{}
)

// https://learn.microsoft.com/en-us/graph/api/identitysecuritydefaultsenforcementpolicy-get?view=graph-rest-beta#permissions
// https://learn.microsoft.com/en-us/graph/api/identitysecuritydefaultsenforcementpolicy-update?view=graph-rest-beta#permissions
func NewSecurityDefaultsEnforcementPolicyResource() resource.Resource {
	return &SecurityDefaultsEnforcementPolicyResource{
		ReadPermissions: []string{
			"Policy.Read.All",
		},
		WritePermissions: []string{
			"Policy.Read.All",
			"Policy.ReadWrite.ConditionalAccess",
		},
	}
}

type SecurityDefaultsEnforcementPolicyResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
}

// Metadata returns the resource type name.
func (r *SecurityDefaultsEnforcementPolicyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

// Configure sets the client for the resource.
func (r *SecurityDefaultsEnforcementPolicyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

// ImportState imports the resource state using the fixed singleton identifier.
func (r *SecurityDefaultsEnforcementPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if req.ID != singletonID {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			"Import this singleton using the literal ID "+singletonID+".",
		)
		return
	}
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// IdentitySchema defines the identity schema for this resource, used by list operations to uniquely identify instances
func (r *SecurityDefaultsEnforcementPolicyResource) IdentitySchema(ctx context.Context, req resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	resp.IdentitySchema = identityschema.Schema{
		Attributes: map[string]identityschema.Attribute{
			"id": identityschema.StringAttribute{
				RequiredForImport: true,
			},
		},
	}
}

// Schema defines the schema for the resource.
func (r *SecurityDefaultsEnforcementPolicyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages the tenant-wide Microsoft Entra ID security defaults enforcement policy using the Microsoft Graph beta " +
			"`/policies/identitySecurityDefaultsEnforcementPolicy` endpoint. Security defaults turn on a baseline of identity protections " +
			"for every user in the tenant, including mandatory MFA registration.\n\n" +
			"This is a **singleton resource**: one policy exists per tenant. `create` adopts the existing policy and `create`/`update` " +
			"apply `is_enabled` with PATCH. On `destroy`, Terraform only removes the resource from state and leaves security defaults " +
			"as they are. Manage only one instance per tenant across all Terraform states.\n\n" +
			"Microsoft Entra ID rejects enabling security defaults while Conditional Access policies are enabled in the tenant.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The fixed singleton identifier `" + singletonID + "`.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"is_enabled": schema.BoolAttribute{
				MarkdownDescription: "Whether Microsoft Entra ID security defaults are enabled for the tenant. Setting this to `true` " +
					"enforces security defaults, including MFA registration, for every user in the tenant. Graph rejects enabling security " +
					"defaults while Conditional Access policies are enabled.",
				Required: true,
			},
			"display_name": schema.StringAttribute{
				MarkdownDescription: "The display name of the policy. Read-only; set by Microsoft Entra ID.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description of the policy. Read-only; set by Microsoft Entra ID.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"timeouts": commonschema.ResourceTimeouts(ctx),
		},
	}
}
