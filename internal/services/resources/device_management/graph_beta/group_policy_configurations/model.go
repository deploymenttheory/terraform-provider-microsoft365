// REF: https://learn.microsoft.com/en-us/graph/api/resources/grouppolicyconfiguration?view=graph-rest-beta
package graphBetaGroupPolicyConfigurations

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// GroupPolicyConfigurationResourceModel represents the schema for the Group Policy Configuration resource
type GroupPolicyConfigurationResourceModel struct {
	ID                   types.String   `tfsdk:"id"`
	DisplayName          types.String   `tfsdk:"display_name"`
	Description          types.String   `tfsdk:"description"`
	CreatedDateTime      types.String   `tfsdk:"created_date_time"`
	LastModifiedDateTime types.String   `tfsdk:"last_modified_date_time"`
	RoleScopeTagIds      types.Set      `tfsdk:"role_scope_tag_ids"`
	DefinitionValues     types.Set      `tfsdk:"definition_values"` // Definition values - these are the policy settings
	Assignments          types.Set      `tfsdk:"assignments"`
	Timeouts             timeouts.Value `tfsdk:"timeouts"`
}

// DefinitionValueModel represents a single policy definition value
type DefinitionValueModel struct {
	ID                   types.String `tfsdk:"id"`
	Enabled              types.Bool   `tfsdk:"enabled"`
	ConfigurationType    types.String `tfsdk:"configuration_type"`
	CreatedDateTime      types.String `tfsdk:"created_date_time"`
	LastModifiedDateTime types.String `tfsdk:"last_modified_date_time"`
	DisplayName          types.String `tfsdk:"display_name"`        // Human-friendly reference to the policy definition e.g., "Allow users to contact Microsoft for feedback and support"
	ClassType            types.String `tfsdk:"class_type"`          // "user" or "machine"
	CategoryPath         types.String `tfsdk:"category_path"`       // Optional category path for disambiguation e.g., "\\OneDrive"
	DefinitionID         types.String `tfsdk:"definition_id"`       // Internal reference to the policy definition (computed/read-only)
	PresentationValues   types.Set    `tfsdk:"presentation_values"` // Presentation values - these contain the actual configuration values
}

// PresentationValueModel represents a single presentation value with polymorphic odata type support
type PresentationValueModel struct {
	ID                   types.String `tfsdk:"id"`
	CreatedDateTime      types.String `tfsdk:"created_date_time"`
	LastModifiedDateTime types.String `tfsdk:"last_modified_date_time"`

	// Reference to the presentation
	PresentationID types.String `tfsdk:"presentation_id"`

	// OData type to determine which fields are used
	ODataType types.String `tfsdk:"odata_type"`

	// Common value field - the actual meaning depends on odata_type
	Value types.String `tfsdk:"value"`

	// Type-specific fields for different presentation value types
	// Text value (used by groupPolicyPresentationValueText)
	TextValue types.String `tfsdk:"text_value"`

	// Decimal value (used by groupPolicyPresentationValueDecimal)
	DecimalValue types.Int64 `tfsdk:"decimal_value"`

	// Boolean value (used by groupPolicyPresentationValueBoolean)
	BooleanValue types.Bool `tfsdk:"boolean_value"`

	// List values (used by groupPolicyPresentationValueList)
	ListValues types.Set `tfsdk:"list_values"`

	// Multi-text values (used by groupPolicyPresentationValueMultiText)
	MultiTextValues types.Set `tfsdk:"multi_text_values"`
}

// ListValueModel represents a single list value entry
type ListValueModel struct {
	Name  types.String `tfsdk:"name"`
	Value types.String `tfsdk:"value"`
}
