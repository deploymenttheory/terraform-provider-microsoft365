package graphBetaNetworkMCPPolicy

// Local response model for the portal-backed endpoint, absent from the generated SDK.
type mcpPolicyResponse struct {
	ID                   *string `json:"id"`
	Name                 *string `json:"name"`
	Description          *string `json:"description"`
	Version              *string `json:"version"`
	LastModifiedDateTime *string `json:"lastModifiedDateTime"`
	Settings             *struct {
		DefaultAction *string `json:"defaultAction"`
	} `json:"settings"`
}
