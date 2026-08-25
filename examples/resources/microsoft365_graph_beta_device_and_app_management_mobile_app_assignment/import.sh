# An assignment is addressed by the app that owns it as well as by its own id, so both are
# needed. The assignment id is reported by Graph as {groupId}_{intentIndex}_{n}.
#
# {mobile_app_id}:{assignment_id}
terraform import microsoft365_graph_beta_device_and_app_management_mobile_app_assignment.example 00000000-0000-0000-0000-000000000000:11111111-1111-1111-1111-111111111111_0_0
