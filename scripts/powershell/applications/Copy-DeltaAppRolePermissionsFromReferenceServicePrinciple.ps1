<#
.SYNOPSIS
    Copies missing app role permissions from a reference service principal to a target service principal.

.DESCRIPTION
    Compares app role assignments between two service principals and copies any missing
    permissions from the reference SP to the target SP. Resolves service principal names
    to their object IDs, applies missing permissions, and validates all assignments.

.PARAMETER TenantId
    The Entra ID tenant ID.

.PARAMETER ClientId
    The application (client) ID for authentication.

.PARAMETER ClientSecret
    The client secret for authentication.

.PARAMETER ReferenceServicePrincipalName
    The display name of the service principal to use as the reference (source of permissions).

.PARAMETER TargetServicePrincipalName
    The display name of the service principal to copy permissions TO.

.PARAMETER DryRun
    If specified, only shows what would be copied without making changes.

.EXAMPLE
    # Copy missing permissions from SP-CPSS-GLBL-AGENTS-C-01 to SP-CPSS-GLBL-APPLICATIONS-C-01
    .\Copy-DeltaAppRolePermissionsFromReferenceServicePrinciple.ps1 `
        -TenantId "your-tenant-id" `
        -ClientId "your-client-id" `
        -ClientSecret "your-secret" `
        -ReferenceServicePrincipalName "SP-CPSS-GLBL-AGENTS-C-01" `
        -TargetServicePrincipalName "SP-CPSS-GLBL-APPLICATIONS-C-01"

.EXAMPLE
    # Dry run to see what would be copied
    .\Copy-DeltaAppRolePermissionsFromReferenceServicePrinciple.ps1 `
        -TenantId "your-tenant-id" `
        -ClientId "your-client-id" `
        -ClientSecret "your-secret" `
        -ReferenceServicePrincipalName "SP-CPSS-GLBL-AGENTS-C-01" `
        -TargetServicePrincipalName "SP-CPSS-GLBL-APPLICATIONS-C-01" `
        -DryRun

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
    [string]$ReferenceServicePrincipalName,
    
    [Parameter(Mandatory=$true)]
    [ValidateNotNullOrEmpty()]
    [string]$TargetServicePrincipalName,
    
    [Parameter(Mandatory=$false)]
    [switch]$DryRun
)

Import-Module Microsoft.Graph.Authentication

# Function to resolve service principal name to object ID
function Resolve-ServicePrincipal {
    param (
        [Parameter(Mandatory=$true)]
        [string]$DisplayName
    )
    
    $uri = "https://graph.microsoft.com/beta/servicePrincipals?`$filter=displayName eq '$DisplayName'"
    $result = Invoke-MgGraphRequest -Method GET -Uri $uri
    
    if (-not $result.value -or $result.value.Count -eq 0) {
        throw "Service principal '$DisplayName' not found"
    }
    
    if ($result.value.Count -gt 1) {
        Write-Host "   ⚠️  Multiple service principals found with name '$DisplayName', using first match" -ForegroundColor Yellow
    }
    
    return $result.value[0]
}

# Function to get all app role assignments for a service principal
function Get-AppRoleAssignments {
    param (
        [Parameter(Mandatory=$true)]
        [string]$ServicePrincipalObjectId
    )
    
    $uri = "https://graph.microsoft.com/beta/servicePrincipals/$ServicePrincipalObjectId/appRoleAssignments"
    $result = Invoke-MgGraphRequest -Method GET -Uri $uri
    
    return $result.value
}

# Function to get app role name from resource service principal
function Get-AppRoleName {
    param (
        [Parameter(Mandatory=$true)]
        [string]$ResourceId,
        
        [Parameter(Mandatory=$true)]
        [string]$AppRoleId
    )
    
    try {
        $uri = "https://graph.microsoft.com/beta/servicePrincipals/$ResourceId"
        $resourceSp = Invoke-MgGraphRequest -Method GET -Uri $uri
        
        $appRole = $resourceSp.appRoles | Where-Object { $_.id -eq $AppRoleId }
        
        if ($appRole) {
            return $appRole.value
        }
        return "Unknown"
    }
    catch {
        return "Unknown"
    }
}

# Function to assign an app role
function Add-AppRoleAssignment {
    param (
        [Parameter(Mandatory=$true)]
        [string]$TargetServicePrincipalObjectId,
        
        [Parameter(Mandatory=$true)]
        [string]$ResourceId,
        
        [Parameter(Mandatory=$true)]
        [string]$AppRoleId
    )
    
    $body = @{
        principalId = $TargetServicePrincipalObjectId
        resourceId = $ResourceId
        appRoleId = $AppRoleId
    } | ConvertTo-Json
    
    $uri = "https://graph.microsoft.com/beta/servicePrincipals/$ResourceId/appRoleAssignedTo"
    
    $result = Invoke-MgGraphRequest -Method POST -Uri $uri -Body $body -ContentType "application/json"
    
    return $result
}

# Function to validate app role assignments
function Test-AppRoleAssignments {
    param (
        [Parameter(Mandatory=$true)]
        [string]$ServicePrincipalObjectId,
        
        [Parameter(Mandatory=$true)]
        [array]$ExpectedAppRoleIds
    )
    
    $currentAssignments = Get-AppRoleAssignments -ServicePrincipalObjectId $ServicePrincipalObjectId
    $currentAppRoleIds = $currentAssignments | ForEach-Object { $_.appRoleId }
    
    $missing = @()
    foreach ($expectedId in $ExpectedAppRoleIds) {
        if ($expectedId -notin $currentAppRoleIds) {
            $missing += $expectedId
        }
    }
    
    return @{
        TotalExpected = $ExpectedAppRoleIds.Count
        TotalAssigned = ($currentAppRoleIds | Where-Object { $_ -in $ExpectedAppRoleIds }).Count
        Missing = $missing
        AllAssigned = ($missing.Count -eq 0)
    }
}

# Main Script Execution
try {
    Write-Host ""
    Write-Host "═══════════════════════════════════════════════════════════════════" -ForegroundColor Magenta
    Write-Host "   📋 Copy Missing App Role Permissions From Reference" -ForegroundColor Magenta
    if ($DryRun) {
        Write-Host "   🔍 DRY RUN MODE - No changes will be made" -ForegroundColor Yellow
    }
    Write-Host "═══════════════════════════════════════════════════════════════════" -ForegroundColor Magenta
    Write-Host ""
    
    # Connect to Microsoft Graph
    Write-Host "🔐 Connecting to Microsoft Graph..." -ForegroundColor Cyan
    
    $secureSecret = ConvertTo-SecureString -String $ClientSecret -AsPlainText -Force
    $credential = New-Object -TypeName System.Management.Automation.PSCredential -ArgumentList $ClientId, $secureSecret
    
    Connect-MgGraph -ClientSecretCredential $credential -TenantId $TenantId -NoWelcome
    
    Write-Host "✅ Connected" -ForegroundColor Green
    Write-Host ""
    
    # Resolve Reference Service Principal
    Write-Host "🔍 Resolving Reference Service Principal: $ReferenceServicePrincipalName" -ForegroundColor Cyan
    $referenceSp = Resolve-ServicePrincipal -DisplayName $ReferenceServicePrincipalName
    Write-Host "   ✅ Found: $($referenceSp.displayName)" -ForegroundColor Green
    Write-Host "      App ID: $($referenceSp.appId)" -ForegroundColor Gray
    Write-Host "      Object ID: $($referenceSp.id)" -ForegroundColor Gray
    Write-Host ""
    
    # Resolve Target Service Principal
    Write-Host "🔍 Resolving Target Service Principal: $TargetServicePrincipalName" -ForegroundColor Cyan
    $targetSp = Resolve-ServicePrincipal -DisplayName $TargetServicePrincipalName
    Write-Host "   ✅ Found: $($targetSp.displayName)" -ForegroundColor Green
    Write-Host "      App ID: $($targetSp.appId)" -ForegroundColor Gray
    Write-Host "      Object ID: $($targetSp.id)" -ForegroundColor Gray
    Write-Host ""
    
    # Get app role assignments from both SPs
    Write-Host "📦 Getting app role assignments..." -ForegroundColor Cyan
    
    $referenceAssignments = Get-AppRoleAssignments -ServicePrincipalObjectId $referenceSp.id
    $targetAssignments = Get-AppRoleAssignments -ServicePrincipalObjectId $targetSp.id
    
    Write-Host "   Reference SP: $($referenceAssignments.Count) assignments" -ForegroundColor White
    Write-Host "   Target SP: $($targetAssignments.Count) assignments" -ForegroundColor White
    Write-Host ""
    
    # Build lookup of target assignments
    $targetAppRoleKeys = @{}
    foreach ($assignment in $targetAssignments) {
        $key = "$($assignment.resourceId)|$($assignment.appRoleId)"
        $targetAppRoleKeys[$key] = $true
    }
    
    # Find missing assignments
    $missingAssignments = @()
    foreach ($assignment in $referenceAssignments) {
        $key = "$($assignment.resourceId)|$($assignment.appRoleId)"
        if (-not $targetAppRoleKeys.ContainsKey($key)) {
            $missingAssignments += $assignment
        }
    }
    
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    Write-Host "📊 Delta Analysis" -ForegroundColor Cyan
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    Write-Host ""
    
    if ($missingAssignments.Count -eq 0) {
        Write-Host "✅ Target SP already has all permissions from Reference SP" -ForegroundColor Green
        Write-Host ""
    }
    else {
        Write-Host "🔄 Missing $($missingAssignments.Count) permission(s) in Target SP:" -ForegroundColor Yellow
        Write-Host ""
        
        # Group by resource for display
        $groupedMissing = $missingAssignments | Group-Object -Property resourceId
        
        foreach ($group in $groupedMissing) {
            # Get resource SP name
            $resourceUri = "https://graph.microsoft.com/beta/servicePrincipals/$($group.Name)"
            $resourceSp = Invoke-MgGraphRequest -Method GET -Uri $resourceUri
            
            Write-Host "   📦 Resource: $($resourceSp.displayName)" -ForegroundColor Yellow
            
            foreach ($assignment in $group.Group) {
                $roleName = Get-AppRoleName -ResourceId $assignment.resourceId -AppRoleId $assignment.appRoleId
                Write-Host "      • $roleName ($($assignment.appRoleId))" -ForegroundColor White
            }
            Write-Host ""
        }
        
        # Apply missing permissions if not dry run
        if (-not $DryRun) {
            Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
            Write-Host "🚀 Applying Missing Permissions" -ForegroundColor Cyan
            Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
            Write-Host ""
            
            $successCount = 0
            $skipCount = 0
            $failCount = 0
            
            foreach ($assignment in $missingAssignments) {
                $roleName = Get-AppRoleName -ResourceId $assignment.resourceId -AppRoleId $assignment.appRoleId
                Write-Host "▶️  Assigning: $roleName" -ForegroundColor White
                
                try {
                    Add-AppRoleAssignment `
                        -TargetServicePrincipalObjectId $targetSp.id `
                        -ResourceId $assignment.resourceId `
                        -AppRoleId $assignment.appRoleId | Out-Null
                    
                    Write-Host "   ✅ Assigned successfully" -ForegroundColor Green
                    $successCount++
                }
                catch {
                    $errorMessage = $_.Exception.Message
                    $errorDetails = $_.ErrorDetails.Message
                    
                    # Check if already exists
                    $alreadyExists = $false
                    if ($errorMessage -like "*already exists*" -or $errorDetails -like "*already exists*") {
                        $alreadyExists = $true
                    }
                    
                    if ($errorDetails) {
                        try {
                            $errorObj = $errorDetails | ConvertFrom-Json
                            if ($errorObj.error.message -like "*already exists*") {
                                $alreadyExists = $true
                            }
                        }
                        catch { }
                    }
                    
                    if ($alreadyExists) {
                        Write-Host "   ⚠️  Already assigned" -ForegroundColor Yellow
                        $skipCount++
                    }
                    else {
                        Write-Host "   ❌ Failed: $errorMessage" -ForegroundColor Red
                        $failCount++
                    }
                }
                Write-Host ""
            }
            
            Write-Host "📊 Assignment Summary" -ForegroundColor Cyan
            Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
            Write-Host "   ✅ Assigned: $successCount" -ForegroundColor Green
            Write-Host "   ⚠️  Already existed: $skipCount" -ForegroundColor Yellow
            Write-Host "   ❌ Failed: $failCount" -ForegroundColor Red
            Write-Host ""
            
            # Validation step
            Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
            Write-Host "🔍 Validation" -ForegroundColor Cyan
            Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
            Write-Host ""
            
            # Wait briefly for propagation
            Write-Host "⏳ Waiting for permission propagation..." -ForegroundColor Gray
            Start-Sleep -Seconds 3
            
            # Get expected app role IDs from reference
            $expectedAppRoleIds = $referenceAssignments | ForEach-Object { $_.appRoleId }
            
            # Validate
            $validation = Test-AppRoleAssignments `
                -ServicePrincipalObjectId $targetSp.id `
                -ExpectedAppRoleIds $expectedAppRoleIds
            
            Write-Host "   Expected permissions: $($validation.TotalExpected)" -ForegroundColor White
            Write-Host "   Assigned permissions: $($validation.TotalAssigned)" -ForegroundColor White
            
            if ($validation.AllAssigned) {
                Write-Host ""
                Write-Host "   ✅ All permissions successfully assigned and consented" -ForegroundColor Green
            }
            else {
                Write-Host ""
                Write-Host "   ⚠️  $($validation.Missing.Count) permission(s) still missing:" -ForegroundColor Yellow
                foreach ($missingId in $validation.Missing) {
                    # Find which resource this belongs to
                    $matchingAssignment = $referenceAssignments | Where-Object { $_.appRoleId -eq $missingId } | Select-Object -First 1
                    if ($matchingAssignment) {
                        $roleName = Get-AppRoleName -ResourceId $matchingAssignment.resourceId -AppRoleId $missingId
                        Write-Host "      • $roleName ($missingId)" -ForegroundColor Red
                    }
                    else {
                        Write-Host "      • $missingId" -ForegroundColor Red
                    }
                }
            }
            Write-Host ""
        }
        else {
            Write-Host "ℹ️  Dry run complete. Use without -DryRun to apply changes." -ForegroundColor Cyan
            Write-Host ""
        }
    }
    
    Write-Host "🎉 Operation completed!" -ForegroundColor Green
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

