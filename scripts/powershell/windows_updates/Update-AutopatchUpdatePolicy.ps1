<#
.SYNOPSIS
    Test script for updating Windows Autopatch Update Policies via Microsoft Graph API

.DESCRIPTION
    This script allows testing different JSON payloads for updating Update Policies
    to determine the correct API schema through trial and error.

.PARAMETER TenantId
    Specify the Entra ID tenant ID (Directory ID)

.PARAMETER ClientId
    Specify the application (client) ID of the Entra ID app registration

.PARAMETER ClientSecret
    Specify the client secret of the Entra ID app registration

.PARAMETER PolicyId
    The ID of the existing update policy to update

.PARAMETER JsonPayloadPath
    Path to JSON file containing the update policy payload

.PARAMETER CreateFirst
    If specified, creates a new policy first before updating it

.PARAMETER AudienceId
    Required if CreateFirst is specified - the deployment audience ID

.PARAMETER TestName
    Optional name for the test run (used in logging)

.EXAMPLE
    .\Update-AutopatchUpdatePolicy.ps1 -TenantId "xxx" -ClientId "xxx" -ClientSecret "xxx" -PolicyId "xxx" -JsonPayloadPath ".\update_payload.json" -TestName "UpdateTest1"

.EXAMPLE
    .\Update-AutopatchUpdatePolicy.ps1 -TenantId "xxx" -ClientId "xxx" -ClientSecret "xxx" -CreateFirst -AudienceId "xxx" -JsonPayloadPath ".\update_payload.json" -TestName "UpdateTest2"

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
    HelpMessage="The ID of the existing update policy to update")]
    [string]$PolicyId,
    
    [Parameter(Mandatory=$true,
    HelpMessage="Path to JSON file containing the update policy payload")]
    [ValidateNotNullOrEmpty()]
    [string]$JsonPayloadPath,
    
    [Parameter(Mandatory=$false,
    HelpMessage="Create a new policy first before updating it")]
    [switch]$CreateFirst,
    
    [Parameter(Mandatory=$false,
    HelpMessage="Required if CreateFirst is specified - the deployment audience ID")]
    [string]$AudienceId,
    
    [Parameter(Mandatory=$false,
    HelpMessage="Optional name for the test run")]
    [string]$TestName = "UnnamedTest"
)

# Set error action preference
$ErrorActionPreference = "Stop"

# Validate parameters
if ($CreateFirst -and [string]::IsNullOrEmpty($AudienceId)) {
    throw "AudienceId is required when CreateFirst is specified"
}

if (-not $CreateFirst -and [string]::IsNullOrEmpty($PolicyId)) {
    throw "PolicyId is required when CreateFirst is not specified"
}

# Initialize results tracking
$script:TestResults = @{
    TestName = $TestName
    Timestamp = Get-Date -Format "yyyy-MM-dd HH:mm:ss"
    JsonPayloadPath = $JsonPayloadPath
    CreateFirst = $CreateFirst.IsPresent
    InitialPolicyId = $PolicyId
    Success = $false
    PolicyId = $PolicyId
    RequestPayload = $null
    ResponseData = $null
    ErrorDetails = $null
    CreatedPolicy = $null
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
        client_id     = $ClientId
        scope         = "https://graph.microsoft.com/.default"
        client_secret = $ClientSecret
        grant_type    = "client_credentials"
    }
    
    try {
        $response = Invoke-RestMethod -Method Post -Uri "https://login.microsoftonline.com/$TenantId/oauth2/v2.0/token" -Body $body -ContentType "application/x-www-form-urlencoded"
        Write-TestLog "Successfully acquired access token" -Level Success
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
        [string]$AudienceId
    )
    
    Write-TestLog "Creating initial update policy..." -Level Info
    
    $createPayload = @{
        "@odata.type" = "#microsoft.graph.windowsUpdates.updatePolicy"
        audience = @{
            id = $AudienceId
        }
        complianceChanges = @(
            @{
                "@odata.type" = "#microsoft.graph.windowsUpdates.contentApproval"
            }
        )
    } | ConvertTo-Json -Depth 10
    
    $headers = @{
        "Authorization" = "Bearer $AccessToken"
        "Content-Type" = "application/json"
    }
    
    try {
        $response = Invoke-RestMethod -Method Post -Uri "https://graph.microsoft.com/beta/admin/windows/updates/updatePolicies" -Headers $headers -Body $createPayload
        Write-TestLog "Successfully created policy with ID: $($response.id)" -Level Success
        return $response
    }
    catch {
        Write-TestLog "Failed to create policy: $_" -Level Error
        throw
    }
}

function Update-Policy {
    param(
        [string]$AccessToken,
        [string]$PolicyId,
        [string]$JsonPayload
    )
    
    Write-TestLog "Updating policy $PolicyId..." -Level Info
    Write-TestLog "Request Payload:" -Level Info
    Write-Host $JsonPayload -ForegroundColor Cyan
    
    $script:TestResults.RequestPayload = $JsonPayload
    
    $headers = @{
        "Authorization" = "Bearer $AccessToken"
        "Content-Type" = "application/json"
    }
    
    try {
        $response = Invoke-RestMethod -Method Patch -Uri "https://graph.microsoft.com/beta/admin/windows/updates/updatePolicies/$PolicyId" -Headers $headers -Body $JsonPayload
        Write-TestLog "Successfully updated policy" -Level Success
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
            
            Write-TestLog "Update failed" -Level Error
            Write-TestLog "Error details:" -Level Error
            Write-Host $responseBody -ForegroundColor Red
            
            $script:TestResults.ErrorDetails = $responseBody
        }
        else {
            Write-TestLog "Update failed: $errorDetails" -Level Error
            $script:TestResults.ErrorDetails = $errorDetails
        }
        
        throw $errorDetails
    }
}

function Get-Policy {
    param(
        [string]$AccessToken,
        [string]$PolicyId
    )
    
    Write-TestLog "Reading policy $PolicyId..." -Level Info
    
    $headers = @{
        "Authorization" = "Bearer $AccessToken"
    }
    
    try {
        $response = Invoke-RestMethod -Method Get -Uri "https://graph.microsoft.com/beta/admin/windows/updates/updatePolicies/$PolicyId" -Headers $headers
        Write-TestLog "Successfully read policy" -Level Success
        return $response
    }
    catch {
        Write-TestLog "Failed to read policy: $_" -Level Error
        throw
    }
}

function Remove-Policy {
    param(
        [string]$AccessToken,
        [string]$PolicyId
    )
    
    Write-TestLog "Deleting policy $PolicyId..." -Level Info
    
    $headers = @{
        "Authorization" = "Bearer $AccessToken"
    }
    
    try {
        Invoke-RestMethod -Method Delete -Uri "https://graph.microsoft.com/beta/admin/windows/updates/updatePolicies/$PolicyId" -Headers $headers
        Write-TestLog "Successfully deleted policy" -Level Success
    }
    catch {
        Write-TestLog "Failed to delete policy: $_" -Level Warning
    }
}

#endregion

#region Main Script

try {
    Write-TestLog "========================================" -Level Info
    Write-TestLog "Starting Update Policy Test: $TestName" -Level Info
    Write-TestLog "========================================" -Level Info
    
    # Validate JSON file exists
    if (-not (Test-Path $JsonPayloadPath)) {
        throw "JSON payload file not found: $JsonPayloadPath"
    }
    
    # Read JSON payload
    Write-TestLog "Reading JSON payload from: $JsonPayloadPath" -Level Info
    $jsonContent = Get-Content -Path $JsonPayloadPath -Raw
    
    # Get access token
    $accessToken = Get-AccessToken -TenantId $TenantId -ClientId $ClientId -ClientSecret $ClientSecret
    
    # Create policy if requested
    if ($CreateFirst) {
        Write-TestLog "CreateFirst flag is set, creating initial policy..." -Level Info
        $createdPolicy = New-UpdatePolicy -AccessToken $accessToken -AudienceId $AudienceId
        $PolicyId = $createdPolicy.id
        $script:TestResults.PolicyId = $PolicyId
        $script:TestResults.CreatedPolicy = $createdPolicy
        
        # Wait a moment for policy to be fully created
        Write-TestLog "Waiting 2 seconds for policy creation to complete..." -Level Info
        Start-Sleep -Seconds 2
    }
    
    # Update the policy
    $updateResponse = Update-Policy -AccessToken $accessToken -PolicyId $PolicyId -JsonPayload $jsonContent
    $script:TestResults.ResponseData = $updateResponse
    $script:TestResults.Success = $true
    
    # Read back the policy to verify
    Write-TestLog "Reading back policy to verify update..." -Level Info
    $verifyResponse = Get-Policy -AccessToken $accessToken -PolicyId $PolicyId
    
    Write-TestLog "Updated Policy Data:" -Level Success
    Write-Host ($verifyResponse | ConvertTo-Json -Depth 10) -ForegroundColor Green
    
    # Cleanup if we created the policy
    if ($CreateFirst) {
        Write-TestLog "Cleaning up created policy..." -Level Info
        Remove-Policy -AccessToken $accessToken -PolicyId $PolicyId
    }
    
    Write-TestLog "========================================" -Level Success
    Write-TestLog "Test completed successfully!" -Level Success
    Write-TestLog "========================================" -Level Success
}
catch {
    Write-TestLog "========================================" -Level Error
    Write-TestLog "Test failed: $_" -Level Error
    Write-TestLog "========================================" -Level Error
    $script:TestResults.Success = $false
    
    # Cleanup if we created a policy
    if ($CreateFirst -and $script:TestResults.PolicyId) {
        Write-TestLog "Attempting cleanup of created policy..." -Level Warning
        try {
            $accessToken = Get-AccessToken -TenantId $TenantId -ClientId $ClientId -ClientSecret $ClientSecret
            Remove-Policy -AccessToken $accessToken -PolicyId $script:TestResults.PolicyId
        }
        catch {
            Write-TestLog "Cleanup failed: $_" -Level Warning
        }
    }
    
    exit 1
}
finally {
    # Output results summary
    Write-Host "`n========================================" -ForegroundColor Cyan
    Write-Host "Test Results Summary" -ForegroundColor Cyan
    Write-Host "========================================" -ForegroundColor Cyan
    Write-Host "Test Name: $($script:TestResults.TestName)"
    Write-Host "Timestamp: $($script:TestResults.Timestamp)"
    Write-Host "JSON Payload: $($script:TestResults.JsonPayloadPath)"
    Write-Host "Policy ID: $($script:TestResults.PolicyId)"
    Write-Host "Success: $($script:TestResults.Success)"
    
    if ($script:TestResults.ErrorDetails) {
        Write-Host "`nError Details:" -ForegroundColor Red
        Write-Host "Status Code: $($script:TestResults.ErrorDetails.StatusCode)"
        Write-Host "Status Description: $($script:TestResults.ErrorDetails.StatusDescription)"
    }
}

#endregion
