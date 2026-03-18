<#
.SYNOPSIS
    Delete a Windows Autopatch Deployment Audience

.DESCRIPTION
    Deletes a deployment audience by ID

.PARAMETER TenantId
    Specify the Entra ID tenant ID (Directory ID)

.PARAMETER ClientId
    Specify the application (client) ID of the Entra ID app registration

.PARAMETER ClientSecret
    Specify the client secret of the Entra ID app registration

.PARAMETER AudienceId
    The ID of the audience to delete

.EXAMPLE
    .\Delete-DeploymentAudience.ps1 -TenantId "xxx" -ClientId "xxx" -ClientSecret "xxx" -AudienceId "xxx"

#>

[CmdletBinding()]
param (
    [Parameter(Mandatory=$true)]
    [string]$TenantId,

    [Parameter(Mandatory=$true)]
    [string]$ClientId,
    
    [Parameter(Mandatory=$true)]
    [string]$ClientSecret,
    
    [Parameter(Mandatory=$true)]
    [string]$AudienceId
)

$ErrorActionPreference = "Stop"

function Get-AccessToken {
    param(
        [string]$TenantId,
        [string]$ClientId,
        [string]$ClientSecret
    )
    
    $body = @{
        client_id     = $ClientId
        scope         = "https://graph.microsoft.com/.default"
        client_secret = $ClientSecret
        grant_type    = "client_credentials"
    }
    
    $response = Invoke-RestMethod -Method Post -Uri "https://login.microsoftonline.com/$TenantId/oauth2/v2.0/token" -Body $body -ContentType "application/x-www-form-urlencoded"
    return $response.access_token
}

try {
    Write-Host "Acquiring access token..." -ForegroundColor Yellow
    $accessToken = Get-AccessToken -TenantId $TenantId -ClientId $ClientId -ClientSecret $ClientSecret
    
    Write-Host "Deleting deployment audience: $AudienceId" -ForegroundColor Yellow
    
    $headers = @{
        "Authorization" = "Bearer $accessToken"
    }
    
    Invoke-RestMethod -Method Delete -Uri "https://graph.microsoft.com/beta/admin/windows/updates/deploymentAudiences/$AudienceId" -Headers $headers
    
    Write-Host "Successfully deleted deployment audience" -ForegroundColor Green
}
catch {
    Write-Host "Failed to delete deployment audience: $_" -ForegroundColor Red
    throw
}
