package graphBetaServicePrincipalTokenSigningCertificate

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"regexp"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

var thumbprintRegex = regexp.MustCompile(`^[0-9a-fA-F]{40}$`)

// MapSelfSignedCertificateToTerraform maps the addTokenSigningCertificate response to the Terraform state
func MapSelfSignedCertificateToTerraform(ctx context.Context, data *ServicePrincipalTokenSigningCertificateResourceModel, certificate graphmodels.SelfSignedCertificateable) {
	if certificate == nil {
		tflog.Warn(ctx, "Received nil self-signed certificate in MapSelfSignedCertificateToTerraform")
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Starting to map self-signed certificate to Terraform state for %s", ResourceName))

	data.DisplayName = convert.GraphToFrameworkString(certificate.GetDisplayName())
	data.EndDateTime = mapTimeKeepingConfiguredValue(data.EndDateTime, convert.GraphToFrameworkTime(certificate.GetEndDateTime()))
	data.StartDateTime = convert.GraphToFrameworkTime(certificate.GetStartDateTime())
	data.Thumbprint = convert.GraphToFrameworkString(certificate.GetThumbprint())

	if keyId := certificate.GetKeyId(); keyId != nil {
		data.KeyId = types.StringValue(keyId.String())
	} else {
		tflog.Warn(ctx, "addTokenSigningCertificate response did not include a keyId")
	}
	if key := certificate.GetKey(); len(key) > 0 {
		data.Value = types.StringValue(base64.StdEncoding.EncodeToString(key))
	}

	if !data.KeyId.IsNull() && !data.KeyId.IsUnknown() {
		data.Id = types.StringValue(fmt.Sprintf("%s/%s", data.ServicePrincipalID.ValueString(), data.KeyId.ValueString()))
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished mapping self-signed certificate for %s", ResourceName))
}

// MapRemoteResourceStateToTerraform locates the signing key credential on the service principal
// by key_id and maps it to the Terraform state. Returns false when the credential no longer exists.
func MapRemoteResourceStateToTerraform(ctx context.Context, data *ServicePrincipalTokenSigningCertificateResourceModel, servicePrincipal graphmodels.ServicePrincipalable) bool {
	if servicePrincipal == nil {
		tflog.Warn(ctx, "Received nil service principal in MapRemoteResourceStateToTerraform")
		return false
	}

	tflog.Debug(ctx, fmt.Sprintf("Starting to map remote state to Terraform state for %s", ResourceName))

	keyId := data.KeyId.ValueString()

	for _, credential := range servicePrincipal.GetKeyCredentials() {
		credentialKeyId := credential.GetKeyId()
		if credentialKeyId == nil || credentialKeyId.String() != keyId {
			continue
		}

		data.DisplayName = convert.GraphToFrameworkString(credential.GetDisplayName())
		data.EndDateTime = mapTimeKeepingConfiguredValue(data.EndDateTime, convert.GraphToFrameworkTime(credential.GetEndDateTime()))
		data.StartDateTime = convert.GraphToFrameworkTime(credential.GetStartDateTime())

		// The GET response omits the certificate thumbprint; recover it from
		// customKeyIdentifier (raw SHA-1 bytes or the ASCII thumbprint string)
		// when it is not already known, e.g. after import.
		if data.Thumbprint.IsNull() || data.Thumbprint.ValueString() == "" {
			if thumbprint := thumbprintFromCustomKeyIdentifier(credential.GetCustomKeyIdentifier()); thumbprint != "" {
				data.Thumbprint = types.StringValue(thumbprint)
			}
		}

		// Key material may be omitted on reads; keep the value captured at creation.
		if key := credential.GetKey(); len(key) > 0 && (data.Value.IsNull() || data.Value.ValueString() == "") {
			data.Value = types.StringValue(base64.StdEncoding.EncodeToString(key))
		}

		data.Id = types.StringValue(fmt.Sprintf("%s/%s", data.ServicePrincipalID.ValueString(), keyId))

		tflog.Debug(ctx, fmt.Sprintf("Finished mapping %s with key_id: %s", ResourceName, keyId))
		return true
	}

	tflog.Debug(ctx, fmt.Sprintf("Key credential %s not found on service principal %s", keyId, data.ServicePrincipalID.ValueString()))
	return false
}

// mapTimeKeepingConfiguredValue keeps the configured timestamp when it represents the same
// instant as the remote value, avoiding perpetual diffs from formatting differences
// (e.g. 2028-01-01T14:59:59Z vs 2028-01-01T14:59:59+00:00).
func mapTimeKeepingConfiguredValue(configured types.String, remote types.String) types.String {
	if configured.IsNull() || configured.IsUnknown() || remote.IsNull() {
		return remote
	}
	configuredTime, configuredErr := time.Parse(time.RFC3339, configured.ValueString())
	remoteTime, remoteErr := time.Parse(time.RFC3339, remote.ValueString())
	if configuredErr == nil && remoteErr == nil && configuredTime.Equal(remoteTime) {
		return configured
	}
	return remote
}

// thumbprintFromCustomKeyIdentifier converts a credential customKeyIdentifier to a SHA-1
// thumbprint string. Graph stores either the raw 20 hash bytes or the ASCII thumbprint.
func thumbprintFromCustomKeyIdentifier(customKeyIdentifier []byte) string {
	if len(customKeyIdentifier) == 0 {
		return ""
	}
	if thumbprintRegex.Match(customKeyIdentifier) {
		return string(customKeyIdentifier)
	}
	if encoded := hex.EncodeToString(customKeyIdentifier); thumbprintRegex.MatchString(encoded) {
		return encoded
	}
	return ""
}
