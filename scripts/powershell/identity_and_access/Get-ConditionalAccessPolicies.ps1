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
    HelpMessage="Specific Conditional Access Policy ID (if not provided, will list all policies)")]
    [string]$PolicyId,
    
    [Parameter(Mandatory=$true,
    HelpMessage="Export results to JSON file")]
    [bool]$ExportToJson
)

# Import required modules
Import-Module Microsoft.Graph.Authentication

# Function to get conditional access policies
function Get-ConditionalAccessPolicies {
    param (
        [Parameter(Mandatory=$false)]
        [string]$SpecificPolicyId
    )
    
    try {
        if ($SpecificPolicyId) {
            # GET specific policy - can use $expand for a single resource
            $uri = "https://graph.microsoft.com/beta/identity/conditionalAccess/policies/$SpecificPolicyId"
            Write-Host "🔍 Getting specific conditional access policy..." -ForegroundColor Cyan
            Write-Host "   Policy ID: $SpecificPolicyId" -ForegroundColor Gray
        } else {
            # GET all policies - do not use $expand for collections
            $uri = "https://graph.microsoft.com/beta/identity/conditionalAccess/policies"
            Write-Host "🔍 Getting all conditional access policies..." -ForegroundColor Cyan
        }
        
        Write-Host "   Endpoint: $uri" -ForegroundColor Gray
        Write-Host ""
        
        $response = Invoke-MgGraphRequest -Method GET -Uri $uri
        
        return $response
    }
    catch {
        Write-Host "❌ Error getting conditional access policies: $_" -ForegroundColor Red
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

# Function to export policies to JSON
function Export-PoliciesToJson {
    param (
        [Parameter(Mandatory=$true)]
        $Policies,
        
        [Parameter(Mandatory=$false)]
        [string]$SpecificPolicyId
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
        
        if ($SpecificPolicyId) {
            # Export single policy
            $policyName = $Policies.displayName -replace '[\\\/\:\*\?\"\<\>\|]', '_'
            if (-not $policyName) { $policyName = $SpecificPolicyId }
            $fileName = "ConditionalAccessPolicy_${policyName}_${timestamp}.json"
            $filePath = Join-Path -Path $outputDir -ChildPath $fileName
            
            $Policies | ConvertTo-Json -Depth 10 | Out-File -FilePath $filePath -Encoding UTF8
            
            Write-Host "💾 Exported policy to: $filePath" -ForegroundColor Green
        } else {
            # Export all policies
            $fileName = "ConditionalAccessPolicies_${timestamp}.json"
            $filePath = Join-Path -Path $outputDir -ChildPath $fileName
            
            $Policies | ConvertTo-Json -Depth 10 | Out-File -FilePath $filePath -Encoding UTF8
            
            Write-Host "💾 Exported policies to: $filePath" -ForegroundColor Green
        }
        
        return $filePath
    }
    catch {
        Write-Host "❌ Error exporting policies to JSON: $_" -ForegroundColor Red
        return $null
    }
}

# Function to display policy details
function Show-PolicyDetails {
    param (
        [Parameter(Mandatory=$true)]
        $Policy
    )
    
    Write-Host "📋 Conditional Access Policy Details:" -ForegroundColor Cyan
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
    
    if ($Policy.state) {
        Write-Host "   • State: $($Policy.state)" -ForegroundColor Green
    }
    
    if ($Policy.createdDateTime) {
        Write-Host "   • Created: $($Policy.createdDateTime)" -ForegroundColor Green
    }
    
    if ($Policy.modifiedDateTime) {
        Write-Host "   • Last Modified: $($Policy.modifiedDateTime)" -ForegroundColor Green
    }
    
    # Display conditions
    if ($Policy.conditions) {
        Write-Host "   • Conditions:" -ForegroundColor Green
        $conditions = $Policy.conditions
        
        # Users
        if ($conditions.users) {
            Write-Host "     - Users:" -ForegroundColor Yellow
            
            if ($conditions.users.includeUsers) {
                Write-Host "       · Include Users: $($conditions.users.includeUsers -join ', ')" -ForegroundColor Yellow
            }
            
            if ($conditions.users.excludeUsers) {
                Write-Host "       · Exclude Users: $($conditions.users.excludeUsers -join ', ')" -ForegroundColor Yellow
            }
            
            if ($conditions.users.includeGroups) {
                Write-Host "       · Include Groups: $($conditions.users.includeGroups -join ', ')" -ForegroundColor Yellow
            }
            
            if ($conditions.users.excludeGroups) {
                Write-Host "       · Exclude Groups: $($conditions.users.excludeGroups -join ', ')" -ForegroundColor Yellow
            }
            
            if ($conditions.users.includeRoles) {
                Write-Host "       · Include Roles: $($conditions.users.includeRoles -join ', ')" -ForegroundColor Yellow
            }
            
            if ($conditions.users.excludeRoles) {
                Write-Host "       · Exclude Roles: $($conditions.users.excludeRoles -join ', ')" -ForegroundColor Yellow
            }
            
            if ($conditions.users.includeGuestsOrExternalUsers) {
                Write-Host "       · Include Guests/External Users:" -ForegroundColor Yellow
                Write-Host "         · Types: $($conditions.users.includeGuestsOrExternalUsers.guestOrExternalUserTypes)" -ForegroundColor Yellow
                if ($conditions.users.includeGuestsOrExternalUsers.externalTenants) {
                    Write-Host "         · External Tenants Kind: $($conditions.users.includeGuestsOrExternalUsers.externalTenants.membershipKind)" -ForegroundColor Yellow
                    if ($conditions.users.includeGuestsOrExternalUsers.externalTenants.members) {
                        Write-Host "         · External Tenants Members: $($conditions.users.includeGuestsOrExternalUsers.externalTenants.members -join ', ')" -ForegroundColor Yellow
                    }
                }
            }
            
            if ($conditions.users.excludeGuestsOrExternalUsers) {
                Write-Host "       · Exclude Guests/External Users:" -ForegroundColor Yellow
                Write-Host "         · Types: $($conditions.users.excludeGuestsOrExternalUsers.guestOrExternalUserTypes)" -ForegroundColor Yellow
                if ($conditions.users.excludeGuestsOrExternalUsers.externalTenants) {
                    Write-Host "         · External Tenants Kind: $($conditions.users.excludeGuestsOrExternalUsers.externalTenants.membershipKind)" -ForegroundColor Yellow
                    if ($conditions.users.excludeGuestsOrExternalUsers.externalTenants.members) {
                        Write-Host "         · External Tenants Members: $($conditions.users.excludeGuestsOrExternalUsers.externalTenants.members -join ', ')" -ForegroundColor Yellow
                    }
                }
            }
        }
        
        # Applications
        if ($conditions.applications) {
            Write-Host "     - Applications:" -ForegroundColor Yellow
            
            if ($conditions.applications.includeApplications) {
                Write-Host "       · Include Applications: $($conditions.applications.includeApplications -join ', ')" -ForegroundColor Yellow
            }
            
            if ($conditions.applications.excludeApplications) {
                Write-Host "       · Exclude Applications: $($conditions.applications.excludeApplications -join ', ')" -ForegroundColor Yellow
            }
            
            if ($conditions.applications.includeUserActions) {
                Write-Host "       · Include User Actions: $($conditions.applications.includeUserActions -join ', ')" -ForegroundColor Yellow
            }
            
            if ($conditions.applications.includeAuthenticationContextClassReferences) {
                Write-Host "       · Include Auth Context Class Refs: $($conditions.applications.includeAuthenticationContextClassReferences -join ', ')" -ForegroundColor Yellow
            }
            
            if ($conditions.applications.applicationFilter) {
                Write-Host "       · Application Filter:" -ForegroundColor Yellow
                Write-Host "         · Mode: $($conditions.applications.applicationFilter.mode)" -ForegroundColor Yellow
                Write-Host "         · Rule: $($conditions.applications.applicationFilter.rule)" -ForegroundColor Yellow
            }
        }
        
        # Platforms
        if ($conditions.platforms) {
            Write-Host "     - Platforms:" -ForegroundColor Yellow
            
            if ($conditions.platforms.includePlatforms) {
                Write-Host "       · Include Platforms: $($conditions.platforms.includePlatforms -join ', ')" -ForegroundColor Yellow
            }
            
            if ($conditions.platforms.excludePlatforms) {
                Write-Host "       · Exclude Platforms: $($conditions.platforms.excludePlatforms -join ', ')" -ForegroundColor Yellow
            }
        }
        
        # Locations
        if ($conditions.locations) {
            Write-Host "     - Locations:" -ForegroundColor Yellow
            
            if ($conditions.locations.includeLocations) {
                Write-Host "       · Include Locations: $($conditions.locations.includeLocations -join ', ')" -ForegroundColor Yellow
            }
            
            if ($conditions.locations.excludeLocations) {
                Write-Host "       · Exclude Locations: $($conditions.locations.excludeLocations -join ', ')" -ForegroundColor Yellow
            }
        }
        
        # Client app types
        if ($conditions.clientAppTypes) {
            Write-Host "     - Client App Types: $($conditions.clientAppTypes -join ', ')" -ForegroundColor Yellow
        }
        
        # Sign-in risk levels
        if ($conditions.signInRiskLevels) {
            Write-Host "     - Sign-in Risk Levels: $($conditions.signInRiskLevels -join ', ')" -ForegroundColor Yellow
        }
        
        # User risk levels
        if ($conditions.userRiskLevels) {
            Write-Host "     - User Risk Levels: $($conditions.userRiskLevels -join ', ')" -ForegroundColor Yellow
        }
        
        # Service principal risk levels
        if ($conditions.servicePrincipalRiskLevels) {
            Write-Host "     - Service Principal Risk Levels: $($conditions.servicePrincipalRiskLevels -join ', ')" -ForegroundColor Yellow
        }
        
        # Client applications
        if ($conditions.clientApplications) {
            Write-Host "     - Client Applications:" -ForegroundColor Yellow
            
            if ($conditions.clientApplications.includeServicePrincipals) {
                Write-Host "       · Include Service Principals: $($conditions.clientApplications.includeServicePrincipals -join ', ')" -ForegroundColor Yellow
            }
            
            if ($conditions.clientApplications.excludeServicePrincipals) {
                Write-Host "       · Exclude Service Principals: $($conditions.clientApplications.excludeServicePrincipals -join ', ')" -ForegroundColor Yellow
            }
            
            if ($conditions.clientApplications.servicePrincipalFilter) {
                Write-Host "       · Service Principal Filter:" -ForegroundColor Yellow
                Write-Host "         · Mode: $($conditions.clientApplications.servicePrincipalFilter.mode)" -ForegroundColor Yellow
                Write-Host "         · Rule: $($conditions.clientApplications.servicePrincipalFilter.rule)" -ForegroundColor Yellow
            }
        }
        
        # Authentication flows
        if ($conditions.authenticationFlows) {
            Write-Host "     - Authentication Flows:" -ForegroundColor Yellow
            if ($conditions.authenticationFlows.transferMethods) {
                Write-Host "       · Transfer Methods: $($conditions.authenticationFlows.transferMethods -join ', ')" -ForegroundColor Yellow
            }
        }
        
        # Insider risk levels
        if ($conditions.insiderRiskLevels) {
            Write-Host "     - Insider Risk Levels: $($conditions.insiderRiskLevels -join ', ')" -ForegroundColor Yellow
        }
        
        # Device states
        if ($conditions.deviceStates) {
            Write-Host "     - Device States:" -ForegroundColor Yellow
            
            if ($conditions.deviceStates.includeStates) {
                Write-Host "       · Include States: $($conditions.deviceStates.includeStates -join ', ')" -ForegroundColor Yellow
            }
            
            if ($conditions.deviceStates.excludeStates) {
                Write-Host "       · Exclude States: $($conditions.deviceStates.excludeStates -join ', ')" -ForegroundColor Yellow
            }
        }
        
        # Devices
        if ($conditions.devices) {
            Write-Host "     - Devices:" -ForegroundColor Yellow
            
            if ($conditions.devices.includeDevices) {
                Write-Host "       · Include Devices: $($conditions.devices.includeDevices -join ', ')" -ForegroundColor Yellow
            }
            
            if ($conditions.devices.excludeDevices) {
                Write-Host "       · Exclude Devices: $($conditions.devices.excludeDevices -join ', ')" -ForegroundColor Yellow
            }
            
            if ($conditions.devices.deviceFilter) {
                Write-Host "       · Device Filter:" -ForegroundColor Yellow
                Write-Host "         · Mode: $($conditions.devices.deviceFilter.mode)" -ForegroundColor Yellow
                Write-Host "         · Rule: $($conditions.devices.deviceFilter.rule)" -ForegroundColor Yellow
            }
        }
    }
    
    # Display grant controls
    if ($Policy.grantControls) {
        Write-Host "   • Grant Controls:" -ForegroundColor Green
        $grantControls = $Policy.grantControls
        
        if ($grantControls.operator) {
            Write-Host "     - Operator: $($grantControls.operator)" -ForegroundColor Yellow
        }
        
        if ($grantControls.builtInControls) {
            Write-Host "     - Built-in Controls: $($grantControls.builtInControls -join ', ')" -ForegroundColor Yellow
        }
        
        if ($grantControls.customAuthenticationFactors) {
            Write-Host "     - Custom Authentication Factors: $($grantControls.customAuthenticationFactors -join ', ')" -ForegroundColor Yellow
        }
        
        if ($grantControls.termsOfUse) {
            Write-Host "     - Terms of Use: $($grantControls.termsOfUse -join ', ')" -ForegroundColor Yellow
        }
        
        if ($grantControls.authenticationStrength) {
            Write-Host "     - Authentication Strength:" -ForegroundColor Yellow
            Write-Host "       · ID: $($grantControls.authenticationStrength.id)" -ForegroundColor Yellow
            Write-Host "       · Display Name: $($grantControls.authenticationStrength.displayName)" -ForegroundColor Yellow
            Write-Host "       · Description: $($grantControls.authenticationStrength.description)" -ForegroundColor Yellow
            Write-Host "       · Policy Type: $($grantControls.authenticationStrength.policyType)" -ForegroundColor Yellow
            Write-Host "       · Requirements Satisfied: $($grantControls.authenticationStrength.requirementsSatisfied)" -ForegroundColor Yellow
            if ($grantControls.authenticationStrength.allowedCombinations) {
                Write-Host "       · Allowed Combinations: $($grantControls.authenticationStrength.allowedCombinations -join ', ')" -ForegroundColor Yellow
            }
        }
    }
    
    # Display session controls
    if ($Policy.sessionControls) {
        Write-Host "   • Session Controls:" -ForegroundColor Green
        $sessionControls = $Policy.sessionControls
        
        if ($sessionControls.applicationEnforcedRestrictions) {
            Write-Host "     - Application Enforced Restrictions:" -ForegroundColor Yellow
            Write-Host "       · Enabled: $($sessionControls.applicationEnforcedRestrictions.isEnabled)" -ForegroundColor Yellow
        }
        
        if ($sessionControls.cloudAppSecurity) {
            Write-Host "     - Cloud App Security:" -ForegroundColor Yellow
            Write-Host "       · Enabled: $($sessionControls.cloudAppSecurity.isEnabled)" -ForegroundColor Yellow
            if ($sessionControls.cloudAppSecurity.cloudAppSecurityType) {
                Write-Host "       · Type: $($sessionControls.cloudAppSecurity.cloudAppSecurityType)" -ForegroundColor Yellow
            }
        }
        
        if ($sessionControls.signInFrequency) {
            Write-Host "     - Sign-in Frequency:" -ForegroundColor Yellow
            Write-Host "       · Enabled: $($sessionControls.signInFrequency.isEnabled)" -ForegroundColor Yellow
            if ($sessionControls.signInFrequency.type) {
                Write-Host "       · Type: $($sessionControls.signInFrequency.type)" -ForegroundColor Yellow
            }
            if ($sessionControls.signInFrequency.value) {
                Write-Host "       · Value: $($sessionControls.signInFrequency.value)" -ForegroundColor Yellow
            }
            if ($sessionControls.signInFrequency.authenticationType) {
                Write-Host "       · Authentication Type: $($sessionControls.signInFrequency.authenticationType)" -ForegroundColor Yellow
            }
            if ($sessionControls.signInFrequency.frequencyInterval) {
                Write-Host "       · Frequency Interval: $($sessionControls.signInFrequency.frequencyInterval)" -ForegroundColor Yellow
            }
        }
        
        if ($sessionControls.persistentBrowser) {
            Write-Host "     - Persistent Browser:" -ForegroundColor Yellow
            Write-Host "       · Enabled: $($sessionControls.persistentBrowser.isEnabled)" -ForegroundColor Yellow
            if ($sessionControls.persistentBrowser.mode) {
                Write-Host "       · Mode: $($sessionControls.persistentBrowser.mode)" -ForegroundColor Yellow
            }
        }
        
        if ($sessionControls.continuousAccessEvaluation) {
            Write-Host "     - Continuous Access Evaluation:" -ForegroundColor Yellow
            Write-Host "       · Mode: $($sessionControls.continuousAccessEvaluation.mode)" -ForegroundColor Yellow
        }
        
        if ($sessionControls.disableResilienceDefaults -ne $null) {
            Write-Host "     - Disable Resilience Defaults: $($sessionControls.disableResilienceDefaults)" -ForegroundColor Yellow
        }
        
        if ($sessionControls.secureSignInSession) {
            Write-Host "     - Secure Sign-in Session:" -ForegroundColor Yellow
            Write-Host "       · Enabled: $($sessionControls.secureSignInSession.isEnabled)" -ForegroundColor Yellow
            if ($sessionControls.secureSignInSession.type) {
                Write-Host "       · Type: $($sessionControls.secureSignInSession.type)" -ForegroundColor Yellow
            }
            if ($sessionControls.secureSignInSession.value) {
                Write-Host "       · Value: $($sessionControls.secureSignInSession.value)" -ForegroundColor Yellow
            }
            if ($sessionControls.secureSignInSession.authenticationType) {
                Write-Host "       · Authentication Type: $($sessionControls.secureSignInSession.authenticationType)" -ForegroundColor Yellow
            }
            if ($sessionControls.secureSignInSession.frequencyInterval) {
                Write-Host "       · Frequency Interval: $($sessionControls.secureSignInSession.frequencyInterval)" -ForegroundColor Yellow
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
    
    # Get the policies
    $policies = Get-ConditionalAccessPolicies -SpecificPolicyId $PolicyId
    
    # Export to JSON if requested
    if ($ExportToJson) {
        $jsonPath = Export-PoliciesToJson -Policies $policies -SpecificPolicyId $PolicyId
    }
    
    if ($PolicyId) {
        # Display single policy
        Show-PolicyDetails -Policy $policies
    } else {
        # Display all policies
        if ($policies.value -and $policies.value.Count -gt 0) {
            Write-Host "📊 Found $($policies.value.Count) conditional access policy(ies)" -ForegroundColor Green
            Write-Host ""
            
            for ($i = 0; $i -lt $policies.value.Count; $i++) {
                Write-Host "Policy $($i + 1):" -ForegroundColor Magenta
                Show-PolicyDetails -Policy $policies.value[$i]
            }
        } elseif ($policies -and -not $policies.value) {
            # Single policy returned (not in a collection)
            Write-Host "📊 Found 1 conditional access policy" -ForegroundColor Green
            Write-Host ""
            Show-PolicyDetails -Policy $policies
        } else {
            Write-Host "📊 No conditional access policies found" -ForegroundColor Yellow
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