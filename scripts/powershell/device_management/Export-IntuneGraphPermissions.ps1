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
    [string]$ClientSecret
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

# Script Setup
Import-Module Microsoft.Graph.Authentication

$secureClientSecret = ConvertTo-SecureString -String $ClientSecret -AsPlainText -Force
$clientSecretCredential = New-Object -TypeName System.Management.Automation.PSCredential -ArgumentList $ClientId, $secureClientSecret

Write-Host "Connecting to Microsoft Graph..."
Connect-MgGraph -ClientSecretCredential $clientSecretCredential -TenantId $TenantId

$scriptPath = $PSScriptRoot
if (-not $scriptPath) {
    $scriptPath = Split-Path -Parent -Path $MyInvocation.MyCommand.Definition
}
if (-not $scriptPath) {
    $scriptPath = Get-Location
}

$exportFolder = Join-Path -Path $scriptPath -ChildPath "Export"
if (-not (Test-Path -Path $exportFolder)) {
    New-Item -ItemType Directory -Path $exportFolder | Out-Null
}

$outputFile = Join-Path -Path $exportFolder -ChildPath "ExportedGraphPermissions.json"

Write-Host "Script started: $(Get-Date)"
Write-Host "Retrieving permissions..."
$allPermissions = Find-MgGraphPermission
Write-Host "Retrieved $($allPermissions.Count) permissions"

Write-Host "Retrieving commands..."
$allCommands = @()
$allCommands += Find-MgGraphCommand -Command "*" -ApiVersion "v1.0"
$allCommands += Find-MgGraphCommand -Command "*" -ApiVersion "beta"
Write-Host "Retrieved $($allCommands.Count) commands"

$topLevelGroups = [ordered]@{}

Write-Host "Processing commands and permissions..."
foreach ($command in $allCommands) {
    $uri = $command.URI
    $method = $command.Method
    $apiVersion = $command.ApiVersion

    $topLevelPath = '/' + ($uri -split '/')[1]
    $referenceId = Get-PathBasedReferenceId -Path $uri
    
    if ($topLevelPath -notin $topLevelGroups.Keys) {
        $topLevelGroups[$topLevelPath] = [ordered]@{
            APIResourceCount = 0
            APIResources = [ordered]@{}
        }
    }
    
    if ($uri -notin $topLevelGroups[$topLevelPath].APIResources.Keys) {
        $topLevelGroups[$topLevelPath].APIResources[$uri] = [ordered]@{
            ReferenceId = $referenceId
        }
        $topLevelGroups[$topLevelPath].APIResourceCount++
    }

    if ($method -notin $topLevelGroups[$topLevelPath].APIResources[$uri].Keys) {
        $topLevelGroups[$topLevelPath].APIResources[$uri][$method] = [ordered]@{
            ApiVersions = @($apiVersion)
            Permissions = [ordered]@{
                Read = @()
                ReadWrite = @()
            }
        }
    } else {
        if ($apiVersion -notin $topLevelGroups[$topLevelPath].APIResources[$uri][$method].ApiVersions) {
            $topLevelGroups[$topLevelPath].APIResources[$uri][$method].ApiVersions += $apiVersion
        }
    }

    foreach ($permission in $command.Permissions) {
        Write-Verbose "Processing $method permission for $uri : $($permission.Name)"
        $permissionType = if ($permission.Name -match "ReadWrite") { "ReadWrite" } else { constants.TfOperationRead }
        $permissionName = $permission.Name

        if ($permissionName -notin $topLevelGroups[$topLevelPath].APIResources[$uri][$method].Permissions[$permissionType]) {
            $topLevelGroups[$topLevelPath].APIResources[$uri][$method].Permissions[$permissionType] += $permissionName
            Write-Verbose "Added permission: $permissionName"
        }
    }
}

Write-Host "Sorting data..."
$sortedTopLevelGroups = [ordered]@{}
$topLevelGroups.GetEnumerator() | Sort-Object Name | ForEach-Object {
    $sortedAPIResources = [ordered]@{}
    $_.Value.APIResources.GetEnumerator() | Sort-Object Name | ForEach-Object {
        $sortedAPIResources[$_.Name] = $_.Value
    }
    $sortedTopLevelGroups[$_.Key] = @{
        APIResourceCount = $_.Value.APIResourceCount
        APIResources = $sortedAPIResources
    }
}

Write-Host "Exporting data to $outputFile"
$sortedTopLevelGroups | ConvertTo-Json -Depth 6 | Out-File -FilePath $outputFile

Write-Host "Script completed: $(Get-Date)"

Disconnect-MgGraph
Write-Host "Disconnected from Microsoft Graph"
