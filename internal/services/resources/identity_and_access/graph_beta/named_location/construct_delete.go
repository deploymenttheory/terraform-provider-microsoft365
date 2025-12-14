package graphBetaNamedLocation

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructResourceForDeletion builds the PATCH request body for preparing a trusted IP named location for deletion.
// Microsoft Graph API requires trusted IP locations to be set to untrusted before deletion.
// The request must include displayName, isTrusted=false, and at least one IP address.
func constructResourceForDeletion(ctx context.Context) (graphmodels.NamedLocationable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s deletion patch body", ResourceName))

	ipLocation := graphmodels.NewIpNamedLocation()

	// Set display name
	displayName := "for_deletion"
	ipLocation.SetDisplayName(&displayName)

	// Set isTrusted to false
	isTrusted := false
	ipLocation.SetIsTrusted(&isTrusted)

	// Set a minimal IP range (required by API)
	ipv4Range := graphmodels.NewIPv4CidrRange()
	cidrAddress := "0.0.0.0/32"
	ipv4Range.SetCidrAddress(&cidrAddress)

	ipRanges := []graphmodels.IpRangeable{ipv4Range}
	ipLocation.SetIpRanges(ipRanges)

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s deletion patch body", ResourceName))

	return ipLocation, nil
}
