package graphBetaRoleAssignment

type RoleAssignmentResourceModel struct {
	ODataType      string   `json:"@odata.type"`
	ID             string   `json:"id"`
	DisplayName    string   `json:"displayName"`
	Description    string   `json:"description"`
	ScopeMembers   []string `json:"scopeMembers"`
	ScopeType      string   `json:"scopeType"`
	ResourceScopes []string `json:"resourceScopes"`
}
