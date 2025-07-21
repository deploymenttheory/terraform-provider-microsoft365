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
    HelpMessage="Optional path to export the device health script details to a JSON file")]
    [string]$ExportPath,
    
    [Parameter(Mandatory=$false,
    HelpMessage="Show decoded script content in the console output")]
    [switch]$ShowScriptContent
)

# Import required modules
Import-Module Microsoft.Graph.Authentication

# Function to get device health script details
function Get-DeviceHealthScript {
    param (
        [Parameter(Mandatory=$true)]
        [string]$ScriptId
    )
    
    try {
        # GET device health script with assignments
        $uri = "https://graph.microsoft.com/beta/deviceManagement/deviceHealthScripts('$ScriptId')?`$expand=assignments"
        Write-Host "🔍 Getting device health script..." -ForegroundColor Cyan
        Write-Host "   Script ID: $ScriptId" -ForegroundColor Gray
        Write-Host "   Endpoint: $uri" -ForegroundColor Gray
        Write-Host ""
        
        $response = Invoke-MgGraphRequest -Method GET -Uri $uri
        
        return $response
    }
    catch {
        Write-Host "❌ Error getting device health script: $_" -ForegroundColor Red
        Write-Host ""
        
        # Enhanced error handling
        if ($_.Exception.Response) {
            $statusCode = $_.Exception.Response.StatusCode
            $statusDescription = $_.Exception.Response.StatusDescription
            Write-Host "   Status Code: $statusCode" -ForegroundColor Red
            Write-Host "   Status Description: $statusDescription" -ForegroundColor Red
            
            # Try to parse the error response JSON
            try {
                $errorResponse = $_.ErrorDetails.Message | ConvertFrom-Json
                if ($errorResponse.error) {
                    Write-Host "   Error Code: $($errorResponse.error.code)" -ForegroundColor Red
                    Write-Host "   Error Message: $($errorResponse.error.message)" -ForegroundColor Red
                    
                    # Handle specific error scenarios
                    if ($errorResponse.error.message -like "*tomb-stoned*") {
                        Write-Host "" -ForegroundColor Red
                        Write-Host "   ℹ️  This device health script appears to be deleted or archived." -ForegroundColor Yellow
                        Write-Host "   ℹ️  The script ID exists in the system but is no longer active." -ForegroundColor Yellow
                    }
                    
                    if ($errorResponse.error.innerError) {
                        Write-Host "   Request ID: $($errorResponse.error.innerError.'request-id')" -ForegroundColor Red
                    }
                }
            }
            catch {
                # Fall back to trying to get the response content the old way
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
        }
        
        throw
    }
}

# Function to display device health script details
function Show-DeviceHealthScriptDetails {
    param (
        [Parameter(Mandatory=$true)]
        $Script
    )
    
    Write-Host "📋 Device Health Script Details:" -ForegroundColor Cyan
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    
    if ($Script.id) {
        Write-Host "   • ID: $($Script.id)" -ForegroundColor Green
    }
    
    if ($Script.displayName) {
        Write-Host "   • Display Name: $($Script.displayName)" -ForegroundColor Green
    }
    
    if ($Script.description) {
        Write-Host "   • Description: $($Script.description)" -ForegroundColor Green
    }
    
    if ($Script.publisher) {
        Write-Host "   • Publisher: $($Script.publisher)" -ForegroundColor Green
    }
    
    if ($Script.version) {
        Write-Host "   • Version: $($Script.version)" -ForegroundColor Green
    }
    
    if ($Script.createdDateTime) {
        Write-Host "   • Created: $($Script.createdDateTime)" -ForegroundColor Green
    }
    
    if ($Script.lastModifiedDateTime) {
        Write-Host "   • Last Modified: $($Script.lastModifiedDateTime)" -ForegroundColor Green
    }
    
    if ($Script.runAsAccount) {
        Write-Host "   • Run As Account: $($Script.runAsAccount)" -ForegroundColor Green
    }
    
    if ($Script.enforceSignatureCheck -ne $null) {
        Write-Host "   • Enforce Signature Check: $($Script.enforceSignatureCheck)" -ForegroundColor Green
    }
    
    if ($Script.runAs32Bit -ne $null) {
        Write-Host "   • Run As 32-bit: $($Script.runAs32Bit)" -ForegroundColor Green
    }
    
    if ($Script.detectionScriptContent) {
        try {
            $decodedDetectionScript = [System.Text.Encoding]::UTF8.GetString([System.Convert]::FromBase64String($Script.detectionScriptContent))
            Write-Host "   • Detection Script Content:" -ForegroundColor Green
            Write-Host "     $decodedDetectionScript" -ForegroundColor Yellow
        } catch {
            Write-Host "   • Detection Script Content: [Unable to decode - may not be base64]" -ForegroundColor Green
        }
    }
    
    if ($Script.remediationScriptContent) {
        try {
            $decodedRemediationScript = [System.Text.Encoding]::UTF8.GetString([System.Convert]::FromBase64String($Script.remediationScriptContent))
            Write-Host "   • Remediation Script Content:" -ForegroundColor Green
            Write-Host "     $decodedRemediationScript" -ForegroundColor Yellow
        } catch {
            Write-Host "   • Remediation Script Content: [Unable to decode - may not be base64]" -ForegroundColor Green
        }
    }
    
    if ($Script.isGlobalScript -ne $null) {
        Write-Host "   • Is Global Script: $($Script.isGlobalScript)" -ForegroundColor Green
    }
    
    if ($Script.deviceHealthScriptType) {
        Write-Host "   • Script Type: $($Script.deviceHealthScriptType)" -ForegroundColor Green
    }
    
    if ($Script.highestAvailableVersion) {
        Write-Host "   • Highest Available Version: $($Script.highestAvailableVersion)" -ForegroundColor Green
    }
    
    if ($Script.roleScopeTagIds -and $Script.roleScopeTagIds.Count -gt 0) {
        Write-Host "   • Role Scope Tag IDs: $($Script.roleScopeTagIds -join ', ')" -ForegroundColor Green
    }
    
    if ($Script.detectionScriptParameters -and $Script.detectionScriptParameters.Count -gt 0) {
        Write-Host "   • Detection Script Parameters: $($Script.detectionScriptParameters.Count) parameter(s)" -ForegroundColor Green
    }
    
    if ($Script.remediationScriptParameters -and $Script.remediationScriptParameters.Count -gt 0) {
        Write-Host "   • Remediation Script Parameters: $($Script.remediationScriptParameters.Count) parameter(s)" -ForegroundColor Green
    }
    
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    Write-Host ""
}

# Function to display device health script assignments details
function Show-DeviceHealthScriptAssignmentsDetails {
    param (
        [Parameter(Mandatory=$true)]
        $Assignments
    )
    
    Write-Host "📋 Device Health Script Assignments Details:" -ForegroundColor Cyan
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    
    if ($Assignments -and $Assignments.Count -gt 0) {
        Write-Host "   Found $($Assignments.Count) assignment(s)" -ForegroundColor Green
        Write-Host ""
        
        for ($i = 0; $i -lt $Assignments.Count; $i++) {
            $assignment = $Assignments[$i]
            Write-Host "   • Assignment $($i + 1):" -ForegroundColor Green
            
            if ($assignment.id) {
                Write-Host "     - ID: $($assignment.id)" -ForegroundColor Yellow
            }
            
            if ($assignment.target) {
                Write-Host "     - Target:" -ForegroundColor Yellow
                $target = $assignment.target
                
                if ($target.'@odata.type') {
                    $targetType = $target.'@odata.type' -replace '#microsoft.graph.', ''
                    Write-Host "       · Type: $targetType" -ForegroundColor Yellow
                    
                    switch ($targetType) {
                        "groupAssignmentTarget" {
                            if ($target.groupId) {
                                Write-Host "       · Group ID: $($target.groupId)" -ForegroundColor Yellow
                            }
                        }
                        "exclusionGroupAssignmentTarget" {
                            if ($target.groupId) {
                                Write-Host "       · Exclusion Group ID: $($target.groupId)" -ForegroundColor Yellow
                            }
                        }
                        "configurationManagerCollectionAssignmentTarget" {
                            if ($target.collectionId) {
                                Write-Host "       · Collection ID: $($target.collectionId)" -ForegroundColor Yellow
                            }
                        }
                        "allDevicesAssignmentTarget" {
                            Write-Host "       · Targets: All Devices" -ForegroundColor Yellow
                        }
                        "allLicensedUsersAssignmentTarget" {
                            Write-Host "       · Targets: All Licensed Users" -ForegroundColor Yellow
                        }
                    }
                }
                
                if ($target.deviceAndAppManagementAssignmentFilterId) {
                    Write-Host "       · Filter ID: $($target.deviceAndAppManagementAssignmentFilterId)" -ForegroundColor Yellow
                }
                
                if ($target.deviceAndAppManagementAssignmentFilterType) {
                    Write-Host "       · Filter Type: $($target.deviceAndAppManagementAssignmentFilterType)" -ForegroundColor Yellow
                }
            }
            
            if ($assignment.runRemediationScript -ne $null) {
                Write-Host "     - Run Remediation Script: $($assignment.runRemediationScript)" -ForegroundColor Yellow
            }
            
            if ($assignment.runSchedule) {
                Write-Host "     - Run Schedule:" -ForegroundColor Yellow
                $runSchedule = $assignment.runSchedule
                
                if ($runSchedule.interval) {
                    Write-Host "       · Interval: $($runSchedule.interval)" -ForegroundColor Yellow
                }
                
                if ($runSchedule.frequency) {
                    Write-Host "       · Frequency: $($runSchedule.frequency)" -ForegroundColor Yellow
                }
                
                if ($runSchedule.time) {
                    Write-Host "       · Time: $($runSchedule.time)" -ForegroundColor Yellow
                }
            }
            
            Write-Host ""
        }
    } else {
        Write-Host "   No assignments found for this device health script" -ForegroundColor Yellow
    }
    
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    Write-Host ""
}

# Function to export device health script to JSON
function Export-DeviceHealthScriptToJson {
    param (
        [Parameter(Mandatory=$true)]
        $Script,
        
        [Parameter(Mandatory=$true)]
        [string]$FilePath
    )
    
    try {
        Write-Host "💾 Exporting device health script to JSON..." -ForegroundColor Cyan
        Write-Host "   Export Path: $FilePath" -ForegroundColor Gray
        Write-Host ""
        
        # Ensure the directory exists
        $directory = Split-Path -Path $FilePath -Parent
        if ($directory -and -not (Test-Path -Path $directory)) {
            New-Item -ItemType Directory -Path $directory -Force | Out-Null
            Write-Host "   Created directory: $directory" -ForegroundColor Yellow
        }
        
        # Convert to JSON with proper depth to capture all nested objects
        $jsonOutput = $Script | ConvertTo-Json -Depth 10 -Compress:$false
        
        # Write to file with UTF8 encoding
        $jsonOutput | Out-File -FilePath $FilePath -Encoding UTF8 -Force
        
        Write-Host "✅ Device health script exported successfully!" -ForegroundColor Green
        Write-Host "   File Size: $([math]::Round((Get-Item $FilePath).Length / 1KB, 2)) KB" -ForegroundColor Green
        Write-Host ""
    }
    catch {
        Write-Host "❌ Error exporting device health script to JSON: $_" -ForegroundColor Red
        Write-Host ""
        throw
    }
}

# Script Setup
try {
    Write-Host "🔐 Connecting to Microsoft Graph..." -ForegroundColor Cyan
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    
    $secureClientSecret = ConvertTo-SecureString -String $ClientSecret -AsPlainText -Force
    $clientSecretCredential = New-Object -TypeName System.Management.Automation.PSCredential -ArgumentList $ClientId, $secureClientSecret
    
    Connect-MgGraph -ClientSecretCredential $clientSecretCredential -TenantId $TenantId -NoWelcome
    
    Write-Host "✅ Connected to Microsoft Graph successfully" -ForegroundColor Green
    Write-Host ""
    
    # Get the device health script with assignments
    $script = Get-DeviceHealthScript -ScriptId $DeviceHealthScriptId
    
    # Display script details
    if ($script) {
        Show-DeviceHealthScriptDetails -Script $script
        
        # Display assignments if they exist
        if ($script.assignments) {
            Show-DeviceHealthScriptAssignmentsDetails -Assignments $script.assignments
        } else {
            Write-Host "📊 No assignments found for this device health script" -ForegroundColor Yellow
            Write-Host ""
        }
        
        # Export to JSON if requested
        if ($ExportPath) {
            Export-DeviceHealthScriptToJson -Script $script -FilePath $ExportPath
        }
    } else {
        Write-Host "📊 No device health script found with ID: $DeviceHealthScriptId" -ForegroundColor Yellow
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