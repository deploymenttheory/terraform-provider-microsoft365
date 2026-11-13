---
page_title: "microsoft365_graph_beta_device_and_app_management_mobile_app_relationship Data Source - terraform-provider-microsoft365"
subcategory: "Device and App Management"

description: |-
  Retrieves a single Intune mobile app relationship (dependency or supersedence) using the /deviceAppManagement/mobileApps/{id}/relationships endpoint. Look the relationship up by relationship_id, or search the relationships of the app named in source_id by target_id, target_display_name, or a custom odata_query. The lookup must match exactly one relationship.
---
# microsoft365_graph_beta_device_and_app_management_mobile_app_relationship (Data Source)

Retrieves a single Intune mobile app relationship (dependency or supersedence) using the `/deviceAppManagement/mobileApps/{id}/relationships` endpoint. Look the relationship up by `relationship_id`, or search the relationships of the app named in `source_id` by `target_id`, `target_display_name`, or a custom `odata_query`. The lookup must match exactly one relationship.

## Microsoft Documentation

- [mobileAppRelationship resource type](https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-mobileapprelationship?view=graph-rest-beta)

## Microsoft Graph API Permissions

The following client `application` permissions are needed in order to use this data source:

**Required:**
- `DeviceManagementApps.Read.All`

**Optional:**
- `None` `[N/A]`

## Version History

| Version | Status | Notes |
|---------|--------|-------|
| v0.14.1-alpha | Experimental | Initial release |
| v1.5.2 | Experimental | Refactored to a single-relationship lookup by `relationship_id`, `target_id`, `target_display_name` or `odata_query` (scoped to `source_id`), reading the per-app relationships endpoint that supports `$filter`. |

## Example Usage

```terraform
# Each lookup must match exactly one relationship (dependency or supersedence).

# Example 1: Look up a relationship by its ID ({source_id}_{target_id})
data "microsoft365_graph_beta_device_and_app_management_mobile_app_relationship" "by_relationship_id" {
  relationship_id = "00000000-0000-0000-0000-000000000001_00000000-0000-0000-0000-000000000002"
}

# Example 2: Look up the relationship between a source app and a target app
data "microsoft365_graph_beta_device_and_app_management_mobile_app_relationship" "by_target_id" {
  source_id = "00000000-0000-0000-0000-000000000001"
  target_id = "00000000-0000-0000-0000-000000000002"
}

# Example 3: Look up a source app's relationship by the target app's display name
data "microsoft365_graph_beta_device_and_app_management_mobile_app_relationship" "by_target_display_name" {
  source_id           = "00000000-0000-0000-0000-000000000001"
  target_display_name = "Microsoft Visual C++ Redistributable"
}

# Example 4: Use an OData filter on the source app's relationships
data "microsoft365_graph_beta_device_and_app_management_mobile_app_relationship" "dependency" {
  source_id   = "00000000-0000-0000-0000-000000000001"
  odata_query = "isof('microsoft.graph.mobileAppDependency')"
}

# Example 5: Find the app that depends on or supersedes this app (an incoming relationship)
data "microsoft365_graph_beta_device_and_app_management_mobile_app_relationship" "parent" {
  source_id   = "00000000-0000-0000-0000-000000000002"
  odata_query = "targetType eq 'parent'"
}

# Example 6: Narrow a lookup with target_type
data "microsoft365_graph_beta_device_and_app_management_mobile_app_relationship" "child_only" {
  source_id   = "00000000-0000-0000-0000-000000000001"
  target_id   = "00000000-0000-0000-0000-000000000002"
  target_type = "child"
}

output "relationship_details" {
  value = {
    id                  = data.microsoft365_graph_beta_device_and_app_management_mobile_app_relationship.by_target_id.id
    relationship_type   = data.microsoft365_graph_beta_device_and_app_management_mobile_app_relationship.by_target_id.relationship_type
    source_display_name = data.microsoft365_graph_beta_device_and_app_management_mobile_app_relationship.by_target_id.source_display_name
    target_display_name = data.microsoft365_graph_beta_device_and_app_management_mobile_app_relationship.by_target_id.target_display_name
    target_type         = data.microsoft365_graph_beta_device_and_app_management_mobile_app_relationship.by_target_id.target_type
    dependency_type     = data.microsoft365_graph_beta_device_and_app_management_mobile_app_relationship.by_target_id.dependency_type
    supersedence_type   = data.microsoft365_graph_beta_device_and_app_management_mobile_app_relationship.by_target_id.supersedence_type
  }
}
```

## Schema

### Optional

- `odata_query` (String) Custom OData `$filter` applied to the relationships of the `source_id` app. Requires `source_id`. Use this for advanced filtering when the standard lookup attributes don't meet your needs. Example: `targetType eq 'child'`
- `relationship_id` (String) The ID of the relationship, in the format `{source_id}_{target_id}`. One of `relationship_id`, `target_id`, `target_display_name`, or `odata_query` must be specified.
- `source_id` (String) The ID of the app whose relationships are searched. Required with `target_id`, `target_display_name`, or `odata_query`. Optional with `relationship_id`, where it is derived from the ID.
- `target_display_name` (String) The display name of the related app (case-insensitive). Requires `source_id`. One of `relationship_id`, `target_id`, `target_display_name`, or `odata_query` must be specified.
- `target_id` (String) The ID of the related app. Requires `source_id`. One of `relationship_id`, `target_id`, `target_display_name`, or `odata_query` must be specified.
- `target_type` (String) Whether the target app is a `child` (the source app depends on or supersedes it) or a `parent` (the target app depends on or supersedes the source app). Can be used as an additional filter when combined with the lookup attributes.
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

### Read-Only

- `dependency_type` (String) For dependencies, how the dependency is handled: `detect` or `autoInstall`.
- `dependent_app_count` (Number) For dependencies, the total number of apps that directly or indirectly depend on the parent app.
- `depends_on_app_count` (Number) For dependencies, the total number of apps the child app directly or indirectly depends on.
- `id` (String) The unique identifier of the relationship, in the format `{source_id}_{target_id}`.
- `relationship_type` (String) The kind of relationship: `dependency` or `supersedence`.
- `source_display_name` (String) The display name of the source app.
- `source_display_version` (String) The display version of the source app.
- `source_publisher_display_name` (String) The publisher display name of the source app.
- `superseded_app_count` (Number) For supersedences, the total number of apps directly or indirectly superseded by the child app.
- `supersedence_type` (String) For supersedences, whether the superseded app is `update`d or `replace`d.
- `superseding_app_count` (Number) For supersedences, the total number of apps directly or indirectly superseding the parent app.
- `target_display_version` (String) The display version of the target app.
- `target_publisher` (String) The publisher of the target app.
- `target_publisher_display_name` (String) The publisher display name of the target app.

<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).
