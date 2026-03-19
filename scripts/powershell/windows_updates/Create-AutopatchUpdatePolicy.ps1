<#
.SYNOPSIS
    Test script for creating Windows Autopatch Update Policies via Microsoft Graph API

.DESCRIPTION
    This script allows testing different JSON payloads for creating Update Policies
    to determine the correct API schema through trial and error.

.PARAMETER TenantId
    Specify the Entra ID tenant ID (Directory ID)

.PARAMETER ClientId
    Specify the application (client) ID of the Entra ID app registration

.PARAMETER ClientSecret
    Specify the client secret of the Entra ID app registration

.PARAMETER JsonPayloadPath
    Path to JSON file containing the update policy payload

.PARAMETER DeleteOnSuccess
    If specified, deletes the created policy after successful creation

.PARAMETER TestName
    Optional name for the test run (used in logging)

.EXAMPLE
    .\Create-AutopatchUpdatePolicy.ps1 -TenantId "xxx" -ClientId "xxx" -ClientSecret "xxx" -JsonPayloadPath ".\payload1.json" -DeleteOnSuccess -TestName "Test1"

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
    
    [Parameter(Mandatory=$true,
    HelpMessage="Path to JSON file containing the update policy payload")]
    [ValidateNotNullOrEmpty()]
    [string]$JsonPayloadPath,
    
    [Parameter(Mandatory=$false,
    HelpMessage="Delete the policy after successful creation")]
    [switch]$DeleteOnSuccess,
    
    [Parameter(Mandatory=$false,
    HelpMessage="Optional name for the test run")]
    [string]$TestName = "UnnamedTest"
)

# Set error action preference
$ErrorActionPreference = "Stop"

# Initialize results tracking
$script:TestResults = @{
    TestName = $TestName
    Timestamp = Get-Date -Format "yyyy-MM-dd HH:mm:ss"
    JsonPayloadPath = $JsonPayloadPath
    Success = $false
    PolicyId = $null
    RequestPayload = $null
    ResponseData = $null
    ErrorDetails = $null
    Deleted = $false
}

#region Helper Functions

function Write-TestLog {
    param(
        [string]$Message,
        [ValidateSet('Info', 'Success', 'Warning', 'Error')]
        [string]$Level = 'Info'
    )
    
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
    
    Write-TestLog "Acquiring access token..." -Level Info
    
    $body = @{
        grant_type    = "client_credentials"
        client_id     = $ClientId
        client_secret = $ClientSecret
        scope         = "https://graph.microsoft.com/.default"
    }
    
    try {
        $response = Invoke-RestMethod -Method Post -Uri "https://login.microsoftonline.com/$TenantId/oauth2/v2.0/token" -Body $body -ContentType "application/x-www-form-urlencoded"
        Write-TestLog "Access token acquired successfully" -Level Success
        return $response.access_token
    }
    catch {
        Write-TestLog "Failed to acquire access token: $_" -Level Error
        throw
    }
}

function New-UpdatePolicy {
    param(
        [string]$AccessToken,
        [string]$JsonPayload
    )
    
    Write-TestLog "Creating Update Policy..." -Level Info
    
    $headers = @{
        "Authorization" = "Bearer $AccessToken"
        "Content-Type"  = "application/json"
    }
    
    $uri = "https://graph.microsoft.com/beta/admin/windows/updates/updatePolicies"
    
    try {
        Write-TestLog "Request URI: $uri" -Level Info
        Write-TestLog "Request Payload:" -Level Info
        Write-Host $JsonPayload -ForegroundColor Cyan
        
        $response = Invoke-RestMethod -Method Post -Uri $uri -Headers $headers -Body $JsonPayload
        
        Write-TestLog "Update Policy created successfully!" -Level Success
        Write-TestLog "Policy ID: $($response.id)" -Level Success
        
        return $response
    }
    catch {
        $errorDetails = $_.Exception.Message
        
        if ($_.Exception.Response) {
            $reader = New-Object System.IO.StreamReader($_.Exception.Response.GetResponseStream())
            $reader.BaseStream.Position = 0
            $reader.DiscardBufferedData()
            $responseBody = $reader.ReadToEnd()
            $errorDetails += "`nResponse Body: $responseBody"
        }
        
        Write-TestLog "Failed to create Update Policy" -Level Error
        Write-TestLog $errorDetails -Level Error
        
        throw $errorDetails
    }
}

function Remove-UpdatePolicy {
    param(
        [string]$AccessToken,
        [string]$PolicyId
    )
    
    Write-TestLog "Deleting Update Policy: $PolicyId" -Level Info
    
    $headers = @{
        "Authorization" = "Bearer $AccessToken"
    }
    
    $uri = "https://graph.microsoft.com/beta/admin/windows/updates/updatePolicies/$PolicyId"
    
    try {
        Invoke-RestMethod -Method Delete -Uri $uri -Headers $headers
        Write-TestLog "Update Policy deleted successfully" -Level Success
        return $true
    }
    catch {
        Write-TestLog "Failed to delete Update Policy: $_" -Level Error
        return $false
    }
}

function Get-UpdatePolicy {
    param(
        [string]$AccessToken,
        [string]$PolicyId
    )
    
    Write-TestLog "Retrieving Update Policy: $PolicyId" -Level Info
    
    $headers = @{
        "Authorization" = "Bearer $AccessToken"
    }
    
    $uri = "https://graph.microsoft.com/beta/admin/windows/updates/updatePolicies/$PolicyId"
    
    try {
        $response = Invoke-RestMethod -Method Get -Uri $uri -Headers $headers
        Write-TestLog "Update Policy retrieved successfully" -Level Success
        return $response
    }
    catch {
        Write-TestLog "Failed to retrieve Update Policy: $_" -Level Error
        throw
    }
}

function Save-TestResults {
    param(
        [hashtable]$Results
    )
    
    $resultsDir = Join-Path $PSScriptRoot "test_results"
    if (-not (Test-Path $resultsDir)) {
        New-Item -ItemType Directory -Path $resultsDir | Out-Null
    }
    
    $timestamp = Get-Date -Format "yyyyMMdd_HHmmss"
    $resultsFile = Join-Path $resultsDir "test_${TestName}_${timestamp}.json"
    
    $Results | ConvertTo-Json -Depth 10 | Out-File -FilePath $resultsFile -Encoding UTF8
    
    Write-TestLog "Test results saved to: $resultsFile" -Level Info
}

#endregion

#region Main Script

try {
    Write-TestLog "========================================" -Level Info
    Write-TestLog "Update Policy API Test: $TestName" -Level Info
    Write-TestLog "========================================" -Level Info
    
    # Validate JSON payload file exists
    if (-not (Test-Path $JsonPayloadPath)) {
        throw "JSON payload file not found: $JsonPayloadPath"
    }
    
    # Read JSON payload
    Write-TestLog "Reading JSON payload from: $JsonPayloadPath" -Level Info
    $jsonContent = Get-Content -Path $JsonPayloadPath -Raw
    $script:TestResults.RequestPayload = $jsonContent | ConvertFrom-Json
    
    # Get access token
    $accessToken = Get-AccessToken -TenantId $TenantId -ClientId $ClientId -ClientSecret $ClientSecret
    
    # Create Update Policy
    $policy = New-UpdatePolicy -AccessToken $accessToken -JsonPayload $jsonContent
    $script:TestResults.Success = $true
    $script:TestResults.PolicyId = $policy.id
    $script:TestResults.ResponseData = $policy
    
    # Retrieve the created policy to verify
    Write-TestLog "Verifying created policy..." -Level Info
    $retrievedPolicy = Get-UpdatePolicy -AccessToken $accessToken -PolicyId $policy.id
    
    Write-TestLog "Retrieved Policy Details:" -Level Info
    Write-Host ($retrievedPolicy | ConvertTo-Json -Depth 10) -ForegroundColor Cyan
    
    # Delete if requested
    if ($DeleteOnSuccess) {
        Write-TestLog "DeleteOnSuccess flag set, cleaning up..." -Level Info
        $deleted = Remove-UpdatePolicy -AccessToken $accessToken -PolicyId $policy.id
        $script:TestResults.Deleted = $deleted
    }
    
    Write-TestLog "========================================" -Level Success
    Write-TestLog "TEST PASSED: $TestName" -Level Success
    Write-TestLog "========================================" -Level Success
}
catch {
    $script:TestResults.Success = $false
    $script:TestResults.ErrorDetails = $_.Exception.Message
    
    Write-TestLog "========================================" -Level Error
    Write-TestLog "TEST FAILED: $TestName" -Level Error
    Write-TestLog "Error: $($_.Exception.Message)" -Level Error
    Write-TestLog "========================================" -Level Error
}
finally {
    # Save test results
    Save-TestResults -Results $script:TestResults
    
    # Exit with appropriate code
    if ($script:TestResults.Success) {
        exit 0
    } else {
        exit 1
    }
}

#endregion
