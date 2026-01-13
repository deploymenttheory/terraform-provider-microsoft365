<#
.SYNOPSIS
    Deletes Conditional Access policies based on pattern matching or deletes all policies.

.DESCRIPTION
    This script connects to Microsoft Graph and deletes Conditional Access policies.
    It supports deletion by regex pattern matching or deletion of all policies.
    Includes WhatIf support for safe dry-run testing.

.PARAMETER TenantId
    The Entra ID tenant ID (Directory ID).

.PARAMETER ClientId
    The application (client) ID for authentication.

.PARAMETER ClientSecret
    The client secret for authentication.

.PARAMETER Pattern
    Regex pattern to match against policy display names.
    Only policies matching this pattern will be deleted.

.PARAMETER DeleteAll
    Switch to delete all Conditional Access policies in the tenant.
    Use with extreme caution!

.PARAMETER WhatIf
    Shows what would be deleted without actually deleting anything.
    Recommended to run first before actual deletion.

.PARAMETER HardDelete
    Permanently deletes policies (removes from deleted items).
    If not specified, policies are only soft-deleted.

.EXAMPLE
    # Dry-run: See what would be deleted with pattern matching
    .\Delete-ConditionalAccessPolicies.ps1 `
        -TenantId "your-tenant-id" `
        -ClientId "your-client-id" `
        -ClientSecret "your-secret" `
        -Pattern "^TEST-.*" `
        -WhatIf

.EXAMPLE
    # Delete all policies matching pattern
    .\Delete-ConditionalAccessPolicies.ps1 `
        -TenantId "your-tenant-id" `
        -ClientId "your-client-id" `
        -ClientSecret "your-secret" `
        -Pattern "^labtest-.*"

.EXAMPLE
    # Delete all policies with hard delete (permanent)
    .\Delete-ConditionalAccessPolicies.ps1 `
        -TenantId "your-tenant-id" `
        -ClientId "your-client-id" `
        -ClientSecret "your-secret" `
        -DeleteAll `
        -HardDelete `
        -WhatIf

.NOTES
    Author: Deployment Theory
    Requires: Microsoft.Graph.Authentication module
    
    Required Permissions:
    - Policy.Read.All (to list policies)
    - Policy.ReadWrite.ConditionalAccess (to delete policies)
#>

[CmdletBinding(SupportsShouldProcess=$true, ConfirmImpact='High')]
param (
    [Parameter(Mandatory=$true,
    HelpMessage="Specify the Entra ID tenant ID (Directory ID)")]
    [ValidateNotNullOrEmpty()]
    [string]$TenantId,

    [Parameter(Mandatory=$true,
    HelpMessage="Specify the application (client) ID")]
    [ValidateNotNullOrEmpty()]
    [string]$ClientId,
    
    [Parameter(Mandatory=$true,
    HelpMessage="Specify the client secret")]
    [ValidateNotNullOrEmpty()]
    [string]$ClientSecret,
    
    [Parameter(Mandatory=$false,
    ParameterSetName="Pattern",
    HelpMessage="Regex pattern to match policy display names")]
    [string]$Pattern,
    
    [Parameter(Mandatory=$false,
    ParameterSetName="DeleteAll",
    HelpMessage="Delete ALL Conditional Access policies")]
    [switch]$DeleteAll,
    
    [Parameter(Mandatory=$false,
    HelpMessage="Permanently delete policies (hard delete)")]
    [switch]$HardDelete
)

# Import required modules
Import-Module Microsoft.Graph.Authentication

# Function to get all Conditional Access policies
function Get-AllConditionalAccessPolicies {
    try {
        Write-Host "🔍 Retrieving Conditional Access policies..." -ForegroundColor Cyan
        
        $uri = "https://graph.microsoft.com/beta/identity/conditionalAccess/policies"
        $response = Invoke-MgGraphRequest -Method GET -Uri $uri
        
        $policies = $response.value
        
        Write-Host "   Found $($policies.Count) policy/policies" -ForegroundColor Gray
        Write-Host ""
        
        return $policies
    }
    catch {
        Write-Host "❌ Error retrieving policies: $_" -ForegroundColor Red
        
        if ($_.Exception.Response) {
            $statusCode = $_.Exception.Response.StatusCode
            Write-Host "   Status Code: $statusCode" -ForegroundColor Red
        }
        
        throw
    }
}

# Function to delete a Conditional Access policy
function Remove-ConditionalAccessPolicy {
    param (
        [Parameter(Mandatory=$true)]
        [string]$PolicyId,
        
        [Parameter(Mandatory=$true)]
        [string]$PolicyName,
        
        [Parameter(Mandatory=$false)]
        [switch]$HardDelete
    )
    
    try {
        # Soft delete
        $uri = "https://graph.microsoft.com/beta/identity/conditionalAccess/policies/$PolicyId"
        
        Write-Host "   🗑️  Deleting: $PolicyName" -ForegroundColor Yellow
        Write-Host "      Policy ID: $PolicyId" -ForegroundColor Gray
        
        Invoke-MgGraphRequest -Method DELETE -Uri $uri
        
        Write-Host "      ✅ Soft deleted successfully" -ForegroundColor Green
        
        # Hard delete if requested
        if ($HardDelete) {
            Write-Host "      ⏳ Waiting 10 seconds before hard delete..." -ForegroundColor Cyan
            Start-Sleep -Seconds 10
            
            $hardDeleteUri = "https://graph.microsoft.com/beta/identity/conditionalAccess/deletedItems/policies/$PolicyId"
            
            Write-Host "      🔥 Performing hard delete (permanent)..." -ForegroundColor Red
            Invoke-MgGraphRequest -Method DELETE -Uri $hardDeleteUri
            
            Write-Host "      ✅ Hard deleted successfully (permanent)" -ForegroundColor Green
        }
        
        Write-Host ""
        return $true
    }
    catch {
        Write-Host "      ❌ Error deleting policy: $_" -ForegroundColor Red
        
        if ($_.Exception.Response) {
            $statusCode = $_.Exception.Response.StatusCode
            Write-Host "         Status Code: $statusCode" -ForegroundColor Red
        }
        
        Write-Host ""
        return $false
    }
}

# Function to display policy summary
function Show-PolicySummary {
    param (
        [Parameter(Mandatory=$true)]
        $Policies,
        
        [Parameter(Mandatory=$false)]
        [string]$FilterPattern
    )
    
    Write-Host "📋 Policies to be deleted:" -ForegroundColor Cyan
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    
    $index = 1
    foreach ($policy in $Policies) {
        Write-Host "   $index. $($policy.displayName)" -ForegroundColor Yellow
        Write-Host "      ID: $($policy.id)" -ForegroundColor Gray
        Write-Host "      State: $($policy.state)" -ForegroundColor Gray
        $index++
    }
    
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    Write-Host "   Total: $($Policies.Count) policy/policies" -ForegroundColor White
    Write-Host ""
}

# Validate parameters
if (-not $Pattern -and -not $DeleteAll) {
    Write-Host "❌ Error: You must specify either -Pattern or -DeleteAll" -ForegroundColor Red
    Write-Host ""
    Write-Host "Examples:" -ForegroundColor Cyan
    Write-Host "  Delete by pattern:  -Pattern '^TEST-.*'" -ForegroundColor Gray
    Write-Host "  Delete all:         -DeleteAll" -ForegroundColor Gray
    exit 1
}

# Warning for DeleteAll
if ($DeleteAll -and -not $WhatIfPreference) {
    Write-Host "⚠️  WARNING: You are about to delete ALL Conditional Access policies!" -ForegroundColor Red
    Write-Host "   This is a destructive operation that affects tenant security." -ForegroundColor Red
    Write-Host "   Consider running with -WhatIf first to preview the changes." -ForegroundColor Yellow
    Write-Host ""
    
    $confirmation = Read-Host "Type 'DELETE ALL' (case-sensitive) to confirm"
    if ($confirmation -ne "DELETE ALL") {
        Write-Host "❌ Operation cancelled. Confirmation text did not match." -ForegroundColor Red
        exit 0
    }
    Write-Host ""
}

# Script execution
try {
    Write-Host "🔐 Connecting to Microsoft Graph..." -ForegroundColor Cyan
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    
    $secureClientSecret = ConvertTo-SecureString -String $ClientSecret -AsPlainText -Force
    $clientSecretCredential = New-Object -TypeName System.Management.Automation.PSCredential -ArgumentList $ClientId, $secureClientSecret
    
    Connect-MgGraph -ClientSecretCredential $clientSecretCredential -TenantId $TenantId -NoWelcome
    
    Write-Host "✅ Connected to Microsoft Graph successfully" -ForegroundColor Green
    Write-Host ""
    
    # Get all policies
    $allPolicies = Get-AllConditionalAccessPolicies
    
    if ($allPolicies.Count -eq 0) {
        Write-Host "📊 No policies found in tenant" -ForegroundColor Yellow
        exit 0
    }
    
    # Filter policies based on parameters
    if ($DeleteAll) {
        $policiesToDelete = $allPolicies
        Write-Host "🎯 Mode: Delete ALL policies" -ForegroundColor Magenta
    }
    else {
        Write-Host "🎯 Mode: Delete by pattern matching" -ForegroundColor Magenta
        Write-Host "   Pattern: $Pattern" -ForegroundColor Gray
        Write-Host ""
        
        $policiesToDelete = $allPolicies | Where-Object { $_.displayName -match $Pattern }
        
        if ($policiesToDelete.Count -eq 0) {
            Write-Host "📊 No policies matched the pattern: $Pattern" -ForegroundColor Yellow
            exit 0
        }
    }
    
    Write-Host ""
    
    # Show summary
    Show-PolicySummary -Policies $policiesToDelete -FilterPattern $Pattern
    
    # WhatIf mode
    if ($WhatIfPreference) {
        Write-Host "💡 WhatIf Mode: No policies will be deleted" -ForegroundColor Cyan
        Write-Host "   Remove -WhatIf to perform actual deletion" -ForegroundColor Gray
        Write-Host ""
        
        if ($HardDelete) {
            Write-Host "⚠️  Hard delete mode enabled (permanent deletion)" -ForegroundColor Red
        } else {
            Write-Host "ℹ️  Soft delete mode (policies can be recovered)" -ForegroundColor Cyan
        }
        
        Write-Host ""
        Write-Host "🎉 WhatIf completed - no changes made" -ForegroundColor Green
        exit 0
    }
    
    # Confirmation for pattern deletion
    if (-not $DeleteAll) {
        Write-Host "⚠️  You are about to delete $($policiesToDelete.Count) policy/policies" -ForegroundColor Yellow
        if ($HardDelete) {
            Write-Host "⚠️  Hard delete is enabled - deletion will be PERMANENT" -ForegroundColor Red
        }
        Write-Host ""
        
        $confirmation = Read-Host "Type 'DELETE' (case-sensitive) to confirm"
        if ($confirmation -ne constants.TfOperationDelete) {
            Write-Host "❌ Operation cancelled. Confirmation text did not match." -ForegroundColor Red
            exit 0
        }
        Write-Host ""
    }
    
    # Perform deletions
    Write-Host "🗑️  Starting deletion process..." -ForegroundColor Cyan
    Write-Host ""
    
    $successCount = 0
    $failCount = 0
    
    foreach ($policy in $policiesToDelete) {
        if ($PSCmdlet.ShouldProcess($policy.displayName, "Delete Conditional Access Policy")) {
            $result = Remove-ConditionalAccessPolicy -PolicyId $policy.id -PolicyName $policy.displayName -HardDelete:$HardDelete
            
            if ($result) {
                $successCount++
            } else {
                $failCount++
            }
        }
    }
    
    # Summary
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    Write-Host "📊 Deletion Summary" -ForegroundColor Cyan
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    Write-Host "   ✅ Successfully deleted: $successCount" -ForegroundColor Green
    Write-Host "   ❌ Failed: $failCount" -ForegroundColor Red
    Write-Host "   📝 Total processed: $($policiesToDelete.Count)" -ForegroundColor White
    
    if ($HardDelete) {
        Write-Host "   🔥 Hard delete mode: Policies permanently removed" -ForegroundColor Red
    } else {
        Write-Host "   ♻️  Soft delete mode: Policies can be recovered" -ForegroundColor Cyan
    }
    Write-Host ""
    
    if ($successCount -gt 0) {
        Write-Host "🎉 Operation completed successfully!" -ForegroundColor Green
    } else {
        Write-Host "⚠️  Operation completed with errors" -ForegroundColor Yellow
    }
}
catch {
    Write-Host "❌ Script execution failed: $_" -ForegroundColor Red
    exit 1
}
finally {
    # Disconnect from Microsoft Graph
    Write-Host "🔌 Disconnecting from Microsoft Graph..." -ForegroundColor Cyan
    try {
        Disconnect-MgGraph 2>$null
        Write-Host "✅ Disconnected from Microsoft Graph" -ForegroundColor Green
    }
    catch {
        # Ignore disconnect errors
    }
}

