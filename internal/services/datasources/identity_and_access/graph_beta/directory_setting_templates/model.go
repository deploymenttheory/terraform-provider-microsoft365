// REF: https://learn.microsoft.com/en-us/graph/api/resources/directorysettingtemplate?view=graph-rest-beta

package graphBetaDirectorySettingTemplates

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// DirectorySettingTemplatesDataSourceModel represents the Terraform data source model for directory setting templates
type DirectorySettingTemplatesDataSourceModel struct {
	ID                        types.String                    `tfsdk:"id"`
	FilterType                types.String                    `tfsdk:"filter_type"`                 // Required field to specify how to filter (all, id, display_name)
	FilterValue               types.String                    `tfsdk:"filter_value"`                // Value to filter by (not used for "all")
	DirectorySettingTemplates []DirectorySettingTemplateModel `tfsdk:"directory_setting_templates"` // List of all templates
	Timeouts                  timeouts.Value                  `tfsdk:"timeouts"`
}

// DirectorySettingTemplateModel represents an individual directory setting template
type DirectorySettingTemplateModel struct {
	ID          types.String                `tfsdk:"id"`
	Description types.String                `tfsdk:"description"`
	DisplayName types.String                `tfsdk:"display_name"`
	Values      []SettingTemplateValueModel `tfsdk:"values"`
}

// SettingTemplateValueModel represents a setting template value
type SettingTemplateValueModel struct {
	Name         types.String `tfsdk:"name"`
	Type         types.String `tfsdk:"type"`
	DefaultValue types.String `tfsdk:"default_value"`
	Description  types.String `tfsdk:"description"`
}
