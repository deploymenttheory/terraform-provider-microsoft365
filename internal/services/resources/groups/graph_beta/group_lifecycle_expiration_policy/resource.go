package graphBetaGroupLifecycleExpirationPolicy

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName  = "microsoft365_graph_beta_groups_group_lifecycle_expiration_policy"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &GroupLifecycleExpirationPolicyResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &GroupLifecycleExpirationPolicyResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &GroupLifecycleExpirationPolicyResource{}
)

func NewGroupLifecycleExpirationPolicyResource() resource.Resource {
	return &GroupLifecycleExpirationPolicyResource{
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

type GroupLifecycleExpirationPolicyResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

// Metadata returns the resource type name.
func (r *GroupLifecycleExpirationPolicyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

// Configure sets the client for the resource.
func (r *GroupLifecycleExpirationPolicyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

// ImportState imports the resource state.
func (r *GroupLifecycleExpirationPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Schema defines the schema for the resource.
func (r *GroupLifecycleExpirationPolicyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages group lifecycle policies for Microsoft 365 groups using the `/groupLifecyclePolicies` endpoint. " +
			"This resource enables administrators to set expiration periods for groups, requiring owners to renew them within specified " +
			"time intervals. When a group reaches its expiration, it can be renewed to extend the expiration date, or if not renewed, it " +
			"expires and is deleted with a 30-day restoration window. Renewal notifications are emailed to group owners 30 days, 15 days, " +
			"and one day prior to group expiration. Group owners must have Exchange licenses to receive notification emails. If a group is " +
			"not renewed, it is deleted along with its associated content from sources such as Outlook, SharePoint, Teams, and Power BI.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "A unique identifier for a policy. Read-only.",
			},
			"alternate_notification_emails": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "List of email address to send notifications for groups without owners. Multiple email address can be defined by separating email address with a semicolon.",
			},
			"group_lifetime_in_days": schema.Int32Attribute{
				Required:            true,
				MarkdownDescription: "Number of days before a group expires and needs to be renewed. Once renewed, the group expiration is extended by the number of days defined. Minimum value is 30 days, maximum is 99999 days.",
				Validators: []validator.Int32{
					int32validator.AtLeast(30),
					int32validator.AtMost(99999),
				},
			},
			"managed_group_types": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The group type for which the expiration policy applies. Possible values are **All**, **Selected** or **None**.",
				Validators: []validator.String{
					stringvalidator.OneOf("All", "Selected", "None"),
				},
			},
			"overwrite_existing_policy": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(false),
				MarkdownDescription: "When set to `true`, Terraform will overwrite the existing tenant-wide group lifecycle expiration policy with the " +
					"configuration specified in this resource. This is useful for first-time adoption when a policy already exists. Since only one lifecycle " +
					"expiration policy is allowed per tenant, setting this to `true` forces a PATCH operation to replace the existing policy settings. Defaults to `false`, which attempts to create (POST) a new policy first.",
			},
			"timeouts": commonschema.ResourceTimeouts(ctx),
		},
	}
}
