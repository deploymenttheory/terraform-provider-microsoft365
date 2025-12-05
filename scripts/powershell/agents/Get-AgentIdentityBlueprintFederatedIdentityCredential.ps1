<#
.SYNOPSIS
    Tests various API endpoints for retrieving federated identity credentials.

.DESCRIPTION
    Performs GET requests to multiple Microsoft Graph API endpoints for federated identity
    credentials on both regular applications and Agent Identity Blueprints. This script
    is useful for testing and comparing API behavior across different endpoint patterns.

    Endpoints tested for regular applications:
    - GET /applications/{id}/federatedIdentityCredentials/{credentialId}
    - GET /applications/{id}/federatedIdentityCredentials/{credentialName}
    - GET /applications(appId='{appId}')/federatedIdentityCredentials/{credentialId}
    - GET /applications(appId='{appId}')/federatedIdentityCredentials/{credentialName}

    Endpoints tested for Agent Identity Blueprints:
    - GET /applications/{id}/microsoft.graph.agentIdentityBlueprint/federatedIdentityCredentials/{credentialId}
    - GET /applications/{id}/microsoft.graph.agentIdentityBlueprint/federatedIdentityCredentials/{credentialName}

.PARAMETER TenantId
    The Entra ID tenant ID.

.PARAMETER ClientId
    The application (client) ID for authentication.

.PARAMETER ClientSecret
    The client secret for authentication.

.PARAMETER ApplicationObjectId
    The Object ID (id) of the application or Agent Identity Blueprint.

.PARAMETER ApplicationAppId
    Optional. The Application (Client) ID (appId) for alternative addressing.

.PARAMETER CredentialId
    Optional. The ID of the federated identity credential to retrieve.

.PARAMETER CredentialName
    Optional. The name of the federated identity credential to retrieve.

.PARAMETER UseSelect
    Optional. Apply $select query parameter to customize response.

.PARAMETER SelectFields
    Optional. Comma-separated list of fields to select. Default: "id,name,issuer,subject,audiences".

.PARAMETER ExportToJson
    Whether to export the results to a JSON file.

.EXAMPLE
    # Test all endpoints with credential ID
    .\Get-AgentIdentityBlueprintFederatedIdentityCredential.ps1 `
        -TenantId "your-tenant-id" `
        -ClientId "your-auth-client-id" `
        -ClientSecret "your-auth-secret" `
        -ApplicationObjectId "target-app-object-id" `
        -CredentialId "your-credential-id" `
        -ExportToJson $true

.EXAMPLE
    # Test all endpoints with credential name
    .\Get-AgentIdentityBlueprintFederatedIdentityCredential.ps1 `
        -TenantId "your-tenant-id" `
        -ClientId "your-auth-client-id" `
        -ClientSecret "your-auth-secret" `
        -ApplicationObjectId "target-app-object-id" `
        -CredentialName "my-federated-credential" `
        -ExportToJson $true

.EXAMPLE
    # Test with both ID and AppId addressing
    .\Get-AgentIdentityBlueprintFederatedIdentityCredential.ps1 `
        -TenantId "your-tenant-id" `
        -ClientId "your-client-id" `
        -ClientSecret "your-secret" `
        -ApplicationObjectId "your-app-object-id" `
        -ApplicationAppId "your-app-client-id" `
        -CredentialId "your-credential-id"

.NOTES
    Author: Deployment Theory
    Requires: Microsoft.Graph.Authentication module
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
    
    [Parameter(Mandatory=$true)]
    [ValidateNotNullOrEmpty()]
    [string]$ApplicationObjectId,
    
    [Parameter(Mandatory=$false)]
    [string]$ApplicationAppId,
    
    [Parameter(Mandatory=$false)]
    [string]$CredentialId,
    
    [Parameter(Mandatory=$false)]
    [string]$CredentialName,
    
    [Parameter(Mandatory=$false)]
    [bool]$UseSelect = $false,
    
    [Parameter(Mandatory=$false)]
    [string]$SelectFields = "id,name,issuer,subject,audiences",
    
    [Parameter(Mandatory=$false)]
    [bool]$ExportToJson = $false
)

Import-Module Microsoft.Graph.Authentication

# Function to make a Graph API request and capture result
function Invoke-GraphApiTest {
    param (
        [Parameter(Mandatory=$true)]
        [string]$EndpointName,
        
        [Parameter(Mandatory=$true)]
        [string]$Uri,
        
        [Parameter(Mandatory=$false)]
        [string]$ApiVersion = "beta"
    )
    
    $fullUri = "https://graph.microsoft.com/$ApiVersion$Uri"
    
    Write-Host ""
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    Write-Host "🧪 Test: $EndpointName" -ForegroundColor Yellow
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    Write-Host "   URI: $fullUri" -ForegroundColor Gray
    Write-Host ""
    
    $result = @{
        EndpointName = $EndpointName
        Uri = $fullUri
        ApiVersion = $ApiVersion
        Success = $false
        StatusCode = $null
        ErrorMessage = $null
        Response = $null
        Timestamp = Get-Date -Format "o"
    }
    
    try {
        $response = Invoke-MgGraphRequest -Method GET -Uri $fullUri -ErrorAction Stop
        
        $result.Success = $true
        $result.StatusCode = 200
        $result.Response = $response
        
        Write-Host "   ✅ SUCCESS (200)" -ForegroundColor Green
        Write-Host ""
        Write-Host "   Response:" -ForegroundColor Cyan
        $response | ConvertTo-Json -Depth 5 | ForEach-Object {
            Write-Host "   $_" -ForegroundColor White
        }
    }
    catch {
        $errorMessage = $_.Exception.Message
        $statusCode = "Unknown"
        
        # Try to extract status code from error
        if ($errorMessage -match "(\d{3})") {
            $statusCode = $Matches[1]
        }
        
        $result.Success = $false
        $result.StatusCode = $statusCode
        $result.ErrorMessage = $errorMessage
        
        Write-Host "   ❌ FAILED ($statusCode)" -ForegroundColor Red
        Write-Host "   Error: $errorMessage" -ForegroundColor Red
    }
    
    return $result
}

# Function to export results to JSON
function Export-ResultsToJson {
    param (
        [Parameter(Mandatory=$true)]
        $Results,
        
        [Parameter(Mandatory=$true)]
        [string]$ApplicationObjectId
    )
    
    try {
        # Create output directory if it doesn't exist
        $outputDir = Join-Path -Path (Get-Location) -ChildPath "output"
        if (-not (Test-Path -Path $outputDir)) {
            New-Item -Path $outputDir -ItemType Directory | Out-Null
            Write-Host "📁 Created output directory: $outputDir" -ForegroundColor Gray
        }
        
        # Generate timestamp for filename
        $timestamp = Get-Date -Format "yyyyMMdd_HHmmss"
        $fileName = "FederatedIdentityCredential_Tests_${ApplicationObjectId}_${timestamp}.json"
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

# Main Script Execution
try {
    Write-Host ""
    Write-Host "═══════════════════════════════════════════════════════════════════" -ForegroundColor Magenta
    Write-Host "   🔍 Federated Identity Credential API Endpoint Tester" -ForegroundColor Magenta
    Write-Host "═══════════════════════════════════════════════════════════════════" -ForegroundColor Magenta
    Write-Host ""
    
    # Validate parameters
    if (-not $CredentialId -and -not $CredentialName) {
        Write-Host "⚠️  Warning: Neither CredentialId nor CredentialName provided." -ForegroundColor Yellow
        Write-Host "   Will attempt to list all federated identity credentials first." -ForegroundColor Yellow
        Write-Host ""
    }
    
    Write-Host "🔐 Connecting to Microsoft Graph..." -ForegroundColor Cyan
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    
    $secureClientSecret = ConvertTo-SecureString -String $ClientSecret -AsPlainText -Force
    $clientSecretCredential = New-Object -TypeName System.Management.Automation.PSCredential -ArgumentList $ClientId, $secureClientSecret
    
    Connect-MgGraph -ClientSecretCredential $clientSecretCredential -TenantId $TenantId
    
    Write-Host "✅ Connected to Microsoft Graph successfully" -ForegroundColor Green
    Write-Host ""
    
    # Display test parameters
    Write-Host "📋 Test Parameters" -ForegroundColor Cyan
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    Write-Host "   Application Object ID: $ApplicationObjectId" -ForegroundColor Green
    if ($ApplicationAppId) {
        Write-Host "   Application App ID:    $ApplicationAppId" -ForegroundColor Green
    }
    if ($CredentialId) {
        Write-Host "   Credential ID:         $CredentialId" -ForegroundColor Green
    }
    if ($CredentialName) {
        Write-Host "   Credential Name:       $CredentialName" -ForegroundColor Green
    }
    if ($UseSelect) {
        Write-Host "   Select Fields:         $SelectFields" -ForegroundColor Green
    }
    Write-Host ""
    
    # Build select query parameter
    $selectParam = ""
    if ($UseSelect) {
        $selectParam = "?`$select=$SelectFields"
    }
    
    # Store all test results
    $allResults = @()
    
    # ═══════════════════════════════════════════════════════════════════
    # TEST 1: List all federated identity credentials (standard endpoint)
    # ═══════════════════════════════════════════════════════════════════
    Write-Host ""
    Write-Host "═══════════════════════════════════════════════════════════════════" -ForegroundColor Blue
    Write-Host "   📋 LISTING ALL FEDERATED IDENTITY CREDENTIALS" -ForegroundColor Blue
    Write-Host "═══════════════════════════════════════════════════════════════════" -ForegroundColor Blue
    
    $listUri = "/applications/$ApplicationObjectId/federatedIdentityCredentials$selectParam"
    $listResult = Invoke-GraphApiTest -EndpointName "List (Standard Endpoint)" -Uri $listUri
    $allResults += $listResult
    
    $listCastUri = "/applications/$ApplicationObjectId/microsoft.graph.agentIdentityBlueprint/federatedIdentityCredentials$selectParam"
    $listCastResult = Invoke-GraphApiTest -EndpointName "List (Agent Blueprint Cast Endpoint)" -Uri $listCastUri
    $allResults += $listCastResult
    
    # ═══════════════════════════════════════════════════════════════════
    # TEST GROUP: Standard Application Endpoints
    # ═══════════════════════════════════════════════════════════════════
    Write-Host ""
    Write-Host "═══════════════════════════════════════════════════════════════════" -ForegroundColor Blue
    Write-Host "   📦 STANDARD APPLICATION ENDPOINTS" -ForegroundColor Blue
    Write-Host "═══════════════════════════════════════════════════════════════════" -ForegroundColor Blue
    
    # Test 2: GET /applications/{id}/federatedIdentityCredentials/{credentialId}
    if ($CredentialId) {
        $uri = "/applications/$ApplicationObjectId/federatedIdentityCredentials/$CredentialId$selectParam"
        $result = Invoke-GraphApiTest -EndpointName "GET by Object ID + Credential ID" -Uri $uri
        $allResults += $result
    }
    
    # Test 3: GET /applications/{id}/federatedIdentityCredentials/{credentialName}
    if ($CredentialName) {
        $uri = "/applications/$ApplicationObjectId/federatedIdentityCredentials/$CredentialName$selectParam"
        $result = Invoke-GraphApiTest -EndpointName "GET by Object ID + Credential Name" -Uri $uri
        $allResults += $result
    }
    
    # Test 4: GET /applications(appId='{appId}')/federatedIdentityCredentials/{credentialId}
    if ($ApplicationAppId -and $CredentialId) {
        $uri = "/applications(appId='$ApplicationAppId')/federatedIdentityCredentials/$CredentialId$selectParam"
        $result = Invoke-GraphApiTest -EndpointName "GET by App ID + Credential ID" -Uri $uri
        $allResults += $result
    }
    
    # Test 5: GET /applications(appId='{appId}')/federatedIdentityCredentials/{credentialName}
    if ($ApplicationAppId -and $CredentialName) {
        $uri = "/applications(appId='$ApplicationAppId')/federatedIdentityCredentials/$CredentialName$selectParam"
        $result = Invoke-GraphApiTest -EndpointName "GET by App ID + Credential Name" -Uri $uri
        $allResults += $result
    }
    
    # ═══════════════════════════════════════════════════════════════════
    # TEST GROUP: Agent Identity Blueprint Endpoints (Cast)
    # ═══════════════════════════════════════════════════════════════════
    Write-Host ""
    Write-Host "═══════════════════════════════════════════════════════════════════" -ForegroundColor Blue
    Write-Host "   🤖 AGENT IDENTITY BLUEPRINT ENDPOINTS (CAST)" -ForegroundColor Blue
    Write-Host "═══════════════════════════════════════════════════════════════════" -ForegroundColor Blue
    
    # Test 6: GET /applications/{id}/microsoft.graph.agentIdentityBlueprint/federatedIdentityCredentials/{credentialId}
    if ($CredentialId) {
        $uri = "/applications/$ApplicationObjectId/microsoft.graph.agentIdentityBlueprint/federatedIdentityCredentials/$CredentialId$selectParam"
        $result = Invoke-GraphApiTest -EndpointName "GET Agent Blueprint + Credential ID" -Uri $uri
        $allResults += $result
    }
    
    # Test 7: GET /applications/{id}/microsoft.graph.agentIdentityBlueprint/federatedIdentityCredentials/{credentialName}
    if ($CredentialName) {
        $uri = "/applications/$ApplicationObjectId/microsoft.graph.agentIdentityBlueprint/federatedIdentityCredentials/$CredentialName$selectParam"
        $result = Invoke-GraphApiTest -EndpointName "GET Agent Blueprint + Credential Name" -Uri $uri
        $allResults += $result
    }
    
    # ═══════════════════════════════════════════════════════════════════
    # SUMMARY
    # ═══════════════════════════════════════════════════════════════════
    Write-Host ""
    Write-Host "═══════════════════════════════════════════════════════════════════" -ForegroundColor Magenta
    Write-Host "   📊 TEST SUMMARY" -ForegroundColor Magenta
    Write-Host "═══════════════════════════════════════════════════════════════════" -ForegroundColor Magenta
    Write-Host ""
    
    $successCount = ($allResults | Where-Object { $_.Success }).Count
    $failedCount = ($allResults | Where-Object { -not $_.Success }).Count
    $totalCount = $allResults.Count
    
    Write-Host "   Total Tests:    $totalCount" -ForegroundColor White
    Write-Host "   ✅ Successful:  $successCount" -ForegroundColor Green
    Write-Host "   ❌ Failed:      $failedCount" -ForegroundColor Red
    Write-Host ""
    
    Write-Host "   Results by Endpoint:" -ForegroundColor Cyan
    Write-Host "   ─────────────────────────────────────────────────────────────" -ForegroundColor Cyan
    
    foreach ($result in $allResults) {
        $statusIcon = if ($result.Success) { "✅" } else { "❌" }
        $statusColor = if ($result.Success) { "Green" } else { "Red" }
        Write-Host "   $statusIcon $($result.EndpointName): $($result.StatusCode)" -ForegroundColor $statusColor
    }
    
    Write-Host ""
    
    # Export results if requested
    if ($ExportToJson) {
        $summaryObject = @{
            Timestamp = Get-Date -Format "o"
            TestParameters = @{
                ApplicationObjectId = $ApplicationObjectId
                ApplicationAppId = $ApplicationAppId
                CredentialId = $CredentialId
                CredentialName = $CredentialName
                UseSelect = $UseSelect
                SelectFields = $SelectFields
            }
            Summary = @{
                TotalTests = $totalCount
                Successful = $successCount
                Failed = $failedCount
            }
            Results = $allResults
        }
        
        Export-ResultsToJson -Results $summaryObject -ApplicationObjectId $ApplicationObjectId
    }
    
    Write-Host "🎉 Testing completed!" -ForegroundColor Green
    Write-Host ""
}
catch {
    Write-Host ""
    Write-Host "❌ Error: $_" -ForegroundColor Red
    Write-Host ""
    exit 1
}
finally {
    Write-Host "🔌 Disconnecting from Microsoft Graph..." -ForegroundColor Cyan
    Disconnect-MgGraph 2>$null
    Write-Host "✅ Disconnected" -ForegroundColor Green
    Write-Host ""
}

