[CmdletBinding()]
param (
    [Parameter(Mandatory=$false,
    HelpMessage="Specific Package Identifier (if not provided, will search based on SearchTerm)")]
    [string]$PackageId,
    
    [Parameter(Mandatory=$false,
    HelpMessage="Search term to find packages (if PackageId is not provided)")]
    [string]$SearchTerm,
    
    [Parameter(Mandatory=$true,
    HelpMessage="Export results to JSON file")]
    [bool]$ExportToJson
)

# Function to get package manifests
function Get-MSStorePackageManifests {
    param (
        [Parameter(Mandatory=$false)]
        [string]$SpecificPackageId,
        
        [Parameter(Mandatory=$false)]
        [string]$SearchKeyword
    )
    
    try {
        if ($SpecificPackageId) {
            # GET specific package manifest by ID
            $uri = "https://storeedgefd.dsx.mp.microsoft.com/v9.0/packageManifests/$SpecificPackageId"
            Write-Host "🔍 Getting specific package manifest..." -ForegroundColor Cyan
            Write-Host "   Package ID: $SpecificPackageId" -ForegroundColor Gray
        } elseif ($SearchKeyword) {
            # Search for packages first, then get manifests
            Write-Host "🔍 Searching for packages matching: $SearchKeyword..." -ForegroundColor Cyan
            $searchResults = Search-MSStorePackages -SearchTerm $SearchKeyword
            
            if (-not $searchResults -or $searchResults.Count -eq 0) {
                Write-Host "⚠️ No packages found matching search term: $SearchKeyword" -ForegroundColor Yellow
                return $null
            }
            
            Write-Host "📦 Found $($searchResults.Count) package(s), retrieving manifests..." -ForegroundColor Cyan
            
            $manifests = @()
            foreach ($package in $searchResults) {
                try {
                    $manifestUri = "https://storeedgefd.dsx.mp.microsoft.com/v9.0/packageManifests/$($package.PackageIdentifier)"
                    Write-Host "   Getting manifest for: $($package.PackageName)" -ForegroundColor Gray
                    
                    $response = Invoke-RestMethodWithRetry -Uri $manifestUri -Method GET
                    if ($response -and $response.Data) {
                        $manifests += $response.Data
                    }
                } catch {
                    Write-Host "   ⚠️ Failed to get manifest for $($package.PackageName): $_" -ForegroundColor Yellow
                }
            }
            
            return $manifests
        } else {
            throw "Either PackageId or SearchKeyword must be provided"
        }
        
        Write-Host "   Endpoint: $uri" -ForegroundColor Gray
        Write-Host ""
        
        $response = Invoke-RestMethodWithRetry -Uri $uri -Method GET
        
        return $response.Data
    }
    catch {
        Write-Host "❌ Error getting package manifests: $_" -ForegroundColor Red
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

# Function to search for packages
function Search-MSStorePackages {
    param(
        [Parameter(Mandatory=$true)]
        [string]$SearchTerm
    )
    
    try {
        $storeSearchUrl = "https://storeedgefd.dsx.mp.microsoft.com/v9.0/manifestSearch"
        $requestBody = @{
            Query = @{
                KeyWord   = $SearchTerm
                MatchType = "Substring"
            }
        } | ConvertTo-Json
        
        $response = Invoke-RestMethodWithRetry -Uri $storeSearchUrl -Method POST -Body $requestBody -ContentType 'application/json'
        
        if ($response.Data) {
            return $response.Data
        } else {
            return @()
        }
    }
    catch {
        Write-Host "❌ Error searching packages: $_" -ForegroundColor Red
        throw
    }
}

# Function with retry logic for API calls
function Invoke-RestMethodWithRetry {
    param(
        [Parameter(Mandatory=$true)]
        [string]$Uri,
        
        [Parameter(Mandatory=$false)]
        [string]$Method = "GET",
        
        [Parameter(Mandatory=$false)]
        [string]$Body = $null,
        
        [Parameter(Mandatory=$false)]
        [string]$ContentType = "application/json"
    )
    
    $retryCount = 0
    $maxRetries = 3
    
    while ($retryCount -le $maxRetries) {
        try {
            if ($Method -eq "POST" -and $Body) {
                $response = Invoke-RestMethod -Uri $Uri -Method $Method -Body $Body -ContentType $ContentType -ErrorAction Stop
            } else {
                $response = Invoke-RestMethod -Uri $Uri -Method $Method -ErrorAction Stop
            }
            
            return $response
        }
        catch {
            $retryCount++
            if ($retryCount -gt $maxRetries) {
                throw "Max retry attempts reached for $Method request to $Uri. Error: $_"
            }
            
            $retryDelay = $retryCount * 2
            Write-Host "   ⚠️ Request attempt $retryCount failed, retrying in $retryDelay seconds..." -ForegroundColor Yellow
            Start-Sleep -Seconds $retryDelay
        }
    }
}

# Function to export manifests to JSON
function Export-ManifestsToJson {
    param (
        [Parameter(Mandatory=$true)]
        $Manifests,
        
        [Parameter(Mandatory=$false)]
        [string]$SpecificPackageId,
        
        [Parameter(Mandatory=$false)]
        [string]$SearchTerm
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
        
        if ($SpecificPackageId) {
            # Export single package manifest
            $packageName = if ($Manifests.PackageIdentifier) { 
                $Manifests.PackageIdentifier 
            } else { 
                $SpecificPackageId 
            }
            $cleanPackageName = $packageName -replace '[\\\/\:\*\?\"\<\>\|]', '_'
            $fileName = "PackageManifest_${cleanPackageName}_${timestamp}.json"
            $filePath = Join-Path -Path $outputDir -ChildPath $fileName
            
            $Manifests | ConvertTo-Json -Depth 20 | Out-File -FilePath $filePath -Encoding UTF8
            
            Write-Host "💾 Exported package manifest to: $filePath" -ForegroundColor Green
        } else {
            # Export multiple package manifests
            $searchTermClean = if ($SearchTerm) { 
                ($SearchTerm -replace '[\\\/\:\*\?\"\<\>\|]', '_') 
            } else { 
                "AllPackages" 
            }
            $fileName = "PackageManifests_${searchTermClean}_${timestamp}.json"
            $filePath = Join-Path -Path $outputDir -ChildPath $fileName
            
            $Manifests | ConvertTo-Json -Depth 20 | Out-File -FilePath $filePath -Encoding UTF8
            
            Write-Host "💾 Exported package manifests to: $filePath" -ForegroundColor Green
        }
        
        return $filePath
    }
    catch {
        Write-Host "❌ Error exporting manifests to JSON: $_" -ForegroundColor Red
        return $null
    }
}

# Function to display package manifest details
function Show-PackageManifestDetails {
    param (
        [Parameter(Mandatory=$true)]
        $Manifest
    )
    
    Write-Host "📦 Microsoft Store Package Manifest Details:" -ForegroundColor Cyan
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    
    if ($Manifest.PackageIdentifier) {
        Write-Host "   • Package Identifier: $($Manifest.PackageIdentifier)" -ForegroundColor Green
    }
    
    if ($Manifest.Versions -and $Manifest.Versions.Count -gt 0) {
        Write-Host "   • Total Versions: $($Manifest.Versions.Count)" -ForegroundColor Green
        
        # Show latest version details
        $latestVersion = $Manifest.Versions | Select-Object -Last 1
        if ($latestVersion) {
            Write-Host "   • Latest Version Details:" -ForegroundColor Green
            
            if ($latestVersion.PackageVersion) {
                Write-Host "     - Version: $($latestVersion.PackageVersion)" -ForegroundColor Yellow
            }
            
            # Default Locale Information
            if ($latestVersion.DefaultLocale) {
                $locale = $latestVersion.DefaultLocale
                Write-Host "     - Default Locale Information:" -ForegroundColor Yellow
                
                if ($locale.PackageLocale) {
                    Write-Host "       · Locale: $($locale.PackageLocale)" -ForegroundColor White
                }
                
                if ($locale.Publisher) {
                    Write-Host "       · Publisher: $($locale.Publisher)" -ForegroundColor White
                }
                
                if ($locale.PublisherUrl) {
                    Write-Host "       · Publisher URL: $($locale.PublisherUrl)" -ForegroundColor White
                }
                
                if ($locale.PublisherSupportUrl) {
                    Write-Host "       · Publisher Support URL: $($locale.PublisherSupportUrl)" -ForegroundColor White
                }
                
                if ($locale.PrivacyUrl) {
                    Write-Host "       · Privacy URL: $($locale.PrivacyUrl)" -ForegroundColor White
                }
                
                if ($locale.Author) {
                    Write-Host "       · Author: $($locale.Author)" -ForegroundColor White
                }
                
                if ($locale.PackageName) {
                    Write-Host "       · Package Name: $($locale.PackageName)" -ForegroundColor White
                }
                
                if ($locale.PackageUrl) {
                    Write-Host "       · Package URL: $($locale.PackageUrl)" -ForegroundColor White
                }
                
                if ($locale.License) {
                    Write-Host "       · License: $($locale.License)" -ForegroundColor White
                }
                
                if ($locale.LicenseUrl) {
                    Write-Host "       · License URL: $($locale.LicenseUrl)" -ForegroundColor White
                }
                
                if ($locale.Copyright) {
                    Write-Host "       · Copyright: $($locale.Copyright)" -ForegroundColor White
                }
                
                if ($locale.CopyrightUrl) {
                    Write-Host "       · Copyright URL: $($locale.CopyrightUrl)" -ForegroundColor White
                }
                
                if ($locale.ShortDescription) {
                    $shortDesc = $locale.ShortDescription
                    if ($shortDesc.Length -gt 100) {
                        $shortDesc = $shortDesc.Substring(0, 100) + "..."
                    }
                    Write-Host "       · Short Description: $shortDesc" -ForegroundColor White
                }
                
                if ($locale.Description) {
                    $desc = $locale.Description
                    if ($desc.Length -gt 150) {
                        $desc = $desc.Substring(0, 150) + "..."
                    }
                    Write-Host "       · Description: $desc" -ForegroundColor White
                }
                
                if ($locale.Moniker) {
                    Write-Host "       · Moniker: $($locale.Moniker)" -ForegroundColor White
                }
                
                if ($locale.Tags) {
                    Write-Host "       · Tags: $($locale.Tags -join ', ')" -ForegroundColor White
                }
                
                if ($locale.Agreements -and $locale.Agreements.Count -gt 0) {
                    Write-Host "       · Agreements: $($locale.Agreements.Count) agreement(s)" -ForegroundColor White
                }
                
                if ($locale.ReleaseNotes) {
                    $releaseNotes = $locale.ReleaseNotes
                    if ($releaseNotes.Length -gt 100) {
                        $releaseNotes = $releaseNotes.Substring(0, 100) + "..."
                    }
                    Write-Host "       · Release Notes: $releaseNotes" -ForegroundColor White
                }
                
                if ($locale.ReleaseNotesUrl) {
                    Write-Host "       · Release Notes URL: $($locale.ReleaseNotesUrl)" -ForegroundColor White
                }
                
                if ($locale.PurchaseUrl) {
                    Write-Host "       · Purchase URL: $($locale.PurchaseUrl)" -ForegroundColor White
                }
                
                if ($locale.InstallationNotes) {
                    Write-Host "       · Installation Notes: $($locale.InstallationNotes)" -ForegroundColor White
                }
                
                if ($locale.Documentations -and $locale.Documentations.Count -gt 0) {
                    Write-Host "       · Documentations: $($locale.Documentations.Count) documentation(s)" -ForegroundColor White
                }
            }
            
            # Installers Information
            if ($latestVersion.Installers -and $latestVersion.Installers.Count -gt 0) {
                Write-Host "     - Installers: $($latestVersion.Installers.Count) installer(s)" -ForegroundColor Yellow
                
                foreach ($installer in $latestVersion.Installers) {
                    Write-Host "       · Installer Details:" -ForegroundColor White
                    
                    if ($installer.InstallerIdentifier) {
                        Write-Host "         - ID: $($installer.InstallerIdentifier)" -ForegroundColor Cyan
                    }
                    
                    if ($installer.Architecture) {
                        Write-Host "         - Architecture: $($installer.Architecture)" -ForegroundColor Cyan
                    }
                    
                    if ($installer.InstallerType) {
                        Write-Host "         - Type: $($installer.InstallerType)" -ForegroundColor Cyan
                    }
                    
                    if ($installer.Scope) {
                        Write-Host "         - Scope: $($installer.Scope)" -ForegroundColor Cyan
                    }
                    
                    if ($installer.InstallerUrl) {
                        Write-Host "         - URL: $($installer.InstallerUrl)" -ForegroundColor Cyan
                    }
                    
                    if ($installer.InstallerSha256) {
                        Write-Host "         - SHA256: $($installer.InstallerSha256)" -ForegroundColor Cyan
                    }
                    
                    if ($installer.SignatureSha256) {
                        Write-Host "         - Signature SHA256: $($installer.SignatureSha256)" -ForegroundColor Cyan
                    }
                    
                    if ($installer.InstallModes) {
                        Write-Host "         - Install Modes: $($installer.InstallModes -join ', ')" -ForegroundColor Cyan
                    }
                    
                    if ($installer.InstallerSwitches) {
                        Write-Host "         - Installer Switches:" -ForegroundColor Cyan
                        $switches = $installer.InstallerSwitches
                        if ($switches.Silent) {
                            Write-Host "           · Silent: $($switches.Silent)" -ForegroundColor DarkCyan
                        }
                        if ($switches.SilentWithProgress) {
                            Write-Host "           · Silent with Progress: $($switches.SilentWithProgress)" -ForegroundColor DarkCyan
                        }
                        if ($switches.Interactive) {
                            Write-Host "           · Interactive: $($switches.Interactive)" -ForegroundColor DarkCyan
                        }
                        if ($switches.InstallLocation) {
                            Write-Host "           · Install Location: $($switches.InstallLocation)" -ForegroundColor DarkCyan
                        }
                        if ($switches.Log) {
                            Write-Host "           · Log: $($switches.Log)" -ForegroundColor DarkCyan
                        }
                        if ($switches.Upgrade) {
                            Write-Host "           · Upgrade: $($switches.Upgrade)" -ForegroundColor DarkCyan
                        }
                        if ($switches.Custom) {
                            Write-Host "           · Custom: $($switches.Custom)" -ForegroundColor DarkCyan
                        }
                    }
                    
                    if ($installer.InstallerSuccessCodes) {
                        Write-Host "         - Success Codes: $($installer.InstallerSuccessCodes -join ', ')" -ForegroundColor Cyan
                    }
                    
                    if ($installer.UpgradeBehavior) {
                        Write-Host "         - Upgrade Behavior: $($installer.UpgradeBehavior)" -ForegroundColor Cyan
                    }
                    
                    if ($installer.Commands) {
                        Write-Host "         - Commands: $($installer.Commands -join ', ')" -ForegroundColor Cyan
                    }
                    
                    if ($installer.Protocols) {
                        Write-Host "         - Protocols: $($installer.Protocols -join ', ')" -ForegroundColor Cyan
                    }
                    
                    if ($installer.FileExtensions) {
                        Write-Host "         - File Extensions: $($installer.FileExtensions -join ', ')" -ForegroundColor Cyan
                    }
                    
                    if ($installer.Dependencies) {
                        Write-Host "         - Dependencies:" -ForegroundColor Cyan
                        $deps = $installer.Dependencies
                        if ($deps.WindowsFeatures) {
                            Write-Host "           · Windows Features: $($deps.WindowsFeatures -join ', ')" -ForegroundColor DarkCyan
                        }
                        if ($deps.WindowsLibraries) {
                            Write-Host "           · Windows Libraries: $($deps.WindowsLibraries -join ', ')" -ForegroundColor DarkCyan
                        }
                        if ($deps.PackageDependencies -and $deps.PackageDependencies.Count -gt 0) {
                            Write-Host "           · Package Dependencies: $($deps.PackageDependencies.Count) dependency(ies)" -ForegroundColor DarkCyan
                        }
                        if ($deps.ExternalDependencies) {
                            Write-Host "           · External Dependencies: $($deps.ExternalDependencies -join ', ')" -ForegroundColor DarkCyan
                        }
                    }
                    
                    if ($installer.PackageFamilyName) {
                        Write-Host "         - Package Family Name: $($installer.PackageFamilyName)" -ForegroundColor Cyan
                    }
                    
                    if ($installer.ProductCode) {
                        Write-Host "         - Product Code: $($installer.ProductCode)" -ForegroundColor Cyan
                    }
                    
                    if ($installer.Capabilities) {
                        Write-Host "         - Capabilities: $($installer.Capabilities -join ', ')" -ForegroundColor Cyan
                    }
                    
                    if ($installer.RestrictedCapabilities) {
                        Write-Host "         - Restricted Capabilities: $($installer.RestrictedCapabilities -join ', ')" -ForegroundColor Cyan
                    }
                    
                    if ($installer.Markets) {
                        Write-Host "         - Markets:" -ForegroundColor Cyan
                        $markets = $installer.Markets
                        if ($markets.AllowedMarkets) {
                            Write-Host "           · Allowed: $($markets.AllowedMarkets -join ', ')" -ForegroundColor DarkCyan
                        }
                        if ($markets.ExcludedMarkets) {
                            Write-Host "           · Excluded: $($markets.ExcludedMarkets -join ', ')" -ForegroundColor DarkCyan
                        }
                    }
                    
                    if ($installer.InstallerAbortsTerminal -ne $null) {
                        Write-Host "         - Aborts Terminal: $($installer.InstallerAbortsTerminal)" -ForegroundColor Cyan
                    }
                    
                    if ($installer.ReleaseDate) {
                        Write-Host "         - Release Date: $($installer.ReleaseDate)" -ForegroundColor Cyan
                    }
                    
                    if ($installer.InstallLocationRequired -ne $null) {
                        Write-Host "         - Install Location Required: $($installer.InstallLocationRequired)" -ForegroundColor Cyan
                    }
                    
                    if ($installer.RequireExplicitUpgrade -ne $null) {
                        Write-Host "         - Require Explicit Upgrade: $($installer.RequireExplicitUpgrade)" -ForegroundColor Cyan
                    }
                    
                    if ($installer.DisplayInstallWarnings -ne $null) {
                        Write-Host "         - Display Install Warnings: $($installer.DisplayInstallWarnings)" -ForegroundColor Cyan
                    }
                    
                    if ($installer.UnsupportedOSArchitectures) {
                        Write-Host "         - Unsupported OS Architectures: $($installer.UnsupportedOSArchitectures -join ', ')" -ForegroundColor Cyan
                    }
                    
                    if ($installer.UnsupportedArguments) {
                        Write-Host "         - Unsupported Arguments: $($installer.UnsupportedArguments -join ', ')" -ForegroundColor Cyan
                    }
                    
                    if ($installer.NestedInstallerType) {
                        Write-Host "         - Nested Installer Type: $($installer.NestedInstallerType)" -ForegroundColor Cyan
                    }
                    
                    if ($installer.NestedInstallerFiles -and $installer.NestedInstallerFiles.Count -gt 0) {
                        Write-Host "         - Nested Installer Files: $($installer.NestedInstallerFiles.Count) file(s)" -ForegroundColor Cyan
                    }
                    
                    Write-Host "" # Add spacing between installers
                }
            }
            
            # Localization Information
            if ($latestVersion.Locales -and $latestVersion.Locales.Count -gt 0) {
                Write-Host "     - Additional Locales: $($latestVersion.Locales.Count) locale(s)" -ForegroundColor Yellow
                foreach ($locale in $latestVersion.Locales) {
                    if ($locale.PackageLocale) {
                        Write-Host "       · $($locale.PackageLocale)" -ForegroundColor White
                    }
                }
            }
        }
    }
    
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    Write-Host ""
}

# Script Setup
try {
    Write-Host "📦 Microsoft Store Package Manifest Explorer" -ForegroundColor Cyan
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    
    # Validate parameters
    if (-not $PackageId -and -not $SearchTerm) {
        throw "Either PackageId or SearchTerm parameter must be provided"
    }
    
    if ($PackageId -and $SearchTerm) {
        Write-Host "⚠️ Both PackageId and SearchTerm provided. Using PackageId only." -ForegroundColor Yellow
        $SearchTerm = $null
    }
    
    # Get the package manifests
    $manifests = Get-MSStorePackageManifests -SpecificPackageId $PackageId -SearchKeyword $SearchTerm
    
    if (-not $manifests) {
        Write-Host "📊 No package manifests found" -ForegroundColor Yellow
        exit 0
    }
    
    # Export to JSON if requested
    if ($ExportToJson) {
        $jsonPath = Export-ManifestsToJson -Manifests $manifests -SpecificPackageId $PackageId -SearchTerm $SearchTerm
    }
    
    if ($PackageId) {
        # Display single package manifest
        Show-PackageManifestDetails -Manifest $manifests
    } else {
        # Display multiple package manifests
        if ($manifests -is [array] -and $manifests.Count -gt 0) {
            Write-Host "📊 Found $($manifests.Count) package manifest(s)" -ForegroundColor Green
            Write-Host ""
            
            for ($i = 0; $i -lt $manifests.Count; $i++) {
                Write-Host "Package $($i + 1):" -ForegroundColor Magenta
                Show-PackageManifestDetails -Manifest $manifests[$i]
            }
        } elseif ($manifests -and -not ($manifests -is [array])) {
            # Single manifest returned (not in a collection)
            Write-Host "📊 Found 1 package manifest" -ForegroundColor Green
            Write-Host ""
            Show-PackageManifestDetails -Manifest $manifests
        } else {
            Write-Host "📊 No package manifests found" -ForegroundColor Yellow
        }
    }
    
    Write-Host "🎉 Operation completed successfully!" -ForegroundColor Green
}
catch {
    Write-Host "❌ Script execution failed: $_" -ForegroundColor Red
    exit 1
}