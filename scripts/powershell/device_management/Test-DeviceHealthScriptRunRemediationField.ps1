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
    HelpMessage="Cleanup resources after testing")]
    [switch]$Cleanup
)

# Script Setup
$ErrorActionPreference = "Stop"
$timestamp = Get-Date -Format "yyyyMMdd-HHmmss"
$testPrefix = "ps-test-runremediation-$timestamp"

Write-Host "=============================================" -ForegroundColor Cyan
Write-Host "Device Health Script runRemediationScript Test" -ForegroundColor Cyan
Write-Host "=============================================" -ForegroundColor Cyan
Write-Host ""

# Connect to Microsoft Graph
$secureClientSecret = ConvertTo-SecureString -String $ClientSecret -AsPlainText -Force
$clientSecretCredential = New-Object -TypeName System.Management.Automation.PSCredential -ArgumentList $ClientId, $secureClientSecret

Write-Host "🔐 Connecting to Microsoft Graph..." -ForegroundColor Yellow
try {
    Connect-MgGraph -ClientSecretCredential $clientSecretCredential -TenantId $TenantId -NoWelcome
    Write-Host "✅ Connected successfully" -ForegroundColor Green
    Write-Host ""
}
catch {
    Write-Error "❌ Failed to connect to Microsoft Graph: $_"
    exit 1
}

# Variables to track created resources
$createdGroupId = $null
$createdScriptId = $null
$createdAssignmentId = $null

try {
    # ============================================
    # Step 1: Create Test Group
    # ============================================
    Write-Host "📋 Step 1: Creating test group..." -ForegroundColor Yellow
    
    $groupParams = @{
        displayName = "$testPrefix-group"
        mailNickname = $testPrefix
        mailEnabled = $false
        securityEnabled = $true
        description = "Test group for runRemediationScript field validation"
    }
    
    $groupJson = $groupParams | ConvertTo-Json -Depth 10
    $group = Invoke-MgGraphRequest `
        -Method POST `
        -Uri "https://graph.microsoft.com/beta/groups" `
        -Body $groupJson `
        -ContentType "application/json"
    $createdGroupId = $group.Id
    
    Write-Host "✅ Group created successfully" -ForegroundColor Green
    Write-Host "   Group ID: $createdGroupId" -ForegroundColor Gray
    Write-Host "   Display Name: $($group.DisplayName)" -ForegroundColor Gray
    Write-Host ""

    # ============================================
    # Step 2: Create Device Health Script
    # ============================================
    Write-Host "📝 Step 2: Creating device health script..." -ForegroundColor Yellow
    
    # Simple detection script (base64 encoded)
    $detectionScript = "Write-Host 'Detection complete'; exit 0"
    $detectionScriptBase64 = [Convert]::ToBase64String([System.Text.Encoding]::UTF8.GetBytes($detectionScript))
    
    # Simple remediation script (base64 encoded)
    $remediationScript = "Write-Host 'Remediation complete'; exit 0"
    $remediationScriptBase64 = [Convert]::ToBase64String([System.Text.Encoding]::UTF8.GetBytes($remediationScript))
    
    $scriptParams = @{
        displayName = "$testPrefix-script"
        description = "Test script for runRemediationScript field validation"
        publisher = "PowerShell Test Script"
        runAsAccount = "system"
        enforceSignatureCheck = $false
        runAs32Bit = $false
        detectionScriptContent = $detectionScriptBase64
        remediationScriptContent = $remediationScriptBase64
        roleScopeTagIds = @("0")
    }
    
    $scriptJson = $scriptParams | ConvertTo-Json -Depth 10
    
    $response = Invoke-MgGraphRequest `
        -Method POST `
        -Uri "https://graph.microsoft.com/beta/deviceManagement/deviceHealthScripts" `
        -Body $scriptJson `
        -ContentType "application/json"
    
    $createdScriptId = $response.id
    
    Write-Host "✅ Device health script created successfully" -ForegroundColor Green
    Write-Host "   Script ID: $createdScriptId" -ForegroundColor Gray
    Write-Host "   Display Name: $($response.displayName)" -ForegroundColor Gray
    Write-Host ""

    # ============================================
    # Step 3: Create Assignment with runRemediationScript = true
    # ============================================
    Write-Host "🎯 Step 3: Creating assignment with runRemediationScript = true..." -ForegroundColor Yellow
    
    $assignmentParams = @{
        deviceHealthScriptAssignments = @(
            @{
                target = @{
                    "@odata.type" = "#microsoft.graph.groupAssignmentTarget"
                    groupId = $createdGroupId
                }
                runRemediationScript = $true
                runSchedule = @{
                    "@odata.type" = "#microsoft.graph.deviceHealthScriptDailySchedule"
                    interval = 1
                    time = "09:00:00.0000000"
                    useUtc = $true
                }
            }
        )
    }
    
    $assignmentJson = $assignmentParams | ConvertTo-Json -Depth 10
    
    Write-Host "📤 Sending assignment request with payload:" -ForegroundColor Cyan
    Write-Host $assignmentJson -ForegroundColor Gray
    Write-Host ""
    
    $assignResponse = Invoke-MgGraphRequest `
        -Method POST `
        -Uri "https://graph.microsoft.com/beta/deviceManagement/deviceHealthScripts/$createdScriptId/assign" `
        -Body $assignmentJson `
        -ContentType "application/json"
    
    Write-Host "✅ Assignment created successfully" -ForegroundColor Green
    Write-Host ""

    # ============================================
    # Step 4: Read back the script with assignments expanded
    # ============================================
    Write-Host "🔍 Step 4: Reading back script with assignments to verify runRemediationScript field..." -ForegroundColor Yellow
    
    Start-Sleep -Seconds 2  # Give API time to process
    
    $scriptWithAssignments = Invoke-MgGraphRequest `
        -Method GET `
        -Uri "https://graph.microsoft.com/beta/deviceManagement/deviceHealthScripts/${createdScriptId}?`$expand=assignments" `
        -ContentType "application/json"
    
    Write-Host "📥 Retrieved script with assignments" -ForegroundColor Cyan
    Write-Host ""
    
    # ============================================
    # Step 5: Analyze the runRemediationScript field
    # ============================================
    Write-Host "=============================================" -ForegroundColor Cyan
    Write-Host "🔬 ANALYSIS RESULTS" -ForegroundColor Cyan
    Write-Host "=============================================" -ForegroundColor Cyan
    Write-Host ""
    
    if ($scriptWithAssignments.assignments -and $scriptWithAssignments.assignments.Count -gt 0) {
        $assignment = $scriptWithAssignments.assignments[0]
        $createdAssignmentId = $assignment.id
        
        Write-Host "Assignment Details:" -ForegroundColor White
        Write-Host "  Assignment ID: $($assignment.id)" -ForegroundColor Gray
        Write-Host "  Target Type: $($assignment.target.'@odata.type')" -ForegroundColor Gray
        Write-Host "  Group ID: $($assignment.target.groupId)" -ForegroundColor Gray
        Write-Host ""
        
        Write-Host "runRemediationScript Field Test:" -ForegroundColor White
        
        if ($null -ne $assignment.PSObject.Properties['runRemediationScript']) {
            $returnedValue = $assignment.runRemediationScript
            Write-Host "  ✅ Field EXISTS in response" -ForegroundColor Green
            Write-Host "  📊 Sent Value: TRUE" -ForegroundColor Yellow
            Write-Host "  📊 Returned Value: $returnedValue" -ForegroundColor Yellow
            Write-Host ""
            
            if ($returnedValue -eq $true) {
                Write-Host "  ✅ SUCCESS: Field value matches (true = true)" -ForegroundColor Green
                Write-Host "  The API correctly persists and returns the runRemediationScript field!" -ForegroundColor Green
            }
            elseif ($returnedValue -eq $false) {
                Write-Host "  ❌ FAILURE: Field value mismatch (sent: true, received: false)" -ForegroundColor Red
                Write-Host "  The API does NOT correctly persist the runRemediationScript field!" -ForegroundColor Red
                Write-Host "  This confirms the field is broken or ignored by the API." -ForegroundColor Red
            }
            else {
                Write-Host "  ⚠️  WARNING: Field has unexpected value: $returnedValue" -ForegroundColor Yellow
            }
        }
        else {
            Write-Host "  ❌ Field DOES NOT EXIST in response" -ForegroundColor Red
            Write-Host "  The API does not return the runRemediationScript field at all!" -ForegroundColor Red
        }
        
        Write-Host ""
        Write-Host "Full Assignment JSON:" -ForegroundColor White
        Write-Host ($assignment | ConvertTo-Json -Depth 10) -ForegroundColor Gray
    }
    else {
        Write-Host "❌ No assignments found in response!" -ForegroundColor Red
    }
    
    Write-Host ""
    Write-Host "=============================================" -ForegroundColor Cyan
    
    # ============================================
    # Step 6: Get specific assignment by ID
    # ============================================
    Write-Host ""
    Write-Host "🔍 Step 6: Getting specific assignment by ID..." -ForegroundColor Yellow
    
    if ($createdAssignmentId) {
        try {
            $specificAssignment = Invoke-MgGraphRequest `
                -Method GET `
                -Uri "https://graph.microsoft.com/beta/deviceManagement/deviceHealthScripts/${createdScriptId}/assignments/${createdAssignmentId}" `
                -ContentType "application/json"
            
            Write-Host "📥 Retrieved specific assignment" -ForegroundColor Cyan
            Write-Host ""
            Write-Host "Specific Assignment Endpoint Test:" -ForegroundColor White
            
            if ($null -ne $specificAssignment.PSObject.Properties['runRemediationScript']) {
                $returnedValue3 = $specificAssignment.runRemediationScript
                Write-Host "  ✅ Field EXISTS in specific assignment response" -ForegroundColor Green
                Write-Host "  📊 Returned Value: $returnedValue3" -ForegroundColor Yellow
                
                if ($returnedValue3 -eq $true) {
                    Write-Host "  ✅ Specific endpoint returns: true" -ForegroundColor Green
                }
                else {
                    Write-Host "  ❌ Specific endpoint returns: $returnedValue3" -ForegroundColor Red
                }
            }
            else {
                Write-Host "  ❌ Field DOES NOT EXIST in specific assignment response" -ForegroundColor Red
            }
            
            Write-Host ""
            Write-Host "Full Specific Assignment JSON:" -ForegroundColor White
            Write-Host ($specificAssignment | ConvertTo-Json -Depth 10) -ForegroundColor Gray
            Write-Host ""
        }
        catch {
            Write-Host "  ⚠️  Failed to get specific assignment: $_" -ForegroundColor Yellow
            Write-Host ""
        }
    }
    
    # ============================================
    # Step 7: Test with runRemediationScript = false
    # ============================================
    Write-Host ""
    Write-Host "🎯 Step 7: Testing with runRemediationScript = false..." -ForegroundColor Yellow
    
    $assignmentParams2 = @{
        deviceHealthScriptAssignments = @(
            @{
                target = @{
                    "@odata.type" = "#microsoft.graph.groupAssignmentTarget"
                    groupId = $createdGroupId
                }
                runRemediationScript = $false
                runSchedule = @{
                    "@odata.type" = "#microsoft.graph.deviceHealthScriptHourlySchedule"
                    interval = 4
                }
            }
        )
    }
    
    $assignmentJson2 = $assignmentParams2 | ConvertTo-Json -Depth 10
    
    $assignResponse2 = Invoke-MgGraphRequest `
        -Method POST `
        -Uri "https://graph.microsoft.com/beta/deviceManagement/deviceHealthScripts/$createdScriptId/assign" `
        -Body $assignmentJson2 `
        -ContentType "application/json"
    
    Start-Sleep -Seconds 2
    
    $scriptWithAssignments2 = Invoke-MgGraphRequest `
        -Method GET `
        -Uri "https://graph.microsoft.com/beta/deviceManagement/deviceHealthScripts/${createdScriptId}?`$expand=assignments" `
        -ContentType "application/json"
    
    if ($scriptWithAssignments2.assignments -and $scriptWithAssignments2.assignments.Count -gt 0) {
        $assignment2 = $scriptWithAssignments2.assignments[0]
        
        Write-Host "runRemediationScript Field Test (false):" -ForegroundColor White
        if ($null -ne $assignment2.PSObject.Properties['runRemediationScript']) {
            $returnedValue2 = $assignment2.runRemediationScript
            Write-Host "  📊 Sent Value: FALSE" -ForegroundColor Yellow
            Write-Host "  📊 Returned Value: $returnedValue2" -ForegroundColor Yellow
            
            if ($returnedValue2 -eq $false) {
                Write-Host "  ✅ Value matches (false = false)" -ForegroundColor Green
            }
            else {
                Write-Host "  ⚠️  Value mismatch (sent: false, received: $returnedValue2)" -ForegroundColor Yellow
            }
        }
    }
    
    Write-Host ""
}
catch {
    Write-Host ""
    Write-Host "❌ ERROR OCCURRED:" -ForegroundColor Red
    Write-Host $_.Exception.Message -ForegroundColor Red
    Write-Host ""
    Write-Host "Full Error Details:" -ForegroundColor Yellow
    Write-Host $_ -ForegroundColor Gray
}
finally {
    # ============================================
    # Cleanup
    # ============================================
    if ($Cleanup) {
        Write-Host ""
        Write-Host "🧹 Cleanup: Removing created resources..." -ForegroundColor Yellow
        Write-Host ""
        
        # Delete Device Health Script (this also deletes assignments)
        if ($createdScriptId) {
            try {
                Write-Host "  Deleting device health script: $createdScriptId" -ForegroundColor Gray
                Invoke-MgGraphRequest `
                    -Method DELETE `
                    -Uri "https://graph.microsoft.com/beta/deviceManagement/deviceHealthScripts/$createdScriptId"
                Write-Host "  ✅ Script deleted" -ForegroundColor Green
            }
            catch {
                Write-Host "  ⚠️  Failed to delete script: $_" -ForegroundColor Yellow
            }
        }
        
        # Delete Group
        if ($createdGroupId) {
            try {
                Write-Host "  Deleting group: $createdGroupId" -ForegroundColor Gray
                Invoke-MgGraphRequest `
                    -Method DELETE `
                    -Uri "https://graph.microsoft.com/beta/groups/$createdGroupId"
                Write-Host "  ✅ Group deleted" -ForegroundColor Green
            }
            catch {
                Write-Host "  ⚠️  Failed to delete group: $_" -ForegroundColor Yellow
            }
        }
        
        Write-Host ""
        Write-Host "✅ Cleanup completed" -ForegroundColor Green
    }
    else {
        Write-Host ""
        Write-Host "ℹ️  Resources NOT cleaned up (use -Cleanup flag to remove)" -ForegroundColor Cyan
        Write-Host ""
        Write-Host "Created Resources:" -ForegroundColor White
        if ($createdGroupId) {
            Write-Host "  Group ID: $createdGroupId" -ForegroundColor Gray
        }
        if ($createdScriptId) {
            Write-Host "  Script ID: $createdScriptId" -ForegroundColor Gray
        }
        if ($createdAssignmentId) {
            Write-Host "  Assignment ID: $createdAssignmentId" -ForegroundColor Gray
        }
    }
    
    # Disconnect from Microsoft Graph
    Write-Host ""
    Write-Host "🔌 Disconnecting from Microsoft Graph..." -ForegroundColor Yellow
    Disconnect-MgGraph | Out-Null
    Write-Host "✅ Disconnected" -ForegroundColor Green
    Write-Host ""
    Write-Host "=============================================" -ForegroundColor Cyan
    Write-Host "Test Complete" -ForegroundColor Cyan
    Write-Host "=============================================" -ForegroundColor Cyan
}

