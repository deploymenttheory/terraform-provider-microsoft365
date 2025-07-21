// All users with daily schedule and include filter with 2 exclusion groups

// Dependencies
import (
	  "context"
	  msgraphsdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
	  graphdevicemanagement "github.com/microsoftgraph/msgraph-beta-sdk-go/devicemanagement"
	  graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
	  //other-imports
)

requestBody := graphdevicemanagement.NewAssignPostRequestBody()


deviceHealthScriptAssignment := graphmodels.NewDeviceHealthScriptAssignment()
id := "11a47806-4491-4596-9690-5b4a4c1656ca:acacacac-9df4-4c7d-9d50-4ef0226f57a9"
deviceHealthScriptAssignment.SetId(&id) 
target := graphmodels.NewAllLicensedUsersAssignmentTarget()
deviceAndAppManagementAssignmentFilterId := "dc20e791-31c9-47d1-8e74-ae7995cabb09"
target.SetDeviceAndAppManagementAssignmentFilterId(&deviceAndAppManagementAssignmentFilterId) 
deviceAndAppManagementAssignmentFilterType := graphmodels.INCLUDE_DEVICEANDAPPMANAGEMENTASSIGNMENTFILTERTYPE 
target.SetDeviceAndAppManagementAssignmentFilterType(&deviceAndAppManagementAssignmentFilterType) 
deviceHealthScriptAssignment.SetTarget(target)
runRemediationScript := true
deviceHealthScriptAssignment.SetRunRemediationScript(&runRemediationScript) 
runSchedule := graphmodels.NewDeviceHealthScriptDailySchedule()
interval := int32(1)
runSchedule.SetInterval(&interval) 
time := 20:0:0
runSchedule.SetTime(&time) 
useUtc := true
runSchedule.SetUseUtc(&useUtc) 
deviceHealthScriptAssignment.SetRunSchedule(runSchedule)
deviceHealthScriptAssignment1 := graphmodels.NewDeviceHealthScriptAssignment()
target := graphmodels.NewExclusionGroupAssignmentTarget()
groupId := "ea8e2fb8-e909-44e6-bae7-56757cf6f347"
target.SetGroupId(&groupId) 
deviceHealthScriptAssignment1.SetTarget(target)
runRemediationScript := true
deviceHealthScriptAssignment1.SetRunRemediationScript(&runRemediationScript) 
runSchedule := null
deviceHealthScriptAssignment1.SetRunSchedule(&runSchedule) 
deviceHealthScriptAssignment2 := graphmodels.NewDeviceHealthScriptAssignment()
target := graphmodels.NewExclusionGroupAssignmentTarget()
groupId := "b15228f4-9d49-41ed-9b4f-0e7c721fd9c2"
target.SetGroupId(&groupId) 
deviceHealthScriptAssignment2.SetTarget(target)
runRemediationScript := true
deviceHealthScriptAssignment2.SetRunRemediationScript(&runRemediationScript) 
runSchedule := null
deviceHealthScriptAssignment2.SetRunSchedule(&runSchedule) 

deviceHealthScriptAssignments := []graphmodels.DeviceHealthScriptAssignmentable {
	deviceHealthScriptAssignment,
	deviceHealthScriptAssignment1,
	deviceHealthScriptAssignment2,
}
requestBody.SetDeviceHealthScriptAssignments(deviceHealthScriptAssignments)

// To initialize your graphClient, see https://learn.microsoft.com/en-us/graph/sdks/create-client?from=snippets&tabs=go
graphClient.DeviceManagement().DeviceHealthScripts().ByDeviceHealthScriptId("deviceHealthScript-id").Assign().Post(context.Background(), requestBody, nil)



// All users with once schedule and exclude filter with 2 exclusion groups


// Dependencies
import (
	  "context"
	  msgraphsdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
	  graphdevicemanagement "github.com/microsoftgraph/msgraph-beta-sdk-go/devicemanagement"
	  graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
	  //other-imports
)

requestBody := graphdevicemanagement.NewAssignPostRequestBody()


deviceHealthScriptAssignment := graphmodels.NewDeviceHealthScriptAssignment()
id := "11a47806-4491-4596-9690-5b4a4c1656ca:acacacac-9df4-4c7d-9d50-4ef0226f57a9"
deviceHealthScriptAssignment.SetId(&id) 
target := graphmodels.NewAllLicensedUsersAssignmentTarget()
deviceAndAppManagementAssignmentFilterId := "dc20e791-31c9-47d1-8e74-ae7995cabb09"
target.SetDeviceAndAppManagementAssignmentFilterId(&deviceAndAppManagementAssignmentFilterId) 
deviceAndAppManagementAssignmentFilterType := graphmodels.EXCLUDE_DEVICEANDAPPMANAGEMENTASSIGNMENTFILTERTYPE 
target.SetDeviceAndAppManagementAssignmentFilterType(&deviceAndAppManagementAssignmentFilterType) 
deviceHealthScriptAssignment.SetTarget(target)
runRemediationScript := true
deviceHealthScriptAssignment.SetRunRemediationScript(&runRemediationScript) 
runSchedule := graphmodels.NewDeviceHealthScriptRunOnceSchedule()
interval := int32(1)
runSchedule.SetInterval(&interval) 
date := 2025-07-19
runSchedule.SetDate(&date) 
time := 20:0:0
runSchedule.SetTime(&time) 
useUtc := true
runSchedule.SetUseUtc(&useUtc) 
deviceHealthScriptAssignment.SetRunSchedule(runSchedule)
deviceHealthScriptAssignment1 := graphmodels.NewDeviceHealthScriptAssignment()
target := graphmodels.NewExclusionGroupAssignmentTarget()
groupId := "ea8e2fb8-e909-44e6-bae7-56757cf6f347"
target.SetGroupId(&groupId) 
deviceHealthScriptAssignment1.SetTarget(target)
runRemediationScript := true
deviceHealthScriptAssignment1.SetRunRemediationScript(&runRemediationScript) 
runSchedule := null
deviceHealthScriptAssignment1.SetRunSchedule(&runSchedule) 
deviceHealthScriptAssignment2 := graphmodels.NewDeviceHealthScriptAssignment()
target := graphmodels.NewExclusionGroupAssignmentTarget()
groupId := "b15228f4-9d49-41ed-9b4f-0e7c721fd9c2"
target.SetGroupId(&groupId) 
deviceHealthScriptAssignment2.SetTarget(target)
runRemediationScript := true
deviceHealthScriptAssignment2.SetRunRemediationScript(&runRemediationScript) 
runSchedule := null
deviceHealthScriptAssignment2.SetRunSchedule(&runSchedule) 

deviceHealthScriptAssignments := []graphmodels.DeviceHealthScriptAssignmentable {
	deviceHealthScriptAssignment,
	deviceHealthScriptAssignment1,
	deviceHealthScriptAssignment2,
}
requestBody.SetDeviceHealthScriptAssignments(deviceHealthScriptAssignments)

// To initialize your graphClient, see https://learn.microsoft.com/en-us/graph/sdks/create-client?from=snippets&tabs=go
graphClient.DeviceManagement().DeviceHealthScripts().ByDeviceHealthScriptId("deviceHealthScript-id").Assign().Post(context.Background(), requestBody, nil)



// All users with hourly schedule and exclude filter with 2 exclusion groups


// Dependencies
import (
	  "context"
	  msgraphsdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
	  graphdevicemanagement "github.com/microsoftgraph/msgraph-beta-sdk-go/devicemanagement"
	  graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
	  //other-imports
)

requestBody := graphdevicemanagement.NewAssignPostRequestBody()


deviceHealthScriptAssignment := graphmodels.NewDeviceHealthScriptAssignment()
id := "11a47806-4491-4596-9690-5b4a4c1656ca:acacacac-9df4-4c7d-9d50-4ef0226f57a9"
deviceHealthScriptAssignment.SetId(&id) 
target := graphmodels.NewAllLicensedUsersAssignmentTarget()
deviceAndAppManagementAssignmentFilterId := "dc20e791-31c9-47d1-8e74-ae7995cabb09"
target.SetDeviceAndAppManagementAssignmentFilterId(&deviceAndAppManagementAssignmentFilterId) 
deviceAndAppManagementAssignmentFilterType := graphmodels.EXCLUDE_DEVICEANDAPPMANAGEMENTASSIGNMENTFILTERTYPE 
target.SetDeviceAndAppManagementAssignmentFilterType(&deviceAndAppManagementAssignmentFilterType) 
deviceHealthScriptAssignment.SetTarget(target)
runRemediationScript := true
deviceHealthScriptAssignment.SetRunRemediationScript(&runRemediationScript) 
runSchedule := graphmodels.NewDeviceHealthScriptHourlySchedule()
interval := int32(1)
runSchedule.SetInterval(&interval) 
deviceHealthScriptAssignment.SetRunSchedule(runSchedule)
deviceHealthScriptAssignment1 := graphmodels.NewDeviceHealthScriptAssignment()
target := graphmodels.NewExclusionGroupAssignmentTarget()
groupId := "ea8e2fb8-e909-44e6-bae7-56757cf6f347"
target.SetGroupId(&groupId) 
deviceHealthScriptAssignment1.SetTarget(target)
runRemediationScript := true
deviceHealthScriptAssignment1.SetRunRemediationScript(&runRemediationScript) 
runSchedule := null
deviceHealthScriptAssignment1.SetRunSchedule(&runSchedule) 
deviceHealthScriptAssignment2 := graphmodels.NewDeviceHealthScriptAssignment()
target := graphmodels.NewExclusionGroupAssignmentTarget()
groupId := "b15228f4-9d49-41ed-9b4f-0e7c721fd9c2"
target.SetGroupId(&groupId) 
deviceHealthScriptAssignment2.SetTarget(target)
runRemediationScript := true
deviceHealthScriptAssignment2.SetRunRemediationScript(&runRemediationScript) 
runSchedule := null
deviceHealthScriptAssignment2.SetRunSchedule(&runSchedule) 

deviceHealthScriptAssignments := []graphmodels.DeviceHealthScriptAssignmentable {
	deviceHealthScriptAssignment,
	deviceHealthScriptAssignment1,
	deviceHealthScriptAssignment2,
}
requestBody.SetDeviceHealthScriptAssignments(deviceHealthScriptAssignments)

// To initialize your graphClient, see https://learn.microsoft.com/en-us/graph/sdks/create-client?from=snippets&tabs=go
graphClient.DeviceManagement().DeviceHealthScripts().ByDeviceHealthScriptId("deviceHealthScript-id").Assign().Post(context.Background(), requestBody, nil)

// All devices with daily schedule and include filter with 2 exclusion groups

// Dependencies
import (
	  "context"
	  msgraphsdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
	  graphdevicemanagement "github.com/microsoftgraph/msgraph-beta-sdk-go/devicemanagement"
	  graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
	  //other-imports
)

requestBody := graphdevicemanagement.NewAssignPostRequestBody()


deviceHealthScriptAssignment := graphmodels.NewDeviceHealthScriptAssignment()
target := graphmodels.NewAllDevicesAssignmentTarget()
deviceAndAppManagementAssignmentFilterId := "2d7956fb-e5b5-4fa3-90b2-5bee9bee7883"
target.SetDeviceAndAppManagementAssignmentFilterId(&deviceAndAppManagementAssignmentFilterId) 
deviceAndAppManagementAssignmentFilterType := graphmodels.INCLUDE_DEVICEANDAPPMANAGEMENTASSIGNMENTFILTERTYPE 
target.SetDeviceAndAppManagementAssignmentFilterType(&deviceAndAppManagementAssignmentFilterType) 
deviceHealthScriptAssignment.SetTarget(target)
runRemediationScript := true
deviceHealthScriptAssignment.SetRunRemediationScript(&runRemediationScript) 
runSchedule := graphmodels.NewDeviceHealthScriptDailySchedule()
interval := int32(1)
runSchedule.SetInterval(&interval) 
time := 1:0:0
runSchedule.SetTime(&time) 
useUtc := false
runSchedule.SetUseUtc(&useUtc) 
deviceHealthScriptAssignment.SetRunSchedule(runSchedule)
deviceHealthScriptAssignment1 := graphmodels.NewDeviceHealthScriptAssignment()
target := graphmodels.NewExclusionGroupAssignmentTarget()
groupId := "ea8e2fb8-e909-44e6-bae7-56757cf6f347"
target.SetGroupId(&groupId) 
deviceHealthScriptAssignment1.SetTarget(target)
runRemediationScript := true
deviceHealthScriptAssignment1.SetRunRemediationScript(&runRemediationScript) 
runSchedule := null
deviceHealthScriptAssignment1.SetRunSchedule(&runSchedule) 
deviceHealthScriptAssignment2 := graphmodels.NewDeviceHealthScriptAssignment()
target := graphmodels.NewExclusionGroupAssignmentTarget()
groupId := "b15228f4-9d49-41ed-9b4f-0e7c721fd9c2"
target.SetGroupId(&groupId) 
deviceHealthScriptAssignment2.SetTarget(target)
runRemediationScript := true
deviceHealthScriptAssignment2.SetRunRemediationScript(&runRemediationScript) 
runSchedule := null
deviceHealthScriptAssignment2.SetRunSchedule(&runSchedule) 

deviceHealthScriptAssignments := []graphmodels.DeviceHealthScriptAssignmentable {
	deviceHealthScriptAssignment,
	deviceHealthScriptAssignment1,
	deviceHealthScriptAssignment2,
}
requestBody.SetDeviceHealthScriptAssignments(deviceHealthScriptAssignments)

// To initialize your graphClient, see https://learn.microsoft.com/en-us/graph/sdks/create-client?from=snippets&tabs=go
graphClient.DeviceManagement().DeviceHealthScripts().ByDeviceHealthScriptId("deviceHealthScript-id").Assign().Post(context.Background(), requestBody, nil)


// All devices with once schedule and exclude filter with 2 exclusion groups

// Dependencies
import (
	  "context"
	  msgraphsdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
	  graphdevicemanagement "github.com/microsoftgraph/msgraph-beta-sdk-go/devicemanagement"
	  graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
	  //other-imports
)

requestBody := graphdevicemanagement.NewAssignPostRequestBody()


deviceHealthScriptAssignment := graphmodels.NewDeviceHealthScriptAssignment()
target := graphmodels.NewAllDevicesAssignmentTarget()
deviceAndAppManagementAssignmentFilterId := "2d7956fb-e5b5-4fa3-90b2-5bee9bee7883"
target.SetDeviceAndAppManagementAssignmentFilterId(&deviceAndAppManagementAssignmentFilterId) 
deviceAndAppManagementAssignmentFilterType := graphmodels.INCLUDE_DEVICEANDAPPMANAGEMENTASSIGNMENTFILTERTYPE 
target.SetDeviceAndAppManagementAssignmentFilterType(&deviceAndAppManagementAssignmentFilterType) 
deviceHealthScriptAssignment.SetTarget(target)
runRemediationScript := true
deviceHealthScriptAssignment.SetRunRemediationScript(&runRemediationScript) 
runSchedule := graphmodels.NewDeviceHealthScriptHourlySchedule()
interval := int32(1)
runSchedule.SetInterval(&interval) 
deviceHealthScriptAssignment.SetRunSchedule(runSchedule)
deviceHealthScriptAssignment1 := graphmodels.NewDeviceHealthScriptAssignment()
target := graphmodels.NewExclusionGroupAssignmentTarget()
groupId := "ea8e2fb8-e909-44e6-bae7-56757cf6f347"
target.SetGroupId(&groupId) 
deviceHealthScriptAssignment1.SetTarget(target)
runRemediationScript := true
deviceHealthScriptAssignment1.SetRunRemediationScript(&runRemediationScript) 
runSchedule := null
deviceHealthScriptAssignment1.SetRunSchedule(&runSchedule) 
deviceHealthScriptAssignment2 := graphmodels.NewDeviceHealthScriptAssignment()
target := graphmodels.NewExclusionGroupAssignmentTarget()
groupId := "b15228f4-9d49-41ed-9b4f-0e7c721fd9c2"
target.SetGroupId(&groupId) 
deviceHealthScriptAssignment2.SetTarget(target)
runRemediationScript := true
deviceHealthScriptAssignment2.SetRunRemediationScript(&runRemediationScript) 
runSchedule := null
deviceHealthScriptAssignment2.SetRunSchedule(&runSchedule) 

deviceHealthScriptAssignments := []graphmodels.DeviceHealthScriptAssignmentable {
	deviceHealthScriptAssignment,
	deviceHealthScriptAssignment1,
	deviceHealthScriptAssignment2,
}
requestBody.SetDeviceHealthScriptAssignments(deviceHealthScriptAssignments)

// To initialize your graphClient, see https://learn.microsoft.com/en-us/graph/sdks/create-client?from=snippets&tabs=go
graphClient.DeviceManagement().DeviceHealthScripts().ByDeviceHealthScriptId("deviceHealthScript-id").Assign().Post(context.Background(), requestBody, nil)

// All devices with hourly schedule and exclude filter with 2 exclusion groups

// Dependencies
import (
	  "context"
	  msgraphsdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
	  graphdevicemanagement "github.com/microsoftgraph/msgraph-beta-sdk-go/devicemanagement"
	  graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
	  //other-imports
)

requestBody := graphdevicemanagement.NewAssignPostRequestBody()


deviceHealthScriptAssignment := graphmodels.NewDeviceHealthScriptAssignment()
target := graphmodels.NewAllDevicesAssignmentTarget()
deviceAndAppManagementAssignmentFilterId := "2d7956fb-e5b5-4fa3-90b2-5bee9bee7883"
target.SetDeviceAndAppManagementAssignmentFilterId(&deviceAndAppManagementAssignmentFilterId) 
deviceAndAppManagementAssignmentFilterType := graphmodels.EXCLUDE_DEVICEANDAPPMANAGEMENTASSIGNMENTFILTERTYPE 
target.SetDeviceAndAppManagementAssignmentFilterType(&deviceAndAppManagementAssignmentFilterType) 
deviceHealthScriptAssignment.SetTarget(target)
runRemediationScript := true
deviceHealthScriptAssignment.SetRunRemediationScript(&runRemediationScript) 
runSchedule := graphmodels.NewDeviceHealthScriptRunOnceSchedule()
interval := int32(1)
runSchedule.SetInterval(&interval) 
date := 2025-07-19
runSchedule.SetDate(&date) 
time := 15:0:0
runSchedule.SetTime(&time) 
useUtc := true
runSchedule.SetUseUtc(&useUtc) 
deviceHealthScriptAssignment.SetRunSchedule(runSchedule)
deviceHealthScriptAssignment1 := graphmodels.NewDeviceHealthScriptAssignment()
target := graphmodels.NewExclusionGroupAssignmentTarget()
groupId := "ea8e2fb8-e909-44e6-bae7-56757cf6f347"
target.SetGroupId(&groupId) 
deviceHealthScriptAssignment1.SetTarget(target)
runRemediationScript := true
deviceHealthScriptAssignment1.SetRunRemediationScript(&runRemediationScript) 
runSchedule := null
deviceHealthScriptAssignment1.SetRunSchedule(&runSchedule) 
deviceHealthScriptAssignment2 := graphmodels.NewDeviceHealthScriptAssignment()
target := graphmodels.NewExclusionGroupAssignmentTarget()
groupId := "b15228f4-9d49-41ed-9b4f-0e7c721fd9c2"
target.SetGroupId(&groupId) 
deviceHealthScriptAssignment2.SetTarget(target)
runRemediationScript := true
deviceHealthScriptAssignment2.SetRunRemediationScript(&runRemediationScript) 
runSchedule := null
deviceHealthScriptAssignment2.SetRunSchedule(&runSchedule) 

deviceHealthScriptAssignments := []graphmodels.DeviceHealthScriptAssignmentable {
	deviceHealthScriptAssignment,
	deviceHealthScriptAssignment1,
	deviceHealthScriptAssignment2,
}
requestBody.SetDeviceHealthScriptAssignments(deviceHealthScriptAssignments)

// To initialize your graphClient, see https://learn.microsoft.com/en-us/graph/sdks/create-client?from=snippets&tabs=go
graphClient.DeviceManagement().DeviceHealthScripts().ByDeviceHealthScriptId("deviceHealthScript-id").Assign().Post(context.Background(), requestBody, nil)

// Include group assignment with 3 groups.  1 with daily schedule and exclude filter, 1 with hourly schedule and exclude filter, 1 with once schedule and exclude filter, and 2 exclusion groups

// Dependencies
import (
	  "context"
	  msgraphsdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
	  graphdevicemanagement "github.com/microsoftgraph/msgraph-beta-sdk-go/devicemanagement"
	  graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
	  //other-imports
)

requestBody := graphdevicemanagement.NewAssignPostRequestBody()


deviceHealthScriptAssignment := graphmodels.NewDeviceHealthScriptAssignment()
target := graphmodels.NewGroupAssignmentTarget()
groupId := "35d09841-af73-43e6-a59f-024fef1b6b95"
target.SetGroupId(&groupId) 
deviceAndAppManagementAssignmentFilterId := "80f8c0a5-f3ec-4936-bcbc-420dc0ca3665"
target.SetDeviceAndAppManagementAssignmentFilterId(&deviceAndAppManagementAssignmentFilterId) 
deviceAndAppManagementAssignmentFilterType := graphmodels.INCLUDE_DEVICEANDAPPMANAGEMENTASSIGNMENTFILTERTYPE 
target.SetDeviceAndAppManagementAssignmentFilterType(&deviceAndAppManagementAssignmentFilterType) 
deviceHealthScriptAssignment.SetTarget(target)
runRemediationScript := true
deviceHealthScriptAssignment.SetRunRemediationScript(&runRemediationScript) 
runSchedule := graphmodels.NewDeviceHealthScriptDailySchedule()
interval := int32(1)
runSchedule.SetInterval(&interval) 
time := 1:0:0
runSchedule.SetTime(&time) 
useUtc := false
runSchedule.SetUseUtc(&useUtc) 
deviceHealthScriptAssignment.SetRunSchedule(runSchedule)
deviceHealthScriptAssignment1 := graphmodels.NewDeviceHealthScriptAssignment()
target := graphmodels.NewGroupAssignmentTarget()
groupId := "410a28bd-9c9f-403f-b1b2-4a0bd04e98d9"
target.SetGroupId(&groupId) 
deviceAndAppManagementAssignmentFilterId := "8333759e-46b1-4ae4-8353-cb9fde130760"
target.SetDeviceAndAppManagementAssignmentFilterId(&deviceAndAppManagementAssignmentFilterId) 
deviceAndAppManagementAssignmentFilterType := graphmodels.EXCLUDE_DEVICEANDAPPMANAGEMENTASSIGNMENTFILTERTYPE 
target.SetDeviceAndAppManagementAssignmentFilterType(&deviceAndAppManagementAssignmentFilterType) 
deviceHealthScriptAssignment1.SetTarget(target)
runRemediationScript := true
deviceHealthScriptAssignment1.SetRunRemediationScript(&runRemediationScript) 
runSchedule := graphmodels.NewDeviceHealthScriptHourlySchedule()
interval := int32(1)
runSchedule.SetInterval(&interval) 
deviceHealthScriptAssignment1.SetRunSchedule(runSchedule)
deviceHealthScriptAssignment2 := graphmodels.NewDeviceHealthScriptAssignment()
target := graphmodels.NewGroupAssignmentTarget()
groupId := "48fe6d79-f045-448a-bd74-716db27f0783"
target.SetGroupId(&groupId) 
deviceAndAppManagementAssignmentFilterId := "99b2823d-a05c-4316-9a82-3efa40ff482d"
target.SetDeviceAndAppManagementAssignmentFilterId(&deviceAndAppManagementAssignmentFilterId) 
deviceAndAppManagementAssignmentFilterType := graphmodels.EXCLUDE_DEVICEANDAPPMANAGEMENTASSIGNMENTFILTERTYPE 
target.SetDeviceAndAppManagementAssignmentFilterType(&deviceAndAppManagementAssignmentFilterType) 
deviceHealthScriptAssignment2.SetTarget(target)
runRemediationScript := true
deviceHealthScriptAssignment2.SetRunRemediationScript(&runRemediationScript) 
runSchedule := graphmodels.NewDeviceHealthScriptRunOnceSchedule()
interval := int32(1)
runSchedule.SetInterval(&interval) 
date := 2025-07-19
runSchedule.SetDate(&date) 
time := 17:45:10
runSchedule.SetTime(&time) 
useUtc := false
runSchedule.SetUseUtc(&useUtc) 
deviceHealthScriptAssignment2.SetRunSchedule(runSchedule)
deviceHealthScriptAssignment3 := graphmodels.NewDeviceHealthScriptAssignment()
target := graphmodels.NewExclusionGroupAssignmentTarget()
groupId := "ea8e2fb8-e909-44e6-bae7-56757cf6f347"
target.SetGroupId(&groupId) 
deviceHealthScriptAssignment3.SetTarget(target)
runRemediationScript := true
deviceHealthScriptAssignment3.SetRunRemediationScript(&runRemediationScript) 
runSchedule := null
deviceHealthScriptAssignment3.SetRunSchedule(&runSchedule) 
deviceHealthScriptAssignment4 := graphmodels.NewDeviceHealthScriptAssignment()
target := graphmodels.NewExclusionGroupAssignmentTarget()
groupId := "b15228f4-9d49-41ed-9b4f-0e7c721fd9c2"
target.SetGroupId(&groupId) 
deviceHealthScriptAssignment4.SetTarget(target)
runRemediationScript := true
deviceHealthScriptAssignment4.SetRunRemediationScript(&runRemediationScript) 
runSchedule := null
deviceHealthScriptAssignment4.SetRunSchedule(&runSchedule) 

deviceHealthScriptAssignments := []graphmodels.DeviceHealthScriptAssignmentable {
	deviceHealthScriptAssignment,
	deviceHealthScriptAssignment1,
	deviceHealthScriptAssignment2,
	deviceHealthScriptAssignment3,
	deviceHealthScriptAssignment4,
}
requestBody.SetDeviceHealthScriptAssignments(deviceHealthScriptAssignments)

// To initialize your graphClient, see https://learn.microsoft.com/en-us/graph/sdks/create-client?from=snippets&tabs=go
graphClient.DeviceManagement().DeviceHealthScripts().ByDeviceHealthScriptId("deviceHealthScript-id").Assign().Post(context.Background(), requestBody, nil)