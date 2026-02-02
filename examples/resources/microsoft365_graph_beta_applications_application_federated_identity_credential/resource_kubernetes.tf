resource "microsoft365_graph_beta_applications_application" "example" {
  display_name = "my-kubernetes-workload-app"
  description  = "Application for Kubernetes workload identity"
}

# Federated credential for Kubernetes workload identity
resource "microsoft365_graph_beta_applications_application_federated_identity_credential" "kubernetes" {
  application_id = microsoft365_graph_beta_applications_application.example.id
  name           = "aks-workload-identity"
  description    = "Azure Kubernetes Service workload identity"
  issuer         = "https://eastus.oic.prod-aks.azure.com/00000000-0000-0000-0000-000000000000/11111111-1111-1111-1111-111111111111/"
  subject        = "system:serviceaccount:default:my-service-account"
  audiences      = ["api://AzureADTokenExchange"]
}
