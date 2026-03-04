<#
.SYNOPSIS
    Discovers required Microsoft Graph API permissions for a specific endpoint using raw REST API calls only.

.DESCRIPTION
    Determines the exact APPLICATION permissions needed for any Microsoft Graph API endpoint by
    fetching the Microsoft Graph PowerShell SDK command metadata directly from GitHub and parsing it
    locally — no PowerShell modules required.

    Optionally compares these requirements against an existing enterprise application's granted
    permissions using raw Microsoft Graph REST API calls (OAuth 2.0 client credentials grant).

    The script provides two modes of operation:
    1. Discovery Mode:  Shows what permissions are required for a given endpoint.
    2. Evaluation Mode: Compares an app's current permissions against the requirements and
                        reports any gaps with remediation guidance.

    Permission metadata is cached to $env:TEMP for 24 hours to avoid repeated GitHub downloads.
    Use -RefreshMetadata to force a fresh download.

.PARAMETER Uri
    The Graph API URI to check (e.g., "/groups", "/groups/{id}", "/users/{id}/memberOf").
    Use {id} or any {placeholder} for resource identifiers.

.PARAMETER Method
    The HTTP method to filter on (GET, POST, PATCH, PUT, DELETE).
    If omitted, all methods for the URI are returned.

.PARAMETER ApiVersion
    The API version to check against. Accepts "v1.0" or "beta". Defaults to "beta".

.PARAMETER AppIdToBeEvaluated
    The Application (Client) ID of the enterprise application to evaluate against the
    endpoint requirements. When provided, TenantId and ClientSecret are also required.

.PARAMETER TenantId
    The Entra ID tenant ID. Required when AppIdToBeEvaluated is provided.

.PARAMETER ClientSecret
    The client secret for the application being evaluated. Required when
    AppIdToBeEvaluated is provided.

.PARAMETER RefreshMetadata
    Forces a fresh download of the Graph SDK command metadata from GitHub,
    ignoring any locally cached copy.

.EXAMPLE
    # Discovery: show all permissions needed to read a group
    .\Get-GraphEndpointPermissionsRaw.ps1 -Uri "/groups/{id}" -Method GET

.EXAMPLE
    # Discovery: show permissions for all methods on an endpoint
    .\Get-GraphEndpointPermissionsRaw.ps1 -Uri "/applications/{id}"

.EXAMPLE
    # Evaluation: compare an app's permissions against what POST /groups requires
    .\Get-GraphEndpointPermissionsRaw.ps1 -Uri "/groups" -Method POST `
        -AppIdToBeEvaluated "00000000-0000-0000-0000-000000000001" `
        -TenantId "00000000-0000-0000-0000-000000000002" `
        -ClientSecret "xxx~xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"

.EXAMPLE
    # Force a metadata refresh before checking permissions
    .\Get-GraphEndpointPermissionsRaw.ps1 -Uri "/servicePrincipals/{id}" -Method PATCH `
        -RefreshMetadata

.NOTES
    Author:     Deployment Theory
    Version:    2.0
    Reference:  https://learn.microsoft.com/en-us/graph/permissions-reference

    No external PowerShell modules required. All Graph API interactions use
    Invoke-RestMethod with raw OAuth 2.0 client credentials tokens.

    Metadata source:
    https://github.com/microsoftgraph/msgraph-sdk-powershell
#>

[CmdletBinding()]
param (
    [Parameter(Mandatory = $true,
        HelpMessage = "The Graph API URI to check (e.g., '/groups', '/groups/{id}')")]
    [ValidateNotNullOrEmpty()]
    [string]$Uri,

    [Parameter(Mandatory = $false,
        HelpMessage = "The HTTP method to filter on (GET, POST, PATCH, PUT, DELETE)")]
    [ValidateSet("GET", "POST", "PATCH", "PUT", "DELETE", "")]
    [string]$Method = "",

    [Parameter(Mandatory = $false,
        HelpMessage = "The API version to check against (v1.0 or beta)")]
    [ValidateSet("v1.0", "beta")]
    [string]$ApiVersion = "beta",

    [Parameter(Mandatory = $false,
        HelpMessage = "The Application (Client) ID of the enterprise app to evaluate")]
    [string]$AppIdToBeEvaluated,

    [Parameter(Mandatory = $false,
        HelpMessage = "The Entra ID Tenant ID (required when AppIdToBeEvaluated is provided)")]
    [string]$TenantId,

    [Parameter(Mandatory = $false,
        HelpMessage = "The client secret for the app being evaluated (required when AppIdToBeEvaluated is provided)")]
    [string]$ClientSecret,

    [Parameter(Mandatory = $false,
        HelpMessage = "Force a fresh metadata download from GitHub, bypassing the local cache")]
    [switch]$RefreshMetadata
)

# =============================================================================
# LOGGING HELPERS
# =============================================================================

function Write-LogHeader {
    param (
        [Parameter(Mandatory = $true)]
        [string]$Message
    )
    Write-Host "`n========================================" -ForegroundColor Cyan
    Write-Host " $Message" -ForegroundColor Cyan
    Write-Host "========================================" -ForegroundColor Cyan
}

function Write-LogSection {
    param (
        [Parameter(Mandatory = $true)]
        [string]$Message
    )
    Write-Host "`n  $Message" -ForegroundColor Yellow
    Write-Host ("  " + "-" * $Message.Length) -ForegroundColor Yellow
}

function Write-LogInfo {
    param (
        [Parameter(Mandatory = $true)]
        [string]$Message,
        [Parameter(Mandatory = $false)]
        [int]$Indent = 0
    )
    Write-Host ("  " * ($Indent + 1) + $Message) -ForegroundColor White
}

function Write-LogSuccess {
    param (
        [Parameter(Mandatory = $true)]
        [string]$Message,
        [Parameter(Mandatory = $false)]
        [int]$Indent = 0
    )
    Write-Host ("  " * ($Indent + 1) + "✓ $Message") -ForegroundColor Green
}

function Write-LogWarning {
    param (
        [Parameter(Mandatory = $true)]
        [string]$Message,
        [Parameter(Mandatory = $false)]
        [int]$Indent = 0
    )
    Write-Host ("  " * ($Indent + 1) + "⚠ $Message") -ForegroundColor Yellow
}

function Write-LogError {
    param (
        [Parameter(Mandatory = $true)]
        [string]$Message,
        [Parameter(Mandatory = $false)]
        [int]$Indent = 0
    )
    Write-Host ("  " * ($Indent + 1) + "✗ $Message") -ForegroundColor Red
}

# =============================================================================
# METADATA FUNCTIONS
# =============================================================================

function Get-GraphCommandMetadata {
    <#
    .SYNOPSIS
        Downloads or retrieves the cached Microsoft Graph SDK command metadata JSON.
    .DESCRIPTION
        The metadata maps Graph API URIs and HTTP methods to their required permissions.
        It is sourced from the public Microsoft Graph PowerShell SDK GitHub repository
        and cached locally in $env:TEMP to avoid repeated downloads.
        The cache is valid for 24 hours; use -ForceRefresh to bypass it.
    #>
    param (
        [Parameter(Mandatory = $false)]
        [switch]$ForceRefresh
    )

    # Resolve a cross-platform temp directory ($env:TEMP is Windows-only; macOS/Linux use $env:TMPDIR or /tmp)
    $tempRoot   = if ($env:TEMP) { $env:TEMP } elseif ($env:TMPDIR) { $env:TMPDIR } else { "/tmp" }
    $cacheDir   = Join-Path $tempRoot "MgGraphPermissionsRaw"
    $cachePath  = Join-Path $cacheDir "MgCommandMetadata.json"
    $cacheHours = 24

    $metadataUrls = @(
        "https://raw.githubusercontent.com/microsoftgraph/msgraph-sdk-powershell/main/src/Authentication/Authentication/custom/common/MgCommandMetadata.json",
        "https://raw.githubusercontent.com/microsoftgraph/msgraph-sdk-powershell/dev/src/Authentication/Authentication/custom/common/MgCommandMetadata.json"
    )

    # Return cached copy if it is still fresh
    if (-not $ForceRefresh -and (Test-Path $cachePath)) {
        $age = (Get-Date) - (Get-Item $cachePath).LastWriteTime
        if ($age.TotalHours -lt $cacheHours) {
            Write-LogSuccess "Using cached metadata (age: $([Math]::Round($age.TotalHours, 1))h, path: $cachePath)" -Indent 1
            try {
                $content = Get-Content $cachePath -Raw
                return $content | ConvertFrom-Json
            }
            catch {
                Write-LogWarning "Cache file is corrupted, re-downloading..." -Indent 1
            }
        }
        else {
            Write-LogInfo "Cache is stale ($([Math]::Round($age.TotalHours, 1))h old), re-downloading..." -Indent 1
        }
    }

    if (-not (Test-Path $cacheDir)) {
        New-Item -ItemType Directory -Path $cacheDir -Force | Out-Null
    }

    $downloaded = $false
    foreach ($url in $metadataUrls) {
        try {
            Write-LogInfo "Downloading Graph command metadata..." -Indent 1
            Write-LogInfo "Source: $url" -Indent 2

            $response = Invoke-RestMethod -Uri $url -Method GET -ErrorAction Stop

            # Persist to cache
            $response | ConvertTo-Json -Depth 20 -Compress | Set-Content -Path $cachePath -Encoding UTF8

            Write-LogSuccess "Metadata downloaded and cached ($($response.Count) entries)" -Indent 1
            $downloaded = $true
            return $response
        }
        catch {
            Write-LogWarning "Failed to download from $url — $($_.Exception.Message)" -Indent 1
        }
    }

    if (-not $downloaded) {
        throw "❌ Could not download Graph command metadata from any known URL. Check network connectivity."
    }
}

function Get-NormalizedUriTemplate {
    <#
    .SYNOPSIS
        Normalises a Graph URI template for comparison by replacing all {param} tokens
        with a canonical placeholder and lowercasing the result.
    .DESCRIPTION
        Both the user-supplied URI and the SDK metadata URIs may use different parameter
        names for the same position (e.g. {id} vs {applicationId}). Normalising both
        to {*} allows reliable equality matching regardless of naming differences.
    #>
    param (
        [Parameter(Mandatory = $true)]
        [string]$UriTemplate
    )

    return ($UriTemplate.TrimEnd('/').ToLower() -replace '\{[^}]+\}', '{*}')
}

function Find-EndpointPermissions {
    <#
    .SYNOPSIS
        Returns all SDK metadata entries that match the supplied URI, method, and API version.
    .DESCRIPTION
        Compares the normalised form of the user URI against each metadata entry's normalised
        URI. An entry matches when both normalised strings are equal, the API version matches,
        and (if specified) the HTTP method matches.
    #>
    param (
        [Parameter(Mandatory = $true)]
        $Metadata,

        [Parameter(Mandatory = $true)]
        [string]$Uri,

        [Parameter(Mandatory = $false)]
        [string]$Method = "",

        [Parameter(Mandatory = $true)]
        [string]$ApiVersion
    )

    $normalizedInput = Get-NormalizedUriTemplate -UriTemplate $Uri
    $matched         = @()

    foreach ($entry in $Metadata) {
        if ($entry.APIVersion -ne $ApiVersion) { continue }
        if ($Method -and $entry.Method -ne $Method.ToUpper()) { continue }

        $normalizedEntry = Get-NormalizedUriTemplate -UriTemplate $entry.URI

        if ($normalizedEntry -eq $normalizedInput) {
            $matched += $entry
        }
    }

    return $matched
}

function Get-PermissionName {
    <#
    .SYNOPSIS
        Extracts the permission name string from a permission object or string entry.
    .DESCRIPTION
        The metadata JSON represents permissions either as plain strings (older SDK versions)
        or as objects with a .Name property (newer SDK versions). This function normalises
        both forms to a plain string.
    #>
    param (
        [Parameter(Mandatory = $true)]
        $PermissionObject
    )

    if ($PermissionObject -is [string]) {
        return $PermissionObject
    }
    if ($PermissionObject.Name) {
        return $PermissionObject.Name
    }
    $str = $PermissionObject.ToString()
    if ($str -ne $PermissionObject.GetType().FullName) {
        return $str
    }
    return $null
}

function Test-IsApplicationPermission {
    <#
    .SYNOPSIS
        Returns $true if the permission name is an Application (not Delegated) permission.
    .DESCRIPTION
        Application permissions end in ".All" and do not contain "AccessAsUser".
        Delegated permissions typically contain "AccessAsUser" or follow other patterns.
    #>
    param (
        [Parameter(Mandatory = $true)]
        [string]$PermissionName
    )

    return ($PermissionName -notlike "*AccessAsUser*" -and $PermissionName -like "*.All")
}

# =============================================================================
# GRAPH REST API FUNCTIONS
# =============================================================================

function Get-GraphAccessToken {
    <#
    .SYNOPSIS
        Acquires a Microsoft Graph access token using the OAuth 2.0 client credentials flow.
    .DESCRIPTION
        Posts to the Entra ID token endpoint with client_id, client_secret, and the
        https://graph.microsoft.com/.default scope. Returns the raw access token string.
    #>
    param (
        [Parameter(Mandatory = $true,
            HelpMessage = "The Entra ID tenant ID")]
        [string]$TenantId,

        [Parameter(Mandatory = $true,
            HelpMessage = "The application (client) ID")]
        [string]$ClientId,

        [Parameter(Mandatory = $true,
            HelpMessage = "The client secret")]
        [string]$ClientSecret
    )

    try {
        Write-LogInfo "Acquiring access token (client credentials grant)..." -Indent 1

        $tokenUri = "https://login.microsoftonline.com/$TenantId/oauth2/v2.0/token"

        $body = @{
            grant_type    = "client_credentials"
            client_id     = $ClientId
            client_secret = $ClientSecret
            scope         = "https://graph.microsoft.com/.default"
        }

        $response = Invoke-RestMethod `
            -Uri         $tokenUri `
            -Method      POST `
            -Body        $body `
            -ContentType "application/x-www-form-urlencoded" `
            -ErrorAction Stop

        Write-LogSuccess "Access token acquired (expires in $($response.expires_in)s)" -Indent 1
        return $response.access_token
    }
    catch {
        throw "Failed to acquire access token for client '$ClientId': $($_.Exception.Message)"
    }
}

function Invoke-GraphRestRequest {
    <#
    .SYNOPSIS
        Makes an authenticated REST request to the Microsoft Graph API.
    .DESCRIPTION
        Constructs the full Graph endpoint URL from the supplied relative path and API version,
        attaches a Bearer token header, and returns the parsed JSON response body.
        Handles both relative paths (e.g. "servicePrincipals/{id}") and full URIs.
    #>
    param (
        [Parameter(Mandatory = $true,
            HelpMessage = "Bearer access token")]
        [string]$AccessToken,

        [Parameter(Mandatory = $true,
            HelpMessage = "Relative Graph path or full URI")]
        [string]$Path,

        [Parameter(Mandatory = $false)]
        [ValidateSet("GET", "POST", "PATCH", "PUT", "DELETE")]
        [string]$Method = "GET",

        [Parameter(Mandatory = $false,
            HelpMessage = "API version to use when building the full URI")]
        [string]$ApiVersion = "v1.0"
    )

    try {
        $fullUri = if ($Path -match '^https://') {
            $Path
        } else {
            "https://graph.microsoft.com/$ApiVersion/$($Path.TrimStart('/'))"
        }

        $headers = @{
            Authorization = "Bearer $AccessToken"
            "Content-Type" = "application/json"
        }

        return Invoke-RestMethod -Uri $fullUri -Method $Method -Headers $headers -ErrorAction Stop
    }
    catch {
        $statusCode = $_.Exception.Response.StatusCode.value__
        throw "Graph API request failed [$statusCode] $Method $Path — $($_.Exception.Message)"
    }
}

function Get-ServicePrincipalByAppId {
    <#
    .SYNOPSIS
        Retrieves a service principal object from Microsoft Graph by its application (client) ID.
    .DESCRIPTION
        Calls GET /servicePrincipals?$filter=appId eq '{appId}' and returns the first match.
        Only the id, appId, and displayName properties are selected to minimise response size.
    #>
    param (
        [Parameter(Mandatory = $true)]
        [string]$AccessToken,

        [Parameter(Mandatory = $true)]
        [string]$AppId
    )

    try {
        Write-LogInfo "Retrieving service principal for AppId: $AppId..." -Indent 1

        $encodedFilter = [uri]::EscapeDataString("appId eq '$AppId'")
        $path          = "servicePrincipals?`$filter=$encodedFilter&`$select=id,appId,displayName"
        $response      = Invoke-GraphRestRequest -AccessToken $AccessToken -Path $path

        $sp = $response.value | Select-Object -First 1
        if (-not $sp) {
            throw "No service principal found with appId '$AppId'"
        }

        return $sp
    }
    catch {
        throw "Failed to retrieve service principal: $($_.Exception.Message)"
    }
}

function Get-ServicePrincipalAppRoleAssignments {
    <#
    .SYNOPSIS
        Returns all app role assignments that have been granted to a service principal.
    .DESCRIPTION
        Calls GET /servicePrincipals/{id}/appRoleAssignments. The returned assignment objects
        contain appRoleId values which must be resolved against the resource service principal's
        appRoles collection to obtain the human-readable permission name (value).
    #>
    param (
        [Parameter(Mandatory = $true)]
        [string]$AccessToken,

        [Parameter(Mandatory = $true)]
        [string]$ServicePrincipalId
    )

    try {
        Write-LogInfo "Retrieving app role assignments for service principal '$ServicePrincipalId'..." -Indent 1

        $response = Invoke-GraphRestRequest -AccessToken $AccessToken -Path "servicePrincipals/$ServicePrincipalId/appRoleAssignments"
        return $response.value
    }
    catch {
        throw "Failed to retrieve app role assignments: $($_.Exception.Message)"
    }
}

function Get-MicrosoftGraphFirstPartyServicePrincipal {
    <#
    .SYNOPSIS
        Retrieves the Microsoft Graph first-party service principal.
    .DESCRIPTION
        Calls GET /servicePrincipals?$filter=appId eq '00000003-0000-0000-c000-000000000000'
        and returns the object including its appRoles collection. This collection is used to
        resolve appRoleId GUIDs from assignment objects back to permission names (values).
    #>
    param (
        [Parameter(Mandatory = $true)]
        [string]$AccessToken
    )

    try {
        Write-LogInfo "Retrieving Microsoft Graph first-party service principal..." -Indent 1

        $graphAppId    = "00000003-0000-0000-c000-000000000000"
        $encodedFilter = [uri]::EscapeDataString("appId eq '$graphAppId'")
        $path          = "servicePrincipals?`$filter=$encodedFilter&`$select=id,appId,appRoles"
        $response      = Invoke-GraphRestRequest -AccessToken $AccessToken -Path $path

        return $response.value | Select-Object -First 1
    }
    catch {
        throw "Failed to retrieve Microsoft Graph service principal: $($_.Exception.Message)"
    }
}

function Get-ServicePrincipalDirectoryRoles {
    <#
    .SYNOPSIS
        Returns the directory roles currently assigned to a service principal.
    .DESCRIPTION
        Calls GET /servicePrincipals/{id}/memberOf to retrieve all group and role memberships,
        then filters to only '#microsoft.graph.directoryRole' objects and fetches the full
        role details (displayName, description) for each via GET /directoryRoles/{id}.
    #>
    param (
        [Parameter(Mandatory = $true)]
        [string]$AccessToken,

        [Parameter(Mandatory = $true)]
        [string]$ServicePrincipalId
    )

    try {
        Write-LogInfo "Retrieving directory role memberships..." -Indent 1

        $memberOfResponse = Invoke-GraphRestRequest -AccessToken $AccessToken -Path "servicePrincipals/$ServicePrincipalId/memberOf"

        $directoryRoles = @()
        foreach ($member in $memberOfResponse.value) {
            if ($member.'@odata.type' -eq '#microsoft.graph.directoryRole') {
                try {
                    $role = Invoke-GraphRestRequest -AccessToken $AccessToken -Path "directoryRoles/$($member.id)"
                    $directoryRoles += $role
                }
                catch {
                    Write-LogWarning "Could not retrieve details for directory role '$($member.id)': $($_.Exception.Message)" -Indent 2
                }
            }
        }

        return $directoryRoles
    }
    catch {
        throw "Failed to retrieve directory roles: $($_.Exception.Message)"
    }
}

# =============================================================================
# MAIN EXECUTION
# =============================================================================

try {
    # -------------------------------------------------------------------------
    # Parameter validation
    # -------------------------------------------------------------------------
    if ($AppIdToBeEvaluated) {
        if (-not $TenantId) {
            Write-LogError "TenantId is required when AppIdToBeEvaluated is provided"
            exit 1
        }
        if (-not $ClientSecret) {
            Write-LogError "ClientSecret is required when AppIdToBeEvaluated is provided"
            exit 1
        }
    }

    # -------------------------------------------------------------------------
    # Step 1: Load Graph SDK command metadata
    # -------------------------------------------------------------------------
    Write-LogHeader "Loading Graph Command Metadata"

    $metadata = Get-GraphCommandMetadata -ForceRefresh:$RefreshMetadata

    if (-not $metadata -or $metadata.Count -eq 0) {
        Write-LogError "Metadata loaded but contained no entries"
        exit 1
    }

    Write-LogInfo "Total command entries available: $($metadata.Count)" -Indent 1

    # -------------------------------------------------------------------------
    # Step 2: Discover matching commands for the requested URI
    # -------------------------------------------------------------------------
    Write-LogHeader "Discovering Commands for URI"

    Write-LogInfo "URI:         $Uri" -Indent 1
    Write-LogInfo "API Version: $ApiVersion" -Indent 1
    if ($Method) {
        Write-LogInfo "Method:      $Method" -Indent 1
    } else {
        Write-LogInfo "Method:      (all methods)" -Indent 1
    }

    $normalizedInput = Get-NormalizedUriTemplate -UriTemplate $Uri
    Write-LogInfo "Normalised:  $normalizedInput" -Indent 1

    $matchedCommands = Find-EndpointPermissions -Metadata $metadata -Uri $Uri -Method $Method -ApiVersion $ApiVersion

    if (-not $matchedCommands -or $matchedCommands.Count -eq 0) {
        Write-LogError "No commands found for: $Method $Uri ($ApiVersion)"
        Write-Host ""
        Write-LogWarning "Troubleshooting suggestions:" -Indent 1
        Write-LogInfo "- Use {id} or any {placeholder} for resource identifiers: /groups/{id}" -Indent 2
        Write-LogInfo "- Try the alternate API version with -ApiVersion v1.0 or -ApiVersion beta" -Indent 2
        Write-LogInfo "- Some deep nested paths (e.g. Application Proxy segments) are not indexed" -Indent 2
        Write-LogInfo "- Use -RefreshMetadata to download the latest SDK metadata" -Indent 2
        exit 1
    }

    Write-LogSuccess "Found $($matchedCommands.Count) matching command(s)" -Indent 1

    # -------------------------------------------------------------------------
    # Step 3: Collect and categorise permissions across all matched commands
    # -------------------------------------------------------------------------
    $allPermissions   = @{}
    $applicationPerms = @()
    $delegatedPerms   = @()

    Write-LogSection "Analysing Commands and Permissions"

    foreach ($cmd in $matchedCommands) {
        Write-Host ""
        Write-LogInfo "Command: $($cmd.Command)" -Indent 1
        Write-LogInfo "Method:  $($cmd.Method)" -Indent 1
        Write-LogInfo "URI:     $($cmd.URI)" -Indent 1

        if ($cmd.Permissions -and $cmd.Permissions.Count -gt 0) {
            Write-LogInfo "Permissions:" -Indent 1

            foreach ($permObj in $cmd.Permissions) {
                $perm = Get-PermissionName -PermissionObject $permObj
                if (-not $perm) { continue }

                $isApp = Test-IsApplicationPermission -PermissionName $perm

                if ($isApp) {
                    if ($perm -notin $applicationPerms) { $applicationPerms += $perm }
                    Write-LogInfo "[App] $perm" -Indent 2
                } else {
                    if ($perm -notin $delegatedPerms) { $delegatedPerms += $perm }
                    Write-LogInfo "[Del] $perm" -Indent 2
                }

                if (-not $allPermissions.ContainsKey($perm)) {
                    $allPermissions[$perm] = @()
                }
                $allPermissions[$perm] += "$($cmd.Method) $($cmd.URI)"
            }
        } else {
            Write-LogWarning "No permissions documented in metadata for this command" -Indent 2
        }
    }

    # -------------------------------------------------------------------------
    # Step 4: Display permission summary
    # -------------------------------------------------------------------------
    Write-LogHeader "Permission Requirements Summary"

    if ($applicationPerms.Count -gt 0) {
        Write-LogSection "Application Permissions (for service principals / daemon apps)"
        foreach ($perm in ($applicationPerms | Sort-Object)) {
            Write-LogSuccess $perm -Indent 1
            foreach ($use in ($allPermissions[$perm] | Select-Object -Unique)) {
                Write-LogInfo "→ $use" -Indent 2
            }
        }
    }

    if ($delegatedPerms.Count -gt 0) {
        Write-LogSection "Delegated Permissions (for user context / interactive apps)"
        foreach ($perm in ($delegatedPerms | Sort-Object)) {
            Write-LogInfo $perm -Indent 1
            foreach ($use in ($allPermissions[$perm] | Select-Object -Unique)) {
                Write-LogInfo "→ $use" -Indent 2
            }
        }
    }

    if ($applicationPerms.Count -eq 0 -and $delegatedPerms.Count -eq 0) {
        Write-LogWarning "No permissions were documented in the metadata for the matched command(s)" -Indent 1
    }

    # -------------------------------------------------------------------------
    # Step 5: Evaluate enterprise application (Evaluation Mode)
    # -------------------------------------------------------------------------
    if ($AppIdToBeEvaluated) {

        Write-LogHeader "Evaluating Enterprise Application"

        # Acquire token
        $accessToken = Get-GraphAccessToken `
            -TenantId     $TenantId `
            -ClientId     $AppIdToBeEvaluated `
            -ClientSecret $ClientSecret

        # Retrieve service principal
        $sp = Get-ServicePrincipalByAppId -AccessToken $accessToken -AppId $AppIdToBeEvaluated
        Write-LogSuccess "Found service principal: $($sp.displayName)" -Indent 1
        Write-LogInfo "Object ID: $($sp.id)" -Indent 2

        # Retrieve granted app role assignments
        $appRoleAssignments = Get-ServicePrincipalAppRoleAssignments -AccessToken $accessToken -ServicePrincipalId $sp.id
        Write-LogSuccess "Retrieved $($appRoleAssignments.Count) app role assignment(s)" -Indent 1

        # Retrieve Microsoft Graph SP to resolve appRoleId → permission name
        $mgSp = Get-MicrosoftGraphFirstPartyServicePrincipal -AccessToken $accessToken

        # Build the list of granted permission values (human-readable names)
        $grantedPerms = @()
        foreach ($assignment in $appRoleAssignments) {
            $appRole = $mgSp.appRoles | Where-Object { $_.id -eq $assignment.appRoleId }
            if ($appRole -and $appRole.value) {
                $grantedPerms += $appRole.value
            }
        }

        Write-LogSuccess "Resolved $($grantedPerms.Count) granted permission name(s)" -Indent 1

        # Compare required vs granted
        Write-LogSection "Permission Comparison"

        $missingPerms = @()
        $coveredPerms = @()

        foreach ($requiredPerm in ($applicationPerms | Sort-Object)) {
            if ($requiredPerm -in $grantedPerms) {
                $coveredPerms += $requiredPerm
                Write-LogSuccess "$requiredPerm  [GRANTED]" -Indent 1
            } else {
                $missingPerms += $requiredPerm
                Write-LogError   "$requiredPerm  [MISSING]" -Indent 1
            }
        }

        # Directory role memberships
        Write-LogSection "Directory Role Assignments"

        $directoryRoles = Get-ServicePrincipalDirectoryRoles -AccessToken $accessToken -ServicePrincipalId $sp.id

        if ($directoryRoles.Count -eq 0) {
            Write-LogError "No directory roles currently assigned to this service principal" -Indent 1
        } else {
            foreach ($role in $directoryRoles) {
                Write-LogSuccess $role.displayName -Indent 1
            }
        }

        # Evaluation summary
        Write-LogHeader "Evaluation Summary"

        Write-LogSection "Coverage Status"
        Write-LogInfo    "Service Principal:              $($sp.displayName)" -Indent 1
        Write-LogInfo    "AppId:                          $AppIdToBeEvaluated" -Indent 1
        Write-LogInfo    "Total Required App Permissions: $($applicationPerms.Count)" -Indent 1
        Write-LogSuccess "Covered: $($coveredPerms.Count)" -Indent 1

        if ($missingPerms.Count -gt 0) {
            Write-LogError "Missing: $($missingPerms.Count)" -Indent 1
        } else {
            Write-LogSuccess "Missing: 0" -Indent 1
        }

        if ($missingPerms.Count -gt 0) {
            Write-LogSection "Missing Permissions — Remediation"
            Write-LogInfo "The following application permissions must be granted:" -Indent 1
            foreach ($perm in $missingPerms) {
                Write-LogError $perm -Indent 2
            }
            Write-Host ""
            Write-LogInfo "To grant these permissions:" -Indent 1
            Write-LogInfo "1. Open Azure Portal → Entra ID → Enterprise Applications" -Indent 2
            Write-LogInfo "2. Locate: $($sp.displayName)  (AppId: $AppIdToBeEvaluated)" -Indent 2
            Write-LogInfo "3. Navigate to: API Permissions" -Indent 2
            Write-LogInfo "4. Add each missing Microsoft Graph Application permission listed above" -Indent 2
            Write-LogInfo "5. Click 'Grant admin consent'" -Indent 2
        } else {
            Write-Host ""
            Write-LogSuccess "All required application permissions are granted — no action needed" -Indent 1
        }

        # Special requirements for role-assignable groups
        if ($Uri -like "*/groups*" -and ($Method -eq "POST" -or $Method -eq "")) {
            Write-LogSection "Special Requirement: Role-Assignable Groups"
            if ($directoryRoles.Count -eq 0) {
                Write-LogWarning "Creating groups with isAssignableToRole=true also requires one of:" -Indent 1
                Write-LogInfo    "- Privileged Role Administrator (directory role)" -Indent 2
                Write-LogInfo    "- Global Administrator (directory role)" -Indent 2
            } else {
                $hasRequiredRole = $directoryRoles | Where-Object {
                    $_.displayName -in @("Privileged Role Administrator", "Global Administrator")
                }
                if ($hasRequiredRole) {
                    Write-LogSuccess "Service principal holds the required directory role for role-assignable group creation" -Indent 1
                } else {
                    Write-LogWarning "Service principal has directory roles but none of the required ones for isAssignableToRole groups" -Indent 1
                    Write-LogInfo    "Required: Privileged Role Administrator OR Global Administrator" -Indent 2
                }
            }
        }

    } else {

        # -------------------------------------------------------------------------
        # Discovery Mode only — no app evaluation
        # -------------------------------------------------------------------------
        Write-LogHeader "Recommendations"

        if ($applicationPerms.Count -gt 0) {
            Write-LogSection "For Service Principal (Application Authentication)"
            Write-LogInfo "Grant these Microsoft Graph Application permissions in Azure Portal:" -Indent 1
            foreach ($perm in ($applicationPerms | Sort-Object)) {
                Write-LogInfo "  • $perm" -Indent 1
            }
            Write-Host ""
            Write-LogInfo "Admin consent is required for all application permissions" -Indent 1

            if ($Uri -like "*/groups*" -and ($Method -eq "POST" -or $Method -eq "")) {
                Write-Host ""
                Write-LogWarning "For groups with isAssignableToRole=true, the service principal also needs:" -Indent 1
                Write-LogInfo    "  • Privileged Role Administrator  OR  Global Administrator  (directory role)" -Indent 1
            }
        }

        if ($delegatedPerms.Count -gt 0) {
            Write-LogSection "For Delegated (User Context) Authentication"
            Write-LogInfo "Grant these Microsoft Graph Delegated permissions:" -Indent 1
            foreach ($perm in ($delegatedPerms | Sort-Object)) {
                Write-LogInfo "  • $perm" -Indent 1
            }
        }
    }

    # -------------------------------------------------------------------------
    # Documentation
    # -------------------------------------------------------------------------
    Write-LogHeader "Documentation References"
    Write-LogInfo "Microsoft Graph Permissions Reference:" -Indent 1
    Write-LogInfo "  https://learn.microsoft.com/en-us/graph/permissions-reference" -Indent 1
    Write-Host ""
    Write-LogInfo "Microsoft Graph API Reference:" -Indent 1
    Write-LogInfo "  https://learn.microsoft.com/en-us/graph/api/overview" -Indent 1
    Write-Host ""
    Write-LogInfo "Graph Explorer (interactive permission discovery):" -Indent 1
    Write-LogInfo "  https://developer.microsoft.com/en-us/graph/graph-explorer" -Indent 1
    Write-Host ""
}
catch {
    Write-LogError "Script execution failed: $($_.Exception.Message)"
    if ($_.ScriptStackTrace) {
        Write-Host "`n  Stack trace:" -ForegroundColor DarkGray
        Write-Host $_.ScriptStackTrace -ForegroundColor DarkGray
    }
    exit 1
}
