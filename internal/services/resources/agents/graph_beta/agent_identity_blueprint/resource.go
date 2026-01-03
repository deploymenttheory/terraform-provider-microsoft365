package graphBetaApplicationsAgentIdentityBlueprint

import (
	"context"
	"fmt"
	"strings"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName  = "microsoft365_graph_beta_agents_agent_identity_blueprint"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &AgentIdentityBlueprintResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &AgentIdentityBlueprintResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &AgentIdentityBlueprintResource{}
)

func NewAgentIdentityBlueprintResource() resource.Resource {
	return &AgentIdentityBlueprintResource{
		ReadPermissions: []string{
			"AgentIdentityBlueprint.Read.All",
			"Application.Read.All",
			"Directory.Read.All",
		},
		WritePermissions: []string{
			"AgentIdentityBlueprint.Create",
			"AgentIdentityBlueprint.ReadWrite.All",
			"Directory.ReadWrite.All",
			"AgentIdentityBlueprint.AddRemoveCreds.All",
			"AgentIdentityBlueprint.UpdateBranding.All",
			"AgentIdentity.DeleteRestore.All", // Needed for hard deletion
		},
		ResourcePath: "/applications",
	}
}

type AgentIdentityBlueprintResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

// Metadata returns the resource type name.
func (r *AgentIdentityBlueprintResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

// Configure sets the client for the resource.
func (r *AgentIdentityBlueprintResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

// ImportState handles importing the resource with an extended ID format.
//
// Supported formats:
//   - Simple:   "resource_id" (hard_delete defaults to false)
//   - Extended: "resource_id:hard_delete=true" or "resource_id:hard_delete=false"
//
// Example:
//
//	terraform import microsoft365_graph_beta_agents_agent_identity_blueprint.example "12345678-1234-1234-1234-123456789012:hard_delete=true"
func (r *AgentIdentityBlueprintResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.Split(req.ID, ":")
	resourceID := idParts[0]
	hardDelete := false // Default to soft delete for safety

	if len(idParts) > 1 {
		for _, part := range idParts[1:] {
			if strings.HasPrefix(part, "hard_delete=") {
				value := strings.TrimPrefix(part, "hard_delete=")
				switch strings.ToLower(value) {
				case "true":
					hardDelete = true
				case "false":
					hardDelete = false
				default:
					resp.Diagnostics.AddError(
						"Invalid Import ID",
						fmt.Sprintf("Invalid hard_delete value '%s'. Must be 'true' or 'false'.", value),
					)
					return
				}
			}
		}
	}

	tflog.Info(ctx, fmt.Sprintf("Importing %s with ID: %s, hard_delete: %t", ResourceName, resourceID, hardDelete))

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), resourceID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("hard_delete"), hardDelete)...)
}

// Schema returns the schema for the resource.
func (r *AgentIdentityBlueprintResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages an Agent Identity Blueprint in Microsoft Entra ID using the `/applications/microsoft.graph.agentIdentityBlueprint` endpoint. " +
			"An agent identity blueprint serves as a template for creating agent identities within the Microsoft Entra ID ecosystem. " +
			"This resource inherits from the application resource type.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier for the agent identity blueprint. This property is referred to as Object ID in the Microsoft Entra admin center. Read-only.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"app_id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier for the application that is assigned to the agent identity blueprint by Microsoft Entra ID. Also known as Application (client) ID. Read-only.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"display_name": schema.StringAttribute{
				MarkdownDescription: "The display name for the agent identity blueprint. Maximum length is 256 characters. Required.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 256),
				},
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Free text field to provide a description of the agent identity blueprint to end users. Maximum length is 1,024 characters.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(1024),
				},
			},
			"sign_in_audience": schema.StringAttribute{
				MarkdownDescription: "Specifies the Microsoft accounts that are supported for the current application. Supported values are: `AzureADMyOrg` (Single tenant), " +
					" the following values from testing don't work: `AzureADMultipleOrgs` (Multi-tenant), `AzureADandPersonalMicrosoftAccount` (Multi-tenant and personal accounts), `PersonalMicrosoftAccount` (Personal accounts only).",
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.OneOf(
						"AzureADMyOrg", // appears to always be this value ?
						// "AzureADMultipleOrgs",
						// "AzureADandPersonalMicrosoftAccount",
						// "PersonalMicrosoftAccount",
					),
				},
			},
			"tags": schema.SetAttribute{
				MarkdownDescription: "Custom strings that can be used to categorize and identify the agent identity blueprint.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
			},
			"sponsor_user_ids": schema.SetAttribute{
				MarkdownDescription: "The user IDs of the sponsors for the agent identity blueprint. At least one sponsor is " +
					"required when creating an agent identity blueprint. Sponsors are users who can approve or oversee the blueprint.",
				Required:    true,
				ElementType: types.StringType,
			},
			"owner_user_ids": schema.SetAttribute{
				MarkdownDescription: "The user IDs of the owners for the agent identity blueprint. At least one owner is required " +
					"when creating an agent identity blueprint. Owners are users who have full control over the blueprint.",
				Required:    true,
				ElementType: types.StringType,
			},
			"hard_delete": schema.BoolAttribute{
				MarkdownDescription: "When set to `true`, the resource will be permanently deleted from the Entra ID (hard delete) " +
					"rather than being moved to deleted items (soft delete). This prevents the resource from being restored " +
					"and immediately frees up the resource name for reuse. When `false` (default), the resource is soft deleted and can be restored within 30 days. " +
					"Note: This field defaults to `false` on import since the API does not return this value.",
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(false),
			},
			"timeouts": commonschema.ResourceTimeouts(ctx),
		},
	}
}
