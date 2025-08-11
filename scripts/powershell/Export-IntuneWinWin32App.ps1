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
    HelpMessage="Specify the Intune application ID to retrieve")]
    [ValidateNotNullOrEmpty()]
    [string]$AppId,
    
    [Parameter(Mandatory=$false,
    HelpMessage="Specify the application version")]
    [string]$AppVersion = "1",
    
    [Parameter(Mandatory=$false,
    HelpMessage="Enable verbose debug output")]
    [bool]$EnableDebug = $false
)

# Usage Examples:
# .\Get-IntuneAppContent.ps1 -TenantId "your-tenant-id" -ClientId "your-client-id" -ClientSecret "your-client-secret" -AppId "736e8c49-a5e8-4079-9719-622700592cb3"
# .\Get-IntuneAppContent.ps1 -TenantId "your-tenant-id" -ClientId "your-client-id" -ClientSecret "your-client-secret" -AppId "5e5eb9ce-32be-4466-a396-8c4fdd63a474" -AppVersion "2"
# .\Get-IntuneAppContent.ps1 -TenantId "your-tenant-id" -ClientId "your-client-id" -ClientSecret "your-client-secret" -AppId "your-app-id" -EnableDebug $true

# Authentication is handled via direct REST API calls

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

function Write-Log {
    param(
        [Parameter(Mandatory=$true)]
        [string]$Message,
        [Parameter(Mandatory=$false)]
        [ValidateSet("Info", "Success", "Warning", "Error")]
        [string]$Level = "Info"
    )
    
    $timestamp = Get-Date -Format "yyyy-MM-dd HH:mm:ss"
    $icon = switch ($Level) {
        "Success" { "✅" }
        "Warning" { "⚠️ " }
        "Error" { "❌" }
        default { "ℹ️ " }
    }
    
    $color = switch ($Level) {
        "Success" { "Green" }
        "Warning" { "Yellow" }
        "Error" { "Red" }
        default { "Cyan" }
    }
    
    Write-Host "$icon $timestamp - $Message" -ForegroundColor $color
}

	##############################
    	# --- Add Cert TO WebRequest    #
	###############################

function Invoke-MtlsRestRequest {
    param(
        [string]$Url,
        [string]$Method = "PUT",
        [string]$Body,
        [System.Security.Cryptography.X509Certificates.X509Certificate2]$Cert,
        [hashtable]$Headers
    )
    $req = [System.Net.HttpWebRequest]::Create($Url)
    $req.Method = $Method
    $req.ContentType = "application/json; charset=utf-8"
    $req.Timeout = 60000
    if ($Cert) { $null = $req.ClientCertificates.Add($Cert) }
    foreach ($key in $Headers.Keys) {
        if ($key -ne "Content-Type") {
            $req.Headers[$key] = $Headers[$key]
        }
    }
    if ($Body) {
        $bytes = [System.Text.Encoding]::UTF8.GetBytes($Body)
        $req.ContentLength = $bytes.Length
        $stream = $req.GetRequestStream()
        $stream.Write($bytes, 0, $bytes.Length)
        $stream.Close()
    }
    try {
        $resp = $req.GetResponse()
        $reader = New-Object IO.StreamReader $resp.GetResponseStream()
        $result = $reader.ReadToEnd()
        $reader.Close()
        $resp.Close()
        return $result
    } catch {
        Write-Host "ERROR: $($_.Exception.Message)"
        if ($_.Exception.Response) {
            $stream = $_.Exception.Response.GetResponseStream()
            $reader = New-Object IO.StreamReader $stream
            $respBody = $reader.ReadToEnd()
            Write-Host "RESPONSE BODY: $respBody"
        }
        return $null
    }
}

	#########################################
    	# Decompress intunesidecarepayload  #
	##########################################
function Decompress-IntuneSidecarPayload {
    param(
        [Parameter(Mandatory)]
        [string]$CompressedBase64
    )
    # 1. Decode from Base64
    $allBytes = [Convert]::FromBase64String($CompressedBase64)
    
    # 2. Extract the original data length (first 4 bytes, little endian)
    $origLen = [BitConverter]::ToInt32($allBytes, 0)
    
    # 3. The actual compressed data is the remaining bytes
    $gzipBytes = $allBytes[4..($allBytes.Length-1)]
    
    # 4. Decompress GZIP
    $ms = New-Object IO.MemoryStream(,$gzipBytes)
    $gs = New-Object IO.Compression.GzipStream($ms, [IO.Compression.CompressionMode]::Decompress)
    $decompressed = New-Object byte[] $origLen
    $read = $gs.Read($decompressed, 0, $origLen)
    $gs.Dispose()
    $ms.Dispose()
    
    # 5. Convert back to string
    [Text.Encoding]::UTF8.GetString($decompressed, 0, $read)
}

	##############################
    	# Get MDM Cert and IDS	   #
	###############################

function Get-IntuneMDMCertAndIDs {
    Write-Log "Searching for Intune MDM device certificate..."
    $mdmOid = '1.2.840.113556.5.6'
    $issuer = 'Microsoft Intune MDM Device CA'
    $store = New-Object System.Security.Cryptography.X509Certificates.X509Store("My", "LocalMachine")
    $store.Open('ReadOnly')
    $cert = $store.Certificates | Where-Object {
        ($_.Issuer -like "*$issuer*") -and ($_.Extensions | Where-Object { $_.Oid.Value -eq $mdmOid }).Count -gt 0
    } | Sort-Object NotAfter -Descending | Select-Object -First 1
    $store.Close()
    if (-not $cert) { throw "No valid Intune MDM certificate found." }
    if ($cert.Subject -notmatch '^CN=([\da-fA-F-]{36})') {
        throw "Could not parse DeviceId from certificate subject."
    }
    $deviceId = [guid]$Matches[1]
    $accountId = $null
    foreach ($ext in $cert.Extensions) {
        if ($ext.Oid.Value -eq $mdmOid) {
            $bytes = if ($ext.RawData.Length -eq 16) {
                $ext.RawData
            } elseif ($ext.RawData.Length -eq 18 -and $ext.RawData[0] -eq 4 -and $ext.RawData[1] -eq 16) {
                $ext.RawData[2..17]
            }
            if ($bytes) { $accountId = [guid][byte[]]$bytes }
            break
        }
    }
    $certBytes = $cert.Export([System.Security.Cryptography.X509Certificates.X509ContentType]::Cert)
    $certBase64 = [Convert]::ToBase64String($certBytes)
    return [PSCustomObject]@{
        Cert         = $cert
        DeviceId     = $deviceId
        AccountId    = $accountId
        CertBlob     = $certBase64
    }
}
	##############################
    	# Get Region/ASU location    #
	###############################

function Get-IntuneLocationServiceUrls {
    $urls = @()
    $key = "HKLM:\SOFTWARE\Microsoft\Provisioning\OMADM\Accounts"
    if (Test-Path $key) {
        foreach ($sub in Get-ChildItem $key) {
            $addrPath = "$($sub.PSPath)\Protected\AddrInfo"
            try {
                $addr = Get-ItemProperty -Path $addrPath -Name Addr -ErrorAction Stop
                if ($addr.Addr -and $addr.Addr -notlike "*checkin.dm.microsoft.com*") {
                    $uri = [uri]$addr.Addr
                    $fqdn = "$($uri.Scheme)://$($uri.Host)"
                    if ($urls -notcontains $fqdn) {
                        $urls += $fqdn
                    }
                }
            } catch { }
        }
    }
    if (-not $urls) {
        $urls = @("https://manage.microsoft.com")
        Write-Log "No registry URLs found, using default: $($urls[0])"
    } else {
        Write-Log "Using LocationService URLs: $($urls -join ', ')"
    }
    return $urls
}

	##############################
    	# QUery Region/ASU location    #
	###############################

function Query-LocationService {
    param (
        [string[]]$LocationServiceUrls,
        [System.Security.Cryptography.X509Certificates.X509Certificate2]$Cert
    )
    $discoPath = "/RestUserAuthLocationService/RestUserAuthLocationService/Certificate/ServiceAddresses"
    foreach ($fqdn in $LocationServiceUrls) {
        $url = "$fqdn$discoPath"
        Write-Log "Querying discovery endpoint: $url"
        try {
            $req = [System.Net.HttpWebRequest]::Create($url)
            $req.Method = "GET"
            $req.ClientCertificates.Add($Cert)
            $req.Timeout = 30000
            $req.Headers.Add("client-request-id", ([guid]::NewGuid()).Guid)
            [System.Net.ServicePointManager]::SecurityProtocol = [System.Net.SecurityProtocolType]::Tls12
            $resp = $req.GetResponse()
            $reader = New-Object IO.StreamReader $resp.GetResponseStream()
            $json = $reader.ReadToEnd()
            $reader.Close()
            $resp.Close()
            Write-Log "Discovery response JSON received."
            $result = $json | ConvertFrom-Json
            foreach ($entry in $result) {
                if ($entry.IsPrimary -eq $true) {
                    $svc = $entry.Services | Where-Object { $_.ServiceName -eq "SideCarGatewayService" } | Select-Object -First 1
                    if ($svc -and $svc.Url -is [string] -and $svc.Url.Trim()) {
                        $cleanUrl = ([string]$svc.Url).Trim()
                        Write-Log "Found SideCarGatewayService URL: $cleanUrl"
                        return $cleanUrl
                    }
                }
            }
        } catch {
            $msg = $_.Exception.Message
            Write-Log ("Discovery request failed for {0}: {1}" -f $fqdn, $msg)
            # If it's an SSL/TLS error, rethrow to be handled by outer catch
            if ($msg -match '(?i)ssl|tls') {
                throw
            }
        }
    }
    Write-Log "No valid SideCarGatewayService URL could be discovered. Returning empty string."
    return [string]::Empty
}
	##############################
    	# GET IME Version    #
	###############################

function Get-IntuneManagementExtensionVersion {
    # Original method - Inventories key
    $inventoryPath = "HKLM:\SOFTWARE\Microsoft\IntuneManagementExtension\Inventories"
    $subKeys = Get-ChildItem -Path $inventoryPath -ErrorAction SilentlyContinue

    foreach ($subKey in $subKeys) {
        $props = Get-ItemProperty -Path $subKey.PSPath
        if ($props.Name -eq "Microsoft Intune Management Extension" -and $props.Version) {
            return $props.Version
        }
    }

    # Fallback - EnterpriseDesktopAppManagement key
    $baseEDAMPath = "HKLM:\SOFTWARE\Microsoft\EnterpriseDesktopAppManagement"
    $sidKeys = Get-ChildItem -Path $baseEDAMPath -ErrorAction SilentlyContinue | Where-Object { $_.PSChildName -match '^S-\d-\d+-\d+.*' }

    foreach ($sidKey in $sidKeys) {
        $msiPath = Join-Path -Path $sidKey.PSPath -ChildPath "MSI"
        $appKeys = Get-ChildItem -Path $msiPath -ErrorAction SilentlyContinue
        foreach ($appKey in $appKeys) {
            $props = Get-ItemProperty -Path $appKey.PSPath -ErrorAction SilentlyContinue
            if ($props.DownloadUrlList) {
                foreach ($url in $props.DownloadUrlList) {
                    if ($url -like "*IntuneWindowsAgent.msi") {
                        return $props.ProductVersion
                    }
                }
            }
        }
    }

    return $null
}


##############################
    	# Get-AvailableApps    #
	###############################

function Get-AvailableApps {
    param (
        [string]$endpoint,
        [string]$accountId,
        [string]$deviceId,
        [string]$bearerToken,
        [System.Security.Cryptography.X509Certificates.X509Certificate2]$cert,
        [string]$AppIdToCheck
    )

    $Imeversion = Get-IntuneManagementExtensionVersion
    if (-not $Imeversion) { $Imeversion = "1.91.102.0" }
    $sessionId = [guid]::NewGuid().Guid
    $url = "$endpoint/SideCarGatewaySessions('$sessionId')?api-version=1.5"

    $clientInfo = @{
        DeviceName = $env:COMPUTERNAME
        OperatingSystemVersion = (Get-CimInstance Win32_OperatingSystem).Version
        SideCarAgentVersion = $Imeversion
    } | ConvertTo-Json -Compress

    $body = @{
        Key = $sessionId
        SessionId = $sessionId
        RequestContentType = "GetAvailableApp"
        RequestPayload = "[]"
        ResponseContentType = $null
        ClientInfo = $clientInfo
        ResponsePayload = $null
        CheckinReasonPayload = '{"NotificationID":"00000000-0000-0000-0000-000000000000","NotificationIntent":""}'
    }
    $bodyJson = $body | ConvertTo-Json -Compress

    $headers = @{
        "Authorization"         = "Bearer $bearerToken"
        "client-request-id"     = ([guid]::NewGuid()).Guid
        "AccountId"             = $accountId
        "DeviceId"              = $deviceId
        "Content-Type"          = "application/json; charset=utf-8"
        "Prefer"                = "return-content"
        "Request-Attempt-Count" = "1"
        "Scenario-Type"         = "Windows-GetAvailableApp"
    }

    $response = Invoke-MtlsRestRequest -Url $url -Method "PUT" -Body $bodyJson -Cert $cert -Headers $headers
    if (!$response) {
        Write-Warning "No response from GetAvailableApps"
        return $null
    }

    try {
        $result = $response | ConvertFrom-Json
    } catch {
        Write-Log "Failed to parse response as JSON" "Error"
        return $null
    }
    $payload = $result.ResponsePayload

    if (-not $payload) {
        Write-Log "No ResponsePayload present" "Warning"
        return $null
    }

    # -- Decompress and parse the payload --
    try {
        $decoded = Decompress-IntuneSidecarPayload -CompressedBase64 $payload
        $apps = $decoded | ConvertFrom-Json
    } catch {
        Write-Log "Failed to decompress or parse available apps payload" "Error"
        return $null
    }

    # -- Search for the AppId --
    $match = $apps | Where-Object { $_.id -eq $AppIdToCheck }
    if ($match) {
        Write-Log "Found AppId [$AppIdToCheck] in available apps:"
        $match | ConvertTo-Json -Depth 8
    } else {
        Write-Log "AppId [$AppIdToCheck] not found in available apps."
    }

    # Return all apps (or change to return $match if you only want that)
    return $apps
}
	##############################
    	# Get required apps    #
	###############################


function Get-RequiredApps {
    param (
        [string]$endpoint,
        [string]$accountId,
        [string]$deviceId,
        [string]$bearerToken,
        [System.Security.Cryptography.X509Certificates.X509Certificate2]$cert,
        [string]$AppIdToCheck
    )

    $Imeversion = Get-IntuneManagementExtensionVersion
    if (-not $Imeversion) { $Imeversion = "1.91.102.0" }
    $sessionId = [guid]::NewGuid().Guid
    $url = "$endpoint/SideCarGatewaySessions('$sessionId')?api-version=1.5"

    $clientInfo = @{
        DeviceName = $env:COMPUTERNAME
        OperatingSystemVersion = (Get-CimInstance Win32_OperatingSystem).Version
        SideCarAgentVersion = $Imeversion
    } | ConvertTo-Json -Compress

    $body = @{
        Key = $sessionId
        SessionId = $sessionId
        RequestContentType = "RequestApplication"
        RequestPayload = "[]"
        ResponseContentType = $null
        ClientInfo = $clientInfo
        ResponsePayload = $null
        CheckinReasonPayload = '{"NotificationID":"00000000-0000-0000-0000-000000000000","NotificationIntent":""}'
    }
    $bodyJson = $body | ConvertTo-Json -Compress

    $headers = @{
        "Authorization"         = "Bearer $bearerToken"
        "client-request-id"     = ([guid]::NewGuid()).Guid
        "AccountId"             = $accountId
        "DeviceId"              = $deviceId
        "Content-Type"          = "application/json; charset=utf-8"
        "Prefer"                = "return-content"
        "Request-Attempt-Count" = "1"
        "Scenario-Type"         = "Windows-RequestApplication"
    }

    $response = Invoke-MtlsRestRequest -Url $url -Method "PUT" -Body $bodyJson -Cert $cert -Headers $headers
    if (!$response) {
        Write-Warning "No response from GetRequiredApps"
        return $null
    }

    try {
        $result = $response | ConvertFrom-Json
    } catch {
        Write-Log "Failed to parse response as JSON" "Error"
        return $null
    }
    $payload = $result.ResponsePayload

    if (-not $payload) {
        Write-Log "No ResponsePayload present" "Warning"
        return $null
    }

    # -- Decompress and parse the payload --
    try {
        $decoded = Decompress-IntuneSidecarPayload -CompressedBase64 $payload
        $apps = $decoded | ConvertFrom-Json
	
    } catch {
        Write-Log "Failed to decompress or parse required apps payload" "Error"
        return $null
    }

    # -- Search for the AppId --
    $match = $apps | Where-Object { $_.id -eq $AppIdToCheck }
    if ($match) {
        Write-Log "Found AppId [$AppIdToCheck] in required apps:"
        $match | ConvertTo-Json -Depth 8
    } else {
        Write-Log "AppId [$AppIdToCheck] not found in required apps."
    }

    return $apps
}

##############################
    	# Get-SelectedApps    #
	###############################
function Get-SelectedApps {
    param (
        [string]$endpoint,
        [string]$accountId,
        [string]$deviceId,
        [string]$bearerToken,
        [System.Security.Cryptography.X509Certificates.X509Certificate2]$cert,
        [string]$AppIdToCheck,
        [string]$UserSid = $null  # Optional: only needed if you want to fetch as a specific user
    )

    $Imeversion = Get-IntuneManagementExtensionVersion
    if (-not $Imeversion) { $Imeversion = "1.91.102.0" }
    $sessionId = [guid]::NewGuid().Guid
    $url = "$endpoint/SideCarGatewaySessions('$sessionId')?api-version=1.5"

    $clientInfo = @{
        DeviceName = $env:COMPUTERNAME
        OperatingSystemVersion = (Get-CimInstance Win32_OperatingSystem).Version
        SideCarAgentVersion = $Imeversion
    } | ConvertTo-Json -Compress

    # For ESP user phase, IME sets the session context to the user's SID
    # If you have it, include it; otherwise, it will still work for device context
    $body = @{
        Key = $sessionId
        SessionId = $sessionId
        RequestContentType = "GetSelectedApp"
        RequestPayload = "[]"
        ResponseContentType = $null
        ClientInfo = $clientInfo
        ResponsePayload = $null
        CheckinReasonPayload = '{"NotificationID":"00000000-0000-0000-0000-000000000000","NotificationIntent":""}'
    }

    if ($UserSid) {
        # Not strictly required, but some IME calls include the user SID for user context
        $body.UserSid = $UserSid
    }

    $bodyJson = $body | ConvertTo-Json -Compress

    $headers = @{
        "Authorization"         = "Bearer $bearerToken"
        "client-request-id"     = ([guid]::NewGuid()).Guid
        "AccountId"             = $accountId
        "DeviceId"              = $deviceId
        "Content-Type"          = "application/json; charset=utf-8"
        "Prefer"                = "return-content"
        "Request-Attempt-Count" = "1"
        "Scenario-Type"         = "Windows-GetSelectedApp"
    }

    $response = Invoke-MtlsRestRequest -Url $url -Method "PUT" -Body $bodyJson -Cert $cert -Headers $headers
    if (!$response) {
        Write-Warning "No response from GetSelectedApps"
        return $null
    }

    try {
        $result = $response | ConvertFrom-Json
    } catch {
        Write-Log "Failed to parse response as JSON" "Error"
        return $null
    }
    $payload = $result.ResponsePayload

    if (-not $payload) {
        Write-Log "No ResponsePayload present" "Warning"
        return $null
    }

    # -- Decompress and parse the payload --
    try {
        $decoded = Decompress-IntuneSidecarPayload -CompressedBase64 $payload
        $apps = $decoded | ConvertFrom-Json
    } catch {
        Write-Log "Failed to decompress or parse selected apps payload" "Error"
        return $null
    }

    # -- Search for the AppId --
    $match = $apps | Where-Object { $_.id -eq $AppIdToCheck }
    if ($match) {
        Write-Log "Found AppId [$AppIdToCheck] in selected apps:"
        $match | ConvertTo-Json -Depth 8
    } else {
        Write-Log "AppId [$AppIdToCheck] not found in selected apps."
    }

    return $apps
}


	##############################
    	# Get Content Info Function     #
	###############################

function Send-GetContentInfoRequest {
    param (
        [string]$endpoint,
        [string]$accountId,
        [string]$deviceId,
        [string]$bearerToken,
        [System.Security.Cryptography.X509Certificates.X509Certificate2]$cert,
        [string]$certBlob,
        [string]$applicationId,
        [string]$applicationVersion
    )
    $Imeversion = Get-IntuneManagementExtensionVersion
    if (-not $Imeversion) {
        $Imeversion = "1.91.102.0"
        Write-Log "IME version not found in registry, falling back to: $Imeversion"
    } else {
        Write-Log "Detected IME version: $Imeversion"
    }
    $sessionId = [guid]::NewGuid().Guid
    $url = "$endpoint/SideCarGatewaySessions('$sessionId')?api-version=1.5"
    Write-Log "Sending PUT to $url (GetContentInfo) for AppId [$applicationId]"

    $reqPayloadObj = @{
        ContentInfo            = $null
        Intent                 = 1
        CertificateBlob        = $certBlob
        DecryptInfo            = $null
        UploadLocation         = $null
        ApplicationVersion     = $applicationVersion
        ApplicationId          = $applicationId
    }
    $reqPayload = $reqPayloadObj | ConvertTo-Json -Compress -Depth 6
    $body = @{
        Key                   = $sessionId
        SessionId             = $sessionId
        RequestContentType    = "GetContentInfo"
        RequestPayload        = $reqPayload
        ResponseContentType   = $null
        ClientInfo            = @{
            DeviceName = $env:COMPUTERNAME
            OperatingSystemVersion = (Get-CimInstance Win32_OperatingSystem).Version
            SideCarAgentVersion = $Imeversion
        } | ConvertTo-Json -Compress
        ResponsePayload = $null
        CheckinReasonPayload = '{"NotificationID":"00000000-0000-0000-0000-000000000000","NotificationIntent":""}'
    }
    $bodyJson = $body | ConvertTo-Json -Compress -Depth 10
    $headers = @{
        "Authorization" = "Bearer $bearerToken"
        "client-request-id" = ([guid]::NewGuid()).Guid
        "AccountId" = $accountId
        "DeviceId" = $deviceId
        "Prefer" = "return-content"
        "Request-Attempt-Count" = "1"
        "Scenario-Type" = "Windows-GetContentInfo"
    }
    return Invoke-MtlsRestRequest -Url $url -Method "PUT" -Body $bodyJson -Cert $cert -Headers $headers
}


	######################################
    	# Get Decryption Info Function    #
	######################################

function Get-DecryptionInfoFromResponse {
    param(
        [Parameter(Mandatory)]
        [object]$ResponseObject
    )

    # Step 1: Parse .ResponsePayload JSON (if string)
    $respPayload = $ResponseObject.ResponsePayload
    if ($respPayload -is [string]) {
        $respPayload = $respPayload | ConvertFrom-Json
    }

    # Step 2: Pull and parse DecryptInfo XML
    $decryptInfoXml = $respPayload.DecryptInfo
    if (-not $decryptInfoXml) {
        Write-Warning "No DecryptInfo in response."
        return $null
    }

    # Parse the XML and get the EncryptedContent node
    [xml]$decryptInfo = $decryptInfoXml
    $encryptedContent = $decryptInfo.EncryptedMessage.EncryptedContent
    if (-not $encryptedContent) {
        Write-Warning "No EncryptedContent found in DecryptInfo XML."
        return $null
    }

    # --- Use classic assembly load
    [Reflection.Assembly]::LoadWithPartialName("System.Security") | Out-Null
    [Reflection.Assembly]::LoadWithPartialName("System.Security.Cryptography.Pkcs") | Out-Null
    $bytes = [Convert]::FromBase64String($encryptedContent)
    $cms = New-Object System.Security.Cryptography.Pkcs.EnvelopedCms
    $cms.Decode($bytes)
    $cms.Decrypt()
    $utf8Json = [System.Text.Encoding]::UTF8.GetString($cms.ContentInfo.Content)
    $decryptResult = $utf8Json | ConvertFrom-Json
    return $decryptResult
}
	######################################
    	# Get BIN File Location Function    #
	######################################

function Get-UploadLocationAndDownload {
    param (
        [Parameter(Mandatory)]
        [object]$Response,
        [Parameter(Mandatory)]
        [string]$OutputPath,
        [int]$MaxRetries = 3,
        [int]$DelaySeconds = 5
    )

    # Parse JSON if needed
    if ($Response -is [string]) {
        $ResponseObj = $Response | ConvertFrom-Json
    } else {
        $ResponseObj = $Response
    }
    if (-not $ResponseObj.ResponsePayload) {
        Write-Host "No ResponsePayload in response!"
        return $null
    }
    $payload = $ResponseObj.ResponsePayload | ConvertFrom-Json
    if (-not $payload.ContentInfo) {
        Write-Host "No ContentInfo in payload!"
        return $null
    }
    $contentInfo = $payload.ContentInfo | ConvertFrom-Json
    $uploadUrl = $contentInfo.UploadLocation
    if (-not $uploadUrl) {
        Write-Host "No UploadLocation found."
        return $null
    }
    Write-Log "Found UploadLocation: $uploadUrl"
    Write-Log "Downloading to: $OutputPath"

    # --- First Try BITS ---
    $attempt = 1
    $bitsFailed = $false
    while ($attempt -le $MaxRetries) {
        try {
            Start-BitsTransfer -Source $uploadUrl -Destination $OutputPath -ErrorAction Stop
            Write-Log "Download complete via BITS on attempt $attempt."
            return $OutputPath
        } catch {
            Write-Log "BITS attempt $attempt failed: $_"
            if ($_ -match "0x80070006" -or $attempt -eq $MaxRetries) {
                $bitsFailed = $true
                break
            }
            Write-Host "BITS retrying in $DelaySeconds seconds..."
            Start-Sleep -Seconds $DelaySeconds
        }
        $attempt++
    }

    # --- Fallback to curl.exe ---
    if ($bitsFailed) {
        try {
            Write-Log "Falling back to curl.exe."
            $curlCmd = "curl.exe -L --fail --output `"$OutputPath`" `"$uploadUrl`""
            $null = Invoke-Expression $curlCmd
            if (Test-Path $OutputPath) {
                Write-Log "Download complete via curl.exe."
                return $OutputPath
            } else {
                Write-Log "curl.exe failed to download the file."
            }
        } catch {
            Write-Log "curl.exe failed: $_"
        }
    }

    Write-Log "All download methods failed. No file downloaded."
    return $null
}

##############################
# --- Graph API Token Function    #
##############################
function Get-GraphApiAccessToken {
    param(
        [Parameter(Mandatory=$true)]
        [string]$TenantId,
        [Parameter(Mandatory=$true)]
        [string]$ClientId,
        [Parameter(Mandatory=$true)]
        [string]$ClientSecret
    )
    
    try {
        Write-Log "🔑 Obtaining access token from Microsoft Graph API" "Info"
        
        $tokenEndpoint = "https://login.microsoftonline.com/$TenantId/oauth2/v2.0/token"
        $body = @{
            client_id     = $ClientId
            client_secret = $ClientSecret
            scope         = "https://graph.microsoft.com/.default"
            grant_type    = "client_credentials"
        }
        
        $response = Invoke-RestMethod -Uri $tokenEndpoint -Method Post -Body $body -ContentType "application/x-www-form-urlencoded"
        
        if ($response.access_token) {
            Write-Log "✅ Access token obtained successfully" "Success"
            Write-DebugInfo "Token type" $response.token_type
            Write-DebugInfo "Expires in" "$($response.expires_in) seconds"
            return $response.access_token
        } else {
            Write-Log "❌ No access token in response" "Error"
            return $null
        }
    }
    catch {
        Write-Log "❌ Failed to obtain access token: $($_.Exception.Message)" "Error"
        if ($_.Exception.Response) {
            $statusCode = $_.Exception.Response.StatusCode
            Write-Log "HTTP Status Code: $statusCode" "Error"
        }
        return $null
    }
}

	##############################
    	# Decrypt Function	    #
	###############################
	function Decrypt-IntuneWinFile {
   	 param(
        [Parameter(Mandatory)] [string]$IntuneWinFile,
        [Parameter(Mandatory)] [string]$OutputZipFile,
        [Parameter(Mandatory)] [string]$Base64Key,
        [Parameter(Mandatory)] [string]$Base64IV
    )

    # Load AES
    Add-Type -AssemblyName System.Security
    [Reflection.Assembly]::LoadWithPartialName("System.Security.Cryptography") | Out-Null

    # Read key/IV
    $key = [Convert]::FromBase64String($Base64Key)
    $iv  = [Convert]::FromBase64String($Base64IV)

    # Open streams
    $inStream  = [System.IO.File]::Open($IntuneWinFile, [System.IO.FileMode]::Open, [System.IO.FileAccess]::Read)
    $outStream = [System.IO.File]::Open($OutputZipFile, [System.IO.FileMode]::Create, [System.IO.FileAccess]::Write)

    try {
        # Skip the first 48 bytes
        $inStream.Seek(48, [System.IO.SeekOrigin]::Begin) | Out-Null

        # Set up AES CBC decryption
        $aes = [System.Security.Cryptography.Aes]::Create()
        $aes.Key = $key
        $aes.IV  = $iv
        $aes.Mode = [System.Security.Cryptography.CipherMode]::CBC
        $aes.Padding = [System.Security.Cryptography.PaddingMode]::PKCS7

        $decryptor = $aes.CreateDecryptor()
        $cryptoStream = New-Object System.Security.Cryptography.CryptoStream($outStream, $decryptor, [System.Security.Cryptography.CryptoStreamMode]::Write)

        # Decrypt in chunks (2MB)
        $buffer = New-Object byte[] 2097152
        while (($read = $inStream.Read($buffer, 0, $buffer.Length)) -gt 0) {
            $cryptoStream.Write($buffer, 0, $read)
        }

        $cryptoStream.FlushFinalBlock()
        Write-Log "Decryption complete: $OutputZipFile"
    }
    finally {
        $cryptoStream.Close()
        $outStream.Close()
        $inStream.Close()
        $aes.Dispose()
    }
}


##############################
# --- Main Execution    #
##############################

try {
    Write-Host "🚀 Starting Intune Application Content Retrieval" -ForegroundColor Cyan
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    Write-Host ""
    
    if ($EnableDebug) {
        Write-Host "🐛 Debug mode enabled - verbose output will be shown" -ForegroundColor Yellow
        Write-Host ""
    }
    
    Write-Log "🔍 Retrieving device certificate and identifiers" "Info"
    $mdm      = Get-IntuneMDMCertAndIDs
    $cert     = $mdm.Cert
    $deviceId = $mdm.DeviceId
    $accountId = $mdm.AccountId
    $certBlob  = $mdm.CertBlob
    Write-Log "📜 Certificate Subject: $($cert.Subject)" "Info"
    Write-Log "📱 Device ID: $deviceId" "Info"
    Write-Log "🏢 Account ID: $accountId" "Info"
    Write-Log "🔑 Obtaining access token from Microsoft Graph" "Info"
    $bearerToken = Get-GraphApiAccessToken -TenantId $TenantId -ClientId $ClientId -ClientSecret $ClientSecret
    if ($bearerToken) {
        Write-Log "✅ Access token obtained successfully" "Success"
    } else {
        Write-Log "❌ Failed to obtain access token" "Error"
        throw "Could not obtain access token"
    }
  


##############################
# ---fixing  endpoint url    #
##############################
$urls = Get-IntuneLocationServiceUrls
try {
    $endpoint = Query-LocationService -LocationServiceUrls $urls -Cert $cert
    if ($endpoint -isnot [string]) { $endpoint = [string]$endpoint }
    $endpoint = $endpoint.Trim() -replace '^\s*0\s*', ''
    if (-not $endpoint.StartsWith("https://")) {
        throw "Invalid endpoint returned: '$endpoint'"
    }
} catch {
    $errMsg = $_.Exception.Message
    if ($errMsg -match '(?i)ssl|tls') {
        Write-Host ""
        Write-Host "************************************************************" -ForegroundColor Yellow
        Write-Host "  ERROR: TLS/SSL channel could not be established." -ForegroundColor Red
        Write-Host "  Run this script in an elevated PowerShell window (Run as Administrator)." -ForegroundColor Yellow
        Write-Host "************************************************************" -ForegroundColor Yellow
        Write-Host ""
    } else {
        Write-Log "ERROR: $errMsg"
    }
    return
}
 
  ##############################
# ---get-availableapps  #
##############################

$availableApps = Get-AvailableApps -endpoint $endpoint -accountId $accountId -deviceId $deviceId `
    -bearerToken $bearerToken -cert $cert -AppIdToCheck $appId

  ##############################
# ---get-required apps #
##############################

$requiredApps = Get-RequiredApps -endpoint $endpoint -accountId $accountId -deviceId $deviceId `
    -bearerToken $bearerToken -cert $cert -AppIdToCheck $appId
 ##############################
# ---get-selected apps #
##############################

$selectedApps = Get-SelectedApps -endpoint $endpoint -accountId $accountId -deviceId $deviceId `
    -bearerToken $bearerToken -cert $cert -AppIdToCheck $appId

 
  ##############################
# --- download content info  #
##############################

$response = Send-GetContentInfoRequest -endpoint $endpoint -accountId $accountId -deviceId $deviceId `
    -bearerToken $bearerToken -cert $cert -certBlob $certBlob `
    -applicationId $appId -applicationVersion $appVer

if (-not $response) {
    Write-Log "No response from server for AppId [$appId]. Check if you specified the right AppId (App must also be assigned!)"
    return
}

try {
    $json = $response | ConvertFrom-Json
} catch {
    Write-Log "Failed to parse response as JSON for AppId [$appId]: $_"
    return
}

if (-not $json.ResponsePayload) {
    Write-Log "No ResponsePayload in response for AppId [$appId]."
    return
}

try {
    $payload = $json.ResponsePayload | ConvertFrom-Json
} catch {
    Write-Log "Failed to parse ResponsePayload as JSON for AppId [$appId]: $_"
    return
}

if (-not $payload.ContentInfo) {
    Write-Log "No ContentInfo found in payload for AppId [$appId]."
    return
}

Write-Log "Successfully found ContentInfo for AppId [$appId]."


    	##############################
    	# --- find decryption info    #
	###############################
    $decryptionInfo = Get-DecryptionInfoFromResponse -ResponseObject $json
    if ($decryptionInfo) {
        Write-Log "Encryption Key: $($decryptionInfo.EncryptionKey)"
        Write-Log "IV: $($decryptionInfo.IV)"
        # Output others as needed
    } else {
        Write-Log "No decryption info found in response!"
    }
   	##############################
    	# --- Download the file    #
	###############################
    		# Create the folder if it doesn't exist
     		 $outputFolder = "C:\Temp"
			if (-not (Test-Path -Path $outputFolder -PathType Container)) {
       		  New-Item -Path $outputFolder -ItemType Directory | Out-Null
		}
   		  $OutputPath = "c:\temp\intunewinfile.bin"
  		  $downloadedFile = Get-UploadLocationAndDownload -Response $json -OutputPath $OutputPath 
   			 if ($downloadedFile) {
    			    Write-Log "Downloaded file: $downloadedFile"
    			} else {
       			 Write-Log "No file downloaded."
  		  }
	##############################################################
          # ------ Decrypt and extract the zip from the Intune Win  #
	##############################################################
	$keyBase64     = $decryptionInfo.EncryptionKey   # Base64 string from previous step
	$ivBase64      = $decryptionInfo.IV              # Base64 string from prus seviotep
        $outputZipFile = "C:\Temp\intunewin_$AppId.decoded.zip"
	Decrypt-IntuneWinFile -IntuneWinFile c:\temp\intunewinfile.bin -OutputZipFile $outputZipFile -Base64Key $keyBase64 -Base64IV $ivBase64

    Write-Host "" 
    Write-Host "🎉 Operation completed successfully!" -ForegroundColor Green
    Write-Host "📦 Decrypted package available at: $outputZipFile" -ForegroundColor Cyan
    Write-Host ""

} catch {
    Write-Log "❌ Script execution failed: $($_.Exception.Message)" "Error"
    exit 1
}