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
    HelpMessage="Specify the ID of the settings catalog policy to retrieve")]
    [ValidateNotNullOrEmpty()]
    [string]$SettingsCatalogItemId,
    
    [Parameter(Mandatory=$true,
    HelpMessage="Export results to JSON file")]
    [bool]$ExportToJson
)

Import-Module Microsoft.Graph.Authentication

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
            
            if ($response.value) {
                $allResults += $response.value
            }
            
            $currentUri = $response.'@odata.nextLink'
        } while ($currentUri)

        Write-Host "   ✅ Retrieved $($allResults.Count) total results from $pageCount page(s)" -ForegroundColor Green
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
        
        $settingsUri = "https://graph.microsoft.com/beta/deviceManagement/configurationPolicies/$PolicyId/settings"
        Write-Host "   Settings Endpoint: $settingsUri" -ForegroundColor Gray
        
        $allSettings = Get-PaginatedResults -InitialUri $settingsUri
        
        # Format settings with sequential IDs
        $formattedSettings = @()
        for($i = 0; $i -lt $allSettings.Count; $i++) {
            $formattedSettings += @{
                id = $i.ToString()
                settingInstance = $allSettings[$i].settingInstance
            }
        }
        
        # Add formatted settings to policy object
        $policy | Add-Member -NotePropertyName 'settings' -NotePropertyValue $formattedSettings
        
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
        $jsonFormatted | Out-File -FilePath $filePath -Encoding UTF8
        
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
    
    # Settings information
    if ($Policy.settings -and $Policy.settings.Count -gt 0) {
        Write-Host "   • settings ($($Policy.settings.Count) setting(s)):" -ForegroundColor Green
        foreach ($setting in $Policy.settings) {
            Write-Host "     · Setting ID: $($setting.id)" -ForegroundColor Yellow
            if ($setting.settingInstance) {
                $instance = $setting.settingInstance
                Write-Host "       - settingDefinitionId: $($instance.settingDefinitionId)" -ForegroundColor Magenta
                Write-Host "       - settingInstanceTemplateId: $($instance.settingInstanceTemplateId)" -ForegroundColor Magenta
                
                if ($instance.simpleSettingValue) {
                    Write-Host "       - simpleSettingValue:" -ForegroundColor Magenta
                    Write-Host "         · value: $($instance.simpleSettingValue.value)" -ForegroundColor White
                    Write-Host "         · valueType: $($instance.simpleSettingValue.valueType)" -ForegroundColor White
                }
                
                if ($instance.choiceSettingValue) {
                    Write-Host "       - choiceSettingValue:" -ForegroundColor Magenta
                    Write-Host "         · value: $($instance.choiceSettingValue.value)" -ForegroundColor White
                    Write-Host "         · settingValueTemplateId: $($instance.choiceSettingValue.settingValueTemplateId)" -ForegroundColor White
                    
                    if ($instance.choiceSettingValue.children) {
                        Write-Host "         · children:" -ForegroundColor White
                        foreach ($child in $instance.choiceSettingValue.children) {
                            Write-Host "           - settingDefinitionId: $($child.settingDefinitionId)" -ForegroundColor Gray
                            Write-Host "           - settingInstanceTemplateId: $($child.settingInstanceTemplateId)" -ForegroundColor Gray
                        }
                    }
                }
                
                if ($instance.groupSettingCollectionValue) {
                    Write-Host "       - groupSettingCollectionValue:" -ForegroundColor Magenta
                    Write-Host "         · children count: $($instance.groupSettingCollectionValue.children.Count)" -ForegroundColor White
                }
            }
        }
    }
    
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    Write-Host ""
}

try {
    Write-Host "🔐 Connecting to Microsoft Graph..." -ForegroundColor Cyan
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    
    $secureClientSecret = ConvertTo-SecureString -String $ClientSecret -AsPlainText -Force
    $clientSecretCredential = New-Object -TypeName System.Management.Automation.PSCredential -ArgumentList $ClientId, $secureClientSecret
    
    Connect-MgGraph -ClientSecretCredential $clientSecretCredential -TenantId $TenantId
    Write-Host "✅ Connected to Microsoft Graph successfully" -ForegroundColor Green
    Write-Host ""
    
    $policyData = Get-SettingsCatalogPolicy -PolicyId $SettingsCatalogItemId
    
    if ($null -ne $policyData) {
        if ($ExportToJson) {
            $exportPath = Export-PolicyToJson -Policy $policyData -PolicyId $SettingsCatalogItemId
            Write-Host ""
        }
        
        Show-PolicyDetails -Policy $policyData
        
        Write-Host "🎉 Operation completed successfully!" -ForegroundColor Green
    } else {
        Write-Host "📊 No data found for the specified settings catalog policy ID" -ForegroundColor Yellow
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