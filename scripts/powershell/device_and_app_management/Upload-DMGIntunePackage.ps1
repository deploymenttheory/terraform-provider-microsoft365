# modified from https://github.com/ugurkocde/IntuneBrew for testing purposes
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
    HelpMessage="Path to the DMG file to upload")]
    [ValidateNotNullOrEmpty()]
    [string]$DmgFilePath,
    
    [Parameter(Mandatory=$true,
    HelpMessage="Display name for the application")]
    [string]$AppDisplayName,
    
    [Parameter(Mandatory=$false,
    HelpMessage="Version number of the application")]
    [string]$AppVersion,
    
    [Parameter(Mandatory=$false,
    HelpMessage="Bundle ID of the application")]
    [string]$AppBundleId,
    
    [Parameter(Mandatory=$false,
    HelpMessage="Description of the application")]
    [string]$AppDescription,
    
    [Parameter(Mandatory=$false,
    HelpMessage="Publisher of the application")]
    [string]$AppPublisher,
    
    [Parameter(Mandatory=$false,
    HelpMessage="Path to a logo PNG file for the application")]
    [string]$LogoFilePath
)

# Helper function to generate reference ID from path
function Get-PathBasedReferenceId {
  param (
      [string]$Path
  )
  
  $segments = $Path -split '/' | Where-Object { $_ -ne '' }
  $refParts = @()
  
  foreach ($segment in $segments) {
      $segmentStr = [string]$segment
      if ($segmentStr -match '{.*}') {
          $paramName = $segmentStr -replace '{|}'
          $paramName = $paramName -replace '-', '_'
          $refParts += "BY_$([string]($paramName.ToUpper()))"
      } else {
          $refParts += [string]($segmentStr.ToUpper())
      }
  }
  
  return $refParts -join '_'
}

# Function to create a new Intune app
function New-IntuneApp {
    param (
        [Parameter(Mandatory=$true)]
        [hashtable]$AppData
    )
    
    try {
        $uri = "https://graph.microsoft.com/beta/deviceAppManagement/mobileApps"
        $response = Invoke-MgGraphRequest -Method POST -Uri $uri -Body ($AppData | ConvertTo-Json -Depth 10)
        
        return $response
    }
    catch {
        Write-Error "Error creating Intune app: $_"
        throw
    }
}

# Function to create a content version for an app
function New-AppContentVersion {
    param (
        [Parameter(Mandatory=$true)]
        [string]$AppId,
        
        [Parameter(Mandatory=$true)]
        [string]$AppType
    )
    
    try {
        $uri = "https://graph.microsoft.com/beta/deviceAppManagement/mobileApps/$AppId/microsoft.graph.$AppType/contentVersions"
        $response = Invoke-MgGraphRequest -Method POST -Uri $uri -Body "{}"
        
        return $response
    }
    catch {
        Write-Error "Error creating content version: $_"
        throw
    }
}

# Function to create a content file for a content version
function New-AppContentFile {
    param (
        [Parameter(Mandatory=$true)]
        [string]$AppId,
        
        [Parameter(Mandatory=$true)]
        [string]$AppType,
        
        [Parameter(Mandatory=$true)]
        [string]$ContentVersionId,
        
        [Parameter(Mandatory=$true)]
        [hashtable]$FileData
    )
    
    try {
        $uri = "https://graph.microsoft.com/beta/deviceAppManagement/mobileApps/$AppId/microsoft.graph.$AppType/contentVersions/$ContentVersionId/files"
        $response = Invoke-MgGraphRequest -Method POST -Uri $uri -Body ($FileData | ConvertTo-Json)
        
        return $response
    }
    catch {
        Write-Error "Error creating content file: $_"
        throw
    }
}

# Function to get content file status
function Get-AppContentFileStatus {
    param (
        [Parameter(Mandatory=$true)]
        [string]$AppId,
        
        [Parameter(Mandatory=$true)]
        [string]$AppType,
        
        [Parameter(Mandatory=$true)]
        [string]$ContentVersionId,
        
        [Parameter(Mandatory=$true)]
        [string]$ContentFileId
    )
    
    try {
        $uri = "https://graph.microsoft.com/beta/deviceAppManagement/mobileApps/$AppId/microsoft.graph.$AppType/contentVersions/$ContentVersionId/files/$ContentFileId"
        $response = Invoke-MgGraphRequest -Method GET -Uri $uri
        
        return $response
    }
    catch {
        Write-Error "Error getting content file status: $_"
        throw
    }
}

# Function to commit a content file
function Commit-AppContentFile {
    param (
        [Parameter(Mandatory=$true)]
        [string]$AppId,
        
        [Parameter(Mandatory=$true)]
        [string]$AppType,
        
        [Parameter(Mandatory=$true)]
        [string]$ContentVersionId,
        
        [Parameter(Mandatory=$true)]
        [string]$ContentFileId,
        
        [Parameter(Mandatory=$true)]
        [hashtable]$CommitData
    )
    
    try {
        $uri = "https://graph.microsoft.com/beta/deviceAppManagement/mobileApps/$AppId/microsoft.graph.$AppType/contentVersions/$ContentVersionId/files/$ContentFileId/commit"
        $response = Invoke-MgGraphRequest -Method POST -Uri $uri -Body ($CommitData | ConvertTo-Json)
        
        return $response
    }
    catch {
        Write-Error "Error committing content file: $_"
        throw
    }
}

# Function to update app with committed content version
function Update-AppWithContentVersion {
    param (
        [Parameter(Mandatory=$true)]
        [string]$AppId,
        
        [Parameter(Mandatory=$true)]
        [string]$AppType,
        
        [Parameter(Mandatory=$true)]
        [string]$ContentVersionId
    )
    
    try {
        $uri = "https://graph.microsoft.com/beta/deviceAppManagement/mobileApps/$AppId"
        $updateData = @{
            "@odata.type"           = "#microsoft.graph.$AppType"
            committedContentVersion = $ContentVersionId
        }
        
        $response = Invoke-MgGraphRequest -Method PATCH -Uri $uri -Body ($updateData | ConvertTo-Json)
        
        return $response
    }
    catch {
        Write-Error "Error updating app with content version: $_"
        throw
    }
}

# Function to update app icon
function Update-AppIcon {
    param (
        [Parameter(Mandatory=$true)]
        [string]$AppId,
        
        [Parameter(Mandatory=$true)]
        [string]$AppType,
        
        [Parameter(Mandatory=$true)]
        [string]$Base64Icon
    )
    
    try {
        $uri = "https://graph.microsoft.com/beta/deviceAppManagement/mobileApps/$AppId"
        $updateData = @{
            "@odata.type" = "#microsoft.graph.$AppType"
            largeIcon     = @{
                "@odata.type" = "#microsoft.graph.mimeContent"
                type          = "image/png"
                value         = $Base64Icon
            }
        }
        
        $response = Invoke-MgGraphRequest -Method PATCH -Uri $uri -Body ($updateData | ConvertTo-Json -Depth 10)
        
        return $response
    }
    catch {
        Write-Error "Error updating app icon: $_"
        throw
    }
}

# Function to encrypt file using AES encryption for Intune upload
# Microsoft Intune uses a specific encryption format for application packages uploaded to the service. 
# This outlines the binary structure and cryptographic algorithms used.
# The encrypted file is structured as follows:
# - HMAC-SHA256 MAC (32 bytes)
# Message Authentication Code for integrity verification
# Computed over the remainder of the file after this field
# - AES Initialization Vector (16 bytes)
# Random IV used for AES-CBC encryption
# - Encrypted Content (variable length)
# The original file encrypted using AES-CBC mode
# Cryptographic Algorithms
# Encryption: AES-256 in CBC mode
# Integrity: HMAC-SHA256
# Keys: Randomly generated for each file. Separate keys for encryption and HMAC operations
function Encrypt-FileForIntune {
    param (
        [Parameter(Mandatory=$true)]
        [string]$SourceFile
    )
    
    function Generate-Key {
        $aesSp = [System.Security.Cryptography.AesCryptoServiceProvider]::new()
        $aesSp.GenerateKey()
        return $aesSp.Key
    }
    
    try {
        $targetFile = "$SourceFile.bin"
        $sha256 = [System.Security.Cryptography.SHA256]::Create()
        $aes = [System.Security.Cryptography.Aes]::Create()
        $aes.Key = Generate-Key
        $hmac = [System.Security.Cryptography.HMACSHA256]::new()
        $hmac.Key = Generate-Key
        $hashLength = $hmac.HashSize / 8
        
        $sourceStream = [System.IO.File]::OpenRead($SourceFile)
        $sourceSha256 = $sha256.ComputeHash($sourceStream)
        $sourceStream.Seek(0, "Begin") | Out-Null
        $targetStream = [System.IO.File]::Open($targetFile, constants.TfOperationCreate)
        
        $targetStream.Write((New-Object byte[] $hashLength), 0, $hashLength)
        $targetStream.Write($aes.IV, 0, $aes.IV.Length)
        $transform = $aes.CreateEncryptor()
        $cryptoStream = [System.Security.Cryptography.CryptoStream]::new($targetStream, $transform, "Write")
        $sourceStream.CopyTo($cryptoStream)
        $cryptoStream.FlushFinalBlock()
        
        $targetStream.Seek($hashLength, "Begin") | Out-Null
        $mac = $hmac.ComputeHash($targetStream)
        $targetStream.Seek(0, "Begin") | Out-Null
        $targetStream.Write($mac, 0, $mac.Length)
        
        $targetStream.Close()
        $cryptoStream.Close()
        $sourceStream.Close()
        
        return [PSCustomObject][ordered]@{
            encryptionKey        = [System.Convert]::ToBase64String($aes.Key)
            fileDigest           = [System.Convert]::ToBase64String($sourceSha256)
            fileDigestAlgorithm  = "SHA256"
            initializationVector = [System.Convert]::ToBase64String($aes.IV)
            mac                  = [System.Convert]::ToBase64String($mac)
            macKey               = [System.Convert]::ToBase64String($hmac.Key)
            profileIdentifier    = "ProfileVersion1"
        }
    }
    catch {
        Write-Error "Error encrypting file: $_"
        throw
    }
}


# Analyze-EncryptedFileHex
# This function reads an encrypted file (for example, one produced by Encrypt-FileForIntune)
# and returns useful insights by analyzing its binary contents as hexadecimal.
# It assumes the file structure is:
#   - Bytes 0 to 31: HMAC-SHA256 MAC,
#   - Bytes 32 to 47: AES-256-CBC Initialization Vector (IV),
#   - Bytes 48 to end: Encrypted content.
#
# The function returns a PSCustomObject containing:
#   - FileLength: total file size in bytes,
#   - HMACHex: hexadecimal representation of the first 32 bytes (MAC),
#   - IVHex: hexadecimal representation of the next 16 bytes (IV),
#   - CiphertextSample: hexadecimal representation of the first 16 bytes of ciphertext (if available),
#   - FullHeaderHex: a hex summary of the first 64 bytes for additional context.
function Analyze-EncryptedFileHex {
  param (
      [Parameter(Mandatory = $true)]
      [string]$EncryptedFilePath
  )
  
  try {
      if (-not (Test-Path $EncryptedFilePath)) {
          throw "Encrypted file not found: $EncryptedFilePath"
      }
      
      $fileBytes = [System.IO.File]::ReadAllBytes($EncryptedFilePath)
      $fileLength = $fileBytes.Length
      
      if ($fileLength -lt 48) {
          throw "File too short to contain valid encryption information. Expected at least 48 bytes, got $fileLength."
      }
      
      # Extract HMAC (first 32 bytes)
      $hmacBytes = $fileBytes[0..31]
      $hmacHex = ($hmacBytes | ForEach-Object { $_.ToString("X2") }) -join ''
      
      # Extract IV (next 16 bytes: bytes 32 to 47)
      $ivBytes = $fileBytes[32..47]
      $ivHex = ($ivBytes | ForEach-Object { $_.ToString("X2") }) -join ''
      
      # If available, extract a sample of the ciphertext (first 16 bytes starting at byte 48)
      if ($fileLength -ge 64) {
          $ciphertextBytes = $fileBytes[48..63]
          $ciphertextSample = ($ciphertextBytes | ForEach-Object { $_.ToString("X2") }) -join ''
      }
      else {
          $ciphertextSample = "N/A"
      }
      
      # Get a full header summary (first 64 bytes, or the entire file if shorter)
      $headerLength = [Math]::Min(64, $fileLength)
      $headerBytes = $fileBytes[0..($headerLength - 1)]
      $fullHeaderHex = ($headerBytes | ForEach-Object { $_.ToString("X2") }) -join ' '
      
      return [PSCustomObject]@{
          FileLength       = $fileLength
          HMACHex          = $hmacHex
          IVHex            = $ivHex
          CiphertextSample = $ciphertextSample
          FullHeaderHex    = $fullHeaderHex
      }
  }
  catch {
      Write-Error "Error analyzing encrypted file: $_"
      throw
  }
}

# Functions to handle Azure Storage uploading
function Upload-FileToAzureStorage {
    param (
        [Parameter(Mandatory=$true)]
        [string]$SasUri,
        
        [Parameter(Mandatory=$true)]
        [string]$FilePath
    )
    
    $blockSize = 8 * 1024 * 1024  # 8 MB block size
    $fileSize = (Get-Item $FilePath).Length
    $totalBlocks = [Math]::Ceiling($fileSize / $blockSize)
    
    $maxRetries = 3
    $retryCount = 0
    $uploadSuccess = $false
    
    while (-not $uploadSuccess -and $retryCount -lt $maxRetries) {
        try {
            $fileStream = [System.IO.File]::OpenRead($FilePath)
            $blockId = 0
            # Initialize block list with proper XML structure
            $blockList = [System.Xml.Linq.XDocument]::Parse(@"
<?xml version="1.0" encoding="utf-8"?>
<BlockList></BlockList>
"@)
            
            # Ensure proper XML namespace
            $blockList.Declaration.Encoding = "utf-8"
            $blockBuffer = [byte[]]::new($blockSize)
            
            Write-Host "`n⬆️  Uploading to Azure Storage..." -ForegroundColor Cyan
            Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
            
            # Show file size with proper formatting
            $fileSizeMB = [Math]::Round($fileSize / 1MB, 2)
            Write-Host "📦 File size: " -NoNewline
            Write-Host "$fileSizeMB MB" -ForegroundColor Yellow
            
            if ($retryCount -gt 0) {
                Write-Host "🔄 Attempt $($retryCount + 1) of $maxRetries" -ForegroundColor Yellow
            }
            Write-Host ""  # Add a blank line before progress bar
            
            while ($bytesRead = $fileStream.Read($blockBuffer, 0, $blockSize)) {
                # Ensure block ID is properly padded and valid base64
                $blockIdBytes = [System.Text.Encoding]::UTF8.GetBytes($blockId.ToString("D6"))
                $id = [System.Convert]::ToBase64String($blockIdBytes)
                $blockList.Root.Add([System.Xml.Linq.XElement]::new("Latest", $id))
                
                $uploadBlockSuccess = $false
                $blockRetries = 3
                while (-not $uploadBlockSuccess -and $blockRetries -gt 0) {
                    try {
                        $blockUri = "$SasUri&comp=block&blockid=$id"
                        try {
                            Invoke-WebRequest -Method Put $blockUri `
                                -Headers @{"x-ms-blob-type" = "BlockBlob" } `
                                -Body ([byte[]]($blockBuffer[0..($bytesRead - 1)])) `
                                -ErrorAction Stop | Out-Null
                                
                            # Block upload successful
                            $uploadBlockSuccess = $true
                        }
                        catch {
                            Write-Host "`nFailed to upload block $blockId" -ForegroundColor Red
                            Write-Host "Error: $_" -ForegroundColor Red
                            throw
                        }
                        $uploadBlockSuccess = $true
                    }
                    catch {
                        $blockRetries--
                        if ($blockRetries -gt 0) {
                            Write-Host "Retrying block upload..." -ForegroundColor Yellow
                            Start-Sleep -Seconds 2
                        }
                        else {
                            Write-Host "Block upload failed: $_" -ForegroundColor Red
                            throw $_
                        }
                    }
                }
                
                $percentComplete = [Math]::Round(($blockId + 1) / $totalBlocks * 100, 1)
                $uploadedMB = [Math]::Min([Math]::Round(($blockId + 1) * $blockSize / 1MB, 1), [Math]::Round($fileSize / 1MB, 1))
                $totalMB = [Math]::Round($fileSize / 1MB, 1)
                
                # Build progress bar
                $progressWidth = 50
                $filledBlocks = [math]::Floor($percentComplete / 2)
                $emptyBlocks = $progressWidth - $filledBlocks
                $progressBar = "[" + ("▓" * $filledBlocks) + ("░" * $emptyBlocks) + "]"
                
                # Clear line and write progress
                [Console]::SetCursorPosition(0, [Console]::CursorTop)
                [Console]::Write((" " * [Console]::WindowWidth))
                [Console]::SetCursorPosition(0, [Console]::CursorTop)
                Write-Host $progressBar -NoNewline
                Write-Host " $percentComplete%" -NoNewline -ForegroundColor Cyan
                Write-Host " ($uploadedMB MB / $totalMB MB)" -NoNewline
                
                $blockId++
            }
            
            Write-Host ""
            
            $fileStream.Close()
            
            Invoke-RestMethod -Method Put "$SasUri&comp=blocklist" -Body $blockList | Out-Null
            $uploadSuccess = $true
        }
        catch {
            $retryCount++
            if ($retryCount -lt $maxRetries) {
                Write-Host "`nUpload failed. Retrying in a few seconds..." -ForegroundColor Yellow
                Start-Sleep -Seconds 5
            }
            else {
                Write-Host "`nFailed to upload file after $maxRetries attempts." -ForegroundColor Red
                Write-Host "Error: $_" -ForegroundColor Red
                throw
            }
        }
        finally {
            if ($fileStream) {
                $fileStream.Close()
            }
        }
    }
    
    Write-Host "✅ Upload completed successfully" -ForegroundColor Green
}

# Function to get app logo for Intune app
function Get-AppLogo {
  param (
      [Parameter(Mandatory=$true)]
      [string]$AppName,
      
      [Parameter(Mandatory=$false)]
      [string]$LocalLogoPath = $null
  )
  
  try {
      $tempLogoPath = $null
      
      if ($LocalLogoPath -and (Test-Path $LocalLogoPath)) {
          # Use the provided local logo file
          $tempLogoPath = $LocalLogoPath
          Write-Host "Using local logo file: $LocalLogoPath" -ForegroundColor Gray
      }
      else {
          Write-Host "⚠️ No valid logo file available" -ForegroundColor Yellow
          return $null
      }
      
      if (-not $tempLogoPath -or -not (Test-Path $tempLogoPath)) {
          Write-Host "⚠️ No valid logo file available" -ForegroundColor Yellow
          return $null
      }
      
      # Convert the logo to base64
      $logoContent = [System.Convert]::ToBase64String([System.IO.File]::ReadAllBytes($tempLogoPath))
      
      # Cleanup temp file if we downloaded it
      if ($tempLogoPath -ne $LocalLogoPath -and (Test-Path $tempLogoPath)) {
          Remove-Item $tempLogoPath -Force
      }
      
      return $logoContent
  }
  catch {
      Write-Host "⚠️ Error processing logo: $_" -ForegroundColor Yellow
      return $null
  }
}


# Script Setup
Import-Module Microsoft.Graph.Authentication

$secureClientSecret = ConvertTo-SecureString -String $ClientSecret -AsPlainText -Force
$clientSecretCredential = New-Object -TypeName System.Management.Automation.PSCredential -ArgumentList $ClientId, $secureClientSecret

Write-Host "Connecting to Microsoft Graph..." -ForegroundColor Cyan
Connect-MgGraph -ClientSecretCredential $clientSecretCredential -TenantId $TenantId

# Add your main script functionality here

# Main function to upload a DMG file to Intune
function Publish-IntunePackage {
    param (
        [Parameter(Mandatory=$true)]
        [string]$DmgFilePath,
        
        [Parameter(Mandatory=$true)]
        [string]$AppDisplayName,
        
        [Parameter(Mandatory=$true)]
        [string]$AppVersion,
        
        [Parameter(Mandatory=$true)]
        [string]$AppBundleId,
        
        [Parameter(Mandatory=$false)]
        [string]$AppDescription,
        
        [Parameter(Mandatory=$false)]
        [string]$AppPublisher,
        
        [Parameter(Mandatory=$false)]
        [string]$LogoFilePath
    )
    
    try {
        # Validate file exists and is DMG
        if (-not (Test-Path $DmgFilePath)) {
            throw "DMG file not found: $DmgFilePath"
        }
        
        if (-not $DmgFilePath.ToLower().EndsWith('.dmg')) {
            throw "File must be a DMG file: $DmgFilePath"
        }
        
        # Set defaults if not provided
        if ([string]::IsNullOrWhiteSpace($AppDescription)) {
            $AppDescription = "$AppDisplayName $AppVersion"
        }
        
        if ([string]::IsNullOrWhiteSpace($AppPublisher)) {
            $AppPublisher = $AppDisplayName
        }
        
        $fileName = [System.IO.Path]::GetFileName($DmgFilePath)
        $appType = "macOSDmgApp"
        
        Write-Host "`n📋 Application Details:" -ForegroundColor Cyan
        Write-Host "   • Display Name: $AppDisplayName" -ForegroundColor Cyan
        Write-Host "   • Version: $AppVersion" -ForegroundColor Cyan
        Write-Host "   • Bundle ID: $AppBundleId" -ForegroundColor Cyan
        Write-Host "   • Publisher: $AppPublisher" -ForegroundColor Cyan
        Write-Host "   • Description: $AppDescription" -ForegroundColor Cyan
        Write-Host "   • File: $fileName" -ForegroundColor Cyan
        
        # Step 1: Create the app in Intune
        Write-Host "`n🔄 Creating macOS app (DMG) in Intune..." -ForegroundColor Yellow
        
        $app = @{
            "@odata.type"                   = "#microsoft.graph.$appType"
            displayName                     = $AppDisplayName
            description                     = $AppDescription
            publisher                       = $AppPublisher
            fileName                        = $fileName
            packageIdentifier               = $AppBundleId
            bundleId                        = $AppBundleId
            versionNumber                   = $AppVersion
            primaryBundleId                 = $AppBundleId
            primaryBundleVersion            = $AppVersion
            minimumSupportedOperatingSystem = @{
                "@odata.type" = "#microsoft.graph.macOSMinimumOperatingSystem"
                v11_0         = $true
            }
            includedApps                    = @(
                @{
                    "@odata.type" = "#microsoft.graph.macOSIncludedApp"
                    bundleId      = $AppBundleId
                    bundleVersion = $AppVersion
                }
            )
        }
        
        $newApp = New-IntuneApp -AppData $app
        $appId = $newApp.id
        Write-Host "✅ App created successfully (ID: $appId)" -ForegroundColor Green
        
        # Step 2: Create content version
        Write-Host "`n🔒 Processing content version..." -ForegroundColor Yellow
        $contentVersion = New-AppContentVersion -AppId $appId -AppType $appType
        $contentVersionId = $contentVersion.id
        Write-Host "✅ Content version created (ID: $contentVersionId)" -ForegroundColor Green
        
        # Step 3: Encrypt the file
        Write-Host "`n🔐 Encrypting application file..." -ForegroundColor Yellow
        $encryptedFilePath = "$DmgFilePath.bin"
        if (Test-Path $encryptedFilePath) {
            Remove-Item $encryptedFilePath -Force
        }
        $fileEncryptionInfo = Encrypt-FileForIntune -SourceFile $DmgFilePath
        Write-Host "✅ Encryption complete" -ForegroundColor Green

        # Analyze the encrypted file and display its hex details
        $analysis = Analyze-EncryptedFileHex -EncryptedFilePath $encryptedFilePath
        Write-Host "`n🔍 Encrypted file analysis:" -ForegroundColor Cyan
        Write-Host "   • File Length: $($analysis.FileLength) bytes" -ForegroundColor Cyan
        Write-Host "   • HMAC (Hex): $($analysis.HMACHex)" -ForegroundColor Cyan
        Write-Host "   • IV (Hex): $($analysis.IVHex)" -ForegroundColor Cyan
        Write-Host "   • Ciphertext Sample (first 16 bytes): $($analysis.CiphertextSample)" -ForegroundColor Cyan
        Write-Host "   • Full Header Hex: $($analysis.FullHeaderHex)" -ForegroundColor Cyan
        
        # Step 4: Create content file
        Write-Host "`n📦 Creating content file..." -ForegroundColor Yellow
        $fileContent = @{
            "@odata.type" = "#microsoft.graph.mobileAppContentFile"
            name          = $fileName
            size          = (Get-Item $DmgFilePath).Length
            sizeEncrypted = (Get-Item "$DmgFilePath.bin").Length
            isDependency  = $false
        }
        
        $contentFile = New-AppContentFile -AppId $appId -AppType $appType -ContentVersionId $contentVersionId -FileData $fileContent
        $contentFileId = $contentFile.id
        
        # Step 5: Wait for Azure Storage Uri
        Write-Host "`n⏳ Waiting for Azure Storage URI..." -ForegroundColor Yellow
        
        $attempts = 0
        $maxAttempts = 30
        $fileStatus = $null
        
        do {
            if ($attempts -gt 0) {
                Write-Host "Waiting for Azure Storage URI... (Attempt $($attempts)/$maxAttempts)" -ForegroundColor Yellow
                Start-Sleep -Seconds 5
            }
            $fileStatus = Get-AppContentFileStatus -AppId $appId -AppType $appType -ContentVersionId $contentVersionId -ContentFileId $contentFileId
            $attempts++
        } while ($fileStatus.uploadState -ne "azureStorageUriRequestSuccess" -and $attempts -lt $maxAttempts)
        
        if ($fileStatus.uploadState -ne "azureStorageUriRequestSuccess") {
            throw "Failed to get Azure Storage URI after $maxAttempts attempts."
        }
        
        Write-Host "✅ Azure Storage URI received" -ForegroundColor Green
        
        # Step 6: Upload file to Azure Storage
        Upload-FileToAzureStorage -SasUri $fileStatus.azureStorageUri -FilePath "$DmgFilePath.bin"
        
        # Step 7: Commit the file
        Write-Host "`n🔄 Committing file..." -ForegroundColor Yellow
        $commitData = @{
            fileEncryptionInfo = $fileEncryptionInfo
        }
        
        Commit-AppContentFile -AppId $appId -AppType $appType -ContentVersionId $contentVersionId -ContentFileId $contentFileId -CommitData $commitData
        
        # Step 8: Wait for commit to complete
        Write-Host "`n⏳ Waiting for file commitment to complete..." -ForegroundColor Yellow
        $retryCount = 0
        $maxRetries = 10
        
        do {
            Start-Sleep -Seconds 10
            $fileStatus = Get-AppContentFileStatus -AppId $appId -AppType $appType -ContentVersionId $contentVersionId -ContentFileId $contentFileId
            
            if ($fileStatus.uploadState -eq "commitFileFailed") {
                $retryCount++
                Write-Host "Commit failed, retrying ($retryCount/$maxRetries)..." -ForegroundColor Yellow
                Commit-AppContentFile -AppId $appId -AppType $appType -ContentVersionId $contentVersionId -ContentFileId $contentFileId -CommitData $commitData
            }
            elseif ($fileStatus.uploadState -eq "commitFileSuccess") {
                Write-Host "✅ File committed successfully" -ForegroundColor Green
                break
            }
            else {
                Write-Host "Current state: $($fileStatus.uploadState). Waiting..." -ForegroundColor Yellow
            }
        } while ($retryCount -lt $maxRetries)
        
        if ($fileStatus.uploadState -ne "commitFileSuccess") {
            throw "Failed to commit file after $maxRetries attempts."
        }
        
        # Step 9: Update app with committed content version
        Write-Host "`n🔄 Updating app with committed content..." -ForegroundColor Yellow
        Update-AppWithContentVersion -AppId $appId -AppType $appType -ContentVersionId $contentVersionId
        Write-Host "✅ App updated successfully" -ForegroundColor Green
        
        # Step 10: Add logo if one was provided
        if ($LogoFilePath -and (Test-Path $LogoFilePath)) {
            Write-Host "`n🖼️  Adding app logo..." -ForegroundColor Yellow
            $logoContent = Get-AppLogo -AppName $AppDisplayName -LocalLogoPath $LogoFilePath
            if ($logoContent) {
                Update-AppIcon -AppId $appId -AppType $appType -Base64Icon $logoContent
                Write-Host "✅ Logo added successfully" -ForegroundColor Green
            }
        }
        
        # Step 11: Clean up temporary files
        Write-Host "`n🧹 Cleaning up temporary files..." -ForegroundColor Yellow
        if (Test-Path "$DmgFilePath.bin") {
            Remove-Item "$DmgFilePath.bin" -Force
        }
        Write-Host "✅ Cleanup complete" -ForegroundColor Green
        
        # Step 12: Final success message
        Write-Host "`n✨ Successfully uploaded $AppDisplayName to Intune" -ForegroundColor Cyan
        Write-Host "🔗 Intune Portal URL: https://intune.microsoft.com/#view/Microsoft_Intune_Apps/SettingsMenu/~/0/appId/$appId" -ForegroundColor Cyan
        
        return $appId
    }
    catch {
        Write-Host "❌ Error publishing package to Intune: $_" -ForegroundColor Red
        throw
    }
}

# Call the Publish-IntunePackage function with the parameters provided
try {
    Write-Host "`n📦 Starting DMG upload process..." -ForegroundColor Cyan
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    
    Publish-IntunePackage -DMGFilePath $DmgFilePath `
                          -AppDisplayName $AppDisplayName `
                          -AppVersion $AppVersion `
                          -AppBundleId $AppBundleId `
                          -AppDescription $AppDescription `
                          -AppPublisher $AppPublisher `
                          -LogoFilePath $LogoFilePath

    Write-Host "`n🎉 DMG upload process completed successfully!" -ForegroundColor Green
}
catch {
    Write-Host "`n❌ DMG upload process failed: $_" -ForegroundColor Red
    Disconnect-MgGraph > $null 2>&1
    exit 1
}

# Disconnect from Microsoft Graph
Write-Host "`nDisconnecting from Microsoft Graph..." -ForegroundColor Cyan
Disconnect-MgGraph > $null 2>&1
Write-Host "✅ Disconnected from Microsoft Graph" -ForegroundColor Green