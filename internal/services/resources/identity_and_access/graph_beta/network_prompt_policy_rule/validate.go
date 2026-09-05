package graphBetaNetworkPromptPolicyRule

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"net/url"
)

func (r *NetworkPromptPolicyRuleResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var data NetworkPromptPolicyRuleResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() || data.ConversationSchemes.IsNull() || data.ConversationSchemes.IsUnknown() {
		return
	}
	for i, element := range data.ConversationSchemes.Elements() {
		p := path.Root("conversation_schemes").AtListIndex(i)
		if element.IsNull() {
			resp.Diagnostics.AddAttributeError(p, "Invalid conversation scheme", "A scheme cannot be null.")
			continue
		}
		if element.IsUnknown() {
			continue
		}
		var scheme ConversationSchemeModel
		object, ok := element.(types.Object)
		if !ok {
			resp.Diagnostics.AddAttributeError(p, "Invalid conversation scheme", "Expected an object.")
			continue
		}
		resp.Diagnostics.Append(object.As(ctx, &scheme, basetypes.ObjectAsOptions{})...)
		if resp.Diagnostics.HasError() {
			return
		}
		switch scheme.Type.ValueString() {
		case schemeTypeCustom:
			if scheme.URL.IsNull() {
				resp.Diagnostics.AddAttributeError(p.AtName("url"), "Missing custom URL", "Custom schemes require url.")
			}
			if !scheme.URL.IsNull() && !scheme.URL.IsUnknown() {
				u, err := url.Parse(scheme.URL.ValueString())
				if err != nil || u.Host == "" || u.Path == "" || (u.Scheme != "http" && u.Scheme != "https") {
					resp.Diagnostics.AddAttributeError(p.AtName("url"), "Invalid custom URL", "Use an absolute HTTP or HTTPS URL with a path, such as /chat or /.")
				}
			}
			if !scheme.SchemeName.IsNull() {
				resp.Diagnostics.AddAttributeError(p.AtName("scheme_name"), "Conflicting scheme attribute", "Custom schemes cannot specify scheme_name.")
			}
		case schemeTypePredefined:
			if scheme.SchemeName.IsNull() {
				resp.Diagnostics.AddAttributeError(p.AtName("scheme_name"), "Missing predefined scheme", "Predefined schemes require scheme_name.")
			}
			if !scheme.URL.IsNull() || !scheme.JSONPath.IsNull() {
				resp.Diagnostics.AddAttributeError(p, "Conflicting scheme attributes", "Predefined schemes cannot specify url or json_path.")
			}
		}
	}
}
