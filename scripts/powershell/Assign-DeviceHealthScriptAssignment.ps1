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
    
    [Parameter(Mandatory=$true,
    HelpMessage="Device Health Script ID")]
    [ValidateNotNullOrEmpty()]
    [string]$DeviceHealthScriptId,
    
    [Parameter(Mandatory=$false,
    HelpMessage="Specific Assignment ID (if not provided, will list all assignments)")]
    [string]$AssignmentId
)

# Import required modules
Import-Module Microsoft.Graph.Authentication

# Function to get device health script assignments
function Get-DeviceHealthScriptAssignments {
    param (
        [Parameter(Mandatory=$true)]
        [string]$ScriptId,
        
        [Parameter(Mandatory=$false)]
        [string]$SpecificAssignmentId
    )
    
    try {
        if ($SpecificAssignmentId) {
            # GET specific assignment
            $uri = "https://graph.microsoft.com/beta/deviceManagement/deviceHealthScripts/$ScriptId/assignments/$SpecificAssignmentId"
            Write-Host "🔍 Getting specific assignment..." -ForegroundColor Cyan
            Write-Host "   Script ID: $ScriptId" -ForegroundColor Gray
            Write-Host "   Assignment ID: $SpecificAssignmentId" -ForegroundColor Gray
        } else {
            # GET all assignments for the script
            $uri = "https://graph.microsoft.com/beta/deviceManagement/deviceHealthScripts/$ScriptId/assignments"
            Write-Host "🔍 Getting all assignments for script..." -ForegroundColor Cyan
            Write-Host "   Script ID: $ScriptId" -ForegroundColor Gray
        }
        
        Write-Host "   Endpoint: $uri" -ForegroundColor Gray
        Write-Host ""
        
        $response = Invoke-MgGraphRequest -Method GET -Uri $uri
        
        return $response
    }
    catch {
        Write-Host "❌ Error getting device health script assignments: $_" -ForegroundColor Red
        Write-Host ""
        
        # Enhanced error handling
        if ($_.Exception.Response) {
            $statusCode = $_.Exception.Response.StatusCode
            $statusDescription = $_.Exception.Response.StatusDescription
            Write-Host "   Status Code: $statusCode" -ForegroundColor Red
            Write-Host "   Status Description: $statusDescription" -ForegroundColor Red
            
            # Try to get the response content
            try {
                $responseContent = $_.Exception.Response.GetResponseStream()
                if ($responseContent) {
                    $reader = [System.IO.StreamReader]::new($responseContent)
                    $errorDetails = $reader.ReadToEnd()
                    $reader.Close()
                    Write-Host "   Error Details: $errorDetails" -ForegroundColor Red
                }
            }
            catch {
                Write-Host "   Could not read error details" -ForegroundColor Red
            }
        }
        
        throw
    }
}

# Function to display assignment details
function Show-AssignmentDetails {
    param (
        [Parameter(Mandatory=$true)]
        $Assignment
    )
    
    Write-Host "📋 Assignment Details:" -ForegroundColor Cyan
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    
    if ($Assignment.id) {
        Write-Host "   • ID: $($Assignment.id)" -ForegroundColor Green
    }
    
    if ($Assignment.runRemediationScript -ne $null) {
        Write-Host "   • Run Remediation Script: $($Assignment.runRemediationScript)" -ForegroundColor Green
    }
    
    if ($Assignment.target) {
        Write-Host "   • Target:" -ForegroundColor Green
        $target = $Assignment.target
        
        if ($target.'@odata.type') {
            $targetType = $target.'@odata.type' -replace '#microsoft.graph.', ''
            Write-Host "     - Type: $targetType" -ForegroundColor Yellow
            
            switch ($targetType) {
                "groupAssignmentTarget" {
                    if ($target.groupId) {
                        Write-Host "     - Group ID: $($target.groupId)" -ForegroundColor Yellow
                    }
                }
                "exclusionGroupAssignmentTarget" {
                    if ($target.groupId) {
                        Write-Host "     - Exclusion Group ID: $($target.groupId)" -ForegroundColor Yellow
                    }
                }
                "configurationManagerCollectionAssignmentTarget" {
                    if ($target.collectionId) {
                        Write-Host "     - Collection ID: $($target.collectionId)" -ForegroundColor Yellow
                    }
                }
                "allDevicesAssignmentTarget" {
                    Write-Host "     - Targets: All Devices" -ForegroundColor Yellow
                }
                "allLicensedUsersAssignmentTarget" {
                    Write-Host "     - Targets: All Licensed Users" -ForegroundColor Yellow
                }
            }
        }
        
        if ($target.deviceAndAppManagementAssignmentFilterId) {
            Write-Host "     - Filter ID: $($target.deviceAndAppManagementAssignmentFilterId)" -ForegroundColor Yellow
        }
        
        if ($target.deviceAndAppManagementAssignmentFilterType) {
            Write-Host "     - Filter Type: $($target.deviceAndAppManagementAssignmentFilterType)" -ForegroundColor Yellow
        }
    }
    
    if ($Assignment.runSchedule) {
        Write-Host "   • Run Schedule:" -ForegroundColor Green
        $schedule = $Assignment.runSchedule
        
        if ($schedule.'@odata.type') {
            $scheduleType = $schedule.'@odata.type' -replace '#microsoft.graph.', ''
            Write-Host "     - Type: $scheduleType" -ForegroundColor Yellow
            
            switch ($scheduleType) {
                "deviceHealthScriptDailySchedule" {
                    if ($schedule.interval) {
                        Write-Host "     - Interval: $($schedule.interval) days" -ForegroundColor Yellow
                    }
                    if ($schedule.time) {
                        Write-Host "     - Time: $($schedule.time)" -ForegroundColor Yellow
                    }
                    if ($schedule.useUtc -ne $null) {
                        Write-Host "     - Use UTC: $($schedule.useUtc)" -ForegroundColor Yellow
                    }
                }
                "deviceHealthScriptHourlySchedule" {
                    if ($schedule.interval) {
                        Write-Host "     - Interval: $($schedule.interval) hours" -ForegroundColor Yellow
                    }
                }
                "deviceHealthScriptRunOnceSchedule" {
                    if ($schedule.date) {
                        Write-Host "     - Date: $($schedule.date)" -ForegroundColor Yellow
                    }
                    if ($schedule.time) {
                        Write-Host "     - Time: $($schedule.time)" -ForegroundColor Yellow
                    }
                    if ($schedule.useUtc -ne $null) {
                        Write-Host "     - Use UTC: $($schedule.useUtc)" -ForegroundColor Yellow
                    }
                }
            }
        }
    }
    
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    Write-Host ""
}

# Script Setup
try {
    Write-Host "🔐 Connecting to Microsoft Graph..." -ForegroundColor Cyan
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    
    $secureClientSecret = ConvertTo-SecureString -String $ClientSecret -AsPlainText -Force
    $clientSecretCredential = New-Object -TypeName System.Management.Automation.PSCredential -ArgumentList $ClientId, $secureClientSecret
    
    Connect-MgGraph -ClientSecretCredential $clientSecretCredential -TenantId $TenantId
    
    Write-Host "✅ Connected to Microsoft Graph successfully" -ForegroundColor Green
    Write-Host ""
    
    # Get the assignments
    $assignments = Get-DeviceHealthScriptAssignments -ScriptId $DeviceHealthScriptId -SpecificAssignmentId $AssignmentId
    
    if ($AssignmentId) {
        # Display single assignment
        Show-AssignmentDetails -Assignment $assignments
    } else {
        # Display all assignments
        if ($assignments.value -and $assignments.value.Count -gt 0) {
            Write-Host "📊 Found $($assignments.value.Count) assignment(s)" -ForegroundColor Green
            Write-Host ""
            
            for ($i = 0; $i -lt $assignments.value.Count; $i++) {
                Write-Host "Assignment $($i + 1):" -ForegroundColor Magenta
                Show-AssignmentDetails -Assignment $assignments.value[$i]
            }
        } elseif ($assignments -and -not $assignments.value) {
            # Single assignment returned (not in a collection)
            Write-Host "📊 Found 1 assignment" -ForegroundColor Green
            Write-Host ""
            Show-AssignmentDetails -Assignment $assignments
        } else {
            Write-Host "📊 No assignments found for this device health script" -ForegroundColor Yellow
        }
    }
    
    Write-Host "🎉 Operation completed successfully!" -ForegroundColor Green
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