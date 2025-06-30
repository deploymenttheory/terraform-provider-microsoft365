dedicated

Request URL
https://graph.microsoft.com/beta/deviceManagement/virtualEndpoint/provisioningPolicies/3bb96ddf-5b64-4843-8335-87ef5105eff6/assign
Request Method
POST

{"assignments":[{"target":{"groupId":"35d09841-af73-43e6-a59f-024fef1b6b95"}},{"target":{"groupId":"410a28bd-9c9f-403f-b1b2-4a0bd04e98d9"}},{"target":{"groupId":"48fe6d79-f045-448a-bd74-716db27f0783"}},{"target":{"groupId":"106483b3-dcad-408b-95b7-9d3b9ad9c71b"}}]}

update

{"assignments":[{"id":"106483b3-dcad-408b-95b7-9d3b9ad9c71b","target":{"groupId":"106483b3-dcad-408b-95b7-9d3b9ad9c71b","servicePlanId":"","allotmentDisplayName":""}},{"id":"48fe6d79-f045-448a-bd74-716db27f0783","target":{"groupId":"48fe6d79-f045-448a-bd74-716db27f0783","servicePlanId":"","allotmentDisplayName":""}}]}

delete

{assignments: []}
assignments
: 
[]

frontline deadicated

create

{"assignments":[{"target":{"groupId":"35d09841-af73-43e6-a59f-024fef1b6b95","servicePlanId":"057efbfe-a95d-4263-acb0-12b4a31fed8d","allotmentLicensesCount":1,"allotmentDisplayName":"1"}}]}

Get with assignments

https://graph.microsoft.com/beta/deviceManagement/virtualEndpoint/provisioningPolicies?$expand=assignments

frontline shared

create

https://graph.microsoft.com/beta/deviceManagement/virtualEndpoint/provisioningPolicies/31092c57-0ef2-49dc-85fd-5170ea7c6cb4/assign

post

{"assignments":[{"target":{"groupId":"35d09841-af73-43e6-a59f-024fef1b6b95","servicePlanId":"057efbfe-a95d-4263-acb0-12b4a31fed8d","allotmentLicensesCount":1,"allotmentDisplayName":"test-1"}}]}



update

POST https://graph.microsoft.com/beta/deviceManagement/virtualEndpoint/provisioningPolicies/31092c57-0ef2-49dc-85fd-5170ea7c6cb4/assign

allotmentLicensesCount = Specify how many Cloud PCs should be provisioned. The number must be between 0 and 900 and it can't be more than the number of shared Cloud PC licenses available.

servicePlanId - front line service plan id

{"assignments":[{"id":"d49415e1-0191-49e3-83a4-7d9bf5f55e85","target":{"groupId":"35d09841-af73-43e6-a59f-024fef1b6b95","servicePlanId":"057efbfe-a95d-4263-acb0-12b4a31fed8d","allotmentLicensesCount":1,"allotmentDisplayName":"test-1-name-update"}}]}