# Microsoft 365 Service Health Monitor
# Following exact patterns from Intune upload script
[CmdletBinding()]
param (
    [Parameter(Mandatory=$true,
    HelpMessage="Specify the Entra ID tenant ID (Directory ID) where the application is registered")]
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
    HelpMessage="Include recently resolved issues from the last 7 days")]
    [switch]$IncludeResolvedIssues,
    
    [Parameter(Mandatory=$false,
    HelpMessage="Export detailed report to JSON file")]
    [string]$ExportPath,
    
    [Parameter(Mandatory=$false,
    HelpMessage="Show only services with issues")]
    [switch]$ProblemsOnly,
    
    [Parameter(Mandatory=$false,
    HelpMessage="Specific issue IDs to monitor (comma-separated)")]
    [string[]]$WatchIssues,
    
    [Parameter(Mandatory=$false,
    HelpMessage="Required Microsoft Graph application permissions to validate")]
    [string[]]$RequiredPermissions = @("ServiceHealth.Read.All", "ServiceMessage.Read.All")
)

# Service status color mapping - Based on actual API responses
$script:StatusColors = @{
    "ServiceOperational"       = "Green"
    "Investigating"            = "Yellow"
    "RestoringService"         = "Cyan"
    "VerifyingService"         = "Cyan"
    "ServiceRestored"          = "Green"
    "PostIncidentReviewPublished" = "Green"
    "ServiceDegradation"       = "Yellow"
    "ServiceInterruption"      = "Red"
    "ExtendedRecovery"         = "Yellow"
    "FalsePositive"            = "Gray"
    "InvestigationSuspended"   = "Magenta"
    "Resolved"                 = "Green"
    "MitigatedExternal"        = "Green"
    "Mitigated"                = "Green"
    "ResolvedExternal"         = "Green"
    "Confirmed"                = "Yellow"
    "Reported"                 = "Yellow"
    # Issue classification types
    "Advisory"                 = "Blue"
    "Incident"                 = "Red"
}

# Function to check if service principal has required permissions
function Test-ServicePrincipalPermissions {
    param (
        [Parameter(Mandatory=$true)]
        [string]$ClientId,
        
        [Parameter(Mandatory=$true)]
        [string[]]$RequiredPermissions
    )
    
    try {
        Write-Host "`n🔐 Validating service principal permissions..." -ForegroundColor Yellow
        
        # Get the service principal for this application
        $servicePrincipalFilter = "appId eq '$ClientId'"
        $encodedFilter = [System.Web.HttpUtility]::UrlEncode($servicePrincipalFilter)
        $spUri = "https://graph.microsoft.com/v1.0/servicePrincipals?`$filter=$encodedFilter"
        
        $spResponse = Invoke-MgGraphRequest -Method GET -Uri $spUri
        
        if (-not $spResponse.value -or $spResponse.value.Count -eq 0) {
            Write-Host "❌ Service principal not found for application: $ClientId" -ForegroundColor Red
            return $false
        }
        
        $servicePrincipal = $spResponse.value[0]
        $spId = $servicePrincipal.id
        
        Write-Host "   ✅ Found service principal: $($servicePrincipal.displayName)" -ForegroundColor Green
        
        # Get Microsoft Graph service principal (resource)
        $graphSpFilter = "appId eq '00000003-0000-0000-c000-000000000000'"
        $encodedGraphFilter = [System.Web.HttpUtility]::UrlEncode($graphSpFilter)
        $graphSpUri = "https://graph.microsoft.com/v1.0/servicePrincipals?`$filter=$encodedGraphFilter"
        
        $graphSpResponse = Invoke-MgGraphRequest -Method GET -Uri $graphSpUri
        $graphServicePrincipal = $graphSpResponse.value[0]
        
        # Get app role assignments for the service principal
        $assignmentsUri = "https://graph.microsoft.com/v1.0/servicePrincipals/$spId/appRoleAssignments"
        $assignmentsResponse = Invoke-MgGraphRequest -Method GET -Uri $assignmentsUri
        
        # Get Microsoft Graph app roles for reference
        $graphAppRoles = $graphServicePrincipal.appRoles
        
        # Build a map of role names to role IDs for Microsoft Graph
        $roleMap = @{}
        foreach ($role in $graphAppRoles) {
            $roleMap[$role.value] = $role.id
        }
        
        Write-Host "`n   📋 Checking required permissions:" -ForegroundColor Cyan
        
        $allPermissionsPresent = $true
        $missingPermissions = @()
        
        foreach ($permission in $RequiredPermissions) {
            $roleId = $roleMap[$permission]
            
            if (-not $roleId) {
                Write-Host "   ❌ " -NoNewline -ForegroundColor Red
                Write-Host "$permission" -NoNewline -ForegroundColor White
                Write-Host " - Unknown permission" -ForegroundColor Red
                $allPermissionsPresent = $false
                $missingPermissions += $permission
                continue
            }
            
            # Check if this permission is assigned
            $assignment = $assignmentsResponse.value | Where-Object { 
                $_.resourceId -eq $graphServicePrincipal.id -and $_.appRoleId -eq $roleId 
            }
            
            if ($assignment) {
                Write-Host "   ✅ " -NoNewline -ForegroundColor Green
                Write-Host "$permission" -NoNewline -ForegroundColor White
                Write-Host " - Granted" -ForegroundColor Green
            }
            else {
                Write-Host "   ❌ " -NoNewline -ForegroundColor Red
                Write-Host "$permission" -NoNewline -ForegroundColor White
                Write-Host " - Not granted" -ForegroundColor Red
                $allPermissionsPresent = $false
                $missingPermissions += $permission
            }
        }
        
        if ($allPermissionsPresent) {
            Write-Host "`n   ✅ All required permissions are present" -ForegroundColor Green
            return $true
        }
        else {
            Write-Host "`n   ❌ Missing required permissions: $($missingPermissions -join ', ')" -ForegroundColor Red
            return $false
        }
        
    }
    catch {
        Write-Host "❌ Failed to validate service principal permissions: $_" -ForegroundColor Red
        return $false
    }
}

# Service ID to friendly name mapping - Based on actual API responses
$script:ServiceNames = @{
    "Exchange"               = "Exchange Online"
    "SharePoint"             = "SharePoint Online"
    "OneDriveForBusiness"    = "OneDrive for Business"
    "MicrosoftTeams"         = "Microsoft Teams"
    "PowerBICloud"           = "Power BI"
    "Dynamics365"            = "Dynamics 365"
    "OSDPPlatform"           = "Microsoft 365 suite"
    "OrgLiveID"              = "Identity Service"
    "Intune"                 = "Microsoft Intune"
    "PowerAppsM365"          = "Power Apps"
    "PowerAutomateM365"      = "Power Automate"
    "WindowsVirtualDesktop"  = "Azure Virtual Desktop"
    "MicrosoftGraphConnectivity" = "Microsoft Graph"
    "PowerVirtualAgents"     = "Power Virtual Agents"
    "Viva"                   = "Microsoft Viva"
    "UniversalPrint"         = "Universal Print"
    "MicrosoftBookings"      = "Microsoft Bookings"
    "MicrosoftForms"         = "Microsoft Forms"
    "MicrosoftStream"        = "Microsoft Stream"
    "Yammer"                 = "Yammer Enterprise"
    "MicrosoftDefenderforOffice365" = "Microsoft Defender for Office 365"
}

# Function to get service health overview with pagination
function Get-ServiceHealthOverview {
    try {
        Write-Host "📊 Retrieving service health overview..." -ForegroundColor Yellow
        
        $allServices = @()
        $uri = "https://graph.microsoft.com/v1.0/admin/serviceAnnouncement/healthOverviews"
        
        do {
            $response = Invoke-MgGraphRequest -Method GET -Uri $uri
            
            if ($response.value) {
                $allServices += $response.value
            }
            
            $uri = $response.'@odata.nextLink'
        } while ($uri)
        
        Write-Host "✅ Retrieved $($allServices.Count) services" -ForegroundColor Green
        return $allServices
    }
    catch {
        Write-Host "❌ Failed to retrieve service health overview: $_" -ForegroundColor Red
        throw
    }
}

# Function to get current service issues with pagination
function Get-CurrentServiceIssues {
    try {
        Write-Host "🚨 Retrieving current service issues..." -ForegroundColor Yellow
        
        $allIssues = @()
        
        # Get issues from last 7 days that are not resolved
        $sevenDaysAgo = (Get-Date).AddDays(-7).ToString("yyyy-MM-ddTHH:mm:ssZ")
        $filter = "lastModifiedDateTime ge $sevenDaysAgo"
        
        if (-not $IncludeResolvedIssues) {
            $filter += " and status ne 'ServiceRestored' and status ne 'FalsePositive' and status ne 'Resolved' and status ne 'Mitigated' and status ne 'MitigatedExternal' and status ne 'ResolvedExternal' and status ne 'PostIncidentReviewPublished'"
        }
        
        $encodedFilter = [System.Web.HttpUtility]::UrlEncode($filter)
        $uri = "https://graph.microsoft.com/v1.0/admin/serviceAnnouncement/issues?`$filter=$encodedFilter&`$orderby=lastModifiedDateTime desc"
        
        do {
            $response = Invoke-MgGraphRequest -Method GET -Uri $uri
            
            if ($response.value) {
                $allIssues += $response.value
            }
            
            $uri = $response.'@odata.nextLink'
        } while ($uri)
        
        Write-Host "✅ Retrieved $($allIssues.Count) service issues" -ForegroundColor Green
        return $allIssues
    }
    catch {
        Write-Host "❌ Failed to retrieve service issues: $_" -ForegroundColor Red
        throw
    }
}

# Function to get specific issues by ID
function Get-SpecificServiceIssues {
    param (
        [Parameter(Mandatory=$true)]
        [string[]]$IssueIds
    )
    
    try {
        Write-Host "🔍 Retrieving specific service issues..." -ForegroundColor Yellow
        
        $specificIssues = @()
        
        foreach ($issueId in $IssueIds) {
            try {
                $uri = "https://graph.microsoft.com/v1.0/admin/serviceAnnouncement/issues/$issueId"
                $issue = Invoke-MgGraphRequest -Method GET -Uri $uri
                $specificIssues += $issue
                Write-Host "   ✅ Retrieved issue: $issueId" -ForegroundColor Green
            }
            catch {
                Write-Host "   ❌ Failed to retrieve issue: $issueId - $_" -ForegroundColor Red
            }
        }
        
        Write-Host "✅ Retrieved $($specificIssues.Count) of $($IssueIds.Count) requested issues" -ForegroundColor Green
        return $specificIssues
    }
    catch {
        Write-Host "❌ Failed to retrieve specific service issues: $_" -ForegroundColor Red
        throw
    }
}

# Function to format service status with colors and emojis - Based on actual API response format
function Format-ServiceStatus {
    param (
        [Parameter(Mandatory=$true)]
        [string]$Status
    )
    
    $color = $script:StatusColors[$Status]
    if (-not $color) { $color = "White" }
    
    $emoji = switch ($Status) {
        "ServiceOperational"           { "✅" }
        "Investigating"                { "🔍" }
        "RestoringService"             { "🔄" }
        "VerifyingService"             { "🔍" }
        "ServiceRestored"              { "✅" }
        "PostIncidentReviewPublished"  { "📋" }
        "ServiceDegradation"           { "⚠️ " }
        "ServiceInterruption"          { "❌" }
        "ExtendedRecovery"             { "🔄" }
        "FalsePositive"                { "✅" }
        "InvestigationSuspended"       { "⏸️ " }
        "Resolved"                     { "✅" }
        "Mitigated"                    { "✅" }
        "Confirmed"                    { "⚠️ " }
        "Reported"                     { "🔍" }
        "Advisory"                     { "💡" }
        "Incident"                     { "🚨" }
        default                        { "❓" }
    }
    
    $statusText = switch ($Status) {
        "ServiceOperational"           { "Healthy" }
        "Investigating"                { "Investigating" }
        "RestoringService"             { "Restoring Service" }
        "VerifyingService"             { "Verifying Service" }
        "ServiceRestored"              { "Service Restored" }
        "PostIncidentReviewPublished"  { "Post-Incident Review Published" }
        "ServiceDegradation"           { "Service Degradation" }
        "ServiceInterruption"          { "Service Interruption" }
        "ExtendedRecovery"             { "Extended Recovery" }
        "FalsePositive"                { "False Positive" }
        "InvestigationSuspended"       { "Investigation Suspended" }
        "Resolved"                     { "Resolved" }
        "Mitigated"                    { "Mitigated" }
        "MitigatedExternal"            { "Mitigated (External)" }
        "ResolvedExternal"             { "Resolved (External)" }
        "Confirmed"                    { "Confirmed" }
        "Reported"                     { "Reported" }
        "Advisory"                     { "Advisory" }
        "Incident"                     { "Incident" }
        default                        { $Status }
    }
    
    return @{
        Emoji = $emoji
        Text = $statusText
        Color = $color
    }
}

# Function to get friendly service name
function Get-FriendlyServiceName {
    param (
        [Parameter(Mandatory=$true)]
        [string]$ServiceKey
    )
    
    $friendlyName = $script:ServiceNames[$ServiceKey]
    if (-not $friendlyName) { 
        $friendlyName = $ServiceKey -replace "([a-z])([A-Z])", '$1 $2'
    }
    return $friendlyName
}

# Function to display service health dashboard - Focus on outages and issues only
function Show-ServiceHealthDashboard {
    param (
        [Parameter(Mandatory=$true)]
        [array]$HealthOverview,
        
        [Parameter(Mandatory=$true)]
        [array]$Issues,
        
        [Parameter(Mandatory=$false)]
        [array]$WatchedIssues = @()
    )
    
    Write-Host "`n🌐 Microsoft 365 Service Health Dashboard" -ForegroundColor Cyan
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    Write-Host "🕒 Last Updated: $(Get-Date -Format 'yyyy-MM-dd HH:mm:ss UTC')" -ForegroundColor Gray
    
    # Filter services for display - API returns both service IDs and friendly names
    $servicesToShow = $HealthOverview
    if ($ProblemsOnly) {
        $servicesToShow = $HealthOverview | Where-Object { $_.status -ne "ServiceOperational" }
    }
    
    # Service Status Overview
    if ($servicesToShow.Count -gt 0) {
        if ($ProblemsOnly) {
            Write-Host "`n🚨 Services with Issues:" -ForegroundColor Red
        } else {
            Write-Host "`n📊 Service Status Overview:" -ForegroundColor Cyan
        }
        
        $healthyCount = 0
        $degradedCount = 0
        $downCount = 0
        
        foreach ($service in $servicesToShow | Sort-Object service) {
            # Use service name as returned by API (already friendly in healthOverviews)
            $serviceName = $service.service
            $status = Format-ServiceStatus -Status $service.status
            
            Write-Host "   $($status.Emoji) " -NoNewline
            Write-Host "$serviceName" -NoNewline -ForegroundColor White
            Write-Host ": " -NoNewline
            Write-Host "$($status.Text)" -ForegroundColor $status.Color
            
            # Count services based on actual API status values (Pascal case)
            switch ($service.status) {
                "ServiceOperational" { $healthyCount++ }
                { $_ -in @("ServiceDegradation", "Investigating", "RestoringService", "VerifyingService", "ExtendedRecovery", "Confirmed", "Reported") } { $degradedCount++ }
                "ServiceInterruption" { $downCount++ }
            }
        }
        
        # Summary stats for all services (even if filtered view) - Pascal case API values
        if (-not $ProblemsOnly) {
            $totalHealthy = ($HealthOverview | Where-Object { $_.status -eq "ServiceOperational" }).Count
            $totalDegraded = ($HealthOverview | Where-Object { $_.status -in @("ServiceDegradation", "Investigating", "RestoringService", "VerifyingService", "ExtendedRecovery", "Confirmed", "Reported") }).Count  
            $totalDown = ($HealthOverview | Where-Object { $_.status -eq "ServiceInterruption" }).Count
            
            Write-Host "`n📈 Summary:" -ForegroundColor Cyan
            Write-Host "   • " -NoNewline
            Write-Host "Healthy: $totalHealthy" -ForegroundColor Green
            Write-Host "   • " -NoNewline  
            Write-Host "Degraded: $totalDegraded" -ForegroundColor Yellow
            Write-Host "   • " -NoNewline
            Write-Host "Down: $totalDown" -ForegroundColor Red
        }
    }
    else {
        Write-Host "`n✅ All services are healthy" -ForegroundColor Green
    }
    
    # Watched Issues - Show specific issues being monitored
    if ($WatchedIssues.Count -gt 0) {
        Write-Host "`n👁️  Watched Issues ($($WatchedIssues.Count)):" -ForegroundColor Cyan
        
        foreach ($issue in $WatchedIssues | Sort-Object lastModifiedDateTime -Descending) {
            $status = Format-ServiceStatus -Status $issue.status
            $classification = Format-ServiceStatus -Status $issue.classification
            $lastUpdated = [DateTime]::Parse($issue.lastModifiedDateTime).ToString("MM/dd HH:mm")
            
            # Determine if this is an active or resolved issue
            $isActive = $issue.status -notin @("ServiceRestored", "FalsePositive", "Resolved", "Mitigated", "MitigatedExternal", "ResolvedExternal", "PostIncidentReviewPublished")
            $statusIcon = if ($isActive) { "🔴" } else { "🟢" }
            
            Write-Host "   $statusIcon " -NoNewline
            Write-Host "$($issue.id)" -NoNewline -ForegroundColor Yellow
            Write-Host ": " -NoNewline
            Write-Host "$($issue.title)" -ForegroundColor White
            Write-Host "     Status: " -NoNewline
            Write-Host "$($status.Text)" -NoNewline -ForegroundColor $status.Color
            Write-Host " | Type: " -NoNewline
            Write-Host "$($classification.Text)" -NoNewline -ForegroundColor $classification.Color
            Write-Host " | Updated: $lastUpdated" -ForegroundColor Gray
            
            # Show affected service
            if ($issue.service -and $issue.service.Trim() -ne "") {
                Write-Host "     Affected: $($issue.service)" -ForegroundColor Gray
            }
            
            # Show feature information if available
            if ($issue.featureGroup -and $issue.featureGroup.Trim() -ne "") {
                $featureInfo = $issue.featureGroup
                if ($issue.feature -and $issue.feature.Trim() -ne "") {
                    $featureInfo += " - $($issue.feature)"
                }
                Write-Host "     Feature: $featureInfo" -ForegroundColor Gray
            }
            
            # Show latest post if available
            if ($issue.posts -and $issue.posts.Count -gt 0) {
                $latestPost = $issue.posts | Sort-Object createdDateTime -Descending | Select-Object -First 1
                if ($latestPost.description -and $latestPost.description.content) {
                    $postContent = $latestPost.description.content -replace '<[^>]+>', '' # Strip HTML
                    $postContent = $postContent -replace '\n', ' ' # Remove line breaks
                    $trimmedPost = $postContent.Substring(0, [Math]::Min(150, $postContent.Length))
                    if ($postContent.Length -gt 150) { $trimmedPost += "..." }
                    Write-Host "     Latest: $trimmedPost" -ForegroundColor DarkGray
                }
            }
        }
    }
    
    # Active Issues - Based on actual API Pascal case status values
    $activeIssues = $Issues | Where-Object { $_.status -notin @("ServiceRestored", "FalsePositive", "Resolved", "Mitigated", "MitigatedExternal", "ResolvedExternal", "PostIncidentReviewPublished") }
    
    if ($activeIssues.Count -gt 0) {
        Write-Host "`n🚨 Active Incidents ($($activeIssues.Count)):" -ForegroundColor Red
        
        foreach ($issue in $activeIssues | Sort-Object lastModifiedDateTime -Descending) {
            $status = Format-ServiceStatus -Status $issue.status
            $classification = Format-ServiceStatus -Status $issue.classification
            $lastUpdated = [DateTime]::Parse($issue.lastModifiedDateTime).ToString("MM/dd HH:mm")
            
            Write-Host "   • " -NoNewline
            Write-Host "$($issue.id)" -NoNewline -ForegroundColor Yellow
            Write-Host ": " -NoNewline
            Write-Host "$($issue.title)" -ForegroundColor White
            Write-Host "     Status: " -NoNewline
            Write-Host "$($status.Text)" -NoNewline -ForegroundColor $status.Color
            Write-Host " | Type: " -NoNewline
            Write-Host "$($classification.Text)" -NoNewline -ForegroundColor $classification.Color
            Write-Host " | Updated: $lastUpdated" -ForegroundColor Gray
            
            # Show affected service (API returns service as string, not array in issues)
            if ($issue.service -and $issue.service.Trim() -ne "") {
                Write-Host "     Affected: $($issue.service)" -ForegroundColor Gray
            }
            
            # Show feature information if available
            if ($issue.featureGroup -and $issue.featureGroup.Trim() -ne "") {
                $featureInfo = $issue.featureGroup
                if ($issue.feature -and $issue.feature.Trim() -ne "") {
                    $featureInfo += " - $($issue.feature)"
                }
                Write-Host "     Feature: $featureInfo" -ForegroundColor Gray
            }
            
            if ($issue.impactDescription -and $issue.impactDescription.Trim() -ne "") {
                $trimmedImpact = $issue.impactDescription.Substring(0, [Math]::Min(100, $issue.impactDescription.Length))
                if ($issue.impactDescription.Length -gt 100) { $trimmedImpact += "..." }
                Write-Host "     Impact: $trimmedImpact" -ForegroundColor Gray
            }
            
            # Show latest post if available
            if ($issue.posts -and $issue.posts.Count -gt 0) {
                $latestPost = $issue.posts | Sort-Object createdDateTime -Descending | Select-Object -First 1
                if ($latestPost.description -and $latestPost.description.content) {
                    $postContent = $latestPost.description.content -replace '<[^>]+>', '' # Strip HTML
                    $postContent = $postContent -replace '\n', ' ' # Remove line breaks
                    $trimmedPost = $postContent.Substring(0, [Math]::Min(120, $postContent.Length))
                    if ($postContent.Length -gt 120) { $trimmedPost += "..." }
                    Write-Host "     Latest: $trimmedPost" -ForegroundColor DarkGray
                }
            }
        }
    }
    else {
        Write-Host "`n✅ No Active Incidents" -ForegroundColor Green
    }
    
    # Recently Resolved Issues - Pascal case API values
    if ($IncludeResolvedIssues) {
        $resolvedIssues = $Issues | Where-Object { $_.status -in @("ServiceRestored", "FalsePositive", "Resolved", "Mitigated", "MitigatedExternal", "ResolvedExternal", "PostIncidentReviewPublished") }
        
        if ($resolvedIssues.Count -gt 0) {
            Write-Host "`n🔧 Recently Resolved Issues ($($resolvedIssues.Count)):" -ForegroundColor Green
            
            foreach ($issue in $resolvedIssues | Sort-Object lastModifiedDateTime -Descending | Select-Object -First 5) {
                $resolvedTime = [DateTime]::Parse($issue.lastModifiedDateTime).ToString("MM/dd HH:mm")
                $classification = Format-ServiceStatus -Status $issue.classification
                
                Write-Host "   • " -NoNewline
                Write-Host "$($issue.id)" -NoNewline -ForegroundColor Yellow
                Write-Host ": " -NoNewline
                Write-Host "$($issue.title)" -ForegroundColor White
                Write-Host "     Type: " -NoNewline
                Write-Host "$($classification.Text)" -NoNewline -ForegroundColor $classification.Color
                Write-Host " | Resolved: $resolvedTime" -ForegroundColor Gray
            }
        }
    }
}

# Function to export detailed report - Focus on health and issues only
function Export-DetailedReport {
    param (
        [Parameter(Mandatory=$true)]
        [array]$HealthOverview,
        
        [Parameter(Mandatory=$true)]
        [array]$Issues,
        
        [Parameter(Mandatory=$false)]
        [array]$WatchedIssues = @(),
        
        [Parameter(Mandatory=$true)]
        [string]$ExportPath
    )
    
    try {
        Write-Host "`n📄 Exporting detailed report..." -ForegroundColor Yellow
        
        $report = @{
            generatedAt = (Get-Date).ToString("yyyy-MM-ddTHH:mm:ssZ")
            summary = @{
                totalServices = $HealthOverview.Count
                healthyServices = ($HealthOverview | Where-Object { $_.status -eq "ServiceOperational" }).Count
                degradedServices = ($HealthOverview | Where-Object { $_.status -in @("ServiceDegradation", "Investigating", "RestoringService", "VerifyingService", "ExtendedRecovery", "Confirmed", "Reported") }).Count
                downServices = ($HealthOverview | Where-Object { $_.status -eq "ServiceInterruption" }).Count
                activeIncidents = ($Issues | Where-Object { $_.status -notin @("ServiceRestored", "FalsePositive", "Resolved", "Mitigated", "MitigatedExternal", "ResolvedExternal", "PostIncidentReviewPublished") }).Count
                resolvedIncidents = ($Issues | Where-Object { $_.status -in @("ServiceRestored", "FalsePositive", "Resolved", "Mitigated", "MitigatedExternal", "ResolvedExternal", "PostIncidentReviewPublished") }).Count
                watchedIssues = $WatchedIssues.Count
                activeWatchedIssues = ($WatchedIssues | Where-Object { $_.status -notin @("ServiceRestored", "FalsePositive", "Resolved", "Mitigated", "MitigatedExternal", "ResolvedExternal", "PostIncidentReviewPublished") }).Count
            }
            serviceHealth = $HealthOverview
            issues = $Issues
            watchedIssues = $WatchedIssues
        }
        
        $report | ConvertTo-Json -Depth 10 | Out-File -FilePath $ExportPath -Encoding UTF8
        
        Write-Host "✅ Report exported to: $ExportPath" -ForegroundColor Green
    }
    catch {
        Write-Host "❌ Failed to export report: $_" -ForegroundColor Red
    }
}

# Script Setup - Same pattern as Intune script
Import-Module Microsoft.Graph.Authentication

$secureClientSecret = ConvertTo-SecureString -String $ClientSecret -AsPlainText -Force
$clientSecretCredential = New-Object -TypeName System.Management.Automation.PSCredential -ArgumentList $ClientId, $secureClientSecret

Write-Host "Connecting to Microsoft Graph..." -ForegroundColor Cyan
try {
    Connect-MgGraph -ClientSecretCredential $clientSecretCredential -TenantId $TenantId -NoWelcome
    Write-Host "✅ Connected to Microsoft Graph successfully" -ForegroundColor Green
}
catch {
    Write-Host "❌ Failed to connect to Microsoft Graph: $_" -ForegroundColor Red
    exit 2
}

# Test permissions before proceeding
function Test-ServiceHealthPermissions {
    try {
        Write-Host "`n🔐 Testing Service Health API permissions..." -ForegroundColor Yellow
        
        # Test the basic endpoint with minimal request
        $testUri = "https://graph.microsoft.com/v1.0/admin/serviceAnnouncement/healthOverviews?`$top=1"
        $testResponse = Invoke-MgGraphRequest -Method GET -Uri $testUri
        
        Write-Host "✅ Service Health permissions verified" -ForegroundColor Green
        return $true
    }
    catch {
        $errorMessage = $_.Exception.Message
        $statusCode = "Unknown"
        
        if ($errorMessage -match "HTTP/\d+\.\d+ (\d+)") {
            $statusCode = $matches[1]
        }
        
        Write-Host "❌ Service Health permission test failed" -ForegroundColor Red
        Write-Host ""
        
        switch ($statusCode) {
            "403" {
                Write-Host "🔒 PERMISSION ISSUE: Application lacks required permissions" -ForegroundColor Red
                Write-Host ""
                Write-Host "Required Azure AD Application Permission:" -ForegroundColor Yellow
                Write-Host "   • ServiceHealth.Read.All" -ForegroundColor White
                Write-Host ""
                Write-Host "To fix this issue:" -ForegroundColor Yellow
                Write-Host "   1. Go to Azure Portal > App Registrations" -ForegroundColor White
                Write-Host "   2. Find your application: $ClientId" -ForegroundColor White
                Write-Host "   3. Go to API Permissions" -ForegroundColor White
                Write-Host "   4. Add Permission > Microsoft Graph > Application permissions" -ForegroundColor White
                Write-Host "   5. Search for and add: ServiceHealth.Read.All" -ForegroundColor White
                Write-Host "   6. Click 'Grant admin consent' (requires Global Admin)" -ForegroundColor White
                Write-Host ""
                Write-Host "Note: Admin consent is REQUIRED for Service Health permissions" -ForegroundColor Red
            }
            "401" {
                Write-Host "🔑 AUTHENTICATION ISSUE: Invalid credentials" -ForegroundColor Red
                Write-Host "   • Verify TenantId: $TenantId" -ForegroundColor White
                Write-Host "   • Verify ClientId: $ClientId" -ForegroundColor White
                Write-Host "   • Verify ClientSecret is valid and not expired" -ForegroundColor White
            }
            "404" {
                Write-Host "🌐 ENDPOINT ISSUE: Service Health API not available" -ForegroundColor Red
                Write-Host "   • Verify tenant supports Service Health API" -ForegroundColor White
                Write-Host "   • Check if tenant has appropriate license" -ForegroundColor White
            }
            default {
                Write-Host "❓ UNKNOWN ISSUE (HTTP $statusCode)" -ForegroundColor Red
                Write-Host "   • Check Azure AD application configuration" -ForegroundColor White
                Write-Host "   • Verify tenant and application settings" -ForegroundColor White
            }
        }
        
        Write-Host ""
        Write-Host "Full error details:" -ForegroundColor Gray
        Write-Host $errorMessage -ForegroundColor Gray
        
        return $false
    }
}

# Main execution function - Focus on platform health and availability only
function Get-M365ServiceHealth {
    try {
        Write-Host "`n🔍 Starting Microsoft 365 Service Health Check..." -ForegroundColor Cyan
        Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
        
        # Get service health data with progress indicators
        $healthOverview = Get-ServiceHealthOverview
        $issues = Get-CurrentServiceIssues
        
        # Get specific watched issues if requested
        $watchedIssues = @()
        if ($WatchIssues -and $WatchIssues.Count -gt 0) {
            $watchedIssues = Get-SpecificServiceIssues -IssueIds $WatchIssues
        }
        
        # Display dashboard
        Show-ServiceHealthDashboard -HealthOverview $healthOverview -Issues $issues -WatchedIssues $watchedIssues
        
        # Export if requested
        if ($ExportPath) {
            Export-DetailedReport -HealthOverview $healthOverview -Issues $issues -WatchedIssues $watchedIssues -ExportPath $ExportPath
        }
        
        Write-Host "`n🎉 Service health check completed successfully!" -ForegroundColor Green
        
        # Return summary for programmatic use - Pascal case API values
        $activeIncidents = ($issues | Where-Object { $_.status -notin @("ServiceRestored", "FalsePositive", "Resolved", "Mitigated", "MitigatedExternal", "ResolvedExternal", "PostIncidentReviewPublished") }).Count
        $degradedServices = ($healthOverview | Where-Object { $_.status -in @("ServiceDegradation", "ServiceInterruption", "Investigating", "RestoringService", "VerifyingService", "ExtendedRecovery", "Confirmed", "Reported") }).Count
        
        # Include watched issues status
        $activeWatchedIssues = ($watchedIssues | Where-Object { $_.status -notin @("ServiceRestored", "FalsePositive", "Resolved", "Mitigated", "MitigatedExternal", "ResolvedExternal", "PostIncidentReviewPublished") }).Count
        
        return @{
            HasIssues = ($activeIncidents -gt 0 -or $degradedServices -gt 0 -or $activeWatchedIssues -gt 0)
            ActiveIncidents = $activeIncidents
            DegradedServices = $degradedServices
            TotalServices = $healthOverview.Count
            HealthyServices = ($healthOverview | Where-Object { $_.status -eq "ServiceOperational" }).Count
            WatchedIssues = $watchedIssues.Count
            ActiveWatchedIssues = $activeWatchedIssues
        }
    }
    catch {
        Write-Host "`n❌ Service health check failed: $_" -ForegroundColor Red
        throw
    }
}

# Execute the main function - Same pattern as Intune script
try {
    # Test permissions first
    if (-not (Test-ServiceHealthPermissions)) {
        Write-Host "`n❌ Cannot proceed due to permission issues" -ForegroundColor Red
        exit 2
    }
    
    $result = Get-M365ServiceHealth
    
    # Set exit code based on service health
    if ($result.HasIssues) {
        Write-Host "`n⚠️ Issues detected in Microsoft 365 services" -ForegroundColor Yellow
        Write-Host "   • Active Incidents: $($result.ActiveIncidents)" -ForegroundColor Yellow
        Write-Host "   • Degraded Services: $($result.DegradedServices)" -ForegroundColor Yellow
        if ($result.WatchedIssues -gt 0) {
            Write-Host "   • Watched Issues: $($result.ActiveWatchedIssues)/$($result.WatchedIssues) active" -ForegroundColor Yellow
        }
        exit 1
    }
    else {
        $statusMessage = "All Microsoft 365 services are healthy ($($result.HealthyServices)/$($result.TotalServices))"
        if ($result.WatchedIssues -gt 0) {
            $statusMessage += " | Watched Issues: $($result.WatchedIssues) (all resolved)"
        }
        Write-Host "`n✅ $statusMessage" -ForegroundColor Green
        exit 0
    }
}
catch {
    Write-Host "`n❌ Script execution failed: $_" -ForegroundColor Red
    exit 2
}

# Disconnect from Microsoft Graph - Same pattern as Intune script
Write-Host "`nDisconnecting from Microsoft Graph..." -ForegroundColor Cyan
Disconnect-MgGraph > $null 2>&1
Write-Host "✅ Disconnected from Microsoft Graph" -ForegroundColor Green