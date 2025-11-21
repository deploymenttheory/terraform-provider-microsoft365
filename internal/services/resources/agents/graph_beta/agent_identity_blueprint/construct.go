package graphBetaAgentsAgentIdentityBlueprint

import (
	"context"
	"fmt"
)

// constructResource maps the Terraform schema to a JSON-serializable map for agent identity blueprint.
// Since we're using the HTTP client with OData type cast URLs, we work with JSON directly.
func constructResource(ctx context.Context, data *AgentIdentityBlueprintResourceModel) (map[string]any, error) {
	// Create a map for the request body
	// Note: agentIdentityBlueprint is a derived type of application
	requestBody := make(map[string]any)

	// Set the OData type to indicate this is an agentIdentityBlueprint
	// This is critical for the API to recognize the resource type
	requestBody["@odata.type"] = "#microsoft.graph.agentIdentityBlueprint"

	// Set display name (required)
	if !data.DisplayName.IsNull() && !data.DisplayName.IsUnknown() {
		requestBody["displayName"] = data.DisplayName.ValueString()
	} else {
		return nil, fmt.Errorf("display_name is required")
	}

	// Set optional string fields
	if !data.Description.IsNull() && !data.Description.IsUnknown() {
		requestBody["description"] = data.Description.ValueString()
	}

	if !data.SignInAudience.IsNull() && !data.SignInAudience.IsUnknown() {
		requestBody["signInAudience"] = data.SignInAudience.ValueString()
	}

	if !data.GroupMembershipClaims.IsNull() && !data.GroupMembershipClaims.IsUnknown() {
		requestBody["groupMembershipClaims"] = data.GroupMembershipClaims.ValueString()
	}

	if !data.TokenEncryptionKeyId.IsNull() && !data.TokenEncryptionKeyId.IsUnknown() {
		requestBody["tokenEncryptionKeyId"] = data.TokenEncryptionKeyId.ValueString()
	}

	if !data.ServiceManagementReference.IsNull() && !data.ServiceManagementReference.IsUnknown() {
		requestBody["serviceManagementReference"] = data.ServiceManagementReference.ValueString()
	}

	// Set collection fields
	if !data.IdentifierUris.IsNull() && !data.IdentifierUris.IsUnknown() {
		var identifierUris []string
		data.IdentifierUris.ElementsAs(ctx, &identifierUris, false)
		requestBody["identifierUris"] = identifierUris
	}

	if !data.Tags.IsNull() && !data.Tags.IsUnknown() {
		var tags []string
		data.Tags.ElementsAs(ctx, &tags, false)
		requestBody["tags"] = tags
	}

	return requestBody, nil
}
