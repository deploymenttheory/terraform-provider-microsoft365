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
    HelpMessage="Specific Settings Catalog Policy ID to check")]
    [string]$PolicyId,
    
    [Parameter(Mandatory=$false,
    HelpMessage="Create a test Settings Catalog policy for import testing")]
    [switch]$CreateTestPolicy,
    
    [Parameter(Mandatory=$false,
    HelpMessage="Check what type of policy a given ID represents")]
    [string]$CheckPolicyType
)

# Import required modules
Import-Module Microsoft.Graph.Authentication

# Function to get Settings Catalog policies
function Get-SettingsCatalogPolicies {
    param (
        [Parameter(Mandatory=$false)]
        [string]$SpecificPolicyId
    )
    
    try {
        if ($SpecificPolicyId) {
            # GET specific Settings Catalog policy
            $uri = "https://graph.microsoft.com/beta/deviceManagement/configurationPolicies/$SpecificPolicyId"
            Write-Host "🔍 Getting specific Settings Catalog policy..." -ForegroundColor Cyan
            Write-Host "   Policy ID: $SpecificPolicyId" -ForegroundColor Gray
        } else {
            # GET all Settings Catalog policies
            $uri = "https://graph.microsoft.com/beta/deviceManagement/configurationPolicies"
            Write-Host "🔍 Getting all Settings Catalog policies..." -ForegroundColor Cyan
        }
        
        Write-Host "   Endpoint: $uri" -ForegroundColor Gray
        Write-Host ""
        
        $response = Invoke-MgGraphRequest -Method GET -Uri $uri
        
        return $response
    }
    catch {
        Write-Host "❌ Error getting Settings Catalog policies: $_" -ForegroundColor Red
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

# Function to check policy type across different endpoints
function Test-PolicyType {
    param (
        [Parameter(Mandatory=$true)]
        [string]$PolicyId
    )
    
    Write-Host "🔍 Testing policy type for ID: $PolicyId" -ForegroundColor Cyan
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    
    # Test 1: Settings Catalog
    Write-Host "1️⃣ Testing Settings Catalog endpoint..." -ForegroundColor Yellow
    try {
        $uri = "https://graph.microsoft.com/beta/deviceManagement/configurationPolicies/$PolicyId"
        $response = Invoke-MgGraphRequest -Method GET -Uri $uri
        Write-Host "   ✅ Found in Settings Catalog!" -ForegroundColor Green
        Write-Host "   📝 Name: $($response.name)" -ForegroundColor Cyan
        Write-Host "   🖥️  Platform: $($response.platforms)" -ForegroundColor Cyan
        Write-Host "   ⚙️  Technology: $($response.technologies)" -ForegroundColor Cyan
        Write-Host "   📊 Settings Count: $($response.settingCount)" -ForegroundColor Cyan
        return "SettingsCatalog"
    }
    catch {
        Write-Host "   ❌ Not found in Settings Catalog" -ForegroundColor Red
    }
    
    # Test 2: Legacy Device Configuration
    Write-Host "2️⃣ Testing Device Configuration endpoint..." -ForegroundColor Yellow
    try {
        $uri = "https://graph.microsoft.com/beta/deviceManagement/deviceConfigurations/$PolicyId"
        $response = Invoke-MgGraphRequest -Method GET -Uri $uri
        Write-Host "   ✅ Found in Device Configuration!" -ForegroundColor Green
        Write-Host "   📝 Name: $($response.displayName)" -ForegroundColor Cyan
        Write-Host "   🏷️  OData Type: $($response.'@odata.type')" -ForegroundColor Cyan
        return "DeviceConfiguration"
    }
    catch {
        Write-Host "   ❌ Not found in Device Configuration" -ForegroundColor Red
    }
    
    # Test 3: Group Policy Configurations
    Write-Host "3️⃣ Testing Group Policy Configuration endpoint..." -ForegroundColor Yellow
    try {
        $uri = "https://graph.microsoft.com/beta/deviceManagement/groupPolicyConfigurations/$PolicyId"
        $response = Invoke-MgGraphRequest -Method GET -Uri $uri
        Write-Host "   ✅ Found in Group Policy Configuration!" -ForegroundColor Green
        Write-Host "   📝 Name: $($response.displayName)" -ForegroundColor Cyan
        return "GroupPolicy"
    }
    catch {
        Write-Host "   ❌ Not found in Group Policy Configuration" -ForegroundColor Red
    }
    
    # Test 4: Intent-based policies
    Write-Host "4️⃣ Testing Intent endpoint..." -ForegroundColor Yellow
    try {
        $uri = "https://graph.microsoft.com/beta/deviceManagement/intents/$PolicyId"
        $response = Invoke-MgGraphRequest -Method GET -Uri $uri
        Write-Host "   ✅ Found in Intents!" -ForegroundColor Green
        Write-Host "   📝 Name: $($response.displayName)" -ForegroundColor Cyan
        return "Intent"
    }
    catch {
        Write-Host "   ❌ Not found in Intents" -ForegroundColor Red
    }
    
    Write-Host "❌ Policy not found in any known endpoint!" -ForegroundColor Red
    return "NotFound"
}

# Function to create test Settings Catalog policy
function New-TestSettingsCatalogPolicy {
    try {
        Write-Host "🔨 Creating test Settings Catalog policy..." -ForegroundColor Cyan
        Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
        
        $policyBody = @{
            name = "Test Settings Catalog Policy - $(Get-Date -Format 'yyyy-MM-dd HH:mm')"
            description = "Simple policy created for testing Terraform import functionality"
            platforms = "windows10"
            technologies = "mdm"
            settings = @(
                @{
                    "@odata.type" = "#microsoft.graph.deviceManagementConfigurationSetting"
                    settingInstance = @{
                        "@odata.type" = "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance"
                        settingDefinitionId = "device_vendor_msft_policy_config_browser_allowaddressbardropdown"
                        simpleSettingValue = @{
                            "@odata.type" = "#microsoft.graph.deviceManagementConfigurationIntegerSettingValue"
                            value = 1
                        }
                    }
                }
            )
        } | ConvertTo-Json -Depth 10
        
        $uri = "https://graph.microsoft.com/beta/deviceManagement/configurationPolicies"
        Write-Host "   Endpoint: $uri" -ForegroundColor Gray
        
        $newPolicy = Invoke-MgGraphRequest -Method POST -Uri $uri -Body $policyBody -ContentType "application/json"
        
        Write-Host "✅ Successfully created Settings Catalog policy!" -ForegroundColor Green
        Write-Host "   📝 Policy ID: $($newPolicy.id)" -ForegroundColor Yellow
        Write-Host "   📝 Policy Name: $($newPolicy.name)" -ForegroundColor Cyan
        Write-Host ""
        Write-Host "🚀 Use this for Terraform import testing:" -ForegroundColor Green
        Write-Host "   terraform import microsoft365_graph_beta_device_management_settings_catalog.imported_policy `"$($newPolicy.id)`"" -ForegroundColor Cyan
        
        return $newPolicy
    }
    catch {
        Write-Host "❌ Error creating policy: $_" -ForegroundColor Red
        Write-Host ""
        Write-Host "💡 Alternative: Create manually through Intune admin center:" -ForegroundColor Yellow
        Write-Host "   1. Go to https://intune.microsoft.com" -ForegroundColor Gray
        Write-Host "   2. Devices > Configuration policies > Create policy" -ForegroundColor Gray
        Write-Host "   3. Platform: Windows 10 and later" -ForegroundColor Gray
        Write-Host "   4. Profile type: Settings catalog" -ForegroundColor Gray
        Write-Host "   5. Add any simple setting (e.g., Browser settings)" -ForegroundColor Gray
        
        throw
    }
}

# Function to display policy details
function Show-PolicyDetails {
    param (
        [Parameter(Mandatory=$true)]
        $Policy
    )
    
    Write-Host "📋 Policy Details:" -ForegroundColor Cyan
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    
    if ($Policy.id) {
        Write-Host "   • ID: $($Policy.id)" -ForegroundColor Green
    }
    
    if ($Policy.name) {
        Write-Host "   • Name: $($Policy.name)" -ForegroundColor Green
    }
    
    if ($Policy.description) {
        Write-Host "   • Description: $($Policy.description)" -ForegroundColor Green
    }
    
    if ($Policy.platforms) {
        Write-Host "   • Platform: $($Policy.platforms)" -ForegroundColor Green
    }
    
    if ($Policy.technologies) {
        Write-Host "   • Technology: $($Policy.technologies)" -ForegroundColor Green
    }
    
    if ($Policy.settingCount -ne $null) {
        Write-Host "   • Settings Count: $($Policy.settingCount)" -ForegroundColor Green
    }
    
    if ($Policy.isAssigned -ne $null) {
        Write-Host "   • Is Assigned: $($Policy.isAssigned)" -ForegroundColor Green
    }
    
    if ($Policy.createdDateTime) {
        Write-Host "   • Created: $($Policy.createdDateTime)" -ForegroundColor Green
    }
    
    if ($Policy.lastModifiedDateTime) {
        Write-Host "   • Last Modified: $($Policy.lastModifiedDateTime)" -ForegroundColor Green
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
    
    # Check policy type if requested
    if ($CheckPolicyType) {
        $policyType = Test-PolicyType -PolicyId $CheckPolicyType
        Write-Host ""
        Write-Host "🏷️  Policy Type Result: $policyType" -ForegroundColor Magenta
        Write-Host ""
        
        if ($policyType -ne "SettingsCatalog") {
            Write-Host "⚠️  This policy is NOT compatible with your Settings Catalog Terraform provider!" -ForegroundColor Red
            Write-Host "   You need a different Terraform resource for this policy type." -ForegroundColor Yellow
        } else {
            Write-Host "✅ This policy IS compatible with your Settings Catalog Terraform provider!" -ForegroundColor Green
        }
        return
    }
    
    # Create test policy if requested
    if ($CreateTestPolicy) {
        $newPolicy = New-TestSettingsCatalogPolicy
        Show-PolicyDetails -Policy $newPolicy
        return
    }
    
    # Get the policies
    $policies = Get-SettingsCatalogPolicies -SpecificPolicyId $PolicyId
    
    if ($PolicyId) {
        # Display single policy
        Show-PolicyDetails -Policy $policies
        Write-Host "🚀 Use this for Terraform import testing:" -ForegroundColor Green
        Write-Host "   terraform import microsoft365_graph_beta_device_management_settings_catalog.imported_policy `"$($policies.id)`"" -ForegroundColor Cyan
    } else {
        # Display all policies
        if ($policies.value -and $policies.value.Count -gt 0) {
            Write-Host "📊 Found $($policies.value.Count) Settings Catalog policy(s)" -ForegroundColor Green
            Write-Host ""
            
            for ($i = 0; $i -lt $policies.value.Count; $i++) {
                Write-Host "Policy $($i + 1):" -ForegroundColor Magenta
                Show-PolicyDetails -Policy $policies.value[$i]
                Write-Host "🚀 Import command:" -ForegroundColor Green
                Write-Host "   terraform import microsoft365_graph_beta_device_management_settings_catalog.imported_policy `"$($policies.value[$i].id)`"" -ForegroundColor Cyan
                Write-Host ""
            }
        } elseif ($policies -and -not $policies.value) {
            # Single policy returned (not in a collection)
            Write-Host "📊 Found 1 Settings Catalog policy" -ForegroundColor Green
            Write-Host ""
            Show-PolicyDetails -Policy $policies
            Write-Host "🚀 Import command:" -ForegroundColor Green
            Write-Host "   terraform import microsoft365_graph_beta_device_management_settings_catalog.imported_policy `"$($policies.id)`"" -ForegroundColor Cyan
        } else {
            Write-Host "📊 No Settings Catalog policies found" -ForegroundColor Yellow
            Write-Host ""
            Write-Host "💡 To create a test policy, run this script with -CreateTestPolicy" -ForegroundColor Yellow
            Write-Host "   Or create one manually through Intune admin center" -ForegroundColor Yellow
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