package graphConditionalAccessTermsOfUse

import (
	"context"
	"regexp"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

const (
	ResourceName  = "microsoft365_graph_identity_and_access_conditional_access_terms_of_use"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &ConditionalAccessTermsOfUseResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &ConditionalAccessTermsOfUseResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &ConditionalAccessTermsOfUseResource{}
)

func NewConditionalAccessTermsOfUseResource() resource.Resource {
	return &ConditionalAccessTermsOfUseResource{
		ReadPermissions: []string{
			"Agreement.Read.All",
			"Agreement.ReadWrite.All",
		},
		WritePermissions: []string{
			"Agreement.ReadWrite.All",
		},
		ResourcePath: "/agreements",
	}
}

type ConditionalAccessTermsOfUseResource struct {
	httpClient       *client.AuthenticatedHTTPClient
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

func (r *ConditionalAccessTermsOfUseResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

func (r *ConditionalAccessTermsOfUseResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.httpClient = client.SetGraphBetaHTTPClientForResource(ctx, req, resp, ResourceName)
}

func (r *ConditionalAccessTermsOfUseResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Schema defines the schema for the resource.
func (r *ConditionalAccessTermsOfUseResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages Microsoft 365 Terms of Use Agreements using the `/agreements` endpoint. " +
			"Terms of use agreements allow organizations to present information that users must accept before accessing data or applications. " +
			"These agreements can be used to ensure compliance with legal or regulatory requirements.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier for the agreement. This is automatically generated when the agreement is created.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"display_name": schema.StringAttribute{
				MarkdownDescription: "Display name of the agreement. The display name is used for internal tracking of the agreement but isn't shown to end users who view the agreement.",
				Required:            true,
			},
			"is_viewing_before_acceptance_required": schema.BoolAttribute{
				MarkdownDescription: "Indicates whether the user has to expand the agreement before accepting.",
				Required:            true,
			},
			"is_per_device_acceptance_required": schema.BoolAttribute{
				MarkdownDescription: "This setting enables you to require end users to accept this agreement on every device that they're accessing it from. The end user is required to register their device in Microsoft Entra ID, if they haven't already done so.",
				Required:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
			"user_reaccept_required_frequency": schema.StringAttribute{
				MarkdownDescription: "The duration after which the user must reaccept the terms of use. " +
					"Must be in ISO 8601 duration format (e.g., `P90D`).\n\n" +
					"**Common values:**\n" +
					"- `P365D` - Annually (365 days)\n" +
					"- `P180D` - Bi-annually (180 days)\n" +
					"- `P90D` - Quarterly (90 days)\n" +
					"- `P30D` - Monthly (30 days)\n" +
					"- `P270D` - 270 days",
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^P\d+D$`),
						"must be in ISO 8601 duration format starting with P and ending with D (e.g., P90D, P365D)",
					),
				},
			},
			"terms_expiration": schema.SingleNestedAttribute{
				MarkdownDescription: "Expiration schedule and frequency of agreement for all users.",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"start_date_time": schema.StringAttribute{
						MarkdownDescription: "The date when the agreement is set to expire for all users. " +
							"Must be in YYYY-MM-DD format (e.g., `2025-12-31`) and is always in UTC time. " +
							"The time portion (T00:00:00Z) will be automatically appended.",
						Optional: true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
						Validators: []validator.String{
							stringvalidator.RegexMatches(
								regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`),
								"must be a valid date in YYYY-MM-DD format (e.g., 2025-12-31)",
							),
						},
					},
					"frequency": schema.StringAttribute{
						MarkdownDescription: "Represents the frequency at which the terms will expire, after its first expiration as set in startDateTime. " +
							"Must be in ISO 8601 duration format.\n\n" +
							"**Accepted values:**\n" +
							"- `P365D` - Annually (365 days)\n" +
							"- `P180D` - Bi-annually (180 days)\n" +
							"- `P90D` - Quarterly (90 days)\n" +
							"- `P30D` - Monthly (30 days)",
						Optional: true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
						Validators: []validator.String{
							stringvalidator.OneOf("P365D", "P180D", "P90D", "P30D"),
						},
					},
				},
			},
			"file": schema.SingleNestedAttribute{
				MarkdownDescription: "Default PDF linked to this agreement. This is required when creating a new agreement.",
				Required:            true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.RequiresReplace(),
				},
				Attributes: map[string]schema.Attribute{
					"localizations": schema.SetNestedAttribute{
						MarkdownDescription: "The localized version of the terms of use agreement files attached to the agreement.",
						Required:            true,
						Validators: []validator.Set{
							setvalidator.SizeAtLeast(1),
						},
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"file_name": schema.StringAttribute{
									MarkdownDescription: "Name of the agreement file (for example, TOU.pdf).",
									Required:            true,
								},
								"display_name": schema.StringAttribute{
									MarkdownDescription: "Localized display name of the policy file of an agreement. The localized display name is shown to end users who view the agreement.",
									Required:            true,
								},
								"language": schema.StringAttribute{
									MarkdownDescription: "The language of the agreement file. When `is_default` is `true`, must be `en-US` (full format with country code). " +
										"When `is_default` is `false`, must use only the two-letter language code (e.g., `en`, `fr`, `de`). " +
										"The language code is a lowercase two-letter code derived from ISO 639-1.",
									Required: true,
									Validators: []validator.String{
										stringvalidator.OneOf("en-US", // Default
											"en",             // English
											"en-GB",          // English (United Kingdom)
											"af",             // Afrikaans
											"am",             // Amharic
											"ar-SA",          // Arabic (Saudi Arabia)
											"hy",             // Armenian
											"as",             // Assamese
											"az",             // Azerbaijani
											"be",             // Belarusian
											"bn",             // Bangla
											"bn-IN",          // Bangla (India)
											"eu",             // Basque
											"ku-Arab",        // Central Kurdish (Arabic)
											"ja",             // Japanese
											"zu",             // Zulu
											"te",             // Telugu
											"th",             // Thai
											"ti",             // Tigrinya
											"tr",             // Turkish
											"tk",             // Turkmen
											"uk",             // Ukrainian
											"ur",             // Urdu
											"ug",             // Uyghur
											"uz",             // Uzbek
											"ca-ES-valencia", // Valencian (Spain)
											"vi",             // Vietnamese
											"cy",             // Welsh
											"wo",             // Wolof
											"yo",             // Yoruba
										),
									},
								},
								"is_default": schema.BoolAttribute{
									MarkdownDescription: "If none of the languages matches the client preference, indicates whether this is the default agreement file. If none of the files are marked as default, the first one is treated as the default. Must be true if the language is 'en-US'.",
									Required:            true,
								},
								"is_major_version": schema.BoolAttribute{
									MarkdownDescription: "Indicates whether the agreement file is a major version update. Major version updates invalidate the agreement's acceptances on the corresponding language.",
									Optional:            true,
									Computed:            true,
									Default:             booldefault.StaticBool(false),
								},
								"file_data": schema.SingleNestedAttribute{
									MarkdownDescription: "Data that represents the terms of use PDF document. Must be provided during creation but is not returned by the API and will always be null in state. This field is intentionally not persisted for security reasons.",
									Required:            true,
									Attributes: map[string]schema.Attribute{
										"data": schema.StringAttribute{
											MarkdownDescription: "Data that represents the terms of use PDF document as raw bytes (base64 encoded).",
											Required:            true,
										},
									},
								},
							},
						},
					},
				},
			},
			"timeouts": commonschema.ResourceTimeouts(ctx),
		},
	}
}
