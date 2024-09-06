package sharedmodels

import "github.com/hashicorp/terraform-plugin-framework/types"

type IdentitySetModel struct {
	ODataType                types.String  `tfsdk:"odata_type"`
	Application              IdentityModel `tfsdk:"application"`
	ApplicationInstance      IdentityModel `tfsdk:"application_instance"`
	Conversation             IdentityModel `tfsdk:"conversation"`
	ConversationIdentityType IdentityModel `tfsdk:"conversation_identity_type"`
	Device                   IdentityModel `tfsdk:"device"`
	Encrypted                IdentityModel `tfsdk:"encrypted"`
	OnPremises               IdentityModel `tfsdk:"on_premises"`
	Guest                    IdentityModel `tfsdk:"guest"`
	Phone                    IdentityModel `tfsdk:"phone"`
	User                     IdentityModel `tfsdk:"user"`
}

type IdentityModel struct {
	ODataType   types.String `tfsdk:"odata_type"`
	DisplayName types.String `tfsdk:"display_name"`
	ID          types.String `tfsdk:"id"`
	TenantID    types.String `tfsdk:"tenant_id"`
}
