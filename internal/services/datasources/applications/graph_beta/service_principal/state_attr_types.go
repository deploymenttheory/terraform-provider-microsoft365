package graphBetaServicePrincipal

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Attribute type maps for nested objects

var samlSingleSignOnSettingsAttrTypes = map[string]attr.Type{
	"relay_state": types.StringType,
}

var verifiedPublisherAttrTypes = map[string]attr.Type{
	"display_name":          types.StringType,
	"verified_publisher_id": types.StringType,
	"added_date_time":       types.StringType,
}

var infoAttrTypes = map[string]attr.Type{
	"terms_of_service_url":  types.StringType,
	"support_url":           types.StringType,
	"privacy_statement_url": types.StringType,
	"marketing_url":         types.StringType,
	"logo_url":              types.StringType,
}
