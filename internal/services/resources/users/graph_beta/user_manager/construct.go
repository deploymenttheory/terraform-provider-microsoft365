package graphBetaUsersUserManager

import (
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructManagerReference creates the reference body for assigning a manager.
// The request body requires an @odata.id pointing to the directoryObject (user or contact).
// REF: https://learn.microsoft.com/en-us/graph/api/agentuser-post-manager?view=graph-rest-beta
func constructManagerReference(managerId string) models.ReferenceUpdateable {
	requestBody := models.NewReferenceUpdate()
	odataId := "https://graph.microsoft.com/beta/users/" + managerId
	requestBody.SetOdataId(&odataId)
	return requestBody
}
