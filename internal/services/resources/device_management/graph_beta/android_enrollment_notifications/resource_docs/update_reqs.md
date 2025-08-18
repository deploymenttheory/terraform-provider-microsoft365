update base resource

PATCH https://graph.microsoft.com/beta/deviceManagement/deviceEnrollmentConfigurations/a8fe8ead-0dfd-4a19-af46-49c5ea89e70d_EnrollmentNotificationsConfiguration
Request
{
  "@odata.type": "#microsoft.graph.deviceEnrollmentNotificationConfiguration",
  "displayName": "test 2",
  "description": "test 2",
  "roleScopeTagIds": [
    "0"
  ],
  "id": "a8fe8ead-0dfd-4a19-af46-49c5ea89e70d_EnrollmentNotificationsConfiguration"
}

// Code snippets are only available for the latest major version. Current major version is $v0.*

// Dependencies
import (
	  "context"
	  msgraphsdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
	  graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
	  //other-imports
)

requestBody := graphmodels.NewDeviceEnrollmentConfiguration()
displayName := "test 2"
requestBody.SetDisplayName(&displayName) 
description := "test 2"
requestBody.SetDescription(&description) 
roleScopeTagIds := []string {
	"0",
}
requestBody.SetRoleScopeTagIds(roleScopeTagIds)
id := "a8fe8ead-0dfd-4a19-af46-49c5ea89e70d_EnrollmentNotificationsConfiguration"
requestBody.SetId(&id) 

// To initialize your graphClient, see https://learn.microsoft.com/en-us/graph/sdks/create-client?from=snippets&tabs=go
deviceEnrollmentConfigurations, err := graphClient.DeviceManagement().DeviceEnrollmentConfigurations().ByDeviceEnrollmentConfigurationId("deviceEnrollmentConfiguration-id").Patch(context.Background(), requestBody, nil)

all of the below for updating notifcation settings

PATCH https://graph.microsoft.com/beta/deviceManagement/notificationMessageTemplates/5a624256-39d8-4071-8f88-c4ac7756e284/localizedNotificationMessages/5a624256-39d8-4071-8f88-c4ac7756e284_en-us
Request
{
  "subject": "test 2",
  "messageTemplate": "test 2"
}
Response
{
  "@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/notificationMessageTemplates('5a624256-39d8-4071-8f88-c4ac7756e284')/localizedNotificationMessages/$entity",
  "id": "5a624256-39d8-4071-8f88-c4ac7756e284_en-us",
  "lastModifiedDateTime": "2025-08-18T10:28:17.5652332Z",
  "locale": "en-us",
  "subject": "test 2",
  "messageTemplate": "test 2",
  "isDefault": true
}

// Code snippets are only available for the latest major version. Current major version is $v0.*

// Dependencies
import (
	  "context"
	  msgraphsdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
	  graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
	  //other-imports
)

requestBody := graphmodels.NewLocalizedNotificationMessage()
subject := "test 2"
requestBody.SetSubject(&subject) 
messageTemplate := "test 2"
requestBody.SetMessageTemplate(&messageTemplate) 

// To initialize your graphClient, see https://learn.microsoft.com/en-us/graph/sdks/create-client?from=snippets&tabs=go
localizedNotificationMessages, err := graphClient.DeviceManagement().NotificationMessageTemplates().ByNotificationMessageTemplateId("notificationMessageTemplate-id").LocalizedNotificationMessages().ByLocalizedNotificationMessageId("localizedNotificationMessage-id").Patch(context.Background(), requestBody, nil)
PATCH https://graph.microsoft.com/beta/deviceManagement/notificationMessageTemplates/ef1f6abc-2c6b-4edf-b0f3-4472e495daa3/localizedNotificationMessages/ef1f6abc-2c6b-4edf-b0f3-4472e495daa3_en-us
Request
{
  "subject": "test 2",
  "messageTemplate": "test more stuff"
}
Response
{
  "@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/notificationMessageTemplates('ef1f6abc-2c6b-4edf-b0f3-4472e495daa3')/localizedNotificationMessages/$entity",
  "id": "ef1f6abc-2c6b-4edf-b0f3-4472e495daa3_en-us",
  "lastModifiedDateTime": "2025-08-18T10:28:17.6458163Z",
  "locale": "en-us",
  "subject": "test 2",
  "messageTemplate": "test more stuff",
  "isDefault": true
}

// Code snippets are only available for the latest major version. Current major version is $v0.*

// Dependencies
import (
	  "context"
	  msgraphsdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
	  graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
	  //other-imports
)

requestBody := graphmodels.NewLocalizedNotificationMessage()
subject := "test 2"
requestBody.SetSubject(&subject) 
messageTemplate := "test more stuff"
requestBody.SetMessageTemplate(&messageTemplate) 

// To initialize your graphClient, see https://learn.microsoft.com/en-us/graph/sdks/create-client?from=snippets&tabs=go
localizedNotificationMessages, err := graphClient.DeviceManagement().NotificationMessageTemplates().ByNotificationMessageTemplateId("notificationMessageTemplate-id").LocalizedNotificationMessages().ByLocalizedNotificationMessageId("localizedNotificationMessage-id").Patch(context.Background(), requestBody, nil)