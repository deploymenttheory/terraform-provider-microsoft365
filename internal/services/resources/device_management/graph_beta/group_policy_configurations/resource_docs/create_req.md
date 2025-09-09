Request URL
https://graph.microsoft.com/beta/deviceManagement/groupPolicyDefinitions('5cdf0ce8-f338-4bb3-b2c2-d00e5a255514')
Request Method
GET

{
    "@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/groupPolicyDefinitions/$entity",
    "classType": "machine",
    "displayName": "\u202FAllow users to contact Microsoft for feedback and support",
    "explainText": "\r\n      This setting specifies whether users in your organization can communicate directly with Microsoft through user experiences in the sync app. Letting users share their thoughts helps us improve OneDrive.\r\n\r\n      If you enable or do not configure this setting, users can use the experiences in the OneDrive sync app to contact Microsoft directly for feedback and support.\r\n\r\n      If you disable this setting, users will be unable to contact Microsoft for support, feedback, or suggestions within the sync app. Users will still have access to help content and self-help tools.\r\n    ",
    "categoryPath": "\\OneDrive",
    "supportedOn": "At least Windows Server 2008 R2 or Windows 7",
    "policyType": "admxIngested",
    "hasRelatedDefinitions": false,
    "groupPolicyCategoryId": "e3269686-b34a-4b03-970a-d1cc37951c27",
    "minDeviceCspVersion": "5.0",
    "minUserCspVersion": "",
    "version": "1.0",
    "id": "5cdf0ce8-f338-4bb3-b2c2-d00e5a255514",
    "lastModifiedDateTime": "2023-03-15T01:31:30.869884Z"
}

Request URL
https://graph.microsoft.com/beta/deviceManagement/groupPolicyDefinitions('5cdf0ce8-f338-4bb3-b2c2-d00e5a255514')/presentations
Request Method
GET

{
    "@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/groupPolicyDefinitions('5cdf0ce8-f338-4bb3-b2c2-d00e5a255514')/presentations",
    "value": [
        {
            "@odata.type": "#microsoft.graph.groupPolicyPresentationText",
            "label": "Allow users to contact Microsoft to: ",
            "id": "3fa73531-d6e9-4ed1-b6a4-803ccc9b2a88",
            "lastModifiedDateTime": "2023-03-15T01:31:33.6980945Z"
        },
        {
            "@odata.type": "#microsoft.graph.groupPolicyPresentationCheckBox",
            "label": "Send Feedback",
            "id": "95e90bbb-f30f-4d15-a3a6-00bcebfa02e9",
            "lastModifiedDateTime": "2023-03-15T01:31:33.6980945Z",
            "defaultChecked": true
        },
        {
            "@odata.type": "#microsoft.graph.groupPolicyPresentationCheckBox",
            "label": "Receive user satisfication surveys",
            "id": "ee6784ea-dca6-4a86-8058-c556add7f60f",
            "lastModifiedDateTime": "2023-03-15T01:31:33.6980945Z",
            "defaultChecked": true
        },
        {
            "@odata.type": "#microsoft.graph.groupPolicyPresentationCheckBox",
            "label": "Contact OneDrive Support​",
            "id": "1626672c-d7cc-400f-bd95-fdedcb90bf10",
            "lastModifiedDateTime": "2023-03-15T01:31:33.6980945Z",
            "defaultChecked": true
        }
    ]
}

Request URL
https://graph.microsoft.com/beta/deviceManagement/groupPolicyDefinitions('5cdf0ce8-f338-4bb3-b2c2-d00e5a255514')?$select=id,explainText,supportedOn
Request Method
GET

{
    "@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/groupPolicyDefinitions(id,explainText,supportedOn)/$entity",
    "id": "5cdf0ce8-f338-4bb3-b2c2-d00e5a255514",
    "explainText": "\r\n      This setting specifies whether users in your organization can communicate directly with Microsoft through user experiences in the sync app. Letting users share their thoughts helps us improve OneDrive.\r\n\r\n      If you enable or do not configure this setting, users can use the experiences in the OneDrive sync app to contact Microsoft directly for feedback and support.\r\n\r\n      If you disable this setting, users will be unable to contact Microsoft for support, feedback, or suggestions within the sync app. Users will still have access to help content and self-help tools.\r\n    ",
    "supportedOn": "At least Windows Server 2008 R2 or Windows 7"
}

Request URL
https://graph.microsoft.com/beta/deviceManagement/groupPolicyConfigurations('59d323cf-756c-46eb-82f8-c9a666c4e5ca')
Request Method
GET

{
    "@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/groupPolicyConfigurations/$entity",
    "createdDateTime": "2025-09-09T06:25:04.3778917Z",
    "displayName": "unit-test-group-policy-presentation-value-text",
    "description": "unit test for group policy presentation value text",
    "roleScopeTagIds": [
        "0",
        "1",
        "2"
    ],
    "policyConfigurationIngestionType": "unknown",
    "id": "59d323cf-756c-46eb-82f8-c9a666c4e5ca",
    "lastModifiedDateTime": "2025-09-09T06:25:04.3778917Z"
}

Request URL
https://graph.microsoft.com/beta/deviceManagement/groupPolicyConfigurations('59d323cf-756c-46eb-82f8-c9a666c4e5ca')/definitionValues?$expand=definition($select=id,classType,displayName,policyType,hasRelatedDefinitions,version,minUserCspVersion,minDeviceCspVersion)
Request Method
GET

{
    "@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/groupPolicyConfigurations('59d323cf-756c-46eb-82f8-c9a666c4e5ca')/definitionValues(definition(id,classType,displayName,policyType,hasRelatedDefinitions,version,minUserCspVersion,minDeviceCspVersion))",
    "value": []
}

Request URL
https://graph.microsoft.com/beta/deviceManagement/groupPolicyConfigurations('59d323cf-756c-46eb-82f8-c9a666c4e5ca')/updateDefinitionValues
Request Method
POST

payload

{"added":[{"enabled":true,"presentationValues":[{"@odata.type":"#microsoft.graph.groupPolicyPresentationValueBoolean","value":true,"presentation@odata.bind":"https://graph.microsoft.com/beta/deviceManagement/groupPolicyDefinitions('5cdf0ce8-f338-4bb3-b2c2-d00e5a255514')/presentations('95e90bbb-f30f-4d15-a3a6-00bcebfa02e9')"},{"@odata.type":"#microsoft.graph.groupPolicyPresentationValueBoolean","value":true,"presentation@odata.bind":"https://graph.microsoft.com/beta/deviceManagement/groupPolicyDefinitions('5cdf0ce8-f338-4bb3-b2c2-d00e5a255514')/presentations('ee6784ea-dca6-4a86-8058-c556add7f60f')"},{"@odata.type":"#microsoft.graph.groupPolicyPresentationValueBoolean","value":true,"presentation@odata.bind":"https://graph.microsoft.com/beta/deviceManagement/groupPolicyDefinitions('5cdf0ce8-f338-4bb3-b2c2-d00e5a255514')/presentations('1626672c-d7cc-400f-bd95-fdedcb90bf10')"}],"definition@odata.bind":"https://graph.microsoft.com/beta/deviceManagement/groupPolicyDefinitions('5cdf0ce8-f338-4bb3-b2c2-d00e5a255514')"}],"updated":[],"deletedIds":[]}

response

Request URL
https://graph.microsoft.com/beta/deviceManagement/groupPolicyConfigurations('59d323cf-756c-46eb-82f8-c9a666c4e5ca')
Request Method
GET

{
    "@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/groupPolicyConfigurations/$entity",
    "createdDateTime": "2025-09-09T06:25:04.3778917Z",
    "displayName": "unit-test-group-policy-presentation-value-text",
    "description": "unit test for group policy presentation value text",
    "roleScopeTagIds": [
        "0",
        "1",
        "2"
    ],
    "policyConfigurationIngestionType": "unknown",
    "id": "59d323cf-756c-46eb-82f8-c9a666c4e5ca",
    "lastModifiedDateTime": "2025-09-09T06:27:11.5793678Z"
}

Request URL
https://graph.microsoft.com/beta/deviceManagement/groupPolicyConfigurations('59d323cf-756c-46eb-82f8-c9a666c4e5ca')/definitionValues?$expand=definition($select=id,classType,displayName,policyType,hasRelatedDefinitions,version,minUserCspVersion,minDeviceCspVersion)
Request Method
GET

{
    "@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/groupPolicyConfigurations('59d323cf-756c-46eb-82f8-c9a666c4e5ca')/definitionValues(definition(id,classType,displayName,policyType,hasRelatedDefinitions,version,minUserCspVersion,minDeviceCspVersion))",
    "value": [
        {
            "createdDateTime": "2025-09-09T06:27:11.5481198Z",
            "enabled": true,
            "configurationType": "policy",
            "id": "40f22a20-87a0-4744-ba77-8839561432c5",
            "lastModifiedDateTime": "2025-09-09T06:27:11.5793678Z",
            "definition": {
                "id": "5cdf0ce8-f338-4bb3-b2c2-d00e5a255514",
                "classType": "machine",
                "displayName": "\u202FAllow users to contact Microsoft for feedback and support",
                "policyType": "admxIngested",
                "hasRelatedDefinitions": false,
                "version": "1.0",
                "minUserCspVersion": "",
                "minDeviceCspVersion": "5.0"
            }
        }
    ]
}

Request URL
https://graph.microsoft.com/beta/deviceManagement/groupPolicyConfigurations('59d323cf-756c-46eb-82f8-c9a666c4e5ca')/assignments
Request Method
GET

{
    "@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/groupPolicyConfigurations('59d323cf-756c-46eb-82f8-c9a666c4e5ca')/assignments",
    "value": [
        {
            "id": "59d323cf-756c-46eb-82f8-c9a666c4e5ca_b15228f4-9d49-41ed-9b4f-0e7c721fd9c2",
            "lastModifiedDateTime": "2025-09-09T06:32:48.1046102Z",
            "target": {
                "@odata.type": "#microsoft.graph.exclusionGroupAssignmentTarget",
                "deviceAndAppManagementAssignmentFilterId": null,
                "deviceAndAppManagementAssignmentFilterType": "none",
                "groupId": "b15228f4-9d49-41ed-9b4f-0e7c721fd9c2"
            }
        },
        {
            "id": "59d323cf-756c-46eb-82f8-c9a666c4e5ca_35d09841-af73-43e6-a59f-024fef1b6b95",
            "lastModifiedDateTime": "2025-09-09T06:32:48.1046102Z",
            "target": {
                "@odata.type": "#microsoft.graph.exclusionGroupAssignmentTarget",
                "deviceAndAppManagementAssignmentFilterId": null,
                "deviceAndAppManagementAssignmentFilterType": "none",
                "groupId": "35d09841-af73-43e6-a59f-024fef1b6b95"
            }
        },
        {
            "id": "59d323cf-756c-46eb-82f8-c9a666c4e5ca_acacacac-9df4-4c7d-9d50-4ef0226f57a9",
            "lastModifiedDateTime": "2025-09-09T06:32:48.1046102Z",
            "target": {
                "@odata.type": "#microsoft.graph.allLicensedUsersAssignmentTarget",
                "deviceAndAppManagementAssignmentFilterId": null,
                "deviceAndAppManagementAssignmentFilterType": "none"
            }
        }
    ]
}


Request URL
https://graph.microsoft.com/beta/deviceManagement/groupPolicyConfigurations('59d323cf-756c-46eb-82f8-c9a666c4e5ca')/updateDefinitionValues
Request Method
POST

{"added":[],"updated":[{"id":"40f22a20-87a0-4744-ba77-8839561432c5","enabled":true,"presentationValues":[{"id":"21151fb5-6597-40de-8f4c-4c30902352c7","@odata.type":"#microsoft.graph.groupPolicyPresentationValueBoolean","value":true,"presentation@odata.bind":"https://graph.microsoft.com/beta/deviceManagement/groupPolicyDefinitions('5cdf0ce8-f338-4bb3-b2c2-d00e5a255514')/presentations('95e90bbb-f30f-4d15-a3a6-00bcebfa02e9')"},{"id":"fd299d43-e2c1-4606-b598-c20511bf83ec","@odata.type":"#microsoft.graph.groupPolicyPresentationValueBoolean","value":true,"presentation@odata.bind":"https://graph.microsoft.com/beta/deviceManagement/groupPolicyDefinitions('5cdf0ce8-f338-4bb3-b2c2-d00e5a255514')/presentations('ee6784ea-dca6-4a86-8058-c556add7f60f')"},{"id":"07780829-cacb-4d48-b7f6-9ddd04c8485b","@odata.type":"#microsoft.graph.groupPolicyPresentationValueBoolean","value":false,"presentation@odata.bind":"https://graph.microsoft.com/beta/deviceManagement/groupPolicyDefinitions('5cdf0ce8-f338-4bb3-b2c2-d00e5a255514')/presentations('1626672c-d7cc-400f-bd95-fdedcb90bf10')"}],"definition@odata.bind":"https://graph.microsoft.com/beta/deviceManagement/groupPolicyDefinitions('5cdf0ce8-f338-4bb3-b2c2-d00e5a255514')"}],"deletedIds":[]}




delete

Request URL
https://graph.microsoft.com/beta/deviceManagement/groupPolicyConfigurations/59d323cf-756c-46eb-82f8-c9a666c4e5ca
Request Method
DELETE



define settings values

enabled

{"added":[{"enabled":true,"presentationValues":[{"@odata.type":"#microsoft.graph.groupPolicyPresentationValueBoolean","value":true,"presentation@odata.bind":"https://graph.microsoft.com/beta/deviceManagement/groupPolicyDefinitions('5cdf0ce8-f338-4bb3-b2c2-d00e5a255514')/presentations('95e90bbb-f30f-4d15-a3a6-00bcebfa02e9')"},{"@odata.type":"#microsoft.graph.groupPolicyPresentationValueBoolean","value":true,"presentation@odata.bind":"https://graph.microsoft.com/beta/deviceManagement/groupPolicyDefinitions('5cdf0ce8-f338-4bb3-b2c2-d00e5a255514')/presentations('ee6784ea-dca6-4a86-8058-c556add7f60f')"},{"@odata.type":"#microsoft.graph.groupPolicyPresentationValueBoolean","value":true,"presentation@odata.bind":"https://graph.microsoft.com/beta/deviceManagement/groupPolicyDefinitions('5cdf0ce8-f338-4bb3-b2c2-d00e5a255514')/presentations('1626672c-d7cc-400f-bd95-fdedcb90bf10')"}],"definition@odata.bind":"https://graph.microsoft.com/beta/deviceManagement/groupPolicyDefinitions('5cdf0ce8-f338-4bb3-b2c2-d00e5a255514')"}],"updated":[],"deletedIds":[]}

disabled looks like this

{"added":[],"updated":[{"id":"92b07061-3e3b-47f9-934b-2f125227fa09","enabled":false,"presentationValues":[],"definition@odata.bind":"https://graph.microsoft.com/beta/deviceManagement/groupPolicyDefinitions('5cdf0ce8-f338-4bb3-b2c2-d00e5a255514')"}],"deletedIds":[]}

and not configured is like this

Request URL
https://graph.microsoft.com/beta/deviceManagement/groupPolicyConfigurations('319bd10c-c79d-43ee-8228-ab52e0c08eb7')/updateDefinitionValues
Request Method
POST

{"added":[],"updated":[],"deletedIds":["92b07061-3e3b-47f9-934b-2f125227fa09"]}



{
    "added":
    [
        {
            "enabled":true,"presentationValues":
            [
                {
                    "@odata.type":"#microsoft.graph.groupPolicyPresentationValueBoolean",
                    "value":true,
                    "presentation@odata.bind":"https://graph.microsoft.com/beta/deviceManagement/groupPolicyDefinitions('5cdf0ce8-f338-4bb3-b2c2-d00e5a255514')/presentations('95e90bbb-f30f-4d15-a3a6-00bcebfa02e9')"
                    },
                    {
                    "@odata.type":"#microsoft.graph.groupPolicyPresentationValueBoolean",
                    "value":true,
                    "presentation@odata.bind":"https://graph.microsoft.com/beta/deviceManagement/groupPolicyDefinitions('5cdf0ce8-f338-4bb3-b2c2-d00e5a255514')/presentations('ee6784ea-dca6-4a86-8058-c556add7f60f')"
                    },
                    {
                        "@odata.type":"#microsoft.graph.groupPolicyPresentationValueBoolean",
                        "value":true,
                        "presentation@odata.bind":"https://graph.microsoft.com/beta/deviceManagement/groupPolicyDefinitions('5cdf0ce8-f338-4bb3-b2c2-d00e5a255514')/presentations('1626672c-d7cc-400f-bd95-fdedcb90bf10')"
                        }
                    ],
                    "definition@odata.bind":"https://graph.microsoft.com/beta/deviceManagement/groupPolicyDefinitions('5cdf0ce8-f338-4bb3-b2c2-d00e5a255514')"
                    }
                ],
            "updated":[],
            "deletedIds":[]
}

Request URL
https://graph.microsoft.com/beta/deviceManagement/groupPolicyConfigurations('07b30785-2ab1-40e2-b060-92e6c117e5d3')/definitionValues?$expand=definition($select=id,classType,displayName,policyType,hasRelatedDefinitions,version,minUserCspVersion,minDeviceCspVersion)
Request Method
GET


Request URL
https://graph.microsoft.com/beta/deviceManagement/groupPolicyDefinitions('5cdf0ce8-f338-4bb3-b2c2-d00e5a255514')
Request Method
GET

{
    "@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/groupPolicyDefinitions/$entity",
    "classType": "machine",
    "displayName": "\u202FAllow users to contact Microsoft for feedback and support",
    "explainText": "\r\n      This setting specifies whether users in your organization can communicate directly with Microsoft through user experiences in the sync app. Letting users share their thoughts helps us improve OneDrive.\r\n\r\n      If you enable or do not configure this setting, users can use the experiences in the OneDrive sync app to contact Microsoft directly for feedback and support.\r\n\r\n      If you disable this setting, users will be unable to contact Microsoft for support, feedback, or suggestions within the sync app. Users will still have access to help content and self-help tools.\r\n    ",
    "categoryPath": "\\OneDrive",
    "supportedOn": "At least Windows Server 2008 R2 or Windows 7",
    "policyType": "admxIngested",
    "hasRelatedDefinitions": false,
    "groupPolicyCategoryId": "e3269686-b34a-4b03-970a-d1cc37951c27",
    "minDeviceCspVersion": "5.0",
    "minUserCspVersion": "",
    "version": "1.0",
    "id": "5cdf0ce8-f338-4bb3-b2c2-d00e5a255514",
    "lastModifiedDateTime": "2023-03-15T01:31:30.869884Z"
}