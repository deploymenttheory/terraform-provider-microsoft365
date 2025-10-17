// https://learn.microsoft.com/en-us/graph/api/resources/customsecurityattributedefinition?view=graph-rest-beta

package graphBetaCustomSecurityAttributeDefinition

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type CustomSecurityAttributeDefinitionResourceModel struct {
	ID                      types.String   `tfsdk:"id"`
	AttributeSet            types.String   `tfsdk:"attribute_set"`
	Name                    types.String   `tfsdk:"name"`
	IsCollection            types.Bool     `tfsdk:"is_collection"`
	IsSearchable            types.Bool     `tfsdk:"is_searchable"`
	Status                  types.String   `tfsdk:"status"`
	Type                    types.String   `tfsdk:"type"`
	UsePreDefinedValuesOnly types.Bool     `tfsdk:"use_pre_defined_values_only"`
	Description             types.String   `tfsdk:"description"`
	Timeouts                timeouts.Value `tfsdk:"timeouts"`
}
