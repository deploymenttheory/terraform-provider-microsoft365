package graphBetaApplicationsOnPremisesIpApplicationSegment

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
)

// constructResource builds the ipApplicationSegment create and update body.
func constructResource(ctx context.Context, data *OnPremisesIpApplicationSegmentResourceModel) (ipApplicationSegmentRequest, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	var requestBody ipApplicationSegmentRequest

	if !data.DestinationHost.IsNull() && !data.DestinationHost.IsUnknown() {
		requestBody.DestinationHost = data.DestinationHost.ValueString()
	}

	// Microsoft Learn documents ipApplicationSegment.destinationType values as
	// "ipAddress", "ipRange", "ipRangeCidr", "fqdn", and "dnsSuffix":
	// https://learn.microsoft.com/en-us/graph/api/resources/ipapplicationsegment?view=graph-rest-beta
	// The create API page also lists "ipAddress" for POST:
	// https://learn.microsoft.com/en-us/graph/api/onpremisespublishingprofile-post-applicationsegments?view=graph-rest-beta
	//
	// The beta endpoint currently behaves differently for a single IP segment:
	// POSTing destinationType "ipAddress" returns 400 InvalidJson_BadRequest
	// ("Valid JSON content expected."), while destinationType "ip" succeeds.
	// Direct API checks (2026-07) showed:
	//   - The endpoint infers the stored destination type from the host format
	//     regardless of the requested destinationType. ipRange stays ipRange
	//     only for IPv4 "start..end" hosts (e.g. "192.168.1.1..192.168.1.10");
	//     CIDR, single-IP, and FQDN hosts are stored as ipRangeCidr, ip, and
	//     fqdn respectively, hyphenated ranges return 400
	//     DestinationHost_InvalidIP, reversed ranges return 500, and IPv6
	//     ranges return 400 DestinationHost_Invalid. ValidateConfig enforces
	//     the "start..end" format at plan time to prevent drift.
	//   - dnsSuffix is accepted, but Graph discards the segment's ports
	//     (returned as []) and protocol (returned as "0"), so it cannot
	//     round-trip through this resource's required ports/protocol schema.
	// The schema only allows values observed to create and read back through
	// this application-scoped endpoint. Wildcard hosts such as
	// "*.internal.example.com" are supported when sent as destinationType fqdn.
	// Keep Terraform's public schema aligned with Learn, but send the literal
	// accepted by Graph.
	if !data.DestinationType.IsNull() && !data.DestinationType.IsUnknown() {
		requestBody.DestinationType = data.DestinationType.ValueString()
		if requestBody.DestinationType == "ipAddress" {
			requestBody.DestinationType = "ip"
		}
	}

	if err := convert.FrameworkToGraphStringSet(ctx, data.Ports, func(ports []string) {
		requestBody.Ports = ports
	}); err != nil {
		return requestBody, fmt.Errorf("failed to set ports: %w", err)
	}

	// Graph takes the protocols as one comma-separated string, e.g. "tcp,udp".
	if err := convert.FrameworkToGraphStringSet(ctx, data.Protocol, func(protocols []string) {
		sort.Strings(protocols)
		requestBody.Protocol = strings.Join(protocols, ",")
	}); err != nil {
		return requestBody, fmt.Errorf("failed to set protocol: %w", err)
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName), map[string]any{"request_body": requestBody})

	return requestBody, nil
}
