Request URL
https://graph.microsoft.com/beta/deviceAppManagement/mobileAppConfigurations
Request Method
POST

{
  "@odata.type":"#microsoft.graph.androidManagedStoreAppConfiguration",
  "displayName":"test",
  "description":"",
  "profileApplicability":"androidDeviceOwner",
  "targetedMobileApps":["9711516a-f6f8-4953-ad1f-45920ef34dda"],
  "roleScopeTagIds":["0"],
  "packageId":"app:com.microsoft.office.officehubrow","payloadJson":"eyJraW5kIjoiYW5kcm9pZGVudGVycHJpc2UjbWFuYWdlZENvbmZpZ3VyYXRpb24iLCJwcm9kdWN0SWQiOiJhcHA6Y29tLm1pY3Jvc29mdC5vZmZpY2Uub2ZmaWNlaHVicm93IiwibWFuYWdlZFByb3BlcnR5IjpbeyJrZXkiOiJjb20ubWljcm9zb2Z0LmludHVuZS5tYW0uQWxsb3dlZEFjY291bnRVUE5zIiwidmFsdWVTdHJpbmciOiJ0aGluZyJ9LHsia2V5IjoiY29tLm1pY3Jvc29mdC5vZmZpY2UuTm90ZXNDcmVhdGlvbkVuYWJsZWQiLCJ2YWx1ZUJvb2wiOnRydWV9LHsia2V5IjoiY29tLm1pY3Jvc29mdC5vZmZpY2Uub2ZmaWNlbW9iaWxlLlRlYW1zQXBwcy5Jc0FsbG93ZWQiLCJ2YWx1ZUJvb2wiOnRydWV9LHsia2V5IjoiY29tLm1pY3Jvc29mdC5vZmZpY2Uub2ZmaWNlbW9iaWxlLkJpbmdDaGF0RW50ZXJwcmlzZS5Jc0FsbG93ZWQiLCJ2YWx1ZUJvb2wiOnRydWV9XX0=","permissionActions":[
    {"permission":"android.permission-group.NEARBY_DEVICES","action":"prompt"},{"permission":"android.permission.NEARBY_WIFI_DEVICES","action":"prompt"},{"permission":"android.permission.BLUETOOTH_CONNECT","action":"prompt"},{"permission":"android.permission.READ_MEDIA_AUDIO","action":"prompt"},{"permission":"android.permission.READ_MEDIA_IMAGES","action":"prompt"},{"permission":"android.permission.READ_MEDIA_VIDEO","action":"prompt"},{"permission":"android.permission.POST_NOTIFICATIONS","action":"prompt"},{"permission":"android.permission.WRITE_EXTERNAL_STORAGE","action":"prompt"},{"permission":"android.permission.READ_EXTERNAL_STORAGE","action":"prompt"},{"permission":"android.permission.RECEIVE_MMS","action":"prompt"},{"permission":"android.permission.RECEIVE_WAP_PUSH","action":"prompt"},{"permission":"android.permission.READ_SMS","action":"prompt"},{"permission":"android.permission.RECEIVE_SMS","action":"prompt"},{"permission":"android.permission.SEND_SMS","action":"prompt"},{"permission":"android.permission.BODY_SENSORS_BACKGROUND","action":"prompt"},{"permission":"android.permission.BODY_SENSORS","action":"prompt"},{"permission":"android.permission.PROCESS_OUTGOING_CALLS","action":"prompt"},{"permission":"android.permission.USE_SIP","action":"prompt"},{"permission":"android.permission.ADD_VOICEMAIL","action":"prompt"},{"permission":"android.permission.WRITE_CALL_LOG","action":"prompt"},{"permission":"android.permission.READ_CALL_LOG","action":"prompt"},{"permission":"android.permission.CALL_PHONE","action":"prompt"},{"permission":"android.permission.READ_PHONE_STATE","action":"prompt"},{"permission":"android.permission.RECORD_AUDIO","action":"prompt"},{"permission":"android.permission.ACCESS_BACKGROUND_LOCATION","action":"prompt"},{"permission":"android.permission.ACCESS_COARSE_LOCATION","action":"prompt"},{"permission":"android.permission.ACCESS_FINE_LOCATION","action":"prompt"},{"permission":"android.permission.GET_ACCOUNTS","action":"autoGrant"},{"permission":"android.permission.WRITE_CONTACTS","action":"prompt"},{"permission":"android.permission.READ_CONTACTS","action":"prompt"},{"permission":"android.permission.CAMERA","action":"autoDeny"},{"permission":"android.permission.WRITE_CALENDAR","action":"prompt"},{"permission":"android.permission.READ_CALENDAR","action":"autoGrant"}],
    "connectedAppsEnabled":true,
    "id":"00000000-0000-0000-0000-000000000000"}


    resp

    {
    "@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceAppManagement/mobileAppConfigurations/$entity",
    "@odata.type": "#microsoft.graph.androidManagedStoreAppConfiguration",
    "id": "8f74ee2e-dcba-40e1-86fd-86aa63f2822f",
    "targetedMobileApps": [
        "9711516a-f6f8-4953-ad1f-45920ef34dda"
    ],
    "roleScopeTagIds": [
        "0"
    ],
    "createdDateTime": "2025-11-14T08:55:03.0638909Z",
    "description": "",
    "lastModifiedDateTime": "2025-11-14T08:55:03.0638909Z",
    "displayName": "test",
    "version": 1,
    "packageId": "app:com.microsoft.office.officehubrow",
    "payloadJson": "eyJraW5kIjoiYW5kcm9pZGVudGVycHJpc2UjbWFuYWdlZENvbmZpZ3VyYXRpb24iLCJwcm9kdWN0SWQiOiJhcHA6Y29tLm1pY3Jvc29mdC5vZmZpY2Uub2ZmaWNlaHVicm93IiwibWFuYWdlZFByb3BlcnR5IjpbeyJrZXkiOiJjb20ubWljcm9zb2Z0LmludHVuZS5tYW0uQWxsb3dlZEFjY291bnRVUE5zIiwidmFsdWVTdHJpbmciOiJ0aGluZyJ9LHsia2V5IjoiY29tLm1pY3Jvc29mdC5vZmZpY2UuTm90ZXNDcmVhdGlvbkVuYWJsZWQiLCJ2YWx1ZUJvb2wiOnRydWV9LHsia2V5IjoiY29tLm1pY3Jvc29mdC5vZmZpY2Uub2ZmaWNlbW9iaWxlLlRlYW1zQXBwcy5Jc0FsbG93ZWQiLCJ2YWx1ZUJvb2wiOnRydWV9LHsia2V5IjoiY29tLm1pY3Jvc29mdC5vZmZpY2Uub2ZmaWNlbW9iaWxlLkJpbmdDaGF0RW50ZXJwcmlzZS5Jc0FsbG93ZWQiLCJ2YWx1ZUJvb2wiOnRydWV9XX0=",
    "appSupportsOemConfig": false,
    "profileApplicability": "androidDeviceOwner",
    "connectedAppsEnabled": true,
    "permissionActions": [
        {
            "permission": "android.permission-group.NEARBY_DEVICES",
            "action": "prompt"
        },
        {
            "permission": "android.permission.NEARBY_WIFI_DEVICES",
            "action": "prompt"
        },
        {
            "permission": "android.permission.BLUETOOTH_CONNECT",
            "action": "prompt"
        },
        {
            "permission": "android.permission.READ_MEDIA_AUDIO",
            "action": "prompt"
        },
        {
            "permission": "android.permission.READ_MEDIA_IMAGES",
            "action": "prompt"
        },
        {
            "permission": "android.permission.READ_MEDIA_VIDEO",
            "action": "prompt"
        },
        {
            "permission": "android.permission.POST_NOTIFICATIONS",
            "action": "prompt"
        },
        {
            "permission": "android.permission.WRITE_EXTERNAL_STORAGE",
            "action": "prompt"
        },
        {
            "permission": "android.permission.READ_EXTERNAL_STORAGE",
            "action": "prompt"
        },
        {
            "permission": "android.permission.RECEIVE_MMS",
            "action": "prompt"
        },
        {
            "permission": "android.permission.RECEIVE_WAP_PUSH",
            "action": "prompt"
        },
        {
            "permission": "android.permission.READ_SMS",
            "action": "prompt"
        },
        {
            "permission": "android.permission.RECEIVE_SMS",
            "action": "prompt"
        },
        {
            "permission": "android.permission.SEND_SMS",
            "action": "prompt"
        },
        {
            "permission": "android.permission.BODY_SENSORS_BACKGROUND",
            "action": "prompt"
        },
        {
            "permission": "android.permission.BODY_SENSORS",
            "action": "prompt"
        },
        {
            "permission": "android.permission.PROCESS_OUTGOING_CALLS",
            "action": "prompt"
        },
        {
            "permission": "android.permission.USE_SIP",
            "action": "prompt"
        },
        {
            "permission": "android.permission.ADD_VOICEMAIL",
            "action": "prompt"
        },
        {
            "permission": "android.permission.WRITE_CALL_LOG",
            "action": "prompt"
        },
        {
            "permission": "android.permission.READ_CALL_LOG",
            "action": "prompt"
        },
        {
            "permission": "android.permission.CALL_PHONE",
            "action": "prompt"
        },
        {
            "permission": "android.permission.READ_PHONE_STATE",
            "action": "prompt"
        },
        {
            "permission": "android.permission.RECORD_AUDIO",
            "action": "prompt"
        },
        {
            "permission": "android.permission.ACCESS_BACKGROUND_LOCATION",
            "action": "prompt"
        },
        {
            "permission": "android.permission.ACCESS_COARSE_LOCATION",
            "action": "prompt"
        },
        {
            "permission": "android.permission.ACCESS_FINE_LOCATION",
            "action": "prompt"
        },
        {
            "permission": "android.permission.GET_ACCOUNTS",
            "action": "autoGrant"
        },
        {
            "permission": "android.permission.WRITE_CONTACTS",
            "action": "prompt"
        },
        {
            "permission": "android.permission.READ_CONTACTS",
            "action": "prompt"
        },
        {
            "permission": "android.permission.CAMERA",
            "action": "autoDeny"
        },
        {
            "permission": "android.permission.WRITE_CALENDAR",
            "action": "prompt"
        },
        {
            "permission": "android.permission.READ_CALENDAR",
            "action": "autoGrant"
        }
    ]
}
