# Find Windows Update product by KB number
# This example shows how to search for a product using a KB article number
# KB numbers are typically known from Microsoft security bulletins or support articles

data "microsoft365_graph_beta_windows_updates_product" "by_kb_number" {
  search_type  = "kb_number"
  search_value = "5029332" # Example: KB5029332
}

output "product_by_kb" {
  description = "Product information retrieved by KB number"
  value = length(data.microsoft365_graph_beta_windows_updates_product.by_kb_number.products) > 0 ? {
    id           = data.microsoft365_graph_beta_windows_updates_product.by_kb_number.products[0].id
    name         = data.microsoft365_graph_beta_windows_updates_product.by_kb_number.products[0].name
    group_name   = data.microsoft365_graph_beta_windows_updates_product.by_kb_number.products[0].group_name
    revisions    = length(data.microsoft365_graph_beta_windows_updates_product.by_kb_number.products[0].revisions)
    known_issues = length(data.microsoft365_graph_beta_windows_updates_product.by_kb_number.products[0].known_issues)
  } : null
}

output "kb_articles" {
  description = "Knowledge Base articles for all revisions"
  value = length(data.microsoft365_graph_beta_windows_updates_product.by_kb_number.products) > 0 ? [
    for revision in data.microsoft365_graph_beta_windows_updates_product.by_kb_number.products[0].revisions : {
      revision_id = revision.id
      kb_id       = try(revision.knowledge_base_article.id, null)
      kb_url      = try(revision.knowledge_base_article.url, null)
    }
  ] : []
}

output "known_issues" {
  description = "All known issues for the product"
  value = length(data.microsoft365_graph_beta_windows_updates_product.by_kb_number.products) > 0 ? [
    for issue in data.microsoft365_graph_beta_windows_updates_product.by_kb_number.products[0].known_issues : {
      id                 = issue.id
      title              = issue.title
      status             = issue.status
      description        = issue.description
      web_view_url       = issue.web_view_url
      start_date_time    = issue.start_date_time
      resolved_date_time = issue.resolved_date_time
    }
  ] : []
}
