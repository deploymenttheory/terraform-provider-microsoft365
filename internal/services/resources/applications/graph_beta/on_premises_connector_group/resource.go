package graphBetaApplicationsOnPremisesConnectorGroup

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/identityschema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName  = "microsoft365_graph_beta_applications_on_premises_connector_group"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	_ resource.Resource                = &OnPremisesConnectorGroupResource{}
	_ resource.ResourceWithConfigure   = &OnPremisesConnectorGroupResource{}
	_ resource.ResourceWithImportState = &OnPremisesConnectorGroupResource{}
	_ resource.ResourceWithIdentity    = &OnPremisesConnectorGroupResource{}
)

func NewOnPremisesConnectorGroupResource() resource.Resource {
	return &OnPremisesConnectorGroupResource{
		ReadPermissions: []string{
			"Directory.ReadWrite.All",
		},
		WritePermissions: []string{
			"Directory.ReadWrite.All",
		},
		ResourcePath: "/onPremisesPublishingProfiles/applicationProxy/connectorGroups",
	}
}

type OnPremisesConnectorGroupResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

func (r *OnPremisesConnectorGroupResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

func (r *OnPremisesConnectorGroupResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

func (r *OnPremisesConnectorGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *OnPremisesConnectorGroupResource) IdentitySchema(ctx context.Context, req resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	resp.IdentitySchema = identityschema.Schema{
		Attributes: map[string]identityschema.Attribute{
			"id": identityschema.StringAttribute{
				RequiredForImport: true,
			},
		},
	}
}

func (r *OnPremisesConnectorGroupResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a Microsoft Entra Application Proxy connector group using the Microsoft Graph beta `/onPremisesPublishingProfiles/applicationProxy/connectorGroups` endpoint.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier of the connector group.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name associated with the connector group.",
				Required:            true,
			},
			"region": schema.StringAttribute{
				MarkdownDescription: "The region the connector group is assigned to and optimizes traffic for. Microsoft Graph beta metadata lists `nam`, `eur`, `aus`, `asia`, `ind`, and `unknownFutureValue`. Direct API verification on 2026-07-05 also observed `japan` on the default connector group, even though that value is absent from beta metadata and the generated SDK enum. Region can only be changed while no connectors or applications are assigned to the connector group.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("nam", "eur", "aus", "asia", "ind", "unknownFutureValue", "japan"),
				},
			},
			"connector_group_type": schema.StringAttribute{
				MarkdownDescription: "The connector group type returned by Microsoft Graph. Direct beta metadata currently lists `applicationProxy`.",
				Computed:            true,
			},
			"is_default": schema.BoolAttribute{
				MarkdownDescription: "Indicates whether this is the default connector group. Only one connector group can be the default connector group, and this value is set by Microsoft Graph.",
				Computed:            true,
			},
			"timeouts": commonschema.ResourceTimeouts(ctx),
		},
	}
}
