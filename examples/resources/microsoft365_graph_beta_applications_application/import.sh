# The import ID format is: {object_id}[:prevent_duplicate_names=true|false][:hard_delete=true|false]
# {object_id} is the application's Object ID, a GUID such as 00000000-0000-0000-0000-000000000000.

# Simple import - defaults to prevent_duplicate_names=false and hard_delete=false
terraform import microsoft365_graph_beta_applications_application.example "00000000-0000-0000-0000-000000000000"

# Extended import - with hard_delete enabled for permanent deletion
terraform import microsoft365_graph_beta_applications_application.example "00000000-0000-0000-0000-000000000000:hard_delete=true"

# Extended import - with both prevent_duplicate_names and hard_delete enabled
terraform import microsoft365_graph_beta_applications_application.example "00000000-0000-0000-0000-000000000000:prevent_duplicate_names=true:hard_delete=true"
