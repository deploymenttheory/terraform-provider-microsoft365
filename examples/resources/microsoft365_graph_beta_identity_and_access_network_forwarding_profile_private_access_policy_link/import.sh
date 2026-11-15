#!/bin/bash
# User forwarding profile: {forwarding_profile_id}/{policy_link_id}
terraform import microsoft365_graph_beta_identity_and_access_network_forwarding_profile_private_access_policy_link.app 00000000-0000-0000-0000-000000000001/00000000-0000-0000-0000-000000000010

# Explicit profile type: {user|device}/{forwarding_profile_id}/{policy_link_id}
terraform import microsoft365_graph_beta_identity_and_access_network_forwarding_profile_private_access_policy_link.device device/00000000-0000-0000-0000-000000000003/00000000-0000-0000-0000-000000000011
