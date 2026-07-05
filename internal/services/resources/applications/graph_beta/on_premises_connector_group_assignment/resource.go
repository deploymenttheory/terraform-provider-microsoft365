package graphBetaApplicationsOnPremisesConnectorGroupAssignment

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
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
	ResourceName  = "microsoft365_graph_beta_applications_on_premises_connector_group_assignment"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	_ resource.Resource                = &OnPremisesConnectorGroupAssignmentResource{}
	_ resource.ResourceWithConfigure   = &OnPremisesConnectorGroupAssignmentResource{}
	_ resource.ResourceWithImportState = &OnPremisesConnectorGroupAssignmentResource{}
	_ resource.ResourceWithIdentity    = &OnPremisesConnectorGroupAssignmentResource{}
)

func NewOnPremisesConnectorGroupAssignmentResource() resource.Resource {
	return &OnPremisesConnectorGroupAssignmentResource{
		ReadPermissions: []string{
			"Directory.ReadWrite.All",
		},
		WritePermissions: []string{
			"Directory.ReadWrite.All",
		},
		ResourcePath: "/applications/{application-id}/connectorGroup/$ref",
	}
}

type OnPremisesConnectorGroupAssignmentResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

func (r *OnPremisesConnectorGroupAssignmentResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

func (r *OnPremisesConnectorGroupAssignmentResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

func (r *OnPremisesConnectorGroupAssignmentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.Split(req.ID, "/")
	if len(parts) != 2 {
		resp.Diagnostics.AddError(
			"Invalid import ID format",
			fmt.Sprintf("Import ID must be in format: application_id/connector_group_id. Got: %s", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("application_id"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("connector_group_id"), parts[1])...)
}

func (r *OnPremisesConnectorGroupAssignmentResource) IdentitySchema(ctx context.Context, req resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	resp.IdentitySchema = identityschema.Schema{
		Attributes: map[string]identityschema.Attribute{
			"id": identityschema.StringAttribute{
				RequiredForImport: true,
			},
		},
	}
}

func (r *OnPremisesConnectorGroupAssignmentResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Assigns a Microsoft Entra Application Proxy connector group to an application using the Microsoft Graph beta `/applications/{id}/connectorGroup/$ref` endpoint.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier for this connector group assignment. This is a composite ID formed by combining the application ID and connector group ID.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"application_id": schema.StringAttribute{
				MarkdownDescription: "The unique object identifier of the application.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(constants.GuidRegex),
						"must be a valid GUID in the format 00000000-0000-0000-0000-000000000000",
					),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"connector_group_id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier of the Application Proxy connector group to assign to the application.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(constants.GuidRegex),
						"must be a valid GUID in the format 00000000-0000-0000-0000-000000000000",
					),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"connector_group_name": schema.StringAttribute{
				MarkdownDescription: "The connector group name returned by Microsoft Graph for the current assignment.",
				Computed:            true,
			},
			"timeouts": commonschema.ResourceTimeouts(ctx),
		},
	}
}
