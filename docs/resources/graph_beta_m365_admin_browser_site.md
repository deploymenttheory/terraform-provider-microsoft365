---
page_title: "microsoft365_graph_beta_m365_admin_browser_site Resource - terraform-provider-microsoft365"
subcategory: "M365 Admin"
description: |-
  Manages browser sites for Internet Explorer mode in Microsoft Edge using the /admin/edge/internetExplorerMode/siteLists/{siteListId}/sites endpoint. This resource configures specific websites to load in Internet Explorer 11 compatibility mode within Microsoft Edge, enabling legacy web application support through controlled compatibility settings.
---

# microsoft365_graph_beta_m365_admin_browser_site (Resource)

Manages browser sites for Internet Explorer mode in Microsoft Edge using the `/admin/edge/internetExplorerMode/siteLists/{siteListId}/sites` endpoint. This resource configures specific websites to load in Internet Explorer 11 compatibility mode within Microsoft Edge, enabling legacy web application support through controlled compatibility settings.

## Microsoft Documentation

- [browserSite resource type](https://learn.microsoft.com/en-us/graph/api/resources/browsersite?view=graph-rest-beta)
- [Create browserSite](https://learn.microsoft.com/en-us/graph/api/browsersitelist-post-sites?view=graph-rest-beta)

## API Permissions

The following API permissions are required in order to use this resource.

### Microsoft Graph

- **Application**: `BrowserSiteLists.ReadWrite.All`

## Version History

| Version | Status | Notes |
|---------|--------|-------|
| v0.15.0-alpha | Experimental | Initial release

## Example Usage

```terraform
resource "microsoft365_graph_beta_device_and_app_management_browser_site" "example_site" {
  browser_site_list_assignment_id = microsoft365_graph_beta_device_and_app_management_browser_site_list.example.id
  web_url                         = "https://example.com"
  allow_redirect                  = true
  compatibility_mode              = "internetExplorer11"
  comment                         = "Example site for IE mode"
  target_environment              = "internetExplorerMode"
  merge_type                      = "noMerge"

  # Optional: Define custom timeouts
  timeouts = {
    create = "30m"
    read   = "10m"
    update = "30m"
    delete = "30m"
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `allow_redirect` (Boolean) Controls the behavior of redirected sites. If `true`, indicates that the site will open in Internet Explorer 11 or Microsoft Edge even if the site is navigated to as part of a HTTP or meta refresh redirection chain.
- `browser_site_list_assignment_id` (String) The browser site list id to assign this browser site to.
- `compatibility_mode` (String) Controls what compatibility setting is used for specific sites or domains.
- `merge_type` (String) The merge type of the site.
- `target_environment` (String) The target environment that the site should open in.
- `web_url` (String) The URL of the site.

### Optional

- `comment` (String) The comment for the site.
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

### Read-Only

- `created_date_time` (String) The date and time when the site was created.
- `deleted_date_time` (String) The date and time when the site was deleted.
- `id` (String) The unique identifier of the browser site.
- `last_modified_date_time` (String) The date and time when the site was last modified.
- `status` (String) Indicates the status of the site.

<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).
- `delete` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Setting a timeout for a Delete operation is only applicable if changes are saved into state before the destroy operation occurs.
- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Read operations occur during any refresh or planning operation when refresh is enabled.
- `update` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).

## Important Notes

- **Browser Site Management**: This resource manages individual sites within Microsoft Edge browser site lists.
- **Internet Explorer Mode**: Sites are commonly used to configure which websites should open in Internet Explorer mode within Microsoft Edge.
- **Compatibility Lists**: Helps manage legacy web applications that require Internet Explorer for proper functionality.
- **Site Types**: Supports different site types including neutral sites, enterprise mode sites, and sites that should open in Microsoft Edge.
- **URL Patterns**: Supports various URL formats including specific URLs, domains, and wildcard patterns.
- **Policy Integration**: Sites integrate with Microsoft Edge administrative templates and Group Policy settings.
- **Centralized Management**: Provides IT administrators with centralized control over browser behavior for specific websites.

## Import

Import is supported using the following syntax:

```shell
# {resource_id}
terraform import microsoft365_graph_beta_m365_admin_browser_site.example browser-site-id
```