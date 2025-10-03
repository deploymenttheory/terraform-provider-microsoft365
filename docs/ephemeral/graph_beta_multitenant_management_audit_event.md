---
page_title: "<no value> <no value> - <no value>"
subcategory: "Multitenant Management"

description: |-
  Retrieves audit events from Microsoft 365 managed tenants as an ephemeral resource.
---

# <no value> (<no value>)

Retrieves audit events from Microsoft 365 managed tenants as an ephemeral resource. This does not persist in state and fetches fresh data on each execution.

## Microsoft Documentation

- [auditEvent resource type](https://learn.microsoft.com/en-us/graph/api/resources/managedtenants-auditevent?view=graph-rest-beta)

## API Permissions

The following API permissions are required in order to use this data source.

### Microsoft Graph

- **Application**: `ManagedTenants.Read.All`, `ManagedTenants.ReadWrite.All`.

## Version History

| Version | Status | Notes |
|---------|--------|-------|
| v0.31.0-alpha | Experimental | Initial release |

## Important Notes
This method supports some of the OData query parameters to help customize the response. For general information, see OData query parameters.
## Example Usage

```terraform
# Example 1: Get all audit events
ephemeral "microsoft365_multitenant_management_audit_events" "all_events" {
  filter_type = "all"
}

# Example 2: Get a specific audit event by ID
ephemeral "microsoft365_multitenant_management_audit_events" "specific_event" {
  filter_type  = "id"
  filter_value = "12345678-1234-1234-1234-123456789abc"
}

# Example 3: Filter audit events by activity display name
ephemeral "microsoft365_multitenant_management_audit_events" "activity_filter" {
  filter_type  = "display_name"
  filter_value = "sign-in"
}

# Example 4: Use OData filtering for recent events
ephemeral "microsoft365_multitenant_management_audit_events" "recent_events" {
  filter_type  = "odata"
  odata_filter = "activityDateTime ge 2024-01-01T00:00:00Z"
  odata_top    = 100
  odata_orderby = "activityDateTime desc"
}

# Example 5: OData filtering with specific categories and pagination
ephemeral "microsoft365_multitenant_management_audit_events" "categorized_events" {
  filter_type   = "odata"
  odata_filter  = "category eq 'ApplicationManagement' or category eq 'UserManagement'"
  odata_top     = 50
  odata_skip    = 0
  odata_select  = "id,activity,category,activityDateTime,initiatedByUpn"
  odata_orderby = "activityDateTime desc"
}

# Example 6: Filter by specific user activities
ephemeral "microsoft365_multitenant_management_audit_events" "user_events" {
  filter_type  = "odata"
  odata_filter = "initiatedByUpn eq 'admin@contoso.com'"
  odata_top    = 25
}

# Output examples to demonstrate accessing the data
output "all_events_count" {
  description = "Total number of audit events"
  value       = length(ephemeral.microsoft365_multitenant_management_audit_events.all_events.items)
}

output "recent_event_activities" {
  description = "Activities from recent events"
  value = [
    for event in ephemeral.microsoft365_multitenant_management_audit_events.recent_events.items : {
      activity = event.activity
      datetime = event.activity_date_time
      user     = event.initiated_by_upn
    }
  ]
}

output "specific_event_details" {
  description = "Details of the specific event"
  value = length(ephemeral.microsoft365_multitenant_management_audit_events.specific_event.items) > 0 ? {
    id       = ephemeral.microsoft365_multitenant_management_audit_events.specific_event.items[0].id
    activity = ephemeral.microsoft365_multitenant_management_audit_events.specific_event.items[0].activity
    category = ephemeral.microsoft365_multitenant_management_audit_events.specific_event.items[0].category
    datetime = ephemeral.microsoft365_multitenant_management_audit_events.specific_event.items[0].activity_date_time
    user     = ephemeral.microsoft365_multitenant_management_audit_events.specific_event.items[0].initiated_by_upn
  } : null
}
```

Schema documentation not available.