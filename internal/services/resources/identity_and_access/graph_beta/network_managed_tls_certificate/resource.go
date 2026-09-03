package graphBetaNetworkManagedTLSCertificate

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	planmodifiers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/plan_modifiers"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/identityschema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName  = "microsoft365_graph_beta_identity_and_access_network_managed_tls_certificate"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	_ resource.Resource                = &NetworkManagedTLSCertificateResource{}
	_ resource.ResourceWithConfigure   = &NetworkManagedTLSCertificateResource{}
	_ resource.ResourceWithImportState = &NetworkManagedTLSCertificateResource{}
	_ resource.ResourceWithIdentity    = &NetworkManagedTLSCertificateResource{}
)

func NewNetworkManagedTLSCertificateResource() resource.Resource {
	return &NetworkManagedTLSCertificateResource{
		ReadPermissions:  []string{"NetworkAccess.Read.All"},
		WritePermissions: []string{"NetworkAccess.ReadWrite.All"},
		ResourcePath:     "/networkaccess/tls/managedCertificateAuthorityCertificates",
	}
}

type NetworkManagedTLSCertificateResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

func (r *NetworkManagedTLSCertificateResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

func (r *NetworkManagedTLSCertificateResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

func (r *NetworkManagedTLSCertificateResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *NetworkManagedTLSCertificateResource) IdentitySchema(ctx context.Context, req resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	resp.IdentitySchema = identityschema.Schema{Attributes: map[string]identityschema.Attribute{
		"id": identityschema.StringAttribute{RequiredForImport: true},
	}}
}

func (r *NetworkManagedTLSCertificateResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a Microsoft-managed certificate authority for Microsoft Entra Global Secure Access TLS inspection using the portal-backed Microsoft Graph beta `/networkaccess/tls/managedCertificateAuthorityCertificates` endpoint.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier for the Microsoft-managed TLS certificate authority.",
				Computed:            true,
				PlanModifiers:       []planmodifier.String{planmodifiers.UseStateForUnknownString()},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The internal name of the certificate authority. When omitted, the provider generates a portal-compatible name in the form `M-TLSi-xxxxx`.",
				Optional:            true,
				Computed:            true,
				Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"common_name": schema.StringAttribute{
				MarkdownDescription: "The common name (CN) of the Microsoft-managed root certificate authority.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("Microsoft Entra TLS Inspection Root CA"),
				Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"organization_name": schema.StringAttribute{
				MarkdownDescription: "The organization name (O) of the Microsoft-managed root certificate authority.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("Microsoft"),
				Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"validity_months": schema.Int32Attribute{
				MarkdownDescription: "The root certificate validity period in months. The Entra portal uses `120` months.",
				Optional:            true,
				Computed:            true,
				Default:             int32default.StaticInt32(120),
				Validators:          []validator.Int32{int32validator.AtLeast(1)},
				PlanModifiers:       []planmodifier.Int32{int32planmodifier.RequiresReplace()},
			},
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "Whether the certificate authority should be enabled for TLS inspection. Enabling sends `status = \"enabled\"` to Microsoft Graph and waits for the observed lifecycle status to become `active`. Disabling waits for `disabled`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"status": schema.StringAttribute{
				MarkdownDescription: "The lifecycle status reported by Microsoft Graph, such as `unknownFutureValue` immediately after a disabled create, `creating`, `enrolling`, `disabled`, or `active`. This differs from the command value sent when enabling the certificate authority.",
				Computed:            true,
			},
			"created_date_time": schema.StringAttribute{
				MarkdownDescription: "The date and time when the certificate authority was created.",
				Computed:            true,
			},
			"validity_start_date_time": schema.StringAttribute{
				MarkdownDescription: "The start of the root certificate validity period.",
				Computed:            true,
			},
			"validity_end_date_time": schema.StringAttribute{
				MarkdownDescription: "The end of the root certificate validity period.",
				Computed:            true,
			},
			"certificate": schema.StringAttribute{
				MarkdownDescription: "The Microsoft-managed root CA certificate returned by Microsoft Graph.",
				Computed:            true,
			},
			"timeouts": commonschema.ResourceTimeouts(ctx),
		},
	}
}
