package graphBetaAdminConsentRequestPolicy

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
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
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
)

const (
	ResourceName  = "microsoft365_graph_beta_identity_and_access_admin_consent_request_policy"
	singletonID   = "adminConsentRequestPolicy"
	CreateTimeout = 180
	ReadTimeout   = 180
	UpdateTimeout = 180
	DeleteTimeout = 180
)

var (
	_ resource.Resource                = &AdminConsentRequestPolicyResource{}
	_ resource.ResourceWithConfigure   = &AdminConsentRequestPolicyResource{}
	_ resource.ResourceWithImportState = &AdminConsentRequestPolicyResource{}
	_ resource.ResourceWithIdentity    = &AdminConsentRequestPolicyResource{}
)

type AdminConsentRequestPolicyResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
}

func NewAdminConsentRequestPolicyResource() resource.Resource {
	return &AdminConsentRequestPolicyResource{ReadPermissions: []string{"Policy.Read.All"}, WritePermissions: []string{"Policy.ReadWrite.ConsentRequest"}}
}

func (r *AdminConsentRequestPolicyResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

func (r *AdminConsentRequestPolicyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

func (r *AdminConsentRequestPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if req.ID != singletonID {
		resp.Diagnostics.AddError("Invalid import ID", "Import this tenant singleton using the literal ID adminConsentRequestPolicy.")
		return
	}
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *AdminConsentRequestPolicyResource) IdentitySchema(_ context.Context, _ resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	resp.IdentitySchema = identityschema.Schema{Attributes: map[string]identityschema.Attribute{"id": identityschema.StringAttribute{RequiredForImport: true}}}
}

func (r *AdminConsentRequestPolicyResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages the tenant-wide admin consent request policy using Microsoft Graph beta `/policies/adminConsentRequestPolicy`. Create and update configure the existing singleton with PUT. Destroy releases Terraform management and retains the tenant policy. Manage only one instance per tenant across all Terraform states.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Fixed singleton identifier `adminConsentRequestPolicy`.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"is_enabled":        schema.BoolAttribute{Required: true, MarkdownDescription: "Whether users can request administrator consent to applications."},
			"notify_reviewers":  schema.BoolAttribute{Required: true, MarkdownDescription: "Whether reviewers receive request notifications."},
			"reminders_enabled": schema.BoolAttribute{Required: true, MarkdownDescription: "Whether reviewers receive reminder notifications."},
			"request_duration_in_days": schema.Int32Attribute{
				Required:            true,
				MarkdownDescription: "Number of days a request remains active. Zero is representable for disabled policies.",
				Validators:          []validator.Int32{int32validator.AtLeast(0)},
			},
			"reviewers": schema.SetNestedAttribute{
				Required:            true,
				MarkdownDescription: "Complete set of query-based reviewer selections. An explicit empty set is supported for a disabled policy.",
				NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
					"query": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "Reviewer selection query, for example `/users/{user-object-id}`.",
						Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
					},
					"query_type": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "Query language, for example `MicrosoftGraph`.",
						Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
					},
					"query_root": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "Root used to resolve relative queries. Omit when the query does not need a root.",
						Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
					},
				}},
			},
			"version":  schema.Int32Attribute{Computed: true, MarkdownDescription: "Read-only policy version returned by Microsoft Graph."},
			"timeouts": commonschema.ResourceTimeouts(ctx),
		},
	}
}
