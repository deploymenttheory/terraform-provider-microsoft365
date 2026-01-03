package graphBetaGroupLifecycleExpirationPolicyAssignment

import (
	"context"
	"regexp"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName  = "microsoft365_graph_beta_groups_group_lifecycle_expiration_policy_assignment"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &GroupLifecycleExpirationPolicyAssignmentResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &GroupLifecycleExpirationPolicyAssignmentResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &GroupLifecycleExpirationPolicyAssignmentResource{}
)

func NewGroupLifecycleExpirationPolicyAssignmentResource() resource.Resource {
	return &GroupLifecycleExpirationPolicyAssignmentResource{
		ReadPermissions: []string{
			"Group.Read.All",
			"Directory.Read.All",
		},
		WritePermissions: []string{
			"Group.ReadWrite.All",
			"Directory.ReadWrite.All",
		},
		ResourcePath: "/groupLifecyclePolicies/{id}/addGroup",
	}
}

type GroupLifecycleExpirationPolicyAssignmentResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

// Metadata returns the resource type name.
func (r *GroupLifecycleExpirationPolicyAssignmentResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

// Configure sets the client for the resource.
func (r *GroupLifecycleExpirationPolicyAssignmentResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

// ImportState imports the resource state.
func (r *GroupLifecycleExpirationPolicyAssignmentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// The import ID is the group ID
	groupID := req.ID

	// Set both group_id and id to the imported ID
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("group_id"), groupID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), groupID)...)
}

// Schema defines the schema for the resource.
func (r *GroupLifecycleExpirationPolicyAssignmentResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages the assignment of a Microsoft 365 group to the tenant's group lifecycle expiration policy. " +
			"This resource adds a group to the lifecycle policy using the `/groupLifecyclePolicies/{id}/addGroup` endpoint. " +
			"This resource is only applicable when the policy's `managed_group_types` is set to **Selected**. " +
			"When `managed_group_types` is set to **Selected**, you can add up to 500 groups to the policy. " +
			"If you need to manage more than 500 groups, set `managed_group_types` to **All** instead. " +
			"Only one lifecycle policy exists per tenant, so this resource automatically finds and uses the tenant's policy.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The unique identifier for this assignment. This is the same as the group_id.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"group_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The unique identifier of the Microsoft 365 group to add to the lifecycle policy. This group must be a Microsoft 365 group (Unified group).",
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile(constants.GuidRegex), "group_id must be a valid UUID"),
				},
			},
			"timeouts": commonschema.ResourceTimeouts(ctx),
		},
	}
}
