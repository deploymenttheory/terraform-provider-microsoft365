package graphBetaDeviceConfigurationTemplatesJson

import (
	"context"
	"encoding/json"
	"fmt"
	"maps"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/types"
	abstractions "github.com/microsoft/kiota-abstractions-go"

	customrequests "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/custom_requests"
	kiotaerrors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
)

// These navigation properties are omitted from profile GET responses. The mapping
// includes inherited certificate relationships from the Graph beta metadata, except
// imported Android owner/work PFX profiles whose rootCertificate routes Graph rejects.
// It does not restrict the profile types accepted by the resource.
var deviceConfigurationRelationships = map[string]map[string]bool{
	"#microsoft.graph.androidDeviceOwnerEnterpriseWiFiConfiguration": {
		"identityCertificateForClientAuthentication": false,
		"rootCertificateForServerValidation":         false,
		"rootCertificatesForServerValidation":        true,
	},
	"#microsoft.graph.androidDeviceOwnerPkcsCertificateProfile": {
		"rootCertificate": false,
	},
	"#microsoft.graph.androidDeviceOwnerScepCertificateProfile": {
		"rootCertificate": false,
	},
	"#microsoft.graph.androidDeviceOwnerVpnConfiguration": {
		"identityCertificate": false,
	},
	"#microsoft.graph.androidEasEmailProfileConfiguration": {
		"identityCertificate":     false,
		"smimeSigningCertificate": false,
	},
	"#microsoft.graph.androidEnterpriseWiFiConfiguration": {
		"identityCertificateForClientAuthentication": false,
		"rootCertificateForServerValidation":         false,
	},
	"#microsoft.graph.androidForWorkEnterpriseWiFiConfiguration": {
		"identityCertificateForClientAuthentication": false,
		"rootCertificateForServerValidation":         false,
	},
	"#microsoft.graph.androidForWorkGmailEasConfiguration": {
		"identityCertificate": false,
	},
	"#microsoft.graph.androidForWorkNineWorkEasConfiguration": {
		"identityCertificate": false,
	},
	"#microsoft.graph.androidForWorkPkcsCertificateProfile": {
		"rootCertificate": false,
	},
	"#microsoft.graph.androidForWorkScepCertificateProfile": {
		"rootCertificate": false,
	},
	"#microsoft.graph.androidForWorkVpnConfiguration": {
		"identityCertificate": false,
	},
	"#microsoft.graph.androidImportedPFXCertificateProfile": {
		"rootCertificate": false,
	},
	"#microsoft.graph.androidPkcsCertificateProfile": {
		"rootCertificate": false,
	},
	"#microsoft.graph.androidScepCertificateProfile": {
		"rootCertificate": false,
	},
	"#microsoft.graph.androidVpnConfiguration": {
		"identityCertificate": false,
	},
	"#microsoft.graph.androidWorkProfileEnterpriseWiFiConfiguration": {
		"identityCertificateForClientAuthentication": false,
		"rootCertificateForServerValidation":         false,
	},
	"#microsoft.graph.androidWorkProfileGmailEasConfiguration": {
		"identityCertificate": false,
	},
	"#microsoft.graph.androidWorkProfileNineWorkEasConfiguration": {
		"identityCertificate": false,
	},
	"#microsoft.graph.androidWorkProfilePkcsCertificateProfile": {
		"rootCertificate": false,
	},
	"#microsoft.graph.androidWorkProfileScepCertificateProfile": {
		"rootCertificate": false,
	},
	"#microsoft.graph.androidWorkProfileVpnConfiguration": {
		"identityCertificate": false,
	},
	"#microsoft.graph.aospDeviceOwnerEnterpriseWiFiConfiguration": {
		"identityCertificateForClientAuthentication": false,
		"rootCertificateForServerValidation":         false,
	},
	"#microsoft.graph.aospDeviceOwnerPkcsCertificateProfile": {
		"rootCertificate": false,
	},
	"#microsoft.graph.aospDeviceOwnerScepCertificateProfile": {
		"rootCertificate": false,
	},
	"#microsoft.graph.iosDeviceFeaturesConfiguration": {
		"identityCertificateForClientAuthentication": false,
		"singleSignOnExtensionPkinitCertificate":     false,
	},
	"#microsoft.graph.iosEasEmailProfileConfiguration": {
		"identityCertificate":        false,
		"smimeEncryptionCertificate": false,
		"smimeSigningCertificate":    false,
	},
	"#microsoft.graph.iosEnterpriseWiFiConfiguration": {
		"identityCertificateForClientAuthentication": false,
		"rootCertificatesForServerValidation":        true,
	},
	"#microsoft.graph.iosScepCertificateProfile": {
		"rootCertificate": false,
	},
	"#microsoft.graph.iosVpnConfiguration": {
		"identityCertificate": false,
	},
	"#microsoft.graph.iosWiredNetworkConfiguration": {
		"identityCertificateForClientAuthentication": false,
		"rootCertificateForServerValidation":         false,
	},
	"#microsoft.graph.iosikEv2VpnConfiguration": {
		"identityCertificate": false,
	},
	"#microsoft.graph.macOSDeviceFeaturesConfiguration": {
		"singleSignOnExtensionPkinitCertificate": false,
	},
	"#microsoft.graph.macOSEnterpriseWiFiConfiguration": {
		"identityCertificateForClientAuthentication": false,
		"rootCertificateForServerValidation":         false,
		"rootCertificatesForServerValidation":        true,
	},
	"#microsoft.graph.macOSScepCertificateProfile": {
		"rootCertificate": false,
	},
	"#microsoft.graph.macOSVpnConfiguration": {
		"identityCertificate": false,
	},
	"#microsoft.graph.macOSWiredNetworkConfiguration": {
		"identityCertificateForClientAuthentication": false,
		"rootCertificateForServerValidation":         false,
	},
	"#microsoft.graph.windows10VpnConfiguration": {
		"identityCertificate": false,
	},
	"#microsoft.graph.windows81SCEPCertificateProfile": {
		"rootCertificate": false,
	},
	"#microsoft.graph.windowsPhone81SCEPCertificateProfile": {
		"rootCertificate": false,
	},
	"#microsoft.graph.windowsPhone81VpnConfiguration": {
		"identityCertificate": false,
	},
	"#microsoft.graph.windowsWifiEnterpriseEAPConfiguration": {
		"identityCertificateForClientAuthentication": false,
		"rootCertificateForClientValidation":         false,
		"rootCertificatesForServerValidation":        true,
	},
	"#microsoft.graph.windowsWiredNetworkConfiguration": {
		"identityCertificateForClientAuthentication":          false,
		"rootCertificateForClientValidation":                  false,
		"rootCertificatesForServerValidation":                 true,
		"secondaryIdentityCertificateForClientAuthentication": false,
		"secondaryRootCertificateForClientValidation":         false,
	},
	"#microsoft.graph.windowsZtdnsConfiguration": {
		"rootCertificatesForClientValidation": true,
		"rootCertificatesForServerValidation": true,
	},
}

func (r *DeviceConfigurationTemplatesJsonResource) readRelationships(
	ctx context.Context,
	id string,
	prior types.String,
	settings map[string]any,
) error {
	var configured map[string]any
	if !prior.IsNull() && !prior.IsUnknown() {
		if err := json.Unmarshal([]byte(prior.ValueString()), &configured); err != nil {
			return fmt.Errorf("decode configured relationships: %w", err)
		}
	}
	typ, _ := settings["@odata.type"].(string)
	for name, collection := range deviceConfigurationRelationships[typ] {
		var response map[string]any
		err := customrequests.JSONRequest(ctx, r.client.GetAdapter(), abstractions.GET,
			r.ResourcePath+"/{id}/"+strings.TrimPrefix(typ, "#")+"/"+name,
			map[string]string{"id": id}, nil, &response)
		if err != nil {
			if kiotaerrors.GraphError(ctx, err).StatusCode == 404 {
				continue
			}
			return fmt.Errorf("read %s relationship: %w", name, err)
		}
		if collection {
			bindings, err := r.readRelationshipCollection(ctx, name, response)
			if err != nil {
				return err
			}
			if len(bindings) > 0 {
				sort.Strings(bindings)
				settings[name+"@odata.bind"] = bindings
			}
		} else {
			binding, err := r.relationshipBinding(response)
			if err != nil {
				return err
			}
			settings[name+"@odata.bind"] = binding
		}
		key := name + "@odata.bind"
		if existing, exists := configured[key]; exists {
			before, err := relationshipIDs(existing, collection)
			if err != nil {
				return err
			}
			after, err := relationshipIDs(settings[key], collection)
			if err != nil {
				return err
			}
			// Navigation collections have no meaningful order. Preserve the user's
			// URL spelling and order only when all actual related IDs still match.
			if maps.EqualFunc(before, after, func(_, _ string) bool { return true }) {
				settings[key] = existing
			}
		}
	}
	return nil
}

func (r *DeviceConfigurationTemplatesJsonResource) relationshipBinding(
	related map[string]any,
) (string, error) {
	id, ok := related["id"].(string)
	if !ok || id == "" {
		return "", fmt.Errorf("%w: relationship response has no id", errInvalidProfileResponse)
	}
	return r.client.GetAdapter().GetBaseUrl() + r.ResourcePath + "('" + id + "')", nil
}

func (r *DeviceConfigurationTemplatesJsonResource) readRelationshipCollection(
	ctx context.Context,
	name string,
	response map[string]any,
) ([]string, error) {
	bindings := []string{}
	seenPages := make(map[string]bool)
	for {
		values, ok := response["value"].([]any)
		if !ok {
			return nil, fmt.Errorf(
				"%w: %s must return a collection",
				errInvalidProfileResponse,
				name,
			)
		}
		for _, value := range values {
			related, ok := value.(map[string]any)
			if !ok {
				return nil, fmt.Errorf(
					"%w: %s contains an invalid reference",
					errInvalidProfileResponse,
					name,
				)
			}
			binding, err := r.relationshipBinding(related)
			if err != nil {
				return nil, err
			}
			bindings = append(bindings, binding)
		}
		next, _ := response["@odata.nextLink"].(string)
		if next == "" {
			return bindings, nil
		}
		base := r.client.GetAdapter().GetBaseUrl()
		if !strings.HasPrefix(next, base+"/") || seenPages[next] {
			return nil, fmt.Errorf("%w: invalid next page for %s", errInvalidProfileResponse, name)
		}
		seenPages[next] = true
		response = nil
		if err := customrequests.JSONRequest(
			ctx,
			r.client.GetAdapter(),
			abstractions.GET,
			strings.TrimPrefix(next, base),
			nil,
			nil,
			&response,
		); err != nil {
			return nil, fmt.Errorf("read next %s relationship page: %w", name, err)
		}
	}
}
