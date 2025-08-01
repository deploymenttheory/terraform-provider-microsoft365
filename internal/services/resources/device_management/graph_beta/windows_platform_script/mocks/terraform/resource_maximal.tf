resource "microsoft365_graph_beta_device_management_windows_platform_script" "maximal" {
  display_name            = "Test Maximal Windows Platform Script - Unique"
  description             = "Maximal platform script for testing with all features"
  file_name               = "test_maximal.ps1"
  script_content          = "# PowerShell Script - Comprehensive test script with all features\n\nParam(\n    [string]$LogPath = \"C:\\Windows\\Temp\\script.log\"\n)\n\n# Start logging\nStart-Transcript -Path $LogPath -Append\n\ntry {\n    Write-Host \"Starting maximal test script execution\" -ForegroundColor Green\n    \n    # Check PowerShell version\n    Write-Host \"PowerShell Version: $($PSVersionTable.PSVersion)\" -ForegroundColor Cyan\n    \n    # System information\n    $OS = Get-CimInstance -ClassName Win32_OperatingSystem\n    Write-Host \"Operating System: $($OS.Caption) $($OS.Version)\" -ForegroundColor Cyan\n    \n    # Check disk space\n    $Disk = Get-CimInstance -ClassName Win32_LogicalDisk -Filter \"DriveType=3\"\n    foreach ($Drive in $Disk) {\n        $FreeSpaceGB = [math]::Round($Drive.FreeSpace / 1GB, 2)\n        $SizeGB = [math]::Round($Drive.Size / 1GB, 2)\n        Write-Host \"Drive $($Drive.DeviceID) - Free: $FreeSpaceGB GB / Total: $SizeGB GB\" -ForegroundColor Cyan\n    }\n    \n    # Network connectivity test\n    if (Test-Connection -ComputerName \"microsoft.com\" -Count 2 -Quiet) {\n        Write-Host \"Network connectivity: OK\" -ForegroundColor Green\n    } else {\n        Write-Host \"Network connectivity: Failed\" -ForegroundColor Yellow\n    }\n    \n    # Registry test\n    $TestPath = \"HKLM:\\SOFTWARE\\TestScript\"\n    New-Item -Path $TestPath -Force | Out-Null\n    Set-ItemProperty -Path $TestPath -Name \"TestValue\" -Value \"ScriptTest\"\n    $TestValue = Get-ItemProperty -Path $TestPath -Name \"TestValue\" -ErrorAction SilentlyContinue\n    \n    if ($TestValue.TestValue -eq \"ScriptTest\") {\n        Write-Host \"Registry test: OK\" -ForegroundColor Green\n        Remove-Item -Path $TestPath -Force\n    } else {\n        Write-Host \"Registry test: Failed\" -ForegroundColor Red\n    }\n    \n    Write-Host \"Maximal test script completed successfully\" -ForegroundColor Green\n    Exit 0\n}\ncatch {\n    Write-Error \"Script failed: $($_.Exception.Message)\"\n    Exit 1\n}\nfinally {\n    Stop-Transcript\n}"
  run_as_account          = "user"
  role_scope_tag_ids      = ["0", "1"]
  enforce_signature_check = true
  run_as_32_bit           = false

  assignments = [
    {
      type     = "groupAssignmentTarget"
      group_id = "44444444-4444-4444-4444-444444444444"
    }
  ]

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}