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
    version_comment = "Allegedly Working Fastly.toml"
    backend {
      address = "https://api.pluralkit.me"
    # DO NOT CHANGE THE NAME OR YOU WILL RECREATE THE BACKEND
      name = "Pluralkit"
      override_host = "api.pluralkit.me"
      port = 443
      use_ssl = true
    }

    domain {
      name = "api.kitsune.gay"
    }


    package {
        filename = "../pkg/WhoisFrontingAtEdge.tar.gz"
    }
    
}

