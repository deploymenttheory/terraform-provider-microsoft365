package graphBetaAuthenticationStrength

import (
	"context"
	"regexp"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	ResourceName  = "graph_beta_identity_and_access_authentication_strength"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &AuthenticationStrengthResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &AuthenticationStrengthResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &AuthenticationStrengthResource{}
)

func NewAuthenticationStrengthResource() resource.Resource {
	return &AuthenticationStrengthResource{
		ReadPermissions: []string{
			"Policy.Read.AuthenticationMethod",
			"Policy.Read.All",
		},
		WritePermissions: []string{
			"Policy.ReadWrite.AuthenticationMethod",
			"Policy.ReadWrite.ConditionalAccess",
		},
		ResourcePath: "/policies/authenticationStrengthPolicies",
	}
}

type AuthenticationStrengthResource struct {
	httpClient       *client.AuthenticatedHTTPClient
	ProviderTypeName string
	TypeName         string
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

// Metadata returns the resource type name.
func (r *AuthenticationStrengthResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	r.ProviderTypeName = req.ProviderTypeName
	r.TypeName = ResourceName
	resp.TypeName = r.FullTypeName()
}

// FullTypeName returns the full resource type name in the format "providername_resourcename".
func (r *AuthenticationStrengthResource) FullTypeName() string {
	return r.ProviderTypeName + "_" + r.TypeName
}

// Configure sets the client for the resource.
func (r *AuthenticationStrengthResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.httpClient = client.SetGraphV1HTTPClientForResource(ctx, req, resp, constants.PROVIDER_NAME+"_"+ResourceName)
}

// ImportState imports the resource state.
func (r *AuthenticationStrengthResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Schema defines the schema for the resource.
func (r *AuthenticationStrengthResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages Microsoft 365 Authentication Strength Policies using the `/identity/conditionalAccess/authenticationStrength/policies` " +
			"endpoint. Authentication Strength Policies define authentication method combinations that can be used in Conditional Access policies. Learn more here: " +
			"https://learn.microsoft.com/en-us/entra/identity/authentication/concept-authentication-strength-advanced-options",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "String (identifier)",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(constants.GuidRegex),
						"must be a valid GUID",
					),
				},
			},
			"display_name": schema.StringAttribute{
				MarkdownDescription: "The display name of the authentication strength policy.",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description of the authentication strength policy.",
				Optional:            true,
			},
			"policy_type": schema.StringAttribute{
				MarkdownDescription: "Indicates whether this is a Microsoft-managed or customer-created policy.",
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("builtIn", "custom"),
				},
			},
			"requirements_satisfied": schema.StringAttribute{
				MarkdownDescription: "Describes the type of authentication method target that this authentication strength satisfies.",
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("mfa", "singleFactor"),
				},
			},
			"created_date_time": schema.StringAttribute{
				MarkdownDescription: "The creation date and time of the authentication strength policy.",
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(constants.TimeFormatRFC3339Regex),
						"must be a valid RFC3339 date-time string",
					),
				},
			},
			"modified_date_time": schema.StringAttribute{
				MarkdownDescription: "The last modified date and time of the authentication strength policy.",
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(constants.TimeFormatRFC3339Regex),
						"must be a valid RFC3339 date-time string",
					),
				},
			},
			"allowed_combinations": schema.SetAttribute{
				MarkdownDescription: "he authentication method combinations allowed by this authentication strength policy. " +
					"The possible values of this are: password, voice, hardwareOath, softwareOath, sms, fido2, windowsHelloForBusiness, " +
					"microsoftAuthenticatorPush, deviceBasedPush, temporaryAccessPassOneTime, temporaryAccessPassMultiUse, email, " +
					"x509CertificateSingleFactor, x509CertificateMultiFactor, federatedSingleFactor, federatedMultiFactor, unknownFutureValue, " +
					"qrCodePin.",
				ElementType: types.StringType,
				Required:    true,
				Validators: []validator.Set{
					setvalidator.ValueStringsAre(
						stringvalidator.OneOf(
							"deviceBasedPush",
							"federatedMultiFactor",
							"federatedSingleFactor",
							"fido2",
							"hardwareOath,federatedSingleFactor",
							"microsoftAuthenticatorPush,federatedSingleFactor",
							"password",
							"password,hardwareOath",
							"password,microsoftAuthenticatorPush",
							"password,sms",
							"password,softwareOath",
							"password,voice",
							"qrCodePin",
							"sms",
							"sms,federatedSingleFactor",
							"softwareOath,federatedSingleFactor",
							"temporaryAccessPassMultiUse",
							"temporaryAccessPassOneTime",
							"voice,federatedSingleFactor",
							"windowsHelloForBusiness",
							"x509CertificateMultiFactor",
							"x509CertificateSingleFactor",
						),
					),
				},
			},

			"timeouts": commonschema.Timeouts(ctx),
		},
	}
}
