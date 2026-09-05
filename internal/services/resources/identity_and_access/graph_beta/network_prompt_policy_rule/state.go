package graphBetaNetworkPromptPolicyRule

import (
	"context"
	"fmt"
	"math"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
)

func MapRemoteStateToTerraform(
	ctx context.Context,
	data *NetworkPromptPolicyRuleResourceModel,
	remote *promptPolicyRuleResponse,
) error {
	if err := validateRemoteRule(remote); err != nil {
		return err
	}
	var prior []ConversationSchemeModel
	if !data.ConversationSchemes.IsNull() && !data.ConversationSchemes.IsUnknown() {
		if diags := data.ConversationSchemes.ElementsAs(ctx, &prior, false); diags.HasError() {
			return fmt.Errorf("%w: previous schemes: %v", errInvalidResponse, diags)
		}
	}
	schemes := make([]attr.Value, 0, len(remote.conditions.schemes))
	for i, scheme := range remote.conditions.schemes {
		if scheme.odataType == nil {
			return fmt.Errorf("%w: scheme missing @odata.type", errInvalidResponse)
		}
		fields := map[string]attr.Value{
			"type":        types.StringNull(),
			"url":         types.StringNull(),
			"json_path":   types.StringNull(),
			"scheme_name": types.StringNull(),
		}
		switch strings.TrimPrefix(*scheme.odataType, "#") {
		case "microsoft.graph.networkaccess.customConversationScheme":
			if scheme.url == nil || *scheme.url == "" || scheme.jsonPath == nil {
				return fmt.Errorf("%w: incomplete custom scheme", errInvalidResponse)
			}
			fields["type"] = types.StringValue(schemeTypeCustom)
			fields["url"] = types.StringValue(*scheme.url)
			fields["json_path"] = types.StringValue(*scheme.jsonPath)
			// Only normalize the observed empty API default for an omitted value on the same scheme.
			if *scheme.jsonPath == "" && i < len(prior) &&
				prior[i].Type.ValueString() == schemeTypeCustom &&
				prior[i].URL.ValueString() == *scheme.url &&
				prior[i].JSONPath.IsNull() {
				fields["json_path"] = types.StringNull()
			}
		case "microsoft.graph.networkaccess.predefinedConversationScheme":
			if scheme.schemeName == nil {
				return fmt.Errorf("%w: predefined scheme missing name", errInvalidResponse)
			}
			fields["type"] = types.StringValue(schemeTypePredefined)
			fields["scheme_name"] = types.StringValue(*scheme.schemeName)
		default:
			return fmt.Errorf(
				"%w: unsupported scheme type %q",
				errInvalidResponse,
				*scheme.odataType,
			)
		}
		schemes = append(
			schemes,
			types.ObjectValueMust(conversationSchemeObjectType().AttrTypes, fields),
		)
	}
	data.ID = convert.GraphToFrameworkString(remote.id)
	data.Name = convert.GraphToFrameworkString(remote.name)
	data.Description = convert.GraphToFrameworkString(remote.description)
	data.Action = convert.GraphToFrameworkString(remote.action)
	data.Priority = types.Int32Value(int32(*remote.priority))
	data.Enabled = types.BoolValue(*remote.status == "enabled")
	data.Status = types.StringValue(*remote.status)
	data.PromptLogging = types.StringValue(*remote.promptLogging)
	data.ScanResult = types.StringValue(*remote.conditions.scanResult)
	data.ConversationSchemes = types.ListValueMust(conversationSchemeObjectType(), schemes)
	return nil
}

func conversationSchemeObjectType() types.ObjectType {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"type":        types.StringType,
			"url":         types.StringType,
			"json_path":   types.StringType,
			"scheme_name": types.StringType,
		},
	}
}

func validateRemoteRule(remote *promptPolicyRuleResponse) error {
	if remote == nil || remote.id == nil || *remote.id == "" || remote.name == nil ||
		remote.action == nil ||
		remote.priority == nil ||
		remote.status == nil ||
		remote.promptLogging == nil ||
		remote.odataType == nil ||
		remote.conditions == nil ||
		remote.conditions.scanResult == nil ||
		!remote.conditions.schemesPresent {
		return fmt.Errorf("%w: missing required fields", errInvalidResponse)
	}
	if strings.TrimPrefix(*remote.odataType, "#") != "microsoft.graph.networkaccess.promptRule" {
		return fmt.Errorf("%w: unsupported rule type %q", errInvalidResponse, *remote.odataType)
	}
	if *remote.priority < 100 || *remote.priority > math.MaxInt32 {
		return fmt.Errorf(
			"%w: priority %d outside supported range",
			errInvalidResponse,
			*remote.priority,
		)
	}
	if *remote.status != "enabled" && *remote.status != "disabled" {
		return fmt.Errorf(
			"%w: unsupported status %q; enabled cannot be determined",
			errInvalidResponse,
			*remote.status,
		)
	}
	return nil
}
