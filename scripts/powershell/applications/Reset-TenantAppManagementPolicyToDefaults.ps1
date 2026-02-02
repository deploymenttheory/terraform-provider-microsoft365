<#
.SYNOPSIS
    Resets the tenant default app management policy to Microsoft defaults.

.DESCRIPTION
    Completely resets the tenant-wide default app management policy to empty/default state.
    This is useful for cleaning up before running tests.

.PARAMETER TenantId
    The Entra ID tenant ID.

.PARAMETER ClientId
    The application (client) ID for authentication.

.PARAMETER ClientSecret
    The client secret for authentication.
#>

[CmdletBinding()]
param(
    [Parameter(Mandatory = $true)]
    [string]$TenantId,

    [Parameter(Mandatory = $true)]
    [string]$ClientId,

    [Parameter(Mandatory = $true)]
    [string]$ClientSecret
)

try {
    Write-Host ""
    Write-Host "═══════════════════════════════════════════════════════════════════" -ForegroundColor Cyan
    Write-Host "   🔄 Reset Tenant App Management Policy to Defaults" -ForegroundColor Cyan
    Write-Host "═══════════════════════════════════════════════════════════════════" -ForegroundColor Cyan
    Write-Host ""
    
    # Connect to Microsoft Graph
    Write-Host "🔐 Connecting to Microsoft Graph..." -ForegroundColor Cyan
    
    $secureSecret = ConvertTo-SecureString -String $ClientSecret -AsPlainText -Force
    $clientSecretCredential = New-Object -TypeName System.Management.Automation.PSCredential -ArgumentList $ClientId, $secureSecret
    
    Connect-MgGraph -TenantId $TenantId -ClientSecretCredential $clientSecretCredential -NoWelcome
    
    Write-Host "✅ Connected" -ForegroundColor Green
    Write-Host ""
    
    # Build default/empty policy
    Write-Host "📝 Building default policy configuration..." -ForegroundColor Cyan
    
    $requestBody = @{
        isEnabled = $true
        displayName = "Default app management tenant policy"
        description = "Default tenant policy that enforces app management restrictions on applications and service principals. To apply policy to targeted resources, create a new policy under appManagementPolicies collection."
        applicationRestrictions = @{
            passwordCredentials = @()
            keyCredentials = @()
        }
        servicePrincipalRestrictions = @{
            passwordCredentials = @()
            keyCredentials = @()
        }
    }
    
    Write-Host "   ✅ Default configuration prepared" -ForegroundColor Green
    Write-Host ""
    
    # Display request body
    Write-Host "📤 Request Body" -ForegroundColor Cyan
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    Write-Host ""
    Write-Host ($requestBody | ConvertTo-Json -Depth 10) -ForegroundColor Gray
    Write-Host ""
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    Write-Host ""
    
    # Reset the policy
    Write-Host "🔄 Resetting policy to defaults..." -ForegroundColor Cyan
    
    $uri = "https://graph.microsoft.com/beta/policies/defaultAppManagementPolicy"
    $response = Invoke-MgGraphRequest -Method PATCH -Uri $uri -Body ($requestBody | ConvertTo-Json -Depth 10) -ContentType "application/json"
    
    Write-Host "✅ Policy reset successfully!" -ForegroundColor Green
    Write-Host "   Note: PATCH returns 204 No Content (no response body)" -ForegroundColor Gray
    Write-Host ""
    
    # Wait for eventual consistency
    Write-Host "⏱️  Waiting 5 seconds for eventual consistency..." -ForegroundColor Gray
    Start-Sleep -Seconds 5
    
    # Verify reset
    Write-Host "🔍 Verifying reset policy..." -ForegroundColor Cyan
    
    $verifyPolicy = Invoke-MgGraphRequest -Method GET -Uri $uri
    
    Write-Host "✅ Verified reset policy" -ForegroundColor Green
    Write-Host ""
    
    # Display verification
    Write-Host "📋 Verified Policy State" -ForegroundColor Cyan
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    Write-Host ""
    
    $appPwdCount = if ($verifyPolicy.applicationRestrictions.passwordCredentials) { $verifyPolicy.applicationRestrictions.passwordCredentials.Count } else { 0 }
    $appKeyCount = if ($verifyPolicy.applicationRestrictions.keyCredentials) { $verifyPolicy.applicationRestrictions.keyCredentials.Count } else { 0 }
    $spPwdCount = if ($verifyPolicy.servicePrincipalRestrictions.passwordCredentials) { $verifyPolicy.servicePrincipalRestrictions.passwordCredentials.Count } else { 0 }
    $spKeyCount = if ($verifyPolicy.servicePrincipalRestrictions.keyCredentials) { $verifyPolicy.servicePrincipalRestrictions.keyCredentials.Count } else { 0 }
    
    Write-Host "   Is Enabled: $($verifyPolicy.isEnabled)" -ForegroundColor White
    Write-Host "   Application Password Restrictions: $appPwdCount" -ForegroundColor White
    Write-Host "   Application Key Restrictions: $appKeyCount" -ForegroundColor White
    Write-Host "   Service Principal Password Restrictions: $spPwdCount" -ForegroundColor White
    Write-Host "   Service Principal Key Restrictions: $spKeyCount" -ForegroundColor White
    Write-Host ""
    
    if ($appPwdCount -eq 0 -and $appKeyCount -eq 0 -and $spPwdCount -eq 0 -and $spKeyCount -eq 0) {
        Write-Host "✅ Policy successfully reset to defaults (no restrictions)" -ForegroundColor Green
    } else {
        Write-Host "⚠️  Warning: Policy still has $($appPwdCount + $appKeyCount + $spPwdCount + $spKeyCount) restriction(s)" -ForegroundColor Yellow
        Write-Host "   This may indicate the API hasn't fully cleared the restrictions yet." -ForegroundColor Yellow
    }
    
    Write-Host ""
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    Write-Host ""
    Write-Host "🎉 Operation completed!" -ForegroundColor Green
    Write-Host ""
}
catch {
    Write-Host ""
    Write-Host "❌ Error occurred:" -ForegroundColor Red
    Write-Host "   $($_.Exception.Message)" -ForegroundColor Red
    Write-Host ""
    exit 1
}
finally {
    # Disconnect from Microsoft Graph
    Write-Host "🔌 Disconnecting from Microsoft Graph..." -ForegroundColor Cyan
    Disconnect-MgGraph | Out-Null
    Write-Host "✅ Disconnected" -ForegroundColor Green
    Write-Host ""
}
