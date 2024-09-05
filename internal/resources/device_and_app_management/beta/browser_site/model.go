package graphbetabrowsersite

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type BrowserSiteResourceModel struct {
	ID                   types.String              `tfsdk:"id"`
	AllowRedirect        types.Bool                `tfsdk:"allow_redirect"`
	Comment              types.String              `tfsdk:"comment"`
	CompatibilityMode    types.String              `tfsdk:"compatibility_mode"`
	CreatedDateTime      types.String              `tfsdk:"created_date_time"`
	DeletedDateTime      types.String              `tfsdk:"deleted_date_time"`
	History              []BrowserSiteHistoryModel `tfsdk:"history"`
	LastModifiedBy       IdentitySetModel          `tfsdk:"last_modified_by"`
	LastModifiedDateTime types.String              `tfsdk:"last_modified_date_time"`
	MergeType            types.String              `tfsdk:"merge_type"`
	Status               types.String              `tfsdk:"status"`
	TargetEnvironment    types.String              `tfsdk:"target_environment"`
	WebUrl               types.String              `tfsdk:"web_url"`
	Timeouts             timeouts.Value            `tfsdk:"timeouts"`
}

type BrowserSiteHistoryModel struct {
	AllowRedirect     types.Bool       `tfsdk:"allow_redirect"`
	Comment           types.String     `tfsdk:"comment"`
	CompatibilityMode types.String     `tfsdk:"compatibility_mode"`
	LastModifiedBy    IdentitySetModel `tfsdk:"last_modified_by"`
	MergeType         types.String     `tfsdk:"merge_type"`
	PublishedDateTime types.String     `tfsdk:"published_date_time"`
	TargetEnvironment types.String     `tfsdk:"target_environment"`
}

type IdentitySetModel struct {
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
	DisplayName types.String `tfsdk:"display_name"`
	ID          types.String `tfsdk:"id"`
	TenantID    types.String `tfsdk:"tenant_id"`
}
