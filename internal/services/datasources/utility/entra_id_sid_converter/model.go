package entra_id_sid_converter

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type EntraIdSidConverterDataSourceModel struct {
	Id       types.String `tfsdk:"id"`
	Sid      types.String `tfsdk:"sid"`
	ObjectId types.String `tfsdk:"object_id"`
}
