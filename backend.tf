// set terraform cloud organization and workspaceterraform
terraform {

  required_providers {
    microsoft365 = {
      source  = "deploymenttheory/microsoft365"
      version = "= 0.20.0-alpha"
    }
  }

  cloud {
    organization = "deploymenttheory"

    workspaces {
      # This is only relevant for CLI calls and is ignored by API calls via pipelines, therefore it can be safely left here.
      tags = ["microsoft_365"]
    }
  }
} 