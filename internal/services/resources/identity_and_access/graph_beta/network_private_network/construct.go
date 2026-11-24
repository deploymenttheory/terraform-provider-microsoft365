package graphBetaNetworkPrivateNetwork

import (
	"context"
	"errors"
	"fmt"
	"net/netip"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const (
	privateNetworkODataType              = "#microsoft.graph.networkaccess.privateNetwork"
	dnsResolutionIdentificationODataType = "#microsoft.graph.networkaccess.dnsResolutionIdentification"
	ipAddressODataType                   = "#microsoft.graph.networkaccess.ipAddress"
	ipSubnetODataType                    = "#microsoft.graph.networkaccess.ipSubnet"
	ipRangeODataType                     = "#microsoft.graph.networkaccess.ipRange"
	fqdnODataType                        = "#microsoft.graph.networkaccess.fqdn"

	expectedIPResolutionTypeIPAddress = "ip_address"
	expectedIPResolutionTypeIPSubnet  = "ip_subnet"
	expectedIPResolutionTypeIPRange   = "ip_range"
)

// errInvalidRequest marks a configuration the private network request cannot be built from.
var errInvalidRequest = errors.New("invalid private network request")

// constructResource builds the portal-observed private network payload for
// /networkaccess/privateNetworks. includeODataType adds the @odata.type the
// portal sends on create.
func constructResource(ctx context.Context, data *NetworkPrivateNetworkResourceModel, includeODataType bool) (privateNetworkRequest, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	var requestBody privateNetworkRequest

	if data.DNSResolutionIdentification == nil {
		return requestBody, fmt.Errorf("%w: dns_resolution_identification must be specified", errInvalidRequest)
	}
	dnsResolution := data.DNSResolutionIdentification

	// A null or unknown set has no elements, so it decodes to no values.
	var dnsServers []string
	if len(dnsResolution.DNSServers.Elements()) > 0 {
		if diags := dnsResolution.DNSServers.ElementsAs(ctx, &dnsServers, false); diags.HasError() {
			return requestBody, fmt.Errorf("%w: failed to read dns_resolution_identification.dns_servers: %s", errInvalidRequest, diags.Errors()[0].Detail())
		}
	}
	if len(dnsServers) == 0 {
		return requestBody, fmt.Errorf("%w: at least one dns_resolution_identification.dns_servers value must be specified", errInvalidRequest)
	}
	serverAddresses := make([]ipResolutionRequest, 0, len(dnsServers))
	for _, server := range dnsServers {
		if _, err := netip.ParseAddr(server); err != nil {
			return requestBody, fmt.Errorf("dns server %q must be a valid IP address: %w", server, err)
		}
		serverAddresses = append(serverAddresses, ipResolutionRequest{ODataType: ipAddressODataType, Value: &server})
	}

	var expectedIPResolutionModels []ExpectedIPResolutionModel
	if len(dnsResolution.ExpectedIPResolutions.Elements()) > 0 {
		if diags := dnsResolution.ExpectedIPResolutions.ElementsAs(ctx, &expectedIPResolutionModels, false); diags.HasError() {
			return requestBody, fmt.Errorf("%w: failed to read dns_resolution_identification.expected_ip_resolutions: %s", errInvalidRequest, diags.Errors()[0].Detail())
		}
	}
	if len(expectedIPResolutionModels) == 0 {
		return requestBody, fmt.Errorf("%w: at least one dns_resolution_identification.expected_ip_resolutions value must be specified", errInvalidRequest)
	}
	// ValueString is "" for a null or unknown value, so an empty string covers every unset field.
	expectedIPResolutions := make([]ipResolutionRequest, 0, len(expectedIPResolutionModels))
	for _, item := range expectedIPResolutionModels {
		resolutionType := item.Type.ValueString()
		switch resolutionType {
		case expectedIPResolutionTypeIPAddress, expectedIPResolutionTypeIPSubnet:
			value := item.Value.ValueString()
			if value == "" {
				return requestBody, fmt.Errorf("%w: expected_ip_resolutions value is required when type is %q", errInvalidRequest, resolutionType)
			}
			if resolutionType == expectedIPResolutionTypeIPSubnet {
				if _, err := netip.ParsePrefix(value); err != nil {
					return requestBody, fmt.Errorf("expected_ip_resolutions value %q must be a valid CIDR prefix: %w", value, err)
				}
				expectedIPResolutions = append(expectedIPResolutions, ipResolutionRequest{ODataType: ipSubnetODataType, Value: &value})
				continue
			}
			if _, err := netip.ParseAddr(value); err != nil {
				return requestBody, fmt.Errorf("expected_ip_resolutions value %q must be a valid IP address: %w", value, err)
			}
			expectedIPResolutions = append(expectedIPResolutions, ipResolutionRequest{ODataType: ipAddressODataType, Value: &value})
		case expectedIPResolutionTypeIPRange:
			begin, end := item.BeginAddress.ValueString(), item.EndAddress.ValueString()
			for _, address := range []struct{ name, value string }{{"begin_address", begin}, {"end_address", end}} {
				if address.value == "" {
					return requestBody, fmt.Errorf("%w: expected_ip_resolutions %s is required when type is %q", errInvalidRequest, address.name, resolutionType)
				}
				if _, err := netip.ParseAddr(address.value); err != nil {
					return requestBody, fmt.Errorf("expected_ip_resolutions %s %q must be a valid IP address: %w", address.name, address.value, err)
				}
			}
			expectedIPResolutions = append(expectedIPResolutions, ipResolutionRequest{ODataType: ipRangeODataType, BeginAddress: &begin, EndAddress: &end})
		default:
			return requestBody, fmt.Errorf("%w: expected_ip_resolutions type must be one of \"ip_address\", \"ip_subnet\", or \"ip_range\"", errInvalidRequest)
		}
	}

	if !data.AppIDs.IsNull() && !data.AppIDs.IsUnknown() {
		appIDs := []string{}
		if diags := data.AppIDs.ElementsAs(ctx, &appIDs, false); diags.HasError() {
			return requestBody, fmt.Errorf("%w: failed to read app_ids: %s", errInvalidRequest, diags.Errors()[0].Detail())
		}
		requestBody.AppIDs = &appIDs
	}

	if includeODataType {
		requestBody.ODataType = privateNetworkODataType
	}
	requestBody.Name = data.Name.ValueStringPointer()
	requestBody.NetworkIdentifications = []dnsResolutionIdentificationRequest{
		{
			ODataType:             dnsResolutionIdentificationODataType,
			ServerAddresses:       serverAddresses,
			FQDNToResolve:         &ipResolutionRequest{ODataType: fqdnODataType, Value: dnsResolution.FQDNToResolve.ValueStringPointer()},
			ExpectedIPResolutions: expectedIPResolutions,
		},
	}

	tflog.Debug(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), map[string]any{"request_body": requestBody})

	return requestBody, nil
}
