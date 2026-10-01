package graphBetaAuthenticationFlowsPolicy

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
)

const (
	ResourceName  = "microsoft365_graph_beta_identity_and_access_authentication_flows_policy"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180

	// singletonID is the fixed identifier for the tenant-wide authentication flows policy.
	singletonID = "authenticationFlowsPolicy"
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &AuthenticationFlowsPolicyResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &AuthenticationFlowsPolicyResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &AuthenticationFlowsPolicyResource{}
)

func NewAuthenticationFlowsPolicyResource() resource.Resource {
	return &AuthenticationFlowsPolicyResource{
		ReadPermissions: []string{
			"Policy.Read.All",
		},
		WritePermissions: []string{
			"Policy.ReadWrite.AuthenticationFlows",
		},
	}
}

type AuthenticationFlowsPolicyResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
}

// Metadata returns the resource type name.
func (r *AuthenticationFlowsPolicyResource) Metadata(
	ctx context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = ResourceName
}

// Configure sets the client for the resource.
func (r *AuthenticationFlowsPolicyResource) Configure(
	ctx context.Context,
	req resource.ConfigureRequest,
	resp *resource.ConfigureResponse,
) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

// ImportState imports the resource state using the fixed singleton identifier.
func (r *AuthenticationFlowsPolicyResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	if req.ID != singletonID {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			"Import this singleton using the literal ID authenticationFlowsPolicy.",
		)
		return
	}
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Schema defines the schema for the resource.
func (r *AuthenticationFlowsPolicyResource) Schema(
	ctx context.Context,
	req resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages the tenant-wide Microsoft Entra authentication flows policy using the Microsoft Graph beta `/policies/authenticationFlowsPolicy` endpoint.\n\n" +
			"This is a **singleton resource** — one policy exists per tenant. Create and update use PATCH to configure the existing policy. " +
			"On destroy, Terraform removes the resource from state only, leaving the policy unchanged. Manage only one instance per tenant across all Terraform states.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The fixed singleton identifier `authenticationFlowsPolicy`.",
				Computed:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"display_name": schema.StringAttribute{
				MarkdownDescription: "The read-only display name returned by Microsoft Graph.",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The read-only description returned by Microsoft Graph.",
				Computed:            true,
			},
			"self_service_sign_up": schema.SingleNestedAttribute{
				MarkdownDescription: "Tenant-wide self-service sign-up configuration for external users.",
				Required:            true,
				Attributes: map[string]schema.Attribute{
					"is_enabled": schema.BoolAttribute{
						MarkdownDescription: "Whether external users can use self-service sign-up. Must be explicitly set to `true` or `false`.",
						Required:            true,
					},
				},
			},
			"timeouts": commonschema.ResourceTimeouts(ctx),
		},
	}
}
