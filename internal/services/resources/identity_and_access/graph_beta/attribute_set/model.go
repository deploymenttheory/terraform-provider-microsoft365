// https://learn.microsoft.com/en-us/graph/api/resources/attributeset?view=graph-rest-beta

package graphBetaAttributeSet

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type AttributeSetResourceModel struct {
	ID                  types.String   `tfsdk:"id"`
	Description         types.String   `tfsdk:"description"`
	MaxAttributesPerSet types.Int32    `tfsdk:"max_attributes_per_set"`
	Timeouts            timeouts.Value `tfsdk:"timeouts"`
}
