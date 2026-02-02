# Simple import - defaults to hard_delete=false (soft delete with 30-day recovery)
terraform import microsoft365_graph_beta_applications_service_principal.example "00000000-0000-0000-0000-000000000000"

# Extended import - with hard_delete enabled for permanent deletion
terraform import microsoft365_graph_beta_applications_service_principal.example "00000000-0000-0000-0000-000000000000:hard_delete=true"
