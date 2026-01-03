package graphBetaCustomSecurityAttributeDefinition

import (
	"context"
	"regexp"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	planmodifiers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/plan_modifiers"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	validate "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/validate/attribute"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName  = "microsoft365_graph_beta_identity_and_access_custom_security_attribute_definition"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &CustomSecurityAttributeDefinitionResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &CustomSecurityAttributeDefinitionResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &CustomSecurityAttributeDefinitionResource{}
)

func NewCustomSecurityAttributeDefinitionResource() resource.Resource {
	return &CustomSecurityAttributeDefinitionResource{
		ReadPermissions: []string{
			"CustomSecAttributeDefinition.Read.All",
		},
		WritePermissions: []string{
			"CustomSecAttributeDefinition.ReadWrite.All",
		},
		ResourcePath: "/directory/customSecurityAttributeDefinitions",
	}
}

type CustomSecurityAttributeDefinitionResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

// Metadata returns the resource type name.
func (r *CustomSecurityAttributeDefinitionResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

// Configure sets the client for the resource.
func (r *CustomSecurityAttributeDefinitionResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

// ImportState imports the resource state.
func (r *CustomSecurityAttributeDefinitionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *CustomSecurityAttributeDefinitionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages Microsoft Entra custom security attribute definitions using the `/directory/customSecurityAttributeDefinitions` endpoint. Custom security attribute definitions define the structure and behavior of custom security attributes that can be assigned to users, groups, and other directory objects.\n\n**Note:** Custom security attribute definitions cannot be deleted once created. When removed from Terraform configuration, the resource will be deactivated by setting its status to 'Deprecated' and then removed from Terraform state. The attribute definition will remain in Microsoft Entra in a deprecated state.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Identifier of the custom security attribute, which is a combination of the attribute set name and the custom security attribute name separated by an underscore (attributeSet_name). The id property is auto generated and cannot be set. Case insensitive. Inherited from entity.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description of the custom security attribute. Can be up to 128 characters long and include Unicode characters. Can be changed later.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(128),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the custom security attribute. Must be unique within an attribute set. Can be up to 32 characters long and include Unicode characters. Cannot contain spaces or special characters. Cannot be changed later. Case insensitive.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 32),
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[^\s]+$`),
						"Name cannot contain spaces or special characters",
					),
				},
				PlanModifiers: []planmodifier.String{
					planmodifiers.RequiresReplaceString(),
				},
			},
			"attribute_set": schema.StringAttribute{
				MarkdownDescription: "Name of the attribute set. Case insensitive.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 32),
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[^\s]+$`),
						"Attribute set cannot contain spaces or special characters",
					),
				},
				PlanModifiers: []planmodifier.String{
					planmodifiers.RequiresReplaceString(),
				},
			},
			"is_collection": schema.BoolAttribute{
				MarkdownDescription: "Indicates whether multiple values can be assigned to the custom security attribute. Cannot be changed later. If type is set to Boolean, isCollection cannot be set to true.",
				Required:            true,
				Validators: []validator.Bool{
					validate.BoolCanOnlyBeFalseWhenStringEquals("type", "Boolean", "When type is set to Boolean, isCollection must be set to false"),
				},
				PlanModifiers: []planmodifier.Bool{
					planmodifiers.NewRequiresReplaceIfChangedBool(),
				},
			},
			"is_searchable": schema.BoolAttribute{
				MarkdownDescription: "Indicates whether custom security attribute values are indexed for searching on objects that are assigned attribute values. Cannot be changed later.",
				Required:            true,
				PlanModifiers: []planmodifier.Bool{
					planmodifiers.NewRequiresReplaceIfChangedBool(),
				},
			},
			"status": schema.StringAttribute{
				MarkdownDescription: "Specifies whether the custom security attribute is active or deactivated. Acceptable values are: Available and Deprecated. Can be changed later.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("Available", "Deprecated"),
				},
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "Data type for the custom security attribute values. Supported types are: Boolean, Integer, and String. Cannot be changed later.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("String", "Integer", "Boolean"),
				},
				PlanModifiers: []planmodifier.String{
					planmodifiers.RequiresReplaceString(),
				},
			},
			"use_pre_defined_values_only": schema.BoolAttribute{
				MarkdownDescription: "Indicates whether only predefined values can be assigned to the custom security attribute. If set to false, free-form values are allowed. Can later be changed from true to false, but cannot be changed from false to true. If type is set to Boolean, usePreDefinedValuesOnly cannot be set to true.",
				Required:            true,
				Validators: []validator.Bool{
					validate.BoolCanOnlyBeFalseWhenStringEquals("type", "Boolean", "When type is set to Boolean, usePreDefinedValuesOnly must be set to false"),
				},
				PlanModifiers: []planmodifier.Bool{
					planmodifiers.NewRequiresReplaceIfFalseToTrue(),
				},
			},
			"timeouts": commonschema.ResourceTimeouts(ctx),
		},
	}
}
