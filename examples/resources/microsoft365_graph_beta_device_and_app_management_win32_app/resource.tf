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