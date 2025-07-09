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