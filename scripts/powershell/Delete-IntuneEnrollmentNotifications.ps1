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
    HelpMessage="Specify the ID of a specific enrollment notification configuration to delete")]
    [ValidateNotNullOrEmpty()]
    [string]$NotificationId,
    
    [Parameter(Mandatory=$false,
    HelpMessage="Delete all enrollment notification configurations")]
    [bool]$DeleteAll = $false,
    
    [Parameter(Mandatory=$false,
    HelpMessage="Filter by device enrollment configuration type")]
    [ValidateSet('enrollmentNotificationsConfiguration', 'limit', 'platformRestrictions', 'windowsHelloForBusiness', 'windows10EnrollmentCompletionPageConfiguration', 'singlePlatformRestriction')]
    [string]$ConfigurationType,
    
    [Parameter(Mandatory=$false,
    HelpMessage="Filter by platform type (e.g., 'ios', 'android', 'windows', 'mac', 'androidForWork')")]
    [ValidateSet('ios', 'android', 'windows', 'mac', 'androidForWork')]
    [string]$PlatformType,
    
    [Parameter(Mandatory=$false,
    HelpMessage="Show what would be deleted without actually deleting")]
    [switch]$WhatIf
)

  # Delete a specific configuration
  #./Delete-IntuneEnrollmentNotifications.ps1 -TenantId "xxx" -ClientId "xxx" -ClientSecret "xxx" -NotificationId "12345678-1234-1234-1234-123456789012"

  # Delete enrollment notification configurations (What-If mode)
  #./Delete-IntuneEnrollmentNotifications.ps1 -TenantId "xxx" -ClientId "xxx" -ClientSecret "xxx" -DeleteAll $true -ConfigurationType "enrollmentNotificationsConfiguration" -WhatIf
  
  # Delete enrollment notification configurations for specific platform (What-If mode)  
  #./Delete-IntuneEnrollmentNotifications.ps1 -TenantId "xxx" -ClientId "xxx" -ClientSecret "xxx" -DeleteAll $true -ConfigurationType "enrollmentNotificationsConfiguration" -PlatformType "androidForWork" -WhatIf

  # Delete all enrollment configurations
  #./Delete-IntuneEnrollmentNotifications.ps1 -TenantId "xxx" -ClientId "xxx" -ClientSecret "xxx" -DeleteAll $true

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

function Get-EnrollmentNotificationConfigurations {
    param (
        [Parameter(Mandatory=$false)]
        [string]$FilterConfiguration,
        [Parameter(Mandatory=$false)]
        [string]$FilterPlatform
    )
    
    try {
        Write-Host "🔍 Getting enrollment notification configurations..." -ForegroundColor Cyan
        
        $baseUri = "https://graph.microsoft.com/beta/deviceManagement/deviceEnrollmentConfigurations"
        
        if ($FilterConfiguration) {
            $filterUri = "$baseUri" + '?$expand=assignments&$filter=' + "deviceEnrollmentConfigurationType eq '$FilterConfiguration'"
            Write-Host "   Filter: deviceEnrollmentConfigurationType eq '$FilterConfiguration'" -ForegroundColor Gray
            Write-Host "   Endpoint: $filterUri" -ForegroundColor Gray
        } else {
            $filterUri = "$baseUri" + '?$expand=assignments'
            Write-Host "   Endpoint: $filterUri" -ForegroundColor Gray
        }
        
        $configurations = Get-PaginatedResults -InitialUri $filterUri
        Write-Host "   ✅ Found $($configurations.Count) configuration(s) from API" -ForegroundColor Green
        
        # Client-side platform filtering if specified
        if ($FilterPlatform -and $configurations.Count -gt 0) {
            Write-Host "   🔍 Applying client-side platform filter: $FilterPlatform" -ForegroundColor Cyan
            $configurations = $configurations | Where-Object { $_.platformType -eq $FilterPlatform }
            Write-Host "   ✅ After platform filtering: $($configurations.Count) configuration(s)" -ForegroundColor Green
        }
        
        # Show full details for ALL configurations
        if ($configurations.Count -gt 0) {
            Write-Host "   📋 All configurations found:" -ForegroundColor Magenta
            for ($i = 0; $i -lt $configurations.Count; $i++) {
                $config = $configurations[$i]
                $assignmentCount = if ($config.assignments) { $config.assignments.Count } else { 0 }
                
                Write-Host "     Configuration $($i + 1) of $($configurations.Count):" -ForegroundColor Cyan
                Write-Host "     • ID: $($config.id)" -ForegroundColor Yellow
                Write-Host "       Name: $($config.displayName)" -ForegroundColor Yellow
                Write-Host "       Type: $($config.deviceEnrollmentConfigurationType)" -ForegroundColor Yellow
                Write-Host "       Platform: $($config.platformType)" -ForegroundColor Yellow
                Write-Host "       Priority: $($config.priority)" -ForegroundColor Yellow
                Write-Host "       Assignments: $assignmentCount" -ForegroundColor Yellow
                
                # Show localized notification messages if template ID exists
                if ($config.notificationMessageTemplateId) {
                    Write-Host "       Notification Template ID: $($config.notificationMessageTemplateId)" -ForegroundColor Yellow
                    
                    # Check if template ID is the null GUID, create localized message if needed
                    if ($config.notificationMessageTemplateId -eq "00000000-0000-0000-0000-000000000000") {
                        Write-Host "       🔧 Null template ID detected, creating localized notification message..." -ForegroundColor Cyan
                        
                        $newMessage = New-LocalizedNotificationMessage -NotificationMessageTemplateId $config.notificationMessageTemplateId -Locale "en-us" -Subject "Default Enrollment Notification" -MessageTemplate "Your device enrollment is in progress." -IsDefault $true
                        
                        if ($newMessage) {
                            Write-Host "       ✅ Created default localized message" -ForegroundColor Green
                        }
                    }
                    
                    try {
                        $template = Get-NotificationTemplateWithMessages -NotificationMessageTemplateId $config.notificationMessageTemplateId
                        $localizedMessages = $template.localizedNotificationMessages
                        
                        if ($localizedMessages -and $localizedMessages.Count -gt 0) {
                            Write-Host "       Localized Messages ($($localizedMessages.Count)):" -ForegroundColor Yellow
                            foreach ($message in $localizedMessages) {
                                Write-Host "         · Message ID: $($message.id)" -ForegroundColor Green
                                Write-Host "         · Locale: $($message.locale)" -ForegroundColor Green
                                Write-Host "         · Subject: $($message.subject)" -ForegroundColor Green
                                Write-Host "         · Is Default: $($message.isDefault)" -ForegroundColor Green
                                Write-Host "         · Last Modified: $($message.lastModifiedDateTime)" -ForegroundColor Green
                                Write-Host ""
                            }
                        } else {
                            Write-Host "       Localized Messages: None found" -ForegroundColor Yellow
                        }
                    } catch {
                        Write-Host "       Localized Messages: Error retrieving - $($_.Exception.Message)" -ForegroundColor Red
                    }
                } else {
                    Write-Host "       Notification Template ID: None" -ForegroundColor Yellow
                }
                
                Write-Host ""
            }
        }
        
        return $configurations
    }
    catch {
        Write-Host "❌ Error retrieving enrollment notification configurations: $_" -ForegroundColor Red
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

function Get-EnrollmentNotificationConfigurationById {
    param (
        [Parameter(Mandatory=$true)]
        [string]$ConfigurationId
    )
    
    try {
        Write-Host "🔍 Getting enrollment notification configuration by ID..." -ForegroundColor Cyan
        Write-Host "   Configuration ID: $ConfigurationId" -ForegroundColor Gray
        
        $configUri = "https://graph.microsoft.com/beta/deviceManagement/deviceEnrollmentConfigurations/$ConfigurationId"
        Write-Host "   Endpoint: $configUri" -ForegroundColor Gray
        
        $configuration = Invoke-MgGraphRequest -Method GET -Uri $configUri
        Write-Host "   ✅ Configuration retrieved successfully" -ForegroundColor Green
        Write-Host ""
        
        return $configuration
    }
    catch {
        Write-Host "❌ Error retrieving enrollment notification configuration: $_" -ForegroundColor Red
        Write-Host ""
        if ($_.Exception.Response) {
            $statusCode = $_.Exception.Response.StatusCode
            $statusDescription = $_.Exception.Response.StatusDescription
            Write-Host "   Status Code: $statusCode" -ForegroundColor Red
            Write-Host "   Status Description: $statusDescription" -ForegroundColor Red
            if ($statusCode -eq 404) {
                Write-Host "   Configuration not found with ID: $ConfigurationId" -ForegroundColor Red
            }
        }
        throw
    }
}

function Get-NotificationTemplateWithMessages {
    param (
        [Parameter(Mandatory=$true)]
        [string]$NotificationMessageTemplateId
    )
    
    try {
        Write-Host "📧 Getting notification template with localized messages..." -ForegroundColor Cyan
        Write-Host "   Template ID: $NotificationMessageTemplateId" -ForegroundColor Gray
        
        $templateUri = "https://graph.microsoft.com/beta/deviceManagement/notificationMessageTemplates/$NotificationMessageTemplateId" + '?$expand=localizedNotificationMessages'
        Write-Host "   Endpoint: $templateUri" -ForegroundColor Gray
        
        $response = Invoke-MgGraphRequest -Method GET -Uri $templateUri
        Write-Host "   ✅ Retrieved template with $($response.localizedNotificationMessages.Count) localized message(s)" -ForegroundColor Green
        
        return $response
    }
    catch {
        Write-Host "❌ Error retrieving notification template: $_" -ForegroundColor Red
        Write-Host "   Template ID: $NotificationMessageTemplateId" -ForegroundColor Red
        if ($_.Exception.Response) {
            $statusCode = $_.Exception.Response.StatusCode
            $statusDescription = $_.Exception.Response.StatusDescription
            Write-Host "   Status Code: $statusCode" -ForegroundColor Red
            Write-Host "   Status Description: $statusDescription" -ForegroundColor Red
            if ($statusCode -eq 404) {
                Write-Host "   Template not found with ID: $NotificationMessageTemplateId" -ForegroundColor Red
            }
        }
        return $null
    }
}

function New-LocalizedNotificationMessage {
    param (
        [Parameter(Mandatory=$true)]
        [string]$NotificationMessageTemplateId,
        [Parameter(Mandatory=$true)]
        [string]$Locale,
        [Parameter(Mandatory=$true)]
        [string]$Subject,
        [Parameter(Mandatory=$true)]
        [string]$MessageTemplate,
        [Parameter(Mandatory=$false)]
        [bool]$IsDefault = $true
    )
    
    try {
        Write-Host "📝 Creating localized notification message..." -ForegroundColor Cyan
        Write-Host "   Template ID: $NotificationMessageTemplateId" -ForegroundColor Gray
        Write-Host "   Locale: $Locale" -ForegroundColor Gray
        Write-Host "   Subject: $Subject" -ForegroundColor Gray
        
        $createUri = "https://graph.microsoft.com/beta/deviceManagement/notificationMessageTemplates/$NotificationMessageTemplateId/localizedNotificationMessages"
        Write-Host "   Endpoint: $createUri" -ForegroundColor Gray
        
        $body = @{
            "@odata.type" = "#microsoft.graph.localizedNotificationMessage"
            "locale" = $Locale
            "subject" = $Subject
            "messageTemplate" = $MessageTemplate
            "isDefault" = $IsDefault
        }
        
        $response = Invoke-MgGraphRequest -Method POST -Uri $createUri -Body ($body | ConvertTo-Json)
        Write-Host "   ✅ Created localized message with ID: $($response.id)" -ForegroundColor Green
        
        return $response
    }
    catch {
        Write-Host "❌ Error creating localized notification message: $_" -ForegroundColor Red
        Write-Host "   Template ID: $NotificationMessageTemplateId" -ForegroundColor Red
        if ($_.Exception.Response) {
            $statusCode = $_.Exception.Response.StatusCode
            $statusDescription = $_.Exception.Response.StatusDescription
            Write-Host "   Status Code: $statusCode" -ForegroundColor Red
            Write-Host "   Status Description: $statusDescription" -ForegroundColor Red
        }
        return $null
    }
}

function Remove-LocalizedNotificationMessage {
    param (
        [Parameter(Mandatory=$true)]
        [string]$NotificationMessageTemplateId,
        [Parameter(Mandatory=$true)]
        [string]$LocalizedNotificationMessageId,
        [Parameter(Mandatory=$false)]
        [bool]$WhatIfPreference = $false
    )
    
    try {
        $deleteUri = "https://graph.microsoft.com/beta/deviceManagement/notificationMessageTemplates/$NotificationMessageTemplateId/localizedNotificationMessages/$LocalizedNotificationMessageId"
        
        if ($WhatIfPreference) {
            Write-Host "🔍 WHAT-IF: Would delete localized message with ID: $LocalizedNotificationMessageId" -ForegroundColor Yellow
            Write-Host "   Endpoint: $deleteUri" -ForegroundColor Gray
            return $true
        } else {
            Write-Host "🗑️ Deleting localized notification message..." -ForegroundColor Cyan
            Write-Host "   Template ID: $NotificationMessageTemplateId" -ForegroundColor Gray
            Write-Host "   Message ID: $LocalizedNotificationMessageId" -ForegroundColor Gray
            Write-Host "   Endpoint: $deleteUri" -ForegroundColor Gray
            
            $response = Invoke-MgGraphRequest -Method DELETE -Uri $deleteUri
            Write-Host "   ✅ Localized message deleted successfully" -ForegroundColor Green
            return $true
        }
    }
    catch {
        Write-Host "❌ Error deleting localized notification message: $_" -ForegroundColor Red
        if ($_.Exception.Response) {
            $statusCode = $_.Exception.Response.StatusCode
            $statusDescription = $_.Exception.Response.StatusDescription
            Write-Host "   Status Code: $statusCode" -ForegroundColor Red
            Write-Host "   Status Description: $statusDescription" -ForegroundColor Red
        }
        return $false
    }
}

function Remove-NotificationMessageTemplate {
    param (
        [Parameter(Mandatory=$true)]
        [string]$NotificationMessageTemplateId,
        [Parameter(Mandatory=$false)]
        [bool]$WhatIfPreference = $false
    )
    
    try {
        $deleteUri = "https://graph.microsoft.com/beta/deviceManagement/notificationMessageTemplates/$NotificationMessageTemplateId"
        
        if ($WhatIfPreference) {
            Write-Host "🔍 WHAT-IF: Would delete notification template with ID: $NotificationMessageTemplateId" -ForegroundColor Yellow
            Write-Host "   Endpoint: $deleteUri" -ForegroundColor Gray
            return $true
        } else {
            Write-Host "🗑️ Deleting notification message template..." -ForegroundColor Cyan
            Write-Host "   Template ID: $NotificationMessageTemplateId" -ForegroundColor Gray
            Write-Host "   Endpoint: $deleteUri" -ForegroundColor Gray
            
            $response = Invoke-MgGraphRequest -Method DELETE -Uri $deleteUri
            Write-Host "   ✅ Notification template deleted successfully" -ForegroundColor Green
            return $true
        }
    }
    catch {
        Write-Host "❌ Error deleting notification message template: $_" -ForegroundColor Red
        if ($_.Exception.Response) {
            $statusCode = $_.Exception.Response.StatusCode
            $statusDescription = $_.Exception.Response.StatusDescription
            Write-Host "   Status Code: $statusCode" -ForegroundColor Red
            Write-Host "   Status Description: $statusDescription" -ForegroundColor Red
        }
        return $false
    }
}

function Remove-EnrollmentNotificationConfigurationComplete {
    param (
        [Parameter(Mandatory=$true)]
        $Configuration,
        [Parameter(Mandatory=$false)]
        [bool]$WhatIfPreference = $false
    )
    
    $success = $true
    
    try {
        Write-Host "🔄 Processing complete deletion for configuration: $($Configuration.displayName)" -ForegroundColor Cyan
        Write-Host "   Configuration ID: $($Configuration.id)" -ForegroundColor Gray
        
        # Step 1: Handle notification template and localized messages if present
        if ($Configuration.notificationMessageTemplateId -and $Configuration.notificationMessageTemplateId -ne "00000000-0000-0000-0000-000000000000") {
            Write-Host "   📧 Processing notification template: $($Configuration.notificationMessageTemplateId)" -ForegroundColor Cyan
            
            try {
                # Get template with localized messages
                $template = Get-NotificationTemplateWithMessages -NotificationMessageTemplateId $Configuration.notificationMessageTemplateId
                
                if ($template -and $template.localizedNotificationMessages) {
                    Write-Host "   🗑️ Deleting $($template.localizedNotificationMessages.Count) localized message(s)..." -ForegroundColor Cyan
                    
                    foreach ($message in $template.localizedNotificationMessages) {
                        $messageSuccess = Remove-LocalizedNotificationMessage -NotificationMessageTemplateId $Configuration.notificationMessageTemplateId -LocalizedNotificationMessageId $message.id -WhatIfPreference $WhatIfPreference
                        if (-not $messageSuccess) {
                            $success = $false
                        }
                    }
                }
                
                # Delete the notification template
                Write-Host "   🗑️ Deleting notification template..." -ForegroundColor Cyan
                $templateSuccess = Remove-NotificationMessageTemplate -NotificationMessageTemplateId $Configuration.notificationMessageTemplateId -WhatIfPreference $WhatIfPreference
                if (-not $templateSuccess) {
                    $success = $false
                }
                
            } catch {
                Write-Host "   ⚠️ Could not process notification template (may not exist): $($_.Exception.Message)" -ForegroundColor Yellow
            }
        } elseif ($Configuration.notificationMessageTemplateId -eq "00000000-0000-0000-0000-000000000000") {
            Write-Host "   📧 Skipping null template ID (00000000-0000-0000-0000-000000000000)" -ForegroundColor Gray
        }
        
        # Step 2: Delete the enrollment configuration
        Write-Host "   🗑️ Deleting enrollment notification configuration..." -ForegroundColor Cyan
        $configSuccess = Remove-EnrollmentNotificationConfiguration -ConfigurationId $Configuration.id -WhatIfPreference $WhatIfPreference
        if (-not $configSuccess) {
            $success = $false
        }
        
        if ($success) {
            if ($WhatIfPreference) {
                Write-Host "   ✅ WHAT-IF: All components would be deleted successfully" -ForegroundColor Green
            } else {
                Write-Host "   ✅ Complete deletion successful" -ForegroundColor Green
            }
        } else {
            Write-Host "   ❌ Some deletion operations failed" -ForegroundColor Red
        }
        
        return $success
        
    } catch {
        Write-Host "❌ Error during complete deletion process: $_" -ForegroundColor Red
        return $false
    }
}

function Remove-EnrollmentNotificationConfiguration {
    param (
        [Parameter(Mandatory=$true)]
        [string]$ConfigurationId,
        [Parameter(Mandatory=$false)]
        [bool]$WhatIfPreference = $false
    )
    
    try {
        $deleteUri = "https://graph.microsoft.com/beta/deviceManagement/deviceEnrollmentConfigurations/$ConfigurationId"
        
        if ($WhatIfPreference) {
            Write-Host "🔍 WHAT-IF: Would delete configuration with ID: $ConfigurationId" -ForegroundColor Yellow
            Write-Host "   Endpoint: $deleteUri" -ForegroundColor Gray
            return $true
        } else {
            Write-Host "🗑️ Deleting enrollment notification configuration..." -ForegroundColor Cyan
            Write-Host "   Configuration ID: $ConfigurationId" -ForegroundColor Gray
            Write-Host "   Endpoint: $deleteUri" -ForegroundColor Gray
            
            $response = Invoke-MgGraphRequest -Method DELETE -Uri $deleteUri
            Write-Host "   ✅ Configuration deleted successfully" -ForegroundColor Green
            return $true
        }
    }
    catch {
        Write-Host "❌ Error deleting enrollment notification configuration: $_" -ForegroundColor Red
        if ($_.Exception.Response) {
            $statusCode = $_.Exception.Response.StatusCode
            $statusDescription = $_.Exception.Response.StatusDescription
            Write-Host "   Status Code: $statusCode" -ForegroundColor Red
            Write-Host "   Status Description: $statusDescription" -ForegroundColor Red
            if ($statusCode -eq 404) {
                Write-Host "   Configuration not found with ID: $ConfigurationId" -ForegroundColor Red
            }
        }
        return $false
    }
}

function Show-ConfigurationDetails {
    param (
        [Parameter(Mandatory=$true)]
        $Configuration
    )
    
    Write-Host "📋 Enrollment Notification Configuration Details:" -ForegroundColor Cyan
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    
    # Basic configuration information
    foreach ($field in @('id', 'displayName', 'description', 'deviceEnrollmentConfigurationType', 'platformType', 'templateType', 'priority', 'createdDateTime', 'lastModifiedDateTime', 'version', 'roleScopeTagIds', 'brandingOptions', 'notificationMessageTemplateId')) {
        if ($Configuration.PSObject.Properties[$field]) {
            if ($field -eq 'roleScopeTagIds' -and $Configuration.$field) {
                Write-Host "   • roleScopeTagIds: $($Configuration.$field -join ', ')" -ForegroundColor Green
            } else {
                Write-Host ("   • {0}: {1}" -f $field, $Configuration.$field) -ForegroundColor Green
            }
        }
    }
    
    # ESP specific fields
    if ($Configuration.deviceEnrollmentConfigurationType -eq 'windows10EnrollmentCompletionPageConfiguration') {
        foreach ($field in @('showInstallationProgress', 'blockDeviceSetupRetryByUser', 'allowDeviceResetOnInstallFailure', 'allowLogCollectionOnInstallFailure', 'customErrorMessage', 'installProgressTimeoutInMinutes', 'allowDeviceUseOnInstallFailure', 'selectedMobileAppIds', 'trackInstallProgressForAutopilotOnly', 'disableUserStatusTrackingAfterFirstUser')) {
            if ($Configuration.PSObject.Properties[$field]) {
                if ($field -eq 'selectedMobileAppIds' -and $Configuration.$field) {
                    Write-Host "   • selectedMobileAppIds: $($Configuration.$field -join ', ')" -ForegroundColor Green
                } else {
                    Write-Host ("   • {0}: {1}" -f $field, $Configuration.$field) -ForegroundColor Green
                }
            }
        }
    }
    
    # Assignment information
    if ($Configuration.assignments -and $Configuration.assignments.Count -gt 0) {
        Write-Host "   • assignments ($($Configuration.assignments.Count) assignment(s)):" -ForegroundColor Green
        foreach ($assignment in $Configuration.assignments) {
            Write-Host "     · Assignment ID: $($assignment.id)" -ForegroundColor Yellow
            Write-Host "     · Source: $($assignment.source)" -ForegroundColor Yellow
            if ($assignment.target) {
                Write-Host "     · Target Type: $($assignment.target.'@odata.type')" -ForegroundColor Yellow
                if ($assignment.target.groupId) {
                    Write-Host "     · Group ID: $($assignment.target.groupId)" -ForegroundColor Yellow
                }
                if ($assignment.target.deviceAndAppManagementAssignmentFilterType -and $assignment.target.deviceAndAppManagementAssignmentFilterType -ne "none") {
                    Write-Host "     · Filter Type: $($assignment.target.deviceAndAppManagementAssignmentFilterType)" -ForegroundColor Yellow
                    Write-Host "     · Filter ID: $($assignment.target.deviceAndAppManagementAssignmentFilterId)" -ForegroundColor Yellow
                }
            }
        }
    } else {
        Write-Host "   • assignments: None" -ForegroundColor Green
    }
    
    # Notification templates information
    if ($Configuration.notificationTemplates -and $Configuration.notificationTemplates.Count -gt 0) {
        Write-Host "   • notificationTemplates: $($Configuration.notificationTemplates -join ', ')" -ForegroundColor Green
    }
    
    # Localized notification messages
    if ($Configuration.notificationMessageTemplateId) {
        Write-Host "   • notificationMessageTemplateId: $($Configuration.notificationMessageTemplateId)" -ForegroundColor Green
        
        # Check if template ID is the null GUID, create localized message if needed
        if ($Configuration.notificationMessageTemplateId -eq "00000000-0000-0000-0000-000000000000") {
            Write-Host "   🔧 Null template ID detected, creating localized notification message..." -ForegroundColor Cyan
            
            $newMessage = New-LocalizedNotificationMessage -NotificationMessageTemplateId $Configuration.notificationMessageTemplateId -Locale "en-us" -Subject "Default Enrollment Notification" -MessageTemplate "Your device enrollment is in progress." -IsDefault $true
            
            if ($newMessage) {
                Write-Host "   ✅ Created default localized message" -ForegroundColor Green
            }
        }
        
        try {
            $template = Get-NotificationTemplateWithMessages -NotificationMessageTemplateId $Configuration.notificationMessageTemplateId
            $localizedMessages = $template.localizedNotificationMessages
            
            if ($localizedMessages -and $localizedMessages.Count -gt 0) {
                Write-Host "   • localizedNotificationMessages ($($localizedMessages.Count) message(s)):" -ForegroundColor Green
                foreach ($message in $localizedMessages) {
                    Write-Host "     · Message ID: $($message.id)" -ForegroundColor Yellow
                    Write-Host "     · Locale: $($message.locale)" -ForegroundColor Yellow
                    Write-Host "     · Subject: $($message.subject)" -ForegroundColor Yellow
                    Write-Host "     · Message Template: $($message.messageTemplate)" -ForegroundColor Yellow
                    Write-Host "     · Is Default: $($message.isDefault)" -ForegroundColor Yellow
                    Write-Host "     · Last Modified: $($message.lastModifiedDateTime)" -ForegroundColor Yellow
                    Write-Host ""
                }
            } else {
                Write-Host "   • localizedNotificationMessages: None found" -ForegroundColor Green
            }
        } catch {
            Write-Host "   • localizedNotificationMessages: Error retrieving messages" -ForegroundColor Red
        }
    }
    
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    Write-Host ""
}

function Show-DeletionSummary {
    param (
        [Parameter(Mandatory=$true)]
        [int]$TotalConfigurations,
        [Parameter(Mandatory=$true)]
        [int]$SuccessfulDeletions,
        [Parameter(Mandatory=$true)]
        [int]$FailedDeletions,
        [Parameter(Mandatory=$false)]
        [bool]$WhatIfPreference = $false
    )
    
    Write-Host "📊 Deletion Summary:" -ForegroundColor Cyan
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    Write-Host ("   • Total configurations found: {0}" -f $TotalConfigurations) -ForegroundColor Green
    
    if ($WhatIfPreference) {
        Write-Host ("   • Configurations that would be deleted: {0}" -f $SuccessfulDeletions) -ForegroundColor Yellow
        Write-Host ("   • Configurations that would fail to delete: {0}" -f $FailedDeletions) -ForegroundColor Red
    } else {
        Write-Host ("   • Successfully deleted: {0}" -f $SuccessfulDeletions) -ForegroundColor Green
        Write-Host ("   • Failed to delete: {0}" -f $FailedDeletions) -ForegroundColor Red
    }
    
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    Write-Host ""
}

# Validation
if (-not $NotificationId -and -not $DeleteAll) {
    Write-Host "❌ Error: You must specify either -NotificationId or -DeleteAll" -ForegroundColor Red
    Write-Host "   Use -NotificationId to delete a specific configuration" -ForegroundColor Yellow
    Write-Host "   Use -DeleteAll to delete all configurations (optionally filtered by type)" -ForegroundColor Yellow
    exit 1
}

if ($NotificationId -and $DeleteAll) {
    Write-Host "❌ Error: You cannot specify both -NotificationId and -DeleteAll" -ForegroundColor Red
    Write-Host "   Use -NotificationId to delete a specific configuration" -ForegroundColor Yellow
    Write-Host "   Use -DeleteAll to delete all configurations (optionally filtered by type)" -ForegroundColor Yellow
    exit 1
}

try {
    Write-Host "🔐 Connecting to Microsoft Graph..." -ForegroundColor Cyan
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    
    $secureClientSecret = ConvertTo-SecureString -String $ClientSecret -AsPlainText -Force
    $clientSecretCredential = New-Object -TypeName System.Management.Automation.PSCredential -ArgumentList $ClientId, $secureClientSecret
    
    Connect-MgGraph -ClientSecretCredential $clientSecretCredential -TenantId $TenantId
    Write-Host "✅ Connected to Microsoft Graph successfully" -ForegroundColor Green
    Write-Host ""
    
    if ($NotificationId) {
        # Delete specific configuration by ID
        Write-Host "🎯 Single Configuration Deletion Mode" -ForegroundColor Cyan
        Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
        
        $configuration = Get-EnrollmentNotificationConfigurationById -ConfigurationId $NotificationId
        
        if ($null -ne $configuration) {
            Show-ConfigurationDetails -Configuration $configuration
            
            if ($WhatIf) {
                Write-Host "🔍 WHAT-IF MODE: The following configuration would be deleted:" -ForegroundColor Yellow
                Write-Host "   ID: $($configuration.id)" -ForegroundColor Yellow
                Write-Host "   Name: $($configuration.displayName)" -ForegroundColor Yellow
                Write-Host "   Type: $($configuration.deviceEnrollmentConfigurationType)" -ForegroundColor Yellow
            } else {
                $userConfirmation = Read-Host "❓ Are you sure you want to delete this configuration? (y/N)"
                if ($userConfirmation -eq 'y' -or $userConfirmation -eq 'Y') {
                    $success = Remove-EnrollmentNotificationConfigurationComplete -Configuration $configuration -WhatIfPreference $false
                    if ($success) {
                        Write-Host "🎉 Configuration deleted successfully!" -ForegroundColor Green
                    }
                } else {
                    Write-Host "❌ Operation cancelled by user" -ForegroundColor Yellow
                }
            }
        } else {
            Write-Host "📊 No configuration found with the specified ID" -ForegroundColor Yellow
        }
        
    } elseif ($DeleteAll) {
        # Delete all configurations (optionally filtered by type)
        Write-Host "🎯 Bulk Deletion Mode" -ForegroundColor Cyan
        Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
        
        $configurations = Get-EnrollmentNotificationConfigurations -FilterConfiguration $ConfigurationType -FilterPlatform $PlatformType
        
        if ($configurations.Count -gt 0) {
            Write-Host "📋 Found $($configurations.Count) configuration(s) to process:" -ForegroundColor Cyan
            foreach ($config in $configurations) {
                Write-Host "   • ID: $($config.id) | Name: $($config.displayName) | Type: $($config.deviceEnrollmentConfigurationType) | Platform: $($config.platformType)" -ForegroundColor Gray
            }
            Write-Host ""
            
            if ($WhatIf) {
                Write-Host "🔍 WHAT-IF MODE: The following configurations would be deleted:" -ForegroundColor Yellow
                foreach ($config in $configurations) {
                    Write-Host "   • $($config.displayName) ($($config.id))" -ForegroundColor Yellow
                    $whatIfSuccess = Remove-EnrollmentNotificationConfigurationComplete -Configuration $config -WhatIfPreference $true
                    Write-Host "" # Add spacing between what-if operations
                }
                Show-DeletionSummary -TotalConfigurations $configurations.Count -SuccessfulDeletions $configurations.Count -FailedDeletions 0 -WhatIfPreference $true
            } else {
                $userConfirmation = Read-Host "❓ Are you sure you want to delete ALL $($configurations.Count) configuration(s)? (y/N)"
                if ($userConfirmation -eq 'y' -or $userConfirmation -eq 'Y') {
                    
                    $successCount = 0
                    $failCount = 0
                    
                    foreach ($config in $configurations) {
                        $success = Remove-EnrollmentNotificationConfigurationComplete -Configuration $config -WhatIfPreference $false
                        if ($success) {
                            $successCount++
                        } else {
                            $failCount++
                        }
                        Write-Host "" # Add spacing between deletions
                    }
                    
                    Show-DeletionSummary -TotalConfigurations $configurations.Count -SuccessfulDeletions $successCount -FailedDeletions $failCount -WhatIfPreference $false
                    
                    if ($successCount -eq $configurations.Count) {
                        Write-Host "🎉 All configurations deleted successfully!" -ForegroundColor Green
                    } elseif ($successCount -gt 0) {
                        Write-Host "⚠️ Some configurations were deleted, but $failCount failed" -ForegroundColor Yellow
                    } else {
                        Write-Host "❌ No configurations were deleted" -ForegroundColor Red
                    }
                } else {
                    Write-Host "❌ Operation cancelled by user" -ForegroundColor Yellow
                }
            }
        } else {
            if ($ConfigurationType) {
                Write-Host "📊 No configurations found with type: $ConfigurationType" -ForegroundColor Yellow
            } else {
                Write-Host "📊 No enrollment notification configurations found" -ForegroundColor Yellow
            }
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