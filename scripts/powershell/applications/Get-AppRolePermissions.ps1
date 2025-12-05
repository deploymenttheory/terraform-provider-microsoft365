<#
.SYNOPSIS
    Gets all app role permissions assigned to a service principal.

.DESCRIPTION
    Retrieves all application permissions (app roles) currently assigned to an enterprise app.
    Displays the permission names, IDs, and resource providers.

.PARAMETER TenantId
    The Entra ID tenant ID.

.PARAMETER ClientId
    The application (client) ID for authentication.

.PARAMETER ClientSecret
    The client secret for authentication.

.PARAMETER PrincipalId
    The Object ID of the service principal to query permissions for.

.PARAMETER ResourceAppId
    Optional. Filter by resource app ID (e.g., '00000003-0000-0000-c000-000000000000' for Microsoft Graph only).

.PARAMETER ExportToJson
    Whether to export the results to a JSON file.

.EXAMPLE
    # Get all permissions for an enterprise app
    .\Get-AppRolePermissions.ps1 `
        -TenantId "your-tenant-id" `
        -ClientId "your-client-id" `
        -ClientSecret "your-secret" `
        -PrincipalId "your-principal-id" `
        -ExportToJson $true

.EXAMPLE
    # Get only Microsoft Graph permissions
    .\Get-AppRolePermissions.ps1 `
        -TenantId "your-tenant-id" `
        -ClientId "your-client-id" `
        -ClientSecret "your-secret" `
        -PrincipalId "your-principal-id" `
        -ResourceAppId "00000003-0000-0000-c000-000000000000" `
        -ExportToJson $true

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
    [string]$PrincipalId,
    
    [Parameter(Mandatory=$false)]
    [string]$ResourceAppId,
    
    [Parameter(Mandatory=$false)]
    [bool]$ExportToJson = $false
)

Import-Module Microsoft.Graph.Authentication

# Function to get app role name from resource service principal
function Get-AppRoleDetails {
    param (
        [Parameter(Mandatory=$true)]
        [string]$ResourceId,
        
        [Parameter(Mandatory=$true)]
        [string]$AppRoleId
    )
    
    try {
        $uri = "https://graph.microsoft.com/v1.0/servicePrincipals/$ResourceId"
        $resourceSp = Invoke-MgGraphRequest -Method GET -Uri $uri
        
        $appRole = $resourceSp.appRoles | Where-Object { $_.id -eq $AppRoleId }
        
        if ($appRole) {
            return @{
                Value = $appRole.value
                DisplayName = $appRole.displayName
                Description = $appRole.description
            }
        } else {
            return @{
                Value = "Unknown"
                DisplayName = "Unknown Permission"
                Description = "App role ID: $AppRoleId"
            }
        }
    }
    catch {
        return @{
            Value = "Unknown"
            DisplayName = "Error retrieving"
            Description = $_.Exception.Message
        }
    }
}

# Function to export results to JSON
function Export-ResultsToJson {
    param (
        [Parameter(Mandatory=$true)]
        $Results,
        
        [Parameter(Mandatory=$true)]
        [string]$PrincipalId
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
        $fileName = "AppRoleAssignments_Get_${PrincipalId}_${timestamp}.json"
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
    Write-Host "   📋 Get App Role Permissions" -ForegroundColor Magenta
    Write-Host "═══════════════════════════════════════════════════════════════════" -ForegroundColor Magenta
    Write-Host ""
    
    Write-Host "🔐 Connecting to Microsoft Graph..." -ForegroundColor Cyan
    
    $secureSecret = ConvertTo-SecureString -String $ClientSecret -AsPlainText -Force
    $credential = New-Object -TypeName System.Management.Automation.PSCredential -ArgumentList $ClientId, $secureSecret
    
    Connect-MgGraph -ClientSecretCredential $credential -TenantId $TenantId -NoWelcome
    
    Write-Host "✅ Connected" -ForegroundColor Green
    Write-Host ""
    
    # Get service principal details
    Write-Host "🔍 Getting service principal details..." -ForegroundColor Cyan
    $uri = "https://graph.microsoft.com/v1.0/servicePrincipals/$PrincipalId"
    $servicePrincipal = Invoke-MgGraphRequest -Method GET -Uri $uri
    
    Write-Host "✅ Found service principal: $($servicePrincipal.displayName)" -ForegroundColor Green
    Write-Host ""
    
    Write-Host "📋 Service Principal" -ForegroundColor Cyan
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    Write-Host "   Display Name: $($servicePrincipal.displayName)" -ForegroundColor Green
    Write-Host "   App ID: $($servicePrincipal.appId)" -ForegroundColor Green
    Write-Host "   Object ID: $($servicePrincipal.id)" -ForegroundColor Green
    Write-Host ""
    
    # Get app role assignments
    Write-Host "🔍 Getting app role assignments..." -ForegroundColor Cyan
    
    $uri = "https://graph.microsoft.com/v1.0/servicePrincipals/$PrincipalId/appRoleAssignments"
    $assignments = Invoke-MgGraphRequest -Method GET -Uri $uri
    
    if (-not $assignments.value -or $assignments.value.Count -eq 0) {
        Write-Host "⚠️  No app role assignments found" -ForegroundColor Yellow
        Write-Host ""
        return
    }
    
    # Filter by resource if specified
    if ($ResourceAppId) {
        Write-Host "🔍 Filtering by resource app ID: $ResourceAppId" -ForegroundColor Cyan
        
        $resourceUri = "https://graph.microsoft.com/v1.0/servicePrincipals?`$filter=appId eq '$ResourceAppId'"
        $resourceSp = Invoke-MgGraphRequest -Method GET -Uri $resourceUri
        
        if ($resourceSp.value -and $resourceSp.value.Count -gt 0) {
            $resourceId = $resourceSp.value[0].id
            $assignments.value = $assignments.value | Where-Object { $_.resourceId -eq $resourceId }
        }
    }
    
    Write-Host "✅ Found $($assignments.value.Count) app role assignment(s)" -ForegroundColor Green
    Write-Host ""
    
    # Process and display assignments
    Write-Host "📦 App Role Assignments" -ForegroundColor Cyan
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    Write-Host ""
    
    $detailedAssignments = @()
    $groupedByResource = $assignments.value | Group-Object -Property resourceId
    
    foreach ($group in $groupedByResource) {
        # Get resource service principal details
        $resourceUri = "https://graph.microsoft.com/v1.0/servicePrincipals/$($group.Name)"
        $resourceSp = Invoke-MgGraphRequest -Method GET -Uri $resourceUri
        
        Write-Host "📦 Resource: $($resourceSp.displayName)" -ForegroundColor Yellow
        Write-Host "   App ID: $($resourceSp.appId)" -ForegroundColor Gray
        Write-Host "   Object ID: $($resourceSp.id)" -ForegroundColor Gray
        Write-Host ""
        
        foreach ($assignment in $group.Group) {
            $roleDetails = Get-AppRoleDetails -ResourceId $assignment.resourceId -AppRoleId $assignment.appRoleId
            
            Write-Host "   ✓ $($roleDetails.Value)" -ForegroundColor Green
            Write-Host "     Display Name: $($roleDetails.DisplayName)" -ForegroundColor Gray
            Write-Host "     App Role ID: $($assignment.appRoleId)" -ForegroundColor Gray
            Write-Host "     Assignment ID: $($assignment.id)" -ForegroundColor Gray
            Write-Host ""
            
            $detailedAssignments += @{
                ResourceDisplayName = $resourceSp.displayName
                ResourceAppId = $resourceSp.appId
                ResourceObjectId = $resourceSp.id
                PermissionName = $roleDetails.Value
                PermissionDisplayName = $roleDetails.DisplayName
                PermissionDescription = $roleDetails.Description
                AppRoleId = $assignment.appRoleId
                AssignmentId = $assignment.id
                CreatedDateTime = $assignment.createdDateTime
            }
        }
    }
    
    # Display summary
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    Write-Host ""
    Write-Host "📊 Summary" -ForegroundColor Cyan
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    Write-Host "   Total Permissions: $($assignments.value.Count)" -ForegroundColor White
    Write-Host "   Resource Providers: $($groupedByResource.Count)" -ForegroundColor White
    Write-Host ""
    
    # Export results if requested
    if ($ExportToJson) {
        $summaryObject = @{
            Timestamp = Get-Date -Format "o"
            ServicePrincipal = @{
                DisplayName = $servicePrincipal.displayName
                AppId = $servicePrincipal.appId
                ObjectId = $servicePrincipal.id
            }
            Summary = @{
                TotalPermissions = $assignments.value.Count
                ResourceProviders = $groupedByResource.Count
            }
            Permissions = $detailedAssignments
        }
        
        Export-ResultsToJson -Results $summaryObject -PrincipalId $PrincipalId
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

