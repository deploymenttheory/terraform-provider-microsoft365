<#
.SYNOPSIS
    Exports Group Policy Definition configurations to HCL for Terraform import

.DESCRIPTION
    This script helps identify and export Group Policy Definition resources from Microsoft Intune
    for importing into Terraform. It performs multi-endpoint resolution to gather complete metadata:
    
    1. Lists all Group Policy Configurations in your tenant
    2. For a selected configuration, retrieves all Definition Values with presentation values
    3. Queries the full policy definition details (policy name, class type, category path)
    4. Generates Terraform import commands and HCL resource blocks
    
    The script mirrors the import logic from the Terraform provider by:
    - Using $expand=definition to get basic definition metadata
    - Making a separate call to /groupPolicyDefinitions/{id} to get categoryPath
    - Fetching presentation values with their types and labels
    
    This ensures the exported HCL matches exactly what the provider expects for successful imports.

.PARAMETER TenantId
    Azure AD Tenant ID

.PARAMETER ClientId
    Application (Client) ID of the Azure AD app registration

.PARAMETER ClientSecret
    Client secret for authentication

.PARAMETER ConfigurationId
    Optional: Specific Group Policy Configuration ID to export from.
    If not provided, the script will list all configurations for selection.

.PARAMETER PolicyName
    Optional: Filter to specific policy definition by display name.
    Supports wildcards (e.g., "*Microsoft Edge*").
    If not provided, all policies in the configuration will be exported.

.PARAMETER OutputDirectory
    Directory where export files will be saved (import commands and HCL)
    Default: Current directory

.PARAMETER GenerateHCL
    Generate Terraform HCL resource blocks in addition to import commands
    Default: true

.EXAMPLE
    # Interactive mode - select configuration from list
    pwsh Export-GroupPolicyDefinitionToHCLForImport.ps1 `
        -TenantId "your-tenant-id" `
        -ClientId "your-client-id" `
        -ClientSecret "your-secret"
    
    Prompts user to select a configuration, then exports all policy definitions

.EXAMPLE
    # Direct mode - specify configuration ID
    pwsh Export-GroupPolicyDefinitionToHCLForImport.ps1 `
        -TenantId "your-tenant-id" `
        -ClientId "your-client-id" `
        -ClientSecret "your-secret" `
        -ConfigurationId "a1b2c3d4-e5f6-g7h8-i9j0-k1l2m3n4o5p6"
    
    Exports all policy definitions from the specified configuration

.EXAMPLE
    # Export a specific policy by exact name
    pwsh Export-GroupPolicyDefinitionToHCLForImport.ps1 `
        -TenantId "your-tenant-id" `
        -ClientId "your-client-id" `
        -ClientSecret "your-secret" `
        -ConfigurationId "a1b2c3d4-e5f6-g7h8-i9j0-k1l2m3n4o5p6" `
        -PolicyName "Browsing Data Lifetime Settings"
    
    Exports only the policy with the exact display name

.EXAMPLE
    # Export policies matching a wildcard pattern
    pwsh Export-GroupPolicyDefinitionToHCLForImport.ps1 `
        -TenantId "your-tenant-id" `
        -ClientId "your-client-id" `
        -ClientSecret "your-secret" `
        -ConfigurationId "a1b2c3d4-e5f6-g7h8-i9j0-k1l2m3n4o5p6" `
        -PolicyName "*Microsoft Edge*"
    
    Exports all policies containing "Microsoft Edge" in their display name

.EXAMPLE
    # Generate import commands only (no HCL)
    pwsh Export-GroupPolicyDefinitionToHCLForImport.ps1 `
        -TenantId "00000000-0000-0000-0000-000000000000" `
        -ClientId "00000000-0000-0000-0000-000000000000" `
        -ClientSecret "your-secret" `
        -GenerateHCL $false

.NOTES
    File Name      : Export-GroupPolicyDefinitionToHCLForImport.ps1
    Prerequisite   : Microsoft.Graph.Authentication PowerShell module
    
    API Permissions Required:
    - DeviceManagementConfiguration.Read.All (minimum)

.LINK
    https://learn.microsoft.com/en-us/graph/api/intune-grouppolicy-grouppolicyconfiguration-list
    https://learn.microsoft.com/en-us/graph/api/intune-grouppolicy-grouppolicydefinitionvalue-list
#>

param(
    [Parameter(Mandatory = $true, HelpMessage = "Azure AD Tenant ID")]
    [string]$TenantId,

    [Parameter(Mandatory = $true, HelpMessage = "Application (Client) ID")]
    [string]$ClientId,

    [Parameter(Mandatory = $true, HelpMessage = "Client Secret")]
    [string]$ClientSecret,

    [Parameter(Mandatory = $false, HelpMessage = "Specific Group Policy Configuration ID")]
    [string]$ConfigurationId,

    [Parameter(Mandatory = $false, HelpMessage = "Filter to specific policy by display name (supports wildcards)")]
    [string]$PolicyName,

    [Parameter(Mandatory = $false, HelpMessage = "Output directory for export files")]
    [string]$OutputDirectory = ".",

    [Parameter(Mandatory = $false, HelpMessage = "Generate Terraform HCL resource blocks")]
    [bool]$GenerateHCL = $true
)

# Function to connect to Microsoft Graph
function Connect-MicrosoftGraph {
    param(
        [string]$TenantId,
        [string]$ClientId,
        [string]$ClientSecret
    )

    Write-Host "🔐 Connecting to Microsoft Graph..." -ForegroundColor Cyan
    
    try {
        $secureSecret = ConvertTo-SecureString $ClientSecret -AsPlainText -Force
        $credential = New-Object System.Management.Automation.PSCredential($ClientId, $secureSecret)
        
        Connect-MgGraph -ClientSecretCredential $credential -TenantId $TenantId -NoWelcome | Out-Null
        
        Write-Host "✅ Connected successfully" -ForegroundColor Green
    }
    catch {
        Write-Host "❌ Failed to connect: $_" -ForegroundColor Red
        exit 1
    }
}

# Function to get all group policy configurations
function Get-AllGroupPolicyConfigurations {
    Write-Host "`n📋 Retrieving Group Policy Configurations..." -ForegroundColor Cyan
    
    $baseUrl = "https://graph.microsoft.com/beta"
    $configurations = @()
    $nextLink = "$baseUrl/deviceManagement/groupPolicyConfigurations"
    
    do {
        $response = Invoke-MgGraphRequest -Uri $nextLink -Method GET
        $configurations += $response.value
        $nextLink = $response.'@odata.nextLink'
    } while ($nextLink)
    
    return $configurations
}

# Function to get definition values for a configuration with multi-endpoint resolution
function Get-DefinitionValuesForConfiguration {
    param(
        [string]$ConfigurationId
    )
    
    Write-Host "`n🔍 Retrieving Definition Values with multi-endpoint resolution..." -ForegroundColor Cyan
    Write-Host "   (Mirroring Terraform provider import logic)" -ForegroundColor Gray
    
    $baseUrl = "https://graph.microsoft.com/beta"
    $definitionValues = @()
    
    # Step 1: Get definition values with expanded definition
    # Note: $expand=definition is required - without it, the definition object is completely missing
    Write-Host "`n   Step 1: GET definitionValues?`$expand=definition" -ForegroundColor DarkGray
    $nextLink = "$baseUrl/deviceManagement/groupPolicyConfigurations/$ConfigurationId/definitionValues?`$expand=definition"
    
    do {
        $response = Invoke-MgGraphRequest -Uri $nextLink -Method GET
        $definitionValues += $response.value
        $nextLink = $response.'@odata.nextLink'
    } while ($nextLink)
    
    Write-Host "   ✅ Found $($definitionValues.Count) definition value(s)" -ForegroundColor Green
    
    # Step 2: For each definition value, get the FULL definition (to get categoryPath)
    # The expanded definition has categoryPath as null, so we need to query the catalog directly
    Write-Host "`n   Step 2: GET groupPolicyDefinitions/{id} for each (to get categoryPath)" -ForegroundColor DarkGray
    
    foreach ($defValue in $definitionValues) {
        if ($defValue.definition -and $defValue.definition.id) {
            $definitionId = $defValue.definition.id
            
            try {
                $fullDefUri = "$baseUrl/deviceManagement/groupPolicyDefinitions/$definitionId"
                $fullDefinition = Invoke-MgGraphRequest -Uri $fullDefUri -Method GET
                
                # Replace the expanded definition with the full definition (includes categoryPath)
                $defValue.definition = $fullDefinition
            }
            catch {
                Write-Host "      ⚠️  Could not retrieve full definition for $definitionId" -ForegroundColor Yellow
            }
        }
        
        # Step 3: Fetch presentation values with expanded presentations
        $presValuesUri = "$baseUrl/deviceManagement/groupPolicyConfigurations/$ConfigurationId/definitionValues/$($defValue.id)/presentationValues?`$expand=presentation"
        try {
            $presResponse = Invoke-MgGraphRequest -Uri $presValuesUri -Method GET
            $defValue | Add-Member -MemberType NoteProperty -Name 'presentationValues' -Value $presResponse.value -Force
        }
        catch {
            Write-Host "      ⚠️  Could not retrieve presentation values for $($defValue.definition.displayName)" -ForegroundColor Yellow
            $defValue | Add-Member -MemberType NoteProperty -Name 'presentationValues' -Value @() -Force
        }
    }
    
    Write-Host "   ✅ Multi-endpoint resolution complete" -ForegroundColor Green
    
    return $definitionValues
}

# Function to determine presentation type display name
function Get-PresentationTypeDisplay {
    param($ODataType)
    
    $typeMap = @{
        '#microsoft.graph.groupPolicyPresentationValueBoolean' = 'CheckBox (Boolean)'
        '#microsoft.graph.groupPolicyPresentationValueText' = 'TextBox'
        '#microsoft.graph.groupPolicyPresentationValueDecimal' = 'Decimal'
        '#microsoft.graph.groupPolicyPresentationValueMultiText' = 'MultiText'
        '#microsoft.graph.groupPolicyPresentationValueList' = 'List/Dropdown'
    }
    
    if ($typeMap.ContainsKey($ODataType)) {
        return $typeMap[$ODataType]
    }
    return $ODataType
}

# Function to generate Terraform import command
function New-TerraformImportCommand {
    param(
        [string]$ResourceName,
        [string]$ConfigurationId,
        [string]$DefinitionValueId
    )
    
    # Use composite ID format: configID/definitionValueID
    $compositeId = "$ConfigurationId/$DefinitionValueId"
    return "terraform import microsoft365_graph_beta_device_management_group_policy_definition.$ResourceName `"$compositeId`""
}

# Function to format value for HCL based on presentation type
function Format-HCLValue {
    param($PresentationValue)
    
    $odataType = $PresentationValue.'@odata.type'
    
    switch ($odataType) {
        '#microsoft.graph.groupPolicyPresentationValueBoolean' {
            if ($PresentationValue.value) { 
                return "`"true`"" 
            } else { 
                return "`"false`"" 
            }
        }
        '#microsoft.graph.groupPolicyPresentationValueText' {
            # Escape quotes in the value
            $escaped = $PresentationValue.value -replace '"', '\"'
            return "`"$escaped`""
        }
        '#microsoft.graph.groupPolicyPresentationValueDecimal' {
            return $PresentationValue.value.ToString()
        }
        '#microsoft.graph.groupPolicyPresentationValueMultiText' {
            if ($PresentationValue.values -and $PresentationValue.values.Count -gt 0) {
                $escaped = $PresentationValue.values[0] -replace '"', '\"'
                return "`"$escaped`""
            }
            return '""'
        }
        '#microsoft.graph.groupPolicyPresentationValueList' {
            if ($PresentationValue.values -and $PresentationValue.values.Count -gt 0) {
                # For list values, use the first value's name or value
                if ($PresentationValue.values[0].name) {
                    $escaped = $PresentationValue.values[0].name -replace '"', '\"'
                    return "`"$escaped`""
                }
                elseif ($PresentationValue.values[0].value) {
                    return $PresentationValue.values[0].value.ToString()
                }
            }
            return '""'
        }
        default {
            return '""'
        }
    }
}

# Function to generate Terraform HCL resource block
function New-TerraformHCL {
    param(
        [string]$ResourceName,
        [string]$ConfigurationId,
        $DefinitionValue
    )
    
    $definition = $DefinitionValue.definition
    $policyName = $definition.displayName -replace '"', '\"'
    $classType = $definition.classType
    $categoryPath = $definition.categoryPath -replace '"', '\"'
    $enabled = if ($DefinitionValue.enabled) { "true" } else { "false" }
    
    # Build values array from presentation values
    $valuesHCL = @()
    
    if ($DefinitionValue.presentationValues -and $DefinitionValue.presentationValues.Count -gt 0) {
        foreach ($presValue in $DefinitionValue.presentationValues) {
            if ($presValue.presentation -and $presValue.presentation.label) {
                $label = $presValue.presentation.label -replace '"', '\"'
                $value = Format-HCLValue -PresentationValue $presValue
                
                $valuesHCL += @"
    {
      label = "$label"
      value = $value
    }
"@
            }
        }
    }
    
    $valuesBlock = if ($valuesHCL.Count -gt 0) { 
        $valuesHCL -join ",`n" 
    } else { 
        "    # No presentation values configured" 
    }
    
    $hcl = @"
resource "microsoft365_graph_beta_device_management_group_policy_definition" "$ResourceName" {
  group_policy_configuration_id = "$ConfigurationId"
  policy_name                   = "$policyName"
  class_type                    = "$classType"
  category_path                 = "$categoryPath"
  enabled                       = $enabled

  values = [
$valuesBlock
  ]

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}
"@
    
    return $hcl
}

# Main execution
try {
    Write-Host "`n╔═══════════════════════════════════════════════════════════════╗" -ForegroundColor Cyan
    Write-Host "║  Export Group Policy Definitions to HCL for Import           ║" -ForegroundColor Cyan
    Write-Host "╚═══════════════════════════════════════════════════════════════╝" -ForegroundColor Cyan
    
    # Ensure output directory exists
    if (-not (Test-Path -Path $OutputDirectory)) {
        New-Item -ItemType Directory -Path $OutputDirectory -Force | Out-Null
    }
    
    # Connect to Microsoft Graph
    Connect-MicrosoftGraph -TenantId $TenantId -ClientId $ClientId -ClientSecret $ClientSecret
    
    # Get configuration ID (either from parameter or user selection)
    if (-not $ConfigurationId) {
        $configurations = Get-AllGroupPolicyConfigurations
        
        if ($configurations.Count -eq 0) {
            Write-Host "`n❌ No Group Policy Configurations found in tenant" -ForegroundColor Red
            exit 1
        }
        
        Write-Host "`n📋 Available Group Policy Configurations:" -ForegroundColor Cyan
        Write-Host ("═" * 80) -ForegroundColor DarkGray
        
        for ($i = 0; $i -lt $configurations.Count; $i++) {
            Write-Host "$($i + 1). $($configurations[$i].displayName)" -ForegroundColor White
            Write-Host "   ID: $($configurations[$i].id)" -ForegroundColor Gray
            Write-Host "   Description: $($configurations[$i].description)" -ForegroundColor Gray
            Write-Host ""
        }
        
        $selection = Read-Host "`nSelect configuration number (1-$($configurations.Count))"
        $selectedIndex = [int]$selection - 1
        
        if ($selectedIndex -lt 0 -or $selectedIndex -ge $configurations.Count) {
            Write-Host "`n❌ Invalid selection" -ForegroundColor Red
            exit 1
        }
        
        $selectedConfig = $configurations[$selectedIndex]
        $ConfigurationId = $selectedConfig.id
        $configName = $selectedConfig.displayName
    }
    else {
        # Get configuration details
        $configUri = "https://graph.microsoft.com/beta/deviceManagement/groupPolicyConfigurations/$ConfigurationId"
        $selectedConfig = Invoke-MgGraphRequest -Uri $configUri -Method GET
        $configName = $selectedConfig.displayName
    }
    
    Write-Host "`n✅ Selected Configuration: $configName" -ForegroundColor Green
    Write-Host "   ID: $ConfigurationId" -ForegroundColor Gray
    
    # Get all definition values for this configuration with multi-endpoint resolution
    $definitionValues = Get-DefinitionValuesForConfiguration -ConfigurationId $ConfigurationId
    
    if ($definitionValues.Count -eq 0) {
        Write-Host "`n⚠️  No policy definitions found in this configuration" -ForegroundColor Yellow
        exit 0
    }
    
    # Filter by PolicyName if specified
    if ($PolicyName) {
        Write-Host "`n🔍 Filtering policies matching: $PolicyName" -ForegroundColor Cyan
        $definitionValues = $definitionValues | Where-Object { $_.definition.displayName -like $PolicyName }
        
        if ($definitionValues.Count -eq 0) {
            Write-Host "`n⚠️  No policy definitions found matching: $PolicyName" -ForegroundColor Yellow
            exit 0
        }
        
        Write-Host "✅ Found $($definitionValues.Count) matching policy definition(s)" -ForegroundColor Green
    }
    else {
        Write-Host "`n✅ Found $($definitionValues.Count) policy definition(s) (exporting all)" -ForegroundColor Green
    }
    
    # Display summary of found definitions
    Write-Host "`n📊 Definition Summary:" -ForegroundColor Cyan
    foreach ($defValue in $definitionValues) {
        $def = $defValue.definition
        $presCount = if ($defValue.presentationValues) { $defValue.presentationValues.Count } else { 0 }
        $types = if ($defValue.presentationValues) { 
            ($defValue.presentationValues | ForEach-Object { Get-PresentationTypeDisplay $_.'@odata.type' } | Select-Object -Unique) -join ', '
        } else { 
            'None' 
        }
        
        Write-Host "   • $($def.displayName)" -ForegroundColor White
        Write-Host "     Class: $($def.classType) | Category: $($def.categoryPath)" -ForegroundColor Gray
        Write-Host "     Values: $presCount ($types)" -ForegroundColor Gray
    }
    
    # Generate export files
    $timestamp = Get-Date -Format "yyyyMMdd_HHmmss"
    $safeConfigName = $configName -replace '[^a-zA-Z0-9_]', '_' -replace '__+', '_'
    $importFile = Join-Path $OutputDirectory "import_${safeConfigName}_$timestamp.sh"
    $hclFile = Join-Path $OutputDirectory "resources_${safeConfigName}_$timestamp.tf"
    
    $importCommands = @()
    $hclBlocks = @()
    
    Write-Host "`n🔄 Generating export files..." -ForegroundColor Cyan
    
    $counter = 1
    foreach ($defValue in $definitionValues) {
        $definition = $defValue.definition
        $policyName = $definition.displayName
        $definitionValueId = $defValue.id
        
        # Create safe resource name
        $resourceName = $policyName -replace '[^a-zA-Z0-9_]', '_' -replace '__+', '_'
        $resourceName = $resourceName.ToLower().Trim('_')
        
        # Add counter suffix if needed to ensure uniqueness
        if ($counter -gt 1) {
            $resourceName = "${resourceName}_$counter"
        }
        $counter++
        
        Write-Host "   $counter. $policyName" -ForegroundColor White
        
        # Generate import command with composite ID format
        $importCmd = New-TerraformImportCommand -ResourceName $resourceName `
                                                 -ConfigurationId $ConfigurationId `
                                                 -DefinitionValueId $definitionValueId
        $importCommands += $importCmd
        
        # Generate HCL if requested
        if ($GenerateHCL) {
            $hcl = New-TerraformHCL -ResourceName $resourceName `
                                     -ConfigurationId $ConfigurationId `
                                     -DefinitionValue $defValue
            $hclBlocks += $hcl
        }
    }
    
    # Write import commands
    $importContent = @"
#!/bin/bash
# Generated: $(Get-Date -Format "yyyy-MM-dd HH:mm:ss")
# Configuration: $configName
# Configuration ID: $ConfigurationId
# Total Resources: $($definitionValues.Count)
#
# Import ID Format: configID/definitionValueID (composite format)
# This matches the resource's internal ID tracking for proper state management

$($importCommands -join "`n")
"@
    
    Set-Content -Path $importFile -Value $importContent -Encoding UTF8
    Write-Host "`n✅ Import commands saved to: $importFile" -ForegroundColor Green
    
    # Write HCL if generated
    if ($GenerateHCL -and $hclBlocks.Count -gt 0) {
        $hclContent = @"
# Generated: $(Get-Date -Format "yyyy-MM-dd HH:mm:ss")
# Configuration: $configName
# Configuration ID: $ConfigurationId
# Total Resources: $($definitionValues.Count)
#
# IMPORTANT: These resources use the new microsoft365_graph_beta_device_management_group_policy_definition
# resource which supports all presentation types (Boolean, TextBox, Decimal, MultiText, Dropdown)
#
# The import was generated using multi-endpoint resolution to ensure:
# - policy_name, class_type, and category_path are populated correctly
# - All presentation values are captured with their labels and types
#
# NOTE: Review and update the group_policy_configuration_id references to match your configuration resource

$($hclBlocks -join "`n`n")
"@
        
        Set-Content -Path $hclFile -Value $hclContent -Encoding UTF8
        Write-Host "✅ HCL resources saved to: $hclFile" -ForegroundColor Green
    }
    
    Write-Host "`n╔═══════════════════════════════════════════════════════════════╗" -ForegroundColor Green
    Write-Host "║                    EXPORT COMPLETE                            ║" -ForegroundColor Green
    Write-Host "╚═══════════════════════════════════════════════════════════════╝" -ForegroundColor Green
    
    Write-Host "`n📋 Summary:" -ForegroundColor Cyan
    Write-Host "   • Configuration: $configName" -ForegroundColor White
    Write-Host "   • Policy Definitions Found: $($definitionValues.Count)" -ForegroundColor White
    Write-Host "   • Import Commands: $importFile" -ForegroundColor White
    if ($GenerateHCL) {
        Write-Host "   • HCL Resources: $hclFile" -ForegroundColor White
    }
    
    Write-Host "`n📝 Next Steps:" -ForegroundColor Cyan
    Write-Host "   1. Review the generated HCL in: $hclFile" -ForegroundColor Gray
    Write-Host "   2. Update the group_policy_configuration_id references" -ForegroundColor Gray
    Write-Host "   3. Run the import commands: bash $importFile" -ForegroundColor Gray
    Write-Host "   4. Run terraform plan to verify" -ForegroundColor Gray
    
    Write-Host "`n💡 Tip: The import uses composite ID format (configID/definitionValueID)" -ForegroundColor Yellow
    Write-Host "   This matches how the Terraform provider tracks resources internally" -ForegroundColor Gray
}
catch {
    Write-Host "`n❌ Script execution failed: $_" -ForegroundColor Red
    Write-Host $_.ScriptStackTrace -ForegroundColor Red
    exit 1
}
finally {
    Write-Host "`n🔓 Disconnecting from Microsoft Graph..." -ForegroundColor Cyan
    Disconnect-MgGraph | Out-Null
    Write-Host "✅ Disconnected" -ForegroundColor Green
}

