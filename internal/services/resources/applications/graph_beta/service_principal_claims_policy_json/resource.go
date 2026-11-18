package graphBetaServicePrincipalClaimsPolicyJSON

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
	ResourceName  = "microsoft365_graph_beta_applications_service_principal_claims_policy_json"
	CreateTimeout = 180
	ReadTimeout   = 180
	UpdateTimeout = 180
	DeleteTimeout = 180
)

var (
	_ resource.Resource                = &ServicePrincipalClaimsPolicyJSONResource{}
	_ resource.ResourceWithConfigure   = &ServicePrincipalClaimsPolicyJSONResource{}
	_ resource.ResourceWithImportState = &ServicePrincipalClaimsPolicyJSONResource{}
	_ resource.ResourceWithIdentity    = &ServicePrincipalClaimsPolicyJSONResource{}
)

type ServicePrincipalClaimsPolicyJSONResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
}

func NewServicePrincipalClaimsPolicyJSONResource() resource.Resource {
	return &ServicePrincipalClaimsPolicyJSONResource{
		ReadPermissions:  []string{"Policy.ReadWrite.ApplicationConfiguration"},
		WritePermissions: []string{"Policy.ReadWrite.ApplicationConfiguration"},
	}
}

func (r *ServicePrincipalClaimsPolicyJSONResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

func (r *ServicePrincipalClaimsPolicyJSONResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

func (r *ServicePrincipalClaimsPolicyJSONResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if !regexp.MustCompile(constants.GuidRegex).MatchString(req.ID) {
		resp.Diagnostics.AddError("Invalid import ID", "Import using the service principal object UUID, not the custom policy UUID.")
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("service_principal_id"), req.ID)...)
}

func (r *ServicePrincipalClaimsPolicyJSONResource) IdentitySchema(_ context.Context, _ resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	resp.IdentitySchema = identityschema.Schema{Attributes: map[string]identityschema.Attribute{"id": identityschema.StringAttribute{RequiredForImport: true}}}
}

func (r *ServicePrincipalClaimsPolicyJSONResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages the complete custom claims policy of a SAML service principal through Microsoft Graph beta `/servicePrincipals/{id}/claimsPolicy`. JSON preserves polymorphic claims, attributes, conditions, ordered transformations, and fields not represented by the generated SDK. Create and update use full PUT. Destroy releases Terraform state and retains the policy because the endpoint does not support DELETE. Manage one policy resource per service principal across all states.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Service principal object ID identifying the scoped policy endpoint.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"service_principal_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Object UUID of the SAML service principal. Changing it replaces Terraform management.",
				Validators:          []validator.String{stringvalidator.RegexMatches(regexp.MustCompile(constants.GuidRegex), "must be a valid UUID")},
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"policy_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Server-generated custom claims policy UUID returned by Graph. Import uses service_principal_id instead.",
			},
			"policy_json": schema.StringAttribute{
				Required:            true,
				CustomType:          PolicyJSONType{},
				MarkdownDescription: "Complete policy object encoded with jsonencode(). Explicit Boolean includeBasicClaimSet/includeApplicationIdInIssuer and a claims array are required to make full-replacement behavior intentional. Other beta fields are preserved. Formatting/object-key order and known service-added nullable fields are ignored; arrays and nested @odata.type discriminators remain significant.",
				Validators:          []validator.String{policyJSONValidator{}},
			},
			"timeouts": commonschema.ResourceTimeouts(ctx),
		},
	}
}
