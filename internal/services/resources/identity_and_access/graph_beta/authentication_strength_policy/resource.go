package graphBetaAuthenticationStrengthPolicy

import (
	"context"
	"regexp"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/validate/attribute"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName  = "microsoft365_graph_beta_identity_and_access_authentication_strength_policy"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &AuthenticationStrengthPolicyResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &AuthenticationStrengthPolicyResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &AuthenticationStrengthPolicyResource{}
)

func NewAuthenticationStrengthPolicyResource() resource.Resource {
	return &AuthenticationStrengthPolicyResource{
		ReadPermissions: []string{
			"Policy.Read.AuthenticationMethod",
			"Policy.Read.All",
		},
		WritePermissions: []string{
			"Policy.ReadWrite.AuthenticationMethod",
			"Policy.ReadWrite.ConditionalAccess",
		},
	}
}

type AuthenticationStrengthPolicyResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
}

// Metadata returns the resource type name.
func (r *AuthenticationStrengthPolicyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

// Configure sets the client for the resource.
func (r *AuthenticationStrengthPolicyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

// ImportState imports the resource state.
func (r *AuthenticationStrengthPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Schema defines the schema for the resource.
func (r *AuthenticationStrengthPolicyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
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
			},
			"display_name": schema.StringAttribute{
				MarkdownDescription: "The display name of the authentication strength policy. Maximum length is 30 characters.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(30),
				},
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description of the authentication strength policy. Maximum length is 100 characters.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(100),
				},
			},
			"policy_type": schema.StringAttribute{
				MarkdownDescription: "Indicates whether this is a Microsoft-managed or customer-created policy.",
				Computed:            true,
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
			},
			"modified_date_time": schema.StringAttribute{
				MarkdownDescription: "The last modified date and time of the authentication strength policy.",
				Computed:            true,
			},
			"allowed_combinations": schema.SetAttribute{
				MarkdownDescription: "The authentication method combinations allowed by this authentication strength policy. " +
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
			"combination_configurations": schema.ListNestedAttribute{
				MarkdownDescription: "Configuration settings that may be required by certain authentication methods. " +
					"For example, configuring which FID02 security keys or which X.509 certificate issuers are allowed.",
				Optional: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "The unique identifier for this configuration.",
							Computed:            true,
							Validators: []validator.String{
								stringvalidator.RegexMatches(
									regexp.MustCompile(constants.GuidRegex),
									"must be a valid GUID",
								),
							},
						},
						"odata_type": schema.StringAttribute{
							MarkdownDescription: "The OData type of the configuration. Automatically inferred from `applies_to_combinations`. Can be either `#microsoft.graph.fido2CombinationConfiguration` or `#microsoft.graph.x509CertificateCombinationConfiguration`. If not specified, it will be determined based on the authentication methods in `applies_to_combinations`.",
							Computed:            true,
							Optional:            true,
							Validators: []validator.String{
								stringvalidator.OneOf(
									"#microsoft.graph.fido2CombinationConfiguration",
									"#microsoft.graph.x509CertificateCombinationConfiguration",
								),
							},
						},
						"applies_to_combinations": schema.StringAttribute{
							MarkdownDescription: "Which authentication method combination this configuration applies to. Must be one of " +
								"`fido2`, `x509CertificateSingleFactor`, or `x509CertificateMultiFactor`.",
							Required: true,
							Validators: []validator.String{
								stringvalidator.OneOf(
									"fido2",
									"x509CertificateSingleFactor",
									"x509CertificateMultiFactor",
								),
							},
						},
						"allowed_aaguids": schema.SetAttribute{
							MarkdownDescription: "(FIDO2 only) A list of AAGUIDs (Authenticator Attestation GUIDs) allowed for FIDO2 security keys. Format: `12345678-1234-1234-1234-123456789012`. Can only be set when `applies_to_combinations` is `fido2`.",
							ElementType:         types.StringType,
							Optional:            true,
							Validators: []validator.Set{
								setvalidator.ValueStringsAre(
									stringvalidator.RegexMatches(
										regexp.MustCompile(constants.GuidRegex),
										"must be a valid GUID format",
									),
								),
								attribute.SetRequiresStringValue("applies_to_combinations", []string{"fido2"}, ""),
							},
						},
						"allowed_issuer_skis": schema.SetAttribute{
							MarkdownDescription: "(X.509 only) A list of Subject Key Identifiers (SKI) in hexadecimal format identifying allowed certificate issuers. Format: `1A2B3C4D5E6F7A8B9C0D1E2F3A4B5C6D7E8F9A0B` (40 hex characters). Maximum of 5 issuers allowed. Can only be set when `applies_to_combinations` is `x509CertificateMultiFactor` or `x509CertificateSingleFactor`.",
							ElementType:         types.StringType,
							Optional:            true,
							Validators: []validator.Set{
								setvalidator.SizeAtMost(5),
								setvalidator.ValueStringsAre(
									stringvalidator.RegexMatches(
										regexp.MustCompile(constants.SubjectKeyIdentifierRegex),
										"must be a 40-character hexadecimal string",
									),
								),
								attribute.SetRequiresStringValue("applies_to_combinations", []string{"x509CertificateMultiFactor", "x509CertificateSingleFactor"}, ""),
							},
						},
						"allowed_issuers": schema.SetAttribute{
							MarkdownDescription: "(X.509 only) A list of allowed certificate issuers. Automatically computed from `allowed_issuer_skis` by adding the `CUSTOMIDENTIFIER:` prefix. Format: `CUSTOMIDENTIFIER:{SKI}` where SKI is the Subject Key Identifier. **Note**: This field is accepted by the API but may not be returned in responses.",
							ElementType:         types.StringType,
							Computed:            true,
						},
						"allowed_policy_oids": schema.SetAttribute{
							MarkdownDescription: "(X.509 only) A list of certificate policy OIDs (Object Identifiers) that are allowed. Format: `1.2.3.4.5` (dotted decimal notation). Maximum of 5 OIDs allowed. Can only be set when `applies_to_combinations` is `x509CertificateMultiFactor` or `x509CertificateSingleFactor`.",
							ElementType:         types.StringType,
							Optional:            true,
							Validators: []validator.Set{
								setvalidator.SizeAtMost(5),
								setvalidator.ValueStringsAre(
									stringvalidator.RegexMatches(
										regexp.MustCompile(constants.OIDRegex),
										"must be a valid OID in dotted decimal notation (e.g., 1.3.6.1.4.1.311.21.8.1.1)",
									),
								),
								attribute.SetRequiresStringValue("applies_to_combinations", []string{"x509CertificateMultiFactor", "x509CertificateSingleFactor"}, ""),
							},
						},
					},
				},
			},
			"timeouts": commonschema.Timeouts(ctx),
		},
	}
}
