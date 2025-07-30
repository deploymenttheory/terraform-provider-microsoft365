{"@odata.type":"#microsoft.graph.windowsUpdateForBusinessConfiguration","id":"168f15f1-92fb-428c-88cc-0ec04e8d2a66","displayName":"Standard Windows Update Ring","description":"Default update ring for standard workstations","roleScopeTagIds":["0"],"microsoftUpdateServiceAllowed":true,"driversExcluded":false,"qualityUpdatesDeferralPeriodInDays":7,"featureUpdatesDeferralPeriodInDays":14,"allowWindows11Upgrade":false,"qualityUpdatesPaused":false,"featureUpdatesPaused":false,"businessReadyUpdatesOnly":"businessReadyOnly","skipChecksBeforeRestart":false,"automaticUpdateMode":"autoInstallAtMaintenanceTime","installationSchedule":{"@odata.type":"#microsoft.graph.windowsUpdateActiveHoursInstall","activeHoursStart":"08:00:00.0000000","activeHoursEnd":"17:00:00.0000000"},"userPauseAccess":"enabled","userWindowsUpdateScanAccess":"enabled","updateNotificationLevel":"restartWarningsOnly","updateWeeks":null,"featureUpdatesRollbackWindowInDays":10,"deadlineForFeatureUpdatesInDays":2,"deadlineForQualityUpdatesInDays":7,"deadlineGracePeriodInDays":1,"postponeRebootUntilAfterDeadline":false,"engagedRestartDeadlineInDays":null,"engagedRestartSnoozeScheduleInDays":null,"engagedRestartTransitionScheduleInDays":null,"engagedRestartSnoozeScheduleForFeatureUpdatesInDays":null,"engagedRestartTransitionScheduleForFeatureUpdatesInDays":null,"autoRestartNotificationDismissal":"notConfigured","scheduleRestartWarningInHours":4,"scheduleImminentRestartWarningInMinutes":30}


{"@odata.type":"#microsoft.graph.windowsUpdateForBusinessConfiguration","id":"168f15f1-92fb-428c-88cc-0ec04e8d2a66","displayName":"Standard Windows Update Ring","description":"Default update ring for standard workstations","roleScopeTagIds":["0"],"microsoftUpdateServiceAllowed":true,"driversExcluded":false,"qualityUpdatesDeferralPeriodInDays":7,"featureUpdatesDeferralPeriodInDays":14,"allowWindows11Upgrade":true,"qualityUpdatesPaused":false,"featureUpdatesPaused":false,"businessReadyUpdatesOnly":"windowsInsiderBuildSlow","skipChecksBeforeRestart":false,"automaticUpdateMode":"autoInstallAndRebootAtScheduledTime","installationSchedule":{"@odata.type":"#microsoft.graph.windowsUpdateScheduledInstall","scheduledInstallDay":"everyday","scheduledInstallTime":"03:00:00.0000000"},"userPauseAccess":"enabled","userWindowsUpdateScanAccess":"enabled","updateNotificationLevel":"disableAllNotifications","updateWeeks":"everyWeek","featureUpdatesRollbackWindowInDays":10,"deadlineForFeatureUpdatesInDays":2,"deadlineForQualityUpdatesInDays":7,"deadlineGracePeriodInDays":1,"postponeRebootUntilAfterDeadline":false,"engagedRestartDeadlineInDays":null,"engagedRestartSnoozeScheduleInDays":null,"engagedRestartTransitionScheduleInDays":null,"engagedRestartSnoozeScheduleForFeatureUpdatesInDays":null,"engagedRestartTransitionScheduleForFeatureUpdatesInDays":null,"autoRestartNotificationDismissal":"notConfigured","scheduleRestartWarningInHours":4,"scheduleImminentRestartWarningInMinutes":30}

{"@odata.type":"#microsoft.graph.windowsUpdateForBusinessConfiguration","id":"168f15f1-92fb-428c-88cc-0ec04e8d2a66","displayName":"Standard Windows Update Ring","description":"Default update ring for standard workstations","roleScopeTagIds":["0"],"microsoftUpdateServiceAllowed":true,"driversExcluded":false,"qualityUpdatesDeferralPeriodInDays":7,"featureUpdatesDeferralPeriodInDays":14,"allowWindows11Upgrade":true,"qualityUpdatesPaused":false,"featureUpdatesPaused":false,"businessReadyUpdatesOnly":"windowsInsiderBuildSlow","skipChecksBeforeRestart":false,"automaticUpdateMode":"notifyDownload","installationSchedule":null,"userPauseAccess":"enabled","userWindowsUpdateScanAccess":"enabled","updateNotificationLevel":"disableAllNotifications","updateWeeks":null,"featureUpdatesRollbackWindowInDays":10,"deadlineForFeatureUpdatesInDays":2,"deadlineForQualityUpdatesInDays":7,"deadlineGracePeriodInDays":1,"postponeRebootUntilAfterDeadline":false,"engagedRestartDeadlineInDays":null,"engagedRestartSnoozeScheduleInDays":null,"engagedRestartTransitionScheduleInDays":null,"engagedRestartSnoozeScheduleForFeatureUpdatesInDays":null,"engagedRestartTransitionScheduleForFeatureUpdatesInDays":null,"autoRestartNotificationDismissal":"notConfigured","scheduleRestartWarningInHours":4,"scheduleImminentRestartWarningInMinutes":30}

{"@odata.type":"#microsoft.graph.windowsUpdateForBusinessConfiguration","id":"bb1e2020-a7ad-4c52-818c-051dcc06ccda","displayName":"Standard Windows Update Ring","description":"Default update ring for standard workstations","roleScopeTagIds":["0"],"microsoftUpdateServiceAllowed":true,"driversExcluded":false,"qualityUpdatesDeferralPeriodInDays":7,"featureUpdatesDeferralPeriodInDays":14,"allowWindows11Upgrade":false,"qualityUpdatesPaused":false,"featureUpdatesPaused":false,"businessReadyUpdatesOnly":"businessReadyOnly","skipChecksBeforeRestart":false,"automaticUpdateMode":"autoInstallAndRebootAtScheduledTime","installationSchedule":{"@odata.type":"#microsoft.graph.windowsUpdateScheduledInstall","scheduledInstallDay":"everyday","scheduledInstallTime":"03:00:00.0000000"},"userPauseAccess":"enabled","userWindowsUpdateScanAccess":"enabled","updateNotificationLevel":"restartWarningsOnly","updateWeeks":"firstWeek,secondWeek,thirdWeek,fourthWeek,everyWeek","featureUpdatesRollbackWindowInDays":10,"deadlineForFeatureUpdatesInDays":2,"deadlineForQualityUpdatesInDays":7,"deadlineGracePeriodInDays":1,"postponeRebootUntilAfterDeadline":false,"engagedRestartDeadlineInDays":null,"engagedRestartSnoozeScheduleInDays":null,"engagedRestartTransitionScheduleInDays":null,"engagedRestartSnoozeScheduleForFeatureUpdatesInDays":null,"engagedRestartTransitionScheduleForFeatureUpdatesInDays":null,"autoRestartNotificationDismissal":"notConfigured","scheduleRestartWarningInHours":4,"scheduleImminentRestartWarningInMinutes":30}

pause feature updates

{"@odata.type":"#microsoft.graph.windowsUpdateForBusinessConfiguration","featureUpdatesPaused":true}

pause quality updates

{"@odata.type":"#microsoft.graph.windowsUpdateForBusinessConfiguration","qualityUpdatesPaused":true}


resume feature updates

{"@odata.type":"#microsoft.graph.windowsUpdateForBusinessConfiguration","featureUpdatesPaused":false}

resume quality updates

{"@odata.type":"#microsoft.graph.windowsUpdateForBusinessConfiguration","qualityUpdatesPaused":false}

extend feature update pause (maximum 35 days)

Request URL
https://graph.microsoft.com/beta/deviceManagement/deviceConfigurations/23c5f5bc-a218-454a-a579-31b8864147fc/microsoft.graph.windowsUpdateForBusinessConfiguration/extendFeatureUpdatesPause
Request Method
POST

feature update uninstall

{"@odata.type":"#microsoft.graph.windowsUpdateForBusinessConfiguration","featureUpdatesWillBeRolledBack":true}

quality update uninstall

{"@odata.type":"#microsoft.graph.windowsUpdateForBusinessConfiguration","qualityUpdatesWillBeRolledBack":true}

