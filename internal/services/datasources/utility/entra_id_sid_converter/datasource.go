package entra_id_sid_converter

import (
	"context"
	"regexp"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var (
	_ datasource.DataSource              = &entraIdSidConverterDataSource{}
	_ datasource.DataSourceWithConfigure = &entraIdSidConverterDataSource{}
)

func NewEntraIdSidConverterDataSource() datasource.DataSource {
	return &entraIdSidConverterDataSource{}
}

type entraIdSidConverterDataSource struct{}

func (d *entraIdSidConverterDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_utility_entra_id_sid_converter"
}

func (d *entraIdSidConverterDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
}

func (d *entraIdSidConverterDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Converts between Microsoft Entra ID (formerly Azure AD) Security Identifiers (SIDs) and Object IDs. " +
			"This utility performs bidirectional conversion - provide either a SID to get an Object ID, or an Object ID to get a SID. " +
			"Useful for hybrid environments where on-premises AD identities are synced to Entra ID.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The ID of this resource.",
			},
			"sid": schema.StringAttribute{
				Optional: true,
				Computed: true,
				MarkdownDescription: "The Security Identifier (SID) in the format `S-1-12-1-<rid1>-<rid2>-<rid3>-<rid4>`. " +
					"Provide this to convert to an Object ID, or leave empty to convert from an Object ID. " +
					"Each RID component must be a valid 32-bit unsigned integer (0 to 4,294,967,295).",
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.MatchRoot("object_id")),
					stringvalidator.RegexMatches(
						regexp.MustCompile(constants.EntraIdSidRegex),
						"SID must be in the format S-1-12-1-<rid1>-<rid2>-<rid3>-<rid4> where each RID is a numeric value",
					),
					ValidateSidRidRange(),
				},
			},
			"object_id": schema.StringAttribute{
				Optional: true,
				Computed: true,
				MarkdownDescription: "The Object ID (GUID) in the format `xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx`. " +
					"Provide this to convert to a SID, or leave empty to convert from a SID.",
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.MatchRoot("sid")),
					stringvalidator.RegexMatches(
						regexp.MustCompile(constants.GuidRegex),
						"Object ID must be a valid GUID in the format xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
					),
				},
			},
		},
	}
}
