# Delete-DirectorySettings.ps1
# Script to get and delete directory settings from Microsoft Graph API
# Based on: https://learn.microsoft.com/en-us/graph/api/directorysetting-delete?view=graph-rest-beta&tabs=http
#
# Example usage:
# pwsh ./Delete-DirectorySettings.ps1 \
#     -TenantId "00000000-0000-0000-0000-000000000000" \
#     -ClientId "11111111-1111-1111-1111-111111111111" \
#     -ClientSecret "22222222-2222-2222-2222-222222222222" \
#     -Force

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
    HelpMessage="Optional ID of a specific directory setting to delete")]
    [string]$SettingId,
    
    [Parameter(Mandatory=$false,
    HelpMessage="Skip confirmation prompt and delete without asking")]
    [switch]$Force,
    
    [Parameter(Mandatory=$false,
    HelpMessage="Path to output JSON file containing deleted settings information")]
    [string]$OutputFile
)

# Function to authenticate and get access token
function Connect-MicrosoftGraph {
    param (
        [Parameter(Mandatory=$true)]
        [string]$TenantId,
        
        [Parameter(Mandatory=$true)]
        [string]$ClientId,
        
        [Parameter(Mandatory=$true)]
        [string]$ClientSecret
    )
    
    try {
        Write-Host "Connecting to Microsoft Graph..." -ForegroundColor Cyan
        
        # Create secure credential
        $secureClientSecret = ConvertTo-SecureString -String $ClientSecret -AsPlainText -Force
        $clientSecretCredential = New-Object -TypeName System.Management.Automation.PSCredential -ArgumentList $ClientId, $secureClientSecret
        
        # Connect to Microsoft Graph
        Import-Module Microsoft.Graph.Authentication
        Connect-MgGraph -ClientSecretCredential $clientSecretCredential -TenantId $TenantId
        
        Write-Host "✅ Connected to Microsoft Graph" -ForegroundColor Green
    }
    catch {
        Write-Host "❌ Error connecting to Microsoft Graph: $_" -ForegroundColor Red
        throw
    }
}

# Function to get directory settings
function Get-DirectorySettings {
    param (
        [Parameter(Mandatory=$false)]
        [string]$SettingId
    )
    
    try {
        $baseUrl = "https://graph.microsoft.com/beta/settings"
        
        if ($SettingId) {
            $url = "$baseUrl/$SettingId"
            Write-Host "Retrieving directory setting: $SettingId..." -ForegroundColor Cyan
            $response = Invoke-MgGraphRequest -Method GET -Uri $url
            return @($response)
        } else {
            Write-Host "Retrieving all directory settings..." -ForegroundColor Cyan
            $response = Invoke-MgGraphRequest -Method GET -Uri $baseUrl
            
            if ($response.value) {
                return $response.value
            } else {
                return @()
            }
        }
    }
    catch {
        Write-Host "❌ Error retrieving directory settings: $_" -ForegroundColor Red
        throw
    }
}

# Function to display settings summary
function Show-SettingsSummary {
    param (
        [Parameter(Mandatory=$true)]
        [array]$Settings
    )
    
    Write-Host "`n📋 Directory Settings Found:" -ForegroundColor Cyan
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    
    if ($Settings.Count -eq 0) {
        Write-Host "   No directory settings found." -ForegroundColor Yellow
        return
    }
    
    foreach ($setting in $Settings) {
        Write-Host "`n📄 Setting: " -NoNewline -ForegroundColor Yellow
        Write-Host "$($setting.displayName)" -ForegroundColor White
        Write-Host "   ID: $($setting.id)" -ForegroundColor Gray
        Write-Host "   Template ID: $($setting.templateId)" -ForegroundColor Gray
        
        if ($setting.values) {
            Write-Host "   Values Count: $($setting.values.Count)" -ForegroundColor Gray
            Write-Host "   Current Values:" -ForegroundColor Gray
            foreach ($value in $setting.values) {
                Write-Host "     • $($value.name): $($value.value)" -ForegroundColor DarkGray
            }
        }
    }
}

# Function to delete directory settings
function Remove-DirectorySettings {
    param (
        [Parameter(Mandatory=$true)]
        [array]$Settings,
        
        [Parameter(Mandatory=$false)]
        [switch]$Force
    )
    
    $deletedSettings = @()
    $failedSettings = @()
    
    foreach ($setting in $Settings) {
        try {
            if (-not $Force) {
                $confirmation = Read-Host "`n⚠️  Delete setting '$($setting.displayName)' (ID: $($setting.id))? [Y/N]"
                if ($confirmation -ne 'Y' -and $confirmation -ne 'y') {
                    Write-Host "   ⏭️  Skipped: $($setting.displayName)" -ForegroundColor Yellow
                    continue
                }
            }
            
            $url = "https://graph.microsoft.com/beta/settings/$($setting.id)"
            Write-Host "   🗑️  Deleting: $($setting.displayName)..." -ForegroundColor Cyan
            
            Invoke-MgGraphRequest -Method DELETE -Uri $url
            
            Write-Host "   ✅ Deleted: $($setting.displayName)" -ForegroundColor Green
            $deletedSettings += $setting
        }
        catch {
            Write-Host "   ❌ Failed to delete: $($setting.displayName) - $_" -ForegroundColor Red
            $failedSettings += @{
                Setting = $setting
                Error = $_.Exception.Message
            }
        }
    }
    
    return @{
        Deleted = $deletedSettings
        Failed = $failedSettings
    }
}

# Function to save deletion results to file
function Save-DeletionResults {
    param (
        [Parameter(Mandatory=$true)]
        [object]$Results,
        
        [Parameter(Mandatory=$true)]
        [string]$OutputFile
    )
    
    try {
        $output = @{
            Timestamp = Get-Date -Format "o"
            DeletedCount = $Results.Deleted.Count
            FailedCount = $Results.Failed.Count
            DeletedSettings = $Results.Deleted
            FailedSettings = $Results.Failed
        }
        
        $prettyJson = ConvertTo-Json -InputObject $output -Depth 10 -Compress:$false
        $prettyJson | Out-File -FilePath $OutputFile -Encoding utf8
        
        Write-Host "   📁 Results saved to: $OutputFile" -ForegroundColor Green
    }
    catch {
        Write-Host "   ⚠️  Warning: Could not save results to file: $_" -ForegroundColor Yellow
    }
}

# Main script execution
try {
    # Connect to Microsoft Graph
    Connect-MicrosoftGraph -TenantId $TenantId -ClientId $ClientId -ClientSecret $ClientSecret
    
    # Get directory settings
    Write-Host "`n🔧 Retrieving directory settings..." -ForegroundColor Cyan
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    
    $settings = Get-DirectorySettings -SettingId $SettingId
    
    # Display summary
    Show-SettingsSummary -Settings $settings
    
    if ($settings.Count -eq 0) {
        Write-Host "`n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Yellow
        Write-Host "✨ No directory settings to delete" -ForegroundColor Yellow
        exit 0
    }
    
    # Confirm deletion
    Write-Host "`n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Yellow
    Write-Host "⚠️  WARNING: You are about to delete $($settings.Count) directory setting(s)" -ForegroundColor Yellow
    
    if (-not $Force) {
        Write-Host "`n💡 Tip: Use -Force parameter to skip individual confirmations" -ForegroundColor Cyan
    }
    
    # Delete settings
    Write-Host "`n🗑️  Starting deletion process..." -ForegroundColor Cyan
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    
    $results = Remove-DirectorySettings -Settings $settings -Force:$Force
    
    # Display results
    Write-Host "`n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Green
    Write-Host "📊 Deletion Summary:" -ForegroundColor Cyan
    Write-Host "   ✅ Successfully deleted: $($results.Deleted.Count)" -ForegroundColor Green
    Write-Host "   ❌ Failed to delete: $($results.Failed.Count)" -ForegroundColor Red
    
    # Save results to file if specified
    if ($OutputFile) {
        Write-Host "`n📁 Saving deletion results..." -ForegroundColor Cyan
        Save-DeletionResults -Results $results -OutputFile $OutputFile
    }
    
    # Final status
    if ($results.Failed.Count -eq 0) {
        Write-Host "`n✨ All directory settings deleted successfully!" -ForegroundColor Green
    } else {
        Write-Host "`n⚠️  Some deletions failed. Check the output above for details." -ForegroundColor Yellow
        if ($OutputFile) {
            Write-Host "   Full results saved to: $OutputFile" -ForegroundColor Gray
        }
    }
}
catch {
    Write-Host "`n❌ Script execution failed: $_" -ForegroundColor Red
    exit 1
}
finally {
    # Disconnect from Microsoft Graph
    Write-Host "`nDisconnecting from Microsoft Graph..." -ForegroundColor Cyan
    Disconnect-MgGraph > $null 2>&1
    Write-Host "✅ Disconnected from Microsoft Graph" -ForegroundColor Green
}


