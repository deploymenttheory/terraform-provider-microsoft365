package sharedmodels

import "github.com/hashicorp/terraform-plugin-framework/types"

type IdentitySetResourceModel struct {
	ODataType                types.String          `tfsdk:"odata_type"`
	Application              IdentityResourceModel `tfsdk:"application"`
	ApplicationInstance      IdentityResourceModel `tfsdk:"application_instance"`
	Conversation             IdentityResourceModel `tfsdk:"conversation"`
	ConversationIdentityType IdentityResourceModel `tfsdk:"conversation_identity_type"`
	Device                   IdentityResourceModel `tfsdk:"device"`
	Encrypted                IdentityResourceModel `tfsdk:"encrypted"`
	OnPremises               IdentityResourceModel `tfsdk:"on_premises"`
	Guest                    IdentityResourceModel `tfsdk:"guest"`
	Phone                    IdentityResourceModel `tfsdk:"phone"`
	User                     IdentityResourceModel `tfsdk:"user"`
}

type IdentityResourceModel struct {
	ODataType   types.String `tfsdk:"odata_type"`
	DisplayName types.String `tfsdk:"display_name"`
	ID          types.String `tfsdk:"id"`
	TenantID    types.String `tfsdk:"tenant_id"`
}
