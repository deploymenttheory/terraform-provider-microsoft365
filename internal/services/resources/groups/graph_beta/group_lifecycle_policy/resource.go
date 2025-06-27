package graphBetaGroupLifecyclePolicy

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName  = "graph_beta_groups_group_lifecycle_policy"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &GroupLifecyclePolicyResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &GroupLifecyclePolicyResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &GroupLifecyclePolicyResource{}
)

func NewGroupLifecyclePolicyResource() resource.Resource {
	return &GroupLifecyclePolicyResource{
		ReadPermissions: []string{
			"Group.Read.All",
			"Directory.Read.All",
		},
		WritePermissions: []string{
			"Group.ReadWrite.All",
			"Directory.ReadWrite.All",
		},
		ResourcePath: "/groupLifecyclePolicies",
	}
}

type GroupLifecyclePolicyResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ProviderTypeName string
	TypeName         string
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

// Metadata returns the resource type name.
func (r *GroupLifecyclePolicyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	r.ProviderTypeName = req.ProviderTypeName
	r.TypeName = ResourceName
	resp.TypeName = req.ProviderTypeName + "_" + ResourceName
}

// FullTypeName returns the full type name of the resource for logging purposes.
func (r *GroupLifecyclePolicyResource) FullTypeName() string {
	return r.ProviderTypeName + "_" + r.TypeName
}

// Configure sets the client for the resource.
func (r *GroupLifecyclePolicyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, constants.PROVIDER_NAME+"_"+ResourceName)
}

// ImportState imports the resource state.
func (r *GroupLifecyclePolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Schema defines the schema for the resource.
func (r *GroupLifecyclePolicyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages group lifecycle policies for Microsoft 365 groups using the `/groupLifecyclePolicies` endpoint. This resource enables administrators to set expiration periods for groups, requiring owners to renew them within specified time intervals. When a group reaches its expiration, it can be renewed to extend the expiration date, or if not renewed, it expires and is deleted with a 30-day restoration window.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "A unique identifier for a policy. Read-only.",
			},
			"alternate_notification_emails": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "List of email address to send notifications for groups without owners. Multiple email address can be defined by separating email address with a semicolon.",
			},
			"group_lifetime_in_days": schema.Int32Attribute{
				Required:            true,
				MarkdownDescription: "Number of days before a group expires and needs to be renewed. Once renewed, the group expiration is extended by the number of days defined.",
				Validators: []validator.Int32{
					int32validator.AtLeast(1),
					int32validator.AtMost(3650), // 10 years max
				},
			},
			"managed_group_types": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The group type for which the expiration policy applies. Possible values are **All**, **Selected** or **None**.",
				Validators: []validator.String{
					stringvalidator.OneOf("All", "Selected", "None"),
				},
			},
			"timeouts": commonschema.Timeouts(ctx),
		},
	}
}
