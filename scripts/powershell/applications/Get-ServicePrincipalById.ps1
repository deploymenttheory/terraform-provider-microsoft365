<#
.SYNOPSIS
    Gets a service principal by Object ID or App ID with all attributes.

.DESCRIPTION
    Retrieves a Microsoft Entra service principal by its Object ID or App ID (Client ID),
    displaying all attributes including tags, service principal names, and API permissions.
    Useful for comparing API values against what appears in the Azure Portal and for
    troubleshooting service principal configurations.

.PARAMETER TenantId
    The Entra ID tenant ID.

.PARAMETER ClientId
    The application (client) ID for authentication.

.PARAMETER ClientSecret
    The client secret for authentication.

.PARAMETER ServicePrincipalId
    The Object ID or App ID (Client ID) of the service principal to retrieve.
    Will automatically detect which type of ID is provided.

.PARAMETER ExportToJson
    Whether to export the results to a JSON file.

.EXAMPLE
    # Get service principal by Object ID
    .\Get-ServicePrincipalById.ps1 `
        -TenantId "your-tenant-id" `
        -ClientId "your-client-id" `
        -ClientSecret "your-secret" `
        -ServicePrincipalId "12345678-1234-1234-1234-123456789012"

.EXAMPLE
    # Get service principal by App ID and export to JSON
    .\Get-ServicePrincipalById.ps1 `
        -TenantId "your-tenant-id" `
        -ClientId "your-client-id" `
        -ClientSecret "your-secret" `
        -ServicePrincipalId "87654321-4321-4321-4321-210987654321" `
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
    [string]$ServicePrincipalId,
    
    [Parameter(Mandatory=$false)]
    [bool]$ExportToJson = $false
)

Import-Module Microsoft.Graph.Authentication

#region Helper Functions

function Format-AttributeValue {
    param (
        [Parameter(Mandatory=$false)]
        $Value
    )
    
    if ($null -eq $Value) {
        return "(null)"
    }
    elseif ($Value -is [array]) {
        if ($Value.Count -eq 0) {
            return "(empty array)"
        }
        return ($Value -join ", ")
    }
    elseif ($Value -is [hashtable] -or $Value -is [System.Collections.IDictionary]) {
        return "(complex object)"
    }
    elseif ($Value -is [bool]) {
        return $Value.ToString().ToLower()
    }
    elseif ($Value -is [string] -and [string]::IsNullOrEmpty($Value)) {
        return "(empty)"
    }
    else {
        return $Value.ToString()
    }
}

function Get-AppRoleDetails {
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

function Export-ResultsToJson {
    param (
        [Parameter(Mandatory=$true)]
        $Results,
        
        [Parameter(Mandatory=$true)]
        [string]$ServicePrincipalId
    )
    
    try {
        $outputDir = Join-Path -Path (Get-Location) -ChildPath "output"
        if (-not (Test-Path -Path $outputDir)) {
            New-Item -Path $outputDir -ItemType Directory | Out-Null
            Write-Host "📁 Created output directory: $outputDir" -ForegroundColor Gray
        }
        
        $timestamp = Get-Date -Format "yyyyMMdd_HHmmss"
        $fileName = "ServicePrincipal_${ServicePrincipalId}_${timestamp}.json"
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
    Write-Host "   🔍 Get Service Principal By ID" -ForegroundColor Magenta
    Write-Host "═══════════════════════════════════════════════════════════════════" -ForegroundColor Magenta
    Write-Host ""
    
    # Connect to Microsoft Graph
    Write-Host "🔐 Connecting to Microsoft Graph..." -ForegroundColor Cyan
    
    $secureSecret = ConvertTo-SecureString -String $ClientSecret -AsPlainText -Force
    $credential = New-Object -TypeName System.Management.Automation.PSCredential -ArgumentList $ClientId, $secureSecret
    
    Connect-MgGraph -ClientSecretCredential $credential -TenantId $TenantId -NoWelcome
    
    Write-Host "✅ Connected" -ForegroundColor Green
    Write-Host ""
    
    # Determine if the provided ID is an Object ID or App ID
    Write-Host "🔍 Retrieving service principal..." -ForegroundColor Cyan
    
    $servicePrincipal = $null
    $searchByAppId = $false
    
    # Try to get by Object ID first
    try {
        $uri = "https://graph.microsoft.com/beta/servicePrincipals/$ServicePrincipalId"
        $servicePrincipal = Invoke-MgGraphRequest -Method GET -Uri $uri
        Write-Host "✅ Found service principal by Object ID" -ForegroundColor Green
    }
    catch {
        # If not found, try searching by App ID
        Write-Host "   Not found by Object ID, trying App ID..." -ForegroundColor Gray
        $searchByAppId = $true
        
        $searchUri = "https://graph.microsoft.com/beta/servicePrincipals?`$filter=appId eq '$ServicePrincipalId'"
        $searchResult = Invoke-MgGraphRequest -Method GET -Uri $searchUri
        
        if ($searchResult.value -and $searchResult.value.Count -gt 0) {
            $servicePrincipal = $searchResult.value[0]
            Write-Host "✅ Found service principal by App ID" -ForegroundColor Green
        }
        else {
            Write-Host "❌ Service principal not found with ID: $ServicePrincipalId" -ForegroundColor Red
            Write-Host ""
            exit 1
        }
    }
    
    # Retrieve full details with specific fields
    Write-Host "🔍 Retrieving full service principal details..." -ForegroundColor Cyan
    
    $detailUri = "https://graph.microsoft.com/beta/servicePrincipals/$($servicePrincipal.id)?`$select=id,appId,displayName,accountEnabled,appRoleAssignmentRequired,description,homepage,loginUrl,logoutUrl,notes,notificationEmailAddresses,preferredSingleSignOnMode,servicePrincipalNames,servicePrincipalType,signInAudience,tags,appOwnerOrganizationId,createdDateTime,disabledByMicrosoftStatus"
    $servicePrincipal = Invoke-MgGraphRequest -Method GET -Uri $detailUri
    
    Write-Host "✅ Retrieved full details" -ForegroundColor Green
    Write-Host ""
    
    # Display basic information
    Write-Host "📋 Service Principal Overview" -ForegroundColor Cyan
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    Write-Host ""
    
    $keyAttributes = @(
        @{ Name = "id"; Label = "Object ID" },
        @{ Name = "appId"; Label = "App ID (Client ID)" },
        @{ Name = "displayName"; Label = "Display Name" },
        @{ Name = "servicePrincipalType"; Label = "Service Principal Type" },
        @{ Name = "accountEnabled"; Label = "Account Enabled" },
        @{ Name = "appRoleAssignmentRequired"; Label = "App Role Assignment Required" },
        @{ Name = "description"; Label = "Description" },
        @{ Name = "homepage"; Label = "Homepage" },
        @{ Name = "loginUrl"; Label = "Login URL" },
        @{ Name = "logoutUrl"; Label = "Logout URL" },
        @{ Name = "notes"; Label = "Notes" },
        @{ Name = "preferredSingleSignOnMode"; Label = "Preferred SSO Mode" },
        @{ Name = "signInAudience"; Label = "Sign-In Audience" },
        @{ Name = "appOwnerOrganizationId"; Label = "App Owner Org ID" },
        @{ Name = "createdDateTime"; Label = "Created Date Time" },
        @{ Name = "disabledByMicrosoftStatus"; Label = "Disabled By Microsoft" }
    )
    
    foreach ($attr in $keyAttributes) {
        $value = $servicePrincipal[$attr.Name]
        $formattedValue = Format-AttributeValue -Value $value
        Write-Host "   $($attr.Label): " -NoNewline -ForegroundColor White
        Write-Host "$formattedValue" -ForegroundColor Green
    }
    
    Write-Host ""
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    Write-Host ""
    
    # Display collection attributes
    Write-Host "📦 Collection Attributes" -ForegroundColor Cyan
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    Write-Host ""
    
    # Tags
    if ($servicePrincipal.tags -and $servicePrincipal.tags.Count -gt 0) {
        Write-Host "   🏷️  Tags ($($servicePrincipal.tags.Count)):" -ForegroundColor Yellow
        foreach ($tag in $servicePrincipal.tags) {
            Write-Host "      • $tag" -ForegroundColor Green
        }
        Write-Host ""
    } else {
        Write-Host "   🏷️  Tags: " -NoNewline -ForegroundColor Yellow
        Write-Host "(none)" -ForegroundColor Gray
        Write-Host ""
    }
    
    # Service Principal Names
    if ($servicePrincipal.servicePrincipalNames -and $servicePrincipal.servicePrincipalNames.Count -gt 0) {
        Write-Host "   🔗 Service Principal Names ($($servicePrincipal.servicePrincipalNames.Count)):" -ForegroundColor Yellow
        foreach ($name in $servicePrincipal.servicePrincipalNames) {
            Write-Host "      • $name" -ForegroundColor Green
        }
        Write-Host ""
    }
    
    # Notification Email Addresses
    if ($servicePrincipal.notificationEmailAddresses -and $servicePrincipal.notificationEmailAddresses.Count -gt 0) {
        Write-Host "   📧 Notification Email Addresses ($($servicePrincipal.notificationEmailAddresses.Count)):" -ForegroundColor Yellow
        foreach ($email in $servicePrincipal.notificationEmailAddresses) {
            Write-Host "      • $email" -ForegroundColor Green
        }
        Write-Host ""
    }
    
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    Write-Host ""
    
    # Get app role assignments
    Write-Host "🔍 Retrieving app role assignments..." -ForegroundColor Cyan
    
    $assignmentsUri = "https://graph.microsoft.com/beta/servicePrincipals/$($servicePrincipal.id)/appRoleAssignments"
    $assignments = Invoke-MgGraphRequest -Method GET -Uri $assignmentsUri
    
    $detailedAssignments = @()
    
    if (-not $assignments.value -or $assignments.value.Count -eq 0) {
        Write-Host "⚠️  No app role assignments found" -ForegroundColor Yellow
        Write-Host ""
    }
    else {
        Write-Host "✅ Found $($assignments.value.Count) app role assignment(s)" -ForegroundColor Green
        Write-Host ""
        
        Write-Host "🔐 Assigned Permissions" -ForegroundColor Cyan
        Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
        Write-Host ""
        
        $groupedByResource = $assignments.value | Group-Object -Property resourceId
        
        foreach ($group in $groupedByResource) {
            # Get resource service principal details
            $resourceUri = "https://graph.microsoft.com/beta/servicePrincipals/$($group.Name)"
            $resourceSp = Invoke-MgGraphRequest -Method GET -Uri $resourceUri
            
            Write-Host "   📦 Resource: $($resourceSp.displayName)" -ForegroundColor Yellow
            Write-Host "      App ID: $($resourceSp.appId)" -ForegroundColor Gray
            Write-Host "      Object ID: $($resourceSp.id)" -ForegroundColor Gray
            Write-Host ""
            
            foreach ($assignment in $group.Group) {
                $roleDetails = Get-AppRoleDetails -ResourceId $assignment.resourceId -AppRoleId $assignment.appRoleId
                
                Write-Host "      ✓ $($roleDetails.Value)" -ForegroundColor Green
                Write-Host "        Display Name: $($roleDetails.DisplayName)" -ForegroundColor Gray
                Write-Host "        App Role ID: $($assignment.appRoleId)" -ForegroundColor Gray
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
        
        Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
        Write-Host ""
    }
    
    # Display summary
    Write-Host "📊 Summary" -ForegroundColor Cyan
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    Write-Host "   Display Name: $($servicePrincipal.displayName)" -ForegroundColor White
    Write-Host "   Object ID: $($servicePrincipal.id)" -ForegroundColor White
    Write-Host "   App ID: $($servicePrincipal.appId)" -ForegroundColor White
    Write-Host "   Type: $($servicePrincipal.servicePrincipalType)" -ForegroundColor White
    Write-Host "   Account Enabled: $(Format-AttributeValue -Value $servicePrincipal.accountEnabled)" -ForegroundColor White
    Write-Host "   App Role Assignment Required: $(Format-AttributeValue -Value $servicePrincipal.appRoleAssignmentRequired)" -ForegroundColor White
    
    $tagsCount = if ($servicePrincipal.tags) { $servicePrincipal.tags.Count } else { 0 }
    Write-Host "   Tags: $tagsCount" -ForegroundColor White
    
    $spNamesCount = if ($servicePrincipal.servicePrincipalNames) { $servicePrincipal.servicePrincipalNames.Count } else { 0 }
    Write-Host "   Service Principal Names: $spNamesCount" -ForegroundColor White
    
    $permissionsCount = if ($assignments.value) { $assignments.value.Count } else { 0 }
    Write-Host "   Total Permissions: $permissionsCount" -ForegroundColor White
    
    Write-Host ""
    
    # Export results if requested
    if ($ExportToJson) {
        $exportObject = @{
            Timestamp = Get-Date -Format "o"
            ServicePrincipal = $servicePrincipal
            Permissions = $detailedAssignments
            Summary = @{
                DisplayName = $servicePrincipal.displayName
                ObjectId = $servicePrincipal.id
                AppId = $servicePrincipal.appId
                ServicePrincipalType = $servicePrincipal.servicePrincipalType
                AccountEnabled = $servicePrincipal.accountEnabled
                AppRoleAssignmentRequired = $servicePrincipal.appRoleAssignmentRequired
                TagsCount = $tagsCount
                ServicePrincipalNamesCount = $spNamesCount
                TotalPermissions = $permissionsCount
            }
        }
        
        Export-ResultsToJson -Results $exportObject -ServicePrincipalId $servicePrincipal.id
    }
    
    Write-Host "🎉 Operation completed!" -ForegroundColor Green
    Write-Host ""
}
catch {
    Write-Host ""
    Write-Host "❌ Error: $_" -ForegroundColor Red
    
    if ($_.Exception.Response) {
        $statusCode = $_.Exception.Response.StatusCode.value__
        Write-Host "   Status Code: $statusCode" -ForegroundColor Gray
        
        if ($_.ErrorDetails) {
            $errorDetails = $_.ErrorDetails.Message | ConvertFrom-Json
            Write-Host "   Error Code: $($errorDetails.error.code)" -ForegroundColor Gray
            Write-Host "   Error Message: $($errorDetails.error.message)" -ForegroundColor Gray
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
