<#
.SYNOPSIS
    Updates the tenant default app management policy.

.DESCRIPTION
    Updates the tenant-wide default app management policy to enforce app management restrictions.
    This policy applies to all applications and service principals unless overridden by a specific
    appManagementPolicy.
    
    The policy ID is always 00000000-0000-0000-0000-000000000000 for the default policy.
    
    Note: The PATCH API returns 204 No Content with no response body on success.

.PARAMETER TenantId
    The Entra ID tenant ID.

.PARAMETER ClientId
    The application (client) ID for authentication.

.PARAMETER ClientSecret
    The client secret for authentication.

.PARAMETER IsEnabled
    Whether to enable the policy. Default is false.

.PARAMETER DisplayName
    The display name of the policy.

.PARAMETER Description
    The description of the policy.

.PARAMETER ConfigFile
    Path to a JSON file containing the full policy configuration.
    If provided, other parameters are ignored except authentication.

.PARAMETER ExportToJson
    Whether to export the results to a JSON file after update.

.EXAMPLE
    # Enable policy with simple password lifetime restriction
    .\Update-AppManagementPolicy.ps1 `
        -TenantId "your-tenant-id" `
        -ClientId "your-client-id" `
        -ClientSecret "your-secret" `
        -IsEnabled $true

.EXAMPLE
    # Update policy from JSON configuration file
    .\Update-AppManagementPolicy.ps1 `
        -TenantId "your-tenant-id" `
        -ClientId "your-client-id" `
        -ClientSecret "your-secret" `
        -ConfigFile "./policy-config.json" `
        -ExportToJson $true

.EXAMPLE
    # Disable the policy
    .\Update-AppManagementPolicy.ps1 `
        -TenantId "your-tenant-id" `
        -ClientId "your-client-id" `
        -ClientSecret "your-secret" `
        -IsEnabled $false

.NOTES
    Author: Deployment Theory
    Requires: Microsoft.Graph.Authentication module
    API Reference: https://learn.microsoft.com/en-us/graph/api/tenantappmanagementpolicy-update?view=graph-rest-beta
    
    Example JSON configuration file format:
    {
        "isEnabled": true,
        "displayName": "Custom Policy",
        "description": "My custom policy",
        "applicationRestrictions": {
            "passwordCredentials": [
                {
                    "restrictionType": "passwordLifetime",
                    "maxLifetime": "P90D",
                    "restrictForAppsCreatedAfterDateTime": "2024-01-01T00:00:00Z"
                }
            ],
            "keyCredentials": [
                {
                    "restrictionType": "asymmetricKeyLifetime",
                    "maxLifetime": "P365D",
                    "restrictForAppsCreatedAfterDateTime": "2024-01-01T00:00:00Z"
                }
            ]
        },
        "servicePrincipalRestrictions": {
            "passwordCredentials": [],
            "keyCredentials": []
        }
    }
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
    [bool]$IsEnabled,
    
    [Parameter(Mandatory=$false)]
    [string]$DisplayName,
    
    [Parameter(Mandatory=$false)]
    [string]$Description,
    
    [Parameter(Mandatory=$false)]
    [string]$ConfigFile,
    
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
        $fileName = "TenantAppManagementPolicy_Update_${timestamp}.json"
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
    Write-Host "   ✏️  Update Tenant App Management Policy" -ForegroundColor Magenta
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
    
    $requestBody = @{}
    
    if ($ConfigFile) {
        # Load configuration from JSON file
        Write-Host "   Loading configuration from file: $ConfigFile" -ForegroundColor Gray
        
        if (-not (Test-Path -Path $ConfigFile)) {
            throw "Configuration file not found: $ConfigFile"
        }
        
        $configContent = Get-Content -Path $ConfigFile -Raw | ConvertFrom-Json
        $requestBody = $configContent | ConvertTo-Json -Depth 10 | ConvertFrom-Json -AsHashtable
        
        Write-Host "   ✅ Configuration loaded from file" -ForegroundColor Green
    }
    else {
        # Build configuration from parameters
        if ($PSBoundParameters.ContainsKey('IsEnabled')) {
            $requestBody['isEnabled'] = $IsEnabled
            Write-Host "   • Is Enabled: $IsEnabled" -ForegroundColor Gray
        }
        
        if (-not [string]::IsNullOrEmpty($DisplayName)) {
            $requestBody['displayName'] = $DisplayName
            Write-Host "   • Display Name: $DisplayName" -ForegroundColor Gray
        }
        
        if (-not [string]::IsNullOrEmpty($Description)) {
            $requestBody['description'] = $Description
            Write-Host "   • Description: $Description" -ForegroundColor Gray
        }
    }
    
    if ($requestBody.Count -eq 0) {
        Write-Host "⚠️  No updates specified. Use -ConfigFile or specify -IsEnabled, -DisplayName, or -Description" -ForegroundColor Yellow
        Write-Host ""
        exit 0
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
    
    # Update the policy
    Write-Host "✏️  Updating policy..." -ForegroundColor Cyan
    
    $uri = "https://graph.microsoft.com/beta/policies/defaultAppManagementPolicy"
    $response = Invoke-MgGraphRequest -Method PATCH -Uri $uri -Body ($requestBody | ConvertTo-Json -Depth 10) -ContentType "application/json"
    
    Write-Host "✅ Policy updated successfully!" -ForegroundColor Green
    Write-Host "   Note: PATCH returns 204 No Content (no response body)" -ForegroundColor Gray
    Write-Host ""
    
    # Wait for eventual consistency
    Write-Host "⏱️  Waiting 5 seconds for eventual consistency..." -ForegroundColor Gray
    Start-Sleep -Seconds 5
    
    # Retrieve updated policy
    Write-Host "🔍 Retrieving updated policy..." -ForegroundColor Cyan
    
    $updatedPolicy = Invoke-MgGraphRequest -Method GET -Uri $uri
    
    Write-Host "✅ Retrieved updated policy" -ForegroundColor Green
    Write-Host ""
    
    # Display updated policy
    Write-Host "📋 Updated Policy" -ForegroundColor Cyan
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    Write-Host ""
    
    Write-Host "   Policy ID: $($updatedPolicy.id)" -ForegroundColor White
    Write-Host "   Display Name: $($updatedPolicy.displayName)" -ForegroundColor White
    Write-Host "   Description: $($updatedPolicy.description)" -ForegroundColor White
    Write-Host "   Is Enabled: " -NoNewline -ForegroundColor White
    Write-Host "$($updatedPolicy.isEnabled)" -ForegroundColor $(if ($updatedPolicy.isEnabled) { "Green" } else { "Yellow" })
    
    Write-Host ""
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    Write-Host ""
    
    # Display full response JSON
    Write-Host "📄 Full Updated Policy JSON" -ForegroundColor Cyan
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    Write-Host ""
    Write-Host ($updatedPolicy | ConvertTo-Json -Depth 10) -ForegroundColor Gray
    Write-Host ""
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    Write-Host ""
    
    # Export results if requested
    if ($ExportToJson) {
        $exportObject = @{
            Timestamp = Get-Date -Format "o"
            Operation = "Update"
            RequestBody = $requestBody
            UpdatedPolicy = $updatedPolicy
        }
        
        Export-ResultsToJson -Results $exportObject
    }
    
    Write-Host "🎉 Operation completed successfully!" -ForegroundColor Green
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
