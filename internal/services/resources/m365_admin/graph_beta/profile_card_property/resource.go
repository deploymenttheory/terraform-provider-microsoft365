package graphBetaM365AdminProfileCardProperty

import (
	"context"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
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
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
)

const (
	ResourceName  = "microsoft365_graph_beta_m365_admin_profile_card_property"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// directoryPropertyNames are the Microsoft Entra attributes that can be surfaced on the profile card
	directoryPropertyNames = []string{
		"UserPrincipalName", "Fax", "StreetAddress", "PostalCode", "StateOrProvince", "Alias",
		"CustomAttribute1", "CustomAttribute2", "CustomAttribute3", "CustomAttribute4", "CustomAttribute5",
		"CustomAttribute6", "CustomAttribute7", "CustomAttribute8", "CustomAttribute9", "CustomAttribute10",
		"CustomAttribute11", "CustomAttribute12", "CustomAttribute13", "CustomAttribute14", "CustomAttribute15",
	}

	// languageTagRegex matches a BCP 47 language tag such as 'de' or 'nb-NO'
	languageTagRegex = regexp.MustCompile(`^[A-Za-z]{2,3}(-[A-Za-z0-9]{2,8})*$`)
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &ProfileCardPropertyResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &ProfileCardPropertyResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &ProfileCardPropertyResource{}

	// Enables plan modification/diff suppression
	_ resource.ResourceWithModifyPlan = &ProfileCardPropertyResource{}

	// Enables identity schema for list resource support
	_ resource.ResourceWithIdentity = &ProfileCardPropertyResource{}
)

// https://learn.microsoft.com/en-us/graph/api/peopleadminsettings-post-profilecardproperties?view=graph-rest-beta
func NewProfileCardPropertyResource() resource.Resource {
	return &ProfileCardPropertyResource{
		ReadPermissions: []string{
			"PeopleSettings.Read.All",
		},
		WritePermissions: []string{
			"PeopleSettings.ReadWrite.All",
		},
	}
}

type ProfileCardPropertyResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
}

// Metadata returns the resource type name.
func (r *ProfileCardPropertyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

// Configure sets the client for the resource.
func (r *ProfileCardPropertyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

// ImportState imports the resource state using the directory property name as the ID.
func (r *ProfileCardPropertyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// IdentitySchema defines the identity schema for this resource, used by list operations to uniquely identify instances
func (r *ProfileCardPropertyResource) IdentitySchema(ctx context.Context, req resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	resp.IdentitySchema = identityschema.Schema{
		Attributes: map[string]identityschema.Attribute{
			"id": identityschema.StringAttribute{
				RequiredForImport: true,
			},
		},
	}
}

// Schema defines the schema for the resource.
func (r *ProfileCardPropertyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a Microsoft 365 profile card property using the `/admin/people/profileCardProperties` endpoint. " +
			"A profile card property surfaces a Microsoft Entra ID attribute (for example `CustomAttribute1` or `Alias`) on users' profile cards, " +
			"optionally with a custom label and localized labels. Deleting the resource removes the profile card configuration only; the directory attribute and its values are not changed.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The identifier of the profile card property. This is the `directory_property_name` exactly as stored by Microsoft Graph.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"directory_property_name": schema.StringAttribute{
				MarkdownDescription: "The Microsoft Entra ID attribute to surface on the profile card. Allowed values are `UserPrincipalName`, `Fax`, `StreetAddress`, `PostalCode`, " +
					"`StateOrProvince`, `Alias` and `CustomAttribute1` to `CustomAttribute15`. Microsoft Graph stores the value exactly as sent and addresses it case-sensitively, " +
					"so use the same casing as any existing property you import. Changing it replaces the property.",
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOfCaseInsensitive(directoryPropertyNames...),
				},
			},
			"is_visible": schema.BoolAttribute{
				MarkdownDescription: "Whether the attribute is shown on users' profile cards. Defaults to `true`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"annotations": schema.SetNestedAttribute{
				MarkdownDescription: "Custom labels for the attribute on the profile card. Each annotation sets a default label and optional localized labels. " +
					"Removing all annotations from configuration clears them in Microsoft Graph.",
				Optional: true,
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"display_name": schema.StringAttribute{
							MarkdownDescription: "The default label shown on the profile card for the attribute, for example `Cost Center`.",
							Required:            true,
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(1),
							},
						},
						"localizations": schema.SetNestedAttribute{
							MarkdownDescription: "Localized labels for the attribute. Users whose client uses one of these languages see the localized label instead of `display_name`.",
							Optional:            true,
							Validators: []validator.Set{
								setvalidator.SizeAtLeast(1),
							},
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"language_tag": schema.StringAttribute{
										MarkdownDescription: "The language culture code of the label, for example `de-DE` or `nb-NO`.",
										Required:            true,
										Validators: []validator.String{
											stringvalidator.RegexMatches(languageTagRegex, "must be a language tag such as 'de' or 'nb-NO'"),
										},
									},
									"display_name": schema.StringAttribute{
										MarkdownDescription: "The localized label for the language in `language_tag`.",
										Required:            true,
										Validators: []validator.String{
											stringvalidator.LengthAtLeast(1),
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
