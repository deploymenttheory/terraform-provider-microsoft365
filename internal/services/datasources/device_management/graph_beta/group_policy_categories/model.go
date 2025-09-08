package graphBetaGroupPolicyCategories

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// GroupPolicyCategoryDataSourceModel represents the main data source model
type GroupPolicyCategoryDataSourceModel struct {
	ID            types.String                   `tfsdk:"id"`
	SettingName   types.String                   `tfsdk:"setting_name"`
	Category      *GroupPolicyCategoryModel      `tfsdk:"category"`
	Definition    *GroupPolicyDefinitionModel    `tfsdk:"definition"`
	Presentations []GroupPolicyPresentationModel `tfsdk:"presentations"`
	Timeouts      timeouts.Value                 `tfsdk:"timeouts"`
}

// GroupPolicyCategoryModel represents a group policy category
type GroupPolicyCategoryModel struct {
	ID              types.String                    `tfsdk:"id"`
	DisplayName     types.String                    `tfsdk:"display_name"`
	IsRoot          types.Bool                      `tfsdk:"is_root"`
	IngestionSource types.String                    `tfsdk:"ingestion_source"`
	Parent          *GroupPolicyCategoryParentModel `tfsdk:"parent"`
}

// GroupPolicyCategoryParentModel represents a parent group policy category
type GroupPolicyCategoryParentModel struct {
	ID          types.String `tfsdk:"id"`
	DisplayName types.String `tfsdk:"display_name"`
	IsRoot      types.Bool   `tfsdk:"is_root"`
}

// GroupPolicyDefinitionModel represents a group policy definition
type GroupPolicyDefinitionModel struct {
	ID                    types.String `tfsdk:"id"`
	DisplayName           types.String `tfsdk:"display_name"`
	CategoryPath          types.String `tfsdk:"category_path"`
	ClassType             types.String `tfsdk:"class_type"`
	PolicyType            types.String `tfsdk:"policy_type"`
	Version               types.String `tfsdk:"version"`
	HasRelatedDefinitions types.Bool   `tfsdk:"has_related_definitions"`
	ExplainText           types.String `tfsdk:"explain_text"`
	SupportedOn           types.String `tfsdk:"supported_on"`
	GroupPolicyCategoryID types.String `tfsdk:"group_policy_category_id"`
	MinDeviceCSPVersion   types.String `tfsdk:"min_device_csp_version"`
	MinUserCSPVersion     types.String `tfsdk:"min_user_csp_version"`
	LastModifiedDateTime  types.String `tfsdk:"last_modified_date_time"`
}

// GroupPolicyPresentationModel represents a group policy presentation
type GroupPolicyPresentationModel struct {
	ID                   types.String `tfsdk:"id"`
	ODataType            types.String `tfsdk:"odata_type"`
	Label                types.String `tfsdk:"label"`
	Required             types.Bool   `tfsdk:"required"`
	LastModifiedDateTime types.String `tfsdk:"last_modified_date_time"`

	// For dropdown lists
	DefaultItem *GroupPolicyPresentationItemModel  `tfsdk:"default_item"`
	Items       []GroupPolicyPresentationItemModel `tfsdk:"items"`

	// For text boxes
	DefaultValue types.String `tfsdk:"default_value"`
	MaxLength    types.Int64  `tfsdk:"max_length"`

	// For checkboxes
	DefaultChecked types.Bool `tfsdk:"default_checked"`

	// For decimal text boxes
	DefaultDecimalValue types.Int64 `tfsdk:"default_decimal_value"`
	MinValue            types.Int64 `tfsdk:"min_value"`
	MaxValue            types.Int64 `tfsdk:"max_value"`
	Spin                types.Bool  `tfsdk:"spin"`
	SpinStep            types.Int64 `tfsdk:"spin_step"`

	// For list boxes
	// REF: https://learn.microsoft.com/en-us/graph/api/intune-grouppolicy-grouppolicypresentationlistbox-get?view=graph-rest-beta
	ExplicitValue types.Bool   `tfsdk:"explicit_value"`
	ValuePrefix   types.String `tfsdk:"value_prefix"`
}

// GroupPolicyPresentationItemModel represents an item in a group policy presentation
type GroupPolicyPresentationItemModel struct {
	DisplayName types.String `tfsdk:"display_name"`
	Value       types.String `tfsdk:"value"`
}
