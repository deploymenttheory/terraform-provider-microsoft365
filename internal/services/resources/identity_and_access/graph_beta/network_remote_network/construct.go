package graphBetaIdentityAndAccessNetworkRemoteNetwork

// profileRefs builds forwarding profile references from their IDs.
func profileRefs(ids []string) []map[string]string {
	refs := make([]map[string]string, 0, len(ids))
	for _, id := range ids {
		refs = append(refs, map[string]string{"id": id})
	}
	return refs
}

// constructResource builds the remote network create payload.
func constructResource(name, region string, profileIDs []string) (map[string]any, error) {
	if len(profileIDs) == 0 {
		return nil, errNoForwardingProfiles
	}
	return map[string]any{"name": name, "region": region, "forwardingProfiles": profileRefs(profileIDs)}, nil
}

// constructForwardingProfileDelta builds the complete desired profile membership.
func constructForwardingProfileDelta(profileIDs []string) (map[string]any, error) {
	if len(profileIDs) == 0 {
		return nil, errNoForwardingProfiles
	}
	return map[string]any{"@context": "#$delta", "value": profileRefs(profileIDs)}, nil
}

// constructNameUpdate builds the remote network name update payload.
func constructNameUpdate(name string) map[string]any { return map[string]any{"name": name} }
