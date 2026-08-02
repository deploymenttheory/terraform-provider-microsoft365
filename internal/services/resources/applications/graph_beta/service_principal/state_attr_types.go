package graphBetaServicePrincipal

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Attribute type maps for computed nested objects

var keyCredentialAttrTypes = map[string]attr.Type{
	"custom_key_identifier": types.StringType,
	"display_name":          types.StringType,
	"end_date_time":         types.StringType,
	"key_id":                types.StringType,
	"start_date_time":       types.StringType,
	"type":                  types.StringType,
	"usage":                 types.StringType,
}

var passwordCredentialAttrTypes = map[string]attr.Type{
	"custom_key_identifier": types.StringType,
	"display_name":          types.StringType,
	"end_date_time":         types.StringType,
	"hint":                  types.StringType,
	"key_id":                types.StringType,
	"start_date_time":       types.StringType,
}
