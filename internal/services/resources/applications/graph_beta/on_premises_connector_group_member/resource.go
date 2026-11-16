package graphBetaApplicationsOnPremisesConnectorGroupMember

import (
	"context"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/identityschema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
)

const (
	ResourceName  = "microsoft365_graph_beta_applications_on_premises_connector_group_member"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &OnPremisesConnectorGroupMemberResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &OnPremisesConnectorGroupMemberResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &OnPremisesConnectorGroupMemberResource{}

	// Enables identity schema for list resource support
	_ resource.ResourceWithIdentity = &OnPremisesConnectorGroupMemberResource{}
)

// NewOnPremisesConnectorGroupMemberResource creates the connector group member resource.
//
// Microsoft Learn lists delegated Directory.ReadWrite.All as the least privileged
// permission for every API this resource calls, and marks application
// permissions as "Not supported":
//   - https://learn.microsoft.com/en-us/graph/api/connector-post-memberof?view=graph-rest-beta
//   - https://learn.microsoft.com/en-us/graph/api/connector-get?view=graph-rest-beta
//   - https://learn.microsoft.com/en-us/graph/api/connector-list-memberof?view=graph-rest-beta
//   - https://learn.microsoft.com/en-us/graph/api/connectorgroup-list?view=graph-rest-beta
func NewOnPremisesConnectorGroupMemberResource() resource.Resource {
	return &OnPremisesConnectorGroupMemberResource{
		ReadPermissions: []string{
			"Directory.ReadWrite.All",
		},
		WritePermissions: []string{
			"Directory.ReadWrite.All",
		},
		ResourcePath: "/onPremisesPublishingProfiles/applicationProxy/connectors/{connector-id}/memberOf/$ref",
	}
}

type OnPremisesConnectorGroupMemberResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

// Metadata returns the resource type name.
func (r *OnPremisesConnectorGroupMemberResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

// Configure sets the client for the resource.
func (r *OnPremisesConnectorGroupMemberResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

// ImportState imports the resource by connector ID. A connector belongs to
// exactly one connector group, so the connector ID identifies the membership and
// Read fills connector_group_id from Microsoft Graph.
func (r *OnPremisesConnectorGroupMemberResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("connector_id"), req.ID)...)
}

// IdentitySchema defines the identity schema for this resource.
func (r *OnPremisesConnectorGroupMemberResource) IdentitySchema(ctx context.Context, req resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	resp.IdentitySchema = identityschema.Schema{
		Attributes: map[string]identityschema.Attribute{
			"id": identityschema.StringAttribute{
				RequiredForImport: true,
			},
		},
	}
}

// Schema returns the schema for the resource.
func (r *OnPremisesConnectorGroupMemberResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Assigns an installed Microsoft Entra private network (Application Proxy) connector to a connector group using the Microsoft Graph beta " +
			"`/onPremisesPublishingProfiles/applicationProxy/connectors/{id}/memberOf/$ref` endpoint. " +
			"A connector always belongs to exactly one connector group, so creating or updating this resource moves the connector into `connector_group_id`. " +
			"Deleting this resource moves the connector back to the tenant's default connector group.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier for this membership. It is the connector ID, because a connector belongs to exactly one connector group.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"connector_id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier of the installed private network connector. Changing this forces a new resource.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile(constants.GuidRegex), "must be a valid GUID in the format 00000000-0000-0000-0000-000000000000"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"connector_group_id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier of the connector group the connector is assigned to. Changing this moves the connector to the new group in place.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile(constants.GuidRegex), "must be a valid GUID in the format 00000000-0000-0000-0000-000000000000"),
				},
			},
			"connector_group_name": schema.StringAttribute{
				MarkdownDescription: "The name of the connector group the connector is currently assigned to.",
				Computed:            true,
			},
			"machine_name": schema.StringAttribute{
				MarkdownDescription: "The name of the computer on which the connector is installed and runs.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"external_ip": schema.StringAttribute{
				MarkdownDescription: "The external IP address as detected by the connector server.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"timeouts": commonschema.ResourceTimeouts(ctx),
		},
	}
}
