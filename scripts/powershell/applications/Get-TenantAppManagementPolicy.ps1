<#
.SYNOPSIS
    Gets the tenant default app management policy with all attributes.

.DESCRIPTION
    Retrieves the tenant-wide default app management policy from Microsoft Entra ID.
    This policy enforces app management restrictions for all applications and service principals
    unless overridden by a specific appManagementPolicy.
    
    The policy ID is always 00000000-0000-0000-0000-000000000000 for the default policy.

.PARAMETER TenantId
    The Entra ID tenant ID.

.PARAMETER ClientId
    The application (client) ID for authentication.

.PARAMETER ClientSecret
    The client secret for authentication.

.PARAMETER ExportToJson
    Whether to export the results to a JSON file.

.EXAMPLE
    # Get the default policy
    .\Get-AppManagementPolicy.ps1 `
        -TenantId "your-tenant-id" `
        -ClientId "your-client-id" `
        -ClientSecret "your-secret"

.EXAMPLE
    # Get policy and export to JSON
    .\Get-AppManagementPolicy.ps1 `
        -TenantId "your-tenant-id" `
        -ClientId "your-client-id" `
        -ClientSecret "your-secret" `
        -ExportToJson $true

.NOTES
    Author: Deployment Theory
    Requires: Microsoft.Graph.Authentication module
    API Reference: https://learn.microsoft.com/en-us/graph/api/tenantappmanagementpolicy-get?view=graph-rest-beta
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
    [bool]$ExportToJson = $false
)

Import-Module Microsoft.Graph.Authentication

#region Helper Functions

function Format-AttributeValue {
    param (
        [Parameter(Mandatory=$false)]
        $Value
    )
    
    if ($null -eq $Value) {
        return "(null)"
    }
    elseif ($Value -is [array]) {
        if ($Value.Count -eq 0) {
            return "(empty array)"
        }
        return ($Value -join ", ")
    }
    elseif ($Value -is [hashtable] -or $Value -is [System.Collections.IDictionary]) {
        return "(complex object)"
    }
    elseif ($Value -is [bool]) {
        return $Value.ToString().ToLower()
    }
    elseif ($Value -is [string] -and [string]::IsNullOrEmpty($Value)) {
        return "(empty)"
    }
    else {
        return $Value.ToString()
    }
}

function Format-CredentialConfiguration {
    param (
        [Parameter(Mandatory=$true)]
        $Credentials,
        
        [Parameter(Mandatory=$true)]
        [string]$Type,
        
        [Parameter(Mandatory=$false)]
        [int]$Indent = 6
    )
    
    $indentStr = " " * $Indent
    $count = 1
    
    foreach ($cred in $Credentials) {
        Write-Host "$indentStr🔹 Restriction $count" -ForegroundColor Yellow
        Write-Host "$indentStr   Restriction Type: $($cred.restrictionType)" -ForegroundColor Gray
        Write-Host "$indentStr   State: $($cred.state)" -ForegroundColor $(if ($cred.state -eq "enabled") { "Green" } else { "Yellow" })
        Write-Host "$indentStr   Restrict From Date: $($cred.restrictForAppsCreatedAfterDateTime)" -ForegroundColor Gray
        
        if ($null -ne $cred.maxLifetime) {
            Write-Host "$indentStr   Max Lifetime: $($cred.maxLifetime)" -ForegroundColor Gray
        }
        
        if ($null -ne $cred.certificateBasedApplicationConfigurationIds -and $cred.certificateBasedApplicationConfigurationIds.Count -gt 0) {
            Write-Host "$indentStr   Cert Config IDs: $($cred.certificateBasedApplicationConfigurationIds.Count)" -ForegroundColor Gray
        }
        
        if ($null -ne $cred.excludeActors) {
            Write-Host "$indentStr   Exclude Actors: (configured)" -ForegroundColor Gray
        }
        
        Write-Host ""
        $count++
    }
}

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
        $fileName = "TenantAppManagementPolicy_${timestamp}.json"
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
    Write-Host "   🔍 Get Tenant App Management Policy" -ForegroundColor Magenta
    Write-Host "═══════════════════════════════════════════════════════════════════" -ForegroundColor Magenta
    Write-Host ""
    
    # Connect to Microsoft Graph
    Write-Host "🔐 Connecting to Microsoft Graph..." -ForegroundColor Cyan
    
    $secureSecret = ConvertTo-SecureString -String $ClientSecret -AsPlainText -Force
    $credential = New-Object -TypeName System.Management.Automation.PSCredential -ArgumentList $ClientId, $secureSecret
    
    Connect-MgGraph -ClientSecretCredential $credential -TenantId $TenantId -NoWelcome
    
    Write-Host "✅ Connected" -ForegroundColor Green
    Write-Host ""
    
    # Get the default policy
    Write-Host "🔍 Retrieving tenant app management policy..." -ForegroundColor Cyan
    
    $uri = "https://graph.microsoft.com/beta/policies/defaultAppManagementPolicy"
    $policy = Invoke-MgGraphRequest -Method GET -Uri $uri
    
    Write-Host "✅ Policy retrieved successfully!" -ForegroundColor Green
    Write-Host ""
    
    # Display policy overview
    Write-Host "📋 Policy Overview" -ForegroundColor Cyan
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    Write-Host ""
    
    Write-Host "   Policy ID: " -NoNewline -ForegroundColor White
    Write-Host "$($policy.id)" -ForegroundColor Green
    
    Write-Host "   Display Name: " -NoNewline -ForegroundColor White
    Write-Host "$($policy.displayName)" -ForegroundColor Green
    
    Write-Host "   Description: " -NoNewline -ForegroundColor White
    Write-Host "$($policy.description)" -ForegroundColor Green
    
    Write-Host "   Is Enabled: " -NoNewline -ForegroundColor White
    Write-Host "$($policy.isEnabled)" -ForegroundColor $(if ($policy.isEnabled) { "Green" } else { "Yellow" })
    
    Write-Host ""
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    Write-Host ""
    
    # Display Application Restrictions
    Write-Host "📱 Application Restrictions" -ForegroundColor Cyan
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    Write-Host ""
    
    $appRestrictions = $policy.applicationRestrictions
    
    if ($appRestrictions.passwordCredentials -and $appRestrictions.passwordCredentials.Count -gt 0) {
        Write-Host "   🔐 Password Credentials Restrictions ($($appRestrictions.passwordCredentials.Count))" -ForegroundColor White
        Write-Host ""
        Format-CredentialConfiguration -Credentials $appRestrictions.passwordCredentials -Type "Password"
    }
    else {
        Write-Host "   🔐 Password Credentials Restrictions: " -NoNewline -ForegroundColor White
        Write-Host "(none)" -ForegroundColor Gray
        Write-Host ""
    }
    
    if ($appRestrictions.keyCredentials -and $appRestrictions.keyCredentials.Count -gt 0) {
        Write-Host "   🔑 Key Credentials Restrictions ($($appRestrictions.keyCredentials.Count))" -ForegroundColor White
        Write-Host ""
        Format-CredentialConfiguration -Credentials $appRestrictions.keyCredentials -Type "Key"
    }
    else {
        Write-Host "   🔑 Key Credentials Restrictions: " -NoNewline -ForegroundColor White
        Write-Host "(none)" -ForegroundColor Gray
        Write-Host ""
    }
    
    if ($appRestrictions.identifierUris -and $appRestrictions.identifierUris.nonDefaultUriAddition) {
        Write-Host "   🔗 Identifier URIs Restrictions: (configured)" -ForegroundColor White
        Write-Host ""
    }
    else {
        Write-Host "   🔗 Identifier URIs Restrictions: " -NoNewline -ForegroundColor White
        Write-Host "(none)" -ForegroundColor Gray
        Write-Host ""
    }
    
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    Write-Host ""
    
    # Display Service Principal Restrictions
    Write-Host "⚙️  Service Principal Restrictions" -ForegroundColor Cyan
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    Write-Host ""
    
    $spRestrictions = $policy.servicePrincipalRestrictions
    
    if ($spRestrictions.passwordCredentials -and $spRestrictions.passwordCredentials.Count -gt 0) {
        Write-Host "   🔐 Password Credentials Restrictions ($($spRestrictions.passwordCredentials.Count))" -ForegroundColor White
        Write-Host ""
        Format-CredentialConfiguration -Credentials $spRestrictions.passwordCredentials -Type "Password"
    }
    else {
        Write-Host "   🔐 Password Credentials Restrictions: " -NoNewline -ForegroundColor White
        Write-Host "(none)" -ForegroundColor Gray
        Write-Host ""
    }
    
    if ($spRestrictions.keyCredentials -and $spRestrictions.keyCredentials.Count -gt 0) {
        Write-Host "   🔑 Key Credentials Restrictions ($($spRestrictions.keyCredentials.Count))" -ForegroundColor White
        Write-Host ""
        Format-CredentialConfiguration -Credentials $spRestrictions.keyCredentials -Type "Key"
    }
    else {
        Write-Host "   🔑 Key Credentials Restrictions: " -NoNewline -ForegroundColor White
        Write-Host "(none)" -ForegroundColor Gray
        Write-Host ""
    }
    
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    Write-Host ""
    
    # Display raw JSON
    Write-Host "📄 Raw JSON Response" -ForegroundColor Cyan
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    Write-Host ""
    Write-Host ($policy | ConvertTo-Json -Depth 10) -ForegroundColor Gray
    Write-Host ""
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    Write-Host ""
    
    # Export results if requested
    if ($ExportToJson) {
        $exportObject = @{
            Timestamp = Get-Date -Format "o"
            Operation = "Get"
            Policy = $policy
        }
        
        Export-ResultsToJson -Results $exportObject
    }
    
    Write-Host "🎉 Operation completed!" -ForegroundColor Green
    Write-Host ""
}
catch {
    Write-Host ""
    Write-Host "❌ Error: $_" -ForegroundColor Red
    
    if ($_.Exception.Response) {
        $statusCode = $_.Exception.Response.StatusCode.value__
        Write-Host "   Status Code: $statusCode" -ForegroundColor Gray
        
        if ($_.ErrorDetails) {
            $errorDetails = $_.ErrorDetails.Message | ConvertFrom-Json
            Write-Host "   Error Code: $($errorDetails.error.code)" -ForegroundColor Gray
            Write-Host "   Error Message: $($errorDetails.error.message)" -ForegroundColor Gray
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
