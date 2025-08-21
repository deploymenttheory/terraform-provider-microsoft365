Request URL
https://graph.microsoft.com/beta/deviceManagement/notificationMessageTemplates
Request Method
POST

{"brandingOptions":"includeCompanyLogo,includeCompanyName,includeContactInformation,includeCompanyPortalLink","displayName":"acc-test","roleScopeTagIds":["0"]}

response

{
    "@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/notificationMessageTemplates/$entity",
    "id": "829c1a4e-fa0f-40fb-bdf2-31e8a7aefee8",
    "lastModifiedDateTime": "2025-08-20T12:01:21.9960945Z",
    "displayName": "acc-test",
    "description": null,
    "defaultLocale": null,
    "brandingOptions": "includeCompanyLogo,includeCompanyName,includeContactInformation,includeCompanyPortalLink",
    "roleScopeTagIds": [
        "0"
    ]
}

localizedNotificationMessages

Request URL
https://graph.microsoft.com/beta/deviceManagement/notificationMessageTemplates/829c1a4e-fa0f-40fb-bdf2-31e8a7aefee8/localizedNotificationMessages
Request Method
POST

{"isDefault":true,"locale":"bg-bg","messageTemplate":"some\n\nmultiline\n\ntext","subject":"acc-test"}

response

{
    "@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/notificationMessageTemplates('829c1a4e-fa0f-40fb-bdf2-31e8a7aefee8')/localizedNotificationMessages/$entity",
    "id": "829c1a4e-fa0f-40fb-bdf2-31e8a7aefee8_bg-bg",
    "lastModifiedDateTime": "2025-08-20T12:01:22.2615406Z",
    "locale": "bg-bg",
    "subject": "acc-test",
    "messageTemplate": "some\n\nmultiline\n\ntext",
    "isDefault": true
}

Request URL
https://graph.microsoft.com/beta/deviceManagement/notificationMessageTemplates/829c1a4e-fa0f-40fb-bdf2-31e8a7aefee8/localizedNotificationMessages
Request Method
POST

request

{"isDefault":false,"locale":"da-dk","messageTemplate":"some\n\nmultiline\n\ntext","subject":"acc-test-3"}

response

{
    "@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/notificationMessageTemplates('829c1a4e-fa0f-40fb-bdf2-31e8a7aefee8')/localizedNotificationMessages/$entity",
    "id": "829c1a4e-fa0f-40fb-bdf2-31e8a7aefee8_da-dk",
    "lastModifiedDateTime": "2025-08-20T12:01:22.2693759Z",
    "locale": "da-dk",
    "subject": "acc-test-3",
    "messageTemplate": "some\n\nmultiline\n\ntext",
    "isDefault": false
}


get by id

Request URL
https://graph.microsoft.com/beta/deviceManagement/notificationMessageTemplates/829c1a4e-fa0f-40fb-bdf2-31e8a7aefee8?$expand=localizedNotificationMessages
Request Method
GET

{
    "@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/notificationMessageTemplates(localizedNotificationMessages())/$entity",
    "id": "829c1a4e-fa0f-40fb-bdf2-31e8a7aefee8",
    "lastModifiedDateTime": "2025-08-20T12:01:22Z",
    "displayName": "acc-test",
    "description": null,
    "defaultLocale": "bg-bg",
    "brandingOptions": "includeCompanyLogo,includeCompanyName,includeContactInformation,includeCompanyPortalLink",
    "roleScopeTagIds": [
        "0"
    ],
    "localizedNotificationMessages@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/notificationMessageTemplates('829c1a4e-fa0f-40fb-bdf2-31e8a7aefee8')/localizedNotificationMessages",
    "localizedNotificationMessages": [
        {
            "id": "829c1a4e-fa0f-40fb-bdf2-31e8a7aefee8_bg-bg",
            "lastModifiedDateTime": "2025-08-20T12:05:43.1275117Z",
            "locale": "bg-bg",
            "subject": "acc-test",
            "messageTemplate": "some\n\nmultiline\n\ntext",
            "isDefault": true
        },
        {
            "id": "829c1a4e-fa0f-40fb-bdf2-31e8a7aefee8_da-dk",
            "lastModifiedDateTime": "2025-08-20T12:05:43.1275117Z",
            "locale": "da-dk",
            "subject": "acc-test-3",
            "messageTemplate": "some\n\nmultiline\n\ntext",
            "isDefault": false
        },
        {
            "id": "829c1a4e-fa0f-40fb-bdf2-31e8a7aefee8_cs-cz",
            "lastModifiedDateTime": "2025-08-20T12:05:43.1275117Z",
            "locale": "cs-cz",
            "subject": "acc-test-2",
            "messageTemplate": "some\n\nmulti\n\nline\n\ntext",
            "isDefault": false
        }
    ]
}