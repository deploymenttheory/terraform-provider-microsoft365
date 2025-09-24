assignments - group assignment. 

// Dependencies
import (
	  "context"
	  msgraphsdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
	  graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
	  //other-imports
)

requestBody := graphmodels.NewWindowsAutopilotDeploymentProfileAssignment()
target := graphmodels.NewGroupAssignmentTarget()
groupId := "b15228f4-9d49-41ed-9b4f-0e7c721fd9c2"
target.SetGroupId(&groupId) 
requestBody.SetTarget(target)

// To initialize your graphClient, see https://learn.microsoft.com/en-us/graph/sdks/create-client?from=snippets&tabs=go
assignments, err := graphClient.DeviceManagement().WindowsAutopilotDeploymentProfiles().ByWindowsAutopilotDeploymentProfileId("windowsAutopilotDeploymentProfile-id").Assignments().Post(context.Background(), requestBody, nil)

delete assignment

// To initialize your graphClient, see https://learn.microsoft.com/en-us/graph/sdks/create-client?from=snippets&tabs=go
graphClient.DeviceManagement().WindowsAutopilotDeploymentProfiles().ByWindowsAutopilotDeploymentProfileId("windowsAutopilotDeploymentProfile-id").Assignments().ByWindowsAutopilotDeploymentProfileAssignmentId("windowsAutopilotDeploymentProfileAssignment-id").Delete(context.Background(), nil)

all device assignment

requestBody := graphmodels.NewWindowsAutopilotDeploymentProfileAssignment()
target := graphmodels.NewAllDevicesAssignmentTarget()
requestBody.SetTarget(target)

// To initialize your graphClient, see https://learn.microsoft.com/en-us/graph/sdks/create-client?from=snippets&tabs=go
assignments, err := graphClient.DeviceManagement().WindowsAutopilotDeploymentProfiles().ByWindowsAutopilotDeploymentProfileId("windowsAutopilotDeploymentProfile-id").Assignments().Post(context.Background(), requestBody, nil)


exclusion group assignment

// Dependencies
import (
	  "context"
	  msgraphsdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
	  graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
	  //other-imports
)

requestBody := graphmodels.NewWindowsAutopilotDeploymentProfileAssignment()
target := graphmodels.NewExclusionGroupAssignmentTarget()
groupId := "ea8e2fb8-e909-44e6-bae7-56757cf6f347"
target.SetGroupId(&groupId) 
requestBody.SetTarget(target)

// To initialize your graphClient, see https://learn.microsoft.com/en-us/graph/sdks/create-client?from=snippets&tabs=go
assignments, err := graphClient.DeviceManagement().WindowsAutopilotDeploymentProfiles().ByWindowsAutopilotDeploymentProfileId("windowsAutopilotDeploymentProfile-id").Assignments().Post(context.Background(), requestBody, nil)