---
page_title: "microsoft365_graph_beta_device_and_app_management_office_suite_app Resource - terraform-provider-microsoft365"
subcategory: "Device and App Management"

description: |-
  Manages Microsoft 365 Apps (Office Suite) applications using the /deviceAppManagement/mobileApps endpoint.Office Suite Apps enable deployment of Microsoft 365 office applications with configuration options including app exclusions, update channels, localization settings, and shared computer activation for enterprise environments. Learn more here 'https://learn.microsoft.com/en-us/intune/intune-service/apps/apps-add-office365'
---

# microsoft365_graph_beta_device_and_app_management_office_suite_app (Resource)

Manages Microsoft 365 Apps (Office Suite) applications using the `/deviceAppManagement/mobileApps` endpoint.Office Suite Apps enable deployment of Microsoft 365 office applications with configuration options including app exclusions, update channels, localization settings, and shared computer activation for enterprise environments. Learn more here 'https://learn.microsoft.com/en-us/intune/intune-service/apps/apps-add-office365'

## Microsoft Documentation

- [office suite app resource type](https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-officesuiteapp?view=graph-rest-beta)
- [Create office suite app](https://learn.microsoft.com/en-us/graph/api/intune-apps-officesuiteapp-create?view=graph-rest-beta)
- [Office Deployment Tool Configuration Options](https://learn.microsoft.com/en-us/microsoft-365-apps/deploy/office-deployment-tool-configuration-options)

## API Permissions

The following API permissions are required in order to use this resource.

### Microsoft Graph

- **Application**: `DeviceManagementApps.ReadWrite.All`

## Example Usage

```terraform
# Example 1: Office Suite App with Configuration Designer (individual settings)
resource "microsoft365_graph_beta_device_and_app_management_office_suite_app" "office_365_config_designer" {
  display_name       = "Microsoft 365 Apps for Enterprise - Configuration Designer"
  description        = "Microsoft 365 Apps deployed with Configuration Designer settings"
  is_featured        = true
  information_url    = "https://support.microsoft.com/office"
  notes              = "Microsoft 365 Apps configured with specific settings using Configuration Designer approach."
  role_scope_tag_ids = ["0"] # Default role scope tag

  categories = [
    "Business",
    "Productivity",
  ]

  app_icon = {
    web_url_source = "https://your-website.com/office-icon.png"
    //or
    icon_file_path_source = "/path/to/office365-icon.png"
  }

  # Configuration Designer block - use this for individual configuration settings
  configuration_designer = {
    auto_accept_eula = true

    excluded_apps = {
      access               = true  # Exclude Microsoft Access
      bing                 = false # Include Microsoft Search in Bing
      excel                = false # Include Excel
      groove               = true  # Exclude OneDrive for Business (Groove)
      info_path            = true  # Exclude InfoPath
      lync                 = false # Include Skype for Business
      one_drive            = false # Include OneDrive
      one_note             = false # Include OneNote
      outlook              = false # Include Outlook
      power_point          = false # Include PowerPoint
      publisher            = true  # Exclude Publisher
      share_point_designer = true  # Exclude SharePoint Designer
      teams                = false # Include Teams
      visio                = true  # Exclude Visio
      word                 = false # Include Word
    }

    locales_to_install = [
      "en-us", # English (United States)
      "fr-fr", # French (France)
      "de-de", # German (Germany)
    ]

    office_platform_architecture         = "x64"
    office_suite_app_default_file_format = "officeOpenXMLFormat"

    product_ids = [
      "o365ProPlusRetail"
    ]

    should_uninstall_older_versions_of_office = true
    target_version                            = "16.0.19029.20244"
    update_channel                            = "current"
    update_version                            = "" // for latest version, use empty string
    use_shared_computer_activation            = false
  }

  timeouts = {
    create = "10m"
    read   = "5m"
    update = "10m"
    delete = "5m"
  }
}

# Example 2: Office Suite App with XML Configuration
resource "microsoft365_graph_beta_device_and_app_management_office_suite_app" "office_365_xml_config" {
  display_name       = "Microsoft 365 Apps for Enterprise - XML Configuration"
  description        = "Microsoft 365 Apps deployed with XML configuration file"
  is_featured        = true
  information_url    = "https://support.microsoft.com/office"
  notes              = "Microsoft 365 Apps configured using Office Deployment Tool (ODT) XML configuration."
  role_scope_tag_ids = ["0"] # Default role scope tag

  categories = [
    "Business",
    "Productivity",
  ]

  app_icon = {
    web_url_source = "https://your-website.com/office-icon.png"
    //or
    icon_file_path_source = "/path/to/office365-icon.png"
  }

  # XML Configuration block - use this for ODT XML-based configuration
  xml_configuration = {
    office_configuration_xml = <<EOF
<Configuration>
  <Add SourcePath="\\Server\Share" 
      OfficeClientEdition="64"
    Channel="MonthlyEnterprise" >
    <Product ID="O365ProPlusRetail">
      <Language ID="en-us" />
      <Language ID="ja-jp" />
    </Product>
    <Product ID="VisioProRetail">
      <Language ID="en-us" />
      <Language ID="ja-jp" />
    </Product>
  </Add>
  <Updates Enabled="TRUE" 
           UpdatePath="\\Server\Share" />
   <Display Level="None" AcceptEULA="TRUE" />  
</Configuration>
EOF
  }

  timeouts = {
    create = "10m"
    read   = "5m"
    update = "10m"
    delete = "5m"
  }
}

# Example 3: Office Suite App with Business Apps (O365BusinessRetail)
resource "microsoft365_graph_beta_device_and_app_management_office_suite_app" "office_365_business" {
  display_name       = "Microsoft 365 Apps for Business"
  description        = "Microsoft 365 Apps for Business with basic configuration"
  is_featured        = false
  information_url    = "https://support.microsoft.com/office"
  notes              = "Basic Microsoft 365 Apps for Business deployment."
  role_scope_tag_ids = ["0"]

  categories = [
    "Business",
    "Productivity",
  ]

  configuration_designer = {
    auto_accept_eula = true

    excluded_apps = {
      access               = true # Exclude Access (not available in Business)
      bing                 = false
      excel                = false
      groove               = false
      info_path            = true # Exclude InfoPath
      lync                 = false
      one_drive            = false
      one_note             = false
      outlook              = false
      power_point          = false
      publisher            = true # Exclude Publisher
      share_point_designer = true
      teams                = false
      visio                = true # Exclude Visio (not available in Business)
      word                 = false
    }

    locales_to_install                   = ["en-us"]
    office_platform_architecture         = "x64"
    office_suite_app_default_file_format = "officeOpenXMLFormat"

    product_ids = [
      "o365BusinessRetail"
    ]

    should_uninstall_older_versions_of_office = false
    update_channel                            = "current"
    update_version                            = "" // for latest version, use empty string
    use_shared_computer_activation            = false
  }

  timeouts = {
    create = "10m"
    read   = "5m"
    update = "10m"
    delete = "5m"
  }
}

# Example 4: Office Suite App with Project and Visio
resource "microsoft365_graph_beta_device_and_app_management_office_suite_app" "office_365_with_project_visio" {
  display_name       = "Microsoft 365 Apps with Project and Visio"
  description        = "Microsoft 365 Apps including Project and Visio applications"
  is_featured        = true
  information_url    = "https://support.microsoft.com/office"
  notes              = "Complete Microsoft 365 suite including Project and Visio for power users."
  role_scope_tag_ids = ["0"]

  categories = [
    "Business",
    "Productivity",
  ]

  configuration_designer = {
    auto_accept_eula = true

    excluded_apps = {
      access               = true # Exclude Access
      bing                 = false
      excel                = false
      groove               = false
      info_path            = true
      lync                 = false
      one_drive            = false
      one_note             = false
      outlook              = false
      power_point          = false
      publisher            = true
      share_point_designer = true
      teams                = false
      visio                = false # Include Visio
      word                 = false
    }

    locales_to_install = [
      "en-us",
      "es-es", # Spanish (Spain)
      "it-it", # Italian (Italy)
    ]

    office_platform_architecture         = "x64"
    office_suite_app_default_file_format = "officeOpenXMLFormat"

    product_ids = [
      "o365ProPlusRetail",
      "projectProRetail",
      "visioProRetail"
    ]

    should_uninstall_older_versions_of_office = true
    update_channel                            = "monthlyEnterprise"
    update_version                            = "" // for latest version, use empty string
    use_shared_computer_activation            = false
  }

  timeouts = {
    create = "10m"
    read   = "5m"
    update = "10m"
    delete = "5m"
  }
}

# Example 5: Office Suite App for Shared Computer Activation (VDI/Terminal Services)
resource "microsoft365_graph_beta_device_and_app_management_office_suite_app" "office_365_shared_activation" {
  display_name       = "Microsoft 365 Apps - Shared Computer Activation"
  description        = "Microsoft 365 Apps configured for shared computer environments (VDI/Terminal Services)"
  is_featured        = false
  information_url    = "https://support.microsoft.com/office"
  notes              = "Microsoft 365 Apps optimized for shared computer activation scenarios."
  role_scope_tag_ids = ["0"]

  categories = [
    "Business",
    "Productivity",
  ]

  configuration_designer = {
    auto_accept_eula = true

    excluded_apps = {
      access               = true
      bing                 = false
      excel                = false
      groove               = true # Exclude OneDrive for shared environments
      info_path            = true
      lync                 = false
      one_drive            = true # Exclude OneDrive for shared environments
      one_note             = false
      outlook              = false
      power_point          = false
      publisher            = true
      share_point_designer = true
      teams                = false
      visio                = true
      word                 = false
    }

    locales_to_install                   = ["en-us"]
    office_platform_architecture         = "x64"
    office_suite_app_default_file_format = "officeOpenXMLFormat"

    product_ids = [
      "o365ProPlusRetail"
    ]

    should_uninstall_older_versions_of_office = true
    update_channel                            = "deferred" # Use deferred channel for stability in shared environments
    update_version                            = ""         // for latest version, use empty string
    use_shared_computer_activation            = true       # Enable shared computer activation
  }

  timeouts = {
    create = "10m"
    read   = "5m"
    update = "10m"
    delete = "5m"
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `description` (String) A detailed description of the Microsoft 365 Apps application.
- `display_name` (String) The title of the Microsoft 365 Apps application.

### Optional

- `app_icon` (Attributes) The source information for the app icon. Supports various image formats (JPEG, PNG, GIF, etc.) which will be automatically converted to PNG as required by Microsoft Intune. (see [below for nested schema](#nestedatt--app_icon))
- `categories` (Set of String) Set of category names to associate with this application. You can use either the predefined Intune category names like 'Business', 'Productivity', etc., or provide specific category UUIDs. Predefined values include: 'Other apps', 'Books & Reference', 'Data management', 'Productivity', 'Business', 'Development & Design', 'Photos & Media', 'Collaboration & Social', 'Computer management'.
- `configuration_designer` (Attributes) Configuration Designer block for Office Suite App. Use this to configure Office applications using individual settings. (see [below for nested schema](#nestedatt--configuration_designer))
- `information_url` (String) The more information URL.
- `is_featured` (Boolean) The value indicating whether the app is marked as featured by the admin. Default is false.
- `notes` (String) Notes for the app.
- `role_scope_tag_ids` (Set of String) Set of scope tag IDs for this Office Suite app.
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))
- `xml_configuration` (Attributes) XML Configuration block for Office Suite App. Use this to configure Office applications using XML configuration. Learn more here'https://learn.microsoft.com/en-us/microsoft-365-apps/deploy/office-deployment-tool-configuration-options'. (see [below for nested schema](#nestedatt--xml_configuration))

### Read-Only

- `created_date_time` (String) The date and time the app was created. This property is read-only.
- `dependent_app_count` (Number) The total number of dependencies the child app has. This property is read-only.
- `developer` (String) The developer of the app.
- `id` (String) The unique identifier for this Microsoft 365 Apps application
- `is_assigned` (Boolean) The value indicating whether the app is assigned to at least one group. This property is read-only.
- `last_modified_date_time` (String) The date and time the app was last modified. This property is read-only.
- `owner` (String) The owner of the app.
- `privacy_information_url` (String) The privacy statement URL. This is automatically set to Microsoft's privacy statement URL.
- `publisher` (String) The publisher of the Microsoft 365 Apps application. Typically 'Microsoft'.
- `publishing_state` (String) The publishing state for the app. The app cannot be assigned unless the app is published. Possible values are: notPublished, processing, published.
- `superseded_app_count` (Number) The total number of apps this app is directly or indirectly superseded by. This property is read-only.
- `superseding_app_count` (Number) The total number of apps this app directly or indirectly supersedes. This property is read-only.
- `upload_state` (Number) The upload state. Possible values are: 0 - Not Ready, 1 - Ready, 2 - Processing. This property is read-only.

<a id="nestedatt--app_icon"></a>
### Nested Schema for `app_icon`

Optional:

- `icon_file_path_source` (String) The file path to the icon file to be uploaded. Supports various image formats which will be automatically converted to PNG.
- `icon_url_source` (String) The web location of the icon file, can be a http(s) URL. Supports various image formats which will be automatically converted to PNG.


<a id="nestedatt--configuration_designer"></a>
### Nested Schema for `configuration_designer`

Optional:

- `auto_accept_eula` (Boolean) The value to accept the EULA automatically on the enduser's device. Default is true.
- `excluded_apps` (Attributes) The Office applications to exclude from the installation. (see [below for nested schema](#nestedatt--configuration_designer--excluded_apps))
- `locales_to_install` (Set of String) By default, Intune will install Office with the default language of the operating system. Choose any additional languages that you want to install.Must be one of the supported Office locale codes in the format 'xx-xx' (e.g., 'en-us', 'ja-jp').
- `office_platform_architecture` (String) The architecture for which to install Office. Possible values are: 'x86', 'x64'. Default is 'x64'. Changing this forces a new resource to be created.
- `office_suite_app_default_file_format` (String) The default file format for Office applications. Possible values are: 'officeOpenXMLFormat', 'officeOpenDocumentFormat'.
- `product_ids` (Set of String) The Product IDs that represent the Office suite app. Example values: 'o365ProPlusRetail', 'o365BusinessRetail','projectProRetail', 'visioProRetail'.
- `should_uninstall_older_versions_of_office` (Boolean) The value to uninstall any existing MSI versions of Office. Default is false.
- `target_version` (String) The specific version of Office to install. Example: '16.0.19029.20244'.
- `update_channel` (String) The Office update channel. Possible values are: 'current', 'deferred', 'firstReleaseCurrent', 'firstReleaseDeferred', 'monthlyEnterprise'.
- `update_version` (String) The specific update version for the Office installation. Example: '2507'. For latest version, use empty string (default).
- `use_shared_computer_activation` (Boolean) The value to enable shared computer activation for Office. Shared computer activation lets you deploy Microsoft 365 Apps to computers that are used by multiple users. Normally, users can only install and activate Microsoft 365 Apps on a limited number of devices, such as 5 PCs. Using Microsoft 365 Apps with shared computer activation doesn't count against that limit. Default is false.

<a id="nestedatt--configuration_designer--excluded_apps"></a>
### Nested Schema for `configuration_designer.excluded_apps`

Optional:

- `access` (Boolean) The value for if MS Office Access should be excluded or not. Default is false.
- `bing` (Boolean) The value for if Microsoft Search as default in Bing should be excluded or not. Default is false.
- `excel` (Boolean) The value for if MS Office Excel should be excluded or not. Default is false.
- `groove` (Boolean) The value for if MS Office OneDrive for Business – Groove should be excluded or not. Default is false.
- `info_path` (Boolean) The value for if MS Office InfoPath should be excluded or not. Default is false.
- `lync` (Boolean) The value for if MS Office Skype for Business – Lync should be excluded or not. Default is false.
- `one_drive` (Boolean) The value for if MS Office OneDrive should be excluded or not. Default is false.
- `one_note` (Boolean) The value for if MS Office OneNote should be excluded or not. Default is false.
- `outlook` (Boolean) The value for if MS Office Outlook should be excluded or not. Default is false.
- `power_point` (Boolean) The value for if MS Office PowerPoint should be excluded or not. Default is false.
- `publisher` (Boolean) The value for if MS Office Publisher should be excluded or not. Default is false.
- `share_point_designer` (Boolean) The value for if MS Office SharePoint Designer should be excluded or not. Default is false.
- `teams` (Boolean) The value for if MS Office Teams should be excluded or not. Default is false.
- `visio` (Boolean) The value for if MS Office Visio should be excluded or not. Default is false.
- `word` (Boolean) The value for if MS Office Word should be excluded or not. Default is false.



<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).
- `delete` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Setting a timeout for a Delete operation is only applicable if changes are saved into state before the destroy operation occurs.
- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Read operations occur during any refresh or planning operation when refresh is enabled.
- `update` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).


<a id="nestedatt--xml_configuration"></a>
### Nested Schema for `xml_configuration`

Required:

- `office_configuration_xml` (String) The XML configuration file for Office deployment. This is base64 encoded XML content that defines the Office installation configuration.

## Important Notes

- **Windows Specific**: This resource is specifically for managing microsoft 365 applications on Windows devices.

## Import

Import is supported using the following syntax:

```shell
# {resource_id}
terraform import microsoft365_graph_beta_device_and_app_management_office_suite_app.example office-suite-app-id
```

