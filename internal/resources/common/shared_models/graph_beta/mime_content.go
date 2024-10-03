package sharedmodels

import "github.com/hashicorp/terraform-plugin-framework/types"

type MimeContentResourceModel struct {
	Type  types.String `tfsdk:"type"`
	Value types.String `tfsdk:"value"`
}
