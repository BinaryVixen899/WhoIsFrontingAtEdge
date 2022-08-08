terraform {
  required_providers {
    fastly = {
        source = "fastly/fastly"
        version = ">= 2.2.1"
    }
  }
}

resource "fastly_service_compute" "api_kitsune_gay" {

    activate = true
    name = "api.kitsune.gay"
    comment = "api sewvice fow a cewtain focks"
    version_comment = "Transitioning to using Terraform"
    backend {
      address = "https://api.pluralkit.me"
    # DO NOT CHANGE THE NAME OR YOU WILL RECREATE THE BACKEND
      name = "Pluralkit"
      override_host = "api.pluralkit.me"
    }

    domain {
      name = "api.kitsune.gay"
    }


    package {
        filename = "../pkg/WhoisFrontingAtEdge.tar.gz"
    }
    
}

