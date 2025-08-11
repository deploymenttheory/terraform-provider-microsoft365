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
    HelpMessage="Specify the ID of the settings catalog policy to retrieve")]
    [ValidateNotNullOrEmpty()]
    [string]$SettingsCatalogItemId,
    
    [Parameter(Mandatory=$false,
    HelpMessage="Export all settings catalog policies (alternative to providing -SettingsCatalogItemId)")]
    [bool]$ExportAll = $false,
    
    [Parameter(Mandatory=$false,
    HelpMessage="Export results to JSON file (optional)")]
    [bool]$ExportToJson = $false,
    
    [Parameter(Mandatory=$false,
    HelpMessage="Export results to HCL/Terraform file (optional)")]
    [bool]$ExportToHcl = $false,
    
    [Parameter(Mandatory=$false,
    HelpMessage="Enable verbose debug output")]
    [bool]$EnableDebug = $false
)

# Usage Examples:
# .\Get-SettingsCatalogPolicy.ps1 -TenantId "your-tenant-id" -ClientId "your-client-id" -ClientSecret "your-client-secret" -SettingsCatalogItemId "policy-id"
# .\Get-SettingsCatalogPolicy.ps1 -TenantId "your-tenant-id" -ClientId "your-client-id" -ClientSecret "your-client-secret" -SettingsCatalogItemId "policy-id" -ExportToJson $true
# .\Get-SettingsCatalogPolicy.ps1 -TenantId "your-tenant-id" -ClientId "your-client-id" -ClientSecret "your-client-secret" -SettingsCatalogItemId "policy-id" -ExportToHcl $true
# .\Get-SettingsCatalogPolicy.ps1 -TenantId "your-tenant-id" -ClientId "your-client-id" -ClientSecret "your-client-secret" -SettingsCatalogItemId "policy-id" -ExportToJson $true -ExportToHcl $true
# .\Get-SettingsCatalogPolicy.ps1 -TenantId "your-tenant-id" -ClientId "your-client-id" -ClientSecret "your-client-secret" -SettingsCatalogItemId "policy-id" -ExportToHcl $true -EnableDebug $true
# .\Get-SettingsCatalogPolicy.ps1 -TenantId "your-tenant-id" -ClientId "your-client-id" -ClientSecret "your-client-secret" -ExportAll $true -ExportToJson $true -ExportToHcl $true

Import-Module Microsoft.Graph.Authentication

function Write-DebugInfo {
    param (
        [Parameter(Mandatory=$true)]
        [string]$Message,
        [Parameter(Mandatory=$false)]
        [object]$Data = $null
    )
    
    if ($EnableDebug) {
        Write-Host "🐛 DEBUG: $Message" -ForegroundColor Yellow
        if ($Data) {
            Write-Host "   Data: $($Data | ConvertTo-Json -Depth 2 -Compress)" -ForegroundColor Gray
        }
    }
}

function Escape-HclString {
    param(
        [Parameter(Mandatory=$true)]
        [string]$InputString
    )
    # Build escaped string so that:
    # - Each single backslash \ becomes \\ in HCL (escaped once)
    # - Double quotes are escaped as \"
    if ($null -eq $InputString) { return "" }
    $builder = New-Object System.Text.StringBuilder
    $length = $InputString.Length
    
    for ($i = 0; $i -lt $length; $i++) {
        $ch = $InputString[$i]
        if ($ch -eq '"') {
            [void]$builder.Append('\"')
        } elseif ($ch -eq '\') {
            [void]$builder.Append('\\')
        } else {
            [void]$builder.Append($ch)
        }
    }
    return $builder.ToString()
}

function Get-PaginatedResults {
    param (
        [Parameter(Mandatory=$true)]
        [string]$InitialUri
    )
    
    try {
        Write-Host "🔄 Retrieving paginated results..." -ForegroundColor Cyan
        Write-Host "   Initial URI: $InitialUri" -ForegroundColor Gray
        
        $allResults = @()
        $currentUri = $InitialUri
        $pageCount = 0

        do {
            $pageCount++
            Write-Host "   📄 Processing page $pageCount..." -ForegroundColor Gray
            
            $response = Invoke-MgGraphRequest -Method GET -Uri $currentUri
            
            Write-DebugInfo "API Response structure" @{
                "hasValue" = ($null -ne $response.value)
                "valueType" = if ($response.value) { $response.value.GetType().Name } else { "null" }
                "valueCount" = if ($response.value) { $response.value.Count } else { 0 }
                "responseKeys" = $response.PSObject.Properties.Name
            }
            
            if ($response.value) {
                # Log each item being added
                for ($i = 0; $i -lt $response.value.Count; $i++) {
                    $item = $response.value[$i]
                    Write-DebugInfo "Adding item $i to results" @{
                        "isNull" = ($null -eq $item)
                        "type" = if ($item) { $item.GetType().Name } else { "null" }
                        "hasId" = if ($item) { ($null -ne $item.id) } else { $false }
                        "hasSettingInstance" = if ($item) { ($null -ne $item.settingInstance) } else { $false }
                    }
                }
                
                $allResults += $response.value
            }
            
            $currentUri = $response.'@odata.nextLink'
        } while ($currentUri)

        Write-Host "   ✅ Retrieved $($allResults.Count) total results from $pageCount page(s)" -ForegroundColor Green
        
        Write-DebugInfo "Final allResults array" @{
            "count" = $allResults.Count
            "firstItemIsNull" = if ($allResults.Count -gt 0) { ($null -eq $allResults[0]) } else { "no items" }
            "types" = if ($allResults.Count -gt 0) { $allResults | ForEach-Object { if ($_ -eq $null) { "null" } else { $_.GetType().Name } } } else { @() }
        }
        
        return $allResults
    }
    catch {
        Write-Host "❌ Error retrieving paginated results: $_" -ForegroundColor Red
        Write-Host ""
        if ($_.Exception.Response) {
            $statusCode = $_.Exception.Response.StatusCode
            $statusDescription = $_.Exception.Response.StatusDescription
            Write-Host "   Status Code: $statusCode" -ForegroundColor Red
            Write-Host "   Status Description: $statusDescription" -ForegroundColor Red
            try {
                $responseContent = $_.Exception.Response.GetResponseStream()
                if ($responseContent) {
                    $reader = [System.IO.StreamReader]::new($responseContent)
                    $errorDetails = $reader.ReadToEnd()
                    $reader.Close()
                    Write-Host "   Error Details: $errorDetails" -ForegroundColor Red
                }
            } catch {
                Write-Host "   Could not read error details" -ForegroundColor Red
            }
        }
        throw
    }
}

function Get-SettingsCatalogPolicy {
    param (
        [Parameter(Mandatory=$true)]
        [string]$PolicyId
    )
    
    try {
        Write-Host "🔍 Getting settings catalog policy..." -ForegroundColor Cyan
        Write-Host "   Policy ID: $PolicyId" -ForegroundColor Gray
        
        $policyUri = "https://graph.microsoft.com/beta/deviceManagement/configurationPolicies/$PolicyId"
        Write-Host "   Policy Endpoint: $policyUri" -ForegroundColor Gray
        
        $policy = Invoke-MgGraphRequest -Method GET -Uri $policyUri
        Write-Host "   ✅ Policy retrieved successfully" -ForegroundColor Green
        Write-DebugInfo "Policy metadata" @{
            "name" = $policy.name
            "platforms" = $policy.platforms
            "technologies" = $policy.technologies
            "settingCount" = $policy.settingCount
        }
        
        # Retrieve settings with definitions expanded for full fidelity (required for template-based policies)
        $settingsUri = "https://graph.microsoft.com/beta/deviceManagement/configurationPolicies/$PolicyId/settings?$expand=settingDefinitions&$top=1000"
        Write-Host "   Settings Endpoint: $settingsUri" -ForegroundColor Gray
        
        $allSettings = Get-PaginatedResults -InitialUri $settingsUri
        Write-DebugInfo "Raw settings count" $allSettings.Count
        
        # Debug: Check the actual contents of $allSettings
        Write-DebugInfo "allSettings array details" @{
            "count" = $allSettings.Count
            "isArray" = $allSettings -is [Array]
            "type" = $allSettings.GetType().Name
        }
        
        # FIXED: Handle both single object and array cases
        # Force $allSettings to be an array to handle single object case
        if ($allSettings -is [Array]) {
            $settingsArray = $allSettings
        } else {
            # Single object - wrap in array
            $settingsArray = @($allSettings)
        }
        
        Write-DebugInfo "Normalized settings array" @{
            "count" = $settingsArray.Count
            "isArray" = $settingsArray -is [Array]
            "type" = $settingsArray.GetType().Name
        }
        
        # Debug: Check each item in the normalized array
        for ($j = 0; $j -lt $settingsArray.Count; $j++) {
            $item = $settingsArray[$j]
            Write-DebugInfo "settingsArray[$j] details" @{
                "isNull" = ($null -eq $item)
                "type" = if ($item) { $item.GetType().Name } else { "null" }
                "hasId" = if ($item) { [bool]($item.PSObject.Properties['id']) } else { $false }
                "hasSettingInstance" = if ($item) { [bool]($item.PSObject.Properties['settingInstance']) } else { $false }
                "properties" = if ($item) { $item.PSObject.Properties.Name } else { @() }
            }
        }
        
        # FIXED: Use the normalized array instead of the raw response
        $formattedSettings = @()
        for($i = 0; $i -lt $settingsArray.Count; $i++) {
            $currentSetting = $settingsArray[$i]
            
            Write-DebugInfo "Processing setting $i" @{
                "id" = $currentSetting.id
                "hasSettingInstance" = ($null -ne $currentSetting.settingInstance)
                "settingInstanceType" = $currentSetting.settingInstance.'@odata.type'
                "isNull" = ($null -eq $currentSetting)
            }
            
            # Skip null settings
            if ($null -eq $currentSetting) {
                Write-Warning "⚠️ Skipping null setting at index $i"
                continue
            }
            
            # Add ID if not present (some APIs include it, some don't)
            if (-not $currentSetting.id) {
                $currentSetting | Add-Member -NotePropertyName 'id' -NotePropertyValue $i.ToString()
            }
            
            # Add the complete setting object (which already contains settingInstance)
            $formattedSettings += $currentSetting
        }
        
        # Add formatted settings to policy object
        $policy | Add-Member -NotePropertyName 'settings' -NotePropertyValue $formattedSettings -Force

        # If this policy is based on a template, also retrieve the setting templates for completeness
        if ($policy.templateReference -and $policy.templateReference.templateId) {
            try {
                $policyTemplateId = $policy.templateReference.templateId
                $settingTemplatesUri = "https://graph.microsoft.com/beta/deviceManagement/configurationPolicyTemplates('$policyTemplateId')/settingTemplates?$expand=settingDefinitions&$top=1000"
                Write-Host "   Setting Templates Endpoint: $settingTemplatesUri" -ForegroundColor Gray
                $settingTemplates = Get-PaginatedResults -InitialUri $settingTemplatesUri
                $policy | Add-Member -NotePropertyName 'settingTemplates' -NotePropertyValue $settingTemplates -Force
            } catch {
                Write-Host "   ⚠️ Unable to retrieve setting templates for templateId ${policyTemplateId}: $_" -ForegroundColor Yellow
            }
        }
        
        Write-Host "   ✅ Settings processed: $($formattedSettings.Count) setting(s)" -ForegroundColor Green
        Write-Host ""
        
        return $policy
    }
    catch {
        Write-Host "❌ Error retrieving settings catalog policy: $_" -ForegroundColor Red
        Write-Host ""
        if ($_.Exception.Response) {
            $statusCode = $_.Exception.Response.StatusCode
            $statusDescription = $_.Exception.Response.StatusDescription
            Write-Host "   Status Code: $statusCode" -ForegroundColor Red
            Write-Host "   Status Description: $statusDescription" -ForegroundColor Red
            try {
                $responseContent = $_.Exception.Response.GetResponseStream()
                if ($responseContent) {
                    $reader = [System.IO.StreamReader]::new($responseContent)
                    $errorDetails = $reader.ReadToEnd()
                    $reader.Close()
                    Write-Host "   Error Details: $errorDetails" -ForegroundColor Red
                }
            } catch {
                Write-Host "   Could not read error details" -ForegroundColor Red
            }
        }
        throw
    }
}

function Export-PolicyToHcl {
    param (
        [Parameter(Mandatory=$true)]
        $Policy,
        [Parameter(Mandatory=$true)]
        [string]$PolicyId
    )
    
    try {
        $outputDir = Join-Path -Path (Get-Location) -ChildPath "output"
        if (-not (Test-Path -Path $outputDir)) {
            New-Item -Path $outputDir -ItemType Directory | Out-Null
            Write-Host "📁 Created output directory: $outputDir" -ForegroundColor Gray
        }
        
        $timestamp = Get-Date -Format "yyyyMMdd_HHmmss"
        $policyName = $Policy.name -replace '[\\\/:\*\?\"\<\>\|]', '_'
        if (-not $policyName) { $policyName = $PolicyId }
        
        $fileName = "SettingsCatalogPolicy_${policyName}_${timestamp}.tf"
        $filePath = Join-Path -Path $outputDir -ChildPath $fileName
        
        $hclContent = Convert-PolicyToHcl -Policy $Policy
        $hclContent | Out-File -LiteralPath $filePath -Encoding UTF8
        
        Write-Host "💾 Exported policy to HCL: $filePath" -ForegroundColor Green
        return $filePath
    }
    catch {
        Write-Host "❌ Error exporting policy to HCL: $_" -ForegroundColor Red
        return $null
    }
}

function Convert-PolicyToHcl {
    param (
        [Parameter(Mandatory=$true)]
        $Policy
    )
    
    # Create a valid Terraform resource name
    $resourceName = ($Policy.name -replace '[^a-zA-Z0-9_]', '_' -replace '__+', '_').ToLower().Trim('_')
    if ($resourceName -match '^[0-9]') {
        $resourceName = "policy_$resourceName"
    }
    
    # Escape description for HCL: convert newlines to \n, escape backslashes and quotes
    $description = if ($Policy.description) {
        ($Policy.description -replace '\\','\\\\' -replace '"','\"' -replace "(\r\n|\n|\r)", "\n")
    } else { "" }
    
    # Handle technologies as array
    $technologies = if ($Policy.technologies) {
        ($Policy.technologies -split ',' | ForEach-Object { '"' + $_.Trim() + '"' }) -join ', '
    } else {
        '""'
    }
    
    # Handle roleScopeTagIds as array
    $roleScopeTagIds = if ($Policy.roleScopeTagIds) {
        ($Policy.roleScopeTagIds | ForEach-Object { '"' + $_ + '"' }) -join ', '
    } else {
        '"0"'
    }
    
    $hcl = @"
resource "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy" "$resourceName" {
  name               = "$($Policy.name)"
  description        = "$description"
  platforms          = "$($Policy.platforms)"
  technologies       = [$technologies]
  role_scope_tag_ids = [$roleScopeTagIds]

  template_reference = {
    template_id = "$($Policy.templateReference.templateId)"
  }

  configuration_policy = {
    settings = [
"@
    
    foreach ($setting in $Policy.settings) {
        $hcl += "`n      {"
        $hcl += Convert-SettingToHcl -Setting $setting -IndentLevel 8
        $hcl += "`n      },"
    }
    
    # Remove the last comma and add closing brackets
    $hcl = $hcl.TrimEnd(',')
    $hcl += @"

    ]
  }
}
"@
    
    return $hcl
}

function Convert-SettingToHcl {
    param (
        [Parameter(Mandatory=$true)]
        $Setting,
        [Parameter(Mandatory=$true)]
        [int]$IndentLevel
    )
    
    $indent = " " * $IndentLevel
    $hcl = ""
    
    if ($Setting.settingInstance) {
        $hcl += "`n$indent" + "setting_instance = {"
        $hcl += Convert-SettingInstanceToHcl -SettingInstance $Setting.settingInstance -IndentLevel ($IndentLevel + 2)
        $hcl += "`n$indent" + "}"
    }
    
    if ($Setting.id) {
        $hcl += "`n$indent" + "id = `"$($Setting.id)`""
    }
    
    return $hcl
}

function Convert-SettingInstanceToHcl {
    param (
        [Parameter(Mandatory=$true)]
        $SettingInstance,
        [Parameter(Mandatory=$true)]
        [int]$IndentLevel
    )
    
    $indent = " " * $IndentLevel
    $hcl = ""
    
    Write-DebugInfo "Converting SettingInstance" @{
        "odata.type" = $SettingInstance.'@odata.type'
        "settingDefinitionId" = $SettingInstance.settingDefinitionId
        "hasSimpleValue" = ($null -ne $SettingInstance.simpleSettingValue)
        "hasChoiceValue" = ($null -ne $SettingInstance.choiceSettingValue)
        "hasGroupCollectionValue" = ($null -ne $SettingInstance.groupSettingCollectionValue)
        "hasSimpleCollectionValue" = ($null -ne $SettingInstance.simpleSettingCollectionValue)
        "hasGroupValue" = ($null -ne $SettingInstance.groupSettingValue)
    }
    
    # Handle @odata.type
    if ($SettingInstance.'@odata.type') {
        $hcl += "`n$indent" + "odata_type = `"$($SettingInstance.'@odata.type')`""
    }
    
    # Handle settingDefinitionId
    if ($SettingInstance.settingDefinitionId) {
        $hcl += "`n$indent" + "setting_definition_id = `"$($SettingInstance.settingDefinitionId)`""
    }
    
    # Handle settingInstanceTemplateReference
    if ($SettingInstance.settingInstanceTemplateReference) {
        $ref = $SettingInstance.settingInstanceTemplateReference
        $hcl += "`n$indent" + "setting_instance_template_reference = {"
        if ($ref.settingInstanceTemplateId) {
            $hcl += "`n$indent  " + "setting_instance_template_id = `"$($ref.settingInstanceTemplateId)`""
        }
        $hcl += "`n$indent" + "}"
    } else {
        $hcl += "`n$indent" + "setting_instance_template_reference = null"
    }
    
    # Handle different setting value types - order matters for proper detection
    if ($SettingInstance.simpleSettingValue) {
        Write-DebugInfo "Processing simpleSettingValue" $SettingInstance.simpleSettingValue
        $hcl += "`n$indent" + "simple_setting_value = {"
        $hcl += Convert-SimpleSettingValueToHcl -SimpleSettingValue $SettingInstance.simpleSettingValue -IndentLevel ($IndentLevel + 2)
        $hcl += "`n$indent" + "}"
    }
    
    if ($SettingInstance.choiceSettingValue) {
        Write-DebugInfo "Processing choiceSettingValue" @{
            "value" = $SettingInstance.choiceSettingValue.value
            "childrenCount" = if ($SettingInstance.choiceSettingValue.children) { $SettingInstance.choiceSettingValue.children.Count } else { 0 }
        }
        $hcl += "`n$indent" + "choice_setting_value = {"
        $hcl += Convert-ChoiceSettingValueToHcl -ChoiceSettingValue $SettingInstance.choiceSettingValue -IndentLevel ($IndentLevel + 2)
        $hcl += "`n$indent" + "}"
    }
    
    if ($SettingInstance.groupSettingCollectionValue) {
        Write-DebugInfo "Processing groupSettingCollectionValue" @{
            "itemCount" = $SettingInstance.groupSettingCollectionValue.Count
        }
        $hcl += "`n$indent" + "group_setting_collection_value = ["
        foreach ($groupValue in $SettingInstance.groupSettingCollectionValue) {
            $hcl += "`n$indent" + "  {"
            $hcl += Convert-GroupSettingValueToHcl -GroupSettingValue $groupValue -IndentLevel ($IndentLevel + 4)
            $hcl += "`n$indent" + "  },"
        }
        $hcl = $hcl.TrimEnd(',')
        $hcl += "`n$indent" + "]"
    }
    
    # Handle simple setting collection value (for arrays like URLs)
    if ($SettingInstance.simpleSettingCollectionValue) {
        Write-DebugInfo "Processing simpleSettingCollectionValue" @{
            "itemCount" = $SettingInstance.simpleSettingCollectionValue.Count
        }
        $hcl += "`n$indent" + "simple_setting_collection_value = ["
        foreach ($simpleValue in $SettingInstance.simpleSettingCollectionValue) {
            $hcl += "`n$indent" + "  {"
            $hcl += Convert-SimpleSettingValueToHcl -SimpleSettingValue $simpleValue -IndentLevel ($IndentLevel + 4)
            $hcl += "`n$indent" + "  },"
        }
        $hcl = $hcl.TrimEnd(',')
        $hcl += "`n$indent" + "]"
    }
    
    # Handle group setting value (single group, not collection)
    if ($SettingInstance.groupSettingValue) {
        Write-DebugInfo "Processing groupSettingValue" $SettingInstance.groupSettingValue
        $hcl += "`n$indent" + "group_setting_value = {"
        $hcl += Convert-GroupSettingValueToHcl -GroupSettingValue $SettingInstance.groupSettingValue -IndentLevel ($IndentLevel + 2)
        $hcl += "`n$indent" + "}"
    }
    
    # Check for any unsupported patterns
    $supportedProperties = @('@odata.type', 'settingDefinitionId', 'settingInstanceTemplateReference', 
                            'simpleSettingValue', 'choiceSettingValue', 'groupSettingCollectionValue', 
                            'simpleSettingCollectionValue', 'groupSettingValue','IsReadOnly','IsFixedSize','IsSynchronized','Keys','Values','SyncRoot')
    
    # Ignore common .NET collection members that can appear on adapted objects
    $ignoredProperties = @('Count')

    foreach ($property in $SettingInstance.PSObject.Properties.Name) {
        if ($property -notin $supportedProperties -and $property -notin $ignoredProperties) {
            Write-Warning "🚨 Unsupported property detected: $property in SettingInstance"
            Write-DebugInfo "Unsupported property" @{
                "property" = $property
                "value" = $SettingInstance.$property
            }
        }
    }
    
    return $hcl
}

function Convert-SimpleSettingValueToHcl {
    param (
        [Parameter(Mandatory=$true)]
        $SimpleSettingValue,
        [Parameter(Mandatory=$true)]
        [int]$IndentLevel
    )
    
    $indent = " " * $IndentLevel
    $hcl = ""
    
    if ($SimpleSettingValue.'@odata.type') {
        $hcl += "`n$indent" + "odata_type = `"$($SimpleSettingValue.'@odata.type')`""
    }
    
    if ($SimpleSettingValue.settingValueTemplateReference) {
        $ref = $SimpleSettingValue.settingValueTemplateReference
        $hcl += "`n$indent" + "setting_value_template_reference = {"
        if ($ref.settingValueTemplateId) {
            $hcl += "`n$indent  " + "setting_value_template_id = `"$($ref.settingValueTemplateId)`""
        }
        if ($null -ne $ref.useTemplateDefault) {
            $hcl += "`n$indent  " + "use_template_default = $($ref.useTemplateDefault.ToString().ToLower())"
        }
        $hcl += "`n$indent" + "}"
    } else {
        $hcl += "`n$indent" + "setting_value_template_reference = null"
    }
    
    # Handle valueState property (for secret settings)
    if ($SimpleSettingValue.valueState) {
        $hcl += "`n$indent" + "value_state = `"$($SimpleSettingValue.valueState)`""
    }
    
    if ($SimpleSettingValue.value -ne $null) {
        if ($SimpleSettingValue.'@odata.type' -match 'Integer') {
            $hcl += "`n$indent" + "value = $($SimpleSettingValue.value)"
        } elseif ($SimpleSettingValue.'@odata.type' -match 'Boolean') {
            $hcl += "`n$indent" + "value = $($SimpleSettingValue.value.ToString().ToLower())"
        } else {
            # Escape string with UNC-aware handling
            $escapedValue = Escape-HclString -InputString $SimpleSettingValue.value
            $hcl += "`n$indent" + "value = `"$escapedValue`""
        }
    }
    
    return $hcl
}

function Convert-ChoiceSettingValueToHcl {
    param (
        [Parameter(Mandatory=$true)]
        $ChoiceSettingValue,
        [Parameter(Mandatory=$true)]
        [int]$IndentLevel
    )
    
    $indent = " " * $IndentLevel
    $hcl = ""
    
    if ($ChoiceSettingValue.settingValueTemplateReference) {
        $ref = $ChoiceSettingValue.settingValueTemplateReference
        $hcl += "`n$indent" + "setting_value_template_reference = {"
        if ($ref.settingValueTemplateId) {
            $hcl += "`n$indent  " + "setting_value_template_id = `"$($ref.settingValueTemplateId)`""
        }
        if ($null -ne $ref.useTemplateDefault) {
            $hcl += "`n$indent  " + "use_template_default = $($ref.useTemplateDefault.ToString().ToLower())"
        }
        $hcl += "`n$indent" + "}"
    } else {
        $hcl += "`n$indent" + "setting_value_template_reference = null"
    }
    
    if ($ChoiceSettingValue.children) {
        $hcl += "`n$indent" + "children = ["
        foreach ($child in $ChoiceSettingValue.children) {
            $hcl += "`n$indent" + "  {"
            $hcl += Convert-SettingInstanceToHcl -SettingInstance $child -IndentLevel ($IndentLevel + 4)
            $hcl += "`n$indent" + "  },"
        }
        $hcl = $hcl.TrimEnd(',')
        $hcl += "`n$indent" + "]"
    } else {
        $hcl += "`n$indent" + "children = []"
    }
    
    if ($ChoiceSettingValue.value) {
        $escapedChoice = Escape-HclString -InputString $ChoiceSettingValue.value
        $hcl += "`n$indent" + "value = `"$escapedChoice`""
    }
    
    return $hcl
}

function Convert-GroupSettingValueToHcl {
    param (
        [Parameter(Mandatory=$true)]
        $GroupSettingValue,
        [Parameter(Mandatory=$true)]
        [int]$IndentLevel
    )
    
    $indent = " " * $IndentLevel
    $hcl = ""
    
    if ($GroupSettingValue.settingValueTemplateReference) {
        $ref = $GroupSettingValue.settingValueTemplateReference
        $hcl += "`n$indent" + "setting_value_template_reference = {"
        if ($ref.settingValueTemplateId) {
            $hcl += "`n$indent  " + "setting_value_template_id = `"$($ref.settingValueTemplateId)`""
        }
        if ($null -ne $ref.useTemplateDefault) {
            $hcl += "`n$indent  " + "use_template_default = $($ref.useTemplateDefault.ToString().ToLower())"
        }
        $hcl += "`n$indent" + "}"
    } else {
        $hcl += "`n$indent" + "setting_value_template_reference = null"
    }
    
    if ($GroupSettingValue.children) {
        $hcl += "`n$indent" + "children = ["
        foreach ($child in $GroupSettingValue.children) {
            $hcl += "`n$indent" + "  {"
            $hcl += Convert-SettingInstanceToHcl -SettingInstance $child -IndentLevel ($IndentLevel + 4)
            $hcl += "`n$indent" + "  },"
        }
        $hcl = $hcl.TrimEnd(',')
        $hcl += "`n$indent" + "]"
    } else {
        $hcl += "`n$indent" + "children = []"
    }
    
    return $hcl
}

function Export-PolicyToJson {
    param (
        [Parameter(Mandatory=$true)]
        $Policy,
        [Parameter(Mandatory=$true)]
        [string]$PolicyId
    )
    
    try {
        $outputDir = Join-Path -Path (Get-Location) -ChildPath "output"
        if (-not (Test-Path -Path $outputDir)) {
            New-Item -Path $outputDir -ItemType Directory | Out-Null
            Write-Host "📁 Created output directory: $outputDir" -ForegroundColor Gray
        }
        
        $timestamp = Get-Date -Format "yyyyMMdd_HHmmss"
        $policyName = $Policy.name -replace '[\\\/:\*\?\"\<\>\|]', '_'
        if (-not $policyName) { $policyName = $PolicyId }
        
        $fileName = "SettingsCatalogPolicy_${policyName}_${timestamp}.json"
        $filePath = Join-Path -Path $outputDir -ChildPath $fileName
        
        $jsonFormatted = $Policy | ConvertTo-Json -Depth 100
        $jsonFormatted | Out-File -LiteralPath $filePath -Encoding UTF8
        
        Write-Host "💾 Exported policy to: $filePath" -ForegroundColor Green
        return $filePath
    }
    catch {
        Write-Host "❌ Error exporting policy to JSON: $_" -ForegroundColor Red
        return $null
    }
}

function Show-PolicyDetails {
    param (
        [Parameter(Mandatory=$true)]
        $Policy
    )
    
    Write-Host "📋 Settings Catalog Policy Details:" -ForegroundColor Cyan
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    
    # Basic policy information
    foreach ($field in @('id', 'name', 'description', 'platforms', 'technologies', 'templateReference', 'settingCount', 'createdDateTime', 'lastModifiedDateTime', 'creationSource', 'priorityMetaData')) {
        if ($Policy.PSObject.Properties[$field]) {
            if ($field -eq 'templateReference' -and $Policy.$field) {
                Write-Host "   • templateReference:" -ForegroundColor Green
                Write-Host "     · templateId: $($Policy.$field.templateId)" -ForegroundColor Yellow
                Write-Host "     · templateFamily: $($Policy.$field.templateFamily)" -ForegroundColor Yellow
                Write-Host "     · templateDisplayName: $($Policy.$field.templateDisplayName)" -ForegroundColor Yellow
                Write-Host "     · templateDisplayVersion: $($Policy.$field.templateDisplayVersion)" -ForegroundColor Yellow
            } elseif ($field -eq 'platforms' -and $Policy.$field) {
                Write-Host "   • platforms: $($Policy.$field -join ', ')" -ForegroundColor Green
            } elseif ($field -eq 'technologies' -and $Policy.$field) {
                Write-Host "   • technologies: $($Policy.$field -join ', ')" -ForegroundColor Green
            } elseif ($field -eq 'priorityMetaData' -and $Policy.$field) {
                Write-Host "   • priorityMetaData:" -ForegroundColor Green
                Write-Host "     · priority: $($Policy.$field.priority)" -ForegroundColor Yellow
            } else {
                Write-Host ("   • {0}: {1}" -f $field, $Policy.$field) -ForegroundColor Green
            }
        }
    }
    
    # Assignment information
    if ($Policy.assignments) {
        Write-Host "   • assignments:" -ForegroundColor Green
        foreach ($assignment in $Policy.assignments) {
            Write-Host "     · id: $($assignment.id)" -ForegroundColor Yellow
            Write-Host "     · target: $($assignment.target)" -ForegroundColor Yellow
        }
    }
    
    # Settings information - ENHANCED to show nested structure
    if ($Policy.settings -and $Policy.settings.Count -gt 0) {
        Write-Host "   • settings ($($Policy.settings.Count) setting(s)):" -ForegroundColor Green
        foreach ($setting in $Policy.settings) {
            Write-Host "     · Setting ID: $($setting.id)" -ForegroundColor Yellow
            if ($setting.settingInstance) {
                Show-SettingInstanceDetails -SettingInstance $setting.settingInstance -IndentLevel 3
            }
        }
    }
    
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    Write-Host ""
}

function Show-SettingInstanceDetails {
    param (
        [Parameter(Mandatory=$true)]
        $SettingInstance,
        [Parameter(Mandatory=$true)]
        [int]$IndentLevel
    )
    
    $indent = "       " + ("  " * $IndentLevel)
    
    if ($SettingInstance) {
        Write-Host "$indent- @odata.type: $($SettingInstance.'@odata.type')" -ForegroundColor Magenta
        Write-Host "$indent- settingDefinitionId: $($SettingInstance.settingDefinitionId)" -ForegroundColor Magenta
        
        if ($SettingInstance.simpleSettingValue) {
            Write-Host "$indent- simpleSettingValue:" -ForegroundColor Magenta
            Write-Host "$indent  · value: $($SettingInstance.simpleSettingValue.value)" -ForegroundColor White
            Write-Host "$indent  · @odata.type: $($SettingInstance.simpleSettingValue.'@odata.type')" -ForegroundColor White
        }
        
        if ($SettingInstance.choiceSettingValue) {
            Write-Host "$indent- choiceSettingValue:" -ForegroundColor Magenta
            Write-Host "$indent  · value: $($SettingInstance.choiceSettingValue.value)" -ForegroundColor White
            
            if ($SettingInstance.choiceSettingValue.children -and $SettingInstance.choiceSettingValue.children.Count -gt 0) {
                Write-Host "$indent  · children ($($SettingInstance.choiceSettingValue.children.Count)):" -ForegroundColor White
                foreach ($child in $SettingInstance.choiceSettingValue.children) {
                    Show-SettingInstanceDetails -SettingInstance $child -IndentLevel ($IndentLevel + 2)
                }
            }
        }
        
        if ($SettingInstance.groupSettingCollectionValue) {
            Write-Host "$indent- groupSettingCollectionValue ($($SettingInstance.groupSettingCollectionValue.Count) items):" -ForegroundColor Magenta
            foreach ($groupValue in $SettingInstance.groupSettingCollectionValue) {
                if ($groupValue.children -and $groupValue.children.Count -gt 0) {
                    Write-Host "$indent  · children ($($groupValue.children.Count)):" -ForegroundColor White
                    foreach ($child in $groupValue.children) {
                        Show-SettingInstanceDetails -SettingInstance $child -IndentLevel ($IndentLevel + 2)
                    }
                }
            }
        }
        
        if ($SettingInstance.simpleSettingCollectionValue) {
            Write-Host "$indent- simpleSettingCollectionValue ($($SettingInstance.simpleSettingCollectionValue.Count) items):" -ForegroundColor Magenta
            foreach ($simpleValue in $SettingInstance.simpleSettingCollectionValue) {
                Write-Host "$indent  · value: $($simpleValue.value)" -ForegroundColor White
                Write-Host "$indent  · @odata.type: $($simpleValue.'@odata.type')" -ForegroundColor White
            }
        }
        
        # Handle single group setting value (not collection)
        if ($SettingInstance.groupSettingValue) {
            Write-Host "$indent- groupSettingValue:" -ForegroundColor Magenta
            if ($SettingInstance.groupSettingValue.children -and $SettingInstance.groupSettingValue.children.Count -gt 0) {
                Write-Host "$indent  · children ($($SettingInstance.groupSettingValue.children.Count)):" -ForegroundColor White
                foreach ($child in $SettingInstance.groupSettingValue.children) {
                    Show-SettingInstanceDetails -SettingInstance $child -IndentLevel ($IndentLevel + 2)
                }
            }
        }
    }
}

try {
    Write-Host "🔐 Connecting to Microsoft Graph..." -ForegroundColor Cyan
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    
    $secureClientSecret = ConvertTo-SecureString -String $ClientSecret -AsPlainText -Force
    $clientSecretCredential = New-Object -TypeName System.Management.Automation.PSCredential -ArgumentList $ClientId, $secureClientSecret
    
    Connect-MgGraph -ClientSecretCredential $clientSecretCredential -TenantId $TenantId
    Write-Host "✅ Connected to Microsoft Graph successfully" -ForegroundColor Green
    Write-Host ""
    
    # Validate that at least one output option is selected
    if (-not $ExportToJson -and -not $ExportToHcl) {
        Write-Host "ℹ️  No export format specified. Policy details will be displayed in console." -ForegroundColor Yellow
        Write-Host "   Use -ExportToJson `$true or -ExportToHcl `$true to export to files." -ForegroundColor Yellow
        Write-Host "   Use -EnableDebug `$true to enable detailed debugging output." -ForegroundColor Yellow
        Write-Host ""
    }
    
    if ($EnableDebug) {
        Write-Host "🐛 Debug mode enabled - verbose output will be shown" -ForegroundColor Yellow
        Write-Host ""
    }
    
    if ($ExportAll) {
        if ($SettingsCatalogItemId) {
            Write-Host "ℹ️  Both -ExportAll and -SettingsCatalogItemId were provided. Proceeding with -ExportAll and ignoring -SettingsCatalogItemId." -ForegroundColor Yellow
        }

        $listUri = "https://graph.microsoft.com/beta/deviceManagement/configurationPolicies"
        Write-Host "🔎 Listing all settings catalog policies..." -ForegroundColor Cyan
        $allPolicies = Get-PaginatedResults -InitialUri $listUri
        $total = if ($allPolicies) { $allPolicies.Count } else { 0 }
        Write-Host "   ✅ Found $total policy item(s)" -ForegroundColor Green

        if ($total -eq 0) {
            Write-Host "📊 No settings catalog policies found." -ForegroundColor Yellow
        } else {
            $index = 0
            foreach ($p in $allPolicies) {
                $index++
                Write-Host "➡️  [$index/$total] Processing policy: $($p.name) ($($p.id))" -ForegroundColor Cyan

                $policyData = Get-SettingsCatalogPolicy -PolicyId $p.id
                if ($null -eq $policyData) { continue }

                $exportPaths = @()

                if ($ExportToJson) {
                    $jsonPath = Export-PolicyToJson -Policy $policyData -PolicyId $p.id
                    if ($jsonPath) { $exportPaths += $jsonPath }
                }

                if ($ExportToHcl) {
                    $hclPath = Export-PolicyToHcl -Policy $policyData -PolicyId $p.id
                    if ($hclPath) { $exportPaths += $hclPath }
                }

                if (-not $ExportToJson -and -not $ExportToHcl) {
                    Show-PolicyDetails -Policy $policyData
                }
            }

            Write-Host "🎉 Completed processing all policies." -ForegroundColor Green
        }
    }
    else {
        $policyData = Get-SettingsCatalogPolicy -PolicyId $SettingsCatalogItemId
        
        if ($null -ne $policyData) {
            $exportPaths = @()
            
            if ($ExportToJson) {
                $jsonPath = Export-PolicyToJson -Policy $policyData -PolicyId $SettingsCatalogItemId
                if ($jsonPath) {
                    $exportPaths += $jsonPath
                }
            }
            
            if ($ExportToHcl) {
                $hclPath = Export-PolicyToHcl -Policy $policyData -PolicyId $SettingsCatalogItemId
                if ($hclPath) {
                    $exportPaths += $hclPath
                }
            }
            
            if ($exportPaths.Count -gt 0) {
                Write-Host "" 
            }
            
            # Always show policy details unless only exporting
            if (-not $ExportToJson -and -not $ExportToHcl) {
                # No export requested, show full details
                Show-PolicyDetails -Policy $policyData
            } else {
                # Export requested, show summary and full details
                Write-Host "📋 Policy Summary:" -ForegroundColor Cyan
                Write-Host "   • Name: $($policyData.name)" -ForegroundColor Green
                Write-Host "   • ID: $($policyData.id)" -ForegroundColor Green
                Write-Host "   • Platform: $($policyData.platforms)" -ForegroundColor Green
                Write-Host "   • Settings Count: $($policyData.settings.Count)" -ForegroundColor Green
                Write-Host "   • Technologies: $($policyData.technologies)" -ForegroundColor Green
                Write-Host ""
                
                # Also show full details
                Show-PolicyDetails -Policy $policyData
            }
            
            Write-Host "🎉 Operation completed successfully!" -ForegroundColor Green
        } else {
            Write-Host "📊 No data found for the specified settings catalog policy ID" -ForegroundColor Yellow
        }
    }
}
catch {
    Write-Host "❌ Script execution failed: $_" -ForegroundColor Red
    exit 1
}
finally {
    Write-Host "🔌 Disconnecting from Microsoft Graph..." -ForegroundColor Cyan
    try {
        Disconnect-MgGraph 2>$null
        Write-Host "✅ Disconnected from Microsoft Graph" -ForegroundColor Green
    } catch {}
}