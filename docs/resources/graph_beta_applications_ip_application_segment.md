---
page_title: "microsoft365_graph_beta_applications_ip_application_segment Resource - terraform-provider-microsoft365"
subcategory: "Applications"

description: |-
  Manages an IP application segment for on-premises publishing. IP application segments define the destination hosts, ports, and protocols for applications published through Azure AD Application Proxy.
---

# microsoft365_graph_beta_applications_ip_application_segment (Resource)

Manages an IP application segment for on-premises publishing. IP application segments define the destination hosts, ports, and protocols for applications published through Azure AD Application Proxy.

## Microsoft Documentation

- [ipApplicationSegment resource type](https://learn.microsoft.com/en-us/graph/api/resources/ipapplicationsegment?view=graph-rest-beta)
- [Create ipApplicationSegment](https://learn.microsoft.com/en-us/graph/api/onpremisespublishingprofile-post-applicationsegments?view=graph-rest-beta&tabs=http)
- [Update ipApplicationSegment](https://learn.microsoft.com/en-us/graph/api/ipapplicationsegment-update?view=graph-rest-beta&tabs=http)
- [Delete ipApplicationSegment](https://learn.microsoft.com/en-us/graph/api/ipapplicationsegment-delete?view=graph-rest-beta&tabs=http)

## API Permissions

The following API permissions are required in order to use this resource.

### Microsoft Graph

- **Application**: `Application.Read.All` and `Application.ReadWrite.All`

## Version History

| Version | Status | Notes |
|---------|--------|-------|
| v0.33.0 | Experimental | Initial release |
| v0.41.0 | Experimental | Renamed `application_id` to `application_object_id` and added more examples|

## Example Usage

### Basic Configuration - IP Address

This example demonstrates a minimal configuration targeting a single IP address.

```terraform
# Basic IP Application Segment with single IP address
# This example demonstrates the minimal configuration for an IP application segment
# targeting a single IP address.

resource "microsoft365_graph_beta_applications_ip_application_segment" "minimal_ip" {
  application_object_id = "00000000-0000-0000-0000-000000000000"
  destination_host      = "192.168.1.100"
  destination_type      = "ipAddress"
  ports                 = ["80-80"]
  protocol              = "tcp"
}
```

### IP Range (CIDR Notation)

This example shows how to configure an application segment for an entire IP subnet using CIDR notation.

```terraform
# IP Application Segment with IP Range (CIDR notation)
# This example demonstrates how to configure an application segment for an entire
# IP subnet using CIDR notation.

resource "microsoft365_graph_beta_applications_ip_application_segment" "ip_range" {
  application_object_id = "00000000-0000-0000-0000-000000000000"
  destination_host      = "192.168.1.0/24"
  destination_type      = "ipRangeCidr"
  ports                 = ["443-443"]
  protocol              = "tcp"

  timeouts = {
    create = "5m"
    read   = "5m"
    update = "5m"
    delete = "5m"
  }
}
```

### Fully Qualified Domain Name (FQDN)

This example demonstrates configuration using a specific hostname with multiple ports.

```terraform
# IP Application Segment with Fully Qualified Domain Name (FQDN)
# This example demonstrates how to configure an application segment using a specific
# hostname with multiple ports.

resource "microsoft365_graph_beta_applications_ip_application_segment" "fqdn" {
  application_object_id = "00000000-0000-0000-0000-000000000000"
  destination_host      = "app.contoso.com"
  destination_type      = "fqdn"
  ports                 = ["443-443", "8443-8443"]
  protocol              = "tcp"

  timeouts = {
    create = "5m"
    read   = "5m"
    update = "5m"
    delete = "5m"
  }
}
```

### DNS Suffix (Wildcard Domain)

This example shows how to use a wildcard domain to match all subdomains.

```terraform
# IP Application Segment with DNS Suffix (Wildcard Domain)
# This example demonstrates how to configure an application segment using a wildcard
# domain to match all subdomains.

resource "microsoft365_graph_beta_applications_ip_application_segment" "dns_suffix" {
  application_object_id = "00000000-0000-0000-0000-000000000000"
  destination_host      = "*.internal.contoso.com"
  destination_type      = "dnsSuffix"
  ports = [
    "80-80",
    "443-443",
    "8080-8080",
    "8443-8443"
  ]
  protocol = "tcp"

  timeouts = {
    create = "5m"
    read   = "5m"
    update = "5m"
    delete = "5m"
  }
}
```

### UDP Protocol

This example demonstrates configuration using UDP protocol, useful for applications like VoIP or video conferencing.

```terraform
# IP Application Segment with UDP Protocol
# This example demonstrates how to configure an application segment using UDP protocol
# instead of TCP, useful for applications like VoIP or video conferencing.

resource "microsoft365_graph_beta_applications_ip_application_segment" "udp_app" {
  application_object_id = "00000000-0000-0000-0000-000000000000"
  destination_host      = "voip.contoso.com"
  destination_type      = "fqdn"
  ports                 = ["5060-5061", "10000-20000"]
  protocol              = "udp"

  timeouts = {
    create = "5m"
    read   = "5m"
    update = "5m"
    delete = "5m"
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `application_object_id` (String) The unique object identifier of the application.
- `destination_host` (String) Either the IP address, IP range, or FQDN of the application segment, with or without wildcards.
- `destination_type` (String) The type of destination for the application segment.The possible values are: `ipAddress`, `ipRange`, `ipRangeCidr`, `fqdn`, `dnsSuffix`, `unknownFutureValue`.
- `ports` (Set of String) List of ports supported for the application segment.
- `protocol` (String) Indicates the protocol of the network traffic acquired for the application segment.The possible values are: `tcp`, `udp`, `unknownFutureValue`.

### Optional

- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

### Read-Only

- `id` (String) The unique identifier of the application segment.

<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).
- `delete` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Setting a timeout for a Delete operation is only applicable if changes are saved into state before the destroy operation occurs.
- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Read operations occur during any refresh or planning operation when refresh is enabled.
- `update` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).

## Import

Import is supported using the following syntax:

```shell
#!/bin/bash

# Import an existing IP application segment by its ID
terraform import microsoft365_graph_beta_applications_ip_application_segment.example_ip_address "00000000-0000-0000-0000-000000000000"

# The ID format is the segment's unique identifier (GUID)
# You can find the segment ID in the Azure Portal or via Microsoft Graph API:
# GET https://graph.microsoft.com/beta/applications/{application-id}/onPremisesPublishing/segmentsConfiguration/microsoft.graph.ipSegmentConfiguration/applicationSegments
```