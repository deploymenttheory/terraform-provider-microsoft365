// REF: include the graph api docs link for the resource type here
package graphVersionResourceTemplate

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ResourceTemplateResourceModel struct {
	ID       types.String   `tfsdk:"id"`
	etc      types.String   `tfsdk:"etc"`
	Timeouts timeouts.Value `tfsdk:"timeouts"`
}
