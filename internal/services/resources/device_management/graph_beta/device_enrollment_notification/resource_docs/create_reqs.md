// create base resource

Request URL
https://graph.microsoft.com/beta/deviceManagement/deviceEnrollmentConfigurations
Request Method
POST

{"@odata.type":"#microsoft.graph.deviceEnrollmentNotificationConfiguration","displayName":"enrollment notification","description":"enrollment notification description","platformType":"androidForWork","defaultLocale":"en-US","roleScopeTagIds":["0"],"brandingOptions":"includeCompanyLogo,includeCompanyName,includeCompanyPortalLink,includeContactInformation,includeDeviceDetails","notificationTemplates":["email_00000000-0000-0000-0000-000000000000","push_00000000-0000-0000-0000-000000000000"]}

response

{
    "@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/deviceEnrollmentConfigurations/$entity",
    "@odata.type": "#microsoft.graph.deviceEnrollmentNotificationConfiguration",
    "id": "ebf4ccac-a764-4254-bf13-9a9765ca8aa9_EnrollmentNotificationsConfiguration",
    "displayName": "enrollment notification",
    "description": "enrollment notification description",
    "priority": 1,
    "createdDateTime": "2025-08-18T12:23:12.4766585Z",
    "lastModifiedDateTime": "2025-08-18T12:23:12.4766585Z",
    "version": 1,
    "roleScopeTagIds": [
        "0"
    ],
    "deviceEnrollmentConfigurationType": "enrollmentNotificationsConfiguration",
    "platformType": "androidForWork",
    "templateType": "0",
    "notificationMessageTemplateId": "00000000-0000-0000-0000-000000000000",
    "notificationTemplates": [
        "Email_3f881c29-a38e-4a83-abeb-b1a6312de85a",
        "Push_6a3edae4-ae98-4fab-b3fe-d3c4203b3056"
    ],
    "brandingOptions": "none",
    "defaultLocale": null
}

// create localized notification message 1

Request URL
https://graph.microsoft.com/beta/deviceManagement/notificationMessageTemplates/6a3edae4-ae98-4fab-b3fe-d3c4203b3056/localizedNotificationMessages
Request Method
POST

{"locale":"en-US","isDefault":true,"subject":"push notification subject","messageTemplate":"push notification message"}

response

{
    "@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/notificationMessageTemplates('6a3edae4-ae98-4fab-b3fe-d3c4203b3056')/localizedNotificationMessages/$entity",
    "id": "6a3edae4-ae98-4fab-b3fe-d3c4203b3056_en-us",
    "lastModifiedDateTime": "2025-08-18T12:23:12.8311541Z",
    "locale": "en-us",
    "subject": "push notification subject",
    "messageTemplate": "push notification message",
    "isDefault": true
}

// create localized notification message 2

Request URL
https://graph.microsoft.com/beta/deviceManagement/notificationMessageTemplates/3f881c29-a38e-4a83-abeb-b1a6312de85a/localizedNotificationMessages
Request Method
POST

{"locale":"en-US","isDefault":true,"subject":"Email Notification subject","messageTemplate":"Email Notification subject message"}

response

{
    "@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/notificationMessageTemplates('3f881c29-a38e-4a83-abeb-b1a6312de85a')/localizedNotificationMessages/$entity",
    "id": "3f881c29-a38e-4a83-abeb-b1a6312de85a_en-us",
    "lastModifiedDateTime": "2025-08-18T12:23:12.8344061Z",
    "locale": "en-us",
    "subject": "Email Notification subject",
    "messageTemplate": "Email Notification subject message",
    "isDefault": true
}

Request URL
https://graph.microsoft.com/beta/deviceManagement/notificationMessageTemplates/3f881c29-a38e-4a83-abeb-b1a6312de85a
Request Method
PATCH

{"brandingOptions":"includeCompanyLogo,includeCompanyName,includeCompanyPortalLink,includeContactInformation,includeDeviceDetails"}

response

{
    "@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/notificationMessageTemplates/$entity",
    "id": "3f881c29-a38e-4a83-abeb-b1a6312de85a",
    "lastModifiedDateTime": "2025-08-18T12:27:34.6986999Z",
    "displayName": "EnrollmentNotificationInternalMEO",
    "description": null,
    "defaultLocale": "en-us",
    "brandingOptions": "includeCompanyName,includeContactInformation,includeCompanyPortalLink,includeDeviceDetails",
    "roleScopeTagIds": [
        "0"
    ]
}
