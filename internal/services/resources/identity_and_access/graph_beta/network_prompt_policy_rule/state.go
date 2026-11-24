package graphBetaNetworkPromptPolicyRule

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
)

// MapRemoteStateToTerraform validates the response before replacing the last known state.
func MapRemoteStateToTerraform(_ context.Context, data *NetworkPromptPolicyRuleResourceModel, remote *promptPolicyRuleResponse) error {
	if remote == nil || remote.ID == nil || *remote.ID == "" || remote.Name == nil || remote.Action == nil ||
		remote.Priority == nil || remote.Settings == nil || remote.Settings.Status == nil ||
		remote.PromptLogging == nil || remote.ODataType == nil || remote.Conditions == nil ||
		remote.Conditions.ScanResult == nil || len(remote.Conditions.ConversationSchemes) == 0 {
		return fmt.Errorf("%w: missing required fields", errInvalidResponse)
	}
	if strings.TrimPrefix(*remote.ODataType, "#") != "microsoft.graph.networkaccess.promptRule" {
		return fmt.Errorf("%w: unsupported rule type %q", errInvalidResponse, *remote.ODataType)
	}
	status := *remote.Settings.Status
	if status != "enabled" && status != "disabled" {
		return fmt.Errorf("%w: unsupported status %q; enabled cannot be determined", errInvalidResponse, status)
	}
	if *remote.Priority < 100 || *remote.Priority > math.MaxInt32 {
		return fmt.Errorf("%w: priority %d outside supported range", errInvalidResponse, *remote.Priority)
	}
	// A null conversationSchemes is present but empty, as it was for the Kiota parser this replaces.
	var remoteSchemes []*conversationSchemeResponse
	if err := json.Unmarshal(remote.Conditions.ConversationSchemes, &remoteSchemes); err != nil {
		return fmt.Errorf("%w: conversationSchemes: %w", errInvalidResponse, err)
	}
	schemes := make([]attr.Value, 0, len(remoteSchemes))
	for _, scheme := range remoteSchemes {
		if scheme == nil {
			return fmt.Errorf("%w: null conversation scheme", errInvalidResponse)
		}
		if scheme.ODataType == nil {
			return fmt.Errorf("%w: scheme missing @odata.type", errInvalidResponse)
		}
		fields := map[string]attr.Value{"type": types.StringNull(), "url": types.StringNull(), "json_path": types.StringNull(), "scheme_name": types.StringNull()}
		switch strings.TrimPrefix(*scheme.ODataType, "#") {
		case "microsoft.graph.networkaccess.customConversationScheme":
			if scheme.URL == nil || *scheme.URL == "" || scheme.JSONPath == nil {
				return fmt.Errorf("%w: incomplete custom scheme", errInvalidResponse)
			}
			fields["type"] = types.StringValue(schemeTypeCustom)
			fields["url"] = types.StringValue(*scheme.URL)
			fields["json_path"] = types.StringValue(*scheme.JSONPath)
		case "microsoft.graph.networkaccess.predefinedConversationScheme":
			if scheme.SchemeName == nil {
				return fmt.Errorf("%w: predefined scheme missing name", errInvalidResponse)
			}
			fields["type"] = types.StringValue(schemeTypePredefined)
			fields["scheme_name"] = types.StringValue(*scheme.SchemeName)
		default:
			return fmt.Errorf("%w: unsupported scheme type %q", errInvalidResponse, *scheme.ODataType)
		}
		schemes = append(schemes, types.ObjectValueMust(conversationSchemeObjectType.AttrTypes, fields))
	}
	data.ID = convert.GraphToFrameworkString(remote.ID)
	data.Name = convert.GraphToFrameworkString(remote.Name)
	data.Description = convert.GraphToFrameworkString(remote.Description)
	data.Action = convert.GraphToFrameworkString(remote.Action)
	data.Priority = types.Int32Value(int32(*remote.Priority))
	data.Enabled = types.BoolValue(status == "enabled")
	data.Status = types.StringValue(status)
	data.PromptLogging = types.StringValue(*remote.PromptLogging)
	data.ScanResult = types.StringValue(*remote.Conditions.ScanResult)
	data.ConversationSchemes = types.ListValueMust(conversationSchemeObjectType, schemes)
	return nil
}
