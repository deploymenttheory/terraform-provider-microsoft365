package graphBetaAttributeSet

import (
	"context"
	"regexp"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	planmodifiers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/plan_modifiers"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName  = "graph_beta_identity_and_access_attribute_set"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &AttributeSetResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &AttributeSetResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &AttributeSetResource{}
)

func NewAttributeSetResource() resource.Resource {
	return &AttributeSetResource{
		ReadPermissions: []string{
			"CustomSecAttributeDefinition.Read.All",
		},
		WritePermissions: []string{
			"CustomSecAttributeDefinition.ReadWrite.All",
		},
		ResourcePath: "/directory/attributeSets",
	}
}

type AttributeSetResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ProviderTypeName string
	TypeName         string
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

// Metadata returns the resource type name.
func (r *AttributeSetResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	r.ProviderTypeName = req.ProviderTypeName
	r.TypeName = ResourceName
	resp.TypeName = r.FullTypeName()
}

// FullTypeName returns the full resource type name in the format "providername_resourcename".
func (r *AttributeSetResource) FullTypeName() string {
	return r.ProviderTypeName + "_" + r.TypeName
}

// Configure sets the client for the resource.
func (r *AttributeSetResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, constants.PROVIDER_NAME+"_"+ResourceName)
}

// ImportState imports the resource state.
func (r *AttributeSetResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *AttributeSetResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages Microsoft Entra custom security attribute sets using the `/directory/attributeSets` endpoint. Attribute sets provide a way to organize and group custom security attributes within a tenant, allowing administrators to define collections of related attributes with configurable limits.\n\n**Note:** Attribute sets cannot be deleted once created. When removed from Terraform configuration, the resource will only be removed from Terraform state but will remain in Microsoft Entra.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Identifier for the attribute set that is unique within a tenant. Can be up to 32 characters long and include Unicode characters. Cannot contain spaces or special characters. Cannot be changed later. Case sensitive. Required.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 32),
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[^\s]+$`),
						"ID cannot contain spaces or special characters",
					),
				},
				PlanModifiers: []planmodifier.String{
					planmodifiers.RequiresReplaceString(),
				},
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description of the attribute set. Can be up to 128 characters long and include Unicode characters. Can be changed later. Optional.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(128),
				},
			},
			"max_attributes_per_set": schema.Int32Attribute{
				MarkdownDescription: "Maximum number of custom security attributes that can be defined in this attribute set. The value must be between 1 and 500. If not specified, the administrator can add up to the maximum of 500 active attributes per tenant. Can be changed later.",
				Optional:            true,
				Validators: []validator.Int32{
					int32validator.Between(1, 500),
				},
			},
			"timeouts": commonschema.Timeouts(ctx),
		},
	}
}
