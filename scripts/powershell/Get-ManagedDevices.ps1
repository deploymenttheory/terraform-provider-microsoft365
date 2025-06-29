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
    HelpMessage="Specific Managed Device ID (if not provided, will list all devices)")]
    [string]$DeviceId,
    
    [Parameter(Mandatory=$true,
    HelpMessage="Export results to JSON file")]
    [bool]$ExportToJson
)

Import-Module Microsoft.Graph.Authentication

function Get-WindowsManagedDevices {
    param (
        [Parameter(Mandatory=$false)]
        [string]$SpecificDeviceId
    )
    try {
        if ($SpecificDeviceId) {
            $uri = "https://graph.microsoft.com/beta/deviceManagement/managedDevices/$SpecificDeviceId"
            Write-Host "🔍 Getting specific Windows managed device..." -ForegroundColor Cyan
            Write-Host "   Device ID: $SpecificDeviceId" -ForegroundColor Gray
        } else {
            $uri = "https://graph.microsoft.com/beta/deviceManagement/managedDevices"
            Write-Host "🔍 Getting all Windows managed devices..." -ForegroundColor Cyan
        }
        Write-Host "   Endpoint: $uri" -ForegroundColor Gray
        Write-Host ""
        $response = Invoke-MgGraphRequest -Method GET -Uri $uri
        return $response
    }
    catch {
        Write-Host "❌ Error getting Windows managed devices: $_" -ForegroundColor Red
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

function Export-DevicesToJson {
    param (
        [Parameter(Mandatory=$true)]
        $Devices,
        [Parameter(Mandatory=$false)]
        [string]$SpecificDeviceId
    )
    try {
        $outputDir = Join-Path -Path (Get-Location) -ChildPath "output"
        if (-not (Test-Path -Path $outputDir)) {
            New-Item -Path $outputDir -ItemType Directory | Out-Null
            Write-Host "📁 Created output directory: $outputDir" -ForegroundColor Gray
        }
        $timestamp = Get-Date -Format "yyyyMMdd_HHmmss"
        if ($SpecificDeviceId) {
            $deviceName = $Devices.deviceName -replace '[\\\/:\*\?\"\<\>\|]', '_'
            if (-not $deviceName) { $deviceName = $SpecificDeviceId }
            $fileName = "WindowsManagedDevice_${deviceName}_${timestamp}.json"
            $filePath = Join-Path -Path $outputDir -ChildPath $fileName
            $Devices | ConvertTo-Json -Depth 10 | Out-File -FilePath $filePath -Encoding UTF8
            Write-Host "💾 Exported device to: $filePath" -ForegroundColor Green
        } else {
            $fileName = "WindowsManagedDevices_${timestamp}.json"
            $filePath = Join-Path -Path $outputDir -ChildPath $fileName
            $Devices | ConvertTo-Json -Depth 10 | Out-File -FilePath $filePath -Encoding UTF8
            Write-Host "💾 Exported devices to: $filePath" -ForegroundColor Green
        }
        return $filePath
    } catch {
        Write-Host "❌ Error exporting devices to JSON: $_" -ForegroundColor Red
        return $null
    }
}

function Show-DeviceDetails {
    param (
        [Parameter(Mandatory=$true)]
        $Device
    )
    Write-Host "📋 Windows Managed Device Details:" -ForegroundColor Cyan
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    
    # Top-level fields
    foreach ($field in @(
        '@odata.type','id','userId','deviceName','ownerType','managedDeviceOwnerType','managementState','enrolledDateTime','lastSyncDateTime','chassisType','operatingSystem','deviceType','complianceState','jailBroken','managementAgent','osVersion','easActivated','easDeviceId','easActivationDateTime','aadRegistered','azureADRegistered','deviceEnrollmentType','lostModeState','activationLockBypassCode','emailAddress','azureActiveDirectoryDeviceId','azureADDeviceId','deviceRegistrationState','deviceCategoryDisplayName','isSupervised','exchangeLastSuccessfulSyncDateTime','exchangeAccessState','exchangeAccessStateReason','remoteAssistanceSessionUrl','remoteAssistanceSessionErrorDetails','isEncrypted','userPrincipalName','model','manufacturer','imei','complianceGracePeriodExpirationDateTime','serialNumber','phoneNumber','androidSecurityPatchLevel','userDisplayName','wiFiMacAddress','subscriberCarrier','meid','totalStorageSpaceInBytes','freeStorageSpaceInBytes','managedDeviceName','partnerReportedThreatState','retireAfterDateTime','preferMdmOverGroupPolicyAppliedDateTime','autopilotEnrolled','requireUserEnrollmentApproval','managementCertificateExpirationDate','iccid','udid','ethernetMacAddress','physicalMemoryInBytes','processorArchitecture','specificationVersion','joinType','skuFamily','securityPatchLevel','skuNumber','managementFeatures','enrollmentProfileName','bootstrapTokenEscrowed','deviceFirmwareConfigurationInterfaceManaged','notes','windowsActiveMalwareCount','windowsRemediatedMalwareCount')) {
        if ($Device.PSObject.Properties[$field]) {
            Write-Host ("   • {0}: {1}" -f $field, $Device.$field) -ForegroundColor Green
        }
    }

    # Arrays
    if ($Device.roleScopeTagIds) {
        Write-Host "   • roleScopeTagIds:" -ForegroundColor Green
        foreach ($item in $Device.roleScopeTagIds) {
            Write-Host "     · $item" -ForegroundColor Yellow
        }
    }
    if ($Device.usersLoggedOn) {
        Write-Host "   • usersLoggedOn:" -ForegroundColor Green
        foreach ($user in $Device.usersLoggedOn) {
            Write-Host "     · userId: $($user.userId), lastLogOnDateTime: $($user.lastLogOnDateTime)" -ForegroundColor Yellow
        }
    }
    if ($Device.chromeOSDeviceInfo) {
        Write-Host "   • chromeOSDeviceInfo:" -ForegroundColor Green
        foreach ($chrome in $Device.chromeOSDeviceInfo) {
            Write-Host "     · name: $($chrome.name), value: $($chrome.value), valueType: $($chrome.valueType), updatable: $($chrome.updatable)" -ForegroundColor Yellow
        }
    }
    if ($Device.deviceActionResults) {
        Write-Host "   • deviceActionResults:" -ForegroundColor Green
        foreach ($action in $Device.deviceActionResults) {
            Write-Host "     · actionName: $($action.actionName), actionState: $($action.actionState), startDateTime: $($action.startDateTime), lastUpdatedDateTime: $($action.lastUpdatedDateTime)" -ForegroundColor Yellow
        }
    }

    # Nested objects
    if ($Device.hardwareInformation) {
        Write-Host "   • hardwareInformation:" -ForegroundColor Green
        $hw = $Device.hardwareInformation
        foreach ($field in @(
            '@odata.type','serialNumber','totalStorageSpace','freeStorageSpace','imei','meid','manufacturer','model','phoneNumber','subscriberCarrier','cellularTechnology','wifiMac','operatingSystemLanguage','isSupervised','isEncrypted','batterySerialNumber','batteryHealthPercentage','batteryChargeCycles','isSharedDevice','tpmSpecificationVersion','operatingSystemEdition','deviceFullQualifiedDomainName','deviceGuardVirtualizationBasedSecurityHardwareRequirementState','deviceGuardVirtualizationBasedSecurityState','deviceGuardLocalSystemAuthorityCredentialGuardState','osBuildNumber','operatingSystemProductType','ipAddressV4','subnetAddress','esimIdentifier','systemManagementBIOSVersion','tpmManufacturer','tpmVersion','batteryLevelPercentage','residentUsersCount','productName','deviceLicensingStatus','deviceLicensingLastErrorCode','deviceLicensingLastErrorDescription')) {
            if ($hw.PSObject.Properties[$field]) {
                Write-Host ("     · {0}: {1}" -f $field, $hw.$field) -ForegroundColor Yellow
            }
        }
        if ($hw.sharedDeviceCachedUsers) {
            Write-Host "     · sharedDeviceCachedUsers:" -ForegroundColor Yellow
            foreach ($user in $hw.sharedDeviceCachedUsers) {
                Write-Host "       - userPrincipalName: $($user.userPrincipalName), dataToSync: $($user.dataToSync), dataQuota: $($user.dataQuota), dataUsed: $($user.dataUsed)" -ForegroundColor Magenta
            }
        }
        if ($hw.wiredIPv4Addresses) {
            Write-Host "     · wiredIPv4Addresses:" -ForegroundColor Yellow
            foreach ($ip in $hw.wiredIPv4Addresses) {
                Write-Host "       - $ip" -ForegroundColor Magenta
            }
        }
    }
    if ($Device.configurationManagerClientEnabledFeatures) {
        Write-Host "   • configurationManagerClientEnabledFeatures:" -ForegroundColor Green
        $cm = $Device.configurationManagerClientEnabledFeatures
        foreach ($field in @('inventory','modernApps','resourceAccess','deviceConfiguration','compliancePolicy','windowsUpdateForBusiness','endpointProtection','officeApps')) {
            if ($cm.PSObject.Properties[$field]) {
                Write-Host ("     · {0}: {1}" -f $field, $cm.$field) -ForegroundColor Yellow
            }
        }
    }
    if ($Device.deviceHealthAttestationState) {
        Write-Host "   • deviceHealthAttestationState:" -ForegroundColor Green
        $dh = $Device.deviceHealthAttestationState
        foreach ($field in @(
            '@odata.type','lastUpdateDateTime','contentNamespaceUrl','deviceHealthAttestationStatus','contentVersion','issuedDateTime','attestationIdentityKey','resetCount','restartCount','dataExcutionPolicy','bitLockerStatus','bootManagerVersion','codeIntegrityCheckVersion','secureBoot','bootDebugging','operatingSystemKernelDebugging','codeIntegrity','testSigning','safeMode','windowsPE','earlyLaunchAntiMalwareDriverProtection','virtualSecureMode','pcrHashAlgorithm','bootAppSecurityVersion','bootManagerSecurityVersion','tpmVersion','pcr0','secureBootConfigurationPolicyFingerPrint','codeIntegrityPolicy','bootRevisionListInfo','operatingSystemRevListInfo','healthStatusMismatchInfo','healthAttestationSupportedStatus','memoryIntegrityProtection','memoryAccessProtection','virtualizationBasedSecurity','firmwareProtection','systemManagementMode','securedCorePC')) {
            if ($dh.PSObject.Properties[$field]) {
                Write-Host ("     · {0}: {1}" -f $field, $dh.$field) -ForegroundColor Yellow
            }
        }
    }
    if ($Device.configurationManagerClientHealthState) {
        Write-Host "   • configurationManagerClientHealthState:" -ForegroundColor Green
        $ch = $Device.configurationManagerClientHealthState
        foreach ($field in @('state','errorCode','lastSyncDateTime')) {
            if ($ch.PSObject.Properties[$field]) {
                Write-Host ("     · {0}: {1}" -f $field, $ch.$field) -ForegroundColor Yellow
            }
        }
    }
    if ($Device.configurationManagerClientInformation) {
        Write-Host "   • configurationManagerClientInformation:" -ForegroundColor Green
        $ci = $Device.configurationManagerClientInformation
        foreach ($field in @('clientIdentifier','isBlocked','clientVersion')) {
            if ($ci.PSObject.Properties[$field]) {
                Write-Host ("     · {0}: {1}" -f $field, $ci.$field) -ForegroundColor Yellow
            }
        }
    }
    if ($Device.deviceIdentityAttestationDetail) {
        Write-Host "   • deviceIdentityAttestationDetail:" -ForegroundColor Green
        $di = $Device.deviceIdentityAttestationDetail
        foreach ($field in @('deviceIdentityAttestationStatus')) {
            if ($di.PSObject.Properties[$field]) {
                Write-Host ("     · {0}: {1}" -f $field, $di.$field) -ForegroundColor Yellow
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
    $devices = Get-WindowsManagedDevices -SpecificDeviceId $DeviceId
    if ($ExportToJson) {
        Export-DevicesToJson -Devices $devices -SpecificDeviceId $DeviceId
    }
    if ($DeviceId) {
        Show-DeviceDetails -Device $devices
    } else {
        if ($devices.value -and $devices.value.Count -gt 0) {
            Write-Host "📊 Found $($devices.value.Count) Windows managed device(s)" -ForegroundColor Green
            Write-Host ""
            for ($i = 0; $i -lt $devices.value.Count; $i++) {
                Write-Host "Device $($i + 1):" -ForegroundColor Magenta
                Show-DeviceDetails -Device $devices.value[$i]
            }
        } elseif ($devices -and -not $devices.value) {
            Write-Host "📊 Found 1 Windows managed device" -ForegroundColor Green
            Write-Host ""
            Show-DeviceDetails -Device $devices
        } else {
            Write-Host "📊 No Windows managed devices found" -ForegroundColor Yellow
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