---
page_title: "microsoft365_graph_beta_device_and_app_management_win32_app Resource - terraform-provider-microsoft365"
subcategory: "Device and App Management"

description: |-
  Manages Win32 applications using the /deviceAppManagement/mobileApps endpoint. Win apps enable deployment of custom Windows applications (.exe, .msi) with advanced installation logic, detection rules, and dependency management for enterprise software distribution. They must be wrapped in the .intunewin file type.'https://learn.microsoft.com/en-us/intune/intune-service/apps/apps-win32-app-management'
---

# microsoft365_graph_beta_device_and_app_management_win32_app (Resource)

Manages Win32 applications using the `/deviceAppManagement/mobileApps` endpoint. Win apps enable deployment of custom Windows applications (.exe, .msi) with advanced installation logic, detection rules, and dependency management for enterprise software distribution. They must be wrapped in the .intunewin file type.'https://learn.microsoft.com/en-us/intune/intune-service/apps/apps-win32-app-management'

## Microsoft Documentation

- [win32LobApp resource type](https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-win32lobapp?view=graph-rest-beta)
- [Create win32LobApp](https://learn.microsoft.com/en-us/graph/api/intune-apps-win32lobapp-create?view=graph-rest-beta)

## API Permissions

The following API permissions are required in order to use this resource.

### Microsoft Graph

- **Application**: `DeviceManagementApps.ReadWrite.All`

## Example Usage

```terraform
resource "microsoft365_graph_beta_device_and_app_management_win32_app" "example" {
  allow_available_uninstall = true

  app_installer = {
    installer_file_path_source = "/Users/dafyddwatkins/Downloads/win_32_lob_app/Firefox_Setup_140.0.4.intunewin"
  }

  // Optional
  app_icon = {
    icon_url_source = "https://images.g2crowd.com/uploads/product/image/large_detail/large_detail_c65522c8f7bacfdc0684fc8d26bba00d/mozilla-firefox.png"
  }

  // Optional
  categories = [
    microsoft365_graph_beta_device_and_app_management_application_category.web_browser.id,
    "Business",
    "Productivity",
  ]

  description     = "Mozilla Firefox 140.0.4 x64 en-US"
  publisher       = "Mozilla"
  developer       = "Mozilla"
  display_name    = "Mozilla Firefox 140.0.4 x64 en-US"
  display_version = "140.0.4.0"
  file_name       = "Firefox_Setup_140.0.4.intunewin"
  information_url = "https://www.mozilla.org/firefox/"

  allowed_architectures             = ["x64", "arm64"]
  minimum_supported_windows_release = "Windows11_23H2"

  install_experience = {
    device_restart_behavior = "allow"
    max_run_time_in_minutes = 60
    run_as_account          = "system"
  }
  setup_file_path        = "Firefox Setup 140.0.4.msi"
  install_command_line   = "msiexec /i \"Firefox Setup 140.0.4.msi\" /qn"
  uninstall_command_line = "msiexec /x {1294A4C5-9977-480F-9497-C0EA1E630130} /qn"

  msi_information = {
    package_type    = "perMachine"
    product_code    = "{1294A4C5-9977-480F-9497-C0EA1E630130}"
    product_name    = "Mozilla Firefox 140.0.4 x64 en-US"
    publisher       = "Mozilla"
    product_version = "140.0.4.0"
    requires_reboot = false
    upgrade_code    = "{3118AB4C-B433-4FBB-B9FA-8F9CA4B5C103}"
  }

  # Detection Rules
  rules = [
    # Rule 0: PowerShell script rule for detection - FILE SYSTEM CHECK
    {
      rule_type                             = "detection"
      rule_sub_type                         = "powershell_script"
      enforce_signature_check               = true
      run_as_32_bit                         = true
      powershell_script_rule_operation_type = "notConfigured"
      lob_app_rule_operator                 = "notConfigured"
      script_content                        = <<EOT
# Detection Method 1: File System Check
# This script checks for Firefox executable in standard installation paths

$firefoxPaths = @(
    "$${env:ProgramFiles}\\Mozilla Firefox\\firefox.exe",
    "$${env:ProgramFiles(x86)}\\Mozilla Firefox\\firefox.exe"
)

$firefoxInstalled = $false
$targetVersion = "140.0.4"

foreach ($path in $firefoxPaths) {
    if (Test-Path -Path $path -PathType Leaf) {
        $fileVersion = (Get-Item $path).VersionInfo.FileVersion
        
        # Check if version meets our requirements (140.0.4)
        if ($fileVersion -like "*$targetVersion*") {
            $firefoxInstalled = $true
            Write-Output "Firefox $fileVersion detected at: $path"
            break
        } else {
            Write-Output "Firefox found at $path but version $fileVersion does not match required version $targetVersion"
        }
    }
}

# Return exit code for detection script
if ($firefoxInstalled) {
    exit 0  # Success - Firefox with correct version is installed
} else {
    Write-Output "Firefox $targetVersion not detected"
    exit 1  # Failure - Firefox with correct version is not installed
}
EOT
    },

    # Rule 1: File system rule for requirement
    {
      rule_type                  = "requirement"
      rule_sub_type              = "file_system"
      check_32_bit_on_64_system  = true
      path                       = "c:\\thing"
      file_or_folder_name        = "folder_name"
      file_system_operation_type = "exists"
      lob_app_rule_operator      = "notConfigured"
    },

    # Rule 2: Registry rule for requirement
    {
      rule_type                 = "requirement"
      rule_sub_type             = "registry"
      check_32_bit_on_64_system = false
      key_path                  = "kay_path"
      value_name                = "key_value"
      operation_type            = "doesNotExist"
      lob_app_rule_operator     = "notConfigured"
    },

    # Rule 3: PowerShell script rule for requirement with dateTime comparison - REGISTRY CHECK
    {
      rule_type                             = "requirement"
      rule_sub_type                         = "powershell_script"
      display_name                          = "Firefox Registry Detection"
      enforce_signature_check               = true
      run_as_32_bit                         = true
      run_as_account                        = "user"
      powershell_script_rule_operation_type = "dateTime"
      lob_app_rule_operator                 = "equal"
      comparison_value                      = "2025-07-08T23:00:00.000Z"
      script_content                        = <<EOT
# Detection Method 2: Registry Check
# This script checks for Firefox in the Windows registry

$firefoxInstalled = $false
$targetVersion = "140.0.4"
$registryPaths = @(
    "HKLM:\\SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\App Paths\\firefox.exe",
    "HKLM:\\SOFTWARE\\Mozilla\\Mozilla Firefox",
    "HKLM:\\SOFTWARE\\WOW6432Node\\Mozilla\\Mozilla Firefox"
)

foreach ($regPath in $registryPaths) {
    if (Test-Path -Path $regPath) {
        Write-Output "Firefox registry entry found at: $regPath"
        
        # Try to get version information from registry
        try {
            if ($regPath -eq "HKLM:\\SOFTWARE\\Mozilla\\Mozilla Firefox") {
                $currentVersion = Get-ItemProperty -Path $regPath -Name "CurrentVersion" -ErrorAction SilentlyContinue
                if ($currentVersion -and $currentVersion.CurrentVersion -like "*$targetVersion*") {
                    $firefoxInstalled = $true
                    Write-Output "Firefox version $($currentVersion.CurrentVersion) confirmed in registry"
                    break
                }
            } elseif ($regPath -like "*App Paths*") {
                $defaultPath = (Get-ItemProperty -Path $regPath -Name "(Default)" -ErrorAction SilentlyContinue)."(Default)"
                if ($defaultPath -and (Test-Path $defaultPath)) {
                    $fileVersion = (Get-Item $defaultPath).VersionInfo.FileVersion
                    if ($fileVersion -like "*$targetVersion*") {
                        $firefoxInstalled = $true
                        Write-Output "Firefox version $fileVersion confirmed via App Paths registry"
                        break
                    }
                }
            }
        } catch {
            Write-Output "Error checking registry: $_"
        }
    }
}

# Return exit code for detection script
if ($firefoxInstalled) {
    exit 0  # Success - Firefox with correct version is installed
} else {
    Write-Output "Firefox $targetVersion not detected in registry"
    exit 1  # Failure - Firefox with correct version is not installed
}
EOT
    },

    # Rule 4: PowerShell script rule for requirement with float comparison
    {
      rule_type                             = "requirement"
      rule_sub_type                         = "powershell_script"
      display_name                          = "script_2.ps1"
      enforce_signature_check               = false
      run_as_32_bit                         = false
      run_as_account                        = "user"
      powershell_script_rule_operation_type = "float"
      lob_app_rule_operator                 = "notEqual"
      comparison_value                      = "3.14159"
      script_content                        = <<EOT
$firefoxPaths = @(
    "$${env:ProgramFiles}\\Mozilla Firefox\\firefox.exe",
    "$${env:ProgramFiles(x86)}\\Mozilla Firefox\\firefox.exe"
)

$firefoxInstalled = $false

foreach ($path in $firefoxPaths) {
    if (Test-Path -Path $path -PathType Leaf) {
        $firefoxInstalled = $true
        $version = (Get-Item $path).VersionInfo.ProductVersion
        Write-Output "Firefox detected: $path (Version: $version)"
        break
    }
}

# Also check registry for installed Firefox
if (-not $firefoxInstalled) {
    $regPaths = @(
        "HKLM:\\SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\App Paths\\firefox.exe",
        "HKLM:\\SOFTWARE\\Mozilla\\Mozilla Firefox"
    )
    
    foreach ($regPath in $regPaths) {
        if (Test-Path -Path $regPath) {
            $firefoxInstalled = $true
            Write-Output "Firefox detected in registry: $regPath"
            break
        }
    }
}

# Return exit code for detection script
if ($firefoxInstalled) {
    exit 0  # Success - Firefox is installed
} else {
    Write-Output "Firefox not detected"
    exit 1  # Failure - Firefox is not installed
}
EOT
    },

    # Rule 5: PowerShell script rule for requirement with dateTime comparison - WMI CHECK
    {
      rule_type                             = "requirement"
      rule_sub_type                         = "powershell_script"
      display_name                          = "Firefox WMI Detection"
      enforce_signature_check               = false
      run_as_32_bit                         = false
      run_as_account                        = "system"
      powershell_script_rule_operation_type = "dateTime"
      lob_app_rule_operator                 = "lessThanOrEqual"
      comparison_value                      = "2025-06-30T23:00:00.000Z"
      script_content                        = <<EOT
# Detection Method 3: WMI/CIM Query
# This script uses WMI/CIM to check for Firefox in installed programs

$firefoxInstalled = $false
$targetVersion = "140.0.4"
$targetPublisher = "Mozilla"

try {
    # Use CIM query to find installed applications (modern approach)
    $installedPrograms = Get-CimInstance -ClassName Win32_Product -ErrorAction SilentlyContinue | 
        Where-Object { $_.Name -like "*Firefox*" -and $_.Vendor -like "*$targetPublisher*" }
    
    if ($installedPrograms) {
        foreach ($program in $installedPrograms) {
            Write-Output "Found Firefox via WMI: $($program.Name), Version: $($program.Version), Vendor: $($program.Vendor)"
            
            # Check if version matches our target
            if ($program.Version -like "*$targetVersion*") {
                $firefoxInstalled = $true
                Write-Output "Firefox $targetVersion confirmed via WMI"
                break
            }
        }
    } else {
        Write-Output "No Firefox installation found via primary WMI query"
        
        # Fallback to registry-based installed programs check (more reliable than WMI sometimes)
        $uninstallKeys = @(
            "HKLM:\\SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Uninstall",
            "HKLM:\\SOFTWARE\\WOW6432Node\\Microsoft\\Windows\\CurrentVersion\\Uninstall"
        )
        
        foreach ($key in $uninstallKeys) {
            if (Test-Path $key) {
                Get-ChildItem $key | ForEach-Object {
                    $programKey = $_
                    $programInfo = Get-ItemProperty $programKey.PSPath
                    
                    if ($programInfo.DisplayName -like "*Firefox*" -and $programInfo.Publisher -like "*$targetPublisher*") {
                        Write-Output "Found Firefox via registry uninstall key: $($programInfo.DisplayName), Version: $($programInfo.DisplayVersion)"
                        
                        if ($programInfo.DisplayVersion -like "*$targetVersion*") {
                            $firefoxInstalled = $true
                            Write-Output "Firefox $targetVersion confirmed via registry uninstall key"
                            break
                        }
                    }
                }
            }
            
            if ($firefoxInstalled) { break }
        }
    }
} catch {
    Write-Output "Error querying installed applications: $_"
}

# Return exit code for detection script
if ($firefoxInstalled) {
    exit 0  # Success - Firefox with correct version is installed
} else {
    Write-Output "Firefox $targetVersion not detected via WMI/registry"
    exit 1  # Failure - Firefox with correct version is not installed
}
EOT
    }
  ]

  return_codes = [
    {
      return_code = 0
      type        = "success"
    },
    {
      return_code = 1707
      type        = "success"
    },
    {
      return_code = 3010
      type        = "softReboot"
    },
    {
      return_code = 1641
      type        = "hardReboot"
    },
    {
      return_code = 1618
      type        = "retry"
    }
  ]
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `allow_available_uninstall` (Boolean) When TRUE, indicates that uninstall is supported from the company portal for the Windows app (Win32) with an Available assignment. When FALSE, indicates that uninstall is not supported for the Windows app (Win32) with an Available assignment. Default value is FALSE.
- `display_name` (String) The admin provided or imported title of the app.
- `file_name` (String) The name of the main Lob application file.
- `install_command_line` (String) The command line to install this app. Typically formatted as 'msiexec /i "application_name.msi" /qn'
- `minimum_supported_windows_release` (String) The value for the minimum supported windows release.
- `msi_information` (Attributes) The MSI details if this Win32 app is an MSI app. (see [below for nested schema](#nestedatt--msi_information))
- `publisher` (String) The publisher of the Intune macOS pkg application.
- `uninstall_command_line` (String) The command line to uninstall this app. Typically formatted as 'msiexec /x {00000000-0000-0000-0000-000000000000} /qn'

### Optional

- `allowed_architectures` (Set of String) The Windows architecture(s) for which this app can run on. Possible values are: none, x64, x86, arm64.
- `app_icon` (Attributes) The source information for the app icon. Supports various image formats (JPEG, PNG, GIF, etc.) which will be automatically converted to PNG as required by Microsoft Intune. (see [below for nested schema](#nestedatt--app_icon))
- `app_installer` (Attributes) Metadata related to the win32 lob app installer file, such as size and checksums. This is automatically computed during app creation and updates. (see [below for nested schema](#nestedatt--app_installer))
- `categories` (Set of String) Set of category names to associate with this application. You can use either thebpredefined Intune category names like 'Business', 'Productivity', etc., or provide specific category UUIDs. Predefined values include: 'Other apps', 'Books & Reference', 'Data management', 'Productivity', 'Business', 'Development & Design', 'Photos & Media', 'Collaboration & Social', 'Computer management'.
- `content_version` (Attributes List) The committed content version of the app, including its files. Only the currently committed version is shown. (see [below for nested schema](#nestedatt--content_version))
- `description` (String) The description of the app.
- `detection_rules` (Attributes List) The detection rules to detect Win32 Line of Business (LoB) app. (see [below for nested schema](#nestedatt--detection_rules))
- `developer` (String) The developer of the app.
- `display_version` (String) The version displayed in the UX for this app.
- `information_url` (String) The more information Url.
- `install_experience` (Attributes) The install experience for this app. (see [below for nested schema](#nestedatt--install_experience))
- `is_featured` (Boolean) The value indicating whether the app is marked as featured by the admin.
- `minimum_cpu_speed_in_mhz` (Number) The value for the minimum CPU speed which is required to install this app.
- `minimum_free_disk_space_in_mb` (Number) The value for the minimum free disk space which is required to install this app.
- `minimum_memory_in_mb` (Number) The value for the minimum physical memory which is required to install this app.
- `minimum_number_of_processors` (Number) The value for the minimum number of processors which is required to install this app.
- `notes` (String) Notes for the app.
- `owner` (String) The owner of the app.
- `privacy_information_url` (String) The privacy statement Url.
- `requirement_rules` (Attributes List) The requirement rules to detect Win32 Line of Business (LoB) app. (see [below for nested schema](#nestedatt--requirement_rules))
- `return_codes` (Attributes List) The return codes for post installation behavior. (see [below for nested schema](#nestedatt--return_codes))
- `role_scope_tag_ids` (Set of String) Set of scope tag IDs for this Settings Catalog template profile.
- `rules` (Attributes List) The detection and requirement rules for this app. (see [below for nested schema](#nestedatt--rules))
- `setup_file_path` (String) The relative path of the setup file in the encrypted Win32LobApp package.
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

### Read-Only

- `committed_content_version` (String) The internal committed content version.
- `created_date_time` (String) The date and time the app was created. This property is read-only.
- `dependent_app_count` (Number) The total number of dependencies the child app has. This property is read-only.
- `id` (String) The unique identifier for this Intune win32 lob application
- `is_assigned` (Boolean) The value indicating whether the app is assigned to at least one group. This property is read-only.
- `last_modified_date_time` (String) The date and time the app was last modified. This property is read-only.
- `publishing_state` (String) The publishing state for the app. The app cannot be assigned unless the app is published. This property is read-only. Possible values are: notPublished, processing, published.
- `size` (Number) The total size, including all uploaded files. This property is read-only.
- `superseded_app_count` (Number) The total number of apps this app is directly or indirectly superseded by. This property is read-only.
- `superseding_app_count` (Number) The total number of apps this app directly or indirectly supersedes. This property is read-only.
- `upload_state` (Number) The upload state. Possible values are: 0 - Not Ready, 1 - Ready, 2 - Processing. This property is read-only.

<a id="nestedatt--msi_information"></a>
### Nested Schema for `msi_information`

Required:

- `package_type` (String) The MSI package type. Possible values are: perMachine, perUser.
- `product_version` (String) The MSI product version.
- `requires_reboot` (Boolean) A value indicating whether the MSI app requires a reboot.
- `upgrade_code` (String) The MSI upgrade code.

Optional:

- `product_code` (String) The MSI product code.
- `product_name` (String) The MSI product name.
- `publisher` (String) The MSI publisher.


<a id="nestedatt--app_icon"></a>
### Nested Schema for `app_icon`

Optional:

- `icon_file_path_source` (String) The file path to the icon file to be uploaded. Supports various image formats which will be automatically converted to PNG.
- `icon_url_source` (String) The web location of the icon file, can be a http(s) URL. Supports various image formats which will be automatically converted to PNG.


<a id="nestedatt--app_installer"></a>
### Nested Schema for `app_installer`

Optional:

- `installer_file_path_source` (String) The path to the win32 lob app installer file to be uploaded. The file must be a valid `.intunewin` file. Value is not returned by API call.
- `installer_url_source` (String) The web location of the win32 lob app installer file, can be a http(s) URL. The file must be a valid `.intunewin` file. Value is not returned by API call.


<a id="nestedatt--content_version"></a>
### Nested Schema for `content_version`

Read-Only:

- `files` (Attributes Set) The files associated with this content version. (see [below for nested schema](#nestedatt--content_version--files))
- `id` (String) The unique identifier for this content version. This ID is assigned during creation of the content version. Read-only.

<a id="nestedatt--content_version--files"></a>
### Nested Schema for `content_version.files`

Read-Only:

- `azure_storage_uri` (String) Indicates the Azure Storage URI that the file is uploaded to. Read-only.
- `azure_storage_uri_expiration` (String) Indicates the date and time when the Azure storage URI expires, in ISO 8601 format. Read-only.
- `created_date_time` (String) Indicates created date and time associated with app content file, in ISO 8601 format. Read-only.
- `is_committed` (Boolean) A value indicating whether the file is committed. A committed app content file has been fully uploaded and validated by the Intune service. Read-only.
- `is_dependency` (Boolean) Indicates whether this content file is a dependency for the main content file.
- `is_framework_file` (Boolean) Indicates whether this content file is a framework file.
- `name` (String) Indicates the name of the file.
- `size` (Number) Indicates the original size of the file, in bytes.
- `size_encrypted` (Number) Indicates the size of the file after encryption, in bytes.
- `upload_state` (String) Indicates the state of the current upload request. This property is read-only.



<a id="nestedatt--detection_rules"></a>
### Nested Schema for `detection_rules`

Required:

- `detection_type` (String) The detection rule type. Possible values are: registry, msi_information, file_system, powershell_script.

Optional:

- `check_32_bit_on_64_system` (Boolean) Whether to check 32-bit registry or file system on 64-bit system. Applicable for registry, file_system, and PowerShell script detection.
- `detection_value` (String) The registry detection value for registry detection.
- `enforce_signature_check` (Boolean) Whether to enforce signature checking for the PowerShell script.
- `file_or_folder_name` (String) The file name for file system detection.
- `file_path` (String) The file path for file system detection.
- `file_system_detection_operator` (String) The filesystem detection operator for filesystem detection. Possible values are: notConfigured, equal, notEqual, greaterThan, greaterThanOrEqual, lessThan, lessThanOrEqual. Used for registry and file_system detection types.
- `file_system_detection_type` (String) The comparison operator for detection. Possible values are: notConfigured, exists,modifiedDate, createdDate, version, sizeInMB, doesNotExist.
- `key_path` (String) The registry key path for registry detection.
- `product_code` (String) The MSI product code for MSI detection.
- `product_version` (String) The MSI product version for MSI detection.
- `product_version_operator` (String) The MSI product version operator for MSI detection.Possible values are: notConfigured, equal, notEqual, greaterThan, greaterThanOrEqual, lessThan, lessThanOrEqual. Used for registry and file_system detection types.
- `registry_detection_operator` (String) The registry detection operator for registry detection. Possible values are: notConfigured, equal, notEqual, greaterThan, greaterThanOrEqual, lessThan, lessThanOrEqual. Used for registry and file_system detection types.
- `registry_detection_type` (String) The registry detection type for registry detection. Possible values are: notConfigured, exists, doesNotExist, string, integer, version.
- `run_as_32_bit` (Boolean) Whether to run the PowerShell script in 32-bit mode on 64-bit systems.
- `script_content` (String) The PowerShell script content to run for script detection.This will be base64-encoded before being sent to the API. Supports PowerShell 5.1 and PowerShell 7.0.
- `value_name` (String) The registry value name for registry detection.


<a id="nestedatt--install_experience"></a>
### Nested Schema for `install_experience`

Optional:

- `device_restart_behavior` (String) The device restart behavior. Possible values are: basedOnReturnCode, allow, suppress, force.
- `max_run_time_in_minutes` (Number) The maximum run time in minutes for the installation.
- `run_as_account` (String) The execution context. Possible values are: system, user.


<a id="nestedatt--requirement_rules"></a>
### Nested Schema for `requirement_rules`

Required:

- `requirement_type` (String) The requirement rule type. Possible values are: registry, file, script.

Optional:

- `check_32_bit_on_64_system` (Boolean) A value indicating whether to check 32-bit on 64-bit system.
- `detection_type` (String) The detection type for registry requirement.
- `detection_value` (String) The value to check for.
- `file_or_folder_name` (String) The file or folder name to check for file requirement.
- `key_path` (String) The registry key path for registry requirement.
- `operator` (String) The operator for the requirement. Possible values are: notConfigured, equal, notEqual, greaterThan, greaterThanOrEqual, lessThan, lessThanOrEqual.
- `value_name` (String) The registry value name for registry requirement.


<a id="nestedatt--return_codes"></a>
### Nested Schema for `return_codes`

Required:

- `return_code` (Number) The return code.
- `type` (String) The return code type. Possible values are: failed, success, softReboot, hardReboot, retry.


<a id="nestedatt--rules"></a>
### Nested Schema for `rules`

Required:

- `rule_sub_type` (String) The rule sub-type. Possible values are: registry, file_system, powershell_script.
- `rule_type` (String) The rule type. Possible values are: detection, requirement.

Optional:

- `check_32_bit_on_64_system` (Boolean) A value indicating whether to check 32-bit on 64-bit system.
- `comparison_value` (String) The value to compare against.
- `display_name` (String) The display name for PowerShell script rules.
- `enforce_signature_check` (Boolean) Whether to enforce signature checking for PowerShell script rules.
- `file_or_folder_name` (String) The file or folder name to check. Required for file_system rules.
- `file_system_operation_type` (String) The operation type for file system rules. Possible values are: notConfigured, exists, modifiedDate, createdDate, version, sizeInMB, doesNotExist.
- `key_path` (String) The registry key path to detect or check. Required for registry rules.
- `lob_app_rule_operator` (String) The operator for the rule. Possible values are: notConfigured, equal, notEqual, greaterThan, greaterThanOrEqual, lessThan, lessThanOrEqual.
- `operation_type` (String) The operation type for registry rules. Possible values are: notConfigured, exists, doesNotExist, string, integer, version.
- `path` (String) The path to the file or folder. Required for file_system rules.
- `powershell_script_rule_operation_type` (String) The operation type for PowerShell script rules. Possible values are: notConfigured, string, dateTime, integer, float, version, boolean.
- `run_as_32_bit` (Boolean) Whether to run PowerShell scripts in 32-bit mode on 64-bit systems.
- `run_as_account` (String) The execution context for PowerShell scripts. Possible values are: system, user.
- `script_content` (String) The PowerShell script content to run for script detection.This will be base64-encoded before being sent to the API. Supports PowerShell 5.1 and PowerShell 7.0.
- `value_name` (String) The registry value name to detect or check. Required for registry rules.


<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).
- `delete` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Setting a timeout for a Delete operation is only applicable if changes are saved into state before the destroy operation occurs.
- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Read operations occur during any refresh or planning operation when refresh is enabled.
- `update` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).

## Important Notes

- **Windows Specific**: This resource is specifically for managing Win32 Line of Business (LOB) applications on Windows devices.
- **App Package Format**: Win32 LOB apps are typically in .msi, .exe, or .appx format and are custom applications developed for the organization.
- **Content Upload**: The resource handles uploading the app content to Intune for distribution to target devices.
- **Assignment Required**: Apps must be assigned to user or device groups to be deployed through Intune.
- **Detection Rules**: Configure detection rules to determine if the app is successfully installed on target devices. Multiple detection rule types are supported:
  - **Registry**: Check registry keys and values
  - **File System**: Check for specific files or folders
  - **MSI Information**: Use MSI product codes and versions
  - **PowerShell Script**: Custom detection using PowerShell scripts
- **Requirement Rules**: Define system requirements that must be met before the app can be installed.
- **Installation Context**: Win32 LOB apps can be installed in user or system context depending on configuration.
- **Return Codes**: Configure custom return codes to handle different installation outcomes.
- **Install Experience**: Control the installation behavior, restart requirements, and user interaction.
- **MSI Information**: For MSI-based apps, specify product codes, versions, and upgrade codes for proper management.
- **Supersedence**: Win32 LOB apps support supersedence relationships to replace older versions.

## Import

Import is supported using the following syntax:

```shell
# {resource_id}
terraform import microsoft365_graph_beta_device_and_app_management_win32_lob_app.example 00000000-0000-0000-0000-000000000000
``` 