#!/bin/bash
# Import using {policy_type}/{policy_id}. policy_type is device (MDM user scope) or app (MAM user scope);
# policy_id is the application ID of the mobility management application (Microsoft Intune shown).
terraform import microsoft365_graph_beta_identity_and_access_mobility_management_policy.mdm device/0000000a-0000-0000-c000-000000000000
terraform import microsoft365_graph_beta_identity_and_access_mobility_management_policy.mam app/0000000a-0000-0000-c000-000000000000
