package graphBetaCustomSecurityAttributeAllowedValue

import (
	"context"
	"fmt"
	"strings"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	planmodifiers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/plan_modifiers"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName  = "microsoft365_graph_beta_identity_and_access_custom_security_attribute_allowed_value"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &CustomSecurityAttributeAllowedValueResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &CustomSecurityAttributeAllowedValueResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &CustomSecurityAttributeAllowedValueResource{}
)

func NewCustomSecurityAttributeAllowedValueResource() resource.Resource {
	return &CustomSecurityAttributeAllowedValueResource{
		ReadPermissions: []string{
			"CustomSecAttributeDefinition.Read.All",
		},
		WritePermissions: []string{
			"CustomSecAttributeDefinition.ReadWrite.All",
		},
		ResourcePath: "/directory/customSecurityAttributeDefinitions/{customSecurityAttributeDefinitionId}/allowedValues",
	}
}

type CustomSecurityAttributeAllowedValueResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

// Metadata returns the resource type name.
func (r *CustomSecurityAttributeAllowedValueResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

// Configure sets the client for the resource.
func (r *CustomSecurityAttributeAllowedValueResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

// ImportState imports the resource state.
func (r *CustomSecurityAttributeAllowedValueResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Import format: {customSecurityAttributeDefinitionId}/{id}
	// Example: Engineering_Project/Alpine
	idParts := strings.Split(req.ID, "/")
	if len(idParts) != 2 {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			fmt.Sprintf("Expected import ID in format 'customSecurityAttributeDefinitionId/id', got: %s", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("custom_security_attribute_definition_id"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), idParts[1])...)
}

func (r *CustomSecurityAttributeAllowedValueResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages Microsoft Entra custom security attribute allowed values using the `/directory/customSecurityAttributeDefinitions/{customSecurityAttributeDefinitionId}/allowedValues` endpoint. Allowed values represent predefined values that can be assigned to custom security attributes.\n\n**Note:** You can define up to 100 allowed values per custom security attribute definition. Allowed values cannot be renamed or deleted, but they can be deactivated.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Identifier for the predefined value. Can be up to 64 characters long and include Unicode characters. Can include spaces, but some special characters aren't allowed. Cannot be changed later. Case sensitive.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.RequiresReplaceString(),
				},
			},
			"custom_security_attribute_definition_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the custom security attribute definition that this allowed value belongs to. Format: 'attributeSet_attributeName'.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.RequiresReplaceString(),
				},
			},
			"is_active": schema.BoolAttribute{
				MarkdownDescription: "Indicates whether the predefined value is active or deactivated. If set to false, this predefined value cannot be assigned to any more supported directory objects. Can be changed later.",
				Required:            true,
			},
			"timeouts": commonschema.ResourceTimeouts(ctx),
		},
	}
}
