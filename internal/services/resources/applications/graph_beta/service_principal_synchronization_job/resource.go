package graphBetaServicePrincipalSynchronizationJob

import (
	"context"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/identityschema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
)

const (
	ResourceName  = "microsoft365_graph_beta_applications_service_principal_synchronization_job"
	CreateTimeout = 300
	ReadTimeout   = 180
	UpdateTimeout = 300
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &ServicePrincipalSynchronizationJobResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &ServicePrincipalSynchronizationJobResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &ServicePrincipalSynchronizationJobResource{}

	// Enables identity schema for list resource support
	_ resource.ResourceWithIdentity = &ServicePrincipalSynchronizationJobResource{}
)

// ServicePrincipalSynchronizationJobResource defines the Microsoft Graph client and permissions for the service principal synchronization job resource.
type ServicePrincipalSynchronizationJobResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
}

// NewServicePrincipalSynchronizationJobResource creates the service principal synchronization job resource.
func NewServicePrincipalSynchronizationJobResource() resource.Resource {
	return &ServicePrincipalSynchronizationJobResource{
		ReadPermissions:  []string{"Synchronization.Read.All"},
		WritePermissions: []string{"Synchronization.ReadWrite.All", "Application.ReadWrite.OwnedBy"},
	}
}

// Metadata returns the resource type name.
func (r *ServicePrincipalSynchronizationJobResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

// Configure sets the client for the resource.
func (r *ServicePrincipalSynchronizationJobResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

// ImportState imports the resource state.
func (r *ServicePrincipalSynchronizationJobResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parent, job, err := parseResourceID(req.ID)
	if err != nil {
		resp.Diagnostics.AddError("Invalid synchronization job import ID", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("service_principal_id"), parent)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("job_id"), job)...)
}

// IdentitySchema defines the identity schema used by list operations to uniquely identify instances.
func (r *ServicePrincipalSynchronizationJobResource) IdentitySchema(_ context.Context, _ resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	resp.IdentitySchema = identityschema.Schema{Attributes: map[string]identityschema.Attribute{"id": identityschema.StringAttribute{RequiredForImport: true}}}
}

// Schema defines the schema for the resource.
func (r *ServicePrincipalSynchronizationJobResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a service principal provisioning job through Microsoft Graph beta `/servicePrincipals/{id}/synchronization/jobs`. Jobs are created disabled by default. Set `enabled = true` only after configuring the target credentials and schema. Changing the parent or synchronization template replaces the job. Destroy deletes the job and its synchronization state, while retaining the parent application and accounts already provisioned to the target.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Composite ID `{service_principal_id}/{job_id}`.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"service_principal_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Object ID of the service principal supporting the selected provisioning connector.",
				Validators:          []validator.String{stringvalidator.RegexMatches(regexp.MustCompile(constants.GuidRegex), "must be a valid UUID")},
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"template_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Synchronization template ID exposed by the service principal, for example `scim`. This is not an application-gallery template GUID.",
				Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"job_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Opaque Graph job ID. Pass this output to `service_principal_synchronization_schema_json.job_id`.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"enabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Whether the schedule is active. Defaults to false. Updates invoke start/pause only on an intentional transition; read never starts or restarts provisioning.",
			},
			"schedule": schema.SingleNestedAttribute{
				Computed:            true,
				MarkdownDescription: "Read-only schedule returned by Graph; this resource does not PATCH schedules.",
				Attributes: map[string]schema.Attribute{
					"expiration": schema.StringAttribute{Computed: true, MarkdownDescription: "Schedule expiration timestamp."},
					"interval":   schema.StringAttribute{Computed: true, MarkdownDescription: "Schedule interval in ISO 8601 duration format."},
					"state":      schema.StringAttribute{Computed: true, MarkdownDescription: "Remote schedule state: Active, Disabled, or Paused."},
				},
			},
			"status_code": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Provisioning status reported by Graph, independently of the schedule state. NotConfigured or Paused is valid for a disabled job.",
			},
			"timeouts": commonschema.ResourceTimeouts(ctx),
		},
	}
}
