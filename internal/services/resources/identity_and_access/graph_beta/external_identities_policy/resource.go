package graphBetaExternalIdentitiesPolicy

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
	ResourceName  = "microsoft365_graph_beta_identity_and_access_external_identities_policy"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180

	// singletonID is the fixed identifier for the tenant-wide external identities policy.
	singletonID = "externalIdentityPolicy"
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &ExternalIdentitiesPolicyResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &ExternalIdentitiesPolicyResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &ExternalIdentitiesPolicyResource{}
)

func NewExternalIdentitiesPolicyResource() resource.Resource {
	return &ExternalIdentitiesPolicyResource{
		ReadPermissions: []string{
			"Policy.Read.All",
		},
		WritePermissions: []string{
			"Policy.ReadWrite.ExternalIdentities",
		},
	}
}

type ExternalIdentitiesPolicyResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
}

// Metadata returns the resource type name.
func (r *ExternalIdentitiesPolicyResource) Metadata(
	ctx context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = ResourceName
}

// Configure sets the client for the resource.
func (r *ExternalIdentitiesPolicyResource) Configure(
	ctx context.Context,
	req resource.ConfigureRequest,
	resp *resource.ConfigureResponse,
) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

// ImportState imports the resource state using the fixed singleton identifier.
func (r *ExternalIdentitiesPolicyResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	if req.ID != singletonID {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			"Import this singleton using the literal ID externalIdentityPolicy.",
		)
		return
	}
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Schema defines the schema for the resource.
func (r *ExternalIdentitiesPolicyResource) Schema(
	ctx context.Context,
	req resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages the tenant-wide Microsoft Entra external identities policy using the Microsoft Graph beta `/policies/externalIdentitiesPolicy` endpoint.\n\n" +
			"This is a **singleton resource** — one policy exists per tenant. Create and update use PATCH to configure the existing policy. " +
			"On destroy, Terraform removes the resource from state only, leaving the policy unchanged. Manage only one instance per tenant across all Terraform states.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The fixed singleton identifier `externalIdentityPolicy`.",
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
			"allow_external_identities_to_leave": schema.BoolAttribute{
				MarkdownDescription: "Whether external users can leave this tenant themselves. When false, an administrator must remove them. Disabling this requires the organization privacy profile to specify both a privacy contact and privacy statement URL. Must be explicitly set to `true` or `false`.",
				Required:            true,
			},
			"allow_deleted_identities_data_removal": schema.BoolAttribute{
				MarkdownDescription: "Reserved for future use by Microsoft Graph. The API accepts and returns this setting, but Microsoft does not document an operational effect. Must be explicitly set to `true` or `false`.",
				Required:            true,
			},
			"timeouts": commonschema.ResourceTimeouts(ctx),
		},
	}
}
