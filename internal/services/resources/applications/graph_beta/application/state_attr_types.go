package graphBetaApplication

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Attribute type maps for nested objects
// These define the structure for types.Object fields

var SignInAudienceRestrictionsAttrTypes = map[string]attr.Type{
	"odata_type":             types.StringType,
	"is_home_tenant_allowed": types.BoolType,
	"allowed_tenant_ids":     types.SetType{ElemType: types.StringType},
}

var ApplicationApiAttrTypes = map[string]attr.Type{
	"accept_mapped_claims":           types.BoolType,
	"known_client_applications":      types.SetType{ElemType: types.StringType},
	"oauth2_permission_scopes":       types.SetType{ElemType: types.ObjectType{AttrTypes: OAuth2PermissionScopeAttrTypes}},
	"pre_authorized_applications":    types.SetType{ElemType: types.ObjectType{AttrTypes: PreAuthorizedApplicationAttrTypes}},
	"requested_access_token_version": types.Int32Type,
}

var OAuth2PermissionScopeAttrTypes = map[string]attr.Type{
	"admin_consent_description":  types.StringType,
	"admin_consent_display_name": types.StringType,
	"id":                         types.StringType,
	"is_enabled":                 types.BoolType,
	"type":                       types.StringType,
	"user_consent_description":   types.StringType,
	"user_consent_display_name":  types.StringType,
	"value":                      types.StringType,
}

var PreAuthorizedApplicationAttrTypes = map[string]attr.Type{
	"app_id":                   types.StringType,
	"delegated_permission_ids": types.SetType{ElemType: types.StringType},
}

var ApplicationInformationalUrlAttrTypes = map[string]attr.Type{
	"logo_url":              types.StringType,
	"marketing_url":         types.StringType,
	"privacy_statement_url": types.StringType,
	"support_url":           types.StringType,
	"terms_of_service_url":  types.StringType,
}

var ApplicationOptionalClaimsAttrTypes = map[string]attr.Type{
	"access_token": types.SetType{ElemType: types.ObjectType{AttrTypes: OptionalClaimAttrTypes}},
	"id_token":     types.SetType{ElemType: types.ObjectType{AttrTypes: OptionalClaimAttrTypes}},
	"saml2_token":  types.SetType{ElemType: types.ObjectType{AttrTypes: OptionalClaimAttrTypes}},
}

var OptionalClaimAttrTypes = map[string]attr.Type{
	"additional_properties": types.SetType{ElemType: types.StringType},
	"essential":             types.BoolType,
	"name":                  types.StringType,
	"source":                types.StringType,
}

var ApplicationParentalControlSettingsAttrTypes = map[string]attr.Type{
	"countries_blocked_for_minors": types.SetType{ElemType: types.StringType},
	"legal_age_group_rule":         types.StringType,
}

var ApplicationPublicClientAttrTypes = map[string]attr.Type{
	"redirect_uris": types.SetType{ElemType: types.StringType},
}

var ApplicationSpaAttrTypes = map[string]attr.Type{
	"redirect_uris": types.SetType{ElemType: types.StringType},
}

var ApplicationWebAttrTypes = map[string]attr.Type{
	"home_page_url":           types.StringType,
	"logout_url":              types.StringType,
	"redirect_uris":           types.SetType{ElemType: types.StringType},
	"implicit_grant_settings": types.ObjectType{AttrTypes: ImplicitGrantSettingsAttrTypes},
	"redirect_uri_settings":   types.SetType{ElemType: types.ObjectType{AttrTypes: RedirectUriSettingsAttrTypes}},
}

var ImplicitGrantSettingsAttrTypes = map[string]attr.Type{
	"enable_access_token_issuance": types.BoolType,
	"enable_id_token_issuance":     types.BoolType,
}

// Attribute types for sets that contain objects

var KeyCredentialAttrTypes = map[string]attr.Type{
	"custom_key_identifier": types.StringType,
	"display_name":          types.StringType,
	"end_date_time":         types.StringType,
	"key":                   types.StringType,
	"key_id":                types.StringType,
	"start_date_time":       types.StringType,
	"type":                  types.StringType,
	"usage":                 types.StringType,
}

var PasswordCredentialAttrTypes = map[string]attr.Type{
	"custom_key_identifier": types.StringType,
	"display_name":          types.StringType,
	"end_date_time":         types.StringType,
	"hint":                  types.StringType,
	"key_id":                types.StringType,
	"secret_text":           types.StringType,
	"start_date_time":       types.StringType,
}

var AppRoleAttrTypes = map[string]attr.Type{
	"allowed_member_types": types.SetType{ElemType: types.StringType},
	"description":          types.StringType,
	"display_name":         types.StringType,
	"id":                   types.StringType,
	"is_enabled":           types.BoolType,
	"origin":               types.StringType,
	"value":                types.StringType,
}

var RequiredResourceAccessAttrTypes = map[string]attr.Type{
	"resource_app_id": types.StringType,
	"resource_access": types.SetType{ElemType: types.ObjectType{AttrTypes: ResourceAccessAttrTypes}},
}

var ResourceAccessAttrTypes = map[string]attr.Type{
	"id":   types.StringType,
	"type": types.StringType,
}

var RedirectUriSettingsAttrTypes = map[string]attr.Type{
	"uri":   types.StringType,
	"index": types.Int32Type,
}
