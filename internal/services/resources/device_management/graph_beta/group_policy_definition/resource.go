package graphBetaGroupPolicyDefinition

import (
	"context"
	"fmt"
	"strings"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	planmodifiers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/plan_modifiers"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/devicemanagement"
)

const (
	ResourceName  = "microsoft365_graph_beta_device_management_group_policy_definition"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &GroupPolicyDefinitionResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &GroupPolicyDefinitionResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &GroupPolicyDefinitionResource{}

	// Enables plan modification/diff suppression
	_ resource.ResourceWithModifyPlan = &GroupPolicyDefinitionResource{}
)

// NewGroupPolicyDefinitionResource returns a new instance of the resource
func NewGroupPolicyDefinitionResource() resource.Resource {
	return &GroupPolicyDefinitionResource{
		ReadPermissions: []string{
			"DeviceManagementConfiguration.Read.All",
		},
		WritePermissions: []string{
			"DeviceManagementConfiguration.ReadWrite.All",
		},
		ResourcePath: "/deviceManagement/groupPolicyDefinitions",
	}
}

// GroupPolicyDefinitionResource defines the resource implementation
type GroupPolicyDefinitionResource struct {
	client *msgraphbetasdk.GraphServiceClient

	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

// Metadata returns the resource type name
func (r *GroupPolicyDefinitionResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

// Configure sets the provider configured client
func (r *GroupPolicyDefinitionResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

// ImportState imports the resource state by performing multi-endpoint resolution to provide a complete
// user experience when importing existing Group Policy definitions.
//
// # Why Multi-Endpoint Resolution
//
// This resource requires users to specify policy_name, class_type, and category_path during creation
// because these fields are used to resolve the policy from Microsoft's catalog of thousands of
// available Group Policy definitions. To maintain consistency and provide the same user experience
// for imported resources, we fetch these same metadata fields during import so that:
//
//  1. Imported resources have the same attribute set as created resources
//  2. Users can immediately manage imported resources without manual configuration
//  3. The resource resolver can continue to function for updates (it requires these fields)
//  4. Changes to policy_name, class_type, or category_path trigger appropriate resource replacements
//     (these fields are marked with RequiresReplace)
//
// # API Endpoint Resolution Strategy
//
// Due to Microsoft Graph API response structure, we query multiple endpoints to gather complete metadata:
//
//  1. GET /groupPolicyConfigurations/{configID}/definitionValues?$expand=definition
//     - Returns the definition value instance (the configured policy)
//     - The $expand=definition parameter includes basic definition metadata
//     - However, the expanded definition object has categoryPath as null (API limitation)
//     - Without $expand, the definition object is absent entirely from the response
//
//  2. GET /groupPolicyDefinitions/{definitionID}
//     - Returns the complete definition from Microsoft's policy catalog
//     - This is the only endpoint that provides categoryPath populated correctly
//     - Also ensures we get the authoritative displayName and classType
//
// The two-call pattern mirrors how the resource's resolver works during normal operations, maintaining
// consistency between import and create workflows.
//
// # ID Format Handling
//
// This function supports two import ID formats for different scenarios:
//
//  1. Simple format (configID): Used when users manually import a resource
//     Example: terraform import resource.name "abc-123-def-456"
//
//  2. Composite format (configID/definitionValueID): Used by Terraform during re-imports or refreshes
//     Example: "abc-123-def-456/xyz-789-uvw-012"
//
// The composite format is required for:
// - Test framework's CheckDestroy to query the specific definition value within a configuration
// - Resource uniqueness (multiple policies can exist in one configuration)
// - Consistency with how the resource tracks itself during CRUD operations
//
// # Metadata Population
//
// After resolution, we populate the state with:
// - id: Composite format (configID/definitionValueID)
// - group_policy_configuration_id: The parent configuration GUID
// - policy_name: Human-readable policy name from catalog
// - class_type: Either "user" or "machine" policy
// - category_path: Policy category hierarchy (e.g., "\Windows Components\...")
// - enabled: Whether the policy is currently enabled
//
// These fields enable users to modify and manage the imported resource exactly as if they had
// created it through Terraform originally.
func (r *GroupPolicyDefinitionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Import ID can be either:
	// 1. Just the configID (when manually importing)
	// 2. Composite format configID/definitionValueID (when re-importing existing resource)
	importID := req.ID
	var configID string

	// Check if it's a composite ID
	if parts := splitCompositeID(importID); len(parts) == 2 {
		configID = parts[0]
		tflog.Debug(ctx, fmt.Sprintf("Importing with composite ID: %s (using configID: %s)", importID, configID))
	} else {
		configID = importID
		tflog.Debug(ctx, fmt.Sprintf("Importing with simple config ID: %s", configID))
	}

	// Fetch definition values with expanded definition metadata
	// The $expand=definition parameter is required to include the definition object in the response;
	// without it, the API returns definition values with no definition reference at all
	definitionValues, err := r.client.
		DeviceManagement().
		GroupPolicyConfigurations().
		ByGroupPolicyConfigurationId(configID).
		DefinitionValues().
		Get(ctx, &devicemanagement.GroupPolicyConfigurationsItemDefinitionValuesRequestBuilderGetRequestConfiguration{
			QueryParameters: &devicemanagement.GroupPolicyConfigurationsItemDefinitionValuesRequestBuilderGetQueryParameters{
				Expand: []string{"definition"},
			},
		})

	if err != nil {
		resp.Diagnostics.AddError(
			"Error importing resource",
			fmt.Sprintf("Could not fetch definition values from configuration: %s", err.Error()),
		)
		return
	}

	if definitionValues == nil || definitionValues.GetValue() == nil || len(definitionValues.GetValue()) == 0 {
		resp.Diagnostics.AddError(
			"Error importing resource",
			"No definition values found in the specified configuration",
		)
		return
	}

	// Get the first definition value (for this resource, there should typically be one per policy)
	firstDefValue := definitionValues.GetValue()[0]
	definition := firstDefValue.GetDefinition()

	if definition == nil || definition.GetId() == nil {
		resp.Diagnostics.AddError(
			"Error importing resource",
			"Definition not found in definition value",
		)
		return
	}

	// Fetch the complete definition from the catalog to get categoryPath
	// The expanded definition from the previous call has categoryPath as null, so we need
	// to query the groupPolicyDefinitions endpoint directly for the full metadata
	definitionID := *definition.GetId()
	fullDefinition, err := r.client.
		DeviceManagement().
		GroupPolicyDefinitions().
		ByGroupPolicyDefinitionId(definitionID).
		Get(ctx, nil)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error importing resource",
			fmt.Sprintf("Could not fetch definition details: %s", err.Error()),
		)
		return
	}

	displayName := fullDefinition.GetDisplayName()
	classType := fullDefinition.GetClassType()
	categoryPath := fullDefinition.GetCategoryPath()

	if displayName == nil || classType == nil || categoryPath == nil {
		resp.Diagnostics.AddError(
			"Error importing resource",
			"Could not extract policy metadata from definition",
		)
		return
	}

	// Construct the composite ID format: configID/definitionValueID
	// This format is used throughout the resource lifecycle and enables:
	// - Unique identification when multiple policies exist in one configuration
	// - Test framework's CheckDestroy to query the specific definition value
	// - Consistency with how the resource manages itself during CRUD operations
	definitionValueID := firstDefValue.GetId()
	if definitionValueID == nil {
		resp.Diagnostics.AddError(
			"Error importing resource",
			"Definition value ID is missing",
		)
		return
	}

	compositeID := fmt.Sprintf("%s/%s", configID, *definitionValueID)

	resp.State.SetAttribute(ctx, path.Root("id"), compositeID)
	resp.State.SetAttribute(ctx, path.Root("group_policy_configuration_id"), configID)
	resp.State.SetAttribute(ctx, path.Root("policy_name"), *displayName)
	resp.State.SetAttribute(ctx, path.Root("class_type"), classType.String())
	resp.State.SetAttribute(ctx, path.Root("category_path"), *categoryPath)
	resp.State.SetAttribute(ctx, path.Root("enabled"), firstDefValue.GetEnabled())

	tflog.Debug(ctx, fmt.Sprintf("Successfully imported: policy_name=%s, class_type=%s, category_path=%s",
		*displayName, classType.String(), *categoryPath))
}

// Schema returns the resource schema
func (r *GroupPolicyDefinitionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a group policy definition with all its presentation values in Microsoft Intune. " +
			"This resource provides a unified interface for configuring any group policy regardless of presentation types (checkboxes, textboxes, dropdowns, etc.). " +
			"Values are provided as label-value pairs, and the resource automatically handles type conversion based on the policy's catalog definition. " +
			"Uses the `deviceManagement/groupPolicyConfigurations('{groupPolicyConfigurationId}')/updateDefinitionValues` endpoint.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
				MarkdownDescription: "The unique identifier for the group policy definition value",
			},
			"group_policy_configuration_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The unique identifier of the group policy configuration that contains this definition",
			},
			"policy_name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The display name of the group policy definition (e.g., 'Remove Default Microsoft Store packages from the system.')",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"class_type": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The class type of the group policy definition. Must be 'user' or 'machine'",
				Validators: []validator.String{
					stringvalidator.OneOf("user", "machine"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"category_path": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The category path of the group policy definition (e.g., '\\Windows Components\\App Package Deployment'). Used to identify the policy in the catalog",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"enabled": schema.BoolAttribute{
				Required:            true,
				MarkdownDescription: "Whether this group policy definition is enabled (true) or disabled (false)",
			},
			"values": schema.SetNestedAttribute{
				Optional:            true,
				MarkdownDescription: "Set of presentation values for this group policy definition. Each value corresponds to a specific presentation (checkbox, textbox, dropdown, etc.) identified by its label. The resource automatically handles type conversion based on the presentation type in the catalog.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The unique identifier of the presentation template (computed from catalog)",
						},
						"label": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: "The human-readable label of the presentation (e.g., 'Xbox Gaming App', 'Feedback Hub'). Must match a label from the policy's catalog definition",
						},
						"value": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: "The value for this presentation as a string. Format depends on presentation type: 'true'/'false' for checkboxes, numeric strings for decimal fields, plain text for textboxes, etc. The resource validates and converts this based on the presentation type",
						},
					},
				},
			},
			"created_date_time": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The date and time when the definition value was created",
			},
			"last_modified_date_time": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The date and time when the definition value was last modified",
			},
			"timeouts": commonschema.ResourceTimeouts(ctx),
		},
	}
}

// splitCompositeID splits a composite ID in format "configID/definitionValueID" into its parts.
// Returns a slice with either 1 element (simple configID) or 2 elements (composite format).
// This allows ImportState to handle both manual imports (simple ID) and re-imports (composite ID).
func splitCompositeID(id string) []string {
	return strings.Split(id, "/")
}
