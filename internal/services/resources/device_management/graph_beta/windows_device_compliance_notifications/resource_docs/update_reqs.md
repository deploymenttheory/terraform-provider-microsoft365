Request URL
https://graph.microsoft.com/beta/deviceManagement/notificationMessageTemplates/829c1a4e-fa0f-40fb-bdf2-31e8a7aefee8
Request Method
PATCH

request body

{"brandingOptions":"includeCompanyLogo,includeCompanyName,includeContactInformation,includeCompanyPortalLink","displayName":"acc-test","roleScopeTagIds":["0"]}

response

{
    "@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/notificationMessageTemplates/$entity",
    "id": "829c1a4e-fa0f-40fb-bdf2-31e8a7aefee8",
    "lastModifiedDateTime": "2025-08-20T12:15:35.7600166Z",
    "displayName": "acc-test",
    "description": null,
    "defaultLocale": "bg-bg",
    "brandingOptions": "includeCompanyLogo,includeCompanyName,includeContactInformation,includeCompanyPortalLink",
    "roleScopeTagIds": [
        "0"
    ]
}

Request URL
https://graph.microsoft.com/beta/deviceManagement/notificationMessageTemplates/829c1a4e-fa0f-40fb-bdf2-31e8a7aefee8/localizedNotificationMessages/829c1a4e-fa0f-40fb-bdf2-31e8a7aefee8_bg-bg
Request Method
PATCH

request
{"subject":"acc-test","messageTemplate":"some\n\nmultiline\n\ntext","isDefault":true}

response

{
    "@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/notificationMessageTemplates('829c1a4e-fa0f-40fb-bdf2-31e8a7aefee8')/localizedNotificationMessages/$entity",
    "id": "829c1a4e-fa0f-40fb-bdf2-31e8a7aefee8_bg-bg",
    "lastModifiedDateTime": "2025-08-20T12:15:36.8414182Z",
    "locale": "bg-bg",
    "subject": "acc-test",
    "messageTemplate": "some\n\nmultiline\n\ntext",
    "isDefault": true
}

equest URL
https://graph.microsoft.com/beta/deviceManagement/notificationMessageTemplates/829c1a4e-fa0f-40fb-bdf2-31e8a7aefee8/localizedNotificationMessages/829c1a4e-fa0f-40fb-bdf2-31e8a7aefee8_da-dk
Request Method
PATCH

{"subject":"acc-test-3","messageTemplate":"some\n\nmultiline\n\ntext"}

response

{
    "@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/notificationMessageTemplates('829c1a4e-fa0f-40fb-bdf2-31e8a7aefee8')/localizedNotificationMessages/$entity",
    "id": "829c1a4e-fa0f-40fb-bdf2-31e8a7aefee8_da-dk",
    "lastModifiedDateTime": "2025-08-20T12:15:36.8901012Z",
    "locale": "da-dk",
    "subject": "acc-test-3",
    "messageTemplate": "some\n\nmultiline\n\ntext",
    "isDefault": false
}

Request URL
https://graph.microsoft.com/beta/deviceManagement/notificationMessageTemplates/829c1a4e-fa0f-40fb-bdf2-31e8a7aefee8/localizedNotificationMessages/829c1a4e-fa0f-40fb-bdf2-31e8a7aefee8_cs-cz
Request Method
PATCH

request

{"subject":"acc-test-2","messageTemplate":"some\n\nmulti\n\nline\n\ntext\n\nmore"}

response

{
    "@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/notificationMessageTemplates('829c1a4e-fa0f-40fb-bdf2-31e8a7aefee8')/localizedNotificationMessages/$entity",
    "id": "829c1a4e-fa0f-40fb-bdf2-31e8a7aefee8_cs-cz",
    "lastModifiedDateTime": "2025-08-20T12:15:36.924864Z",
    "locale": "cs-cz",
    "subject": "acc-test-2",
    "messageTemplate": "some\n\nmulti\n\nline\n\ntext\n\nmore",
    "isDefault": false
}


Request URL
https://graph.microsoft.com/beta/deviceManagement/notificationMessageTemplates/829c1a4e-fa0f-40fb-bdf2-31e8a7aefee8?$expand=localizedNotificationMessages
Request Method
GET

{
    "@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/notificationMessageTemplates(localizedNotificationMessages())/$entity",
    "id": "829c1a4e-fa0f-40fb-bdf2-31e8a7aefee8",
    "lastModifiedDateTime": "2025-08-20T12:15:36Z",
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
            "lastModifiedDateTime": "2025-08-20T12:15:36.6727967Z",
            "locale": "bg-bg",
            "subject": "acc-test",
            "messageTemplate": "some\n\nmultiline\n\ntext",
            "isDefault": true
        },
        {
            "id": "829c1a4e-fa0f-40fb-bdf2-31e8a7aefee8_da-dk",
            "lastModifiedDateTime": "2025-08-20T12:15:36.6727967Z",
            "locale": "da-dk",
            "subject": "acc-test-3",
            "messageTemplate": "some\n\nmultiline\n\ntext",
            "isDefault": false
        },
        {
            "id": "829c1a4e-fa0f-40fb-bdf2-31e8a7aefee8_cs-cz",
            "lastModifiedDateTime": "2025-08-20T12:15:36.6727967Z",
            "locale": "cs-cz",
            "subject": "acc-test-2",
            "messageTemplate": "some\n\nmulti\n\nline\n\ntext",
            "isDefault": false
        }
    ]
}