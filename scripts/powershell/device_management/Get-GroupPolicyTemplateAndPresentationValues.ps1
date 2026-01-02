<#
.SYNOPSIS
    Looks up Group Policy definitions and presentations from Microsoft Graph API catalog

.DESCRIPTION
    This script queries the Microsoft Graph API to retrieve Group Policy definitions
    and their supported presentation types (value types) from the built-in catalog.
    
    It supports two modes:
    - LIST mode: Retrieves the full catalog of available Group Policy templates
    - MATCH mode: Looks up a specific policy name and returns its definition + presentations
    
    The script performs these API call steps:
    1. Query groupPolicyDefinitions - the Microsoft catalog of available policies
    2. Query presentations for the definition - what value types are supported (text, boolean, dropdown, etc.)

.PARAMETER TenantId
    Azure AD Tenant ID

.PARAMETER ClientId
    Application (Client) ID of the Azure AD app registration

.PARAMETER ClientSecret
    Client secret for authentication

.PARAMETER PolicyName
    Exact policy display name to search for. Leave empty for LIST mode.
    Example: "Allow users to connect remotely by using Remote Desktop Services"

.PARAMETER ClassType
    Policy class type filter: 'user' or 'machine'. Optional.

.PARAMETER CategoryPath
    Category path for disambiguation when multiple policies have the same name.
    Example: "\Windows Components\Remote Desktop Services\..."

.PARAMETER PresentationIndex
    Which presentation to use if multiple exist (default: 0)

.PARAMETER OutputDirectory
    Directory path where JSON responses will be saved

.EXAMPLE
    # LIST MODE - Get all Group Policy definitions (full catalog)
    pwsh Get-GroupPolicyTemplateAndPresentationValues `
        -TenantId "your-tenant-id" `
        -ClientId "your-client-id" `
        -ClientSecret "your-secret" `
        -OutputDirectory "C:\temp\gpo_catalog"
    
    Returns: Complete catalog of all available Group Policy templates saved to JSON

.EXAMPLE
    # LIST MODE - Get only MACHINE policies
    pwsh Get-GroupPolicyTemplateAndPresentationValues `
        -TenantId "your-tenant-id" `
        -ClientId "your-client-id" `
        -ClientSecret "your-secret" `
        -ClassType "machine" `
        -OutputDirectory "C:\temp\gpo_catalog"
    
    Returns: All machine (computer) configuration policies

.EXAMPLE
    # LIST MODE - Get only USER policies
    pwsh Get-GroupPolicyTemplateAndPresentationValues `
        -TenantId "your-tenant-id" `
        -ClientId "your-client-id" `
        -ClientSecret "your-secret" `
        -ClassType "user" `
        -OutputDirectory "C:\temp\gpo_catalog"
    
    Returns: All user configuration policies

.EXAMPLE
    # MATCH MODE - Look up specific policy definition and presentations
    pwsh Get-GroupPolicyTemplateAndPresentationValues `
        -TenantId "your-tenant-id" `
        -ClientId "your-client-id" `
        -ClientSecret "your-secret" `
        -PolicyName "Allow users to connect remotely by using Remote Desktop Services" `
        -ClassType "machine" `
        -OutputDirectory "C:\temp\rdp_policy"
    
    Returns: 
    - step1_definitions.json: Policy definition (ID, name, category, class type)
    - step2_presentation_templates.json: Supported value types (boolean, text, dropdown, etc.)

.EXAMPLE
    # MATCH MODE - Look up with category path for disambiguation
    pwsh Get-GroupPolicyTemplateAndPresentationValues `
        -TenantId "your-tenant-id" `
        -ClientId "your-client-id" `
        -ClientSecret "your-secret" `
        -PolicyName "Enable Profile Containers" `
        -ClassType "machine" `
        -CategoryPath "\FSLogix\Profile Containers" `
        -OutputDirectory "C:\temp\fslogix_policy"
    
    Returns: Specific policy definition when multiple policies share the same name

.EXAMPLE
    # Real-world example with actual credentials
    pwsh Get-GroupPolicyTemplateAndPresentationValues `
        -TenantId "00000000-0000-0000-0000-000000000000" `
        -ClientId "00000000-0000-0000-0000-000000000000" `
        -ClientSecret "your-secret" `
        -PolicyName "Allow users to connect remotely by using Remote Desktop Services" `
        -ClassType "machine" `
        -OutputDirectory "/Users/dafyddwatkins/localtesting/group_policy_api_test"

.NOTES
    File Name      : Get-GroupPolicyTemplateAndPresentationValues
    Prerequisite   : Microsoft.Graph.Authentication PowerShell module
    Copyright      : Based on id_resolver.go from terraform-provider-microsoft365
    
    API Permissions Required:
    - DeviceManagementConfiguration.Read.All (minimum)
    - DeviceManagementConfiguration.ReadWrite.All (for full access)

.LINK
    https://learn.microsoft.com/en-us/graph/api/intune-grouppolicy-grouppolicydefinition-list
    https://learn.microsoft.com/en-us/graph/api/intune-grouppolicy-grouppolicydefinitionvalue-list
#>

[CmdletBinding()]
param (
    [Parameter(Mandatory=$true)]
    [string]$TenantId,

    [Parameter(Mandatory=$true)]
    [string]$ClientId,
    
    [Parameter(Mandatory=$true)]
    [string]$ClientSecret,
    
    [Parameter(Mandatory=$false,
    HelpMessage="Exact policy name to search for (leave empty to list all)")]
    [string]$PolicyName = "",
    
    [Parameter(Mandatory=$false,
    HelpMessage="Class type: 'user' or 'machine' (optional filter)")]
    [ValidateSet("user", "machine", "")]
    [string]$ClassType = "",
    
    [Parameter(Mandatory=$false,
    HelpMessage="Optional category path for disambiguation")]
    [string]$CategoryPath = "",
    
    [Parameter(Mandatory=$false,
    HelpMessage="Presentation index to use (default 0)")]
    [int]$PresentationIndex = 0,
    
    [Parameter(Mandatory=$true)]
    [string]$OutputDirectory
)

# Import module
Import-Module Microsoft.Graph.Authentication

# Function to authenticate
function Connect-MicrosoftGraph {
    param (
        [string]$TenantId,
        [string]$ClientId,
        [string]$ClientSecret
    )
    
    try {
        Write-Host "🔐 Authenticating to Microsoft Graph..." -ForegroundColor Cyan
        
        $secureClientSecret = ConvertTo-SecureString -String $ClientSecret -AsPlainText -Force
        $clientSecretCredential = New-Object -TypeName System.Management.Automation.PSCredential -ArgumentList $ClientId, $secureClientSecret
        
        Connect-MgGraph -ClientSecretCredential $clientSecretCredential -TenantId $TenantId -NoWelcome
        
        Write-Host "✅ Connected successfully" -ForegroundColor Green
    }
    catch {
        Write-Host "❌ Authentication failed: $_" -ForegroundColor Red
        throw
    }
}

# Step 1: Get all definitions or resolve specific policy
function Get-GroupPolicyDefinitions {
    param (
        [string]$PolicyName,
        [string]$ClassType,
        [string]$CategoryPath,
        [string]$OutputFile
    )
    
    try {
        # Determine mode
        $isListMode = [string]::IsNullOrWhiteSpace($PolicyName)
        
        if ($isListMode) {
            Write-Host "`n📍 STEP 1: LIST ALL Group Policy Definition Templates" -ForegroundColor Yellow
            Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Yellow
            Write-Host "Mode: LIST (retrieving full catalog)" -ForegroundColor Cyan
        } else {
            Write-Host "`n📍 STEP 1: MATCH Policy Name → Definition Template ID" -ForegroundColor Yellow
            Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Yellow
            Write-Host "Mode: MATCH (searching for specific policy)" -ForegroundColor Cyan
        }
        
        # Build URL
        $queryParams = @()
        
        if (-not $isListMode) {
            # MATCH mode - build filter
            $filterParts = @()
            
            $filterParts += "displayName eq '$PolicyName'"
            
            if ($ClassType) {
                $normalizedClassType = $ClassType.ToLower()
                $filterParts += "classType eq '$normalizedClassType'"
            }
            
            $filter = $filterParts -join " and "
            $queryParams += "`$filter=$([System.Web.HttpUtility]::UrlEncode($filter))"
            
            Write-Host "Search: PolicyName='$PolicyName'" -ForegroundColor Gray
            if ($ClassType) {
                Write-Host "Search: ClassType='$ClassType'" -ForegroundColor Gray
            }
            if ($CategoryPath) {
                Write-Host "Search: CategoryPath='$CategoryPath'" -ForegroundColor Gray
            }
        } else {
            # LIST mode - optionally filter by class type
            if ($ClassType) {
                $normalizedClassType = $ClassType.ToLower()
                $filter = "classType eq '$normalizedClassType'"
                $queryParams += "`$filter=$([System.Web.HttpUtility]::UrlEncode($filter))"
                Write-Host "Filter: ClassType='$ClassType'" -ForegroundColor Gray
            }
            
            # Get more results for list mode
            $queryParams += "`$top=999"
        }
        
        # Always select these fields
        $queryParams += "`$select=id,displayName,classType,categoryPath,policyType,version"
        
        $queryString = $queryParams -join "&"
        $url = "https://graph.microsoft.com/beta/deviceManagement/groupPolicyDefinitions?$queryString"
        
        Write-Host "API Call: $url" -ForegroundColor DarkGray
        
        $response = Invoke-MgGraphRequest -Method GET -Uri $url
        
        # Save response
        $prettyJson = ConvertTo-Json -InputObject $response -Depth 10 -Compress:$false
        $prettyJson | Out-File -FilePath $OutputFile -Encoding utf8
        
        if (-not $response.value -or $response.value.Count -eq 0) {
            if ($isListMode) {
                throw "No definitions found in catalog"
            } else {
                throw "No definition found for policy '$PolicyName'"
            }
        }
        
        $count = $response.value.Count
        Write-Host "✅ Found $count definition(s)" -ForegroundColor Green
        
        if ($isListMode) {
            # LIST mode - show summary
            Write-Host "`nShowing first 10 definitions:" -ForegroundColor Cyan
            $preview = $response.value | Select-Object -First 10
            foreach ($def in $preview) {
                Write-Host "  • $($def.displayName)" -ForegroundColor Gray
                Write-Host "    ID: $($def.id) | Type: $($def.classType) | Path: $($def.categoryPath)" -ForegroundColor DarkGray
            }
            if ($count -gt 10) {
                Write-Host "  ... and $($count - 10) more (see JSON file for full list)" -ForegroundColor DarkGray
            }
            
            return $response.value
        } else {
            # MATCH mode - find specific match
            $matchingDef = $null
            if ($CategoryPath) {
                # Filter by category path
                $matchingDef = $response.value | Where-Object { $_.categoryPath -eq $CategoryPath } | Select-Object -First 1
                if (-not $matchingDef) {
                    throw "No definition found matching category path '$CategoryPath'"
                }
            } elseif ($response.value.Count -gt 1) {
                Write-Host "⚠️  Multiple definitions found, using first match" -ForegroundColor Yellow
                $matchingDef = $response.value[0]
            } else {
                $matchingDef = $response.value[0]
            }
            
            Write-Host "✅ Matched Definition Template ID: $($matchingDef.id)" -ForegroundColor Green
            Write-Host "   Display Name: $($matchingDef.displayName)" -ForegroundColor Gray
            Write-Host "   Class Type: $($matchingDef.classType)" -ForegroundColor Gray
            Write-Host "   Category Path: $($matchingDef.categoryPath)" -ForegroundColor Gray
            Write-Host "   Policy Type: $($matchingDef.policyType)" -ForegroundColor Gray
            
            return $matchingDef
        }
    }
    catch {
        Write-Host "❌ Failed: $_" -ForegroundColor Red
        throw
    }
}

# Step 2: Get presentation templates and resolve to presentation template ID
function Resolve-PresentationTemplateID {
    param (
        [string]$DefinitionTemplateID,
        [int]$PresentationIndex,
        [string]$OutputFile
    )
    
    try {
        Write-Host "`n📍 STEP 2: Get Presentations → Presentation Template ID" -ForegroundColor Yellow
        Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Yellow
        
        $url = "https://graph.microsoft.com/beta/deviceManagement/groupPolicyDefinitions/$DefinitionTemplateID/presentations"
        
        Write-Host "Input: Definition Template ID='$DefinitionTemplateID'" -ForegroundColor Gray
        Write-Host "Input: Presentation Index=$PresentationIndex" -ForegroundColor Gray
        Write-Host "API Call: $url" -ForegroundColor DarkGray
        
        $response = Invoke-MgGraphRequest -Method GET -Uri $url
        
        # Save response
        $prettyJson = ConvertTo-Json -InputObject $response -Depth 10 -Compress:$false
        $prettyJson | Out-File -FilePath $OutputFile -Encoding utf8
        
        if (-not $response.value -or $response.value.Count -eq 0) {
            Write-Host "⚠️  No presentations found (this is valid for some policies)" -ForegroundColor Yellow
            return $null
        }
        
        if ($PresentationIndex -ge $response.value.Count) {
            throw "Presentation index $PresentationIndex out of range (found $($response.value.Count) presentations)"
        }
        
        $presentation = $response.value[$PresentationIndex]
        
        Write-Host "✅ Found Presentation Template ID: $($presentation.id)" -ForegroundColor Green
        Write-Host "   Label: $($presentation.label)" -ForegroundColor Gray
        Write-Host "   OData Type: $($presentation.'@odata.type')" -ForegroundColor Gray
        Write-Host "   Total Presentations Available: $($response.value.Count)" -ForegroundColor Gray
        
        return $presentation
    }
    catch {
        Write-Host "❌ Failed to resolve presentation template ID: $_" -ForegroundColor Red
        throw
    }
}

# Step 3: Get all group policy configurations
function Get-AllConfigurations {
    param (
        [string]$OutputFile
    )
    
    try {
        Write-Host "`n📍 STEP 3: Get All Group Policy Configurations" -ForegroundColor Yellow
        Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Yellow
        
        $url = "https://graph.microsoft.com/beta/deviceManagement/groupPolicyConfigurations"
        
        Write-Host "API Call: $url" -ForegroundColor DarkGray
        
        $response = Invoke-MgGraphRequest -Method GET -Uri $url
        
        # Save response
        $prettyJson = ConvertTo-Json -InputObject $response -Depth 10 -Compress:$false
        $prettyJson | Out-File -FilePath $OutputFile -Encoding utf8
        
        $configCount = if ($response.value) { $response.value.Count } else { 0 }
        Write-Host "✅ Found $configCount configuration(s)" -ForegroundColor Green
        
        return $response.value
    }
    catch {
        Write-Host "❌ Failed to get configurations: $_" -ForegroundColor Red
        throw
    }
}

# Step 4: Resolve template IDs to instance IDs
function Resolve-InstanceIDs {
    param (
        [string]$ConfigurationID,
        [string]$ConfigurationName,
        [string]$DefinitionTemplateID,
        [string]$PresentationTemplateID,
        [string]$DefinitionValuesOutputFile,
        [string]$PresentationValuesOutputFile
    )
    
    try {
        Write-Host "`n  🔍 Checking configuration: $ConfigurationName" -ForegroundColor Cyan
        
        # Get definition values with expanded definition
        $url = "https://graph.microsoft.com/beta/deviceManagement/groupPolicyConfigurations/$ConfigurationID/definitionValues?`$expand=definition(`$select=id,classType,displayName,policyType,hasRelatedDefinitions,version,minUserCspVersion,minDeviceCspVersion)"
        
        Write-Host "     API Call: $url" -ForegroundColor DarkGray
        
        $defValuesResponse = Invoke-MgGraphRequest -Method GET -Uri $url
        
        # Save definition values response
        $prettyJson = ConvertTo-Json -InputObject $defValuesResponse -Depth 10 -Compress:$false
        $prettyJson | Out-File -FilePath $DefinitionValuesOutputFile -Encoding utf8
        
        # Find matching definition value
        $matchingDefValue = $null
        foreach ($defValue in $defValuesResponse.value) {
            if ($defValue.definition -and $defValue.definition.id -eq $DefinitionTemplateID) {
                $matchingDefValue = $defValue
                break
            }
        }
        
        if (-not $matchingDefValue) {
            Write-Host "     ℹ️  No instance found for this policy in this configuration" -ForegroundColor DarkGray
            return $null
        }
        
        Write-Host "     ✅ MATCH FOUND! Definition Value Instance ID: $($matchingDefValue.id)" -ForegroundColor Green
        
        # Get presentation values
        $presValuesUrl = "https://graph.microsoft.com/beta/deviceManagement/groupPolicyConfigurations/$ConfigurationID/definitionValues/$($matchingDefValue.id)/presentationValues"
        
        Write-Host "     API Call: $presValuesUrl" -ForegroundColor DarkGray
        
        $presValuesResponse = Invoke-MgGraphRequest -Method GET -Uri $presValuesUrl
        
        # Save presentation values response
        $prettyJson = ConvertTo-Json -InputObject $presValuesResponse -Depth 10 -Compress:$false
        $prettyJson | Out-File -FilePath $PresentationValuesOutputFile -Encoding utf8
        
        # Find matching presentation value
        $matchingPresValue = $null
        if ($PresentationTemplateID) {
            foreach ($presValue in $presValuesResponse.value) {
                # Match by OData type since we might not have direct presentation ID binding
                $matchingPresValue = $presValue
                break # For now, take first one
            }
        }
        
        if ($matchingPresValue) {
            Write-Host "     ✅ Presentation Value Instance ID: $($matchingPresValue.id)" -ForegroundColor Green
        }
        
        return @{
            ConfigurationID = $ConfigurationID
            ConfigurationName = $ConfigurationName
            DefinitionValueInstanceID = $matchingDefValue.id
            PresentationValueInstanceID = if ($matchingPresValue) { $matchingPresValue.id } else { $null }
            Enabled = $matchingDefValue.enabled
            DefinitionValue = $matchingDefValue
            PresentationValues = $presValuesResponse.value
        }
    }
    catch {
        Write-Host "     ⚠️  Error checking configuration: $_" -ForegroundColor Yellow
        return $null
    }
}

# Main execution
try {
    # Ensure output directory exists
    if (-not (Test-Path -Path $OutputDirectory)) {
        New-Item -ItemType Directory -Path $OutputDirectory -Force | Out-Null
    }
    
    Write-Host "`n╔═══════════════════════════════════════════════════════════════╗" -ForegroundColor Cyan
    Write-Host "║   Group Policy ID Resolution - id_resolver.go Replication   ║" -ForegroundColor Cyan
    Write-Host "╚═══════════════════════════════════════════════════════════════╝" -ForegroundColor Cyan
    
    # Determine mode
    $isListMode = [string]::IsNullOrWhiteSpace($PolicyName)
    
    # Connect
    Connect-MicrosoftGraph -TenantId $TenantId -ClientId $ClientId -ClientSecret $ClientSecret
    
    # Step 1: Get definitions
    $step1File = Join-Path -Path $OutputDirectory -ChildPath "step1_definitions.json"
    $definitionsResult = Get-GroupPolicyDefinitions -PolicyName $PolicyName -ClassType $ClassType -CategoryPath $CategoryPath -OutputFile $step1File
    
    if ($isListMode) {
        # LIST mode - just show the catalog and exit
        Write-Host "`n╔═══════════════════════════════════════════════════════════════╗" -ForegroundColor Green
        Write-Host "║                    CATALOG LISTING                            ║" -ForegroundColor Green
        Write-Host "╚═══════════════════════════════════════════════════════════════╝" -ForegroundColor Green
        
        Write-Host "`n📋 Total Definitions: $($definitionsResult.Count)" -ForegroundColor Cyan
        
        if ($ClassType) {
            $userCount = ($definitionsResult | Where-Object { $_.classType -eq 'user' }).Count
            $machineCount = ($definitionsResult | Where-Object { $_.classType -eq 'machine' }).Count
            Write-Host "   • User policies: $userCount" -ForegroundColor Gray
            Write-Host "   • Machine policies: $machineCount" -ForegroundColor Gray
        }
        
        # Group by policy type
        $byType = $definitionsResult | Group-Object -Property policyType
        Write-Host "`n📊 By Policy Type:" -ForegroundColor Cyan
        foreach ($group in $byType) {
            Write-Host "   • $($group.Name): $($group.Count)" -ForegroundColor Gray
        }
        
        Write-Host "`n📁 Output File:" -ForegroundColor Cyan
        Write-Host "   step1_definitions.json - Complete catalog with $($definitionsResult.Count) definitions" -ForegroundColor Gray
        
        Write-Host "`n✨ Catalog listing complete! Use -PolicyName to resolve specific policy." -ForegroundColor Green
        
        # Exit early - finally block will handle disconnect
        return
    }
    
    # MATCH mode - continue with resolution
    $definition = $definitionsResult
    $definitionTemplateID = $definition.id
    
    # Step 2: Resolve presentation template ID
    $step2File = Join-Path -Path $OutputDirectory -ChildPath "step2_presentation_templates.json"
    $presentation = Resolve-PresentationTemplateID -DefinitionTemplateID $definitionTemplateID -PresentationIndex $PresentationIndex -OutputFile $step2File
    $presentationTemplateID = if ($presentation) { $presentation.id } else { $null }
    
    # Get presentation details for summary
    $presentationsResponse = Invoke-MgGraphRequest -Uri "https://graph.microsoft.com/beta/deviceManagement/groupPolicyDefinitions/$definitionTemplateID/presentations" -Method GET
    $allPresentations = $presentationsResponse.value
    
    # Final summary
    Write-Host "`n╔═══════════════════════════════════════════════════════════════╗" -ForegroundColor Green
    Write-Host "║                    CATALOG LOOKUP COMPLETE                    ║" -ForegroundColor Green
    Write-Host "╚═══════════════════════════════════════════════════════════════╝" -ForegroundColor Green
    
    Write-Host "`n📋 Policy Definition:" -ForegroundColor Cyan
    Write-Host "   Display Name:     $($definition.displayName)" -ForegroundColor White
    Write-Host "   Definition ID:    $definitionTemplateID" -ForegroundColor White
    Write-Host "   Class Type:       $($definition.classType)" -ForegroundColor White
    Write-Host "   Category Path:    $($definition.categoryPath)" -ForegroundColor White
    Write-Host "   Policy Type:      $($definition.policyType)" -ForegroundColor White
    
    if ($allPresentations -and $allPresentations.Count -gt 0) {
        Write-Host "`n📋 Presentation Types (Value Types Supported):" -ForegroundColor Cyan
        foreach ($pres in $allPresentations) {
            $presType = $pres.'@odata.type' -replace '#microsoft.graph.groupPolicyPresentation', ''
            Write-Host "   • $($pres.label)" -ForegroundColor White
            Write-Host "     Type: $presType" -ForegroundColor Gray
            Write-Host "     ID: $($pres.id)" -ForegroundColor Gray
            Write-Host "     Required: $($pres.required)" -ForegroundColor Gray
        }
    } else {
        Write-Host "`n⚠️  No presentations found (this policy has no configurable values)" -ForegroundColor Yellow
    }
    
    Write-Host "`n📁 Output Files:" -ForegroundColor Cyan
    Write-Host "   step1_definitions.json - Policy definition from catalog" -ForegroundColor Gray
    Write-Host "   step2_presentation_templates.json - Value types supported" -ForegroundColor Gray
    
    Write-Host "`n✨ Catalog lookup complete! Files saved to: $OutputDirectory" -ForegroundColor Green
}
catch {
    Write-Host "`n❌ Script execution failed: $_" -ForegroundColor Red
    exit 1
}
finally {
    Write-Host "`n🔓 Disconnecting from Microsoft Graph..." -ForegroundColor Cyan
    Disconnect-MgGraph | Out-Null
    Write-Host "✅ Disconnected" -ForegroundColor Green
}

