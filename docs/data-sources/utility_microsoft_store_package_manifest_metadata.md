---
page_title: "microsoft365_utility_microsoft_store_package_manifest_metadata Data Source - terraform-provider-microsoft365"
subcategory: "Utilities"

description: |-
  Retrieves Microsoft Store package manifests by package identifier or search term. Used for winget packages.
---

# microsoft365_utility_microsoft_store_package_manifest_metadata (Data Source)

The `graph_beta_macos_pkg_app_metadata` data source extracts metadata from macOS PKG installer files (.pkg), providing essential information needed for deploying macOS applications in Microsoft Intune. It can extract metadata from either local files or remote URLs.

This data source is particularly useful when creating macOS PKG app resources in Intune, as it automatically extracts critical information such as bundle identifier, version, package IDs, and other metadata required for proper app configuration and management.

## Example Usage

```terraform
# Data source to retrieve WPS Office 2022 package manifest by Package ID
data "microsoft365_utility_microsoft_store_package_manifest_metadata" "wps_office" {
  package_identifier = "XP8M1ZJCZ99QJW"

  timeouts = {
    read = "5m"
  }
}

# Data source to retrieve Microsoft PC Manager package manifest
data "microsoft365_utility_microsoft_store_package_manifest_metadata" "pc_manager" {
  package_identifier = "9PM860492SZD"

  timeouts = {
    read = "3m"
  }
}

# Output examples showing how to access the retrieved data

output "wps_office_package_name" {
  description = "The package name from the default locale"
  value       = data.microsoft365_utility_microsoft_store_package_manifest_metadata.wps_office.manifests[0].versions[0].default_locale.package_name
}

output "wps_office_publisher" {
  description = "The publisher information"
  value       = data.microsoft365_utility_microsoft_store_package_manifest_metadata.wps_office.manifests[0].versions[0].default_locale.publisher
}

output "wps_office_version" {
  description = "The package version"
  value       = data.microsoft365_utility_microsoft_store_package_manifest_metadata.wps_office.manifests[0].versions[0].package_version
}

output "wps_office_description" {
  description = "Short description of the package"
  value       = data.microsoft365_utility_microsoft_store_package_manifest_metadata.wps_office.manifests[0].versions[0].default_locale.short_description
}

output "wps_office_tags" {
  description = "Tags associated with the package"
  value       = data.microsoft365_utility_microsoft_store_package_manifest_metadata.wps_office.manifests[0].versions[0].default_locale.tags
}

output "wps_office_supported_architectures" {
  description = "List of supported architectures"
  value = [
    for installer in data.microsoft365_utility_microsoft_store_package_manifest_metadata.wps_office.manifests[0].versions[0].installers :
    installer.architecture
  ]
}

output "wps_office_installer_types" {
  description = "List of installer types available"
  value = [
    for installer in data.microsoft365_utility_microsoft_store_package_manifest_metadata.wps_office.manifests[0].versions[0].installers :
    installer.installer_type
  ]
}

output "wps_office_available_locales" {
  description = "List of available locales"
  value = [
    for locale in data.microsoft365_utility_microsoft_store_package_manifest_metadata.wps_office.manifests[0].versions[0].locales :
    locale.package_locale
  ]
}

output "wps_office_agreements" {
  description = "Package agreements"
  value = {
    for agreement in data.microsoft365_utility_microsoft_store_package_manifest_metadata.wps_office.manifests[0].versions[0].default_locale.agreements :
    agreement.agreement_label => {
      agreement = agreement.agreement
      url       = agreement.agreement_url
    }
  }
}

# Example of conditional logic based on installer type
locals {
  wps_exe_installers = [
    for installer in data.microsoft365_utility_microsoft_store_package_manifest_metadata.wps_office.manifests[0].versions[0].installers :
    installer if installer.installer_type == "exe"
  ]

  wps_store_installers = [
    for installer in data.microsoft365_utility_microsoft_store_package_manifest_metadata.wps_office.manifests[0].versions[0].installers :
    installer if installer.installer_type == "msstore"
  ]
}

output "wps_exe_installer_info" {
  description = "Information about EXE installers"
  value = length(local.wps_exe_installers) > 0 ? {
    url    = local.wps_exe_installers[0].installer_url
    sha256 = local.wps_exe_installers[0].installer_sha256
    locale = local.wps_exe_installers[0].installer_locale
    size   = local.wps_exe_installers[0].minimum_os_version
  } : null
}

output "wps_store_installer_info" {
  description = "Information about Microsoft Store installers"
  value = length(local.wps_store_installers) > 0 ? {
    product_identifier  = local.wps_store_installers[0].ms_store_product_identifier
    package_family_name = local.wps_store_installers[0].package_family_name
    architecture        = local.wps_store_installers[0].architecture
  } : null
}

# Example for PC Manager (simpler output)
output "pc_manager_info" {
  description = "Basic information about Microsoft PC Manager"
  value = {
    package_id  = data.microsoft365_utility_microsoft_store_package_manifest_metadata.pc_manager.manifests[0].package_identifier
    name        = data.microsoft365_utility_microsoft_store_package_manifest_metadata.pc_manager.manifests[0].versions[0].default_locale.package_name
    publisher   = data.microsoft365_utility_microsoft_store_package_manifest_metadata.pc_manager.manifests[0].versions[0].default_locale.publisher
    description = data.microsoft365_utility_microsoft_store_package_manifest_metadata.pc_manager.manifests[0].versions[0].default_locale.short_description
  }
}

# Example of using the data in other resources (hypothetical)
resource "local_file" "wps_office_manifest" {
  filename = "wps_office_manifest.json"
  content = jsonencode({
    package_identifier = data.microsoft365_utility_microsoft_store_package_manifest_metadata.wps_office.manifests[0].package_identifier
    package_name       = data.microsoft365_utility_microsoft_store_package_manifest_metadata.wps_office.manifests[0].versions[0].default_locale.package_name
    version            = data.microsoft365_utility_microsoft_store_package_manifest_metadata.wps_office.manifests[0].versions[0].package_version
    publisher          = data.microsoft365_utility_microsoft_store_package_manifest_metadata.wps_office.manifests[0].versions[0].default_locale.publisher
    description        = data.microsoft365_utility_microsoft_store_package_manifest_metadata.wps_office.manifests[0].versions[0].default_locale.description
    tags               = data.microsoft365_utility_microsoft_store_package_manifest_metadata.wps_office.manifests[0].versions[0].default_locale.tags
    installers = [
      for installer in data.microsoft365_utility_microsoft_store_package_manifest_metadata.wps_office.manifests[0].versions[0].installers : {
        type         = installer.installer_type
        architecture = installer.architecture
        url          = installer.installer_url
        sha256       = installer.installer_sha256
      }
    ]
  })
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `package_identifier` (String) The specific package identifier to retrieve manifest for. Either this or search_term must be provided.
- `search_term` (String) Search term to find packages. Either this or package_identifier must be provided.
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

### Read-Only

- `manifests` (Attributes List) List of package manifests retrieved. (see [below for nested schema](#nestedatt--manifests))

<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).
- `delete` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Setting a timeout for a Delete operation is only applicable if changes are saved into state before the destroy operation occurs.
- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Read operations occur during any refresh or planning operation when refresh is enabled.
- `update` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).


<a id="nestedatt--manifests"></a>
### Nested Schema for `manifests`

Read-Only:

- `package_identifier` (String) The package identifier.
- `type` (String) The type identifier for the manifest.
- `versions` (Attributes List) List of package versions. (see [below for nested schema](#nestedatt--manifests--versions))

<a id="nestedatt--manifests--versions"></a>
### Nested Schema for `manifests.versions`

Read-Only:

- `default_locale` (Attributes) Default locale information for the package. (see [below for nested schema](#nestedatt--manifests--versions--default_locale))
- `installers` (Attributes List) List of installer information. (see [below for nested schema](#nestedatt--manifests--versions--installers))
- `locales` (Attributes List) List of locale-specific information. (see [below for nested schema](#nestedatt--manifests--versions--locales))
- `package_version` (String) The package version number.
- `type` (String) The type identifier for the version.

<a id="nestedatt--manifests--versions--default_locale"></a>
### Nested Schema for `manifests.versions.default_locale`

Read-Only:

- `agreements` (Attributes List) List of agreements for the package. (see [below for nested schema](#nestedatt--manifests--versions--default_locale--agreements))
- `copyright` (String) The copyright information.
- `description` (String) Detailed description of the package.
- `license` (String) The license information.
- `package_locale` (String) The locale code (e.g., en-us).
- `package_name` (String) The package name.
- `privacy_url` (String) The privacy policy URL.
- `publisher` (String) The publisher name.
- `publisher_support_url` (String) The publisher support URL.
- `publisher_url` (String) The publisher website URL.
- `short_description` (String) Short description of the package.
- `tags` (List of String) List of tags associated with the package.
- `type` (String) The type identifier for the locale.

<a id="nestedatt--manifests--versions--default_locale--agreements"></a>
### Nested Schema for `manifests.versions.default_locale.agreements`

Read-Only:

- `agreement` (String) The agreement text.
- `agreement_label` (String) The agreement label.
- `agreement_url` (String) The agreement URL.
- `type` (String) The type identifier for the agreement.



<a id="nestedatt--manifests--versions--installers"></a>
### Nested Schema for `manifests.versions.installers`

Read-Only:

- `apps_and_features_entries` (Attributes List) List of Apps and Features entries. (see [below for nested schema](#nestedatt--manifests--versions--installers--apps_and_features_entries))
- `architecture` (String) The target architecture (e.g., x86, x64, arm64).
- `download_command_prohibited` (Boolean) Whether download command is prohibited.
- `expected_return_codes` (Attributes List) List of expected return codes. (see [below for nested schema](#nestedatt--manifests--versions--installers--expected_return_codes))
- `installer_locale` (String) Locale for the installer.
- `installer_sha256` (String) SHA256 hash of the installer.
- `installer_success_codes` (List of Number) List of success codes for the installer.
- `installer_switches` (Attributes) Installer switches information. (see [below for nested schema](#nestedatt--manifests--versions--installers--installer_switches))
- `installer_type` (String) The installer type (e.g., msstore, exe).
- `installer_url` (String) URL to download the installer.
- `markets` (Attributes) Market information for the installer. (see [below for nested schema](#nestedatt--manifests--versions--installers--markets))
- `minimum_os_version` (String) Minimum OS version required.
- `ms_store_product_identifier` (String) Microsoft Store product identifier.
- `package_family_name` (String) The package family name.
- `scope` (String) The installation scope (user or machine).
- `type` (String) The type identifier for the installer.

<a id="nestedatt--manifests--versions--installers--apps_and_features_entries"></a>
### Nested Schema for `manifests.versions.installers.apps_and_features_entries`

Read-Only:

- `display_name` (String) Display name in Apps and Features.
- `display_version` (String) Display version in Apps and Features.
- `installer_type` (String) Installer type.
- `product_code` (String) Product code.
- `publisher` (String) Publisher name in Apps and Features.
- `type` (String) The type identifier for the entry.


<a id="nestedatt--manifests--versions--installers--expected_return_codes"></a>
### Nested Schema for `manifests.versions.installers.expected_return_codes`

Read-Only:

- `installer_return_code` (Number) The return code.
- `return_response` (String) The response description for the return code.
- `type` (String) The type identifier for return code.


<a id="nestedatt--manifests--versions--installers--installer_switches"></a>
### Nested Schema for `manifests.versions.installers.installer_switches`

Read-Only:

- `silent` (String) Silent installation switch.
- `type` (String) The type identifier for installer switches.


<a id="nestedatt--manifests--versions--installers--markets"></a>
### Nested Schema for `manifests.versions.installers.markets`

Read-Only:

- `allowed_markets` (List of String) List of allowed markets.
- `type` (String) The type identifier for markets.



<a id="nestedatt--manifests--versions--locales"></a>
### Nested Schema for `manifests.versions.locales`

Read-Only:

- `copyright` (String) The copyright information.
- `description` (String) Detailed description of the package.
- `license` (String) The license information.
- `package_locale` (String) The locale code (e.g., en-us).
- `package_name` (String) The package name.
- `privacy_url` (String) The privacy policy URL.
- `publisher` (String) The publisher name.
- `publisher_support_url` (String) The publisher support URL.
- `publisher_url` (String) The publisher website URL.
- `short_description` (String) Short description of the package.
- `tags` (List of String) List of tags associated with the package.
- `type` (String) The type identifier for the locale.