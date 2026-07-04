---
page_title: "microsoft365_graph_beta_applications_on_premises_ip_application_segment Resource - terraform-provider-microsoft365"
subcategory: "Applications"

description: |-
  Manages an IP application segment for on-premises publishing. IP application segments define the destination hosts, ports, and protocols for applications published through Azure AD Application Proxy.
---

# microsoft365_graph_beta_applications_on_premises_ip_application_segment (Resource)

Manages an IP application segment for on-premises publishing. IP application segments define the destination hosts, ports, and protocols for applications published through Azure AD Application Proxy.

## Microsoft Documentation

- [ipApplicationSegment resource type](https://learn.microsoft.com/en-us/graph/api/resources/ipapplicationsegment?view=graph-rest-beta)
- [Create ipApplicationSegment](https://learn.microsoft.com/en-us/graph/api/onpremisespublishingprofile-post-applicationsegments?view=graph-rest-beta&tabs=http)
- [Read ipApplicationSegment](https://learn.microsoft.com/en-us/graph/api/ipapplicationsegment-get?view=graph-rest-beta&tabs=http)
- [Update ipApplicationSegment](https://learn.microsoft.com/en-us/graph/api/ipapplicationsegment-update?view=graph-rest-beta&tabs=http)
- [Delete ipApplicationSegment](https://learn.microsoft.com/en-us/graph/api/ipapplicationsegment-delete?view=graph-rest-beta&tabs=http)

## Microsoft Graph API Permissions

The following client `application` permissions are needed in order to use this resource:

**Required:**
- `Application.Read.All`
- `Application.ReadWrite.All`

**Optional:**
- `None` `[N/A]`

## Version History

| Version | Status | Notes |
|---------|--------|-------|
| v0.33.0 | Experimental | Initial release |
| v0.41.0 | Experimental | Renamed `application_id` to `application_object_id` and added more examples|
| v0.42.0 | Experimental | Renamed resource from `ip_application_segment` to `on_premises_ip_application_segment`|

## Example Usage

### Basic Configuration - IP Address

This example demonstrates a minimal configuration targeting a single IP address.

```terraform
# Basic IP Application Segment with single IP address
# This example demonstrates the minimal configuration for an IP application segment
# targeting a single IP address.

resource "microsoft365_graph_beta_applications_on_premises_ip_application_segment" "minimal_ip" {
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

resource "microsoft365_graph_beta_applications_on_premises_ip_application_segment" "ip_range" {
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

resource "microsoft365_graph_beta_applications_on_premises_ip_application_segment" "fqdn" {
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

### Wildcard FQDN

This example shows how to use a wildcard hostname with `destination_type = "fqdn"`.

```terraform
# IP Application Segment with wildcard FQDN
# The application-scoped Graph endpoint accepts wildcard hosts when
# destination_type is fqdn. dnsSuffix is reserved for Quick Access configuration.

resource "microsoft365_graph_beta_applications_on_premises_ip_application_segment" "wildcard_fqdn" {
  application_object_id = "00000000-0000-0000-0000-000000000000"
  destination_host      = "*.internal.contoso.com"
  destination_type      = "fqdn"
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

resource "microsoft365_graph_beta_applications_on_premises_ip_application_segment" "udp_app" {
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
- `destination_type` (String) The type of destination for the application segment.The supported values are: `ipAddress`, `ipRangeCidr`, and `fqdn`. Microsoft Learn lists additional enum members for `ipApplicationSegment`, but this application-scoped Graph endpoint currently rejects `dnsSuffix` for nonweb applications and does not create a usable address range for `ipRange`.
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

# Import an existing IP application segment by application object ID and segment ID
terraform import microsoft365_graph_beta_applications_on_premises_ip_application_segment.example_ip_address "11111111-1111-1111-1111-111111111111/00000000-0000-0000-0000-000000000000"

# The ID format is: {application_object_id}/{ip_application_segment_id}
# You can find both IDs in the Azure Portal or via Microsoft Graph API:
# GET https://graph.microsoft.com/beta/applications/{application-id}/onPremisesPublishing/segmentsConfiguration/microsoft.graph.ipSegmentConfiguration/applicationSegments
```
