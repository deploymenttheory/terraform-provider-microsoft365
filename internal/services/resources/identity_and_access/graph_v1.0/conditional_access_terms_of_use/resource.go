package graphConditionalAccessTermsOfUse

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
)

const (
	ResourceName  = "graph_identity_and_access_conditional_access_terms_of_use"
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
	client           *msgraphsdk.GraphServiceClient
	ProviderTypeName string
	TypeName         string
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

func (r *ConditionalAccessTermsOfUseResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	r.ProviderTypeName = req.ProviderTypeName
	r.TypeName = ResourceName
	resp.TypeName = r.FullTypeName()
}

// FullTypeName returns the full resource type name in the format "providername_resourcename".
func (r *ConditionalAccessTermsOfUseResource) FullTypeName() string {
	return r.ProviderTypeName + "_" + r.TypeName
}

func (r *ConditionalAccessTermsOfUseResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphStableClientForResource(ctx, req, resp, constants.PROVIDER_NAME+"_"+ResourceName)
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
			},
			"user_reaccept_required_frequency": schema.StringAttribute{
				MarkdownDescription: "The duration after which the user must reaccept the terms of use. Accepted values: `P365D` for annually, `P180D` for bi-annually, `P90D` for quarterly, `P30D` for monthly.",
				Optional:            true,
			},
			"terms_expiration": schema.SingleNestedAttribute{
				MarkdownDescription: "Expiration schedule and frequency of agreement for all users.",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"start_date_time": schema.StringAttribute{
						MarkdownDescription: "The DateTime when the agreement is set to expire for all users. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time.",
						Optional:            true,
					},
					"frequency": schema.StringAttribute{
						MarkdownDescription: "Represents the frequency at which the terms will expire, after its first expiration as set in startDateTime. Accepted values: `P365D` for annually, `P180D` for bi-annually, `P90D` for quarterly, `P30D` for monthly.",
						Optional:            true,
					},
				},
			},
			"file": schema.SingleNestedAttribute{
				MarkdownDescription: "Default PDF linked to this agreement. This is required when creating a new agreement.",
				Required:            true,
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
									MarkdownDescription: "The language of the agreement file in the format 'languagecode2-country/regioncode2'. 'languagecode2' is a lowercase two-letter code derived from ISO 639-1, while 'country/regioncode2' is derived from ISO 3166 and usually consists of two uppercase letters, or a BCP-47 language tag. For example, U.S. English is `en-US`.",
									Required:            true,
								},
								"is_default": schema.BoolAttribute{
									MarkdownDescription: "If none of the languages matches the client preference, indicates whether this is the default agreement file. If none of the files are marked as default, the first one is treated as the default. Read-only.",
									Optional:            true,
									Computed:            true,
									Default:             booldefault.StaticBool(false),
								},
								"is_major_version": schema.BoolAttribute{
									MarkdownDescription: "Indicates whether the agreement file is a major version update. Major version updates invalidate the agreement's acceptances on the corresponding language.",
									Optional:            true,
									Computed:            true,
									Default:             booldefault.StaticBool(false),
								},
								"file_data": schema.SingleNestedAttribute{
									MarkdownDescription: "Data that represents the terms of use PDF document.",
									Required:            true,
									Attributes: map[string]schema.Attribute{
										"data": schema.StringAttribute{
											MarkdownDescription: "Data that represents the terms of use PDF document as raw bytes.",
											Required:            true,
										},
									},
								},
							},
						},
					},
				},
			},
			"timeouts": commonschema.Timeouts(ctx),
		},
	}
}
