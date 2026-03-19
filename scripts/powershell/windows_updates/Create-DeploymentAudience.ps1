<#
.SYNOPSIS
    Creates a Windows Autopatch Deployment Audience for testing

.DESCRIPTION
    This script creates a deployment audience and returns its ID for use in Update Policy tests.
    Can optionally delete the audience after use.

.PARAMETER TenantId
    Specify the Entra ID tenant ID (Directory ID)

.PARAMETER ClientId
    Specify the application (client) ID of the Entra ID app registration

.PARAMETER ClientSecret
    Specify the client secret of the Entra ID app registration

.PARAMETER AudienceName
    Optional name for the deployment audience (default: "Test-Audience-{timestamp}")

.PARAMETER DeleteOnExit
    If specified, deletes the created audience when script exits

.PARAMETER OutputIdOnly
    If specified, only outputs the audience ID (useful for piping to other scripts)

.EXAMPLE
    .\Create-DeploymentAudience.ps1 -TenantId "xxx" -ClientId "xxx" -ClientSecret "xxx"

.EXAMPLE
    $audienceId = .\Create-DeploymentAudience.ps1 -TenantId "xxx" -ClientId "xxx" -ClientSecret "xxx" -OutputIdOnly

.EXAMPLE
    .\Create-DeploymentAudience.ps1 -TenantId "xxx" -ClientId "xxx" -ClientSecret "xxx" -DeleteOnExit

#>

[CmdletBinding()]
param (
    [Parameter(Mandatory=$true,
    HelpMessage="Specify the Entra ID tenant ID (Directory ID)")]
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
    HelpMessage="Name for the deployment audience")]
    [string]$AudienceName,
    
    [Parameter(Mandatory=$false,
    HelpMessage="Delete the audience when script exits")]
    [switch]$DeleteOnExit,
    
    [Parameter(Mandatory=$false,
    HelpMessage="Only output the audience ID")]
    [switch]$OutputIdOnly
)

# Set error action preference
$ErrorActionPreference = "Stop"

# Generate audience name if not provided
if (-not $AudienceName) {
    $timestamp = Get-Date -Format "yyyyMMdd-HHmmss"
    $AudienceName = "Test-Audience-$timestamp"
}

# Track created audience for cleanup
$script:CreatedAudienceId = $null

#region Helper Functions

function Write-Log {
    param(
        [string]$Message,
        [ValidateSet('Info', 'Success', 'Warning', 'Error')]
        [string]$Level = 'Info'
    )
    
    if ($OutputIdOnly) {
        return
    }
    
    $timestamp = Get-Date -Format "HH:mm:ss"
    $color = switch ($Level) {
        'Success' { 'Green' }
        'Warning' { 'Yellow' }
        'Error' { 'Red' }
        default { 'White' }
    }
    
    Write-Host "[$timestamp] [$Level] $Message" -ForegroundColor $color
}

function Get-AccessToken {
    param(
        [string]$TenantId,
        [string]$ClientId,
        [string]$ClientSecret
    )
    
    Write-Log "Acquiring access token..." -Level Info
    
    $body = @{
        grant_type    = "client_credentials"
        client_id     = $ClientId
        client_secret = $ClientSecret
        scope         = "https://graph.microsoft.com/.default"
    }
    
    try {
        $response = Invoke-RestMethod -Method Post -Uri "https://login.microsoftonline.com/$TenantId/oauth2/v2.0/token" -Body $body -ContentType "application/x-www-form-urlencoded"
        Write-Log "Access token acquired successfully" -Level Success
        return $response.access_token
    }
    catch {
        Write-Log "Failed to acquire access token: $_" -Level Error
        throw
    }
}

function New-DeploymentAudience {
    param(
        [string]$AccessToken,
        [string]$Name
    )
    
    Write-Log "Creating Deployment Audience: $Name" -Level Info
    
    $headers = @{
        "Authorization" = "Bearer $AccessToken"
        "Content-Type"  = "application/json"
    }
    
    $uri = "https://graph.microsoft.com/beta/admin/windows/updates/deploymentAudiences"
    
    # Minimal payload - just create an empty audience
    $payload = @{
        "@odata.type" = "#microsoft.graph.windowsUpdates.deploymentAudience"
    } | ConvertTo-Json
    
    try {
        Write-Log "Request URI: $uri" -Level Info
        Write-Log "Request Payload: $payload" -Level Info
        
        $response = Invoke-RestMethod -Method Post -Uri $uri -Headers $headers -Body $payload
        
        Write-Log "Deployment Audience created successfully!" -Level Success
        Write-Log "Audience ID: $($response.id)" -Level Success
        Write-Log "Audience @odata.context: $($response.'@odata.context')" -Level Info
        
        return $response
    }
    catch {
        $errorDetails = $_.Exception.Message
        
        if ($_.Exception.Response) {
            try {
                $reader = New-Object System.IO.StreamReader($_.Exception.Response.GetResponseStream())
                $reader.BaseStream.Position = 0
                $reader.DiscardBufferedData()
                $responseBody = $reader.ReadToEnd()
                $errorDetails += "`nResponse Body: $responseBody"
            }
            catch {
                # Ignore errors reading response body
            }
        }
        
        Write-Log "Failed to create Deployment Audience" -Level Error
        Write-Log $errorDetails -Level Error
        
        throw $errorDetails
    }
}

function Remove-DeploymentAudience {
    param(
        [string]$AccessToken,
        [string]$AudienceId
    )
    
    Write-Log "Deleting Deployment Audience: $AudienceId" -Level Info
    
    $headers = @{
        "Authorization" = "Bearer $AccessToken"
    }
    
    $uri = "https://graph.microsoft.com/beta/admin/windows/updates/deploymentAudiences/$AudienceId"
    
    try {
        Invoke-RestMethod -Method Delete -Uri $uri -Headers $headers
        Write-Log "Deployment Audience deleted successfully" -Level Success
        return $true
    }
    catch {
        Write-Log "Failed to delete Deployment Audience: $_" -Level Error
        return $false
    }
}

function Get-DeploymentAudience {
    param(
        [string]$AccessToken,
        [string]$AudienceId
    )
    
    Write-Log "Retrieving Deployment Audience: $AudienceId" -Level Info
    
    $headers = @{
        "Authorization" = "Bearer $AccessToken"
    }
    
    $uri = "https://graph.microsoft.com/beta/admin/windows/updates/deploymentAudiences/$AudienceId"
    
    try {
        $response = Invoke-RestMethod -Method Get -Uri $uri -Headers $headers
        Write-Log "Deployment Audience retrieved successfully" -Level Success
        return $response
    }
    catch {
        Write-Log "Failed to retrieve Deployment Audience: $_" -Level Error
        throw
    }
}

#endregion

#region Main Script

try {
    if (-not $OutputIdOnly) {
        Write-Log "========================================" -Level Info
        Write-Log "Creating Deployment Audience" -Level Info
        Write-Log "========================================" -Level Info
    }
    
    # Get access token
    $accessToken = Get-AccessToken -TenantId $TenantId -ClientId $ClientId -ClientSecret $ClientSecret
    
    # Create Deployment Audience
    $audience = New-DeploymentAudience -AccessToken $accessToken -Name $AudienceName
    $script:CreatedAudienceId = $audience.id
    
    # Verify the created audience
    Write-Log "Verifying created audience..." -Level Info
    $retrievedAudience = Get-DeploymentAudience -AccessToken $accessToken -AudienceId $audience.id
    
    if (-not $OutputIdOnly) {
        Write-Log "Retrieved Audience Details:" -Level Info
        Write-Host ($retrievedAudience | ConvertTo-Json -Depth 10) -ForegroundColor Cyan
        
        Write-Log "========================================" -Level Success
        Write-Log "Audience Created Successfully" -Level Success
        Write-Log "Audience ID: $($audience.id)" -Level Success
        Write-Log "========================================" -Level Success
    }
    
    # Output the ID (for piping to other scripts)
    if ($OutputIdOnly) {
        Write-Output $audience.id
    } else {
        Write-Host ""
        Write-Host "To use this audience ID in tests, run:" -ForegroundColor Yellow
        Write-Host "  `$audienceId = `"$($audience.id)`"" -ForegroundColor Cyan
        Write-Host ""
    }
    
    # If DeleteOnExit is not set, wait for user input before cleaning up
    if (-not $DeleteOnExit -and -not $OutputIdOnly) {
        Write-Host "Press Enter to delete the audience, or Ctrl+C to keep it..." -ForegroundColor Yellow
        Read-Host
        $DeleteOnExit = $true
    }
}
catch {
    Write-Log "========================================" -Level Error
    Write-Log "Failed to create Deployment Audience" -Level Error
    Write-Log "Error: $($_.Exception.Message)" -Level Error
    Write-Log "========================================" -Level Error
    
    exit 1
}
finally {
    # Clean up if requested
    if ($DeleteOnExit -and $script:CreatedAudienceId) {
        Write-Log "Cleaning up..." -Level Info
        $deleted = Remove-DeploymentAudience -AccessToken $accessToken -AudienceId $script:CreatedAudienceId
        
        if ($deleted) {
            Write-Log "Cleanup completed successfully" -Level Success
        } else {
            Write-Log "Cleanup failed - you may need to manually delete audience: $script:CreatedAudienceId" -Level Warning
        }
    } elseif ($script:CreatedAudienceId -and -not $OutputIdOnly) {
        Write-Log "Audience not deleted. To delete manually, run:" -Level Warning
        Write-Host "  Remove-MgBetaWindowsUpdatesDeploymentAudience -DeploymentAudienceId `"$script:CreatedAudienceId`"" -ForegroundColor Yellow
    }
}

#endregion
