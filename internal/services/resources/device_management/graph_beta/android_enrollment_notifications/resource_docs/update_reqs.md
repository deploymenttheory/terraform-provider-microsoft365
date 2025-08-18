update base resource

Request URL
https://graph.microsoft.com/beta/deviceManagement/deviceEnrollmentConfigurations/3ca5db68-ea51-4432-ba00-0e54d115b6c5_EnrollmentNotificationsConfiguration
Request Method
PATCH

{"@odata.type":"#microsoft.graph.deviceEnrollmentNotificationConfiguration","displayName":"test 2","description":"test 2","roleScopeTagIds":["0"],"id":"3ca5db68-ea51-4432-ba00-0e54d115b6c5_EnrollmentNotificationsConfiguration"}

update localized notification messages

Request URL
https://graph.microsoft.com/beta/deviceManagement/notificationMessageTemplates/a5c117c2-24a7-400e-bff9-e0c411a363e3/localizedNotificationMessages/a5c117c2-24a7-400e-bff9-e0c411a363e3_en-us
Request Method
PATCH

{"subject":"test 2","messageTemplate":"test 2"}

Request URL
https://graph.microsoft.com/beta/deviceManagement/notificationMessageTemplates/93d8ee09-7d8a-4e4f-a3d2-f8b4c3d70b8d/localizedNotificationMessages/93d8ee09-7d8a-4e4f-a3d2-f8b4c3d70b8d_en-us
Request Method
PATCH

{"subject":"test 2","messageTemplate":"test 2"}

update branding options for email

Request URL
https://graph.microsoft.com/beta/deviceManagement/notificationMessageTemplates/a5c117c2-24a7-400e-bff9-e0c411a363e3
Request Method
PATCH

{"brandingOptions":"includeCompanyLogo,includeCompanyPortalLink,includeContactInformation,includeDeviceDetails"}