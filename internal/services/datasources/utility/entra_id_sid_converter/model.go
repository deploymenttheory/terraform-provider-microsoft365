package utilityEntraIdSidConverter

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type EntraIdSidConverterDataSourceModel struct {
	Id       types.String   `tfsdk:"id"`
	Sid      types.String   `tfsdk:"sid"`
	ObjectId types.String   `tfsdk:"object_id"`
	Timeouts timeouts.Value `tfsdk:"timeouts"`
}
