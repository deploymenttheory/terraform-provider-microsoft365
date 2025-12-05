[CmdletBinding()]
param (
    [Parameter(Mandatory=$true,
    HelpMessage="Specify the Entra ID tenant ID (Directory ID) where the application is registered")]
    [ValidateNotNullOrEmpty()]
    [string]$TenantId,

    [Parameter(Mandatory=$true,
    HelpMessage="Specify the application (client) ID of the Entra ID app registration")]
    [ValidateNotNullOrEmpty()]
    [string]$ClientId,
    
    [Parameter(Mandatory=$true,
    HelpMessage="Specify the client secret of the Entra ID app registration")]
    [ValidateNotNullOrEmpty()]
    [string]$ClientSecret,

    [Parameter(Mandatory=$false,
    HelpMessage="Specify the desired state for the managed installer ('Enable' or 'Disable'). If not specified, the current state will be displayed.")]
    [ValidateSet('Enable', 'Disable')]
    [string]$Action,
    
    [Parameter(Mandatory=$false,
    HelpMessage="Show what would be changed without actually making changes")]
    [switch]$WhatIf
)

# Usage examples:
# Get current status:
# ./Patch-AppControlForBusinessManagedInstaller.ps1 -TenantId "xxx" -ClientId "xxx" -ClientSecret "xxx"

# Enable managed installer (What-If mode):
# ./Patch-AppControlForBusinessManagedInstaller.ps1 -TenantId "xxx" -ClientId "xxx" -ClientSecret "xxx" -Action "Enable" -WhatIf

# Enable managed installer:
# ./Patch-AppControlForBusinessManagedInstaller.ps1 -TenantId "xxx" -ClientId "xxx" -ClientSecret "xxx" -Action "Enable"

# Disable managed installer:
# ./Patch-AppControlForBusinessManagedInstaller.ps1 -TenantId "xxx" -ClientId "xxx" -ClientSecret "xxx" -Action "Disable"

Import-Module Microsoft.Graph.Authentication

function Get-WindowsManagementAppStatus {
    try {
        Write-Host "🔍 Getting Windows Management App status..." -ForegroundColor Cyan
        
        $appUri = "https://graph.microsoft.com/beta/deviceAppManagement/windowsManagementApp"
        Write-Host "   Endpoint: $appUri" -ForegroundColor Gray
        
        $response = Invoke-MgGraphRequest -Method GET -Uri $appUri
        Write-Host "   ✅ Windows Management App status retrieved successfully" -ForegroundColor Green
        Write-Host ""
        
        return $response
    }
    catch {
        Write-Host "❌ Error retrieving Windows Management App status: $_" -ForegroundColor Red
        Write-Host ""
        if ($_.Exception.Response) {
            $statusCode = $_.Exception.Response.StatusCode
            $statusDescription = $_.Exception.Response.StatusDescription
            Write-Host "   Status Code: $statusCode" -ForegroundColor Red
            Write-Host "   Status Description: $statusDescription" -ForegroundColor Red
            try {
                $responseContent = $_.Exception.Response.GetResponseStream()
                if ($responseContent) {
                    $reader = [System.IO.StreamReader]::new($responseContent)
                    $errorDetails = $reader.ReadToEnd()
                    $reader.Close()
                    Write-Host "   Error Details: $errorDetails" -ForegroundColor Red
                }
            } catch {
                Write-Host "   Could not read error details" -ForegroundColor Red
            }
        }
        throw
    }
}

function Set-ManagedInstallerStatus {
    param (
        [Parameter(Mandatory=$false)]
        [bool]$WhatIfPreference = $false
    )
    
    try {
        $setUri = "https://graph.microsoft.com/beta/deviceAppManagement/windowsManagementApp/setAsManagedInstaller"
        
        if ($WhatIfPreference) {
            Write-Host "🔍 WHAT-IF: Would toggle managed installer status" -ForegroundColor Yellow
            Write-Host "   Endpoint: $setUri" -ForegroundColor Gray
            Write-Host "   Method: POST" -ForegroundColor Gray
            Write-Host "   Body: Empty (toggle operation)" -ForegroundColor Gray
            return $true
        } else {
            Write-Host "🔄 Setting managed installer status..." -ForegroundColor Cyan
            Write-Host "   Endpoint: $setUri" -ForegroundColor Gray
            Write-Host "   Method: POST" -ForegroundColor Gray
            Write-Host "   Body: Empty (toggle operation)" -ForegroundColor Gray
            
            $response = Invoke-MgGraphRequest -Method POST -Uri $setUri -Body ""
            Write-Host "   ✅ Managed installer status updated successfully" -ForegroundColor Green
            return $true
        }
    }
    catch {
        Write-Host "❌ Error setting managed installer status: $_" -ForegroundColor Red
        if ($_.Exception.Response) {
            $statusCode = $_.Exception.Response.StatusCode
            $statusDescription = $_.Exception.Response.StatusDescription
            Write-Host "   Status Code: $statusCode" -ForegroundColor Red
            Write-Host "   Status Description: $statusDescription" -ForegroundColor Red
            try {
                $responseContent = $_.Exception.Response.GetResponseStream()
                if ($responseContent) {
                    $reader = [System.IO.StreamReader]::new($responseContent)
                    $errorDetails = $reader.ReadToEnd()
                    $reader.Close()
                    Write-Host "   Error Details: $errorDetails" -ForegroundColor Red
                }
            } catch {
                Write-Host "   Could not read error details" -ForegroundColor Red
            }
        }
        return $false
    }
}

function Show-WindowsManagementAppDetails {
    param (
        [Parameter(Mandatory=$true)]
        $WindowsManagementApp
    )
    
    Write-Host "📋 Windows Management App Details:" -ForegroundColor Cyan
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    
    # Basic app information
    Write-Host ("   • ID: {0}" -f $WindowsManagementApp.id) -ForegroundColor Green
    Write-Host ("   • Available Version: {0}" -f $WindowsManagementApp.availableVersion) -ForegroundColor Green
    Write-Host ("   • Managed Installer Status: {0}" -f $WindowsManagementApp.managedInstaller) -ForegroundColor Green
    
    # Managed installer configuration date/time
    if ($WindowsManagementApp.managedInstallerConfiguredDateTime) {
        Write-Host ("   • Managed Installer Configured: {0}" -f $WindowsManagementApp.managedInstallerConfiguredDateTime) -ForegroundColor Green
    } else {
        Write-Host "   • Managed Installer Configured: Never" -ForegroundColor Green
    }
    
    # Health states if available
    if ($WindowsManagementApp.healthStates -and $WindowsManagementApp.healthStates.Count -gt 0) {
        Write-Host ("   • Health States: {0} state(s) available" -f $WindowsManagementApp.healthStates.Count) -ForegroundColor Green
    } else {
        Write-Host "   • Health States: None available" -ForegroundColor Green
    }
    
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    Write-Host ""
}

function Show-ManagedInstallerInfo {
    Write-Host "ℹ️  App Control for Business - Managed Installer Information:" -ForegroundColor Cyan
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    Write-Host ""
    Write-Host "The Intune Management Extension as a Managed Installer feature allows:" -ForegroundColor White
    Write-Host "• The Intune Management Extension to be trusted for app installations" -ForegroundColor White
    Write-Host "• Applications deployed through Intune to bypass App Control for Business policies" -ForegroundColor White
    Write-Host "• Streamlined application deployment in secured environments" -ForegroundColor White
    Write-Host ""
    Write-Host "Status Values:" -ForegroundColor White
    Write-Host "• 'enabled'  - Intune Management Extension is configured as a managed installer" -ForegroundColor Green
    Write-Host "• 'disabled' - Intune Management Extension is not configured as a managed installer" -ForegroundColor Yellow
    Write-Host ""
    Write-Host "Learn more:" -ForegroundColor White
    Write-Host "https://learn.microsoft.com/en-us/windows/security/application-security/application-control/app-control-for-business/design/configure-authorized-apps-deployed-with-a-managed-installer" -ForegroundColor Blue
    Write-Host ""
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    Write-Host ""
}

function Show-ActionSummary {
    param (
        [Parameter(Mandatory=$true)]
        [string]$RequestedAction,
        [Parameter(Mandatory=$true)]
        [string]$CurrentStatus,
        [Parameter(Mandatory=$true)]
        [bool]$ActionTaken,
        [Parameter(Mandatory=$false)]
        [bool]$WhatIfPreference = $false
    )
    
    Write-Host "📊 Action Summary:" -ForegroundColor Cyan
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    Write-Host ("   • Requested Action: {0}" -f $RequestedAction) -ForegroundColor Green
    Write-Host ("   • Current Status: {0}" -f $CurrentStatus) -ForegroundColor Green
    
    if ($ActionTaken) {
        if ($WhatIfPreference) {
            Write-Host "   • Action Result: Would toggle managed installer status" -ForegroundColor Yellow
        } else {
            $newStatus = if ($CurrentStatus -eq "enabled") { "disabled" } else { "enabled" }
            Write-Host ("   • Action Result: Status toggled to '{0}'" -f $newStatus) -ForegroundColor Green
        }
    } else {
        Write-Host "   • Action Result: No action needed - already in desired state" -ForegroundColor Yellow
    }
    
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    Write-Host ""
}

try {
    Write-Host "🔐 Connecting to Microsoft Graph..." -ForegroundColor Cyan
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    
    $secureClientSecret = ConvertTo-SecureString -String $ClientSecret -AsPlainText -Force
    $clientSecretCredential = New-Object -TypeName System.Management.Automation.PSCredential -ArgumentList $ClientId, $secureClientSecret
    
    Connect-MgGraph -ClientSecretCredential $clientSecretCredential -TenantId $TenantId
    Write-Host "✅ Connected to Microsoft Graph successfully" -ForegroundColor Green
    Write-Host ""

    # Show information about the feature
    Show-ManagedInstallerInfo
    
    # Get current status
    Write-Host "🔍 Retrieving Current Status" -ForegroundColor Cyan
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    
    $windowsManagementApp = Get-WindowsManagementAppStatus
    
    if ($null -ne $windowsManagementApp) {
        Show-WindowsManagementAppDetails -WindowsManagementApp $windowsManagementApp
        
        $currentStatus = $windowsManagementApp.managedInstaller
        Write-Host ("Current managed installer status: {0}" -f $currentStatus) -ForegroundColor $(if ($currentStatus -eq "enabled") { "Green" } else { "Yellow" })
        Write-Host ""
        
        if ($Action) {
            Write-Host "🎯 Requested Action Processing" -ForegroundColor Cyan
            Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
            
            $desiredState = if ($Action -eq "Enable") { "enabled" } else { "disabled" }
            $actionNeeded = $currentStatus -ne $desiredState
            
            Write-Host ("Requested action: {0}" -f $Action) -ForegroundColor Cyan
            Write-Host ("Desired state: {0}" -f $desiredState) -ForegroundColor Cyan
            Write-Host ("Action needed: {0}" -f $(if ($actionNeeded) { "Yes" } else { "No" })) -ForegroundColor $(if ($actionNeeded) { "Yellow" } else { "Green" })
            Write-Host ""
            
            if ($actionNeeded) {
                if ($WhatIf) {
                    Write-Host "🔍 WHAT-IF MODE: The following action would be performed:" -ForegroundColor Yellow
                    Write-Host ("   Toggle managed installer from '{0}' to '{1}'" -f $currentStatus, $desiredState) -ForegroundColor Yellow
                    
                    $success = Set-ManagedInstallerStatus -WhatIfPreference $true
                    Show-ActionSummary -RequestedAction $Action -CurrentStatus $currentStatus -ActionTaken $success -WhatIfPreference $true
                } else {
                    $userConfirmation = Read-Host "❓ Are you sure you want to $($Action.ToLower()) the managed installer? (y/N)"
                    if ($userConfirmation -eq 'y' -or $userConfirmation -eq 'Y') {
                        $success = Set-ManagedInstallerStatus -WhatIfPreference $false
                        if ($success) {
                            Write-Host ""
                            Write-Host "🔄 Retrieving updated status..." -ForegroundColor Cyan
                            Start-Sleep -Seconds 2  # Brief pause to allow status to update
                            
                            $updatedApp = Get-WindowsManagementAppStatus
                            if ($null -ne $updatedApp) {
                                Show-WindowsManagementAppDetails -WindowsManagementApp $updatedApp
                                Show-ActionSummary -RequestedAction $Action -CurrentStatus $currentStatus -ActionTaken $true -WhatIfPreference $false
                                Write-Host "🎉 Managed installer status updated successfully!" -ForegroundColor Green
                            }
                        }
                    } else {
                        Write-Host "❌ Operation cancelled by user" -ForegroundColor Yellow
                    }
                }
            } else {
                Show-ActionSummary -RequestedAction $Action -CurrentStatus $currentStatus -ActionTaken $false -WhatIfPreference $false
                Write-Host "ℹ️  No action performed - managed installer is already in the desired state" -ForegroundColor Blue
            }
        } else {
            Write-Host "ℹ️  No action specified. Use -Action parameter to enable or disable the managed installer." -ForegroundColor Blue
            Write-Host "   Example: -Action Enable   (to enable managed installer)" -ForegroundColor Gray
            Write-Host "   Example: -Action Disable  (to disable managed installer)" -ForegroundColor Gray
        }
    } else {
        Write-Host "❌ Could not retrieve Windows Management App status" -ForegroundColor Red
    }
}
catch {
    Write-Host "❌ Script execution failed: $_" -ForegroundColor Red
    exit 1
}
finally {
    Write-Host "🔌 Disconnecting from Microsoft Graph..." -ForegroundColor Cyan
    try {
        Disconnect-MgGraph 2>$null
        Write-Host "✅ Disconnected from Microsoft Graph" -ForegroundColor Green
    } catch {}
}