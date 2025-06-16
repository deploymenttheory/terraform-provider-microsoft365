package sharedmodels

import "github.com/hashicorp/terraform-plugin-framework/types"

// MimeContentResourceModel struct to hold the mime content configuration
type MimeContentResourceModel struct {
	Type  types.String `tfsdk:"type"`
	Value types.String `tfsdk:"value"`
}
