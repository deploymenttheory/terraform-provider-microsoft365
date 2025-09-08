Request URL
https://graph.microsoft.com/beta/deviceManagement/groupPolicyUploadedDefinitionFiles
Request Method
POST

req body

{
  "content":"PD94bWwgdmVyc2lvbj0iMS4wIiBlbmNvZGluZz0idXRmLTgiPz4KPHBvbGljeURlZmluaXRpb25zIHJldmlzaW9uPSI3LjIiIHNjaGVtYVZlcnNpb249IjEuMCI+CiAgPHBvbGljeU5hbWVzcGFjZXM+CiAgICA8dGFyZ2V0IHByZWZpeD0iZmlyZWZveCIgbmFtZXNwYWNlPSJNb3ppbGxhLlBvbGljaWVzLkZpcmVmb3giLz4KICAgIDx1c2luZyBwcmVmaXg9Ik1vemlsbGEiIG5hbWVzcGFjZT0iTW96aWxsYS5Qb2xpY2llcyIvPgogIDwvcG9saWN5TmFtZXNwYWNlcz4KICA8cmVzb3VyY2VzIG1pblJl+CiAgICAgICAgICAgICAgPHN0cmluZz50b29sYmFyPC9zdHJpbmc+CiAgICAgICAgPGRlY2ltYWwgdmFsdWU9IjEiLz4KICAgICAgPC9lbmFibGVkVmFsdWU+CiAgICAgIDxkaXNhYmxlZFZhbHVlPgogICAgICAgIDxkZWNpbWFsIHZhbHVlPSIwIi8+CiAgICAgIDwvZGlzYWJsZWRWYWx1ZT4KICAgIDwvcG9saWN5PgogIDwvcG9saWNpZXM+CjwvcG9saWN5RGVmaW5pdGlvbnM+Cg==",
  "fileName":"firefox.admx",
  "defaultLanguageCode":"",
  "groupPolicyUploadedLanguageFiles":[
    {"fileName":"firefox.adml",
    "languageCode":"en-US","content":"PD94bWwgdmVyc2lvbj0iMS4wIiBlbmNvZGluZz0idXRmLTgiPz4KPHBvbGljeURlZmluaXRpb25SZXNvdXJjZXMgcmV2aXNpb249IjcuMiIgc2NoZW1hVmVyc2lvbj0iMS4wIj4KICA8ZGlzcGxheU5hbWUvPgogIDxkZXNjcmlwdGlvbi8+CiAgICAgIDwvcHJlc2VudGF0aW9uPgogICAgPC9wcmVzZW50YXRpb25UYWJsZT4KICA8L3Jlc291cmNlcz4KPC9wb2xpY3lEZWZpbml0aW9uUmVzb3VyY2VzPgo="
    }
  ]
}

response

{
    "@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/groupPolicyUploadedDefinitionFiles/$entity",
    "displayName": null,
    "description": null,
    "languageCodes": [],
    "targetPrefix": null,
    "targetNamespace": "Mozilla.Policies.Firefox",
    "policyType": "admxIngested",
    "revision": null,
    "fileName": "firefox.admx",
    "id": "976886d5-f758-4cd9-9c80-b30b7355bfcb",
    "lastModifiedDateTime": "2025-09-07T08:32:57.5081686Z",
    "status": "uploadInProgress",
    "content": null,
    "uploadDateTime": "0001-01-01T00:00:00Z",
    "defaultLanguageCode": "en-US",
    "groupPolicyUploadedLanguageFiles": []
}

then performs a get to evaluate the status. takes a while

Request URL
https://graph.microsoft.com/beta/deviceManagement/groupPolicyUploadedDefinitionFiles
Request Method
GET

{
    "@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/groupPolicyUploadedDefinitionFiles",
    "value": [
        {
            "displayName": null,
            "description": null,
            "languageCodes": [],
            "targetPrefix": null,
            "targetNamespace": null,
            "policyType": "admxIngested",
            "revision": null,
            "fileName": "SecGuide.admx",
            "id": "b7b1f088-5cd5-4c9b-a19d-6e89409e0336",
            "lastModifiedDateTime": "2023-12-27T18:49:28.9386413Z",
            "status": "uploadFailed", <---- checking this
            "content": null,
            "uploadDateTime": "2023-12-27T18:49:28.9386413Z",
            "defaultLanguageCode": "en-US",
            "groupPolicyUploadedLanguageFiles": []
        },
        {
            "displayName": null,
            "description": null,
            "languageCodes": [],
            "targetPrefix": null,
            "targetNamespace": "Mozilla.Policies.Firefox",
            "policyType": "admxIngested",
            "revision": null,
            "fileName": "firefox.admx",
            "id": "976886d5-f758-4cd9-9c80-b30b7355bfcb",
            "lastModifiedDateTime": "2025-09-07T08:32:57.5081686Z",
            "status": "uploadInProgress", <---- checking this
            "content": null,
            "uploadDateTime": "0001-01-01T00:00:00Z",
            "defaultLanguageCode": "en-US",
            "groupPolicyUploadedLanguageFiles": []
        },
        {
            "displayName": null,
            "description": null,
            "languageCodes": [],
            "targetPrefix": "Mozilla47690255-69eb-47cf-bad8-91a56c11df9c",
            "targetNamespace": "Mozilla.Policies",
            "policyType": "admxIngested",
            "revision": "4.8",
            "fileName": "mozilla.admx",
            "id": "47690255-69eb-47cf-bad8-91a56c11df9c",
            "lastModifiedDateTime": "2025-09-07T08:42:58.653605Z",
            "status": "available", <---- success
            "content": null,
            "uploadDateTime": "2025-09-07T08:42:58.3508438Z",
            "defaultLanguageCode": "en-US",
            "groupPolicyUploadedLanguageFiles": []
        }
    ]
}


delete

Request URL
https://graph.microsoft.com/beta/deviceManagement/groupPolicyUploadedDefinitionFiles('976886d5-f758-4cd9-9c80-b30b7355bfcb')/remove
Request Method
POST