terraform {
  required_providers {
    fastly = {
        source = "fastly/fastly"
        version = ">= 2.2.1"
    }
  }
}


variable "honeycomb_token" {
  description = "Honeycomb API token"
  type = string
  sensitive = true
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
      ssl_cert_hostname = "api.pluralkit.me"
      ssl_check_cert = true
      ssl_sni_hostname = "api.pluralkit.me"
      use_ssl = true
    }
    logging_honeycomb {
      name = "API.Kitsune.Gay"
      token = var.honeycomb_token
      dataset = "API.Kitsune.Gay"
    }

    domain {
      name = "api.kitsune.gay"
    }
    


    package {
        filename = "../pkg/WhoisFrontingAtEdge.tar.gz"
    }
    
}

