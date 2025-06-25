terraform {
  required_version = "1.12.2"
  required_providers {
    keycloak = {
      source  = "keycloak/keycloak"
      version = "5.2.0"
    }
  }
}

provider "keycloak" {
  client_id = "admin-cli"
  url       = "http://auth:18080"
}
