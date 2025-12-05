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
    HelpMessage="Inventory Policy ID")]
    [ValidateNotNullOrEmpty()]
    [string]$InventoryPolicyId
)

# Import required modules
Import-Module Microsoft.Graph.Authentication

# Function to get inventory policy details
function Get-InventoryPolicy {
    param (
        [Parameter(Mandatory=$true)]
        [string]$PolicyId
    )
    
    try {
        # GET inventory policy
        $uri = "https://graph.microsoft.com/beta/deviceManagement/inventoryPolicies('$PolicyId')"
        Write-Host "🔍 Getting inventory policy..." -ForegroundColor Cyan
        Write-Host "   Policy ID: $PolicyId" -ForegroundColor Gray
        Write-Host "   Endpoint: $uri" -ForegroundColor Gray
        Write-Host ""
        
        $response = Invoke-MgGraphRequest -Method GET -Uri $uri
        
        return $response
    }
    catch {
        Write-Host "❌ Error getting inventory policy: $_" -ForegroundColor Red
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

# Function to get inventory policy settings
function Get-InventoryPolicySettings {
    param (
        [Parameter(Mandatory=$true)]
        [string]$PolicyId
    )
    
    try {
        # GET inventory policy settings
        $uri = "https://graph.microsoft.com/beta/deviceManagement/inventoryPolicies('$PolicyId')/settings"
        Write-Host "🔍 Getting inventory policy settings..." -ForegroundColor Cyan
        Write-Host "   Policy ID: $PolicyId" -ForegroundColor Gray
        Write-Host "   Endpoint: $uri" -ForegroundColor Gray
        Write-Host ""
        
        $response = Invoke-MgGraphRequest -Method GET -Uri $uri
        
        return $response
    }
    catch {
        Write-Host "❌ Error getting inventory policy settings: $_" -ForegroundColor Red
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

# Function to get inventory policy assignments
function Get-InventoryPolicyAssignments {
    param (
        [Parameter(Mandatory=$true)]
        [string]$PolicyId
    )
    
    try {
        # GET inventory policy assignments
        $uri = "https://graph.microsoft.com/beta/deviceManagement/inventoryPolicies('$PolicyId')/assignments"
        Write-Host "🔍 Getting inventory policy assignments..." -ForegroundColor Cyan
        Write-Host "   Policy ID: $PolicyId" -ForegroundColor Gray
        Write-Host "   Endpoint: $uri" -ForegroundColor Gray
        Write-Host ""
        
        $response = Invoke-MgGraphRequest -Method GET -Uri $uri
        
        return $response
    }
    catch {
        Write-Host "❌ Error getting inventory policy assignments: $_" -ForegroundColor Red
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

# Function to display inventory policy details
function Show-InventoryPolicyDetails {
    param (
        [Parameter(Mandatory=$true)]
        $Policy
    )
    
    Write-Host "📋 Inventory Policy Details:" -ForegroundColor Cyan
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    
    if ($Policy.id) {
        Write-Host "   • ID: $($Policy.id)" -ForegroundColor Green
    }
    
    if ($Policy.displayName) {
        Write-Host "   • Display Name: $($Policy.displayName)" -ForegroundColor Green
    }
    
    if ($Policy.description) {
        Write-Host "   • Description: $($Policy.description)" -ForegroundColor Green
    }
    
    if ($Policy.createdDateTime) {
        Write-Host "   • Created: $($Policy.createdDateTime)" -ForegroundColor Green
    }
    
    if ($Policy.lastModifiedDateTime) {
        Write-Host "   • Last Modified: $($Policy.lastModifiedDateTime)" -ForegroundColor Green
    }
    
    if ($Policy.version) {
        Write-Host "   • Version: $($Policy.version)" -ForegroundColor Green
    }
    
    # Display any additional properties
    $Policy.PSObject.Properties | Where-Object { 
        $_.Name -notin 'id', 'displayName', 'description', 'createdDateTime', 'lastModifiedDateTime', 'version', '@odata.context'
    } | ForEach-Object {
        $propertyName = $_.Name
        $propertyValue = $_.Value
        
        # Handle complex objects
        if ($propertyValue -is [System.Management.Automation.PSCustomObject]) {
            Write-Host "   • ${propertyName}:" -ForegroundColor Green
            $propertyValue.PSObject.Properties | ForEach-Object {
                Write-Host "     - $($_.Name): $($_.Value)" -ForegroundColor Yellow
            }
        } else {
            Write-Host "   • ${propertyName}: $propertyValue" -ForegroundColor Green
        }
    }
    
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    Write-Host ""
}

# Function to display inventory policy settings details
function Show-InventoryPolicySettingsDetails {
    param (
        [Parameter(Mandatory=$true)]
        $Settings
    )
    
    Write-Host "📋 Inventory Policy Settings Details:" -ForegroundColor Cyan
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    
    if ($Settings.value -and $Settings.value.Count -gt 0) {
        for ($i = 0; $i -lt $Settings.value.Count; $i++) {
            $setting = $Settings.value[$i]
            Write-Host "   • Setting $($i + 1):" -ForegroundColor Green
            
            if ($setting.'@odata.type') {
                $settingType = $setting.'@odata.type' -replace '#microsoft.graph.', ''
                Write-Host "     - Type: $settingType" -ForegroundColor Yellow
            }
            
            if ($setting.id) {
                Write-Host "     - ID: $($setting.id)" -ForegroundColor Yellow
            }
            
            if ($setting.displayName) {
                Write-Host "     - Display Name: $($setting.displayName)" -ForegroundColor Yellow
            }
            
            if ($setting.description) {
                Write-Host "     - Description: $($setting.description)" -ForegroundColor Yellow
            }
            
            if ($setting.state -ne $null) {
                Write-Host "     - State: $($setting.state)" -ForegroundColor Yellow
            }
            
            if ($setting.category) {
                Write-Host "     - Category: $($setting.category)" -ForegroundColor Yellow
            }
            
            if ($setting.platform) {
                Write-Host "     - Platform: $($setting.platform)" -ForegroundColor Yellow
            }
            
            # Display any additional properties specific to the setting type
            $setting.PSObject.Properties | Where-Object { 
                $_.Name -notin '@odata.type', 'id', 'displayName', 'description', 'state', 'category', 'platform'
            } | ForEach-Object {
                $propertyName = $_.Name
                $propertyValue = $_.Value
                
                # Handle complex objects
                if ($propertyValue -is [System.Management.Automation.PSCustomObject]) {
                    Write-Host "     - ${propertyName}:" -ForegroundColor Yellow
                    $propertyValue.PSObject.Properties | ForEach-Object {
                        Write-Host "       · $($_.Name): $($_.Value)" -ForegroundColor Yellow
                    }
                } else {
                    Write-Host "     - ${propertyName}: $propertyValue" -ForegroundColor Yellow
                }
            }
            
            Write-Host ""
        }
    } else {
        Write-Host "   No settings found for this inventory policy" -ForegroundColor Yellow
    }
    
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    Write-Host ""
}

# Function to display inventory policy assignments details
function Show-InventoryPolicyAssignmentsDetails {
    param (
        [Parameter(Mandatory=$true)]
        $Assignments
    )
    
    Write-Host "📋 Inventory Policy Assignments Details:" -ForegroundColor Cyan
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    
    if ($Assignments.value -and $Assignments.value.Count -gt 0) {
        Write-Host "   Found $($Assignments.value.Count) assignment(s)" -ForegroundColor Green
        Write-Host ""
        
        for ($i = 0; $i -lt $Assignments.value.Count; $i++) {
            $assignment = $Assignments.value[$i]
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
            
            Write-Host ""
        }
    } else {
        Write-Host "   No assignments found for this inventory policy" -ForegroundColor Yellow
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
    
    # Get the inventory policy
    $policy = Get-InventoryPolicy -PolicyId $InventoryPolicyId
    
    # Get the inventory policy settings
    $settings = Get-InventoryPolicySettings -PolicyId $InventoryPolicyId
    
    # Get the inventory policy assignments
    $assignments = Get-InventoryPolicyAssignments -PolicyId $InventoryPolicyId
    
    # Display policy details
    if ($policy) {
        Show-InventoryPolicyDetails -Policy $policy
    } else {
        Write-Host "📊 No policy found with ID: $InventoryPolicyId" -ForegroundColor Yellow
    }
    
    # Display settings
    if ($settings) {
        Show-InventoryPolicySettingsDetails -Settings $settings
    } else {
        Write-Host "📊 No settings found for this inventory policy" -ForegroundColor Yellow
    }
    
    # Display assignments
    if ($assignments) {
        Show-InventoryPolicyAssignmentsDetails -Assignments $assignments
    } else {
        Write-Host "📊 No assignments found for this inventory policy" -ForegroundColor Yellow
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
