package graphBetaServicePrincipalTokenSigningCertificate

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/serviceprincipals"
)

// constructResource constructs the addTokenSigningCertificate POST request body
func constructResource(ctx context.Context, data *ServicePrincipalTokenSigningCertificateResourceModel) (*serviceprincipals.ItemAddTokenSigningCertificatePostRequestBody, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	requestBody := serviceprincipals.NewItemAddTokenSigningCertificatePostRequestBody()

	convert.FrameworkToGraphString(data.DisplayName, requestBody.SetDisplayName)

	if err := convert.FrameworkToGraphTime(data.EndDateTime, requestBody.SetEndDateTime); err != nil {
		return nil, fmt.Errorf("failed to set end_date_time: %w", err)
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return requestBody, nil
}

// constructDeleteBody builds the raw ServicePrincipal PATCH body that removes the
// certificate's key credentials (Sign and Verify) and the associated password credential.
// addTokenSigningCertificate stamps all three credentials with the certificate thumbprint
// as customKeyIdentifier and gives the password credential the same keyId as the Sign key
// credential, so the certificate's credentials are identified primarily by the exact
// keyId/customKeyIdentifier from the Sign credential (with the thumbprint heuristic as a
// fallback for imported state where the Sign credential may no longer resolve).
//
// The body is assembled from the raw GET JSON rather than the typed SDK models: Graph
// rejects any perceived modification of a retained credential with "Update to existing
// credential with KeyId ... is not allowed", and the SDK's serialization loses the
// fractional-second precision of the original timestamps (and requires the key material
// to be sent back as an explicit null), both of which read as modifications. Passing the
// original JSON through untouched — with only the key material nulled out — matches what
// the API requires.
func constructDeleteBody(ctx context.Context, rawServicePrincipal map[string]any, keyId string, thumbprint string) ([]byte, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing delete request for %s", ResourceName))

	asCredentialMaps := func(rawCredentials any) []map[string]any {
		credentials, _ := rawCredentials.([]any)
		result := make([]map[string]any, 0, len(credentials))
		for _, rawCredential := range credentials {
			if credential, ok := rawCredential.(map[string]any); ok {
				result = append(result, credential)
			}
		}
		return result
	}

	keyCredentials := asCredentialMaps(rawServicePrincipal["keyCredentials"])
	passwordCredentials := asCredentialMaps(rawServicePrincipal["passwordCredentials"])

	// Resolve the certificate's exact customKeyIdentifier from the Sign credential's keyId,
	// which links the Verify credential and the password credential to the same certificate.
	targetCustomKeyIdentifier := ""
	for _, credential := range keyCredentials {
		if credentialKeyId, _ := credential["keyId"].(string); credentialKeyId != "" && credentialKeyId == keyId {
			targetCustomKeyIdentifier, _ = credential["customKeyIdentifier"].(string)
			break
		}
	}

	belongsToCertificate := func(credential map[string]any) bool {
		if credentialKeyId, _ := credential["keyId"].(string); credentialKeyId != "" && credentialKeyId == keyId {
			return true
		}
		customKeyIdentifier, _ := credential["customKeyIdentifier"].(string)
		if targetCustomKeyIdentifier != "" && customKeyIdentifier == targetCustomKeyIdentifier {
			return true
		}
		return base64CustomKeyIdentifierMatchesThumbprint(customKeyIdentifier, thumbprint)
	}

	filterCredentials := func(credentials []map[string]any, nullifyFields []string) []any {
		retained := make([]any, 0)
		for _, credential := range credentials {
			if belongsToCertificate(credential) {
				continue
			}
			// Graph requires existing credentials to be "sent back with null value"
			// for their secret material.
			for _, field := range nullifyFields {
				credential[field] = nil
			}
			retained = append(retained, credential)
		}
		return retained
	}

	body := map[string]any{
		"keyCredentials":      filterCredentials(keyCredentials, []string{"key"}),
		"passwordCredentials": filterCredentials(passwordCredentials, []string{"secretText"}),
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing delete request for %s", ResourceName))

	return json.Marshal(body)
}

// base64CustomKeyIdentifierMatchesThumbprint decodes a raw JSON customKeyIdentifier
// (base64) and reports whether it refers to the given SHA-1 thumbprint.
func base64CustomKeyIdentifierMatchesThumbprint(customKeyIdentifier string, thumbprint string) bool {
	if customKeyIdentifier == "" {
		return false
	}
	decoded, err := base64.StdEncoding.DecodeString(customKeyIdentifier)
	if err != nil {
		return false
	}
	return customKeyIdentifierMatchesThumbprint(decoded, thumbprint)
}

// customKeyIdentifierMatchesThumbprint reports whether a credential's customKeyIdentifier
// refers to the given SHA-1 thumbprint. Graph stores the identifier either as the raw hash
// bytes (20 bytes, hex-encodes to the thumbprint) or as the ASCII thumbprint string.
func customKeyIdentifierMatchesThumbprint(customKeyIdentifier []byte, thumbprint string) bool {
	if len(customKeyIdentifier) == 0 {
		return false
	}
	return strings.EqualFold(hex.EncodeToString(customKeyIdentifier), thumbprint) ||
		strings.EqualFold(string(customKeyIdentifier), thumbprint)
}
