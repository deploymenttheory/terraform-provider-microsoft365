package graphBetaAuthorizationPolicy

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
)

const (
	ResourceName  = "microsoft365_graph_beta_identity_and_access_authorization_policy"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180

	// singletonID is the fixed identifier for the tenant-wide authorization policy.
	singletonID = "authorizationPolicy"
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &AuthorizationPolicyResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &AuthorizationPolicyResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &AuthorizationPolicyResource{}
)

func NewAuthorizationPolicyResource() resource.Resource {
	return &AuthorizationPolicyResource{
		ReadPermissions: []string{
			"Policy.Read.All",
		},
		WritePermissions: []string{
			"Policy.ReadWrite.Authorization",
		},
	}
}

type AuthorizationPolicyResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
}

// Metadata returns the resource type name.
func (r *AuthorizationPolicyResource) Metadata(
	ctx context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = ResourceName
}

// Configure sets the client for the resource.
func (r *AuthorizationPolicyResource) Configure(
	ctx context.Context,
	req resource.ConfigureRequest,
	resp *resource.ConfigureResponse,
) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

// ImportState imports the resource state using the fixed singleton identifier.
func (r *AuthorizationPolicyResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	if req.ID != singletonID {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			"Import this singleton using the literal ID authorizationPolicy.",
		)
		return
	}
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Schema defines the schema for the resource.
func (r *AuthorizationPolicyResource) Schema(
	ctx context.Context,
	req resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages the tenant-wide Microsoft Entra authorization policy using the Microsoft Graph beta `/policies/authorizationPolicy` endpoint.\n\n" +
			"This is a **singleton resource** — one policy exists per tenant. " +
			"The `create` and `update` operations use PATCH to configure the existing policy. " +
			"On `destroy`, Terraform removes the resource from state only, leaving the policy unchanged. " +
			"Manage only one instance per tenant across all Terraform states. Omitted optional settings retain their existing service values.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The fixed singleton identifier `authorizationPolicy`.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"allowed_to_sign_up_email_based_subscriptions": schema.BoolAttribute{
				MarkdownDescription: "Whether users can sign up for email-based subscriptions.",
				Required:            true,
			},
			"allowed_to_use_sspr": schema.BoolAttribute{
				MarkdownDescription: "Whether tenant administrators can use self-service password reset.",
				Required:            true,
			},
			"allow_email_verified_users_to_join_organization": schema.BoolAttribute{
				MarkdownDescription: "Whether users can join the tenant through email verification.",
				Required:            true,
			},
			"allow_invites_from": schema.StringAttribute{
				MarkdownDescription: "Who can invite guests: `none`, `adminsAndGuestInviters`, `adminsGuestInvitersAndAllMembers`, or `everyone`.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.OneOf(
						"none",
						"adminsAndGuestInviters",
						"adminsGuestInvitersAndAllMembers",
						"everyone",
					),
				},
			},
			"allow_user_consent_for_risky_apps": schema.BoolAttribute{
				MarkdownDescription: "Whether users can consent to risky applications.",
				Required:            true,
			},
			"block_msol_powershell": schema.BoolAttribute{
				MarkdownDescription: "Whether user-based access to the legacy MSOnline PowerShell service endpoint is blocked.",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description of the authorization policy.",
				Optional:            true,
				Computed:            true,
			},
			"display_name": schema.StringAttribute{
				MarkdownDescription: "Display name of the authorization policy.",
				Optional:            true,
				Computed:            true,
			},
			"enabled_preview_features": schema.SetAttribute{
				MarkdownDescription: "Features enabled for private preview on the tenant. Use an empty set to clear the list.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
			},
			"guest_user_role_id": schema.StringAttribute{
				MarkdownDescription: "Role template ID granted to guests: User (`a0b1b346-4d3e-4e8b-98f8-753987be4970`), Guest User (`10dae51f-b6af-4016-8d66-8c2a99b929b3`), or Restricted Guest User (`2af84b1e-32c8-42b7-82bc-daa82404023b`).",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.OneOf(
						"a0b1b346-4d3e-4e8b-98f8-753987be4970",
						"10dae51f-b6af-4016-8d66-8c2a99b929b3",
						"2af84b1e-32c8-42b7-82bc-daa82404023b",
					),
				},
			},
			"permission_grant_policy_ids_assigned_to_default_user_role": schema.SetAttribute{
				MarkdownDescription: "Permission grant policies assigned to the default user role. Values use `managePermissionGrantsForSelf.{id}` or `managePermissionGrantsForOwnedResource.{id}`. An empty set disables user consent to apps.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
			},
			"default_user_role_permissions": schema.SingleNestedAttribute{
				MarkdownDescription: "Customizable permissions for the default user role. All Boolean permissions must be explicitly configured.",
				Required:            true,
				Attributes: map[string]schema.Attribute{
					"allowed_to_create_apps": schema.BoolAttribute{
						MarkdownDescription: "Whether users can register applications.",
						Required:            true,
					},
					"allowed_to_create_agent_identity_blueprints": schema.BoolAttribute{
						MarkdownDescription: "Whether users can create agent identity blueprints.",
						Required:            true,
					},
					"allowed_to_create_security_groups": schema.BoolAttribute{
						MarkdownDescription: "Whether users can create security groups.",
						Required:            true,
					},
					"allowed_to_create_tenants": schema.BoolAttribute{
						MarkdownDescription: "Whether users can create Microsoft Entra tenants.",
						Required:            true,
					},
					"allowed_to_read_bitlocker_keys_for_owned_device": schema.BoolAttribute{
						MarkdownDescription: "Whether users can read BitLocker recovery keys for their owned devices.",
						Required:            true,
					},
					"allowed_to_read_other_users": schema.BoolAttribute{
						MarkdownDescription: "Whether users can read other users. Microsoft advises keeping this permission enabled.",
						Required:            true,
					},
				},
			},
			"timeouts": commonschema.ResourceTimeouts(ctx),
		},
	}
}
