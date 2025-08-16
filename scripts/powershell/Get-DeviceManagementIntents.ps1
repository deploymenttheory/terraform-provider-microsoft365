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
    HelpMessage="Specific Device Management Intent ID (if not provided, will list all intents)")]
    [string]$IntentId,
    
    [Parameter(Mandatory=$true,
    HelpMessage="Export results to JSON file")]
    [bool]$ExportToJson
)

Import-Module Microsoft.Graph.Authentication

function Get-DeviceManagementIntents {
    param (
        [Parameter(Mandatory=$false)]
        [string]$DeviceManagementIntentId
    )
    try {
        if ($DeviceManagementIntentId) {
            $uri = "https://graph.microsoft.com/beta/deviceManagement/intents/$DeviceManagementIntentId"
            Write-Host "🔍 Getting specific device management intent..." -ForegroundColor Cyan
            Write-Host "   Intent ID: $DeviceManagementIntentId" -ForegroundColor Gray
        } else {
            $uri = "https://graph.microsoft.com/beta/deviceManagement/intents"
            Write-Host "🔍 Getting all device management intents..." -ForegroundColor Cyan
        }
        Write-Host "   Endpoint: $uri" -ForegroundColor Gray
        Write-Host ""
        $response = Invoke-MgGraphRequest -Method GET -Uri $uri
        return $response
    }
    catch {
        Write-Host "❌ Error getting device management intents: $_" -ForegroundColor Red
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

function Export-IntentsToJson {
    param (
        [Parameter(Mandatory=$true)]
        $Intents,
        [Parameter(Mandatory=$false)]
        [string]$DeviceManagementIntentId
    )
    try {
        $outputDir = Join-Path -Path (Get-Location) -ChildPath "output"
        if (-not (Test-Path -Path $outputDir)) {
            New-Item -Path $outputDir -ItemType Directory | Out-Null
            Write-Host "📁 Created output directory: $outputDir" -ForegroundColor Gray
        }
        $timestamp = Get-Date -Format "yyyyMMdd_HHmmss"
        if ($DeviceManagementIntentId) {
            $intentName = $Intents.displayName -replace '[\\\/:\*\?\"\<\>\|]', '_'
            if (-not $intentName) { $intentName = $DeviceManagementIntentId }
            $fileName = "DeviceManagementIntent_${intentName}_${timestamp}.json"
            $filePath = Join-Path -Path $outputDir -ChildPath $fileName
            $Intents | ConvertTo-Json -Depth 10 | Out-File -FilePath $filePath -Encoding UTF8
            Write-Host "💾 Exported intent to: $filePath" -ForegroundColor Green
        } else {
            $fileName = "DeviceManagementIntents_${timestamp}.json"
            $filePath = Join-Path -Path $outputDir -ChildPath $fileName
            $Intents | ConvertTo-Json -Depth 10 | Out-File -FilePath $filePath -Encoding UTF8
            Write-Host "💾 Exported intents to: $filePath" -ForegroundColor Green
        }
        return $filePath
    } catch {
        Write-Host "❌ Error exporting intents to JSON: $_" -ForegroundColor Red
        return $null
    }
}

function Show-IntentDetails {
    param (
        [Parameter(Mandatory=$true)]
        $Intent
    )
    Write-Host "📋 Device Management Intent Details:" -ForegroundColor Cyan
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    
    # Top-level fields
    foreach ($field in @(
        '@odata.type','id','displayName','description','isAssigned','isMigratingToConfigurationPolicy','lastModifiedDateTime','templateId')) {
        if ($Intent.PSObject.Properties[$field]) {
            Write-Host ("   • {0}: {1}" -f $field, $Intent.$field) -ForegroundColor Green
        }
    }

    # Arrays
    if ($Intent.roleScopeTagIds) {
        Write-Host "   • roleScopeTagIds:" -ForegroundColor Green
        foreach ($item in $Intent.roleScopeTagIds) {
            Write-Host "     · $item" -ForegroundColor Yellow
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
    $intents = Get-DeviceManagementIntents -DeviceManagementIntentId $IntentId
    if ($ExportToJson) {
        Export-IntentsToJson -Intents $intents -DeviceManagementIntentId $IntentId
    }
    if ($IntentId) {
        Show-IntentDetails -Intent $intents
    } else {
        if ($intents.value -and $intents.value.Count -gt 0) {
            Write-Host "📊 Found $($intents.value.Count) device management intent(s)" -ForegroundColor Green
            Write-Host ""
            for ($i = 0; $i -lt $intents.value.Count; $i++) {
                Write-Host "Intent $($i + 1):" -ForegroundColor Magenta
                Show-IntentDetails -Intent $intents.value[$i]
            }
        } elseif ($intents -and -not $intents.value) {
            Write-Host "📊 Found 1 device management intent" -ForegroundColor Green
            Write-Host ""
            Show-IntentDetails -Intent $intents
        } else {
            Write-Host "📊 No device management intents found" -ForegroundColor Yellow
        }
    }
    Write-Host "🎉 Operation completed successfully!" -ForegroundColor Green
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
