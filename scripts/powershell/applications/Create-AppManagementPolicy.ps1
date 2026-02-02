<#
.SYNOPSIS
    Configures the tenant default app management policy (initial setup).

.DESCRIPTION
    Configures the tenant-wide default app management policy for the first time.
    Since the default policy always exists (ID: 00000000-0000-0000-0000-000000000000),
    this script uses PATCH to configure it with restrictions.
    
    This script is useful for initial policy setup and testing different restriction types.

.PARAMETER TenantId
    The Entra ID tenant ID.

.PARAMETER ClientId
    The application (client) ID for authentication.

.PARAMETER ClientSecret
    The client secret for authentication.

.PARAMETER RestrictPasswordAddition
    Block addition of new passwords/secrets to apps created after the specified date.

.PARAMETER RestrictPasswordLifetime
    Restrict maximum password lifetime to specified duration (e.g., "P90D" for 90 days).

.PARAMETER RestrictSymmetricKeyLifetime
    Restrict maximum symmetric key lifetime to specified duration (e.g., "P30D" for 30 days).

.PARAMETER RestrictKeyLifetime
    Restrict maximum asymmetric key/certificate lifetime to specified duration (e.g., "P365D" for 365 days).

.PARAMETER RestrictFromDate
    The date from which restrictions apply to newly created applications.
    Format: YYYY-MM-DDTHH:MM:SSZ (e.g., "2024-01-01T00:00:00Z")
    Default: Current date.

.PARAMETER ApplyToServicePrincipals
    If true, also apply the same restrictions to service principals.

.PARAMETER ExportToJson
    Whether to export the results to a JSON file.

.EXAMPLE
    # Block password addition for all new apps
    .\Create-AppManagementPolicy.ps1 `
        -TenantId "your-tenant-id" `
        -ClientId "your-client-id" `
        -ClientSecret "your-secret" `
        -RestrictPasswordAddition $true

.EXAMPLE
    # Enforce 90-day password lifetime for apps created from today
    .\Create-AppManagementPolicy.ps1 `
        -TenantId "your-tenant-id" `
        -ClientId "your-client-id" `
        -ClientSecret "your-secret" `
        -RestrictPasswordLifetime "P90D" `
        -RestrictFromDate "2024-01-01T00:00:00Z"

.EXAMPLE
    # Comprehensive policy with multiple restrictions
    .\Create-AppManagementPolicy.ps1 `
        -TenantId "your-tenant-id" `
        -ClientId "your-client-id" `
        -ClientSecret "your-secret" `
        -RestrictPasswordLifetime "P90D" `
        -RestrictSymmetricKeyLifetime "P30D" `
        -RestrictKeyLifetime "P365D" `
        -ApplyToServicePrincipals $true `
        -ExportToJson $true

.NOTES
    Author: Deployment Theory
    Requires: Microsoft.Graph.Authentication module
    API Reference: https://learn.microsoft.com/en-us/graph/api/tenantappmanagementpolicy-update?view=graph-rest-beta
    
    Valid restriction types:
    - passwordAddition (blocks new passwords, maxLifetime must be null)
    - passwordLifetime (restricts lifetime, maxLifetime required)
    - symmetricKeyAddition (blocks new symmetric keys, maxLifetime must be null)
    - symmetricKeyLifetime (restricts lifetime, maxLifetime required)
    - customPasswordAddition (blocks custom passwords, maxLifetime must be null)
    - asymmetricKeyLifetime (restricts cert lifetime, maxLifetime required)
    - trustedCertificateAuthority (restricts to specific CAs)
#>

[CmdletBinding()]
param (
    [Parameter(Mandatory=$true)]
    [ValidateNotNullOrEmpty()]
    [string]$TenantId,

    [Parameter(Mandatory=$true)]
    [ValidateNotNullOrEmpty()]
    [string]$ClientId,
    
    [Parameter(Mandatory=$true)]
    [ValidateNotNullOrEmpty()]
    [string]$ClientSecret,
    
    [Parameter(Mandatory=$false)]
    [bool]$RestrictPasswordAddition = $false,
    
    [Parameter(Mandatory=$false)]
    [string]$RestrictPasswordLifetime,
    
    [Parameter(Mandatory=$false)]
    [string]$RestrictSymmetricKeyLifetime,
    
    [Parameter(Mandatory=$false)]
    [string]$RestrictKeyLifetime,
    
    [Parameter(Mandatory=$false)]
    [string]$RestrictFromDate,
    
    [Parameter(Mandatory=$false)]
    [bool]$ApplyToServicePrincipals = $false,
    
    [Parameter(Mandatory=$false)]
    [bool]$ExportToJson = $false
)

Import-Module Microsoft.Graph.Authentication

#region Helper Functions

function Export-ResultsToJson {
    param (
        [Parameter(Mandatory=$true)]
        $Results
    )
    
    try {
        $outputDir = Join-Path -Path (Get-Location) -ChildPath "output"
        if (-not (Test-Path -Path $outputDir)) {
            New-Item -Path $outputDir -ItemType Directory | Out-Null
            Write-Host "📁 Created output directory: $outputDir" -ForegroundColor Gray
        }
        
        $timestamp = Get-Date -Format "yyyyMMdd_HHmmss"
        $fileName = "TenantAppManagementPolicy_Create_${timestamp}.json"
        $filePath = Join-Path -Path $outputDir -ChildPath $fileName
        
        $Results | ConvertTo-Json -Depth 10 | Out-File -FilePath $filePath -Encoding UTF8
        
        Write-Host "💾 Exported results to: $filePath" -ForegroundColor Green
        return $filePath
    }
    catch {
        Write-Host "❌ Error exporting results to JSON: $_" -ForegroundColor Red
        return $null
    }
}

#endregion Helper Functions

#region Main Script Execution

try {
    Write-Host ""
    Write-Host "═══════════════════════════════════════════════════════════════════" -ForegroundColor Magenta
    Write-Host "   ➕ Create/Configure Tenant App Management Policy" -ForegroundColor Magenta
    Write-Host "═══════════════════════════════════════════════════════════════════" -ForegroundColor Magenta
    Write-Host ""
    
    # Connect to Microsoft Graph
    Write-Host "🔐 Connecting to Microsoft Graph..." -ForegroundColor Cyan
    
    $secureSecret = ConvertTo-SecureString -String $ClientSecret -AsPlainText -Force
    $credential = New-Object -TypeName System.Management.Automation.PSCredential -ArgumentList $ClientId, $secureSecret
    
    Connect-MgGraph -ClientSecretCredential $credential -TenantId $TenantId -NoWelcome
    
    Write-Host "✅ Connected" -ForegroundColor Green
    Write-Host ""
    
    # Build request body
    Write-Host "📝 Building policy configuration..." -ForegroundColor Cyan
    
    # Set default restrict from date if not provided
    if ([string]::IsNullOrEmpty($RestrictFromDate)) {
        $RestrictFromDate = (Get-Date).ToUniversalTime().ToString("yyyy-MM-ddTHH:mm:ssZ")
        Write-Host "   Using current date for restrictions: $RestrictFromDate" -ForegroundColor Gray
    }
    
    $requestBody = @{
        isEnabled = $true
    }
    
    $passwordCreds = @()
    $keyCreds = @()
    
    # Add password restriction configurations
    if ($RestrictPasswordAddition) {
        $passwordCreds += @{
            restrictionType = "passwordAddition"
            maxLifetime = $null
            restrictForAppsCreatedAfterDateTime = $RestrictFromDate
        }
        Write-Host "   • Blocking password addition from: $RestrictFromDate" -ForegroundColor Gray
    }
    
    if (-not [string]::IsNullOrEmpty($RestrictPasswordLifetime)) {
        $passwordCreds += @{
            restrictionType = "passwordLifetime"
            maxLifetime = $RestrictPasswordLifetime
            restrictForAppsCreatedAfterDateTime = $RestrictFromDate
        }
        Write-Host "   • Restricting password lifetime to: $RestrictPasswordLifetime" -ForegroundColor Gray
    }
    
    if (-not [string]::IsNullOrEmpty($RestrictSymmetricKeyLifetime)) {
        $passwordCreds += @{
            restrictionType = "symmetricKeyLifetime"
            maxLifetime = $RestrictSymmetricKeyLifetime
            restrictForAppsCreatedAfterDateTime = $RestrictFromDate
        }
        Write-Host "   • Restricting symmetric key lifetime to: $RestrictSymmetricKeyLifetime" -ForegroundColor Gray
    }
    
    # Add key restriction configurations
    if (-not [string]::IsNullOrEmpty($RestrictKeyLifetime)) {
        $keyCreds += @{
            restrictionType = "asymmetricKeyLifetime"
            maxLifetime = $RestrictKeyLifetime
            restrictForAppsCreatedAfterDateTime = $RestrictFromDate
        }
        Write-Host "   • Restricting asymmetric key lifetime to: $RestrictKeyLifetime" -ForegroundColor Gray
    }
    
    # Build application restrictions
    if ($passwordCreds.Count -gt 0 -or $keyCreds.Count -gt 0) {
        $requestBody['applicationRestrictions'] = @{}
        
        if ($passwordCreds.Count -gt 0) {
            $requestBody['applicationRestrictions']['passwordCredentials'] = $passwordCreds
        }
        
        if ($keyCreds.Count -gt 0) {
            $requestBody['applicationRestrictions']['keyCredentials'] = $keyCreds
        }
    }
    
    # Apply same restrictions to service principals if requested
    if ($ApplyToServicePrincipals -and ($passwordCreds.Count -gt 0 -or $keyCreds.Count -gt 0)) {
        $requestBody['servicePrincipalRestrictions'] = @{}
        
        if ($passwordCreds.Count -gt 0) {
            $requestBody['servicePrincipalRestrictions']['passwordCredentials'] = $passwordCreds
        }
        
        if ($keyCreds.Count -gt 0) {
            $requestBody['servicePrincipalRestrictions']['keyCredentials'] = $keyCreds
        }
        
        Write-Host "   • Applying same restrictions to service principals" -ForegroundColor Gray
    }
    
    Write-Host ""
    
    # Display request body
    Write-Host "📤 Request Body" -ForegroundColor Cyan
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    Write-Host ""
    Write-Host ($requestBody | ConvertTo-Json -Depth 10) -ForegroundColor Gray
    Write-Host ""
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    Write-Host ""
    
    # Create/Configure the policy
    Write-Host "➕ Configuring policy..." -ForegroundColor Cyan
    
    $uri = "https://graph.microsoft.com/beta/policies/defaultAppManagementPolicy"
    $response = Invoke-MgGraphRequest -Method PATCH -Uri $uri -Body ($requestBody | ConvertTo-Json -Depth 10) -ContentType "application/json"
    
    Write-Host "✅ Policy configured successfully!" -ForegroundColor Green
    Write-Host "   Note: PATCH returns 204 No Content (no response body)" -ForegroundColor Gray
    Write-Host ""
    
    # Wait for eventual consistency
    Write-Host "⏱️  Waiting 5 seconds for eventual consistency..." -ForegroundColor Gray
    Start-Sleep -Seconds 5
    
    # Retrieve configured policy
    Write-Host "🔍 Retrieving configured policy..." -ForegroundColor Cyan
    
    $configuredPolicy = Invoke-MgGraphRequest -Method GET -Uri $uri
    
    Write-Host "✅ Retrieved configured policy" -ForegroundColor Green
    Write-Host ""
    
    # Display configured policy
    Write-Host "📋 Configured Policy" -ForegroundColor Cyan
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    Write-Host ""
    
    Write-Host "   Policy ID: $($configuredPolicy.id)" -ForegroundColor White
    Write-Host "   Display Name: $($configuredPolicy.displayName)" -ForegroundColor White
    Write-Host "   Description: $($configuredPolicy.description)" -ForegroundColor White
    Write-Host "   Is Enabled: " -NoNewline -ForegroundColor White
    Write-Host "$($configuredPolicy.isEnabled)" -ForegroundColor $(if ($configuredPolicy.isEnabled) { "Green" } else { "Yellow" })
    
    Write-Host ""
    
    # Display application restrictions summary
    $appRestrictions = $configuredPolicy.applicationRestrictions
    $pwdCount = if ($appRestrictions.passwordCredentials) { $appRestrictions.passwordCredentials.Count } else { 0 }
    $keyCount = if ($appRestrictions.keyCredentials) { $appRestrictions.keyCredentials.Count } else { 0 }
    
    Write-Host "   Application Restrictions:" -ForegroundColor White
    Write-Host "     • Password Credential Restrictions: $pwdCount" -ForegroundColor Gray
    Write-Host "     • Key Credential Restrictions: $keyCount" -ForegroundColor Gray
    
    Write-Host ""
    
    # Display service principal restrictions summary
    $spRestrictions = $configuredPolicy.servicePrincipalRestrictions
    $spPwdCount = if ($spRestrictions.passwordCredentials) { $spRestrictions.passwordCredentials.Count } else { 0 }
    $spKeyCount = if ($spRestrictions.keyCredentials) { $spRestrictions.keyCredentials.Count } else { 0 }
    
    Write-Host "   Service Principal Restrictions:" -ForegroundColor White
    Write-Host "     • Password Credential Restrictions: $spPwdCount" -ForegroundColor Gray
    Write-Host "     • Key Credential Restrictions: $spKeyCount" -ForegroundColor Gray
    
    Write-Host ""
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    Write-Host ""
    
    # Display full response JSON
    Write-Host "📄 Full Policy JSON" -ForegroundColor Cyan
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    Write-Host ""
    Write-Host ($configuredPolicy | ConvertTo-Json -Depth 10) -ForegroundColor Gray
    Write-Host ""
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    Write-Host ""
    
    # Export results if requested
    if ($ExportToJson) {
        $exportObject = @{
            Timestamp = Get-Date -Format "o"
            Operation = "Create/Configure"
            RequestBody = $requestBody
            ConfiguredPolicy = $configuredPolicy
        }
        
        Export-ResultsToJson -Results $exportObject
    }
    
    Write-Host "🎉 Operation completed successfully!" -ForegroundColor Green
    Write-Host ""
    Write-Host "💡 Policy ID is always: 00000000-0000-0000-0000-000000000000" -ForegroundColor Cyan
    Write-Host ""
}
catch {
    Write-Host ""
    Write-Host "❌ Error: $_" -ForegroundColor Red
    
    if ($_.Exception.Response) {
        $statusCode = $_.Exception.Response.StatusCode.value__
        Write-Host "   Status Code: $statusCode" -ForegroundColor Gray
        
        if ($_.ErrorDetails) {
            try {
                $errorDetails = $_.ErrorDetails.Message | ConvertFrom-Json
                Write-Host "   Error Code: $($errorDetails.error.code)" -ForegroundColor Gray
                Write-Host "   Error Message: $($errorDetails.error.message)" -ForegroundColor Gray
                
                if ($errorDetails.error.innerError) {
                    Write-Host "   Request ID: $($errorDetails.error.innerError.'request-id')" -ForegroundColor Gray
                    Write-Host "   Date: $($errorDetails.error.innerError.date)" -ForegroundColor Gray
                }
            }
            catch {
                Write-Host "   Error Details: $($_.ErrorDetails.Message)" -ForegroundColor Gray
            }
        }
    }
    
    Write-Host ""
    exit 1
}
finally {
    Write-Host "🔌 Disconnecting from Microsoft Graph..." -ForegroundColor Cyan
    Disconnect-MgGraph 2>$null
    Write-Host "✅ Disconnected" -ForegroundColor Green
    Write-Host ""
}

#endregion Main Script Execution
