<#
.SYNOPSIS
    Retrieves and permanently deletes all Agent Identities from Microsoft Entra ID.

.DESCRIPTION
    This script connects to Microsoft Graph API and:
    1. Lists all Agent Identities in the tenant
    2. Soft deletes each Agent Identity (moves to deleted items)
    3. Permanently deletes each Agent Identity from deleted items

.PARAMETER TenantId
    The Azure AD tenant ID.

.PARAMETER ClientId
    The application (client) ID used for authentication.

.PARAMETER ClientSecret
    The client secret for the application.

.PARAMETER WhatIf
    Shows what would be deleted without actually deleting.

.EXAMPLE
    .\Remove-AllAgentIdentities.ps1 -TenantId "xxx" -ClientId "xxx" -ClientSecret "xxx"

.EXAMPLE
    .\Remove-AllAgentIdentities.ps1 -TenantId "xxx" -ClientId "xxx" -ClientSecret "xxx" -WhatIf

.NOTES
    Requires Microsoft.Graph.Authentication module.
    Requires AgentIdentity.Read.All and AgentIdentity.DeleteRestore.All permissions.
#>

[CmdletBinding(SupportsShouldProcess)]
param(
    [Parameter(Mandatory = $true)]
    [string]$TenantId,

    [Parameter(Mandatory = $true)]
    [string]$ClientId,

    [Parameter(Mandatory = $true)]
    [string]$ClientSecret
)

#region Functions

function Connect-ToGraph {
    param(
        [string]$TenantId,
        [string]$ClientId,
        [string]$ClientSecret
    )

    Write-Host "🔐 Connecting to Microsoft Graph..." -ForegroundColor Cyan

    Import-Module Microsoft.Graph.Authentication -ErrorAction Stop

    $secureSecret = ConvertTo-SecureString -String $ClientSecret -AsPlainText -Force
    $credential = New-Object -TypeName System.Management.Automation.PSCredential -ArgumentList $ClientId, $secureSecret

    Connect-MgGraph -ClientSecretCredential $credential -TenantId $TenantId -NoWelcome

    Write-Host "✅ Connected" -ForegroundColor Green
    Write-Host ""
}

function Get-AllAgentIdentities {
    Write-Host "📋 Retrieving all Agent Identities..." -ForegroundColor Cyan

    $uri = "https://graph.microsoft.com/beta/servicePrincipals/microsoft.graph.agentIdentity"
    $allIdentities = @()

    try {
        do {
            $response = Invoke-MgGraphRequest -Method GET -Uri $uri
            if ($response.value) {
                $allIdentities += $response.value
            }
            $uri = $response.'@odata.nextLink'
        } while ($uri)

        Write-Host "   Found $($allIdentities.Count) Agent Identities" -ForegroundColor Yellow
        return $allIdentities
    }
    catch {
        Write-Host "   ❌ Failed to retrieve Agent Identities: $($_.Exception.Message)" -ForegroundColor Red
        return @()
    }
}

function Remove-AgentIdentitySoft {
    param(
        [string]$Id,
        [string]$DisplayName
    )

    $uri = "https://graph.microsoft.com/beta/servicePrincipals/$Id"

    try {
        Invoke-MgGraphRequest -Method DELETE -Uri $uri
        Write-Host "      ✅ Soft deleted" -ForegroundColor Green
        return $true
    }
    catch {
        $errorMessage = $_.Exception.Message
        if ($errorMessage -like "*does not exist*" -or $errorMessage -like "*404*") {
            Write-Host "      ⚠️  Already deleted or not found" -ForegroundColor Yellow
            return $true
        }
        Write-Host "      ❌ Soft delete failed: $errorMessage" -ForegroundColor Red
        return $false
    }
}

function Remove-AgentIdentityPermanent {
    param(
        [string]$Id,
        [string]$DisplayName,
        [int]$MaxRetries = 5,
        [int]$RetryDelaySeconds = 2
    )

    $uri = "https://graph.microsoft.com/beta/directory/deletedItems/$Id"

    for ($attempt = 1; $attempt -le $MaxRetries; $attempt++) {
        try {
            Invoke-MgGraphRequest -Method DELETE -Uri $uri
            Write-Host "      ✅ Permanently deleted" -ForegroundColor Green
            return $true
        }
        catch {
            $errorMessage = $_.Exception.Message
            $isNotFound = $errorMessage -like "*does not exist*" -or 
                          $errorMessage -like "*404*" -or 
                          $errorMessage -like "*NotFound*"

            if ($isNotFound -and $attempt -lt $MaxRetries) {
                Write-Host "      ⏳ Not yet in deleted items, retrying in ${RetryDelaySeconds}s... (attempt $attempt/$MaxRetries)" -ForegroundColor DarkGray
                Start-Sleep -Seconds $RetryDelaySeconds
                continue
            }
            elseif ($isNotFound) {
                # After all retries, assume it was auto-purged or already deleted
                Write-Host "      ⚠️  Not found in deleted items (likely auto-purged)" -ForegroundColor Yellow
                return $true
            }

            Write-Host "      ❌ Permanent delete failed: $errorMessage" -ForegroundColor Red
            return $false
        }
    }

    return $false
}

function Show-Summary {
    param(
        [int]$Total,
        [int]$SoftDeleted,
        [int]$PermanentlyDeleted,
        [int]$Failed
    )

    Write-Host ""
    Write-Host "📊 Summary" -ForegroundColor Cyan
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor DarkGray
    Write-Host "   📋 Total Agent Identities: $Total"
    Write-Host "   🗑️  Soft Deleted: $SoftDeleted"
    Write-Host "   💀 Permanently Deleted: $PermanentlyDeleted"
    Write-Host "   ❌ Failed: $Failed"
    Write-Host ""
}

#endregion

#region Main

try {
    # Connect to Graph
    Connect-ToGraph -TenantId $TenantId -ClientId $ClientId -ClientSecret $ClientSecret

    # Get all Agent Identities
    $identities = Get-AllAgentIdentities

    if ($identities.Count -eq 0) {
        Write-Host "✅ No Agent Identities found. Nothing to delete." -ForegroundColor Green
        Disconnect-MgGraph | Out-Null
        exit 0
    }

    # Display identities to be deleted
    Write-Host ""
    Write-Host "🔍 Agent Identities to delete:" -ForegroundColor Cyan
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor DarkGray
    foreach ($identity in $identities) {
        Write-Host "   • $($identity.displayName) ($($identity.id))" -ForegroundColor White
    }
    Write-Host ""

    # Counters
    $softDeleted = 0
    $permanentlyDeleted = 0
    $failed = 0

    # Process deletions
    if ($WhatIfPreference) {
        Write-Host "⚠️  WhatIf mode - no deletions will be performed" -ForegroundColor Yellow
        Write-Host ""
        Show-Summary -Total $identities.Count -SoftDeleted 0 -PermanentlyDeleted 0 -Failed 0
    }
    else {
        Write-Host "🗑️  Deleting Agent Identities..." -ForegroundColor Cyan
        Write-Host ""

        foreach ($identity in $identities) {
            $id = $identity.id
            $displayName = $identity.displayName

            Write-Host "   ▶️  Processing: $displayName ($id)" -ForegroundColor White

            # Step 1: Soft delete
            Write-Host "      Step 1: Soft delete..." -ForegroundColor DarkGray
            $softResult = Remove-AgentIdentitySoft -Id $id -DisplayName $displayName

            if ($softResult) {
                $softDeleted++

                # Step 2: Permanent delete (with built-in retry for replication latency)
                Write-Host "      Step 2: Permanent delete..." -ForegroundColor DarkGray
                $permanentResult = Remove-AgentIdentityPermanent -Id $id -DisplayName $displayName

                if ($permanentResult) {
                    $permanentlyDeleted++
                }
                else {
                    $failed++
                }
            }
            else {
                $failed++
            }

            Write-Host ""
        }

        Show-Summary -Total $identities.Count -SoftDeleted $softDeleted -PermanentlyDeleted $permanentlyDeleted -Failed $failed
    }
}
catch {
    Write-Host "❌ Error: $($_.Exception.Message)" -ForegroundColor Red
    exit 1
}
finally {
    Disconnect-MgGraph | Out-Null
}

#endregion

